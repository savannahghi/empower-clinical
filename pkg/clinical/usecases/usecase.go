package usecases

import (
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/base"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/foundation"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/specialized"
)

type Usecases struct {
	foundation.FoundationImpl
	base.BaseImpl
	specialized.SpecializedImpl
	clinical.ClinicalImpl
}

func NewUsecasesImpl(
	foundation foundation.FoundationImpl,
	base base.BaseImpl,
	specialized *specialized.SpecializedImpl,
	clinical clinical.ClinicalImpl,
) (*Usecases, error) {
	return &Usecases{
		foundation,
		base,
		*specialized,
		clinical,
	}, nil
}
