package foundation

import (
	"context"
	"fmt"
	"strings"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// GetTenantMetaTags is a helper to create tags that are used to identify which tenant a resource belongs to
// and are saved in a resources `Meta` attribute
func (c *FoundationImpl) GetTenantMetaTags(ctx context.Context) ([]domain.FHIRCodingInput, error) {
	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	return c.CreateTenantMetaTags(ctx, identifiers.OrganizationID, identifiers.FacilityID)
}

// GetTenantMetaTags is a helper to create tags that are used to identify which tenant a resource belongs to
// and are saved in a resources `Meta` attribute
func (c *FoundationImpl) CreateTenantMetaTags(ctx context.Context, organisationID, facilityID string) ([]domain.FHIRCodingInput, error) {
	organisation, err := c.FHIR.GetFHIROrganization(ctx, organisationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tenant organisation: %w", err)
	}

	facility, err := c.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tenant organisation: %w", err)
	}

	userSelected := false
	organisationTagVersion, facilityTagVersion := "1.0", "1.0"

	tags := []domain.FHIRCodingInput{
		{
			System:       &common.OrganisationSystem,
			Version:      &organisationTagVersion,
			Code:         scalarutils.Code(*organisation.Resource.ID),
			Display:      *organisation.Resource.Name,
			UserSelected: &userSelected,
		},
		{
			System:       &common.FacilitySystem,
			Version:      &facilityTagVersion,
			Code:         scalarutils.Code(*facility.Resource.ID),
			Display:      *facility.Resource.Name,
			UserSelected: &userSelected,
		},
	}

	return tags, nil
}

// maxPatientReferenceMatches caps how many patients a free-text search resolves to, keeping the
// downstream reference IN-list within FHIR/URL length limits.
const maxPatientReferenceMatches = 100

// patientReferenceSearch pairs a Patient search parameter with the value to match against it.
type patientReferenceSearch struct {
	param string
	value string
}

// SearchPatientReferences resolves a free-text term into a list of "Patient/<id>" references used
// as a comma-separated IN-list filter on resources that reference a patient. The term is matched
// against a deliberately narrow set of Patient fields — name, identifier and phone — and the
// results are unioned (OR semantics). This is intentionally narrower than a `_content` full-text
// search, which would also match address, notes, etc. Name and identifier use the raw term; phone
// is only searched when the term is a valid, normalisable phone number. It returns an empty slice
// when nothing matches, and logs when the match set hits the cap.
func (c *FoundationImpl) SearchPatientReferences(ctx context.Context, term string, tenant dto.TenantIdentifiers) ([]string, error) {
	searches := []patientReferenceSearch{
		{param: "name", value: term},
		{param: "identifier", value: term},
	}

	// A phone search only makes sense for a normalisable phone number, and it must be normalised so
	// it matches how numbers are stored.
	if normalized, err := converterandformatter.NormalizeMSISDN(term); err == nil {
		searches = append(searches, patientReferenceSearch{param: "phone", value: *normalized})
	}

	seen := make(map[string]bool)
	references := make([]string, 0)

	for _, search := range searches {
		params := map[string]any{
			search.param: search.value,
			"_count":     maxPatientReferenceMatches,
		}

		patients, err := c.FHIR.SearchFHIRResource(ctx, "", "Patient", params, tenant, dto.Pagination{})
		if err != nil {
			return nil, err
		}

		for _, resource := range patients.Resources {
			patient, err := domain.ConvertMapToFHIRResource[domain.FHIRPatient](resource)
			if err != nil {
				return nil, err
			}

			if patient.ID == nil || seen[*patient.ID] {
				continue
			}

			seen[*patient.ID] = true
			references = append(references, fmt.Sprintf("Patient/%s", *patient.ID))

			if len(references) >= maxPatientReferenceMatches {
				c.Warn("patient search reached the match cap; some patients may be omitted from the filter",
					"term", term, "cap", maxPatientReferenceMatches)

				return references, nil
			}
		}
	}

	return references, nil
}

