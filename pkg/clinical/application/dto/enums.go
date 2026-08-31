package dto

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
)

// Organization
type OrganizationIdentifierType string

const (
	SladeCode OrganizationIdentifierType = "slade-code"
	MFLCode   OrganizationIdentifierType = "mfl-code"
	FRID      OrganizationIdentifierType = "facility-registry-id"
	SLADEID   OrganizationIdentifierType = "SLADEID"
	FHIRID    OrganizationIdentifierType = "fhir-id"
	FIDCode   OrganizationIdentifierType = "fid-code"
)

func (o OrganizationIdentifierType) ToString() string {
	switch o {
	case SladeCode:
		return "Slade360 Code"
	case MFLCode:
		return "Master Facility List Code"
	case FRID:
		return "Facility Registry ID"
	case SLADEID:
		return "Universal System Identifier"
	case FHIRID:
		return "FHIR Id"
	case FIDCode:
		return "Facility ID Code"
	default:
		return fmt.Sprintf("Unknown OrganizationIdentifierType: %s", o)
	}
}

type OrganisationSource string

const (
	OrganisationSourceWHO  OrganisationSource = "WHO"
	OrganisationSourceCIEL OrganisationSource = "CIEL"
)

func (o OrganizationIdentifierType) IsValid() bool {
	switch o {
	case SladeCode, SLADEID, FRID, MFLCode, FHIRID:
		return true
	}

	return false
}

type EpisodeOfCareStatusEnum string

const (
	EpisodeOfCareStatusEnumPlanned        EpisodeOfCareStatusEnum = "PLANNED"
	EpisodeOfCareStatusEnumActive         EpisodeOfCareStatusEnum = "ACTIVE"
	EpisodeOfCareStatusEnumFinished       EpisodeOfCareStatusEnum = "FINISHED"
	EpisodeOfCareStatusEnumCancelled      EpisodeOfCareStatusEnum = "CANCELLED"
	EpisodeOfCareStatusEnumEnteredInError EpisodeOfCareStatusEnum = "ENTERED_IN_ERROR"
)

// EncounterStatusEnum is a FHIR enum
type EncounterStatusEnum string

const (
	// EncounterStatusEnumPlanned ...
	EncounterStatusEnumPlanned EncounterStatusEnum = "PLANNED"
	// EncounterStatusEnumOnHold ...
	EncounterStatusEnumOnHold EncounterStatusEnum = "ON_HOLD"
	// EncounterStatusEnumInProgress ...
	EncounterStatusEnumInProgress EncounterStatusEnum = "IN_PROGRESS"
	// EncounterStatusEnumCompleted ...
	EncounterStatusEnumCompleted EncounterStatusEnum = "COMPLETED"
	// EncounterStatusEnumCancelled ...
	EncounterStatusEnumCancelled EncounterStatusEnum = "CANCELLED"
	// EncounterStatusEnumDiscontinued...
	EncounterStatusEnumDiscontinued EncounterStatusEnum = "DISCONTINUED"
)

// EncounterClassEnum ...
type EncounterClassEnum string

const (
	// Also referred to as outpatient - For now we'll start with outpatient only
	EncounterClassEnumAmbulatory EncounterClassEnum = "AMBULATORY"
)

type ResourceType string

const (
	ResourceTypeAllergyIntolerance  ResourceType = "AllergyIntolerance"
	ResourceTypeObservation         ResourceType = "Observation"
	ResourceTypeCondition           ResourceType = "Condition"
	ResourceTypeMedicationStatement ResourceType = "MedicationStatement"
	ResourceTypeRiskAssessment      ResourceType = "RiskAssessment"
	ResourceTypeTask                ResourceType = "Task"
	ResourceTypeServiceRequest      ResourceType = "ServiceRequest"
)

type AllergyIntoleranceReactionSeverityEnum string

const (
	AllergyIntoleranceReactionSeverityEnumMild     AllergyIntoleranceReactionSeverityEnum = "MILD"
	AllergyIntoleranceReactionSeverityEnumModerate AllergyIntoleranceReactionSeverityEnum = "MODERATE"
	AllergyIntoleranceReactionSeverityEnumSevere   AllergyIntoleranceReactionSeverityEnum = "SEVERE"
)

type ObservationStatus string

const (
	ObservationStatusFinal     ObservationStatus = "FINAL"
	ObservationStatusCancelled ObservationStatus = "CANCELLED"
)

type MedicationStatementStatusEnum string

const (
	MedicationStatementStatusEnumActive     MedicationStatementStatusEnum = "ACTIVE"
	MedicationStatementStatusEnumInActive   MedicationStatementStatusEnum = "INACTIVE"
	MedicationStatementStatusEnumUnknown    MedicationStatementStatusEnum = "UNKNOWN"
	MedicationStatementStatusEnumRecurrence MedicationStatementStatusEnum = "RECURRENCE"
	MedicationStatementStatusEnumRelapse    MedicationStatementStatusEnum = "RELAPSE"
	MedicationStatementStatusEnumRemission  MedicationStatementStatusEnum = "REMISSSION"
)

type IdentifierType string

const (
	IdentifierTypeNationalID IdentifierType = "national_id"
	IdentifierTypePassport   IdentifierType = "passport_number"
	IdentifierTypeAlienID    IdentifierType = "alien_id"
	IdentifierTypeCCCNumber  IdentifierType = "ccc_number"
	IdentifierTypeHealthID   IdentifierType = "slade_health_id"
)

type ContactType string

const (
	ContactTypePhoneNumber ContactType = "PHONE_NUMBER"
)

// Gender is a FHIR enum
type Gender string

const (
	// GenderMale ...
	GenderMale Gender = "male"
	// GenderFemale ...
	GenderFemale Gender = "female"
	// GenderOther ...
	GenderOther Gender = "other"
)

