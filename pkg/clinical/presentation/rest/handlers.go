package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"

	silgotel "github.com/savannahghi/sil-gotel"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/advantage"
)

const name = "github.com/savannahghi/empower-clinical-service/pkg/clinical/presentation/rest/handlers"

type usecases interface {
	CreatePubsubPatient(ctx context.Context, payload dto.PatientPubSubMessage) error
	CreatePubsubOrganization(ctx context.Context, data dto.FacilityPubSubMessage) error
	CreatePubsubVitals(ctx context.Context, data dto.VitalSignPubSubMessage) error
	CreatePubsubMedicationStatement(ctx context.Context, data dto.MedicationPubSubMessage) error
	CreatePubsubTenant(ctx context.Context, data dto.OrganizationInput) error
	CreatePubsubTestResult(ctx context.Context, data dto.PatientTestResultPubSubMessage) error
	CreateReferralTask(ctx context.Context, tags *dto.MetaInput, serviceRequest *dto.ServiceRequest) (*domain.FHIRTask, error)
	CreateTask(ctx context.Context, task *domain.FHIRTaskInput) (*domain.FHIRTask, error)
	RegisterTenant(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error)
	RegisterFacility(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error)
	ProvisionTenant(ctx context.Context, input dto.ProvisionTenantInput) (*dto.ProvisionTenantOutput, error)
	GetTenantProvisioningStatus(ctx context.Context, tenantID string) (*dto.ProvisionTenantOutput, error)
	UploadMedia(ctx context.Context, encounterID, serviceRequestID string, file io.Reader, contentType string) (*dto.Media, error)
	ListRiskAssessment(ctx context.Context, searchID string, filter *dto.RiskAssessmentFilterInput, pagination serverutils.PaginationInput) (*dto.RiskAssessmentConnection, error)
	ListTasks(ctx context.Context, filter *dto.TaskFilterInput, pagination dto.Pagination) (*dto.TaskOutputConnection, error)
	GenerateReferralReportPDF(ctx context.Context, serviceRequestID string) ([]byte, error)
	ReferPatient(ctx context.Context, input *dto.ReferralInput) (*dto.ServiceRequest, error)
	CreateReferral(ctx context.Context, input *dto.CreateReferralInput) (*dto.ServiceRequest, error)
	CreateTestOrder(ctx context.Context, input *dto.TestOrder) (*dto.ServiceRequest, error)
	ShareReferralForm(ctx context.Context, input *dto.ShareReferralFormInput) (bool, error)
	GetPatientReferrals(ctx context.Context, searchInput *dto.ReferralSearchInput) (*dto.ReferralDetailConnection, error)
	ListLabOrders(ctx context.Context, filter *dto.ServiceRequestFilterInput, pagination dto.Pagination) (*dto.ServiceRequestOutputConnection, error)
	CreateEpisodeOfCare(ctx context.Context, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error)
	CreateCondition(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error)
	RecordTreatmentEnrollment(ctx context.Context, input *dto.TreatmentEnrollmentInput) (*dto.Condition, error)
	UpdateTreatmentEnrollment(ctx context.Context, id string, input *dto.UpdateTreatmentEnrollmentInput) (*dto.Condition, error)
	CreateAllergyIntolerance(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error)
	PatchEpisodeOfCare(ctx context.Context, id string, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error)
	EndEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error)
	CreatePatient(ctx context.Context, input dto.PatientInput) (*dto.Patient, error)
	PatchPatient(ctx context.Context, id string, input dto.PatientInput) (*dto.Patient, error)
	StartEncounter(ctx context.Context, episodeID string) (string, error)
	CreateComposition(ctx context.Context, input dto.CompositionInput) (*dto.Composition, error)
	AppendNoteToComposition(ctx context.Context, id string, input dto.PatchCompositionInput) (*dto.Composition, error)
	UpdateTask(ctx context.Context, taskID string, updateData *dto.PatchTaskInput) (bool, error)
	PatchEncounter(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error)
	GetScreeningReport(ctx context.Context, encounterID string, status domain.ServiceRequestStatusEnum) (*dto.ScreeningReport, error)
	DeletePatient(ctx context.Context, id string) (bool, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	EndScreening(ctx context.Context, encounterID string) (bool, error)
	GetEncounterAssociatedResources(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResourceOutput, error)
	RecordConsent(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error)
	GetPatientReferralDetails(ctx context.Context, serviceRequestID string) (*domain.PatientReferralDetails, error)
	RecordObservationV2(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordTestResult(ctx context.Context, input dto.TestResultInput) (*dto.DiagnosticReport, error)
	CreateLabOrderResult(ctx context.Context, input *dto.TestOrderResult) (*dto.DiagnosticReport, error)
	ScheduleAppointment(ctx context.Context, input *dto.ScheduleAppointmentInput, headers *dto.AdvantageHeaders) (bool, error)
	RecordTests(ctx context.Context, payload dto.TestInput) (*dto.DiagnosticReport, error)
	CreatePrescription(ctx context.Context, input dto.PrescriptionInput) ([]*dto.MedicationRequestOutput, error)
	PatchMedicationRequests(ctx context.Context, id string, value domain.MedicationRequestStatus) (*dto.MedicationRequestOutput, error)
	GetLabOrder(ctx context.Context, serviceRequestID string) (*dto.ServiceRequest, error)
	GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error)
	GetPatientTimeline(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.HealthTimeline, error)
	GetPatientBanner(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.PatientBanner, error)
	GetEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error)
	SimpleQuestionnaireResponse(ctx context.Context, questionnaireResponseID string) ([]*domain.SimpleQuestionnaireResponse, error)
	ListPatientConditions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, strategy string, pagination dto.Pagination) (*dto.ConditionConnection, error)
	FetchMedicationRequestByID(ctx context.Context, medicationRequestID string) (*dto.MedicationRequestOutput, error)
	ListPatientCompositions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.CompositionConnection, error)
	SearchAllergy(ctx context.Context, name string, pagination dto.Pagination) (*dto.TerminologyConnection, error)
	GetAllergyIntolerance(ctx context.Context, id string) (*dto.Allergy, error)
	ListMedicationRequests(ctx context.Context, filter *dto.MedicationRequestFilterInput, pagination dto.Pagination) (*dto.MedicationRequestConnection, error)
	GetTaskByID(ctx context.Context, taskID string) (*dto.TaskOutput, error)
	GetPatientObservations(ctx context.Context, payload *dto.FetchObservationPayload) (*dto.ObservationConnection, error)
	ListPatientEncounters(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.EncounterConnection, error)
	CreatePubsubAllergyIntolerance(ctx context.Context, data dto.PatientAllergyPubSubMessage) error
	CreateQuestionnaire(ctx context.Context, questionnaireInput *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error)
	ListQuestionnaires(ctx context.Context, searchParam string, pagination *dto.Pagination) (*dto.Questionnaire, error)
	GenerateRiskAssessment(ctx context.Context, questionnaireID string, questionnaireResponse *dto.QuestionnaireResponse) (*dto.RiskAssessmentResult, error)
	FetchQuestionnaire(ctx context.Context, searchParam string, pagination *dto.Pagination) (*dto.Questionnaire, error)
	ConceptMapper(concept dto.ObservationConceptEnum) (string, string, error)
	ListPatientAllergies(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.AllergyConnection, error)
	ListPatientMedia(ctx context.Context, encounterID, serviceRequestID string, pagination dto.Pagination) (*dto.MediaConnection, error)
	PatchPatientObservation(ctx context.Context, observationID string, input *dto.PatchObservationInput) (*dto.Observation, error)
	CreateQuestionnaireResponse(ctx context.Context, questionnaireID string, encounterID string, input dto.QuestionnaireResponse) (*dto.QuestionnaireReviewSummary, error)
	CreatePlanDefinition(ctx context.Context, questionnaireInput *dto.PlanDefinitionInput) (*domain.FHIRPlanDefinition, error)
	RetrievePlanDefinition(ctx context.Context, name string) (*dto.PlanDefinitionOutputConnection, error)
	RecordOncologicalDiagnosis(ctx context.Context, input *dto.OncologyDiagnosisInput) (*dto.Condition, error)
	DeleteTest(ctx context.Context, observationID string) (bool, error)
	RecordMedication(ctx context.Context, medications []*dto.MedicationInput) ([]*dto.MedicationOutput, error)
	FetchMedicationByID(ctx context.Context, id string) (*dto.MedicationOutput, error)
	CreatePatientCarePlan(ctx context.Context, input *dto.CarePlanInput) error
	PatientCarePlan(ctx context.Context, input *dto.CarePlanPayload) (*domain.FHIRCarePlan, error)
	FetchPatientCarePlan(ctx context.Context, encounterID string) (*dto.CarePlanOutput, error)
	GetObservationByID(ctx context.Context, id string) (*dto.Observation, error)
	GetRiskAssessmentByID(ctx context.Context, id string) (*dto.RiskAssessment, error)
}

type baseExtension interface {
	VerifyPubSubJWTAndDecodePayload(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error)
	GetPubSubTopic(m *pubsubtools.PubSubPayload) (string, error)
}

// PresentationHandlersImpl represents the usecase implementation object
type PresentationHandlersImpl struct {
	usecases
	baseExtension
	advantage.AdvantageService
}

// NewPresentationHandlers initializes a new rest handlers usecase
func NewPresentationHandlers(usecases usecases, baseExt baseExtension, advantageSvc advantage.AdvantageService) *PresentationHandlersImpl {
	return &PresentationHandlersImpl{
		usecases,
		baseExt,
		advantageSvc,
	}
}

// ReceivePubSubPushMessage receives and processes a pubsub message
func (p PresentationHandlersImpl) ReceivePubSubPushMessage(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ReceivePubSubPushMessage")
	defer span.End()

	message, err := p.VerifyPubSubJWTAndDecodePayload(c.Writer, c.Request)
	if err != nil {
		serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	topicID, err := p.GetPubSubTopic(message)
	if err != nil {
		serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	switch topicID {
	case utils.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName):
		var data dto.PatientPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.CreatePubsubPatient(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName):
		var data dto.FacilityPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.CreatePubsubOrganization(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName):
		var data dto.VitalSignPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.CreatePubsubVitals(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName):
		var data dto.PatientAllergyPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.CreatePubsubAllergyIntolerance(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName):
		var data dto.MedicationPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.CreatePubsubMedicationStatement(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.TenantTopicName, common.ClinicalServiceName):
		var data dto.OrganizationInput

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.CreatePubsubTenant(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName):
		var data dto.PatientTestResultPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.CreatePubsubTestResult(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.SegmentationTopicName, common.ClinicalServiceName):
		var data dto.SegmentationPayload

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.SegmentPatient(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.ReferralTopicName, common.ClinicalServiceName):
		var data dto.PatientReferralTaskPayload

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		_, err = p.CreateReferralTask(ctx, data.Meta, data.Referral)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.ReferralReportNotificationTopic, common.ClinicalServiceName):
		// Do the whole process of notification
		var data dto.ReferralReportNotification

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.CreateCarePlanTopic, common.ClinicalServiceName):
		var data dto.CarePlanPayload

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		_, err = p.PatientCarePlan(ctx, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.FollowUpTaskTopic, common.ClinicalServiceName):
		var data domain.FHIRTaskInput

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		_, err = p.CreateTask(ctx, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	default:
		err := fmt.Errorf("unknown topic ID: %v", topicID)

		serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	resp := map[string]string{"Status": "Success"}
	c.JSON(http.StatusOK, resp)
}

func (p PresentationHandlersImpl) RegisterTenant(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "RegisterTenant")
	defer span.End()

	input := dto.OrganizationInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	organization, err := p.usecases.RegisterTenant(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, organization)
}

func jsonErrorResponse(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: err.Error(),
	})
}

