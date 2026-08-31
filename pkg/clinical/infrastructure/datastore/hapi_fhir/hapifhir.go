package hapifhir

import (
	"context"
	"fmt"

	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
)

// HapiFHIRSDKImplementation defines the interface for interacting with a HAPI FHIR server.
// This interface allows for duck typing, enabling various implementations to be injected.
type HapiFHIRSDKImplementation interface {
	CreateFHIRResource(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error
	PatchFHIRResource(ctx context.Context, resourceType string, resourceID string, payload interface{}, resource interface{}) error
	GetFHIRResource(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error
	DeleteFHIRResource(ctx context.Context, resourceType, fhirResourceID string) error
	SearchFHIRResource(ctx context.Context, bundleID, resourceType string, params map[string]any, bundle interface{}) error
	ValidateResource(ctx context.Context, resourcType string, params map[string]interface{}) error
	GetPatientEverything(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error
	FHIRPathPatch(ctx context.Context, resourceType string, resourceID string, payload map[string]interface{}, resource interface{}) error
	PostFHIRBundle(ctx context.Context, payload interface{}, response interface{}) error
	PutFHIRResource(ctx context.Context, resourceType, resourceID string, payload map[string]any, resource any, useCREnabledServer bool) error
}

// HAPIImpl is a concrete implementation of the HapiFHIRSDKImplementation interface.
// It serves as a wrapper around an actual HAPI FHIR client.
type HAPIImpl struct {
	client HapiFHIRSDKImplementation
}

// NewHAPIFHIRImplementation initializes a new instance of HAPIImpl.
// It takes a client that adheres to the HapiFHIRSDKImplementation interface.
func NewHAPIFHIRImplementation(client HapiFHIRSDKImplementation) *HAPIImpl {
	return &HAPIImpl{
		client: client,
	}
}

// CreateFHIRResource delegates the task of creating a FHIR resource to the injected client
func (s *HAPIImpl) CreateFHIRResource(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
	payload = AddRequiredFieldToPayload(payload, resourceType)

	return s.client.CreateFHIRResource(ctx, resourceType, payload, resource)
}

// UpdateFHIRResource delegates the tasks of updating a FHIR resource to the injected client
func (s *HAPIImpl) PatchFHIRResource(ctx context.Context, resourceType string, resourceID string, payload interface{}, resource interface{}) error {
	return s.client.PatchFHIRResource(ctx, resourceType, resourceID, payload, resource)
}

// GetFHIRResource delegates the task for fetching a FHIR resource to the injected client
func (s *HAPIImpl) GetFHIRResource(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
	return s.client.GetFHIRResource(ctx, resourceType, fhirResourceID, resource)
}

// DeleteFHIRResource delegates the task of deleting a FHIR resource to the injected client
func (s *HAPIImpl) DeleteFHIRResource(ctx context.Context, resourceType, fhirResourceID string) error {
	return s.client.DeleteFHIRResource(ctx, resourceType, fhirResourceID)
}

// SearchFHIRResource delegates the task for searching a FHIR resource to the injeted client
func (s *HAPIImpl) SearchFHIRResource(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
	if tenant.OrganizationID != "" && tenant.FacilityID != "" {
		params["_tag"] = fmt.Sprintf("%s|%s", common.OrganisationSystem, tenant.OrganizationID)
		params["_tag"] = fmt.Sprintf("%s|%s", common.FacilitySystem, tenant.FacilityID)
	}

	if bundleID != "" {
		return s.client.SearchFHIRResource(ctx, bundleID, "", params, bundle)
	}

	return s.client.SearchFHIRResource(ctx, "", resourceType, params, bundle)
}

// GetFHIRPatientEverything is used to fetch all patient information.
func (s *HAPIImpl) GetPatientEverything(ctx context.Context, id string, searchParams map[string]interface{}, bundle interface{}) error {
	return s.client.GetPatientEverything(ctx, id, searchParams, bundle)
}

// ValidateFHIRResource is used to validate a FHIR resource.
func (s *HAPIImpl) ValidateResource(ctx context.Context, resourceType string, params map[string]interface{}) error {
	return s.client.ValidateResource(ctx, resourceType, params)
}

// FHIRPathPatch is a mechanism where elements to be manipulated by the patch interaction are described using their FHIRpath (https://www.hl7.org/fhir/fhirpath.html) names and navigation.
func (s *HAPIImpl) FHIRPathPatch(ctx context.Context, resourceType string, resourceID string, payload map[string]interface{}, resource interface{}) error {
	return s.client.FHIRPathPatch(ctx, resourceType, resourceID, payload, resource)
}

// PostFHIRBundle creates a FHIR bundle.
func (s *HAPIImpl) PostFHIRBundle(ctx context.Context, payload interface{}, response interface{}) error {
	return s.client.PostFHIRBundle(ctx, payload, response)
}

// PutFHIRResource updates a FHIR resource.
func (s *HAPIImpl) PutFHIRResource(ctx context.Context, resourceType, resourceID string, payload map[string]any, resource any, useCREnabledServer bool) error {
	return s.client.PutFHIRResource(ctx, resourceType, resourceID, payload, resource, useCREnabledServer)
}
