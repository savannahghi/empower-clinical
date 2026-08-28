package dto

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

type OrganizationIdentifier struct {
	Type  OrganizationIdentifierType `json:"type,omitempty"`
	Value string                     `json:"value,omitempty"`
}

type OrganizationInput struct {
	ID          string                   `json:"id,omitempty"`
	Name        string                   `json:"name,omitempty"`
	PhoneNumber string                   `json:"phoneNumber,omitempty"`
	Identifiers []OrganizationIdentifier `json:"identifiers,omitempty"`
}

type EpisodeOfCareInput struct {
	Status    EpisodeOfCareStatusEnum `json:"status" validate:"required"`
	PatientID string                  `json:"patientID" validate:"required"`
}

type EncounterInput struct {
	Status EncounterStatusEnum `json:"status" validate:"required"`
}

// StartEncounterInput models the iput for StartEncounter
type StartEncounterInput struct {
	EpisodeOfCareID string `json:"episodeOfCareID" validate:"required"`
}

// ObservationInput models the observation input
type ObservationInput struct {
	Status             ObservationStatus      `json:"status,omitempty" validate:"required"`
	Concept            ObservationConceptEnum `json:"concept,omitempty"`
	EncounterID        string                 `json:"encounterID,omitempty" validate:"required"`
	Note               string                 `json:"note,omitempty"`
	Value              string                 `json:"value,omitempty"`
	Date               *scalarutils.Date      `json:"date,omitempty" swaggertype:"primitive,string"`
	ObservationSubType ObservationSubtype     `json:"observationSubtype,omitempty"`
	UsageContext       ScreeningTypeEnum      `json:"usageContext,omitempty"`
}
type PatientInput struct {
	FirstName   string            `json:"firstName" validate:"required"`
	LastName    string            `json:"lastName" validate:"required"`
	OtherNames  *string           `json:"otherNames"`
	BirthDate   *scalarutils.Date `json:"birthDate,omitempty" swaggertype:"primitive,string"`
	Gender      Gender            `json:"gender"`
	Identifiers []IdentifierInput `json:"identifiers"`
	Contacts    []ContactInput    `json:"contacts"`
}

type IdentifierInput struct {
	Type  IdentifierType `json:"type"`
	Value string         `json:"value"`
}

type ContactInput struct {
	Type  ContactType `json:"type"`
	Value string      `json:"value"`
}

// ConditionInput represents input for creating a FHIR condition
type ConditionInput struct {
	Code        string                   `json:"code" example:"2A21"`
	Name        string                   `json:"name" example:"Mastocytosis"`
	System      domain.TerminologySource `json:"system" example:"LOINC"`
	Status      ConditionStatus          `json:"status" example:"ACTIVE"`
	Category    ConditionCategory        `json:"category" example:"PROBLEM_LIST_ITEM"`
	EncounterID string                   `json:"encounterID" example:"59260d90-80f4-4dc3-ae3a-c7d06c63b1b5"`
	Note        string                   `json:"note" example:"this is test condition"`
	OnsetDate   *scalarutils.Date        `json:"onsetDate" swaggertype:"primitive,string" example:"2025-04-13"`
	Severity    ConditionSeverity        `json:"severity" example:"severe"`
}

type OncologyDiagnosisInput struct {
	EncounterID           string       `json:"encounterId,omitempty" validate:"required"`
	Condition             ValueSetData `json:"condition,omitempty" validate:"required"`
	ICDO3PrimaryTumorCode string       `json:"ICDO3PrimaryTumorCode,omitempty" validate:"required"`
	ICDO3MorphologyCode   string       `json:"ICDO3MorphologyCode,omitempty" validate:"required"`
	Behavior              ValueSetData `json:"behavior,omitempty" validate:"required"`
	Grade                 ValueSetData `json:"grade,omitempty" validate:"required"`
	Stage                 ValueSetData `json:"stage,omitempty" validate:"required"`
	Notes                 string       `json:"notes,omitempty" validate:"required"`
}

