package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// CreateQuestionnaireResponse creates a questionnaire response
func (c *ClinicalImpl) CreateQuestionnaireResponse(ctx context.Context, questionnaireID string, encounterID string, input dto.QuestionnaireResponse) (*dto.QuestionnaireReviewSummary, error) {
	questionnaireResponse := &domain.FHIRQuestionnaireResponse{}

	err := mapstructure.Decode(input, questionnaireResponse)
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot create a questionnaire response in a completed encounter")
	}

	// TODO: Ensure user cannot submit the same risk assessment twice in the same encounter
	questionnaireResponse.Source = &domain.FHIRReference{
		Reference: encounter.Resource.Subject.Reference,
		Display:   encounter.Resource.Subject.Display,
	}

	if encounter.Resource.Subject.ID != nil {
		questionnaireResponse.Source.ID = encounter.Resource.Subject.ID
	} else {
		questionnaireResponse.Source.ID = &encounter.Resource.Subject.Display
	}

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	questionnaireResponse.Encounter = &domain.FHIRReference{
		ID:        encounter.Resource.ID,
		Reference: &encounterReference,
		Display:   *encounter.Resource.ID,
	}

	questionnaireResponse.Questionnaire = &questionnaireID

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	questionnaireResponse.Meta.Tag = tags
	questionnaireResponse.Authored = &input.Authored
	questionnaireResponse.Status = input.Status

	resp, err := c.FHIR.CreateFHIRQuestionnaireResponse(ctx, questionnaireResponse)
	if err != nil {
		return nil, err
	}

	output := &dto.QuestionnaireResponse{}

	err = mapstructure.Decode(resp, output)
	if err != nil {
		return nil, err
	}

	// TODO: This will affect the API performance. Optimize it
	riskLevel, err := c.generateQuestionnaireReviewSummary(
		ctx,
		questionnaireID,
		*resp.ID,
		encounter,
		output,
	)
	if err != nil {
		return nil, err
	}

	return &dto.QuestionnaireReviewSummary{
		RiskLevel:               riskLevel,
		QuestionnaireResponseID: *resp.ID,
	}, nil
}

