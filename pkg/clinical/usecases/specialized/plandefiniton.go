package specialized

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

func (c *SpecializedImpl) CreatePlanDefinition(ctx context.Context, input *dto.PlanDefinitionInput) (*domain.FHIRPlanDefinition, error) {
	if err := helpers.Validate(input); err != nil {
		return nil, fmt.Errorf("incomplete plan definition input: %w", err)
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant meta tags: %w", err)
	}

	var (
		planTypeCode                     = "order-set"
		planSystem                       = "http://terminology.hl7.org/CodeSystem/plan-definition-type"
		useContextSystem                 = "http://terminology.hl7.org/CodeSystem/v2-0265"
		useContextCode                   = "CAN"
		planDefinitionUsageContextSystem = "http://terminology.hl7.org/CodeSystem/usage-context-type"
		planDefinitionUsageContextCode   = "focus"
	)

	var topLevelActions []domain.PlanDefinitionAction

	for _, action := range input.Action {
		builtAction, err := c.buildAction(ctx, tags, &action)
		if err != nil {
			return nil, fmt.Errorf("failed to build action: %w", err)
		}

		topLevelActions = append(topLevelActions, *builtAction)
	}

	planDefinitionName := strings.ReplaceAll(input.Title, " ", "")
	currentTime := time.Now().Format(time.RFC3339)

	planDef := &domain.FHIRPlanDefinition{
		Name:  &planDefinitionName,
		Title: &input.Title,
		Type: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  (*scalarutils.URI)(&planSystem),
					Code:    (*scalarutils.Code)(&planTypeCode),
					Display: planTypeCode,
				},
			},
			Text: planTypeCode,
		},
		Status:      domain.PublicationStatusActive,
		Date:        &currentTime,
		Description: &input.Description,
		UseContext: []domain.FHIRUsageContext{
			{
				Code: &domain.FHIRCoding{
					System:  (*scalarutils.URI)(&planDefinitionUsageContextSystem),
					Code:    (*scalarutils.Code)(&planDefinitionUsageContextCode),
					Display: "Clinical Focus",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{
						{
							System:  (*scalarutils.URI)(&useContextSystem),
							Code:    (*scalarutils.Code)(&useContextCode),
							Display: "Cancer",
						},
					},
					Text: "Cancer",
				},
			},
		},
		Action: topLevelActions,
		Meta: &domain.FHIRMetaInput{
			Tag: tags,
		},
	}

	createdPlanDef, err := c.FHIR.CreateFHIRPlanDefinition(ctx, planDef)
	if err != nil {
		return nil, fmt.Errorf("failed to create FHIR PlanDefinition: %w", err)
	}

	return createdPlanDef, nil
}

