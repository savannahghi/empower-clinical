package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/keycloak"
)

// fakeIntrospector stands in for the Keycloak client.
type fakeIntrospector struct {
	claims *keycloak.TokenClaims
	err    error
}

func (i *fakeIntrospector) IntrospectTokenWithClaims(_ context.Context, _ string) (*keycloak.TokenClaims, error) {
	if i.err != nil {
		return nil, i.err
	}

	return i.claims, nil
}

func testRouter(mw gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(mw)
	r.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}

func TestAuthenticationGinMiddleware(t *testing.T) {
	validClaims := &keycloak.TokenClaims{
		Subject: gofakeit.UUID(),
		Exp:     time.Now().Add(10 * time.Minute).Unix(),
	}

	tests := []struct {
		name         string
		authHeader   string
		introspector *fakeIntrospector
		wantStatus   int
	}{
		{
			name:         "valid token is accepted",
			authHeader:   "Bearer " + gofakeit.UUID(),
			introspector: &fakeIntrospector{claims: validClaims},
			wantStatus:   http.StatusOK,
		},
		{
			name:         "introspection failure is rejected",
			authHeader:   "Bearer " + gofakeit.UUID(),
			introspector: &fakeIntrospector{err: errors.New("introspection error")},
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "nil claims are rejected",
			authHeader:   "Bearer " + gofakeit.UUID(),
			introspector: &fakeIntrospector{},
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:       "expired token is rejected",
			authHeader: "Bearer " + gofakeit.UUID(),
			introspector: &fakeIntrospector{claims: &keycloak.TokenClaims{
				Subject: gofakeit.UUID(),
				Exp:     time.Now().Add(-5 * time.Minute).Unix(),
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:         "missing authorization header is rejected",
			authHeader:   "",
			introspector: &fakeIntrospector{claims: validClaims},
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "malformed authorization header is rejected",
			authHeader:   "Token missingBearerPrefix",
			introspector: &fakeIntrospector{claims: validClaims},
			wantStatus:   http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := testRouter(AuthenticationGinMiddleware(tt.introspector))

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusUnauthorized {
				var body map[string]any
				_ = json.Unmarshal(rec.Body.Bytes(), &body)
				assert.Equal(t, "Unauthorized. Authentication failed", body["message"])
			}
		})
	}
}

func TestExtractBearerToken_Helper(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	const token = "abc123"

	r.Header.Set("Authorization", "Bearer "+token)

	got, err := firebasetools.ExtractBearerToken(r)
	assert.NoError(t, err)
	assert.Equal(t, token, got)
}
