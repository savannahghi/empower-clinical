package foundation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_RecordConsent(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	ID := gofakeit.UUID()
	input := dto.ConsentInput{
		EncounterID: ID,
		Decision:    domain.ConsentDecisionPermit,
		DenyReason:  "",
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID: &ID,
			Subject: &domain.FHIRReference{
				ID:      &ID,
				Display: ID,
			},
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	type args struct {
		ctx   context.Context
		input dto.ConsentInput
	}

	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create a fhir consent",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().CreateFHIRConsent(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRConsent) (*domain.FHIRConsent, error) {
						status := domain.ConsentStatusActive
						return &domain.FHIRConsent{
							Status: &status,
						}, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to create consent",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Decision = domain.ConsentDecisionDeny

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().CreateFHIRConsent(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRConsent) (*domain.FHIRConsent, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid encounter id",
			setup: func(mh *usecaseMock.Mocks) args {
				input.EncounterID = ""

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid context",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed on finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Decision = domain.ConsentDecisionDeny

				encounter.Resource.Status = domain.EncounterStatusEnumCompleted
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
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

			_, err := clinicalUsecase.RecordConsent(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordConsent() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}

}