// ConditionStatus represents status of a FHIR condition
type ConditionStatus string

const (
	ConditionStatusActive   ConditionStatus = "ACTIVE"
	ConditionStatusInactive ConditionStatus = "INACTIVE"
	ConditionStatusResolved ConditionStatus = "RESOLVED"
	ConditionStatusUnknown  ConditionStatus = "UNKNOWN"
)

// ConditionSeverity represents hoe severe a condition is
type ConditionSeverity string

const (
	ConditionSeveritySevere   ConditionSeverity = "severe"
	ConditionSeverityMild     ConditionSeverity = "mild"
	ConditionSeverityModerate ConditionSeverity = "moderate"
)

func (cs ConditionSeverity) ToString() string {
	switch cs {
	case ConditionSeveritySevere:
		return "Severe"
	case ConditionSeverityMild:
		return "Mild"
	case ConditionSeverityModerate:
		return "Moderate"
	default:
		return fmt.Sprintf("Unknown Condition Severity: %s", cs)
	}
}

// ConditionCategory represents status of a FHIR condition
type ConditionCategory string

const (
	ConditionCategoryProblemList ConditionCategory = "PROBLEM_LIST_ITEM"
	ConditionCategoryDiagnosis   ConditionCategory = "ENCOUNTER_DIAGNOSIS"
)

func (cs ConditionCategory) ToCategoryCode() string {
	switch cs {
	case ConditionCategoryProblemList:
		return "problem-list-item"
	case ConditionCategoryDiagnosis:
		return "encounter-diagnosis"
	default:
		return fmt.Sprintf("Unknown Condition Category: %s", cs)
	}
}

func (cs ConditionCategory) ToString() string {
	switch cs {
	case ConditionCategoryProblemList:
		return "Problem List Item"
	case ConditionCategoryDiagnosis:
		return "Encounter Diagnosis"
	default:
		return fmt.Sprintf("Unknown Condition Category: %s", cs)
	}
}

// LOINCCodes represents LOINC assessment codes
type LOINCCodes string

const (
	LOINCPlanOfCareCode     LOINCCodes = "18776-5"
	LOINCAssessmentPlanCode LOINCCodes = "51847-2"
)

// CompositionCategory enum represents category composition attribute
type CompositionCategory string

const (
	AssessmentAndPlan               CompositionCategory = "ASSESSMENT_PLAN"
	HistoryOfPresentingIllness      CompositionCategory = "HISTORY_OF_PRESENTING_ILLNESS"
	SocialHistory                   CompositionCategory = "SOCIAL_HISTORY"
	FamilyHistory                   CompositionCategory = "FAMILY_HISTORY"
	Examination                     CompositionCategory = "EXAMINATION"
	PlanOfCare                      CompositionCategory = "PLAN_OF_CARE"
	ProviderUnspecifiedProgressNote CompositionCategory = "PROGRESS_NOTE"
	PastMedicalSurgeryHistory       CompositionCategory = "PAST_MEDICAL_SURGERY_HISTORY"
	ChiefComplaint                  CompositionCategory = "CHIEF_COMPLAINT"
)

// Type enum represents type composition attribute
type CompositionType string

const (
	ProgressNote CompositionType = "PROGRESS_NOTE"
)

// CompositionStatus enum represents status composition attribute
type CompositionStatusEnum string

const (
	CompositionStatusEnumPreliminary               CompositionStatusEnum = "PRELIMINARY"
	CompositionStatusEnumFinal                     CompositionStatusEnum = "FINAL"
	CompositionStatusEnumAmended                   CompositionStatusEnum = "AMENDED"
	CompositionStatusEnumEnteredInErrorPreliminary CompositionStatusEnum = "ENTERED_IN_ERROR"
)

// ToCode is used to convert to FHIR specific codes
func (t CompositionStatusEnum) ToCode() string {
	switch t {
	case CompositionStatusEnumPreliminary:
		return "preliminary"
	case CompositionStatusEnumFinal:
		return "final"
	case CompositionStatusEnumAmended:
		return "amended"
	case CompositionStatusEnumEnteredInErrorPreliminary:
		return "entered-in-error"
	default:
		return fmt.Sprintf("Unknown CompositionStatusEnum: %s", t)
	}
}

// QuantityComparatorEnum is a FHIR enum
type QuantityComparatorEnum string

const (
	QuantityComparatorEnumLessThan             QuantityComparatorEnum = "less_than"
	QuantityComparatorEnumLessThanOrEqualTo    QuantityComparatorEnum = "less_than_or_equal_to"
	QuantityComparatorEnumGreaterThanOrEqualTo QuantityComparatorEnum = "greater_than_or_equal_to"
	QuantityComparatorEnumGreaterThan          QuantityComparatorEnum = "greater_than"
)

// VIAOutcomeEnum a type enum that represents the results of a VIA test
// VIA (Visual Inspection with Acetic Acid)
type VIAOutcomeEnum string

const (
	VIAOutcomeNegative               VIAOutcomeEnum = "negative"
	VIAOutcomePositive               VIAOutcomeEnum = "positive"
	VIAOutcomePositiveInvasiveCancer VIAOutcomeEnum = "suspicious_for_cancer"
)

// IsValid checks validity of a VIAOutcomeEnum enum
func (c VIAOutcomeEnum) IsValid() bool {
	switch c {
	case VIAOutcomeNegative, VIAOutcomePositive, VIAOutcomePositiveInvasiveCancer:
		return true
	}

	return false
}

// String converts VIA to string
func (c VIAOutcomeEnum) String() string {
	return string(c)
}

