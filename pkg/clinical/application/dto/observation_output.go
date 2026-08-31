package dto

import "github.com/savannahghi/serverutils"

// Observation is a minimal representation of a fhir Observation
type Observation struct {
	ID               string            `json:"id,omitempty"`
	Name             string            `json:"name,omitempty"`
	Value            string            `json:"value,omitempty"`
	Status           ObservationStatus `json:"status,omitempty"`
	Category         string            `json:"category,omitempty"`
	Code             string            `json:"code,omitempty"`
	PatientID        string            `json:"patientID,omitempty"`
	PatientName      string            `json:"patientName,omitempty"`
	EncounterID      string            `json:"encounterID,omitempty"`
	TimeRecorded     string            `json:"timeRecorded,omitempty"`
	Interpretation   []string          `json:"interpretation,omitempty"`
	Note             string            `json:"note,omitempty"`
	UsageContext     ScreeningTypeEnum `json:"usageContext,omitempty"`
	ServiceRequestID string            `json:"serviceRequestID,omitempty"`
}

// NodeID implements the NodeID() method that returns ID of a Node
func (o Observation) NodeID() string {
	return o.ID
}

// ConvertObservationEdges converts the generic edges to a observation edges (strict go types ⛓️ )
func ConvertObservationEdges(edges []serverutils.Edge[Observation]) []ObservationEdge {
	result := make([]ObservationEdge, len(edges))
	for i, edge := range edges {
		result[i] = ObservationEdge{
			Cursor: edge.Cursor,
			Node:   edge.Node,
		}
	}

	return result
}

// ObservationEdge is an observation edge
type ObservationEdge struct {
	Node   Observation `json:"node"`
	Cursor string      `json:"cursor"`
}

// ObservationConnection  is an Observation Connection Type
type ObservationConnection struct {
	SearchID   string               `json:"searchID"`
	TotalCount int                  `json:"totalCount"`
	Edges      []ObservationEdge    `json:"edges"`
	PageInfo   serverutils.PageInfo `json:"pageInfo"`
}

// CreateObservationConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateObservationConnection(observations []*Observation, pageInfo serverutils.PageInfo, total int) ObservationConnection {
	connection := ObservationConnection{
		TotalCount: total,
		Edges:      []ObservationEdge{},
		PageInfo:   pageInfo,
	}

	for _, observation := range observations {
		edge := ObservationEdge{
			Node:   *observation,
			Cursor: observation.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
