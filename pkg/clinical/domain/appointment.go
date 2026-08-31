package domain

import "github.com/savannahghi/scalarutils"

// FHIRAppointment is documented here http://hl7.org/fhir/StructureDefinition/Appointment
type FHIRAppointment struct {
	ID                    *string                       `json:"id,omitempty"`
	Meta                  *FHIRMeta                     `json:"meta,omitempty"`
	ImplicitRules         *string                       `json:"implicitRules,omitempty"`
	Language              *string                       `json:"language,omitempty"`
	Text                  *FHIRNarrative                `json:"text,omitempty"`
	Extension             []*FHIRExtension              `json:"extension,omitempty"`
	ModifierExtension     []*FHIRExtension              `json:"modifierExtension,omitempty"`
	Identifier            []*FHIRIdentifier             `json:"identifier,omitempty"`
	Status                *scalarutils.Code             `json:"status"`
	CancellationReason    *FHIRCodeableConcept          `json:"cancelationReason,omitempty"`
	ServiceCategory       []*FHIRCodeableConcept        `json:"serviceCategory,omitempty"`
	ServiceType           []*FHIRCodeableConcept        `json:"serviceType,omitempty"`
	Specialty             []*FHIRCodeableConcept        `json:"specialty,omitempty"`
	AppointmentType       *FHIRCodeableConcept          `json:"appointmentType,omitempty"`
	Reason                []*FHIRCodeableReference      `json:"reason,omitempty"`
	ReasonReference       []*FHIRReference              `json:"reasonReference,omitempty"`
	Priority              *FHIRCodeableConcept          `json:"priority,omitempty"`
	Description           *string                       `json:"description,omitempty"`
	SupportingInformation []*FHIRReference              `json:"supportingInformation,omitempty"`
	Start                 *scalarutils.Instant          `json:"start,omitempty"`
	End                   *scalarutils.Instant          `json:"end,omitempty"`
	MinutesDuration       *int                          `json:"minutesDuration,omitempty"`
	Slot                  []*FHIRReference              `json:"slot,omitempty"`
	Created               *scalarutils.Date             `json:"created,omitempty"`
	Comment               *string                       `json:"comment,omitempty"`
	PatientInstruction    *string                       `json:"patientInstruction,omitempty"`
	BasedOn               []*FHIRReference              `json:"basedOn,omitempty"`
	Subject               FHIRReference                 `json:"subject,omitempty"`
	Participant           []*FHIRAppointmentParticipant `json:"participant"`
	RequestedPeriod       []*FHIRPeriod                 `json:"requestedPeriod,omitempty"`
}

// FHIRAppointmentParticipant describes the participants involves in the appointment
type FHIRAppointmentParticipant struct {
	ID                *string                `json:"id,omitempty"`
	Extension         []*FHIRExtension       `json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension       `json:"modifierExtension,omitempty"`
	Type              []*FHIRCodeableConcept `json:"type,omitempty"`
	Actor             *FHIRReference         `json:"actor,omitempty"`
	Required          *scalarutils.Code      `json:"required,omitempty"`
	Status            *scalarutils.Code      `json:"status"`
	Period            *FHIRPeriod            `json:"period,omitempty"`
}

// FHIRAppointmentInput is the input data model for FHIR resource appointment
type FHIRAppointmentInput struct {
	ID                    *string                       `json:"id,omitempty"`
	Meta                  *FHIRMetaInput                `json:"meta,omitempty"`
	ImplicitRules         *string                       `json:"implicitRules,omitempty"`
	Language              *string                       `json:"language,omitempty"`
	Text                  *FHIRNarrative                `json:"text,omitempty"`
	Extension             []*FHIRExtension              `json:"extension,omitempty"`
	ModifierExtension     []*FHIRExtension              `json:"modifierExtension,omitempty"`
	Identifier            []*FHIRIdentifierInput        `json:"identifier,omitempty"`
	Status                *scalarutils.Code             `json:"status"`
	CancellationReason    *FHIRCodeableConceptInput     `json:"cancelationReason,omitempty"`
	Class                 []*FHIRCodeableConceptInput   `json:"class,omitempty"`
	ServiceCategory       []*FHIRCodeableConceptInput   `json:"serviceCategory,omitempty"`
	ServiceType           []*FHIRCodeableConceptInput   `json:"serviceType,omitempty"`
	Specialty             []*FHIRCodeableConceptInput   `json:"specialty,omitempty"`
	AppointmentType       *FHIRCodeableConceptInput     `json:"appointmentType,omitempty"`
	Reason                []*FHIRCodeableReferenceInput `json:"reason,omitempty"`
	Priority              *FHIRCodeableConceptInput     `json:"priority,omitempty"`
	Replaces              []*FHIRCodeableReferenceInput `json:"replaces,omitempty"`
	Description           *string                       `json:"description,omitempty"`
	SupportingInformation []*FHIRReferenceInput         `json:"supportingInformation,omitempty"`
	Start                 *scalarutils.Instant          `json:"start,omitempty"`
	End                   *scalarutils.Instant          `json:"end,omitempty"`
	MinutesDuration       *int                          `json:"minutesDuration,omitempty"`
	Slot                  []*FHIRReferenceInput         `json:"slot,omitempty"`
	Created               *scalarutils.DateTime         `json:"created,omitempty"`
	PatientInstruction    *string                       `json:"patientInstruction,omitempty"`
	BasedOn               []*FHIRReferenceInput         `json:"basedOn,omitempty"`
	Participant           []*FHIRAppointmentParticipant `json:"participant"`
	RequestedPeriod       []*FHIRPeriod                 `json:"requestedPeriod,omitempty"`
	Subject               *FHIRReferenceInput           `json:"subject,omitempty"`
	Note                  []*FHIRAnnotationInput        `json:"note,omitempty"`
}
