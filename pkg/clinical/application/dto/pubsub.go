package dto

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// This file will be used to hold the payload for data that has been published to
// different pubsub topics

// PatientPubSubMessage models the payload that is published to the `patient.create`
// topic
type PatientPubSubMessage struct {
	UserID      string           `json:"userID"`
	ClientID    string           `json:"clientID"`
	Name        string           `json:"name"`
	DateOfBirth time.Time        `json:"dateOfBirth"`
	Gender      enumutils.Gender `json:"gender"`
	Active      bool             `json:"active"`
	PhoneNumber string           `json:"phoneNumber"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`
}

// FacilityPubSubMessage models the details of healthcare facilities that are on myCareHub platform.
// This will be used to create a FHIR organization
type FacilityPubSubMessage struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	Code        int     `json:"code"`
	Phone       string  `json:"phone"`
	Active      bool    `json:"active"`
	County      string  `json:"county"`
	Description string  `json:"description"`
}

// VitalSignPubSubMessage models the details that will be posted to the vitals pub/sub topic
type VitalSignPubSubMessage struct {
	Name      string    `json:"name"`
	ConceptID *string   `json:"conceptId"`
	Value     string    `json:"value"`
	Date      time.Time `json:"date"`

	PatientID string `json:"patientID"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`
}

// PatientAllergyPubSubMessage contains allergy details for a patient
type PatientAllergyPubSubMessage struct {
	Name      string          `json:"name"`
	ConceptID *string         `json:"conceptID"`
	Date      time.Time       `json:"date"`
	Reaction  AllergyReaction `json:"reaction"`
	Severity  AllergySeverity `json:"severity"`

	PatientID string `json:"patientID"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`
}

// AllergyReaction ...
type AllergyReaction struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptID"`
}

// AllergySeverity ...
type AllergySeverity struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptID"`
}

// MedicationPubSubMessage contains details for medication that a patient/client is prescribed or using
type MedicationPubSubMessage struct {
	Name      string          `json:"medication"`
	ConceptID *string         `json:"conceptId"`
	Date      time.Time       `json:"date"`
	Value     string          `json:"value"`
	Drug      *MedicationDrug `json:"drug"`

	PatientID string `json:"patientID"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`
}

// MedicationDrug ...
type MedicationDrug struct {
	ConceptID *string `json:"conceptId"`
}

// PatientTestResultPubSubMessage models details that are published to the test results topic
type PatientTestResultPubSubMessage struct {
	Name      string     `json:"name"`
	ConceptID *string    `json:"conceptId"`
	Date      time.Time  `json:"date"`
	Result    TestResult `json:"result"`

	PatientID string `json:"patientID"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`
}

// TestResult ...
type TestResult struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptId"`
}

// UpdatePatientFHIRID represents the data structure used for updating the FHIR ID in a user's profile.
type UpdatePatientFHIRID struct {
	FhirID   string `json:"fhirID"`
	ClientID string `json:"clientID"`
}

// UpdateProgramFHIRID represents the data structure used for updating the fhir id field in a program
type UpdateProgramFHIRID struct {
	ProgramID    string `json:"programID"`
	FHIRTenantID string `json:"fhirTenantID"`
}

// UpdateFacilityFHIRID represents the data structure used for updating the fhir id field in a facility
type UpdateFacilityFHIRID struct {
	FacilityID string `json:"facilityID"`
	FhirID     string `json:"fhirOrganisationID"`
}

// PatientReferralTaskPayload is used to create a patient referral task
type PatientReferralTaskPayload struct {
	Meta     *MetaInput      `json:"meta"`
	Referral *ServiceRequest `json:"referral"`
}

// ReferralReportNotification models the data that will be used to send a notification (SMS/Email) to the patient
// and the receiving facility
type ReferralReportNotification struct {
	DocumentReference domain.FHIRDocumentReference `json:"documentReference,omitempty"`
	WorkstationID     string                       `json:"workstationID,omitempty"`
	BranchID          string                       `json:"branchID,omitempty"`
	ServiceRequestID  string                       `json:"serviceRequestID,omitempty"`
}