// MarshalGQL writes the VIA as a quoted string
func (c VIAOutcomeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a VIA enum
func (c *VIAOutcomeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be a string")
	}

	*c = VIAOutcomeEnum(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid VIAOutcomeEnum Enum", str)
	}

	return nil
}

// ObservationStatusEnum is a FHIR enum
type ObservationStatusEnum string

const (
	// ObservationStatusEnumRegistered ...
	ObservationStatusEnumRegistered ObservationStatusEnum = "registered"
	// ObservationStatusEnumPreliminary ...
	ObservationStatusEnumPreliminary ObservationStatusEnum = "preliminary"
	// ObservationStatusEnumFinal ...
	ObservationStatusEnumFinal ObservationStatusEnum = "final"
	// ObservationStatusEnumAmended ...
	ObservationStatusEnumAmended ObservationStatusEnum = "amended"
	// ObservationStatusEnumCorrected ...
	ObservationStatusEnumCorrected ObservationStatusEnum = "corrected"
	// ObservationStatusEnumCancelled ...
	ObservationStatusEnumCancelled ObservationStatusEnum = "cancelled"
	// ObservationStatusEnumEnteredInError ...
	ObservationStatusEnumEnteredInError ObservationStatusEnum = "entered_in_error"
	// ObservationStatusEnumUnknown ...
	ObservationStatusEnumUnknown ObservationStatusEnum = "unknown"
)

// AllObservationStatusEnum ...
var AllObservationStatusEnum = []ObservationStatusEnum{
	ObservationStatusEnumRegistered,
	ObservationStatusEnumPreliminary,
	ObservationStatusEnumFinal,
	ObservationStatusEnumAmended,
	ObservationStatusEnumCorrected,
	ObservationStatusEnumCancelled,
	ObservationStatusEnumEnteredInError,
	ObservationStatusEnumUnknown,
}

// IsValid ...
func (e ObservationStatusEnum) IsValid() bool {
	switch e {
	case ObservationStatusEnumRegistered, ObservationStatusEnumPreliminary, ObservationStatusEnumFinal, ObservationStatusEnumAmended, ObservationStatusEnumCorrected, ObservationStatusEnumCancelled, ObservationStatusEnumEnteredInError, ObservationStatusEnumUnknown:
		return true
	}

	return false
}

