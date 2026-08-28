package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// ListMedicationRequest(s) is used to retrieve a paginated list of available medication requests given some parameters.
//
// If none are provided, it returns all the available medication requests
func (c *ClinicalImpl) ListMedicationRequests(ctx context.Context, filter *dto.MedicationRequestFilterInput, pagination dto.Pagination) (*dto.MedicationRequestConnection, error) {
	medicationRequestSearchParams := map[string]interface{}{
		"_sort":  "_lastUpdated",
		"_total": "accurate",
	}

	err := filter.MedicationRequestFilters(medicationRequestSearchParams)
	if err != nil {
		return nil, err
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, err
	}

	results, err := c.FHIR.SearchFHIRMedicationRequest(ctx, medicationRequestSearchParams, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	medicationRequestList := []*dto.MedicationRequestOutput{}

	for _, medicationRequestItem := range results.MedicationRequests {
		medicationRequestList = append(medicationRequestList, mapFHIRMedicationRequestToMedicationRequestDTO(medicationRequestItem))
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     results.HasNextPage,
		EndCursor:       &results.NextCursor,
		HasPreviousPage: results.HasPreviousPage,
		StartCursor:     &results.PreviousCursor,
	}

	connection := dto.CreateMedicationRequestConnection(medicationRequestList, pageInfo, results.TotalCount)

	return &connection, nil
}

// CreatePrescription creates a prescription entry as a FHIR MedicationRequest resource
func (c *ClinicalImpl) CreatePrescription(ctx context.Context, input dto.PrescriptionInput) ([]*dto.MedicationRequestOutput, error) {
	if err := helpers.Validate(input); err != nil {
		return nil, err
	}

	var output []*dto.MedicationRequestOutput

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record an observation in a finished encounter")
	}

	for _, medicationInput := range input.Medications {
		medication, err := c.FHIR.FetchMedicationByID(ctx, medicationInput.MedicationID)
		if err != nil {
			return nil, err
		}

		var (
			patientReference   = fmt.Sprintf("Patient/%s", *encounter.Resource.Subject.ID)
			encounterReference = fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
			status             = domain.ActiveMedicationStatus
			medicationRef      = fmt.Sprintf("Medication/%s", *medication.ID)
		)

		dosageInstructions := []*domain.FHIRDosageInput{}

		for _, dose := range medicationInput.DosageInstructions {
			dosageInstruction := domain.FHIRDosageInput{
				PatientInstruction: &dose.PatientInstruction,
				AsNeeded:           &dose.AsNeeded,
				Route: &domain.FHIRCodeableConceptInput{
					Text: dose.Route.Code,
					Coding: []*domain.FHIRCodingInput{
						{
							System:  helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier),
							Code:    scalarutils.Code(dose.Route.Code),
							Display: dose.Route.Display,
						},
					},
				},
				DoseAndRate: []*domain.FHIRDosageDoseandrateInput{
					{
						DoseQuantity: &domain.FHIRQuantityInput{
							Value:  dose.DoseQuantity,
							Unit:   dose.DoseUnit,
							System: (scalarutils.URI)(helpers.QuantitySystem),
							Code:   (*scalarutils.Code)(&dose.DoseUnit),
						},
					},
				},
				Timing: &domain.FHIRTimingInput{
					Repeat: &domain.FHIRTimingRepeatInput{
						BoundsPeriod: &domain.FHIRPeriodInput{
							Start: dose.StartDate,
							End:   dose.EndDate,
						},
						Duration:     (*json.Number)(&dose.Duration),
						DurationUnit: dose.DurationUnit,
						Frequency:    &dose.Frequency,
						Period:       (*json.Number)(&dose.Period),
						PeriodUnit:   dose.PeriodUnit,
					},
				},
			}

			dosageInstructions = append(dosageInstructions, &dosageInstruction)
		}

		today := time.Now().Format(time.RFC3339)
		authoredOn := (*scalarutils.DateTime)(&today)

		medicationRequest := domain.FHIRMedicationRequestInput{
			Status: &status,
			Intent: domain.PlanMedicationIntent,
			Subject: &domain.FHIRReferenceInput{
				ID:        encounter.Resource.Subject.ID,
				Reference: &patientReference,
				Display:   encounter.Resource.Subject.Display,
			},
			Encounter: &domain.FHIRReferenceInput{
				ID:        encounter.Resource.ID,
				Reference: &encounterReference,
				Display:   *encounter.Resource.ID,
			},
			Priority:   (*scalarutils.Code)(&medicationInput.Priority),
			AuthoredOn: authoredOn,
			Requester: &domain.FHIRReferenceInput{
				Reference: encounter.Resource.ServiceProvider.Reference,
				ID:        &encounter.Resource.ServiceProvider.Display,
				Display:   encounter.Resource.ServiceProvider.Display,
			},
			Medication: &domain.FHIRCodeableReferenceInput{
				Reference: &domain.FHIRReferenceInput{
					ID:        medication.ID,
					Reference: &medicationRef,
					Display:   *medication.ID,
				},
			},
			DosageInstruction: dosageInstructions,
		}

		tags, err := c.GetTenantMetaTags(ctx)
		if err != nil {
			return nil, err
		}

		medicationRequest.Meta = domain.FHIRMetaInput{
			Tag: tags,
		}

		medRequest, err := c.FHIR.CreateFHIRMedicationRequest(ctx, medicationRequest)
		if err != nil {
			return nil, err
		}

		output = append(output, mapFHIRMedicationRequestToMedicationRequestDTO(*medRequest.Resource))
	}

	return output, nil
}

