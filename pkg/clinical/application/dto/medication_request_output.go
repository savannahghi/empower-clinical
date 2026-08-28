package dto

import (
	"errors"
	"strconv"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

type MedicationRequestOutput struct {
	ID                 string                         `json:"id,omitempty"`
	EncounterID        string                         `json:"encounterID,omitempty"`
	Status             domain.MedicationRequestStatus `json:"status,omitempty"`
	Notes              []*Annotation                  `json:"notes,omitempty"`
	Subject            *Reference                     `json:"subject,omitempty"`
	DosageInstructions []DosageInstruction            `json:"dosageInstructions,omitempty"`
	Medication         string                         `json:"medication,omitempty"`
	AuthoredOn         *scalarutils.DateTime          `json:"authoredOn,omitempty" swaggertype:"primitive,string"`
	Diagnosis          string                         `json:"diagnosis,omitempty"`
	Priority           string                         `json:"priority,omitempty"`
	OrderedBy          string                         `json:"orderedBy,omitempty"`
	FacilityName       string                         `json:"facilityName,omitempty"`
}

type MedicationRequestEdge struct {
	Node   MedicationRequestOutput `json:"node,omitempty"`
	Cursor string                  `json:"cursor,omitempty"`
}

type MedicationRequestConnection struct {
	TotalCount int                     `json:"totalCount"`
	Edges      []MedicationRequestEdge `json:"edges,omitempty"`
	PageInfo   PageInfo                `json:"pageInfo,omitempty"`
}

// DosageInstruction models a dosage instruction for a medication
type DosageInstruction struct {
	Route                 ValueSetData            `json:"route,omitempty"`
	DoseQuantity          float64                 `json:"doseQuantity,omitempty"`
	DoseUnit              string                  `json:"doseUnit,omitempty"`
	Period                string                  `json:"period,omitempty"`
	PeriodUnit            *domain.UnitsOfTimeEnum `json:"periodUnit,omitempty"`
	Frequency             int                     `json:"frequency,omitempty"`
	Duration              string                  `json:"duration,omitempty"`
	DurationUnit          *domain.UnitsOfTimeEnum `json:"durationUnit,omitempty"`
	StartDate             *scalarutils.DateTime   `json:"startDate,omitempty" swaggertype:"primitive,string"`
	EndDate               *scalarutils.DateTime   `json:"endDate,omitempty" swaggertype:"primitive,string"`
	Condition             string                  `json:"condition,omitempty"`
	PatientInstruction    string                  `json:"patientInstruction,omitempty"`
	AdditionalInstruction []string                `json:"additionalInstruction,omitempty"`
	AsNeeded              bool                    `json:"asNeeded,omitempty"`

	// Free text dosage instructions can be used for cases where the instructions are too complex to code.
	FreeTextInstruction *string `json:"freeTextInstruction,omitempty"`
}

type ValueSetData struct {
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

// Validate is used to validate the dosage instruction inputs based on certain constraints
func (d *DosageInstruction) Validate() error {
	// tim-1: if there's a duration, there needs to be duration units
	if d.Duration != "" && d.DurationUnit == nil {
		return errors.New("duration is set but durationUnit is missing")
	}

	// tim-2: if there's a period, there needs to be period units
	if d.Period != "" && d.PeriodUnit == nil {
		return errors.New("period is set but periodUnit is missing")
	}

	// tim-4: duration SHALL be a non-negative value
	if d.Duration != "" {
		durationValue, err := strconv.ParseFloat(d.Duration, 64)
		if err != nil || durationValue < 0 {
			return errors.New("duration must be a non-negative number")
		}
	}

	// tim-5: period SHALL be a non-negative value
	if d.Period != "" {
		periodValue, err := strconv.ParseFloat(d.Period, 64)
		if err != nil || periodValue < 0 {
			return errors.New("period must be a non-negative number")
		}
	}

	return nil
}

// CreateMedicationRequestConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateMedicationRequestConnection(medicationRequests []*MedicationRequestOutput, pageInfo PageInfo, total int) MedicationRequestConnection {
	connection := MedicationRequestConnection{
		TotalCount: total,
		Edges:      []MedicationRequestEdge{},
		PageInfo:   pageInfo,
	}

	for _, medicationRequest := range medicationRequests {
		edge := MedicationRequestEdge{
			Node:   *medicationRequest,
			Cursor: medicationRequest.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
