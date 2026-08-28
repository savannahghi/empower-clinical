package domain

import "github.com/savannahghi/scalarutils"

// FHIRDiagnosticReport is documented here http://hl7.org/fhir/StructureDefinition/DiagnosticReport
type FHIRDiagnosticReport struct {
	ID                 *string                          `json:"id,omitempty"`
	Meta               *FHIRMeta                        `json:"meta,omitempty"`
	ImplicitRules      *string                          `json:"implicitRules,omitempty"`
	Language           *string                          `json:"language,omitempty"`
	Text               *FHIRNarrative                   `json:"text,omitempty"`
	Extension          []*Extension                     `json:"extension,omitempty"`
	ModifierExtension  []*Extension                     `json:"modifierExtension,omitempty"`
	Identifier         []*FHIRIdentifier                `json:"identifier,omitempty"`
	BasedOn            []*FHIRReference                 `json:"basedOn,omitempty"`
	Status             DiagnosticReportStatusEnum       `json:"status"`
	Category           []*FHIRCodeableConcept           `json:"category,omitempty"`
	Code               FHIRCodeableConcept              `json:"code"`
	Subject            *FHIRReference                   `json:"subject,omitempty"`
	Encounter          *FHIRReference                   `json:"encounter,omitempty"`
	EffectiveDateTime  *scalarutils.DateTime            `json:"effectiveDateTime,omitempty"`
	EffectivePeriod    *FHIRPeriod                      `json:"effectivePeriod,omitempty"`
	Issued             *string                          `json:"issued,omitempty"`
	Performer          []*FHIRReference                 `json:"performer,omitempty"`
	ResultsInterpreter []*FHIRReference                 `json:"resultsInterpreter,omitempty"`
	Specimen           []*FHIRReference                 `json:"specimen,omitempty"`
	Result             []*FHIRReference                 `json:"result,omitempty"`
	ImagingStudy       []*FHIRReference                 `json:"imagingStudy,omitempty"`
	Media              []*FHIRDiagnosticReportMedia     `json:"media,omitempty"`
	Conclusion         *string                          `json:"conclusion,omitempty"`
	ConclusionCode     []*FHIRCodeableConcept           `json:"conclusionCode,omitempty"`
	PresentedForm      []*FHIRAttachment                `json:"presentedForm,omitempty"`
	SupportingInfo     []DiagnosticReportSupportingInfo `json:"supportingInfo,omitempty"`
}

// FHIRDiagnosticReportMedia represents the key images associated with this report
type FHIRDiagnosticReportMedia struct {
	ID                *string        `json:"id,omitempty"`
	Extension         []*Extension   `json:"extension,omitempty"`
	ModifierExtension []*Extension   `json:"modifierExtension,omitempty"`
	Comment           *string        `json:"comment,omitempty"`
	Link              *FHIRReference `json:"link"`
}

// FHIRDiagnosticReportInput is documented here http://hl7.org/fhir/StructureDefinition/DiagnosticReport
type FHIRDiagnosticReportInput struct {
	ID                 *string                           `json:"id,omitempty"`
	Meta               *FHIRMetaInput                    `json:"meta,omitempty"`
	ImplicitRules      *string                           `json:"implicitRules,omitempty"`
	Language           *string                           `json:"language,omitempty"`
	Text               *FHIRNarrative                    `json:"text,omitempty"`
	Extension          []*FHIRExtension                  `json:"extension,omitempty"`
	ModifierExtension  []*FHIRExtension                  `json:"modifierExtension,omitempty"`
	Identifier         []*FHIRIdentifierInput            `json:"identifier,omitempty"`
	BasedOn            []*FHIRReferenceInput             `json:"basedOn,omitempty"`
	Status             DiagnosticReportStatusEnum        `json:"status"`
	Category           []*FHIRCodeableConceptInput       `json:"category,omitempty"`
	Code               FHIRCodeableConceptInput          `json:"code"`
	Subject            *FHIRReferenceInput               `json:"subject,omitempty"`
	Encounter          *FHIRReferenceInput               `json:"encounter,omitempty"`
	EffectiveDateTime  *scalarutils.DateTime             `json:"effectiveDateTime,omitempty"`
	EffectivePeriod    *FHIRPeriodInput                  `json:"effectivePeriod,omitempty"`
	Issued             *string                           `json:"issued,omitempty"`
	Performer          []*FHIRReferenceInput             `json:"performer,omitempty"`
	ResultsInterpreter []*FHIRReferenceInput             `json:"resultsInterpreter,omitempty"`
	Specimen           []*FHIRReferenceInput             `json:"specimen,omitempty"`
	Result             []*FHIRReferenceInput             `json:"result,omitempty"`
	ImagingStudy       []*FHIRReferenceInput             `json:"imagingStudy,omitempty"`
	Media              []*FHIRDiagnosticReportMediaInput `json:"media,omitempty"`
	Conclusion         *string                           `json:"conclusion,omitempty"`
	ConclusionCode     []*FHIRCodeableConceptInput       `json:"conclusionCode,omitempty"`
	PresentedForm      []*FHIRAttachmentInput            `json:"presentedForm,omitempty"`
	SupportingInfo     []DiagnosticReportSupportingInfo  `json:"supportingInfo,omitempty"`
}

type DiagnosticReportSupportingInfo struct {
	ID                *string             `json:"id,omitempty"`
	Extension         []Extension         `json:"extension,omitempty"`
	ModifierExtension []Extension         `json:"modifierExtension,omitempty"`
	Type              FHIRCodeableConcept `json:"type"`
	Reference         FHIRReferenceInput  `json:"reference"`
}

// FHIRDiagnosticReportMediaInput represents the key images associated with this report
type FHIRDiagnosticReportMediaInput struct {
	ID                *string             `json:"id,omitempty"`
	Extension         []*FHIRExtension    `json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension    `json:"modifierExtension,omitempty"`
	Comment           *string             `json:"comment,omitempty"`
	Link              *FHIRReferenceInput `json:"link"`
}

// PagedFHIRDiagnosticReport is used to return paginated list of risk diagnostic report
type PagedFHIRDiagnosticReport struct {
	DiagnosticReport []FHIRDiagnosticReport `mapstructure:"diagnosticReport"`
	HasNextPage      bool                   `mapstructure:"hasNextPage"`
	NextCursor       string                 `mapstructure:"nextCursor"`
	HasPreviousPage  bool                   `mapstructure:"hasPreviousPage"`
	PreviousCursor   string                 `mapstructure:"previousCursor"`
	TotalCount       int                    `mapstructure:"totalCount"`
}
