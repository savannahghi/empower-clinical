package clinical

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/foundation"
)

// ObservationInputMutatorFunc modifies a FHIR observation input resource
// Used by methods to add logic not provided by the general RecordObservation method
// Example: add an interpretation to an observation which varies by input
type ObservationInputMutatorFunc func(context.Context, *domain.FHIRObservationInput) error

// RecordTemperature is used to record a patient's temperature and saves it as a FHIR observation
func (c *ClinicalImpl) RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.TemperatureLOINCTerminologyCode,
	}

	temperatureObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return temperatureObservation, nil
}

// RecordMuac is used to record a patient's Muac
func (c *ClinicalImpl) RecordMuac(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.MuacLOINCTerminologyCode,
	}

	muacObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return muacObservation, nil
}

// RecordOxygenSaturation is used to record a patient's oxygen saturation
func (c *ClinicalImpl) RecordOxygenSaturation(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.OxygenSaturationLOINCTerminologyCode,
	}

	oxygenSaturationObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return oxygenSaturationObservation, nil
}

// GetPatientTemperatureEntries returns all the temperature entries for a patient, they are automatically sorted in chronological order
func (c *ClinicalImpl) GetPatientTemperatureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.TemperatureLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordHeight records a patient's height and saves it to fhir
func (c *ClinicalImpl) RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.HeightLOINCTerminologyCode,
	}

	heightObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return heightObservation, nil
}

// GetPatientHeightEntries gets the height records of a patient
func (c *ClinicalImpl) GetPatientHeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.HeightLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// PatchPatientHeight patches the height record of a patient
func (c *ClinicalImpl) PatchPatientHeight(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientWeight patches the weight record of a patient
func (c *ClinicalImpl) PatchPatientWeight(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// UpdateTestResults is used to update the result of a given test
func (c *ClinicalImpl) UpdateTestResults(ctx context.Context, observationID string, result string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, observationID, result)
}

// PatchPatientBMI patches the BMI record of a patient
func (c *ClinicalImpl) PatchPatientBMI(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientTemperature patches the temperature record of a patient
func (c *ClinicalImpl) PatchPatientTemperature(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientDiastolicBloodPressure patches the diastolic blood pressure record of a patient
func (c *ClinicalImpl) PatchPatientDiastolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientSystolicBloodPressure patches the Systolic blood pressure record of a patient
func (c *ClinicalImpl) PatchPatientSystolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientRespiratoryRate patches the respiration rate record of a patient
func (c *ClinicalImpl) PatchPatientRespiratoryRate(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientOxygenSaturation patches the oxygen saturation record of a patient
func (c *ClinicalImpl) PatchPatientOxygenSaturation(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientPulseRate patches the pulse rate record of a patient
func (c *ClinicalImpl) PatchPatientPulseRate(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientViralLoad patches the viral load record of a patient
func (c *ClinicalImpl) PatchPatientViralLoad(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientMuac patches the muac record of a patient
func (c *ClinicalImpl) PatchPatientMuac(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientLastMenstrualPeriod patches the last menstrual record of a patient
func (c *ClinicalImpl) PatchPatientLastMenstrualPeriod(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientBloodSugar patches the blood sugar record of a patient
func (c *ClinicalImpl) PatchPatientBloodSugar(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// RecordWeight records a patient's weight
func (c *ClinicalImpl) RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.WeightLOINCTerminologyCode,
	}

	weightObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return weightObservation, nil
}

// RecordViralLoad records the patient viral load
func (c *ClinicalImpl) RecordViralLoad(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.ViralLoadLOINCTerminologyCode,
	}

	viralLoadObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return viralLoadObservation, nil
}

// GetPatientWeightEntries gets the weight records of a patient
func (c *ClinicalImpl) GetPatientWeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.WeightLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// GetPatientMuacEntries gets the patient's muac
func (c *ClinicalImpl) GetPatientMuacEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.MuacLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// GetPatientOxygenSaturationEntries gets the patient's oxygen saturation
func (c *ClinicalImpl) GetPatientOxygenSaturationEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.OxygenSaturationLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// GetPatientViralLoad gets the patient's viral load
func (c *ClinicalImpl) GetPatientViralLoad(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.ViralLoadLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordRespiratoryRate records a patient's respiratory rate
func (c *ClinicalImpl) RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.RespiratoryRateLOINCTerminologyCode,
	}

	respiratoryRateObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return respiratoryRateObservation, nil
}

// GetPatientRespiratoryRateEntries gets a patient's respiratory rate entries
func (c *ClinicalImpl) GetPatientRespiratoryRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.RespiratoryRateLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordPulseRate records a patient's pulse rate
func (c *ClinicalImpl) RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.PulseLOINCTerminologyCode,
	}

	pulseRateObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return pulseRateObservation, nil
}

// GetPatientPulseRateEntries gets the pulse rate records of a patient
func (c *ClinicalImpl) GetPatientPulseRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.PulseLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordBloodPressure records a patient's blood pressure
func (c *ClinicalImpl) RecordBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.BloodPressureLOINCTerminologyCode,
	}

	bloodPressureObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return bloodPressureObservation, nil
}

// GetPatientBloodPressureEntries retrieves all blood pressure entries for a patient
func (c *ClinicalImpl) GetPatientBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.BloodPressureLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordBMI records a patient's BMI
func (c *ClinicalImpl) RecordBMI(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.BMILOINCTerminologyCode,
	}

	bmiObservation, err := c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
	if err != nil {
		return nil, err
	}

	return bmiObservation, nil
}

// GetPatientBMIEntries retrieves all BMI entries for a patient
func (c *ClinicalImpl) GetPatientBMIEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.BMILOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordBloodSugar records a patient's blood sugar level (Serum glucose)
func (c *ClinicalImpl) RecordBloodSugar(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.BloodSugarLOINCTerminologyCode,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
}

// GetPatientBloodSugarEntries retrieves all blood sugar entries for a patient
func (c *ClinicalImpl) GetPatientBloodSugarEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.BloodSugarLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordLastMenstrualPeriod records last menstrual period
func (c *ClinicalImpl) RecordLastMenstrualPeriod(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.LastMenstrualPeriodLOINCTerminologyCode,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
}

// GetPatientLastMenstrualPeriodEntries retrieves all blood sugar entries for a patient
func (c *ClinicalImpl) GetPatientLastMenstrualPeriodEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.LastMenstrualPeriodLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordDiastolicBloodPressure records diastolic blood pressure
func (c *ClinicalImpl) RecordDiastolicBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.DiastolicBloodPressureLOINCTerminologyCode,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
}

