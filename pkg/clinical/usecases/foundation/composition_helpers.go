package foundation

import (
	"context"
	"fmt"

	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// CompositionConcept is used to map composition concepts
type CompositionConcept struct {
	CompositionCategoryConcept *domain.Concept
	CompositionTypeConcept     *domain.Concept
}

// mapCompositionConcepts composes a unified representation of composition concepts and types into a single model
func (c *FoundationImpl) mapCompositionConcepts(ctx context.Context, compositionCategoryCode, conceptID string) (*CompositionConcept, error) {
	compositionCategoryConcept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, compositionCategoryCode)
	if err != nil {
		return nil, err
	}

	compositionTypeConcept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, conceptID)
	if err != nil {
		return nil, err
	}

	return &CompositionConcept{
		CompositionCategoryConcept: compositionCategoryConcept,
		CompositionTypeConcept:     compositionTypeConcept,
	}, nil
}

// mapCategoryEnumToCode is used to map various composition categories to respective LOINC codes
func (*FoundationImpl) mapCategoryEnumToCode(category dto.CompositionCategory) (string, error) {
	var compositionCategoryCode string

	switch category {
	case "ASSESSMENT_PLAN":
		compositionCategoryCode = common.LOINCAssessmentPlanCode
	case "HISTORY_OF_PRESENTING_ILLNESS":
		compositionCategoryCode = common.LOINCHistoryOfPresentingIllness
	case "SOCIAL_HISTORY":
		compositionCategoryCode = common.LOINCSocialHistory
	case "FAMILY_HISTORY":
		compositionCategoryCode = common.LOINCFamilyHistory
	case "EXAMINATION":
		compositionCategoryCode = common.LOINCExamination
	case "PLAN_OF_CARE":
		compositionCategoryCode = common.LOINCPlanOfCare
	case "PROGRESS_NOTE":
		compositionCategoryCode = common.LOINCProviderUnspecifiedProgressNote
	case "PAST_MEDICAL_SURGERY_HISTORY":
		compositionCategoryCode = common.LOINCPastSurgeryHistory
	case "CHIEF_COMPLAINT":
		compositionCategoryCode = common.LOINCChiefComplaint
	default:
		return "", fmt.Errorf("category is needed")
	}

	return compositionCategoryCode, nil
}