// String ...
func (e ObservationStatusEnum) String() string {
	if e == ObservationStatusEnumEnteredInError {
		return "entered-in-error"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *ObservationStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ObservationStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ObservationStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the observation status to the supplied writer as a quoted string
func (e ObservationStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// SegmentationCategory models advantage segmentation categories for clients
type SegmentationCategory string

const (
	SegmentationCategoryNoRisk            SegmentationCategory = "CERVICAL_CANCER_TIPS"
	SegmentationCategoryLowRisk           SegmentationCategory = "CERVICAL_CANCER_LOW_RISK"
	SegmentationCategoryHighRiskPositive  SegmentationCategory = "CERVICAL_CANCER_POSITIVE"
	SegmentationCategoryHighRiskNegative  SegmentationCategory = "CERVICAL_CANCER_HIGH_RISK"
	SegmentationBreastCategoryHighRisk    SegmentationCategory = "BREAST_CANCER_HIGH_RISK"
	SegmentationBreastCategoryAverageRisk SegmentationCategory = "BREAST_CANCER_AVERAGE_RISK"
	SegmentationBreastCategoryLowRisk     SegmentationCategory = "BREAST_CANCER_LOW_RISK"
)

// IsValid checks validity of a SegmentationCategory enum
func (c SegmentationCategory) IsValid() bool {
	switch c {
	case SegmentationCategoryNoRisk, SegmentationCategoryLowRisk,
		SegmentationCategoryHighRiskPositive, SegmentationCategoryHighRiskNegative,
		SegmentationBreastCategoryHighRisk, SegmentationBreastCategoryAverageRisk,
		SegmentationBreastCategoryLowRisk:
		return true
	}

	return false
}

// String converts segmentation to string
func (c SegmentationCategory) String() string {
	return string(c)
}

// MarshalGQL writes the segmentation as a quoted string
func (c SegmentationCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a segmentation enum
func (c *SegmentationCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be a string")
	}

	*c = SegmentationCategory(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid SegmentationCategory Enum", str)
	}

	return nil
}

// ReferralTypeEnum defines the type of referral
type ReferralTypeEnum string

const (
	// DIAGNOSTICS refers to diagnostics and testing
	DIAGNOSTICS ReferralTypeEnum = "DIAGNOSTICS"
	// SPECIALIST refers to a specialist consultation
	SPECIALIST ReferralTypeEnum = "SPECIALIST"
	// TREATMENT refers to treatment
	TREATMENT ReferralTypeEnum = "TREATMENT"
)

// AppointmentStatus models allowed FHIR appointment statuses
type AppointmentStatus string

const (
	AppointmentStatusProposed       AppointmentStatus = "proposed"
	AppointmentStatusPending        AppointmentStatus = "pending"
	AppointmentStatusBooked         AppointmentStatus = "booked"
	AppointmentStatusArrived        AppointmentStatus = "arrived"
	AppointmentStatusFulfilled      AppointmentStatus = "fulfilled"
	AppointmentStatusNoShow         AppointmentStatus = "noshow"
	AppointmentStatusEnteredInError AppointmentStatus = "entered_in_error"
	AppointmentStatusCheckedIn      AppointmentStatus = "checked_in"
	AppointmentStatusWaitList       AppointmentStatus = "waitlist"
)

// IsValid checks validity of a AppointmentStatus enum
func (c AppointmentStatus) IsValid() bool {
	switch c {
	case AppointmentStatusProposed, AppointmentStatusPending, AppointmentStatusBooked,
		AppointmentStatusArrived, AppointmentStatusFulfilled, AppointmentStatusNoShow,
		AppointmentStatusEnteredInError, AppointmentStatusCheckedIn, AppointmentStatusWaitList:
		return true
	}

	return false
}

// String converts appointment status to string
func (c AppointmentStatus) String() string {
	switch c {
	case AppointmentStatusEnteredInError:
		return "entered-in-error"
	case AppointmentStatusCheckedIn:
		return "checked-in"
	default:
		return string(c)
	}
}

// MarshalGQL writes the appointment as a quoted string
func (c AppointmentStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a appointment enum
func (c *AppointmentStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be a string")
	}

	*c = AppointmentStatus(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid Appointment Status Enum", str)
	}

	return nil
}

// ParticipationStatus models allowed FHIR participation statuses
type ParticipationStatus string

const (
	ParticipationStatusAccepted    ParticipationStatus = "accepted"
	ParticipationStatusDeclined    ParticipationStatus = "declined"
	ParticipationStatusTentative   ParticipationStatus = "tentative"
	ParticipationStatusNeedsAction ParticipationStatus = "needs_action"
)

// IsValid checks validity of a ParticipationStatus enum
func (c ParticipationStatus) IsValid() bool {
	switch c {
	case ParticipationStatusAccepted, ParticipationStatusDeclined, ParticipationStatusTentative,
		ParticipationStatusNeedsAction:
		return true
	}

	return false
}

// String converts Participation status to string
func (c ParticipationStatus) String() string {
	switch c {
	case ParticipationStatusNeedsAction:
		return "needs-action"
	default:
		return string(c)
	}
}

// MarshalGQL writes the Participation as a quoted string
func (c ParticipationStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a Participation enum
func (c *ParticipationStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be a string")
	}

	*c = ParticipationStatus(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid Participation Status Enum", str)
	}

	return nil
}

// ParticipantRequired models allowed FHIR participant states
type ParticipantRequired string

const (
	ParticipantRequiredRequired        ParticipantRequired = "required"
	ParticipantRequiredOptional        ParticipantRequired = "optional"
	ParticipantRequiredInformationOnly ParticipantRequired = "information_only"
)

// IsValid checks validity of a ParticipantRequired enum
func (c ParticipantRequired) IsValid() bool {
	switch c {
	case ParticipantRequiredRequired, ParticipantRequiredOptional, ParticipantRequiredInformationOnly:
		return true
	}

	return false
}

// String converts ParticipantRequired state to string
func (c ParticipantRequired) String() string {
	switch c {
	case ParticipantRequiredInformationOnly:
		return "information-only"
	default:
		return string(c)
	}
}

// MarshalGQL writes the participant required as a quoted string
func (c ParticipantRequired) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a participant required enum
func (c *ParticipantRequired) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be a string")
	}

	*c = ParticipantRequired(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid ParticipantRequired Status Enum", str)
	}

	return nil
}

// ScreeningTypeEnum is used to define the type of Screening
type ScreeningTypeEnum string

const (
	// CervicalCancerScreeningTypeEnum is cervical Screening type
	CervicalCancerScreeningTypeEnum ScreeningTypeEnum = "CERVICAL_CANCER_SCREENING"

	// BreastCancerScreeningTypeEnum is represents the breast Screening type
	BreastCancerScreeningTypeEnum ScreeningTypeEnum = "BREAST_CANCER_SCREENING"

	ProstateCancerScreeningTypeEnum ScreeningTypeEnum = "PROSTATE_CANCER_SCREENING"
)

// AllScreeningTypeEnum is a list of all possible Screening types
var AllScreeningTypeEnum = []ScreeningTypeEnum{
	CervicalCancerScreeningTypeEnum,
	BreastCancerScreeningTypeEnum,
	ProstateCancerScreeningTypeEnum,
}

// IsValid ...
func (e ScreeningTypeEnum) IsValid() bool {
	switch e {
	case CervicalCancerScreeningTypeEnum,
		BreastCancerScreeningTypeEnum,
		ProstateCancerScreeningTypeEnum:
		return true
	}

	return false
}

// String ...
func (e ScreeningTypeEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *ScreeningTypeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ScreeningTypeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Screening type enum", str)
	}

	return nil
}

