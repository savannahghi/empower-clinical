package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	hapifhirmodels "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

func (c *ClinicalImpl) CreatePatientCarePlan(ctx context.Context, input *dto.CarePlanInput) error {
	if err := helpers.Validate(input); err != nil {
		return err
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return err
	}

	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return err
	}

	payload := dto.CarePlanPayload{
		Data:       *input,
		Tags:       tags,
		FacilityID: facilityID,
	}

	err = c.Pubsub.NotifyCreatePatientCarePlan(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

// PatientCarePlan creates a FHIR CarePlan and all its associated Task resources
// (for phases and cycles) in a single, atomic FHIR transaction.
func (c *ClinicalImpl) PatientCarePlan(ctx context.Context, input *dto.CarePlanPayload) (*domain.FHIRCarePlan, error) {
	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.Data.EncounterID)
	if err != nil {
		return nil, err
	}

	planDefinition, err := c.FHIR.FetchPlanDefinitionByID(ctx, input.Data.PlanDefinitionID)
	if err != nil {
		return nil, err
	}

	allTaskEntries, err := c.prepareCarePlanActivityEntries(ctx, planDefinition, input.Tags, encounter.Resource, input.FacilityID)
	if err != nil {
		return nil, fmt.Errorf("error preparing care plan activities: %w", err)
	}

	carePlanName := fmt.Sprintf("%s Care Plan", encounter.Resource.Subject.Display)
	subjectReference := encounter.Resource.Subject.Reference
	today := time.Now().Format(time.RFC3339)
	orgRef := fmt.Sprintf("Organization/%s", input.FacilityID)
	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	chemotherapyCode := common.CancerChemoTerminologyCode

	carePlan := &domain.FHIRCarePlan{
		ResourceType: "CarePlan",
		Meta: &domain.FHIRMetaInput{
			Tag: input.Tags,
		},
		InstantiatesCanonical: []string{
			fmt.Sprintf("%s/PlanDefinition/%s", serverutils.MustGetEnvVar("HAPI_FHIR_BASE_URL"), *planDefinition.ID),
		},
		Status: "active",
		Intent: "plan",
		Encounter: &domain.FHIRReference{
			ID:        encounter.Resource.ID,
			Reference: &encounterRef,
			Display:   *encounter.Resource.ID,
		},
		Category: []domain.FHIRCodeableConcept{
			{
				Coding: []*domain.FHIRCoding{
					{
						System:  &common.LoincSystemURL,
						Code:    (*scalarutils.Code)(&chemotherapyCode),
						Display: "Chemotherapy Cancer",
					},
				},
				Text: "Chemotherapy Cancer",
			},
		},
		Title:       &carePlanName,
		Description: planDefinition.Description,
		Subject: domain.FHIRReference{
			ID:        encounter.Resource.Subject.ID,
			Reference: subjectReference,
			Display:   encounter.Resource.Subject.Display,
		},
		Created: &today,
		Custodian: &domain.FHIRReference{
			Reference: &orgRef,
			Display:   input.FacilityID,
		},
	}

	if input.Data.Notes != "" {
		carePlan.Note = []domain.FHIRAnnotation{
			{
				AuthorReference: &domain.FHIRReference{
					ID:        &input.FacilityID,
					Reference: &orgRef,
					Display:   input.FacilityID,
				},
				Text: (*scalarutils.Markdown)(&input.Data.Notes),
			},
		}
	}

	// Link the CarePlan to its Phase Tasks Using URNs.
	carePlanActivities := make([]domain.CarePlanActivity, 0)

	for phaseIndex, phase := range planDefinition.Action {
		// Recreate the same deterministic URN that prepareCarePlanActivityEntries created for the phase task.
		phaseTaskIdentifier := fmt.Sprintf("phase-%d-%s", phaseIndex+1, *planDefinition.ID)
		phaseTaskFullURL := generateFullURL(phaseTaskIdentifier, "Task")

		activity := domain.CarePlanActivity{
			PerformedActivity: []domain.FHIRCodeableReference{
				{
					Reference: &domain.FHIRReference{
						Reference: &phaseTaskFullURL,
						Display:   *phase.Title,
					},
				},
			},
		}

		carePlanActivities = append(carePlanActivities, activity)
	}

	carePlan.Activity = carePlanActivities

	carePlanFullURL := "urn:uuid:" + uuid.New().String()

	carePlanEntry := makeEntry(carePlanFullURL, "CarePlan", carePlan)
	if carePlanEntry.Resource == nil {
		return nil, fmt.Errorf("failed to create bundle entry for CarePlan")
	}

	finalBundleEntries := make([]hapifhirmodels.BundleEntry, 0, 1+len(allTaskEntries))
	finalBundleEntries = append(finalBundleEntries, carePlanEntry)
	finalBundleEntries = append(finalBundleEntries, allTaskEntries...)

	finalBundle := &hapifhirmodels.Bundle{
		Type:  hapifhirmodels.BundleTypeTransaction,
		Entry: finalBundleEntries,
	}

	bundlepayload, _ := json.MarshalIndent(finalBundle, "", " ")
	fmt.Println(string(bundlepayload))

	_, err = c.FHIR.PostFHIRBundle(ctx, finalBundle)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// prepareCarePlanActivityEntries orchestrates the creation of all phase and cycle Tasks in memory.
// It returns a slice of all prepared BundleEntry objects, ready for inclusion in the final transaction bundle.
func (c *ClinicalImpl) prepareCarePlanActivityEntries(
	ctx context.Context,
	planDefinition *domain.FHIRPlanDefinition,
	tags []domain.FHIRCodingInput,
	encounter *domain.FHIREncounter,
	facilityID string,
) ([]hapifhirmodels.BundleEntry, error) {
	var allTaskEntries []hapifhirmodels.BundleEntry

	if len(planDefinition.Action) > 0 {
		dateTime := scalarutils.DateTime(time.Now().Format(time.RFC3339))
		subjectReference := encounter.Subject.Reference
		encounterReference := fmt.Sprintf("Encounter/%s", *encounter.ID)
		orgRef := fmt.Sprintf("Organization/%s", facilityID)
		planDefinitionRef := fmt.Sprintf("PlanDefinition/%s", *planDefinition.ID)
		taskReadyStatus := "ready"
		taskIntent := "plan"
		taskPriority := "urgent"
		taskTypeSystem := "http://hl7.org/fhir/ValueSet/task-code"

		for phaseIndex, phase := range planDefinition.Action {
			phaseTaskInput := &domain.FHIRTaskInput{
				ResourceType: "Task",
				Focus: &domain.FHIRReferenceInput{
					Reference: &planDefinitionRef,
					Display:   *planDefinition.Title,
				},
				Meta: &domain.FHIRMetaInput{
					Tag: tags,
				},
				Status: (*scalarutils.Code)(&taskReadyStatus),
				BusinessStatus: &domain.FHIRCodeableConceptInput{
					Text: *phase.Title,
				},
				Intent:   (*scalarutils.Code)(&taskIntent),
				Priority: (*scalarutils.Code)(&taskPriority),
				Code: &domain.FHIRCodeableConceptInput{
					Coding: []*domain.FHIRCodingInput{
						{
							System:  (*scalarutils.URI)(&taskTypeSystem),
							Code:    "fulfill",
							Display: "Fulfill the focal request",
						},
					},
					Text: "Fulfill the focal request",
				},
				Description: phase.Description,
				For: &domain.FHIRReferenceInput{
					ID:        encounter.Subject.ID,
					Reference: subjectReference,
					Display:   encounter.Subject.Display,
				},
				Encounter: &domain.FHIRReferenceInput{
					ID:        encounter.ID,
					Reference: &encounterReference,
					Display:   *encounter.ID,
				},
				AuthoredOn: &dateTime,
				Requester: &domain.FHIRReferenceInput{
					ID:        &facilityID,
					Reference: &orgRef,
					Display:   facilityID,
				},
				Owner: &domain.FHIRReferenceInput{
					Reference: &orgRef,
					Display:   facilityID,
				},
				Note: []*domain.FHIRAnnotationInput{},
				RequestedPerformer: []*domain.FHIRCodeableReference{
					{
						Reference: &domain.FHIRReference{
							Reference: &orgRef,
							Display:   facilityID,
						},
					},
				},
			}

			phaseTaskIdentifier := fmt.Sprintf("phase-%d-%s", phaseIndex+1, *planDefinition.ID)
			phaseTaskFullURL := generateFullURL(phaseTaskIdentifier, "Task")

			phaseTaskEntry := makeEntry(phaseTaskFullURL, "Task", phaseTaskInput)
			if phaseTaskEntry.Resource == nil {
				err := fmt.Errorf("failed to create bundle entry for phase task %d", phaseIndex+1)
				return nil, err
			}

			allTaskEntries = append(allTaskEntries, phaseTaskEntry)

			if phase.TimingTiming != nil && phase.TimingTiming.Repeat != nil && phase.TimingTiming.Repeat.Count != nil && *phase.TimingTiming.Repeat.Count > 0 {
				cycleTaskEntries, err := c.prepareChemoPhaseCycleEntries(
					ctx,
					phaseTaskFullURL,
					phaseTaskInput,
					*phase.TimingTiming.Repeat.Count,
					phase.Action,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to prepare cycle tasks for phase '%s': %w", *phase.Title, err)
				}

				allTaskEntries = append(allTaskEntries, cycleTaskEntries...)
			}
		}
	}

	return allTaskEntries, nil
}

// prepareChemoPhaseCycleEntries prepares a slice of FHIR BundleEntry objects for each cycle task.
// It returns the prepared entries, ready to be added to a larger transaction bundle.
func (c *ClinicalImpl) prepareChemoPhaseCycleEntries(
	_ context.Context,
	parentTaskIDURN string,
	task *domain.FHIRTaskInput,
	cycleCount int,
	cycleAction []domain.PlanDefinitionAction,
) ([]hapifhirmodels.BundleEntry, error) {
	var cycleTaskEntries []hapifhirmodels.BundleEntry

	for i := 0; i < cycleCount; i++ {
		cycleTask := *task

		cycleDescription := fmt.Sprintf("Cycle %d of %d", i+1, cycleCount)

		cycleTask.Description = &cycleDescription

		cycleTask.PartOf = []domain.FHIRReference{
			{
				Reference: &parentTaskIDURN,
			},
		}

		if len(cycleAction) > 0 {
			generateTaskInput := func(actions []domain.PlanDefinitionAction) string {
				for _, act := range actions {
					if act.Title != nil {
						return *act.Title
					}
				}

				return ""
			}

			for _, act := range cycleAction {
				if act.Title != nil {
					cycleTask.BusinessStatus.Text = *act.Title
				}

				cycleTask.Input = append(cycleTask.Input, &domain.TaskInput{
					Type: &domain.FHIRCodeableConcept{
						Text: generateTaskInput(act.Action),
					},
					ValueCodeableReference: &domain.FHIRCodeableReference{
						Reference: &domain.FHIRReference{
							Reference: &parentTaskIDURN,
						},
					},
				})
			}
		}

		cycleTaskIdentifier := parentTaskIDURN + "-cycle-" + strconv.Itoa(i+1)

		cycleTaskFullURL := generateFullURL(cycleTaskIdentifier, "Task")

		entry := makeEntry(cycleTaskFullURL, "Task", &cycleTask)
		if entry.Resource == nil {
			err := fmt.Errorf("failed to create bundle entry for cycle task %d, could not marshal resource", i+1)
			return nil, err
		}

		cycleTaskEntries = append(cycleTaskEntries, entry)
	}

	return cycleTaskEntries, nil
}

// generateFullURL returns a full URN for use within a bundle.
//
// We use this to show the relationships between the resources in the bundle.
func generateFullURL(identifier, resourceType string) string {
	return "urn:uuid:" + resourceType + "-" + identifier
}

// makeEntry creates a standard POST entry for a transaction bundle.
func makeEntry(
	fullURL string,
	resourceType string,
	resource interface{},
) hapifhirmodels.BundleEntry {
	rawResource, err := json.Marshal(resource)
	if err != nil {
		return hapifhirmodels.BundleEntry{}
	}

	return hapifhirmodels.BundleEntry{
		FullURL:  (*hapifhirmodels.BundleEntryFullURL)(&fullURL),
		Resource: rawResource,
		Request: &hapifhirmodels.BundleEntryRequest{
			Method: hapifhirmodels.HTTPVerbPOST,
			URL:    resourceType,
		},
	}
}

// CreateCarePlan is used to create a care plan
func (c *ClinicalImpl) CreateCarePlan(ctx context.Context, payload *domain.FHIRCarePlan) (*domain.FHIRCarePlan, error) {
	return c.FHIR.CreateFHIRCarePlan(ctx, payload)
}

func (c *ClinicalImpl) SearchCarePlan(ctx context.Context, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error) {
	taskSearchParams := map[string]interface{}{
		"_sort":  "_lastUpdated",
		"_total": "accurate",
	}

	results, err := c.FHIR.SearchFHIRCarePlan(ctx, taskSearchParams, pagination)
	if err != nil {
		return nil, err
	}

	return results, nil
}
