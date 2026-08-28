package domain

// Specimen is documented here http://hl7.org/fhir/StructureDefinition/Specimen
type Specimen struct {
	ID                  *string               `json:"id,omitempty"`
	Meta                *FHIRMetaInput        `json:"meta,omitempty"`
	ImplicitRules       *string               `json:"implicitRules,omitempty"`
	Language            *string               `json:"language,omitempty"`
	Text                *FHIRNarrative        `json:"text,omitempty"`
	Extension           []FHIRExtension       `json:"extension,omitempty"`
	ModifierExtension   []FHIRExtension       `json:"modifierExtension,omitempty"`
	Identifier          []FHIRIdentifier      `json:"identifier,omitempty"`
	AccessionIdentifier *FHIRIdentifier       `json:"accessionIdentifier,omitempty"`
	Status              *SpecimenStatus       `json:"status,omitempty"`
	Type                *FHIRCodeableConcept  `json:"type,omitempty"`
	Subject             *Reference            `json:"subject,omitempty"`
	ReceivedTime        *string               `json:"receivedTime,omitempty"`
	Parent              []Reference           `json:"parent,omitempty"`
	Request             []Reference           `json:"request,omitempty"`
	Combined            *SpecimenCombined     `json:"combined,omitempty"`
	Role                []FHIRCodeableConcept `json:"role,omitempty"`
	Feature             []SpecimenFeature     `json:"feature,omitempty"`
	Collection          *SpecimenCollection   `json:"collection,omitempty"`
	Processing          []SpecimenProcessing  `json:"processing,omitempty"`
	Container           []SpecimenContainer   `json:"container,omitempty"`
	Condition           []FHIRCodeableConcept `json:"condition,omitempty"`
	Note                []FHIRAnnotation      `json:"note,omitempty"`
}

type SpecimenFeature struct {
	ID                *string             `json:"id,omitempty"`
	Extension         []FHIRExtension     `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension     `json:"modifierExtension,omitempty"`
	Type              FHIRCodeableConcept `json:"type"`
	Description       string              `json:"description"`
}

type SpecimenCollection struct {
	ID                           *string                `json:"id,omitempty"`
	Extension                    []FHIRExtension        `json:"extension,omitempty"`
	ModifierExtension            []FHIRExtension        `json:"modifierExtension,omitempty"`
	Collector                    *Reference             `json:"collector,omitempty"`
	CollectedDateTime            *string                `json:"collectedDateTime,omitempty"`
	CollectedPeriod              *FHIRPeriod            `json:"collectedPeriod,omitempty"`
	Duration                     *FHIRDuration          `json:"duration,omitempty"`
	Quantity                     *FHIRQuantity          `json:"quantity,omitempty"`
	Method                       *FHIRCodeableConcept   `json:"method,omitempty"`
	Device                       *FHIRCodeableReference `json:"device,omitempty"`
	Procedure                    *FHIRReference         `json:"procedure,omitempty"`
	BodySite                     *FHIRCodeableReference `json:"bodySite,omitempty"`
	FastingStatusCodeableConcept *FHIRCodeableConcept   `json:"fastingStatusCodeableConcept,omitempty"`
	FastingStatusDuration        *FHIRDuration          `json:"fastingStatusDuration,omitempty"`
}

type SpecimenProcessing struct {
	ID                *string              `json:"id,omitempty"`
	Extension         []FHIRExtension      `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension      `json:"modifierExtension,omitempty"`
	Description       *string              `json:"description,omitempty"`
	Method            *FHIRCodeableConcept `json:"method,omitempty"`
	Additive          []FHIRReference      `json:"additive,omitempty"`
	TimeDateTime      *string              `json:"timeDateTime,omitempty"`
	TimePeriod        *FHIRPeriod          `json:"timePeriod,omitempty"`
}

type SpecimenContainer struct {
	Id                *string         `json:"id,omitempty"`
	Extension         []FHIRExtension `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension `json:"modifierExtension,omitempty"`
	Device            FHIRReference   `json:"device"`
	Location          *FHIRReference  `json:"location,omitempty"`
	SpecimenQuantity  *FHIRQuantity   `json:"specimenQuantity,omitempty"`
}

type OtherSpecimen Specimen

// FHIRSpecimenStatus is documented here http://hl7.org/fhir/ValueSet/specimen-status
type SpecimenStatus string

const (
	SpecimenStatusAvailable      SpecimenStatus = "available"
	SpecimenStatusUnavailable    SpecimenStatus = "unavailable"
	SpecimenStatusUnsatisfactory SpecimenStatus = "unsatisfactory"
	SpecimenStatusEnteredInError SpecimenStatus = "entered-in-error"
)

type SpecimenCombined string

const (
	SpecimenCombinedGrouped SpecimenCombined = "grouped"
	SpecimenCombinedPooled  SpecimenCombined = "combined"
)
