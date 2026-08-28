package clinical

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/empower-clinical/pkg/clinical/usecases")

// CreateServiceRequest creates a service request resource in FHIR
func (c *ClinicalImpl) CreateServiceRequest(ctx context.Context, input *domain.FHIRServiceRequestInput) (*dto.ServiceRequest, error) {
	ctx, span := tracer.Start(ctx, "CreateServiceRequest")

	defer func() { span.End() }()

	serviceRequest, err := c.FHIR.CreateFHIRServiceRequest(ctx, *input)
	if err != nil {
		return nil, err
	}

	output := &dto.ServiceRequest{
		ID:       *serviceRequest.Resource.ID,
		Status:   dto.ServiceRequestStatus(serviceRequest.Resource.Status),
		Intent:   string(serviceRequest.Resource.Intent),
		Priority: string(serviceRequest.Resource.Priority),
		Subject: dto.Reference{
			ID:        *serviceRequest.Resource.Subject.ID,
			Reference: *serviceRequest.Resource.Subject.Reference,
			Display:   serviceRequest.Resource.Subject.Display,
		},
		Encounter: &dto.Reference{
			ID:        *serviceRequest.Resource.Encounter.ID,
			Reference: *serviceRequest.Resource.Encounter.Reference,
			Display:   serviceRequest.Resource.Encounter.Display,
		},
		ReceivingFacility: serviceRequest.Resource.GetPerformerID(),
	}

	if len(serviceRequest.Resource.Note) > 0 {
		output.Note = append(output.Note, dto.Annotation{
			Text: *serviceRequest.Resource.Note[0].Text,
		})
	}

	if serviceRequest.Resource.Text != nil {
		output.UsageContext = dto.ScreeningTypeEnum(serviceRequest.Resource.Text.Div)
	}

	if len(serviceRequest.Resource.Code.Concept.Coding) > 0 {
		output.OrderDetails.Name = serviceRequest.Resource.GetPatientReferralReason()
	}

	return output, nil
}

// CreateLaboratoryOrder is used to create a lab order when a patient is being referred to a lab for tests
func (c *ClinicalImpl) CreateLaboratoryOrder(ctx context.Context, serviceRequestID string) (*dto.ServiceRequest, error) {
	userSelected := false
	today := time.Now().Format(time.RFC3339)
	authoredOn := (*scalarutils.DateTime)(&today)
	codeSystem := helpers.CodeSystem("service-request-cs")

	serviceRequest, err := c.FHIR.GetFHIRServiceRequest(ctx, serviceRequestID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	var meta *domain.FHIRMetaInput
	if err := mapstructure.Decode(serviceRequest.Resource.Meta, &meta); err != nil {
		return nil, err
	}

	serviceRequestReference := fmt.Sprintf("ServiceRequest/%s", *serviceRequest.Resource.ID)

	requestingFacilityID, err := domain.GetFacilityIDFromResource(serviceRequest.Resource)
	if err != nil {
		return nil, err
	}

	facilityReference := fmt.Sprintf("Organization/%s", requestingFacilityID)

	patient, err := c.FHIR.GetFHIRPatient(ctx, *serviceRequest.Resource.Subject.ID)
	if err != nil {
		return nil, err
	}

	var code domain.FHIRCodeableReferenceInput

	err = mapstructure.Decode(serviceRequest.Resource.Code, &code)
	if err != nil {
		return nil, err
	}

	var performer []*domain.FHIRReferenceInput

	err = mapstructure.Decode(serviceRequest.Resource.Performer, &performer)
	if err != nil {
		return nil, err
	}

	labOrder := &domain.FHIRServiceRequestInput{
		AuthoredOn: authoredOn,
		Status:     domain.ServiceRequestStatusActive,
		Intent:     domain.ServiceRequestIntentInstanceOrder,
		Priority:   domain.ServiceRequestPriorityUrgent,
		Meta:       *meta,
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       codeSystem,
						Code:         scalarutils.Code(domain.LaboratoryProcedureCategoryType),
						Display:      domain.LaboratoryProcedureCategoryType.Display(),
						UserSelected: &userSelected,
					},
				},
				Text: domain.LaboratoryProcedureCategoryType.Display(),
			},
		},
		Code: &code,
		BasedOn: []*domain.FHIRReferenceInput{
			{
				ID:        &serviceRequestID,
				Display:   serviceRequestID,
				Reference: &serviceRequestReference,
			},
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        serviceRequest.Resource.Encounter.ID,
			Display:   serviceRequest.Resource.Encounter.Display,
			Reference: serviceRequest.Resource.Encounter.Reference,
		},
		Subject: &domain.FHIRReferenceInput{
			ID:        serviceRequest.Resource.Subject.ID,
			Reference: serviceRequest.Resource.Subject.Reference,
			Display:   serviceRequest.Resource.Subject.Display,
		},
		Requester: &domain.FHIRReferenceInput{
			ID:        &requestingFacilityID,
			Display:   requestingFacilityID,
			Reference: &facilityReference,
		},
		Performer: performer,
	}

	patientHealthIDIdentifier := patient.Resource.GetPatientHealthIDIdentifier()

	if patientHealthIDIdentifier != "" {
		labOrder.Subject.Identifier = &domain.FHIRIdentifierInput{
			Value: patientHealthIDIdentifier,
		}
	}

	if serviceRequest.Resource.Text != nil {
		labOrder.Text = &domain.FHIRNarrativeInput{
			Status: serviceRequest.Resource.Text.Status,
			Div:    serviceRequest.Resource.Text.Div,
		}
	}

	if len(serviceRequest.Resource.Reason) > 0 {
		var reason []*domain.FHIRCodeableReferenceInput

		err := mapstructure.Decode(serviceRequest.Resource.Reason, &reason)
		if err != nil {
			return nil, err
		}

		labOrder.Reason = reason
	}

	labOrderRequest, err := c.CreateServiceRequest(ctx, labOrder)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return labOrderRequest, nil
}

