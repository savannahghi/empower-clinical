package dto

import (
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// Allergy represents an allergy containing minimal FHIR resources
type Allergy struct {
	ID                string                   `json:"ID"`
	PatientID         string                   `json:"patientID"`
	Code              string                   `json:"code"`
	Name              string                   `json:"name"`
	System            string                   `json:"system"`
	TerminologySource domain.TerminologySource `json:"terminologySource" swaggertype:"primitive,string"`
	OnsetDateTime     scalarutils.DateTime     `json:"onsetDateTime,omitempty" swaggertype:"primitive,string"`
	EncounterID       string                   `json:"encounterID"`
	Reaction          Reaction                 `json:"reaction"`
}

// Reaction represents a reaction containing minimal FHIR resources
type Reaction struct {
	Code     string                                 `json:"code"`
	Name     string                                 `json:"name"`
	System   string                                 `json:"system"`
	Severity AllergyIntoleranceReactionSeverityEnum `json:"severity" swaggertype:"primitive,string"`
}

// AllergyEdge is an allergy edge
type AllergyEdge struct {
	Node   Allergy
	Cursor string
}

// AllergyConnection  is an Allergy Connection Type
type AllergyConnection struct {
	TotalCount int
	Edges      []AllergyEdge
	PageInfo   PageInfo
}

// CreateAllergyConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateAllergyConnection(allergies []*Allergy, pageInfo PageInfo, total int) AllergyConnection {
	connection := AllergyConnection{
		TotalCount: total,
		Edges:      []AllergyEdge{},
		PageInfo:   pageInfo,
	}

	for _, allergy := range allergies {
		edge := AllergyEdge{
			Node:   *allergy,
			Cursor: allergy.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
