package specialized_test

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

func TestUseCasesClinicalImpl_CreateQuestionnaire(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	ID := gofakeit.UUID()
	input := &domain.FHIRQuestionnaire{
		ID: &ID,
		Meta: &domain.FHIRMetaInput{
			VersionID: ID,
			Source:    "",
			Tag: []domain.FHIRCodingInput{
				{
					ID:           ID,
					Version:      &ID,
					Display:      "",
					UserSelected: new(bool),
				},
			},
		},
	}

	type args struct {
		ctx                context.Context
		questionnaireInput *domain.FHIRQuestionnaire
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRQuestionnaire(mock.Anything, mock.AnythingOfType("*domain.FHIRQuestionnaire")).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
						return &domain.FHIRQuestionnaire{ID: &ID}, nil
					})
				return args{ctx: ctx, questionnaireInput: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRQuestionnaire(mock.Anything, mock.AnythingOfType("*domain.FHIRQuestionnaire")).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, questionnaireInput: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreateQuestionnaire(args.ctx, args.questionnaireInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateQuestionnaire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListQuestionnaires(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	type args struct {
		ctx        context.Context
		title      string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: list questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().ListFHIRQuestionnaire(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error) {
						return &domain.PagedFHIRQuestionnaires{}, nil
					})
				return args{ctx: ctx, title: "Cervical Cancer Screening Form", pagination: &dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to list questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().ListFHIRQuestionnaire(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, title: "Cervical Cancer Screening Form", pagination: &dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListQuestionnaires(args.ctx, args.title, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListQuestionnaires() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSpecializedImpl_GenerateRiskAssessment(t *testing.T) {
	cacx := "Cervical Cancer Screening"
	breast := "Breast Cancer Screening"
	prostate := "Prostate Cancer Screening"
	invalid := "Invalid Cancer Screening"
	type args struct {
		ctx                   context.Context
		questionnaireID       string
		questionnaireResponse *dto.QuestionnaireResponse
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: stratify patient - Cervical",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						return &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &id,
								Title: &cacx,
							},
						}, nil
					})

				return args{ctx: context.Background(), questionnaireID: gofakeit.UUID(), questionnaireResponse: &dto.QuestionnaireResponse{
					Authored: gofakeit.TimeZone(),
				}}
			},
			wantErr: false,
		},
		{
			name: "Happy case: stratify patient - Breast",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						return &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &id,
								Title: &breast,
							},
						}, nil
					})

				return args{ctx: context.Background(), questionnaireID: gofakeit.UUID(), questionnaireResponse: &dto.QuestionnaireResponse{
					Authored: gofakeit.TimeZone(),
				}}
			},
			wantErr: false,
		},
		{
			name: "Happy case: stratify patient - Prostate",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						return &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &id,
								Title: &prostate,
							},
						}, nil
					})

				return args{ctx: context.Background(), questionnaireID: gofakeit.UUID(), questionnaireResponse: &dto.QuestionnaireResponse{
					Authored: gofakeit.TimeZone(),
				}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable stratify patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRQuestionnaire(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
						return &domain.FHIRQuestionnaireRelayPayload{
							Resource: &domain.FHIRQuestionnaire{
								ID:    &id,
								Title: &invalid,
							},
						}, nil
					})

				return args{ctx: context.Background(), questionnaireID: gofakeit.UUID(), questionnaireResponse: &dto.QuestionnaireResponse{
					Authored: gofakeit.TimeZone(),
				}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GenerateRiskAssessment(args.ctx, args.questionnaireID, args.questionnaireResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpecializedImpl.GenerateRiskAssessment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSpecializedImpl_FetchQuestionnaire(t *testing.T) {
	id := gofakeit.UUID()
	type args struct {
		ctx         context.Context
		searchParam string
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: list FHIR Questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().ListFHIRQuestionnaire(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error) {
						return &domain.PagedFHIRQuestionnaires{
							Questionnaires: []domain.FHIRQuestionnaire{
								{
									ID: &id,
								},
							},
						}, nil
					})

				return args{ctx: context.Background(), searchParam: gofakeit.UUID(), pagination: &dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to list FHIR Questionnaire",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().ListFHIRQuestionnaire(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error) {
						return nil, fmt.Errorf("error")
					})

				return args{ctx: context.Background(), searchParam: gofakeit.UUID(), pagination: &dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.FetchQuestionnaire(args.ctx, args.searchParam, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpecializedImpl.FetchQuestionnaire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