// generateQuestionnaireReviewSummary takes a questionnaire response and
// analyzes it to determine the risk stratification based on three distinct groups:
// symptoms, family history, and risk factors. The assumption is that the
// questionnaire has groups with linkIds: symptoms, family_history, and risk-factors.
// The function looks into the responses saved under the tags <group_name>-score,
// calculates the total scores for each group, and returns a summary indicating
// whether the individual is high risk, low risk, or average risk.
func (c *ClinicalImpl) generateQuestionnaireReviewSummary(
	ctx context.Context,
	questionnaireID,
	questionnaireResponseID string,
	encounter *domain.FHIREncounterRelayPayload,
	questionnaireResponse *dto.QuestionnaireResponse,
) (string, error) {
	riskLevel := ""

	questionnaire, err := c.FHIR.GetFHIRQuestionnaire(ctx, questionnaireID)
	if err != nil {
		return "", err
	}

	// Some reference resources do not have the ID field,
	// but have the ID under the Dispaly field
	var id *string

	if encounter.Resource.Subject.ID != nil {
		id = encounter.Resource.Subject.ID
	} else {
		id = &encounter.Resource.Subject.Display
	}

	patient, err := c.FHIR.GetFHIRPatient(ctx, *id)
	if err != nil {
		return "", err
	}

	switch *questionnaire.Resource.Title {
	case dto.CervicalCancerScreeningTypeEnum.Text():
		scores := helpers.CalculateScoresFromResponses(questionnaire.Resource, questionnaireResponse)

		totalScore := scores["symptoms"] + scores["risk-factors"]

		switch {
		case totalScore >= 2:
			riskLevel, err = c.recordAssessmentAndSegmentPatient(
				ctx,
				*patient.Resource.ID,
				encounter,
				questionnaireResponseID,
				string(domain.HighRiskProbability),
				domain.HighRiskProbability.Display(),
				dto.CervicalCancerScreeningTypeEnum.Text(),
				dto.SegmentationCategoryHighRiskNegative,
			)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return "", err
			}

		case totalScore == 1:
			riskLevel, err = c.recordAssessmentAndSegmentPatient(
				ctx,
				*patient.Resource.ID,
				encounter,
				questionnaireResponseID,
				string(domain.LowRiskRiskProbability),
				domain.LowRiskRiskProbability.Display(),
				dto.CervicalCancerScreeningTypeEnum.Text(),
				dto.SegmentationCategoryLowRisk,
			)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return "", err
			}

		case totalScore == 0:
			riskLevel, err = c.recordAssessmentAndSegmentPatient(
				ctx,
				*patient.Resource.ID,
				encounter,
				questionnaireResponseID,
				string(domain.NegligibleRiskProbability),
				domain.NegligibleRiskProbability.Display(),
				dto.CervicalCancerScreeningTypeEnum.Text(),
				dto.SegmentationCategoryNoRisk,
			)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return "", err
			}
		}

	case dto.BreastCancerScreeningTypeEnum.Text():
		scores := helpers.CalculateScoresFromResponses(questionnaire.Resource, questionnaireResponse)

		if scores["risk-assessment"] >= 1 {
			riskLevel, err = c.recordAssessmentAndSegmentPatient(
				ctx,
				*patient.Resource.ID,
				encounter,
				questionnaireResponseID,
				string(domain.HighRiskProbability),
				domain.HighRiskProbability.Display(),
				dto.BreastCancerScreeningTypeEnum.Text(),
				dto.SegmentationBreastCategoryHighRisk,
			)
			if err != nil {
				return "", err
			}
		} else {
			riskLevel, err = c.recordAssessmentAndSegmentPatient(
				ctx,
				*patient.Resource.ID,
				encounter,
				questionnaireResponseID,
				string(domain.ModerateRiskProbability),
				domain.ModerateRiskProbability.Display(),
				dto.BreastCancerScreeningTypeEnum.Text(),
				dto.SegmentationBreastCategoryAverageRisk,
			)
			if err != nil {
				return "", err
			}
		}

	case dto.ProstateCancerScreeningTypeEnum.Text():
		scores := helpers.CalculateScoresFromResponses(questionnaire.Resource, questionnaireResponse)

		if scores["risk-assessment"] >= 1 {
			riskLevel, err = c.recordAssessmentAndSegmentPatient(
				ctx,
				*patient.Resource.ID,
				encounter,
				questionnaireResponseID,
				string(domain.HighRiskProbability),
				domain.HighRiskProbability.Display(),
				dto.ProstateCancerScreeningTypeEnum.Text(),
				dto.SegmentationBreastCategoryHighRisk,
			)
			if err != nil {
				return "", err
			}
		} else {
			riskLevel, err = c.recordAssessmentAndSegmentPatient(
				ctx,
				*patient.Resource.ID,
				encounter,
				questionnaireResponseID,
				string(domain.ModerateRiskProbability),
				domain.ModerateRiskProbability.Display(),
				dto.ProstateCancerScreeningTypeEnum.Text(),
				dto.SegmentationBreastCategoryAverageRisk,
			)
			if err != nil {
				return "", err
			}
		}

	default:
		return "", fmt.Errorf("questionnaire does not exist")
	}

	return riskLevel, nil
}