// RegisterFacility creates a facility in fhir.
func (p PresentationHandlersImpl) RegisterFacility(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "RegisterFacility")
	defer span.End()

	input := dto.OrganizationInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	organization, err := p.usecases.RegisterFacility(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, organization)
}

// ProvisionTenant idempotently provisions (creates or retrieves) a tenant.
//
//	@Summary		Provision a tenant
//	@Description	Idempotently provisions a tenant. If the tenant already exists it is returned unchanged.
//	@Tags			Tenants
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.ProvisionTenantInput	true	"Tenant provisioning payload"
//	@Success		200		{object}	dto.ProvisionTenantOutput
//	@Failure		400		{object}	APIResponse
//	@Router			/v1/provision/ [post]
func (p PresentationHandlersImpl) ProvisionTenant(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ProvisionTenant")
	defer span.End()

	input := dto.ProvisionTenantInput{}

	if err := c.BindJSON(&input); err != nil {
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	output, err := p.usecases.ProvisionTenant(ctx, input)
	if err != nil {
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, output)
}

// GetTenantProvisioningStatus retrieves the provisioning status for a tenant.
//
//	@Summary		Get tenant provisioning status
//	@Description	Returns the provisioning status for the specified tenant.
//	@Tags			Tenants
//	@Produce		json
//	@Param			tenant-id	path		string	true	"Tenant ID"
//	@Success		200			{object}	dto.ProvisionTenantOutput
//	@Failure		400			{object}	APIResponse
//	@Failure		404			{object}	APIResponse
//	@Router			/v1/provision/{tenant-id}/ [get]
func (p PresentationHandlersImpl) GetTenantProvisioningStatus(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "GetTenantProvisioningStatus")
	defer span.End()

	tenantID := c.Param("tenant-id")
	if tenantID == "" {
		jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("tenant-id path parameter is required"))

		return
	}

	output, err := p.usecases.GetTenantProvisioningStatus(ctx, tenantID)
	if err != nil {
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusNotFound, err)

		return
	}

	c.JSON(http.StatusOK, output)
}

func (p PresentationHandlersImpl) CreatePlanDefinition(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateCareplanDefinition")
	defer span.End()

	input := &dto.PlanDefinitionInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	output, err := p.usecases.CreatePlanDefinition(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, output)
}

func (p PresentationHandlersImpl) CreatePatientCarePlan(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreatePatientCareplan")
	defer span.End()

	input := &dto.CarePlanInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	err = p.usecases.CreatePatientCarePlan(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (p PresentationHandlersImpl) FetchCarePlan(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "FetchCarePlan")
	defer span.End()

	searchParam := c.Request.URL.Query().Get("encounterID")

	output, err := p.FetchPatientCarePlan(ctx, searchParam)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, output)
}

// ListQuestionnaire is used to provide params used to fetch questionnaires
func (p PresentationHandlersImpl) ListPlanDefinition(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListPlanDefinition")
	defer span.End()

	searchParam := c.Request.URL.Query().Get("name")

	planDefinitions, err := p.RetrievePlanDefinition(ctx, searchParam)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, planDefinitions)
}

// UploadMedia uploads media to GCS and stores the URL in FHIR attachment
func (p PresentationHandlersImpl) UploadMedia(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "UploadMedia")
	defer span.End()

	input := &dto.MediaInput{
		EncounterID:      c.Request.FormValue("encounterID"),
		ServiceRequestID: c.Request.FormValue("serviceRequestID"),
		File:             c.Request.MultipartForm.File,
	}

	if err := c.Request.ParseMultipartForm(500 << 20); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	var response []*dto.Media

	for _, fileHeaders := range input.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				silgotel.NewLogger(name)
				silgotel.RecordError(span, err)
				jsonErrorResponse(c, http.StatusBadRequest, err)

				return
			}

			defer file.Close()

			contentType := fileHeader.Header.Get("Content-Type")

			output, err := p.usecases.UploadMedia(ctx, input.EncounterID, input.ServiceRequestID, file, contentType)
			if err != nil {
				silgotel.NewLogger(name)
				silgotel.RecordError(span, err)
				jsonErrorResponse(c, http.StatusInternalServerError, err)

				return
			}

			response = append(response, output)
		}
	}

	c.JSON(http.StatusOK, response)
}

// LoadQuestionnaire is used to upload a user defined questionnaire for the purpose of soliciting client data.
func (p PresentationHandlersImpl) LoadQuestionnaire(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "LoadQuestionnaire")
	defer span.End()

	input := domain.FHIRQuestionnaire{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	questionnaire, err := p.CreateQuestionnaire(ctx, &input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, questionnaire)
}

// ListQuestionnaire is used to provide params used to fetch questionnaires
func (p PresentationHandlersImpl) ListQuestionnaire(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListQuestionnaire")
	span.End()

	queryParams := c.Request.URL.Query()

	searchParam := queryParams.Get("searchParam")

	questionnaire, err := p.ListQuestionnaires(ctx, searchParam, &dto.Pagination{})
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, questionnaire)
}

// FetchQuestionnaire is used to provide params used to fetch questionnaire
func (p PresentationHandlersImpl) FetchQuestionnaire(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "FetchQuestionnaire")
	defer span.End()

	queryParams := c.Request.URL.Query()

	searchParam := queryParams.Get("searchParam")

	questionnaire, err := p.usecases.FetchQuestionnaire(ctx, searchParam, &dto.Pagination{})
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, questionnaire)
}

func (p PresentationHandlersImpl) Assessment(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "Assessment")
	defer span.End()

	queryParams := c.Request.URL.Query()

	questionnaireID := queryParams.Get("questionnaire_id")

	if questionnaireID == "" {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, fmt.Errorf("questionnaireID required"))
		jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("questionnaireID required"))

		return
	}

	input := &dto.QuestionnaireResponse{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	questionnaire, err := p.GenerateRiskAssessment(ctx, questionnaireID, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, questionnaire)
}

// RiskAssessment is used to fetch risk assessment details
//
//	@Summary		Used to fetch risk assessment details
//	@Description	RiskAssessment is get the details of a risk assessment
//	@Tags			RiskAssessment
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id	path		string				true	"Risk Assessment ID"
//	@Success		200	{object}	dto.RiskAssessment	"OK"
//	@Failure		400	{object}	APIResponse			"Error: Bad Request"
//	@Failure		500	{object}	APIResponse			"Error: InternalServerError"
//	@Failure		404	{object}	APIResponse			"Error: Not Found"
//	@Failure		401	{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/risk-assessment/{id} [get]
func (p PresentationHandlersImpl) RiskAssessment(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "RiskAssessment")
	defer span.End()

	id := c.Param("id")
	if id == "" {
		err := fmt.Errorf("missing observation id")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	riskAssessment, err := p.GetRiskAssessmentByID(ctx, id)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, riskAssessment)
}

// ListRiskAssessment returns the list of risk assessments
func (p PresentationHandlersImpl) ListRiskAssessment(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListRiskAssessment")
	defer span.End()

	queryParams := c.Request.URL.Query()

	var (
		dateString    = queryParams.Get("date")
		patientID     = queryParams.Get("patientID")
		encounterID   = queryParams.Get("encounterID")
		first         = queryParams.Get("first")
		last          = queryParams.Get("last")
		before        = queryParams.Get("before")
		after         = queryParams.Get("after")
		screeningType = queryParams.Get("screeningType")
		result        = queryParams.Get("result")
		searchID      = queryParams.Get("searchID")
	)

	var (
		formattedDate *scalarutils.Date
		err           error
	)

	if dateString != "" {
		formattedDate, err = utils.ConvertDateStringToDateScalar(time.DateOnly, dateString)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}
	}

	var pagination serverutils.PaginationInput

	if first != "" {
		firstVal, err := strconv.Atoi(first)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid 'first' value: %w", err))

			return
		}

		pagination.First = &firstVal

		if after != "" {
			if searchID == "" {
				silgotel.NewLogger(name)
				silgotel.RecordError(span, err)
				jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("'searchID' is required when using 'after'"))

				return
			}

			pagination.After = &after
		}
	}

	if last != "" {
		lastVal, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid 'last' value: %w", err))

			return
		}

		pagination.Last = &lastVal

		if before != "" {
			if searchID == "" {
				silgotel.NewLogger(name)
				silgotel.RecordError(span, err)
				jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("'searchID' is required when using 'before'"))

				return
			}

			pagination.Before = &before
		}
	}

	filter := &dto.RiskAssessmentFilterInput{
		FilterInput: dto.FilterInput{
			PatientID:   patientID,
			EncounterID: encounterID,
			Date:        formattedDate,
		},
		Result: result,
	}

	if screeningType != "" {
		if !dto.ScreeningTypeEnum(screeningType).IsValid() {
			err := fmt.Errorf("invalid screening type: %s", screeningType)

			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		filter.ScreeningType = dto.ScreeningTypeEnum(dto.ScreeningTypeEnum(screeningType))
	}

	riskAssessments, err := p.usecases.ListRiskAssessment(ctx, searchID, filter, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, riskAssessments)
}

