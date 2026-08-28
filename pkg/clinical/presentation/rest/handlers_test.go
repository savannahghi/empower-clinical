package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	infraMocks "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/mock"
	advantageMocks "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/advantage/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/presentation/rest"
	restMocks "github.com/savannahghi/empower-clinical/pkg/clinical/presentation/rest/mock"
)

func newGinTestCtx(
	method string,
	qs map[string]string,
	pathParams []gin.Param,
	body []byte,
	headers map[string]string,
) (*gin.Context, *httptest.ResponseRecorder) {
	if method == "" {
		method = http.MethodGet
	}

	req := httptest.NewRequest(method, "/", bytes.NewReader(body))
	if qs != nil {
		q := req.URL.Query()
		for k, v := range qs {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)

	c.Request = req

	c.Params = pathParams

	return c, rec
}

func multipartCtx(encID string, addFile bool) (*gin.Context, *httptest.ResponseRecorder, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	_ = w.WriteField("encounterID", encID)
	_ = w.WriteField("serviceRequestID", encID)

	if addFile {
		fw, _ := w.CreateFormFile("file", "test.txt")
		io.Copy(fw, strings.NewReader("test content"))
	}
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	return ctx, rec, w.FormDataContentType()
}

func TestPresentationHandlersImpl_ListRiskAssessment(t *testing.T) {
	today := time.Now().Format(time.DateOnly)

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: invalid date format",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"date": "99-99-9999"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid limit parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"first": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"last": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid screeningType enum",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"screeningType": "NOPE"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Usecase Returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ListRiskAssessment(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.RiskAssessmentFilterInput"), mock.AnythingOfType("serverutils.PaginationInput")).
					RunAndReturn(func(ctx context.Context, bundleID string, filter *dto.RiskAssessmentFilterInput, pagination serverutils.PaginationInput) (*dto.RiskAssessmentConnection, error) {
						return nil, errors.New("failed to get risk assessment")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully get a risk assessment",
			setup: func(uc *restMocks.Mockusecases) args {
				patientID, encounterID := uuid.NewString(), uuid.NewString()
				after := "aBC123"
				first := 10

				uc.EXPECT().ListRiskAssessment(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.RiskAssessmentFilterInput"), mock.AnythingOfType("serverutils.PaginationInput")).
					RunAndReturn(func(ctx context.Context, searchID string, filter *dto.RiskAssessmentFilterInput, pagination serverutils.PaginationInput) (*dto.RiskAssessmentConnection, error) {
						return &dto.RiskAssessmentConnection{
							TotalCount: 10,
						}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"date":          today,
					"patientID":     patientID,
					"encounterID":   encounterID,
					"first":         strconv.Itoa(first),
					"searchID":      patientID,
					"after":         after,
					"screeningType": string(dto.BreastCancerScreeningTypeEnum),
					"result":        "POSITIVE",
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.ListRiskAssessment(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListTasks(t *testing.T) {
	count := 10

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: Invalid date format",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"date": "99-99-9999"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid limit parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"first": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"last": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid task status enum",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"status": "NOPE"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Usecase Returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ListTasks(mock.Anything, mock.AnythingOfType("*dto.TaskFilterInput"), mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, filter *dto.TaskFilterInput, pagination dto.Pagination) (*dto.TaskOutputConnection, error) {
						return nil, errors.New("failed to get risk assessment")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully get a task",
			setup: func(uc *restMocks.Mockusecases) args {
				patientID, encounterID := uuid.NewString(), uuid.NewString()
				after := "aBC123"
				first := 10
				today := time.Now().Format(time.DateOnly)

				uc.EXPECT().ListTasks(mock.Anything, mock.AnythingOfType("*dto.TaskFilterInput"), mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, filter *dto.TaskFilterInput, pagination dto.Pagination) (*dto.TaskOutputConnection, error) {
						return &dto.TaskOutputConnection{TotalCount: &count}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patientID":   patientID,
					"encounterID": encounterID,
					"first":       strconv.Itoa(first),
					"after":       after,
					"status":      string(dto.CompletedTasksStatus),
					"type":        "any",
					"date":        today,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.ListTasks(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_GenerateReferralReport(t *testing.T) {

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: Fail to generate referral report",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().GenerateReferralReportPDF(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, serviceRequestID string) ([]byte, error) {
						return nil, errors.New("no pdf generated")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"servicerequest": gofakeit.UUID(),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully generate referral report",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().GenerateReferralReportPDF(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, serviceRequestID string) ([]byte, error) {
						return []byte("fake pdf"), nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"servicerequest": gofakeit.UUID(),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.GenerateReferralReport(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ReferPatient(t *testing.T) {

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: Fail to refer patient",
			setup: func(uc *restMocks.Mockusecases) args {
				in := &dto.ReferralInput{}
				b, _ := json.Marshal(in)
				uc.EXPECT().ReferPatient(mock.Anything, mock.AnythingOfType("*dto.ReferralInput")).
					RunAndReturn(func(ctx context.Context, input *dto.ReferralInput) (*dto.ServiceRequest, error) {
						return nil, errors.New("fail to refer patient")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, b, nil)

				return args{ctx: ctx, rec: rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Invalid json",
			setup: func(_ *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{bad}`), nil)

				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully refer patient",
			setup: func(uc *restMocks.Mockusecases) args {
				in := &dto.ReferralInput{}
				out := &dto.ServiceRequest{ID: "SR123"}
				body, _ := json.Marshal(in)
				uc.EXPECT().ReferPatient(mock.Anything, mock.AnythingOfType("*dto.ReferralInput")).
					RunAndReturn(func(ctx context.Context, input *dto.ReferralInput) (*dto.ServiceRequest, error) {
						return out, nil
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{c, r}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.ReferPatient(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateTestOrder(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: Invalid json",
			setup: func(uc *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{bad}`), nil)

				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Fail to create test order",
			setup: func(uc *restMocks.Mockusecases) args {
				input := &dto.TestOrder{}
				body, _ := json.Marshal(input)
				uc.EXPECT().CreateTestOrder(mock.Anything, mock.AnythingOfType("*dto.TestOrder")).
					RunAndReturn(func(ctx context.Context, input *dto.TestOrder) (*dto.ServiceRequest, error) {
						return nil, errors.New("failed to create test order")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx: ctx, rec: rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully create test order",
			setup: func(uc *restMocks.Mockusecases) args {
				in := &dto.TestOrder{}
				out := &dto.ServiceRequest{ID: "SR123"}
				body, _ := json.Marshal(in)
				uc.EXPECT().CreateTestOrder(mock.Anything, mock.AnythingOfType("*dto.TestOrder")).
					RunAndReturn(func(ctx context.Context, input *dto.TestOrder) (*dto.ServiceRequest, error) {
						return out, nil
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{c, r}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreateTestOrder(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RecordTreatmentEnrollment(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: Invalid json",
			setup: func(uc *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{bad}`), nil)

				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Fail to link diagnosis to treatment",
			setup: func(uc *restMocks.Mockusecases) args {
				input := &dto.TreatmentEnrollmentInput{EncounterID: "encounter-123"}
				body, _ := json.Marshal(input)
				uc.EXPECT().RecordTreatmentEnrollment(mock.Anything, mock.AnythingOfType("*dto.TreatmentEnrollmentInput")).
					RunAndReturn(func(ctx context.Context, input *dto.TreatmentEnrollmentInput) (*dto.Condition, error) {
						return nil, errors.New("failed to link diagnosis to treatment")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx: ctx, rec: rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Happy Case: Successfully link diagnosis to treatment",
			setup: func(uc *restMocks.Mockusecases) args {
				in := &dto.TreatmentEnrollmentInput{EncounterID: "encounter-123"}
				out := &dto.Condition{ID: "condition-123"}
				body, _ := json.Marshal(in)
				uc.EXPECT().RecordTreatmentEnrollment(mock.Anything, mock.AnythingOfType("*dto.TreatmentEnrollmentInput")).
					RunAndReturn(func(ctx context.Context, input *dto.TreatmentEnrollmentInput) (*dto.Condition, error) {
						return out, nil
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{c, r}
			},
			want: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.RecordTreatmentEnrollment(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_UpdateTreatmentEnrollment(t *testing.T) {
	pathParams := []gin.Param{{Key: "id", Value: gofakeit.UUID()}}

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: Invalid json",
			setup: func(uc *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPatch, nil, pathParams, []byte(`{bad}`), nil)

				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Fail to update treatment enrollment",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().UpdateTreatmentEnrollment(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.UpdateTreatmentEnrollmentInput")).
					RunAndReturn(func(ctx context.Context, id string, input *dto.UpdateTreatmentEnrollmentInput) (*dto.Condition, error) {
						return nil, errors.New("failed to update treatment enrollment")
					})

				c, r := newGinTestCtx(http.MethodPatch, nil, pathParams, []byte(`{"date":"2027-01-15"}`), nil)

				return args{c, r}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Happy Case: Successfully update treatment enrollment",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().UpdateTreatmentEnrollment(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.UpdateTreatmentEnrollmentInput")).
					RunAndReturn(func(ctx context.Context, id string, input *dto.UpdateTreatmentEnrollmentInput) (*dto.Condition, error) {
						return &dto.Condition{ID: "condition-123"}, nil
					})

				c, r := newGinTestCtx(http.MethodPatch, nil, pathParams, []byte(`{"condition":{"code":"2A22","display":"Updated"},"date":"2027-01-15","enrollment_date":"2027-02-20"}`), nil)

				return args{c, r}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.UpdateTreatmentEnrollment(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ShareReferralForm(t *testing.T) {

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: Fail to share referral form",
			setup: func(uc *restMocks.Mockusecases) args {
				in := &dto.ShareReferralFormInput{}
				b, _ := json.Marshal(in)
				uc.EXPECT().ShareReferralForm(mock.Anything, mock.AnythingOfType("*dto.ShareReferralFormInput")).
					RunAndReturn(func(ctx context.Context, input *dto.ShareReferralFormInput) (bool, error) {
						return false, errors.New("fail to share referral form")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, b, nil)

				return args{ctx: ctx, rec: rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Invalid json",
			setup: func(_ *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{bad}`), nil)

				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully share referral form",
			setup: func(uc *restMocks.Mockusecases) args {
				in := &dto.ShareReferralFormInput{}
				body, _ := json.Marshal(in)
				uc.EXPECT().ShareReferralForm(mock.Anything, mock.AnythingOfType("*dto.ShareReferralFormInput")).
					RunAndReturn(func(ctx context.Context, input *dto.ShareReferralFormInput) (bool, error) {
						return true, nil
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{c, r}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.ShareReferralForm(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_LoadQuestionnaire(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: Invalid json",
			setup: func(uc *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{bad}`), nil)

				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Fail to load questionnaire",
			setup: func(uc *restMocks.Mockusecases) args {
				input := &domain.FHIRQuestionnaire{}
				body, _ := json.Marshal(input)
				uc.EXPECT().CreateQuestionnaire(mock.Anything, mock.AnythingOfType("*domain.FHIRQuestionnaire")).
					RunAndReturn(func(ctx context.Context, questionnaireInput *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
						return nil, errors.New("failed to create test order")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx: ctx, rec: rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully load questionnaire",
			setup: func(uc *restMocks.Mockusecases) args {
				id := uuid.NewString()
				in := &domain.FHIRQuestionnaire{ID: &id}
				out := &domain.FHIRQuestionnaire{ID: in.ID}
				b, _ := json.Marshal(in)

				uc.EXPECT().CreateQuestionnaire(mock.Anything, mock.AnythingOfType("*domain.FHIRQuestionnaire")).
					RunAndReturn(func(ctx context.Context, questionnaireInput *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
						return out, nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, b, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.LoadQuestionnaire(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListQuestionnaire(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: Fail to list questionnaire",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ListQuestionnaires(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParam string, pagination *dto.Pagination) (*dto.Questionnaire, error) {
						return nil, errors.New("fetch failed")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"searchParam": "covid"}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully list questionnaires",
			setup: func(uc *restMocks.Mockusecases) args {
				searchParam := "test"
				id := uuid.New().String()

				uc.EXPECT().ListQuestionnaires(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParam string, pagination *dto.Pagination) (*dto.Questionnaire, error) {
						return &dto.Questionnaire{
							ID: id,
						}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"searchParam": searchParam,
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.ListQuestionnaire(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RegisterFacility(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: invalid JSON",
			setup: func(_ *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{bad}`), nil)
				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: use-case error",
			setup: func(uc *restMocks.Mockusecases) args {
				in := dto.OrganizationInput{Name: "Bad Org"}
				body, _ := json.Marshal(in)

				uc.EXPECT().RegisterFacility(mock.Anything, mock.AnythingOfType("dto.OrganizationInput")).
					RunAndReturn(func(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error) {
						return nil, errors.New("db down")
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)
				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully create facility",
			setup: func(uc *restMocks.Mockusecases) args {
				in := dto.OrganizationInput{Name: "Good Org"}
				out := &dto.Organization{ID: uuid.New().String(), Name: gofakeit.Name()}
				body, _ := json.Marshal(in)

				uc.EXPECT().RegisterFacility(mock.Anything, mock.AnythingOfType("dto.OrganizationInput")).
					RunAndReturn(func(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error) {
						return out, nil
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)
				return args{c, r}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.RegisterFacility(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RegisterTenant(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: invalid JSON",
			setup: func(_ *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{bad}`), nil)
				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: use-case error",
			setup: func(uc *restMocks.Mockusecases) args {
				in := dto.OrganizationInput{Name: "Bad Org"}
				body, _ := json.Marshal(in)

				uc.EXPECT().RegisterTenant(mock.Anything, mock.AnythingOfType("dto.OrganizationInput")).
					RunAndReturn(func(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error) {
						return nil, errors.New("db down")
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)
				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully create Tenant",
			setup: func(uc *restMocks.Mockusecases) args {
				in := dto.OrganizationInput{Name: "Good Org"}
				out := &dto.Organization{ID: uuid.New().String(), Name: gofakeit.Name()}
				body, _ := json.Marshal(in)

				uc.EXPECT().RegisterTenant(mock.Anything, mock.AnythingOfType("dto.OrganizationInput")).
					RunAndReturn(func(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error) {
						return out, nil
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)
				return args{c, r}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.RegisterTenant(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_EncounterAssociatedResources(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: missing id parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"id": ""}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: encounterID},
				}

				uc.EXPECT().GetEncounterAssociatedResources(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResourceOutput, error) {
						return nil, errors.New("failed to get encounter associated resources")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully get encounter associated resources",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: encounterID},
				}

				uc.EXPECT().GetEncounterAssociatedResources(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResourceOutput, error) {
						return &dto.EncounterAssociatedResourceOutput{}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.EncounterAssociatedResources(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_DeletePatient(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: Usecase Returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				patientID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: patientID},
				}

				uc.EXPECT().DeletePatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (bool, error) {
						return false, errors.New("failed to get delete patient")
					})

				ctx, rec := newGinTestCtx(http.MethodDelete, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully delete a patient",
			setup: func(uc *restMocks.Mockusecases) args {
				patientID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: patientID},
				}

				uc.EXPECT().DeletePatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (bool, error) {
						return true, nil
					})

				ctx, rec := newGinTestCtx(http.MethodDelete, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.DeletePatient(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ScreeningReport(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().GetScreeningReport(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string, status domain.ServiceRequestStatusEnum) (*dto.ScreeningReport, error) {
						return nil, errors.New("failed to get screening report")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully get a risk assessment",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				uc.EXPECT().GetScreeningReport(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string, status domain.ServiceRequestStatusEnum) (*dto.ScreeningReport, error) {
						return &dto.ScreeningReport{}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"encounterID": encounterID,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.ScreeningReport(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_EndEncounter(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: encounterID},
				}

				uc.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return false, errors.New("failed to end encounter")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully end an encounter",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: encounterID},
				}

				uc.EXPECT().EndEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return true, nil
					})

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.EndEncounter(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_EndScreening(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: missing id parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: encounterID},
				}

				uc.EXPECT().EndScreening(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return false, errors.New("failed to end screening")
					})

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully end a screening",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: encounterID},
				}

				uc.EXPECT().EndScreening(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (bool, error) {
						return true, nil
					})

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.EndScreening(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListPatientEncounters(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: invalid limit parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"limit": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"last": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Unable to list encounters",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ListPatientEncounters(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.Pagination")).
					RunAndReturn(func(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.EncounterConnection, error) {
						return nil, errors.New("failed to fetch encounters")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Successfully get a list encounters",
			setup: func(uc *restMocks.Mockusecases) args {
				id := uuid.NewString()
				after := "aBC123"
				first := 10

				uc.EXPECT().ListPatientEncounters(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.Pagination")).
					RunAndReturn(func(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.EncounterConnection, error) {
						return &dto.EncounterConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id": id,
					"first":      strconv.Itoa(first),
					"after":      after,
					"limit":      "10",
					"before":     id,
					"last":       "10",
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.ListPatientEncounters(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_GetObservations(t *testing.T) {
	today := time.Now().Format(time.DateOnly)

	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: Successfully get observations (forward pagination)",
			setup: func(uc *restMocks.Mockusecases) args {
				id := uuid.NewString()
				concept := dto.ObservationConceptEnumBMI.String()
				first := 10

				uc.EXPECT().ConceptMapper(mock.Anything).
					RunAndReturn(func(concept dto.ObservationConceptEnum) (string, string, error) {
						return dto.ObservationConceptEnumBMI.String(), "vital-signs", nil
					})

				uc.EXPECT().GetPatientObservations(mock.Anything, mock.AnythingOfType("*dto.FetchObservationPayload")).
					RunAndReturn(func(ctx context.Context, payload *dto.FetchObservationPayload) (*dto.ObservationConnection, error) {
						return &dto.ObservationConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id":   id,
					"encounter_id": id,
					"date":         today,
					"first":        strconv.Itoa(first),
					"searchID":     id,
					"after":        id,
					"concept":      concept,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Happy Case: Successfully get observations (backward pagination)",
			setup: func(uc *restMocks.Mockusecases) args {
				id := uuid.NewString()
				concept := dto.ObservationConceptEnumBMI.String()
				first := 10

				uc.EXPECT().ConceptMapper(mock.Anything).
					RunAndReturn(func(concept dto.ObservationConceptEnum) (string, string, error) {
						return dto.ObservationConceptEnumBMI.String(), "vital-signs", nil
					})

				uc.EXPECT().GetPatientObservations(mock.Anything, mock.AnythingOfType("*dto.FetchObservationPayload")).
					RunAndReturn(func(ctx context.Context, payload *dto.FetchObservationPayload) (*dto.ObservationConnection, error) {
						return &dto.ObservationConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id":   id,
					"encounter_id": id,
					"date":         today,
					"searchID":     id,
					"before":       id,
					"last":         strconv.Itoa(first),
					"concept":      concept,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: invalid date parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"date": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"last": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: unable to map concept",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ConceptMapper(mock.Anything).
					RunAndReturn(func(concept dto.ObservationConceptEnum) (string, string, error) {
						return "", "vital-signs", fmt.Errorf("unable to map concept")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"concept": dto.ObservationConceptEnumBMI.String()}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: unable to get obs",
			setup: func(uc *restMocks.Mockusecases) args {
				id := uuid.NewString()
				concept := dto.ObservationConceptEnumBMI.String()
				first := 10

				uc.EXPECT().ConceptMapper(mock.Anything).
					RunAndReturn(func(concept dto.ObservationConceptEnum) (string, string, error) {
						return dto.ObservationConceptEnumBMI.String(), "vital-signs", nil
					})

				uc.EXPECT().GetPatientObservations(mock.Anything, mock.AnythingOfType("*dto.FetchObservationPayload")).
					RunAndReturn(func(ctx context.Context, payload *dto.FetchObservationPayload) (*dto.ObservationConnection, error) {
						return nil, fmt.Errorf("error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id":   id,
					"encounter_id": id,
					"date":         today,
					"first":        strconv.Itoa(first),
					"limit":        strconv.Itoa(first),
					"searchID":     id,
					"before":       id,
					"last":         strconv.Itoa(first),
					"concept":      concept,
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.GetObservations(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_GetTaskByID(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully fetch task by ID",
			setup: func(uc *restMocks.Mockusecases) args {
				taskID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: taskID},
				}

				task := &dto.TaskOutput{
					ID: taskID,
				}

				uc.EXPECT().GetTaskByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, taskID string) (*dto.TaskOutput, error) {
						return task, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Missing task ID in path",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Error from GetTaskByID",
			setup: func(uc *restMocks.Mockusecases) args {
				taskID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: taskID},
				}

				uc.EXPECT().GetTaskByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, taskID string) (*dto.TaskOutput, error) {
						return nil, fmt.Errorf("task not found")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.GetTaskByID(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListMedicationRequests(t *testing.T) {
	today := time.Now().Format(time.DateOnly)

	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully list medication requests",
			setup: func(uc *restMocks.Mockusecases) args {
				id := uuid.NewString()
				first := 10

				uc.EXPECT().ListMedicationRequests(mock.Anything, mock.AnythingOfType("*dto.MedicationRequestFilterInput"), mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, filter *dto.MedicationRequestFilterInput, pagination dto.Pagination) (*dto.MedicationRequestConnection, error) {
						return &dto.MedicationRequestConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patientID":   id,
					"encounterID": id,
					"date":        today,
					"status":      "active",
					"first":       strconv.Itoa(first),
					"last":        strconv.Itoa(first),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: invalid first parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"first": "abc"}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"last": "xyz"}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid date",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"date": "xyz"}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid status parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"status": "UNKNOWN"}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: ListMedicationRequests usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				id := uuid.NewString()
				first := 10

				uc.EXPECT().ListMedicationRequests(mock.Anything, mock.AnythingOfType("*dto.MedicationRequestFilterInput"), mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, filter *dto.MedicationRequestFilterInput, pagination dto.Pagination) (*dto.MedicationRequestConnection, error) {
						return nil, fmt.Errorf("usecase error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patientID":   id,
					"encounterID": id,
					"date":        today,
					"first":       strconv.Itoa(first),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.ListMedicationRequests(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_GetAllergyIntolerance(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully fetch allergy intolerance",
			setup: func(uc *restMocks.Mockusecases) args {
				alleryIntolerranceID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: alleryIntolerranceID},
				}

				task := &dto.Allergy{
					ID: alleryIntolerranceID,
				}

				uc.EXPECT().GetAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.Allergy, error) {
						return task, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Missing allergyIntolerance ID in path",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Unable to get allergy intolerance",
			setup: func(uc *restMocks.Mockusecases) args {
				allergyIntoleranceID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: allergyIntoleranceID},
				}

				uc.EXPECT().GetAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.Allergy, error) {
						return nil, fmt.Errorf("allergy intolerance not found")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.GetAllergyIntolerance(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_SearchAllergy(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully search allergy",
			setup: func(uc *restMocks.Mockusecases) args {
				limit := 5
				id := gofakeit.UUID()
				allergyName := "penicillin"

				uc.EXPECT().SearchAllergy(mock.Anything, mock.Anything, mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, name string, pagination dto.Pagination) (*dto.TerminologyConnection, error) {
						return &dto.TerminologyConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"name":   allergyName,
					"limit":  strconv.Itoa(limit),
					"last":   strconv.Itoa(limit),
					"before": id,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: invalid limit value",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"limit": "invalid"}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last value",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"last": "badnum"}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				limit := 10
				allergyName := "aspirin"

				uc.EXPECT().SearchAllergy(mock.Anything, mock.Anything, mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, name string, pagination dto.Pagination) (*dto.TerminologyConnection, error) {
						return nil, fmt.Errorf("database error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"name":  allergyName,
					"limit": strconv.Itoa(limit),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.SearchAllergy(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListPatientCompositions(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully list patient compositions",
			setup: func(uc *restMocks.Mockusecases) args {
				limit := 5
				patientID := gofakeit.UUID()
				encounterID := gofakeit.UUID()
				dateStr := "2023-10-01"

				uc.EXPECT().ListPatientCompositions(mock.Anything, mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*scalarutils.Date"), mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.CompositionConnection, error) {
						return &dto.CompositionConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id":   patientID,
					"encounter_id": encounterID,
					"date":         dateStr,
					"limit":        strconv.Itoa(limit),
					"last":         strconv.Itoa(limit),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: invalid limit",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"limit": "notanint",
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid date",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"date": "invalid-date",
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last value",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"last": "nope",
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				limit := 3
				patientID := gofakeit.UUID()

				uc.EXPECT().ListPatientCompositions(mock.Anything, mock.Anything, mock.AnythingOfType("*string"), mock.AnythingOfType("*scalarutils.Date"), mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.CompositionConnection, error) {
						return nil, fmt.Errorf("unexpected error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id": patientID,
					"limit":      strconv.Itoa(limit),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.ListPatientCompositions(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_FetchMedicationRequestByID(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully fetch medication request",
			setup: func(uc *restMocks.Mockusecases) args {
				medicationRequestID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: medicationRequestID},
				}

				medicationReq := &dto.MedicationRequestOutput{
					ID: medicationRequestID,
				}

				uc.EXPECT().FetchMedicationRequestByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, medicationRequestID string) (*dto.MedicationRequestOutput, error) {
						return medicationReq, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Missing medication request ID in path",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Unable to fetch medication request by id",
			setup: func(uc *restMocks.Mockusecases) args {
				medicationRequestID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: medicationRequestID},
				}

				uc.EXPECT().FetchMedicationRequestByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, medicationRequestID string) (*dto.MedicationRequestOutput, error) {
						return nil, fmt.Errorf("medication request not found")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.FetchMedicationRequestByID(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListPatientConditions(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully list conditions",
			setup: func(uc *restMocks.Mockusecases) args {
				limit := 5
				patientID := gofakeit.UUID()
				encounterID := gofakeit.UUID()

				uc.EXPECT().ListPatientConditions(mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*scalarutils.Date"), mock.Anything, mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, strategy string, pagination dto.Pagination) (*dto.ConditionConnection, error) {
						return &dto.ConditionConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id":   patientID,
					"encounter_id": encounterID,
					"date":         "2023-11-01",
					"limit":        strconv.Itoa(limit),
					"last":         strconv.Itoa(limit),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: invalid date format",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id": "123",
					"date":       "invalid-date",
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid limit format",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id": "123",
					"limit":      "not-a-number",
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last format",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id": "123",
					"last":       "not-a-number",
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				limit := 10
				patientID := "456"

				uc.EXPECT().ListPatientConditions(mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("*scalarutils.Date"), mock.Anything, mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, strategy string, pagination dto.Pagination) (*dto.ConditionConnection, error) {
						return nil, fmt.Errorf("some db error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id": patientID,
					"limit":      strconv.Itoa(limit),
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.ListPatientConditions(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_SimpleQuestionnaireResponse(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully fetch questionnaire response",
			setup: func(uc *restMocks.Mockusecases) args {
				questionnaireResponseID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: questionnaireResponseID},
				}

				questionnaireResponse := []*domain.SimpleQuestionnaireResponse{
					{
						Group:     new(string),
						Questions: []domain.Questions{},
					},
				}

				uc.EXPECT().SimpleQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, questionnaireResponseID string) ([]*domain.SimpleQuestionnaireResponse, error) {
						return questionnaireResponse, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Missing questionnaire ID in path",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Unable to fetch questionnaire response by id",
			setup: func(uc *restMocks.Mockusecases) args {
				questionnaireResponseID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: questionnaireResponseID},
				}

				uc.EXPECT().SimpleQuestionnaireResponse(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, questionnaireResponseID string) ([]*domain.SimpleQuestionnaireResponse, error) {
						return nil, fmt.Errorf("questionnaire response not found")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.SimpleQuestionnaireResponse(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_GetEpisodeOfCare(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully get episode of care",
			setup: func(uc *restMocks.Mockusecases) args {
				episodeOfCareID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: episodeOfCareID},
				}

				episodeOfCare := &dto.EpisodeOfCare{
					ID: episodeOfCareID,
				}

				uc.EXPECT().GetEpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.EpisodeOfCare, error) {
						return episodeOfCare, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Missing episode of care ID in path",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Unable to retrieve episode of care by id",
			setup: func(uc *restMocks.Mockusecases) args {
				episodeOfCareID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: episodeOfCareID},
				}

				uc.EXPECT().GetEpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.EpisodeOfCare, error) {
						return nil, fmt.Errorf("episode of care not found")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.GetEpisodeOfCare(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_GetMedicalData(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully get medical data by patient id",
			setup: func(uc *restMocks.Mockusecases) args {
				patientID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: patientID},
				}

				medicalData := &dto.MedicalData{
					Regimen: []*dto.MedicationStatement{
						{
							ID: patientID,
						},
					},
				}

				uc.EXPECT().GetMedicalData(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientID string) (*dto.MedicalData, error) {
						return medicalData, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Missing patient ID in path",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Unable to retrieve medical data by patient id",
			setup: func(uc *restMocks.Mockusecases) args {
				patientID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: patientID},
				}

				uc.EXPECT().GetMedicalData(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientID string) (*dto.MedicalData, error) {
						return nil, fmt.Errorf("patient medical data not found")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.GetMedicalData(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_GetLabOrderByID(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully fetch lab order by id",
			setup: func(uc *restMocks.Mockusecases) args {
				labOrder := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: labOrder},
				}

				labOrderOutput := &dto.ServiceRequest{
					ID: labOrder,
				}

				uc.EXPECT().GetLabOrder(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, serviceRequestID string) (*dto.ServiceRequest, error) {
						return labOrderOutput, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Missing service request ID in path",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Unable to retrieve lab order by id",
			setup: func(uc *restMocks.Mockusecases) args {
				serviceRequestID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: serviceRequestID},
				}

				uc.EXPECT().GetLabOrder(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, serviceRequestID string) (*dto.ServiceRequest, error) {
						return nil, fmt.Errorf("lab order not found")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.GetLabOrderByID(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RecordConsent(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Bad input data",
			setup: func(_ *restMocks.Mockusecases) args {
				body := []byte(`{"decision": 10}`)

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Invalid Decision enum",
			setup: func(uc *restMocks.Mockusecases) args {
				consentInput := dto.ConsentInput{
					Decision:      mock.Anything,
					EncounterID:   uuid.NewString(),
					ScreeningType: dto.BreastCancerScreeningTypeEnum,
				}

				uc.EXPECT().RecordConsent(mock.Anything, mock.AnythingOfType("dto.ConsentInput")).
					RunAndReturn(func(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error) {
						return nil, errors.New("failed to record consent")
					})

				body, err := json.Marshal(consentInput)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Invalid screeningType enum",
			setup: func(uc *restMocks.Mockusecases) args {
				consentInput := dto.ConsentInput{
					Decision:      domain.ConsentDecisionPermit,
					EncounterID:   uuid.NewString(),
					ScreeningType: mock.Anything,
				}

				uc.EXPECT().RecordConsent(mock.Anything, mock.AnythingOfType("dto.ConsentInput")).
					RunAndReturn(func(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error) {
						return nil, errors.New("failed to record consent")
					})

				body, err := json.Marshal(consentInput)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				consentInput := dto.ConsentInput{
					Decision:      domain.ConsentDecisionPermit,
					EncounterID:   uuid.NewString(),
					ScreeningType: dto.BreastCancerScreeningTypeEnum,
				}

				uc.EXPECT().RecordConsent(mock.Anything, mock.AnythingOfType("dto.ConsentInput")).
					RunAndReturn(func(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error) {
						return nil, errors.New("failed to record consent")
					})

				body, err := json.Marshal(consentInput)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully record consent",
			setup: func(uc *restMocks.Mockusecases) args {
				consentInput := dto.ConsentInput{}

				uc.EXPECT().RecordConsent(mock.Anything, mock.AnythingOfType("dto.ConsentInput")).
					RunAndReturn(func(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error) {
						return &dto.ConsentOutput{}, nil
					})

				body, err := json.Marshal(consentInput)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.RecordConsent(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RecordObservation(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Bad input data",
			setup: func(uc *restMocks.Mockusecases) args {
				body := []byte(`{"status": 100}`)
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Missing required field",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.ObservationInput{
					Status: mock.Anything,
				}

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Missing concept field",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.ObservationInput{
					Status:      mock.Anything,
					EncounterID: mock.Anything,
				}

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Invalid status value",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.ObservationInput{
					Status:      mock.Anything,
					EncounterID: mock.Anything,
				}

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.ObservationInput{
					Status:      dto.ObservationStatusCancelled,
					EncounterID: uuid.NewString(),
					Concept:     dto.ObservationConceptEnumBMI,
					Value:       strconv.Itoa(gofakeit.Number(18, 24)),
				}

				uc.EXPECT().RecordObservationV2(mock.Anything, mock.AnythingOfType("dto.ObservationInput")).
					RunAndReturn(func(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
						return nil, errors.New("failed to record observation")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully record observation",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.ObservationInput{
					Status:      dto.ObservationStatusCancelled,
					Concept:     dto.ObservationConceptEnumBMI,
					EncounterID: uuid.NewString(),
					Value:       strconv.Itoa(gofakeit.Number(18, 24)),
				}

				uc.EXPECT().RecordObservationV2(mock.Anything, mock.AnythingOfType("dto.ObservationInput")).
					RunAndReturn(func(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
						return &dto.Observation{}, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.RecordObservation(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RecordTestResult(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Bad input data",
			setup: func(uc *restMocks.Mockusecases) args {
				body := []byte(`{"entry": 10, "servicerequestID": 12}`)

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.TestResultInput{
					Entry: dto.DiagnosticReportInput{
						EncounterID: uuid.NewString(),
						Note:        mock.Anything,
					},
					ServiceRequestID: uuid.NewString(),
				}

				uc.EXPECT().RecordTestResult(mock.Anything, mock.AnythingOfType("dto.TestResultInput")).
					RunAndReturn(func(ctx context.Context, input dto.TestResultInput) (*dto.DiagnosticReport, error) {
						return nil, errors.New("failed to record test result")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully record test results",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.TestResultInput{
					Entry: dto.DiagnosticReportInput{
						EncounterID: uuid.NewString(),
						Note:        mock.Anything,
					},
					ServiceRequestID: uuid.NewString(),
				}

				uc.EXPECT().RecordTestResult(mock.Anything, mock.AnythingOfType("dto.TestResultInput")).
					RunAndReturn(func(ctx context.Context, input dto.TestResultInput) (*dto.DiagnosticReport, error) {
						return &dto.DiagnosticReport{}, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.RecordTestResult(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateLabOrderResult(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Bad input data",
			setup: func(uc *restMocks.Mockusecases) args {
				body := []byte(`"servicerequestID": 21}`)

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.TestOrderResult{
					ServiceRequestID: uuid.NewString(),
					Test: []dto.TestOrderObservation{
						{
							Test:    mock.Anything,
							Value:   mock.Anything,
							Finding: mock.Anything,
						},
					},
				}

				uc.EXPECT().CreateLabOrderResult(mock.Anything, mock.AnythingOfType("*dto.TestOrderResult")).
					RunAndReturn(func(ctx context.Context, input *dto.TestOrderResult) (*dto.DiagnosticReport, error) {
						return nil, errors.New("failed to create lab order result")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully create lab order result",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.TestOrderResult{
					ServiceRequestID: uuid.NewString(),
					Test: []dto.TestOrderObservation{
						{
							Test:    mock.Anything,
							Value:   mock.Anything,
							Finding: mock.Anything,
						},
					},
				}

				uc.EXPECT().CreateLabOrderResult(mock.Anything, mock.AnythingOfType("*dto.TestOrderResult")).
					RunAndReturn(func(ctx context.Context, input *dto.TestOrderResult) (*dto.DiagnosticReport, error) {
						return &dto.DiagnosticReport{}, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreateLabOrderResult(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ScheduleAppointment(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Bad input data",
			setup: func(uc *restMocks.Mockusecases) args {
				body := []byte(`{"appointmentInput": 21}`)

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.ScheduleAppointmentPayload{
					AppointmentInput: dto.ScheduleAppointmentInput{
						EncounterID: uuid.NewString(),
						PatientID:   uuid.NewString(),
						Reason:      mock.Anything,
					},
					HeadersInput: dto.AdvantageHeaders{
						Organisation: mock.Anything,
						Cluster:      mock.Anything,
						Department:   mock.Anything,
						Branch:       mock.Anything,
						Workstation:  mock.Anything,
						Variant:      mock.Anything,
					},
				}

				uc.EXPECT().ScheduleAppointment(mock.Anything, mock.AnythingOfType("*dto.ScheduleAppointmentInput"), mock.AnythingOfType("*dto.AdvantageHeaders")).
					RunAndReturn(func(ctx context.Context, input *dto.ScheduleAppointmentInput, headers *dto.AdvantageHeaders) (bool, error) {
						return false, errors.New("failed to schedule appointment")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully schedule appointment",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.ScheduleAppointmentPayload{
					AppointmentInput: dto.ScheduleAppointmentInput{
						EncounterID: uuid.NewString(),
						PatientID:   uuid.NewString(),
						Reason:      mock.Anything,
					},
					HeadersInput: dto.AdvantageHeaders{
						Organisation: mock.Anything,
						Cluster:      mock.Anything,
						Department:   mock.Anything,
						Branch:       mock.Anything,
						Workstation:  mock.Anything,
						Variant:      mock.Anything,
					},
				}

				uc.EXPECT().ScheduleAppointment(mock.Anything, mock.AnythingOfType("*dto.ScheduleAppointmentInput"), mock.AnythingOfType("*dto.AdvantageHeaders")).
					RunAndReturn(func(ctx context.Context, input *dto.ScheduleAppointmentInput, headers *dto.AdvantageHeaders) (bool, error) {
						return true, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.ScheduleAppointment(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RecordTests(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Bad input data",
			setup: func(_ *restMocks.Mockusecases) args {
				body := []byte(`{"testType": 12}`)

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.TestInput{
					TestType: dto.MammogramTest,
					Input: dto.DiagnosticReportInput{
						EncounterID: uuid.NewString(),
						Note:        mock.Anything,
						Findings:    mock.Anything,
					},
				}

				uc.EXPECT().RecordTests(mock.Anything, mock.AnythingOfType("dto.TestInput")).
					RunAndReturn(func(ctx context.Context, payload dto.TestInput) (*dto.DiagnosticReport, error) {
						return nil, errors.New("failed to record tests")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully record tests",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.TestInput{
					TestType: dto.MammogramTest,
					Input: dto.DiagnosticReportInput{
						EncounterID: uuid.NewString(),
						Note:        mock.Anything,
						Findings:    mock.Anything,
					},
				}

				uc.EXPECT().RecordTests(mock.Anything, mock.AnythingOfType("dto.TestInput")).
					RunAndReturn(func(ctx context.Context, payload dto.TestInput) (*dto.DiagnosticReport, error) {
						return &dto.DiagnosticReport{}, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.RecordTests(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreatePrescription(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case:Bad input data",
			setup: func(uc *restMocks.Mockusecases) args {
				body := []byte(`{encounterID: 12345}`)

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.PrescriptionInput{
					EncounterID: uuid.NewString(),
					Medications: []dto.PrescriptionMedicationInput{
						{
							MedicationID: mock.Anything,
						},
					},
				}

				uc.EXPECT().CreatePrescription(mock.Anything, mock.AnythingOfType("dto.PrescriptionInput")).
					RunAndReturn(func(ctx context.Context, input dto.PrescriptionInput) ([]*dto.MedicationRequestOutput, error) {
						return nil, errors.New("failed to create prescription")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully create prescription",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.PrescriptionInput{
					EncounterID: uuid.NewString(),
					Medications: []dto.PrescriptionMedicationInput{
						{
							MedicationID: mock.Anything,
							DosageInstructions: []dto.DosageInstruction{
								{
									Route: dto.ValueSetData{
										Code:    mock.Anything,
										Display: mock.Anything,
									},
									DoseQuantity: 23.45,
									DoseUnit:     mock.Anything,
									Period:       mock.Anything,
									Frequency:    2,
									Duration:     mock.Anything,
									Condition:    mock.Anything,
								},
							},
							Priority: dto.MedicationRequestPriorityAsap,
						},
					},
				}

				uc.EXPECT().CreatePrescription(mock.Anything, mock.AnythingOfType("dto.PrescriptionInput")).
					RunAndReturn(func(ctx context.Context, input dto.PrescriptionInput) ([]*dto.MedicationRequestOutput, error) {
						return []*dto.MedicationRequestOutput{}, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreatePrescription(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_UpdateMedication(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Bad input data",
			setup: func(uc *restMocks.Mockusecases) args {
				body := []byte(`{"status": 213}`)

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				input := dto.PatchMedicationInput{
					Status: domain.ActiveMedicationStatus,
				}

				uc.EXPECT().PatchMedicationRequests(mock.Anything, mock.Anything, mock.AnythingOfType("domain.MedicationRequestStatus")).
					RunAndReturn(func(ctx context.Context, id string, value domain.MedicationRequestStatus) (*dto.MedicationRequestOutput, error) {
						return nil, errors.New("failed to update medication")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				pathParams := []gin.Param{
					{Key: "id", Value: encounterID},
				}

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy case: Successfully update medication",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				input := dto.PatchMedicationInput{
					Status: domain.ActiveMedicationStatus,
				}

				uc.EXPECT().PatchMedicationRequests(mock.Anything, mock.Anything, mock.AnythingOfType("domain.MedicationRequestStatus")).
					RunAndReturn(func(ctx context.Context, id string, value domain.MedicationRequestStatus) (*dto.MedicationRequestOutput, error) {
						return &dto.MedicationRequestOutput{}, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				pathParams := []gin.Param{
					{Key: "id", Value: encounterID},
				}

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.UpdateMedication(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateEpisodeOfCare(t *testing.T) {
	t.Parallel()
	input := dto.EpisodeOfCareInput{
		Status:    dto.EpisodeOfCareStatusEnumActive,
		PatientID: gofakeit.UUID(),
	}

	payload, _ := json.Marshal(input)
	fakeInput := dto.EpisodeOfCareInput{}
	fakePayload, _ := json.Marshal(fakeInput)

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull creates episode of care",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreateEpisodeOfCare(mock.Anything, mock.AnythingOfType("EpisodeOfCareInput")).
					RunAndReturn(func(ctx context.Context, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error) {
						return &dto.EpisodeOfCare{
							ID:        gofakeit.UUID(),
							Status:    "PLANNED",
							PatientID: gofakeit.UUID(),
						}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad Case: unable to create episode of care",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreateEpisodeOfCare(mock.Anything, mock.AnythingOfType("EpisodeOfCareInput")).
					RunAndReturn(func(ctx context.Context, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid input payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, fakePayload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreateEpisodeOfCare(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateCondition(t *testing.T) {
	t.Parallel()
	input := dto.ConditionInput{
		Code: gofakeit.Name(),
		Name: gofakeit.StreetName(),
	}

	payload, _ := json.Marshal(input)

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull creates condition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreateCondition(mock.Anything, mock.AnythingOfType("ConditionInput")).
					RunAndReturn(func(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error) {
						return &dto.Condition{
							ID:   gofakeit.UUID(),
							Name: gofakeit.Name(),
						}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad Case: failed to create condition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreateCondition(mock.Anything, mock.AnythingOfType("ConditionInput")).
					RunAndReturn(func(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreateCondition(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateAllergyIntolerance(t *testing.T) {
	t.Parallel()
	input := dto.AllergyInput{
		PatientID:   gofakeit.UUID(),
		EncounterID: gofakeit.UUID(),
		Code:        "test code",
	}

	payload, _ := json.Marshal(input)
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull creates allergy intollerance",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreateAllergyIntolerance(mock.Anything, mock.AnythingOfType("AllergyInput")).
					RunAndReturn(func(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error) {
						return &dto.Allergy{
							ID:          gofakeit.UUID(),
							PatientID:   gofakeit.UUID(),
							EncounterID: gofakeit.UUID(),
						}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Happy Case: successfull creates allergy intollerance",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreateAllergyIntolerance(mock.Anything, mock.AnythingOfType("AllergyInput")).
					RunAndReturn(func(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreateAllergyIntolerance(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreatePatient(t *testing.T) {
	t.Parallel()

	input := dto.PatientInput{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
	}
	fakeInput := dto.PatientInput{
		Gender: "male",
	}
	payload, _ := json.Marshal(input)
	fakePayload, _ := json.Marshal(fakeInput)

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull creates patient",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreatePatient(mock.Anything, mock.AnythingOfType("PatientInput")).
					RunAndReturn(func(ctx context.Context, input dto.PatientInput) (*dto.Patient, error) {
						return &dto.Patient{
							ID:     gofakeit.UUID(),
							Active: true,
							Gender: "Female",
						}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad Case: fails to create patient",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreatePatient(mock.Anything, mock.AnythingOfType("PatientInput")).
					RunAndReturn(func(ctx context.Context, input dto.PatientInput) (*dto.Patient, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: invalid input fields",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, fakePayload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreatePatient(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateComposition(t *testing.T) {
	t.Parallel()
	input := dto.CompositionInput{
		EncounterID: gofakeit.UUID(),
	}

	payload, _ := json.Marshal(input)

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull creates a composition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreateComposition(mock.Anything, mock.AnythingOfType("CompositionInput")).
					RunAndReturn(func(ctx context.Context, input dto.CompositionInput) (*dto.Composition, error) {
						return &dto.Composition{
							ID: gofakeit.UUID(),
						}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad Case: fails to create composition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreateComposition(mock.Anything, mock.AnythingOfType("CompositionInput")).
					RunAndReturn(func(ctx context.Context, input dto.CompositionInput) (*dto.Composition, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreateComposition(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_PatchEncounter(t *testing.T) {
	t.Parallel()
	testID := gofakeit.UUID()

	input := dto.EncounterInput{
		Status: "PLANNED",
	}

	fakeInput := dto.EncounterInput{}
	fakePayload, _ := json.Marshal(fakeInput)
	payload, _ := json.Marshal(input)

	pathParams := []gin.Param{
		{
			Key:   "id",
			Value: gofakeit.UUID(),
		},
	}

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull updates an encounter",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().PatchEncounter(mock.Anything, mock.Anything, mock.AnythingOfType("EncounterInput")).
					RunAndReturn(func(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error) {
						return &dto.Encounter{
							ID:        &testID,
							PatientID: &testID,
						}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: fails to update an encounter",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().PatchEncounter(mock.Anything, mock.Anything, mock.AnythingOfType("EncounterInput")).
					RunAndReturn(func(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid input field",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, fakePayload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.PatchEncounter(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_PatchTask(t *testing.T) {
	t.Parallel()
	input := dto.PatchTaskInput{
		Status: "completed",
	}

	payload, _ := json.Marshal(input)
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}

	pathParams := []gin.Param{
		{
			Key:   "id",
			Value: gofakeit.UUID(),
		},
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull creates a composition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().UpdateTask(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.PatchTaskInput")).
					RunAndReturn(func(ctx context.Context, taskID string, updateData *dto.PatchTaskInput) (bool, error) {
						return true, nil
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: fails to create composition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().UpdateTask(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.PatchTaskInput")).
					RunAndReturn(func(ctx context.Context, taskID string, updateData *dto.PatchTaskInput) (bool, error) {
						return false, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: invalid payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.PatchTask(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_PatchComposition(t *testing.T) {
	t.Parallel()

	input := dto.PatchCompositionInput{
		Note: "example notes",
	}
	payload, _ := json.Marshal(input)

	pathParams := []gin.Param{
		{
			Key:   "id",
			Value: gofakeit.UUID(),
		},
	}
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull creates a composition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().AppendNoteToComposition(mock.Anything, mock.Anything, mock.AnythingOfType("PatchCompositionInput")).
					RunAndReturn(func(ctx context.Context, id string, input dto.PatchCompositionInput) (*dto.Composition, error) {
						return &dto.Composition{}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad Case: fails to create composition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().AppendNoteToComposition(mock.Anything, mock.Anything, mock.AnythingOfType("PatchCompositionInput")).
					RunAndReturn(func(ctx context.Context, id string, input dto.PatchCompositionInput) (*dto.Composition, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.PatchComposition(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_StartEncounter(t *testing.T) {
	t.Parallel()
	input := dto.StartEncounterInput{
		EpisodeOfCareID: gofakeit.UUID(),
	}
	fakeInput := dto.StartEncounterInput{}

	payload, _ := json.Marshal(input)
	fakePayload, _ := json.Marshal(fakeInput)

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfully starts an encounter",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().StartEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeID string) (string, error) {
						return gofakeit.UUID(), nil
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad Case: failed to start an encounter",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().StartEncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, episodeID string) (string, error) {
						return "", fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: invalid payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, fakePayload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.StartEncounter(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_PatchPatient(t *testing.T) {
	t.Parallel()

	input := dto.PatientInput{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
	}

	fakeInput := dto.PatientInput{
		Gender: "male",
	}
	payload, _ := json.Marshal(input)
	fakePayload, _ := json.Marshal(fakeInput)

	pathParams := []gin.Param{
		{
			Key:   "id",
			Value: gofakeit.UUID(),
		},
	}

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfully updates patient's record",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().PatchPatient(mock.Anything, mock.Anything, mock.AnythingOfType("dto.PatientInput")).
					RunAndReturn(func(ctx context.Context, id string, input dto.PatientInput) (*dto.Patient, error) {
						return &dto.Patient{ID: gofakeit.UUID()}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: failed to update a patient's record",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().PatchPatient(mock.Anything, mock.Anything, mock.AnythingOfType("dto.PatientInput")).
					RunAndReturn(func(ctx context.Context, id string, input dto.PatientInput) (*dto.Patient, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, fakePayload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.PatchPatient(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_EndEpisodeOfCare(t *testing.T) {
	t.Parallel()
	pathParams := []gin.Param{
		{
			Key:   "id",
			Value: gofakeit.UUID(),
		},
	}

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfully ends an episode of care",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().EndEpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.EpisodeOfCare, error) {
						return &dto.EpisodeOfCare{ID: gofakeit.UUID()}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: failed to end an episode of care",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().EndEpisodeOfCare(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.EpisodeOfCare, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.EndEpisodeOfCare(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_PatchEpisodeOfCare(t *testing.T) {
	t.Parallel()

	input := dto.EpisodeOfCareInput{
		PatientID: gofakeit.UUID(),
		Status:    "PLANNED",
	}
	pathParams := []gin.Param{
		{
			Key:   "id",
			Value: gofakeit.UUID(),
		},
	}
	fakeInput := dto.EpisodeOfCareInput{
		PatientID: gofakeit.UUID(),
	}

	payload, _ := json.Marshal(input)
	fakePayload, _ := json.Marshal(fakeInput)

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfully updates an episode of care",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().PatchEpisodeOfCare(mock.Anything, mock.Anything, mock.AnythingOfType("EpisodeOfCareInput")).
					RunAndReturn(func(ctx context.Context, id string, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error) {
						return &dto.EpisodeOfCare{ID: gofakeit.UUID()}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: empty payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid payload",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, fakePayload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: failed to update an episode of care",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().PatchEpisodeOfCare(mock.Anything, mock.Anything, mock.AnythingOfType("EpisodeOfCareInput")).
					RunAndReturn(func(ctx context.Context, id string, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.PatchEpisodeOfCare(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListLabOrders(t *testing.T) {
	t.Parallel()

	today := time.Now().Format(time.DateOnly)

	patientID := gofakeit.UUID()
	encounterID := gofakeit.UUID()
	facilityID := gofakeit.UUID()
	last := 10

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfully lists lab orders",
			setup: func(r *restMocks.Mockusecases) args {
				r.On(
					"ListLabOrders",
					mock.Anything,
					mock.AnythingOfType("*dto.ServiceRequestFilterInput"),
					mock.AnythingOfType("Pagination"),
				).Return(
					&dto.ServiceRequestOutputConnection{
						TotalCount: &last,
					}, nil,
				)
				r.EXPECT().ListLabOrders(mock.Anything, mock.AnythingOfType("*dto.ServiceRequestFilterInput"), mock.AnythingOfType("Pagination")).
					RunAndReturn(func(ctx context.Context, filter *dto.ServiceRequestFilterInput, pagination dto.Pagination) (*dto.ServiceRequestOutputConnection, error) {
						return &dto.ServiceRequestOutputConnection{TotalCount: &last}, nil
					})
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"date":        today,
						"patientID":   patientID,
						"encounterID": encounterID,
						"first":       strconv.Itoa(10),
						"last":        strconv.Itoa(last),
						"before":      "before",
						"after":       "testAfter",
						"status":      "draft",
						"type":        "LP94892-4",
						"facilityID":  facilityID,
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: failed to list lab orders",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().ListLabOrders(mock.Anything, mock.AnythingOfType("*dto.ServiceRequestFilterInput"), mock.AnythingOfType("Pagination")).
					RunAndReturn(func(ctx context.Context, filter *dto.ServiceRequestFilterInput, pagination dto.Pagination) (*dto.ServiceRequestOutputConnection, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"date":        today,
						"patientID":   patientID,
						"encounterID": encounterID,
						"first":       strconv.Itoa(10),
						"last":        strconv.Itoa(last),
						"before":      "before",
						"after":       "testAfter",
						"status":      "draft",
						"type":        "LP94892-4",
						"facilityID":  facilityID,
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: invalid status",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"status": "testing status",
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid limit",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"last": "testing last",
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid first",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"first": "testing first",
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid date",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"date": "testing time",
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.ListLabOrders(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_PatientReferralDetailsByID(t *testing.T) {
	t.Parallel()
	pathParams := []gin.Param{
		{
			Key:   "serviceRequestID",
			Value: gofakeit.UUID(),
		},
	}
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfully gets a patient's record by ID",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().GetPatientReferralDetails(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, serviceRequestID string) (*domain.PatientReferralDetails, error) {
						return &domain.PatientReferralDetails{PatientName: gofakeit.FirstName()}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: failed to get a patient's record by ID",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().GetPatientReferralDetails(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, serviceRequestID string) (*domain.PatientReferralDetails, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.PatientReferralDetailsByID(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ReferralDetails(t *testing.T) {
	t.Parallel()

	today := time.Now().Format(time.DateOnly)

	patientID := gofakeit.UUID()
	encounterID := gofakeit.UUID()
	last := 10

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfully list referral details",
			setup: func(r *restMocks.Mockusecases) args {
				r.On(
					"GetPatientReferrals",
					mock.Anything,
					mock.AnythingOfType("*dto.ReferralSearchInput"),
				).Return(&dto.ReferralDetailConnection{TotalCount: 10}, nil)

				ctx, rec := newGinTestCtx(
					http.MethodGet,
					map[string]string{
						"date":        today,
						"patientID":   patientID,
						"encounterID": encounterID,
						"first":       strconv.Itoa(10),
						"last":        strconv.Itoa(last),
						"before":      "before",
						"after":       "testAfter",
						"status":      "active",
					},
					nil,
					nil,
					nil,
				)

				return args{ctx: ctx, rec: rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: failed to list lab orders",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().GetPatientReferrals(mock.Anything, mock.AnythingOfType("*dto.ReferralSearchInput")).
					RunAndReturn(func(ctx context.Context, searchInput *dto.ReferralSearchInput) (*dto.ReferralDetailConnection, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"patientID":   patientID,
						"encounterID": encounterID,
						"first":       strconv.Itoa(10),
						"last":        strconv.Itoa(last),
						"before":      "before",
						"after":       "testAfter",
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad Case: invalid limit",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"last": "testing last",
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid first",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"first": "testing first",
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid date",
			setup: func(r *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(
					http.MethodPatch,
					map[string]string{
						"date": "testing time",
					},
					nil,
					nil,
					nil,
				)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.ReferralDetails(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_UploadMedia(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: fail to upload media",
			setup: func(uc *restMocks.Mockusecases) args {
				ctx, rec, _ := multipartCtx("ENC1", true)

				uc.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID, serviceRequestID string, file io.Reader, contentType string) (*dto.Media, error) {
						return nil, errors.New("gcs down")
					})

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Happy Case: Successfully upload file",
			setup: func(uc *restMocks.Mockusecases) args {
				ctx, rec, _ := multipartCtx("ENC1", true)

				uc.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID, serviceRequestID string, file io.Reader, contentType string) (*dto.Media, error) {
						return &dto.Media{MediaLink: "https://example.com/file.txt"}, nil
					})

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.UploadMedia(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateReferral(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Sad Case: Invalid json",
			setup: func(uc *restMocks.Mockusecases) args {
				c, r := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{bad}`), nil)

				return args{c, r}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Fail to create a referral",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.CreateReferralInput{
					EncounterID:     uuid.New().String(),
					Diagnosis:       gofakeit.Name(),
					ReferralType:    "INBOUND",
					Urgency:         "URGENT",
					ClinicalHistory: "",
					ReferralDate: scalarutils.Date{
						Year:  2025,
						Month: 4,
						Day:   4,
					},
				}
				body, _ := json.Marshal(input)

				uc.On("CreateReferral", mock.Anything, &input).Return(
					nil, errors.New("failed to create referral"),
				)

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)
				return args{c, r}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Happy Case: Successfully create a referral",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.CreateReferralInput{
					EncounterID:     uuid.New().String(),
					Diagnosis:       gofakeit.Name(),
					ReferralType:    "INBOUND",
					Urgency:         "URGENT",
					ClinicalHistory: "",
					ReferralDate: scalarutils.Date{
						Year:  2025,
						Month: 4,
						Day:   4,
					},
				}
				body, _ := json.Marshal(input)

				output := &dto.ServiceRequest{
					ID:     uuid.New().String(),
					Status: "",
					Intent: "",
				}

				uc.On("CreateReferral", mock.Anything, &input).Return(
					output, nil,
				)

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)
				return args{c, r}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.CreateReferral(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ReceivePubSubPushMessage(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, adv *advantageMocks.MockAdvantageService) args
		want  int
	}{
		{
			name: "Sad Case: JWT verification fails",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return nil, errors.New("jwt err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Cannot get pubsub topic",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return "", errors.New("topic err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: unknown topic",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return "nonexistent-topic", nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: CreatePubsubPatient returns error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				patientMsg := dto.PatientPubSubMessage{UserID: uuid.NewString()}
				patientBytes, _ := json.Marshal(patientMsg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: patientBytes},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubPatient(mock.Anything, mock.AnythingOfType("dto.PatientPubSubMessage")).
					RunAndReturn(func(ctx context.Context, payload dto.PatientPubSubMessage) error {
						return errors.New("fail")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Patient message processed",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				patientMsg := dto.PatientPubSubMessage{
					UserID:         gofakeit.UUID(),
					ClientID:       gofakeit.UUID(),
					Name:           gofakeit.Name(),
					Gender:         "male",
					Active:         true,
					PhoneNumber:    gofakeit.Phone(),
					OrganizationID: gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
				}
				patientBytes, _ := json.Marshal(patientMsg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: patientBytes},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubPatient(mock.Anything, mock.AnythingOfType("dto.PatientPubSubMessage")).
					RunAndReturn(func(ctx context.Context, payload dto.PatientPubSubMessage) error {
						return nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Patient topic but invalid JSON payload",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				bad := []byte(`{bad`)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: bad},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Organization topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: CreatePubsubOrganization error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.FacilityPubSubMessage{Name: "Bad Facility"}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubOrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.FacilityPubSubMessage) error {
						return errors.New("org err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Organization topic",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.FacilityPubSubMessage{Name: "Facility A"}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubOrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.FacilityPubSubMessage) error {
						return nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Vitals topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: CreatePubsubVitals error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.VitalSignPubSubMessage{PatientID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubVitals(mock.Anything, mock.AnythingOfType("dto.VitalSignPubSubMessage")).
					RunAndReturn(func(ctx context.Context, data dto.VitalSignPubSubMessage) error {
						return errors.New("vitals err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Vitals topic",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.VitalSignPubSubMessage{PatientID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubVitals(mock.Anything, mock.AnythingOfType("dto.VitalSignPubSubMessage")).
					RunAndReturn(func(ctx context.Context, data dto.VitalSignPubSubMessage) error {
						return nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Allergy topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: CreatePubsubAllergyIntolerance error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.PatientAllergyPubSubMessage{PatientID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubAllergyIntolerance(mock.Anything, mock.AnythingOfType("dto.PatientAllergyPubSubMessage")).
					RunAndReturn(func(ctx context.Context, data dto.PatientAllergyPubSubMessage) error {
						return errors.New("allergy err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Allergy topic",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.PatientAllergyPubSubMessage{PatientID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubAllergyIntolerance(mock.Anything, mock.AnythingOfType("dto.PatientAllergyPubSubMessage")).
					RunAndReturn(func(ctx context.Context, data dto.PatientAllergyPubSubMessage) error {
						return nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Medication topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: CreatePubsubMedicationStatement error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.MedicationPubSubMessage{PatientID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubMedicationStatement(mock.Anything, mock.AnythingOfType("dto.MedicationPubSubMessage")).
					RunAndReturn(func(ctx context.Context, data dto.MedicationPubSubMessage) error {
						return errors.New("med err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Medication topic",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.MedicationPubSubMessage{PatientID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubMedicationStatement(mock.Anything, mock.AnythingOfType("dto.MedicationPubSubMessage")).
					RunAndReturn(func(ctx context.Context, data dto.MedicationPubSubMessage) error {
						return nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Tenant topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.TenantTopicName, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: CreatePubsubTenant error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.OrganizationInput{Name: "TenantX"}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.TenantTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubTenant(mock.Anything, mock.AnythingOfType("dto.OrganizationInput")).
					RunAndReturn(func(ctx context.Context, data dto.OrganizationInput) error {
						return errors.New("tenant err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Tenant topic",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.OrganizationInput{Name: "TenantA"}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.TenantTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubTenant(mock.Anything, mock.AnythingOfType("dto.OrganizationInput")).
					RunAndReturn(func(ctx context.Context, data dto.OrganizationInput) error {
						return nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: TestResult topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: CreatePubsubTestResult error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.PatientTestResultPubSubMessage{PatientID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubTestResult(mock.Anything, mock.AnythingOfType("dto.PatientTestResultPubSubMessage")).
					RunAndReturn(func(ctx context.Context, data dto.PatientTestResultPubSubMessage) error {
						return errors.New("result err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: TestResult topic",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.PatientTestResultPubSubMessage{PatientID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreatePubsubTestResult(mock.Anything, mock.AnythingOfType("dto.PatientTestResultPubSubMessage")).
					RunAndReturn(func(ctx context.Context, data dto.PatientTestResultPubSubMessage) error {
						return nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Segmentation topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.SegmentationTopicName, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Segmentation error",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, adv *advantageMocks.MockAdvantageService) args {
				msg := dto.SegmentationPayload{
					ClinicalID:   gofakeit.UUID(),
					SegmentLabel: dto.SegmentationBreastCategoryAverageRisk,
				}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: b},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.SegmentationTopicName, common.ClinicalServiceName), nil
					})
				adv.EXPECT().SegmentPatient(mock.Anything, mock.AnythingOfType("dto.SegmentationPayload")).
					RunAndReturn(func(ctx context.Context, payload dto.SegmentationPayload) error {
						return errors.New("failed to segment patient")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: ReferralTask topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.ReferralTopicName, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: CreateReferralTask error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.PatientReferralTaskPayload{
					Referral: &dto.ServiceRequest{ID: uuid.NewString()},
					Meta:     &dto.MetaInput{},
				}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.ReferralTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreateReferralTask(mock.Anything, mock.AnythingOfType("*dto.MetaInput"), mock.AnythingOfType("*dto.ServiceRequest")).
					RunAndReturn(func(ctx context.Context, tags *dto.MetaInput, serviceRequest *dto.ServiceRequest) (*domain.FHIRTask, error) {
						return nil, errors.New("task err")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: ReferralTask topic",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.PatientReferralTaskPayload{
					Referral: &dto.ServiceRequest{ID: uuid.NewString()},
					Meta:     &dto.MetaInput{},
				}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.ReferralTopicName, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreateReferralTask(mock.Anything, mock.AnythingOfType("*dto.MetaInput"), mock.AnythingOfType("*dto.ServiceRequest")).
					RunAndReturn(func(ctx context.Context, tags *dto.MetaInput, serviceRequest *dto.ServiceRequest) (*domain.FHIRTask, error) {
						return &domain.FHIRTask{}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: ReferralReportNotification topic invalid JSON",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.ReferralReportNotificationTopic, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: ReferralReportNotification topic",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.ReferralReportNotification{WorkstationID: uuid.NewString()}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.ReferralReportNotificationTopic, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: CreateCarePlan error",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.CarePlanInput{
					EncounterID:      uuid.NewString(),
					PlanDefinitionID: uuid.NewString(),
				}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.CreateCarePlanTopic, common.ClinicalServiceName), nil
					})
				uc.EXPECT().PatientCarePlan(mock.Anything, mock.AnythingOfType("*dto.CarePlanPayload")).
					RunAndReturn(func(ctx context.Context, input *dto.CarePlanPayload) (*domain.FHIRCarePlan, error) {
						return nil, errors.New("careplan error")
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: Care plan topic",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.CarePlanInput{
					EncounterID:      uuid.NewString(),
					PlanDefinitionID: uuid.NewString(),
				}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.CreateCarePlanTopic, common.ClinicalServiceName), nil
					})
				uc.EXPECT().PatientCarePlan(mock.Anything, mock.AnythingOfType("*dto.CarePlanPayload")).
					RunAndReturn(func(ctx context.Context, input *dto.CarePlanPayload) (*domain.FHIRCarePlan, error) {
						return &domain.FHIRCarePlan{}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: bad careplan data",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.CreateCarePlanTopic, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: unable to create task",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				active := "active"
				msg := domain.FHIRTaskInput{
					Status: (*scalarutils.Code)(&active),
				}

				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.FollowUpTaskTopic, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreateTask(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, task *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
					return nil, fmt.Errorf("error")
				})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: create task",
			setup: func(ext *infraMocks.MockBaseExtension, uc *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				msg := dto.CarePlanInput{
					EncounterID:      uuid.NewString(),
					PlanDefinitionID: uuid.NewString(),
				}
				b, _ := json.Marshal(msg)

				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{Message: pubsubtools.PubSubMessage{Data: b}}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.FollowUpTaskTopic, common.ClinicalServiceName), nil
					})
				uc.EXPECT().CreateTask(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, task *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
					return &domain.FHIRTask{}, nil
				})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: invalid task topic data",
			setup: func(ext *infraMocks.MockBaseExtension, _ *restMocks.Mockusecases, _ *advantageMocks.MockAdvantageService) args {
				ext.EXPECT().VerifyPubSubJWTAndDecodePayload(mock.Anything, mock.Anything).
					RunAndReturn(func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
						return &pubsubtools.PubSubPayload{
							Message: pubsubtools.PubSubMessage{Data: []byte(`{bad`)},
						}, nil
					})
				ext.EXPECT().GetPubSubTopic(mock.Anything).
					RunAndReturn(func(m *pubsubtools.PubSubPayload) (string, error) {
						return utils.AddPubSubNamespace(common.FollowUpTaskTopic, common.ClinicalServiceName), nil
					})

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockExt, mockUC, mockAdv)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.ReceivePubSubPushMessage(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListPatientAllergies(t *testing.T) {

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: invalid limit parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"limit": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"last": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Usecase Returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ListPatientAllergies(mock.Anything, mock.Anything, mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.AllergyConnection, error) {
						return nil, errors.New("failed to get patient allergies")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Happy Case: Successfully get a search patient allergies",
			setup: func(uc *restMocks.Mockusecases) args {
				patientID := uuid.NewString()
				after := "aBC123"
				first := 10

				uc.EXPECT().ListPatientAllergies(mock.Anything, mock.Anything, mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.AllergyConnection, error) {
						return &dto.AllergyConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id": patientID,
					"limit":      strconv.Itoa(first),
					"after":      after,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.ListPatientAllergies(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_ListPatientMedia(t *testing.T) {

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad Case: invalid limit parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"limit": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: invalid last parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{"last": "abc"}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Usecase Returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ListPatientMedia(mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, encounterID, serviceRequestID string, pagination dto.Pagination) (*dto.MediaConnection, error) {
						return nil, errors.New("failed to get patient media")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Happy Case: Successfully get patient media",
			setup: func(uc *restMocks.Mockusecases) args {
				patientID := uuid.NewString()
				after := "aBC123"
				first := 10

				uc.EXPECT().ListPatientMedia(mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("dto.Pagination")).
					RunAndReturn(func(ctx context.Context, encounterID, serviceRequestID string, pagination dto.Pagination) (*dto.MediaConnection, error) {
						return &dto.MediaConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"patient_id": patientID,
					"limit":      strconv.Itoa(first),
					"after":      after,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)

			hdl.ListPatientMedia(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_PatchPatientObservations(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Missing id parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				pathParams := []gin.Param{
					{Key: "id", Value: ""},
				}
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Bad input data",
			setup: func(uc *restMocks.Mockusecases) args {
				observationID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: observationID},
				}

				body := []byte(`{"value": 123}`)

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				observationID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: observationID},
				}

				input := dto.PatchObservationInput{
					Value:           mock.Anything,
					ObservationType: dto.PatientBloodSugar,
				}

				uc.EXPECT().PatchPatientObservation(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.PatchObservationInput")).
					RunAndReturn(func(ctx context.Context, observationID string, input *dto.PatchObservationInput) (*dto.Observation, error) {
						return nil, errors.New("an error occurred")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Happy case: Successfully patch patient observation",
			setup: func(uc *restMocks.Mockusecases) args {
				observationID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: observationID},
				}

				input := dto.PatchObservationInput{
					Value:           mock.Anything,
					ObservationType: dto.PatientBloodSugar,
				}

				uc.EXPECT().PatchPatientObservation(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.PatchObservationInput")).
					RunAndReturn(func(ctx context.Context, observationID string, input *dto.PatchObservationInput) (*dto.Observation, error) {
						return &dto.Observation{}, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPatch, nil, pathParams, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.PatchPatientObservations(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateQuestionnaireResponses(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Sad case: Missing questionnaire id parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPost, map[string]string{"questionnaireID": ""}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Missing encounter id parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPost, map[string]string{"questionnaireID": mock.Anything}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Bad input data",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()
				questionnaireID := uuid.NewString()

				body := []byte(`{"resourceType": 123}`)

				ctx, rec := newGinTestCtx(http.MethodPost, map[string]string{
					"encounterID":     encounterID,
					"questionnaireID": questionnaireID,
				}, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: Usecase returns error",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()
				questionnaireID := uuid.NewString()

				input := dto.QuestionnaireResponse{}

				uc.EXPECT().CreateQuestionnaireResponse(mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("dto.QuestionnaireResponse")).
					RunAndReturn(func(ctx context.Context, questionnaireID, encounterID string, input dto.QuestionnaireResponse) (*dto.QuestionnaireReviewSummary, error) {
						return nil, errors.New("an error occurred")
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, map[string]string{
					"encounterID":     encounterID,
					"questionnaireID": questionnaireID,
				}, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Happy case: Successfully create questionnaire response",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()
				questionnaireID := uuid.NewString()

				dummyText := mock.Anything
				input := dto.QuestionnaireResponse{
					ResourceType: mock.Anything,
					Meta: dto.MetaInput{
						VersionID: uuid.NewString(),
						Source:    mock.Anything,
					},
					Status:   domain.DiagnosticReportStatusAmended,
					Authored: mock.Anything,
					Item: []dto.QuestionnaireResponseItem{
						{
							LinkID: mock.Anything,
							Text:   &dummyText,
						},
					},
				}

				uc.EXPECT().CreateQuestionnaireResponse(mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("dto.QuestionnaireResponse")).
					RunAndReturn(func(ctx context.Context, questionnaireID, encounterID string, input dto.QuestionnaireResponse) (*dto.QuestionnaireReviewSummary, error) {
						return &dto.QuestionnaireReviewSummary{}, nil
					})

				body, err := json.Marshal(input)
				if err != nil {
					return args{nil, nil}
				}

				ctx, rec := newGinTestCtx(http.MethodPost, map[string]string{
					"encounterID":     encounterID,
					"questionnaireID": questionnaireID,
				}, nil, body, nil)

				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreateQuestionnaireResponses(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreatePlanDefinition(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: Successfully create a plan definition",
			setup: func(uc *restMocks.Mockusecases) args {
				input := dto.PlanDefinitionInput{
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

				body, _ := json.Marshal(input)

				uid := gofakeit.UUID()
				output := &domain.FHIRPlanDefinition{
					ID: &uid,
				}

				uc.EXPECT().CreatePlanDefinition(mock.Anything, mock.AnythingOfType("*dto.PlanDefinitionInput")).
					RunAndReturn(func(ctx context.Context, questionnaireInput *dto.PlanDefinitionInput) (*domain.FHIRPlanDefinition, error) {
						return output, nil
					})

				c, r := newGinTestCtx(http.MethodPost, nil, nil, body, nil)
				return args{c, r}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad Case: fails to create plan definition",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().CreatePlanDefinition(mock.Anything, mock.AnythingOfType("*dto.PlanDefinitionInput")).
					RunAndReturn(func(ctx context.Context, questionnaireInput *dto.PlanDefinitionInput) (*domain.FHIRPlanDefinition, error) {
						return nil, errors.New("failed to create plan definition")
					})

				input := dto.PlanDefinitionInput{}

				body, _ := json.Marshal(input)

				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, body, nil)

				return args{ctx: ctx, rec: rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.CreatePlanDefinition(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RetrievePlanDefinition(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: Successfully search plan definition",
			setup: func(uc *restMocks.Mockusecases) args {
				name := uuid.NewString()

				uc.EXPECT().RetrievePlanDefinition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, name string) (*dto.PlanDefinitionOutputConnection, error) {
						return &dto.PlanDefinitionOutputConnection{TotalCount: 1}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"name": name,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Unable to search plan definition",
			setup: func(uc *restMocks.Mockusecases) args {
				name := uuid.NewString()

				uc.EXPECT().RetrievePlanDefinition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, name string) (*dto.PlanDefinitionOutputConnection, error) {
						return nil, fmt.Errorf("some db error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"name": name,
				}, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)
			hdl := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			hdl.ListPlanDefinition(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_OncologyDiagnosis(t *testing.T) {
	input := &dto.OncologyDiagnosisInput{
		EncounterID: gofakeit.UUID(),
		Condition: dto.ValueSetData{
			Code:    gofakeit.UUID(),
			Display: gofakeit.UUID(),
		},
		ICDO3PrimaryTumorCode: gofakeit.UUID(),
		ICDO3MorphologyCode:   gofakeit.UUID(),
		Behavior: dto.ValueSetData{
			Code:    gofakeit.UUID(),
			Display: gofakeit.UUID(),
		},
		Grade: dto.ValueSetData{
			Code:    gofakeit.UUID(),
			Display: gofakeit.UUID(),
		},
		Stage: dto.ValueSetData{
			Code:    gofakeit.UUID(),
			Display: gofakeit.UUID(),
		},
		Notes: gofakeit.UUID(),
	}

	payload, err := json.Marshal(input)
	if err != nil {
		t.Errorf("unable to marshal input")
	}

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull create oncology condition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().RecordOncologicalDiagnosis(mock.Anything, mock.AnythingOfType("*dto.OncologyDiagnosisInput")).
					RunAndReturn(func(ctx context.Context, input *dto.OncologyDiagnosisInput) (*dto.Condition, error) {
						return &dto.Condition{
							ID:   gofakeit.UUID(),
							Name: gofakeit.Name(),
						}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad case: unable to create oncology condition",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().RecordOncologicalDiagnosis(mock.Anything, mock.AnythingOfType("*dto.OncologyDiagnosisInput")).
					RunAndReturn(func(ctx context.Context, input *dto.OncologyDiagnosisInput) (*dto.Condition, error) {
						return nil, fmt.Errorf("error")
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.OncologyDiagnosis(arg.ctx)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreateMedication(t *testing.T) {
	input := []*dto.MedicationInput{
		{
			Name: gofakeit.BeerName(),
			DoseForm: dto.ValueSetData{
				Code:    gofakeit.UUID(),
				Display: gofakeit.UUID(),
			},
		},
	}

	payload, err := json.Marshal(input)
	if err != nil {
		t.Errorf("unable to marshal input")
	}

	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		args  args
		want  int
	}{
		{
			name: "Happy Case: successfull create medication",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().RecordMedication(mock.Anything, mock.AnythingOfType("[]*dto.MedicationInput")).
					RunAndReturn(func(ctx context.Context, medications []*dto.MedicationInput) ([]*dto.MedicationOutput, error) {
						return []*dto.MedicationOutput{
							{
								ID: gofakeit.UUID(),
							},
						}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad case: unable to create medication",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().RecordMedication(mock.Anything, mock.AnythingOfType("[]*dto.MedicationInput")).
					RunAndReturn(func(ctx context.Context, medications []*dto.MedicationInput) ([]*dto.MedicationOutput, error) {
						return nil, fmt.Errorf("error")
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, payload, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.CreateMedication(arg.ctx)
		})
	}
}

func TestPresentationHandlersImpl_FetchMedicationByID(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}

	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: Successfully fetch medication",
			setup: func(uc *restMocks.Mockusecases) args {
				medicationID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: medicationID},
				}

				medicationReq := &dto.MedicationOutput{
					ID:   medicationID,
					Name: mock.Anything,
					DoseForm: dto.ValueSetData{
						Code:    mock.Anything,
						Display: mock.Anything,
					},
					Status:    mock.Anything,
					LotNumber: mock.Anything,
				}

				uc.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, medicationID string) (*dto.MedicationOutput, error) {
						return medicationReq, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad Case: Missing medication ID in path",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: Unable to fetch medication by id",
			setup: func(uc *restMocks.Mockusecases) args {
				medicationID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "id", Value: medicationID},
				}

				uc.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, medicationID string) (*dto.MedicationOutput, error) {
						return nil, fmt.Errorf("medication not found")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			ph := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			ph.FetchMedicationByID(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_CreatePatientCarePlan(t *testing.T) {
	input1 := &dto.CarePlanInput{
		EncounterID:      gofakeit.UUID(),
		PlanDefinitionID: gofakeit.UUID(),
		Notes:            mock.Anything,
	}

	payload1, err := json.Marshal(input1)
	if err != nil {
		t.Errorf("unable to marshal input")
	}

	input2 := &dto.CarePlanInput{
		EncounterID: gofakeit.UUID(),
		Notes:       mock.Anything,
	}

	payload2, err := json.Marshal(input2)
	if err != nil {
		t.Errorf("unable to marshal input")
	}

	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy Case: successfull create patient treatment plan",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreatePatientCarePlan(mock.Anything, mock.AnythingOfType("*dto.CarePlanInput")).
					RunAndReturn(func(ctx context.Context, input *dto.CarePlanInput) error {
						return nil
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, payload1, nil)
				return args{ctx, rec}
			},
			want: http.StatusCreated,
		},
		{
			name: "Sad case: unable to create create patient treatment plan",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreatePatientCarePlan(mock.Anything, mock.AnythingOfType("*dto.CarePlanInput")).
					RunAndReturn(func(ctx context.Context, input *dto.CarePlanInput) error {
						return fmt.Errorf("error")
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, payload1, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad case: invalid payload",
			setup: func(r *restMocks.Mockusecases) args {
				r.EXPECT().CreatePatientCarePlan(mock.Anything, mock.AnythingOfType("*dto.CarePlanInput")).
					RunAndReturn(func(ctx context.Context, input *dto.CarePlanInput) error {
						return fmt.Errorf("error")
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, payload2, nil)
				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.CreatePatientCarePlan(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_FetchCarePlan(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy case: Fetch care plan",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "encounterID", Value: encounterID},
				}

				uc.EXPECT().FetchPatientCarePlan(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (*dto.CarePlanOutput, error) {
						return &dto.CarePlanOutput{
							ID: gofakeit.UUID(),
						}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: unable to fetch care plan",
			setup: func(uc *restMocks.Mockusecases) args {
				encounterID := uuid.NewString()

				pathParams := []gin.Param{
					{Key: "encounterID", Value: encounterID},
				}

				uc.EXPECT().FetchPatientCarePlan(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, encounterID string) (*dto.CarePlanOutput, error) {
						return nil, fmt.Errorf("error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.FetchCarePlan(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_Observation(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy case: get observation details",
			setup: func(uc *restMocks.Mockusecases) args {
				pathParams := []gin.Param{
					{Key: "id", Value: gofakeit.UUID()},
				}

				uc.EXPECT().GetObservationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.Observation, error) {
						return &dto.Observation{
							ID: gofakeit.UUID(),
						}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: unable to get observation details",
			setup: func(uc *restMocks.Mockusecases) args {
				pathParams := []gin.Param{
					{Key: "id", Value: gofakeit.UUID()},
				}

				uc.EXPECT().GetObservationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.Observation, error) {
						return nil, fmt.Errorf("error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad case: missing id parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.Observation(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_RiskAssessment(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy case: get risk assessmemt details",
			setup: func(uc *restMocks.Mockusecases) args {
				pathParams := []gin.Param{
					{Key: "id", Value: gofakeit.UUID()},
				}

				uc.EXPECT().GetRiskAssessmentByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.RiskAssessment, error) {
						uid := gofakeit.UUID()
						return &dto.RiskAssessment{
							ID: &uid,
						}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: unable to get risk assessment details",
			setup: func(uc *restMocks.Mockusecases) args {
				pathParams := []gin.Param{
					{Key: "id", Value: gofakeit.UUID()},
				}

				uc.EXPECT().GetRiskAssessmentByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*dto.RiskAssessment, error) {
						return nil, fmt.Errorf("error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Sad case: missing id parameter",
			setup: func(_ *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPatch, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			arg := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.RiskAssessment(arg.c)

			require.Equal(t, tt.want, arg.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_PatientBanner(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy case: get patient bannerr",
			setup: func(uc *restMocks.Mockusecases) args {
				pathParams := []gin.Param{
					{Key: "id", Value: gofakeit.UUID()},
				}

				uc.EXPECT().GetPatientBanner(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.PatientEverythingFilterParams")).
					RunAndReturn(func(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.PatientBanner, error) {
						return &dto.PatientBanner{
							Conditions: []dto.TimelineResource{
								{
									ResourceType: dto.ResourceTypeCondition,
								},
							},
						}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: unable to get patient bannerr",
			setup: func(uc *restMocks.Mockusecases) args {
				pathParams := []gin.Param{
					{Key: "id", Value: gofakeit.UUID()},
				}

				uc.EXPECT().GetPatientBanner(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.PatientEverythingFilterParams")).
					RunAndReturn(func(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.PatientBanner, error) {
						return nil, fmt.Errorf("error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, nil, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: no patient id",
			setup: func(uc *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			args := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.PatientBanner(args.c)
		})
	}
}

func TestPresentationHandlersImpl_GetPatientTimeline(t *testing.T) {
	queryParams := map[string]string{
		"count":      "10",
		"page_token": "10",
		"start":      "10",
		"end":        "10",
		"since":      "10",
		"type":       "Observation,Encounter",
	}
	pathParams := []gin.Param{
		{Key: "id", Value: gofakeit.UUID()},
	}
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy case: get patient timelime",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().GetPatientTimeline(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.PatientEverythingFilterParams")).
					RunAndReturn(func(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.HealthTimeline, error) {
						return &dto.HealthTimeline{
							TotalCount: 2,
						}, nil
					})

				ctx, rec := newGinTestCtx(http.MethodGet, queryParams, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: missing type query param",
			setup: func(uc *restMocks.Mockusecases) args {
				invalidQueryParams := map[string]string{
					"count":      "10",
					"page_token": "10",
					"start":      "10",
					"end":        "10",
					"since":      "10",
				}
				ctx, rec := newGinTestCtx(http.MethodGet, invalidQueryParams, pathParams, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: unable to get patient timeline",
			setup: func(uc *restMocks.Mockusecases) args {
				pathParams := []gin.Param{
					{Key: "id", Value: gofakeit.UUID()},
				}

				uc.EXPECT().GetPatientTimeline(mock.Anything, mock.Anything, mock.AnythingOfType("*dto.PatientEverythingFilterParams")).
					RunAndReturn(func(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.HealthTimeline, error) {
						return nil, fmt.Errorf("error")
					})

				ctx, rec := newGinTestCtx(http.MethodGet, queryParams, pathParams, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: no patient id",
			setup: func(uc *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, nil)

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			args := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.GetPatientTimeline(args.c)
		})
	}
}

func TestPresentationHandlersImpl_ProvisionTenant(t *testing.T) {
	validInput := &dto.ProvisionTenantInput{
		TenantID:              gofakeit.UUID(),
		ParentID:              gofakeit.UUID(),
		Name:                  gofakeit.UUID(),
		LegacyIdentifierType:  dto.LegacyIdentifierTypeSladeCode,
		LegacyIdentifierValue: gofakeit.UUID(),
		Status:                dto.TenantStatusActive,
	}

	validPayload, err := json.Marshal(validInput)
	if err != nil {
		t.Fatalf("unable to marshal input: %v", err)
	}

	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy case: successfully provision tenant",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ProvisionTenant(mock.Anything, mock.AnythingOfType("dto.ProvisionTenantInput")).
					RunAndReturn(func(ctx context.Context, input dto.ProvisionTenantInput) (*dto.ProvisionTenantOutput, error) {
						return &dto.ProvisionTenantOutput{}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, validPayload, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: unable to provision tenant",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().ProvisionTenant(mock.Anything, mock.AnythingOfType("dto.ProvisionTenantInput")).
					RunAndReturn(func(ctx context.Context, input dto.ProvisionTenantInput) (*dto.ProvisionTenantOutput, error) {
						return nil, fmt.Errorf("provision failed")
					})
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, validPayload, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: invalid payload",
			setup: func(uc *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodPost, nil, nil, []byte(`{invalid`), nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			args := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.ProvisionTenant(args.c)

			require.Equal(t, tt.want, args.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestPresentationHandlersImpl_GetTenantProvisioningStatus(t *testing.T) {
	type args struct {
		c   *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(uc *restMocks.Mockusecases) args
		want  int
	}{
		{
			name: "Happy case: successfully get tenant provisioning status",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().GetTenantProvisioningStatus(mock.Anything, mock.AnythingOfType("string")).
					RunAndReturn(func(ctx context.Context, tenantID string) (*dto.ProvisionTenantOutput, error) {
						return &dto.ProvisionTenantOutput{}, nil
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, gin.Params{
					{Key: "tenant-id", Value: gofakeit.UUID()},
				}, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
		{
			name: "Sad case: missing tenant-id path parameter",
			setup: func(uc *restMocks.Mockusecases) args {
				ctx, rec := newGinTestCtx(http.MethodGet, nil, gin.Params{
					{Key: "tenant-id", Value: ""},
				}, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad case: tenant not found",
			setup: func(uc *restMocks.Mockusecases) args {
				uc.EXPECT().GetTenantProvisioningStatus(mock.Anything, mock.AnythingOfType("string")).
					RunAndReturn(func(ctx context.Context, tenantID string) (*dto.ProvisionTenantOutput, error) {
						return nil, fmt.Errorf("tenant not found")
					})
				ctx, rec := newGinTestCtx(http.MethodGet, nil, gin.Params{
					{Key: "tenant-id", Value: gofakeit.UUID()},
				}, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockusecases(t)
			mockExt := infraMocks.NewMockBaseExtension(t)
			mockAdv := advantageMocks.NewMockAdvantageService(t)

			args := tt.setup(mockUC)

			p := rest.NewPresentationHandlers(mockUC, mockExt, mockAdv)
			p.GetTenantProvisioningStatus(args.c)

			require.Equal(t, tt.want, args.rec.Code)
			mockUC.AssertExpectations(t)
		})
	}
}
