package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

const timeFormatStr = "2006-01-02T15:04:05+03:00"

// resolveEffectiveTime returns the supplied date (when the test was performed)
// if it is provided, otherwise it falls back to the current time.
func resolveEffectiveTime(date *scalarutils.Date) time.Time {
	if date != nil {
		return date.AsTime()
	}

	return time.Now()
}

var (
	diagnosticReportCategoryCodeSystem = "http://terminology.hl7.org/CodeSystem/v2-0074"
)

// RecordMammographyResult is used to record mammography diagnostic reports as specified in https://hl7.org/fhir/R4/diagnosticreport.html#scope
func (c *ClinicalImpl) RecordMammographyResult(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  input.EncounterID,
			Value:        input.Findings,
			UsageContext: input.UsageContext,
		},
		VitalSignsConceptID: common.MammogramTerminologyCode,
	}

	// !NOTE: The terminology code used here is used TEMPORARILY. Pending discussion about how to represent BI-RADs conclusions/observation
	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("imaging")})
	if err != nil {
		return nil, err
	}

	if err := c.handleSelfLabOrder(ctx, input, observationOutput); err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, common.MammogramTerminologyCode, input, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory("NMR")})
}

// handleSelfLabOrder is a helper method handling common logic used to build intra-referral lab-order payload
func (c *ClinicalImpl) handleSelfLabOrder(ctx context.Context, input dto.DiagnosticReportInput, observationOutput *dto.Observation) error {
	if input.UsageContext.IsValid() && input.Findings == "" {
		payload := dto.IntraLabOrderInput{
			Code:        observationOutput.Code,
			Name:        observationOutput.Name,
			EncounterID: observationOutput.EncounterID,
			Patient: dto.Patient{
				Name: observationOutput.PatientName,
				ID:   observationOutput.PatientID,
			},
			ObservationID: observationOutput.ID,
			UsageContext:  input.UsageContext,
			Date:          input.Date,
		}

		_, err := c.CreateLabOrder(ctx, &payload)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateLabOrder is a helper function to create lab order in the case of internal referrals where results are EMPTY. This allows us to have a list of lab orders where test
// results can be uploaded against
func (c *ClinicalImpl) CreateLabOrder(ctx context.Context, labOrderInput *dto.IntraLabOrderInput) (*dto.ServiceRequest, error) {
	orderDate := resolveEffectiveTime(labOrderInput.Date).Format(time.RFC3339)
	authoredOn := (*scalarutils.DateTime)(&orderDate)
	userSelected := false
	encounterRef := fmt.Sprintf("Encounter/%s", labOrderInput.EncounterID)

	patient, err := c.FHIR.GetFHIRPatient(ctx, labOrderInput.Patient.ID)
	if err != nil {
		return nil, err
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.Resource.ID)

	var (
		serviceRequestCategory = helpers.CodeSystem("service-request-cs")
		observationReference   = fmt.Sprintf("Observation/%s", labOrderInput.ObservationID)
	)

	organizationID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	organizationRef := fmt.Sprintf("Organization/%s", organizationID)

	labOrder := &domain.FHIRServiceRequestInput{
		AuthoredOn: authoredOn,
		Status:     domain.ServiceRequestStatusActive,
		Intent:     domain.ServiceRequestIntentInstanceOrder,
		Priority:   domain.ServiceRequestPriorityUrgent,
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       serviceRequestCategory,
						Code:         scalarutils.Code(domain.LaboratoryProcedureCategoryType),
						Display:      domain.LaboratoryProcedureCategoryType.Display(),
						UserSelected: &userSelected,
					},
				},
				Text: domain.LaboratoryProcedureCategoryType.Display(),
			},
		},
		Code: &domain.FHIRCodeableReferenceInput{
			Concept: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						Code:         scalarutils.Code(labOrderInput.Code),
						Display:      labOrderInput.Name,
						System:       (*scalarutils.URI)(&common.LoincSystemURL),
						UserSelected: &userSelected,
					},
				},
				Text: fmt.Sprintf("%s test(s)", labOrderInput.Name),
			},
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &labOrderInput.EncounterID,
			Reference: &encounterRef,
			Display:   encounterRef,
		},
		Subject: &domain.FHIRReferenceInput{
			ID:        patient.Resource.ID,
			Display:   patient.Resource.Names(),
			Reference: &patientRef,
		},
		Requester: &domain.FHIRReferenceInput{
			Reference: &organizationRef,
			Display:   organizationID,
		},
		Performer: []*domain.FHIRReferenceInput{
			{
				Display:   organizationID,
				Reference: &organizationRef,
			},
		},
		Reason: []*domain.FHIRCodeableReferenceInput{
			{
				Reference: &domain.FHIRReferenceInput{
					Display:   labOrderInput.ObservationID,
					Reference: &observationReference,
				},
			},
		},
	}

	patientHealthIDIdentifier := patient.Resource.GetPatientHealthIDIdentifier()

	if patientHealthIDIdentifier != "" {
		labOrder.Subject.Identifier = &domain.FHIRIdentifierInput{
			Value: patientHealthIDIdentifier,
		}
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	labOrder.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	if labOrder.Text == nil {
		generatedTextStatus := "generated"
		div := fmt.Sprintf("<div xmlns=\"http://www.w3.org/1999/xhtml\">%s</div>", labOrderInput.UsageContext)
		labOrder.Text = &domain.FHIRNarrativeInput{
			Status: (*domain.NarrativeStatusEnum)(&generatedTextStatus),
			Div:    scalarutils.XHTML(div),
		}
	}

	labOrderOutput, err := c.CreateServiceRequest(ctx, labOrder)
	if err != nil {
		return nil, err
	}

	var meta *dto.MetaInput

	err = mapstructure.Decode(labOrder.Meta, &meta)
	if err != nil {
		return nil, err
	}

	payload := &dto.PatientReferralTaskPayload{
		Meta:     meta,
		Referral: labOrderOutput,
	}

	err = c.Pubsub.NotifyCreatePatientReferralTask(ctx, *payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return labOrderOutput, nil
}

// RecordBiopsy is used to record biopsy test results as a diagnostic report
// FHIR recommends use of diagnostic resource to record the findings and interpretation of biopsy test results
// performed on patients, groups of patients, devices, and locations, and/or specimens.
func (c *ClinicalImpl) RecordBiopsy(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  input.EncounterID,
			Value:        input.Findings,
			UsageContext: input.UsageContext,
		},
		VitalSignsConceptID: common.BiopsyTerminologySystem,
	}

	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("procedure")})
	if err != nil {
		return nil, err
	}

	if err := c.handleSelfLabOrder(ctx, input, observationOutput); err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, common.BiopsyTerminologySystem, input, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory("CP")})
}

