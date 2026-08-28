package base

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// StartEncounter starts an encounter within an episode of care
func (c *BaseImpl) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	if episodeID == "" {
		c.Warn("episode ID is required")
		return "", fmt.Errorf("an episode of care ID is required")
	}

	episodeOfCare, err := c.FHIR.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		c.Error("failed to get episode of care", "error", err.Error())
		utils.ReportErrorToSentry(err)

		return "", err
	}

	encounterClassCode := scalarutils.Code("AMB")
	encounterClassSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassDisplay := string(dto.EncounterClassEnumAmbulatory)
	encounterClassUserSelected := false
	now := time.Now()
	startTime := scalarutils.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	episodeReference := fmt.Sprintf("EpisodeOfCare/%s", *episodeOfCare.Resource.ID)
	encounterPayload := domain.FHIREncounterInput{
		Status: domain.EncounterStatusEnum(domain.EncounterStatusEnumInProgress.String()),
		Class: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &encounterClassSystem,
						Code:         encounterClassCode,
						Display:      strings.ToLower(encounterClassDisplay),
						UserSelected: &encounterClassUserSelected,
					},
				},
			},
		},
		Subject: &domain.FHIRReferenceInput{
			Reference: episodeOfCare.Resource.Patient.Reference,
			Type:      episodeOfCare.Resource.Patient.Type,
			Display:   episodeOfCare.Resource.Patient.Display,
			ID:        episodeOfCare.Resource.Patient.ID,
		},
		EpisodeOfCare: []*domain.FHIRReferenceInput{
			{
				ID:        episodeOfCare.Resource.ID,
				Reference: &episodeReference,
				Display:   *episodeOfCare.Resource.ID,
			},
		},
		ServiceProvider: &domain.FHIRReferenceInput{
			Reference: episodeOfCare.Resource.ManagingOrganization.Reference,
			Type:      episodeOfCare.Resource.ManagingOrganization.Type,
			Display:   *episodeOfCare.Resource.ManagingOrganization.ID,
		},
		ActualPeriod: &domain.FHIRPeriodInput{
			Start: &startTime,
		},
		Participant: []*domain.FHIREncounterParticipantInput{
			{
				Actor: &domain.FHIRReferenceInput{
					Reference: episodeOfCare.Resource.Patient.Reference,
					Type:      episodeOfCare.Resource.Patient.Type,
					Display:   episodeOfCare.Resource.Patient.Display,
				},
			},
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return "", err
	}

	encounterPayload.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	additionalNarrativeStatus := domain.NarrativeStatusEnumGenerated
	narrative := fmt.Sprintf(`<div xmlns="http://www.w3.org/1999/xhtml"><h1>%s</h1></div>`, "Generated Encounter.")
	encounterPayload.Text = &domain.FHIRNarrative{
		Status: &additionalNarrativeStatus,
		Div:    scalarutils.XHTML(narrative),
	}

	encounter, err := c.FHIR.CreateFHIREncounter(ctx, encounterPayload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return "", err
	}

	return *encounter.Resource.ID, nil
}

func (c *BaseImpl) PatchEncounter(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error) {
	if encounterID == "" {
		return nil, fmt.Errorf("an encounterID is required")
	}

	status := domain.EncounterStatusEnum(strings.ToLower(string(input.Status)))
	encounterInput := domain.FHIREncounterInput{
		Status: domain.EncounterStatusEnum(status.String()),
	}

	if status.IsFinal() {
		encounter, err := c.FHIR.GetFHIREncounter(ctx, encounterID)
		if err != nil {
			return nil, fmt.Errorf("unable to get encounter with ID %s: %w", encounterID, err)
		}

		var startTime scalarutils.DateTime
		if encounter.Resource.ActualPeriod == nil {
			startTime = scalarutils.DateTime(time.Now().Format(timeFormatStr))
		} else {
			startTime = encounter.Resource.ActualPeriod.Start
		}

		// workaround for odd date comparison behavior on the Google Cloud Healthcare API
		end := startTime.Time().Add(time.Hour * 24)
		endTime := scalarutils.DateTime(end.Format(timeFormatStr))

		encounterInput.ActualPeriod = &domain.FHIRPeriodInput{Start: &startTime, End: &endTime}
	}

	fhirEncounter, err := c.FHIR.PatchFHIREncounter(ctx, encounterID, encounterInput)
	if err != nil {
		return nil, err
	}

	encounters := []*dto.Encounter{}

	err = mapstructure.Decode([]domain.FHIREncounter{*fhirEncounter}, &encounters)
	if err != nil {
		return nil, err
	}

	return encounters[0], nil
}

