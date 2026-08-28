package domain

import "github.com/savannahghi/scalarutils"

type FHIRMedicationDispense struct {
	ID                        *string                             `bson:"id,omitempty" json:"id,omitempty"`
	Meta                      *FHIRMeta                           `bson:"meta,omitempty" json:"meta,omitempty"`
	ImplicitRules             *string                             `bson:"implicitRules,omitempty" json:"implicitRules,omitempty"`
	Language                  *string                             `bson:"language,omitempty" json:"language,omitempty"`
	Text                      *FHIRNarrative                      `bson:"text,omitempty" json:"text,omitempty"`
	Extension                 []*FHIRExtension                    `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension         []*FHIRExtension                    `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Identifier                []*FHIRIdentifier                   `bson:"identifier,omitempty" json:"identifier,omitempty"`
	BasedOn                   []*FHIRReference                    `bson:"basedOn,omitempty" json:"basedOn,omitempty"`
	PartOf                    []*FHIRReference                    `bson:"partOf,omitempty" json:"partOf,omitempty"`
	Status                    MedicationDispenseStatusCodes       `bson:"status" json:"status"`
	NotPerformedReason        *FHIRCodeableReference              `bson:"notPerformedReason,omitempty" json:"notPerformedReason,omitempty"`
	StatusChanged             *string                             `bson:"statusChanged,omitempty" json:"statusChanged,omitempty"`
	Category                  []*FHIRCodeableConcept              `bson:"category,omitempty" json:"category,omitempty"`
	Medication                *FHIRCodeableReference              `bson:"medication" json:"medication"`
	Subject                   FHIRReference                       `bson:"subject" json:"subject"`
	Encounter                 *FHIRReference                      `bson:"encounter,omitempty" json:"encounter,omitempty"`
	SupportingInformation     []*FHIRReference                    `bson:"supportingInformation,omitempty" json:"supportingInformation,omitempty"`
	Performer                 []*FHIRMedicationDispensePerformer  `bson:"performer,omitempty" json:"performer,omitempty"`
	Location                  *FHIRReference                      `bson:"location,omitempty" json:"location,omitempty"`
	AuthorizingPrescription   []*FHIRReference                    `bson:"authorizingPrescription,omitempty" json:"authorizingPrescription,omitempty"`
	Type                      *FHIRCodeableConcept                `bson:"type,omitempty" json:"type,omitempty"`
	Quantity                  *FHIRQuantity                       `bson:"quantity,omitempty" json:"quantity,omitempty"`
	DaysSupply                *FHIRQuantity                       `bson:"daysSupply,omitempty" json:"daysSupply,omitempty"`
	Recorded                  *string                             `bson:"recorded,omitempty" json:"recorded,omitempty"`
	WhenPrepared              *string                             `bson:"whenPrepared,omitempty" json:"whenPrepared,omitempty"`
	WhenHandedOver            *scalarutils.DateTime               `bson:"whenHandedOver,omitempty" json:"whenHandedOver,omitempty"`
	Destination               *FHIRReference                      `bson:"destination,omitempty" json:"destination,omitempty"`
	Receiver                  []*FHIRReference                    `bson:"receiver,omitempty" json:"receiver,omitempty"`
	Note                      []*FHIRAnnotation                   `bson:"note,omitempty" json:"note,omitempty"`
	RenderedDosageInstruction *string                             `bson:"renderedDosageInstruction,omitempty" json:"renderedDosageInstruction,omitempty"`
	DosageInstruction         []*FHIRDosage                       `bson:"dosageInstruction,omitempty" json:"dosageInstruction,omitempty"`
	Substitution              *FHIRMedicationDispenseSubstitution `bson:"substitution,omitempty" json:"substitution,omitempty"`
	EventHistory              []*FHIRReference                    `bson:"eventHistory,omitempty" json:"eventHistory,omitempty"`
}

type FHIRMedicationDispensePerformer struct {
	ID                *string              `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []*FHIRExtension     `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension     `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Function          *FHIRCodeableConcept `bson:"function,omitempty" json:"function,omitempty"`
	Actor             *FHIRReference       `bson:"actor" json:"actor"`
}

type FHIRMedicationDispenseSubstitution struct {
	ID                *string                `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []*FHIRExtension       `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension       `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	WasSubstituted    bool                   `bson:"wasSubstituted" json:"wasSubstituted"`
	Type              *FHIRCodeableConcept   `bson:"type,omitempty" json:"type,omitempty"`
	Reason            []*FHIRCodeableConcept `bson:"reason,omitempty" json:"reason,omitempty"`
	ResponsibleParty  *FHIRReference         `bson:"responsibleParty,omitempty" json:"responsibleParty,omitempty"`
}

type PagedFHIRMedicationDispense struct {
	MedicationDispense []*FHIRMedicationDispense
	HasNextPage        bool
	NextCursor         string
	HasPreviousPage    bool
	PreviousCursor     string
	TotalCount         int
}
