package domain

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

type FHIRPractitioner struct {
	ID                *string                         `json:"id,omitempty"`
	Meta              *FHIRMeta                       `json:"meta,omitempty"`
	ImplicitRules     *string                         `json:"implicitRules,omitempty"`
	Language          *string                         `json:"language,omitempty"`
	Text              *FHIRNarrative                  `json:"text,omitempty"`
	Extension         []FHIRExtension                 `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension                 `json:"modifierExtension,omitempty"`
	Identifier        []*FHIRIdentifier               `json:"identifier,omitempty"`
	Active            *bool                           `json:"active,omitempty"`
	Name              []*FHIRHumanName                `json:"name,omitempty"`
	Telecom           []*FHIRContactPoint             `json:"telecom,omitempty"`
	Address           []FHIRAddress                   `json:"address,omitempty"`
	Gender            PatientGenderEnum               `json:"gender,omitempty"`
	BirthDate         *scalarutils.Date               `json:"birthDate,omitempty"`
	Photo             []FHIRAttachment                `json:"photo,omitempty"`
	Qualification     []FHIRPractitionerQualification `json:"qualification,omitempty"`
}

type FHIRPractitionerQualification struct {
	ID         *string             `json:"id,omitempty"`
	Identifier []FHIRIdentifier    `json:"identifier,omitempty"`
	Code       FHIRCodeableConcept `json:"code"`
	Period     *FHIRPeriod         `json:"period,omitempty"`
	Issuer     *FHIRReference      `json:"issuer,omitempty"`
}

type PractitionerConnection struct {
	Edges    []*PatientEdge          `json:"edges"`
	PageInfo *firebasetools.PageInfo `json:"pageInfo"`
}
