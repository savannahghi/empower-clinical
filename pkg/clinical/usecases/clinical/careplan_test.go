package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	r5hapifhirmodel "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestClinicalImpl_SearchCarePlan(t *testing.T) {
	ten := 10
	type args struct {
		ctx        context.Context
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: search care plan",
			setup: func(mh *usecaseMock.Mocks) args {
				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
				pagination := dto.Pagination{
					First: &ten,
				}

				mh.FHIR.EXPECT().SearchFHIRCarePlan(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error) {
						ID := uuid.NewString()
						careplan := &domain.PagedFHIRCarePlan{
							CarePlan: []domain.FHIRCarePlan{
								{ID: &ID},
							},
						}
						return careplan, nil
					})

				return args{ctx: ctx, pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search care plan",
			setup: func(mh *usecaseMock.Mocks) args {
				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
				pagination := dto.Pagination{
					First: &ten,
				}

				mh.FHIR.EXPECT().SearchFHIRCarePlan(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, pagination: pagination}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.SearchCarePlan(args.ctx, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.SearchCarePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalImpl_CreateCarePlan(t *testing.T) {
	uid := gofakeit.UUID()

	type args struct {
		ctx   context.Context
		input *domain.FHIRCarePlan
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create care plan",
			setup: func(mh *usecaseMock.Mocks) args {
				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
				input := &domain.FHIRCarePlan{
					ID: &uid,
				}

				mh.FHIR.EXPECT().CreateFHIRCarePlan(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRCarePlan) (*domain.FHIRCarePlan, error) {
						return nil, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create care plan",
			setup: func(mh *usecaseMock.Mocks) args {
				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
				input := &domain.FHIRCarePlan{
					ID: &uid,
				}

				mh.FHIR.EXPECT().CreateFHIRCarePlan(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRCarePlan) (*domain.FHIRCarePlan, error) {
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

			_, err := clinicalUsecase.CreateCarePlan(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.CreateCarePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalImpl_PatientCarePlan(t *testing.T) {
	data := dto.CarePlanInput{
		EncounterID:      gofakeit.UUID(),
		PlanDefinitionID: gofakeit.UUID(),
		Notes:            gofakeit.UUID(),
	}
	payload := dto.CarePlanPayload{
		Data:       data,
		Tags:       []domain.FHIRCodingInput{},
		FacilityID: uuid.NewString(),
	}

	id := uuid.NewString()
	patientRef := "Patient/" + uuid.NewString()
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:            &id,
			Text:          &domain.FHIRNarrative{},
			Identifier:    []*domain.FHIRIdentifier{},
			Status:        domain.EncounterStatusEnum(domain.EncounterStatusEnumDischarged),
			StatusHistory: []*domain.FHIREncounterStatushistory{},
			Class:         []domain.FHIRCodeableConcept{},
			ClassHistory:  []*domain.FHIREncounterClasshistory{},
			Type:          []*domain.FHIRCodeableConcept{},
			ServiceType:   &domain.FHIRCodeableConcept{},
			Priority:      &domain.FHIRCodeableConcept{},
			Subject: &domain.FHIRReference{
				ID:        &id,
				Reference: &patientRef,
			},
			EpisodeOfCare:   []*domain.FHIRReference{},
			BasedOn:         []*domain.FHIRReference{},
			Participant:     []*domain.FHIREncounterParticipant{},
			Appointment:     []*domain.FHIRReference{},
			ActualPeriod:    &domain.FHIRPeriod{},
			Length:          &domain.FHIRDuration{},
			ReasonReference: []*domain.FHIRReference{},
			Diagnosis:       []*domain.FHIREncounterDiagnosis{},
			Account:         []*domain.FHIRReference{},
			Hospitalization: &domain.FHIREncounterHospitalization{},
			Location:        []*domain.FHIREncounterLocation{},
			ServiceProvider: &domain.FHIRReference{},
			PartOf:          &domain.FHIRReference{},
		},
	}

	type args struct {
		ctx   context.Context
		input *dto.CarePlanPayload
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create patient care plan",
			setup: func(mh *usecaseMock.Mocks) args {
				implicitRules := mock.Anything
				language := mock.Anything
				planID := uuid.NewString()
				title := mock.Anything
				description := mock.Anything
				count := 3
				planDefiniton := &domain.FHIRPlanDefinition{
					ID:            &planID,
					Language:      &language,
					ImplicitRules: &implicitRules,
					Text:          &domain.FHIRNarrative{},
					Title:         &title,
					Identifier:    []domain.FHIRIdentifier{},
					Action: []domain.PlanDefinitionAction{
						{
							Id:          &id,
							Title:       &title,
							Description: &description,
							TimingTiming: &domain.FHIRTiming{
								ID: &id,
								Repeat: &domain.FHIRTimingRepeat{
									Count: &count,
								},
							},
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().FetchPlanDefinitionByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPlanDefinition, error) {
						return planDefiniton, nil
					})
				mh.FHIR.EXPECT().PostFHIRBundle(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *r5hapifhirmodel.Bundle) (*r5hapifhirmodel.Bundle, error) {
						return nil, nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: &payload}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: &payload}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to fetch plan definition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().FetchPlanDefinitionByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPlanDefinition, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: &payload}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create FHIR task",
			setup: func(mh *usecaseMock.Mocks) args {
				implicitRules := mock.Anything
				language := mock.Anything
				planID := uuid.NewString()
				planDefiniton := &domain.FHIRPlanDefinition{
					ID:            &planID,
					Language:      &language,
					ImplicitRules: &implicitRules,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []domain.FHIRIdentifier{},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().FetchPlanDefinitionByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPlanDefinition, error) {
						return planDefiniton, nil
					})
				mh.FHIR.EXPECT().PostFHIRBundle(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *r5hapifhirmodel.Bundle) (*r5hapifhirmodel.Bundle, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: &payload}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.PatientCarePlan(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.PatientCarePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalImpl_CreatePatientCarePlan(t *testing.T) {
	data := dto.CarePlanInput{
		EncounterID:      gofakeit.UUID(),
		PlanDefinitionID: gofakeit.UUID(),
		Notes:            gofakeit.UUID(),
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	orgName := mock.Anything
	ID := uuid.NewString()
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	type args struct {
		ctx   context.Context
		input *dto.CarePlanInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create patient care plan",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientCarePlan(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.CarePlanPayload) error {
						return nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: &data}
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to create patient care plan",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientCarePlan(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.CarePlanPayload) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: &data}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facility id frm ctx",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: context.Background(), input: &data}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: &data}
			},
			wantErr: true,
		},
		{
			name: "Sad case: bad input",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: &dto.CarePlanInput{EncounterID: uuid.NewString()}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			err := clinicalUsecase.CreatePatientCarePlan(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.CreatePatientCarePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
