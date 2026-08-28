package domain

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIREncounter definition: an interaction between a patient and healthcare provider(s) for the purpose of providing healthcare service(s) or assessing the health status of a patient.
type FHIREncounter struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// Identifier(s) by which this encounter is known.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// planned | on-hold | discharged | discontinued | in-progress | entered-in-error | completed | cancelled | unknown +.
	Status EncounterStatusEnum `json:"status,omitempty"`

	// The status history permits the encounter resource to contain the status history without needing to read through the historical versions of the resource, or even have the server store them.
	StatusHistory []*FHIREncounterStatushistory `json:"statusHistory,omitempty"`

	// Concepts representing classification of patient encounter such as ambulatory (outpatient), inpatient, emergency, home health or others due to local variations.
	Class []FHIRCodeableConcept `json:"class,omitempty"`

	// The class history permits the tracking of the encounters transitions without needing to go  through the resource history.  This would be used for a case where an admission starts of as an emergency encounter, then transitions into an inpatient scenario. Doing this and not restarting a new encounter ensures that any lab/diagnostic results can more easily follow the patient and not require re-processing and not get lost or cancelled during a kind of discharge from emergency to inpatient.
	ClassHistory []*FHIREncounterClasshistory `json:"classHistory,omitempty"`

	// Specific type of encounter (e.g. e-mail consultation, surgical day-care, skilled nursing, rehabilitation).
	Type []*FHIRCodeableConcept `json:"type,omitempty"`

	// Broad categorization of the service that is to be provided (e.g. cardiology).
	ServiceType *FHIRCodeableConcept `json:"serviceType,omitempty"`

	// Indicates the urgency of the encounter.
	Priority *FHIRCodeableConcept `json:"priority,omitempty"`

	// The patient or group present at the encounter.
	Subject *FHIRReference `json:"subject,omitempty"`

	// Where a specific encounter should be classified as a part of a specific episode(s) of care this field should be used. This association can facilitate grouping of related encounters together for a specific purpose, such as government reporting, issue tracking, association via a common problem.  The association is recorded on the encounter as these are typically created after the episode of care and grouped on entry rather than editing the episode of care to append another encounter to it (the episode of care could span years).
	EpisodeOfCare []*FHIRReference `json:"episodeOfCare,omitempty"`

	// The request this encounter satisfies (e.g. incoming referral or procedure request).
	BasedOn []*FHIRReference `json:"basedOn,omitempty"`

	// The list of people responsible for providing the service.
	Participant []*FHIREncounterParticipant `json:"participant,omitempty"`

	// The appointment that scheduled this encounter.
	Appointment []*FHIRReference `json:"appointment,omitempty"`

	// Quantity of time the encounter lasted. This excludes the time during leaves of absence.
	Length *FHIRDuration `json:"length,omitempty"`

	// Reason the encounter takes place, expressed as a code. For admissions, this can be used for a coded admission diagnosis.
	ReasonCode *scalarutils.Code `json:"reasonCode,omitempty"`

	// Reason the encounter takes place, expressed as a code. For admissions, this can be used for a coded admission diagnosis.
	ReasonReference []*FHIRReference `json:"reasonReference,omitempty"`

	// The list of diagnosis relevant to this encounter.
	Diagnosis []*FHIREncounterDiagnosis `json:"diagnosis,omitempty"`

	// The set of accounts that may be used for billing for this Encounter.
	Account []*FHIRReference `json:"account,omitempty"`

	// Details about the admission to a healthcare service.
	Hospitalization *FHIREncounterHospitalization `json:"hospitalization,omitempty"`

	// List of locations where  the patient has been during this encounter.
	Location []*FHIREncounterLocation `json:"location,omitempty"`

	// The organization that is primarily responsible for this Encounter's services. This MAY be the same as the organization on the Patient record, however it could be different, such as if the actor performing the services was from an external organization (which may be billed separately) for an external consultation.  Refer to the example bundle showing an abbreviated set of Encounters for a colonoscopy.
	ServiceProvider *FHIRReference `json:"serviceProvider,omitempty"`

	// Another Encounter of which this encounter is a part of (administratively or in time).
	PartOf *FHIRReference `json:"partOf,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`

	ImplicitRules *string `json:"implicitRules,omitempty"`

	Language *string `json:"language,omitempty"`

	ModifierExtension []Extension `json:"modifierExtension,omitempty"`

	SubjectStatus *FHIRCodeableConcept `json:"subjectStatus,omitempty"`

	CareTeam []Reference `json:"careTeam,omitempty"`

	VirtualService []FHIRVirtualServiceDetail `json:"virtualService,omitempty"`

	ActualPeriod *FHIRPeriod `json:"actualPeriod,omitempty"`

	PlannedStartDate *string `json:"plannedStartDate,omitempty"`

	PlannedEndDate *string `json:"plannedEndDate,omitempty"`

	Reason []FHIREncounterReason `json:"reason,omitempty"`

	DietPreference []FHIRCodeableConcept `json:"dietPreference,omitempty"`

	SpecialArrangement []FHIRCodeableConcept `json:"specialArrangement,omitempty"`

	SpecialCourtesy []FHIRCodeableConcept `json:"specialCourtesy,omitempty"`

	Admission *FHIREncounterAdmission `json:"admission,omitempty"`
}

