package clinical

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// diagnosisTreatmentLinkageExtensionURL identifies the FHIR extension that stores the
// retrospective diagnosis-to-treatment linkage captured via RecordTreatmentEnrollment. The nested
// sub-extension URLs key the individual data points recorded within it.
var (
	diagnosisTreatmentLinkageExtensionURL = "http://savannahghi.org/fhir/StructureDefinition/diagnosis-treatment-linkage"
	linkedToTreatmentExtensionURL         = "linkedToTreatment"
	treatmentFacilityExtensionURL         = "treatmentFacility"
	treatmentProgramExtensionURL          = "treatmentProgram"
	enrollmentDateExtensionURL            = "enrollmentDate"

	// linkedRecordURL is set as a Condition's Meta.source to mark records captured
	// retrospectively via RecordTreatmentEnrollment. Meta.source is searchable through the FHIR `_source`
	// parameter, so listing can filter to these records without colliding with the tenant `_tag`.
	linkedRecordURL = "http://savannahghi.org/fhir/record-source/linkage"

	// sourceStrategy is the ListPatientConditions strategy value that restricts results
	// to records captured retrospectively (matched via linkedRecordURL).
	sourceStrategy = "linkage"
	severitySystem = scalarutils.URI("http://terminology.hl7.org/CodeSystem/adverse-event-severity")
)

