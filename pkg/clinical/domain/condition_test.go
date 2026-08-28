package domain

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
)

func TestFHIRConditionInput_ComposeConditionCodeField(t *testing.T) {

	UUID := uuid.New().String()
	statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
	status := "active"
	note := scalarutils.Markdown("Fever Fever")

	type args struct {
		code           string
		displayName    string
		system         string
		allowedSources map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad Case: Missing code",
			args: args{
				code:        "",
				displayName: "Malaria",
				system:      "ICHI",
				allowedSources: map[string]string{
					"ICHI": "ICHI",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Missing display name",
			args: args{
				code:        "1F4Z",
				displayName: "",
				system:      "ICD-11",
				allowedSources: map[string]string{
					"ICD-11": "ICD-11",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Missing system",
			args: args{
				code:        "1F4Z",
				displayName: "Malaria",
				system:      "",
				allowedSources: map[string]string{
					"ICD-11": "ICD-11",
				},
			},
			wantErr: true,
		},
		{
			name: "Happy case: successful mutation of the condition code field",
			args: args{
				code:        "1F4Z",
				displayName: "Malaria",
				system:      "ICD-11",
				allowedSources: map[string]string{
					"ICD-11": "ICD-11",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FHIRConditionInput{
				ID:         &UUID,
				Text:       &FHIRNarrativeInput{},
				Identifier: []*FHIRIdentifierInput{},
				ClinicalStatus: &FHIRCodeableConceptInput{
					Coding: []*FHIRCodingInput{
						{
							System:  &statusSystem,
							Display: string(status),
						},
					},
					Text: string(status),
				},
				OnsetDateTime: &scalarutils.Date{},
				RecordedDate:  &scalarutils.Date{},
				Note: []*FHIRAnnotationInput{
					{
						Text: &note,
					},
				},
				Subject: &FHIRReferenceInput{
					ID: &UUID,
				},
				Encounter: &FHIRReferenceInput{
					ID: &UUID,
				},
				Category: []*FHIRCodeableConceptInput{
					{
						ID: &UUID,
						Coding: []*FHIRCodingInput{
							{
								ID:           UUID,
								System:       (*scalarutils.URI)(&UUID),
								Version:      &UUID,
								Display:      gofakeit.BeerAlcohol(),
								UserSelected: new(bool),
							},
						},
						Text: "PROBLEM_LIST_ITEM",
					},
				},
			}
			if err := condition.ComposeConditionCodeField(tt.args.code, tt.args.displayName, tt.args.system, tt.args.allowedSources); (err != nil) != tt.wantErr {
				t.Errorf("FHIRConditionInput.ComposeConditionCodeField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
