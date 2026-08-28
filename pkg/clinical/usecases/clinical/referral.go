package clinical

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	dto "github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// resolveReferralDate determines the effective referral timestamp from the optional referral date
// on the input. When no date is supplied the current time is used, preserving the previous
// behaviour. When a date is supplied it may not be later than today — a referral cannot be
// future-dated — and is returned as the start of that day.
func resolveReferralDate(date *scalarutils.Date) (time.Time, error) {
	if date == nil {
		return time.Now(), nil
	}

	referralDate := date.AsTime()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if referralDate.After(today) {
		return time.Time{}, fmt.Errorf("referral date cannot be in the future")
	}

	return referralDate, nil
}

// ReferPatient creates a new patient referral based on the specified inputs.
// This method handles the logic for generating a referral for different purposes,
// such as diagnostics and testing, specialist consultation, or treatment.
// The method constructs a ServiceRequest object following the FHIR R4 standards
func (c *ClinicalImpl) ReferPatient(
	ctx context.Context,
	input *dto.ReferralInput,
) (*dto.ServiceRequest, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	referralInstant, err := resolveReferralDate(input.ReferralDate)
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record a referral in a completed encounter")
	}

	patientReference := fmt.Sprintf("Patient/%s", *encounter.Resource.Subject.ID)
	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	startTime := scalarutils.DateTime(referralInstant.Format("2006-01-02T15:04:05+03:00"))

	var performerReference string

	if input.Facility.FHIROrganisationID == "" {
		return nil, fmt.Errorf("please provide the facility receiving the referral")
	} else {
		performerReference = fmt.Sprintf("Organization/%s", input.Facility.FHIROrganisationID)
	}

	serviceRequestSystem := helpers.CodeSystem("service-request-cs")

	serviceRequest := domain.FHIRServiceRequestInput{
		Status:   domain.ServiceRequestStatusActive,
		Intent:   domain.ServiceRequestIntentOrder,
		Priority: domain.ServiceRequestPriorityUrgent,
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  serviceRequestSystem,
						Code:    scalarutils.Code(domain.ReferralCategoryType),
						Display: domain.ReferralCategoryType.Display(),
					},
				},
				Text: domain.ReferralCategoryType.Display(),
			},
		},
		Code: &domain.FHIRCodeableReferenceInput{
			Concept: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  &common.LoincSystemURL,
						Code:    common.ReferralReasonLOINCode,
						Display: "Reason for referral (narrative)",
					},
				},
			},
		},
		AuthoredOn: &startTime,
		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: &patientReference,
			Display:   encounter.Resource.Subject.Display,
		},
		Performer: []*domain.FHIRReferenceInput{
			{
				Display:   input.Facility.FHIROrganisationID,
				Reference: &performerReference,
			},
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Display:   *encounter.Resource.ID,
			Reference: &encounterReference,
		},
		// TODO: This should be the actual diagnosis code e.g Breast Cancer Screening etc
		// Papercut: Assign Adrian
		Reason: []*domain.FHIRCodeableReferenceInput{
			{
				Concept: &domain.FHIRCodeableConceptInput{
					Coding: []*domain.FHIRCodingInput{
						{
							System:  &common.LoincSystemURL,
							Code:    common.ReferralLOINCTerminologySystem,
							Display: "Referral note",
						},
					},
					Text: "Referral note",
				},
			},
		},
	}

	if input.ReferralNote != "" {
		serviceRequest.Note = []*domain.FHIRAnnotationInput{
			{
				Time: &startTime,
				Text: (*scalarutils.Markdown)(&input.ReferralNote),
			},
		}
	}

	if len(input.Tests) > 0 {
		var validTests []string

		for _, test := range input.Tests {
			if strings.TrimSpace(test) == "" {
				continue
			}

			testName := dto.LabTestTypeEnum(test)

			referredTest := domain.FHIRCodingInput{
				System:  &common.LoincSystemURL,
				Code:    scalarutils.Code(testName.Code()),
				Display: testName.Display(),
			}

			serviceRequest.Code.Concept.Coding = append(serviceRequest.Code.Concept.Coding, &referredTest)

			validTests = append(validTests, testName.Display())
		}

		if len(validTests) > 0 {
			serviceRequest.Code.Concept.Text = fmt.Sprintf("%s test(s)", strings.Join(validTests, ", "))
		}
	}

	switch input.ReferralType {
	case dto.SPECIALIST:
		serviceRequest.Code.Concept.Text = fmt.Sprintf("Referred to %s specialist for further evaluation", input.Specialist)
	case dto.TREATMENT:
		serviceRequest.Code.Concept.Text = "Referred for treatment"
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	requesterRef := fmt.Sprintf("Organization/%s", identifiers.OrganizationID)

	serviceRequest.Requester = &domain.FHIRReferenceInput{
		Display:   identifiers.OrganizationID,
		Reference: &requesterRef,
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

	var meta *dto.MetaInput

	err = mapstructure.Decode(serviceRequest.Meta, &meta)
	if err != nil {
		return nil, err
	}

	payload := &dto.PatientReferralTaskPayload{
		Meta:     meta,
		Referral: output,
	}

	err = c.Pubsub.NotifyCreatePatientReferralTask(ctx, *payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return output, nil
}

// GetPatientReferrals lists all the referrals for a patient
func (c *ClinicalImpl) GetPatientReferrals(ctx context.Context, searchInput *dto.ReferralSearchInput) (*dto.ReferralDetailConnection, error) {
	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	serviceRequestSearchParams := map[string]interface{}{
		"_sort":       "-_lastUpdated",
		"category":    string(domain.ReferralCategoryType),
		"_revinclude": "DocumentReference:based-on",
		"_total":      "accurate",
		"status":      searchInput.SetStatus(),
	}

	if searchInput.PatientID != nil {
		patientReference := fmt.Sprintf("Patient/%s", *searchInput.PatientID)
		serviceRequestSearchParams["subject"] = patientReference
	}

	if searchInput.EncounterID != nil {
		encounterReference := fmt.Sprintf("Encounter/%s", *searchInput.EncounterID)
		serviceRequestSearchParams["encounter"] = encounterReference
	}

	if searchInput.Date != nil {
		serviceRequestSearchParams["authored"] = searchInput.Date.AsTime().Format(time.DateOnly)
	}

	if searchInput.PatientSearch != "" {
		references, err := c.SearchPatientReferences(ctx, searchInput.PatientSearch, *identifiers)
		if err != nil {
			return nil, err
		}

		if len(references) == 0 {
			return &dto.ReferralDetailConnection{}, nil
		}

		serviceRequestSearchParams["subject"] = strings.Join(references, ",")
	}

	pagedResource, err := c.FHIR.SearchFHIRResource(ctx, "", "ServiceRequest", serviceRequestSearchParams, *identifiers, *searchInput.Pagination)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	var connection dto.ReferralDetailConnection

	var serviceRequests []*domain.FHIRServiceRequest

	var documentReferences []*domain.FHIRDocumentReference

	for _, resource := range pagedResource.Resources {
		if resourceType, ok := resource["resourceType"].(string); ok {
			switch resourceType {
			case "ServiceRequest":
				serviceRequest, err := domain.ConvertMapToFHIRResource[domain.FHIRServiceRequest](resource)
				if err != nil {
					return nil, err
				}

				serviceRequests = append(serviceRequests, serviceRequest)

			case "DocumentReference":
				documentReference, err := domain.ConvertMapToFHIRResource[domain.FHIRDocumentReference](resource)
				if err != nil {
					return nil, err
				}

				documentReferences = append(documentReferences, documentReference)
			}
		}
	}

	docRefMap := make(map[string]string)

	for _, docRef := range documentReferences {
		for _, related := range docRef.Context {
			if related.Reference != nil {
				document, err := docRef.GetDocumentAttachment()
				if err != nil {
					return nil, err
				}

				docRefMap[*related.Reference] = string(*document.URL)
			}
		}
	}

	for _, serviceRequest := range serviceRequests {
		referral := &dto.ReferralDetail{
			ID:           *serviceRequest.ID,
			ReferralDate: serviceRequest.AuthoredOn.Time().Format("January 2, 2006"),
			PatientName:  serviceRequest.Subject.Display,
			Status:       string(serviceRequest.Status),
		}
		if serviceRequest.Subject.ID != nil {
			referral.PatientID = *serviceRequest.Subject.ID
		} else {
			referral.PatientID = serviceRequest.Subject.Display
		}

		receivingFacility, err := c.FHIR.GetFHIROrganization(ctx, serviceRequest.GetPerformerID())
		if err != nil {
			utils.ReportErrorToSentry(err)
			return nil, err
		}

		referral.ReferredTo = *receivingFacility.Resource.Name

		referral.ReferredFor = serviceRequest.GetPatientReferralReason()

		referral.ReferredTests = serviceRequest.GetRequestedServices(common.TestsOrderedLOINCCode)

		serviceRequestReference := fmt.Sprintf("ServiceRequest/%v", *serviceRequest.ID)

		if serviceRequest.Text != nil {
			referral.UsageContext = dto.ScreeningTypeEnum(helpers.ExtractTextFromHTML(string(serviceRequest.Text.Div)))
		}

		if url, found := docRefMap[serviceRequestReference]; found {
			referral.ReferralReportLink = &url
		}

		connection.Edges = append(connection.Edges, dto.ReferralDetailEdge{
			Node: *referral,
		})
	}

	connection.TotalCount = pagedResource.TotalCount
	connection.PageInfo.HasNextPage = pagedResource.HasNextPage
	connection.PageInfo.HasPreviousPage = pagedResource.HasPreviousPage
	connection.PageInfo.StartCursor = &pagedResource.NextCursor
	connection.PageInfo.EndCursor = &pagedResource.PreviousCursor

	return &connection, nil
}

// GetPatientReferralDetails fetches the patient's referral details
func (c *ClinicalImpl) GetPatientReferralDetails(ctx context.Context, serviceRequestID string) (*domain.PatientReferralDetails, error) {
	if serviceRequestID == "" {
		return nil, fmt.Errorf("servicerequest id is required")
	}

	if _, err := uuid.Parse(serviceRequestID); err != nil {
		return nil, fmt.Errorf("invalid service request ID: %s", serviceRequestID)
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	var serviceRequest domain.FHIRServiceRequest

	var organization domain.FHIROrganization

	searchParams := map[string]any{
		"_id":      serviceRequestID,
		"_include": "ServiceRequest:performer",
	}

	resourceBundle, err := c.FHIR.SearchFHIRResource(ctx, "", "ServiceRequest", searchParams, *identifiers, dto.Pagination{})
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, resource := range resourceBundle.Resources {
		if resourceType, ok := resource["resourceType"].(string); ok {
			switch resourceType {
			case "ServiceRequest":
				output, err := domain.ConvertMapToFHIRResource[domain.FHIRServiceRequest](resource)
				if err != nil {
					return nil, err
				}

				serviceRequest = *output

			case "Organization":
				output, err := domain.ConvertMapToFHIRResource[domain.FHIROrganization](resource)
				if err != nil {
					return nil, err
				}

				organization = *output
			}
		}
	}

	receivingFacility, err := organization.GetReceivingFacilityContactDetails()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	category := dto.ObservationCategoryLab
	patientID := domain.SafeStringPointerDeference(serviceRequest.Subject.ID)

	encounterID := domain.SafeStringPointerDeference(serviceRequest.Encounter.ID)

	if encounterID == "" && serviceRequest.Encounter != nil && serviceRequest.Encounter.Reference != nil {
		encounterID = strings.TrimPrefix(*serviceRequest.Encounter.Reference, "Encounter/")
	}

	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     &encounterID,
		Date:            nil,
		ObservationCode: "",
		Category:        &category,
		Pagination:      &dto.Pagination{},
	}

	observations, err := c.GetPatientObservations(ctx, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	var tests []domain.Test

	for _, observation := range observations.Edges {
		// If observation name or value of the test is empty, skip the observation and move to the next
		if observation.Node.Value == "" || observation.Node.Name == "" {
			continue
		}

		obs := domain.Test{
			Name:    observation.Node.Name,
			Results: observation.Node.Value,
		}

		if observation.Node.TimeRecorded != "" {
			date, err := utils.ConvertDateStringToDateScalar(time.RFC3339, observation.Node.TimeRecorded)
			if err != nil {
				return nil, err
			}

			obs.Date = date.String()
		}

		tests = append(tests, obs)
	}

	referral := &domain.PatientReferralDetails{
		PatientName: serviceRequest.Subject.Display,
		ReceivingFacility: domain.ReceivingFacility{
			FacilityName:    receivingFacility.Location,
			FacilityCounty:  receivingFacility.Name,
			FacilityContact: receivingFacility.Phone,
			FacilityEmail:   receivingFacility.Email,
		},
		ReferralReason: domain.ReferralReason{
			Reason: serviceRequest.GetPatientReferralReason(),
			Test:   serviceRequest.GetPatientReferralTest(),
			Date:   serviceRequest.AuthoredOn.Time().Format("Jan 2, 2006"),
			Note:   serviceRequest.GetPractitionersNotes(),
		},
		MedicalHistory: domain.MedicalHistory{
			Tests: tests,
		},
	}

	serviceRequestReference := fmt.Sprintf("ServiceRequest/%s", *serviceRequest.ID)

	documentReferenceSearchParams := map[string]interface{}{
		"contenttype": "application/pdf",
		"based-on":    serviceRequestReference,
		"_sort":       "_lastUpdated",
		"_total":      "accurate",
	}

	if serviceRequest.Subject != nil && serviceRequest.Subject.ID != nil {
		documentReferenceSearchParams["subject"] = fmt.Sprintf("Patient/%s", *serviceRequest.Subject.ID)
	}

	documentReferences, err := c.FHIR.SearchFHIRDocumentReference(ctx, documentReferenceSearchParams, *identifiers, dto.Pagination{})
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	if len(documentReferences.DocumentReferences) == 0 {
		referral.ReferralReportLink = ""
	} else {
		for _, documentReference := range documentReferences.DocumentReferences {
			if len(documentReference.Content) > 0 {
				attachment, err := documentReference.GetDocumentAttachment()
				if err != nil {
					utils.ReportErrorToSentry(err)
				}

				referral.ReferralReportLink = string(*attachment.URL)
			}
		}
	}

	return referral, nil
}

// CreateReferral is used to create a new patient referral in FHIR as a service request.
// This method will mainly be used by advantage backend to save the referral details as part of the
// patient profile
//
// NB: This has been added for backwards compatibility.
// TODO: Purge the ReferPatient method and APIs
func (c *ClinicalImpl) CreateReferral(ctx context.Context, input *dto.CreateReferralInput) (*dto.ServiceRequest, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record a referral in a completed encounter")
	}

	patientReference := encounter.Resource.Subject.Reference
	encounterReference := fmt.Sprintf("Encounter/%s", input.EncounterID)

	startTime := scalarutils.DateTime(time.Now().Format("2006-01-02T15:04:05+03:00"))
	performer := fmt.Sprintf("Organization/%s", input.ReferredTo.FHIROrganisationID)

	codeConcept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, common.ReferralReasonLOINCode)
	if err != nil {
		return nil, err
	}

	var requester *domain.FHIRReferenceInput

	if input.ReferredFrom != nil {
		ref := fmt.Sprintf("Organization/%s", input.ReferredFrom.FHIROrganisationID)
		requester = &domain.FHIRReferenceInput{
			Display:   input.ReferredFrom.FHIROrganisationID,
			Reference: &ref,
		}
	} else {
		identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
		if err != nil {
			utils.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
		}

		ref := fmt.Sprintf("Organization/%s", identifiers.OrganizationID)

		requester = &domain.FHIRReferenceInput{
			Display:   identifiers.OrganizationID,
			Reference: &ref,
		}
	}

	if !domain.ServiceRequestPriorityEnum(input.Urgency).IsValid() {
		return nil, fmt.Errorf("unknown RequestPriority code: %s", input.Urgency)
	}

	serviceRequestSystem := helpers.CodeSystem("service-request-cs")

	serviceRequest := domain.FHIRServiceRequestInput{
		Status:   domain.ServiceRequestStatusActive,
		Intent:   domain.ServiceRequestIntentOrder,
		Priority: domain.ServiceRequestPriorityEnum(input.Urgency),
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  serviceRequestSystem,
						Code:    scalarutils.Code(domain.ReferralCategoryType),
						Display: domain.ReferralCategoryType.Display(),
					},
				},
				Text: domain.ReferralCategoryType.Display(),
			},
		},
		Code: &domain.FHIRCodeableReferenceInput{
			Concept: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&codeConcept.URL),
						Code:    scalarutils.Code(codeConcept.ID),
						Display: codeConcept.GetConceptDisplay(),
					},
				},
				Text: codeConcept.GetConceptDisplay(),
			},
		},
		AuthoredOn: &startTime,
		Subject: &domain.FHIRReferenceInput{
			Reference: patientReference,
			Display:   encounter.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &input.EncounterID,
			Display:   input.EncounterID,
			Reference: &encounterReference,
		},
		OccurrenceDateTime: &input.ReferralDate,
		Requester:          requester,
		Performer: []*domain.FHIRReferenceInput{
			{
				Display:   input.ReferredTo.FHIROrganisationID,
				Reference: &performer,
			},
		},
		Note: []*domain.FHIRAnnotationInput{
			{
				Time: &startTime,
				Text: (*scalarutils.Markdown)(&input.ClinicalHistory),
			},
		},
		// TODO: This should be the actual diagnosis code e.g Breast Cancer Screening etc
		// Papercut: Assign Adrian
		Reason: []*domain.FHIRCodeableReferenceInput{
			{
				Concept: &domain.FHIRCodeableConceptInput{
					Coding: []*domain.FHIRCodingInput{
						{
							System:  &common.LoincSystemURL,
							Code:    common.ReferralLOINCTerminologySystem,
							Display: "Referral note",
						},
					},
					Text: "Referral note",
				},
			},
		},
	}

	if input.Diagnosis != "" {
		condition, err := c.FHIR.GetFHIRCondition(ctx, input.Diagnosis)
		if err != nil {
			return nil, err
		}

		condtionCode := condition.Resource.Code.ConditionCode()
		serviceRequest.Reason = append(serviceRequest.Reason, condtionCode)
	}

	if len(input.Tests) > 0 {
		for _, test := range input.Tests {
			if strings.TrimSpace(test) == "" {
				continue
			}

			referredTest := domain.FHIRCodingInput{
				System:  &common.LoincSystemURL,
				Code:    common.TestsOrderedLOINCCode,
				Display: test,
			}

			serviceRequest.Code.Concept.Coding = append(serviceRequest.Code.Concept.Coding, &referredTest)
		}
	}

	if input.UsageContext != "" {
		additionalStatus := "additional"
		serviceRequest.Text = (*domain.FHIRNarrativeInput)(utils.NarrativeGenerator(string(input.UsageContext), &additionalStatus))
	}

	if input.Specialist != "" {
		specialistExtension := &domain.FHIRExtension{
			URL: "http://savannahghi.org/fhir/StructureDefinition/referred-specialist",
			Extension: []domain.Extension{
				{
					URL:         "specialistName",
					ValueString: input.Specialist,
				},
			},
		}

		serviceRequest.Extension = append(serviceRequest.Extension, specialistExtension)
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

	var meta *dto.MetaInput

	err = mapstructure.Decode(serviceRequest.Meta, &meta)
	if err != nil {
		return nil, err
	}

	payload := &dto.PatientReferralTaskPayload{
		Meta:     meta,
		Referral: output,
	}

	err = c.Pubsub.NotifyCreatePatientReferralTask(ctx, *payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return output, nil
}