// TreatmentEnrollmentInput is a minimal input used to retrospectively capture a patient's known
// diagnosis and their enrollment into treatment, for easier reporting and analytics.
type TreatmentEnrollmentInput struct {
	EncounterID       string            `json:"encounter_id,omitempty" validate:"required"`
	Condition         ValueSetData      `json:"condition,omitempty" validate:"required"`
	CancerStage       ValueSetData      `json:"cancer_stage,omitempty"`
	Date              *scalarutils.Date `json:"date,omitempty" validate:"required"`
	LinkedToTreatment bool              `json:"linked_to_treatment"`
	TreatmentFacility string            `json:"treatment_facility,omitempty"`
	TreatmentProgram  string            `json:"treatment_program,omitempty"`
	EnrollmentDate    *scalarutils.Date `json:"enrollment_date,omitempty"`
	Severity          ConditionSeverity `json:"severity" validate:"required"`
}

// UpdateTreatmentEnrollmentInput models the mutable fields of a treatment enrollment (recorded as a
// FHIR Condition via RecordTreatmentEnrollment). All fields are optional; only the fields that are
// supplied are updated and the rest are left unchanged.
type UpdateTreatmentEnrollmentInput struct {
	Condition      *ValueSetData     `json:"condition,omitempty"`
	Date           *scalarutils.Date `json:"date,omitempty"`
	EnrollmentDate *scalarutils.Date `json:"enrollment_date,omitempty"`
}

// AllergyInput models the allergy input
type AllergyInput struct {
	PatientID         string                   `json:"patientID"`
	Code              string                   `json:"code" validate:"required"`
	TerminologySource domain.TerminologySource `json:"terminologySource" validate:"required" swaggertype:"primitive,string"`
	EncounterID       string                   `json:"encounterID" validate:"required,uuid4"`
	Reaction          *ReactionInput           `json:"reaction"`
}

// ReactionInput models the reaction input
type ReactionInput struct {
	Code     string                                 `json:"code"`
	System   string                                 `json:"system"`
	Severity AllergyIntoleranceReactionSeverityEnum `json:"severity" swaggertype:"primitive,string"`
}

// MediaInput models the dataclass to upload media to FHIR
type MediaInput struct {
	EncounterID      string                             `json:"encounterID"`
	ServiceRequestID string                             `json:"serviceRequestID"`
	File             map[string][]*multipart.FileHeader `form:"file" json:"file"`
}

// CompositionInput models the composition input
type CompositionInput struct {
	EncounterID string                `json:"encounterID"`
	Type        CompositionType       `json:"type"`
	Category    CompositionCategory   `json:"category"`
	Status      CompositionStatusEnum `json:"status"`
	Note        string                `json:"note"`
}

// PatchCompositionInput models the patch composition input
type PatchCompositionInput struct {
	Type     CompositionType       `json:"type"`
	Category CompositionCategory   `json:"category"`
	Status   CompositionStatusEnum `json:"status"`
	Note     string                `json:"note"`
	Section  []*SectionInput       `json:"section"`
}

// SectionInput models the composition section input
type SectionInput struct {
	ID      string          `json:"id,omitempty"`
	Title   string          `json:"title"`
	Code    string          `json:"code"`
	Author  string          `json:"author"`
	Text    string          `json:"text"`
	Section []*SectionInput `json:"section"`
}

// ConsentInput models the consent input
type ConsentInput struct {
	Decision      domain.ConsentDecisionEnum `json:"decision"`
	EncounterID   string                     `json:"encounterID,omitempty"`
	DenyReason    string                     `json:"denyReason,omitempty"`
	ScreeningType ScreeningTypeEnum          `json:"screeningType,omitempty"`
}

// QuestionnaireResponse models input for questionnaire response resource in fhir
type QuestionnaireResponse struct {
	ResourceType string                                 `json:"resourceType,omitempty"`
	Meta         MetaInput                              `json:"meta,omitempty"`
	Status       domain.QuestionnaireResponseStatusEnum `json:"status"`
	Authored     string                                 `json:"authored,omitempty"`
	Item         []QuestionnaireResponseItem            `json:"item,omitempty"`
}

// MetaInput represents the data class model of a metadata input
type MetaInput struct {
	VersionID   string            `json:"versionId,omitempty"`
	LastUpdated time.Time         `json:"lastUpdated,omitempty"`
	Source      string            `json:"source,omitempty"`
	Tag         []Coding          `json:"tag,omitempty"`
	Security    []Coding          `json:"security,omitempty"`
	Profile     []scalarutils.URI `json:"profile,omitempty" swaggertype:"primitive,string"`
}

