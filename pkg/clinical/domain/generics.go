package domain

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/savannahghi/scalarutils"
)

// HasFHIRSubject is an interface that should be implemented by any FHIR resource type
// that includes a subject reference. The subject typically points to a patient
// associated with the resource. Implementing this interface allows for generic
// operations to be performed on various FHIR resource types that have a subject.
//
// The GetSubject method returns a reference to the subject of the FHIR resource,
// which can be used to retrieve the associated patient's ID or other information.
//
// Example types that might implement this interface include FHIRDocumentReference,
// FHIRDiagnosticReport, FHIRRiskAssessment etc
type HasFHIRSubject interface {
	GetSubject() *FHIRReference
	GetFacilityFromMeta() *FHIRMeta
}

// GetPatientID retrieves the patient ID from any FHIR resource that implements
// the HasFHIRSubject interface. It returns an error if the subject or patient ID
// is not found.
func GetPatientID[T HasFHIRSubject](resource T) (string, error) {
	subject := resource.GetSubject()

	if subject != nil && subject.ID != nil {
		return *subject.ID, nil
	}

	return "", errors.New("no subject found")
}

func GetFacilityIDFromResource[T HasFHIRSubject](resource T) (string, error) {
	meta := resource.GetFacilityFromMeta()

	facilitySystem := scalarutils.URI("http://mycarehub/tenant-identification/facility")

	if meta != nil && len(meta.Tag) > 0 {
		for _, tag := range meta.Tag {
			if string(*tag.System) == string(facilitySystem) {
				return string(*tag.Code), nil
			}
		}
	}

	return "", errors.New("no facility ID found")
}

// ConvertMapToFHIRResource converts a slice of map[string]interface{} to a slice of a specific FHIR resource type.
func ConvertMapToFHIRResource[T any](maps map[string]interface{}) (*T, error) {
	jsonData, err := json.Marshal(maps)
	if err != nil {
		return nil, fmt.Errorf("error marshalling map to JSON: %w", err)
	}

	var resource T

	err = json.Unmarshal(jsonData, &resource)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON to target type: %w", err)
	}

	return &resource, nil
}
