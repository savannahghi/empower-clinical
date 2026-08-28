package domain

// Procedure is documented here http://hl7.org/fhir/StructureDefinition/Procedure
type FHIRProcedure struct {
	ID                 *string               `json:"id,omitempty"`
	Meta               *FHIRMeta             `json:"meta,omitempty"`
	Language           *string               `json:"language,omitempty"`
	Text               *FHIRNarrative        `json:"text,omitempty"`
	Identifier         []FHIRIdentifier      `json:"identifier,omitempty"`
	BasedOn            []Reference           `json:"basedOn,omitempty"`
	PartOf             []Reference           `json:"partOf,omitempty"`
	Status             ProcedureStatusEnum   `json:"status,omitempty"`
	StatusReason       *FHIRCodeableConcept  `json:"statusReason,omitempty"`
	Category           *FHIRCodeableConcept  `json:"category,omitempty"`
	Code               *FHIRCodeableConcept  `json:"code,omitempty"`
	Subject            FHIRReference         `json:"subject,omitempty"`
	Encounter          *FHIRReference        `json:"encounter,omitempty"`
	PerformedDateTime  *string               `json:"performedDateTime,omitempty"`
	PerformedPeriod    *FHIRPeriod           `json:"performedPeriod,omitempty"`
	PerformedString    *string               `json:"performedString,omitempty"`
	PerformedAge       *FHIRAge              `json:"performedAge,omitempty"`
	PerformedRange     *FHIRRange            `json:"performedRange,omitempty"`
	Recorder           *FHIRReference        `json:"recorder,omitempty"`
	Asserter           *FHIRReference        `json:"asserter,omitempty"`
	Performer          []ProcedurePerformer  `json:"performer,omitempty"`
	Location           *FHIRReference        `json:"location,omitempty"`
	ReasonCode         []FHIRCodeableConcept `json:"reasonCode,omitempty"`
	ReasonReference    []FHIRReference       `json:"reasonReference,omitempty"`
	BodySite           []FHIRCodeableConcept `json:"bodySite,omitempty"`
	Outcome            *FHIRCodeableConcept  `json:"outcome,omitempty"`
	Report             []FHIRReference       `json:"report,omitempty"`
	Complication       []FHIRCodeableConcept `json:"complication,omitempty"`
	ComplicationDetail []FHIRReference       `json:"complicationDetail,omitempty"`
	FollowUp           []FHIRCodeableConcept `json:"followUp,omitempty"`
	Note               []FHIRAnnotation      `json:"note,omitempty"`
	UsedReference      []FHIRReference       `json:"usedReference,omitempty"`
	UsedCode           []FHIRCodeableConcept `json:"usedCode,omitempty"`
}
type ProcedurePerformer struct {
	ID                *string              `json:"id,omitempty"`
	Extension         []FHIRExtension      `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension      `json:"modifierExtension,omitempty"`
	Function          *FHIRCodeableConcept `json:"function,omitempty"`
	Actor             FHIRReference        `json:"actor"`
	OnBehalfOf        *FHIRReference       `json:"onBehalfOf,omitempty"`
}
type PagedFHIRProcedure struct {
	Procedure       []*FHIRProcedure
	HasNextPage     bool
	NextCursor      string
	HasPreviousPage bool
	PreviousCursor  string
	TotalCount      int
}
