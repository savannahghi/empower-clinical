package base_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_StartEncounter(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully start an encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						episodeID := uuid.NewString()
						patientID := uuid.NewString()
						orgRefID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumActive
						ref := mock.Anything
						uri := scalarutils.URI(gofakeit.URL())
						return &domain.FHIREpisodeOfCareRelayPayload{
							Resource: &domain.FHIREpisodeOfCare{
								ID:     &episodeID,
								Status: &status,
								Patient: &domain.FHIRReference{
									ID:        &patientID,
									Reference: &ref,
									Type:      &uri,
									Display:   mock.Anything,
								},
								ManagingOrganization: &domain.FHIRReference{
									ID:        &orgRefID,
									Reference: &ref,
									Type:      &uri,
									Display:   mock.Anything,
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						orgName := gofakeit.Name()
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:   &orgID,
								Name: &orgName,
							},
						}, nil
					})
				mh.FHIR.EXPECT().CreateFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID:     &encounterID,
								Status: domain.EncounterStatusEnumInProgress,
							},
						}, nil
					})

				return args{ctx: ctx, episodeID: uuid.NewString()}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Empty episode of care ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, episodeID: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						episodeID := uuid.NewString()
						patientID := uuid.NewString()
						orgRefID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumActive
						ref := mock.Anything
						uri := scalarutils.URI(gofakeit.URL())
						return &domain.FHIREpisodeOfCareRelayPayload{
							Resource: &domain.FHIREpisodeOfCare{
								ID:     &episodeID,
								Status: &status,
								Patient: &domain.FHIRReference{
									ID:        &patientID,
									Reference: &ref,
									Type:      &uri,
									Display:   mock.Anything,
								},
								ManagingOrganization: &domain.FHIRReference{
									ID:        &orgRefID,
									Reference: &ref,
									Type:      &uri,
									Display:   mock.Anything,
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						orgName := gofakeit.Name()
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:   &orgID,
								Name: &orgName,
							},
						}, nil
					})
				mh.FHIR.EXPECT().CreateFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, episodeID: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						episodeID := uuid.NewString()
						patientID := uuid.NewString()
						orgRefID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumActive
						ref := mock.Anything
						uri := scalarutils.URI(gofakeit.URL())
						return &domain.FHIREpisodeOfCareRelayPayload{
							Resource: &domain.FHIREpisodeOfCare{
								ID:     &episodeID,
								Status: &status,
								Patient: &domain.FHIRReference{
									ID:        &patientID,
									Reference: &ref,
									Type:      &uri,
									Display:   mock.Anything,
								},
								ManagingOrganization: &domain.FHIRReference{
									ID:        &orgRefID,
									Reference: &ref,
									Type:      &uri,
									Display:   mock.Anything,
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, episodeID: uuid.NewString()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.StartEncounter(args.ctx, args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == "" {
					t.Errorf("expected an episode of care ID but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_PatchEncounter(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		encounterID string
		input       dto.EncounterInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully patch encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID: &encounterID,
							},
						}, nil
					})
				mh.FHIR.EXPECT().PatchFHIREncounter(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string, input domain.FHIREncounterInput) (*domain.FHIREncounter, error) {
						id := uuid.NewString()
						return &domain.FHIREncounter{
							ID: &id,
						}, nil
					})

				return args{ctx: ctx, encounterID: uuid.NewString(), input: dto.EncounterInput{Status: dto.EncounterStatusEnumCancelled}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid encounterID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, input: dto.EncounterInput{Status: dto.EncounterStatusEnumInProgress}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: uuid.NewString(), input: dto.EncounterInput{Status: dto.EncounterStatusEnumCancelled}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to patch encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().PatchFHIREncounter(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string, input domain.FHIREncounterInput) (*domain.FHIREncounter, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, encounterID: uuid.NewString(), input: dto.EncounterInput{Status: dto.EncounterStatusEnumInProgress}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.PatchEncounter(args.ctx, args.encounterID, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_EndEncounter(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		encounterID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully end encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return true, nil
					})

				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Missing encounter ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to end encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return false, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.EndEncounter(args.ctx, args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.EndEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesClinicalImpl.EndEncounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListPatientEncounters(t *testing.T) {
	ctx := context.Background()
	first := 3
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list patient encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := uuid.NewString()
						return &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &patientID,
								Name: []*domain.FHIRHumanName{
									{
										Text: gofakeit.Name(),
									},
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchPatientEncounters(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
						encounterID := uuid.NewString()
						return &domain.PagedFHIREncounter{
							Encounters: []domain.FHIREncounter{
								{
									ID:     &encounterID,
									Status: domain.EncounterStatusEnumInProgress,
								},
							},
						}, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), pagination: &dto.Pagination{First: &first, Skip: false}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing patient ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, pagination: &dto.Pagination{First: &first, Skip: false}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get fhir patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), pagination: &dto.Pagination{First: &first, Skip: false}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := uuid.NewString()
						return &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &patientID,
								Name: []*domain.FHIRHumanName{
									{
										Text: gofakeit.Name(),
									},
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), pagination: &dto.Pagination{First: &first, Skip: false}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get patient encounters",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := uuid.NewString()
						return &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &patientID,
								Name: []*domain.FHIRHumanName{
									{
										Text: gofakeit.Name(),
									},
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchPatientEncounters(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), pagination: &dto.Pagination{First: &first, Skip: false}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid pagination",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, patientID: uuid.NewString(), pagination: &dto.Pagination{First: &first, Last: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.ListPatientEncounters(args.ctx, args.patientID, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListPatientEncounters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetEncounterAssociatedResources(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		encounterID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list encounter's all data",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIREncounterAllData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: []map[string]interface{}{
								{
									"resourceType": "Observation",
									"id":           "9012",
									"status":       "final",
									"valueString":  "positive",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Biopsy of cervix",
									},
									"effectiveInstant": "2024-02-13T10:22:54+03:00",
									"note": []map[string]interface{}{
										{
											"text": "This is a note",
										},
									},
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "laboratory",
													"display":      "Laboratory",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "Laboratory",
										},
									},
								},
								{
									"resourceType": "Encounter",
									"id":           "5678",
									"status":       "in-progress",
								},
								{
									"basis": []map[string]interface{}{
										{
											"reference": "QuestionnaireResponse/af03246c-eac9-4d84-bec5-509d8c064918",
										},
									},
									"encounter": map[string]interface{}{
										"id":        "a54e4681-5f5c-455a-9feb-cedf84f0ba39",
										"reference": "Encounter/a54e4681-5f5c-455a-9feb-cedf84f0ba39",
									},
									"id":           "b72a852b-dcf5-48da-ac25-d343b0e43172",
									"resourceType": "RiskAssessment",
									"status":       "final",
									"subject": map[string]interface{}{
										"display":   "Omar, Lula ",
										"id":        "f6693702-3629-471c-959e-3cd3eaaec6d9",
										"reference": "Patient/f6693702-3629-471c-959e-3cd3eaaec6d9",
									},
								},
								{
									"businessStatus": map[string]interface{}{
										"text": "VIA test",
									},
									"description": "A VIA test",
									"encounter": map[string]interface{}{
										"id":        "62d8d969-e9f9-4e5c-95e3-3ab0fd3034aa",
										"reference": "Encounter/62d8d969-e9f9-4e5c-95e3-3ab0fd3034aa",
									},
									"id":       "2ce15512-e97e-40cf-a2ab-9a5df9e92404",
									"intent":   "order",
									"priority": "routine",
									"reasonCode": map[string]interface{}{
										"text": "Breast Cancer Screening",
									},
									"resourceType": "Task",
									"status":       "requested",
								},
								{
									"businessStatus": map[string]interface{}{
										"text": "VIA test",
									},
									"description": "A VIA test",
									"encounter": map[string]interface{}{
										"reference": "Encounter/62d8d969-e9f9-4e5c-95e3-3ab0fd3034aa",
									},
									"id":       "2ce15512-e97e-40cf-a2ab-9a5df9e92404",
									"intent":   "order",
									"priority": "routine",
									"reasonCode": map[string]interface{}{
										"text": "Breast Cancer Screening",
									},
									"resourceType": "Task",
									"status":       "requested",
								},
								{
									"resourceType": "Consent",
									"id":           "5678",
									"provision": []map[string]interface{}{
										{
											"type": "permit",
										},
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})

				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing encounter ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to search all fhir encounter data",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIREncounterAllData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetEncounterAssociatedResources(args.ctx, args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetEncounterAssociatedResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetScreeningReport(t *testing.T) {
	name := "Nairobi Hospital"
	pagedFHIRResource := &domain.PagedFHIRResource{
		Resources: []map[string]interface{}{
			{
				"authoredOn": "2024-06-10T08:09:03+08:00",
				"category": []map[string]interface{}{
					{
						"coding": []map[string]interface{}{
							{
								"code":    "167731",
								"display": "Referral",
							},
						},
						"text": "Referral",
					},
				},
				"code": map[string]interface{}{
					"concept": map[string]interface{}{
						"coding": []map[string]interface{}{
							{
								"code":    "159623",
								"display": "Diagnostics",
							},
							{
								"code":    "TEST",
								"display": "Mammogram",
							},
						},
						"text": "Facility Referral Reason",
					},
				},
				"encounter": map[string]interface{}{
					"id":        "65d22133-2f17-4599-99ab-9f3ee0ef020a",
					"reference": "Encounter/65d22133-2f17-4599-99ab-9f3ee0ef020a",
				},
				"extension": []map[string]interface{}{
					{
						"extension": []map[string]interface{}{
							{
								"url":         "facilityName",
								"valueString": "One pad more",
							},
							{
								"url":         "facilityContact",
								"valueString": "+254727645367",
							},
							{
								"url":         "facilityCounty",
								"valueString": "NAIROBI",
							},
						},
						"url": "http://savannahghi.org/fhir/StructureDefinition/referred-facility",
					},
				},
				"id":       "551ea02e-22cd-4916-bb31-6ffccce9b203",
				"intent":   "order",
				"language": "EN",
				"meta": map[string]interface{}{
					"lastUpdated": "2024-06-10T08:09:03.793059+00:00",
					"tag": []map[string]interface{}{
						{
							"code":         "e303711e-d24a-4fda-87d4-a4c177c7e90d",
							"display":      "Roche",
							"system":       "http://mycarehub/tenant-identification/organisation",
							"userSelected": false,
							"version":      "1.0",
						},
						{
							"code":         "ad031755-9149-44a8-a702-367bafa2ed40",
							"display":      "Main Branch",
							"system":       "http://mycarehub/tenant-identification/facility",
							"userSelected": false,
							"version":      "1.0",
						},
					},
					"versionId": "MTcxODAwNjk0Mzc5MzA1OTAwMA",
				},
				"note": []map[string]interface{}{
					{
						"text": "Testing",
						"time": "2024-06-10T08:09:03+08:00",
					},
				},
				"priority":     "urgent",
				"resourceType": "ServiceRequest",
				"status":       "active",
				"subject": map[string]interface{}{
					"display":   "Talisha, Idah ",
					"id":        "4d808b60-5149-443b-9d03-a6016e5af1b5",
					"reference": "Patient/4d808b60-5149-443b-9d03-a6016e5af1b5",
				},
			},
			{
				"content": []map[string]interface{}{
					{
						"attachment": map[string]interface{}{
							"contentType": "application/pdf",
							"title":       "Doe, Jane 's Referral report",
							"url":         "https://example.invalid/fixtures/report.pdf",
						},
					},
				},
				"context": []map[string]interface{}{
					{
						"reference": "ServiceRequest/551ea02e-22cd-4916-bb31-6ffccce9b203",
					},
				},
				"date":      "2024-06-11T07:57:44Z",
				"docStatus": "final",
				"id":        "9100df95-fbfd-457d-808e-6844cd16ffc2",
				"language":  "EN",
				"meta": map[string]interface{}{
					"lastUpdated": "2024-06-11T07:57:45.480128+00:00",
					"tag": []map[string]interface{}{
						{
							"code":         "e303711e-d24a-4fda-87d4-a4c177c7e90d",
							"display":      "Roche",
							"system":       "http://mycarehub/tenant-identification/organisation",
							"userSelected": false,
							"version":      "1.0",
						},
						{
							"code":         "ad031755-9149-44a8-a702-367bafa2ed40",
							"display":      "Main Branch",
							"system":       "http://mycarehub/tenant-identification/facility",
							"userSelected": false,
							"version":      "1.0",
						},
					},
					"versionId": "MTcxODA5MjY2NTQ4MDEyODAwMA",
				},
				"resourceType": "DocumentReference",
				"status":       "current",
				"subject": map[string]interface{}{
					"id":        "e8131505-b8da-4ad9-b991-f5fc43e5ab48",
					"reference": "Patient/e8131505-b8da-4ad9-b991-f5fc43e5ab48",
				},
				"type": map[string]interface{}{
					"coding": []map[string]interface{}{
						{
							"system": "",
						},
					},
				},
			},
		},
		HasNextPage:     false,
		NextCursor:      "",
		HasPreviousPage: false,
		PreviousCursor:  "",
		TotalCount:      0,
	}

	type args struct {
		ctx         context.Context
		encounterID string
		status      domain.ServiceRequestStatusEnum
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list encounter's all data (vital signs)",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID: &encounterID,
								Subject: &domain.FHIRReference{
									Display: mock.Anything,
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIREncounterAllData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: []map[string]interface{}{
								{
									"resourceType": "Observation",
									"id":           "9012",
									"status":       "final",
									"valueString":  "positive",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Biopsy of cervix",
									},
									"effectiveInstant": "2024-02-13T10:22:54+03:00",
									"note": []map[string]interface{}{
										{
											"text": "This is a note",
										},
									},
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "vital-signs",
													"display":      "Vital Signs",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "Vital Signs",
										},
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return pagedFHIRResource, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								Name: &name,
							},
						}, nil
					},
				)
				return args{
					ctx:         usecaseMock.AddTenantIdentifierContext(context.Background()),
					encounterID: uuid.NewString(),
					status:      domain.ServiceRequestStatusActive,
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully list encounter's all data (exam)",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID: &encounterID,
								Subject: &domain.FHIRReference{
									Display: mock.Anything,
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIREncounterAllData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: []map[string]interface{}{
								{
									"resourceType": "Observation",
									"id":           "9012",
									"status":       "final",
									"valueString":  "positive",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Biopsy of cervix",
									},
									"effectiveInstant": "2024-02-13T10:22:54+03:00",
									"note": []map[string]interface{}{
										{
											"text": "This is a note",
										},
									},
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "exam",
													"display":      "exam",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "Exam",
										},
									},
								},
								{
									"resourceType": "Encounter",
									"id":           "5678",
									"status":       "in-progress",
								},
								{
									"basis": []map[string]interface{}{
										{
											"reference": "QuestionnaireResponse/af03246c-eac9-4d84-bec5-509d8c064918",
										},
									},
									"encounter": map[string]interface{}{
										"id":        "a54e4681-5f5c-455a-9feb-cedf84f0ba39",
										"reference": "Encounter/a54e4681-5f5c-455a-9feb-cedf84f0ba39",
									},
									"id":           "b72a852b-dcf5-48da-ac25-d343b0e43172",
									"resourceType": "RiskAssessment",
									"status":       "final",
									"subject": map[string]interface{}{
										"display":   "Omar, Lula ",
										"id":        "f6693702-3629-471c-959e-3cd3eaaec6d9",
										"reference": "Patient/f6693702-3629-471c-959e-3cd3eaaec6d9",
									},
								},
								{
									"businessStatus": map[string]interface{}{
										"text": "VIA test",
									},
									"description": "A VIA test",
									"encounter": map[string]interface{}{
										"id":        "62d8d969-e9f9-4e5c-95e3-3ab0fd3034aa",
										"reference": "Encounter/62d8d969-e9f9-4e5c-95e3-3ab0fd3034aa",
									},
									"id":       "2ce15512-e97e-40cf-a2ab-9a5df9e92404",
									"intent":   "order",
									"priority": "routine",
									"reasonCode": map[string]interface{}{
										"text": "Breast Cancer Screening",
									},
									"resourceType": "Task",
									"status":       "requested",
								},
								{
									"resourceType": "Consent",
									"id":           "5678",
									"provision": []map[string]interface{}{
										{
											"type": "permit",
										},
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return pagedFHIRResource, nil
					})

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								Name: &name,
							},
						}, nil
					},
				)
				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), encounterID: uuid.NewString(), status: domain.ServiceRequestStatusActive}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully list encounter's all data (lab, imaging, procedure)",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID: &encounterID,
								Subject: &domain.FHIRReference{
									Display: mock.Anything,
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIREncounterAllData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: []map[string]interface{}{
								{
									"resourceType": "Observation",
									"id":           "9012",
									"status":       "final",
									"valueString":  "positive",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Biopsy of cervix",
									},
									"effectiveInstant": "2024-02-13T10:22:54+03:00",
									"note": []map[string]interface{}{
										{
											"text": "This is a note",
										},
									},
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "laboratory",
													"display":      "Laboratory",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "Laboratory",
										},
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return pagedFHIRResource, nil
					})

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								Name: &name,
							},
						}, nil
					},
				)
				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), encounterID: uuid.NewString(), status: domain.ServiceRequestStatusActive}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - unable to list encounter's all data",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID: &encounterID,
								Subject: &domain.FHIRReference{
									Display: mock.Anything,
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIREncounterAllData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), encounterID: uuid.NewString(), status: domain.ServiceRequestStatusActive}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), encounterID: uuid.NewString(), status: domain.ServiceRequestStatusActive}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get patient referrals",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID: &encounterID,
								Subject: &domain.FHIRReference{
									Display: mock.Anything,
								},
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIREncounterAllData(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: []map[string]interface{}{
								{
									"resourceType": "Observation",
									"id":           "9012",
									"status":       "final",
									"valueString":  "positive",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Biopsy of cervix",
									},
									"effectiveInstant": "2024-02-13T10:22:54+03:00",
									"note": []map[string]interface{}{
										{
											"text": "This is a note",
										},
									},
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "laboratory",
													"display":      "Laboratory",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "Laboratory",
										},
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), encounterID: uuid.NewString(), status: domain.ServiceRequestStatusActive}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetScreeningReport(args.ctx, args.encounterID, args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetScreeningReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_EndScreening(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		encounterID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully end screening",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID:     &encounterID,
								Status: domain.EncounterStatusEnumInProgress,
							},
						}, nil
					})
				mh.FHIR.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return true, nil
					})
				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Missing encounter ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to end encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID:     &encounterID,
								Status: domain.EncounterStatusEnumInProgress,
							},
						}, nil
					})
				mh.FHIR.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return false, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Try to end a finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounterID := uuid.NewString()
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID:     &encounterID,
								Status: domain.EncounterStatusEnumCompleted,
							},
						}, nil
					})

				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Get encounter with nil ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID:     new(string),
								Status: domain.EncounterStatusEnumInProgress,
							},
						}, nil
					})
				return args{ctx: ctx, encounterID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.EndScreening(args.ctx, args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.EndScreening() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesClinicalImpl.EndScreening() = %v, want %v", got, tt.want)
			}
		})
	}
}