// ListLabOrders is used to retrieve a list of all lab orders
func (c *ClinicalImpl) ListLabOrders(ctx context.Context, filter *dto.ServiceRequestFilterInput, pagination dto.Pagination) (*dto.ServiceRequestOutputConnection, error) {
	labOrderSearchParams := map[string]any{
		"_sort":    "-_lastUpdated",
		"category": "laboratory-procedure",
		"_total":   "accurate",
	}

	if pagination.First != nil {
		labOrderSearchParams["_count"] = *pagination.First
	}

	err := filter.ServiceRequestFilters(labOrderSearchParams)
	if err != nil {
		return nil, err
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, err
	}

	if filter.PatientSearch != "" {
		references, err := c.SearchPatientReferences(ctx, filter.PatientSearch, *identifiers)
		if err != nil {
			return nil, err
		}

		if len(references) == 0 {
			emptyTotal := 0
			connection := dto.CreateServiceRequestConnection([]*dto.ServiceRequest{}, dto.PageInfo{}, &emptyTotal)

			return &connection, nil
		}

		labOrderSearchParams["subject"] = strings.Join(references, ",")
	}

	results, err := c.FHIR.SearchFHIRServiceRequest(ctx, labOrderSearchParams, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	serviceRequestList := make([]*dto.ServiceRequest, 0, len(results.ServiceRequests))

	for _, serviceRequest := range results.ServiceRequests {
		// A referral with no ordered tests carries only the "Reason for referral (narrative)"
		// coding. Such records are not actual lab orders, so exclude them from the listing.
		if hasOnlyReferralNarrativeCoding(serviceRequest) {
			continue
		}

		serviceRequestList = append(serviceRequestList, mapFHIRServiceRequestToServiceRequestDTO(serviceRequest, nil))
	}

	// Reconcile the total with anything the guard dropped so the count never over-reports the lab
	// orders actually returned. When the search is honoured server-side nothing is dropped and the
	// count is unchanged.
	totalCount := results.TotalCount
	if totalCount != nil {
		reconciled := max(*totalCount-(len(results.ServiceRequests)-len(serviceRequestList)), 0)
		totalCount = &reconciled
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     results.HasNextPage,
		EndCursor:       &results.NextCursor,
		HasPreviousPage: results.HasPreviousPage,
		StartCursor:     &results.PreviousCursor,
	}

	connection := dto.CreateServiceRequestConnection(serviceRequestList, pageInfo, totalCount)

	return &connection, nil
}

// GetLabOrder returns a lab order by ID. A lab order is represented as a FHIR service request
func (c *ClinicalImpl) GetLabOrder(ctx context.Context, serviceRequestID string) (*dto.ServiceRequest, error) {
	if serviceRequestID == "" {
		return nil, fmt.Errorf("servicerequestID is required")
	}

	servicerequest, err := c.FHIR.GetFHIRServiceRequest(ctx, serviceRequestID)
	if err != nil {
		return nil, err
	}

	var encounterID string

	if servicerequest.Resource.Encounter.ID != nil {
		encounterID = *servicerequest.Resource.Encounter.ID
	} else {
		encounterID = strings.TrimPrefix(*servicerequest.Resource.Encounter.Reference, "Encounter/")
	}

	requestID := *servicerequest.Resource.ID
	params := map[string]interface{}{
		"based-on":  requestID,
		"_sort":     "-_lastUpdated",
		"encounter": encounterID,
		"category":  "laboratory",
		"_total":    "accurate",
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, err
	}

	observation, err := c.FHIR.SearchFHIRObservation(ctx, params, *identifiers, dto.Pagination{Skip: true})
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	request := mapFHIRServiceRequestToServiceRequestDTO(*servicerequest.Resource, observation.Observations)

	return request, nil
}

func mapFHIRServiceRequestToServiceRequestDTO(serviceRequest domain.FHIRServiceRequest, observation []domain.FHIRObservation) *dto.ServiceRequest {
	authored := (*scalarutils.DateTime)(serviceRequest.AuthoredOn)

	output := &dto.ServiceRequest{
		Status:            dto.ServiceRequestStatus(serviceRequest.Status),
		Intent:            string(serviceRequest.Intent),
		Priority:          string(serviceRequest.Priority),
		ReceivingFacility: serviceRequest.GetFacilityName(),
		Category:          serviceRequest.GetCategory(),
		Date:              authored,
	}

	if serviceRequest.ID != nil {
		output.ID = *serviceRequest.ID
	}

	if serviceRequest.Subject != nil {
		output.Subject.Display = serviceRequest.Subject.Display

		if serviceRequest.Subject.ID != nil {
			output.Subject.ID = *serviceRequest.Subject.ID
		}

		if serviceRequest.Subject.Identifier != nil {
			output.Subject.Identifier = &dto.Identifier{
				Value: &serviceRequest.Subject.Identifier.Value,
			}
		}
	}

	if serviceRequest.Encounter != nil && serviceRequest.Encounter.ID != nil {
		output.Encounter = &dto.Reference{
			ID: *serviceRequest.Encounter.ID,
		}
	}

	if serviceRequest.Code != nil {
		var code, orderName string

		if testCode, testName, ok := getOrderedTestCoding(serviceRequest); ok {
			code, orderName = testCode, testName
		}

		output.OrderDetails = dto.OrderDetails{
			Name: orderName,
			Code: code,
		}
	}

	if serviceRequest.Text != nil {
		output.UsageContext = dto.ScreeningTypeEnum(helpers.ExtractTextFromHTML(string(serviceRequest.Text.Div)))
	}

	if len(observation) > 0 {
		for _, observation := range observation {
			obs := dto.Observation{
				ID:           *observation.ID,
				Status:       dto.ObservationStatus(*observation.Status),
				Value:        *observation.ValueString,
				Name:         output.OrderDetails.Name,
				PatientID:    *observation.Subject.ID,
				EncounterID:  *observation.Encounter.ID,
				TimeRecorded: string(*observation.EffectiveInstant),
			}

			if observation.Note != nil {
				obs.Note = string(*observation.Note[0].Text)
			}

			if observation.Text != nil {
				obs.UsageContext = dto.ScreeningTypeEnum(observation.Text.Div)
			}

			if observation.BasedOn != nil {
				basedOn := observation.GetBasedOn()
				obs.ServiceRequestID = *basedOn.ID
			}

			output.Results = append(output.Results, obs)
		}
	}

	return output
}

// getOrderedTestCoding returns the code and name of the actual ordered lab test from a lab
// order's Code block. Referral-derived orders seed the Code with the "Reason for referral
// (narrative)" coding (common.ReferralReasonLOINCode); that coding is skipped so the actual test
// is surfaced rather than the referral reason. The ordered test is then resolved from the first
// remaining coding:
//
//   - When recorded under the generic common.TestsOrderedLOINCCode marker, the test type is
//     carried in Display and resolved to its real LOINC code (defined in defaults).
//   - Otherwise the coding already carries the real LOINC code and display directly (self/intra
//     lab orders, and referrals that store the resolved test code).
//
// ok is false when no test coding is present (e.g. a referral with no ordered tests), in which
// case callers fall back to the default coding.
func getOrderedTestCoding(serviceRequest domain.FHIRServiceRequest) (code string, name string, ok bool) {
	if serviceRequest.Code == nil || serviceRequest.Code.Concept == nil {
		return "", "", false
	}

	for _, coding := range serviceRequest.Code.Concept.Coding {
		if coding == nil || coding.Code == nil {
			continue
		}

		codeValue := string(*coding.Code)

		// The referral-narrative coding records the reason for referral, not a test; skip it so
		// listings surface the ordered test rather than "Reason for referral (narrative)".
		if codeValue == common.ReferralReasonLOINCode {
			continue
		}

		// Some referral flows record the ordered test under the generic "tests ordered" marker,
		// carrying the test type in Display. Resolve it to the real LOINC code so the frontend
		// receives the concrete test code rather than the marker. Display may hold either the raw
		// enum ("IHC_HER2") or its humanised label ("IHC HER2"); LabTestTypeFromString handles both.
		if codeValue == common.TestsOrderedLOINCCode {
			testType := dto.LabTestTypeFromString(coding.Display)
			if testCode := testType.Code(); testCode != "" {
				return testCode, testType.Display(), true
			}

			// Unknown test type with no mapped code; surface the display verbatim.
			return "", coding.Display, true
		}

		// Any other coding already carries the real LOINC test code and its display directly.
		return codeValue, coding.Display, true
	}

	return "", "", false
}

// hasOnlyReferralNarrativeCoding reports whether a service request's Code carries nothing but the
// "Reason for referral (narrative)" coding (common.ReferralReasonLOINCode). Referrals seed their
// Code with this narrative coding and append a coding per ordered test; a referral with no tests
// therefore ends up with only the narrative coding. When such a referral is turned into a lab
// order it has no actual test to display, so listings use this to filter it out. It returns false
// as soon as any non-narrative (i.e. real test) coding is present, keeping genuine lab orders.
func hasOnlyReferralNarrativeCoding(serviceRequest domain.FHIRServiceRequest) bool {
	if serviceRequest.Code == nil || serviceRequest.Code.Concept == nil {
		return false
	}

	codings := serviceRequest.Code.Concept.Coding
	if len(codings) == 0 {
		return false
	}

	sawNarrative := false

	for _, coding := range codings {
		if coding == nil || coding.Code == nil {
			continue
		}

		if string(*coding.Code) == common.ReferralReasonLOINCode {
			sawNarrative = true
			continue
		}

		// A non-narrative coding represents an actual ordered test: keep the record.
		return false
	}

	return sawNarrative
}

// CreateTestOrder creates a new diagnostic test order based on the provided input.
// This function is responsible for generating a FHIR ServiceRequest resource
// to represent the requested test and linking it to the appropriate clinical context.
func (c *ClinicalImpl) CreateTestOrder(ctx context.Context, input *dto.TestOrder) (*dto.ServiceRequest, error) {
	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record a referral in a completed encounter")
	}

	patientReference := fmt.Sprintf("Patient/%s", *encounter.Resource.Subject.ID)
	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	startTime := scalarutils.DateTime(time.Now().Format("2006-01-02T15:04:05+03:00"))
	requester := fmt.Sprintf("Organization/%s", input.Facility.FHIROrganisationID)
	system := helpers.CodeSystem("service-request-cs")

	testOrderCode, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, input.LoincCode)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	serviceRequest := domain.FHIRServiceRequestInput{
		Status:   domain.ServiceRequestStatusActive,
		Intent:   domain.ServiceRequestIntentOrder,
		Priority: domain.ServiceRequestPriorityRoutine,
		Code: &domain.FHIRCodeableReferenceInput{
			Concept: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						Code:    scalarutils.Code(testOrderCode.ID),
						Display: testOrderCode.GetConceptDisplay(),
					},
				},
				Text: testOrderCode.GetConceptDisplay(),
			},
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  system,
						Code:    scalarutils.Code(domain.ReferralCategoryType),
						Display: domain.ReferralCategoryType.Display(),
					},
				},
				Text: domain.ReferralCategoryType.Display(),
			},
		},
		AuthoredOn: &startTime,
		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: &patientReference,
			Display:   encounter.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Reference: &encounterReference,
		},
		Requester: &domain.FHIRReferenceInput{
			ID:        &input.Facility.FHIROrganisationID,
			Reference: &requester,
		},
		Note: []*domain.FHIRAnnotationInput{
			{
				Time: &startTime,
				Text: (*scalarutils.Markdown)(&input.ClinicalNote),
			},
		},
	}

	if input.Diagnosis != "" {
		condition, err := c.FHIR.GetFHIRCondition(ctx, input.Diagnosis)
		if err != nil {
			return nil, err
		}

		if condition.Resource != nil {
			conditionCode := condition.Resource.Code.ConditionCode()
			serviceRequest.Reason = append(serviceRequest.Reason, conditionCode)
		}
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	serviceRequest.Meta.Tag = tags

	output, err := c.CreateServiceRequest(ctx, &serviceRequest)
	if err != nil {
		return nil, err
	}

	// TODO: Create some FHIR Task?
	return output, nil
}