// ListTasks returns the list of tasks
//
// ListTasks is used to list available task depending on the provided filters
//
//	@Summary		read-instances: List Task is used to list available task depending on the provided filters
//	@Description	List Task is used to list available task depending on the provided filters
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//
//	@Param			date		query		string						false	"Date"
//	@Param			patientID	query		string						false	"Patient ID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Param			encounterID	query		string						false	"Encounter ID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Param			first		query		string						false	"First"
//	@Param			last		query		string						false	"Last"
//	@Param			before		query		string						false	"Before"
//	@Param			after		query		string						false	"After"
//	@Param			type		query		string						false	"Type"
//	@Param			status		query		string						false	"Status"
//
//	@Success		200			{object}	dto.TaskOutputConnection	"Tasks"
//	@Failure		400			{object}	APIResponse					"Error: Bad Request"
//	@Failure		404			{object}	APIResponse					"Error: Not Found"
//	@Failure		401			{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/task [get]
func (p PresentationHandlersImpl) ListTasks(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListTasks")
	defer span.End()

	queryParams := c.Request.URL.Query()

	var (
		dateString  = queryParams.Get("date")
		patientID   = queryParams.Get("patientID")
		encounterID = queryParams.Get("encounterID")
		limit       = queryParams.Get("first")
		last        = queryParams.Get("last")
		before      = queryParams.Get("before")
		after       = queryParams.Get("after")
		category    = queryParams.Get("type")
		status      = queryParams.Get("status")

		patientSearch = queryParams.Get("patient")
	)

	var formattedDate *scalarutils.Date

	if dateString != "" {
		parsedTime, err := time.Parse(time.DateOnly, dateString)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		year, month, day := parsedTime.Date()

		formattedDate = &scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		}
	}

	var pagination dto.Pagination

	if limit != "" {
		count, err := strconv.Atoi(limit)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	baseFilters := &dto.TaskFilterInput{
		FilterInput: dto.FilterInput{
			PatientID:   patientID,
			EncounterID: encounterID,
			Date:        formattedDate,
		},
	}

	if category != "" {
		baseFilters.Type = category
	}

	if status != "" {
		if !dto.TaskStatus(status).IsValid() {
			err := fmt.Errorf("invalid task status: %s", status)

			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		baseFilters.Status = dto.TaskStatus(dto.TaskStatus(status).String())
	}

	baseFilters.PatientSearch = patientSearch

	tasks, err := p.usecases.ListTasks(ctx, baseFilters, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GenerateReferralReport handles incoming HTTP requests to generate referral reports in PDF format.
//
//	@Summary		read-instance: Generate referral report
//	@Description	Retrieve referral report by service request ID as PDF file
//	@Tags			Service Request
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			servicerequest	query		string		true	"Service Request ID"	Example(ABC12345)
//	@Success		200				{file}		binary		"Report generated successfully!"
//	@Failure		400				{object}	APIResponse	"Error: Bad Request"
//	@Failure		404				{object}	APIResponse	"Error: Not Found"
//	@Failure		401				{object}	APIResponse	"Error: Not Authorized"
//	@Router			/api/v1/referral-report [get]
func (p PresentationHandlersImpl) GenerateReferralReport(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "GenerateReferralReport")
	defer span.End()

	queryParams := c.Request.URL.Query()

	serviceRequestID := queryParams.Get("servicerequest")

	pdfBytes, err := p.GenerateReferralReportPDF(ctx, serviceRequestID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	fileName := "referral_report.pdf"

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ReferPatient is used to refer patient from one facility to another
//
//	@Summary		create-type: Used to refer patient from one facility to another
//	@Description	Refers a patient from on facility to another
//	@Tags			Service Request
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.ReferralInput	true	"ReferPatientInput"
//	@Success		200		{object}	dto.ServiceRequest	"Service Request"
//	@Failure		400		{object}	APIResponse			"Error: Bad Request"
//	@Failure		404		{object}	APIResponse			"Error: Not Found"
//	@Failure		401		{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/refer [post]
func (p PresentationHandlersImpl) ReferPatient(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ReferPatient")
	defer span.End()

	payload := dto.ReferralInput{}

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	questionnaire, err := p.usecases.ReferPatient(ctx, &payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, questionnaire)
}

// CreateReferral is used to create a new patient referral in FHIR as a service request. This method will mainly be used by advantage backend to save the referral details as part of the patient profile
//
//	@Summary		create-type: Used to create a new patient referral in FHIR as a service request. This method will mainly be used by advantage backend to save the referral details as part of the patient profile
//	@Description	CreateReferral is used to create a new patient referral in FHIR as a service request. This method will mainly be used by advantage backend to save the referral details as part of the patient profile
//	@Tags			Service Request
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.CreateReferralInput	true	"CreateReferralInput"
//	@Success		200		{object}	dto.ServiceRequest		"Service Request"
//	@Failure		400		{object}	APIResponse				"Error: Bad Request"
//	@Failure		500		{object}	APIResponse				"Error: InternalServerError"
//	@Failure		404		{object}	APIResponse				"Error: Not Found"
//	@Failure		401		{object}	APIResponse				"Error: Not Authorized"
//	@Router			/api/v1/referV2 [post]
func (p PresentationHandlersImpl) CreateReferral(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateReferral")
	defer span.End()

	payload := dto.CreateReferralInput{}

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	referral, err := p.usecases.CreateReferral(ctx, &payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, referral)
}

// CreateTestOrder creates a new diagnostic test order based on the provided input.
//
//	@Summary		create-type: Used to create a new diagnostic test order based on the provided input.
//	@Description	CreateTestOrder is used to create a new diagnostic test order based on the provided input.
//	@Tags			Service Request
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.TestOrder		true	"TestOrder"
//	@Success		200		{object}	dto.ServiceRequest	"Service Request"
//	@Failure		400		{object}	APIResponse			"Error: Bad Request"
//	@Failure		404		{object}	APIResponse			"Error: Not Found"
//	@Failure		401		{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/lab-orders/test [post]
func (p PresentationHandlersImpl) CreateTestOrder(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateTestOrder")
	defer span.End()

	payload := dto.TestOrder{}

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	testOrder, err := p.usecases.CreateTestOrder(ctx, &payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, testOrder)
}

// ShareReferralForm is searches for a document reference using the service request ID associated and retrieves the document URL
// (which is the url of the referral form that we want to share) and sends it to the patient via SMS
//
//	@Summary		create-type: Used to search for a document reference using the service request ID associated and retrieves the document URL (which is the url of the referral form that we want to share) and sends it to the patient via SMS
//	@Description	ShareReferralForm searches for a document reference using the service request ID associated and retrieves the document URL (which is the url of the referral form that we want to share) and sends it to the patient via SMS
//	@Tags			Service Request
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.ShareReferralFormInput	true	"ShareReferralFormInput"
//	@Success		200		{boolean}	bool						"Success"
//	@Failure		400		{object}	APIResponse					"Error: Bad Request"
//	@Failure		404		{object}	APIResponse					"Error: Not Found"
//	@Failure		401		{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/referral-form [post]
func (p PresentationHandlersImpl) ShareReferralForm(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ShareReferralForm")
	defer span.End()

	payload := dto.ShareReferralFormInput{}

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	shared, err := p.usecases.ShareReferralForm(ctx, &payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"results": shared})
}

// ReferralDetails returns the referral details of a patient if the patientID is passed
// If a patient ID is not provided, it returns referral details of all patients
func (p PresentationHandlersImpl) ReferralDetails(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ReferralDetails")
	defer span.End()

	queryParams := c.Request.URL.Query()

	patientID := queryParams.Get("patientID")
	status := queryParams.Get("status")
	patientSearch := queryParams.Get("patient")
	dateString := queryParams.Get("date")

	var patientIDPtr *string
	if patientID != "" {
		patientIDPtr = &patientID
	}

	var formattedDate *scalarutils.Date

	if dateString != "" {
		parsedTime, err := time.Parse(time.DateOnly, dateString)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		year, month, day := parsedTime.Date()

		formattedDate = &scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		}
	}

	pagination := &dto.Pagination{}

	if first := c.Query("first"); first != "" {
		firstInt, err := strconv.Atoi(first)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.First = &firstInt
	}

	if last := c.Query("last"); last != "" {
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	}

	pagination.After = c.Query("after")
	pagination.Before = c.Query("before")

	referralInput := &dto.ReferralSearchInput{
		PatientID:     patientIDPtr,
		EncounterID:   nil,
		Date:          formattedDate,
		Pagination:    pagination,
		Status:        domain.ServiceRequestStatusEnum(status),
		PatientSearch: patientSearch,
	}

	referralDetails, err := p.GetPatientReferrals(ctx, referralInput)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, referralDetails)
}

// ListLabOrders is used to fetch a list of available lab orders
func (p PresentationHandlersImpl) ListLabOrders(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListLabOrders")
	defer span.End()

	queryParams := c.Request.URL.Query()

	var (
		dateString  = queryParams.Get("date")
		patientID   = queryParams.Get("patientID")
		encounterID = queryParams.Get("encounterID")
		limit       = queryParams.Get("first")
		last        = queryParams.Get("last")
		before      = queryParams.Get("before")
		after       = queryParams.Get("after")
		category    = queryParams.Get("type")
		status      = queryParams.Get("status")
		facilityID  = queryParams.Get("facilityID")

		patientSearch = queryParams.Get("patient")
	)

	var formattedDate *scalarutils.Date

	if dateString != "" {
		parsedTime, err := time.Parse(time.DateOnly, dateString)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		year, month, day := parsedTime.Date()

		formattedDate = &scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		}
	}

	var pagination dto.Pagination

	if limit != "" {
		count, err := strconv.Atoi(limit)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	baseFilters := &dto.ServiceRequestFilterInput{
		FilterInput: dto.FilterInput{
			PatientID:   patientID,
			EncounterID: encounterID,
			Date:        formattedDate,
		},
	}

	if category != "" {
		baseFilters.Type = category
	}

	if status != "" {
		if !dto.ServiceRequestStatus(status).IsValid() {
			err := fmt.Errorf("invalid service request status: %s", status)

			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		baseFilters.Status = dto.ServiceRequestStatus(status)
	}

	if facilityID != "" {
		baseFilters.FacilityID = facilityID
	}

	baseFilters.PatientSearch = patientSearch

	tasks, err := p.usecases.ListLabOrders(ctx, baseFilters, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, tasks)
}

// DeleteTests provides a room to remove tests that are mistakenly recorded
func (p PresentationHandlersImpl) DeleteTests(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "DeleteTests")
	defer span.End()

	queryParams := c.Request.URL.Query()

	searchParam := queryParams.Get("observation-id")
	if searchParam == "" {
		errMsg := fmt.Errorf("unable to delete test")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, errMsg)
		jsonErrorResponse(c, http.StatusBadRequest, errMsg)

		return
	}

	ok, err := p.DeleteTest(ctx, searchParam)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, ok)
}

// CreateEpisodeOfCare handles incoming http requests to create a new episode for a patient.
//
//	@Summary		create-type: Create a new episode of care
//	@Description	Creates a new episode of care and returns it to the client
//	@Tags			EpisodeOfCare
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			episodeOfCare	body		dto.EpisodeOfCareInput	true	"EpisodeOfCare Input"
//	@Success		201				{object}	dto.EpisodeOfCare		"OK"
//	@Failure		400				{object}	APIResponse				"Error: Bad Request"
//	@Failure		404				{object}	APIResponse				"Error: Not Found"
//	@Failure		401				{object}	APIResponse				"Error: Not Authorized"
//	@Router			/api/v1/episode-of-care [post]
func (p PresentationHandlersImpl) CreateEpisodeOfCare(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ScheduleAppointment")
	defer span.End()

	input := dto.EpisodeOfCareInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := helpers.Validate(input); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	episodeOfCare, err := p.usecases.CreateEpisodeOfCare(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusCreated, episodeOfCare)
}

// CreateCondition handles the creation of Condition FHIR Resource
//
//	@Summary		create-type: Records patient's condition(s)
//	@Description	Creates a FHIR Condition resource
//	@Tags			Conditions
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.ConditionInput	true	"Condition payload"
//	@Success		200		{object}	dto.Condition		"Ok"
//	@Failure		400		{object}	APIResponse			"Error: Bad Request"
//	@Failure		404		{object}	APIResponse			"Error: Not Found"
//	@Failure		401		{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/condition [post]
func (p PresentationHandlersImpl) CreateCondition(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateCondition")
	defer span.End()

	var payload dto.ConditionInput

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	condition, err := p.usecases.CreateCondition(ctx, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusCreated, condition)
}

// OncologyDiagnosis is used to record an oncological condition
//
//	@Summary		Records patient's oncological condition(s)
//	@Description	Creates a FHIR Condition resource
//	@Tags			Conditions
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.OncologyDiagnosisInput	true	"Condition payload"
//	@Success		200		{object}	dto.Condition				"Ok"
//	@Failure		400		{object}	APIResponse					"Error: Bad Request"
//	@Failure		404		{object}	APIResponse					"Error: Not Found"
//	@Failure		401		{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/condition/encounter-diagnosis [post]
func (p PresentationHandlersImpl) OncologyDiagnosis(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "OncologyDiagnosis")
	defer span.End()

	payload := &dto.OncologyDiagnosisInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	condition, err := p.RecordOncologicalDiagnosis(ctx, payload)
	if err != nil {
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, condition)
}

// RecordTreatmentEnrollment is used to retrospectively record confirmed oncological condition and its treatment.
//
//	@Summary		Records patient's oncological condition(s) and its treatment.
//	@Description	Creates a FHIR Condition resource and its treatment.
//	@Tags			Conditions
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.TreatmentEnrollmentInput	true	"Condition payload"
//	@Success		200		{object}	dto.Condition						"Ok"
//	@Failure		400		{object}	APIResponse							"Error: Bad Request"
//	@Failure		404		{object}	APIResponse							"Error: Not Found"
//	@Failure		401		{object}	APIResponse							"Error: Not Authorized"
//	@Router			/api/v1/condition/treatment-enrollment [post]
func (p PresentationHandlersImpl) RecordTreatmentEnrollment(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "RecordTreatmentEnrollment")
	defer span.End()

	payload := &dto.TreatmentEnrollmentInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	condition, err := p.usecases.RecordTreatmentEnrollment(ctx, payload)
	if err != nil {
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, condition)
}