// CheckPatientExistenceUsingPhoneNumber checks whether a patient with the phone number they're trying to register with exists
func (c *FoundationImpl) CheckPatientExistenceUsingPhoneNumber(ctx context.Context, patientInput domain.SimplePatientRegistrationInput) (bool, error) {
	exists := false

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	for _, phone := range patientInput.PhoneNumbers {
		phoneNumber := phone.Msisdn

		search, err := converterandformatter.NormalizeMSISDN(phoneNumber)
		if err != nil {
			return false, fmt.Errorf("can't normalize contact: %w", err)
		}

		patient, err := c.FHIR.SearchFHIRPatient(ctx, *search, *identifiers, dto.Pagination{})
		if err != nil {
			return false, fmt.Errorf("unable to find patient by phonenumber: %s", phoneNumber)
		}

		if len(patient.Edges) >= 1 {
			exists = true
			break
		}
	}

	return exists, nil
}

// SimplePatientRegistrationInputToPatientInput transforms a patient input into
func (c *FoundationImpl) SimplePatientRegistrationInputToPatientInput(ctx context.Context, input domain.SimplePatientRegistrationInput, organizationID string) (*domain.FHIRPatientInput, error) {
	contacts, err := c.ContactsToContactPointInput(ctx, input.PhoneNumbers, input.Emails)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't register patient with invalid contacts: %w", err)
	}

	ids, err := helpers.IDToIdentifier(input.IdentificationDocuments, organizationID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't register patient with invalid identifiers: %w", err)
	}

	// fullPatientInput is to be filled up by processing the simple patient input
	gender := domain.PatientGenderEnum(input.Gender)
	patientInput := domain.FHIRPatientInput{
		BirthDate: input.BirthDate,
		Gender:    &gender,
		Active:    &input.Active,
	}
	patientInput.Identifier = ids
	patientInput.Telecom = contacts
	patientInput.Name = helpers.NameToHumanName(input.Names)
	// patientInput.Photo = photos
	patientInput.Address = helpers.PhysicalPostalAddressesToFHIRAddresses(
		input.PhysicalAddresses, input.PostalAddresses)
	patientInput.Communication = helpers.LanguagesToCommunicationInputs(input.Languages)

	return &patientInput, nil
}

// ContactsToContactPointInput translates phone and email contacts to
// FHIR contact points
func (c *FoundationImpl) ContactsToContactPointInput(_ context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
	if phones == nil && emails == nil {
		return nil, nil
	}

	output := []*domain.FHIRContactPointInput{}
	rank := int64(1)
	phoneSystem := domain.ContactPointSystemEnumPhone
	use := domain.ContactPointUseEnumHome

	for _, phone := range phones {
		normalized, err := converterandformatter.NormalizeMSISDN(phone.Msisdn)
		if err != nil {
			utils.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to normalize phonenumber")
		}

		phoneContact := &domain.FHIRContactPointInput{
			System: &phoneSystem,
			Use:    &use,
			Rank:   &rank,
			Period: common.DefaultPeriodInput(),
			Value:  normalized,
		}
		output = append(output, phoneContact)
		rank++
	}

	emailSystem := domain.ContactPointSystemEnumEmail

	for _, email := range emails {
		emailErr := utils.ValidateEmail(*email.Email)
		if emailErr != nil {
			return nil, fmt.Errorf("invalid email: %w", emailErr)
		}

		emailContact := &domain.FHIRContactPointInput{
			System: &emailSystem,
			Use:    &use,
			Rank:   &rank,
			Period: common.DefaultPeriodInput(),
			Value:  email.Email,
		}
		output = append(output, emailContact)
		rank++
	}

	return output, nil
}