// RecordColposcopy records colposcopy findings
func (c *ClinicalImpl) RecordColposcopy(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.ColposcopyLOINCTerminologyCode,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("vital-signs")})
}

// RecordHPV is used to record HPV test results. We record it as observations as specified in https://build.fhir.org/ig/HL7/cqf-measures/Measure-EXM124-FHIR.html
// Check whether the gender of the patient is valid and that the patient is within the acceptable age range
func (c *ClinicalImpl) RecordHPV(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	var terminologyCode string

	switch input.ObservationSubType {
	case dto.HPV_ONCOPROTEIN:
		terminologyCode = common.HPV_OncoproteinTerminologyCode
	case dto.HPV_PCR_DNA:
		terminologyCode = common.HPV_PCR_DNATerminologyCode
	default:
		terminologyCode = common.HPVLOINCTerminologyCode
	}

	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: terminologyCode,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("exam")})
}

// RecordVIA records Visual Inspection with Acetic Acid results
func (c *ClinicalImpl) RecordVIA(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	if !dto.VIAOutcomeEnum(input.Value).IsValid() {
		return nil, fmt.Errorf("invalid value for VIA: %s", input.Value)
	}

	// observation mutator func to add a VIA interpretation to a FHIR Observation
	addInterpretation := func(ctx context.Context, observation *domain.FHIRObservationInput) error {
		var conceptCode, conceptDiplay, interpretationText string

		system := "http://terminology.hl7.org/CodeSystem/v3-ObservationInterpretation"

		switch dto.VIAOutcomeEnum(input.Value) {
		case dto.VIAOutcomeNegative:
			conceptCode = "NEG"
			conceptDiplay = "Negative"
			interpretationText = "Negative"

		case dto.VIAOutcomePositive:
			conceptCode = "POS"
			conceptDiplay = "Positive"
			interpretationText = "Patient is at risk of cancer. Please enroll/refer for treatment"

		case dto.VIAOutcomePositiveInvasiveCancer:
			conceptCode = "E"
			conceptDiplay = "Equivocal"
			interpretationText = "Suspicious of cancer"
		}

		userSelected := false
		concept := &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       (*scalarutils.URI)(&system),
					Code:         scalarutils.Code(conceptCode),
					Display:      conceptDiplay,
					UserSelected: &userSelected,
				},
			},
			Text: interpretationText,
		}

		observation.Interpretation = append(observation.Interpretation, concept)

		return nil
	}

	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.VIALOINCCode,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addInterpretation, addObservationCategory("exam")})
}