// CreateCondition creates a new conditions
func (c *ClinicalImpl) CreateCondition(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error) {
	today := time.Now()

	date, err := scalarutils.NewDate(today.Day(), int(today.Month()), today.Year())
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record a condition in a completed encounter")
	}

	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
	categorySystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-category")
	verificationSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-ver-status")
	userSelectedFalse := false
	conditionInput := domain.FHIRConditionInput{
		ClinicalStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  &statusSystem,
					Code:    scalarutils.Code(strings.ToLower(string(input.Status))),
					Display: string(input.Status),
				},
			},
			Text: string(input.Status),
		},
		VerificationStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &verificationSystem,
					Code:         scalarutils.Code("confirmed"),
					Display:      "confirmed",
					UserSelected: &userSelectedFalse,
				},
			},
			Text: "confirmed",
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &categorySystem,
						Code:         scalarutils.Code(input.Category.ToCategoryCode()),
						Display:      input.Category.ToString(),
						UserSelected: &userSelectedFalse,
					},
				},
				Text: input.Category.ToString(),
			},
		},
		RecordedDate: date,
		Subject: &domain.FHIRReferenceInput{
			Reference: encounter.Resource.Subject.Reference,
			Display:   encounter.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterRef,
			Display:   *encounter.Resource.ID,
		},
		Severity: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  &severitySystem,
					Code:    scalarutils.Code(input.Severity),
					Display: input.Severity.ToString(),
				},
			},
			Text: input.Severity.ToString(),
		},
		Participant: []*domain.FHIRConditionParticipantInput{
			{
				Actor: &domain.FHIRReferenceInput{
					Reference: encounter.Resource.Subject.Reference,
					Display:   encounter.Resource.Subject.Display,
				},
			},
		},
	}

	var allowedTerminologySources = map[string]string{
		string(domain.TerminologySourceICD11WHO): common.SystemURLICD11,
		string(domain.TerminologySourceICHI):     common.SystemURLICHI,
	}

	err = conditionInput.ComposeConditionCodeField(input.Code, input.Name, string(input.System), allowedTerminologySources)
	if err != nil {
		return nil, err
	}

	if input.OnsetDate != nil {
		conditionInput.OnsetDateTime = input.OnsetDate
	}

	if input.Note != "" {
		note := scalarutils.Markdown(input.Note)
		noteTime := scalarutils.DateTime(time.Now().Format(scalarutils.DateTimeFormatLayout))

		conditionInput.Note = []*domain.FHIRAnnotationInput{
			{
				Time: &noteTime,
				Text: &note,
			},
		}
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	conditionInput.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	condition, err := c.FHIR.CreateFHIRCondition(ctx, conditionInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRConditionToConditionDTO(*condition.Resource), nil
}

func mapFHIRConditionToConditionDTO(condition domain.FHIRCondition) *dto.Condition {
	output := dto.Condition{
		ID:           *condition.ID,
		Status:       dto.ConditionStatus(condition.ClinicalStatus.Text),
		Name:         condition.Code.Text,
		RecordedDate: condition.RecordedDate,
		EncounterID:  helpers.StripIDFromReference(condition.Encounter),
		PatientID:    helpers.StripIDFromReference(condition.Subject),
	}

	oncology := dto.OncologyConditionOutput{}

	if len(condition.Category) > 0 {
		switch condition.Category[0].Text {
		case dto.ConditionCategoryProblemList.ToString():
			output.Category = "problem-list-item"
		case dto.ConditionCategoryDiagnosis.ToString():
			output.Category = "encounter-diagnosis"
		}

		for _, cat := range condition.Category {
			if strings.Contains(cat.Text, "Primary Tumor Code") {
				oncology.ICDO3PrimaryTumorCode = cat.Text
			} else if strings.Contains(cat.Text, "Morphology Code") {
				oncology.ICDO3MorphologyCode = cat.Text
			}
		}
	}

	if condition.Code != nil && condition.Code.Coding != nil {
		for _, coding := range condition.Code.Coding {
			if coding.Code != nil {
				output.Code = string(*coding.Code)
			}

			if coding.System != nil {
				output.System = string(*coding.System)
			}
		}
	}

	if len(condition.Note) > 0 && condition.Note[0].Text != nil {
		output.Note = string(*condition.Note[0].Text)
	}

	if condition.OnsetDateTime != nil {
		output.OnsetDate = condition.OnsetDateTime
	}

	if len(condition.Stage) > 0 {
		oncology.Stage = condition.GetStage()
	}

	// Only attach the oncology block when it actually carries data, so non-oncology conditions omit
	// it (the field is a pointer with omitempty).
	if oncology != (dto.OncologyConditionOutput{}) {
		output.OncologyCondition = &oncology
	}

	output.TreatmentLinkage = treatmentLinkageFromExtensions(condition.Extension)

	return &output
}

// treatmentLinkageFromExtensions extracts the diagnosis-to-treatment linkage recorded by
// RecordTreatmentEnrollment from a condition's extensions. It returns nil when the condition does not
// carry the linkage extension.
func treatmentLinkageFromExtensions(extensions []*domain.FHIRExtension) *dto.TreatmentLinkageOutput {
	for _, extension := range extensions {
		if extension == nil || extension.URL != diagnosisTreatmentLinkageExtensionURL {
			continue
		}

		linkage := &dto.TreatmentLinkageOutput{}

		for _, field := range extension.Extension {
			switch field.URL {
			case linkedToTreatmentExtensionURL:
				if field.ValueBoolean != nil {
					linkage.LinkedToTreatment = *field.ValueBoolean
				}
			case treatmentFacilityExtensionURL:
				linkage.TreatmentFacility = field.ValueString
			case treatmentProgramExtensionURL:
				linkage.TreatmentProgram = field.ValueString
			case enrollmentDateExtensionURL:
				if enrollmentDate, err := time.Parse(time.DateOnly, field.ValueDate); err == nil {
					if date, err := scalarutils.NewDate(enrollmentDate.Day(), int(enrollmentDate.Month()), enrollmentDate.Year()); err == nil {
						linkage.EnrollmentDate = date
					}
				}
			}
		}

		return linkage
	}

	return nil
}

// RecordTreatmentEnrollment links a patient diagnosis to treatment
func (c ClinicalImpl) RecordTreatmentEnrollment(ctx context.Context, input *dto.TreatmentEnrollmentInput) (*dto.Condition, error) {
	if err := helpers.Validate(input); err != nil {
		return nil, err
	}

	// When a diagnosis is linked to treatment the treating facility is mandatory. The treatment
	// program and enrolment date remain optional and are only recorded when supplied.
	if input.LinkedToTreatment && input.TreatmentFacility == "" {
		return nil, fmt.Errorf("treatment facility must be provided when the diagnosis is linked to treatment")
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted || encounter.Resource.Status == domain.EncounterStatusEnumCancelled {
		return nil, fmt.Errorf("cannot record a condition in a completed or cancelled encounter")
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	var (
		encounterRef       = fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
		statusSystem       = scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
		categorySystem     = scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-category")
		verificationSystem = scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-ver-status")
		userSelectedFalse  = false
	)

	conditionInput := domain.FHIRConditionInput{
		Meta: domain.FHIRMetaInput{
			Tag:    tags,
			Source: linkedRecordURL,
		},
		ClinicalStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  &statusSystem,
					Code:    scalarutils.Code("active"),
					Display: "Active",
				},
			},
			Text: "Active",
		},
		VerificationStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &verificationSystem,
					Code:         scalarutils.Code("confirmed"),
					Display:      "Confirmed",
					UserSelected: &userSelectedFalse,
				},
			},
			Text: "Confirmed",
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &categorySystem,
						Code:         scalarutils.Code("problem-list-item"),
						Display:      "Problem List Item",
						UserSelected: &userSelectedFalse,
					},
				},
				Text: "Problem List Item",
			},
		},
		Severity: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  &severitySystem,
					Code:    scalarutils.Code(input.Severity),
					Display: input.Severity.ToString(),
				},
			},
			Text: input.Severity.ToString(),
		},
		Code:         c.defaultConditionCode(input.Condition, &userSelectedFalse),
		RecordedDate: input.Date,
		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: encounter.Resource.Subject.Reference,
		},
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterRef,
			ID:        encounter.Resource.ID,
		},
		Participant: []*domain.FHIRConditionParticipantInput{
			{
				Actor: &domain.FHIRReferenceInput{
					Reference: encounter.Resource.Subject.Reference,
					Display:   encounter.Resource.Subject.Display,
				},
			},
		},
	}

	// Record the cancer stage as a Condition.stage summary when supplied. The retrospective flow does
	// not create supporting diagnostic evidence, so no stage assessment references are attached.
	if input.CancerStage.Code != "" || input.CancerStage.Display != "" {
		version := "0.1.0"
		conditionInput.Stage = []*domain.FHIRConditionStageInput{
			{
				Summary: &domain.FHIRCodeableConceptInput{
					Coding: []*domain.FHIRCodingInput{
						{
							System:       helpers.CodeSystem(common.CancerStagesCodeSystem),
							Code:         scalarutils.Code(input.CancerStage.Code),
							Display:      input.CancerStage.Display,
							UserSelected: &userSelectedFalse,
							Version:      &version,
						},
					},
					Text: input.CancerStage.Display,
				},
			},
		}
	}

	// Persist the diagnosis-treatment linkage as a FHIR extension so it can be retrospectively
	// queried for reporting and analytics. Only fields that were supplied are recorded.
	linkage := []domain.Extension{
		{
			URL:          linkedToTreatmentExtensionURL,
			ValueBoolean: &input.LinkedToTreatment,
		},
	}

	if input.TreatmentFacility != "" {
		linkage = append(linkage, domain.Extension{
			URL:         treatmentFacilityExtensionURL,
			ValueString: input.TreatmentFacility,
		})
	}

	if input.TreatmentProgram != "" {
		linkage = append(linkage, domain.Extension{
			URL:         treatmentProgramExtensionURL,
			ValueString: input.TreatmentProgram,
		})
	}

	if input.EnrollmentDate != nil {
		linkage = append(linkage, domain.Extension{
			URL:       enrollmentDateExtensionURL,
			ValueDate: input.EnrollmentDate.AsTime().Format(time.DateOnly),
		})
	}

	conditionInput.Extension = []*domain.FHIRExtension{
		{
			URL:       diagnosisTreatmentLinkageExtensionURL,
			Extension: linkage,
		},
	}

	condition, err := c.FHIR.CreateFHIRCondition(ctx, conditionInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRConditionToConditionDTO(*condition.Resource), nil
}

