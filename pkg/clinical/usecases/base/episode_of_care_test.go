package base_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_CreateEpisodeOfCare(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.EpisodeOfCareInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: create an episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
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
				mh.FHIR.EXPECT().SearchFHIREpisodeOfCare(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
						return &domain.FHIREpisodeOfCareRelayConnection{
							Edges:    []*domain.FHIREpisodeOfCareRelayEdge{},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})
				mh.FHIR.EXPECT().CreateEpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error) {
						episodeID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumActive
						return &domain.EpisodeOfCarePayload{
							EpisodeOfCare: &domain.FHIREpisodeOfCare{
								ID:     &episodeID,
								Status: &status,
							},
						}, nil
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: false,
		},
		{
			name: "sad case: create an episode of care already exists",
			setup: func(mh *usecaseMock.Mocks) args {
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
				mh.FHIR.EXPECT().SearchFHIREpisodeOfCare(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
						PatientRef := "Patient/1"
						OrgRef := "Organization/1"
						return &domain.FHIREpisodeOfCareRelayConnection{
							Edges: []*domain.FHIREpisodeOfCareRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIREpisodeOfCare{
										ID:            new(string),
										Text:          &domain.FHIRNarrative{},
										Identifier:    []*domain.FHIRIdentifier{},
										StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
										Type:          []*domain.FHIRCodeableConcept{},
										Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
										Patient:       &domain.FHIRReference{Reference: &PatientRef},
										ManagingOrganization: &domain.FHIRReference{
											Reference: &OrgRef,
										},
										Period:          &domain.FHIRPeriod{},
										ReferralRequest: []*domain.FHIRReference{},
										CareManager:     &domain.FHIRReference{},
										Team:            []*domain.FHIRReference{},
										Account:         []*domain.FHIRReference{},
										Meta:            &domain.FHIRMeta{},
										Extension:       []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "sad case: missing facility identifier in context",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: context.Background(),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "sad case: error fetching facility",
			setup: func(mh *usecaseMock.Mocks) args {
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

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "sad case: error fetching organization",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "sad case: error fetching patient",
			setup: func(mh *usecaseMock.Mocks) args {
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
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
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

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to search episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
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
				mh.FHIR.EXPECT().SearchFHIREpisodeOfCare(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "sad Case - Fail to get tenant meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
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
					}).Once()
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
					}).Times(2)
				mh.FHIR.EXPECT().SearchFHIREpisodeOfCare(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
						return &domain.FHIREpisodeOfCareRelayConnection{
							Edges:    []*domain.FHIREpisodeOfCareRelayEdge{},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					}).Once()

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to create episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
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
				mh.FHIR.EXPECT().SearchFHIREpisodeOfCare(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
						return &domain.FHIREpisodeOfCareRelayConnection{
							Edges:    []*domain.FHIREpisodeOfCareRelayEdge{},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})
				mh.FHIR.EXPECT().CreateEpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumActive,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.CreateEpisodeOfCare(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_PatchEpisodeOfCare(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		id    string
		input dto.EpisodeOfCareInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Patch an episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						episodeID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumActive
						return &domain.FHIREpisodeOfCareRelayPayload{
							Resource: &domain.FHIREpisodeOfCare{
								ID:     &episodeID,
								Status: &status,
							},
						}, nil
					})
				mh.FHIR.EXPECT().PatchFHIREpisodeOfCare(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, episode domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error) {
						episodeID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumCancelled
						return &domain.FHIREpisodeOfCare{
							ID:     &episodeID,
							Status: &status,
						}, nil
					})

				return args{
					ctx: ctx,
					id:  uuid.NewString(),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumCancelled,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid episode of care ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: ctx,
					id:  "123",
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumCancelled,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: ctx,
					id:  uuid.NewString(),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumCancelled,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to patch episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						episodeID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumActive
						return &domain.FHIREpisodeOfCareRelayPayload{
							Resource: &domain.FHIREpisodeOfCare{
								ID:     &episodeID,
								Status: &status,
							},
						}, nil
					})
				mh.FHIR.EXPECT().PatchFHIREpisodeOfCare(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, episode domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: ctx,
					id:  uuid.NewString(),
					input: dto.EpisodeOfCareInput{
						Status:    dto.EpisodeOfCareStatusEnumCancelled,
						PatientID: gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.PatchEpisodeOfCare(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_EndEpisodeOfCare(t *testing.T) {

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: end episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchEpisodeEncounter(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
						encounterID := uuid.NewString()
						return &domain.PagedFHIREncounter{
							Encounters: []domain.FHIREncounter{
								{
									ID: &encounterID,
								},
							},
						}, nil
					})
				mh.FHIR.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return true, nil
					})
				mh.FHIR.EXPECT().EndEpisode(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeID string) (bool, error) {
						return true, nil
					})
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						episodeID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumActive
						return &domain.FHIREpisodeOfCareRelayPayload{
							Resource: &domain.FHIREpisodeOfCare{
								ID:     &episodeID,
								Status: &status,
							},
						}, nil
					})

				return args{ctx: context.Background(), id: uuid.NewString()}
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid episode of care id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), id: "123"}
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchEpisodeEncounter(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
						encounterID := uuid.NewString()
						return &domain.PagedFHIREncounter{
							Encounters: []domain.FHIREncounter{
								{
									ID: &encounterID,
								},
							},
						}, nil
					})
				mh.FHIR.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return true, nil
					})
				mh.FHIR.EXPECT().EndEpisode(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeID string) (bool, error) {
						return true, nil
					})
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), id: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), id: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to search encounters",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchEpisodeEncounter(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), id: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to end encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchEpisodeEncounter(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
						encounterID := uuid.NewString()
						return &domain.PagedFHIREncounter{
							Encounters: []domain.FHIREncounter{
								{
									ID: &encounterID,
								},
							},
						}, nil
					})
				mh.FHIR.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return false, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), id: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to end episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().SearchEpisodeEncounter(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
						encounterID := uuid.NewString()
						return &domain.PagedFHIREncounter{
							Encounters: []domain.FHIREncounter{
								{
									ID: &encounterID,
								},
							},
						}, nil
					})
				mh.FHIR.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return true, nil
					})
				mh.FHIR.EXPECT().EndEpisode(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeID string) (bool, error) {
						return false, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), id: uuid.NewString()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.EndEpisodeOfCare(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetEpisodeOfCare(t *testing.T) {

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: get episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						episodeID := uuid.NewString()
						status := domain.EpisodeOfCareStatusEnumActive
						return &domain.FHIREpisodeOfCareRelayPayload{
							Resource: &domain.FHIREpisodeOfCare{
								ID:     &episodeID,
								Status: &status,
							},
						}, nil
					})

				return args{ctx: context.Background(), id: uuid.NewString()}
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid episode of care id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), id: "invalid"}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get episode of care",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), id: uuid.NewString()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetEpisodeOfCare(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}
