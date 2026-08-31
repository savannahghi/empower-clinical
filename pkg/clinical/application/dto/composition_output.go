package dto

import (
	"github.com/savannahghi/scalarutils"
)

// Composition is a minimal representation of a fhir Composition
type Composition struct {
	ID          string                `json:"id,omitempty"`
	Text        string                `json:"text,omitempty"`
	Type        CompositionType       `json:"type,omitempty"`
	Category    CompositionCategory   `json:"category,omitempty"`
	Status      CompositionStatusEnum `json:"status,omitempty"`
	PatientID   string                `json:"patientID,omitempty"`
	EncounterID string                `json:"encounterID,omitempty"`
	Date        *scalarutils.Date     `json:"date" swaggertype:"primitive,string"`
	Author      string                `json:"author,omitempty"`
	Section     []*Section            `json:"section"`
}

// CompositionEdge is a composition edge
type CompositionEdge struct {
	Node   Composition
	Cursor string
}

// CompositionConnection  is a Composition Connection Type
type CompositionConnection struct {
	TotalCount int
	Edges      []CompositionEdge
	PageInfo   PageInfo
}

// CreateCompositionConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateCompositionConnection(compositions []Composition, pageInfo PageInfo, total int) CompositionConnection {
	connection := CompositionConnection{
		TotalCount: total,
		Edges:      []CompositionEdge{},
		PageInfo:   pageInfo,
	}

	for _, composition := range compositions {
		edge := CompositionEdge{
			Node:   composition,
			Cursor: composition.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