// GetPatientDiastolicBloodPressureEntries retrieves all diastolic blood pressure entries for a patient
func (c *ClinicalImpl) GetPatientDiastolicBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.DiastolicBloodPressureLOINCTerminologyCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordObservation is an extracted function that takes any observation input and saves it to FHIR.
// A concept ID is also passed so that we can get the concept code of the passed observation
func (c *ClinicalImpl) RecordObservation(
	ctx context.Context,
	obsPayload dto.ObservationPayload,
	mutators []ObservationInputMutatorFunc,
) (*dto.Observation, error) {
	ctx, span := tracer.Start(ctx, "RecordObservation")

	defer func() { span.End() }()

	err := helpers.Validate(obsPayload.ObservationInput)
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, obsPayload.ObservationInput.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record an observation in a completed encounter")
	}

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	performerReference := fmt.Sprintf("Organization/%s", facilityID)

	currentTime := time.Now().Format(time.RFC3339)
	instant := scalarutils.Instant(currentTime)

	// EffectiveInstant reflects when the observation is clinically true (e.g. when the
	// test was performed). Use the supplied date when present, otherwise fall back to
	// the current time. Issued remains the time the observation is recorded. This is to protect
	// callers who don't supply a date
	effectiveInstant := scalarutils.Instant(resolveEffectiveTime(obsPayload.ObservationInput.Date).Format(time.RFC3339))

	observationTextStatus := "additional"
	status := strings.ToLower(string(obsPayload.ObservationInput.Status))

	observation := domain.FHIRObservationInput{
		Status:           (*domain.ObservationStatusEnum)(&status),
		Category:         []*domain.FHIRCodeableConceptInput{},
		EffectiveInstant: &effectiveInstant,

		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: encounter.Resource.Subject.Reference,
			Display:   encounter.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Reference: &encounterReference,
			Display:   *encounter.Resource.ID,
		},
		Issued:         (*scalarutils.Instant)(&instant),
		Interpretation: []*domain.FHIRCodeableConceptInput{},
		Performer: []*domain.FHIRReferenceInput{
			{
				Reference: &performerReference,
				Display:   facilityID,
			},
		},
	}

	if strings.TrimSpace(obsPayload.ObservationInput.Value) != "" {
		testResultValue := strings.TrimSpace(obsPayload.ObservationInput.Value)
		observation.ValueString = &testResultValue
	}

	if obsPayload.ObservationInput.Note != "" {
		note := domain.FHIRAnnotationInput{
			Text: (*scalarutils.Markdown)(&obsPayload.ObservationInput.Note),
		}
		observation.Note = append(observation.Note, &note)
	}

	if obsPayload.VitalSignsConceptID != "" {
		vitalsConcept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, obsPayload.VitalSignsConceptID)
		if err != nil {
			return nil, err
		}

		coding := domain.FHIRCodingInput{
			System:  (*scalarutils.URI)(&common.LoincSystemURL),
			Code:    scalarutils.Code(vitalsConcept.ID),
			Display: vitalsConcept.GetConceptDisplay(),
		}

		var codingInput []*domain.FHIRCodingInput

		codeableConcept := &domain.FHIRCodeableConceptInput{
			Coding: append(codingInput, &coding),
			Text:   vitalsConcept.GetConceptDisplay(),
		}

		observation.Code = codeableConcept
	}

	if obsPayload.ObservationInput.UsageContext != "" {
		observation.Text = utils.NarrativeGenerator(string(obsPayload.ObservationInput.UsageContext), &observationTextStatus)
	}

	if obsPayload.ServiceRequestID != "" {
		serviceRequestReference := fmt.Sprintf("ServiceRequest/%s", obsPayload.ServiceRequestID)

		basedOn := &domain.FHIRReferenceInput{
			ID:        &obsPayload.ServiceRequestID,
			Reference: &serviceRequestReference,
			Display:   obsPayload.ServiceRequestID,
		}

		observation.BasedOn = append(observation.BasedOn, basedOn)
	}

	for _, mutator := range mutators {
		err = mutator(ctx, &observation)
		if err != nil {
			return nil, err
		}
	}

	if len(observation.Category) < 1 {
		return nil, fmt.Errorf("observation category (i.e laboratory, vital-signs etc.) must be specified")
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	observation.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	fhirObservation, err := c.FHIR.CreateFHIRObservation(ctx, observation)
	if err != nil {
		return nil, err
	}

	return foundation.MapFHIRObservationToObservationDTO(*fhirObservation), nil
}