// MarshalGQL writes the Screening type to the supplied writer as a quoted string
func (e ScreeningTypeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Text returns a human friendly representation of the screening type
func (e ScreeningTypeEnum) Text() string {
	switch e {
	case BreastCancerScreeningTypeEnum:
		return "Breast Cancer Screening"
	case CervicalCancerScreeningTypeEnum:
		return "Cervical Cancer Screening"
	case ProstateCancerScreeningTypeEnum:
		return "Prostate Cancer Screening"
	}

	return "unknown screening type"
}

// TaskStatus is used to define the allowed statuses of a task
type TaskStatus string

// NOTE: This list is not exhaustive. More status can be found here: https://hl7.org/fhir/R4/valueset-task-status.html
const (
	// CompletedTasksStatus is the status of a completed task
	CompletedTasksStatus TaskStatus = "completed"

	// RequestedTasksStatus is the status of a task in a requested state
	RequestedTasksStatus TaskStatus = "requested"

	// CancelledTasksStatus is the status of a task in a cancelled state
	CancelledTasksStatus TaskStatus = "cancelled"

	// ReceivedTasksStatus is the status of a task in a received state
	ReceivedTasksStatus TaskStatus = "received"

	// AcceptedTasksStatus is the status of a task in a accepted state
	AcceptedTasksStatus TaskStatus = "accepted"

	// ReadyTasksStatus is the status of a task in a ready state
	ReadyTasksStatus TaskStatus = "ready"
)

// AllTaskStatus is a list of all possible task statuses
var AllTaskStatus = []TaskStatus{
	CompletedTasksStatus,
	RequestedTasksStatus,
	CancelledTasksStatus,
	ReceivedTasksStatus,
	AcceptedTasksStatus,
	ReadyTasksStatus,
}

// IsValid ...
func (t TaskStatus) IsValid() bool {
	switch t {
	case CompletedTasksStatus, RequestedTasksStatus, CancelledTasksStatus,
		AcceptedTasksStatus, ReceivedTasksStatus, ReadyTasksStatus:
		return true
	}

	return false
}

// String ...
func (t TaskStatus) String() string {
	return string(t)
}

// UnmarshalGQL ...
func (t *TaskStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*t = TaskStatus(str)
	if !t.IsValid() {
		return fmt.Errorf("%s is not a valid task status", str)
	}

	return nil
}

// MarshalGQL writes the task status to the supplied writer as a quoted string
func (t TaskStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}

// ServiceRequestStatus is used to define the allowed statuses of a task
type ServiceRequestStatus string

const (
	ServiceRequestStatusDraft          ServiceRequestStatus = "draft"
	ServiceRequestStatusActive         ServiceRequestStatus = "active"
	ServiceRequestStatusOnHold         ServiceRequestStatus = "on_hold"
	ServiceRequestStatusRevoked        ServiceRequestStatus = "revoked"
	ServiceRequestStatusCompleted      ServiceRequestStatus = "completed"
	ServiceRequestStatusEnteredInError ServiceRequestStatus = "entered_in_error"
	ServiceRequestStatusUnknown        ServiceRequestStatus = "unknown"
)

// AllServiceRequestStatus is a list of all possible service request statuses
var AllServiceRequestStatus = []ServiceRequestStatus{
	ServiceRequestStatusDraft,
	ServiceRequestStatusActive,
	ServiceRequestStatusOnHold,
	ServiceRequestStatusRevoked,
	ServiceRequestStatusCompleted,
	ServiceRequestStatusEnteredInError,
	ServiceRequestStatusEnteredInError,
}

// IsValid ...
func (t ServiceRequestStatus) IsValid() bool {
	switch t {
	case ServiceRequestStatusDraft, ServiceRequestStatusActive,
		ServiceRequestStatusOnHold, ServiceRequestStatusRevoked,
		ServiceRequestStatusCompleted, ServiceRequestStatusEnteredInError,
		ServiceRequestStatusUnknown:
		return true
	}

	return false
}

// String ...
func (t ServiceRequestStatus) String() string {
	if t == ServiceRequestStatusEnteredInError {
		return "entered-in-error"
	}

	if t == ServiceRequestStatusOnHold {
		return "on-hold"
	}

	return string(t)
}

// UnmarshalGQL ...
func (t *ServiceRequestStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*t = ServiceRequestStatus(str)
	if !t.IsValid() {
		return fmt.Errorf("%s is not a valid service request status", str)
	}

	return nil
}

// MarshalGQL writes the task status to the supplied writer as a quoted string
func (t ServiceRequestStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}

// DayOfWeek is a list of all the 7 days in a week.
type DayOfWeek string

// DayOfWeek constants
const (
	DayOfWeekMonday    DayOfWeek = "MONDAY"
	DayOfWeekTuesday   DayOfWeek = "TUESDAY"
	DayOfWeekWednesday DayOfWeek = "WEDNESDAY"
	DayOfWeekThursday  DayOfWeek = "THURSDAY"
	DayOfWeekFriday    DayOfWeek = "FRIDAY"
	DayOfWeekSaturday  DayOfWeek = "SATURDAY"
	DayOfWeekSunday    DayOfWeek = "SUNDAY"
)

// IsValid returns true if a user is valid
func (m DayOfWeek) IsValid() bool {
	switch m {
	case DayOfWeekMonday, DayOfWeekTuesday, DayOfWeekWednesday, DayOfWeekThursday, DayOfWeekFriday, DayOfWeekSaturday, DayOfWeekSunday:
		return true
	}

	return false
}

func (m DayOfWeek) Code() string {
	switch m {
	case DayOfWeekMonday:
		return "mon"
	case DayOfWeekTuesday:
		return "tue"
	case DayOfWeekWednesday:
		return "wed"
	case DayOfWeekThursday:
		return "thu"
	case DayOfWeekFriday:
		return "fri"
	case DayOfWeekSaturday:
		return "sat"
	case DayOfWeekSunday:
		return "sun"
	}

	return "<unknown>"
}
func (m DayOfWeek) Display() string {
	switch m {
	case DayOfWeekMonday:
		return "Monday"
	case DayOfWeekTuesday:
		return "Tuesday"
	case DayOfWeekWednesday:
		return "Wednesday"
	case DayOfWeekThursday:
		return "Thursday"
	case DayOfWeekFriday:
		return "Friday"
	case DayOfWeekSaturday:
		return "Saturday"
	case DayOfWeekSunday:
		return "Sunday"
	}

	return "<unknown>"
}

// UnmarshalGQL converts the supplied value to a user type.
func (m *DayOfWeek) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = DayOfWeek(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid day of the week", str)
	}

	return nil
}

// MarshalGQL writes the user type to the supplied
func (m DayOfWeek) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.Display()))
}

// Terminologies is a list of all the clinical terminologies available
type Terminologies string

// Terminologies constants
const (
	TerminologiesCIEL Terminologies = "CIEL"
)

// IsValid returns true if a user is valid
func (m Terminologies) IsValid() bool {
	return m == TerminologiesCIEL
}

func (m Terminologies) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a user type.
func (m *Terminologies) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = Terminologies(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid terminology", str)
	}

	return nil
}

// MarshalGQL writes the user type to the supplied
func (m Terminologies) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}

