package dto

import "time"

// SegmentationPayload is used to stratify clients in advantage EMR.
type SegmentationPayload struct {
	// ClinicalID represents the patient's ID in this service
	ClinicalID   string               `json:"clinical_id,omitempty"`
	SegmentLabel SegmentationCategory `json:"segment_label,omitempty"`
}

// SMSPayload is used to model data used to send sms to patients through the advantage service
type SMSPayload struct {
	Intention  string   `json:"intention"`
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
}

// AdvantageResponse is the generalized model for all the responses from advantage emr.
type AdvantageResponse struct {
	Count       *int  `json:"count"`
	Next        any   `json:"next"`
	Previous    any   `json:"previous"`
	PageSize    *int  `json:"page_size"`
	CurrentPage *int  `json:"current_page"`
	TotalPages  *int  `json:"total_pages"`
	StartIndex  *int  `json:"start_index"`
	EndIndex    *int  `json:"end_index"`
	Results     []any `json:"results"`
}

type PersonContacts struct {
	ID                    string `json:"id"`
	CreatedByName         string `json:"created_by_name"`
	UpdatedByName         string `json:"updated_by_name"`
	Active                bool   `json:"active"`
	ContactType           string `json:"contact_type"`
	Contact               string `json:"contact"`
	Verified              bool   `json:"verified"`
	ConsentToContactGiven bool   `json:"consent_to_contact_given"`
	IsPrimaryContact      bool   `json:"is_primary_contact"`
	Person                string `json:"person"`
}

type Person struct {
	ID               string           `json:"id"`
	PersonDisplay    string           `json:"person_display"`
	PersonContacts   []PersonContacts `json:"person_contacts"`
	PersonIds        []any            `json:"person_ids"`
	PhoneNumber      string           `json:"phone_number"`
	Email            any              `json:"email"`
	Age              any              `json:"age"`
	GlobalHealthID   any              `json:"global_health_id"`
	Active           bool             `json:"active"`
	FirstName        string           `json:"first_name"`
	LastName         string           `json:"last_name"`
	OtherNames       any              `json:"other_names"`
	DateOfBirth      any              `json:"date_of_birth"`
	Title            string           `json:"title"`
	Gender           any              `json:"gender"`
	Deceased         bool             `json:"deceased"`
	Language         any              `json:"language"`
	Metadata         any              `json:"metadata"`
	AssociatedRegion any              `json:"associated_region"`
}

type PractitionerData struct {
	ID            string `json:"id"`
	Person        Person `json:"person"`
	Active        bool   `json:"active"`
	Qualification string `json:"qualification"`
}

type ScheduleResults struct {
	ID               string           `json:"id"`
	PractitionerData PractitionerData `json:"practitioner_data"`
	Specialty        string           `json:"specialty"`
	Description      string           `json:"description"`
}

// Slot is used to models available slots data for a given schedule time
type Slot struct {
	ID    string `json:"id" mapstructure:"id"`
	Start string `json:"start" mapstructure:"start"`
	End   string `json:"end" mapstructure:"end"`
}

// AdvantageHeaders model is used to represent advantage's headers that are required for every request
type AdvantageHeaders struct {
	Organisation string `json:"organisation"`
	Cluster      string `json:"cluster"`
	Department   string `json:"department"`
	Branch       string `json:"branch"`
	Workstation  string `json:"workstation"`
	Variant      string `json:"variant,omitempty"`
}

// Checkin is is used to model the payload used to create check-in
type Checkin struct {
	Slot    string `json:"slot,omitempty"`
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Patient string `json:"patient,omitempty"`
}

// Schedule model is used to represent the response from the advantage's schedule API
type Schedule struct {
	ID               string           `json:"id"`
	CreatedByName    string           `json:"created_by_name"`
	UpdatedByName    string           `json:"updated_by_name"`
	Availability     any              `json:"availability"`
	Queue            string           `json:"queue"`
	PractitionerData PractitionerData `json:"practitioner_data"`
	WorkstationID    string           `json:"workstation_id"`
	DepartmentID     string           `json:"department_id"`
	BranchID         string           `json:"branch_id"`
	ClusterID        string           `json:"cluster_id"`
	Active           bool             `json:"active"`
	Created          time.Time        `json:"created"`
	CreatedBy        string           `json:"created_by"`
	Updated          time.Time        `json:"updated"`
	UpdatedBy        string           `json:"updated_by"`
	Actor            string           `json:"actor"`
	Specialty        string           `json:"specialty"`
	SlotDuration     int              `json:"slot_duration"`
	Description      string           `json:"description"`
	Organisation     string           `json:"organisation"`
	Practitioner     string           `json:"practitioner"`
}
