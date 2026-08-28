package domain

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRMedicationRequest definition: an order or request for both supply of the medication and the instructions for administration of the medication to a patient. the resource is called "medicationrequest" rather than "medicationprescription" or "medicationorder" to generalize the use across inpatient and outpatient settings, including care plans, etc., and to harmonize with workflow patterns.
type FHIRMedicationRequest struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// Identifiers associated with this medication request that are defined by business processes and/or used to refer to it when a direct URL reference to the resource itself is not appropriate. They are business identifiers assigned to this resource by the performer or other systems and remain constant as the resource is updated and propagates from server to server.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// A code specifying the current state of the order.  Generally, this will be active or completed state.
	Status *MedicationRequestStatus `json:"status,omitempty"`

	// Captures the reason for the current state of the MedicationRequest.
	StatusReason *FHIRCodeableConcept `json:"statusReason,omitempty"`

	// Whether the request is a proposal, plan, or an original order.
	Intent MedicationRequestIntent `json:"intent,omitempty"`

	// Indicates the type of medication request (for example, where the medication is expected to be consumed or administered (i.e. inpatient or outpatient)).
	Category []*FHIRCodeableConcept `json:"category,omitempty"`

	// Indicates how quickly the Medication Request should be addressed with respect to other requests.
	Priority *scalarutils.Code `json:"priority,omitempty"`

	// If true indicates that the provider is asking for the medication request not to occur.
	DoNotPerform *bool `json:"doNotPerform,omitempty"`

	// Indicates if this record was captured as a secondary 'reported' record rather than as an original primary source-of-truth record.  It may also indicate the source of the report.
	ReportedBoolean *bool `json:"reportedBoolean,omitempty"`

	// Indicates if this record was captured as a secondary 'reported' record rather than as an original primary source-of-truth record.  It may also indicate the source of the report.
	ReportedReference *FHIRReference `json:"reportedReference,omitempty"`

	// Identifies the medication being requested. This is a link to a resource that represents the medication which may be the details of the medication or simply an attribute carrying a code that identifies the medication from a known list of medications.
	MedicationCodeableConcept *FHIRCodeableConcept `json:"medicationCodeableConcept,omitempty"`

	// Identifies the medication being requested. This is a link to a resource that represents the medication which may be the details of the medication or simply an attribute carrying a code that identifies the medication from a known list of medications.
	MedicationReference *FHIRReference `json:"medicationReference,omitempty"`

	// A link to a resource representing the person or set of individuals to whom the medication will be given.
	Subject *FHIRReference `json:"subject,omitempty"`

	// The Encounter during which this [x] was created or to which the creation of this record is tightly associated.
	Encounter *FHIRReference `json:"encounter,omitempty"`

	// Include additional information (for example, patient height and weight) that supports the ordering of the medication.
	SupportingInformation []*FHIRReference `json:"supportingInformation,omitempty"`

	// The date (and perhaps time) when the prescription was initially written or authored on.
	AuthoredOn *scalarutils.DateTime `json:"authoredOn,omitempty"`

	// The individual, organization, or device that initiated the request and has responsibility for its activation.
	Requester *FHIRReference `json:"requester,omitempty"`

	// The specified desired performer of the medication treatment (e.g. the performer of the medication administration).
	Performer *FHIRReference `json:"performer,omitempty"`

	// Indicates the type of performer of the administration of the medication.
	PerformerType *FHIRCodeableConcept `json:"performerType,omitempty"`

	// The person who entered the order on behalf of another individual for example in the case of a verbal or a telephone order.
	Recorder *FHIRReference `json:"recorder,omitempty"`

	// The reason or the indication for ordering or not ordering the medication.
	ReasonCode *scalarutils.Code `json:"reasonCode,omitempty"`

	// Condition or observation that supports why the medication was ordered.
	ReasonReference []*FHIRReference `json:"reasonReference,omitempty"`

	// The URL pointing to a protocol, guideline, orderset, or other definition that is adhered to in whole or in part by this MedicationRequest.
	InstantiatesCanonical *scalarutils.Canonical `json:"instantiatesCanonical,omitempty"`

	// The URL pointing to an externally maintained protocol, guideline, orderset or other definition that is adhered to in whole or in part by this MedicationRequest.
	InstantiatesURI *scalarutils.Instant `json:"instantiatesURI,omitempty"`

	// A plan or request that is fulfilled in whole or in part by this medication request.
	BasedOn []*FHIRReference `json:"basedOn,omitempty"`

	// A shared identifier common to all requests that were authorized more or less simultaneously by a single author, representing the identifier of the requisition or prescription.
	GroupIdentifier *FHIRIdentifier `json:"groupIdentifier,omitempty"`

	// The description of the overall patte3rn of the administration of the medication to the patient.
	CourseOfTherapyType *FHIRCodeableConcept `json:"courseOfTherapyType,omitempty"`

	// Insurance plans, coverage extensions, pre-authorizations and/or pre-determinations that may be required for delivering the requested service.
	Insurance []*FHIRReference `json:"insurance,omitempty"`

	// Extra information about the prescription that could not be conveyed by the other attributes.
	Note []*FHIRAnnotation `json:"note,omitempty"`

	// Indicates how the medication is to be used by the patient.
	DosageInstruction []*FHIRDosage `json:"dosageInstruction,omitempty"`

	// Indicates the specific details for the dispense or medication supply part of a medication request (also known as a Medication Prescription or Medication Order).  Note that this information is not always sent with the order.  There may be in some settings (e.g. hospitals) institutional or system support for completing the dispense details in the pharmacy department.
	DispenseRequest *FHIRMedicationrequestDispenserequest `json:"dispenseRequest,omitempty"`

	// Indicates whether or not substitution can or should be part of the dispense. In some cases, substitution must happen, in other cases substitution must not happen. This block explains the prescriber's intent. If nothing is specified substitution may be done.
	Substitution *FHIRMedicationrequestSubstitution `json:"substitution,omitempty"`

	// A link to a resource representing an earlier order related order or prescription.
	PriorPrescription *FHIRReference `json:"priorPrescription,omitempty"`

	// Indicates an actual or potential clinical issue with or between one or more active or proposed clinical actions for a patient; e.g. Drug-drug interaction, duplicate therapy, dosage alert etc.
	DetectedIssue []*FHIRReference `json:"detectedIssue,omitempty"`

	// Links to Provenance records for past versions of this resource or fulfilling request or event resources that identify key state transitions or updates that are likely to be relevant to a user looking at the current version of the resource.
	EventHistory []*FHIRReference `json:"eventHistory,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`

	// Indicates that actual medication tied to this order
	Medication *FHIRCodeableReference `json:"medication,omitempty"`
}

