package base

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// AddTestResultsLater creates a task to add test results later.
// Note: AddTestResultsLater needs to be re-evaluated.
func (c *BaseImpl) AddTestResultsLater(ctx context.Context, task *dto.TaskInput) (*dto.TaskOutput, error) {
	encounter, err := c.FHIR.GetFHIREncounter(ctx, task.EncounterID)
	if err != nil {
		return nil, err
	}

	organizationID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	taskMeta := &domain.FHIRMetaInput{
		Tag: tags,
	}

	requestedStatus := dto.RequestedTasksStatus.String()
	orderIntent := "order"
	priority := "routine"
	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	requesterReference := fmt.Sprintf("Organization/%s", organizationID)
	today := time.Now().Format(time.RFC3339)
	authoredOn := (*scalarutils.DateTime)(&today)

	taskPayload := &domain.FHIRTaskInput{
		BusinessStatus: &domain.FHIRCodeableConceptInput{
			Text: task.Workflow.Text(),
		},
		Meta:     taskMeta,
		Status:   (*scalarutils.Code)(&requestedStatus),
		Intent:   (*scalarutils.Code)(&orderIntent),
		Priority: (*scalarutils.Code)(&priority),
		For: &domain.FHIRReferenceInput{
			ID:        &encounter.Resource.Subject.Display,
			Reference: encounter.Resource.Subject.Reference,
			Display:   encounter.Resource.Subject.Display,
		},
		StatusReason: &domain.FHIRCodeableReferenceInput{
			Concept: &domain.FHIRCodeableConceptInput{
				Text: task.Workflow.Text(),
			},
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &task.EncounterID,
			Reference: &encounterReference,
			Display:   task.EncounterID,
		},
		Requester: &domain.FHIRReferenceInput{
			ID:        &organizationID,
			Reference: &requesterReference,
			Display:   organizationID,
		},
		Owner: &domain.FHIRReferenceInput{
			Reference: &requesterReference,
			Display:   organizationID,
		},

		Reason: []*domain.FHIRCodeableReference{
			{
				Concept: &domain.FHIRCodeableConcept{
					Text: fmt.Sprintf("%s test", task.Task),
				},
			},
		},
		AuthoredOn: authoredOn,

		RequestedPerformer: []*domain.FHIRCodeableReference{
			{
				Reference: &domain.FHIRReference{
					Reference: &requesterReference,
					Display:   organizationID,
				},
			},
		},
	}

	if task.Description != "" {
		taskPayload.Description = &task.Description
	}

	taskOutput, err := c.CreateTask(ctx, taskPayload)
	if err != nil {
		return nil, err
	}

	return &dto.TaskOutput{
		ID:          *taskOutput.ID,
		EncounterID: taskOutput.Encounter.ResourceID(),
		Task:        taskOutput.BusinessStatus.Text,
		Description: taskOutput.Description,
		Status:      dto.TaskStatus(*taskOutput.Status),
		Workflow:    taskOutput.StatusReason.Concept.Text,
	}, nil
}

// CreateTask is used tp create FHIR task
func (c *BaseImpl) CreateTask(ctx context.Context, task *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
	return c.FHIR.CreateFHIRTask(ctx, task)
}

