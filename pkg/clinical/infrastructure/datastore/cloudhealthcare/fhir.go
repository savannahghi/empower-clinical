package fhir

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"github.com/savannahghi/converterandformatter"
	hapifhirmodels "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/cache"
)

// constants and defaults
const (
	internalError = "an error occurred on our end. Please try again later"
	timeFormatStr = "2006-01-02T15:04:05+03:00"
)

// resource types
const (
	organizationResource              = "Organization"
	patientResourceType               = "Patient"
	episodeOfCareResourceType         = "EpisodeOfCare"
	observationResourceType           = "Observation"
	allergyIntoleranceResourceType    = "AllergyIntolerance"
	serviceRequestResourceType        = "ServiceRequest"
	medicationRequestResourceType     = "MedicationRequest"
	conditionResourceType             = "Condition"
	encounterResourceType             = "Encounter"
	compositionResourceType           = "Composition"
	medicationStatementResourceType   = "MedicationStatement"
	medicationResourceType            = "Medication"
	mediaResourceType                 = "Media"
	questionnaireResourceType         = "Questionnaire"
	consentResourceType               = "Consent"
	questionnaireResponseResourceType = "QuestionnaireResponse"
	riskAssessmentResourceType        = "RiskAssessment"
	diagnosticReportResourceType      = "DiagnosticReport"
	subscriptionResourceType          = "Subscription"
	documentReferenceResourceType     = "DocumentReference"
	appointmentResourceType           = "Appointment"
	taskResourceType                  = "Task"
	locationResourceType              = "Location"
	practitionerRoleResourceType      = "PractitionerRole"
	practitionerResourceType          = "Practitioner"
	substanceResourceType             = "Substance"
	procedureResourceType             = "Procedure"
	medicationDispenseResourceType    = "MedicationDispense"
	planDefinitionResourceType        = "PlanDefinition"
	activityDefinitionResourceType    = "ActivityDefinition"
	carePlanResourceType              = "CarePlan"
)

type HapiFHIRImplementation interface {
	CreateFHIRResource(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error
	PatchFHIRResource(ctx context.Context, resourceType string, resourceID string, payload interface{}, resource interface{}) error
	GetFHIRResource(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error
	PutFHIRResource(ctx context.Context, resourceType, resourceID string, payload map[string]any, resource any, useCREnabledServer bool) error
	DeleteFHIRResource(ctx context.Context, resourceType, fhirResourceID string) error
	SearchFHIRResource(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error
	GetPatientEverything(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error
	ValidateResource(ctx context.Context, resourceType string, params map[string]interface{}) error
	FHIRPathPatch(ctx context.Context, resourceType string, resourceID string, payload map[string]interface{}, resource interface{}) error
	PostFHIRBundle(ctx context.Context, payload interface{}, response interface{}) error
}

// StoreImpl represents the FHIR infrastructure implementation
type StoreImpl struct {
	cache        cache.CacheService
	HapiFHIRImpl HapiFHIRImplementation
}

// NewFHIRStoreImpl initializes the new FHIR implementation
func NewFHIRStoreImpl(
	cache cache.CacheService,
	hapiFHIR HapiFHIRImplementation,
) *StoreImpl {
	return &StoreImpl{
		cache:        cache,
		HapiFHIRImpl: hapiFHIR,
	}
}

// FHIROperation defines a function type that represents an operation to be performed on FHIR data.
type FHIROperation func() (interface{}, error)

// GetOrSetCache is a helper function to handle caching logic. It attempts to retrieve data
// from the cache. If the data is not found (cache miss), it fetches the data using the provided
// FHIROperation, caches the fetched data, and then returns it. This function helps reduce
// latency by serving data from the cache for subsequent requests
func (fh StoreImpl) GetOrSetCache(ctx context.Context, key string, op FHIROperation) ([]byte, error) {
	cachedData := fh.cache.Get(ctx, key)

	cachedResult, err := cachedData.Result()
	if errors.Is(err, redis.Nil) {
		data, err := op()
		if err != nil {
			return nil, err
		}

		marshaledData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		err = fh.cache.Set(ctx, key, marshaledData, time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return marshaledData, nil
	} else if err != nil {
		return nil, err
	}

	return []byte(cachedResult), nil
}

// SetCache marshals the value and sets it in the cache with the specified expiration time.
func (fh StoreImpl) SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	marshaledData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("unable to marshal value for cache: %w", err)
	}

	err = fh.cache.Set(ctx, key, marshaledData, expiration).Err()
	if err != nil {
		return fmt.Errorf("unable to set cache: %w", err)
	}

	return nil
}

// SearchPatientObservations fetches all observations that belong to a specific patient
func (fh StoreImpl) SearchPatientObservations(
	ctx context.Context,
	bundleID string,
	searchParameters map[string]interface{},
	tenant dto.TenantIdentifiers,
	pagination serverutils.PaginationInput,
) (*domain.PagedFHIRObservations, error) {
	observationsBundle := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, bundleID, observationResourceType, searchParameters, tenant, observationsBundle)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextURL, previousURL string

	for _, link := range observationsBundle.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextURL = link.URL
		case "previous":
			hasPreviousPage = true
			previousURL = link.URL
		}
	}

	observationOutput := domain.PagedFHIRObservations{
		BundleID:        *observationsBundle.ID,
		Observations:    []domain.FHIRObservation{},
		HasNextPage:     hasNextPage,
		NextPageURL:     nextURL,
		HasPreviousPage: hasPreviousPage,
		PreviousPageURL: previousURL,
		TotalCount:      *observationsBundle.Total,
	}

	for _, obs := range observationsBundle.Entry {
		var observation domain.FHIRObservation

		err = json.Unmarshal(obs.Resource, &observation)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		observationOutput.Observations = append(observationOutput.Observations, observation)
	}

	return &observationOutput, nil
}

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (fh StoreImpl) SearchPatientEncounters(
	ctx context.Context,
	patientReference string,
	status *domain.EncounterStatusEnum,
	tenant dto.TenantIdentifiers,
	pagination dto.Pagination,
) (*domain.PagedFHIREncounter, error) {
	params := map[string]interface{}{
		"patient": patientReference,
		"_total":  "accurate",
	}
	if status != nil {
		params["status:exact"] = status.String()
	}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", encounterResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	encounterOutput := domain.PagedFHIREncounter{
		Encounters:      []domain.FHIREncounter{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, resource := range resources.Entry {
		var encounter domain.FHIREncounter

		err = json.Unmarshal(resource.Resource, &encounter)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		encounterOutput.Encounters = append(encounterOutput.Encounters, encounter)
	}

	return &encounterOutput, nil
}

// SearchPatentMedia searches all the patients media resources
func (fh StoreImpl) SearchPatientMedia(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedia, error) {
	params := map[string]interface{}{
		"patient": patientReference,
	}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", mediaResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			nextCursor = link.URL
			hasNextPage = true
		case "previous":
			previousCursor = link.URL
			hasPreviousPage = true
		}
	}

	mediaOutput := domain.PagedFHIRMedia{
		Media:           []domain.FHIRMedia{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, resource := range resources.Entry {
		var media domain.FHIRMedia

		resourceBs, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &media)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		mediaOutput.Media = append(mediaOutput.Media, media)
	}

	return &mediaOutput, nil
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (fh StoreImpl) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	output := domain.FHIREpisodeOfCareRelayConnection{}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", episodeOfCareResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Entry {
		var resource domain.FHIREpisodeOfCare

		err = json.Unmarshal(result.Resource, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", episodeOfCareResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIREpisodeOfCareRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateEpisodeOfCare is the final common pathway for creation of episodes of
// care.
func (fh StoreImpl) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error) {
	payload, err := converterandformatter.StructToMap(episode)
	if err != nil {
		return nil, fmt.Errorf("unable to turn episode of care input into a map: %w", err)
	}

	fhirEpisode := &domain.FHIREpisodeOfCare{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, episodeOfCareResourceType, payload, fhirEpisode)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to create episode of care resource: %w", err)
	}

	output := &domain.EpisodeOfCarePayload{
		EpisodeOfCare: fhirEpisode,
	}

	return output, nil
}

// CreateFHIRCondition creates a FHIRCondition instance
func (fh StoreImpl) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", conditionResourceType, err)
	}

	resource := &domain.FHIRCondition{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, conditionResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", conditionResourceType, err)
	}

	output := &domain.FHIRConditionRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// GetFHIRCondition retrieves an instance of a fhir condition by its ID
