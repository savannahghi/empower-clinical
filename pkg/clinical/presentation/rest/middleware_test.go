package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/presentation/rest"
	restMocks "github.com/savannahghi/empower-clinical/pkg/clinical/presentation/rest/mock"
)

func TestOrganisationValidator(t *testing.T) {
	type args struct {
		v          rest.Validators
		identifier string
	}
	tests := []struct {
		name    string
		setup   func(v *restMocks.MockValidators) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get organisation",
			setup: func(v *restMocks.MockValidators) args {
				id := uuid.NewString()
				v.On("GetFHIROrganization", mock.Anything, id).
					Return(&domain.FHIROrganizationRelayPayload{}, nil)

				return args{v, id}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get organisation",
			setup: func(v *restMocks.MockValidators) args {
				id := uuid.NewString()
				v.On("GetFHIROrganization", mock.Anything, id).
					Return(nil, errors.New("failed to get organisation"))

				return args{v, id}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := restMocks.NewMockValidators(t)
			arg := tt.setup(mockUC)

			if err := rest.OrganisationValidator(arg.v, arg.identifier); (err != nil) != tt.wantErr {
				t.Errorf("OrganisationValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTenantIdentifierExtractionMiddleware(t *testing.T) {
	type args struct {
		ctx *gin.Context
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name  string
		setup func(v *restMocks.MockValidators) args
		want  int
	}{
		{
			name: "Sad Case: organisation header missing",
			setup: func(v *restMocks.MockValidators) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"Clinical-Facility-ID": "FAC",
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: facility header missing",
			setup: func(v *restMocks.MockValidators) args {
				ctx, rec := newGinTestCtx(http.MethodGet, map[string]string{
					"Clinical-Organization-ID": "ORG",
				}, nil, nil, nil)
				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Sad Case: validator rejects header",
			setup: func(v *restMocks.MockValidators) args {
				orgID := uuid.New().String()
				facilityID := uuid.New().String()

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, map[string]string{
					"Clinical-Organization-ID": orgID,
					"Clinical-Facility-ID":     facilityID,
				})

				v.On("GetFHIROrganization", mock.Anything, orgID).
					Return(nil, errors.New("oops"))

				return args{ctx, rec}
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Happy Case: valid headers",
			setup: func(v *restMocks.MockValidators) args {
				orgID := uuid.New().String()
				facilityID := uuid.New().String()

				ctx, rec := newGinTestCtx(http.MethodGet, nil, nil, nil, map[string]string{
					"Clinical-Organization-ID": orgID,
					"Clinical-Facility-ID":     facilityID,
				})

				v.On("GetFHIROrganization", mock.Anything, orgID).
					Return(&domain.FHIROrganizationRelayPayload{}, nil)
				v.On("GetFHIROrganization", mock.Anything, facilityID).
					Return(&domain.FHIROrganizationRelayPayload{}, nil)

				return args{ctx, rec}
			},
			want: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUC := restMocks.NewMockValidators(t)
			arg := tt.setup(mockUC)

			mw := rest.TenantIdentifierExtractionMiddleware(mockUC)

			mw(arg.ctx)
			require.Equal(t, tt.want, arg.rec.Code)

			if tt.want == http.StatusOK {
				_, ok1 := arg.ctx.Get(string(utils.OrganizationIDContextKey))
				_, ok2 := arg.ctx.Get(string(utils.FacilityIDContextKey))
				require.True(t, ok1 && ok2)
			}

			mockUC.AssertExpectations(t)
		})
	}
}
