package fhir_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	hapifhirmodels "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	fhir "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"

	fakeCache "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/cache/mock"
	fakeHapiFHIR "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
)

type mockHandler struct {
	cache    *fakeCache.MockCacheService
	hapiFHIR *fakeHapiFHIR.MockHapiFHIRImplementation
}

func TestStoreImpl_SearchFHIRObservation(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"_sort": "-date",
	}
	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}

	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: search observation",
			setup: func(mh *mockHandler) args {
				// build a fake bundle we want the mock to return
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRObservation{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})

				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "sad case: search resource error",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						return errors.New("failed to search fhir resource")
					})

				return args{ctx: ctx, params: params}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)

			mh := mockHandler{
				cache:    mockCache,
				hapiFHIR: mockhapiFHIR,
			}

			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mh)

			got, err := fh.SearchFHIRObservation(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Fatalf("expected a response but got nil")
				}
			}
		})
	}

}

func TestStoreImpl_DeleteFHIRObservation(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: delete observation",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return nil
					})

				return args{ctx: ctx, id: id}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: error deleting resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return errors.New("failed to delete resource")
					})

				return args{ctx: ctx, id: id}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.DeleteFHIRObservation(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.DeleteFHIRObservation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_CreateFHIRObservation(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	input := domain.FHIRObservationInput{
		ID: &ID,
	}
	type args struct {
		ctx   context.Context
		input domain.FHIRObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: create fhir observation",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return errors.New("failed to create resource")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.CreateFHIRObservation(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_PatchFHIRObservation(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.NewString()
	id := ksuid.New().String()
	input := domain.FHIRObservationInput{ID: &UUID}

	type args struct {
		ctx   context.Context
		id    string
		input domain.FHIRObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *domain.FHIRObservationRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case - successfully patch observation",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, id: id, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case - fail to patch observation",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: id, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.PatchFHIRObservation(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.PatchFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRPatient(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: get patient: cache miss",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						return cmd
					})

				return args{ctx: ctx, id: uuid.NewString()}
			},
			wantErr: false,
		},
		{
			name: "happy case: get patient: cache hit",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)

						id := gofakeit.UUID()
						patient := &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &id,
							},
						}

						bs, err := json.Marshal(patient)
						if err != nil {
							cmd.SetErr(err)
						}

						cmd.SetVal(string(bs))
						return cmd
					})

				return args{ctx: ctx, id: uuid.NewString()}
			},
			wantErr: false,
		},
		{
			name: "Sad case: error retrieving fhir resource",
			setup: func(mh *mockHandler) args {
				ctx := context.Background()

				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("failed to get resource")
					})
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})

				return args{ctx: ctx, id: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to set value in cache",
			setup: func(mh *mockHandler) args {
				ctx := context.Background()

				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetErr(errors.New("failed to set record in cache"))
						return cmd
					})
				return args{ctx: ctx, id: uuid.NewString()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.GetFHIRPatient(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRResourceType(t *testing.T) {
	ctx := context.Background()
	results := []map[string]string{
		{"Patient": uuid.NewString()},
	}
	type args struct {
		ctx     context.Context
		results []map[string]string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: delete resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return nil
					})

				return args{ctx: ctx, results: results}
			},
			wantErr: false,
		},
		{
			name: "sad case: delete resource error",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return fmt.Errorf("failed to delete resource")
					})

				return args{ctx: ctx, results: results}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			if err := fh.DeleteFHIRResourceType(args.ctx, args.results); (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRResourceType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: delete service request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return nil
					})

				return args{ctx: ctx, id: id}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: delete resource error",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return fmt.Errorf("failed to delete resource")
					})

				return args{ctx: ctx, id: id}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.DeleteFHIRServiceRequest(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.DeleteFHIRServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_CreateFHIRMedicationStatement(t *testing.T) {
	ctx := context.Background()
	input := domain.FHIRMedicationStatementInput{
		Category: &domain.FHIRCodeableConceptInput{
			Text: "dawa",
		},
	}
	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationStatementInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: create medication statement",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create resource")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.CreateFHIRMedicationStatement(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedicationStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRMedication(t *testing.T) {
	ctx := context.Background()
	input := domain.FHIRMedicationInput{
		Code: &domain.FHIRCodeableConceptInput{
			Text: "ARV",
		},
	}
	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: create medication",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create resource")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.CreateFHIRMedication(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRPatient(t *testing.T) {
	ctx := context.Background()
	input := domain.FHIRPatientInput{
		MaritalStatus: &domain.FHIRCodeableConceptInput{
			Text: "single",
		},
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRPatientInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: create patient",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create resource")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.CreateFHIRPatient(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_PatchFHIRPatient(t *testing.T) {
	ctx := context.Background()
	is_active := false
	input := domain.FHIRPatientInput{Active: &is_active}
	id := uuid.NewString()

	type args struct {
		ctx   context.Context
		id    string
		input domain.FHIRPatientInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: patch patient",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, input: input, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: error patching resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create resource")
					})

				return args{ctx: ctx, input: input, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.PatchFHIRPatient(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.PatchFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_PatchFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	cancelled := domain.EpisodeOfCareStatusEnumCancelled
	startTime := scalarutils.DateTime(time.Now().Format(scalarutils.DateTimeFormatLayout))
	endDateTime := scalarutils.DateTime(time.Now().Add(2 * time.Hour).Format(scalarutils.DateTimeFormatLayout))
	input := domain.FHIREpisodeOfCareInput{
		Status: &cancelled,
		Period: &domain.FHIRPeriodInput{
			Start: &startTime,
			End:   &endDateTime,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		input domain.FHIREpisodeOfCareInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully patch an episode of care",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, input: input, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable to patch episode of care",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.PatchFHIREpisodeOfCare(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.PatchFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIREpisodeOfCare(t *testing.T) {
	id := uuid.NewString()
	ctx := context.Background()
	payload := domain.FHIREpisodeOfCareInput{
		ID: &id,
	}

	type args struct {
		ctx            context.Context
		fhirResourceID string
		payload        domain.FHIREpisodeOfCareInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: update episode of care",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, payload: payload, fhirResourceID: id}
			},
			wantErr: false,
		},
		{
			name: "sad case: fhirResourceID nil",
			setup: func(mh *mockHandler) args {
				return args{ctx: ctx, payload: payload}
			},
			wantErr: true,
		},
		{
			name: "sad case: error updating resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, payload: payload, fhirResourceID: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.UpdateFHIREpisodeOfCare(args.ctx, args.fhirResourceID, args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.UpdateFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRMedicationStatement(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"name": "ARVs",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: search medication statement",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRMedication{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})

				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "sad case: search resource error",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return errors.New("error")
					})

				return args{ctx: ctx, params: params}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIRMedicationStatement(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRMedicationStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRPatient(t *testing.T) {
	ctx := context.Background()
	params := uuid.NewString()

	type args struct {
		ctx          context.Context
		searchParams string
		tenant       dto.TenantIdentifiers
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: search patient",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRPatient{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})

				return args{ctx: ctx, searchParams: params}
			},
			wantErr: false,
		},
		{
			name: "sad case: search patient error",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to search patient")
					})

				return args{ctx: ctx, searchParams: params}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIRPatient(args.ctx, args.searchParams, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRPatient(t *testing.T) {
	ID := uuid.NewString()
	total := 7
	language := "en"
	implicitrules := "http://example.org/rules"

	// Patient resource
	patientID := uuid.NewString()
	patient := map[string]any{
		"resourceType": "Patient",
		"id":           patientID,
		"active":       true,
		"name": []map[string]any{
			{
				"use":    "official",
				"family": "Doe",
				"given":  []string{"John"},
			},
		},
	}
	patientPayload, err := json.Marshal(patient)
	if err != nil {
		t.Fatalf("failed to marshal patient: %v", err)
	}

	// Encounter resource
	encounterID := uuid.NewString()
	encounter := map[string]any{
		"resourceType": "Encounter",
		"id":           encounterID,
		"status":       "finished",
		"class": map[string]any{
			"system": "http://terminology.hl7.org/CodeSystem/v3-ActCode",
			"code":   "AMB",
		},
		"subject": map[string]any{
			"reference": "Patient/" + patientID,
		},
	}
	encounterPayload, err := json.Marshal(encounter)
	if err != nil {
		t.Fatalf("failed to marshal encounter: %v", err)
	}

	// Organization resource
	orgID := uuid.NewString()
	organization := map[string]any{
		"resourceType": "Organization",
		"id":           orgID,
		"active":       true,
		"name":         "Test Organization",
	}
	organizationPayload, err := json.Marshal(organization)
	if err != nil {
		t.Fatalf("failed to marshal organization: %v", err)
	}

	// EpisodeOfCare resource
	episodeID := uuid.NewString()
	episode := map[string]any{
		"resourceType": "EpisodeOfCare",
		"id":           episodeID,
		"status":       "active",
		"patient": map[string]any{
			"reference": "Patient/" + patientID,
		},
		"managingOrganization": map[string]any{
			"reference": "Organization/" + orgID,
		},
		"period": map[string]any{
			"start": time.Now().Format(time.RFC3339),
		},
	}
	episodePayload, err := json.Marshal(episode)
	if err != nil {
		t.Fatalf("failed to marshal episode of care: %v", err)
	}

	documentReferenceID := uuid.NewString()
	documentReference := map[string]any{
		"resourceType": "DocumentReference",
		"id":           documentReferenceID,
		"status":       "active",
		"period": map[string]any{
			"start": time.Now().Format(time.RFC3339),
		},
	}

	documentReferencePayload, err := json.Marshal(documentReference)
	if err != nil {
		t.Fatalf("failed to marshal document reference: %v", err)
	}

	// Observation resource
	observationID := uuid.NewString()
	observation := map[string]any{
		"resourceType": "Observation",
		"id":           observationID,
		"status":       "final",
		"code": map[string]any{
			"coding": []map[string]any{
				{
					"system":  "http://loinc.org",
					"code":    "718-7",
					"display": "Hemoglobin",
				},
			},
		},
		"subject": map[string]any{
			"reference": "Patient/" + patientID,
		},
		"encounter": map[string]any{
			"reference": "Encounter/" + encounterID,
		},
		"effectiveDateTime": time.Now().Format(time.RFC3339),
	}
	observationPayload, err := json.Marshal(observation)
	if err != nil {
		t.Fatalf("failed to marshal observation: %v", err)
	}

	// MedicationRequest resource
	medRequestID := uuid.NewString()
	medRequest := map[string]any{
		"resourceType": "MedicationRequest",
		"id":           medRequestID,
		"status":       "active",
		"intent":       "order",
		"medicationCodeableConcept": map[string]any{
			"coding": []map[string]any{
				{
					"system":  "http://www.nlm.nih.gov/research/umls/rxnorm",
					"code":    "860975",
					"display": "Amoxicillin 500 MG Oral Capsule",
				},
			},
		},
		"subject": map[string]any{
			"reference": "Patient/" + patientID,
		},
		"encounter": map[string]any{
			"reference": "Encounter/" + encounterID,
		},
	}
	medRequestPayload, err := json.Marshal(medRequest)
	if err != nil {
		t.Fatalf("failed to marshal medication request: %v", err)
	}

	// Task resource
	taskID := uuid.NewString()
	task := map[string]any{
		"resourceType": "Task",
		"id":           taskID,
		"status":       "requested",
		"intent":       "order",
		"for": map[string]any{
			"reference": "Patient/" + patientID,
		},
		"encounter": map[string]any{
			"reference": "Encounter/" + encounterID,
		},
	}
	taskPayload, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("failed to marshal task: %v", err)
	}

	id := uuid.NewString()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: delete all patient data",
			setup: func(mh *mockHandler) args {
				ctx := context.Background()

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: scalarutils.DateTime(currentTime),
							End:   scalarutils.DateTime(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/next"},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/previous"},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &patientID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-patient"},
							},
							Resource: patientPayload,
						},
						{
							ID: &encounterID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-encounter"},
							},
							Resource: encounterPayload,
						},
						{
							ID: &observationID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-observation"},
							},
							Resource: observationPayload,
						},
						{
							ID: &medRequestID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-medication-request"},
							},
							Resource: medRequestPayload,
						},
						{
							ID: &taskID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-task"},
							},
							Resource: taskPayload,
						},
						{
							ID: &orgID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-organization"},
							},
							Resource: organizationPayload,
						},
						{
							ID: &episodeID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-episode-of-care"},
							},
							Resource: episodePayload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})

				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: all patient data error",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						return errors.New("error")
					})
				return args{ctx: context.Background(), id: id}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: all patient data invalid entry",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						return errors.New("errors")
					})
				return args{ctx: context.Background(), id: id}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: all patient data invalid entry type",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						return errors.New("error")
					})
				return args{ctx: context.Background(), id: id}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: all patient data entry invalid resource type",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						return errors.New("error")
					})
				return args{ctx: context.Background(), id: id}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: error deleting medication request",
			setup: func(mh *mockHandler) args {
				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: scalarutils.DateTime(currentTime),
							End:   scalarutils.DateTime(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/next"},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/previous"},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &medRequestID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-medication-request"},
							},
							Resource: medRequestPayload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						if resourceType == "MedicationRequest" {
							return fmt.Errorf("failed")
						}
						return nil
					})
				return args{ctx: context.Background(), id: id}
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting encounters",
			setup: func(mh *mockHandler) args {
				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: scalarutils.DateTime(currentTime),
							End:   scalarutils.DateTime(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/next"},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/previous"},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &encounterID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-encounter"},
							},
							Resource: encounterPayload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						if resourceType == "Encounter" {
							return fmt.Errorf("failed")
						}
						return nil
					})
				return args{ctx: context.Background(), id: id}
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting episode of care",
			setup: func(mh *mockHandler) args {
				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: scalarutils.DateTime(currentTime),
							End:   scalarutils.DateTime(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/next"},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/previous"},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &episodeID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-episode-of-care"},
							},
							Resource: episodePayload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						if resourceType == "EpisodeOfCare" {
							return fmt.Errorf("failed")
						}
						return nil
					})
				return args{ctx: context.Background(), id: id}
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting observation",
			setup: func(mh *mockHandler) args {
				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: scalarutils.DateTime(currentTime),
							End:   scalarutils.DateTime(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/next"},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/previous"},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &observationID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-observation"},
							},
							Resource: observationPayload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						if resourceType == "Observation" {
							return fmt.Errorf("failed")
						}
						return nil
					})
				return args{ctx: context.Background(), id: id}
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting patient",
			setup: func(mh *mockHandler) args {
				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: scalarutils.DateTime(currentTime),
							End:   scalarutils.DateTime(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/next"},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/previous"},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &patientID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-patient"},
							},
							Resource: patientPayload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						if resourceType == "Patient" {
							return fmt.Errorf("failed")
						}
						return nil
					})
				return args{ctx: context.Background(), id: id}
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting other types",
			setup: func(mh *mockHandler) args {
				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: scalarutils.DateTime(currentTime),
							End:   scalarutils.DateTime(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/next"},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/previous"},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &episodeID,
							Extension: []hapifhirmodels.Extension{
								{URL: "http://example.org/entry-ext-episode-of-care"},
							},
							Resource: documentReferencePayload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						if resourceType == "DocumentReference" {
							return fmt.Errorf("failed")
						}
						return nil
					})
				return args{ctx: context.Background(), id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.DeleteFHIRPatient(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.DeleteFHIRPatient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_CreateFHIRCondition(t *testing.T) {
	ID := uuid.NewString()
	ctx := context.Background()
	input := domain.FHIRConditionInput{
		ID:                 &ID,
		Identifier:         []*domain.FHIRIdentifierInput{},
		ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
		VerificationStatus: &domain.FHIRCodeableConceptInput{},
		Category:           []*domain.FHIRCodeableConceptInput{},
		Severity:           &domain.FHIRCodeableConceptInput{},
		Code:               &domain.FHIRCodeableConceptInput{},
		BodySite:           []*domain.FHIRCodeableConceptInput{},
		Subject:            &domain.FHIRReferenceInput{},
		Encounter:          &domain.FHIRReferenceInput{},
		OnsetDateTime: &scalarutils.Date{
			Year:  2000,
			Month: 3,
			Day:   30,
		},
		OnsetAge:    &domain.FHIRAgeInput{},
		OnsetPeriod: &domain.FHIRPeriodInput{},
		OnsetRange:  &domain.FHIRRangeInput{},
		OnsetString: new(string),
		AbatementDateTime: &scalarutils.Date{
			Year:  2000,
			Month: 3,
			Day:   30,
		},
		AbatementAge:    &domain.FHIRAgeInput{},
		AbatementPeriod: &domain.FHIRPeriodInput{},
		AbatementRange:  &domain.FHIRRangeInput{},
		AbatementString: new(string),
		RecordedDate: &scalarutils.Date{
			Year:  2000,
			Month: 3,
			Day:   30,
		},
		Recorder: &domain.FHIRReferenceInput{},
		Asserter: &domain.FHIRReferenceInput{},
		Stage:    []*domain.FHIRConditionStageInput{},
		Evidence: []*domain.FHIRConditionEvidenceInput{},
		Note:     []*domain.FHIRAnnotationInput{},
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRConditionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create FHIR condition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create FHIR condition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRCondition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIROrganization_Unittest(t *testing.T) {
	ctx := context.Background()
	ID := ksuid.New().String()
	active := true
	testname := gofakeit.FirstName()

	orgInput := domain.FHIROrganizationInput{
		ID:         &ID,
		Active:     &active,
		Identifier: []*domain.FHIRIdentifierInput{},
		Type:       []*domain.FHIRCodeableConceptInput{},
		Name:       &testname,
		Alias:      []string{"alias test"},
		Contact:    []domain.FHIROrganizationContactInput{},
		Address:    []*domain.FHIRAddressInput{},
	}

	type args struct {
		ctx   context.Context
		input domain.FHIROrganizationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create FHIR organization",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: orgInput}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create FHIR organization",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return errors.New("an error occurred")
					})
				return args{ctx: ctx, input: orgInput}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIROrganization(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIROrganization_Unittest(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	type args struct {
		ctx            context.Context
		organizationID string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *domain.FHIROrganizationRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case: find organization by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, organizationID: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: missing organization ID",
			setup: func(mh *mockHandler) args {
				return args{ctx: ctx, organizationID: ""}
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to find organization by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: context.Background(), organizationID: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.GetFHIROrganization(args.ctx, args.organizationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected response to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIREncounter(t *testing.T) {
	id := uuid.NewString()
	ctx := context.Background()

	input := domain.FHIREncounterInput{
		ID: &id,
	}
	type args struct {
		ctx   context.Context
		input domain.FHIREncounterInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *domain.FHIREncounterRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case: create FHIR encounter",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetVal("")
						return cmd
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create fhir service request")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to set data in cache",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						if encounter, ok := resource.(*domain.FHIREncounter); ok {
							encounter.ID = &id
						}
						return nil
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetErr(errors.New("failed to set record in cache"))
						return cmd
					})
				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIREncounter(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestStoreImpl_PatchFHIREncounter(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	startTime := scalarutils.DateTime(time.Now().Format(scalarutils.DateTimeFormatLayout))
	endDateTime := scalarutils.DateTime(time.Now().Add(2 * time.Hour).Format(scalarutils.DateTimeFormatLayout))
	input := domain.FHIREncounterInput{
		Status: domain.EncounterStatusEnumDischarged,
		ActualPeriod: &domain.FHIRPeriodInput{
			Start: &startTime,
			End:   &endDateTime,
		},
	}

	type args struct {
		ctx         context.Context
		encounterID string
		input       domain.FHIREncounterInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully patch an encounter",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, encounterID: id, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable to patch encounter",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, encounterID: id, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.PatchFHIREncounter(args.ctx, args.encounterID, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.PatchFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	id := ksuid.New().String()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: get FHIR episode of care",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIREpisodeOfCare(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestStoreImpl_GetFHIREncounter(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get encounter: cache miss",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						return cmd
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Happy case: get encounter: cache hit",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)

						id := gofakeit.UUID()
						patient := &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID: &id,
							},
						}

						bs, err := json.Marshal(patient)
						if err != nil {
							cmd.SetErr(err)
						}

						cmd.SetVal(string(bs))
						return cmd
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to retrieve encounter",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				return args{ctx: ctx, id: ""}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to set value in cache",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetErr(errors.New("failed to set record in cache"))
						return cmd
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIREncounter(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	input := domain.FHIRServiceRequestInput{
		ID: &id,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRServiceRequestInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						if servicerequest, ok := resource.(*domain.FHIRServiceRequest); ok {
							servicerequest.ID = &id
						}
						return nil
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetVal("")
						return cmd
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create FHIR service request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create fhir service request")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to set data in cache",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						if servicerequest, ok := resource.(*domain.FHIRServiceRequest); ok {
							servicerequest.ID = &id
						}

						return nil
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetErr(errors.New("failed to set cache"))
						return cmd
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRServiceRequest(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.NewString()
	input := domain.FHIRAllergyIntoleranceInput{
		ID: &UUID,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRAllergyIntoleranceInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create FHIR allergy intolerance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create fhir service request")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRAllergyIntolerance(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.New().String()
	input := domain.FHIRAllergyIntoleranceInput{
		ID: &UUID,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRAllergyIntoleranceInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update allergy intolerance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - missing input",
			setup: func(mh *mockHandler) args {
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.UpdateFHIRAllergyIntolerance(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRComposition(t *testing.T) {
	ctx := context.Background()
	input := domain.FHIRCompositionInput{}

	type args struct {
		ctx   context.Context
		input domain.FHIRCompositionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create FHIR composition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create FHIR composition")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRComposition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestStoreImpl_GetFHIRComposition(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "happy case: get composition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "sad case: error retrieving fhir resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.GetFHIRComposition(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRComposition(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.NewString()
	input := domain.FHIRCompositionInput{
		ID: &UUID,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRCompositionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing user ID",
			setup: func(mh *mockHandler) args {
				return args{ctx: ctx, input: domain.FHIRCompositionInput{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.UpdateFHIRComposition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_PatchFHIRComposition(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.NewString()
	id := ksuid.New().String()
	input := domain.FHIRCompositionInput{ID: &UUID}

	type args struct {
		ctx   context.Context
		id    string
		input domain.FHIRCompositionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case - succesfully patch composition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case - fail to patch composition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.PatchFHIRComposition(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.PatchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRComposition(t *testing.T) {
	ctx := context.Background()
	id := ksuid.New().String()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.DeleteFHIRComposition(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_CreateFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	input := domain.FHIRMedicationRequestInput{
		ID: &id,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationRequestInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create medication request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to create medication request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to create fhir service request")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.CreateFHIRMedicationRequest(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRServiceRequests(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	params := map[string]interface{}{
		"id": ID,
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully search fhir service request",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRServiceRequest{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a service request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to search resource")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRServiceRequest(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"id": "1234",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully search fhir allergy intolerance",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				allergy := domain.FHIRAllergyIntolerance{
					ID: &obsID,
				}

				payload, err := json.Marshal(allergy)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search an allergy intolerance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to search resource")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIRAllergyIntolerance(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRComposition(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"id": "1234",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully search fhir composition",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRComposition{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a composition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to search resource")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIRComposition(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRCondition(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"id": "1234",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully search fhir condition",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRCondition{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a condition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to search resource")
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIRCondition(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIREncounter(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"id": "1234",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *domain.FHIREncounterRelayConnection
		wantErr bool
	}{
		{
			name: "Happy Case - successfully search fhir encounter",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIREncounter{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search an encounter",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to search resource")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIREncounter(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"id": "1234",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully search fhir medication request",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRMedicationRequest{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a medication request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("error")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRMedicationRequest(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRCondition(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	input := domain.FHIRConditionInput{ID: &id}

	type args struct {
		ctx   context.Context
		input domain.FHIRConditionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully update fhir condition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to update fhir condition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - missing ID",
			setup: func(mh *mockHandler) args {
				return args{ctx: ctx, input: domain.FHIRConditionInput{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.UpdateFHIRCondition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.UpdateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRMedicationRequest(t *testing.T) {
	id := uuid.NewString()
	ctx := context.Background()
	input := domain.FHIRMedicationRequestInput{ID: &id}

	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationRequestInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully update fhir medication request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to update fhir medication request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return errors.New("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - missing ID",
			setup: func(mh *mockHandler) args {
				return args{ctx: ctx, input: domain.FHIRMedicationRequestInput{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.UpdateFHIRMedicationRequest(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.UpdateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - successfully delete a medication request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - fail to delete a medication request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().DeleteFHIRResource(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string) error {
						return fmt.Errorf("failed to update resource")
					})
				return args{ctx: ctx, id: id}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.DeleteFHIRMedicationRequest(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.DeleteFHIRMedicationRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_SearchPatientEncounters(t *testing.T) {
	ctx := context.Background()

	status := domain.EncounterStatusEnumPlanned
	type args struct {
		ctx              context.Context
		patientReference string
		status           *domain.EncounterStatusEnum
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: get encounters",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				enc := domain.FHIREncounter{
					ID: &obsID,
				}

				payload, err := json.Marshal(enc)
				if err != nil {
					t.Fatalf("failed to marshal encounter: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, patientReference: gofakeit.BS(), status: &status}
			},
			wantErr: false,
		},
		{
			name: "Happy case: nil status",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRObservation{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, patientReference: gofakeit.BS()}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, patientReference: gofakeit.BS(), status: &status}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchPatientEncounters(args.ctx, args.patientReference, args.status, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.Encounters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"test": "search",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case: search FHIR episode of care",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIREpisodeOfCare{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIREpisodeOfCare(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_CreateEpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	status := domain.EpisodeOfCareStatusEnumPlanned
	UUID := gofakeit.UUID()
	PatientRef := "Patient/1"
	OrgRef := "Organization/"
	episode := domain.FHIREpisodeOfCareInput{
		ID:     &UUID,
		Status: &status,
		Patient: &domain.FHIRReferenceInput{
			Reference: &PatientRef,
		},
		ManagingOrganization: &domain.FHIRReferenceInput{
			Reference: &OrgRef,
		},
	}

	type args struct {
		ctx     context.Context
		episode domain.FHIREpisodeOfCareInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create episode of care",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: context.Background(), episode: episode}
			},
			wantErr: false,
		},
		{
			name: "Happy case: create episode of care, episode does not exist",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, episode: episode}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, episode: episode}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.CreateEpisodeOfCare(args.ctx, args.episode)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIROrganization(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"test": "params",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search FHIR organisation",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIROrganization{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR organisation",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIROrganization(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchEpisodesByParam(t *testing.T) {
	ctx := context.Background()
	searchParams := map[string]interface{}{
		"period": map[string]interface{}{
			"start": time.February.String(),
			"end":   time.February.String(),
		},
	}

	type args struct {
		ctx          context.Context
		searchParams map[string]interface{}
		tenant       dto.TenantIdentifiers
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search episode by param",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"
				currentTime := time.Now().Format(time.DateOnly)

				obsID := uuid.NewString()
				eoc := domain.FHIREpisodeOfCare{
					ID: &obsID,
					Period: &domain.FHIRPeriod{
						ID:    new(string),
						Start: (scalarutils.DateTime)(currentTime),
						End:   (scalarutils.DateTime)(currentTime),
					},
				}

				payload, err := json.Marshal(eoc)
				if err != nil {
					t.Fatalf("failed to marshal episode of care: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, searchParams: searchParams}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, searchParams: searchParams}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchEpisodesByParam(args.ctx, args.searchParams, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchEpisodesByParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_OpenEpisodes(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx              context.Context
		patientReference string
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: open episodes",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				openEpisodes := domain.FHIREpisodeOfCare{
					ID: &obsID,
				}

				payload, err := json.Marshal(openEpisodes)
				if err != nil {
					t.Fatalf("failed to marshal episodes: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, patientReference: gofakeit.BS()}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, patientReference: gofakeit.BS()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.OpenEpisodes(args.ctx, args.patientReference, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.OpenEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_HasOpenEpisode(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	patient := domain.FHIRPatient{
		ID: &ID,
	}

	type args struct {
		ctx        context.Context
		patient    domain.FHIRPatient
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: has open episodes",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRObservation{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, patient: patient}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, patient: patient}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.HasOpenEpisode(args.ctx, args.patient, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.HasOpenEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.HasOpenEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_SearchEpisodeEncounter(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx              context.Context
		episodeReference string
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search episode encounter",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				encounter := domain.FHIREncounter{
					ID: &obsID,
				}

				payload, err := json.Marshal(encounter)
				if err != nil {
					t.Fatalf("failed to marshal encounter: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, episodeReference: gofakeit.BS()}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, episodeReference: gofakeit.BS()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchEpisodeEncounter(args.ctx, args.episodeReference, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchEpisodeEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_GetActiveEpisode(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()

	type args struct {
		ctx        context.Context
		episodeID  string
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: get active episodes",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				eoc := domain.FHIREpisodeOfCare{
					ID: &obsID,
				}

				payload, err := json.Marshal(eoc)
				if err != nil {
					t.Fatalf("failed to marshal episode of care: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, episodeID: ID}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, episodeID: ID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: empty FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, episodeID: ID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: FHIR resource has more than one entry",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return nil
					})
				return args{ctx: ctx, episodeID: ID}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.GetActiveEpisode(args.ctx, args.episodeID, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetActiveEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_StartEncounter(t *testing.T) {
	ctx := context.Background()
	status := domain.EpisodeOfCareStatusEnumActive
	finishedStatus := domain.EpisodeOfCareStatusEnumFinished
	UUID := uuid.NewString()
	dummyString := gofakeit.BS()
	uri := "foo://example.com:8042"

	episode := domain.FHIREpisodeOfCare{
		ID:     &UUID,
		Status: &status,
		Patient: &domain.FHIRReference{
			ID:        &UUID,
			Reference: &dummyString,
			Type:      (*scalarutils.URI)(&uri),
			Display:   dummyString,
		},
		ManagingOrganization: &domain.FHIRReference{
			ID:        &UUID,
			Reference: &dummyString,
			Type:      (*scalarutils.URI)(&uri),
			Display:   dummyString,
		},
	}

	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: start encounter",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return nil
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetVal("")
						return cmd
					})
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						episode := domain.FHIREpisodeOfCare{
							ID: &UUID,
						}
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return nil
					})
				return args{ctx: ctx, episodeID: UUID}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get FHIR episode of care",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, episodeID: UUID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: episode  not active",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						episode := domain.FHIREpisodeOfCare{
							ID:     &UUID,
							Status: &finishedStatus,
							Patient: &domain.FHIRReference{
								ID:        &UUID,
								Reference: &dummyString,
								Type:      (*scalarutils.URI)(&uri),
								Display:   dummyString,
							},
							ManagingOrganization: &domain.FHIRReference{
								ID:        &UUID,
								Reference: &dummyString,
								Type:      (*scalarutils.URI)(&uri),
								Display:   dummyString,
							},
						}
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return nil
					})
				return args{ctx: context.Background(), episodeID: gofakeit.UUID()}
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to create encounter",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return nil
					})
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						episode := domain.FHIREpisodeOfCare{
							ID: &UUID,
						}
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: context.Background(), episodeID: gofakeit.UUID()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.StartEncounter(args.ctx, args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_EndEncounter(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.NewString()
	encounter := domain.FHIREncounter{
		ID: &UUID,
		ActualPeriod: &domain.FHIRPeriod{
			ID:    &UUID,
			Start: scalarutils.DateTime(time.February.String()),
			End:   scalarutils.DateTime(time.March.String()),
		},
	}

	type args struct {
		ctx         context.Context
		encounterID string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: end encounter",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetVal("")
						return cmd
					})
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(encounter)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return nil
					})
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, encounterID: UUID}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: failed to get encounter",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(encounter)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, encounterID: UUID}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to update resource",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetVal("")
						return cmd
					})
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(encounter)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return nil
					})
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return errors.New("an error occurred")
					})
				return args{ctx: context.Background(), encounterID: gofakeit.UUID()}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.EndEncounter(args.ctx, args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.EndEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.EndEncounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_EndEpisode(t *testing.T) {
	ctx := context.Background()
	status := domain.EpisodeOfCareStatusEnumActive
	UUID := gofakeit.UUID()
	episode := domain.FHIREpisodeOfCare{
		ID:     &UUID,
		Status: &status,
		Period: &domain.FHIRPeriod{
			ID:    &UUID,
			Start: scalarutils.DateTime(time.February.String()),
			End:   scalarutils.DateTime(time.March.String()),
		},
	}

	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: end episode",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return nil
					})
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, episodeID: UUID}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: failed to get encounter",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return errors.New("an error occurred")
					})
				return args{ctx: ctx, episodeID: UUID}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to update resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						bs, err := json.Marshal(episode)
						if err != nil {
							return err
						}

						err = json.Unmarshal(bs, resource)
						if err != nil {
							return err
						}
						return nil
					})
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, episodeID: UUID}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.EndEpisode(args.ctx, args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.EndEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.EndEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_SearchPatientObservations(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	params := map[string]interface{}{
		"patient":         fmt.Sprintf("Patient/%s", ID),
		"encounter":       fmt.Sprintf("Encounter/%s", ID),
		"observationCode": "5088",
	}

	type args struct {
		ctx        context.Context
		bundleID   string
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination serverutils.PaginationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully search patient observation",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRObservation{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search fhir resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to search observation resource")
					})
				return args{ctx: ctx, params: params}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchPatientObservations(args.ctx, args.bundleID, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchPatientObservations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: get allergy intolerance by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: ID}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get allergy intolerance by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("error")
					})
				return args{ctx: ctx, id: ID}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRAllergyIntolerance(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchPatientAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	tenantIDs := dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	patientReference := fmt.Sprintf("Patient/%s", gofakeit.UUID())

	type args struct {
		ctx              context.Context
		patientReference string
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search allergy intolerance",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				ai := domain.FHIRAllergyIntolerance{
					ID: &obsID,
				}

				payload, err := json.Marshal(ai)
				if err != nil {
					t.Fatalf("failed to marshal allergy intolerance: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, patientReference: patientReference, tenant: tenantIDs, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search allergy intolerance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("some error")
					})
				return args{ctx: ctx, patientReference: patientReference, tenant: tenantIDs, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchPatientAllergyIntolerance(args.ctx, args.patientReference, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchPatientAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRMedia(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	input := domain.FHIRMedia{
		Subject: &domain.FHIRReferenceInput{
			ID: &id,
		},
	}
	type args struct {
		ctx   context.Context
		input domain.FHIRMedia
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create FHIR media",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create FHIR media",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error ocurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRMedia(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchPatientMedia(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	tenant := dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}
	first := 1
	pagination := dto.Pagination{
		First: &first,
	}

	type args struct {
		ctx              context.Context
		patientReference string
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search patient media",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				media := domain.FHIRMedia{
					ID: &obsID,
				}

				payload, err := json.Marshal(media)
				if err != nil {
					t.Fatalf("failed to marshal media: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, patientReference: gofakeit.BS(), tenant: tenant, pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search patient media",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("unable to search patient media")
					})
				return args{ctx: ctx, patientReference: gofakeit.BS(), tenant: tenant, pagination: pagination}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchPatientMedia(args.ctx, args.patientReference, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchPatentMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_ListFHIRQuestionnaire(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"status": "active",
		"name":   "title",
	}
	tenant := dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search questionnaire",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				questionnaire := domain.FHIRQuestionnaire{
					ID: &obsID,
				}

				payload, err := json.Marshal(questionnaire)
				if err != nil {
					t.Fatalf("failed to marshal questionnaire: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params, tenant: tenant, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search questionnaire",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, params: params, tenant: tenant, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.ListFHIRQuestionnaire(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.ListFHIRQuestionnaire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRQuestionnaire(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	input := &domain.FHIRQuestionnaire{
		ID:                &ID,
		Meta:              &domain.FHIRMetaInput{},
		ImplicitRules:     new(string),
		Language:          new(string),
		Text:              &domain.FHIRNarrative{},
		Extension:         []*domain.Extension{},
		ModifierExtension: []*domain.Extension{},
		Identifier:        []*domain.FHIRIdentifier{},
		Version:           new(string),
		Name:              new(string),
		Title:             new(string),
		DerivedFrom:       []*string{},
		Experimental:      new(bool),
		Publisher:         new(string),
		Description:       new(string),
		UseContext:        &domain.FHIRUsageContext{},
		Jurisdiction:      []*domain.FHIRCodeableConcept{},
		Purpose:           new(string),
		EffectivePeriod:   &domain.FHIRPeriod{},
		Code:              []*domain.FHIRCoding{},
		Item:              []*domain.FHIRQuestionnaireItem{},
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRQuestionnaire
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: successfully create a questionnaire resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create a questionnaire resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error ocurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRQuestionnaire(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRQuestionnaire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRConsent(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	status := domain.ConsentStatusActive
	decision := domain.ConsentDecisionDeny
	consent := domain.FHIRConsent{
		Status:   (*domain.ConsentStatusEnum)(&status),
		Decision: &decision,
		Provision: []*domain.FHIRConsentProvision{
			{
				ID: &ID,
			},
		},
		Subject: &domain.FHIRReference{ID: &ID},
		Grantor: []*domain.FHIRReference{
			{
				ID: &ID,
			},
		},
		Grantee: []*domain.FHIRReference{
			{
				ID: &ID,
			},
		},
		Manager: []*domain.FHIRReference{
			{
				ID: &ID,
			},
		},
		Controller: []*domain.FHIRReference{
			{
				ID: &ID,
			},
		},
		Meta: &domain.FHIRMetaInput{Tag: []domain.FHIRCodingInput{}},
	}

	type args struct {
		ctx     context.Context
		consent domain.FHIRConsent
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create a fhir consent",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, consent: consent}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create consent",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return errors.New("unable to create fhir consent")
					})
				return args{ctx: ctx, consent: consent}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRConsent(args.ctx, args.consent)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRConsent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRQuestionnaireResponse(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	input := &domain.FHIRQuestionnaireResponse{
		ID:            &ID,
		Meta:          &domain.FHIRMetaInput{},
		ImplicitRules: new(string),
		Language:      new(string),
		Text:          &domain.FHIRNarrative{},
		Item:          []domain.FHIRQuestionnaireResponseItem{},
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRQuestionnaireResponse
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: successfully create a questionnaire response resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create a questionnaire response resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return errors.New("unable to create fhir questionnair response")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRQuestionnaireResponse(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRQuestionnaireResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRRiskAssessment(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	input := &domain.FHIRRiskAssessmentInput{}

	type args struct {
		ctx   context.Context
		input *domain.FHIRRiskAssessmentInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create a risk assessment",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						if riskAssessment, ok := resource.(*domain.FHIRRiskAssessment); ok {
							riskAssessment.ID = &id
						}

						return nil
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						return cmd
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create a risk assessment",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return errors.New("unable to create fhir consent")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to set data in cache",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						if riskAssessment, ok := resource.(*domain.FHIRRiskAssessment); ok {
							riskAssessment.ID = &id
						}

						return nil
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetErr(errors.New("an error occurred"))
						return cmd
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.CreateFHIRRiskAssessment(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRRiskAssessment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRQuestionnaire(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully get questionnaire",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Fail to get questionnaire by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: ""}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRQuestionnaire(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRQuestionnaire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRDiagnosticReport(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	input := &domain.FHIRDiagnosticReportInput{ID: &ID}

	type args struct {
		ctx   context.Context
		input *domain.FHIRDiagnosticReportInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create diagnostic report",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create diagnostic report",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error ocurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRDiagnosticReport(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRDiagnosticReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRRiskAssessment(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		bundleID   string
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination serverutils.PaginationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully search a fhir risk assessment",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				ra := domain.FHIRRiskAssessment{
					ID: &obsID,
				}

				payload, err := json.Marshal(ra)
				if err != nil {
					t.Fatalf("failed to marshal risk assessment: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a fhir risk assessment",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to search risk assessment")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			got, err := fh.SearchFHIRRiskAssessment(args.ctx, args.bundleID, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRRiskAssessment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRPatientEverything(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"_count": 10,
	}

	type args struct {
		ctx    context.Context
		id     string
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: fetch patient everything",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := "123"
				patient := domain.FHIRPatient{
					ID: &obsID,
				}

				payload, err := json.Marshal(patient)
				if err != nil {
					t.Fatalf("failed to marshal patient: %v", err)
				}

				currentTime := time.Now().Format(time.RFC3339)

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Identifier: &hapifhirmodels.Identifier{
						Period: &hapifhirmodels.Period{
							ID:    new(string),
							Start: (scalarutils.DateTime)(currentTime),
							End:   (scalarutils.DateTime)(currentTime),
						},
					},
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, id: "123", params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to fetch patient everything",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetPatientEverything(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientFhirID string, searchParams map[string]interface{}, bundle interface{}) error {
						return errors.New("error")
					})
				return args{ctx: ctx, id: "1"}
			},
			wantErr: true,
		},
		{
			name: "Sad case: no patient id provided",
			setup: func(mh *mockHandler) args {
				return args{ctx: ctx, id: ""}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRPatientEverything(args.ctx, args.id, args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRPatientEverything() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func unMarshaller(data map[string]interface{}) (hapifhirmodels.Bundle, error) {
	var bundle hapifhirmodels.Bundle
	outpur, err := json.Marshal(data)
	if err != nil {
		return hapifhirmodels.Bundle{}, err
	}

	err = json.Unmarshal(outpur, &bundle)
	if err != nil {
		return hapifhirmodels.Bundle{}, err
	}
	return bundle, nil
}

func TestStoreImpl_CreateFHIRSubscription(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	subscription := &domain.FHIRSubscriptionInput{ID: &id}

	type args struct {
		ctx          context.Context
		subscription *domain.FHIRSubscriptionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create subscription",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, subscription: subscription}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create subscription",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, subscription: subscription}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRSubscription(args.ctx, args.subscription)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_PatchFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	input := domain.FHIRServiceRequestInput{
		ID:                 new(string),
		BodySite:           []*domain.FHIRCodeableConceptInput{},
		Note:               []*domain.FHIRAnnotationInput{},
		PatientInstruction: new(string),
		RelevantHistory:    []*domain.FHIRReferenceInput{},
		Meta:               domain.FHIRMetaInput{},
		Extension:          []*domain.FHIRExtension{},
	}

	type args struct {
		ctx   context.Context
		id    string
		input domain.FHIRServiceRequestInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: patch service request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to patch service request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.PatchFHIRServiceRequest(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.PatchFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRDocumentReference(t *testing.T) {
	ctx := context.Background()
	input := &domain.FHIRDocumentReferenceInput{}
	type args struct {
		ctx   context.Context
		input *domain.FHIRDocumentReferenceInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create a document reference",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create a document reference",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRDocumentReference(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRDocumentReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRDocumentReference(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	params := map[string]interface{}{"related": fmt.Sprintf("ServiceRequest/%s", ID)}

	type args struct {
		ctx          context.Context
		searchParams map[string]interface{}
		tenant       dto.TenantIdentifiers
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search for document reference",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				docref := domain.FHIRDocumentReference{
					ID: obsID,
				}

				payload, err := json.Marshal(docref)
				if err != nil {
					t.Fatalf("failed to marshal document reference: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, searchParams: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search for document reference",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, searchParams: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRDocumentReference(args.ctx, args.searchParams, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRDocumentReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRAppointment(t *testing.T) {
	ctx := context.Background()
	id := gofakeit.UUID()
	input := &domain.FHIRAppointmentInput{ID: &id}

	type args struct {
		ctx   context.Context
		input *domain.FHIRAppointmentInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create appointment",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create appointment",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRAppointment(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRAppointment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRTask(t *testing.T) {
	ctx := context.Background()
	status := "requested"
	input := &domain.FHIRTaskInput{Status: (*scalarutils.Code)(&status)}

	type args struct {
		ctx   context.Context
		input *domain.FHIRTaskInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create a task",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create a task",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRTask(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRTask(t *testing.T) {
	ctx := context.Background()
	status := "requested"
	id := uuid.NewString()
	endDateTime := scalarutils.DateTime(time.Now().Add(2 * time.Hour).Format(scalarutils.DateTimeFormatLayout))
	input := domain.FHIRTaskInput{
		ID:     &id,
		Status: (*scalarutils.Code)(&status),
		ExecutionPeriod: &domain.FHIRPeriodInput{
			End: &endDateTime,
		},
		Note: []*domain.FHIRAnnotationInput{
			{
				Text: (*scalarutils.Markdown)(&status),
			},
		},
		Reason: []*domain.FHIRCodeableReference{
			{
				Concept: &domain.FHIRCodeableConcept{
					Text: status,
				},
			},
		},
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRTaskInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: update task",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: no id",
			setup: func(mh *mockHandler) args {
				return args{ctx: context.Background(), input: domain.FHIRTaskInput{Status: (*scalarutils.Code)(&status)}}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to update  task",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.PatchFHIRTask(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.PatchFHIRTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRTask(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"status": "completed",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search for task",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				task := domain.FHIRTask{
					ID: &obsID,
				}

				payload, err := json.Marshal(task)
				if err != nil {
					t.Fatalf("failed to marshal task: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search for task",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRTask(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetOrSetCache(t *testing.T) {
	ctx := context.Background()
	mockValue := "mock value"
	mockJSON, _ := json.Marshal(mockValue)
	op := func() (any, error) { return mockValue, nil }
	op1 := func() (interface{}, error) { return nil, fmt.Errorf("failed to fetch data") }
	key := uuid.NewString()

	type args struct {
		ctx context.Context
		key string
		op  fhir.FHIROperation
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    string
		wantErr bool
		op      fhir.FHIROperation
	}{
		{
			name: "Happy Case: Cache hit",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetVal(mockValue)
						return cmd
					})

				return args{ctx: ctx, key: key, op: op}
			},
			want:    string(mockJSON),
			wantErr: false,
		},
		{
			name: "Happy Case: Cache miss",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						return cmd
					})

				return args{ctx: ctx, key: key, op: op}
			},
			want:    string(mockJSON),
			wantErr: false,
		},
		{
			name: "Sad Case: Cache miss, fail to set data in cache",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetErr(fmt.Errorf("failed to set data in cache"))
						return cmd
					})

				return args{ctx: ctx, key: key, op: op}
			},
			want:    string(mockJSON),
			wantErr: true,
		},
		{
			name: "Sad Case: Cache miss, fail to fetch data",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})

				return args{ctx: ctx, key: key, op: op1}
			},
			want:    string(mockJSON),
			wantErr: true,
		},
		{
			name: "Sad Case: Cache hit, return error",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(fmt.Errorf("an error occured"))
						return cmd
					})

				return args{ctx: ctx, key: key, op: op}
			},
			want:    string(mockJSON),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetOrSetCache(args.ctx, args.key, args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetOrSetCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIREncounterAllData(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *domain.PagedFHIRResource
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get all encounter data",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				enc := domain.FHIREncounter{
					ID: &obsID,
				}

				payload, err := json.Marshal(enc)
				if err != nil {
					t.Fatalf("failed to marshal encounter: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get all encounter data",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to get resource")
					})
				return args{ctx: ctx, params: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIREncounterAllData(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIREncounterAllData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRObservation(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get observations",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get observations",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("failed to get observation")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRObservation(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRQuestionnaireResponse(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get questionnaire response",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get questionnaire response",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("failed to get questionnaire response")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRQuestionnaireResponse(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRQuestionnaireResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get service request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get service request",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("failed to get service request")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)
			_, err := fh.GetFHIRServiceRequest(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRTask(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: get fhir task by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get fhir task by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("failed to get task by id")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRTask(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRResource(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"_sort":       "_lastUpdated",
		"category":    "57133-8",
		"_revinclude": "DocumentReference:relatesto",
	}

	type args struct {
		ctx          context.Context
		bundleID     string
		resourceType string
		params       map[string]interface{}
		tenant       dto.TenantIdentifiers
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully search FHIR resource",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				observation := domain.FHIRObservation{
					ID: &obsID,
				}

				payload, err := json.Marshal(observation)
				if err != nil {
					t.Fatalf("failed to marshal observation: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, resourceType: "ServiceRequest", params: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get FHIR resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("failed to get resource")
					})
				return args{ctx: ctx, resourceType: "ServiceRequest", params: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)
			got, err := fh.SearchFHIRResource(args.ctx, args.bundleID, args.resourceType, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_PatchFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.NewString()
	id := ksuid.New().String()
	status := domain.CompletedMedicationStatus
	input := domain.FHIRMedicationRequestInput{
		ID:     &UUID,
		Status: &status,
	}

	type args struct {
		ctx   context.Context
		id    string
		input domain.FHIRMedicationRequestInput
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case - successfully patch medication resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case - fail to patch medication resource",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().FHIRPathPatch(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("failed to patch medication")
					})
				return args{ctx: ctx, id: id, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.PatchFHIRMedicationRequest(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.PatchFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: retrieve medication request by its ID (Cache hit)",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)

						id := gofakeit.UUID()
						patient := &domain.FHIRMedicationRequestRelayPayload{
							Resource: &domain.FHIRMedicationRequest{
								ID: &id,
							},
						}

						bs, err := json.Marshal(patient)
						if err != nil {
							cmd.SetErr(err)
						}

						cmd.SetVal(string(bs))
						return cmd
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Happy case: retrieve medication request by its ID (Cache miss)",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})

				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						return cmd
					})
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to retrieve medication request",
			setup: func(mh *mockHandler) args {
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("failed to get medication")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Unable to set value in cache",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				mh.cache.EXPECT().Get(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
						cmd := redis.NewStringCmd(ctx)
						cmd.SetErr(redis.Nil)
						return cmd
					})

				mh.cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
						cmd := redis.NewStatusCmd(ctx)
						cmd.SetErr(errors.New("failed to set record in cache"))
						return cmd
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRMedicationRequest(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRCondition(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully get condition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Fail to get condition by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: ""}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRCondition(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRPractitionerRole(t *testing.T) {
	ctx := context.Background()
	id := gofakeit.UUID()
	type args struct {
		ctx   context.Context
		input *domain.FHIRPractitionerRole
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create a practitioner role",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{
					ctx: ctx,
					input: &domain.FHIRPractitionerRole{
						ID: &id,
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create a practitioner role",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{
					ctx: ctx,
					input: &domain.FHIRPractitionerRole{
						ID: &id,
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRPractitionerRole(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRPractitionerRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRPractitionerRoles(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"lastUpdated": time.Now().UTC(),
	}

	type args struct {
		ctx          context.Context
		searchParams map[string]interface{}
		tenant       dto.TenantIdentifiers
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search practitioner roles",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				role := domain.FHIRPractitionerRole{
					ID: &obsID,
				}

				payload, err := json.Marshal(role)
				if err != nil {
					t.Fatalf("failed to marshal practitioner role: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, searchParams: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search practitioner roles",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, searchParams: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRPractitionerRoles(args.ctx, args.searchParams, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRPractitionerRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRPractitionerRole(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: get practitioner role by id",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get practitioner role by id",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRPractitionerRole(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRPractitionerRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRLocation(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	input := &domain.FHIRLocation{
		ID: &ID,
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRLocation
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create a location",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create a location",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRLocation(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRLocation(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: successfully get location",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to get location by ID",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRLocation(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestStoreImpl_GetFHIRPractitioner(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully gets a practitioner",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Fails to get a practitioner",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: ""}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRPractitioner(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRPractitioner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRPractioner(t *testing.T) {
	ctx := context.Background()
	phone := gofakeit.Phone()
	use := domain.ContactPointUseEnumHome
	phoneSystem := domain.ContactPointSystemEnumPhone
	id := uuid.NewString()
	name := gofakeit.Username()
	practitioner := &domain.FHIRPractitioner{
		ID: &id,
		Name: []*domain.FHIRHumanName{
			{
				Given: []*string{&name},
			},
		},
		Telecom: []*domain.FHIRContactPoint{
			{
				System: &phoneSystem,
				Use:    &use,
				Period: common.DefaultPeriod(),
				Value:  &phone,
			},
		},
		Gender: domain.PatientGenderEnumMale,
		BirthDate: &scalarutils.Date{
			Year:  gofakeit.Year(),
			Month: int(gofakeit.Date().Month()),
			Day:   gofakeit.Date().Day(),
		},
	}
	type args struct {
		ctx          context.Context
		resourceType string
		input        *domain.FHIRPractitioner
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully creates a practitioner",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, resourceType: "Practitioner", input: practitioner}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create a practitioner",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, resourceType: "Practitioner", input: practitioner}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRPractioner(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRPractioner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRPractitioner(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"_sort": "-date",
	}

	type args struct {
		ctx          context.Context
		params       map[string]interface{}
		tenant       dto.TenantIdentifiers
		resourceType string
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully gets practitioners",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				practitioner := domain.FHIRPractitioner{
					ID: &obsID,
				}

				payload, err := json.Marshal(practitioner)
				if err != nil {
					t.Fatalf("failed to marshal practitioner: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, resourceType: "Practitioner", params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to get practitioners",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, resourceType: "Practitioner", params: params}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRPractitioner(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRPractitioner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRSubstance(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	input := &domain.FHIRSubstance{
		ID: &ID,
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRSubstance
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create substance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable create substance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRSubstance(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRSubstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRSubstance(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: get substance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get substance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRSubstance(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRSubstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRSubstance(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"lastUpdated": time.Now().UTC(),
	}

	type args struct {
		ctx          context.Context
		searchParams map[string]interface{}
		tenant       dto.TenantIdentifiers
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: search substance",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				substance := domain.FHIRSubstance{
					ID: &obsID,
				}

				payload, err := json.Marshal(substance)
				if err != nil {
					t.Fatalf("failed to marshal substance: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, searchParams: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search substance",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, searchParams: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRSubstance(args.ctx, args.searchParams, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRSubstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRProcedure(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	input := &domain.FHIRProcedure{
		ID: &ID,
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRProcedure
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully creates a procedure",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to create a procedure",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRProcedure(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRProcedure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRProcedure(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully fetches a procedure",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to fetch a procedure",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRProcedure(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRProcedure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRProcedure(t *testing.T) {
	ctx := context.Background()

	params := map[string]interface{}{
		"_sort": "-date",
	}

	type args struct {
		ctx          context.Context
		params       map[string]interface{}
		tenant       dto.TenantIdentifiers
		resourceType string
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully gets procedures",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				procedure := domain.FHIRProcedure{
					ID: &obsID,
				}

				payload, err := json.Marshal(procedure)
				if err != nil {
					t.Fatalf("failed to marshal procedure: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, resourceType: "Procedure", params: params}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to get procedures",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, resourceType: "Procedure", params: params}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRProcedure(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRProcedure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestStoreImpl_CreateFHIRMedicationDispense(t *testing.T) {
	ID := uuid.NewString()
	ctx := context.Background()

	input := &domain.FHIRMedicationDispense{
		ID: &ID,
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRMedicationDispense
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully creates a medication dispense",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unaable to create a medication dispense",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRMedicationDispense(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedicationDispense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRMedicationDispense(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"_sort": "-date",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully search medication dispenses",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				md := domain.FHIRMedicationDispense{
					ID: &obsID,
				}

				payload, err := json.Marshal(md)
				if err != nil {
					t.Fatalf("failed to marshal medication dispense: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to search for medication dispenses",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, params: params, tenant: dto.TenantIdentifiers{}, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRMedicationDispense(args.ctx, args.params, args.tenant, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRMedicationDispense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRMedicationDispense(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully get a medication dispense",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to get a medication dispense",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRMedicationDispense(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRMedicationDispense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRPlanDefinition(t *testing.T) {
	ctx := context.Background()
	uid := uuid.NewString()
	input := &domain.FHIRPlanDefinition{
		ID: &uid,
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRPlanDefinition
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create plan definition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: create plan definition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRPlanDefinition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRPlanDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRActivityDefinition(t *testing.T) {
	ctx := context.Background()
	uid := uuid.NewString()
	input := &domain.FHIRActivityDefinition{
		ID: &uid,
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRActivityDefinition
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create activity definition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: create activity definition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRActivityDefinition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRActivityDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRCarePlan(t *testing.T) {
	ctx := context.Background()
	uid := uuid.NewString()
	input := &domain.FHIRCarePlan{
		ID: &uid,
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRCarePlan
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create care plan",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create care plan",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().CreateFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType string, payload map[string]interface{}, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.CreateFHIRCarePlan(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRCarePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_FetchMedicationByID(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: fetch mediccation by id",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to fetch mediccation by id",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.FetchMedicationByID(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.FetchMedicationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRPlanDefinition(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"_sort": "_lastUpdated",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully search plan definition",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				pd := domain.FHIRPlanDefinition{
					ID: &obsID,
				}

				payload, err := json.Marshal(pd)
				if err != nil {
					t.Fatalf("failed to marshal plan definition: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to search for plan definition",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, params: params, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRPlanDefinition(args.ctx, args.params, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRPlanDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRCarePlan(t *testing.T) {
	ctx := context.Background()
	params := map[string]interface{}{
		"_sort": "_lastUpdated",
	}

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully fetch care plan",
			setup: func(mh *mockHandler) args {
				ID := uuid.NewString()
				total := 1
				language := "en"
				implicitrules := "http://example.org/rules"

				obsID := uuid.NewString()
				cp := domain.FHIRCarePlan{
					ID: &obsID,
				}

				payload, err := json.Marshal(cp)
				if err != nil {
					t.Fatalf("failed to marshal care plan: %v", err)
				}

				expectedBundle := &hapifhirmodels.Bundle{
					ID:            &ID,
					Total:         &total,
					Language:      &language,
					ImplicitRules: &implicitrules,
					Link: []hapifhirmodels.BundleLink{
						{
							ID:  &ID,
							URL: "next",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/next",
								},
							},
						},
						{
							ID:  &ID,
							URL: "previous",
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/previous",
								},
							},
						},
					},
					Entry: []hapifhirmodels.BundleEntry{
						{
							ID: &ID,
							Extension: []hapifhirmodels.Extension{
								{
									URL: "http://example.org/entry-ext",
								},
							},
							Resource: payload,
						},
					},
				}

				mh.hapiFHIR.EXPECT().
					SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						bundleID string,
						resourceType string,
						p map[string]any,
						tenant dto.TenantIdentifiers,
						bundle any,
					) error {
						b, ok := bundle.(*hapifhirmodels.Bundle)
						if !ok {
							return fmt.Errorf("unexpected bundle type %T", bundle)
						}

						*b = *expectedBundle
						return nil
					})
				return args{ctx: ctx, params: params, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to search for care plan",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, bundle interface{}) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, params: params, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.SearchFHIRCarePlan(args.ctx, args.params, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRCarePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRRiskAssessment(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get risk assessment",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return nil
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get risk assessment",
			setup: func(mh *mockHandler) args {
				mh.hapiFHIR.EXPECT().GetFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, fhirResourceID string, resource interface{}) error {
						return fmt.Errorf("failed to get risk assessment")
					})
				return args{ctx: ctx, id: id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}
			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			_, err := fh.GetFHIRRiskAssessment(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRRiskAssessment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_PutFHIRResource(t *testing.T) {
	type args struct {
		ctx                context.Context
		resourceType       string
		resourceID         string
		payload            map[string]any
		resource           any
		useCREnabledServer bool
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully put a FHIR resource",
			setup: func(mh *mockHandler) args {
				ctx := context.Background()
				id := uuid.NewString()
				mh.hapiFHIR.EXPECT().PutFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]any, resource any, useCREnabledServer bool) error {
						return nil
					})
				return args{ctx: ctx, resourceType: "Observation", resourceID: id, payload: map[string]any{"key": "value"}, resource: &domain.FHIRObservation{}, useCREnabledServer: false}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to put a FHIR resource",
			setup: func(mh *mockHandler) args {
				ctx := context.Background()
				id := uuid.NewString()
				mh.hapiFHIR.EXPECT().PutFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string, payload map[string]any, resource any, useCREnabledServer bool) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, resourceType: "Observation", resourceID: id, payload: map[string]any{"key": "value"}, resource: &domain.FHIRObservation{}, useCREnabledServer: false}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := fakeCache.NewMockCacheService(t)
			mockhapiFHIR := fakeHapiFHIR.NewMockHapiFHIRImplementation(t)
			mockHandler := mockHandler{
				mockCache,
				mockhapiFHIR,
			}

			fh := fhir.NewFHIRStoreImpl(mockCache, mockhapiFHIR)

			args := tt.setup(&mockHandler)

			if err := fh.PutFHIRResource(args.ctx, args.resourceType, args.resourceID, args.payload, args.resource, args.useCREnabledServer); (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.PutFHIRResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
