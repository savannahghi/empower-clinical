package clinical

import (
	"context"
	"fmt"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

func (c *ClinicalImpl) RecordMedication(ctx context.Context, medications []*dto.MedicationInput) ([]*dto.MedicationOutput, error) {
	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	var medicationsOutput []*dto.MedicationOutput

	for _, medication := range medications {
		payload := &domain.FHIRMedicationInput{
			Extension: []*domain.Extension{
				{
					URL:         fmt.Sprintf("%s/ValueSet/medication-codes", serverutils.MustGetEnvVar("HAPI_FHIR_BASE_URL")),
					ValueString: medication.Name,
				},
			},
			Code: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier),
						Code:    (scalarutils.Code("sghidefaultcode")),
						Display: "SGHI Default Code",
					},
				},
			},
			DoseForm: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  helpers.CodeSystem(common.UnspecifiedCodeSystemIdentifier),
						Code:    scalarutils.Code(medication.DoseForm.Code),
						Display: medication.DoseForm.Display,
					},
				},
			},
			Status: domain.MedicationStatusEnumActive,
		}

		payload.Meta = domain.FHIRMetaInput{
			Tag: tags,
		}

		medication, err := c.FHIR.CreateFHIRMedication(ctx, *payload)
		if err != nil {
			return nil, fmt.Errorf("an error occurred: %w", err)
		}

		code, display := medication.Resource.GetDoseForm()

		output := &dto.MedicationOutput{
			ID:   *medication.Resource.ID,
			Name: medication.Resource.GetText(),
			DoseForm: dto.ValueSetData{
				Code:    code,
				Display: display,
			},
			Status: medication.Resource.Status.String(),
		}

		if medication.Resource.Batch != nil && medication.Resource.Batch.LotNumber != nil {
			output.LotNumber = *medication.Resource.Batch.LotNumber
		}

		if medication.Resource.Batch != nil && medication.Resource.Batch.ExpirationDate != nil {
			output.ExpiryDate = (*scalarutils.DateTime)(medication.Resource.Batch.ExpirationDate)
		}

		medicationsOutput = append(medicationsOutput, output)
	}

	return medicationsOutput, nil
}
