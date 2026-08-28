package dto

type DiagnosticReport struct {
	ID          string            `json:"id,omitempty"`
	Status      ObservationStatus `json:"status,omitempty"`
	PatientID   string            `json:"patientID,omitempty"`
	EncounterID string            `json:"encounterID,omitempty"`
	Issued      string            `json:"issued,omitempty"`
	Result      []*Observation    `json:"result,omitempty"`
	Media       []*Media          `json:"media,omitempty"`
	Conclusion  string            `json:"conclusion,omitempty"`
}
