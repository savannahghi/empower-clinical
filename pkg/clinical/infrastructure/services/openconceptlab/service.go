// Package openconceptlab provides APIs to interact with an OpenConceptLab API
// server
package openconceptlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/serverutils"
)

// constants used to configure the OCL service
const (
	OCLAPIURLEnvVarName  = "OPENCONCEPTLAB_API_URL"
	OCLTokenEnvVarName   = "OPENCONCEPTLAB_TOKEN"
	OCLAPITimeoutSeconds = 30
)

// NewServiceOCL creates a new open conceptlab ServiceOCL
func NewServiceOCL() *ServiceOCL {
	baseURL := serverutils.MustGetEnvVar(OCLAPIURLEnvVarName)
	token := serverutils.MustGetEnvVar(OCLTokenEnvVarName)
	header := fmt.Sprintf("Token %s", token)

	srv := &ServiceOCL{
		baseURL: baseURL,
		header:  header,
	}
	srv.enforcePreconditions()

	return srv
}

// ServiceOCL is an OpenConceptLab service
type ServiceOCL struct {
	baseURL string
	header  string
}

func (s ServiceOCL) enforcePreconditions() {
	if s.baseURL == "" {
		log.Panicf("Open Concept Lab API Base URL not set in service")
	}

	if s.header == "" {
		log.Panicf("Open Concept Lab API Token header not set in service")
	}
}

// MakeRequest composes an authenticated OCL request that has the correct content type
func (s ServiceOCL) MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/?%s", s.baseURL, path, params.Encode())

	req, reqErr := http.NewRequest(method, url, body)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.header)

	httpClient := &http.Client{Timeout: time.Second * OCLAPITimeoutSeconds}

	resp, respErr := httpClient.Do(req)
	if respErr != nil {
		return nil, respErr
	}

	return resp, nil
}

// GetConcept retrieves a single concept from OpenConceptLab.
// The URL that is composed follows this pattern: GET /orgs/:org/sources/:source/[:sourceVersion/]concepts/:concept/
// e.g GET /orgs/WHO/sources/ICD-10-2010/concepts/A15.1/?includeInverseMappings=true
func (s ServiceOCL) GetConcept(
	_ context.Context, org string, source string, concept string,
	includeMappings bool, includeInverseMappings bool) (*domain.Concept, error) {
	s.enforcePreconditions()

	path := fmt.Sprintf("orgs/%s/sources/%s/concepts/%s", org, source, concept)

	params := url.Values{}
	params.Add("includeMappings", strconv.FormatBool(includeMappings))
	params.Add("includeMappings", strconv.FormatBool(includeInverseMappings))

	resp, err := s.MakeRequest("GET", path, params, nil)
	if err != nil {
		return nil, fmt.Errorf("OCL API request error: %w", err)
	}

	defer resp.Body.Close()

	output := make(map[string]interface{})

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read OCL API response body: %w", err)
	}

	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to marshal OCL get concept response %s to JSON: %w", string(data), err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"OCL returned %d for %s concept %s: %s", resp.StatusCode, source, concept, string(data))
	}

	if id, ok := output["id"].(string); !ok || id == "" {
		return nil, fmt.Errorf("failed to get %v concept with id %v", source, concept)
	}

	var terminologyConcept *domain.Concept

	err = mapstructure.Decode(output, &terminologyConcept)
	if err != nil {
		return nil, err
	}

	return terminologyConcept, nil
}

// ListConcepts searches for matching concepts on OpenConceptLab
// The URL that is composed follows this pattern: GET /orgs/:org/sources/:source/[:sourceVersion/]concepts/
// e.g GET /orgs/PEPFAR-Test7/sources/MER/concepts/?conceptClass="Symptom"+OR+"Diagnosis"
func (s ServiceOCL) ListConcepts(
	_ context.Context, org []string, source []string, verbose bool, q *string,
	sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
	locale *string, includeRetired *bool,
	includeMappings *bool, includeInverseMappings *bool, paginationInput *dto.Pagination) (*domain.ConceptPage, error) {
	s.enforcePreconditions()

	path := fmt.Sprintf("orgs/%s/sources/%s/concepts", strings.Join(org, ","), strings.Join(source, ","))

	params := url.Values{}
	params.Add("verbose", strconv.FormatBool(verbose))

	if paginationInput.After != "" {
		params.Add("page", paginationInput.After)
	} else {
		params.Add("page", "1")
	}

	params.Add("limit", strconv.Itoa(*paginationInput.First))

	if q != nil {
		params.Add("q", *q)
	}

	if sortAsc != nil {
		params.Add("sortAsc", *sortAsc)
	}

	if sortDesc != nil {
		params.Add("sortDesc", *sortDesc)
	}

	if conceptClass != nil {
		params.Add("conceptClass", *conceptClass)
	}

	if dataType != nil {
		params.Add("dataType", *dataType)
	}

	if locale != nil {
		params.Add("locale", *locale)
	}

	if includeRetired != nil {
		if *includeRetired {
			params.Add("includeRetired", "1")
		} else {
			params.Add("includeRetired", "0")
		}
	}

	if includeMappings != nil && *includeMappings {
		params.Add("includeMappings", "true")
	}

	if includeInverseMappings != nil && *includeInverseMappings {
		params.Add("includeReverseMappings", "true")
	}

	resp, err := s.MakeRequest("GET", path, params, nil)
	if err != nil {
		return nil, fmt.Errorf("OCL API request error: %w", err)
	}

	defer resp.Body.Close()

	terminologyConcepts := []map[string]interface{}{}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read OCL API response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OCL returned %d listing concepts: %s", resp.StatusCode, string(data))
	}

	err = json.Unmarshal(data, &terminologyConcepts)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to marshal OCL get concept response %s to JSON: %w", string(data), err)
	}

	totalCount, err := strconv.Atoi(resp.Header.Get("num_found"))
	if err != nil {
		return nil, err
	}

	conceptsPage := &domain.ConceptPage{
		Count: totalCount,
	}

	nextConcept := resp.Header.Get("next")

	if nextConcept != "" {
		params, err := url.ParseQuery(nextConcept)
		if err != nil {
			return nil, fmt.Errorf("server unable to parse url params: %w", err)
		}

		cursor := params["page"][0]
		conceptsPage.Next = &cursor
	}

	previousConcept := resp.Header.Get("previous")

	if previousConcept != "" {
		params, err := url.ParseQuery(previousConcept)
		if err != nil {
			return nil, fmt.Errorf("server unable to parse url params: %w", err)
		}

		cursor := params["page"][0]
		conceptsPage.Previous = &cursor
	}

	for _, terminologyConcept := range terminologyConcepts {
		var concept *domain.Concept

		err := mapstructure.Decode(terminologyConcept, &concept)
		if err != nil {
			return nil, err
		}

		conceptsPage.Results = append(conceptsPage.Results, concept)
	}

	return conceptsPage, nil
}
