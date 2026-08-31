package clinical_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestClinicalImpl_CreateMedication(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := []*dto.MedicationInput{
		{
			Name: gofakeit.Name(),
			DoseForm: dto.ValueSetData{
				Code:    gofakeit.UUID(),
				Display: gofakeit.UUID(),
			},
		},
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	ID := uuid.NewString()
	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	status := domain.MedicationStatusEnumActive
	code := scalarutils.Code(gofakeit.BeerName())
	lotNumber := "12"
	expirationDate := time.Now().GoString()
	medication := &domain.FHIRMedication{
		ID:     &ID,
		Status: &status,
		DoseForm: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					Code:    &code,
					Display: gofakeit.BeerName(),
				},
			},
		},
		Extension: []domain.Extension{
			{
				ValueString: gofakeit.BeerName(),
			},
		},
		Batch: &domain.FHIRMedicationBatch{
			LotNumber:      &lotNumber,
			ExpirationDate: &expirationDate,
		},
	}

	type args struct {
		ctx   context.Context
		input []*dto.MedicationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: creates medication successful",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedication(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
						return &domain.FHIRMedicationRelayPayload{
							Resource: medication,
						}, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create medication",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedication(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RecordMedication(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.CreateMedication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
