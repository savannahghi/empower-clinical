package specialized_test

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

func TestSpecializedImpl_CreatePlanDefinition(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	ID := uuid.NewString()
	tenantIDs := dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}
	input := &dto.PlanDefinitionInput{
		Title:       gofakeit.BeerName(),
		Description: gofakeit.BeerName(),
		Action: []dto.PlanAction{
			{
				Title:       gofakeit.BeerName(),
				Description: gofakeit.BeerName(),
				TimingTiming: &dto.Timing{
					Repeat: &dto.Repeat{
						Frequency:  1,
						Period:     2,
						PeriodUnit: "wk",
						Count:      4,
						Offset:     0,
					},
				},
				Medications: []dto.PlanMedication{
					{
						MedicationID: gofakeit.BeerName(),
						Dosage: dto.DosageAdministrationInput{
							AdministrationInstructions: gofakeit.UUID(),
						},
					},
				},
				Action: []dto.PlanAction{
					{
						Title:        gofakeit.BeerName(),
						Description:  gofakeit.BeerName(),
						TimingTiming: &dto.Timing{},
						Medications:  []dto.PlanMedication{},
						Action:       []dto.PlanAction{},
					},
				},
			},
		},
	}
	medication := &domain.FHIRMedication{
		ID: &ID,
		Code: &domain.FHIRCodeableConcept{
			Text: gofakeit.Name(),
		},
	}
	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	type args struct {
		ctx   context.Context
		input *dto.PlanDefinitionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create plan definition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &tenantIDs, nil
					})
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return medication, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRActivityDefinition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRActivityDefinition) (*domain.FHIRActivityDefinition, error) {
						return &domain.FHIRActivityDefinition{ID: &ID}, nil
					})
				mh.FHIR.EXPECT().CreateFHIRPlanDefinition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRPlanDefinition) (*domain.FHIRPlanDefinition, error) {
						return &domain.FHIRPlanDefinition{ID: &ID}, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("error")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: incomplete plan definition",
			setup: func(mh *usecaseMock.Mocks) args {
				input := &dto.PlanDefinitionInput{
					Title:       gofakeit.BeerName(),
					Description: gofakeit.BeerName(),
				}
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create plan definition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &tenantIDs, nil
					})
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return medication, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRActivityDefinition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRActivityDefinition) (*domain.FHIRActivityDefinition, error) {
						return &domain.FHIRActivityDefinition{ID: &ID}, nil
					})
				mh.FHIR.EXPECT().CreateFHIRPlanDefinition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRPlanDefinition) (*domain.FHIRPlanDefinition, error) {
						return nil, fmt.Errorf("an error occured")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to fetch medication",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &tenantIDs, nil
					})
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return nil, fmt.Errorf("error")
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create activity definition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &tenantIDs, nil
					})
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return medication, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRActivityDefinition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRActivityDefinition) (*domain.FHIRActivityDefinition, error) {
						return nil, fmt.Errorf("error")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreatePlanDefinition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpecializedImpl.CreatePlanDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSpecializedImpl_RetrievePlanDefinition(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: Search for plan definition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().SearchFHIRPlanDefinition(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRPlanDefinition, error) {
						return &domain.PagedFHIRPlanDefinition{}, nil
					})
				return args{ctx: ctx, name: "Breast"}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to search for plan definition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().SearchFHIRPlanDefinition(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRPlanDefinition, error) {
						return nil, fmt.Errorf("unable top search for plan definition")
					})
				return args{ctx: ctx, name: "Breast"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RetrievePlanDefinition(args.ctx, args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpecializedImpl.RetrievePlanDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
