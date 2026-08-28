package domain

import "strings"

// Concept models a concept type from OpenConceptLab
type Concept struct {
	ConceptClass     string          `mapstructure:"concept_class" json:"concept_class"`
	DataType         string          `mapstructure:"datatype" json:"datatype"`
	DisplayLocale    string          `mapstructure:"display_locale" json:"display_locale"`
	DisplayName      string          `mapstructure:"display_name" json:"display_name"`
	ExternalID       string          `mapstructure:"external_id" json:"external_id"`
	ID               string          `mapstructure:"id" json:"id"`
	IsLatestVersion  bool            `mapstructure:"is_latest_version" json:"is_latest_version"`
	Locale           *string         `mapstructure:"locale" json:"locale"`
	Owner            string          `mapstructure:"owner" json:"owner"`
	OwnerType        string          `mapstructure:"owner_type" json:"owner_type"`
	OwnerURL         string          `mapstructure:"owner_url" json:"owner_url"`
	Retired          bool            `mapstructure:"retired" json:"retired"`
	Source           string          `mapstructure:"source" json:"source"`
	Type             string          `mapstructure:"type" json:"type"`
	UpdateComment    string          `mapstructure:"update_comment" json:"update_comment"`
	URL              string          `mapstructure:"url" json:"url"`
	UUID             string          `mapstructure:"uuid" json:"uuid"`
	Version          string          `mapstructure:"version" json:"version"`
	VersionCreatedBy string          `mapstructure:"version_created_by" json:"version_created_by"`
	VersionCreatedOn string          `mapstructure:"version_created_on" json:"version_created_on"`
	VersionURL       string          `mapstructure:"version_url" json:"version_url"`
	VersionsURL      string          `mapstructure:"versions_url" json:"versions_url"`
	Descriptions     []*Descriptions `mapstructure:"descriptions" json:"descriptions,omitempty"`
}

type Descriptions struct {
	UUID            string `json:"uuid,omitempty"`
	Description     string `json:"description,omitempty"`
	ExternalID      any    `json:"external_id,omitempty"`
	Type            string `json:"type,omitempty"`
	Locale          string `json:"locale,omitempty"`
	LocalePreferred bool   `json:"locale_preferred,omitempty"`
	DescriptionType string `json:"description_type,omitempty"`
	Checksum        string `json:"checksum,omitempty"`
}

// TerminologySource represents various concept sources
type TerminologySource string

const (
	TerminologySourceICD10    TerminologySource = "ICD10"
	TerminologySourceICD10WHO TerminologySource = "ICD_10_WHO"
	TerminologySourceICD11WHO TerminologySource = "ICD_11_WHO"
	TerminologySourceCIEL     TerminologySource = "CIEL"
	TerminologySourceLOINC    TerminologySource = "LOINC"
	TerminologySourceICHI     TerminologySource = "ICHI"
)

// Hyphenated returns the TerminologySource string with underscores replaced by hyphens
func (t TerminologySource) Hyphenated() string {
	return strings.ReplaceAll(string(t), "_", "-")
}

func (t TerminologySource) Underscore() string {
	return strings.ReplaceAll(string(t), "-", "_")
}

func (d *Concept) GetConceptDisplay() string {
	switch d.Source {
	case TerminologySourceLOINC.Hyphenated():
		if len(d.Descriptions) > 0 {
			for _, d := range d.Descriptions {
				return d.Description
			}
		} else {
			return d.DisplayName
		}

	default:
		return d.DisplayName
	}

	return ""
}

// ConceptPage models the output of ocl concepts with pagination
type ConceptPage struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []*Concept `json:"results"`
}
