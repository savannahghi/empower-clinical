package clinical

import (
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

func mapFHIRRiskAssessmentToRiskAssessmentDTO(riskAssessment domain.FHIRRiskAssessment) dto.RiskAssessment {
	occurrenceTime := (*scalarutils.DateTime)(riskAssessment.OccurrenceDateTime)
	output := &dto.RiskAssessment{
		ID: riskAssessment.ID,
		Subject: dto.Reference{
			Display: riskAssessment.Subject.Display,
			ID:      riskAssessment.Subject.ResourceID(),
		},
		Date: occurrenceTime,
	}

	if riskAssessment.Encounter != nil && riskAssessment.Encounter.ID != nil {
		output.Encounter = &dto.Reference{
			ID: *riskAssessment.Encounter.ID,
		}
	}

	if riskAssessment.Text != nil {
		output.UsageContext = helpers.ExtractTextFromHTML(string(riskAssessment.Text.Div))
	}

	var outcome string

	for _, prediction := range riskAssessment.Prediction {
		if prediction.QualitativeRisk != nil && prediction.QualitativeRisk.Text != "" {
			outcome = prediction.QualitativeRisk.Text
			break
		} else if prediction.Outcome != nil && prediction.Outcome.Text != "" {
			outcome = prediction.Outcome.Text
			break
		}
	}

	output.Prediction = append(output.Prediction, dto.RiskAssessmentPrediction{
		Outcome: &dto.CodeableConcept{
			Text: outcome,
		},
	})

	return *output
}