// FHIREncounterClasshistory definition: an interaction between a patient and healthcare provider(s) for the purpose of providing healthcare service(s) or assessing the health status of a patient.
type FHIREncounterClasshistory struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// inpatient | outpatient | ambulatory | emergency +.
	Class *FHIRCoding `json:"class,omitempty"`

	// The time that the episode was in the specified class.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIREncounterClasshistoryInput is the input type for EncounterClasshistory
type FHIREncounterClasshistoryInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// inpatient | outpatient | ambulatory | emergency +.
	Class *FHIRCodingInput `json:"class,omitempty"`

	// The time that the episode was in the specified class.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// FHIREncounterDiagnosis definition: an interaction between a patient and healthcare provider(s) for the purpose of providing healthcare service(s) or assessing the health status of a patient.
type FHIREncounterDiagnosis struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Reason the encounter takes place, as specified using information from another resource. For admissions, this is the admission diagnosis. The indication will typically be a Condition (with other resources referenced in the evidence.detail), or a Procedure.
	Condition *FHIRReference `json:"condition,omitempty"`

	// Role that this diagnosis has within the encounter (e.g. admission, billing, discharge …).
	Use *FHIRCodeableConcept `json:"use,omitempty"`

	// Ranking of the diagnosis (for each role type).
	Rank *string `json:"rank,omitempty"`
}

// FHIREncounterDiagnosisInput is the input type for EncounterDiagnosis
type FHIREncounterDiagnosisInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Reason the encounter takes place, as specified using information from another resource. For admissions, this is the admission diagnosis. The indication will typically be a Condition (with other resources referenced in the evidence.detail), or a Procedure.
	Condition *FHIRReferenceInput `json:"condition,omitempty"`

	// Role that this diagnosis has within the encounter (e.g. admission, billing, discharge …).
	Use *FHIRCodeableConceptInput `json:"use,omitempty"`

	// Ranking of the diagnosis (for each role type).
	Rank *string `json:"rank,omitempty"`
}

// FHIREncounterHospitalization definition: an interaction between a patient and healthcare provider(s) for the purpose of providing healthcare service(s) or assessing the health status of a patient.
type FHIREncounterHospitalization struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Pre-admission identifier.
	PreAdmissionIdentifier *FHIRIdentifier `json:"preAdmissionIdentifier,omitempty"`

	// The location/organization from which the patient came before admission.
	Origin *FHIRReference `json:"origin,omitempty"`

	// From where patient was admitted (physician referral, transfer).
	AdmitSource *FHIRCodeableConcept `json:"admitSource,omitempty"`

	// Whether this hospitalization is a readmission and why if known.
	ReAdmission *FHIRCodeableConcept `json:"reAdmission,omitempty"`

	// Diet preferences reported by the patient.
	DietPreference []*FHIRCodeableConcept `json:"dietPreference,omitempty"`

	// Special courtesies (VIP, board member).
	SpecialCourtesy []*FHIRCodeableConcept `json:"specialCourtesy,omitempty"`

	// Any special requests that have been made for this hospitalization encounter, such as the provision of specific equipment or other things.
	SpecialArrangement []*FHIRCodeableConcept `json:"specialArrangement,omitempty"`

	// Location/organization to which the patient is discharged.
	Destination *FHIRReference `json:"destination,omitempty"`

	// Category or kind of location after discharge.
	DischargeDisposition *FHIRCodeableConcept `json:"dischargeDisposition,omitempty"`
}

