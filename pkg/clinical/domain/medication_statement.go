package domain

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRMedicationStatement definition: A record of a medication that is being consumed by a patient.
// A MedicationStatement may indicate that the patient may be taking the medication now or in the past or in the future.
type FHIRMedicationStatement struct {
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrative `json:"text,omitempty"`

	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	BasedOn []*FHIRReference `json:"basedOn,omitempty"`

	PartOf []*FHIRReference `json:"partOf,omitempty"`

	Status *MedicationStatementStatusEnum `json:"status,omitempty"`

	StatusReason []*FHIRCodeableConcept `json:"statusReason,omitempty"`

	Category []*FHIRCodeableConcept `json:"category,omitempty"`

	MedicationCodeableConcept *FHIRCodeableConcept `json:"medicationCodeableConcept,omitempty"`

	Medication *FHIRCodeableReference `json:"medication,omitempty"`

	Subject *FHIRReference `json:"subject,omitempty"`

	Context *FHIRReference `json:"context,omitempty"`

	EffectiveDateTime *scalarutils.Date `json:"effectiveDateTime,omitempty"`

	EffectivePeriod *FHIRPeriod `json:"effectivePeriod,omitempty"`

	DateAsserted *scalarutils.Date `json:"dateAsserted,omitempty"`

	InformationSource *FHIRReference `json:"informationSource,omitempty"`

	DerivedFrom []*FHIRReference `json:"derivedFrom,omitempty"`

	ReasonCode []*FHIRCodeableConcept `json:"reasonCode,omitempty"`

	ReasonReference []*FHIRReference `json:"reasonReference,omitempty"`

	Note []*FHIRAnnotation `json:"note,omitempty"`

	Dosage []*FHIRDosage `json:"dosage,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// FHIRMedicationStatementInput ...
type FHIRMedicationStatementInput struct {
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrativeInput `json:"text,omitempty"`

	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	BasedOn []*FHIRReferenceInput `json:"basedOn,omitempty"`

	PartOf []*FHIRReferenceInput `json:"partOf,omitempty"`

	Status *MedicationStatementStatusEnum `json:"status,omitempty"`

	StatusReason []*FHIRCodeableConceptInput `json:"statusReason,omitempty"`

	Category *FHIRCodeableConceptInput `json:"category,omitempty"`

	MedicationCodeableConcept *FHIRCodeableConceptInput `json:"medicationCodeableConcept,omitempty"`

	MedicationReference *FHIRMedicationInput `json:"medicationReference,omitempty"`

	Subject *FHIRReferenceInput `json:"subject,omitempty"`

	Context *FHIRReferenceInput `json:"context,omitempty"`

	EffectiveDateTime *scalarutils.Date `json:"effectiveDateTime,omitempty"`

	EffectivePeriod *FHIRPeriodInput `json:"effectivePeriod,omitempty"`

	DateAsserted *scalarutils.Date `json:"dateAsserted,omitempty"`

	InformationSource *FHIRReferenceInput `json:"informationSource,omitempty"`

	DerivedFrom []*FHIRReferenceInput `json:"derivedFrom,omitempty"`

	ReasonCode []*FHIRCodeableConceptInput `json:"reasonCode,omitempty"`

	ReasonReference []*FHIRReferenceInput `json:"reasonReference,omitempty"`

	Note []*FHIRAnnotationInput `json:"note,omitempty"`

	Dosage []*FHIRDosageInput `json:"dosage,omitempty"`

	// Meta stores more information about the resource
	Meta FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// FHIRMedicationStatementRelayConnection is a Relay connection for MedicationStatement
type FHIRMedicationStatementRelayConnection struct {
	Edges []*FHIRMedicationStatementRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRMedicationStatementRelayEdge is a Relay edge for MedicationStatement
type FHIRMedicationStatementRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRMedicationStatement `json:"node,omitempty"`
}

// FHIRMedicationStatementRelayPayload is used to return single instances of MedicationStatement
type FHIRMedicationStatementRelayPayload struct {
	Resource *FHIRMedicationStatement `json:"resource,omitempty"`
}