// UpdateTreatmentEnrollment updates the condition, recorded date and enrollment date of an existing
// treatment enrollment.
//
//	@Summary		Updates a patient's treatment enrollment.
//	@Description	Updates the condition, date and/or enrollment date of an existing treatment enrollment (FHIR Condition).
//	@Tags			Conditions
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organization ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id		path		string								true	"Treatment enrollment (Condition) ID"
//	@Param			Payload	body		dto.UpdateTreatmentEnrollmentInput	true	"Fields to update"
//	@Success		200		{object}	dto.Condition						"Ok"
//	@Failure		400		{object}	APIResponse							"Error: Bad Request"
//	@Failure		404		{object}	APIResponse							"Error: Not Found"
//	@Failure		401		{object}	APIResponse							"Error: Not Authorized"
//	@Router			/api/v1/condition/treatment-enrollment/{id} [patch]
func (p PresentationHandlersImpl) UpdateTreatmentEnrollment(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "UpdateTreatmentEnrollment")
	defer span.End()

	id := c.Param("id")

	payload := &dto.UpdateTreatmentEnrollmentInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	condition, err := p.usecases.UpdateTreatmentEnrollment(ctx, id, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, condition)
}

// CreateCondition handles the creation of Condition FHIR Resource
//
//	@Summary		create-type: Creates FHIR AllergyIntolerance
//	@Description	Creates an AllergyIntolerance FHIR resource
//	@Tags			AllergyIntolerance
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.AllergyInput	true	"AllergyIntolerance payload"
//	@Success		200		{object}	dto.Allergy			"Ok"
//	@Failure		400		{object}	APIResponse			"Error: Bad Request"
//	@Failure		404		{object}	APIResponse			"Error: Not Found"
//	@Failure		401		{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/allergyintolerance [post]
func (p PresentationHandlersImpl) CreateAllergyIntolerance(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateAllergyIntolerance")
	defer span.End()

	payload := dto.AllergyInput{}
	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	alleryIntolerrance, err := p.usecases.CreateAllergyIntolerance(ctx, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusCreated, alleryIntolerrance)
}

// PatchEpisodeOfCare handles incoming http requests to update an episode for a patient.
//
//	@Summary		instance-patch: Update an episode of care
//	@Description	Updates an episode of care and returns it to the client
//	@Tags			EpisodeOfCare
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Param			id							path	string	true	"episodeOfCare ID"
//	@Security		OAuth2Password
//	@Param			episodeOfCare	body		dto.EpisodeOfCareInput	true	"EpisodeOfCare Input"
//	@Success		200				{object}	dto.EpisodeOfCare		"OK"
//	@Failure		400				{object}	APIResponse				"Error: Bad Request"
//	@Failure		500				{object}	APIResponse				"Error: InternalServerError"
//	@Failure		404				{object}	APIResponse				"Error: Not Found"
//	@Failure		401				{object}	APIResponse				"Error: Not Authorized"
//	@Router			/api/v1/episode-of-care/{id} [patch]
func (p PresentationHandlersImpl) PatchEpisodeOfCare(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "PatchEpisodeOfCare")
	defer span.End()

	id := c.Param("id")
	input := dto.EpisodeOfCareInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := helpers.Validate(input); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	episodeOfCare, err := p.usecases.PatchEpisodeOfCare(ctx, id, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, episodeOfCare)
}

// EndEpisodeOfCare handles incoming http requests to end an episode for a patient.
//
//	@Summary		instance-patch: End an episode of care
//	@Description	Sets an episode of care status to finished and returns it to the client
//	@Tags			EpisodeOfCare
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Param			id							path	string	true	"episodeOfCare ID"
//	@Security		OAuth2Password
//	@Success		200	{object}	dto.EpisodeOfCare	"OK"
//	@Failure		400	{object}	APIResponse			"Error: Bad Request"
//	@Failure		500	{object}	APIResponse			"Error: InternalServerError"
//	@Failure		404	{object}	APIResponse			"Error: Not Found"
//	@Failure		401	{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/episode-of-care/{id}/end [patch]
func (p PresentationHandlersImpl) EndEpisodeOfCare(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "EndEpisodeOfCare")
	defer span.End()

	input := c.Param("id")

	episodeOfCare, err := p.usecases.EndEpisodeOfCare(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, episodeOfCare)
}

// Create Patient handles incoming HTTP requests to create a new patient
//
//	@Summary		create-type: Create a new patient
//	@Description	Creates a new patient with all of their details
//	@Tags			Patient
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			PatientInput	body		dto.PatientInput	true	"Patient input"
//	@Success		201				{object}	dto.Patient			"Patient"
//	@Failure		400				{object}	APIResponse			"Error: Bad Request"
//	@Failure		500				{object}	APIResponse			"Error: InternalServerError"
//	@Failure		404				{object}	APIResponse			"Error: Not Found"
//	@Failure		401				{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/patient [post]
func (p PresentationHandlersImpl) CreatePatient(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreatePatient")
	defer span.End()

	input := dto.PatientInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := helpers.Validate(input); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	patient, err := p.usecases.CreatePatient(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, patient)
}

// PatchPatient is used to modify a patient basic information
//
//	@Summary		PatchPatient is used to modify a patient basic information
//	@Description	Modifies patient basic information
//	@Tags			Patient
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string				true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string				true	"Clinical-Facility-ID"
//
//	@Param			id							path	string				true	"Patient ID"
//
//	@Param			updateData					body	dto.PatientInput	true	"Patient Input"
//
//	@Security		OAuth2Password
//	@Success		200	{object}	dto.Patient	"Patient"
//	@Failure		400	{object}	APIResponse	"Error: Bad Request"
//	@Failure		500	{object}	APIResponse	"Error: InternalServerError"
//	@Failure		404	{object}	APIResponse	"Error: Not Found"
//	@Failure		401	{object}	APIResponse	"Error: Not Authorized"
//	@Router			/api/v1/patient/{id} [patch]
func (p PresentationHandlersImpl) PatchPatient(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "PatchPatient")
	span.End()

	patientID := c.Param("id")
	payload := dto.PatientInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := helpers.Validate(payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	result, err := p.usecases.PatchPatient(ctx, patientID, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, result)
}

// StartEncounter handles incoming http requests to starts an encounter within an episode and returns its ID
//
//	@Summary		create-type: Start an encounter for an episode
//	@Description	Starts an encounter within an episode of care and returns the encounter ID
//	@Tags			Encounter
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string					true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string					true	"Clinical-Facility-ID"
//	@Param			episodeOfCareID				body	dto.StartEncounterInput	true	"ID of the associated episode"
//	@Security		OAuth2Password
//	@Success		201	{object}	string		"OK"
//	@Failure		400	{object}	APIResponse	"Error: Bad Request"
//	@Failure		500	{object}	APIResponse	"Error: InternalServerError"
//	@Failure		404	{object}	APIResponse	"Error: Not Found"
//	@Failure		401	{object}	APIResponse	"Error: Not Authorized"
//	@Router			/api/v1/encounter [post]
func (p PresentationHandlersImpl) StartEncounter(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "StartEncounter")
	span.End()

	input := dto.StartEncounterInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	err = helpers.Validate(input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	encounterID, err := p.usecases.StartEncounter(ctx, input.EpisodeOfCareID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, gin.H{"results": encounterID})
}

// CreateComposition handles the creation of FHIR composition
//
//	@Summary		create-type: Records a FHIR Composition
//	@Description	Creates FHIR Composition resource
//	@Tags			Composition
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			Payload	body		dto.CompositionInput	true	"CompositionInput payload"
//	@Success		200		{object}	dto.Composition			"Ok"
//	@Failure		400		{object}	APIResponse				"Error: Bad Request"
//	@Failure		500		{object}	APIResponse				"Error: InternalServerError"
//	@Failure		404		{object}	APIResponse				"Error: Not Found"
//	@Failure		401		{object}	APIResponse				"Error: Not Authorized"
//	@Router			/api/v1/compositions [post]
func (p PresentationHandlersImpl) CreateComposition(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateComposition")
	span.End()

	payload := dto.CompositionInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	composition, err := p.usecases.CreateComposition(ctx, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, composition)
}

// PatchComposition is used to modify a compostion resource
//
//	@Summary		instance-patch: Updates a FHIR Composition
//	@Description	Updates a FHIR Composition resource
//	@Tags			Composition
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string						true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string						true	"Clinical-Facility-ID"
//	@Param			id							path	string						true	"composition ID"
//	@Param			updateData					body	dto.PatchCompositionInput	true	"Patch Compostion Payload"
//	@Security		OAuth2Password
//	@Success		200	{object}	dto.Composition	"OK"
//	@Failure		400	{object}	APIResponse		"Error: Bad Request"
//	@Failure		404	{object}	APIResponse		"Error: Not Found"
//	@Failure		500	{object}	APIResponse		"Error: InternalServerError"
//	@Failure		401	{object}	APIResponse		"Error: Not Authorized"
//	@Router			/api/v1/compositions/{id} [patch]
func (p PresentationHandlersImpl) PatchComposition(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "PatchComposition")
	span.End()

	compositionID := c.Param("id")

	payload := dto.PatchCompositionInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	composition, err := p.AppendNoteToComposition(ctx, compositionID, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, composition)
}

// PatchTask is used to modify a task resource
//
//	@Summary		instance-patch: Updates a FHIR Task
//	@Description	Updates a FHIR Task resource
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string				true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string				true	"Clinical-Facility-ID"
//	@Param			id							path	string				true	"Task ID"
//	@Param			updateData					body	dto.PatchTaskInput	true	"Patch Task Payload"
//	@Security		OAuth2Password
//	@Success		200	{bool}		true		"Ok"
//	@Failure		400	{object}	APIResponse	"Error: Bad Request"
//	@Failure		404	{object}	APIResponse	"Error: Not Found"
//	@Failure		500	{object}	APIResponse	"Error: InternalServerError"
//	@Failure		401	{object}	APIResponse	"Error: Not Authorized"
//	@Router			/api/v1/task/{id} [patch]
func (p PresentationHandlersImpl) PatchTask(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "PatchTask")
	span.End()

	taskID := c.Param("id")

	payload := dto.PatchTaskInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	ok, err := p.UpdateTask(ctx, taskID, &payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"results": ok})
}

// PatchEncounter handles incoming http requests to update the status of an encounter
//
//	@Summary		instance-patch: Updates the status of an encounter
//	@Description	Updates the status of an encounter and returns the encounter
//	@Tags			Encounter
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string				true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string				true	"Clinical-Facility-ID"
//	@Param			id							path	string				true	"Encounter ID"
//	@Param			EncounterInput				body	dto.EncounterInput	true	"Encounter Input"
//	@Security		OAuth2Password
//	@Success		200	{object}	dto.Encounter	"OK"
//	@Failure		400	{object}	APIResponse		"Error: Bad Request"
//	@Failure		404	{object}	APIResponse		"Error: Not Found"
//	@Failure		500	{object}	APIResponse		"Error: InternalServerError"
//	@Failure		401	{object}	APIResponse		"Error: Not Authorized"
//	@Router			/api/v1/encounter/{id} [patch]
func (p PresentationHandlersImpl) PatchEncounter(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "PatchEncounter")
	defer span.End()

	id := c.Param("id")
	input := dto.EncounterInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	err = helpers.Validate(input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	encounter, err := p.usecases.PatchEncounter(ctx, id, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, encounter)
}

// @Summary		read-instance: Used to generate a screening report
// @Description	Refers to a screening report
// @Tags			Encounter
// @Accept			json
// @Produce		json
// @Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
// @Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
// @Security		OAuth2Password
// @Param			encounterID	query		string				true	"encounterID"	Example(ABC12345)
// @Success		200			{object}	dto.ScreeningReport	"Screening Report"
// @Failure		400			{object}	APIResponse			"Error: Bad Request"
// @Failure		404			{object}	APIResponse			"Error: Not Found"
// @Failure		401			{object}	APIResponse			"Error: Not Authorized"
// @Router			/api/v1/screening-report [get]
func (p PresentationHandlersImpl) ScreeningReport(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ScreeningReport")
	defer span.End()

	queryParams := c.Request.URL.Query()

	encounterID := queryParams.Get("encounterID")
	status := queryParams.Get("status")

	questionnaire, err := p.GetScreeningReport(ctx, encounterID, domain.ServiceRequestStatusEnum(status))
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, questionnaire)
}