// RecordMRI is used to record MRI scan results as a diagnostic report
func (c *ClinicalImpl) RecordMRI(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  input.EncounterID,
			Value:        input.Findings,
			UsageContext: input.UsageContext,
		},
		VitalSignsConceptID: common.MRITerminologySystem,
	}

	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("imaging")})
	if err != nil {
		return nil, err
	}

	if err := c.handleSelfLabOrder(ctx, input, observationOutput); err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, common.MRITerminologySystem, input, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory("NMR")})
}

// RecordUltrasound is used to record the breast ultrasound diagnostic reports
func (c *ClinicalImpl) RecordUltrasound(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  input.EncounterID,
			Value:        input.Findings,
			UsageContext: input.UsageContext,
		},
		VitalSignsConceptID: common.ChestUltrasoundTerminologySystem,
	}

	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("imaging")})
	if err != nil {
		return nil, err
	}

	if err := c.handleSelfLabOrder(ctx, input, observationOutput); err != nil {
		return nil, err
	}

	// TODO: `BilateralConceptTerminologySystem` is a `PLACE HOLDER`. It should be adjusted accordingly when the breast cancer designs are revamped
	return c.RecordDiagnosticReport(ctx, common.ChestUltrasoundTerminologySystem, input, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory("RUS")})
}

// RecordCBE is used to record clinical based examination test results for a patient
func (c *ClinicalImpl) RecordCBE(ctx context.Context, input *dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  input.EncounterID,
			Value:        input.Findings,
			UsageContext: input.UsageContext,
		},
		VitalSignsConceptID: common.BreastExaminationLOINCTerminologySystem,
	}

	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("exam")})
	if err != nil {
		return nil, err
	}

	if err := c.handleSelfLabOrder(ctx, *input, observationOutput); err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, common.BreastExaminationLOINCTerminologySystem, *input, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory("OTH")})
}