// FHIREncounterHospitalizationInput is the input type for EncounterHospitalization
type FHIREncounterHospitalizationInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Pre-admission identifier.
	PreAdmissionIdentifier *FHIRIdentifierInput `json:"preAdmissionIdentifier,omitempty"`

	// The location/organization from which the patient came before admission.
	Origin *FHIRReferenceInput `json:"origin,omitempty"`

	// From where patient was admitted (physician referral, transfer).
	AdmitSource *FHIRCodeableConceptInput `json:"admitSource,omitempty"`

	// Whether this hospitalization is a readmission and why if known.
	ReAdmission *FHIRCodeableConceptInput `json:"reAdmission,omitempty"`

	// Diet preferences reported by the patient.
	DietPreference []*FHIRCodeableConceptInput `json:"dietPreference,omitempty"`

	// Special courtesies (VIP, board member).
	SpecialCourtesy []*FHIRCodeableConceptInput `json:"specialCourtesy,omitempty"`

	// Any special requests that have been made for this hospitalization encounter, such as the provision of specific equipment or other things.
	SpecialArrangement []*FHIRCodeableConceptInput `json:"specialArrangement,omitempty"`

	// Location/organization to which the patient is discharged.
	Destination *FHIRReferenceInput `json:"destination,omitempty"`

	// Category or kind of location after discharge.
	DischargeDisposition *FHIRCodeableConceptInput `json:"dischargeDisposition,omitempty"`
}

// FHIREncounterInput is the input type for Encounter
type FHIREncounterInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrative `json:"text,omitempty"`

	// Identifier(s) by which this encounter is known.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// planned | on-hold | discharged | discontinued | in-progress | entered-in-error | completed | cancelled | unknown +.
	Status EncounterStatusEnum `json:"status,omitempty"`

	// The status history permits the encounter resource to contain the status history without needing to read through the historical versions of the resource, or even have the server store them.
	StatusHistory []*FHIREncounterStatushistoryInput `json:"statusHistory,omitempty"`

	// Concepts representing classification of patient encounter such as ambulatory (outpatient), inpatient, emergency, home health or others due to local variations.
	Class []*FHIRCodeableConceptInput `json:"class,omitempty"`

	// The class history permits the tracking of the encounters transitions without needing to go  through the resource history.  This would be used for a case where an admission starts of as an emergency encounter, then transitions into an inpatient scenario. Doing this and not restarting a new encounter ensures that any lab/diagnostic results can more easily follow the patient and not require re-processing and not get lost or cancelled during a kind of discharge from emergency to inpatient.
	ClassHistory []*FHIREncounterClasshistoryInput `json:"classHistory,omitempty"`

	// Specific type of encounter (e.g. e-mail consultation, surgical day-care, skilled nursing, rehabilitation).
	Type []*FHIRCodeableConceptInput `json:"type,omitempty"`

	// Broad categorization of the service that is to be provided (e.g. cardiology).
	ServiceType *FHIRCodeableConceptInput `json:"serviceType,omitempty"`

	// Indicates the urgency of the encounter.
	Priority *FHIRCodeableConceptInput `json:"priority,omitempty"`

	// The patient or group present at the encounter.
	Subject *FHIRReferenceInput `json:"subject,omitempty"`

	// Where a specific encounter should be classified as a part of a specific episode(s) of care this field should be used. This association can facilitate grouping of related encounters together for a specific purpose, such as government reporting, issue tracking, association via a common problem.  The association is recorded on the encounter as these are typically created after the episode of care and grouped on entry rather than editing the episode of care to append another encounter to it (the episode of care could span years).
	EpisodeOfCare []*FHIRReferenceInput `json:"episodeOfCare,omitempty"`

	// The request this encounter satisfies (e.g. incoming referral or procedure request).
	BasedOn []*FHIRReferenceInput `json:"basedOn,omitempty"`

	// The list of people responsible for providing the service.
	Participant []*FHIREncounterParticipantInput `json:"participant,omitempty"`

	// The appointment that scheduled this encounter.
	Appointment []*FHIRReferenceInput `json:"appointment,omitempty"`

	// Quantity of time the encounter lasted. This excludes the time during leaves of absence.
	Length *FHIRDurationInput `json:"length,omitempty"`

	// Reason the encounter takes place, expressed as a code. For admissions, this can be used for a coded admission diagnosis.
	ReasonCode *scalarutils.Code `json:"reasonCode,omitempty"`

	// Reason the encounter takes place, expressed as a code. For admissions, this can be used for a coded admission diagnosis.
	ReasonReference []*FHIRReferenceInput `json:"reasonReference,omitempty"`

	// The list of diagnosis relevant to this encounter.
	Diagnosis []*FHIREncounterDiagnosisInput `json:"diagnosis,omitempty"`

	// The set of accounts that may be used for billing for this Encounter.
	Account []*FHIRReferenceInput `json:"account,omitempty"`

	// Details about the admission to a healthcare service.
	Hospitalization *FHIREncounterHospitalizationInput `json:"hospitalization,omitempty"`

	// List of locations where  the patient has been during this encounter.
	Location []*FHIREncounterLocationInput `json:"location,omitempty"`

	// The organization that is primarily responsible for this Encounter's services. This MAY be the same as the organization on the Patient record, however it could be different, such as if the actor performing the services was from an external organization (which may be billed separately) for an external consultation.  Refer to the example bundle showing an abbreviated set of Encounters for a colonoscopy.
	ServiceProvider *FHIRReferenceInput `json:"serviceProvider,omitempty"`

	// Another Encounter of which this encounter is a part of (administratively or in time).
	PartOf *FHIRReferenceInput `json:"partOf,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`

	ImplicitRules *string `json:"implicitRules,omitempty"`

	Language *string `json:"language,omitempty"`

	ModifierExtension []FHIRExtension `json:"modifierExtension,omitempty"`

	SubjectStatus *FHIRCodeableConceptInput `json:"subjectStatus,omitempty"`

	CareTeam []FHIRReferenceInput `json:"careTeam,omitempty"`

	VirtualService []FHIRVirtualServiceDetail `json:"virtualService,omitempty"`

	ActualPeriod *FHIRPeriodInput `json:"actualPeriod,omitempty"`

	PlannedStartDate *string `json:"plannedStartDate,omitempty"`

	PlannedEndDate *string `json:"plannedEndDate,omitempty"`

	Reason []FHIREncounterReason `json:"reason,omitempty"`

	DietPreference []FHIRCodeableConceptInput `json:"dietPreference,omitempty"`

	SpecialArrangement []FHIRCodeableConceptInput `json:"specialArrangement,omitempty"`

	SpecialCourtesy []FHIRCodeableConceptInput `json:"specialCourtesy,omitempty"`

	Admission *FHIREncounterAdmission `json:"admission,omitempty"`
}