func (fh StoreImpl) GetFHIRCondition(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
	resource := &domain.FHIRCondition{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, conditionResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", conditionResourceType, id, err)
	}

	payload := &domain.FHIRConditionRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// CreateFHIROrganization creates a FHIROrganization instance
func (fh StoreImpl) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", organizationResource, err)
	}

	resource := &domain.FHIROrganization{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, organizationResource, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", organizationResource, err)
	}

	output := &domain.FHIROrganizationRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// SearchFHIROrganization provides a search API for FHIROrganization
func (fh StoreImpl) SearchFHIROrganization(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIROrganizationRelayConnection, error) {
	output := domain.FHIROrganizationRelayConnection{}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", organizationResource, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Entry {
		var resource domain.FHIROrganization

		err = json.Unmarshal(result.Resource, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", organizationResource, err)
		}

		output.Edges = append(output.Edges, &domain.FHIROrganizationRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// GetFHIROrganization finds and retrieves organization details using the specified organization ID
func (fh StoreImpl) GetFHIROrganization(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}

	organization := &domain.FHIROrganization{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, organizationResource, organizationID, organization)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve organization: %w", err)
	}

	return &domain.FHIROrganizationRelayPayload{
		Resource: organization,
	}, nil
}

// PutFHIROrganization updates an existing FHIROrganization instance with the provided input data. It identifies the organization to update using the specified ID and applies the changes defined in the input. If the organization with the given ID does not exist, it returns an error indicating that the resource was not found.
func (fh StoreImpl) PutFHIROrganization(ctx context.Context, id string, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", organizationResource, err)
	}

	resource := &domain.FHIROrganization{}

	err = fh.HapiFHIRImpl.PutFHIRResource(ctx, organizationResource, id, payload, resource, false)
	if err != nil {
		return nil, fmt.Errorf("unable to PUT %s resource: %w", organizationResource, err)
	}

	output := &domain.FHIROrganizationRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// GetFHIRAllergyIntolerance fetches the allergy from FHIR repository using its id
func (fh StoreImpl) GetFHIRAllergyIntolerance(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	allergyIntoleranace := &domain.FHIRAllergyIntolerance{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, allergyIntoleranceResourceType, id, allergyIntoleranace)
	if err != nil {
		return nil, err
	}

	return &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: allergyIntoleranace,
	}, nil
}

// SearchEpisodesByParam search episodes by params
func (fh StoreImpl) SearchEpisodesByParam(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", episodeOfCareResourceType, searchParams, tenant, resources)
	if err != nil {
		return nil, err
	}

	output := []*domain.FHIREpisodeOfCare{}

	var episode domain.FHIREpisodeOfCare

	err = mapstructure.Decode(resources.Identifier.Period, &episode)
	if err != nil {
		return nil, err
	}

	output = append(output, &episode)

	return output, nil
}

// OpenEpisodes returns the IDs of a patient's open episodes
func (fh StoreImpl) OpenEpisodes(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error) {
	params := map[string]interface{}{
		"status:exact": domain.EpisodeOfCareStatusEnumActive.String(),
		"patient":      patientReference,
	}

	return fh.SearchEpisodesByParam(ctx, params, tenant, pagination)
}

// HasOpenEpisode determines if a patient has an open episode
func (fh StoreImpl) HasOpenEpisode(
	ctx context.Context,
	patient domain.FHIRPatient,
	tenant dto.TenantIdentifiers,
	pagination dto.Pagination,
) (bool, error) {
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)

	episodes, err := fh.OpenEpisodes(ctx, patientReference, tenant, pagination)
	if err != nil {
		return false, err
	}

	return len(episodes) > 0, nil
}

// CreateFHIREncounter creates a FHIREncounter instance
func (fh StoreImpl) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", encounterResourceType, err)
	}

	resource := &domain.FHIREncounter{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, encounterResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", encounterResourceType, err)
	}

	output := &domain.FHIREncounterRelayPayload{
		Resource: resource,
	}

	cacheKey := ""

	if output.Resource.ID != nil {
		cacheKey = fmt.Sprintf("%s:%s", encounterResourceType, *output.Resource.ID)
	}

	err = fh.SetCache(ctx, cacheKey, output, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetFHIREpisodeOfCare retrieves instances of FHIREpisodeOfCare by ID
func (fh StoreImpl) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	resource := &domain.FHIREpisodeOfCare{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, episodeOfCareResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", episodeOfCareResourceType, id, err)
	}

	payload := &domain.FHIREpisodeOfCareRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// StartEncounter starts an encounter within an episode of care
func (fh StoreImpl) StartEncounter(
	ctx context.Context, episodeID string) (string, error) {
	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return "", fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}

	activeEpisodeStatus := domain.EpisodeOfCareStatusEnumActive
	activeEncounterStatus := domain.EncounterStatusEnumInProgress

	if episodePayload.Resource.Status.String() != activeEpisodeStatus.String() {
		return "", fmt.Errorf("an encounter can only be started for an active episode")
	}

	episodeRef := fmt.Sprintf("EpisodeOfCare/%s", *episodePayload.Resource.ID)
	now := time.Now()
	startTime := scalarutils.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	encounterClassCode := scalarutils.Code("AMB")
	encounterClassSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassVersion := "2018-08-12"
	encounterClassDisplay := "ambulatory"
	encounterClassUserSelected := false

	encounterInput := domain.FHIREncounterInput{
		Status: activeEncounterStatus,
		Class: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &encounterClassSystem,
						Version:      &encounterClassVersion,
						Code:         encounterClassCode,
						Display:      encounterClassDisplay,
						UserSelected: &encounterClassUserSelected,
					},
				},
			},
		},
		Subject: &domain.FHIRReferenceInput{
			Reference: episodePayload.Resource.Patient.Reference,
			Display:   episodePayload.Resource.Patient.Display,
			Type:      episodePayload.Resource.Patient.Type,
		},
		EpisodeOfCare: []*domain.FHIRReferenceInput{
			{
				Reference: &episodeRef,
			},
		},
		ServiceProvider: &domain.FHIRReferenceInput{
			Display: episodePayload.Resource.ManagingOrganization.Display,
			Type:    episodePayload.Resource.ManagingOrganization.Type,
		},
		ActualPeriod: &domain.FHIRPeriodInput{
			Start: &startTime,
		},
	}

	encPl, err := fh.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		return "", fmt.Errorf("unable to start encounter: %w", err)
	}

	return *encPl.Resource.ID, nil
}

// PatchFHIREncounter is used to patch an encounter resource
func (fh StoreImpl) PatchFHIREncounter(
	ctx context.Context,
	encounterID string,
	input domain.FHIREncounterInput,
) (*domain.FHIREncounter, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", encounterResourceType, err)
	}

	resource := &domain.FHIREncounter{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, encounterResourceType, encounterID, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// SearchEpisodeEncounter returns all encounters in a visit
func (fh StoreImpl) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
	tenant dto.TenantIdentifiers,
	pagination dto.Pagination,
) (*domain.PagedFHIREncounter, error) {
	episodeRef := fmt.Sprintf("EpisodeOfCare/%s", episodeReference)
	encounterFilterParams := map[string]interface{}{
		"episode-of-care": episodeRef,
		"status":          "in_progress",
		"_total":          "accurate",
	}

	encounterConn, err := fh.SearchFHIREncounter(ctx, encounterFilterParams, tenant, pagination)
	if err != nil {
		return nil, fmt.Errorf("unable to search encounter: %w", err)
	}

	return encounterConn, nil
}

// EndEncounter ends an encounter
func (fh StoreImpl) EndEncounter(
	ctx context.Context, encounterID string) (bool, error) {
	_, err := fh.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return false, err
	}

	updatedStatus := domain.EncounterStatusEnumCompleted
	// workaround for odd date comparison behavior on the Google Cloud Healthcare API
	// the end time must be at least 24 hours after the start time
	// so: if the time now is less than 24 hours after start, set the end to be
	// 24 hours after the start of the visit. If the time now is more than 24 hours
	// after the start, use the current time as the end of the visit
	end := time.Now().Add(time.Hour * 24)
	endTime := scalarutils.DateTime(end.Format(timeFormatStr))

	updateData := &domain.FHIREncounterInput{
		Status: updatedStatus,
		ActualPeriod: &domain.FHIRPeriodInput{
			End: &endTime,
		},
	}

	payload, err := converterandformatter.StructToMap(updateData)
	if err != nil {
		return false, fmt.Errorf("unable to turn %s input into a map: %w", encounterResourceType, err)
	}

	resource := &domain.FHIREncounter{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, encounterResourceType, encounterID, payload, resource)
	if err != nil {
		return false, err
	}

	return true, nil
}

