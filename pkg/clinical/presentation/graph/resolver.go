package graph

import (
	"context"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

//go:generate go run github.com/99designs/gqlgen

type usecases interface {
	PatchEpisodeOfCare(ctx context.Context, id string, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error)
	CreateEpisodeOfCare(ctx context.Context, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error)
	EndEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error)
	StartEncounter(ctx context.Context, episodeID string) (string, error)
	PatchEncounter(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordBMI(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordViralLoad(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordMuac(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordOxygenSaturation(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordLastMenstrualPeriod(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordDiastolicBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordColposcopy(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordVIA(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordImmunoHistoChemistry(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordPostCoitalBleeding(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordHistoryOfPresentIllness(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordPastMedicalAndSurgicalHistory(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordFamilyAndSocialHistory(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordHPV(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	CreatePatient(ctx context.Context, input dto.PatientInput) (*dto.Patient, error)

	RecordBloodSugar(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	PatchPatient(ctx context.Context, id string, input dto.PatientInput) (*dto.Patient, error)
	DeletePatient(ctx context.Context, id string) (bool, error)
	CreateCondition(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error)

	CreateAllergyIntolerance(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error)
	PatchPatientHeight(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientWeight(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientTemperature(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientDiastolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientSystolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientRespiratoryRate(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientOxygenSaturation(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientPulseRate(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientViralLoad(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientMuac(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientLastMenstrualPeriod(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientBloodSugar(ctx context.Context, id string, value string) (*dto.Observation, error)
	RecordConsent(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error)
	CreateQuestionnaireResponse(ctx context.Context, questionnaireID string, encounterID string, input dto.QuestionnaireResponse) (*dto.QuestionnaireReviewSummary, error)
	RecordMammographyResult(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error)
	RecordBiopsy(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error)
	RecordUltrasound(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error)
	RecordPapSmear(ctx context.Context, input *dto.DiagnosticReportInput) (*dto.DiagnosticReport, error)

	UpdateTestResults(ctx context.Context, id string, value string) (*dto.Observation, error)
	GetEncounterAssociatedResources(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResourceOutput, error)
	ReferPatient(ctx context.Context, input *dto.ReferralInput) (*dto.ServiceRequest, error)
	CreateReferral(ctx context.Context, input *dto.CreateReferralInput) (*dto.ServiceRequest, error)
	CreateTestOrder(ctx context.Context, input *dto.TestOrder) (*dto.ServiceRequest, error)
	CreateLabOrderResult(ctx context.Context, input *dto.TestOrderResult) (*dto.DiagnosticReport, error)
	ShareReferralForm(ctx context.Context, input *dto.ShareReferralFormInput) (bool, error)
	EndScreening(ctx context.Context, encounterID string) (bool, error)
	ScheduleAppointment(ctx context.Context, input *dto.ScheduleAppointmentInput, headersInput *dto.AdvantageHeaders) (bool, error)
	AddTestResultsLater(ctx context.Context, taskInput *dto.TaskInput) (*dto.TaskOutput, error)
	UpdateTask(ctx context.Context, id string, input *dto.PatchTaskInput) (bool, error)
	CreatePrescription(ctx context.Context, input dto.PrescriptionInput) ([]*dto.MedicationRequestOutput, error)
	PatchMedicationRequests(ctx context.Context, id string, value domain.MedicationRequestStatus) (*dto.MedicationRequestOutput, error)
	RecordTestResult(ctx context.Context, input dto.TestResultInput) (*dto.DiagnosticReport, error)
	GetPatientTimeline(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error)
	GetEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error)
	ListPatientConditions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, strategy string, pagination dto.Pagination) (*dto.ConditionConnection, error)
	ListPatientCompositions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.CompositionConnection, error)
	ListPatientEncounters(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.EncounterConnection, error)
	GetPatientTemperatureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientHeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientRespiratoryRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)

	GetPatientPulseRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientBMIEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientWeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientMuacEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientOxygenSaturationEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientViralLoad(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientBloodSugarEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientLastMenstrualPeriodEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientDiastolicBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientImmunoHistoChemistryRecords(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientPostCoitalBleedingRecords(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	ListObservations(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, category *dto.ObservationCategory, pagination dto.Pagination) (*dto.ObservationConnection, error)
	SearchAllergy(ctx context.Context, name string, pagination dto.Pagination) (*dto.TerminologyConnection, error)
	ListPatientAllergies(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.AllergyConnection, error)
	ListPatientMedia(ctx context.Context, encounterID, serviceRequestID string, pagination dto.Pagination) (*dto.MediaConnection, error)
	GetScreeningReport(ctx context.Context, encounterID string, status domain.ServiceRequestStatusEnum) (*dto.ScreeningReport, error)
	ListMedicationRequests(ctx context.Context, filterInput *dto.MedicationRequestFilterInput, pagination dto.Pagination) (*dto.MedicationRequestConnection, error)
	GetTaskByID(ctx context.Context, taskID string) (*dto.TaskOutput, error)
	FetchMedicationRequestByID(ctx context.Context, medicationRequestID string) (*dto.MedicationRequestOutput, error)
	GetFamilyAndSocialHistory(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetHistoryOfPresentIllness(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPastMedicalAndSurgicalHistory(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	ListLabOrders(ctx context.Context, filterInput *dto.ServiceRequestFilterInput, pagination dto.Pagination) (*dto.ServiceRequestOutputConnection, error)
	SimpleQuestionnaireResponse(ctx context.Context, questionnaireResponseID string) ([]*domain.SimpleQuestionnaireResponse, error)
	PatchPatientBMI(ctx context.Context, id string, value string) (*dto.Observation, error)
	RecordCBE(ctx context.Context, input *dto.DiagnosticReportInput) (*dto.DiagnosticReport, error)
	RecordMRI(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error)
	RecordPSA(ctx context.Context, input *dto.PSAInput) (*dto.DiagnosticReport, error)
	ListRiskAssessment(ctx context.Context, searchID string, input *dto.RiskAssessmentFilterInput, pagination serverutils.PaginationInput) (*dto.RiskAssessmentConnection, error)
	ListTasks(ctx context.Context, input *dto.TaskFilterInput, pagination dto.Pagination) (*dto.TaskOutputConnection, error)
	GetPatientReferrals(ctx context.Context, input *dto.ReferralSearchInput) (*dto.ReferralDetailConnection, error)
	GetPatientReferralDetails(ctx context.Context, serviceRequestID string) (*domain.PatientReferralDetails, error)
	GetLabOrder(ctx context.Context, serviceRequestID string) (*dto.ServiceRequest, error)
	GetAllergyIntolerance(ctx context.Context, id string) (*dto.Allergy, error)
	PatchPatientObservations(ctx context.Context, id string, value string) (*dto.Observation, error)
}

// Resolver wires up the resolvers needed for the clinical services
type Resolver struct {
	usecases usecases
}

// NewResolver initializes a working top leve Resolver that has been initialized
// with all necessary dependencies
func NewResolver(usecases usecases) (*Resolver, error) {
	return &Resolver{
		usecases: usecases,
	}, nil
}

// CheckDependencies ensures that the resolver has what it needs in order to work
func (r *Resolver) CheckDependencies() {
}
