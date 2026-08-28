package domain

// Extension is an optional element that provides additional information not captured in the basic resource definition.
// Extensions allow the definition of new data elements or the modification of existing data elements in the FHIR data model.
type Extension struct {
	URL                  string               `json:"url,omitempty"`
	ValueBoolean         *bool                `json:"valueBoolean,omitempty"`
	ValueInteger         *int                 `json:"valueInteger,omitempty"`
	ValueDecimal         *float64             `json:"valueDecimal,omitempty"`
	ValueBase64Binary    string               `json:"valueBase64Binary,omitempty"`
	ValueInstant         string               `json:"valueInstant,omitempty"`
	ValueString          string               `json:"valueString,omitempty"`
	ValueURI             string               `json:"valueURI,omitempty"`
	ValueDate            string               `json:"valueDate,omitempty"`
	ValueDateTime        string               `json:"valueDateTime,omitempty"`
	ValueTime            string               `json:"valueTime,omitempty"`
	ValueCode            string               `json:"valueCode,omitempty"`
	ValueOid             string               `json:"valueOid,omitempty"`
	ValueUUID            string               `json:"valueUUID,omitempty"`
	ValueID              string               `json:"valueID,omitempty"`
	ValueUnsignedInt     int                  `json:"valueUnsignedInt,omitempty"`
	ValuePositiveInt     int                  `json:"valuePositiveInt,omitempty"`
	ValueMarkdown        string               `json:"valueMarkdown,omitempty"`
	ValueAnnotation      *FHIRAnnotation      `json:"valueAnnotation,omitempty"`
	ValueAttachment      *FHIRAttachment      `json:"valueAttachment,omitempty"`
	ValueIdentifier      *FHIRIdentifier      `json:"valueIdentifier,omitempty"`
	ValueCodeableConcept *FHIRCodeableConcept `json:"valueCodeableConcept,omitempty"`
	ValueCoding          *FHIRCoding          `json:"valueCoding,omitempty"`
	ValueQuantity        *FHIRQuantity        `json:"valueQuantity,omitempty"`
	ValueRange           *FHIRRange           `json:"valueRange,omitempty"`
	ValuePeriod          *FHIRPeriod          `json:"valuePeriod,omitempty"`
	ValueRatio           *FHIRRatio           `json:"valueRatio,omitempty"`
	ValueReference       *FHIRReference       `json:"valueReference,omitempty"`
	ValueExpression      *FHIRExpression      `json:"valueExpression,omitempty"`
}

// FHIRExtension contains child elements to represent additional information that is not part of the basic definition of the resource.
type FHIRExtension struct {
	URL       string      `json:"url,omitempty"`
	Extension []Extension `json:"extension,omitempty"`
}
