package dto

import "github.com/savannahghi/empower-clinical/pkg/clinical/domain"

type CarePlanInput struct {
	EncounterID      string `json:"encounterID" validate:"required"`
	PlanDefinitionID string `json:"planDefinitionID" validate:"required"`
	Notes            string `json:"notes"`
}

type CarePlanOutput struct {
	ID              string         `json:"id,omitempty"`
	Title           string         `json:"title,omitempty"`
	Description     string         `json:"description,omitempty"`
	EncounterID     string         `json:"encounterId,omitempty"`
	Patient         Patient        `json:"patient,omitempty"`
	TreatmentPhases []*ChemoPhases `json:"treatmentPhases,omitempty"`
}

type ChemoPhases struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Cycles      []*Cycles `json:"cycles,omitempty"`
	Status      string    `json:"status,omitempty"`
}

type Cycles struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type CarePlanPayload struct {
	Data       CarePlanInput
	Tags       []domain.FHIRCodingInput
	FacilityID string
}