// FHIRMedicationRequestInput is the input type for MedicationRequest
type FHIRMedicationRequestInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// Identifiers associated with this medication request that are defined by business processes and/or used to refer to it when a direct URL reference to the resource itself is not appropriate. They are business identifiers assigned to this resource by the performer or other systems and remain constant as the resource is updated and propagates from server to server.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// A code specifying the current state of the order.  Generally, this will be active or completed state.
	Status *MedicationRequestStatus `json:"status,omitempty"`

	// Captures the reason for the current state of the MedicationRequest.
	StatusReason *FHIRCodeableConceptInput `json:"statusReason,omitempty"`

	// Whether the request is a proposal, plan, or an original order.
	Intent MedicationRequestIntent `json:"intent,omitempty"`

	// Indicates the type of medication request (for example, where the medication is expected to be consumed or administered (i.e. inpatient or outpatient)).
	Category []*FHIRCodeableConceptInput `json:"category,omitempty"`

	// Indicates how quickly the Medication Request should be addressed with respect to other requests.
	Priority *scalarutils.Code `json:"priority,omitempty"`

	// If true indicates that the provider is asking for the medication request not to occur.
	DoNotPerform *bool `json:"doNotPerform,omitempty"`

	// Indicates if this record was captured as a secondary 'reported' record rather than as an original primary source-of-truth record.  It may also indicate the source of the report.
	ReportedBoolean *bool `json:"reportedBoolean,omitempty"`

	// Indicates if this record was captured as a secondary 'reported' record rather than as an original primary source-of-truth record.  It may also indicate the source of the report.
	ReportedReference *FHIRReferenceInput `json:"reportedReference,omitempty"`

	// Identifies the medication being requested. This is a link to a resource that represents the medication which may be the details of the medication or simply an attribute carrying a code that identifies the medication from a known list of medications.
	MedicationCodeableConcept *FHIRCodeableConceptInput `json:"medicationCodeableConcept,omitempty"`

	// Identifies the medication being requested. This is a link to a resource that represents the medication which may be the details of the medication or simply an attribute carrying a code that identifies the medication from a known list of medications.
	Medication *FHIRCodeableReferenceInput `json:"medication,omitempty"`

	// A link to a resource representing the person or set of individuals to whom the medication will be given.
	Subject *FHIRReferenceInput `json:"subject,omitempty"`

	// The Encounter during which this [x] was created or to which the creation of this record is tightly associated.
	Encounter *FHIRReferenceInput `json:"encounter,omitempty"`

	// Include additional information (for example, patient height and weight) that supports the ordering of the medication.
	SupportingInformation []*FHIRReferenceInput `json:"supportingInformation,omitempty"`

	// The date (and perhaps time) when the prescription was initially written or authored on.
	AuthoredOn *scalarutils.DateTime `json:"authoredOn,omitempty"`

	// The individual, organization, or device that initiated the request and has responsibility for its activation.
	Requester *FHIRReferenceInput `json:"requester,omitempty"`

	// The specified desired performer of the medication treatment (e.g. the performer of the medication administration).
	Performer *FHIRReferenceInput `json:"performer,omitempty"`

	// Indicates the type of performer of the administration of the medication.
	PerformerType *FHIRCodeableConceptInput `json:"performerType,omitempty"`

	// The person who entered the order on behalf of another individual for example in the case of a verbal or a telephone order.
	Recorder *FHIRReferenceInput `json:"recorder,omitempty"`

	// The reason or the indication for ordering or not ordering the medication.
	ReasonCode *scalarutils.Code `json:"reasonCode,omitempty"`

	// Condition or observation that supports why the medication was ordered.
	ReasonReference []*FHIRReferenceInput `json:"reasonReference,omitempty"`

	// The URL pointing to a protocol, guideline, orderset, or other definition that is adhered to in whole or in part by this MedicationRequest.
	InstantiatesCanonical *scalarutils.Canonical `json:"instantiatesCanonical,omitempty"`

	// The URL pointing to an externally maintained protocol, guideline, orderset or other definition that is adhered to in whole or in part by this MedicationRequest.
	InstantiatesURI *scalarutils.Instant `json:"instantiatesURI,omitempty"`

	// A plan or request that is fulfilled in whole or in part by this medication request.
	BasedOn []*FHIRReferenceInput `json:"basedOn,omitempty"`

	// A shared identifier common to all requests that were authorized more or less simultaneously by a single author, representing the identifier of the requisition or prescription.
	GroupIdentifier *FHIRIdentifierInput `json:"groupIdentifier,omitempty"`

	// The description of the overall patte3rn of the administration of the medication to the patient.
	CourseOfTherapyType *FHIRCodeableConceptInput `json:"courseOfTherapyType,omitempty"`

	// Insurance plans, coverage extensions, pre-authorizations and/or pre-determinations that may be required for delivering the requested service.
	Insurance []*FHIRReferenceInput `json:"insurance,omitempty"`

	// Extra information about the prescription that could not be conveyed by the other attributes.
	Note []*FHIRAnnotationInput `json:"note,omitempty"`

	// Indicates how the medication is to be used by the patient.
	DosageInstruction []*FHIRDosageInput `json:"dosageInstruction,omitempty"`

	// Indicates the specific details for the dispense or medication supply part of a medication request (also known as a Medication Prescription or Medication Order).  Note that this information is not always sent with the order.  There may be in some settings (e.g. hospitals) institutional or system support for completing the dispense details in the pharmacy department.
	DispenseRequest *FHIRMedicationrequestDispenserequestInput `json:"dispenseRequest,omitempty"`

	// Indicates whether or not substitution can or should be part of the dispense. In some cases, substitution must happen, in other cases substitution must not happen. This block explains the prescriber's intent. If nothing is specified substitution may be done.
	Substitution *FHIRMedicationrequestSubstitutionInput `json:"substitution,omitempty"`

	// A link to a resource representing an earlier order related order or prescription.
	PriorPrescription *FHIRReferenceInput `json:"priorPrescription,omitempty"`

	// Indicates an actual or potential clinical issue with or between one or more active or proposed clinical actions for a patient; e.g. Drug-drug interaction, duplicate therapy, dosage alert etc.
	DetectedIssue []*FHIRReferenceInput `json:"detectedIssue,omitempty"`

	// Links to Provenance records for past versions of this resource or fulfilling request or event resources that identify key state transitions or updates that are likely to be relevant to a user looking at the current version of the resource.
	EventHistory []*FHIRReferenceInput `json:"eventHistory,omitempty"`

	// Meta stores more information about the resource
	Meta FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
	Period    *FHIRPeriodInput `json:"period,omitempty"`
}