func mapFHIRMedicationRequestToMedicationRequestDTO(medicationRequest domain.FHIRMedicationRequest) *dto.MedicationRequestOutput {
	authoredOn := (*scalarutils.DateTime)(medicationRequest.AuthoredOn)
	output := &dto.MedicationRequestOutput{
		ID:           *medicationRequest.ID,
		EncounterID:  medicationRequest.Encounter.Display,
		Status:       domain.MedicationRequestStatus(*medicationRequest.Status),
		AuthoredOn:   authoredOn,
		Diagnosis:    "",
		Priority:     string(*medicationRequest.Priority),
		OrderedBy:    medicationRequest.Meta.GetFacilityName(), //TODO: Use facility for now, SHOULD be a practitioner
		FacilityName: medicationRequest.Meta.GetFacilityName(),
	}

	if medicationRequest.Subject != nil {
		output.Subject = &dto.Reference{
			ID:      medicationRequest.Subject.Display,
			Display: medicationRequest.Subject.Display,
		}
	}

	if medicationRequest.Encounter != nil {
		output.EncounterID = medicationRequest.Encounter.Display
	}

	if len(medicationRequest.Note) > 0 {
		for _, note := range medicationRequest.Note {
			output.Notes = append(output.Notes, &dto.Annotation{
				Text: *note.Text,
			})
		}
	}

	if len(medicationRequest.DosageInstruction) > 0 {
		for i, dosageInstruction := range medicationRequest.DosageInstruction {
			output.DosageInstructions = append(output.DosageInstructions, dto.DosageInstruction{
				Route: dto.ValueSetData{
					Display: dosageInstruction.Route.Text,
				},
				DoseQuantity:       dosageInstruction.DoseAndRate[i].DoseQuantity.Value,
				DoseUnit:           dosageInstruction.DoseAndRate[i].DoseQuantity.Unit,
				Period:             dosageInstruction.Timing.Repeat.Period.String(),
				PeriodUnit:         dosageInstruction.Timing.Repeat.PeriodUnit,
				Frequency:          *dosageInstruction.Timing.Repeat.Frequency,
				Duration:           dosageInstruction.Timing.Repeat.Duration.String(),
				DurationUnit:       dosageInstruction.Timing.Repeat.DurationUnit,
				StartDate:          &dosageInstruction.Timing.Repeat.BoundsPeriod.Start,
				EndDate:            &dosageInstruction.Timing.Repeat.BoundsPeriod.End,
				PatientInstruction: *dosageInstruction.PatientInstruction,
				AsNeeded:           *dosageInstruction.AsNeeded,
			})
		}
	}

	return output
}

// PatchMedicationRequests update a medication request resource
func (c *ClinicalImpl) PatchMedicationRequests(ctx context.Context, id string, value domain.MedicationRequestStatus) (*dto.MedicationRequestOutput, error) {
	if id == "" {
		return nil, fmt.Errorf("medication request ID is required")
	}

	medicationRequestID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid medication request id: %s", id)
	}

	// TODO: For starters, we will only patch the status. As the work evolves, consider using a
	// data class(struct) for the patch payload
	if !value.IsValid() {
		return nil, fmt.Errorf("invalid medication request status: %v", medicationRequestID)
	}

	observationInput := &domain.FHIRMedicationRequestInput{
		Status: &value,
	}

	output, err := c.FHIR.PatchFHIRMedicationRequest(ctx, id, *observationInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRMedicationRequestToMedicationRequestDTO(*output), nil
}

// FetchMedicationRequestByID is used to retrieve the details a medication request givens its ID
func (c *ClinicalImpl) FetchMedicationRequestByID(ctx context.Context, medicationRequestID string) (*dto.MedicationRequestOutput, error) {
	if medicationRequestID == "" {
		return nil, fmt.Errorf("medication request ID is required")
	}

	output, err := c.FHIR.GetFHIRMedicationRequest(ctx, medicationRequestID)
	if err != nil {
		return nil, err
	}

	return mapFHIRMedicationRequestToMedicationRequestDTO(*output.Resource), nil
}

func (c *ClinicalImpl) FetchMedicationByID(ctx context.Context, id string) (*dto.MedicationOutput, error) {
	medication, err := c.FHIR.FetchMedicationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	code, display := medication.GetDoseForm()
	output := &dto.MedicationOutput{
		ID:   *medication.ID,
		Name: medication.GetText(),
		DoseForm: dto.ValueSetData{
			Code:    code,
			Display: display,
		},
	}

	if medication.Batch != nil && medication.Batch.LotNumber != nil {
		output.LotNumber = *medication.Batch.LotNumber
	}

	if medication.Status != nil {
		output.Status = medication.Status.String()
	}

	if medication.Batch != nil && medication.Batch.ExpirationDate != nil {
		output.ExpiryDate = (*scalarutils.DateTime)(medication.Batch.ExpirationDate)
	}

	return output, nil
}