// UpdateTreatmentEnrollment updates the mutable fields of a treatment enrollment previously recorded
// as a FHIR Condition by RecordTreatmentEnrollment. Only the condition code, recorded date and
// enrollment date may change, and only the fields supplied on the input are touched; the rest of the
// record is left as-is. The target must be a treatment enrollment — it must carry the
// diagnosis-treatment linkage extension — otherwise the update is refused.
func (c ClinicalImpl) UpdateTreatmentEnrollment(ctx context.Context, id string, input *dto.UpdateTreatmentEnrollmentInput) (*dto.Condition, error) {
	if id == "" {
		return nil, fmt.Errorf("treatment enrollment id is required")
	}

	if input == nil || (input.Condition == nil && input.Date == nil && input.EnrollmentDate == nil) {
		return nil, fmt.Errorf("at least one of condition, date or enrollment_date must be provided")
	}

	existing, err := c.FHIR.GetFHIRCondition(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing == nil || existing.Resource == nil {
		return nil, fmt.Errorf("treatment enrollment %s was not found", id)
	}

	// Only records captured as treatment enrollments carry the diagnosis-treatment linkage
	// extension. Refuse to mutate any other condition through this endpoint.
	if treatmentLinkageFromExtensions(existing.Resource.Extension) == nil {
		return nil, fmt.Errorf("condition %s is not a treatment enrollment", id)
	}

	// The update is applied as a FHIR path patch that replaces each supplied top-level field. Meta is
	// a value field that always serialises, so it must be repopulated to avoid clearing the record's
	// tenant tags; rebuild it exactly as RecordTreatmentEnrollment does.
	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	userSelectedFalse := false

	conditionInput := domain.FHIRConditionInput{
		ID: &id,
		Meta: domain.FHIRMetaInput{
			Tag:    tags,
			Source: linkedRecordURL,
		},
	}

	if input.Condition != nil {
		conditionInput.Code = c.defaultConditionCode(*input.Condition, &userSelectedFalse)
	}

	if input.Date != nil {
		conditionInput.RecordedDate = input.Date
	}

	if input.EnrollmentDate != nil {
		conditionInput.Extension = withEnrollmentDate(existing.Resource.Extension, *input.EnrollmentDate)
	}

	updated, err := c.FHIR.UpdateFHIRCondition(ctx, conditionInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRConditionToConditionDTO(*updated.Resource), nil
}

// withEnrollmentDate returns the condition's extensions with the diagnosis-treatment linkage's
// enrollment date set to the supplied date, preserving every other linkage sub-extension. The
// enrollment date sub-extension is added when it was not previously recorded. The existing slice is
// mutated in place and returned for convenience.
func withEnrollmentDate(extensions []*domain.FHIRExtension, date scalarutils.Date) []*domain.FHIRExtension {
	formatted := date.AsTime().Format(time.DateOnly)

	for _, extension := range extensions {
		if extension == nil || extension.URL != diagnosisTreatmentLinkageExtensionURL {
			continue
		}

		for i := range extension.Extension {
			if extension.Extension[i].URL == enrollmentDateExtensionURL {
				extension.Extension[i].ValueDate = formatted

				return extensions
			}
		}

		extension.Extension = append(extension.Extension, domain.Extension{
			URL:       enrollmentDateExtensionURL,
			ValueDate: formatted,
		})
	}

	return extensions
}

func (c ClinicalImpl) RecordOncologicalDiagnosis(ctx context.Context, input *dto.OncologyDiagnosisInput) (*dto.Condition, error) {
	if err := helpers.Validate(input); err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// observation for grading
	observation, err := c.createGradingObservation(ctx, input, encounter, facilityID, tags)
	if err != nil {
		return nil, err
	}

	// diagnostic report for behavior
	diagnosticReport, err := c.createDiagnosticReport(ctx, input, encounter.Resource, observation, facilityID, tags)
	if err != nil {
		return nil, err
	}

	// condition
	condition, err := c.createCondition(ctx, input, encounter.Resource, observation, diagnosticReport, facilityID, tags)
	if err != nil {
		return nil, err
	}

	return mapFHIRConditionToConditionDTO(*condition.Resource), nil
}

func (c ClinicalImpl) createGradingObservation(ctx context.Context, input *dto.OncologyDiagnosisInput, encounter *domain.FHIREncounterRelayPayload, facilityID string, tags []domain.FHIRCodingInput) (*domain.FHIRObservation, error) {
	var (
		instant           = scalarutils.Instant(time.Now().Format(time.RFC3339))
		encounterRef      = fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
		orgRef            = fmt.Sprintf("Organization/%s", facilityID)
		observationSystem = scalarutils.URI("http://terminology.hl7.org/CodeSystem/observation-category")
		userSelectedFalse = false
		observationStatus = "final"
	)

	observation := domain.FHIRObservationInput{
		Status: (*domain.ObservationStatusEnum)(&observationStatus),
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &observationSystem,
						Code:         "laboratory",
						Display:      "Laboratory",
						UserSelected: &userSelectedFalse,
					},
				},
				Text: "Laboratory",
			},
		},
		Code: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier),
					Code:         scalarutils.Code(input.Grade.Code),
					Display:      input.Grade.Display,
					UserSelected: &userSelectedFalse,
				},
			},
			Text: input.Grade.Display,
		},
		EffectiveInstant: &instant,
		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: encounter.Resource.Subject.Reference,
			Display:   encounter.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Reference: &encounterRef,
			Display:   *encounter.Resource.ID,
		},
		Issued: &instant,
		Performer: []*domain.FHIRReferenceInput{
			{
				Reference: &orgRef,
				Display:   facilityID,
			},
		},
		ValueCodeableConcept: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier),
					Code:         scalarutils.Code(input.Grade.Code),
					Display:      input.Grade.Display,
					UserSelected: &userSelectedFalse,
				},
			},
			Text: input.Grade.Display,
		},
		Meta: &domain.FHIRMetaInput{
			Tag: tags,
		},
	}

	return c.FHIR.CreateFHIRObservation(ctx, observation)
}