// FHIREncounterLocation definition: an interaction between a patient and healthcare provider(s) for the purpose of providing healthcare service(s) or assessing the health status of a patient.
type FHIREncounterLocation struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The location where the encounter takes place.
	Location *FHIRReference `json:"location,omitempty"`

	// The status of the participants' presence at the specified location during the period specified. If the participant is no longer at the location, then the period will have an end date/time.
	Status *EncounterLocationStatusEnum `json:"status,omitempty"`

	// This will be used to specify the required levels (bed/ward/room/etc.) desired to be recorded to simplify either messaging or query.
	PhysicalType *FHIRCodeableConcept `json:"physicalType,omitempty"`

	// Time period during which the patient was present at the location.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIREncounterLocationInput is the input type for EncounterLocation
type FHIREncounterLocationInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The location where the encounter takes place.
	Location *FHIRReferenceInput `json:"location,omitempty"`

	// The status of the participants' presence at the specified location during the period specified. If the participant is no longer at the location, then the period will have an end date/time.
	Status *EncounterLocationStatusEnum `json:"status,omitempty"`

	// This will be used to specify the required levels (bed/ward/room/etc.) desired to be recorded to simplify either messaging or query.
	PhysicalType *FHIRCodeableConceptInput `json:"physicalType,omitempty"`

	// Time period during which the patient was present at the location.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// FHIREncounterParticipant definition: an interaction between a patient and healthcare provider(s) for the purpose of providing healthcare service(s) or assessing the health status of a patient.
type FHIREncounterParticipant struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Role of participant in encounter.
	Type []*FHIRCodeableConcept `json:"type,omitempty"`

	// The period of time that the specified participant participated in the encounter. These can overlap or be sub-sets of the overall encounter's period.
	Period *FHIRPeriod `json:"period,omitempty"`

	// Persons involved in the encounter other than the patient.
	Actor *FHIRReference `json:"actor,omitempty"`
}

// FHIREncounterParticipantInput is the input type for EncounterParticipant
type FHIREncounterParticipantInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Role of participant in encounter.
	Type []*FHIRCodeableConceptInput `json:"type,omitempty"`

	// The period of time that the specified participant participated in the encounter. These can overlap or be sub-sets of the overall encounter's period.
	Period *FHIRPeriodInput `json:"period,omitempty"`

	// Persons involved in the encounter other than the patient.
	Actor *FHIRReferenceInput `json:"actor,omitempty"`
}