// FHIRMedicationrequestDispenserequest definition: an order or request for both supply of the medication and the instructions for administration of the medication to a patient. the resource is called "medicationrequest" rather than "medicationprescription" or "medicationorder" to generalize the use across inpatient and outpatient settings, including care plans, etc., and to harmonize with workflow patterns.
type FHIRMedicationrequestDispenserequest struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Indicates the quantity or duration for the first dispense of the medication.
	InitialFill *FHIRMedicationrequestInitialfill `json:"initialFill,omitempty"`

	// The minimum period of time that must occur between dispenses of the medication.
	DispenseInterval *FHIRDuration `json:"dispenseInterval,omitempty"`

	// This indicates the validity period of a prescription (stale dating the Prescription).
	ValidityPeriod *FHIRPeriod `json:"validityPeriod,omitempty"`

	// An integer indicating the number of times, in addition to the original dispense, (aka refills or repeats) that the patient can receive the prescribed medication. Usage Notes: This integer does not include the original order dispense. This means that if an order indicates dispense 30 tablets plus "3 repeats", then the order can be dispensed a total of 4 times and the patient can receive a total of 120 tablets.  A prescriber may explicitly say that zero refills are permitted after the initial dispense.
	NumberOfRepeatsAllowed *string `json:"numberOfRepeatsAllowed,omitempty"`

	// The amount that is to be dispensed for one fill.
	Quantity *FHIRQuantity `json:"quantity,omitempty"`

	// Identifies the period time over which the supplied product is expected to be used, or the length of time the dispense is expected to last.
	ExpectedSupplyDuration *FHIRDuration `json:"expectedSupplyDuration,omitempty"`

	// Indicates the intended dispensing Organization specified by the prescriber.
	Performer *FHIRReference `json:"performer,omitempty"`
}