func (c ClinicalImpl) createDiagnosticReport(ctx context.Context, input *dto.OncologyDiagnosisInput, encounter *domain.FHIREncounter,
	observation *domain.FHIRObservation, facilityID string, tags []domain.FHIRCodingInput) (*domain.FHIRDiagnosticReport, error,
) {
	var (
		patientRef                         = fmt.Sprintf("Patient/%s", *encounter.Subject.ID)
		encounterRef                       = fmt.Sprintf("Encounter/%s", *encounter.ID)
		obsRef                             = fmt.Sprintf("Observation/%s", *observation.ID)
		orgRef                             = fmt.Sprintf("Organization/%s", facilityID)
		instant                            = scalarutils.Instant(time.Now().Format(time.RFC3339))
		dateTime                           = scalarutils.DateTime(time.Now().Format(time.RFC3339))
		diagnosticReportCategoryCodeSystem = scalarutils.URI("http://terminology.hl7.org/CodeSystem/v2-0074")
		observationTypeSystem              = scalarutils.URI("http://terminology.hl7.org/CodeSystem/v2-0936")
		observationTypeCode                = "SCI"
		userSelectedFalse                  = false
	)

	diagnosticReport := &domain.FHIRDiagnosticReportInput{
		Status: domain.DiagnosticReportStatusFinal,
		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Subject.ID,
			Reference: &patientRef,
			Display:   *encounter.Subject.ID,
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &diagnosticReportCategoryCodeSystem,
						Code:         "LAB",
						Display:      "Laboratory",
						UserSelected: &userSelectedFalse,
					},
				},
				Text: "Laboratory",
			},
		},
		Code: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &common.LoincSystemURL,
					Code:         scalarutils.Code(common.CancerProgressionLOINCCode),
					Display:      "Cancer disease progression",
					UserSelected: &userSelectedFalse,
				},
			},
			Text: "Cancer disease progression",
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.ID,
			Reference: &encounterRef,
			Display:   *encounter.ID,
		},
		Issued: (*string)(&instant),
		Performer: []*domain.FHIRReferenceInput{
			{
				ID:        &facilityID,
				Reference: &orgRef,
				Display:   facilityID,
			},
		},
		ResultsInterpreter: []*domain.FHIRReferenceInput{
			{
				Reference: &orgRef,
				Display:   orgRef,
			},
		},
		EffectiveDateTime: &dateTime,
		Conclusion:        &input.Stage.Display,
		SupportingInfo: []domain.DiagnosticReportSupportingInfo{
			{
				Type: domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{
						{
							System:       &observationTypeSystem,
							Code:         (*scalarutils.Code)(&observationTypeCode),
							Display:      "Supporting Clinical Information",
							UserSelected: &userSelectedFalse,
						},
					},
					Text: input.Behavior.Display,
				},
				Reference: domain.FHIRReferenceInput{
					ID:        observation.ID,
					Reference: &obsRef,
					Display:   *observation.ID,
				},
			},
		},
		Meta: &domain.FHIRMetaInput{
			Tag: tags,
		},
	}

	return c.FHIR.CreateFHIRDiagnosticReport(ctx, diagnosticReport)
}

