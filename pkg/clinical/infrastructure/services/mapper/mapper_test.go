package mapper

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/assert"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

func TestDefaultTimelineMapper_ToTimeline(t *testing.T) {
	fakeID := gofakeit.UUID()
	fakeDate := scalarutils.Date{Year: 2023, Month: 5, Day: 10}
	fakeTime := time.Date(2023, 5, 10, 12, 0, 0, 0, time.UTC)
	fakeStatus := domain.ObservationStatusEnumFinal
	fakeText := "Test Name"
	fakeValue := "Test Value"

	obs := &domain.FHIRObservation{
		ID:                &fakeID,
		Code:              &domain.FHIRCodeableConcept{Text: fakeText},
		Status:            &fakeStatus,
		ValueString:       &fakeValue,
		EffectiveDateTime: func() *scalarutils.DateTime { d := scalarutils.DateTime(fakeTime.Format(time.RFC3339)); return &d }(),
	}

	cond := &domain.FHIRCondition{
		ID:             &fakeID,
		Code:           &domain.FHIRCodeableConcept{Text: fakeText},
		ClinicalStatus: &domain.FHIRCodeableConcept{Text: "Active"},
		Category:       []*domain.FHIRCodeableConcept{{Text: "Diagnosis"}},
		OnsetDateTime:  &fakeDate,
	}

	risk := &domain.FHIRRiskAssessment{
		ID:     &fakeID,
		Code:   &domain.FHIRCodeableConcept{Text: fakeText},
		Status: domain.ObservationStatusEnumFinal,
		Prediction: []domain.FHIRRiskAssessmentPrediction{{
			Outcome: &domain.FHIRCodeableConcept{Text: "Low"},
		}},
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(dto.CervicalCancerScreeningTypeEnum),
		},
		OccurrenceDateTime: func() *string { s := fakeTime.Format(time.RFC3339); return &s }(),
	}

	allergy := &domain.FHIRAllergyIntolerance{
		ID:           &fakeID,
		Code:         &domain.FHIRCodeableConcept{Text: fakeText},
		RecordedDate: func() *scalarutils.DateTime { d := scalarutils.DateTime(fakeTime.Format(time.RFC3339)); return &d }(),
	}

	invalid := struct{ Foo string }{Foo: "bar"}

	tests := []struct {
		name     string
		input    interface{}
		wantType dto.ResourceType
		wantErr  bool
	}{
		{
			name:     "Observation happy case",
			input:    obs,
			wantType: dto.ResourceTypeObservation,
			wantErr:  false,
		},
		{
			name:     "Condition happy case",
			input:    cond,
			wantType: dto.ResourceTypeCondition,
			wantErr:  false,
		},
		{
			name:     "RiskAssessment happy case",
			input:    risk,
			wantType: dto.ResourceTypeRiskAssessment,
			wantErr:  false,
		},
		{
			name:     "AllergyIntolerance happy case",
			input:    allergy,
			wantType: dto.ResourceTypeAllergyIntolerance,
			wantErr:  false,
		},
		{
			name:    "Unsupported type returns error",
			input:   invalid,
			wantErr: true,
		},
	}

	mapper := NewTimelineMapper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := mapper.ToTimeline(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.wantType, res.ResourceType)
				assert.NotEmpty(t, res.ID)
			}
		})
	}
}