// EndEpisode ends an episode of care by patching its status to "finished"
func (fh StoreImpl) EndEpisode(
	ctx context.Context, episodeID string) (bool, error) {
	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return false, fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}

	startTime := scalarutils.DateTime(time.Now().Format(timeFormatStr))
	if episodePayload.Resource.Period != nil {
		startTime = episodePayload.Resource.Period.Start
	}

	// workaround for odd date comparison behavior on the Google Cloud Healthcare API
	// the end time must be at least 24 hours after the start time
	// so: if the time now is less than 24 hours after start, set the end to be
	// 24 hours after the start of the visit. If the time now is more than 24 hours
	// after the start, use the current time as the end of the visit
	end := time.Now().Add(time.Hour * 24)
	endTime := scalarutils.DateTime(end.Format(timeFormatStr))
	updatedStatus := domain.EpisodeOfCareStatusEnumFinished

	episode := &domain.FHIREpisodeOfCare{}

	updateData := &domain.FHIREpisodeOfCareInput{
		Status: &updatedStatus,
		Period: &domain.FHIRPeriodInput{
			Start: &startTime,
			End:   &endTime,
		},
	}

	payload, err := converterandformatter.StructToMap(updateData)
	if err != nil {
		return false, fmt.Errorf("unable to turn %s input into a map: %w", episodeOfCareResourceType, err)
	}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, episodeOfCareResourceType, episodeID, payload, episode)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetActiveEpisode returns any ACTIVE episode that has to the indicated ID
func (fh StoreImpl) GetActiveEpisode(ctx context.Context, episodeID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCare, error) {
	params := map[string]interface{}{
		"status:exact": domain.EpisodeOfCareStatusEnumActive.String(),
		"_id":          episodeID,
		"_total":       "accurate",
	}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", episodeOfCareResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	if len(resources.Entry) != 1 {
		return nil, fmt.Errorf(
			"expected exactly one ACTIVE episode for episode ID %s, got %d", episodeID, len(resources.Entry))
	}

	var episode domain.FHIREpisodeOfCare

	resourceBs, err := json.Marshal(resources.Entry[0])
	if err != nil {
		return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
	}

	err = json.Unmarshal(resourceBs, &episode)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
	}

	return &episode, nil
}

// SearchFHIRServiceRequest provides a search API for FHIRServiceRequests
func (fh StoreImpl) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRServiceRequest, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", serviceRequestResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	serviceRequestOutput := domain.PagedFHIRServiceRequest{
		ServiceRequests: make([]domain.FHIRServiceRequest, 0, len(resources.Entry)),
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRServiceRequest

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", serviceRequestResourceType, err)
		}

		serviceRequestOutput.ServiceRequests = append(serviceRequestOutput.ServiceRequests, resource)
	}

	return &serviceRequestOutput, nil
}

// CreateFHIRServiceRequest creates a FHIRServiceRequest instance
func (fh StoreImpl) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", serviceRequestResourceType, err)
	}

	resource := &domain.FHIRServiceRequest{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, serviceRequestResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", serviceRequestResourceType, err)
	}

	output := &domain.FHIRServiceRequestRelayPayload{
		Resource: resource,
	}

	cacheKey := fmt.Sprintf("%s:%s", serviceRequestResourceType, *output.Resource.ID)

	err = fh.SetCache(ctx, cacheKey, output, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SearchFHIRAllergyIntolerance provides a search API for FHIRAllergyIntolerance
func (fh StoreImpl) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", allergyIntoleranceResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	output := domain.PagedFHIRAllergy{
		Allergies:       []domain.FHIRAllergyIntolerance{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRAllergyIntolerance

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", allergyIntoleranceResourceType, err)
		}

		output.Allergies = append(output.Allergies, resource)
	}

	return &output, nil
}

// CreateFHIRAllergyIntolerance creates a FHIRAllergyIntolerance instance
func (fh StoreImpl) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", allergyIntoleranceResourceType, err)
	}

	resource := &domain.FHIRAllergyIntolerance{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, allergyIntoleranceResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", allergyIntoleranceResourceType, err)
	}

	output := &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// UpdateFHIRAllergyIntolerance updates a FHIRAllergyIntolerance instance
// The resource must have its ID set.
func (fh StoreImpl) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", allergyIntoleranceResourceType, err)
	}

	resource := &domain.FHIRAllergyIntolerance{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, allergyIntoleranceResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, err
	}

	output := &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// SearchFHIRComposition provides a search API for FHIRComposition
func (fh StoreImpl) SearchFHIRComposition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRComposition, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", compositionResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	output := domain.PagedFHIRComposition{
		Compositions:    []domain.FHIRComposition{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRComposition

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", compositionResourceType, err)
		}

		output.Compositions = append(output.Compositions, resource)
	}

	return &output, nil
}

// CreateFHIRComposition creates a FHIRComposition instance
func (fh StoreImpl) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", compositionResourceType, err)
	}

	resource := &domain.FHIRComposition{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, compositionResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", compositionResourceType, err)
	}

	output := &domain.FHIRCompositionRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// UpdateFHIRComposition updates a FHIRComposition instance
// The resource must have its ID set.
func (fh StoreImpl) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", compositionResourceType, err)
	}

	resource := &domain.FHIRComposition{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, compositionResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// PatchFHIRTask is used to update task resource
func (fh StoreImpl) PatchFHIRTask(ctx context.Context, input domain.FHIRTaskInput) (*domain.FHIRTask, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", taskResourceType, err)
	}

	resource := &domain.FHIRTask{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, taskResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// DeleteFHIRComposition deletes the FHIRComposition identified by the supplied ID
func (fh StoreImpl) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	err := fh.HapiFHIRImpl.DeleteFHIRResource(ctx, compositionResourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			compositionResourceType, err,
		)
	}

	return true, nil
}

// SearchFHIRCondition provides a search API for FHIRCondition
func (fh StoreImpl) SearchFHIRCondition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", conditionResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	output := domain.PagedFHIRCondition{
		Conditions:      []domain.FHIRCondition{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRCondition

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", conditionResourceType, err)
		}

		output.Conditions = append(output.Conditions, resource)
	}

	return &output, nil
}

// SearchPatientAllergyIntolerance searches for a patient's FHIR allergy intolerance using patient ID
func (fh StoreImpl) SearchPatientAllergyIntolerance(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
	params := map[string]interface{}{
		"patient": patientReference,
	}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", allergyIntoleranceResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	HasNextPage := false
	HasPreviousPage := false

	var NextCursor string

	var PreviousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			NextCursor = link.URL
			HasNextPage = true
		case "previous":
			PreviousCursor = link.URL
			HasPreviousPage = true
		}
	}

	allergyOutput := domain.PagedFHIRAllergy{
		Allergies:       []domain.FHIRAllergyIntolerance{},
		HasNextPage:     HasNextPage,
		NextCursor:      NextCursor,
		HasPreviousPage: HasPreviousPage,
		PreviousCursor:  PreviousCursor,
		TotalCount:      *resources.Total,
	}

	for _, resource := range resources.Entry {
		var allergyIntolerance domain.FHIRAllergyIntolerance

		resourceBs, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &allergyIntolerance)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", allergyIntoleranceResourceType, err)
		}

		allergyOutput.Allergies = append(allergyOutput.Allergies, allergyIntolerance)
	}

	return &allergyOutput, nil
}

// UpdateFHIRCondition updates a FHIRCondition instance
// The resource must have its ID set.
func (fh StoreImpl) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", conditionResourceType, err)
	}

	resource := &domain.FHIRCondition{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, conditionResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, err
	}

	output := &domain.FHIRConditionRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// GetFHIREncounter retrieves instances of FHIREncounter by ID
func (fh StoreImpl) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	cacheKey := fmt.Sprintf("encounter:%s", id)

	fetchFunc := func() (interface{}, error) {
		resource := &domain.FHIREncounter{}

		err := fh.HapiFHIRImpl.GetFHIRResource(ctx, encounterResourceType, id, resource)
		if err != nil {
			return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", encounterResourceType, id, err)
		}

		payload := &domain.FHIREncounterRelayPayload{
			Resource: resource,
		}

		return payload, nil
	}

	data, err := fh.GetOrSetCache(ctx, cacheKey, fetchFunc)
	if err != nil {
		return nil, err
	}

	var encounterPayload *domain.FHIREncounterRelayPayload

	err = json.Unmarshal(data, &encounterPayload)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal cached data: %w", err)
	}

	return encounterPayload, nil
}