// UpdateTask is used to update the state of a task
func (c *BaseImpl) UpdateTask(ctx context.Context, taskID string, updateData *dto.PatchTaskInput) (bool, error) {
	_, err := uuid.Parse(taskID)
	if err != nil {
		return false, fmt.Errorf("invalid task id: %s", taskID)
	}

	if updateData.Status.IsValid() && updateData.UpdateReason == "" {
		return false, fmt.Errorf("the reason for changing the task status is required")
	}

	if updateData.UpdateReason != "" && updateData.Status == "" {
		return false, fmt.Errorf("you need to provide task status since you have given a reason for the update")
	}

	updatePayload := &domain.FHIRTaskInput{
		ID:     &taskID,
		Status: (*scalarutils.Code)(&updateData.Status),
	}

	_, err = c.FHIR.PatchFHIRTask(ctx, *updatePayload)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ListTasks is used to list available tasks
func (c *BaseImpl) ListTasks(ctx context.Context, filter *dto.TaskFilterInput, pagination dto.Pagination) (*dto.TaskOutputConnection, error) {
	taskSearchParams := map[string]interface{}{
		"_sort":  "-_lastUpdated",
		"_total": "accurate",
		"status": filter.SetStatus(),
	}

	err := filter.TaskFilters(taskSearchParams)
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
			connection := dto.CreateTaskConnection([]*dto.TaskOutput{}, dto.PageInfo{}, &emptyTotal)

			return &connection, nil
		}

		taskSearchParams["patient"] = strings.Join(references, ",")
	}

	results, err := c.FHIR.SearchFHIRTask(ctx, taskSearchParams, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	taskList := []*dto.TaskOutput{}

	for _, taskItem := range results.Tasks {
		if taskItem.Status == nil || taskItem.BusinessStatus == nil {
			continue
		}

		taskList = append(taskList, mapFHIRTaskToTaskDTO(taskItem))
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     results.HasNextPage,
		EndCursor:       &results.NextCursor,
		HasPreviousPage: results.HasPreviousPage,
		StartCursor:     &results.PreviousCursor,
	}

	connection := dto.CreateTaskConnection(taskList, pageInfo, results.TotalCount)

	return &connection, nil
}

func mapFHIRTaskToTaskDTO(task domain.FHIRTask) *dto.TaskOutput {
	output := &dto.TaskOutput{
		ID:       *task.ID,
		Status:   dto.TaskStatus(*task.Status),
		Priority: string(*task.Priority),
		Subject: &dto.Reference{
			ID:        task.For.ResourceID(),
			Reference: *task.For.Reference,
			Display:   task.For.Display,
			Identifier: &dto.Identifier{
				Value: new(string),
			},
		},
	}

	if task.For != nil && task.For.Identifier != nil {
		output.Subject.Identifier = &dto.Identifier{
			Value: &task.For.Identifier.Value,
		}
	}

	if task.AuthoredOn != nil {
		authoredOn := (*scalarutils.DateTime)(task.AuthoredOn)
		lastUpdateTime := task.Meta.LastUpdated.Format(time.RFC3339)

		output.AuthoredOn = authoredOn
		output.LastUpdated = (*scalarutils.DateTime)(&lastUpdateTime)
	}

	if task.Encounter != nil {
		output.EncounterID = task.Encounter.ResourceID()
	}

	if task.Reason != nil {
		output.Task = task.Reason[0].Concept.Text
	}

	if task.ExecutionPeriod != nil {
		output.DueDate = &task.ExecutionPeriod.End
	}

	if len(task.Note) > 0 {
		output.Notes = &dto.Annotation{
			Text: *task.Note[0].Text,
		}
	}

	if task.Description != "" {
		output.Description = task.Description
	} else if task.BusinessStatus != nil {
		output.Description = task.BusinessStatus.Text
	}

	if task.StatusReason != nil {
		output.StatusReason = task.StatusReason.Concept.Text
	}

	return output
}

// getPatientPhoneTelecom returns the patient's phone contact value. It prefers a contact
// point explicitly marked as a phone, falling back to the first non-empty telecom value.
func getPatientPhoneTelecom(patient *domain.FHIRPatient) string {
	if patient == nil {
		return ""
	}

	var fallback string

	for _, contact := range patient.Telecom {
		if contact == nil || contact.Value == nil || *contact.Value == "" {
			continue
		}

		if contact.System != nil && *contact.System == domain.ContactPointSystemEnumPhone {
			return *contact.Value
		}

		if fallback == "" {
			fallback = *contact.Value
		}
	}

	return fallback
}

