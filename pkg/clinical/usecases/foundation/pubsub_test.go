package foundation_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_CreatePubsubPatient(t *testing.T) {
	ctx := context.Background()
	payload := dto.PatientPubSubMessage{
		UserID:         gofakeit.UUID(),
		ClientID:       gofakeit.UUID(),
		Name:           gofakeit.Name(),
		DateOfBirth:    time.Now(),
		Gender:         "male",
		Active:         true,
		PhoneNumber:    gofakeit.Phone(),
		OrganizationID: gofakeit.UUID(),
		FacilityID:     gofakeit.UUID(),
	}
	ID := gofakeit.UUID()
	patient := domain.PatientPayload{
		PatientRecord: &domain.FHIRPatient{
			ID: &ID,
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
		ctx     context.Context
		payload dto.PatientPubSubMessage
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create pubsub patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Twice()
				mh.FHIR.EXPECT().CreateFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
						return &patient, nil
					})
				mh.PubSub.EXPECT().NotifyPatientFHIRIDUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.UpdatePatientFHIRID) error {
						return nil
					})
				return args{ctx: ctx, payload: payload}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Twice()
				mh.FHIR.EXPECT().CreateFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, payload: payload}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to add FHIR ID to profile",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Twice()
				mh.FHIR.EXPECT().CreateFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
						return &patient, nil
					})
				mh.PubSub.EXPECT().NotifyPatientFHIRIDUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.UpdatePatientFHIRID) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, payload: payload}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get organisation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, payload: payload}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.CreatePubsubPatient(args.ctx, args.payload); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubPatient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubOrganization(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	data := dto.FacilityPubSubMessage{
		ID:          &ID,
		Name:        "Test Facility",
		Code:        0,
		Phone:       "",
		Active:      false,
		County:      "",
		Description: "",
	}
	orgName := mock.Anything
	organization := domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	type args struct {
		ctx  context.Context
		data dto.FacilityPubSubMessage
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create pubsub organization",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return &organization, nil
					})
				mh.PubSub.EXPECT().NotifyFacilityFHIRIDUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.UpdateFacilityFHIRID) error {
						return nil
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create pubsub organization",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to add fhir id to facility",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return &organization, nil
					})
				mh.PubSub.EXPECT().NotifyFacilityFHIRIDUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.UpdateFacilityFHIRID) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.CreatePubsubOrganization(args.ctx, args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubVitals(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()

	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	data := dto.VitalSignPubSubMessage{
		PatientID:      ID,
		OrganizationID: "",
		Name:           "",
		ConceptID:      new(string),
		Value:          "",
		Date:           time.Time{},
	}

	name := gofakeit.Name()
	patient := domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
			Name: []*domain.FHIRHumanName{
				{
					Given: []*string{
						&name,
					},
				},
			},
		},
	}

	observation := &domain.FHIRObservation{
		ID: &ID,
	}

	payload := &domain.Concept{
		DisplayName: "Malaria",
		URL:         gofakeit.URL(),
		ID:          gofakeit.UUID(),
	}

	type args struct {
		ctx  context.Context
		data dto.VitalSignPubSubMessage
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create pubsub vitals",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Twice()
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create pubsub vitals - available organizationID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Twice()
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create pubsub vitals with facilityID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Times(3)
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				data.FacilityID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to create pubsub vitals with facilityID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Times(3)
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.FacilityID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to find patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to find organisation using org ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Times(3)
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.CreatePubsubVitals(args.ctx, args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubVitals() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubTestResult(t *testing.T) {
	ctx := context.Background()

	ID := uuid.NewString()
	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	data := dto.PatientTestResultPubSubMessage{
		PatientID:      uuid.New().String(),
		OrganizationID: "",
		FacilityID:     "",
		Name:           "",
		ConceptID:      new(string),
		Date:           time.Time{},
		Result:         dto.TestResult{},
	}

	name := gofakeit.Name()
	patient := domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
			Name: []*domain.FHIRHumanName{
				{
					Given: []*string{
						&name,
					},
				},
			},
		},
	}

	observation := &domain.FHIRObservation{
		ID: &ID,
	}

	payload := &domain.Concept{
		DisplayName: "Malaria",
		URL:         gofakeit.URL(),
		ID:          gofakeit.UUID(),
	}

	type args struct {
		ctx  context.Context
		data dto.PatientTestResultPubSubMessage
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create test result",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Twice()
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create test result - with organisation ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Twice()
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create pubsub vitals with facilityID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Times(3)
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				data.FacilityID = uuid.NewString()
				data.Result = dto.TestResult{
					Name:      "",
					ConceptID: new(string),
				}
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail create pubsub vitals with facilityID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.FacilityID = uuid.NewString()
				data.Result = dto.TestResult{
					Name:      "",
					ConceptID: new(string),
				}
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get fhir patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get organisation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to create observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.CreatePubsubTestResult(args.ctx, args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubTestResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubMedicationStatement(t *testing.T) {
	ctx := context.Background()
	conceptID := "12345"

	ID := uuid.NewString()
	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	data := dto.MedicationPubSubMessage{
		PatientID:      uuid.New().String(),
		OrganizationID: "",
		FacilityID:     gofakeit.UUID(),
		Name:           "",
		ConceptID:      &conceptID,
		Date:           time.Time{},
		Value:          "",
		Drug: &dto.MedicationDrug{
			ConceptID: &conceptID,
		},
	}

	name := gofakeit.Name()
	patient := domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
			Name: []*domain.FHIRHumanName{
				{
					Given: []*string{
						&name,
					},
				},
			},
		},
	}

	medication := &domain.FHIRMedicationStatementRelayPayload{
		Resource: &domain.FHIRMedicationStatement{
			ID: &ID,
		},
	}

	payload := &domain.Concept{
		DisplayName: "Malaria",
		URL:         gofakeit.URL(),
		ID:          gofakeit.UUID(),
	}

	type args struct {
		ctx  context.Context
		data dto.MedicationPubSubMessage
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create medication statement",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Twice()
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedicationStatement(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
						return medication, nil
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create medication statement - with organisation ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Twice()
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedicationStatement(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
						return medication, nil
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create medication statement - with facilityID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Twice()
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedicationStatement(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
						return medication, nil
					})
				data.OrganizationID = uuid.NewString()
				data.FacilityID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to create medication statement with facilityID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Twice()
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedicationStatement(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.FacilityID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Twice()
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get organisation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Twice()
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to create medication statement",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Twice()
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedicationStatement(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				data.OrganizationID = uuid.NewString()
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.CreatePubsubMedicationStatement(args.ctx, args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubMedicationStatement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubAllergyIntolerance(t *testing.T) {
	ctx := context.Background()

	ID := uuid.NewString()
	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	data := dto.PatientAllergyPubSubMessage{
		PatientID:      uuid.New().String(),
		OrganizationID: "",
		Name:           "",
		ConceptID:      &ID,
		Date:           time.Time{},
		Reaction:       dto.AllergyReaction{},
		Severity: dto.AllergySeverity{
			ConceptID: &ID,
		},
	}

	name := gofakeit.Name()
	patient := domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
			Name: []*domain.FHIRHumanName{
				{
					Given: []*string{
						&name,
					},
				},
			},
		},
	}
	allergyintolerance := &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: &domain.FHIRAllergyIntolerance{
			ID: &ID,
		},
	}
	payload := &domain.Concept{
		DisplayName: "Malaria",
		URL:         gofakeit.URL(),
		ID:          gofakeit.UUID(),
	}

	type args struct {
		ctx  context.Context
		data dto.PatientAllergyPubSubMessage
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create allergy intolerance",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Times(3)
				mh.FHIR.EXPECT().CreateFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
						return allergyintolerance, nil
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create allergy with reaction",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					}).Times(3)
				mh.FHIR.EXPECT().CreateFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
						return allergyintolerance, nil
					})

				data.FacilityID = uuid.NewString()
				data.Reaction = dto.AllergyReaction{
					Name:      "",
					ConceptID: new(string),
				}
				data.Severity = dto.AllergySeverity{
					Name:      "",
					ConceptID: new(string),
				}
				return args{ctx: ctx, data: data}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail create allergy with reaction",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().CreateFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				data.FacilityID = uuid.NewString()
				data.Reaction = dto.AllergyReaction{
					Name:      "",
					ConceptID: new(string),
				}
				data.Severity = dto.AllergySeverity{
					Name:      "",
					ConceptID: new(string),
				}
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create allergy intolerance",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().CreateFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get organisation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: data}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.CreatePubsubAllergyIntolerance(args.ctx, args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_getConcept(t *testing.T) {
	ctx := context.Background()
	payload := &domain.Concept{
		ConceptClass:  "DRUG",
		DataType:      "TEXT",
		DisplayLocale: "en",
		DisplayName:   "Malaria",
	}
	type args struct {
		ctx       context.Context
		source    domain.TerminologySource
		conceptID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Sad case: failed to get icd10 concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occured")
					})
				return args{ctx: ctx, source: domain.TerminologySourceICD10, conceptID: gofakeit.BS()}
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get ciel concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occured")
					})
				return args{ctx: ctx, source: domain.TerminologySourceICD10, conceptID: gofakeit.BS()}
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to get loinc concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occured")
					})
				return args{ctx: ctx, source: domain.TerminologySourceICD10, conceptID: gofakeit.BS()}
			},
			wantErr: true,
		},

		{
			name: "Sad case: invalid concept source",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, source: domain.TerminologySource("invalid"), conceptID: gofakeit.BS()}
			},
			wantErr: true,
		},
		{
			name: "Happy Case: Successfully return who content",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				return args{ctx: ctx, source: domain.TerminologySourceICD10WHO, conceptID: gofakeit.BS()}
			},
			wantErr: false,
		},
		{
			name: "Happy Case: Successfully return loinc content",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				return args{ctx: ctx, source: domain.TerminologySourceLOINC, conceptID: gofakeit.BS()}
			},
			wantErr: false,
		},
		{
			name: "Happy Case: Successfully return ciel content",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				return args{ctx: ctx, source: domain.TerminologySourceCIEL, conceptID: gofakeit.BS()}
			},
			wantErr: false,
		},
		{
			name: "Happy Case: Successfully return icd 11 content",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return payload, nil
					})
				return args{ctx: ctx, source: domain.TerminologySourceICD11WHO, conceptID: gofakeit.BS()}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetConcept(args.ctx, args.source, args.conceptID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.getConcept() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubTenant(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	orgInput := dto.OrganizationInput{
		Name:        "test",
		PhoneNumber: "test",
		Identifiers: []dto.OrganizationIdentifier{
			{
				Value: gofakeit.UUID(),
				Type:  dto.OrganizationIdentifierType("MCHProgram"),
			},
		},
	}

	orgName := mock.Anything
	active := true
	organization := domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:         &ID,
			Name:       &orgName,
			Active:     &active,
			Identifier: []*domain.FHIRIdentifier{},
		},
	}

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
			name: "Happy case: create tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return &organization, nil
					})
				mh.PubSub.EXPECT().NotifyProgramFHIRIDUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.UpdateProgramFHIRID) error {
						return nil
					})
				return args{ctx: ctx, data: orgInput}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create tenant",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: orgInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to update fhir patient id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
						return &organization, nil
					})
				mh.PubSub.EXPECT().NotifyProgramFHIRIDUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.UpdateProgramFHIRID) error {
						return fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, data: orgInput}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.CreatePubsubTenant(args.ctx, args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubTenant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
