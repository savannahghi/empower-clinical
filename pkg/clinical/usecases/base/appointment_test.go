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
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/base"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_CreateAppointment(t *testing.T) {
	id := gofakeit.UUID()

	type args struct {
		ctx         context.Context
		appointment *domain.FHIRAppointmentInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create appointment",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRAppointment(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRAppointmentInput) (*domain.FHIRAppointment, error) {
						return nil, nil
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					appointment: &domain.FHIRAppointmentInput{
						ID: &id,
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create appointment",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRAppointment(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRAppointmentInput) (*domain.FHIRAppointment, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					appointment: &domain.FHIRAppointmentInput{
						ID: &id,
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreateAppointment(args.ctx, args.appointment)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateAppointment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateCheckIn(t *testing.T) {
	patientID := gofakeit.UUID()
	patientRef := fmt.Sprintf("Patient/%s", patientID)

	type args struct {
		ctx                context.Context
		patientID          string
		date               *scalarutils.Date
		headers            *dto.AdvantageHeaders
		appointmentPayload *base.AppointmentPayload
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create check-in",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Advantage.EXPECT().GetSchedules(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error) {
						return []*dto.Schedule{
							{
								ID:            uuid.NewString(),
								CreatedByName: gofakeit.Name(),
								UpdatedByName: gofakeit.Name(),
								Queue:         mock.Anything,
								WorkstationID: uuid.NewString(),
								BranchID:      uuid.NewString(),
								DepartmentID:  uuid.NewString(),
								Active:        true,
							},
						}, nil
					})
				mh.Advantage.EXPECT().GetSlots(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, startDate, scheduleID string, headers *dto.AdvantageHeaders) ([]*dto.Slot, error) {
						return []*dto.Slot{
							{
								ID:    uuid.NewString(),
								Start: gofakeit.Date().GoString(),
								End:   gofakeit.Date().GoString(),
							},
						}, nil
					})
				mh.Advantage.EXPECT().CreateCheckin(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, checkIn *dto.Checkin, headers *dto.AdvantageHeaders) error {
						return nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{
							ID:          uuid.NewString(),
							DisplayName: gofakeit.Name(),
							URL:         gofakeit.URL(),
						}, nil
					})
				mh.FHIR.EXPECT().CreateFHIRAppointment(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRAppointmentInput) (*domain.FHIRAppointment, error) {
						return nil, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: uuid.NewString(),
					date:      &scalarutils.Date{},
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
					appointmentPayload: &base.AppointmentPayload{
						Tags: &domain.FHIRMetaInput{},
						Subject: &domain.FHIRReference{
							ID:        &patientID,
							Reference: &patientRef,
						},
						Reason: "refer for test",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get schedules",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Advantage.EXPECT().GetSchedules(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: uuid.NewString(),
					date:      &scalarutils.Date{},
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get slots",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Advantage.EXPECT().GetSchedules(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error) {
						return []*dto.Schedule{
							{
								ID:            uuid.NewString(),
								CreatedByName: gofakeit.Name(),
								UpdatedByName: gofakeit.Name(),
								Queue:         mock.Anything,
								WorkstationID: uuid.NewString(),
								BranchID:      uuid.NewString(),
								DepartmentID:  uuid.NewString(),
								Active:        true,
							},
						}, nil
					})
				mh.Advantage.EXPECT().GetSlots(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, startDate, scheduleID string, headers *dto.AdvantageHeaders) ([]*dto.Slot, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: uuid.NewString(),
					date:      &scalarutils.Date{},
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create check-in",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Advantage.EXPECT().GetSchedules(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error) {
						return []*dto.Schedule{
							{
								ID:            uuid.NewString(),
								CreatedByName: gofakeit.Name(),
								UpdatedByName: gofakeit.Name(),
								Queue:         mock.Anything,
								WorkstationID: uuid.NewString(),
								BranchID:      uuid.NewString(),
								DepartmentID:  uuid.NewString(),
								Active:        true,
							},
						}, nil
					})
				mh.Advantage.EXPECT().GetSlots(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, startDate, scheduleID string, headers *dto.AdvantageHeaders) ([]*dto.Slot, error) {
						return []*dto.Slot{
							{
								ID:    uuid.NewString(),
								Start: gofakeit.Date().GoString(),
								End:   gofakeit.Date().GoString(),
							},
						}, nil
					})
				mh.Advantage.EXPECT().CreateCheckin(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, checkIn *dto.Checkin, headers *dto.AdvantageHeaders) error {
						return fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: uuid.NewString(),
					date:      &scalarutils.Date{},
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: no schedules found",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Advantage.EXPECT().GetSchedules(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error) {
						return []*dto.Schedule{}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: uuid.NewString(),
					date:      &scalarutils.Date{},
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create appointment",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Advantage.EXPECT().GetSchedules(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error) {
						return []*dto.Schedule{
							{
								ID:            uuid.NewString(),
								CreatedByName: gofakeit.Name(),
								UpdatedByName: gofakeit.Name(),
								Queue:         mock.Anything,
								WorkstationID: uuid.NewString(),
								BranchID:      uuid.NewString(),
								DepartmentID:  uuid.NewString(),
								Active:        true,
							},
						}, nil
					})
				mh.Advantage.EXPECT().GetSlots(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, startDate, scheduleID string, headers *dto.AdvantageHeaders) ([]*dto.Slot, error) {
						return []*dto.Slot{
							{
								ID:    uuid.NewString(),
								Start: gofakeit.Date().GoString(),
								End:   gofakeit.Date().GoString(),
							},
						}, nil
					})
				mh.Advantage.EXPECT().CreateCheckin(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, checkIn *dto.Checkin, headers *dto.AdvantageHeaders) error {
						return nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{
							ID:          uuid.NewString(),
							DisplayName: gofakeit.Name(),
							URL:         gofakeit.URL(),
						}, nil
					})
				mh.FHIR.EXPECT().CreateFHIRAppointment(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRAppointmentInput) (*domain.FHIRAppointment, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: uuid.NewString(),
					date:      &scalarutils.Date{},
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
					appointmentPayload: &base.AppointmentPayload{
						Tags: &domain.FHIRMetaInput{},
						Subject: &domain.FHIRReference{
							ID:        &patientID,
							Reference: &patientRef,
						},
						Reason: "refer for test",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Advantage.EXPECT().GetSchedules(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, headers *dto.AdvantageHeaders) ([]*dto.Schedule, error) {
						return []*dto.Schedule{
							{
								ID:            uuid.NewString(),
								CreatedByName: gofakeit.Name(),
								UpdatedByName: gofakeit.Name(),
								Queue:         mock.Anything,
								WorkstationID: uuid.NewString(),
								BranchID:      uuid.NewString(),
								DepartmentID:  uuid.NewString(),
								Active:        true,
							},
						}, nil
					})
				mh.Advantage.EXPECT().GetSlots(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, startDate, scheduleID string, headers *dto.AdvantageHeaders) ([]*dto.Slot, error) {
						return []*dto.Slot{
							{
								ID:    uuid.NewString(),
								Start: gofakeit.Date().GoString(),
								End:   gofakeit.Date().GoString(),
							},
						}, nil
					})
				mh.Advantage.EXPECT().CreateCheckin(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, checkIn *dto.Checkin, headers *dto.AdvantageHeaders) error {
						return nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: uuid.NewString(),
					date:      &scalarutils.Date{},
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
					appointmentPayload: &base.AppointmentPayload{
						Tags: &domain.FHIRMetaInput{},
						Subject: &domain.FHIRReference{
							ID:        &patientID,
							Reference: &patientRef,
						},
						Reason: "refer for test",
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

			if _, err := clinicalUsecase.CreateCheckIn(args.ctx, args.patientID, args.appointmentPayload, args.date, args.headers); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateCheckIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