// QuestionnaireResponseItem models input for item object of questionnaire response resource
type QuestionnaireResponseItem struct {
	LinkID string                            `json:"linkId"`
	Text   *string                           `json:"text,omitempty"`
	Answer []QuestionnaireResponseItemAnswer `json:"answer,omitempty"`
	Item   []QuestionnaireResponseItem       `json:"item,omitempty"`
}

// FHIRQuestionnaireResponseItemAnswer models item answer object of questionnaire response resource
type QuestionnaireResponseItemAnswer struct {
	ValueBoolean    *bool                       `json:"valueBoolean,omitempty"`
	ValueDecimal    *float64                    `json:"valueDecimal,omitempty"`
	ValueInteger    *int                        `json:"valueInteger,omitempty"`
	ValueDate       *string                     `json:"valueDate,omitempty"`
	ValueDateTime   *string                     `json:"valueDateTime,omitempty"`
	ValueTime       *string                     `json:"valueTime,omitempty"`
	ValueString     *string                     `json:"valueString,omitempty"`
	ValueURI        *string                     `json:"valueUri,omitempty"`
	ValueAttachment *Attachment                 `json:"valueAttachment,omitempty"`
	ValueCoding     *Coding                     `json:"valueCoding,omitempty"`
	ValueQuantity   *Quantity                   `json:"valueQuantity,omitempty"`
	ValueReference  *Reference                  `json:"valueReference,omitempty"`
	Item            []QuestionnaireResponseItem `json:"item,omitempty"`
}

// ToString converts the existing item answer to a string
func (q *QuestionnaireResponseItemAnswer) ToString() string {
	switch {
	case q.ValueBoolean != nil:
		return strconv.FormatBool(*q.ValueBoolean)
	case q.ValueDecimal != nil:
		return strconv.FormatFloat(*q.ValueDecimal, 'f', -1, 64)
	case q.ValueInteger != nil:
		return strconv.Itoa(*q.ValueInteger)
	case q.ValueDate != nil:
		return *q.ValueDate
	case q.ValueDateTime != nil:
		return *q.ValueDateTime
	case q.ValueTime != nil:
		return *q.ValueTime
	case q.ValueString != nil:
		return *q.ValueString
	case q.ValueURI != nil:
		return *q.ValueURI
	case q.ValueCoding != nil:
		return q.ValueCoding.ToString()

	default:
		return ""
	}
}

// Coding : an input for a code defined by a terminology system.
type Coding struct {
	ID           string            `json:"id,omitempty"`
	System       scalarutils.URI   `json:"system,omitempty" swaggertype:"primitive,string"`
	Version      string            `json:"version,omitempty"`
	Code         *scalarutils.Code `json:"code,omitempty" swaggertype:"primitive,string"`
	Display      string            `json:"display,omitempty"`
	UserSelected bool              `json:"userSelected,omitempty"`
}

// ToString returns the Display field of the Coding struct as a string
func (c *Coding) ToString() string {
	return c.Display
}

// Attachment definition: input for referring to data content defined in other formats.
type Attachment struct {
	ID          string                   `json:"id,omitempty"`
	ContentType scalarutils.Code         `json:"contentType,omitempty" swaggertype:"primitive,string"`
	Language    scalarutils.Code         `json:"language,omitempty" swaggertype:"primitive,string"`
	Data        scalarutils.Base64Binary `json:"data,omitempty" swaggertype:"primitive,string"`
	URL         scalarutils.URL          `json:"url,omitempty" swaggertype:"primitive,string"`
	Size        int                      `json:"size,omitempty"`
	Hash        scalarutils.Base64Binary `json:"hash,omitempty" swaggertype:"primitive,string"`
	Title       string                   `json:"title,omitempty"`
	Creation    scalarutils.DateTime     `json:"creation,omitempty" swaggertype:"primitive,string"`
}

// Quantity definition: input for measured amount (or an amount that can potentially be measured). note that measured amounts include amounts that are not precisely quantified, including amounts involving arbitrary units and floating currencies.
type Quantity struct {
	ID         string                  `json:"id,omitempty"`
	Value      float64                 `json:"value"`
	Comparator *QuantityComparatorEnum `json:"comparator,omitempty"`
	Unit       string                  `json:"unit"`
	System     scalarutils.URI         `json:"system" swaggertype:"primitive,string"`
	Code       scalarutils.Code        `json:"code" swaggertype:"primitive,string"`
}