// GetPatientObservations is a helper function used to fetch patient's observations based on the passed LOINC
// terminology code. The observations will be sorted in a chronological order
func (c *ClinicalImpl) GetPatientObservations(ctx context.Context, payload *dto.FetchObservationPayload) (*dto.ObservationConnection, error) {
	searchParams := map[string]any{
		"_sort":  "-_lastUpdated",
		"_total": "accurate",
		"status": payload.SetStatus(),
	}

	if payload.PatientID != "" {
		_, err := uuid.Parse(payload.PatientID)
		if err != nil {
			return nil, fmt.Errorf("invalid patient id: %s", payload.PatientID)
		}

		patientReference := fmt.Sprintf("Patient/%s", payload.PatientID)
		searchParams["patient"] = patientReference
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	if payload.Usage == "SCREENING_EXAMINATIONS" { //TODO: Make an enum
		codes := []string{
			common.BreastExaminationLOINCTerminologySystem,
			common.HPV_OncoproteinTerminologyCode,
			common.HPV_PCR_DNATerminologyCode,
			common.ProstateCancerTerminologyCode,
			common.VIALOINCCode,
		}
		searchParams["code"] = strings.Join(codes, ",")
	}

	if payload.ObservationCode != "" {
		searchParams["code"] = payload.ObservationCode
	}

	if payload.EncounterID != nil {
		encounterReference := fmt.Sprintf("Encounter/%s", *payload.EncounterID)
		searchParams["encounter"] = encounterReference
	}

	if payload.Date != nil {
		searchParams["date"] = payload.Date.AsTime().Format(time.DateOnly)
	}

	if payload.Category != nil {
		if *payload.Category == dto.ObservationCategoryVitalSigns {
			searchParams["category"] = payload.Category.Text()
		} else {
			searchParams["category"] = payload.Category.String()
		}
	}

	var limit, offset int

	switch {
	case payload.PaginationV2 != nil:
		limit, offset, err = payload.PaginationV2.ToLimitOffset()
		if err != nil {
			return nil, fmt.Errorf("invalid pagination input: %w", err)
		}

		searchParams["_count"] = strconv.Itoa(limit)
	case payload.Pagination != nil && payload.Pagination.First != nil:
		searchParams["_count"] = strconv.Itoa(*payload.Pagination.First)
	}

	var patientObs *domain.PagedFHIRObservations

	if payload.SearchID != "" {
		searchParams["_getpages"] = payload.SearchID
		searchParams["_getpagesoffset"] = strconv.Itoa(offset)
		searchParams["_pretty"] = "true"
		searchParams["_bundletype"] = "searchset"

		patientObs, err = c.FHIR.SearchPatientObservations(ctx, payload.SearchID, searchParams, *identifiers, *payload.PaginationV2)
		if err != nil {
			return nil, err
		}
	} else {
		var pagination *serverutils.PaginationInput

		// This if clause here is meant to support the pagination that has been there. Basically, existing pagination wont't be affected
		if payload.Pagination != nil {
			pagination = &serverutils.PaginationInput{
				First:  payload.Pagination.First,
				After:  &payload.Pagination.After,
				Last:   payload.Pagination.Last,
				Before: &payload.Pagination.Before,
			}
		} else {
			pagination = payload.PaginationV2
		}

		patientObs, err = c.FHIR.SearchPatientObservations(ctx, "", searchParams, *identifiers, *pagination)
		if err != nil {
			return nil, err
		}
	}

	observations := []dto.Observation{}

	for _, obs := range patientObs.Observations {
		if obs.Subject == nil {
			continue
		}

		if obs.Subject.Display == "" {
			continue
		}

		observations = append(observations, *foundation.MapFHIRObservationToObservationDTO(obs))
	}

	connection := serverutils.BuildLimitOffsetConnection(observations, offset, limit, patientObs.TotalCount)

	return &dto.ObservationConnection{
		SearchID:   patientObs.BundleID,
		TotalCount: connection.TotalCount,
		Edges:      dto.ConvertObservationEdges(connection.Edges),
		PageInfo:   connection.PageInfo,
	}, nil
}

// GetObservationByID is used to get the details of an observation
func (c *ClinicalImpl) GetObservationByID(ctx context.Context, id string) (*dto.Observation, error) {
	observation, err := c.FHIR.GetFHIRObservation(ctx, id)
	if err != nil {
		return nil, err
	}

	return foundation.MapFHIRObservationToObservationDTO(*observation.Resource), nil
}

// PatchPatientObservations update a patient's observation resource
func (c *ClinicalImpl) PatchPatientObservations(ctx context.Context, id string, value string) (*dto.Observation, error) {
	if value == "" {
		return nil, fmt.Errorf("observation value required")
	}

	if id == "" {
		return nil, fmt.Errorf("an observation id is required")
	}

	observation, err := c.FHIR.GetFHIRObservation(ctx, id)
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, *observation.Resource.Encounter.ID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot patch an observation in a completed encounter")
	}

	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	observationInput := &domain.FHIRObservationInput{
		EffectiveInstant: &instant,
		ValueString:      &value,
	}

	output, err := c.FHIR.PatchFHIRObservation(ctx, id, *observationInput)
	if err != nil {
		return nil, err
	}

	return foundation.MapFHIRObservationToObservationDTO(*output), nil
}