func (c ClinicalImpl) createCondition(ctx context.Context, input *dto.OncologyDiagnosisInput,
	encounter *domain.FHIREncounter, observation *domain.FHIRObservation, diagnosticReport *domain.FHIRDiagnosticReport,
	facilityID string, tags []domain.FHIRCodingInput) (*domain.FHIRConditionRelayPayload, error,
) {
	var (
		encounterRef        = fmt.Sprintf("Encounter/%s", *encounter.ID)
		diagnosticReportRef = fmt.Sprintf("DiagnosticReport/%s", *diagnosticReport.ID)
		obsRef              = fmt.Sprintf("Observation/%s", *observation.ID)
		orgRef              = fmt.Sprintf("Organization/%s", facilityID)
		dateTime            = scalarutils.DateTime(time.Now().Format(time.RFC3339))
		statusSystem        = scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
		categorySystem      = scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-category")
		verificationSystem  = scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-ver-status")
		userSelectedFalse   = false
	)

	today := time.Now()

	date, err := scalarutils.NewDate(today.Day(), int(today.Month()), today.Year())
	if err != nil {
		return nil, err
	}

	conditionInput := domain.FHIRConditionInput{
		Meta: domain.FHIRMetaInput{Tag: tags},
		ClinicalStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  &statusSystem,
					Code:    scalarutils.Code("active"),
					Display: "Active",
				},
			},
			Text: "Active",
		},
		VerificationStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &verificationSystem,
					Code:         scalarutils.Code("confirmed"),
					Display:      "Confirmed",
					UserSelected: &userSelectedFalse,
				},
			},
			Text: "Confirmed",
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  &categorySystem,
						Code:    scalarutils.Code("problem-list-item"),
						Display: "Problem List Item",
					},
				},
				Text: fmt.Sprintf("Primary Tumor Code - %s", input.ICDO3PrimaryTumorCode),
			},
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  &categorySystem,
						Code:    scalarutils.Code("problem-list-item"),
						Display: "Problem List Item",
					},
				},
				Text: fmt.Sprintf("Morphology Code - %s", input.ICDO3MorphologyCode),
			},
		},
		Code:         c.defaultConditionCode(input.Condition, &userSelectedFalse),
		RecordedDate: date,
		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Subject.ID,
			Reference: encounter.Subject.Reference,
			Display:   encounter.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.ID,
			Reference: &encounterRef,
			Display:   *encounter.ID,
		},
		Severity: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  helpers.CodeSystem("SGHIConditionSeverityCodeSystem"),
					Code:    scalarutils.Code("severe"),
					Display: "Severe",
				},
			},
			Text: "Severe",
		},
		Participant: []*domain.FHIRConditionParticipantInput{
			{
				Actor: &domain.FHIRReferenceInput{
					Reference: encounter.Subject.Reference,
					Display:   encounter.Subject.Display,
				},
			},
		},
		Stage: []*domain.FHIRConditionStageInput{
			{
				Summary: &domain.FHIRCodeableConceptInput{
					Coding: []*domain.FHIRCodingInput{
						{
							System:       helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier),
							Code:         scalarutils.Code(input.Stage.Code),
							Display:      input.Stage.Display,
							UserSelected: &userSelectedFalse,
						},
					},
					Text: input.Stage.Display,
				},
				Assessment: []*domain.FHIRReferenceInput{
					{
						ID:        diagnosticReport.ID,
						Reference: &diagnosticReportRef,
						Display:   *diagnosticReport.ID,
					},
					{
						ID:        observation.ID,
						Reference: &obsRef,
						Display:   *observation.ID,
					},
				},
			},
		},
		Note: []*domain.FHIRAnnotationInput{
			{
				AuthorReference: &domain.FHIRReferenceInput{
					ID:        &facilityID,
					Reference: &orgRef,
				},
				Time: &dateTime,
				Text: (*scalarutils.Markdown)(&input.Notes),
			},
		},
	}

	return c.FHIR.CreateFHIRCondition(ctx, conditionInput)
}

