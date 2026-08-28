package openconceptlab_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/openconceptlab"
)

func TestService_GetConcept(t *testing.T) {

	type args struct {
		ctx                    context.Context
		org                    string
		source                 string
		concept                string
		includeMappings        bool
		includeInverseMappings bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get concept",
			args: args{
				ctx:                    context.Background(),
				org:                    "CIEL",
				source:                 "CIEL",
				concept:                "1234",
				includeMappings:        false,
				includeInverseMappings: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tt.name == "happy case: get concept" {
				httpmock.RegisterResponder(http.MethodGet, "/orgs/CIEL/sources/CIEL/concepts/1234/",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{
							"concept_class":       "Diagnosis",
							"datatype":            "N/A",
							"display_locale":      "en",
							"display_name":        "Acute Coryza",
							"external_id":         "106AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
							"id":                  "106",
							"is_latest_version":   true,
							"locale":              nil,
							"owner":               "CIEL",
							"owner_type":          "Organization",
							"owner_url":           "/orgs/CIEL/",
							"retired":             false,
							"source":              "CIEL",
							"type":                "Concept",
							"update_comment":      nil,
							"url":                 "/orgs/CIEL/sources/CIEL/concepts/106/",
							"uuid":                "2828492",
							"version":             "2828492",
							"version_created_by":  "ocladmin",
							"version_created_on":  "2023-01-27T13:31:22.143729Z",
							"version_url":         "/orgs/CIEL/sources/CIEL/concepts/106/2828492/",
							"versioned_object_id": 1776,
							"versions_url":        "/orgs/CIEL/sources/CIEL/concepts/106/versions/",
						})

					},
				)
			}

			s := openconceptlab.NewServiceOCL()
			got, err := s.GetConcept(tt.args.ctx, tt.args.org, tt.args.source, tt.args.concept, tt.args.includeMappings, tt.args.includeInverseMappings)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetConcept() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected result not to be nil for %v", tt.name)
				return
			}

		})
	}
}

func toPointer[S any](s S) *S {
	return &s
}

func TestService_ListConcepts(t *testing.T) {
	limit := 10
	next := "page=2"
	previous := "page=1"
	type args struct {
		ctx                    context.Context
		org                    []string
		source                 []string
		verbose                bool
		q                      *string
		sortAsc                *string
		sortDesc               *string
		conceptClass           *string
		dataType               *string
		locale                 *string
		includeRetired         *bool
		includeMappings        *bool
		includeInverseMappings *bool
		paginationInput        *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: list concepts",
			args: args{
				ctx:                    context.Background(),
				org:                    []string{"CIEL", "WHO"},
				source:                 []string{"CIEL", "WHO"},
				verbose:                false,
				q:                      toPointer("panadol"),
				sortAsc:                toPointer("id"),
				sortDesc:               toPointer("owner"),
				conceptClass:           toPointer("Diagnosis"),
				dataType:               toPointer("N/A"),
				locale:                 toPointer("en/us"),
				includeRetired:         toPointer(true),
				includeMappings:        toPointer(true),
				includeInverseMappings: toPointer(true),
				paginationInput: &dto.Pagination{
					First: &limit,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tt.name == "happy case: list concepts" {
				httpmock.RegisterResponder(http.MethodGet, "=~^https://api.openconceptlab.org/orgs/CIEL,WHO/sources/CIEL,WHO/concepts",
					func(req *http.Request) (*http.Response, error) {
						response, err := httpmock.NewJsonResponse(200, []map[string]interface{}{
							{
								"concept_class":       "Diagnosis",
								"datatype":            "N/A",
								"display_locale":      "en",
								"display_name":        "Acute Coryza",
								"external_id":         "106AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
								"id":                  "106",
								"is_latest_version":   true,
								"locale":              nil,
								"owner":               "CIEL",
								"owner_type":          "Organization",
								"owner_url":           "/orgs/CIEL/",
								"retired":             false,
								"source":              "CIEL",
								"type":                "Concept",
								"update_comment":      nil,
								"url":                 "/orgs/CIEL/sources/CIEL/concepts/106/",
								"uuid":                "2828492",
								"version":             "2828492",
								"version_created_by":  "ocladmin",
								"version_created_on":  "2023-01-27T13:31:22.143729Z",
								"version_url":         "/orgs/CIEL/sources/CIEL/concepts/106/2828492/",
								"versioned_object_id": 1776,
								"versions_url":        "/orgs/CIEL/sources/CIEL/concepts/106/versions/",
							},
						})

						response.Header.Set("num_found", "1")
						response.Header.Set("next", next)
						response.Header.Set("previous", previous)

						return response, err
					},
				)
			}

			s := openconceptlab.NewServiceOCL()

			got, err := s.ListConcepts(tt.args.ctx, tt.args.org, tt.args.source, tt.args.verbose, tt.args.q, tt.args.sortAsc, tt.args.sortDesc, tt.args.conceptClass, tt.args.dataType, tt.args.locale, tt.args.includeRetired, tt.args.includeMappings, tt.args.includeInverseMappings, tt.args.paginationInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListConcepts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected result not to be nil for %v", tt.name)
				return
			}
		})
	}
}