// GetFHIRMedicationRequest is used to retrieve medication request information
func (fh StoreImpl) GetFHIRMedicationRequest(ctx context.Context, id string) (*domain.FHIRMedicationRequestRelayPayload, error) {
	cacheKey := fmt.Sprintf("medicationrequest:%s", id)
	fetchFunc := func() (interface{}, error) {
		resource := &domain.FHIRMedicationRequest{}

		err := fh.HapiFHIRImpl.GetFHIRResource(ctx, medicationRequestResourceType, id, resource)
		if err != nil {
			return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", medicationRequestResourceType, id, err)
		}

		payload := &domain.FHIRMedicationRequestRelayPayload{
			Resource: resource,
		}

		return payload, nil
	}

	data, err := fh.GetOrSetCache(ctx, cacheKey, fetchFunc)
	if err != nil {
		return nil, err
	}

	var medicationRequest *domain.FHIRMedicationRequestRelayPayload

	err = json.Unmarshal(data, &medicationRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal cached data: %w", err)
	}

	return medicationRequest, nil
}

// SearchFHIREncounter provides a search API for FHIREncounter
func (fh StoreImpl) SearchFHIREncounter(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", encounterResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	encounterOutput := domain.PagedFHIREncounter{
		Encounters:      []domain.FHIREncounter{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIREncounter

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", encounterResourceType, err)
		}

		encounterOutput.Encounters = append(encounterOutput.Encounters, resource)
	}

	return &encounterOutput, nil
}

// SearchFHIREncounterAllData provides a search API for a FHIREncounter and all other resources that reference the encounter
func (fh StoreImpl) SearchFHIREncounterAllData(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", encounterResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	encounterAllDataOutput := domain.PagedFHIRResource{
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, entry := range resources.Entry {
		resource := map[string]interface{}{}

		entryBytes, err := entry.Resource.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("server error: failed to marshal resource to Bytes: %w", err)
		}

		if err = json.Unmarshal(entryBytes, &resource); err != nil {
			return nil, fmt.Errorf("server error: failed to convert resource to a map")
		}

		encounterAllDataOutput.Resources = append(encounterAllDataOutput.Resources, resource)
	}

	return &encounterAllDataOutput, nil
}

// SearchFHIRMedicationRequest provides a search API for FHIRMedicationRequest
func (fh StoreImpl) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedicationRequest, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", medicationRequestResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	medicationRequestOutput := domain.PagedFHIRMedicationRequest{
		MedicationRequests: []domain.FHIRMedicationRequest{},
		HasNextPage:        hasNextPage,
		NextCursor:         nextCursor,
		HasPreviousPage:    hasPreviousPage,
		PreviousCursor:     previousCursor,
		TotalCount:         *resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRMedicationRequest

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", medicationRequestResourceType, err)
		}

		medicationRequestOutput.MedicationRequests = append(medicationRequestOutput.MedicationRequests, resource)
	}

	return &medicationRequestOutput, nil
}

// CreateFHIRMedicationRequest creates a FHIRMedicationRequest instance
func (fh StoreImpl) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationRequestResourceType, err)
	}

	resource := &domain.FHIRMedicationRequest{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, medicationRequestResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", medicationRequestResourceType, err)
	}

	output := &domain.FHIRMedicationRequestRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// UpdateFHIRMedicationRequest updates a FHIRMedicationRequest instance
// The resource must have its ID set.
func (fh StoreImpl) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	resource := &domain.FHIRMedicationRequest{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationRequestResourceType, err)
	}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, medicationRequestResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, err
	}

	output := &domain.FHIRMedicationRequestRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// DeleteFHIRMedicationRequest deletes the FHIRMedicationRequest identified by the supplied ID
func (fh StoreImpl) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	err := fh.HapiFHIRImpl.DeleteFHIRResource(ctx, medicationRequestResourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			medicationRequestResourceType, err,
		)
	}

	return true, nil
}

// SearchFHIRObservation provides a search API for FHIRObservation
func (fh StoreImpl) SearchFHIRObservation(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", observationResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextURL, previousURL string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextURL = link.URL
		case "previous":
			hasPreviousPage = true
			previousURL = link.URL
		}
	}

	observationOutput := domain.PagedFHIRObservations{
		Observations:    []domain.FHIRObservation{},
		HasNextPage:     hasNextPage,
		NextPageURL:     nextURL,
		HasPreviousPage: hasPreviousPage,
		PreviousPageURL: previousURL,
		TotalCount:      *resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRObservation

		err = json.Unmarshal(result.Resource, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", observationResourceType, err)
		}

		observationOutput.Observations = append(observationOutput.Observations, resource)
	}

	return &observationOutput, nil
}

func (fh StoreImpl) SearchFHIRDiagnosticReport(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDiagnosticReport, error) {
	diagnosticReport := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", diagnosticReportResourceType, params, tenant, diagnosticReport)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range diagnosticReport.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	output := domain.PagedFHIRDiagnosticReport{
		DiagnosticReport: []domain.FHIRDiagnosticReport{},
		HasNextPage:      hasNextPage,
		NextCursor:       nextCursor,
		HasPreviousPage:  hasPreviousPage,
		PreviousCursor:   previousCursor,
		TotalCount:       *diagnosticReport.Total,
	}

	for _, entry := range diagnosticReport.Entry {
		var report domain.FHIRDiagnosticReport

		reportEntry, err := json.Marshal(&entry.Resource)
		if err != nil {
			return nil, fmt.Errorf("could not marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(reportEntry, &report)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal resource: %w", err)
		}

		output.DiagnosticReport = append(output.DiagnosticReport, report)
	}

	return &output, nil
}

// CreateFHIRObservation creates a FHIRObservation instance
func (fh StoreImpl) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", observationResourceType, err)
	}

	resource := &domain.FHIRObservation{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, observationResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", observationResourceType, err)
	}

	return resource, nil
}

// DeleteFHIRObservation deletes the FHIRObservation identified by the passed ID
func (fh StoreImpl) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	err := fh.HapiFHIRImpl.DeleteFHIRResource(ctx, observationResourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			observationResourceType, err,
		)
	}

	return true, nil
}

func (fh StoreImpl) DeleteFHIRResource(ctx context.Context, resourceType, resourceID string) (bool, error) {
	err := fh.HapiFHIRImpl.DeleteFHIRResource(ctx, resourceType, resourceID)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			resourceType, err,
		)
	}

	return true, nil
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (fh StoreImpl) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	cacheKey := fmt.Sprintf("patient:%s", id)

	fetchPatientFunc := func() (interface{}, error) {
		resource := &domain.FHIRPatient{}

		err := fh.HapiFHIRImpl.GetFHIRResource(ctx, patientResourceType, id, resource)
		if err != nil {
			return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", patientResourceType, id, err)
		}

		payload := &domain.FHIRPatientRelayPayload{
			Resource: resource,
		}

		return payload, nil
	}

	data, err := fh.GetOrSetCache(ctx, cacheKey, fetchPatientFunc)
	if err != nil {
		return nil, err
	}

	var patientPayload *domain.FHIRPatientRelayPayload

	err = json.Unmarshal(data, &patientPayload)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal cached data: %w", err)
	}

	return patientPayload, nil
}

