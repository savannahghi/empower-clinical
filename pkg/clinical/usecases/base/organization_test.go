package base_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_RegisterTenant(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.OrganizationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully register tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						orgName := gofakeit.Name()
						active := true
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:     &orgID,
								Name:   &orgName,
								Active: &active,
								Identifier: []*domain.FHIRIdentifier{
									{
										Value: mock.Anything,
									},
								},
							},
						}, nil
					})

				return args{
					ctx: ctx,
					input: dto.OrganizationInput{
						Name:        "Test facility",
						PhoneNumber: "0700000000",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("Other"),
								Value: "001",
							},
						},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: ctx,
					input: dto.OrganizationInput{
						Name:        "Test facility",
						PhoneNumber: "0700000000",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing name",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: ctx,
					input: dto.OrganizationInput{
						PhoneNumber: "0700000000",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("Other"),
								Value: "001",
							},
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create organisation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: ctx,
					input: dto.OrganizationInput{
						Name:        "Test facility",
						PhoneNumber: "0700000000",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("Other"),
								Value: "001",
							},
						},
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

			got, err := clinicalUsecase.RegisterTenant(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RegisterTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got: %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RegisterFacility(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.OrganizationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case - successfully register facility",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						orgName := gofakeit.Name()
						active := true
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:     &orgID,
								Name:   &orgName,
								Active: &active,
								Identifier: []*domain.FHIRIdentifier{
									{
										Value: mock.Anything,
									},
								},
							},
						}, nil
					})

				return args{
					ctx: ctx,
					input: dto.OrganizationInput{
						Name:        "Test",
						PhoneNumber: "Number",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("slade-code"),
								Value: "1234",
							},
						},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case - fail to register facility",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: ctx,
					input: dto.OrganizationInput{
						Name:        "Test",
						PhoneNumber: "Number",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("slade-code"),
								Value: "1234",
							},
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing slade code / mfl code identifier",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: ctx,
					input: dto.OrganizationInput{
						Name:        "Test",
						PhoneNumber: "Number",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("Other"),
								Value: "1234",
							},
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing name",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: ctx,
					input: dto.OrganizationInput{
						PhoneNumber: "Number",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("slade-code"),
								Value: "1234",
							},
						},
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

			got, err := clinicalUsecase.RegisterFacility(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RegisterFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got: %v", got)
					return
				}
			}
		})
	}
}

