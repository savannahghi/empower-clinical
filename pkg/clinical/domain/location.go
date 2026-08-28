package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FHIRLocation models the data used to create and/or modify a FHIR location
type FHIRLocation struct {
	ID                   *string                     `json:"id,omitempty"`
	Meta                 *FHIRMeta                   `json:"meta,omitempty"`
	ImplicitRules        *string                     `json:"implicitRules,omitempty"`
	Language             *string                     `json:"language,omitempty"`
	Text                 *FHIRNarrative              `json:"text,omitempty"`
	Extension            []FHIRExtension             `json:"extension,omitempty"`
	ModifierExtension    []FHIRExtension             `json:"modifierExtension,omitempty"`
	Identifier           []FHIRIdentifier            `json:"identifier,omitempty"`
	Status               *FHIRLocationStatus         `json:"status,omitempty"`
	OperationalStatus    *FHIRCoding                 `json:"operationalStatus,omitempty"`
	Name                 *string                     `json:"name,omitempty"`
	Alias                []string                    `json:"alias,omitempty"`
	Description          *string                     `json:"description,omitempty"`
	Mode                 *FHIRLocationMode           `json:"mode,omitempty"`
	Type                 []FHIRCodeableConcept       `json:"type,omitempty"`
	Contact              []FHIRExtendedContactDetail `json:"contact,omitempty"`
	Address              *FHIRAddress                `json:"address,omitempty"`
	Form                 *FHIRCodeableConcept        `json:"form,omitempty"`
	Position             *LocationPosition           `json:"position,omitempty"`
	ManagingOrganization *Reference                  `json:"managingOrganization,omitempty"`
	PartOf               *Reference                  `json:"partOf,omitempty"`
	Characteristic       []FHIRCodeableConcept       `json:"characteristic,omitempty"`
	HoursOfOperation     []FHIRAvailability          `json:"hoursOfOperation,omitempty"`
	VirtualService       []FHIRVirtualServiceDetail  `json:"virtualService,omitempty"`
	Endpoint             []Reference                 `json:"endpoint,omitempty"`
}

type LocationPosition struct {
	ID                *string     `json:"id,omitempty"`
	Extension         []Extension `json:"extension,omitempty"`
	ModifierExtension []Extension `json:"modifierExtension,omitempty"`
	Longitude         *float64    `json:"longitude,omitempty"`
	Latitude          *float64    `json:"latitude,omitempty"`
	Altitude          *float64    `json:"altitude,omitempty"`
}

type FHIRAvailability struct {
	ID               *string                         `json:"id,omitempty"`
	Extension        []FHIRExtension                 `json:"extension,omitempty"`
	AvailableTime    []FHIRAvailabilityAvailableTime `json:"availableTime,omitempty"`
	NotAvailableTime []AvailabilityNotAvailableTime  `json:"notAvailableTime,omitempty"`
}

type FHIRAvailabilityAvailableTime struct {
	ID                 *string         `json:"id,omitempty"`
	Extension          []FHIRExtension `json:"extension,omitempty"`
	DaysOfWeek         []DaysOfWeek    `json:"daysOfWeek,omitempty"`
	AllDay             *bool           `json:"allDay,omitempty"`
	AvailableStartTime *string         `json:"availableStartTime,omitempty"`
	AvailableEndTime   *string         `json:"availableEndTime,omitempty"`
}

type AvailabilityNotAvailableTime struct {
	ID          *string         `json:"id,omitempty"`
	Extension   []FHIRExtension `json:"extension,omitempty"`
	Description *string         `json:"description,omitempty"`
	During      *FHIRPeriod     `json:"during,omitempty"`
}

type DaysOfWeek int

const (
	DaysOfWeekMon DaysOfWeek = iota
	DaysOfWeekTue
	DaysOfWeekWed
	DaysOfWeekThu
	DaysOfWeekFri
	DaysOfWeekSat
	DaysOfWeekSun
)

