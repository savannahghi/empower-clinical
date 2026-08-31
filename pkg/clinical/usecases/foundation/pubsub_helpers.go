package foundation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// GetConcept is a helper function that returns a concept associated the terminology source passed
func (c *FoundationImpl) GetConcept(ctx context.Context, terminologySource domain.TerminologySource, conceptID string) (*domain.Concept, error) {
	var (
		organisation string
		source       string
	)

	switch terminologySource {
	case domain.TerminologySourceICD10WHO:
		organisation = "WHO"
		source = "ICD-10-WHO"

	case domain.TerminologySourceICD10:
		organisation = "WHO"
		source = "ICD-10-WHO"

	case domain.TerminologySourceICD11WHO:
		organisation = "WHO"
		source = "ICD-11-WHO"

	case domain.TerminologySourceCIEL:
		organisation = "CIEL"
		source = "CIEL"

	case domain.TerminologySourceLOINC:
		organisation = "Regenstrief"
		source = "LOINC"

	default:
		return nil, fmt.Errorf("terminology source %s not supported", source)
	}

	response, err := c.OpenConceptLab.GetConcept(
		ctx,
		organisation,
		source,
		conceptID,
		false,
		false,
	)
	if err != nil {
		return nil, err
	}

	var concept *domain.Concept

	err = mapstructure.Decode(response, &concept)
	if err != nil {
		return nil, err
	}

	return concept, nil
}

// ComposeVitalsInput composes a vitals observation from data received
func (c *FoundationImpl) ComposeVitalsInput(ctx context.Context, input dto.VitalSignPubSubMessage) (*domain.FHIRObservationInput, error) {
	vitalsConcept, err := c.GetConcept(ctx, domain.TerminologySourceCIEL, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	system := "http://terminology.hl7.org/CodeSystem/observation-category"
	status := domain.ObservationStatusEnumPreliminary
	instant := scalarutils.Instant(input.Date.Format(time.RFC3339))
	observation := domain.FHIRObservationInput{
		Status: &status,
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&system),
						Code:    "vital-signs",
						Display: "Vital Signs",
					},
				},
				Text: "Vital Signs",
			},
		},
		EffectiveInstant: &instant,
		Code: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&vitalsConcept.URL),
					Code:    scalarutils.Code(vitalsConcept.ID),
					Display: vitalsConcept.DisplayName,
				},
			},
			Text: vitalsConcept.DisplayName,
		},
		ValueString: &input.Value,
	}

	patient, err := c.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	patientName := *patient.Resource.Name[0].Given[0]
	observation.Subject = &domain.FHIRReferenceInput{
		ID:        patient.Resource.ID,
		Reference: &patientReference,
		Display:   patientName,
	}

	if input.FacilityID != "" {
		facility, err := c.FHIR.GetFHIROrganization(ctx, input.FacilityID)
		if err != nil {
			// Should not fail if facility is not found
			log.Printf("the error is: %v", err)
		}

		if facility != nil {
			performerReference := fmt.Sprintf("Organization/%s", *facility.Resource.ID)
			referenceInput := &domain.FHIRReferenceInput{
				ID:        facility.Resource.ID,
				Reference: &performerReference,
				Display:   *facility.Resource.Name,
			}

			observation.Performer = append(observation.Performer, referenceInput)
		}
	}

	return &observation, nil
}

// ComposeAllergyIntoleranceInput composes an allergy intolerance input from the data received
func (c *FoundationImpl) ComposeAllergyIntoleranceInput(ctx context.Context, input dto.PatientAllergyPubSubMessage) (*domain.FHIRAllergyIntoleranceInput, error) {
	allergyType := domain.AllergyIntoleranceTypeEnumAllergy
	allergyCategory := domain.AllergyIntoleranceCategoryEnumMedication
	allergy := &domain.FHIRAllergyIntoleranceInput{
		Type: &domain.FHIRCodeableConceptInput{
			Text: allergyType.String(),
		},
		Category: []*domain.AllergyIntoleranceCategoryEnum{&allergyCategory},
		ClinicalStatus: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&FhirAllergyIntoleranceClinicalStatusURL),
					Code:    scalarutils.Code("active"),
					Display: "Active",
				},
			},
			Text: "Active",
		},
		VerificationStatus: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&fhirAllergyIntoleranceVerificationStatusURL),
					Code:    scalarutils.Code("confirmed"),
					Display: "Confirmed",
				},
			},
			Text: "Confirmed",
		},
		Reaction: []*domain.FHIRAllergyintoleranceReactionInput{},
	}

	date := input.Date.Format(time.DateOnly)
	allergy.RecordedDate = (*scalarutils.DateTime)(&date)

	patient, err := c.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	subjectReference := fmt.Sprintf("Patient/%s", input.PatientID)
	patientName := *patient.Resource.Name[0].Given[0]

	allergy.Patient = &domain.FHIRReferenceInput{
		ID:        patient.Resource.ID,
		Reference: &subjectReference,
		Display:   patientName,
	}

	allergenConcept, err := c.GetConcept(ctx, domain.TerminologySourceCIEL, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	allergy.Code = domain.FHIRCodeableConceptInput{
		Coding: []*domain.FHIRCodingInput{
			{
				System:  (*scalarutils.URI)(&allergenConcept.URL),
				Code:    scalarutils.Code(allergenConcept.ID),
				Display: allergenConcept.DisplayName,
			},
		},
		Text: allergenConcept.DisplayName,
	}

	// create the allergy reaction
	var reaction domain.FHIRAllergyintoleranceReactionInput

	// reaction manifestation is required
	//
	// check if there is a reaction manifestation,
	// if no reaction use unknown
	var manifestationConcept *domain.Concept
	if input.Reaction.ConceptID != nil {
		manifestationConcept, err = c.GetConcept(ctx, domain.TerminologySourceCIEL, *input.Reaction.ConceptID)
		if err != nil {
			return nil, err
		}
	} else {
		manifestationConcept, err = c.GetConcept(ctx, domain.TerminologySourceCIEL, unknownConceptID)
		if err != nil {
			return nil, err
		}
	}

	manifestation := &domain.FHIRCodeableReferenceInput{
		Concept: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&manifestationConcept.URL),
					Code:    scalarutils.Code(manifestationConcept.ID),
					Display: manifestationConcept.DisplayName,
				},
			},
			Text: manifestationConcept.DisplayName,
		},
	}

	// add reaction manifestation
	reaction.Manifestation = append(reaction.Manifestation, manifestation)

	if input.Severity.ConceptID != nil {
		severityConcept, err := c.GetConcept(ctx, domain.TerminologySourceCIEL, *input.Severity.ConceptID)
		if err != nil {
			return nil, err
		}

		reaction.Description = &severityConcept.DisplayName
	}

	// add allergy reaction
	allergy.Reaction = append(allergy.Reaction, &reaction)

	return allergy, nil
}