func TestBaseImpl_CreatePubsubTenant(t *testing.T) {

	type args struct {
		ctx  context.Context
		data dto.OrganizationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: Create pubsub patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						orgName := gofakeit.Name()
						active := true
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:     &orgID,
								Name:   &orgName,
								Active: &active,
								Identifier: []*domain.FHIRIdentifier{
									{
										Value: mock.Anything,
									},
								},
							},
						}, nil
					})
				mh.PubSub.EXPECT().NotifyProgramFHIRIDUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.UpdateProgramFHIRID) error {
						return nil
					})

				return args{
					ctx: context.Background(),
					data: dto.OrganizationInput{
						Name:        "Test facility",
						PhoneNumber: "0700000000",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("MCHProgram"),
								Value: "001",
							},
						},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Fail to create pubsub patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						orgName := gofakeit.Name()
						active := true
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:     &orgID,
								Name:   &orgName,
								Active: &active,
								Identifier: []*domain.FHIRIdentifier{
									{
										Value: mock.Anything,
									},
								},
							},
						}, nil
					})
				mh.PubSub.EXPECT().NotifyProgramFHIRIDUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.UpdateProgramFHIRID) error {
						return fmt.Errorf("an error occurred")
					})

				return args{
					ctx: context.Background(),
					data: dto.OrganizationInput{
						Name:        "Test facility",
						PhoneNumber: "0700000000",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("MCHProgram"),
								Value: "001",
							},
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: Fail to register tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: context.Background(),
					data: dto.OrganizationInput{
						Name:        "Test facility",
						PhoneNumber: "0700000000",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("MCHProgram"),
								Value: "001",
							},
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid organization identifier",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: context.Background(),
					data: dto.OrganizationInput{
						Name:        "Test facility",
						PhoneNumber: "0700000000",
						Identifiers: []dto.OrganizationIdentifier{
							{
								Type:  dto.OrganizationIdentifierType("Other"),
								Value: "001",
							},
						},
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

			if err := clinicalUsecase.CreatePubsubTenant(args.ctx, args.data); (err != nil) != tt.wantErr {
				t.Errorf("BaseImpl.CreatePubsubTenant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseImpl_ProvisionTenant(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.ProvisionTenantInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: successfully provision new tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{
					TenantID:              gofakeit.UUID(),
					Name:                  gofakeit.Company(),
					Status:                dto.TenantStatusActive,
					LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
					LegacyIdentifierValue: gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, input.TenantID).
					Return(nil, fmt.Errorf("FHIR error (HTTP 404): resource not found"))

				mh.FHIR.EXPECT().PutFHIROrganization(mock.Anything, input.TenantID, mock.Anything).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID: &input.TenantID,
						},
					}, nil)

				return args{ctx: context.Background(), input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: provision inactive tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{
					TenantID:              gofakeit.UUID(),
					Name:                  gofakeit.Company(),
					Status:                dto.TenantStatusInactive,
					LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
					LegacyIdentifierValue: gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, input.TenantID).
					Return(nil, fmt.Errorf("FHIR error (HTTP 404): resource not found"))

				mh.FHIR.EXPECT().PutFHIROrganization(mock.Anything, input.TenantID, mock.Anything).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID: &input.TenantID,
						},
					}, nil)

				return args{ctx: context.Background(), input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: provision tenant with parent organization",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{
					TenantID:              gofakeit.UUID(),
					Name:                  gofakeit.Company(),
					Status:                dto.TenantStatusActive,
					LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
					LegacyIdentifierValue: gofakeit.UUID(),
					ParentID:              gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, input.TenantID).
					Return(nil, fmt.Errorf("FHIR error (HTTP 404): resource not found"))

				mh.FHIR.EXPECT().PutFHIROrganization(mock.Anything, input.TenantID, mock.Anything).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID: &input.TenantID,
						},
					}, nil)

				return args{ctx: context.Background(), input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully update existing tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{
					TenantID:              gofakeit.UUID(),
					Name:                  gofakeit.Company(),
					Status:                dto.TenantStatusActive,
					LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
					LegacyIdentifierValue: gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, input.TenantID).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID: &input.TenantID,
						},
					}, nil)

				mh.FHIR.EXPECT().PutFHIROrganization(mock.Anything, input.TenantID, mock.Anything).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID: &input.TenantID,
						},
					}, nil)

				return args{ctx: context.Background(), input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid input fails validation",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{}

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid status",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{
					TenantID:              gofakeit.UUID(),
					Name:                  gofakeit.Company(),
					Status:                "INVALID_STATUS",
					LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
					LegacyIdentifierValue: gofakeit.UUID(),
				}

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: non-404 error checking existing tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{
					TenantID:              gofakeit.UUID(),
					Name:                  gofakeit.Company(),
					Status:                dto.TenantStatusActive,
					LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
					LegacyIdentifierValue: gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, input.TenantID).
					Return(nil, fmt.Errorf("FHIR error (HTTP 500): internal server error"))

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to update existing tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{
					TenantID:              gofakeit.UUID(),
					Name:                  gofakeit.Company(),
					Status:                dto.TenantStatusActive,
					LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
					LegacyIdentifierValue: gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, input.TenantID).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID: &input.TenantID,
						},
					}, nil)

				mh.FHIR.EXPECT().PutFHIROrganization(mock.Anything, input.TenantID, mock.Anything).
					Return(nil, fmt.Errorf("update failed"))

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to create new tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ProvisionTenantInput{
					TenantID:              gofakeit.UUID(),
					Name:                  gofakeit.Company(),
					Status:                dto.TenantStatusActive,
					LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
					LegacyIdentifierValue: gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, input.TenantID).
					Return(nil, fmt.Errorf("FHIR error (HTTP 404): resource not found"))

				mh.FHIR.EXPECT().PutFHIROrganization(mock.Anything, input.TenantID, mock.Anything).
					Return(nil, fmt.Errorf("create failed"))

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ProvisionTenant(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseImpl.ProvisionTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBaseImpl_GetTenantProvisioningStatus(t *testing.T) {
	type args struct {
		ctx      context.Context
		tenantID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: successfully get active tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				tenantID := gofakeit.UUID()
				name := gofakeit.Company()
				active := true

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, tenantID).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID:     &tenantID,
							Name:   &name,
							Active: &active,
						},
					}, nil)

				return args{ctx: context.Background(), tenantID: tenantID}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully get inactive tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				tenantID := gofakeit.UUID()
				name := gofakeit.Company()
				active := false

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, tenantID).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID:     &tenantID,
							Name:   &name,
							Active: &active,
						},
					}, nil)

				return args{ctx: context.Background(), tenantID: tenantID}
			},
			wantErr: false,
		},
		{
			name: "Happy case: tenant with nil resource",
			setup: func(mh *usecaseMock.Mocks) args {
				tenantID := gofakeit.UUID()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, tenantID).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: nil,
					}, nil)

				return args{ctx: context.Background(), tenantID: tenantID}
			},
			wantErr: false,
		},
		{
			name: "Happy case: tenant with nil ID and Name",
			setup: func(mh *usecaseMock.Mocks) args {
				tenantID := gofakeit.UUID()
				active := true

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, tenantID).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							Active: &active,
						},
					}, nil)

				return args{ctx: context.Background(), tenantID: tenantID}
			},
			wantErr: false,
		},
		{
			name: "Happy case: tenant with nil Active",
			setup: func(mh *usecaseMock.Mocks) args {
				tenantID := gofakeit.UUID()
				name := gofakeit.Company()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, tenantID).
					Return(&domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID:   &tenantID,
							Name: &name,
						},
					}, nil)

				return args{ctx: context.Background(), tenantID: tenantID}
			},
			wantErr: false,
		},
		{
			name: "Sad case: tenant not found",
			setup: func(mh *usecaseMock.Mocks) args {
				tenantID := gofakeit.UUID()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, tenantID).
					Return(nil, fmt.Errorf("not found"))

				return args{ctx: context.Background(), tenantID: tenantID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: empty tenant ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), tenantID: ""}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mocks := usecaseMock.SetupMocks(t)
			args := tt.setup(&mocks)

			_, err := clinicalUsecase.GetTenantProvisioningStatus(args.ctx, args.tenantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseImpl.GetTenantProvisioningStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