// ConceptMapper maps ObservationConceptEnum to a terminology code
func (c *FoundationImpl) ConceptMapper(concept dto.ObservationConceptEnum) (string, string, error) {
	var code, category string

	switch concept {
	case "TEMPERATURE":
		code = common.TemperatureLOINCTerminologyCode
		category = "vital-signs"
	case "HEIGHT":
		code = common.HeightLOINCTerminologyCode
		category = "vital-signs"
	case "WEIGHT":
		code = common.WeightLOINCTerminologyCode
		category = "vital-signs"
	case "RESPIRATORY_RATE":
		code = common.RespiratoryRateLOINCTerminologyCode
		category = "vital-signs"
	case "PULSE_RATE":
		code = common.PulseLOINCTerminologyCode
		category = "vital-signs"
	case "BLOOD_PRESSURE":
		code = common.BloodPressureLOINCTerminologyCode
		category = "vital-signs"
	case "BMI":
		code = common.BMILOINCTerminologyCode
		category = "vital-signs"
	case "VIRAL_LOAD":
		code = common.ViralLoadLOINCTerminologyCode
		category = "vital-signs"
	case "MUAC":
		code = common.MuacLOINCTerminologyCode
		category = "vital-signs"
	case "OXYGEN_SATURATION":
		code = common.OxygenSaturationLOINCTerminologyCode
		category = "vital-signs"
	case "BLOOD_SUGAR":
		code = common.BloodSugarLOINCTerminologyCode
		category = "vital-signs"
	case "LAST_MENSTRUAL_PERIOD":
		code = common.LastMenstrualPeriodLOINCTerminologyCode
		category = "vital-signs"
	case "DIASTOLIC_BLOOD_PRESSURE":
		code = common.DiastolicBloodPressureLOINCTerminologyCode
		category = "vital-signs"
	case "COLPOSCOPY":
		code = common.ColposcopyLOINCTerminologyCode
		category = "vital-signs"
	case "VIA":
		code = common.VIALOINCCode
		category = "exam"
	case "IMMUNO_HISTO_CHEMISTRY":
		code = common.ImmunoHistoChemistryLOINCCode
		category = "laboratory"
	case "POST_COITAL_BLEEDING":
		code = common.PostCoitalBleedingCIELCode
		category = "laboratory"
	case "HISTORY_OF_PRESENTING_ILLNESS":
		code = common.LOINCHistoryOfPresentingIllness
		category = "social-history"
	case "PAST_MEDICAL_AND_SURGICAL_HISTORY":
		code = common.LOINCPastMedicalSurgeryHistory
		category = "social-history"
	case "CHIEF_COMPLAINT":
		code = common.LOINCChiefComplaint
		category = "social-history"
	case "FAMILY_AND_SOCIAL_HISTORY":
		code = common.LOINCFamilyHistory
		category = "social-history"
	case "HPV":
		code = common.HPVLOINCTerminologyCode
		category = "exam"
	case "CBE":
		code = common.BreastExaminationLOINCTerminologySystem
		category = "exam"
	case "GENERAL_EXAMINATION":
		code = common.LOINCGeneralStatusNarrative
		category = "exam"
	default:
		code = "UNKNOWN"
	}

	if code == "UNKNOWN" {
		return "", "", fmt.Errorf("invalid concept provided")
	}

	return code, category, nil
}

func MapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(fhirAllergyIntolerance domain.FHIRAllergyIntolerance) *dto.Allergy {
	allergyIntolerance := &dto.Allergy{}
	if fhirAllergyIntolerance.Code != nil {
		allergyIntolerance.Code = string(*fhirAllergyIntolerance.GetConceptData().Code)
		allergyIntolerance.Name = fhirAllergyIntolerance.GetConceptData().Display
		allergyIntolerance.System = string(*fhirAllergyIntolerance.GetConceptData().System)
	}

	if fhirAllergyIntolerance.ID != nil {
		allergyIntolerance.ID = *fhirAllergyIntolerance.ID
	}

	if fhirAllergyIntolerance.Patient != nil && fhirAllergyIntolerance.Patient.ID != nil {
		allergyIntolerance.PatientID = *fhirAllergyIntolerance.Patient.ID
	}

	if fhirAllergyIntolerance.Encounter != nil && fhirAllergyIntolerance.Encounter.ID != nil {
		allergyIntolerance.EncounterID = *fhirAllergyIntolerance.Encounter.ID
	}

	if fhirAllergyIntolerance.OnsetPeriod != nil {
		allergyIntolerance.OnsetDateTime = fhirAllergyIntolerance.OnsetPeriod.Start
	}

	if len(fhirAllergyIntolerance.Reaction) > 0 {
		reaction := fhirAllergyIntolerance.Reaction[0]
		if reaction.Severity != nil {
			allergyIntolerance.Reaction.Severity = dto.AllergyIntoleranceReactionSeverityEnum(*reaction.Severity)
		}

		if len(reaction.Manifestation) > 0 {
			manifestation := reaction.Manifestation[0]
			if len(manifestation.Concept.Coding) > 0 {
				coding := manifestation.Concept.Coding[0]
				if coding.System != nil {
					allergyIntolerance.Reaction.System = string(*coding.System)
				}

				allergyIntolerance.Reaction.Code = string(*coding.Code)
				allergyIntolerance.Reaction.Name = string(coding.Display)
			}
		}
	}

	return allergyIntolerance
}