// EndEncounter marks an encounter as finished and updates the endtime field
func (c *BaseImpl) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	if encounterID == "" {
		return false, fmt.Errorf("an encounterID is required")
	}

	ok, err := c.FHIR.EndEncounter(ctx, encounterID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, err
	}

	return ok, nil
}

// ListPatientEncounters lists all the encounters that a patient has been part of
func (c *BaseImpl) ListPatientEncounters(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.EncounterConnection, error) {
	if patientID == "" {
		return nil, fmt.Errorf("a patient ID is required")
	}

	err := pagination.Validate()
	if err != nil {
		return nil, err
	}

	_, err = c.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", patientID)

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	encounterResponses, err := c.FHIR.SearchPatientEncounters(ctx, patientReference, nil, *identifiers, *pagination)
	if err != nil {
		return nil, err
	}

	encounters := []*dto.Encounter{}

	err = mapstructure.Decode(encounterResponses.Encounters, &encounters)
	if err != nil {
		return nil, err
	}

	pagedInfo := dto.PageInfo{
		HasNextPage:     encounterResponses.HasNextPage,
		EndCursor:       &encounterResponses.NextCursor,
		HasPreviousPage: encounterResponses.HasPreviousPage,
		StartCursor:     &encounterResponses.PreviousCursor,
	}

	connection := dto.CreateEncounterConnection(encounters, pagedInfo, encounterResponses.TotalCount)

	return &connection, nil
}

