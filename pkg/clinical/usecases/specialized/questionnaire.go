package specialized

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/foundation"
)

type SpecializedImpl struct {
	infrastructure.Infrastructure
	foundation.FoundationImpl
	*slog.Logger
}

func NewSpecializedImpl(infra infrastructure.Infrastructure, foundation foundation.FoundationImpl, log *slog.Logger) *SpecializedImpl {
	return &SpecializedImpl{
		infra,
		foundation,
		log,
	}
}

// CreateQuestionnaire is used to create a new Questionnaire.
// These questionnaire are used to solicit various types of information from patients to server organisation usecases.
func (c *SpecializedImpl) CreateQuestionnaire(ctx context.Context, questionnaireInput *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
	questionnaire, err := c.FHIR.CreateFHIRQuestionnaire(ctx, questionnaireInput)
	if err != nil {
		return nil, err
	}

	return questionnaire, nil
}

// ListQuestionnaires is used to list questionnaires from FHIR repository.
// This search is performed using the name or the title of the questionnaire and returns the available questionnaire(s).
func (c *SpecializedImpl) ListQuestionnaires(ctx context.Context, searchParam string, pagination *dto.Pagination) (*dto.Questionnaire, error) {
	params := map[string]interface{}{
		"status": "active",
		"_sort":  "-_lastUpdated",
		"_count": "1",
	}

	if searchParam != "" {
		params["title:exact"] = searchParam
	}

	questionnaire, err := c.FHIR.ListFHIRQuestionnaire(ctx, params, dto.TenantIdentifiers{}, *pagination)
	if err != nil {
		return nil, err
	}

	var dtoQuestionnaire *dto.Questionnaire

	for _, questionnaire := range questionnaire.Questionnaires {
		err := mapstructure.Decode(questionnaire, &dtoQuestionnaire)
		if err != nil {
			return nil, err
		}
	}

	return dtoQuestionnaire, nil
}

func (c *SpecializedImpl) FetchQuestionnaire(ctx context.Context, searchParam string, pagination *dto.Pagination) (*dto.Questionnaire, error) {
	params := map[string]interface{}{
		"status": "active",
		"_sort":  "-_lastUpdated",
		"_count": "1",
	}

	if searchParam != "" {
		params["title:exact"] = searchParam
	}

	questionnaire, err := c.FHIR.ListFHIRQuestionnaire(ctx, params, dto.TenantIdentifiers{}, *pagination)
	if err != nil {
		return nil, err
	}

	var dtoQuestionnaire *dto.Questionnaire

	for _, questionnaire := range questionnaire.Questionnaires {
		err := mapstructure.Decode(questionnaire, &dtoQuestionnaire)
		if err != nil {
			return nil, err
		}
	}

	return dtoQuestionnaire, nil
}

func (c *SpecializedImpl) GenerateRiskAssessment(ctx context.Context, questionnaireID string, questionnaireResponse *dto.QuestionnaireResponse) (*dto.RiskAssessmentResult, error) {
	questionnaire, err := c.FHIR.GetFHIRQuestionnaire(ctx, questionnaireID)
	if err != nil {
		return nil, err
	}

	switch *questionnaire.Resource.Title {
	case "Cervical Cancer Screening":
		scores := helpers.CalculateScoresFromResponses(questionnaire.Resource, questionnaireResponse)

		totalScore := scores["symptoms"] + scores["risk-factors"]

		switch {
		case totalScore >= 2:
			return &dto.RiskAssessmentResult{
				RiskLevel: string(domain.HighRiskProbability),
			}, nil
		case totalScore == 1:
			return &dto.RiskAssessmentResult{
				RiskLevel: string(domain.LowRiskRiskProbability),
			}, nil
		case totalScore == 0:
			return &dto.RiskAssessmentResult{
				RiskLevel: string(domain.NegligibleRiskProbability),
			}, nil
		}
	case "Breast Cancer Screening":
		scores := helpers.CalculateScoresFromResponses(questionnaire.Resource, questionnaireResponse)

		if scores["risk-assessment"] >= 1 {
			return &dto.RiskAssessmentResult{
				RiskLevel: string(domain.HighRiskProbability),
			}, nil
		} else {
			return &dto.RiskAssessmentResult{
				RiskLevel: string(domain.ModerateRiskProbability),
			}, nil
		}
	case "Prostate Cancer Screening":
		scores := helpers.CalculateScoresFromResponses(questionnaire.Resource, questionnaireResponse)

		if scores["risk-assessment"] >= 1 {
			return &dto.RiskAssessmentResult{
				RiskLevel: string(domain.HighRiskProbability),
			}, nil
		} else {
			return &dto.RiskAssessmentResult{
				RiskLevel: string(domain.ModerateRiskProbability),
			}, nil
		}
	default:
		return nil, fmt.Errorf("unknown questionnaire")
	}

	return nil, nil
}