// Reference definition: input for reference from one resource to another.
type Reference struct {
	ID         string          `json:"id,omitempty"`
	Reference  string          `json:"reference,omitempty"`
	Type       scalarutils.URI `json:"type,omitempty" swaggertype:"primitive,string"`
	Identifier *Identifier     `json:"identifier,omitempty"`
	Display    string          `json:"display,omitempty"`
}

// ResourceID safely extracts the logical ID from the Reference.reference field (e.g., from "ResourceType/12345").
//
// Note:
// Since we dont want the codebase to be full of such code,
//
// resourceID, err := activity.Reference.GetResourceID()
//
//	if err != nil {
//		return nil, err
//	}
//
// We shall leave the responsibility of handling nil references to the function caller.
// e.g (something like below).
//
//	if r == nil || r.Reference == nil || *r.Reference == "" {
//		return "", fmt.Errorf("cannot get resource ID from a nil or empty reference")
//	}.
func (r *Reference) ResourceID() string {
	if r == nil || r.Reference == "" {
		return ""
	}

	refStr := r.Reference

	lastSlashIndex := strings.LastIndex(refStr, "/")
	if lastSlashIndex == -1 || lastSlashIndex == len(refStr)-1 {
		return ""
	}

	return refStr[lastSlashIndex+1:]
}

// Expression is documented here http://hl7.org/fhir/StructureDefinition/Expression
type Expression struct {
	ID          *string     `json:"id,omitempty"`
	Extension   []Extension `json:"extension,omitempty"`
	Description *string     `json:"description,omitempty"`
	Name        *string     `json:"name,omitempty"`
	Language    string      `json:"language,omitempty"`
	Expression  *string     `json:"expression,omitempty"`
	Reference   *string     `json:"reference,omitempty"`
}

// DiagnosticReportInput represents the data class used to provide diagnostic report information
type DiagnosticReportInput struct {
	EncounterID  string            `json:"encounterID,omitempty" validate:"required"`
	Date         *scalarutils.Date `json:"date" validate:"required"`
	Note         string            `json:"note,omitempty"`
	Media        []*Media          `json:"media"`
	Findings     string            `json:"findings,omitempty"`
	UsageContext ScreeningTypeEnum `json:"usageContext,omitempty"`
}

// TestResultInput is specifically used to provide test result information for a ordered set of test(s)
type TestResultInput struct {
	Entry            DiagnosticReportInput `json:"entry,omitempty"`
	ServiceRequestID string                `json:"serviceRequestID,omitempty"`
}

// ObservationPayload is a custom data type for creating an observation thats customized for specific use-cases
type ObservationPayload struct {
	ObservationInput    ObservationInput `json:"observationInput"`
	VitalSignsConceptID string           `json:"vitalSignsConceptID,omitempty"`
	ServiceRequestID    string           `json:"serviceRequestID,omitempty"`
}

type PSAInput struct {
	DiagnosticInput *DiagnosticReportInput `json:"diagnosticInput,omitempty"`
	PSAType         PSAType                `json:"psaType"`
}

// PatientEverythingFilterParams provides filter parameters that can be combined to filter compartments when retrieving patient information
type PatientEverythingFilterParams struct {
	Count          string `json:"count,omitempty"`
	PageToken      string `json:"pageToken,omitempty"`
	GetPages       string `json:"_getpages,omitempty"`
	GetPagesOffset string `json:"_getpagesoffset,omitempty"`
	Since          string `json:"since,omitempty"`
	Type           string `json:"type,omitempty"`
	End            string `json:"end,omitempty"`
	Start          string `json:"start,omitempty"`
	Fields         string `json:"fields,omitempty"`
}

// ReferralInput represents the input for referring a patient
type ReferralInput struct {
	EncounterID  string            `json:"encounterID" validate:"required"`
	ReferralType ReferralTypeEnum  `json:"referralType"`
	Tests        []string          `json:"tests,omitempty"`
	Specialist   string            `json:"specialist,omitempty"`
	Facility     *FacilityInput    `json:"facility"`
	ReferralNote string            `json:"notes"`
	UsageContext ScreeningTypeEnum `json:"usageContext,omitempty"`
	// ReferralDate is the date the referral is made. It is optional: when omitted the current
	// date is used, and when supplied it may not be later than today.
	ReferralDate *scalarutils.Date `json:"referralDate,omitempty"`
}