func (code DaysOfWeek) Code() string {
	switch code {
	case DaysOfWeekMon:
		return "mon"
	case DaysOfWeekTue:
		return "tue"
	case DaysOfWeekWed:
		return "wed"
	case DaysOfWeekThu:
		return "thu"
	case DaysOfWeekFri:
		return "fri"
	case DaysOfWeekSat:
		return "sat"
	case DaysOfWeekSun:
		return "sun"
	}

	return "<unknown>"
}

// FHIRLocationMode is documented here http://hl7.org/fhir/ValueSet/location-mode
type FHIRLocationMode string

const (
	FHIRLocationModeInstance FHIRLocationMode = "instance"
	FHIRLocationModeKind     FHIRLocationMode = "kind"
)

func (code FHIRLocationMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(code.Code())
}
func (code *FHIRLocationMode) UnmarshalJSON(json []byte) error {
	s := strings.Trim(string(json), "\"")
	switch s {
	case "instance":
		*code = FHIRLocationModeInstance
	case "kind":
		*code = FHIRLocationModeKind
	default:
		return fmt.Errorf("unknown FHIRLocationMode code `%s`", s)
	}

	return nil
}
func (code FHIRLocationMode) String() string {
	return code.Code()
}
func (code FHIRLocationMode) Code() string {
	switch code {
	case FHIRLocationModeInstance:
		return "instance"
	case FHIRLocationModeKind:
		return "kind"
	}

	return "<unknown>"
}
func (code FHIRLocationMode) Display() string {
	switch code {
	case FHIRLocationModeInstance:
		return "Instance"
	case FHIRLocationModeKind:
		return "Kind"
	}

	return "<unknown>"
}
func (code FHIRLocationMode) Definition() string {
	switch code {
	case FHIRLocationModeInstance:
		return "The Location resource represents a specific instance of a location (e.g. Operating Theatre 1A)."
	case FHIRLocationModeKind:
		return "The Location represents a class of locations (e.g. Any Operating Theatre) although this class of locations could be constrained within a specific boundary (such as organization, or parent location, address etc.)."
	}

	return "<unknown>"
}

// LocationStatus is documented here http://hl7.org/fhir/ValueSet/location-status
type FHIRLocationStatus string

const (
	FHIRLocationStatusActive    FHIRLocationStatus = "active"
	FHIRLocationStatusSuspended FHIRLocationStatus = "suspended"
	FHIRLocationStatusInactive  FHIRLocationStatus = "inactive"
)

func (code FHIRLocationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(code.Code())
}

func (code *FHIRLocationStatus) UnmarshalJSON(json []byte) error {
	s := strings.Trim(string(json), "\"")
	switch s {
	case "active":
		*code = FHIRLocationStatusActive
	case "suspended":
		*code = FHIRLocationStatusSuspended
	case "inactive":
		*code = FHIRLocationStatusInactive
	default:
		return fmt.Errorf("unknown FHIRLocationStatus code `%s`", s)
	}

	return nil
}

func (code FHIRLocationStatus) String() string {
	return code.Code()
}

func (code FHIRLocationStatus) Code() string {
	switch code {
	case FHIRLocationStatusActive:
		return "active"
	case FHIRLocationStatusSuspended:
		return "suspended"
	case FHIRLocationStatusInactive:
		return "inactive"
	}

	return "<unknown>"
}

func (code FHIRLocationStatus) Display() string {
	switch code {
	case FHIRLocationStatusActive:
		return "Active"
	case FHIRLocationStatusSuspended:
		return "Suspended"
	case FHIRLocationStatusInactive:
		return "Inactive"
	}

	return "<unknown>"
}

func (code FHIRLocationStatus) Definition() string {
	switch code {
	case FHIRLocationStatusActive:
		return "The location is operational."
	case FHIRLocationStatusSuspended:
		return "The location is temporarily closed."
	case FHIRLocationStatusInactive:
		return "The location is no longer used."
	}

	return "<unknown>"
}