// GetEncounterAssociatedResources get all resources assocuated with an encounter
func (c *BaseImpl) GetEncounterAssociatedResources(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResourceOutput, error) {
	if encounterID == "" {
		return nil, fmt.Errorf("an encounterID is required")
	}

	result, err := c.EncounterAssociatedResources(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	output := &dto.EncounterAssociatedResourceOutput{}
	if len(result.RiskAssessment) > 0 {
		output.RiskAssessment = result.RiskAssessment[0]
	}

	if len(result.Consent) > 0 {
		output.Consent = result.Consent[0]
	}

	if len(result.Observation) > 0 {
		output.Observation = result.Observation[0]
	}

	output.Encounter = result.Encounter

	if len(result.Task) > 0 {
		output.Tasks = result.Task
	}

	return output, nil
}

// EncounterAssociatedResources is a helper method that abstract the fetching of specified resources that are associated with a given encounter
func (c *BaseImpl) EncounterAssociatedResources(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResources, error) {
	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, err
	}

	includedResources := []string{
		"RiskAssessment:encounter",
		"Consent:data",
		"Observation:encounter",
		"Task:encounter",
	}

	encounterSearchParams := map[string]interface{}{
		"_id":         encounterID,
		"_sort":       "_lastUpdated",
		"_revinclude": includedResources,
		"_total":      "accurate",
	}

	encounterAllData, err := c.FHIR.SearchFHIREncounterAllData(ctx, encounterSearchParams, *identifiers, dto.Pagination{})
	if err != nil {
		return nil, err
	}

	result := dto.EncounterAssociatedResources{}

	for _, encounterData := range encounterAllData.Resources {
		switch encounterData["resourceType"] {
		case "RiskAssessment":
			riskAssessmentBytes, err := json.Marshal(encounterData)
			if err != nil {
				return nil, err
			}

			var riskAssessment dto.RiskAssessment

			if err := json.Unmarshal(riskAssessmentBytes, &riskAssessment); err != nil {
				return nil, err
			}

			var fhirRiskAssessment domain.FHIRRiskAssessment

			if err := json.Unmarshal(riskAssessmentBytes, &fhirRiskAssessment); err != nil {
				return nil, err
			}

			if fhirRiskAssessment.Text != nil {
				riskAssessment.UsageContext = helpers.ExtractTextFromHTML(string(fhirRiskAssessment.Text.Div))
			}

			// In the cases where questionnaire response ID is null, append on from the reference value
			riskAssessment.AppendQuestionnaireResponseReferenceID()

			result.RiskAssessment = append(result.RiskAssessment, &riskAssessment)
		case "Consent":
			consentBytes, err := json.Marshal(encounterData)
			if err != nil {
				return nil, err
			}

			var consent dto.Consent

			if err := json.Unmarshal(consentBytes, &consent); err != nil {
				return nil, err
			}

			var domainConsent domain.FHIRConsent

			if err := json.Unmarshal(consentBytes, &domainConsent); err != nil {
				return nil, err
			}

			if domainConsent.Text != nil {
				consent.UsageContext = dto.ScreeningTypeEnum(helpers.ExtractTextFromHTML(string(domainConsent.Text.Div)))
			}

			result.Consent = append(result.Consent, &consent)

		case "Observation":
			var observation domain.FHIRObservation

			var observationValue, observationNote, category, instant string

			observationBytes, err := json.Marshal(encounterData)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(observationBytes, &observation); err != nil {
				return nil, err
			}

			if observation.ValueString != nil {
				observationValue = *observation.ValueString
			}

			if observation.Note != nil {
				observationNote = string(*observation.Note[0].Text)
			}

			if observation.EffectiveInstant != nil {
				instant = string(*observation.EffectiveInstant)
			}

			if len(observation.Category) > 0 {
				category = string(*observation.Category[0].Coding[0].Code)
			}

			output := &dto.Observation{
				ID:           *observation.ID,
				Name:         observation.Code.Text,
				Value:        observationValue,
				Status:       dto.ObservationStatus(*observation.Status),
				TimeRecorded: instant,
				Code:         string(*observation.Coding().Code),
				Category:     category,
				Note:         observationNote,
			}

			if observation.Text != nil {
				output.UsageContext = dto.ScreeningTypeEnum(observation.Text.Div)
			}

			result.Observation = append(result.Observation, output)

		case "Encounter":
			var encounter dto.Encounter

			encounterBytes, err := json.Marshal(encounterData)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(encounterBytes, &encounter); err != nil {
				return nil, err
			}

			result.Encounter = &encounter

		case "Task":
			var task domain.FHIRTask

			taskBytes, err := json.Marshal(encounterData)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(taskBytes, &task); err != nil {
				return nil, err
			}

			if string(*task.Status) == "requested" {
				taskOutput := &dto.TaskOutput{
					ID:          *task.ID,
					Description: task.Description,
					Status:      dto.TaskStatus(*task.Status),
				}

				if task.Reason != nil && task.Reason[0].Reference != nil && task.Reason[0].Reference.Display != "" {
					taskOutput.Workflow = task.Reason[0].Reference.Display
				}

				if task.BusinessStatus != nil && task.BusinessStatus.Text != "" {
					taskOutput.Task = task.BusinessStatus.Text
				}

				if task.Encounter.ID != nil {
					taskOutput.EncounterID = *task.Encounter.ID
				} else if task.Encounter.Reference != nil {
					// In some cases, the encounter ID might be null, so we extract it from the reference value
					taskOutput.EncounterID = strings.TrimPrefix(*task.Encounter.Reference, "Encounter/")
				}

				result.Task = append(result.Task, taskOutput)
			}
		}
	}

	return &result, nil
}