// FacilityIdentifierType is a list of all the facility identifier types.
type FacilityIdentifierType string

const (
	// FacilityIdentifierTypeMFLCode represents the mfl facility identifier type
	FacilityIdentifierTypeMFLCode FacilityIdentifierType = "MFL_CODE"

	// FacilityIdentifierTypeHealthCRM represents the health crm facility identifier type
	FacilityIdentifierTypeHealthCRM FacilityIdentifierType = "HEALTH_CRM"

	// FacilityIdentifierTypeSladeCode represents the health crm slade code identifier type
	FacilityIdentifierTypeSladeCode FacilityIdentifierType = "SLADE_CODE"
)

// IsValid returns true if a facility identifier type is valid
func (f FacilityIdentifierType) IsValid() bool {
	switch f {
	case FacilityIdentifierTypeMFLCode, FacilityIdentifierTypeHealthCRM, FacilityIdentifierTypeSladeCode:
		return true
	}

	return false
}

func (f FacilityIdentifierType) String() string {
	return string(f)
}

// UnmarshalGQL converts the supplied value to a facility identifier type.
func (f *FacilityIdentifierType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*f = FacilityIdentifierType(str)
	if !f.IsValid() {
		return fmt.Errorf("%s is not a valid FacilityIdentifierType", str)
	}

	return nil
}

// MarshalGQL writes the facility identifier type to the supplied
func (f FacilityIdentifierType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}

// Country is is a list of allowed countries
type Country string

const (
	// CountryKenya represents Kenya
	CountryKenya Country = "KE"
)

// IsValid returns true if a country is valid
func (f Country) IsValid() bool {
	return f == CountryKenya
}

func (f Country) String() string {
	return string(f)
}

// UnmarshalGQL converts the supplied value to a country
func (f *Country) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*f = Country(str)
	if !f.IsValid() {
		return fmt.Errorf("%s is not a valid Country", str)
	}

	return nil
}

// MarshalGQL writes the country to the supplied
func (f Country) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}

// MedicationRequestPriority represents the possible priorities of a MedicationRequest Priority
type MedicationRequestPriority string

const (
	MedicationRequestPriorityRoutine MedicationRequestPriority = "routine"
	MedicationRequestPriorityUrgent  MedicationRequestPriority = "urgent"
	MedicationRequestPriorityAsap    MedicationRequestPriority = "asap"
	MedicationRequestPriorityStat    MedicationRequestPriority = "stat"
)

// AllMedicationRequestPriority is a list of all possible medication request statuses
var AllMedicationRequestPriority = []MedicationRequestPriority{
	MedicationRequestPriorityRoutine,
	MedicationRequestPriorityUrgent,
	MedicationRequestPriorityAsap,
	MedicationRequestPriorityStat,
}

// IsValid ...
func (m MedicationRequestPriority) IsValid() bool {
	switch m {
	case MedicationRequestPriorityRoutine, MedicationRequestPriorityUrgent,
		MedicationRequestPriorityAsap, MedicationRequestPriorityStat:
		return true
	}

	return false
}

// UnmarshalGQL ...
func (m *MedicationRequestPriority) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = MedicationRequestPriority(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid task status", str)
	}

	return nil
}

// String ...
func (m MedicationRequestPriority) String() string {
	return string(m)
}

// MarshalGQL writes the medication status to the supplied writer as a quoted string
func (m MedicationRequestPriority) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}

// ObservationCategory represents the categories which an observation may fall under
type ObservationCategory string

const (
	ObservationCategoryLab        ObservationCategory = "laboratory"
	ObservationCategoryExam       ObservationCategory = "exam"
	ObservationCategoryVitalSigns ObservationCategory = "vital_signs"
)

// AllObservationCategory is a list of allowed observation categories
var AllObservationCategory = []ObservationCategory{
	ObservationCategoryExam,
	ObservationCategoryLab,
	ObservationCategoryVitalSigns,
}

// IsValid ...
func (o ObservationCategory) IsValid() bool {
	switch o {
	case ObservationCategoryExam, ObservationCategoryLab,
		ObservationCategoryVitalSigns:
		return true
	}

	return false
}

// UnmarshalGQL ...
func (o *ObservationCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*o = ObservationCategory(str)
	if !o.IsValid() {
		return fmt.Errorf("%s is not a valid obs category", str)
	}

	return nil
}

// String ...
func (o ObservationCategory) String() string {
	return string(o)
}

// String ...
func (o ObservationCategory) Text() string {
	if o == ObservationCategoryVitalSigns {
		return "vital-signs"
	}

	return string(o)
}

// MarshalGQL writes the medication status to the supplied writer as a quoted string
func (o ObservationCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(o.String()))
}

// PSAType defines the type of PSA test
type PSAType string

const (
	// WholeBlood ...
	WholeBlood PSAType = "WHOLE_BLOOD"
	// ProstaticSerumAntigen ...
	ProstaticSerumAntigen PSAType = "PROSTATIC_SERUM_ANTIGEN"
)

// ObservationSubtype defines the type of HPV test
type ObservationSubtype string

const (
	// HPV_PCR_DNA test ...
	HPV_PCR_DNA ObservationSubtype = "HPV_PCR_DNA"
	// HPV_ONCOPROTEIN ...
	HPV_ONCOPROTEIN ObservationSubtype = "HPV_ONCOPROTEIN"
)

// SpecialityEnum
type SpecialityEnum string

const (
	GeneralPractioner SpecialityEnum = "general-practitioner"
	Pediatrician      SpecialityEnum = "pediatrician"
	Obstetrician      SpecialityEnum = "obstetrician"
	Surgeon           SpecialityEnum = "surgeon"
	Anesthesiologist  SpecialityEnum = "anesthesiologist"
	Radiologist       SpecialityEnum = "radiologist"
	Dentist           SpecialityEnum = "dentist"
)

