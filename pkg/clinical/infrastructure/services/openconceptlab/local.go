package openconceptlab

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// ConceptsFileEnvVarName points at a JSON file whose concepts are merged over the
// built-in seed, keyed by "<org>/<source>". Use it to add terminology the seed does
// not carry, such as an ICD-11 allergen list.
const ConceptsFileEnvVarName = "OPENCONCEPTLAB_CONCEPTS_FILE"

//go:embed concepts.json
var seedConcepts []byte

// LocalService resolves concepts from an in-process table instead of a remote
// OpenConceptLab server, so the stack runs without an OCL account.
type LocalService struct {
	concepts map[string][]*domain.Concept
}

// NewLocalService builds a LocalService from the embedded seed, overlaying the file
// named by OPENCONCEPTLAB_CONCEPTS_FILE when it is set.
func NewLocalService() (*LocalService, error) {
	concepts := map[string][]*domain.Concept{}

	err := json.Unmarshal(seedConcepts, &concepts)
	if err != nil {
		return nil, fmt.Errorf("unable to parse the embedded concept seed: %w", err)
	}

	path := os.Getenv(ConceptsFileEnvVarName)
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("unable to read %s %q: %w", ConceptsFileEnvVarName, path, err)
		}

		overlay := map[string][]*domain.Concept{}

		err = json.Unmarshal(data, &overlay)
		if err != nil {
			return nil, fmt.Errorf("unable to parse %q: %w", path, err)
		}

		for key, added := range overlay {
			concepts[key] = append(concepts[key], added...)
		}
	}

	return &LocalService{concepts: concepts}, nil
}

func key(org, source string) string {
	return fmt.Sprintf("%s/%s", org, source)
}

// MakeRequest is not supported: there is no server to call.
func (s LocalService) MakeRequest(_ string, _ string, _ url.Values, _ io.Reader) (*http.Response, error) {
	return nil, fmt.Errorf("the local terminology backend does not make HTTP requests")
}

// GetConcept returns the seeded concept, or an error naming what is missing so the
// caller fails loudly rather than writing an uncoded resource.
func (s LocalService) GetConcept(
	_ context.Context, org string, source string, concept string,
	_ bool, _ bool,
) (*domain.Concept, error) {
	for _, candidate := range s.concepts[key(org, source)] {
		if candidate.ID == concept {
			found := *candidate
			return &found, nil
		}
	}

	return nil, fmt.Errorf(
		"concept %q is not in the local %s terminology; add it via %s",
		concept, key(org, source), ConceptsFileEnvVarName)
}

// ListConcepts does a case-insensitive substring match over the seeded display names
// of every requested source.
func (s LocalService) ListConcepts(
	_ context.Context, org []string, source []string, _ bool, q *string,
	_ *string, _ *string, _ *string, _ *string,
	_ *string, includeRetired *bool,
	_ *bool, _ *bool, paginationInput *dto.Pagination,
) (*domain.ConceptPage, error) {
	var matches []*domain.Concept

	term := ""
	if q != nil {
		term = strings.ToLower(strings.TrimSpace(*q))
	}

	for _, o := range org {
		for _, src := range source {
			for _, candidate := range s.concepts[key(o, src)] {
				if candidate.Retired && (includeRetired == nil || !*includeRetired) {
					continue
				}

				if term != "" && !strings.Contains(strings.ToLower(candidate.DisplayName), term) {
					continue
				}

				matches = append(matches, candidate)
			}
		}
	}

	page := &domain.ConceptPage{Count: len(matches)}

	start := 0
	if paginationInput != nil && paginationInput.After != "" {
		_, err := fmt.Sscanf(paginationInput.After, "%d", &start)
		if err != nil {
			return nil, fmt.Errorf("invalid pagination cursor %q: %w", paginationInput.After, err)
		}
	}

	if start > len(matches) {
		start = len(matches)
	}

	end := len(matches)
	if paginationInput != nil && paginationInput.First != nil && start+*paginationInput.First < end {
		end = start + *paginationInput.First
	}

	if end < len(matches) {
		cursor := fmt.Sprintf("%d", end)
		page.Next = &cursor
	}

	page.Results = matches[start:end]

	return page, nil
}