// GetTaskByID is used to retrieve a task by given its ID.
func (c *BaseImpl) GetTaskByID(ctx context.Context, taskID string) (*dto.TaskOutput, error) {
	_, err := uuid.Parse(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %s", taskID)
	}

	results, err := c.FHIR.GetFHIRTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	taskOutput := mapFHIRTaskToTaskDTO(*results.Resource)

	if taskOutput != nil {
		// When the task's subject reference carries no business identifier, look up the
		// referenced patient and fall back to their phone (telecom) so the subject
		// identifier is still populated.
		if subject := results.Resource.For; subject != nil && subject.Identifier == nil {
			if patientID := subject.ResourceID(); patientID != "" {
				patient, err := c.FHIR.GetFHIRPatient(ctx, patientID)
				if err != nil {
					return nil, err
				}

				if phone := getPatientPhoneTelecom(patient.Resource); phone != "" {
					taskOutput.Subject.Identifier = &dto.Identifier{
						Value: &phone,
					}
				}
			}
		}

		terminologySystemCodes := []string{common.ReferralLOINCTerminologySystem, common.IntraReferralLOINCCode}

		documentReferenceParams := map[string]any{
			"type":   strings.Join(terminologySystemCodes, ","),
			"_total": "accurate",
		}

		basedOn, err := results.Resource.GetServiceRequestIDFromTask()
		if err != nil {
			return nil, err
		}

		if basedOn != "" {
			documentReferenceParams["based-on"] = basedOn

			identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
			if err != nil {
				return nil, err
			}

			documentReferences, err := c.fetchDocumentReferences(ctx, documentReferenceParams, *identifiers)
			if err != nil {
				return nil, err
			}

			taskOutput.Attachment = documentReferences
		}
	}

	return taskOutput, nil
}

// fetchDocumentReferences retrieves document references based on the given parameters.
func (c *BaseImpl) fetchDocumentReferences(ctx context.Context, params map[string]interface{}, identifiers dto.TenantIdentifiers) ([]*dto.AttachmentResponse, error) {
	documentReferenceResults, err := c.FHIR.SearchFHIRDocumentReference(ctx, params, identifiers, dto.Pagination{})
	if err != nil {
		return nil, fmt.Errorf("failed to search FHIR document references: %w", err)
	}

	var response []*dto.AttachmentResponse

	for _, documentRef := range documentReferenceResults.DocumentReferences {
		attachment, err := documentRef.GetDocumentAttachment()
		if err != nil {
			return nil, err
		}

		output := &dto.Attachment{
			URL: scalarutils.URL(*attachment.URL),
		}

		if attachment.Title != nil {
			output.Title = *attachment.Title
		}

		if attachment.ContentType != nil {
			output.ContentType = *attachment.ContentType
		}

		attachmentResponse := &dto.AttachmentResponse{
			Attachment: output,
		}

		terminologyCode, err := documentRef.GetDocumentType()
		if err != nil {
			return nil, err
		}

		switch terminologyCode {
		case common.ReferralLOINCTerminologySystem:
			attachmentResponse.Type = "EXTERNAL FACILITY REFERRAL"
		case common.IntraReferralLOINCCode:
			attachmentResponse.Type = "INTERNAL FACILITY REFERRAL"
		}

		response = append(response, attachmentResponse)
	}

	return response, nil
}