// FHIREncounterStatushistory definition: an interaction between a patient and healthcare provider(s) for the purpose of providing healthcare service(s) or assessing the health status of a patient.
type FHIREncounterStatushistory struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// planned | arrived | triaged | in-progress | onleave | finished | cancelled +.
	Status *EncounterStatusHistoryStatusEnum `json:"status,omitempty"`

	// The time that the episode was in the specified status.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIREncounterStatushistoryInput is the input type for EncounterStatushistory
type FHIREncounterStatushistoryInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// planned | arrived | triaged | in-progress | onleave | finished | cancelled +.
	Status *EncounterStatusHistoryStatusEnum `json:"status,omitempty"`

	// The time that the episode was in the specified status.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// FHIREncounterRelayConnection is a Relay connection for Encounter
type FHIREncounterRelayConnection struct {
	Edges []*FHIREncounterRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIREncounterRelayEdge is a Relay edge for Encounter
type FHIREncounterRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIREncounter `json:"node,omitempty"`
}

// FHIREncounterRelayPayload is used to return single instances of Encounter
type FHIREncounterRelayPayload struct {
	Resource *FHIREncounter `json:"resource,omitempty"`
}

// PagedFHIREncounter is an encounter's paginaton dataclass
type PagedFHIREncounter struct {
	Encounters      []FHIREncounter
	HasNextPage     bool
	NextCursor      string
	HasPreviousPage bool
	PreviousCursor  string
	TotalCount      int
}

// VirtualServiceDetail is documented here http://hl7.org/fhir/StructureDefinition/VirtualServiceDetail
type FHIRVirtualServiceDetail struct {
	ID                           *string                    `json:"id,omitempty"`
	Extension                    []Extension                `json:"extension,omitempty"`
	ChannelType                  *FHIRCoding                `json:"channelType,omitempty"`
	AddressURL                   *string                    `json:"addressUrl,omitempty"`
	AddressString                *string                    `json:"addressString,omitempty"`
	AddressContactPoint          *FHIRContactPoint          `json:"addressContactPoint,omitempty"`
	AddressExtendedContactDetail *FHIRExtendedContactDetail `json:"addressExtendedContactDetail,omitempty"`
	AdditionalInfo               []string                   `json:"additionalInfo,omitempty"`
	MaxParticipants              *int                       `json:"maxParticipants,omitempty"`
	SessionKey                   *string                    `json:"sessionKey,omitempty"`
}

// ExtendedContactDetail is documented here http://hl7.org/fhir/StructureDefinition/ExtendedContactDetail
type FHIRExtendedContactDetail struct {
	ID           *string              `json:"id,omitempty"`
	Extension    []Extension          `json:"extension,omitempty"`
	Purpose      *FHIRCodeableConcept `json:"purpose,omitempty"`
	Name         []FHIRHumanName      `json:"name,omitempty"`
	Telecom      []FHIRContactPoint   `json:"telecom,omitempty"`
	Address      *FHIRAddress         `json:"address,omitempty"`
	Organization *Reference           `json:"organization,omitempty"`
	Period       *FHIRPeriod          `json:"period,omitempty"`
}

type FHIREncounterReason struct {
	ID                *string                 `json:"id,omitempty"`
	Extension         []FHIRExtension         `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension         `json:"modifierExtension,omitempty"`
	Use               []FHIRCodeableConcept   `json:"use,omitempty"`
	Value             []FHIRCodeableReference `json:"value,omitempty"`
}

type FHIREncounterAdmission struct {
	ID                     *string              `json:"id,omitempty"`
	Extension              []Extension          `json:"extension,omitempty"`
	ModifierExtension      []Extension          `json:"modifierExtension,omitempty"`
	PreAdmissionIdentifier *FHIRIdentifier      `json:"preAdmissionIdentifier,omitempty"`
	Origin                 *FHIRReference       `json:"origin,omitempty"`
	AdmitSource            *FHIRCodeableConcept `json:"admitSource,omitempty"`
	ReAdmission            *FHIRCodeableConcept `json:"reAdmission,omitempty"`
	Destination            *FHIRReference       `json:"destination,omitempty"`
	DischargeDisposition   *FHIRCodeableConcept `json:"dischargeDisposition,omitempty"`
}
