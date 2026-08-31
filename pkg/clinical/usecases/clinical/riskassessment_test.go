package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_RecordRiskAssessment(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	riskassessment := &domain.FHIRRiskAssessmentRelayPayload{
		Resource: &domain.FHIRRiskAssessment{
			ID: &ID,
		},
	}

	type args struct {
		ctx            context.Context
		riskAssessment *domain.FHIRRiskAssessmentInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record a risk assessment",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return riskassessment, nil
					})
				return args{ctx: ctx}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to create fhir risk assessment",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.RecordRiskAssessment(args.ctx, args.riskAssessment)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordRiskAssessment() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_ListRiskAssessment(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	ID := uuid.NewString()
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	filters := &dto.RiskAssessmentFilterInput{
		FilterInput: dto.FilterInput{
			PatientID:   gofakeit.UUID(),
			EncounterID: gofakeit.UUID(),
			Date: &scalarutils.Date{
				Year:  2023,
				Month: 10,
				Day:   4,
			},
		},
		ScreeningType: dto.BreastCancerScreeningTypeEnum,
		Result:        "High Risk",
	}

	firstTen := 10
	pagination := serverutils.PaginationInput{
		First: &firstTen,
	}

	screeningType := dto.BreastCancerScreeningTypeEnum.String()
	occurrenceTime := "2024-05-11T13:08:16Z"
	riskassessment := &domain.PagedFHIRRiskAssessment{
		RiskAssessment: []domain.FHIRRiskAssessment{
			{
				ID: &ID,
				Prediction: []domain.FHIRRiskAssessmentPrediction{
					{
						Outcome: &domain.FHIRCodeableConcept{
							Text: "High Risk",
						},
					},
				},
				Subject: domain.FHIRReference{
					ID:      &ID,
					Display: "Test Subject",
				},
				Encounter: &domain.FHIRReference{
					ID:      &ID,
					Display: "Tes Subject",
				},
				Text: &domain.FHIRNarrative{
					Div: scalarutils.XHTML(screeningType),
				},
				OccurrenceDateTime: &occurrenceTime,
			},
		},
	}

	type args struct {
		ctx        context.Context
		searchID   string
		filters    *dto.RiskAssessmentFilterInput
		pagination serverutils.PaginationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get risk level - from outcome",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRRiskAssessment(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRRiskAssessment, error) {
						return riskassessment, nil
					})
				return args{ctx: ctx, filters: filters, searchID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get risk level - no search ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRRiskAssessment(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRRiskAssessment, error) {
						return riskassessment, nil
					})
				return args{ctx: ctx, filters: filters, pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable to list risk assessment (invalid pagination)",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				return args{
					ctx: ctx,
					filters: &dto.RiskAssessmentFilterInput{
						FilterInput: dto.FilterInput{
							PatientID:   gofakeit.UUID(),
							EncounterID: gofakeit.UUID(),
							Date: &scalarutils.Date{
								Year:  2023,
								Month: 10,
								Day:   4,
							},
						},
						ScreeningType: dto.BreastCancerScreeningTypeEnum,
					},
					searchID: gofakeit.UUID(),
					pagination: serverutils.PaginationInput{
						First: &firstTen,
						Last:  &firstTen,
					}}
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Successfully get risk level - from qualitative risk",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRRiskAssessment(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRRiskAssessment, error) {
						return riskassessment, nil
					})
				return args{ctx: ctx, filters: filters, searchID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - invalid patient id",
			setup: func(mh *usecaseMock.Mocks) args {
				filters := &dto.RiskAssessmentFilterInput{
					FilterInput: dto.FilterInput{
						PatientID:   "id",
						EncounterID: gofakeit.UUID(),
						Date: &scalarutils.Date{
							Year:  2023,
							Month: 10,
							Day:   4,
						},
					},
					ScreeningType: dto.CervicalCancerScreeningTypeEnum,
					Result:        "Low Risk",
				}
				return args{ctx: ctx, filters: filters, searchID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid encounter id",
			setup: func(mh *usecaseMock.Mocks) args {
				filters := &dto.RiskAssessmentFilterInput{
					FilterInput: dto.FilterInput{
						PatientID:   gofakeit.UUID(),
						EncounterID: "id",
						Date: &scalarutils.Date{
							Year:  2023,
							Month: 10,
							Day:   4,
						},
					},
					ScreeningType: dto.CervicalCancerScreeningTypeEnum,
					Result:        "Low Risk",
				}
				return args{ctx: ctx, filters: filters, searchID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, filters: filters, searchID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search risk assessment",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRRiskAssessment(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRRiskAssessment, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, filters: filters, searchID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListRiskAssessment(args.ctx, args.searchID, args.filters, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetQuestionnaireResponseRiskLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalImpl_GetRiskAssessmentByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully get observation by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						riskID := gofakeit.UUID()
						occurrenceDate := gofakeit.Date().GoString()
						return &domain.FHIRRiskAssessmentRelayPayload{
							Resource: &domain.FHIRRiskAssessment{
								ID:                 &riskID,
								OccurrenceDateTime: &occurrenceDate,
								Subject: domain.FHIRReference{
									ID:      &riskID,
									Display: riskID,
								},
								Encounter: &domain.FHIRReference{
									ID:      &riskID,
									Display: riskID,
								},
								Text: &domain.FHIRNarrative{
									Div: scalarutils.XHTML(gofakeit.BeerName()),
								},
								Prediction: []domain.FHIRRiskAssessmentPrediction{
									{
										QualitativeRisk: &domain.FHIRCodeableConcept{
											Text: "High",
										},
									},
								},
							},
						}, nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), id: gofakeit.UUID()}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get observation by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRRiskAssessment(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRRiskAssessmentRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), id: gofakeit.UUID()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := c.GetRiskAssessmentByID(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.GetRiskAssessmentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
