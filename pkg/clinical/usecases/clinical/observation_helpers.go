package clinical

import (
	"context"
	"strings"
	"time"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

var (
	observationCategorySystem = "http://terminology.hl7.org/CodeSystem/observation-category"
)

// addObservationEffectiveDate overrides the observation's effective time with the
// supplied date (when the test was performed). When no date is provided the
// observation keeps the default effective time set by RecordObservation.
func addObservationEffectiveDate(date *scalarutils.Date) ObservationInputMutatorFunc {
	return func(ctx context.Context, observation *domain.FHIRObservationInput) error {
		if date == nil {
			return nil
		}

		instant := scalarutils.Instant(date.AsTime().Format(time.RFC3339))
		observation.EffectiveInstant = &instant

		return nil
	}
}

func addObservationCategory(code string) ObservationInputMutatorFunc {
	var display string

	switch strings.TrimSpace(code) {
	case "exam":
		display = "Exam"
	case "procedure":
		display = "Procedure"
	case "imaging":
		display = "Imaging"
	case "laboratory":
		display = "Laboratory"
	case "vital-signs":
		display = "Vital Signs"
	case "social-history":
		display = "Social History"
	}

	return func(ctx context.Context, observation *domain.FHIRObservationInput) error {
		userSelected := false
		category := []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       (*scalarutils.URI)(&observationCategorySystem),
						Code:         scalarutils.Code(code),
						Display:      display,
						UserSelected: &userSelected,
					},
				},
				Text: display,
			},
		}

		observation.Category = append(observation.Category, category...)

		return nil
	}
}