// FacilityInput represents the facility payload passed when referring a patient
type FacilityInput struct {
	Name    string `json:"name,omitempty"`
	County  string `json:"county,omitempty"`
	Contact string `json:"contact,omitempty"`
	Email   string `json:"email,omitempty"`

	Phone              string                  `json:"phone"`
	Active             bool                    `json:"active"`
	Country            Country                 `json:"country"`
	Address            string                  `json:"address"`
	Description        string                  `json:"description"`
	FHIROrganisationID string                  `json:"fhirOrganisationID"`
	Identifier         FacilityIdentifierInput `json:"identifier"`
	Coordinates        CoordinatesInput        `json:"coordinates"`
	Services           []FacilityServiceInput  `json:"services"`
	BusinessHours      []BusinessHoursInput    `json:"businessHours"`
}

type ScheduleAppointmentInput struct {
	EncounterID string            `json:"encounterID,omitempty" validate:"required"`
	PatientID   string            `json:"patientID,omitempty" validate:"required"`
	Reason      string            `json:"reason,omitempty"`
	Date        *scalarutils.Date `json:"date,omitempty" swaggertype:"primitive,string"`
}

type ScheduleAppointmentPayload struct {
	AppointmentInput ScheduleAppointmentInput `json:"appointmentInput,omitempty"`
	HeadersInput     AdvantageHeaders         `json:"headersInput,omitempty" validate:"required"`
}

// FilterInput contains fields that are commonly used for searching/filtering resource data
type FilterInput struct {
	PatientID   string            `json:"patientID,omitempty"`
	EncounterID string            `json:"encounterID,omitempty"`
	Date        *scalarutils.Date `json:"date,omitempty"`
}

// TaskFilterInput is used to specifically filter tasks
type TaskFilterInput struct {
	FilterInput
	Type   string     `json:"type"`
	Status TaskStatus `json:"status"`

	// PatientSearch is a free-text term matched across the whole Patient resource (name, any
	// identifier, telecom, etc.) via the FHIR `_content` full-text index — the general search box.
	PatientSearch string `json:"patient"`
}

func (f *TaskFilterInput) SetStatus() string {
	status := f.Status

	if status == "" {
		f.Status = RequestedTasksStatus
	}

	return status.String()
}

// MedicationRequestFilterInput is used to filter medication requests
type MedicationRequestFilterInput struct {
	FilterInput
	Status domain.MedicationRequestStatus
}

type PatchMedicationInput struct {
	Status domain.MedicationRequestStatus `json:"status,omitempty"`
}

// MedicationRequestFilters validates the whether the filter parameters used to search for medication requests are valid
// for filters and set them as search parameters
func (mr *MedicationRequestFilterInput) MedicationRequestFilters(params map[string]interface{}) error {
	err := mr.ValidateBaseFilters(params)
	if err != nil {
		return err
	}

	if mr.Status != "" {
		params["status"] = mr.Status.String()
	}

	return nil
}

// TaskFilters validates the whether the filter parameters used to search for tasks are valid
// for filters and set them as search parameters
func (f *TaskFilterInput) TaskFilters(params map[string]interface{}) error {
	err := f.ValidateBaseFilters(params)
	if err != nil {
		return err
	}

	if f.Date != nil {
		params["authored-on"] = f.Date.AsTime().Format(time.DateOnly)
	}

	if f.Type != "" {
		params["business-status:text"] = f.Type
	}

	if f.Status != "" {
		params["status"] = f.Status.String()
	}

	return nil
}

// ValidateBaseFilters is used to check the validity of commonly shared filters for almost every resource
func (f *FilterInput) ValidateBaseFilters(params map[string]any) error {
	if f.PatientID != "" {
		_, err := uuid.Parse(f.PatientID)
		if err != nil {
			return fmt.Errorf("invalid patient id: %s", f.PatientID)
		}

		patientReference := fmt.Sprintf("Patient/%s", f.PatientID)
		params["patient"] = patientReference
	}

	if f.EncounterID != "" {
		_, err := uuid.Parse(f.EncounterID)
		if err != nil {
			return fmt.Errorf("invalid encounter id: %s", f.EncounterID)
		}

		encounterReference := fmt.Sprintf("Encounter/%s", f.EncounterID)
		params["encounter"] = encounterReference
	}

	return nil
}

type RiskAssessmentFilterInput struct {
	FilterInput
	ScreeningType ScreeningTypeEnum
	Result        string
}