// @Summary		instance-delete: Deletes a Patient
// @Description	Deletes a Patient
// @Tags			Patient
// @Accept			json
// @Produce		json
// @Produce		json
// @Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
// @Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
// @Security		OAuth2Password
// @Param			id	path		string		true	"Patient ID"
// @Success		200	{object}	boolean		"OK"
// @Failure		400	{object}	APIResponse	"Error: Bad Request"
// @Failure		404	{object}	APIResponse	"Error: Not Found"
// @Failure		401	{object}	APIResponse	"Error: Not Authorized"
// @Success		200	{object}	boolean		"OK"
// @Failure		400	{object}	APIResponse	"Error: Bad Request"
// @Failure		404	{object}	APIResponse	"Error: Not Found"
// @Failure		401	{object}	APIResponse	"Error: Not Authorized"
// @Router			/api/v1/patient/{id} [delete]
func (p PresentationHandlersImpl) DeletePatient(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ScheduleAppointment")
	defer span.End()

	patientID := c.Param("id")

	output, err := p.usecases.DeletePatient(ctx, patientID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusNoContent, gin.H{"results": output})
}

// EndEncounter handles incoming http requests to end an encounter
//
//	@Summary		instance-patch: End an encounter within an episode
//	@Description	Ends an encounter within an episode and returns a boolean of the status
//	@Tags			Encounter
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Param			id							path	string	true	"ID of the encounter"
//	@Security		OAuth2Password
//	@Success		200	{object}	boolean		"OK"
//	@Failure		400	{object}	APIResponse	"Error: Bad Request"
//	@Failure		404	{object}	APIResponse	"Error: Not Found"
//	@Failure		401	{object}	APIResponse	"Error: Not Authorized"
//	@Router			/api/v1/encounter/{id}/end [patch]
func (p PresentationHandlersImpl) EndEncounter(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "EndEncounter")
	defer span.End()

	id := c.Param("id")

	status, err := p.usecases.EndEncounter(ctx, id)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"results": status})
}

// @Summary		instance-patch: Ends Screening
// @Description	Ends a Screening
// @Tags			Encounter
// @Accept			json
// @Produce		json
// @Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
// @Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
// @Security		OAuth2Password
// @Param			id	path		string		true	"encounter ID"
// @Success		200	{bool}		true		"Ok"
// @Failure		400	{object}	APIResponse	"Error: Bad Request"
// @Failure		404	{object}	APIResponse	"Error: Not Found"
// @Failure		401	{object}	APIResponse	"Error: Not Authorized"
// @Router			/api/v1/end-screening/{id} [patch]
func (p PresentationHandlersImpl) EndScreening(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "EndScreening")
	defer span.End()

	id := c.Param("id")

	if id == "" {
		err := fmt.Errorf("missing encounter ID")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	status, err := p.usecases.EndScreening(ctx, id)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"results": status})
}

// EncounterAssociatedResources is used to retrieve all resources assiociated with a given encounter
//
//	@Summary		read-instance: Fetch encounter associated resources
//	@Description	Retrives the resources associated to a specific encounter
//	@Tags			Encounter
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id	path		string		true	"encounter ID"
//	@Success		200	{bool}		true		"Ok"
//	@Failure		400	{object}	APIResponse	"Error: Bad Request"
//	@Failure		404	{object}	APIResponse	"Error: Not Found"
//	@Failure		401	{object}	APIResponse	"Error: Not Authorized"
//	@Router			/api/v1/encounter/{id}/associated-resources [get]
func (p PresentationHandlersImpl) EncounterAssociatedResources(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "EncounterAssociatedResources")
	span.End()

	id := c.Param("id")

	if id == "" {
		err := fmt.Errorf("missing encounter ID")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	status, err := p.GetEncounterAssociatedResources(ctx, id)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, status)
}

// RecordConsent handles incoming http requests to record a user's consent
//
//	@Summary		create-type: Record user consent
//	@Description	Record user consent and return the consent output
//	@Tags			Consent
//	@Accept			json
//	@Produce		json
//	@Summary		Record user consent
//	@Description	Record user consent and return the consent output
//	@Tags			Consent
//	@Accept			json
//
// Produce 			json
//
//	@Param			Clinical-Organization-ID	header	string				true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string				true	"Clinical-Facility-ID"
//	@Param			consentInput				body	dto.ConsentInput	true	"consent input"
//	@Security		OAuth2Password
//	@Success		201	{object}	dto.ConsentOutput	"OK"
//	@Failure		400	{object}	APIResponse			"Error: Bad Request"
//	@Failure		404	{object}	APIResponse			"Error: Not Found"
//	@Failure		401	{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/consent [post]
func (p PresentationHandlersImpl) RecordConsent(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "RecordConsent")
	span.End()

	input := dto.ConsentInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	err = helpers.Validate(input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	consentOuput, err := p.usecases.RecordConsent(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, consentOuput)
}

// PatientReferralDetails is used to retrieve patient referral by ID
//
//	@Summary		read-instance: Get Patient Referral by service request ID
//	@Description	Fetch Patient Referral details given service request ID
//	@Tags			Service Request
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			serviceRequestID	query		string							true	"Service Request ID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Success		200					{object}	domain.PatientReferralDetails	"OK"
//	@Failure		400					{object}	APIResponse						"Error: Bad Request"
//	@Failure		500					{object}	APIResponse						"Error: InternalServerError"
//	@Failure		404					{object}	APIResponse						"Error: Not Found"
//	@Failure		401					{object}	APIResponse						"Error: Not Authorized"
//	@Router			/api/v1/referral [get]
func (p PresentationHandlersImpl) PatientReferralDetailsByID(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "PatientReferralDetailsByID")
	span.End()

	queryParams := c.Request.URL.Query()
	serviceRequestID := queryParams.Get("serviceRequestID")

	referralDetails, err := p.GetPatientReferralDetails(ctx, serviceRequestID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, referralDetails)
}

// RecordObservation handles incoming http requests to record an observation

//	@Summary		create-type: General endpoint to record different type of observation concepts
//	@Description	Record a type of observation defined by the concept
//	@Tags			Observations
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string					true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string					true	"Clinical-Facility-ID"
//	@Param			observationInput			body	dto.ObservationInput	true	"observation input"
//	@Security		OAuth2Password
//	@Success		201	{object}	dto.Observation	"Observation output"
//	@Failure		400	{object}	APIResponse		"Error: Bad Request"
//	@Failure		404	{object}	APIResponse		"Error: Not Found"
//	@Failure		401	{object}	APIResponse		"Error: Not Authorized"
//	@Router			/api/v1/observations [post]
//	@Summary		General endpoint to record different type of observation concepts
//	@Description	Record a type of observation defined by the concept
//	@Tags			Observations
//	@Accept			json
//
// Produce 			json
//
//	@Param			Clinical-Organization-ID	header	string					true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string					true	"Clinical-Facility-ID"
//	@Param			observationInput			body	dto.ObservationInput	true	"observation input"
//	@Security		OAuth2Password
//	@Success		201	{object}	dto.Observation	"Observation output"
//	@Failure		400	{object}	APIResponse		"Error: Bad Request"
//	@Failure		404	{object}	APIResponse		"Error: Not Found"
//	@Failure		401	{object}	APIResponse		"Error: Not Authorized"
//	@Router			/api/v1/observations [post]
func (p PresentationHandlersImpl) RecordObservation(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "RecordObservation")
	span.End()

	input := dto.ObservationInput{}

	err := c.BindJSON(&input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	err = helpers.Validate(input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if input.Concept == "" {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, fmt.Errorf("concept is missing"))
		jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("concept is missing"))

		return
	}

	observation, err := p.RecordObservationV2(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, observation)
}

// @Summary		create-type: Records test results
// @Description	Record test results
// @Tags			Diagnostics
// @Accept			json
// @Produce		json
// @Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
// @Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
// @Security		OAuth2Password
// @Param			TestResultInput	body		dto.TestResultInput		true	"Record test results"
// @Success		200				{object}	dto.DiagnosticReport	"Diagnostic report"
// @Failure		400				{object}	APIResponse				"Error: Bad Request"
// @Failure		404				{object}	APIResponse				"Error: Not Found"
// @Failure		401				{object}	APIResponse				"Error: Not Authorized"
// @Success		200				{object}	dto.DiagnosticReport	"Diagnostic report"
// @Failure		400				{object}	APIResponse				"Error: Bad Request"
// @Failure		404				{object}	APIResponse				"Error: Not Found"
// @Failure		401				{object}	APIResponse				"Error: Not Authorized"
// @Router			/api/v1/test-results [post]
func (p PresentationHandlersImpl) RecordTestResult(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "RecordTestResult")
	span.End()

	input := dto.TestResultInput{}

	if err := c.BindJSON(&input); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	results, err := p.usecases.RecordTestResult(ctx, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, results)
}

// @Summary		create-type: Create lab order results
// @Description	Creates lab order results
// @Tags			Diagnostics
// @Accept			json
// @Produce		json
// @Produce		json
// @Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
// @Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
// @Security		OAuth2Password
// @Param			TestOrderResult	body		dto.TestOrderResult		true	"Create lab order results"
// @Success		200				{object}	dto.DiagnosticReport	"Diagnostic report"
// @Failure		400				{object}	APIResponse				"Error: Bad Request"
// @Failure		404				{object}	APIResponse				"Error: Not Found"
// @Failure		401				{object}	APIResponse				"Error: Not Authorized"
// @Success		200				{object}	dto.DiagnosticReport	"Diagnostic report"
// @Failure		400				{object}	APIResponse				"Error: Bad Request"
// @Failure		404				{object}	APIResponse				"Error: Not Found"
// @Failure		401				{object}	APIResponse				"Error: Not Authorized"
// @Router			/api/v1/lab-order [post]
func (p PresentationHandlersImpl) CreateLabOrderResult(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateLabOrderResult")
	span.End()

	input := dto.TestOrderResult{}

	if err := c.BindJSON(&input); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	results, err := p.usecases.CreateLabOrderResult(ctx, &input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, results)
}

// ScheduleAppointment is used to schedule an appointment for a given patient
//
//	@Summary		create-type: Sets up a new patient schedule
//	@Description	Creates new schedule for a patient. It'll show up in advantage
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			AppointmentPayload	body		dto.ScheduleAppointmentPayload	true	"Schedule an appointment"
//	@Success		200					{boolean}	bool							"Success"
//	@Failure		400					{object}	APIResponse						"Error: Bad Request"
//	@Failure		404					{object}	APIResponse						"Error: Not Found"
//	@Failure		401					{object}	APIResponse						"Error: Not Authorized"
//	@Router			/api/v1/appointment [post]
func (p PresentationHandlersImpl) ScheduleAppointment(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ScheduleAppointment")
	span.End()

	input := dto.ScheduleAppointmentPayload{}
	if err := c.BindJSON(&input); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	results, err := p.usecases.ScheduleAppointment(ctx, &input.AppointmentInput, &input.HeadersInput)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// @Summary		create-type: Recording lab test
// @Description	Records lab test results
// @Tags			Diagnostics
// @Accept			json
// @Produce		json
// @Produce		json
// @Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
// @Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
// @Security		OAuth2Password
// @Param			PSAInput	body		dto.TestInput			true	"Record Lab tests"
// @Success		200			{object}	dto.DiagnosticReport	"Recording lab tests"
// @Failure		400			{object}	APIResponse				"Error: Bad Request"
// @Failure		404			{object}	APIResponse				"Error: Not Found"
// @Failure		401			{object}	APIResponse				"Error: Not Authorized"
// @Param			PSAInput	body		dto.TestInput			true	"Record Lab tests"
// @Success		200			{object}	dto.DiagnosticReport	"Recording lab tests"
// @Failure		400			{object}	APIResponse				"Error: Bad Request"
// @Failure		404			{object}	APIResponse				"Error: Not Found"
// @Failure		401			{object}	APIResponse				"Error: Not Authorized"
// @Router			/api/v1/tests [post]
func (p PresentationHandlersImpl) RecordTests(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "RecordTests")
	span.End()

	payload := dto.TestInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	results, err := p.usecases.RecordTests(ctx, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, results)
}