// DeleteFHIRPatient deletes the FHIRPatient identified by the supplied ID
func (fh StoreImpl) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	includedResources := []string{
		"EpisodeOfCare:patient",
		"Observation:subject",
		"Encounter:subject",
		"MedicationRequest:patient",
	}
	params := map[string]interface{}{
		"_id":         id,
		"_revinclude": includedResources,
	}

	patientEverythingBs := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.GetPatientEverything(ctx, id, params, patientEverythingBs)
	if err != nil {
		return false, fmt.Errorf("unable to get patient's compartment: %w", err)
	}

	// This list stores assorted ResourceTypes and ResourceIDs found in an Encounter
	// i.e. Observations, Medication Request etc
	assortedResourceTypes := []map[string]string{}

	// This list stores all the encounters ResourceType and ResourceID in a patient's compartment
	encounters := []map[string]string{}

	// This list stores all the Episodesofcare ResourceType and ResourceIDs in a patient's compartment
	episodesOfCare := []map[string]string{}

	// This list stores the patient ResourceType and ResourceID
	patient := []map[string]string{}

	// This list stores the task ResourceType and ResourceID
	task := []map[string]string{}

	// This list stores all the observations resource types
	observations := []map[string]string{}

	medicationRequests := []map[string]string{}

	for _, en := range patientEverythingBs.Entry {
		var resourceMap map[string]interface{}

		err := json.Unmarshal(en.Resource, &resourceMap)
		if err != nil {
			return false, err
		}

		resourceType := resourceMap["resourceType"].(string)

		resourceTypeIDMap := map[string]string{
			"resourceType": resourceType,
			"resourceID":   resourceMap["id"].(string),
		}

		switch resourceType {
		case encounterResourceType:
			encounters = append(
				encounters,
				resourceTypeIDMap,
			)

			continue

		case episodeOfCareResourceType:
			episodesOfCare = append(
				episodesOfCare,
				resourceTypeIDMap,
			)

			continue

		case patientResourceType:
			patient = append(
				patient,
				resourceTypeIDMap,
			)

			continue

		case observationResourceType:
			observations = append(
				observations,
				resourceTypeIDMap,
			)

			continue

		case organizationResource:
			continue

		case taskResourceType:
			task = append(task, resourceTypeIDMap)
			continue

		case medicationRequestResourceType:
			medicationRequests = append(
				medicationRequests,
				resourceTypeIDMap,
			)

			continue
		}

		assortedResourceTypes = append(
			assortedResourceTypes,
			resourceTypeIDMap,
		)
	}

	// Special case, a medication request causes the failure for deleting a FHIR Condition
	if err = fh.DeleteFHIRResourceType(ctx, medicationRequests); err != nil {
		return false, err
	}

	// Order of deletion matters to avoid conflicts
	// First delete the ResourceTypes found in an encounter
	if err = fh.DeleteFHIRResourceType(ctx, assortedResourceTypes); err != nil {
		return false, err
	}

	// Secondly, delete the encounters. This will bring no conflict
	// as it ensures ResourceType that refers to the encounter is not found
	if err = fh.DeleteFHIRResourceType(ctx, encounters); err != nil {
		return false, err
	}

	// Thirdly, delete the episodes of care. This will bring no conflict
	// as it ensures Encounter that refers to the EpisodeOfCare is not found
	if err = fh.DeleteFHIRResourceType(ctx, episodesOfCare); err != nil {
		return false, err
	}

	// Fourthly, delete the observaion. This will bring no conflict
	// as it ensures Encounter that refers to the Observation is not found
	if err = fh.DeleteFHIRResourceType(ctx, observations); err != nil {
		return false, err
	}

	// Fifthly, delete the task. This will bring no conflict
	if err = fh.DeleteFHIRResourceType(ctx, task); err != nil {
		return false, err
	}
	// Finally delete the patient ResourceType
	if err = fh.DeleteFHIRResourceType(ctx, patient); err != nil {
		return false, err
	}

	return true, nil
}

// DeleteFHIRResourceType takes a ResourceType and ID and deletes them from FHIR
func (fh StoreImpl) DeleteFHIRResourceType(ctx context.Context, results []map[string]string) error {
	for _, result := range results {
		resourceType := result["resourceType"]
		resourceID := result["resourceID"]

		err := fh.HapiFHIRImpl.DeleteFHIRResource(ctx, resourceType, resourceID)
		if err != nil {
			return fmt.Errorf(
				"unable to delete %s:%s, error: %w",
				resourceType, resourceID, err,
			)
		}
	}

	return nil
}

// DeleteFHIRServiceRequest deletes the FHIRServiceRequest identified by the supplied ID
func (fh StoreImpl) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	err := fh.HapiFHIRImpl.DeleteFHIRResource(ctx, serviceRequestResourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			serviceRequestResourceType, err,
		)
	}

	return true, nil
}

// CreateFHIRMedicationStatement creates a new FHIR Medication statement instance
func (fh StoreImpl) CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationStatementResourceType, err)
	}

	resource := &domain.FHIRMedicationStatement{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, medicationStatementResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", medicationStatementResourceType, err)
	}

	output := &domain.FHIRMedicationStatementRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// CreateFHIRMedication creates a new FHIR Medication instance
func (fh StoreImpl) CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationResourceType, err)
	}

	resource := &domain.FHIRMedication{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, medicationResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", medicationResourceType, err)
	}

	output := &domain.FHIRMedicationRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// CreateFHIRMedia creates a FHIR media resource
func (fh StoreImpl) CreateFHIRMedia(ctx context.Context, input domain.FHIRMedia) (*domain.FHIRMedia, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, err
	}

	resource := &domain.FHIRMedia{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, mediaResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", mediaResourceType, err)
	}

	return resource, nil
}

// SearchFHIRMedicationStatement used to search for a fhir medication statement
func (fh StoreImpl) SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
	output := domain.FHIRMedicationStatementRelayConnection{}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", medicationStatementResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRMedicationStatement

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", medicationStatementResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRMedicationStatementRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateFHIRPatient creates a patient on FHIR
func (fh StoreImpl) CreateFHIRPatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", patientResourceType, err)
	}

	resource := &domain.FHIRPatient{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, patientResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", patientResourceType, err)
	}

	output := &domain.PatientPayload{
		PatientRecord: resource,
	}

	return output, nil
}

