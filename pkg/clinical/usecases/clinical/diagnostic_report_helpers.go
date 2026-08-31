package clinical

import (
	"context"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// DiagnosticReportMutatorFunc is a helper function that is used by the "caller" of "RecordDiagnosticReport" to modify the diagnostic report category field
// with the aapropriate data to suit its use case.
type DiagnosticReportMutatorFunc func(context.Context, *domain.FHIRDiagnosticReportInput) error

// addDiagnosticReportCategory mutates the category field, based on code, of a diagnostic report.
func addDiagnosticReportCategory(code string) DiagnosticReportMutatorFunc {
	var display, categoryCode string

	switch code {
	case "CP":
		display = "Cytopathology"
		categoryCode = code
	case "NMR":
		display = "Nuclear Magnetic Resonance"
		categoryCode = code
	case "RUS":
		display = "Radiology Ultrasound"
		categoryCode = code
	case "OTH":
		display = "Other"
		categoryCode = code
	default:
		display = "Laboratory"
		categoryCode = "LAB"
	}

	return func(ctx context.Context, input *domain.FHIRDiagnosticReportInput) error {
		userSelected := false
		category := []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       (*scalarutils.URI)(&diagnosticReportCategoryCodeSystem),
						Code:         scalarutils.Code(categoryCode),
						Display:      display,
						UserSelected: &userSelected,
					},
				},
				Text: display,
			},
		}

		input.Category = category

		return nil
	}
}

func testTypeCategory(testType dto.LabTestTypeEnum) (string, string, string) {
	var (
		observationCategory, diagnosticReportCategory, conceptID string
	)

	switch testType {
	case "MAMMOGRAM":
		conceptID = common.MammogramTerminologyCode
		observationCategory = "imaging"
		diagnosticReportCategory = ""
	case "MRI":
		conceptID = common.MRITerminologySystem
		observationCategory = "imaging"
		diagnosticReportCategory = "NMR"
	case "CBE":
		conceptID = common.BreastExaminationLOINCTerminologySystem
		observationCategory = "exam"
		diagnosticReportCategory = "OTH"
	case "PAPSMEAR":
		conceptID = common.PapSmearTerminologyCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"
	case "BIOPSY":
		conceptID = common.BiopsyTerminologySystem
		observationCategory = "procedure"
		diagnosticReportCategory = "CP"

	case "ULTRASOUND":
		conceptID = common.ChestUltrasoundTerminologySystem
		observationCategory = "imaging"
		diagnosticReportCategory = "RUS"

	case "WHOLE_BLOOD":
		conceptID = common.WholeBloodTerminologyCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"

	case "PROSTATIC_SERUM_ANTIGEN":
		conceptID = common.ProstateCancerTerminologyCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"

	case "IHC_PROGESTERONE_RECEPTOR":
		conceptID = common.IHCProgesteroneReceptorLOINCCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"

	case "IHC_ESTROGEN_RECEPTOR":
		conceptID = common.IHCEstrogenReceptorLOINCCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"

	case "IHC_HER2":
		conceptID = common.HER2LOINCCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"

	case "IHC_KI67":
		conceptID = common.Ki67LOINCCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"

	case "HPV_PCR_DNA":
		conceptID = common.HPV_PCR_DNATerminologyCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"

	case "HPV_ONCOPROTEIN":
		conceptID = common.HPV_OncoproteinTerminologyCode
		observationCategory = "laboratory"
		diagnosticReportCategory = "OTH"
	}

	return observationCategory, diagnosticReportCategory, conceptID
}