// CreatePrescription is used to create a new prescription
//
//	@Summary		create-type: CreatePrescription is used to create a new prescription
//	@Description	Create a new prescription
//	@Tags			Medication
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			PrescriptionInput	body		dto.PrescriptionInput		true	"Prescription Input"
//	@Success		200					{object}	dto.MedicationRequestOutput	"Array of Medication Request"
//	@Failure		400					{object}	APIResponse					"Error: Bad Request"
//	@Failure		404					{object}	APIResponse					"Error: Not Found"
//	@Failure		401					{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/medication/prescription [post]
func (p PresentationHandlersImpl) CreatePrescription(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreatePrescription")
	span.End()

	payload := dto.PrescriptionInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	results, err := p.usecases.CreatePrescription(ctx, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, results)
}

// UpdateMedication is used to modify a patient medications
//
//	@Summary		instance-patch: Updates a patient medication
//	@Description	Updates a medication resource
//	@Tags			Medication
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string						true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string						true	"Clinical-Facility-ID"
//	@Param			id							path	string						true	"Medication ID"
//	@Param			updateData					body	dto.PatchMedicationInput	true	"Medication Request"
//	@Security		OAuth2Password
//	@Success		200	{object}	dto.MedicationRequestOutput	"Medication Request"
//	@Failure		400	{object}	APIResponse					"Error: Bad Request"
//	@Failure		404	{object}	APIResponse					"Error: Not Found"
//	@Failure		401	{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/medication/{id} [patch]
func (p PresentationHandlersImpl) UpdateMedication(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "UpdateMedication")
	span.End()

	medicationID := c.Param("id")

	payload := dto.PatchMedicationInput{}

	if err := helpers.Validate(payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	ok, err := p.PatchMedicationRequests(ctx, medicationID, payload.Status)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusCreated, ok)
}

// GetLabOrderByID is used to retrieve lab order given its id
//
//	@Summary		read-instance: Retrieves a lab order given its ID (Service Request ID)
//	@Description	Retrieves a lab order given its service request ID
//	@Tags			Service Request
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id	path		string				true	"Service Request ID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Success		200	{object}	dto.ServiceRequest	"OK"
//	@Failure		400	{object}	APIResponse			"Error: Bad Request"
//	@Failure		404	{object}	APIResponse			"Error: Not Found"
//	@Failure		401	{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/lab-orders/{id} [get]
func (p PresentationHandlersImpl) GetLabOrderByID(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ScheduleAppointment")
	span.End()

	id := c.Param("id")

	if id == "" {
		err := fmt.Errorf("missing service request ID")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	labOrder, err := p.GetLabOrder(ctx, id)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, labOrder)
}

// GetMedicalData is used to retrieve patient's medical data
//
//	@Summary		read-instance: Get Patient's medical data
//	@Description	Fetches a patient's medical data
//	@Tags			Patient
//
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			patientID	query		string			true	"Patient ID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Success		200			{object}	dto.MedicalData	"OK"
//	@Failure		400			{object}	APIResponse		"Error: Bad Request"
//	@Failure		404			{object}	APIResponse		"Error: Not Found"
//	@Failure		401			{object}	APIResponse		"Error: Not Authorized"
//	@Router			/api/v1/medical-data/{id} [get]
func (p PresentationHandlersImpl) GetMedicalData(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "GetMedicalData")
	defer span.End()

	patientID := c.Param("id")
	if patientID == "" {
		silgotel.NewLogger(name)

		errMsg := fmt.Errorf("patient id cannot be null")
		silgotel.RecordError(span, errMsg)
		jsonErrorResponse(c, http.StatusBadRequest, errMsg)

		return
	}

	medicalData, err := p.usecases.GetMedicalData(ctx, patientID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, medicalData)
}

// GetEpisodeOfCare handles incoming http requests to get an episode's details for a patient.
//
//	@Summary		read-instance: Get details about an episode of care
//	@Description	Get details about an episode of care and returns it
//	@Tags			EpisodeOfCare
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Param			id							path	string	true	"episodeOfCare ID"
//	@Security		OAuth2Password
//	@Success		200	{object}	dto.EpisodeOfCare	"OK"
//	@Failure		400	{object}	APIResponse			"Error: Bad Request"
//	@Failure		404	{object}	APIResponse			"Error: Not Found"
//	@Failure		401	{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/episode-of-care/{id} [get]
func (p PresentationHandlersImpl) GetEpisodeOfCare(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "GetEpisodeOfCare")
	span.End()

	ID := c.Param("id")

	if ID == "" {
		err := fmt.Errorf("episodeOfCare id is missing")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	episodeOfCare, err := p.usecases.GetEpisodeOfCare(ctx, ID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, episodeOfCare)
}

// SimpleQuestionnaireResponse handles the request to retrieve a simplified version of a QuestionnaireResponse
//
//	@Summary		read-instance: handles the request to retrieve a simplified version of a QuestionnaireResponse
//	@Description	This endpoint processes the full questionnaire response data and returns a minimal, streamlined version that includes only essential information
//	@Tags			QuestionnaireResponse
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id	path		string									true	"QuestionnaireResponse ID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Success		200	{object}	[]domain.SimpleQuestionnaireResponse	"OK"
//	@Failure		400	{object}	APIResponse								"Error: Bad Request"
//	@Failure		404	{object}	APIResponse								"Error: Not Found"
//	@Failure		401	{object}	APIResponse								"Error: Not Authorized"
//	@Router			/api/v1/questionnaire-response/{id} [get]
func (p PresentationHandlersImpl) SimpleQuestionnaireResponse(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "SimpleQuestionnaireResponse")
	span.End()

	var questionnaireResponseID string
	if questionnaireResponseID = c.Param("id"); questionnaireResponseID == "" {
		err := fmt.Errorf("questionnaire response ID is missing")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	responses, err := p.usecases.SimpleQuestionnaireResponse(ctx, questionnaireResponseID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadGateway, err)

		return
	}

	c.JSON(http.StatusOK, responses)
}

// ListPatientConditions handles incoming http requests to list patient's conditions
//
//	@Summary		read-instances: List patient conditions
//	@Description	List the conditions of the patient given by the ID
//	@Tags			Conditions
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			patient_id		query		string					true	"Patient ID"
//	@Param			encounter_id	query		string					false	"Encounter ID"
//	@Param			date			query		string					false	"Date"	Example(2025-01-31)
//	@Param			limit			query		string					false	"Forward pagination argument"
//	@Param			after			query		string					false	"Forward pagination argument"
//	@Param			last			query		string					false	"Backward pagination argument"
//	@Param			before			query		string					false	"Backward pagination argument"
//	@Param			strategy		query		string					false	"Listing strategy. Supported value: 'linkage' returns only records captured via treatment-enrollment linkage"
//	@Success		200				{object}	dto.ConditionConnection	"OK"
//	@Failure		400				{object}	APIResponse				"Error: Bad Request"
//	@Failure		404				{object}	APIResponse				"Error: Not Found"
//	@Failure		401				{object}	APIResponse				"Error: Not Authorized"
//	@Router			/api/v1/conditions [get]
func (p PresentationHandlersImpl) ListPatientConditions(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListPatientConditions")
	defer span.End()

	queryParams := c.Request.URL.Query()

	var (
		patientID   = queryParams.Get("patient_id")
		encounterID = queryParams.Get("encounter_id")
		date        = queryParams.Get("date")
		limit       = queryParams.Get("limit")
		before      = queryParams.Get("before")
		after       = queryParams.Get("after")
		last        = queryParams.Get("last")
		strategy    = queryParams.Get("strategy")
	)

	var encounterInput *string
	if encounterID != "" {
		encounterInput = &encounterID
	}

	var formattedDate *scalarutils.Date

	if date != "" {
		parsedTime, err := time.Parse(time.DateOnly, date)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		year, month, day := parsedTime.Date()

		formattedDate = &scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		}
	}

	var pagination dto.Pagination

	if limit != "" {
		count, err := strconv.Atoi(limit)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	conditions, err := p.usecases.ListPatientConditions(ctx, patientID, encounterInput, formattedDate, strategy, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, conditions)
}

// FetchMedicationRequestByID is used to retrieve patient's medical data
//
//	@Summary		read-instance: Get Medication request by ID
//	@Description	Fetches medication request
//	@Tags			Medication
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id	path		string						true	"MedicationRequestID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Success		200	{object}	dto.MedicationRequestOutput	"OK"
//	@Failure		400	{object}	APIResponse					"Error: Bad Request"
//	@Failure		404	{object}	APIResponse					"Error: Not Found"
//	@Failure		401	{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/medication-request/{id} [get]
func (p PresentationHandlersImpl) FetchMedicationRequestByID(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "FetchMedicationRequestByID")
	span.End()

	var medicationRequestID string
	if medicationRequestID = c.Param("id"); medicationRequestID == "" {
		err := fmt.Errorf("missing medication request ID")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	medicationRequest, err := p.usecases.FetchMedicationRequestByID(ctx, medicationRequestID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, medicationRequest)
}

// ListPatientCompositions handles incoming http requests to list patient's compositions
//
//	@Summary		read-instances: List patient compositions
//	@Description	List the compositions of the patient given by the ID
//	@Tags			Composition
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			patient_id		query		string						true	"Patient ID"
//	@Param			encounter_id	query		string						false	"Encounter ID"
//	@Param			date			query		string						false	"Date"	Example(2025-01-31)
//	@Param			limit			query		string						false	"Forward pagination argument"
//	@Param			after			query		string						false	"Forward pagination argument"
//	@Param			last			query		string						false	"Backward pagination argument"
//	@Param			before			query		string						false	"Backward pagination argument"
//	@Success		200				{object}	dto.CompositionConnection	"OK"
//	@Failure		400				{object}	APIResponse					"Error: Bad Request"
//	@Failure		404				{object}	APIResponse					"Error: Not Found"
//	@Failure		401				{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/compositions [get]
func (p PresentationHandlersImpl) ListPatientCompositions(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListPatientCompositions")
	span.End()

	queryParams := c.Request.URL.Query()

	var (
		patientID   = queryParams.Get("patient_id")
		encounterID = queryParams.Get("encounter_id")
		date        = queryParams.Get("date")
		limit       = queryParams.Get("limit")
		before      = queryParams.Get("before")
		after       = queryParams.Get("after")
		last        = queryParams.Get("last")
	)

	var encounterInput *string
	if encounterID != "" {
		encounterInput = &encounterID
	}

	var formattedDate *scalarutils.Date

	if date != "" {
		parsedTime, err := time.Parse(time.DateOnly, date)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		year, month, day := parsedTime.Date()

		formattedDate = &scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		}
	}

	var pagination dto.Pagination

	if limit != "" {
		count, err := strconv.Atoi(limit)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	compositions, err := p.usecases.ListPatientCompositions(ctx, patientID, encounterInput, formattedDate, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, compositions)
}

