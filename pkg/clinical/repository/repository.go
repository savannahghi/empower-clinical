package repository

import (
	"context"

	hapifhirmodels "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// FHIR represents all the FHIR logic
type FHIR interface {
	FHIROrganization
	FHIRBundle
	FHIRPatient
	FHIREpisodeOfCare
	FHIRObservation
	FHIRAllergyIntolerance
	FHIRServiceRequest
	FHIRMedicationRequest
	FHIRCondition
	FHIREncounter
	FHIRComposition
	FHIRMedicationStatement
	FHIRMedication
	FHIRMedia
	FHIRQuestionnaire
	FHIRConsent
	FHIRQuestionnaireResponse
	FHIRRiskAssessment
	FHIRDiagnosticReport
	FHIRSubscription
	FHIRDocumentReference
	FHIRAppointment
	FHIRTask
	FHIRPractitioner
	FHIRPractitionerRole
	FHIRProcedure
	FHIRResource
	FHIRPlanDefinition
	FHIRActivityDefinition
	FHIRCarePlan
}

type FHIRBundle interface {
	PostFHIRBundle(ctx context.Context, payload *hapifhirmodels.Bundle) (*hapifhirmodels.Bundle, error)
}

type FHIROrganization interface {
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	SearchFHIROrganization(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIROrganizationRelayConnection, error)
	GetFHIROrganization(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error)
	PutFHIROrganization(ctx context.Context, id string, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
}

type FHIRPatient interface {
	GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error)
	DeleteFHIRPatient(ctx context.Context, id string) (bool, error)
	CreateFHIRPatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	PatchFHIRPatient(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error)
	SearchFHIRPatient(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error)
	GetFHIRPatientEverything(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error)
	ValidateResource(ctx context.Context, resourceType string, params map[string]interface{}) error
}
type FHIREpisodeOfCare interface {
	SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error)
	SearchEpisodesByParam(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error)
	GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error)
	PatchFHIREpisodeOfCare(ctx context.Context, id string, episode domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error)
	UpdateFHIREpisodeOfCare(ctx context.Context, fhirResourceID string, input domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error)
	HasOpenEpisode(ctx context.Context, patient domain.FHIRPatient, tenant dto.TenantIdentifiers, pagination dto.Pagination) (bool, error)
	OpenEpisodes(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error)
	EndEpisode(ctx context.Context, episodeID string) (bool, error)
	GetActiveEpisode(ctx context.Context, episodeID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCare, error)
}
type FHIRObservation interface {
	GetFHIRObservation(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error)
	SearchFHIRObservation(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error)
	CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error)
	DeleteFHIRObservation(ctx context.Context, id string) (bool, error)
	SearchPatientObservations(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error)
	PatchFHIRObservation(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error)
}
type FHIRAllergyIntolerance interface {
	SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error)
	CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	GetFHIRAllergyIntolerance(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	SearchPatientAllergyIntolerance(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error)
}
type FHIRServiceRequest interface {
	SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRServiceRequest, error)
	CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error)
	DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error)
	GetFHIRServiceRequest(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error)
	PatchFHIRServiceRequest(ctx context.Context, id string, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error)
}
type FHIRMedicationRequest interface {
	SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedicationRequest, error)
	CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error)
	PatchFHIRMedicationRequest(ctx context.Context, id string, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequest, error)
	GetFHIRMedicationRequest(ctx context.Context, id string) (*domain.FHIRMedicationRequestRelayPayload, error)
}
type FHIRCondition interface {
	SearchFHIRCondition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error)
	CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	GetFHIRCondition(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error)
}
type FHIREncounter interface {
	CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error)
	SearchPatientEncounters(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
	StartEncounter(ctx context.Context, episodeID string) (string, error)
	SearchEpisodeEncounter(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error)
	SearchFHIREncounter(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
	PatchFHIREncounter(ctx context.Context, encounterID string, input domain.FHIREncounterInput) (*domain.FHIREncounter, error)
	SearchFHIREncounterAllData(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error)
}
type FHIRComposition interface {
	GetFHIRComposition(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error)
	SearchFHIRComposition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRComposition, error)
	CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error)
	PatchFHIRComposition(ctx context.Context, id string, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error)
	DeleteFHIRComposition(ctx context.Context, id string) (bool, error)
}
type FHIRMedicationStatement interface {
	CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error)
}

type FHIRMedication interface {
	CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error)
	FetchMedicationByID(ctx context.Context, id string) (*domain.FHIRMedication, error)
}

type FHIRMedia interface {
	CreateFHIRMedia(ctx context.Context, input domain.FHIRMedia) (*domain.FHIRMedia, error)
	SearchPatientMedia(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedia, error)
}

type FHIRQuestionnaire interface {
	GetFHIRQuestionnaire(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error)
	CreateFHIRQuestionnaire(ctx context.Context, input *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error)
	ListFHIRQuestionnaire(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error)
}
type FHIRConsent interface {
	CreateFHIRConsent(ctx context.Context, input domain.FHIRConsent) (*domain.FHIRConsent, error)
}