// FHIRMedicationrequestDispenserequestInput is the input type for MedicationrequestDispenserequest
type FHIRMedicationrequestDispenserequestInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Indicates the quantity or duration for the first dispense of the medication.
	InitialFill *FHIRMedicationrequestInitialfillInput `json:"initialFill,omitempty"`

	// The minimum period of time that must occur between dispenses of the medication.
	DispenseInterval *FHIRDurationInput `json:"dispenseInterval,omitempty"`

	// This indicates the validity period of a prescription (stale dating the Prescription).
	ValidityPeriod *FHIRPeriodInput `json:"validityPeriod,omitempty"`

	// An integer indicating the number of times, in addition to the original dispense, (aka refills or repeats) that the patient can receive the prescribed medication. Usage Notes: This integer does not include the original order dispense. This means that if an order indicates dispense 30 tablets plus "3 repeats", then the order can be dispensed a total of 4 times and the patient can receive a total of 120 tablets.  A prescriber may explicitly say that zero refills are permitted after the initial dispense.
	NumberOfRepeatsAllowed *string `json:"numberOfRepeatsAllowed,omitempty"`

	// The amount that is to be dispensed for one fill.
	Quantity *FHIRQuantityInput `json:"quantity,omitempty"`

	// Identifies the period time over which the supplied product is expected to be used, or the length of time the dispense is expected to last.
	ExpectedSupplyDuration *FHIRDurationInput `json:"expectedSupplyDuration,omitempty"`

	// Indicates the intended dispensing Organization specified by the prescriber.
	Performer *FHIRReferenceInput `json:"performer,omitempty"`
}