// SearchAllergy handles http requests to search for an allergy
//
//	@Summary		search-type: Search for allergy
//	@Description	Search for allergy by name
//	@Tags			AllergyIntolerance
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			name	query		string						true	"Allergy Name"
//	@Param			limit	query		string						false	"Forward pagination argument"
//	@Param			after	query		string						false	"Forward pagination argument"
//	@Param			last	query		string						false	"Backward pagination argument"
//	@Param			before	query		string						false	"Backward pagination argument"
//	@Success		200		{object}	dto.TerminologyConnection	"OK"
//	@Failure		400		{object}	APIResponse					"Error: Bad Request"
//	@Failure		404		{object}	APIResponse					"Error: Not Found"
//	@Failure		401		{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/allergyintolerance/search [get]
func (p PresentationHandlersImpl) SearchAllergy(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "SearchAllergy")
	span.End()

	queryParams := c.Request.URL.Query()

	var (
		name   = queryParams.Get("name")
		limit  = queryParams.Get("limit")
		before = queryParams.Get("before")
		after  = queryParams.Get("after")
		last   = queryParams.Get("last")
	)

	var pagination dto.Pagination

	if limit != "" {
		count, err := strconv.Atoi(limit)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	allergy, err := p.usecases.SearchAllergy(ctx, name, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, allergy)
}

// GetAllergyIntolerance handles http requests to get a patient's allergy intolerance
//
//	@Summary		read-instance: Get patient's allergy information
//	@Description	Fetches a patient's allergy intolerance
//	@Tags			AllergyIntolerance
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id	path		string		true	"AllergyIntolerance ID"
//	@Success		200	{object}	dto.Allergy	"OK"
//	@Failure		400	{object}	APIResponse	"Error: Bad Request"
//	@Failure		404	{object}	APIResponse	"Error: Not Found"
//	@Failure		401	{object}	APIResponse	"Error: Not Authorized"
//	@Router			/api/v1/allergyintolerance/{id} [get]
func (p PresentationHandlersImpl) GetAllergyIntolerance(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "GetAllergyIntolerance")
	span.End()

	allergyID := c.Param("id")
	if allergyID == "" {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, fmt.Errorf("allergyID is missing"))
		jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("allergyID is missing"))

		return
	}

	allergy, err := p.usecases.GetAllergyIntolerance(ctx, allergyID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, allergy)
}

// ListMedicationRequests is used to retrieve a paginated list of available medication requests given some parameters
//
//	@Summary		read-instances: Retrieves a list of all available medication requests
//	@Description	Retrieves is used to retrieve a paginated list of available medication requests given some parameters
//	@Tags			Medication
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			date		path		string							false	"Date"
//	@Param			patientID	path		string							false	"Patient ID"
//	@Param			encounterID	path		string							false	"encounter ID"
//	@Param			first		path		string							false	"first"
//	@Param			last		path		string							false	"last"
//	@Param			before		path		string							false	"before"
//	@Param			after		path		string							false	"after"
//	@Param			status		path		string							false	"status"
//	@Success		200			{object}	dto.MedicationRequestConnection	"OK"
//	@Failure		400			{object}	APIResponse						"Error: Bad Request"
//	@Failure		404			{object}	APIResponse						"Error: Not Found"
//	@Failure		401			{object}	APIResponse						"Error: Not Authorized"
//	@Router			/api/v1/medication-request/ [get]
func (p PresentationHandlersImpl) ListMedicationRequests(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListMedicationRequests")
	span.End()

	queryParams := c.Request.URL.Query()

	var (
		date        = queryParams.Get("date")
		patientID   = queryParams.Get("patientID")
		encounterID = queryParams.Get("encounterID")
		first       = queryParams.Get("first")
		last        = queryParams.Get("last")
		before      = queryParams.Get("before")
		after       = queryParams.Get("after")
		status      = queryParams.Get("status")
	)

	var formattedDate *scalarutils.Date

	if date != "" {
		parsedTime, err := utils.ConvertDateStringToDateScalar(time.DateOnly, date)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		formattedDate = parsedTime
	}

	var pagination dto.Pagination

	if first != "" {
		count, err := strconv.Atoi(first)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	searchFilter := dto.MedicationRequestFilterInput{
		FilterInput: dto.FilterInput{
			PatientID:   patientID,
			EncounterID: encounterID,
			Date:        formattedDate,
		},
	}

	if status != "" {
		if !domain.MedicationRequestStatus(status).IsValid() {
			err := fmt.Errorf("invalid medication request status: %s", status)

			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		searchFilter.Status = domain.MedicationRequestStatus(status)
	}

	referralDetails, err := p.usecases.ListMedicationRequests(ctx, &searchFilter, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, referralDetails)
}

// GetTaskByID is used to task by ID
//
//	@Summary		read-instance: Retrieves a Task by its ID
//	@Description	Retrieves a Task given a task ID
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			taskID	path		string			true	"Task ID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Success		200		{object}	dto.TaskOutput	"OK"
//	@Failure		400		{object}	APIResponse		"Error: Bad Request"
//	@Failure		404		{object}	APIResponse		"Error: Not Found"
//	@Failure		401		{object}	APIResponse		"Error: Not Authorized"
//	@Router			/api/v1/task/{id} [get]
func (p PresentationHandlersImpl) GetTaskByID(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "GetTaskByID")
	span.End()

	var taskID string

	if taskID = c.Param("id"); taskID == "" {
		err := fmt.Errorf("missing task id")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	task, err := p.usecases.GetTaskByID(ctx, taskID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, task)
}

// GetObservations handles incoming http requests to get entries for an observation category
//
//	@Summary		General endpoint to get entries for an observation category
//	@Description	Get entries for a category of observation
//	@Tags			Observations
//	@Accept			json
//
// Produce json
//
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			patient_id		query		string						true	"Patient ID"
//	@Param			encounter_id	query		string						false	"Encounter ID"
//	@Param			concept			query		dto.ObservationConceptEnum	false	"Observation concept"
//	@Param			date			query		string						false	"Date"	Example(2025-01-31)
//	@Param			limit			query		string						false	"Forward pagination argument"
//	@Param			after			query		string						false	"Forward pagination argument"
//	@Param			last			query		string						false	"Backward pagination argument"
//	@Param			before			query		string						false	"Backward pagination argument"
//
//	@Success		200				{object}	dto.ObservationConnection	"OK"
//	@Failure		400				{object}	APIResponse					"Error: Bad Request"
//	@Failure		404				{object}	APIResponse					"Error: Not Found"
//	@Failure		401				{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/observations [get]
func (p PresentationHandlersImpl) GetObservations(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "GetObservations")
	span.End()

	queryParams := c.Request.URL.Query()

	var (
		patientID   = queryParams.Get("patient_id")
		concept     = queryParams.Get("concept")
		encounterID = queryParams.Get("encounter_id")
		date        = queryParams.Get("date")
		first       = queryParams.Get("first")
		after       = queryParams.Get("after")
		before      = queryParams.Get("before")
		last        = queryParams.Get("last")

		// This param has special case. The purpose is to enable BE to re-use this functionality and show completed "Empower" examinations. Should be an ENUM incase other integrations need similar adjustment
		useContext = queryParams.Get("use_context")

		// Specifically used to help in search operation for pagination purposes. The client MUST provide this ID so that FHIR server is able to build next and previous page urls
		searchID = queryParams.Get("searchID")
		status   = queryParams.Get("status")
	)

	var encounterInput *string
	if encounterID != "" {
		encounterInput = &encounterID
	}

	var formattedDate *scalarutils.Date

	if date != "" {
		parsedTime, err := time.Parse(time.DateOnly, date)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		year, month, day := parsedTime.Date()

		formattedDate = &scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		}
	}

	var pagination serverutils.PaginationInput

	if first != "" {
		firstVal, err := strconv.Atoi(first)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid 'first' value: %w", err))

			return
		}

		pagination.First = &firstVal

		if after != "" {
			if searchID == "" {
				silgotel.NewLogger(name)
				silgotel.RecordError(span, fmt.Errorf("'searchID' is required when using 'after'"))
				jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("'searchID' is required when using 'after'"))

				return
			}

			pagination.After = &after
		}
	}

	if last != "" {
		lastVal, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid 'last' value: %w", err))

			return
		}

		pagination.Last = &lastVal

		if before != "" {
			if searchID == "" {
				silgotel.NewLogger(name)
				silgotel.RecordError(span, fmt.Errorf("'searchID' is required when using 'before'"))
				jsonErrorResponse(c, http.StatusBadRequest, fmt.Errorf("'searchID' is required when using 'before'"))

				return
			}

			pagination.Before = &before
		}
	}

	var (
		conceptID string
		err       error
	)

	if concept != "" {
		conceptID, _, err = p.ConceptMapper(dto.ObservationConceptEnum(concept))
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}
	}

	payload := &dto.FetchObservationPayload{
		PatientID:       patientID,
		EncounterID:     encounterInput,
		Date:            formattedDate,
		ObservationCode: conceptID,
		Category:        nil,
		PaginationV2:    &pagination,
		Usage:           useContext,
		SearchID:        searchID,
		Status:          dto.ObservationStatusEnum(status),
	}

	observation, err := p.GetPatientObservations(ctx, payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, observation)
}

// ListPatientEncounters handles incoming http requests to list patient's encounters.
//
//	@Summary		read-instances: List patient encounters
//	@Description	List the encounters of the patient given by the ID
//	@Tags			Encounter
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			patient_id	query		string					true	"Patient ID"
//	@Param			limit		query		string					false	"Forward pagination argument"
//	@Param			after		query		string					false	"Forward pagination argument"
//	@Param			last		query		string					false	"Backward pagination argument"
//	@Param			before		query		string					false	"Backward pagination argument"
//	@Success		200			{object}	dto.EncounterConnection	"OK"
//	@Failure		400			{object}	APIResponse				"Error: Bad Request"
//	@Failure		404			{object}	APIResponse				"Error: Not Found"
//	@Failure		401			{object}	APIResponse				"Error: Not Authorized"
//	@Router			/api/v1/encounter [get]
func (p PresentationHandlersImpl) ListPatientEncounters(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListPatientEncounters")
	span.End()

	queryParams := c.Request.URL.Query()

	var (
		patientID = queryParams.Get("patient_id")
		limit     = queryParams.Get("limit")
		before    = queryParams.Get("before")
		after     = queryParams.Get("after")
		last      = queryParams.Get("last")
	)

	var pagination dto.Pagination

	if limit != "" {
		count, err := strconv.Atoi(limit)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	encounters, err := p.usecases.ListPatientEncounters(ctx, patientID, &pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, encounters)
}

// ListPatientAllergies handles incoming http requests to list patient's allergies
//
//	@Summary		read-instances: List patient allergies
//	@Description	List the allergies of the patient given by the ID
//	@Tags			AllergyIntolerance
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			patient_id	query		string					true	"Patient ID"
//	@Param			limit		query		string					false	"Forward pagination argument"
//	@Param			after		query		string					false	"Forward pagination argument"
//	@Param			last		query		string					false	"Backward pagination argument"
//	@Param			before		query		string					false	"Backward pagination argument"
//	@Success		200			{object}	dto.AllergyConnection	"OK"
//	@Failure		400			{object}	APIResponse				"Error: Bad Request"
//	@Failure		500			{object}	APIResponse				"Error: InternalServerError"
//	@Failure		404			{object}	APIResponse				"Error: Not Found"
//	@Failure		401			{object}	APIResponse				"Error: Not Authorized"
//	@Router			/api/v1/allergyintolerance [get]
func (p PresentationHandlersImpl) ListPatientAllergies(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListPatientAllergies")
	defer span.End()

	queryParams := c.Request.URL.Query()

	var (
		patientID = queryParams.Get("patient_id")
		limit     = queryParams.Get("limit")
		before    = queryParams.Get("before")
		after     = queryParams.Get("after")
		last      = queryParams.Get("last")
	)

	var pagination dto.Pagination

	if limit != "" {
		count, err := strconv.Atoi(limit)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	allergies, err := p.usecases.ListPatientAllergies(ctx, patientID, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, allergies)
}

// ListPatientMedia handles incoming http requests to list patient's media
//
//	@Summary		List patient media
//	@Description	List the media of the patient given by the ID
//	@Tags			Media
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			patient_id	query		string				true	"Patient ID"
//	@Param			limit		query		string				false	"Forward pagination argument"
//	@Param			after		query		string				false	"Forward pagination argument"
//	@Param			last		query		string				false	"Backward pagination argument"
//	@Param			before		query		string				false	"Backward pagination argument"
//	@Success		200			{object}	dto.MediaConnection	"OK"
//	@Failure		400			{object}	APIResponse			"Error: Bad Request"
//	@Failure		500			{object}	APIResponse			"Error: InternalServerError"
//	@Failure		404			{object}	APIResponse			"Error: Not Found"
//	@Failure		401			{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/media [get]
func (p PresentationHandlersImpl) ListPatientMedia(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "ListPatientMedia")
	span.End()

	queryParams := c.Request.URL.Query()

	var (
		encounterID      = queryParams.Get("encounter_id")
		serviceRequestID = queryParams.Get("service_request_id")
		limit            = queryParams.Get("limit")
		before           = queryParams.Get("before")
		after            = queryParams.Get("after")
		last             = queryParams.Get("last")
	)

	var pagination dto.Pagination

	if limit != "" {
		count, err := strconv.Atoi(limit)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination = dto.Pagination{
			First: &count,
		}
	}

	switch {
	case last != "":
		lastInt, err := strconv.Atoi(last)
		if err != nil {
			silgotel.NewLogger(name)
			silgotel.RecordError(span, err)
			jsonErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		pagination.Last = &lastInt
	case before != "":
		pagination.Before = before
	case after != "":
		pagination.After = after
	}

	media, err := p.usecases.ListPatientMedia(ctx, encounterID, serviceRequestID, pagination)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, media)
}

// PatchPatientObservations is used to patch a patient's vital signs
//
//	@Summary		Updates a Patient's Vitals signss
//	@Description	PatchPatientObservations is used to patch a patient's vital signs
//	@Tags			Observations
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			PatchObservationInput	body		dto.PatchObservationInput	true	"PatchObservationInput"
//	@Param			observationID			path		string						true	"Observation ID"
//	@Success		200						{object}	dto.Observation				"OK"
//	@Failure		400						{object}	APIResponse					"Error: Bad Request"
//	@Failure		500						{object}	APIResponse					"Error: InternalServerError"
//	@Failure		404						{object}	APIResponse					"Error: Not Found"
//	@Failure		401						{object}	APIResponse					"Error: Not Authorized"
//	@Router			/api/v1/observations/{id} [patch]
func (p PresentationHandlersImpl) PatchPatientObservations(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "PatchPatientObservations")
	defer span.End()

	id := c.Param("id")
	if id == "" {
		err := fmt.Errorf("missing observation id")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	payload := dto.PatchObservationInput{}

	if err := c.BindJSON(&payload); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	observation, err := p.PatchPatientObservation(ctx, id, &payload)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, observation)
}

