package dto

import "github.com/savannahghi/empower-clinical/pkg/clinical/domain"

type PlanDefinitionInput struct {
	Title       string       `json:"title,omitempty" validate:"required"`
	Description string       `json:"description,omitempty" validate:"required"`
	Action      []PlanAction `json:"action,omitempty" validate:"required"`
}

type PlanAction struct {
	Title        string           `json:"title,omitempty"`
	Description  string           `json:"description,omitempty"`
	TimingTiming *Timing          `json:"timingTiming,omitempty"`
	Medications  []PlanMedication `json:"medications,omitempty"` // For medication
	Action       []PlanAction     `json:"action,omitempty"`      // for nested actions (e.g., cycles)
}

type PlanMedication struct {
	MedicationID string                    `json:"medicationID,omitempty"`
	Dosage       DosageAdministrationInput `json:"dosage,omitempty"`
	DoseQuantity float64                   `json:"doseQuantity,omitempty"`
	DoseUnit     string                    `json:"doseUnit,omitempty"`
}

type DosageAdministrationInput struct {
	Route                      Coding `json:"route,omitempty"`
	Method                     Coding `json:"method,omitempty"`
	AdministrationInstructions string `json:"administrationInstructions,omitempty"`
}

type Timing struct {
	Repeat *Repeat `json:"repeat,omitempty"`
}

type Repeat struct {
	Frequency  int    `json:"frequency,omitempty"`
	Period     int    `json:"period,omitempty"`
	PeriodUnit string `json:"periodUnit,omitempty"` // e.g., "wk"
	Count      int    `json:"count,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

type PlanDefinitionOutputConnection struct {
	TotalCount int                  `json:"totalCount"`
	Edges      []PlanDefinitionEdge `json:"edges,omitempty"`
	PageInfo   PageInfo             `json:"pageInfo,omitempty"`
}

type PlanDefinitionEdge struct {
	Node   domain.FHIRPlanDefinition `json:"node,omitempty"`
	Cursor string                    `json:"cursor,omitempty"`
}

func CreatePlanDefinitionConnection(planDefinitions []*domain.FHIRPlanDefinition, pageInfo PageInfo, total int) PlanDefinitionOutputConnection {
	connection := PlanDefinitionOutputConnection{
		TotalCount: total,
		Edges:      []PlanDefinitionEdge{},
		PageInfo:   pageInfo,
	}

	for _, pd := range planDefinitions {
		edge := PlanDefinitionEdge{
			Node: *pd,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