func (s SpecialityEnum) IsValid() bool {
	switch s {
	case GeneralPractioner, Dentist, Pediatrician, Obstetrician, Surgeon, Anesthesiologist, Radiologist:
		return true
	}

	return false
}

func (s SpecialityEnum) Code() string {
	switch s {
	case Dentist:
		return "dentist"
	case Anesthesiologist:
		return "anesthesiologist"
	case Radiologist:
		return "radiologist"
	case Pediatrician:
		return "pediatrician"
	case GeneralPractioner:
		return "general-practitioner"
	case Surgeon:
		return "surgeon"
	case Obstetrician:
		return "obstetrician"
	}

	return "unknown"
}

func (s SpecialityEnum) Display() string {
	switch s {
	case GeneralPractioner:
		return "General Practice"
	case Dentist:
		return "Dentist"
	case Pediatrician:
		return "Pediatrician"
	case Obstetrician:
		return "Obstetrician"
	case Surgeon:
		return "Surgeon"
	case Radiologist:
		return "Radiologist"
	}

	return "Unknown"
}

type ObservationConceptEnum string

const (
	ObservationConceptEnumTemperature                  ObservationConceptEnum = "TEMPERATURE"
	ObservationConceptEnumHeight                       ObservationConceptEnum = "HEIGHT"
	ObservationConceptEnumWeight                       ObservationConceptEnum = "WEIGHT"
	ObservationConceptEnumRespiratoryRate              ObservationConceptEnum = "RESPIRATORY_RATE"
	ObservationConceptEnumPulseRate                    ObservationConceptEnum = "PULSE_RATE"
	ObservationConceptEnumBloodPressure                ObservationConceptEnum = "BLOOD_PRESSURE"
	ObservationConceptEnumBMI                          ObservationConceptEnum = "BMI"
	ObservationConceptEnumViralLoad                    ObservationConceptEnum = "VIRAL_LOAD"
	ObservationConceptEnumMUAC                         ObservationConceptEnum = "MUAC"
	ObservationConceptEnumOxygenSaturation             ObservationConceptEnum = "OXYGEN_SATURATION"
	ObservationConceptEnumBloodSugar                   ObservationConceptEnum = "BLOOD_SUGAR"
	ObservationConceptEnumLastMenstrualPeriod          ObservationConceptEnum = "LAST_MENSTRUAL_PERIOD"
	ObservationConceptEnumDiastolicBloodPressure       ObservationConceptEnum = "DIASTOLIC_BLOOD_PRESSURE"
	ObservationConceptEnumColposcopy                   ObservationConceptEnum = "COLPOSCOPY"
	ObservationConceptEnumVIA                          ObservationConceptEnum = "VIA"
	ObservationConceptEnumImmunoHistoChemistry         ObservationConceptEnum = "IMMUNO_HISTO_CHEMISTRY"
	ObservationConceptEnumPostCoitalBleeding           ObservationConceptEnum = "POST_COITAL_BLEEDING"
	ObservationConceptEnumHistoryOfPresentingIllnesses ObservationConceptEnum = "HISTORY_OF_PRESENTING_ILLNESS"
	ObservationConceptEnumPastMedicalSurgicalHistory   ObservationConceptEnum = "PAST_MEDICAL_AND_SURGICAL_HISTORY"
	ObservationConceptEnumFamilySocialHistory          ObservationConceptEnum = "FAMILY_AND_SOCIAL_HISTORY"
	ObservationConceptEnumHPV                          ObservationConceptEnum = "HPV"
	ObservationConceptEnumChiefComplaint               ObservationConceptEnum = "CHIEF_COMPLAINT"
)

func (o ObservationConceptEnum) String() string {
	return string(o)
}

type LabTestTypeEnum string

const (
	MammogramTest               LabTestTypeEnum = "MAMMOGRAM"
	BiopsyTest                  LabTestTypeEnum = "BIOPSY"
	MRITest                     LabTestTypeEnum = "MRI"
	UltrasoundTest              LabTestTypeEnum = "ULTRASOUND"
	CBETest                     LabTestTypeEnum = "CBE"
	PapSmearTest                LabTestTypeEnum = "PAPSMEAR"
	WholeBloodTest              LabTestTypeEnum = "WHOLE_BLOOD"
	ProstaticSerumAntigenTest   LabTestTypeEnum = "PROSTATIC_SERUM_ANTIGEN"
	HPV_PCR_DNA_Test            LabTestTypeEnum = "HPV_PCR_DNA"
	HPV_ONCOPROTEIN_Test        LabTestTypeEnum = "HPV_ONCOPROTEIN"
	IHCProgesteroneReceptorTest LabTestTypeEnum = "IHC_PROGESTERONE_RECEPTOR"
	IHCEstrogenReceptorTest     LabTestTypeEnum = "IHC_ESTROGEN_RECEPTOR"
	IHCHER2ReceptorTest         LabTestTypeEnum = "IHC_HER2"
	IHCKi67Test                 LabTestTypeEnum = "IHC_KI67"
)