// CreateAppointmentTask creates an appointment and a follow-up task that is associated with the appointment
func (c *BaseImpl) CreateAppointmentTask(ctx context.Context, encounterID string, tags *domain.FHIRMetaInput, appointment *domain.FHIRAppointment) (*domain.FHIRTask, error) {
	requestedStatus := dto.RequestedTasksStatus.String()
	orderIntent := "plan"
	priority := "routine"

	organizationID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	requesterReference := fmt.Sprintf("Organization/%s", organizationID)
	appointmentReference := fmt.Sprintf("Appointment/%s", *appointment.ID)
	today := time.Now().Format(time.RFC3339)
	authoredOn := (*scalarutils.DateTime)(&today)
	encounterReference := fmt.Sprintf("Encounter/%s", encounterID)

	taskPayload := &domain.FHIRTaskInput{
		BusinessStatus: &domain.FHIRCodeableConceptInput{
			Text: "Appointment",
		},
		Meta:     tags,
		Status:   (*scalarutils.Code)(&requestedStatus),
		Intent:   (*scalarutils.Code)(&orderIntent),
		Priority: (*scalarutils.Code)(&priority),
		Focus: &domain.FHIRReferenceInput{
			ID:        appointment.ID,
			Reference: &appointmentReference,
			Display:   *appointment.ID,
		},
		Requester: &domain.FHIRReferenceInput{
			ID:        &organizationID,
			Reference: &requesterReference,
			Display:   organizationID,
		},
		Owner: &domain.FHIRReferenceInput{
			Reference: &requesterReference,
			Display:   organizationID,
		},
		AuthoredOn: authoredOn,

		Encounter: &domain.FHIRReferenceInput{
			ID:        &encounterID,
			Reference: &encounterReference,
			Display:   encounterID,
		},
		ExecutionPeriod: &domain.FHIRPeriodInput{
			End: (*scalarutils.DateTime)(appointment.End),
		},
		RequestedPerformer: []*domain.FHIRCodeableReference{
			{
				Reference: &domain.FHIRReference{
					Reference: &requesterReference,
					Display:   organizationID,
				},
			},
		},
		Reason: []*domain.FHIRCodeableReference{},
	}

	if appointment.Reason != nil {
		taskPayload.Reason = []*domain.FHIRCodeableReference{
			{
				Concept: &domain.FHIRCodeableConcept{
					Text: appointment.Reason[0].Concept.Text,
				},
			},
		}
	}

	if appointment.Participant != nil {
		taskPayload.For = &domain.FHIRReferenceInput{
			ID:        appointment.Participant[0].Actor.ID,
			Reference: appointment.Participant[0].Actor.Reference,
			Display:   appointment.Participant[0].Actor.Display,
		}
	}

	if appointment.Reason != nil {
		taskPayload.Reason[0].Concept.Text = appointment.Reason[0].Concept.Text
	}

	return c.CreateTask(ctx, taskPayload)
}

// CreateReferralTask is responsible for creating a task associated with patient referral
func (c *BaseImpl) CreateReferralTask(
	ctx context.Context,
	tags *dto.MetaInput,
	serviceRequest *dto.ServiceRequest,
) (*domain.FHIRTask, error) {
	var (
		orderIntent = "plan"
		priority    = "routine"
	)

	now := time.Now()
	authoredOn := scalarutils.DateTime(now.Format(time.RFC3339))
	executionEnd := scalarutils.DateTime(now.AddDate(0, 0, 7).Format(time.RFC3339)) // 7 days from now

	encounterReference := fmt.Sprintf("Encounter/%s", serviceRequest.Encounter.ResourceID())
	requesterReference := fmt.Sprintf("Organization/%s", *tags.Tag[0].Code)

	labOrder, err := c.CreateLaboratoryOrder(ctx, serviceRequest.ID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	labOrderReference := fmt.Sprintf("ServiceRequest/%s", labOrder.ID)
	serviceRequestReference := fmt.Sprintf("ServiceRequest/%s", serviceRequest.ID)

	var meta *domain.FHIRMetaInput

	if err := mapstructure.Decode(tags, &meta); err != nil {
		return nil, err
	}

	labOrderType := domain.LabOrderServiceRequestType.String()
	patientReferral := domain.ReferralServiceRequestType.String()
	requestedStatus := dto.RequestedTasksStatus.String()

	receivingFacilityRef := fmt.Sprintf("Organization/%s", labOrder.ReceivingFacility)

	taskPayload := &domain.FHIRTaskInput{
		Meta:     meta,
		Status:   (*scalarutils.Code)(&requestedStatus),
		Intent:   (*scalarutils.Code)(&orderIntent),
		Priority: (*scalarutils.Code)(&priority),
		For: &domain.FHIRReferenceInput{
			Reference: &serviceRequest.Subject.Reference,
			Display:   serviceRequest.Subject.Display,
		},
		BasedOn: []*domain.FHIRReferenceInput{
			{
				ID:        &labOrder.ID, // reference to lab order
				Reference: &labOrderReference,
				Display:   labOrderType,
			},
			{
				ID:        &serviceRequest.ID, // reference to referral service request
				Reference: &serviceRequestReference,
				Display:   patientReferral,
			},
		},
		Requester: &domain.FHIRReferenceInput{
			ID:        (*string)(&meta.Tag[0].Code),
			Reference: &requesterReference,
			Display:   requesterReference,
		},
		Owner: &domain.FHIRReferenceInput{
			Reference: &requesterReference,
		},
		AuthoredOn: &authoredOn,
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterReference,
		},
		Focus: &domain.FHIRReferenceInput{
			Reference: &serviceRequestReference,
		},
		ExecutionPeriod: &domain.FHIRPeriodInput{
			End: &executionEnd,
		},
		RequestedPerformer: []*domain.FHIRCodeableReference{
			{
				Reference: &domain.FHIRReference{
					Reference: &receivingFacilityRef,
				},
			},
		},
	}

	if labOrder.OrderDetails.Name != "" {
		taskPayload.BusinessStatus = &domain.FHIRCodeableConceptInput{
			Text: labOrder.OrderDetails.Name,
		}
	}

	return c.CreateTask(ctx, taskPayload)
}

