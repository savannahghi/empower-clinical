package domain

import (
	"time"

	"github.com/savannahghi/scalarutils"
)

// FHIRMetaInput is the input to resource meta
type FHIRMetaInput struct {
	VersionID string            `json:"versionId,omitempty"`
	Source    string            `json:"source,omitempty"`
	Tag       []FHIRCodingInput `json:"tag,omitempty"`
	Security  []FHIRCodingInput `json:"security,omitempty"`
	Profile   []scalarutils.URI `json:"profile,omitempty"`
}

func (c *FHIRMetaInput) GetOrganization() string {
	facilitySystem := scalarutils.URI("http://mycarehub/tenant-identification/facility")

	if len(c.Tag) > 0 {
		for _, tag := range c.Tag {
			if string(*tag.System) == string(facilitySystem) {
				return string(tag.Code)
			}
		}
	}

	return ""
}

// FHIRMeta is a set of metadata that provides technical and workflow context to a resource.
type FHIRMeta struct {
	VersionID   string       `json:"versionId,omitempty"`
	LastUpdated time.Time    `json:"lastUpdated,omitempty"`
	Source      string       `json:"source,omitempty"`
	Tag         []FHIRCoding `json:"tag,omitempty"`
	Security    []FHIRCoding `json:"security,omitempty"`
}

// GetFacilityName is used to get the name of the facility
func (f *FHIRMeta) GetFacilityName() string {
	facilitySystem := scalarutils.URI("http://mycarehub/tenant-identification/facility")

	if f.Tag != nil {
		for _, meta := range f.Tag {
			if string(*meta.System) == string(facilitySystem) {
				return meta.Display
			}
		}
	}

	return ""
}
