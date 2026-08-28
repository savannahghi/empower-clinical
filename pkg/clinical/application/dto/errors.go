package dto

// ErrorDetails contains more details about the error that occurred while making a REST API call to FHIR servers
type ErrorDetails struct {
	Text string `json:"text"`
}

// ErrorIssue models the error issue returned from FHIR
type ErrorIssue struct {
	Details     ErrorDetails `json:"details"`
	Diagnostics string       `json:"diagnostics"`
}

// ErrorResponse models the json object returned when an error occurs while calling the
// FHIR server REST APIs
type ErrorResponse struct {
	Issue []ErrorIssue `json:"issue"`
}