// GetScreeningReport serves a specific usecase to retrieve screening reports
func (c *BaseImpl) GetScreeningReport(ctx context.Context, encounterID string, status domain.ServiceRequestStatusEnum) (*dto.ScreeningReport, error) {
	encounter, err := c.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	result, err := c.EncounterAssociatedResources(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	output := &dto.ScreeningReport{
		RiskAssessment:  nil,
		Consent:         nil,
		Observation:     dto.ObservationOutput{},
		ReferralDetails: []dto.ReferralDetail{},
	}

	if len(result.RiskAssessment) > 0 {
		output.RiskAssessment = result.RiskAssessment
	}

	if len(result.Consent) > 0 {
		output.Consent = result.Consent
	}

	c.processObservations(result.Observation, output)

	reportStatus := status

	if reportStatus == "" {
		reportStatus = domain.ServiceRequestStatusActive
	}

	if err := c.processPatientReferrals(
		ctx,
		&dto.ReferralSearchInput{
			PatientID:   encounter.Resource.Subject.ID,
			EncounterID: &encounterID,
			Pagination:  &dto.Pagination{},
			Status:      reportStatus,
		},
		output,
	); err != nil {
		return nil, err
	}

	c.filterOutReferrals(output)

	return output, nil
}

// processObservations is a helper function that processes observations and categorizes them
func (c *BaseImpl) processObservations(observations []*dto.Observation, output *dto.ScreeningReport) {
	for _, observation := range observations {
		if observation.Category == "vital-signs" {
			continue
		}

		obs := dto.Observation{
			ID:           observation.ID,
			Name:         observation.Name,
			Value:        observation.Value,
			Status:       observation.Status,
			Code:         observation.Code,
			TimeRecorded: observation.TimeRecorded,
			Category:     observation.Category,
			Note:         observation.Note,
			UsageContext: dto.ScreeningTypeEnum(helpers.ExtractTextFromHTML(string(observation.UsageContext))),
		}

		switch observation.Category {
		case "exam":
			output.Observation.Examinations = append(output.Observation.Examinations, obs)
		case "laboratory", "imaging", "procedure":
			output.Observation.Tests = append(output.Observation.Tests, obs)
		}
	}
}

// processPatientReferrals is a helper function that processes patient referrals
func (c *BaseImpl) processPatientReferrals(ctx context.Context, searchInput *dto.ReferralSearchInput, output *dto.ScreeningReport) error {
	referralInput := &dto.ReferralSearchInput{
		PatientID:   searchInput.PatientID,
		EncounterID: searchInput.EncounterID,
		Pagination:  searchInput.Pagination,
		Status:      searchInput.Status,
	}

	patientReferrals, err := c.GetPatientReferrals(ctx, referralInput)
	if err != nil {
		return err
	}

	if len(patientReferrals.Edges) > 0 {
		for _, edge := range patientReferrals.Edges {
			output.ReferralDetails = append(output.ReferralDetails, edge.Node)
		}
	}

	return nil
}

// filterOutReferrals is a helper function that separates diagnostics & non-diagnostics referrals.
// This is influenced the need to add (test) lab orders (intra & extra referrals) as part of the tests
// when showing screening report.
func (c *BaseImpl) filterOutReferrals(output *dto.ScreeningReport) {
	filteredReferrals := make([]dto.ReferralDetail, 0, len(output.ReferralDetails))

	for _, ref := range output.ReferralDetails {
		if ref.ReferredFor != "Diagnostics" {
			// append non-diagnostics referrals to the filtered list
			filteredReferrals = append(filteredReferrals, ref)
			continue
		}

		// diagnostic referrals (lab orders) are mapped into the tests section
		if len(ref.ReferredTests) > 0 {
			output.Observation.ReferredTest = append(output.Observation.ReferredTest, dto.ReferralDetail{
				ID:                 ref.ID,
				ReferralDate:       ref.ReferralDate,
				ReferredFor:        ref.ReferredFor,
				ReferredTests:      ref.ReferredTests,
				ReferredTo:         ref.ReferredTo,
				ReferralReportLink: ref.ReferralReportLink,
				UsageContext:       ref.UsageContext,
			})
		}
	}

	output.ReferralDetails = filteredReferrals
}

// EndScreening ends the screening process by marking the encounter as `finished`.
func (c *BaseImpl) EndScreening(ctx context.Context, encounterID string) (bool, error) {
	if encounterID == "" {
		return false, fmt.Errorf("encounter ID is required")
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return false, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return false, fmt.Errorf("encounter is already completed")
	}

	ok, err := c.EndEncounter(ctx, *encounter.Resource.ID)
	if err != nil {
		return false, err
	}

	return ok, nil
}
