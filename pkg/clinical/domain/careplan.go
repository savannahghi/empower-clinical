package domain

// FHIRCarePlan is documented here http://hl7.org/fhir/StructureDefinition/CarePlan
type FHIRCarePlan struct {
	ResourceType          string                  `json:"resourceType,omitempty"`
	ID                    *string                 `json:"id,omitempty"`
	Meta                  *FHIRMetaInput          `json:"meta,omitempty"`
	ImplicitRules         *string                 `json:"implicitRules,omitempty"`
	Language              *string                 `json:"language,omitempty"`
	Text                  *FHIRNarrativeInput     `json:"text,omitempty"`
	Extension             []Extension             `json:"extension,omitempty"`
	ModifierExtension     []Extension             `json:"modifierExtension,omitempty"`
	Identifier            []FHIRIdentifierInput   `json:"identifier,omitempty"`
	InstantiatesCanonical []string                `json:"instantiatesCanonical,omitempty"`
	InstantiatesUri       []string                `json:"instantiatesUri,omitempty"`
	BasedOn               []FHIRReference         `json:"basedOn,omitempty"`
	Replaces              []FHIRReference         `json:"replaces,omitempty"`
	PartOf                []FHIRReference         `json:"partOf,omitempty"`
	Status                RequestStatus           `json:"status"`
	Intent                CarePlanIntent          `json:"intent"`
	Category              []FHIRCodeableConcept   `json:"category,omitempty"`
	Title                 *string                 `json:"title,omitempty"`
	Description           *string                 `json:"description,omitempty"`
	Subject               FHIRReference           `json:"subject"`
	Encounter             *FHIRReference          `json:"encounter,omitempty"`
	Period                *FHIRPeriod             `json:"period,omitempty"`
	Created               *string                 `json:"created,omitempty"`
	Custodian             *FHIRReference          `json:"custodian,omitempty"`
	Contributor           []FHIRReference         `json:"contributor,omitempty"`
	CareTeam              []FHIRReference         `json:"careTeam,omitempty"`
	Addresses             []FHIRCodeableReference `json:"addresses,omitempty"`
	SupportingInfo        []FHIRReference         `json:"supportingInfo,omitempty"`
	Goal                  []FHIRReference         `json:"goal,omitempty"`
	Activity              []CarePlanActivity      `json:"activity,omitempty"`
	Note                  []FHIRAnnotation        `json:"note,omitempty"`
}

type CarePlanActivity struct {
	ID                       *string                 `json:"id,omitempty"`
	Extension                []Extension             `json:"extension,omitempty"`
	ModifierExtension        []Extension             `json:"modifierExtension,omitempty"`
	PerformedActivity        []FHIRCodeableReference `json:"performedActivity,omitempty"`
	Progress                 []FHIRAnnotation        `json:"progress,omitempty"`
	PlannedActivityReference *FHIRReference          `json:"plannedActivityReference,omitempty"`
}

type OtherCarePlan FHIRCarePlan

type PagedFHIRCarePlan struct {
	CarePlan        []FHIRCarePlan `json:"carePlan,omitempty"`
	HasNextPage     bool           `json:"hasNextPage,omitempty"`
	NextCursor      string         `json:"nextCursor"`
	HasPreviousPage bool           `json:"hasPreviousPage"`
	PreviousCursor  string         `json:"previousCursor"`
	TotalCount      int            `json:"totalCount"`
}