func MapFHIRObservationToObservationDTO(fhirObservation domain.FHIRObservation) *dto.Observation {
	var value string

	if fhirObservation.ValueQuantity != nil {
		value = fmt.Sprintf("%v %v", fhirObservation.ValueQuantity.Value, fhirObservation.ValueQuantity.Unit)
	}

	if fhirObservation.ValueCodeableConcept != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueCodeableConcept)
	}

	if fhirObservation.ValueString != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueString)
	}

	if fhirObservation.ValueBoolean != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueBoolean)
	}

	if fhirObservation.ValueInteger != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueInteger)
	}

	if fhirObservation.ValueRange != nil {
		value = fmt.Sprintf("%v %v - %v %v", fhirObservation.ValueRange.High.Value, fhirObservation.ValueRange.High.Unit, fhirObservation.ValueRange.Low.Value, fhirObservation.ValueRange.Low.Unit)
	}

	if fhirObservation.ValueRatio != nil {
		value = fmt.Sprintf("%v %v : %v %v", fhirObservation.ValueRatio.Numerator.Value, fhirObservation.ValueRatio.Numerator.Unit, fhirObservation.ValueRatio.Denominator.Value, fhirObservation.ValueRatio.Denominator.Unit)
	}

	if fhirObservation.ValueSampledData != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueSampledData.ID)
	}

	if fhirObservation.ValueTime != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueTime)
	}

	if fhirObservation.ValueDateTime != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueDateTime)
	}

	if fhirObservation.ValuePeriod != nil {
		value = fmt.Sprintf("%v - %v", fhirObservation.ValuePeriod.Start, fhirObservation.ValuePeriod.End)
	}

	obs := &dto.Observation{
		ID:     *fhirObservation.ID,
		Status: dto.ObservationStatus(strings.ToUpper(fhirObservation.Status.String())),
		Value:  value,
	}

	// Coding() returns nil when Code carries an empty Coding slice, so guarding on
	// Code alone is not enough - that combination panicked here whenever an
	// observation was stored without a resolved coding.
	if coding := fhirObservation.Coding(); coding != nil {
		if coding.Code != nil {
			obs.Code = string(*coding.Code)
		}

		obs.Name = coding.Display
	}

	if len(fhirObservation.Category) > 0 {
		obs.Category = fhirObservation.Category[0].Text
	}

	if fhirObservation.Text != nil {
		obs.UsageContext = dto.ScreeningTypeEnum(helpers.ExtractTextFromHTML(string(fhirObservation.Text.Div)))
	}

	if fhirObservation.Subject != nil && fhirObservation.Subject.ID != nil {
		obs.PatientID = *fhirObservation.Subject.ID
	}

	if fhirObservation.Encounter != nil && fhirObservation.Encounter.Display != "" {
		obs.EncounterID = fhirObservation.Encounter.Display
	}

	if fhirObservation.Note != nil {
		obs.Note = string(*fhirObservation.Note[0].Text)
	}

	if fhirObservation.EffectiveInstant != nil {
		obs.TimeRecorded = string(*fhirObservation.EffectiveInstant)
	}

	for _, interpretation := range fhirObservation.Interpretation {
		obs.Interpretation = append(obs.Interpretation, interpretation.Text)
	}

	return obs
}
