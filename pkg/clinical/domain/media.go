package domain

import "time"

// FHIRMedia is the domain representation of FHIR media records
type FHIRMedia struct {
	ID         *string              `json:"id,omitempty"`
	Identifier []*FHIRIdentifier    `json:"identifier,omitempty"`
	BaseOn     *FHIRReference       `json:"basedOn,omitempty"`
	PartOf     *FHIRReference       `json:"partOf,omitempty"`
	Status     MediaStatusEnum      `json:"status,omitempty"`
	Type       *FHIRCodeableConcept `json:"type,omitempty"`
	Modality   *FHIRCodeableConcept `json:"modality,omitempty"`
	View       *FHIRCodeableConcept `json:"view,omitempty"`
	Subject    *FHIRReferenceInput  `json:"subject,omitempty"`
	Encounter  *FHIRReferenceInput  `json:"encounter,omitempty"`
	Issued     *time.Time           `json:"issued,omitempty"`
	Operator   *FHIRReferenceInput  `json:"operator,omitempty"`
	ReasonCode *FHIRCodeableConcept `json:"reasonCode,omitempty"`
	BodySite   *FHIRCodeableConcept `json:"bodySite,omitempty"`
	DeviceName string               `json:"deviceName,omitempty"`
	Device     *FHIRCodeableConcept `json:"device,omitempty"`
	Height     int64                `json:"height,omitempty"`
	Width      int64                `json:"width,omitempty"`
	Frames     int64                `json:"frames,omitempty"`
	Duration   int64                `json:"duration,omitempty"`
	Content    *FHIRAttachmentInput `json:"content,omitempty"`
	Meta       *FHIRMetaInput       `json:"meta,omitempty"`
}

// PagedFHIRMedia is a media's paginaton dataclass
type PagedFHIRMedia struct {
	Media           []FHIRMedia
	HasNextPage     bool
	NextCursor      string
	HasPreviousPage bool
	PreviousCursor  string
	TotalCount      int
}