func (c *BaseImpl) FetchPatientCarePlan(ctx context.Context, encounterID string) (*dto.CarePlanOutput, error) {
	encounter, err := c.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", encounter.Resource.Subject.ResourceID())

	// Fetch care plan for that patient
	searchParams := map[string]interface{}{
		"patient":  patientReference,
		"status":   "active",
		"intent":   "plan",
		"category": "21967-5", // Chemo concept code
		"_sort":    "_lastUpdated",
	}

	patientCarePlan, err := c.FHIR.SearchFHIRCarePlan(ctx, searchParams, dto.Pagination{})
	if err != nil {
		return nil, err
	}

	if len(patientCarePlan.CarePlan) > 0 {
		output := dto.CarePlanOutput{
			ID:          *patientCarePlan.CarePlan[0].ID,
			Title:       *patientCarePlan.CarePlan[0].Title,
			Description: *patientCarePlan.CarePlan[0].Description,
			EncounterID: *encounter.Resource.ID,
			Patient: dto.Patient{
				ID:   patientCarePlan.CarePlan[0].Subject.ResourceID(),
				Name: patientCarePlan.CarePlan[0].Subject.Display,
			},
		}

		// Fetch chemo phases
		for _, chemoPhase := range patientCarePlan.CarePlan[0].Activity {
			_, err := c.fetchChemoPhaseActivities(ctx, chemoPhase, &output)
			if err != nil {
				return nil, err
			}
		}

		return &output, nil
	} else {
		return nil, fmt.Errorf("patient has no care plan")
	}
}

func (c *BaseImpl) fetchChemoPhaseActivities(
	ctx context.Context,
	chemoPhase domain.CarePlanActivity,
	carePlan *dto.CarePlanOutput,
) (*dto.CarePlanOutput, error) {
	if len(chemoPhase.PerformedActivity) > 0 {
		for _, activity := range chemoPhase.PerformedActivity {
			task, err := c.FHIR.GetFHIRTask(ctx, activity.Reference.ResourceID())
			if err != nil {
				return nil, err
			}

			phases := &dto.ChemoPhases{
				ID:          *task.Resource.ID,
				Title:       task.Resource.BusinessStatus.Text,
				Description: task.Resource.Description,
				Status:      string(*task.Resource.Status),
			}

			phaseCycles, err := c.fetchChemoCycles(ctx, *task.Resource.ID)
			if err != nil {
				return nil, err
			}

			phases.Cycles = phaseCycles

			carePlan.TreatmentPhases = append(carePlan.TreatmentPhases, phases)
		}

		return carePlan, nil
	}

	return nil, nil
}

func (c *BaseImpl) fetchChemoCycles(
	ctx context.Context,
	taskID string,
) ([]*dto.Cycles, error) {
	searchParam := map[string]interface{}{
		"part-of": fmt.Sprintf("Task/%s", taskID),
	}

	tasks, err := c.FHIR.SearchFHIRTask(ctx, searchParam, dto.TenantIdentifiers{}, dto.Pagination{})
	if err != nil {
		return nil, err
	}

	var cycles []*dto.Cycles

	for _, cycle := range tasks.Tasks {
		cy := &dto.Cycles{
			ID:          *cycle.ID,
			Title:       cycle.BusinessStatus.Text,
			Description: cycle.Description,
			Status:      string(*cycle.Status),
		}

		cycles = append(cycles, cy)
	}

	return cycles, nil
}