// ComposeTestResultInput composes a test result input from data received
func (c *FoundationImpl) ComposeTestResultInput(ctx context.Context, input dto.PatientTestResultPubSubMessage) (*domain.FHIRObservationInput, error) {
	var patientName string

	patient, err := c.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	patientName = *patient.Resource.Name[0].Given[0]

	observationConcept, err := c.GetConcept(ctx, domain.TerminologySourceCIEL, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	system := "http://terminology.hl7.org/CodeSystem/observation-category"
	subjectReference := fmt.Sprintf("Patient/%s", input.PatientID)
	status := domain.ObservationStatusEnumPreliminary
	instant := scalarutils.Instant(input.Date.Format(time.RFC3339))

	observation := domain.FHIRObservationInput{
		Status: &status,
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&system),
						Code:    "laboratory",
						Display: "Laboratory",
					},
				},
				Text: "Laboratory",
			},
		},
		Code: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&observationConcept.URL),
					Code:    scalarutils.Code(observationConcept.ID),
					Display: observationConcept.DisplayName,
				},
			},
			Text: observationConcept.DisplayName,
		},
		ValueString:      &input.Result.Name,
		EffectiveInstant: &instant,
		Subject: &domain.FHIRReferenceInput{
			ID:        patient.Resource.ID,
			Reference: &subjectReference,
			Display:   patientName,
		},
	}

	if input.FacilityID != "" {
		facility, err := c.FHIR.GetFHIROrganization(ctx, input.FacilityID)
		if err != nil {
			// Should not fail if the facility is not found
			log.Printf("the error is: %v", err)
		}

		if facility != nil {
			performer := fmt.Sprintf("Organization/%s", *facility.Resource.ID)

			referenceInput := &domain.FHIRReferenceInput{
				ID:        facility.Resource.ID,
				Reference: &performer,
				Display:   *facility.Resource.Name,
			}

			observation.Performer = append(observation.Performer, referenceInput)
		}
	}

	return &observation, nil
}

// ComposeMedicationStatementInput composes a medication statement input from received data
func (c *FoundationImpl) ComposeMedicationStatementInput(ctx context.Context, input dto.MedicationPubSubMessage) (*domain.FHIRMedicationStatementInput, error) {
	medicationConcept, err := c.GetConcept(ctx, domain.TerminologySourceCIEL, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	drugConcept, err := c.GetConcept(ctx, domain.TerminologySourceCIEL, *input.Drug.ConceptID)
	if err != nil {
		return nil, err
	}

	year, month, day := input.Date.Date()
	status := domain.MedicationStatementStatusEnumUnknown
	medicationStatement := domain.FHIRMedicationStatementInput{
		Status: &status,
		Category: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&medicationConcept.URL),
					Code:    scalarutils.Code(medicationConcept.ID),
					Display: medicationConcept.DisplayName,
				},
			},
			Text: medicationConcept.DisplayName,
		},
		MedicationCodeableConcept: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&drugConcept.URL),
					Code:    scalarutils.Code(drugConcept.ID),
					Display: drugConcept.DisplayName,
				},
			},
			Text: drugConcept.DisplayName,
		},
		EffectiveDateTime: &scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		},
	}

	patient, err := c.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	patientName := *patient.Resource.Name[0].Given[0]
	medicationStatement.Subject = &domain.FHIRReferenceInput{
		ID:        patient.Resource.ID,
		Reference: &patientReference,
		Display:   patientName,
	}

	if input.FacilityID != "" {
		facility, err := c.FHIR.GetFHIROrganization(ctx, input.FacilityID)
		if err != nil {
			log.Printf("the error is: %v", err)
		}

		if facility != nil {
			informationSourceReference := fmt.Sprintf("Organization/%s", *facility.Resource.ID)

			reference := &domain.FHIRReferenceInput{
				ID:        facility.Resource.ID,
				Reference: &informationSourceReference,
				Display:   *facility.Resource.Name,
			}

			medicationStatement.InformationSource = reference
		}
	}

	return &medicationStatement, nil
}