// Display returns a human-readable label for a lab test type. Known enum values map to curated
// labels that preserve medical acronyms (IHC, HER2, HPV, PCR, DNA, MRI, CBE) which naive
// title-casing would mangle. Any other value falls back to a generic humanisation: underscores
// become spaces and each word is title-cased (e.g. "SOME_NEW_TEST" -> "Some New Test").
func (l LabTestTypeEnum) Display() string {
	switch l {
	case MammogramTest:
		return "Mammogram"
	case BiopsyTest:
		return "Biopsy"
	case MRITest:
		return "MRI"
	case UltrasoundTest:
		return "Ultrasound"
	case CBETest:
		return "CBE"
	case PapSmearTest:
		return "Pap Smear"
	case WholeBloodTest:
		return "Whole Blood"
	case ProstaticSerumAntigenTest:
		return "Prostatic Serum Antigen"
	case HPV_PCR_DNA_Test:
		return "HPV PCR DNA"
	case HPV_ONCOPROTEIN_Test:
		return "HPV Oncoprotein"
	case IHCProgesteroneReceptorTest:
		return "IHC Progesterone Receptor"
	case IHCEstrogenReceptorTest:
		return "IHC Estrogen Receptor"
	case IHCHER2ReceptorTest:
		return "IHC HER2"
	case IHCKi67Test:
		return "IHC Ki67"
	}

	return humanizeToken(string(l))
}

// Code returns the LOINC code associated with a lab test type, as defined in the application
// defaults. It is the code companion to Display: given a test type it yields the concrete code
// used to identify the test in FHIR. Unmapped values return an empty string.
func (l LabTestTypeEnum) Code() string {
	switch l {
	case MammogramTest:
		return common.MammogramTerminologyCode
	case BiopsyTest:
		return common.BiopsyTerminologySystem
	case MRITest:
		return common.MRITerminologySystem
	case UltrasoundTest:
		return common.ChestUltrasoundTerminologySystem
	case CBETest:
		return common.BreastExaminationLOINCTerminologySystem
	case PapSmearTest:
		return common.PapSmearTerminologyCode
	case WholeBloodTest:
		return common.WholeBloodTerminologyCode
	case ProstaticSerumAntigenTest:
		return common.ProstateCancerTerminologyCode
	case HPV_PCR_DNA_Test:
		return common.HPV_PCR_DNATerminologyCode
	case HPV_ONCOPROTEIN_Test:
		return common.HPV_OncoproteinTerminologyCode
	case IHCProgesteroneReceptorTest:
		return common.IHCProgesteroneReceptorLOINCCode
	case IHCEstrogenReceptorTest:
		return common.IHCEstrogenReceptorLOINCCode
	case IHCHER2ReceptorTest:
		return common.HER2LOINCCode
	case IHCKi67Test:
		return common.Ki67LOINCCode
	}

	return ""
}

// labTestTypes enumerates every known lab test type. It backs LabTestTypeFromString so the
// resolver stays in sync with the constants above without a second hand-maintained mapping.
var labTestTypes = []LabTestTypeEnum{
	MammogramTest, BiopsyTest, MRITest, UltrasoundTest, CBETest, PapSmearTest,
	WholeBloodTest, ProstaticSerumAntigenTest, HPV_PCR_DNA_Test, HPV_ONCOPROTEIN_Test,
	IHCProgesteroneReceptorTest, IHCEstrogenReceptorTest, IHCHER2ReceptorTest, IHCKi67Test,
}

// LabTestTypeFromString resolves a raw enum value (e.g. "IHC_HER2") or the human-readable label
// produced by Display (e.g. "IHC HER2") back to its canonical LabTestTypeEnum. This lets callers
// recover the test type regardless of whether the source stored the enum verbatim or its display
// label. Comparison is case-insensitive. When no known test type matches, the input is returned
// cast verbatim so unmapped values still surface to callers unchanged.
func LabTestTypeFromString(value string) LabTestTypeEnum {
	for _, testType := range labTestTypes {
		if strings.EqualFold(string(testType), value) || strings.EqualFold(testType.Display(), value) {
			return testType
		}
	}

	return LabTestTypeEnum(value)
}

// humanizeToken converts a raw snake/kebab-cased token (e.g. "SOME_NEW_TEST") into a readable
// label ("Some New Test"). It is the fallback for values not covered by a curated mapping.
func humanizeToken(value string) string {
	words := strings.FieldsFunc(value, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	for i, word := range words {
		lower := strings.ToLower(word)
		words[i] = strings.ToUpper(lower[:1]) + lower[1:]
	}

	return strings.Join(words, " ")
}

type PatientVitalSignsEnum string

const (
	PatientHeight                 PatientVitalSignsEnum = "HEIGHT"
	PatientWeight                 PatientVitalSignsEnum = "WEIGHT"
	PatientTemperature            PatientVitalSignsEnum = "TEMPERATURE"
	PatientDiastolicBloodPressure PatientVitalSignsEnum = "DIASTOLIC_BLOOD_PRESSURE"
	PatientSystolicBloodPressure  PatientVitalSignsEnum = "SYSTOLIC_BLOOD_PRESSURE"
	PatientRespiratoryRate        PatientVitalSignsEnum = "RESPIRATORY_RATE"
	PatientOxygenSaturation       PatientVitalSignsEnum = "OXYGEN_SATURATION"
	PatientPulseRate              PatientVitalSignsEnum = "PULSE_RATE"
	PatientViralLoad              PatientVitalSignsEnum = "VIRAL_LOAD"
	PatientMidUpperCircumference  PatientVitalSignsEnum = "MUAC"
	PatientMenstrualPeriod        PatientVitalSignsEnum = "MENSTRUAL_PERIOD"
	PatientBloodSugar             PatientVitalSignsEnum = "BLOOD_SUGAR"
	PatientBMI                    PatientVitalSignsEnum = "BMI"
)

func (p PatientVitalSignsEnum) IsValid() bool {
	switch p {
	case PatientHeight, PatientMenstrualPeriod, PatientWeight,
		PatientBloodSugar, PatientDiastolicBloodPressure, PatientPulseRate,
		PatientTemperature, PatientMidUpperCircumference, PatientOxygenSaturation,
		PatientSystolicBloodPressure, PatientRespiratoryRate, PatientBMI, PatientViralLoad:
		return true
	}

	return false
}