// PatchFHIRPatient is used to patch a patient resource
func (fh StoreImpl) PatchFHIRPatient(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error) {
	resource := &domain.FHIRPatient{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", patientResourceType, err)
	}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, patientResourceType, id, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// PatchFHIREpisodeOfCare patches a FHIR episode of care
func (fh StoreImpl) PatchFHIREpisodeOfCare(ctx context.Context, id string, input domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error) {
	resource := &domain.FHIREpisodeOfCare{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", episodeOfCareResourceType, err)
	}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, episodeOfCareResourceType, id, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// UpdateFHIREpisodeOfCare updates a fhir episode of care
func (fh StoreImpl) UpdateFHIREpisodeOfCare(ctx context.Context, fhirResourceID string, input domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error) {
	if fhirResourceID == "" {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", episodeOfCareResourceType, err)
	}

	resource := &domain.FHIREpisodeOfCare{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, episodeOfCareResourceType, fhirResourceID, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// SearchFHIRPatient searches for a FHIR patient
func (fh StoreImpl) SearchFHIRPatient(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error) {
	params := map[string]interface{}{
		"_content": searchParams,
		"_total":   "accurate",
	}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", patientResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	output := domain.PatientConnection{}

	for _, resource := range resources.Entry {
		var patient domain.FHIRPatient

		err := mapstructure.Decode(resource, &patient)
		if err != nil {
			return nil, fmt.Errorf("%s, error:%w", internalError, err)
		}

		output.Edges = append(output.Edges, &domain.PatientEdge{
			Node: &patient,
		})
	}

	return &output, nil
}

// GetFHIRComposition retrieves instances of FHIRComposition by ID
func (fh StoreImpl) GetFHIRComposition(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
	resource := &domain.FHIRComposition{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, compositionResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", compositionResourceType, id, err)
	}

	payload := &domain.FHIRCompositionRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// GetFHIRTask is used to retrieve a task given its ID
func (fh StoreImpl) GetFHIRTask(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
	resource := &domain.FHIRTask{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, taskResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", taskResourceType, id, err)
	}

	payload := &domain.FHIRTaskRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// PatchFHIRComposition is used to patch a composition resource
func (fh StoreImpl) PatchFHIRComposition(ctx context.Context, id string, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", compositionResourceType, err)
	}

	resource := &domain.FHIRComposition{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, compositionResourceType, id, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetFHIRObservation retrieves instances of FHIRObservation by ID
func (fh StoreImpl) GetFHIRObservation(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
	resource := &domain.FHIRObservation{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, observationResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", observationResourceType, id, err)
	}

	payload := &domain.FHIRObservationRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// PatchFHIRObservation is used to patch an observation resource
func (fh StoreImpl) PatchFHIRObservation(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", observationResourceType, err)
	}

	resource := &domain.FHIRObservation{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, observationResourceType, id, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// PatchFHIRMedicationRequest is used to patch a medication resource
func (fh StoreImpl) PatchFHIRMedicationRequest(ctx context.Context, id string, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequest, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationRequestResourceType, err)
	}

	resource := &domain.FHIRMedicationRequest{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, medicationRequestResourceType, id, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// ListFHIRQuestionnaire is used to list questionnaire resource using the name or the title of the resource.
func (fh StoreImpl) ListFHIRQuestionnaire(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error) {
	results := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", questionnaireResourceType, params, tenant, results)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range results.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	questionnaireOutput := domain.PagedFHIRQuestionnaires{
		Questionnaires:  []domain.FHIRQuestionnaire{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *results.Total,
	}

	for _, result := range results.Entry {
		var questionnaire domain.FHIRQuestionnaire

		resourceBytes, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBytes, &questionnaire)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		questionnaireOutput.Questionnaires = append(questionnaireOutput.Questionnaires, questionnaire)
	}

	return &questionnaireOutput, nil
}

// CreateFHIRQuestionnaire is used to create a FHIR Questionnaire resource
func (fh StoreImpl) CreateFHIRQuestionnaire(ctx context.Context, input *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", questionnaireResourceType, err)
	}

	resource := &domain.FHIRQuestionnaire{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, questionnaireResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", questionnaireResourceType, err)
	}

	return resource, nil
}

// CreateFHIRConsent creates a FHIRConsent instance
func (fh StoreImpl) CreateFHIRConsent(ctx context.Context, input domain.FHIRConsent) (*domain.FHIRConsent, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", consentResourceType, err)
	}

	resource := &domain.FHIRConsent{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, consentResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", consentResourceType, err)
	}

	return resource, nil
}

// CreateFHIRQuestionnaireResponse is used to create a FHIR Questionnaire response resource
func (fh StoreImpl) CreateFHIRQuestionnaireResponse(ctx context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", questionnaireResponseResourceType, err)
	}

	resource := &domain.FHIRQuestionnaireResponse{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, questionnaireResponseResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", questionnaireResponseResourceType, err)
	}

	return resource, nil
}

// CreateFHIRRiskAssessment creates a RiskAssessment on FHIR
// The RiskAssessment resource represents an assessment of the likely outcome(s) for a patient's health over
// a period of time, considering various factors.
func (fh StoreImpl) CreateFHIRRiskAssessment(ctx context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", riskAssessmentResourceType, err)
	}

	resource := &domain.FHIRRiskAssessment{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, riskAssessmentResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", riskAssessmentResourceType, err)
	}

	output := &domain.FHIRRiskAssessmentRelayPayload{
		Resource: resource,
	}

	cacheKey := fmt.Sprintf("%s:%s", riskAssessmentResourceType, *output.Resource.ID)

	err = fh.SetCache(ctx, cacheKey, output, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SearchFHIRRiskAssessment searches for a fhir risk assessment
func (fh StoreImpl) SearchFHIRRiskAssessment(ctx context.Context, bundleID string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRRiskAssessment, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, bundleID, riskAssessmentResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	riskAssessmentOutput := domain.PagedFHIRRiskAssessment{
		BundleID:        *resources.ID,
		RiskAssessment:  []domain.FHIRRiskAssessment{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRRiskAssessment

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", riskAssessmentResourceType, err)
		}

		riskAssessmentOutput.RiskAssessment = append(riskAssessmentOutput.RiskAssessment, resource)
	}

	return &riskAssessmentOutput, nil
}

// SearchFHIRTask is used to search for available tasks
func (fh StoreImpl) SearchFHIRTask(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRTask, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", taskResourceType, params, tenant, resources)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	taskOutput := domain.PagedFHIRTask{
		Tasks:           []domain.FHIRTask{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      resources.Total,
	}

	for _, result := range resources.Entry {
		var resource domain.FHIRTask

		resourceBs, err := json.Marshal(result.Resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", taskResourceType, err)
		}

		taskOutput.Tasks = append(taskOutput.Tasks, resource)
	}

	return &taskOutput, nil
}

// GetFHIRQuestionnaire retrieves instances of FHIRQuestionnaire by ID
func (fh StoreImpl) GetFHIRQuestionnaire(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
	resource := &domain.FHIRQuestionnaire{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, questionnaireResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", questionnaireResourceType, id, err)
	}

	payload := &domain.FHIRQuestionnaireRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// GetFHIRQuestionnaireResponse retrieves an instance of FHIRQuestionnaireResponse by ID
func (fh StoreImpl) GetFHIRQuestionnaireResponse(ctx context.Context, id string) (*domain.FHIRQuestionnaireResponseRelayPayload, error) {
	resource := &domain.FHIRQuestionnaireResponse{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, questionnaireResponseResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", questionnaireResponseResourceType, id, err)
	}

	payload := &domain.FHIRQuestionnaireResponseRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// CreateFHIRDiagnosticReport is used to create a diagnostic report resource for a patient
func (fh StoreImpl) CreateFHIRDiagnosticReport(ctx context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", diagnosticReportResourceType, err)
	}

	resource := &domain.FHIRDiagnosticReport{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, diagnosticReportResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", diagnosticReportResourceType, err)
	}

	return resource, nil
}

// GetFHIRPatientEverything is used to retrieve all patient related information
func (fh StoreImpl) GetFHIRPatientEverything(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error) {
	if id == "" {
		return nil, fmt.Errorf("patient ID cannot be empty")
	}

	results := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.GetPatientEverything(ctx, id, params, results)
	if err != nil {
		return nil, fmt.Errorf("unable to get patient's compartment for ID %s: %w", id, err)
	}

	paginationInfo := extractPaginationInfo(results.Link)

	totalCount := 0
	if results.Total != nil {
		totalCount = *results.Total
	}

	response := domain.PagedFHIRResource{
		Resources:       []map[string]interface{}{},
		HasNextPage:     paginationInfo.hasNextPage,
		NextCursor:      paginationInfo.nextCursor,
		HasPreviousPage: paginationInfo.hasPreviousPage,
		PreviousCursor:  paginationInfo.previousCursor,
		TotalCount:      totalCount,
		BundleID:        "",
	}

	if results.ID != nil {
		response.BundleID = *results.ID
	}

	if results.Entry != nil {
		for _, entry := range results.Entry {
			if len(entry.Resource) == 0 {
				continue
			}

			var resourceMap map[string]interface{}
			if err := json.Unmarshal(entry.Resource, &resourceMap); err != nil {
				continue
			}

			response.Resources = append(response.Resources, resourceMap)
		}
	}

	return &response, nil
}

// paginationInfo holds extracted pagination details
type paginationInfo struct {
	hasNextPage     bool
	nextCursor      string
	hasPreviousPage bool
	previousCursor  string
}

// extractPaginationInfo extracts pagination information from FHIR bundle links
func extractPaginationInfo(links []hapifhirmodels.BundleLink) paginationInfo {
	info := paginationInfo{}

	if links == nil {
		return info
	}

	for _, link := range links {
		switch link.Relation {
		case "next":
			info.hasNextPage = true
			info.nextCursor = link.URL
		case "previous", "prev":
			info.hasPreviousPage = true
			info.previousCursor = link.URL
		}
	}

	return info
}

// GetFHIRServiceRequest retrieves a FHIR service request using its primary ID
func (fh StoreImpl) GetFHIRServiceRequest(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
	resource := &domain.FHIRServiceRequest{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, serviceRequestResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", serviceRequestResourceType, id, err)
	}

	payload := &domain.FHIRServiceRequestRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// CreateFHIRSubscription is responsible for creating a subscription resource in FHIR repository
func (fh StoreImpl) CreateFHIRSubscription(ctx context.Context, subscription *domain.FHIRSubscriptionInput) (*domain.FHIRSubscription, error) {
	payload, err := converterandformatter.StructToMap(subscription)
	if err != nil {
		return nil, fmt.Errorf("unable to convert subscription input into a map: %w", err)
	}

	fhirSubscription := &domain.FHIRSubscription{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, subscriptionResourceType, payload, fhirSubscription)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to create episode of care resource: %w", err)
	}

	return fhirSubscription, nil
}

// CreateFHIRDocumentReference method is used to create a document reference resource that provides a reference to a document of any kind for any purpose.
func (fh StoreImpl) CreateFHIRDocumentReference(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", documentReferenceResourceType, err)
	}

	resource := &domain.FHIRDocumentReference{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, documentReferenceResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", documentReferenceResourceType, err)
	}

	return resource, nil
}

// PatchFHIRServiceRequest is used to update the specified fhir service request resource
func (fh StoreImpl) PatchFHIRServiceRequest(ctx context.Context, id string, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", serviceRequestResourceType, err)
	}

	resource := &domain.FHIRServiceRequest{}

	err = fh.HapiFHIRImpl.FHIRPathPatch(ctx, serviceRequestResourceType, id, payload, resource)
	if err != nil {
		return nil, err
	}

	return &domain.FHIRServiceRequestRelayPayload{
		Resource: resource,
	}, nil
}

// SearchFHIRDocumentReference is used to search for FHIR document reference resource using the client provided parameters. Some of these parameters include (but not limited) to `categories`, 'subject` etc.`
func (fh StoreImpl) SearchFHIRDocumentReference(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
	documentReferences := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", documentReferenceResourceType, searchParams, tenant, documentReferences)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range documentReferences.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	documentReferenceOutput := domain.PagedFHIRDocumentReference{
		DocumentReferences: []domain.FHIRDocumentReference{},
		HasNextPage:        hasNextPage,
		NextCursor:         nextCursor,
		HasPreviousPage:    hasPreviousPage,
		PreviousCursor:     previousCursor,
		TotalCount:         *documentReferences.Total,
	}

	for _, reference := range documentReferences.Entry {
		var documentReference domain.FHIRDocumentReference

		resourceBs, err := json.Marshal(reference)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &documentReference)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		documentReferenceOutput.DocumentReferences = append(documentReferenceOutput.DocumentReferences, documentReference)
	}

	return &documentReferenceOutput, nil
}

// CreateFHIRAppointment is used to create a new FHIR appointment resource
func (fh StoreImpl) CreateFHIRAppointment(ctx context.Context, input *domain.FHIRAppointmentInput) (*domain.FHIRAppointment, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", appointmentResourceType, err)
	}

	resource := &domain.FHIRAppointment{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, appointmentResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", appointmentResourceType, err)
	}

	return resource, nil
}

// CreateFHIRTask creates a task resource
func (fh StoreImpl) CreateFHIRTask(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", taskResourceType, err)
	}

	resource := &domain.FHIRTask{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, taskResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", taskResourceType, err)
	}

	return resource, nil
}

// CreateFHIRLocation is use to create a location resource
func (fh StoreImpl) CreateFHIRLocation(ctx context.Context, input *domain.FHIRLocation) (*domain.FHIRLocation, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", locationResourceType, err)
	}

	resource := &domain.FHIRLocation{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, locationResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", locationResourceType, err)
	}

	return resource, nil
}