func (c ClinicalImpl) defaultConditionCode(condition dto.ValueSetData, userSelected *bool) *domain.FHIRCodeableConceptInput {
	var (
		system         = "http://id.who.int/icd/release/11/mms"
		defaultCode    = "sghidefaultcode"
		defaultDisplay = "SGHI Default Code"
		codeSystem     *scalarutils.URI
		codeValue      scalarutils.Code
		display        string
	)

	if condition.Code != "" {
		codeSystem = (*scalarutils.URI)(&system)
		codeValue = scalarutils.Code(condition.Code)
		display = condition.Display
	} else {
		codeSystem = helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier)
		codeValue = scalarutils.Code(defaultCode)
		display = defaultDisplay
	}

	return &domain.FHIRCodeableConceptInput{
		Coding: []*domain.FHIRCodingInput{
			{
				UserSelected: userSelected,
				System:       codeSystem,
				Code:         codeValue,
				Display:      display,
			},
		},
		Text: condition.Display,
	}
}

// ListPatientConditions lists a patients conditions
// TODO: pagination
func (c ClinicalImpl) ListPatientConditions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, strategy string, pagination dto.Pagination) (*dto.ConditionConnection, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	err = pagination.Validate()
	if err != nil {
		return nil, err
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	patient, err := c.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	params := map[string]interface{}{
		"subject": patientRef,
		"_sort":   "-_lastUpdated",
		"_total":  "accurate",
	}

	if encounterID != nil {
		encounterReference := fmt.Sprintf("Encounter/%s", *encounterID)
		params["encounter"] = encounterReference
	}

	if date != nil {
		params["recorded-date"] = date.AsTime().Format(time.DateOnly)
	}

	if strategy == sourceStrategy {
		params["_source"] = linkedRecordURL
	}

	conditionsResponse, err := c.FHIR.SearchFHIRCondition(ctx, params, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	conditions := []dto.Condition{}

	for _, resource := range conditionsResponse.Conditions {
		condition := mapFHIRConditionToConditionDTO(resource)
		conditions = append(conditions, *condition)
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     conditionsResponse.HasNextPage,
		EndCursor:       &conditionsResponse.NextCursor,
		HasPreviousPage: conditionsResponse.HasPreviousPage,
		StartCursor:     &conditionsResponse.PreviousCursor,
	}

	connection := dto.CreateConditionConnection(conditions, pageInfo, conditionsResponse.TotalCount)

	return &connection, nil
}