func (c *SpecializedImpl) buildAction(ctx context.Context, tags []domain.FHIRCodingInput, input *dto.PlanAction) (*domain.PlanDefinitionAction, error) {
	action := &domain.PlanDefinitionAction{
		Title: &input.Title,
	}

	if input.Description != "" {
		action.Description = &input.Description
	}

	if input.TimingTiming != nil && input.TimingTiming.Repeat != nil {
		repeat := input.TimingTiming.Repeat
		repeatPeriod := strconv.Itoa(repeat.Period)
		immediateWhen := []domain.TimingRepeatWhenEnum{
			domain.TimingRepeatWhenEnumImmediate,
		}

		action.TimingTiming = &domain.FHIRTiming{
			Repeat: &domain.FHIRTimingRepeat{
				Frequency:  &repeat.Frequency,
				Period:     (*json.Number)(&repeatPeriod),
				PeriodUnit: (*domain.UnitsOfTimeEnum)(&repeat.PeriodUnit),
				Count:      &repeat.Count,
				Offset:     &repeat.Offset,
				When:       immediateWhen,
			},
		}
	}

	// handle medications by linking them to ActivityDefinitions
	for _, med := range input.Medications {
		medication, err := c.FHIR.FetchMedicationByID(ctx, med.MedicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch medication %s: %w", med.MedicationID, err)
		}

		var (
			medicationReference     = fmt.Sprintf("Medication/%s", *medication.ID)
			activityDefinitionName  = "SGHIMedicationActivityDefinition"
			activityDefinitionTitle = fmt.Sprintf("%s Administration", medication.Code.Text)
			currentTime             = scalarutils.DateTime(time.Now().Format(time.RFC3339))
		)

		activityDef := &domain.FHIRActivityDefinition{
			Name:        &activityDefinitionName,
			Title:       &activityDefinitionTitle,
			Description: &med.Dosage.AdministrationInstructions,
			Kind:        "ServiceRequest",
			Meta: &domain.FHIRMetaInput{
				Tag: tags,
			},
			ProductReference: &domain.FHIRReference{
				ID:        medication.ID,
				Reference: &medicationReference,
				Display:   medication.Code.Text,
			},
			Dosage: []domain.FHIRDosage{
				{
					Route: &domain.FHIRCodeableConcept{
						Coding: []*domain.FHIRCoding{
							{
								System:  helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier),
								Code:    med.Dosage.Route.Code,
								Display: med.Dosage.Route.Display,
							},
						},
						Text: med.Dosage.Route.Display,
					},
					Method: &domain.FHIRCodeableConcept{
						Coding: []*domain.FHIRCoding{
							{
								System:  helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier),
								Code:    med.Dosage.Method.Code,
								Display: med.Dosage.Method.Display,
							},
						},
						Text: med.Dosage.Method.Display,
					},
					Timing: &domain.FHIRTiming{
						Event: []*scalarutils.DateTime{
							&currentTime,
						},
					},
					DoseAndRate: []*domain.FHIRDosageDoseandrate{
						{
							DoseQuantity: &domain.FHIRQuantity{
								Value:  med.DoseQuantity,
								Unit:   med.DoseUnit,
								Code:   (*scalarutils.Code)(&med.DoseUnit),
								System: (scalarutils.URI)(helpers.QuantitySystem),
							},
						},
					},
				},
			},
			Status: domain.PublicationStatusActive,
		}

		activityDefOutput, err := c.FHIR.CreateFHIRActivityDefinition(ctx, activityDef)
		if err != nil {
			return nil, fmt.Errorf("failed to create activity definition for %s: %w", *medication.ID, err)
		}

		// add sub-action referencing this ActivityDefinition
		activityDefinitionCanonicalURL := fmt.Sprintf("%s/ActivityDefinition/%s", serverutils.MustGetEnvVar("HAPI_FHIR_BASE_URL"), *activityDefOutput.ID)
		subAction := &domain.PlanDefinitionAction{
			Title:               &medication.Code.Text,
			DefinitionCanonical: &activityDefinitionCanonicalURL,
		}

		action.Action = append(action.Action, *subAction)
	}

	// recursively build nested sub-actions
	for _, subInput := range input.Action {
		childAction, err := c.buildAction(ctx, tags, &subInput)
		if err != nil {
			return nil, err
		}

		action.Action = append(action.Action, *childAction)
	}

	return action, nil
}

func (c *SpecializedImpl) RetrievePlanDefinition(ctx context.Context, name string) (*dto.PlanDefinitionOutputConnection, error) {
	filterParams := map[string]any{
		"_sort":  "_lastUpdated",
		"_total": "accurate",
	}

	if name != "" {
		filterParams["name"] = name
	}

	results, err := c.FHIR.SearchFHIRPlanDefinition(ctx, filterParams, dto.Pagination{})
	if err != nil {
		return nil, err
	}

	planDefinitionList := []*domain.FHIRPlanDefinition{}

	for _, planDefinition := range results.PlanDefinition {
		planDefinitionList = append(planDefinitionList, &planDefinition)
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     results.HasNextPage,
		EndCursor:       &results.NextCursor,
		HasPreviousPage: results.HasPreviousPage,
		StartCursor:     &results.PreviousCursor,
	}

	connection := dto.CreatePlanDefinitionConnection(planDefinitionList, pageInfo, results.TotalCount)

	return &connection, nil
}