type FHIRQuestionnaireResponse interface {
	CreateFHIRQuestionnaireResponse(_ context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error)
	GetFHIRQuestionnaireResponse(ctx context.Context, id string) (*domain.FHIRQuestionnaireResponseRelayPayload, error)
}

type FHIRRiskAssessment interface {
	CreateFHIRRiskAssessment(_ context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error)
	GetFHIRRiskAssessment(ctx context.Context, id string) (*domain.FHIRRiskAssessmentRelayPayload, error)
	SearchFHIRRiskAssessment(ctx context.Context, bundleID string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRRiskAssessment, error)
}

type FHIRDiagnosticReport interface {
	CreateFHIRDiagnosticReport(_ context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error)
	SearchFHIRDiagnosticReport(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDiagnosticReport, error)
}

// FHIRSubscription contains method signatures of the FHIRSubscription interface
type FHIRSubscription interface {
	CreateFHIRSubscription(_ context.Context, input *domain.FHIRSubscriptionInput) (*domain.FHIRSubscription, error)
}

// FHIRDocumentReference interface contains the method signatures for the various action to be performed against the FHIR Document Reference resource
type FHIRDocumentReference interface {
	CreateFHIRDocumentReference(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error)
	SearchFHIRDocumentReference(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error)
}

// FHIRAppointment holds appointment related methods and their signatures
type FHIRAppointment interface {
	CreateFHIRAppointment(ctx context.Context, input *domain.FHIRAppointmentInput) (*domain.FHIRAppointment, error)
}

// FHIRTask holds task related methods and their signatures
type FHIRTask interface {
	CreateFHIRTask(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error)
	PatchFHIRTask(ctx context.Context, input domain.FHIRTaskInput) (*domain.FHIRTask, error)
	SearchFHIRTask(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRTask, error)
	GetFHIRTask(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error)
}

type FHIRLocation interface {
	CreateFHIRLocation(ctx context.Context, input *domain.FHIRLocation) (*domain.FHIRLocation, error)
	GetFHIRLocation(ctx context.Context, id string) (*domain.FHIRLocation, error)
}

type FHIRPractitionerRole interface {
	CreateFHIRPractitionerRole(ctx context.Context, input *domain.FHIRPractitionerRole) (*domain.FHIRPractitionerRole, error)
	SearchFHIRPractitionerRoles(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRPractitionerRole, error)
	GetFHIRPractitionerRole(ctx context.Context, id string) (*domain.FHIRPractitionerRoleRelayPayload, error)
}

type FHIRSubstance interface {
	CreateFHIRSubstance(ctx context.Context, input *domain.FHIRSubstance) (*domain.FHIRSubstance, error)
	SearchFHIRSubstances(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRSubstance, error)
	GetFHIRSubstance(ctx context.Context, id string) (*domain.FHIRSubstanceRelayPayload, error)
}

type FHIRPractitioner interface {
	CreateFHIRPractioner(ctx context.Context, input *domain.FHIRPractitioner) (*domain.FHIRPractitioner, error)
	SearchFHIRPractitioner(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PractitionerConnection, error)
	GetFHIRPractitioner(ctx context.Context, id string) (*domain.FHIRPractitioner, error)
}

type FHIRProcedure interface {
	CreateFHIRProcedure(ctx context.Context, input *domain.FHIRProcedure) (*domain.FHIRProcedure, error)
	SearchFHIRProcedure(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRProcedure, error)
	GetFHIRProcedure(ctx context.Context, id string) (*domain.FHIRProcedure, error)
}

type FHIRMedicationDispense interface {
	CreateFHIRMedicationDispense(ctx context.Context, input *domain.FHIRMedicationDispense) (*domain.FHIRMedicationDispense, error)
	SearchFHIRMedicationDispense(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedicationDispense, error)
	GetFHIRMedicationDispense(ctx context.Context, id string) (*domain.FHIRMedicationDispense, error)
}
type FHIRResource interface {
	SearchFHIRResource(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error)
	DeleteFHIRResource(ctx context.Context, resourceType, resourceID string) (bool, error)
}

type FHIRPlanDefinition interface {
	CreateFHIRPlanDefinition(ctx context.Context, input *domain.FHIRPlanDefinition) (*domain.FHIRPlanDefinition, error)
	SearchFHIRPlanDefinition(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRPlanDefinition, error)
	FetchPlanDefinitionByID(ctx context.Context, id string) (*domain.FHIRPlanDefinition, error)
}

type FHIRActivityDefinition interface {
	CreateFHIRActivityDefinition(ctx context.Context, input *domain.FHIRActivityDefinition) (*domain.FHIRActivityDefinition, error)
}

type FHIRCarePlan interface {
	CreateFHIRCarePlan(ctx context.Context, input *domain.FHIRCarePlan) (*domain.FHIRCarePlan, error)
	SearchFHIRCarePlan(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error)
}
