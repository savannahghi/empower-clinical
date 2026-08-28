package domain

// DataRequirement is documented here http://hl7.org/fhir/StructureDefinition/DataRequirement
type DataRequirement struct {
	Id                     *string                      `json:"id,omitempty"`
	Extension              []Extension                  `json:"extension,omitempty"`
	SubjectCodeableConcept *FHIRCodeableConcept         `json:"subjectCodeableConcept,omitempty"`
	SubjectReference       *Reference                   `json:"subjectReference,omitempty"`
	MustSupport            []string                     `json:"mustSupport,omitempty"`
	CodeFilter             []DataRequirementCodeFilter  `json:"codeFilter,omitempty"`
	DateFilter             []DataRequirementDateFilter  `json:"dateFilter,omitempty"`
	ValueFilter            []DataRequirementValueFilter `json:"valueFilter,omitempty"`
	Limit                  *int                         `json:"limit,omitempty"`
	Sort                   []DataRequirementSort        `json:"sort,omitempty"`
}

type DataRequirementCodeFilter struct {
	Id          *string      `json:"id,omitempty"`
	Extension   []Extension  `json:"extension,omitempty"`
	Path        *string      `json:"path,omitempty"`
	SearchParam *string      `json:"searchParam,omitempty"`
	ValueSet    *string      `json:"valueSet,omitempty"`
	Code        []FHIRCoding `json:"code,omitempty"`
}

type DataRequirementDateFilter struct {
	Id            *string       `json:"id,omitempty"`
	Extension     []Extension   `json:"extension,omitempty"`
	Path          *string       `json:"path,omitempty"`
	SearchParam   *string       `json:"searchParam,omitempty"`
	ValueDateTime *string       `json:"valueDateTime,omitempty"`
	ValuePeriod   *FHIRPeriod   `json:"valuePeriod,omitempty"`
	ValueDuration *FHIRDuration `json:"valueDuration,omitempty"`
}

type DataRequirementValueFilter struct {
	Id            *string       `json:"id,omitempty"`
	Extension     []Extension   `json:"extension,omitempty"`
	Path          *string       `json:"path,omitempty"`
	SearchParam   *string       `json:"searchParam,omitempty"`
	ValueDateTime *string       `json:"valueDateTime,omitempty"`
	ValuePeriod   *FHIRPeriod   `json:"valuePeriod,omitempty"`
	ValueDuration *FHIRDuration `json:"valueDuration,omitempty"`
}

type DataRequirementSort struct {
	Id        *string     `json:"id,omitempty"`
	Extension []Extension `json:"extension,omitempty"`
	Path      string      `json:"path"`
}