// RecordPapSmear is used to record Pap Smear test results.
func (c *ClinicalImpl) RecordPapSmear(ctx context.Context, input *dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  input.EncounterID,
			Value:        input.Findings,
			UsageContext: input.UsageContext,
		},
		VitalSignsConceptID: common.PapSmearTerminologyCode,
	}

	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("laboratory")})
	if err != nil {
		return nil, err
	}

	if err := c.handleSelfLabOrder(ctx, *input, observationOutput); err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, common.PapSmearTerminologyCode, *input, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory("OTH")})
}

// RecordPSA is used to record prostate cancer test results.
func (c *ClinicalImpl) RecordPSA(ctx context.Context, payload *dto.PSAInput) (*dto.DiagnosticReport, error) {
	if err := helpers.Validate(payload.DiagnosticInput); err != nil {
		return nil, err
	}

	var terminologyCode string

	switch payload.PSAType {
	case dto.ProstaticSerumAntigen:
		terminologyCode = common.ProstateCancerTerminologyCode
	case dto.WholeBlood:
		terminologyCode = common.WholeBloodTerminologyCode
	default:
		return nil, fmt.Errorf("invalid PSA type: %s", payload.PSAType)
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  payload.DiagnosticInput.EncounterID,
			Value:        payload.DiagnosticInput.Findings,
			UsageContext: payload.DiagnosticInput.UsageContext,
		},
		VitalSignsConceptID: terminologyCode,
	}

	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("laboratory")})
	if err != nil {
		return nil, fmt.Errorf("failed to record observation for PSA type %s: %w", payload.PSAType, err)
	}

	if err := c.handleSelfLabOrder(ctx, *payload.DiagnosticInput, observationOutput); err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, terminologyCode, *payload.DiagnosticInput, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory("OTH")})
}

// CreateLabOrderResult creates a result of the test ordered.
func (c *ClinicalImpl) CreateLabOrderResult(ctx context.Context, input *dto.TestOrderResult) (*dto.DiagnosticReport, error) {
	if input.ServiceRequestID == "" {
		return nil, fmt.Errorf("service request ID is required")
	}

	serviceRequest, err := c.FHIR.GetFHIRServiceRequest(ctx, input.ServiceRequestID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	encounterID := serviceRequest.Resource.Encounter.Display
	encounterReference := fmt.Sprintf("Encounter/%s", serviceRequest.Resource.Encounter.Display)
	patientReference := fmt.Sprintf("Patient/%s", serviceRequest.Resource.Subject.Display)
	serviceRequestReference := fmt.Sprintf("ServiceRequest/%s", *serviceRequest.Resource.ID)

	// The LOINC code of the test has been saved in the ServiceRequest Code.Coding field
	labOrderLOINCCode := serviceRequest.Resource.GetCode()

	var observationReferences []*domain.FHIRReferenceInput

	// We need to create an observation - to store the results of the test, the observation should be linked to the
	// service request..
	for _, test := range input.Test {
		testResult := &dto.ObservationPayload{
			ObservationInput: dto.ObservationInput{
				Status:      dto.ObservationStatusFinal,
				EncounterID: encounterID,
				Note:        test.Finding,
				Value:       test.Value,
			},
			VitalSignsConceptID: labOrderLOINCCode,
		}

		observation, err := c.RecordObservation(ctx, *testResult, []ObservationInputMutatorFunc{addObservationCategory("laboratory")})
		if err != nil {
			return nil, err
		}

		observationsReference := fmt.Sprintf("Observation/%s", observation.ID)
		observationReferences = append(observationReferences, &domain.FHIRReferenceInput{
			ID:        &observation.ID,
			Reference: &observationsReference,
			Display:   observation.ID,
		})
	}

	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))
	dateTime := scalarutils.DateTime(time.Now().Format(timeFormatStr))

	diagnosticReport := domain.FHIRDiagnosticReportInput{
		BasedOn: []*domain.FHIRReferenceInput{
			{
				ID:        &input.ServiceRequestID,
				Reference: &serviceRequestReference,
				Display:   input.ServiceRequestID,
			},
		},
		Status: domain.DiagnosticReportStatusFinal,
		Subject: &domain.FHIRReferenceInput{
			ID:        &serviceRequest.Resource.Subject.Display,
			Reference: &patientReference,
			Display:   serviceRequest.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &encounterID,
			Reference: &encounterReference,
			Display:   encounterID,
		},
		Result: observationReferences,
		Issued: (*string)(&instant),

		EffectiveDateTime: &dateTime,
	}

	err = mapstructure.Decode(serviceRequest.Resource.Performer, &diagnosticReport.Performer)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(serviceRequest.Resource.Performer, &diagnosticReport.ResultsInterpreter)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(serviceRequest.Resource.Category, &diagnosticReport.Category)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(serviceRequest.Resource.Code.Concept, &diagnosticReport.Code)
	if err != nil {
		return nil, err
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	diagnosticReport.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	report, err := c.FHIR.CreateFHIRDiagnosticReport(ctx, &diagnosticReport)
	if err != nil {
		return nil, err
	}

	return &dto.DiagnosticReport{
		ID:          *report.ID,
		Status:      dto.ObservationStatus(report.Status),
		PatientID:   *report.Subject.ID,
		EncounterID: *report.Encounter.ID,
	}, nil
}