// RiskAssessmentFilter returns valid filter parameters for  risk assessment
func (f *RiskAssessmentFilterInput) RiskAssessmentFilter(params map[string]interface{}) error {
	err := f.ValidateBaseFilters(params)
	if err != nil {
		return err
	}

	if f.Date != nil {
		params["date"] = f.Date.AsTime().Format(time.DateOnly)
	}

	if f.ScreeningType != "" {
		params["_text"] = f.ScreeningType.Text()
	}

	if f.Result != "" {
		params["risk:text"] = f.Result
	}

	return nil
}

// PatchTaskInput used to update a given task
type PatchTaskInput struct {
	Status       TaskStatus
	DueDate      scalarutils.DateTime `json:"dueDate" swaggertype:"primitive,string"`
	Author       string
	Notes        string
	UpdateReason string
}

// ServiceRequestFilterInput is used to filter service requests
type ServiceRequestFilterInput struct {
	FilterInput
	FacilityID string
	Type       string
	Status     ServiceRequestStatus

	// PatientSearch is a free-text term matched across the whole Patient resource (name, any
	// identifier, telecom, etc.) via the FHIR `_content` full-text index — the general search box.
	PatientSearch string
}

// ServiceRequestFilters validates the whether the filter parameters used to search for service requests are valid
// filters and set them as search parameters
func (f *ServiceRequestFilterInput) ServiceRequestFilters(params map[string]any) error {
	err := f.ValidateBaseFilters(params)
	if err != nil {
		return err
	}

	if f.Date != nil {
		params["authored"] = f.Date.AsTime().Format(time.DateOnly)
	}

	// Type narrows results to a specific order/test code (e.g. the lab order LOINC code). It maps to
	// the FHIR ServiceRequest `code` search parameter and is orthogonal to the category filter.
	if f.Type != "" {
		params["code"] = f.Type
	}

	status := f.Status

	if status == "" {
		status = ServiceRequestStatusActive
	}

	params["status"] = string(status)

	if f.FacilityID != "" {
		params["performer"] = f.FacilityID
	}

	return nil
}

// FacilityIdentifierInput is the identifier of the facility
type FacilityIdentifierInput struct {
	Type       FacilityIdentifierType `json:"type"`
	Value      string                 `json:"value"`
	FacilityID string                 `json:"facilityID"`
}

// CoordinatesInput is used to get the coordinates of a facility
type CoordinatesInput struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// FacilityServiceInput is used to get the services offered in a facility
type FacilityServiceInput struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Identifiers []ServiceIdentifierInput `json:"identifiers"`
}

// ServiceIdentifierInput is used to hold the identifier values of the service that is offered in a facility
type ServiceIdentifierInput struct {
	IdentifierType  Terminologies `json:"identifierType"`
	IdentifierValue string        `json:"identifierValue"`
}

// BusinessHoursInput is used to model business hours data input
type BusinessHoursInput struct {
	Day         DayOfWeek `json:"day"`
	OpeningTime string    `json:"openingTime"`
	ClosingTime string    `json:"closingTime"`
}

type PrescriptionInput struct {
	EncounterID string                        `json:"encounterID,omitempty" validate:"required"`
	Medications []PrescriptionMedicationInput `json:"medications" validate:"required"`
}

type PrescriptionMedicationInput struct {
	MedicationID       string                    `json:"medicationID,omitempty"`
	DosageInstructions []DosageInstruction       `json:"dosageInstructions,omitempty"`
	Priority           MedicationRequestPriority `json:"priority,omitempty"`
}

// CreateReferralInput models the input to the API used to create a referral.
type CreateReferralInput struct {
	EncounterID string `json:"encounterID,omitempty"`
	// This field will be used to hold the diagnosis, since it will be created as a
	// condition in FHIR
	Diagnosis       string            `json:"conditionID,omitempty"`
	ReferralType    string            `json:"referralType,omitempty"`
	Urgency         string            `json:"urgency,omitempty"`
	ClinicalHistory string            `json:"clinicalHistory,omitempty"`
	ReferralDate    scalarutils.Date  `json:"referralDate,omitempty" swaggertype:"primitive,string"`
	ReferredFrom    *ReferralFacility `json:"referredFrom,omitempty"`
	ReferredTo      ReferralFacility  `json:"referredTo,omitempty"`

	// Empower related fields
	Tests        []string          `json:"tests,omitempty"`
	Specialist   string            `json:"specialist,omitempty"`
	UsageContext ScreeningTypeEnum `json:"usageContext,omitempty"`
}

