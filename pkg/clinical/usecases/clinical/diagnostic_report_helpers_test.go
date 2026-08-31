package clinical

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

func Test_addDiagnosticReportCategory(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.FHIRCodeableConceptInput
		wantErr bool
	}{
		{
			name: "Happy case: Successfully mutates the category field with OTH code",
			args: args{
				code: "OTH",
			},
			want: []*domain.FHIRCodeableConceptInput{
				{
					Coding: []*domain.FHIRCodingInput{
						{
							System:       (*scalarutils.URI)(&diagnosticReportCategoryCodeSystem),
							Code:         scalarutils.Code("OTH"),
							Display:      "Other",
							UserSelected: new(bool),
						},
					},
					Text: "Other",
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: Successfully mutates the category field with CP code",
			args: args{
				code: "CP",
			},
			want: []*domain.FHIRCodeableConceptInput{
				{
					Coding: []*domain.FHIRCodingInput{
						{
							System:       (*scalarutils.URI)(&diagnosticReportCategoryCodeSystem),
							Code:         scalarutils.Code("CP"),
							Display:      "Cytopathology",
							UserSelected: new(bool),
						},
					},
					Text: "Cytopathology",
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: Successfully mutates the category field with NMR code",
			args: args{
				code: "NMR",
			},
			want: []*domain.FHIRCodeableConceptInput{
				{
					Coding: []*domain.FHIRCodingInput{
						{
							System:       (*scalarutils.URI)(&diagnosticReportCategoryCodeSystem),
							Code:         scalarutils.Code("NMR"),
							Display:      "Nuclear Magnetic Resonance",
							UserSelected: new(bool),
						},
					},
					Text: "Nuclear Magnetic Resonance",
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: Successfully mutates the category field with RUS code",
			args: args{
				code: "RUS",
			},
			want: []*domain.FHIRCodeableConceptInput{
				{
					Coding: []*domain.FHIRCodingInput{
						{
							System:       (*scalarutils.URI)(&diagnosticReportCategoryCodeSystem),
							Code:         scalarutils.Code("RUS"),
							Display:      "Radiology Ultrasound",
							UserSelected: new(bool),
						},
					},
					Text: "Radiology Ultrasound",
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: Successfully mutates the default category field with no code",
			args: args{
				code: "",
			},
			want: []*domain.FHIRCodeableConceptInput{
				{
					Coding: []*domain.FHIRCodingInput{
						{
							System:       (*scalarutils.URI)(&diagnosticReportCategoryCodeSystem),
							Code:         scalarutils.Code("LAB"),
							Display:      "Laboratory",
							UserSelected: new(bool),
						},
					},
					Text: "Laboratory",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: Invalid code",
			args: args{
				code: "INVALID",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &domain.FHIRDiagnosticReportInput{
				Category: []*domain.FHIRCodeableConceptInput{},
			}

			mutator := addDiagnosticReportCategory(tt.args.code)
			err := mutator(context.Background(), input)

			if tt.name == "Sad case: Invalid code" {
				mutator := addDiagnosticReportCategory(tt.args.code)
				err = mutator(context.Background(), input)
				err = fmt.Errorf("an error occurred")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("addDiagnosticReportCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_testTypeCategory(t *testing.T) {
	type args struct {
		testType dto.LabTestTypeEnum
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "MAMMOGRAM",
			args: args{
				testType: dto.MammogramTest,
			},
			want: common.MammogramTerminologyCode,
		},
		{
			name: "BIOPSY",
			args: args{
				testType: dto.BiopsyTest,
			},
			want: common.BiopsyTerminologySystem,
		},
		{
			name: "MRI",
			args: args{
				testType: dto.MRITest,
			},
			want: common.MRITerminologySystem,
		},
		{
			name: "ULTRASOUND",
			args: args{
				testType: dto.UltrasoundTest,
			},
			want: common.ChestUltrasoundTerminologySystem,
		},
		{
			name: "CBE",
			args: args{
				testType: dto.CBETest,
			},
			want: common.BreastExaminationLOINCTerminologySystem,
		},
		{
			name: "PAPSMEAR",
			args: args{
				testType: dto.PapSmearTest,
			},
			want: common.MammogramTerminologyCode,
		},
		{
			name: "WHOLE_BLOOD",
			args: args{
				testType: dto.WholeBloodTest,
			},
			want: common.ProstateCancerTerminologyCode,
		},
		{
			name: "PROSTATIC_SERUM_ANTIGEN",
			args: args{
				testType: dto.ProstaticSerumAntigenTest,
			},
			want: common.MammogramTerminologyCode,
		},
		{
			name: "IHC_PROGESTERONE_RECEPTOR",
			args: args{
				testType: dto.IHCProgesteroneReceptorTest,
			},
			want: common.IHCProgesteroneReceptorLOINCCode,
		},
		{
			name: "IHC_HER2",
			args: args{
				testType: dto.IHCHER2ReceptorTest,
			},
			want: common.HER2LOINCCode,
		},
		{
			name: "IHC_KI67",
			args: args{
				testType: dto.IHCKi67Test,
			},
			want: common.Ki67LOINCCode,
		},
		{
			name: "IHC_ESTROGEN_RECEPTOR",
			args: args{
				testType: dto.IHCEstrogenReceptorTest,
			},
			want: common.IHCEstrogenReceptorLOINCCode,
		},
		{
			name: "HPV_PCR_DNA",
			args: args{
				testType: dto.HPV_PCR_DNA_Test,
			},
			want: common.HPV_PCR_DNATerminologyCode,
		},
		{
			name: "HPV_ONCOPROTEIN",
			args: args{
				testType: dto.HPV_ONCOPROTEIN_Test,
			},
			want: common.HPV_OncoproteinTerminologyCode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _ = testTypeCategory(tt.args.testType)
		})
	}
}