// RecordDiagnosticReport is a re-usable method to help with diagnostic report recording
func (c *ClinicalImpl) RecordDiagnosticReport(ctx context.Context, conceptID string,
	input dto.DiagnosticReportInput, observation *dto.Observation, mutators []DiagnosticReportMutatorFunc) (*dto.DiagnosticReport, error) {
	observationsReference := fmt.Sprintf("Observation/%s", observation.ID)
	observationType := scalarutils.URI("Observation")
	encounterReference := fmt.Sprintf("Encounter/%s", observation.EncounterID)
	encounterType := scalarutils.URI("Encounter")
	patientReference := fmt.Sprintf("Patient/%s", observation.PatientID)
	patientType := scalarutils.URI("Patient")

	// Use the supplied date when present, otherwise fall back to the current time.
	dateTime := scalarutils.DateTime(resolveEffectiveTime(input.Date).Format(timeFormatStr))

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	facility, err := c.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	orgRef := fmt.Sprintf("Organization/%s", *facility.Resource.ID)
	orgType := scalarutils.URI("Organization")

	id := uuid.New().String()
	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	diagnosticReport := &domain.FHIRDiagnosticReportInput{
		ID:     &id,
		Status: domain.DiagnosticReportStatusFinal,
		Subject: &domain.FHIRReferenceInput{
			ID:        &observation.PatientID,
			Reference: &patientReference,
			Type:      &patientType,
			Display:   observation.PatientID,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &observation.EncounterID,
			Reference: &encounterReference,
			Type:      &encounterType,
			Display:   observation.EncounterID,
		},
		Issued: (*string)(&instant),
		Performer: []*domain.FHIRReferenceInput{
			{
				ID:        facility.Resource.ID,
				Reference: &orgRef,
				Type:      &orgType,
				Display:   *facility.Resource.ID,
			},
		},
		ResultsInterpreter: []*domain.FHIRReferenceInput{
			{
				Reference: &orgRef,
				Type:      &orgType,
				Display:   *facility.Resource.ID,
			},
		},
		Result: []*domain.FHIRReferenceInput{
			{
				ID:        &observation.ID,
				Reference: &observationsReference,
				Type:      &observationType,
				Display:   observation.ID,
			},
		},
		EffectiveDateTime: (*scalarutils.DateTime)(&dateTime),
	}

	if conceptID != "" {
		concept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, conceptID)
		if err != nil {
			return nil, err
		}

		code := domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&common.LoincSystemURL),
					Code:    scalarutils.Code(concept.ID),
					Display: concept.GetConceptDisplay(),
				},
			},
			Text: concept.GetConceptDisplay(),
		}

		diagnosticReport.Code = code
	}

	diagnosticReport.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	if input.Note != "" {
		diagnosticReport.Conclusion = &input.Note
	}

	if input.UsageContext != "" {
		diagnosticReportTextStatus := domain.NarrativeStatusEnumAdditional.String()
		text := (dto.ScreeningTypeEnum)(input.UsageContext).String()
		diagnosticReport.Text = utils.NarrativeGenerator(text, &diagnosticReportTextStatus)
	}

	if len(input.Media) > 0 {
		for _, media := range input.Media {
			mediaReference := fmt.Sprintf("DocumentReference/%s", media.ID)

			attachment := []*domain.FHIRAttachmentInput{
				{
					ID:    &media.ID,
					URL:   (*scalarutils.URL)(&media.MediaLink),
					Title: &mediaReference,
				},
			}

			diagnosticReport.PresentedForm = append(diagnosticReport.PresentedForm, attachment...)
		}
	}

	if len(mutators) > 0 {
		for _, mutator := range mutators {
			err = mutator(ctx, diagnosticReport)
			if err != nil {
				return nil, err
			}
		}
	}

	result, err := c.FHIR.CreateFHIRDiagnosticReport(ctx, diagnosticReport)
	if err != nil {
		return nil, err
	}

	output := &dto.DiagnosticReport{
		ID:          *result.ID,
		Status:      dto.ObservationStatus(result.Status),
		PatientID:   *result.Subject.ID,
		EncounterID: *result.Encounter.ID,
		Issued:      *result.Issued,
	}

	if result.Conclusion != nil {
		output.Conclusion = *result.Conclusion
	}

	return output, nil
}