// GetFHIRLocation implements repository.FHIR.GetFHIRLocation method used to retrieve a location given its logical id
func (fh *StoreImpl) GetFHIRLocation(ctx context.Context, id string) (*domain.FHIRLocation, error) {
	resource := &domain.FHIRLocation{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, locationResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s resource: %w", locationResourceType, err)
	}

	return resource, nil
}

func (fh *StoreImpl) FetchMedicationByID(ctx context.Context, id string) (*domain.FHIRMedication, error) {
	resource := &domain.FHIRMedication{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, medicationResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s resource: %w", medicationResourceType, err)
	}

	return resource, nil
}

// CreateFHIRPractitionerRole implements repository.FHIR.CreateFHIRPractitionerRole method used to create practitioners role
func (fh *StoreImpl) CreateFHIRPractitionerRole(ctx context.Context, input *domain.FHIRPractitionerRole) (*domain.FHIRPractitionerRole, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", practitionerRoleResourceType, err)
	}

	resource := &domain.FHIRPractitionerRole{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, practitionerRoleResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", practitionerRoleResourceType, err)
	}

	return resource, nil
}

// SearchFHIRPractitionerRoles is used to search for available practitioners roles
func (fh *StoreImpl) SearchFHIRPractitionerRoles(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRPractitionerRole, error) {
	practitionerRoles := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", practitionerRoleResourceType, searchParams, tenant, practitionerRoles)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range practitionerRoles.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	practitionerRoleOutput := domain.PagedFHIRPractitionerRole{
		PractitionerRoles: []*domain.FHIRPractitionerRole{},
		HasNextPage:       hasNextPage,
		NextCursor:        nextCursor,
		HasPreviousPage:   hasPreviousPage,
		PreviousCursor:    previousCursor,
		TotalCount:        *practitionerRoles.Total,
	}

	for _, role := range practitionerRoles.Entry {
		var practitionerRole domain.FHIRPractitionerRole

		resourceBs, err := json.Marshal(&role)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &practitionerRole)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		practitionerRoleOutput.PractitionerRoles = append(practitionerRoleOutput.PractitionerRoles, &practitionerRole)
	}

	return &practitionerRoleOutput, nil
}

// GetFHIRPractitionerRole is used to retrieve practitioner role information
func (fh *StoreImpl) GetFHIRPractitionerRole(ctx context.Context, id string) (*domain.FHIRPractitionerRoleRelayPayload, error) {
	resource := &domain.FHIRPractitionerRole{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, practitionerRoleResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s resource: %w", practitionerRoleResourceType, err)
	}

	return &domain.FHIRPractitionerRoleRelayPayload{
		Resource: resource,
	}, nil
}

// SearchFHIRResource implements repository.FHIR.
func (fh *StoreImpl) SearchFHIRResource(ctx context.Context, bundleID, resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, bundleID, resourceType, params, tenant, resources)
	if err != nil {
		return nil, fmt.Errorf("failed to search resource: %w", err)
	}

	var previousCursor, nextCursor string

	var hasNextPage, hasPreviousPage bool

	for _, link := range resources.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL

		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	pagedResults := domain.PagedFHIRResource{
		Resources:       []map[string]interface{}{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *resources.Total,
	}

	for _, resource := range resources.Entry {
		var resourceMap map[string]interface{}

		jsonData, err := resource.Resource.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("server error: serialization to JSON bytes failed: %w", err)
		}

		if err = json.Unmarshal(jsonData, &resourceMap); err != nil {
			return nil, fmt.Errorf("deserialization to a map failed: %w", err)
		}

		pagedResults.Resources = append(pagedResults.Resources, resourceMap)
	}

	return &pagedResults, nil
}

// CreateFHIRPractitioner creates a Practitioner resource
func (fh StoreImpl) CreateFHIRPractioner(ctx context.Context, input *domain.FHIRPractitioner) (*domain.FHIRPractitioner, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("failed to turn %s into a map: %w", practitionerResourceType, err)
	}

	resource := &domain.FHIRPractitioner{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, practitionerResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return resource, nil
}

// GetFHIRPractitioner fetches a Practitioner from the server.
func (fh StoreImpl) GetFHIRPractitioner(ctx context.Context, id string) (*domain.FHIRPractitioner, error) {
	resource := &domain.FHIRPractitioner{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, practitionerResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("could not get Practitioner: %w", err)
	}

	return resource, nil
}

// SearchFHIRPractitioner searches for Practitioners
func (fh StoreImpl) SearchFHIRPractitioner(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PractitionerConnection, error) {
	searchParam := map[string]interface{}{
		"_content": params,
	}

	resources := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", practitionerResourceType, searchParam, tenant, resources)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources: %w", err)
	}

	output := domain.PractitionerConnection{}

	for _, resource := range resources.Entry {
		var patient domain.FHIRPatient
		// This will be reviewed.
		practitioner := domain.FHIRPractitioner{
			ID:        patient.ID,
			Text:      patient.Text,
			Active:    patient.Active,
			BirthDate: patient.BirthDate,
		}

		practitioner.Telecom = append(practitioner.Telecom, patient.Telecom...)
		practitioner.Identifier = append(practitioner.Identifier, patient.Identifier...)
		practitioner.Name = append(practitioner.Name, patient.Name...)

		err := mapstructure.Decode(resource, &practitioner)
		if err != nil {
			return nil, fmt.Errorf("%s, erro:%w", internalError, err)
		}

		output.Edges = append(output.Edges, &domain.PatientEdge{
			Node: &patient,
		})
	}

	return &output, nil
}

// CreateFHIRSubstance implements repository.FHIR.CreateFHIRSubstance method that creates a substance
func (fh *StoreImpl) CreateFHIRSubstance(ctx context.Context, input *domain.FHIRSubstance) (*domain.FHIRSubstance, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", substanceResourceType, err)
	}

	resource := &domain.FHIRSubstance{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, substanceResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", substanceResourceType, err)
	}

	return resource, nil
}

// GetFHIRSubstance is used to retrieve a substance records
func (fh *StoreImpl) GetFHIRSubstance(ctx context.Context, id string) (*domain.FHIRSubstance, error) {
	resource := &domain.FHIRSubstance{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, substanceResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s resource: %w", substanceResourceType, err)
	}

	return resource, nil
}

// SearchFHIRSubstance is used to search for substances
func (fh *StoreImpl) SearchFHIRSubstance(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRSubstance, error) {
	substances := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", substanceResourceType, searchParams, tenant, substances)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range substances.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	substanceOutput := domain.PagedFHIRSubstance{
		Substance:       []*domain.FHIRSubstance{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *substances.Total,
	}

	for _, entry := range substances.Entry {
		var substance domain.FHIRSubstance

		resourceBs, err := json.Marshal(&entry)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &substance)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		substanceOutput.Substance = append(substanceOutput.Substance, &substance)
	}

	return &substanceOutput, nil
}

