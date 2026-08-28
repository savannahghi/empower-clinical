package mapper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// TimelineMapper defines the interface for mapping FHIR resources to TimelineResource DTOs.
type TimelineMapper interface {
	ToTimeline(resource interface{}) (*dto.TimelineResource, error)
}

// DefaultTimelineMapper is the default implementation of TimelineMapper.
type DefaultTimelineMapper struct{}

// NewTimelineMapper returns the default TimelineMapper implementation.
//
//nolint:ireturn
func NewTimelineMapper() *DefaultTimelineMapper {
	return &DefaultTimelineMapper{}
}

// ToTimeline maps a generic FHIR resource to a TimelineResource DTO.
func (m *DefaultTimelineMapper) ToTimeline(resource interface{}) (*dto.TimelineResource, error) {
	if resourceMap, ok := resource.(map[string]interface{}); ok {
		return m.mapFromGenericMap(resourceMap)
	}

	switch r := resource.(type) {
	case *domain.FHIRObservation:
		return mapObservationToTimeline(r)
	case *domain.FHIRCondition:
		return mapConditionToTimeline(r)
	case *domain.FHIRRiskAssessment:
		return mapRiskAssessmentToTimeline(r)
	case *domain.FHIRAllergyIntolerance:
		return mapAllergyIntoleranceToTimeline(r)
	case *domain.FHIRMedicationStatement:
		return mapMedicationsToTimeline(r)
	default:
		return nil, fmt.Errorf("ToTimeline not implemented for resource type: %T", resource)
	}
}

