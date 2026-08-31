package dto

import "github.com/savannahghi/scalarutils"

// Condition represents a FHIR condition
type Condition struct {
	ID     string          `json:"id" example:"10414684-7ec5-4ad2-9aa6-31e35c51e5bb"`
	Status ConditionStatus `json:"status" example:"ACTIVE"`
	Name   string          `json:"condition" example:"Mastocytosis"`
	Code   string          `json:"code" example:"2A21"`
	System string          `json:"system" example:"http://id.who.int/icd/release/11/mms"`

	Category ConditionCategory `json:"category" example:""`

	OnsetDate    *scalarutils.Date `json:"onsetDate" swaggertype:"primitive,string" example:"2025-03-12"`
	RecordedDate *scalarutils.Date `json:"recordedDate" swaggertype:"primitive,string" example:"2025-04-03"`

	Note string `json:"note" example:"An example note"`

	PatientID         string                   `json:"patientID" example:"5f4a922d-abcd-405a-11cc-a88451195bfe"`
	EncounterID       string                   `json:"encounterID" example:"640bd94d-51ad-4eee-b496-bb0fb4b87d38"`
	OncologyCondition *OncologyConditionOutput `json:"oncologyCondition,omitempty"`
	TreatmentLinkage  *TreatmentLinkageOutput  `json:"treatmentLinkage,omitempty"`
}

type OncologyConditionOutput struct {
	ICDO3PrimaryTumorCode string `json:"ICDO3PrimaryTumorCode,omitempty"`
	ICDO3MorphologyCode   string `json:"ICDO3MorphologyCode,omitempty"`
	Stage                 string `json:"stage,omitempty"`
}

// TreatmentLinkageOutput captures the retrospective diagnosis-to-treatment linkage recorded
// against a condition. It is only populated for conditions that carry the linkage information.
type TreatmentLinkageOutput struct {
	LinkedToTreatment bool              `json:"linkedToTreatment,omitempty"`
	TreatmentFacility string            `json:"treatmentFacility"`
	TreatmentProgram  string            `json:"treatmentProgram"`
	EnrollmentDate    *scalarutils.Date `json:"enrollmentDate" swaggertype:"primitive,string"`
}

// ConditionEdge is a condition edge
type ConditionEdge struct {
	Node   Condition
	Cursor string
}

// ConditionConnection  is a Condition Connection Type
type ConditionConnection struct {
	TotalCount int
	Edges      []ConditionEdge
	PageInfo   PageInfo
}

// CreateConditionConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateConditionConnection(conditions []Condition, pageInfo PageInfo, total int) ConditionConnection {
	connection := ConditionConnection{
		TotalCount: total,
		Edges:      []ConditionEdge{},
		PageInfo:   pageInfo,
	}

	for _, condition := range conditions {
		edge := ConditionEdge{
			Node:   condition,
			Cursor: condition.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