// CreateFHIRProcedure is used to create procedure resource
func (fh *StoreImpl) CreateFHIRProcedure(ctx context.Context, input *domain.FHIRProcedure) (*domain.FHIRProcedure, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", procedureResourceType, err)
	}

	resource := &domain.FHIRProcedure{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, procedureResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s: %w", procedureResourceType, err)
	}

	return resource, nil
}

// GetFHIRProcedure is used to fetch procedure
func (fh *StoreImpl) GetFHIRProcedure(ctx context.Context, id string) (*domain.FHIRProcedure, error) {
	resource := &domain.FHIRProcedure{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, procedureResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("could not fetch %s resource: %w", procedureResourceType, err)
	}

	return resource, nil
}

// SearchFHIRProcedure searches for Procedures
func (fh *StoreImpl) SearchFHIRProcedure(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRProcedure, error) {
	procedure := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", procedureResourceType, searchParams, tenant, procedure)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range procedure.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	output := domain.PagedFHIRProcedure{
		Procedure:       []*domain.FHIRProcedure{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *procedure.Total,
	}

	for _, entry := range procedure.Entry {
		var procedure domain.FHIRProcedure

		procedureEntry, err := json.Marshal(&entry)
		if err != nil {
			return nil, fmt.Errorf("could not marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(procedureEntry, &procedure)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal resource: %w", err)
		}

		output.Procedure = append(output.Procedure, &procedure)
	}

	return &output, nil
}

// CreateFHIRMedicationDispense is used to create medication dispense
func (fh *StoreImpl) CreateFHIRMedicationDispense(ctx context.Context, input *domain.FHIRMedicationDispense) (*domain.FHIRMedicationDispense, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationDispenseResourceType, err)
	}

	resource := &domain.FHIRMedicationDispense{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, medicationDispenseResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s: %w", medicationDispenseResourceType, err)
	}

	return resource, nil
}

// SearchFHIRMedicationDispense is used to search for medication dispense record
func (fh *StoreImpl) SearchFHIRMedicationDispense(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedicationDispense, error) {
	medicationDispenses := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", medicationDispenseResourceType, params, tenant, medicationDispenses)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range medicationDispenses.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	output := domain.PagedFHIRMedicationDispense{
		MedicationDispense: []*domain.FHIRMedicationDispense{},
		HasNextPage:        hasNextPage,
		NextCursor:         nextCursor,
		HasPreviousPage:    hasPreviousPage,
		PreviousCursor:     previousCursor,
		TotalCount:         *medicationDispenses.Total,
	}

	for _, entry := range medicationDispenses.Entry {
		var medicationDispense domain.FHIRMedicationDispense

		medicationDispenseEntry, err := json.Marshal(&entry)
		if err != nil {
			return nil, fmt.Errorf("could not marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(medicationDispenseEntry, &medicationDispense)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal resource: %w", err)
		}

		output.MedicationDispense = append(output.MedicationDispense, &medicationDispense)
	}

	return &output, nil
}

func (fh *StoreImpl) GetFHIRMedicationDispense(ctx context.Context, id string) (*domain.FHIRMedicationDispense, error) {
	resource := &domain.FHIRMedicationDispense{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, medicationDispenseResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("could not fetch %s resource: %w", medicationDispenseResourceType, err)
	}

	return resource, nil
}

func (fh *StoreImpl) ValidateResource(ctx context.Context, resourceType string, params map[string]interface{}) error {
	return fh.HapiFHIRImpl.ValidateResource(ctx, resourceType, params)
}

// CreateFHIRPlanDefinition is used to create plan definition
func (fh *StoreImpl) CreateFHIRPlanDefinition(ctx context.Context, input *domain.FHIRPlanDefinition) (*domain.FHIRPlanDefinition, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", planDefinitionResourceType, err)
	}

	resource := &domain.FHIRPlanDefinition{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, planDefinitionResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s: %w", planDefinitionResourceType, err)
	}

	return resource, nil
}

// CreateFHIRActivityDefinition is used to create an activity definition
func (fh *StoreImpl) CreateFHIRActivityDefinition(ctx context.Context, input *domain.FHIRActivityDefinition) (*domain.FHIRActivityDefinition, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", activityDefinitionResourceType, err)
	}

	resource := &domain.FHIRActivityDefinition{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, activityDefinitionResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s: %w", activityDefinitionResourceType, err)
	}

	return resource, nil
}

// CreateFHIRCarePlan is used to store information describing the intention of how one or more practitioners intend to deliver care for a particular patient, group or community for a period of time, possibly limited to care for a specific condition or set of conditions.
func (fh *StoreImpl) CreateFHIRCarePlan(ctx context.Context, input *domain.FHIRCarePlan) (*domain.FHIRCarePlan, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", carePlanResourceType, err)
	}

	resource := &domain.FHIRCarePlan{}

	err = fh.HapiFHIRImpl.CreateFHIRResource(ctx, carePlanResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s: %w", carePlanResourceType, err)
	}

	return resource, nil
}

func (fh *StoreImpl) SearchFHIRPlanDefinition(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRPlanDefinition, error) {
	planDefinitions := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", planDefinitionResourceType, params, dto.TenantIdentifiers{}, planDefinitions)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range planDefinitions.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	output := domain.PagedFHIRPlanDefinition{
		PlanDefinition:  []domain.FHIRPlanDefinition{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *planDefinitions.Total,
	}

	for _, entry := range planDefinitions.Entry {
		var planDefinition domain.FHIRPlanDefinition

		planDefinitionEntry, err := json.Marshal(&entry.Resource)
		if err != nil {
			return nil, fmt.Errorf("could not marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(planDefinitionEntry, &planDefinition)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal resource: %w", err)
		}

		output.PlanDefinition = append(output.PlanDefinition, planDefinition)
	}

	return &output, nil
}

func (fh *StoreImpl) SearchFHIRCarePlan(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error) {
	carePlan := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.SearchFHIRResource(ctx, "", carePlanResourceType, params, dto.TenantIdentifiers{}, carePlan)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPreviousPage bool

	var nextCursor, previousCursor string

	for _, link := range carePlan.Link {
		switch link.Relation {
		case "next":
			hasNextPage = true
			nextCursor = link.URL
		case "previous":
			hasPreviousPage = true
			previousCursor = link.URL
		}
	}

	output := domain.PagedFHIRCarePlan{
		CarePlan:        []domain.FHIRCarePlan{},
		HasNextPage:     hasNextPage,
		NextCursor:      nextCursor,
		HasPreviousPage: hasPreviousPage,
		PreviousCursor:  previousCursor,
		TotalCount:      *carePlan.Total,
	}

	for _, entry := range carePlan.Entry {
		var carePlan domain.FHIRCarePlan

		carePlanEntry, err := json.Marshal(&entry.Resource)
		if err != nil {
			return nil, fmt.Errorf("could not marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(carePlanEntry, &carePlan)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal resource: %w", err)
		}

		output.CarePlan = append(output.CarePlan, carePlan)
	}

	return &output, nil
}

func (fh *StoreImpl) FetchPlanDefinitionByID(ctx context.Context, id string) (*domain.FHIRPlanDefinition, error) {
	resource := &domain.FHIRPlanDefinition{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, planDefinitionResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("could not fetch %s resource: %w", planDefinitionResourceType, err)
	}

	return resource, nil
}

// PostFHIRBundle posts a FHIR bundle to hapi fhir.
func (fh *StoreImpl) PostFHIRBundle(ctx context.Context, payload *hapifhirmodels.Bundle) (*hapifhirmodels.Bundle, error) {
	resource := &hapifhirmodels.Bundle{}

	err := fh.HapiFHIRImpl.PostFHIRBundle(ctx, payload, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (fh *StoreImpl) GetFHIRRiskAssessment(ctx context.Context, id string) (*domain.FHIRRiskAssessmentRelayPayload, error) {
	resource := &domain.FHIRRiskAssessment{}

	err := fh.HapiFHIRImpl.GetFHIRResource(ctx, riskAssessmentResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", riskAssessmentResourceType, id, err)
	}

	payload := &domain.FHIRRiskAssessmentRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

func (fh *StoreImpl) PutFHIRResource(ctx context.Context, resourceType, resourceID string, payload map[string]any, resource any, useCREnabledServer bool) error {
	return fh.HapiFHIRImpl.PutFHIRResource(ctx, resourceType, resourceID, payload, resource, useCREnabledServer)
}