// RecordTestResult is used to record the test result of an ordered test
func (c *ClinicalImpl) RecordTestResult(ctx context.Context, input dto.TestResultInput) (*dto.DiagnosticReport, error) {
	if err := helpers.Validate(input.Entry); err != nil {
		return nil, err
	}

	if input.Entry.Findings == "" {
		return nil, fmt.Errorf("test findings must be specified")
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  input.Entry.EncounterID,
			Value:        input.Entry.Findings,
			UsageContext: input.Entry.UsageContext,
			Note:         input.Entry.Note,
			Date:         input.Entry.Date,
		},
		VitalSignsConceptID: common.LabResultsInterpretationLOINCTerminologyCode,
		ServiceRequestID:    input.ServiceRequestID,
	}

	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("laboratory")})
	if err != nil {
		return nil, err
	}

	diagnosticReportOutput, err := c.RecordDiagnosticReport(ctx, common.LabResultsInterpretationLOINCTerminologyCode, input.Entry, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory("OTH")})
	if err != nil {
		return nil, err
	}

	err = c.TestReviewFollowUpTask(ctx, diagnosticReportOutput, observation)
	if err != nil {
		return nil, err
	}

	return diagnosticReportOutput, nil
}

// TestReviewFollowUpTask is used to create a follow-up task for a healthcare worker (HCW) to reach out to a patient and notify them to return for a review of their test results.
func (c *ClinicalImpl) TestReviewFollowUpTask(ctx context.Context, diagnosticReportOutput *dto.DiagnosticReport, observation *dto.ObservationPayload) error {
	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return err
	}

	serviceRequest, err := c.FHIR.GetFHIRServiceRequest(ctx, observation.ServiceRequestID)
	if err != nil {
		return err
	}

	meta := &domain.FHIRMetaInput{
		Tag: tags,
	}

	var (
		orderIntent             = "plan"
		priority                = "routine"
		requestedStatus         = dto.RequestedTasksStatus.String()
		patientRef              = fmt.Sprintf("Patient/%s", diagnosticReportOutput.PatientID)
		diagnosticReportRef     = fmt.Sprintf("DiagnosticReport/%s", diagnosticReportOutput.ID)
		encounterReference      = fmt.Sprintf("Encounter/%s", observation.ObservationInput.EncounterID)
		requesterReference      = fmt.Sprintf("Organization/%s", meta.GetOrganization())
		serviceRequestReference = fmt.Sprintf("ServiceRequest/%s", observation.ServiceRequestID)
	)

	now := time.Now()
	authoredOn := scalarutils.DateTime(now.Format(time.RFC3339))
	defaultPeriod := scalarutils.DateTime(now.Add(5 * 24 * time.Hour).Format(time.RFC3339))

	var code *domain.FHIRCodeableConceptInput

	err = mapstructure.Decode(serviceRequest.Resource.Code.Concept, &code)
	if err != nil {
		return err
	}

	followUpTask := &domain.FHIRTaskInput{
		Meta:   meta,
		Status: (*scalarutils.Code)(&requestedStatus),
		Reason: []*domain.FHIRCodeableReference{},
		BusinessStatus: &domain.FHIRCodeableConceptInput{
			Text: "Test results review",
		},
		Code:     code,
		Intent:   (*scalarutils.Code)(&orderIntent),
		Priority: (*scalarutils.Code)(&priority),
		BasedOn: []*domain.FHIRReferenceInput{
			{
				ID:        &observation.ServiceRequestID,
				Reference: &serviceRequestReference,
				Display:   observation.ServiceRequestID,
			},
		},
		For: &domain.FHIRReferenceInput{
			ID:        &diagnosticReportOutput.PatientID,
			Reference: &patientRef,
			Display:   serviceRequest.Resource.Subject.Display,
		},
		Focus: &domain.FHIRReferenceInput{
			ID:        &diagnosticReportOutput.ID,
			Reference: &diagnosticReportRef,
			Display:   diagnosticReportOutput.ID,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &observation.ObservationInput.EncounterID,
			Reference: &encounterReference,
			Display:   observation.ObservationInput.EncounterID,
		},
		AuthoredOn: &authoredOn,
		Requester: &domain.FHIRReferenceInput{
			Reference: &requesterReference,
			Display:   meta.GetOrganization(),
		},
		Owner: &domain.FHIRReferenceInput{
			Reference: &requesterReference,
			Display:   meta.GetOrganization(),
		},
		RequestedPeriod: &domain.FHIRPeriodInput{
			Start: &authoredOn,
			End:   &defaultPeriod,
		},
		RequestedPerformer: []*domain.FHIRCodeableReference{
			{
				Reference: &domain.FHIRReference{
					Display:   meta.GetOrganization(),
					Reference: &requesterReference,
				},
			},
		},
	}

	if observation.ObservationInput.UsageContext.IsValid() {
		status := "additional"
		divContent := fmt.Sprintf("<div xmlns=\"http://www.w3.org/1999/xhtml\">%s</div>", observation.ObservationInput.UsageContext)

		followUpTask.Text = &domain.FHIRNarrative{
			Status: (*domain.NarrativeStatusEnum)(&status),
			Div:    scalarutils.XHTML(divContent),
		}
	}

	err = c.Pubsub.NotifyCreateFollowUpTask(ctx, followUpTask)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return err
	}

	_, err = c.FHIR.PatchFHIRServiceRequest(ctx, observation.ServiceRequestID, domain.FHIRServiceRequestInput{
		Status: domain.ServiceRequestStatusCompleted,
	})
	if err != nil {
		return err
	}

	return nil
}