// Observation is used to retrieve observation details
//
//	@Summary		Used to retrieve observations
//	@Description	Observation is retrieve patient observations
//	@Tags			Observations
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id	path		string			true	"Observation ID"
//	@Success		200	{object}	dto.Observation	"OK"
//	@Failure		400	{object}	APIResponse		"Error: Bad Request"
//	@Failure		500	{object}	APIResponse		"Error: InternalServerError"
//	@Failure		404	{object}	APIResponse		"Error: Not Found"
//	@Failure		401	{object}	APIResponse		"Error: Not Authorized"
//	@Router			/api/v1/observations/{id} [get]
func (p PresentationHandlersImpl) Observation(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "Observation")
	span.End()

	id := c.Param("id")
	if id == "" {
		err := fmt.Errorf("missing observation id")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	observation, err := p.GetObservationByID(ctx, id)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, observation)
}

// CreateQuestionnaireResponse is used to create a questionnaire response
//
//	@Summary		Creates a questionnaire response
//	@Description	CreateQuestionnaireResponses is used to create a questionnaire response
//	@Tags			QuestionnaireResponse
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			QuestionnaireResponseInput	body		dto.QuestionnaireResponse		true	"QuestionnaireResponse"
//	@Param			encounterID					query		string							true	"Encounter ID"
//	@Param			questionnaireID				query		string							true	"Questionnaire ID"
//	@Success		201							{object}	dto.QuestionnaireReviewSummary	"OK"
//	@Failure		400							{object}	APIResponse						"Error: Bad Request"
//	@Failure		404							{object}	APIResponse						"Error: Not Found"
//	@Failure		500							{object}	APIResponse						"Error: InternalServerError"
//	@Failure		401							{object}	APIResponse						"Error: Not Authorized"
//	@Router			/api/v1/questionnaire-response [post]
func (p PresentationHandlersImpl) CreateQuestionnaireResponses(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateQuestionnaireResponses")
	span.End()

	queryParams := c.Request.URL.Query()
	questionnaireID := queryParams.Get("questionnaireID")
	encounterID := queryParams.Get("encounterID")

	if questionnaireID == "" {
		err := fmt.Errorf("missing questionnaire id")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if encounterID == "" {
		err := fmt.Errorf("missing encounter id")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	input := dto.QuestionnaireResponse{}

	if err := c.BindJSON(&input); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	response, err := p.CreateQuestionnaireResponse(ctx, questionnaireID, encounterID, input)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, response)
}

// CreateMedication is used to create a new medication
func (p PresentationHandlersImpl) CreateMedication(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "CreateMedication")
	span.End()

	var paylaod []*dto.MedicationInput

	if err := c.BindJSON(&paylaod); err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	medication, err := p.RecordMedication(ctx, paylaod)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusCreated, medication)
}

// FetchMedicationByID is used to retrieve patient's medical data
//
//	@Summary		read-instance: Get Medication by ID
//	@Description	Fetches medication
//	@Tags			Medication
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical Organizatin ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Security		OAuth2Password
//	@Param			id	path		string					true	"MedicationID"	Example(5f4a944d-abcd-405a-91cc-a88451195bfe)
//	@Success		200	{object}	dto.MedicationOutput	"OK"
//	@Failure		400	{object}	APIResponse				"Error: Bad Request"
//	@Failure		404	{object}	APIResponse				"Error: Not Found"
//	@Failure		401	{object}	APIResponse				"Error: Not Authorized"
//	@Router			/api/v1/medication/{id} [get]
func (p PresentationHandlersImpl) FetchMedicationByID(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "FetchMedicationByID")
	span.End()

	var medicationID string
	if medicationID = c.Param("id"); medicationID == "" {
		err := fmt.Errorf("missing medication ID")

		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)

		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	medication, err := p.usecases.FetchMedicationByID(ctx, medicationID)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, medication)
}

// GetPatientTimeline handles incoming http requests to get a patient's timeline
//
//	@Summary		read-instance: Get patient timeline
//	@Description	Get a patient's timeline as a minimal, chronologically sorted list of timeline resources
//	@Tags			Patient
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Param			id							path	string	true	"Patient ID"
//	@Param			count						query	int		false	"Maximum number of resources in a page (max 1000)"
//	@Param			type						query	string	false	"Comma-delimited FHIR resource types to include"
//	@Param			page_token					query	string	false	"Page token for pagination"
//	@Param			since						query	string	false	"Only include resources updated after this time (RFC3339 format)"
//	@Param			start						query	string	false	"Start date for filtering (YYYY-MM-DD format)"
//	@Param			end							query	string	false	"End date for filtering (YYYY-MM-DD format)"
//	@Security		OAuth2Password
//	@Success		200	{object}	dto.HealthTimeline	"OK"
//	@Failure		400	{object}	APIResponse			"Error: Bad Request"
//	@Failure		404	{object}	APIResponse			"Error: Not Found"
//	@Failure		401	{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/patients/{id}/timeline [get]
func (p PresentationHandlersImpl) GetPatientTimeline(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "GetPatientTimeline")
	span.End()

	patientID := c.Param("id")

	if patientID == "" {
		silgotel.NewLogger(name)

		errMsg := fmt.Errorf("patient id cannot be null")

		silgotel.RecordError(span, errMsg)
		jsonErrorResponse(c, http.StatusBadRequest, errMsg)

		return
	}

	queryParams := c.Request.URL.Query()
	if queryParams.Get("type") == "" {
		silgotel.NewLogger(name)

		errMsg := fmt.Errorf("type is missing in query param")
		silgotel.RecordError(span, errMsg)
		jsonErrorResponse(c, http.StatusBadRequest, errMsg)

		return
	}

	params := timelineParams(queryParams)

	timeline, err := p.usecases.GetPatientTimeline(ctx, patientID, params)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, timeline)
}

func timelineParams(queryParams url.Values) *dto.PatientEverythingFilterParams {
	var params *dto.PatientEverythingFilterParams

	if queryParams.Get("count") != "" || queryParams.Get("type") != "" ||
		queryParams.Get("page_token") != "" || queryParams.Get("since") != "" ||
		queryParams.Get("start") != "" || queryParams.Get("end") != "" ||
		queryParams.Get("_getpages") != "" || queryParams.Get("_getpagesoffset") != "" {
		params = &dto.PatientEverythingFilterParams{}

		if countStr := queryParams.Get("count"); countStr != "" {
			params.Count = countStr
		}

		if typeStr := queryParams.Get("type"); typeStr != "" {
			params.Type = typeStr
		}

		if pageToken := queryParams.Get("page_token"); pageToken != "" {
			params.PageToken = pageToken
		}

		if getPages := queryParams.Get("_getpages"); getPages != "" {
			params.GetPages = getPages
		}

		if getPagesOffset := queryParams.Get("_getpagesoffset"); getPagesOffset != "" {
			params.GetPagesOffset = getPagesOffset
		}

		if since := queryParams.Get("since"); since != "" {
			params.Since = since
		}

		if start := queryParams.Get("start"); start != "" {
			params.Start = start
		}

		if end := queryParams.Get("end"); end != "" {
			params.End = end
		}
	}

	return params
}

// PatientBanner is used to display patient banner data (specifically, most recent 3 Allergies, Conditions, Medications)
//
//	@Summary		read-instance: Get patient banner info
//	@Description	Get a patient's associated data that powers their banner
//	@Tags			Patient
//	@Accept			json
//	@Produce		json
//	@Param			Clinical-Organization-ID	header	string	true	"Clinical-Organization-ID"
//	@Param			Clinical-Facility-ID		header	string	true	"Clinical-Facility-ID"
//	@Param			id							path	string	true	"Patient ID"
//	@Param			count						query	int		false	"Maximum number of resources in a page (max 1000)"
//	@Param			type						query	string	false	"Comma-delimited FHIR resource types to include"
//	@Param			page_token					query	string	false	"Page token for pagination"
//	@Param			since						query	string	false	"Only include resources updated after this time (RFC3339 format)"
//	@Param			start						query	string	false	"Start date for filtering (YYYY-MM-DD format)"
//	@Param			end							query	string	false	"End date for filtering (YYYY-MM-DD format)"
//	@Security		OAuth2Password
//	@Success		200	{object}	dto.HealthTimeline	"OK"
//	@Failure		400	{object}	APIResponse			"Error: Bad Request"
//	@Failure		404	{object}	APIResponse			"Error: Not Found"
//	@Failure		401	{object}	APIResponse			"Error: Not Authorized"
//	@Router			/api/v1/patients/{id}/banner [get]
func (p PresentationHandlersImpl) PatientBanner(c *gin.Context) {
	ctx, span := silgotel.Trace(c.Request.Context(), name, "PatientBanner")
	defer span.End()

	patientID := c.Param("id")
	if patientID == "" {
		silgotel.NewLogger(name)

		errMsg := fmt.Errorf("patient id cannot be null")
		silgotel.RecordError(span, errMsg)
		jsonErrorResponse(c, http.StatusBadRequest, errMsg)

		return
	}

	queryParams := c.Request.URL.Query()

	params := timelineParams(queryParams)

	timeline, err := p.GetPatientBanner(ctx, patientID, params)
	if err != nil {
		silgotel.NewLogger(name)
		silgotel.RecordError(span, err)
		jsonErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, timeline)
}