// RecordImmunoHistoChemistry creates an immunohistochemistry observation record.
func (c *ClinicalImpl) RecordImmunoHistoChemistry(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.ImmunoHistoChemistryLOINCCode,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("laboratory")})
}

// GetPatientImmunoHistoChemistryRecords returns all patients IHC records
func (c *ClinicalImpl) GetPatientImmunoHistoChemistryRecords(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.ImmunoHistoChemistryLOINCCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordPostCoitalBleeding creates a post-coital bleeding observation record.
func (c *ClinicalImpl) RecordPostCoitalBleeding(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.PostCoitalBleedingCIELCode,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("laboratory")})
}

// GetPatientPostCoitalBleedingRecords returns all patient's PCB records
func (c *ClinicalImpl) GetPatientPostCoitalBleedingRecords(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.PostCoitalBleedingCIELCode,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordHistoryOfPresentIllness creates a history of a present illness
func (c *ClinicalImpl) RecordHistoryOfPresentIllness(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.LOINCHistoryOfPresentingIllness,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("exam")})
}

// RecordPastMedicalAndSurgicalHistory creates a record of past medical and surgical history
func (c *ClinicalImpl) RecordPastMedicalAndSurgicalHistory(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.LOINCHistoryOfPresentingIllness,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("exam")})
}

// RecordFamilyAndSocialHistory creates family and social history record.
func (c *ClinicalImpl) RecordFamilyAndSocialHistory(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: common.LOINCFamilyHistory,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory("exam")})
}

// GetHistoryOfPresentIllness returns all history of present illness
func (c *ClinicalImpl) GetHistoryOfPresentIllness(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.LOINCHistoryOfPresentingIllness,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// GetPastMedicalAndSurgicalHistory returns all history of present illness
func (c *ClinicalImpl) GetPastMedicalAndSurgicalHistory(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.LOINCPastMedicalSurgeryHistory,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// GetFamilyAndSocialHistory returns all history of present illness
func (c *ClinicalImpl) GetFamilyAndSocialHistory(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: common.LOINCFamilyHistory,
		Category:        nil,
		Pagination:      pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// ListObservations is used list observations
func (c *ClinicalImpl) ListObservations(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, category *dto.ObservationCategory, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterID,
		Date:            date,
		ObservationCode: "",
		Category:        nil,
		Pagination:      &pagination,
	}

	return c.GetPatientObservations(ctx, payload)
}

// RecordObservationV2 is used by the REST endpoint which exposes one endpoint
// for all POST operations on Observations
// ObservationInput is passed with a concept string specifies the Observation category being recorded
// It will eventually replace RecordObservation when migration to REST API is finalized
func (c *ClinicalImpl) RecordObservationV2(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	conceptID, category, err := c.ConceptMapper(input.Concept)
	if err != nil {
		return nil, err
	}

	observation := &dto.ObservationPayload{
		ObservationInput:    input,
		VitalSignsConceptID: conceptID,
	}

	return c.RecordObservation(ctx, *observation, []ObservationInputMutatorFunc{addObservationCategory(category)})
}

func (c *ClinicalImpl) PatchPatientObservation(ctx context.Context, observationID string, input *dto.PatchObservationInput) (*dto.Observation, error) {
	if err := helpers.Validate(input); err != nil {
		return nil, err
	}

	if input.ObservationType != "" && !input.ObservationType.IsValid() {
		return nil, fmt.Errorf("an invalid observation type")
	}

	return c.PatchPatientObservations(ctx, observationID, input.Value)
}