// recordAssessmentAndSegmentPatient records a cancer screening risk assessment for a patient,
// then segments the patient for targeted engagement based on the assessment's risk level.
func (c *ClinicalImpl) recordAssessmentAndSegmentPatient(
	ctx context.Context,
	patientID string,
	encounter *domain.FHIREncounterRelayPayload,
	questionnaireResponseID, outcomeCode string,
	outcomeDisplay, usageContext string,
	segmentLabel dto.SegmentationCategory,
) (string, error) {
	riskLevel, err := c.recordRiskAssessment(
		ctx,
		encounter,
		questionnaireResponseID,
		outcomeCode,
		outcomeDisplay,
		usageContext,
	)
	if err != nil {
		return "", err
	}

	err = c.Pubsub.NotifySegmentation(ctx, dto.SegmentationPayload{
		ClinicalID:   patientID,
		SegmentLabel: segmentLabel,
	})
	if err != nil {
		return "", err
	}

	return riskLevel, nil
}

func (c *ClinicalImpl) recordRiskAssessment(
	ctx context.Context,
	encounter *domain.FHIREncounterRelayPayload,
	questionnaireResponseID, outcomeCode string,
	outcomeDisplay, usageContext string,
) (string, error) {
	ProbabilityRiskAssessmentSystem := scalarutils.URI(common.RiskAssessmentCodeSystem)

	codingCode := scalarutils.Code(outcomeCode)

	risk := domain.FHIRCodeableConcept{
		Coding: []*domain.FHIRCoding{
			{
				System:  &ProbabilityRiskAssessmentSystem,
				Code:    &codingCode,
				Display: outcomeDisplay,
			},
		},
		Text: outcomeDisplay,
	}

	patientReference := encounter.Resource.Subject.Reference

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	questionnaireResponseReference := fmt.Sprintf("QuestionnaireResponse/%s", questionnaireResponseID)

	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	div := fmt.Sprintf("<div xmlns=\"http://www.w3.org/1999/xhtml\">%s</div>", usageContext)

	textStatus := domain.NarrativeStatusEnumAdditional
	riskAssessment := domain.FHIRRiskAssessmentInput{
		Status: domain.ObservationStatusEnumFinal,
		Subject: domain.FHIRReferenceInput{
			Reference: patientReference,
			Display:   encounter.Resource.Subject.Display,
		},

		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Reference: &encounterReference,
			Display:   *encounter.Resource.ID,
		},
		OccurrenceDateTime: (*string)(&instant),
		Prediction: []domain.FHIRRiskAssessmentPrediction{
			{
				QualitativeRisk: &risk,
			},
		},
		Basis: []domain.FHIRReferenceInput{
			{
				ID:        &questionnaireResponseID,
				Reference: &questionnaireResponseReference,
				Display:   questionnaireResponseID,
			},
		},
		Text: &domain.FHIRNarrativeInput{
			Status: &textStatus,
			Div:    scalarutils.XHTML(div),
		},
	}

	if encounter.Resource.Subject.ID != nil {
		riskAssessment.Subject.ID = encounter.Resource.Subject.ID
	} else {
		riskAssessment.Subject.ID = &encounter.Resource.Subject.Display
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return "", err
	}

	riskAssessment.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	assessment, err := c.RecordRiskAssessment(ctx, &riskAssessment)
	if err != nil {
		return "", err
	}

	riskLevel := assessment.Prediction[0].QualitativeRisk.Text

	return riskLevel, nil
}

// SimpleQuestionnaireResponse handles the request to retrieve a simplified version
// of a QuestionnaireResponse. This endpoint processes the full questionnaire response data
// and returns a minimal, streamlined version that includes only essential information.
func (c *ClinicalImpl) SimpleQuestionnaireResponse(ctx context.Context, questionnaireResponseID string) ([]*domain.SimpleQuestionnaireResponse, error) {
	if questionnaireResponseID == "" {
		return nil, fmt.Errorf("questionnaireResponseID is required")
	}

	questionnaireResponse, err := c.FHIR.GetFHIRQuestionnaireResponse(ctx, questionnaireResponseID)
	if err != nil {
		return nil, err
	}

	responses := questionnaireResponse.Resource.GetFHIRQuestionnaireResponse()

	return responses, nil
}