// RecordLabTests is used to record lab test
func (c *ClinicalImpl) RecordTests(ctx context.Context, payload dto.TestInput) (*dto.DiagnosticReport, error) {
	if err := helpers.Validate(payload); err != nil {
		return nil, err
	}

	observation := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  payload.Input.EncounterID,
			Value:        payload.Input.Findings,
			UsageContext: payload.Input.UsageContext,
			Note:         payload.Input.Note,
		},
	}

	obsCategory, diagnosticReportCategory, conceptID := testTypeCategory(payload.TestType)

	observation.VitalSignsConceptID = conceptID

	observationOutput, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory(obsCategory), addObservationEffectiveDate(payload.Input.Date)})
	if err != nil {
		return nil, err
	}

	if err := c.handleSelfLabOrder(ctx, payload.Input, observationOutput); err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, conceptID, payload.Input, observationOutput, []DiagnosticReportMutatorFunc{addDiagnosticReportCategory(diagnosticReportCategory)})
}

func (c *ClinicalImpl) DeleteTest(ctx context.Context, observationID string) (bool, error) {
	searchParams := map[string]any{
		"result": fmt.Sprintf("%s/%s", "Observation", observationID),
		"_total": "accurate",
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	diagnosticReportResults, err := c.FHIR.SearchFHIRDiagnosticReport(ctx, searchParams, *identifiers, dto.Pagination{})
	if err != nil {
		return false, err
	}

	if len(diagnosticReportResults.DiagnosticReport) > 0 {
		for _, diagnosticReport := range diagnosticReportResults.DiagnosticReport {
			_, err := c.FHIR.DeleteFHIRResource(ctx, "DiagnosticReport", *diagnosticReport.ID)
			if err != nil {
				return false, err
			}
		}
	}

	_, err = c.FHIR.DeleteFHIRResource(ctx, "Observation", observationID)
	if err != nil {
		return false, err
	}

	return true, nil
}
