package advantage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/advantage"
	authMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/authutils/mock"
)

type mockHandler struct {
	auth *authMock.MockOAuthClientService
}

func TestServiceAdvantageImpl_PatientSegmentation(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload dto.SegmentationPayload
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: segment patients",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return &authutils.OAUTHResponse{
							AccessToken: gofakeit.UUID(),
						}, nil
					})

				return args{
					ctx: context.Background(),
					payload: dto.SegmentationPayload{
						ClinicalID:   gofakeit.UUID(),
						SegmentLabel: dto.SegmentationCategoryHighRiskPositive,
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to segment patients",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: context.Background(),
					payload: dto.SegmentationPayload{
						ClinicalID:   gofakeit.UUID(),
						SegmentLabel: dto.SegmentationCategoryHighRiskPositive,
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAuth := authMock.NewMockOAuthClientService(t)
			s := advantage.NewServiceAdvantage(fakeAuth)

			args := tt.setup(&mockHandler{auth: fakeAuth})

			if err := s.SegmentPatient(args.ctx, args.payload); (err != nil) != tt.wantErr {
				t.Errorf("ServiceAdvantageImpl.SegmentPatient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceAdvantageImpl_SendSMS(t *testing.T) {
	type args struct {
		ctx           context.Context
		payload       dto.SMSPayload
		workstationID string
		branchID      string
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: send SMS",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return &authutils.OAUTHResponse{
							AccessToken: gofakeit.UUID(),
						}, nil
					})

				return args{
					ctx: context.Background(),
					payload: dto.SMSPayload{
						Intention:  "DIRECT_MESSAGE",
						Message:    "message",
						Recipients: []string{},
					},
					workstationID: gofakeit.UUID(),
					branchID:      gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to send SMS",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: context.Background(),
					payload: dto.SMSPayload{
						Intention:  "DIRECT_MESSAGE",
						Message:    "message",
						Recipients: []string{},
					},
					workstationID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAuth := authMock.NewMockOAuthClientService(t)
			s := advantage.NewServiceAdvantage(fakeAuth)

			args := tt.setup(&mockHandler{auth: fakeAuth})

			if err := s.SendSMS(args.ctx, args.workstationID, args.branchID, args.payload); (err != nil) != tt.wantErr {
				t.Errorf("ServiceAdvantageImpl.SendSMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceAdvantageImpl_GetSchedules(t *testing.T) {
	type args struct {
		ctx     context.Context
		headers *dto.AdvantageHeaders
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *dto.Schedule
		wantErr bool
	}{
		{
			name: "Happy case: get schedule",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return &authutils.OAUTHResponse{
							AccessToken: gofakeit.UUID(),
						}, nil
					})

				return args{
					ctx: context.Background(),
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get schedule",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: context.Background(),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAuth := authMock.NewMockOAuthClientService(t)
			s := advantage.NewServiceAdvantage(fakeAuth)

			args := tt.setup(&mockHandler{auth: fakeAuth})

			_, err := s.GetSchedules(args.ctx, args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceAdvantageImpl.GetSchedules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceAdvantageImpl_GetSlots(t *testing.T) {
	type args struct {
		ctx        context.Context
		startDate  string
		scheduleID string
		headers    *dto.AdvantageHeaders
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *dto.Slot
		wantErr bool
	}{
		{
			name: "Happy case: get available slots",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return &authutils.OAUTHResponse{
							AccessToken: gofakeit.UUID(),
						}, nil
					})

				return args{
					ctx:        context.Background(),
					startDate:  "2024-04-18",
					scheduleID: gofakeit.UUID(),
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get available slots",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:        context.Background(),
					startDate:  "2024-04-18",
					scheduleID: gofakeit.UUID(),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAuth := authMock.NewMockOAuthClientService(t)
			s := advantage.NewServiceAdvantage(fakeAuth)

			args := tt.setup(&mockHandler{auth: fakeAuth})

			_, err := s.GetSlots(args.ctx, args.startDate, args.scheduleID, args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceAdvantageImpl.GetSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceAdvantageImpl_CreateCheckin(t *testing.T) {
	id := gofakeit.UUID()
	type args struct {
		ctx     context.Context
		checkIn *dto.Checkin
		headers *dto.AdvantageHeaders
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		wantErr bool
	}{
		{
			name: "Happy case: create checkin successfully",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return &authutils.OAUTHResponse{
							AccessToken: gofakeit.UUID(),
						}, nil
					})

				return args{
					ctx: context.Background(),
					checkIn: &dto.Checkin{
						Slot:    gofakeit.UUID(),
						Start:   "2024-04-18 15:20:59",
						End:     time.Now().Format("2006-01-02T15:04:05+03:00"),
						Patient: id,
					},
					headers: &dto.AdvantageHeaders{
						Organisation: gofakeit.UUID(),
						Cluster:      gofakeit.UUID(),
						Department:   gofakeit.UUID(),
						Branch:       gofakeit.UUID(),
						Workstation:  gofakeit.UUID(),
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create checkin",
			setup: func(mh *mockHandler) args {
				mh.auth.EXPECT().Authenticate().
					RunAndReturn(func() (*authutils.OAUTHResponse, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: context.Background(),
					checkIn: &dto.Checkin{
						Slot:    gofakeit.UUID(),
						Start:   "2024-04-18 15:20:59",
						End:     time.Now().Format("2006-01-02T15:04:05+03:00"),
						Patient: id,
					},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAuth := authMock.NewMockOAuthClientService(t)
			s := advantage.NewServiceAdvantage(fakeAuth)

			args := tt.setup(&mockHandler{auth: fakeAuth})

			if err := s.CreateCheckin(args.ctx, args.checkIn, args.headers); (err != nil) != tt.wantErr {
				t.Errorf("ServiceAdvantageImpl.CreateCheckin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