// mapFromGenericMap handles resources coming from GetPatientEverything
func (m *DefaultTimelineMapper) mapFromGenericMap(resourceMap map[string]interface{}) (*dto.TimelineResource, error) {
	resourceType, ok := resourceMap["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("resource missing resourceType field")
	}

	resourceBytes, err := json.Marshal(resourceMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	switch resourceType {
	case "Observation":
		var obs domain.FHIRObservation
		if err := json.Unmarshal(resourceBytes, &obs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Observation: %w", err)
		}

		return mapObservationToTimeline(&obs)

	case "Condition":
		var condition domain.FHIRCondition
		if err := json.Unmarshal(resourceBytes, &condition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Condition: %w", err)
		}

		return mapConditionToTimeline(&condition)

	case "RiskAssessment":
		var risk domain.FHIRRiskAssessment
		if err := json.Unmarshal(resourceBytes, &risk); err != nil {
			return nil, fmt.Errorf("failed to unmarshal RiskAssessment: %w", err)
		}

		return mapRiskAssessmentToTimeline(&risk)

	case "AllergyIntolerance":
		var allergy domain.FHIRAllergyIntolerance
		if err := json.Unmarshal(resourceBytes, &allergy); err != nil {
			return nil, fmt.Errorf("failed to unmarshal AllergyIntolerance: %w", err)
		}

		return mapAllergyIntoleranceToTimeline(&allergy)

	case "MedicationStatement":
		var medicationStatement domain.FHIRMedicationStatement
		if err := json.Unmarshal(resourceBytes, &medicationStatement); err != nil {
			return nil, fmt.Errorf("failed to unmarshal MedicationStatement: %w", err)
		}

		return mapMedicationsToTimeline(&medicationStatement)

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// mapObservationToTimeline maps a FHIR Observation to a TimelineResource DTO.
//
//nolint:gocritic
func mapObservationToTimeline(obs *domain.FHIRObservation) (*dto.TimelineResource, error) {
	if obs.ID == nil {
		return nil, fmt.Errorf("observation missing ID")
	}

	var name string
	if obs.Code != nil && obs.Code.Text != "" {
		name = obs.Code.Text
	}

	var value string
	if obs.ValueString != nil {
		value = *obs.ValueString
	}

	var status string
	if obs.Status != nil {
		status = string(*obs.Status)
	}

	var date scalarutils.Date

	var timeRecorded time.Time

	if obs.EffectiveDateTime != nil {
		instant := helpers.ParseDate(string(*obs.EffectiveDateTime))

		d, err := scalarutils.NewDate(instant.Day(), int(instant.Month()), instant.Year())
		if err != nil {
			return nil, fmt.Errorf("failed to create date: %w", err)
		}

		date = *d

		timeRecorded = instant
	} else if obs.EffectiveInstant != nil {
		instant := helpers.ParseDate(string(*obs.EffectiveInstant))

		d, err := scalarutils.NewDate(instant.Day(), int(instant.Month()), instant.Year())
		if err != nil {
			return nil, fmt.Errorf("failed to create date: %w", err)
		}

		date = *d

		timeRecorded = instant
	} else if obs.Issued != nil {
		instant := helpers.ParseDate(string(*obs.Issued))

		d, err := scalarutils.NewDate(instant.Day(), int(instant.Month()), instant.Year())
		if err != nil {
			return nil, fmt.Errorf("failed to create date: %w", err)
		}

		date = *d

		timeRecorded = instant
	}

	var category string
	if len(obs.Category) > 0 && obs.Category[0].Text != "" {
		category = obs.Category[0].Text
	}

	timeline := &dto.TimelineResource{
		ID:           *obs.ID,
		ResourceType: dto.ResourceTypeObservation,
		Name:         name,
		Value:        value,
		Status:       status,
		Date:         date,
		TimeRecorded: timeRecorded,
		Category:     category,
	}

	return timeline, nil
}

// mapConditionToTimeline maps a FHIR Condition to a TimelineResource DTO.
//
//nolint:gocritic
func mapConditionToTimeline(condition *domain.FHIRCondition) (*dto.TimelineResource, error) {
	if condition.ID == nil {
		return nil, fmt.Errorf("condition missing ID")
	}

	var name string
	if condition.Code != nil && condition.Code.Text != "" {
		name = condition.Code.Text
	} else if condition.Code != nil && len(condition.Code.Coding) > 0 && condition.Code.Coding[0].Display != "" {
		name = condition.Code.Coding[0].Display
	} else {
		name = "Condition"
	}

	var value string
	if condition.ClinicalStatus != nil && condition.ClinicalStatus.Text != "" {
		value = condition.ClinicalStatus.Text
	}

	var status string
	if len(condition.Category) > 0 && condition.Category[0].Text != "" {
		status = condition.Category[0].Text
	}

	var date scalarutils.Date

	var timeRecorded time.Time

	if condition.OnsetDateTime != nil {
		date = *condition.OnsetDateTime
		timeRecorded = date.AsTime()
	} else if condition.AbatementDateTime != nil {
		date = *condition.AbatementDateTime
		timeRecorded = date.AsTime()
	} else if condition.RecordedDate != nil {
		date = *condition.RecordedDate
		timeRecorded = date.AsTime()
	}

	timeline := &dto.TimelineResource{
		ID:           *condition.ID,
		ResourceType: dto.ResourceTypeCondition,
		Name:         name,
		Value:        value,
		Status:       status,
		Date:         date,
		TimeRecorded: timeRecorded,
	}

	return timeline, nil
}

// mapRiskAssessmentToTimeline maps a FHIR RiskAssessment to a TimelineResource DTO.
//
//nolint:gocritic
func mapRiskAssessmentToTimeline(risk *domain.FHIRRiskAssessment) (*dto.TimelineResource, error) {
	if risk.ID == nil {
		return nil, fmt.Errorf("risk assessment missing ID")
	}

	var name string
	if risk.Code != nil && risk.Code.Text != "" {
		name = risk.Code.Text
	} else if risk.Code != nil && len(risk.Code.Coding) > 0 && risk.Code.Coding[0].Display != "" {
		name = risk.Code.Coding[0].Display
	} else {
		name = "RiskAssessment"
	}

	value := risk.GetPredictionValue()

	status := string(risk.Status)

	var date scalarutils.Date

	var timeRecorded time.Time

	if risk.OccurrenceDateTime != nil {
		instant := helpers.ParseDate(*risk.OccurrenceDateTime)

		d, err := scalarutils.NewDate(instant.Day(), int(instant.Month()), instant.Year())
		if err != nil {
			return nil, fmt.Errorf("failed to create date: %w", err)
		}

		date = *d
		timeRecorded = instant
	} else if risk.OccurrencePeriod != nil {
		start := risk.OccurrencePeriod.Start
		if start != "" {
			t := start.Time()

			d, err := scalarutils.NewDate(t.Day(), int(t.Month()), t.Year())
			if err != nil {
				return nil, fmt.Errorf("failed to create date: %w", err)
			}

			date = *d
			timeRecorded = t
		}
	}

	category := ""
	if risk.Method != nil && risk.Method.Text != "" {
		category = risk.Method.Text
	}

	timeline := &dto.TimelineResource{
		ID:           *risk.ID,
		ResourceType: dto.ResourceTypeRiskAssessment,
		Name:         name,
		Value:        value,
		Status:       status,
		Date:         date,
		TimeRecorded: timeRecorded,
		Category:     category,
		UsageContext: helpers.ExtractTextFromHTML(string(risk.Text.Div)),
	}

	return timeline, nil
}

// mapAllergyIntoleranceToTimeline maps a FHIR AllergyIntolerance to a TimelineResource DTO.
//
//nolint:gocritic
func mapAllergyIntoleranceToTimeline(allergy *domain.FHIRAllergyIntolerance) (*dto.TimelineResource, error) {
	if allergy.ID == nil {
		return nil, fmt.Errorf("allergy intolerance missing ID")
	}

	var name string
	if allergy.Code != nil && allergy.Code.Text != "" {
		name = allergy.Code.Text
	} else if allergy.Code != nil && len(allergy.Code.Coding) > 0 && allergy.Code.Coding[0].Display != "" {
		name = allergy.Code.Coding[0].Display
	} else {
		name = "AllergyIntolerance"
	}

	var value string
	if len(allergy.Reaction) > 0 && len(allergy.Reaction[0].Manifestation) > 0 && allergy.Reaction[0].Manifestation[0].Concept.Text != "" {
		value = allergy.Reaction[0].Manifestation[0].Concept.Text
	}

	var status string
	if allergy.ClinicalStatus.Text != "" {
		status = allergy.ClinicalStatus.Text
	}

	var date scalarutils.Date

	var timeRecorded time.Time

	if allergy.RecordedDate != nil {
		instant := helpers.ParseDate(string(*allergy.RecordedDate))

		d, err := scalarutils.NewDate(instant.Day(), int(instant.Month()), instant.Year())
		if err != nil {
			return nil, fmt.Errorf("failed to create date: %w", err)
		}

		date = *d
		timeRecorded = instant
	}

	timeline := &dto.TimelineResource{
		ID:           *allergy.ID,
		ResourceType: dto.ResourceTypeAllergyIntolerance,
		Name:         name,
		Value:        value,
		Status:       status,
		Date:         date,
		TimeRecorded: timeRecorded,
	}

	return timeline, nil
}

// mapMedicationsToTimeline maps a FHIR MedicationDispense to a TimelineResource DTO.
//
//nolint:gocritic
func mapMedicationsToTimeline(medication *domain.FHIRMedicationStatement) (*dto.TimelineResource, error) {
	if medication.ID == nil {
		return nil, fmt.Errorf("medication ID missing")
	}

	timeline := &dto.TimelineResource{
		ID:           *medication.ID,
		ResourceType: dto.ResourceTypeMedicationStatement,
	}

	if medication.Medication != nil && medication.Medication.Concept != nil {
		timeline.Name = medication.Medication.Concept.Text
	} else if medication.Medication != nil && medication.Medication.Reference != nil {
		timeline.Name = medication.Medication.Reference.Display
	} else {
		timeline.Name = "Medication"
	}

	if medication.Status != nil {
		timeline.Status = string(*medication.Status)
	}

	if len(medication.Category) > 0 && medication.Category[0].Coding != nil {
		if len(medication.Category[0].Coding) > 0 {
			timeline.Category = medication.Category[0].Coding[0].Display
		}
	} else if medication.Medication != nil && len(medication.Category) > 0 {
		timeline.Category = medication.Category[0].Text
	}

	if medication.DateAsserted != nil {
		instant := helpers.ParseDate(medication.DateAsserted.String())

		date, err := scalarutils.NewDate(instant.Day(), int(instant.Month()), instant.Year())
		if err != nil {
			return nil, fmt.Errorf("failed parse instant to date: %w", err)
		}

		timeline.Date = *date
		timeline.TimeRecorded = instant
	}

	return timeline, nil
}
