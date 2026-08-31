package dto

import (
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// Consent models a fhir consent resource.
type Consent struct {
	ID           *string                    `json:"id,omitempty"`
	Status       *domain.ConsentStatusEnum  `json:"status"`
	Decision     domain.ConsentDecisionEnum `json:"decision,omitempty"`
	Patient      *Reference                 `json:"patient,omitempty"`
	DateTime     *scalarutils.DateTime      `json:"dateTime,omitempty" swaggertype:"primitive,string"`
	UsageContext ScreeningTypeEnum          `json:"usageContext,omitempty"`
}