// ReferralFacility holds the details of the facilities in the referral workflow
// Accounts for both the sending and the receiving facility
type ReferralFacility struct {
	FHIROrganisationID string `json:"fhirOrganisationID,omitempty"`
	Name               string `json:"name,omitempty"`
	PhoneNumber        string `json:"phoneNumber,omitempty"`
	Email              string `json:"email,omitempty"`
	Branch             string `json:"branch,omitempty"`
}

// TestOrder represents a diagnostic test order
type TestOrder struct {
	EncounterID string `json:"encounterID,omitempty"`
	// This field will be used to hold the diagnosis, since it will be created as a
	// condition in FHIR
	Diagnosis    string           `json:"conditionID,omitempty"`
	Name         string           `json:"name,omitempty"`
	LoincCode    string           `json:"loincCode,omitempty"`
	Status       string           `json:"status,omitempty"`
	Facility     ReferralFacility `json:"facility,omitempty"`
	ClinicalNote string           `json:"clinicalNote,omitempty"`
}

// TestOrderResult models the fields needed when recording a test order result
// They will eventually be saved in FHIR Diagnostic Report & Observation resources.
type TestOrderResult struct {
	ServiceRequestID string                 `json:"serviceRequestID,omitempty"`
	Test             []TestOrderObservation `json:"test,omitempty"`
}

// TestOrderObservation represents the tests that have been done
type TestOrderObservation struct {
	Test    string `json:"test,omitempty"`
	Value   string `json:"value,omitempty"`
	Finding string `json:"finding,omitempty"`
}

// IntraLabOrderInput model is used lto model a minimal payload that is sufficient enough to create a test order on intra-referrals
type IntraLabOrderInput struct {
	Code          string            `json:"code,omitempty"`
	Name          string            `json:"name,omitempty"`
	EncounterID   string            `json:"encounterID,omitempty"`
	Patient       Patient           `json:"patientID,omitempty"`
	ObservationID string            `json:"observationID,omitempty"`
	UsageContext  ScreeningTypeEnum `json:"usageContext,omitempty"`
	// Date indicates when the test was performed. When nil, callers fall back to the current time.
	Date *scalarutils.Date `json:"date,omitempty"`
}

// ShareReferralFormInput models the data class used to initiate sharing a referral form to a patient
type ShareReferralFormInput struct {
	ServiceRequestID string `json:"serviceRequestID,omitempty"`
	WorkstationID    string `json:"workstationID,omitempty"`
	BranchID         string `json:"branchID,omitempty"`
}
type TestInput struct {
	TestType LabTestTypeEnum       `json:"testType,omitempty" validate:"required"`
	Input    DiagnosticReportInput `json:"input,omitempty" validate:"required"`
}

type PatchObservationInput struct {
	Value           string                `json:"value,omitempty" validate:"required"`
	ObservationType PatientVitalSignsEnum `json:"observationType,omitempty"`
}

type MedicationInput struct {
	Name     string       `json:"name,omitempty"`
	DoseForm ValueSetData `json:"doseForm,omitempty"`
}

type FetchObservationPayload struct {
	PatientID       string
	EncounterID     *string
	Date            *scalarutils.Date
	ObservationCode string
	Category        *ObservationCategory
	Pagination      *Pagination
	Status          ObservationStatusEnum `json:"status"`

	// PaginationV2 will eventually transition to `Pagination` once all list pagination work is done
	PaginationV2 *serverutils.PaginationInput
	Usage        string
	SearchID     string
}

func (f *FetchObservationPayload) SetStatus() string {
	status := f.Status

	if status == "" {
		status = ObservationStatusEnumFinal
	}

	return status.String()
}

type ReferralSearchInput struct {
	PatientID   *string                         `json:"patientID"`
	EncounterID *string                         `json:"encounterID"`
	Date        *scalarutils.Date               `json:"date,omitempty" swaggertype:"primitive,string"`
	Pagination  *Pagination                     `json:"pagination"`
	Status      domain.ServiceRequestStatusEnum `json:"status"`

	// PatientSearch is a free-text term matched across the whole Patient resource (name, any
	// identifier, telecom, etc.) via the FHIR `_content` full-text index — the general search box.
	PatientSearch string `json:"patient"`
}

func (r *ReferralSearchInput) SetStatus() string {
	status := r.Status

	if status == "" {
		status = domain.ServiceRequestStatusActive
	}

	return string(status)
}