// FHIRMedicationrequestInitialfill definition: an order or request for both supply of the medication and the instructions for administration of the medication to a patient. the resource is called "medicationrequest" rather than "medicationprescription" or "medicationorder" to generalize the use across inpatient and outpatient settings, including care plans, etc., and to harmonize with workflow patterns.
type FHIRMedicationrequestInitialfill struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The amount or quantity to provide as part of the first dispense.
	Quantity *FHIRQuantity `json:"quantity,omitempty"`

	// The length of time that the first dispense is expected to last.
	Duration *FHIRDuration `json:"duration,omitempty"`
}

// FHIRMedicationrequestInitialfillInput is the input type for MedicationrequestInitialfill
type FHIRMedicationrequestInitialfillInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The amount or quantity to provide as part of the first dispense.
	Quantity *FHIRQuantityInput `json:"quantity,omitempty"`

	// The length of time that the first dispense is expected to last.
	Duration *FHIRDurationInput `json:"duration,omitempty"`
}

// FHIRMedicationrequestSubstitution definition: an order or request for both supply of the medication and the instructions for administration of the medication to a patient. the resource is called "medicationrequest" rather than "medicationprescription" or "medicationorder" to generalize the use across inpatient and outpatient settings, including care plans, etc., and to harmonize with workflow patterns.
type FHIRMedicationrequestSubstitution struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// True if the prescriber allows a different drug to be dispensed from what was prescribed.
	AllowedBoolean *bool `json:"allowedBoolean,omitempty"`

	// True if the prescriber allows a different drug to be dispensed from what was prescribed.
	AllowedCodeableConcept *scalarutils.Code `json:"allowedCodeableConcept,omitempty"`

	// Indicates the reason for the substitution, or why substitution must or must not be performed.
	Reason *FHIRCodeableConcept `json:"reason,omitempty"`
}

// FHIRMedicationrequestSubstitutionInput is the input type for MedicationrequestSubstitution
type FHIRMedicationrequestSubstitutionInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// True if the prescriber allows a different drug to be dispensed from what was prescribed.
	AllowedBoolean *bool `json:"allowedBoolean,omitempty"`

	// True if the prescriber allows a different drug to be dispensed from what was prescribed.
	AllowedCodeableConcept *scalarutils.Code `json:"allowedCodeableConcept,omitempty"`

	// Indicates the reason for the substitution, or why substitution must or must not be performed.
	Reason *FHIRCodeableConceptInput `json:"reason,omitempty"`
}

// FHIRMedicationRequestRelayConnection is a Relay connection for MedicationRequest
type FHIRMedicationRequestRelayConnection struct {
	Edges []*FHIRMedicationRequestRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRMedicationRequestRelayEdge is a Relay edge for MedicationRequest
type FHIRMedicationRequestRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRMedicationRequest `json:"node,omitempty"`
}

// FHIRMedicationRequestRelayPayload is used to return single instances of MedicationRequest
type FHIRMedicationRequestRelayPayload struct {
	Resource *FHIRMedicationRequest `json:"resource,omitempty"`
}

// PagedFHIRMedicationRequest models a paginated data class for a FHIR Medication request resource
type PagedFHIRMedicationRequest struct {
	MedicationRequests []FHIRMedicationRequest
	HasNextPage        bool
	NextCursor         string
	HasPreviousPage    bool
	PreviousCursor     string
	TotalCount         int
}
