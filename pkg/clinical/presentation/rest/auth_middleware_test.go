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
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
	auth "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/authutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/keycloak"
)

// in-memory cache store for tests
var memoryStore = persist.NewMemoryStore(60 * time.Minute)

// fake primary OAuth client service
type fakePrimary struct {
	ok     bool
	resp   *authutils.TokenIntrospectionResponse
	errMap map[string]string
}

func (p *fakePrimary) HasValidSlade360BearerToken(ctx context.Context, r *http.Request) (bool, map[string]string, *authutils.TokenIntrospectionResponse) {
	if p.ok {
		return true, nil, p.resp
	}
	return false, p.errMap, nil
}

// fake Authenticate to satisfy authutils interface
func (p *fakePrimary) Authenticate() (*authutils.OAUTHResponse, error) {
	return nil, errors.New("unauthenticated")
}

var _ auth.OAuthClientService = (*fakePrimary)(nil)

// fake Keycloak introspector
type fakeIntrospector struct {
	claims *keycloak.TokenClaims
	err    error
}

func (i *fakeIntrospector) IntrospectTokenWithClaims(ctx context.Context, accessToken string) (*keycloak.TokenClaims, error) {
	if i.err != nil {
		return nil, i.err
	}
	return i.claims, nil
}

func withAuthHeader(r *http.Request, token string) {
	r.Header.Set("Authorization", "Bearer "+token)
}

func withMalformedAuthHeader(r *http.Request, raw string) {
	r.Header.Set("Authorization", raw)
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

// ensure clean state per test by resetting keys since memoryStore is shared
func cacheSetToken(t *testing.T, token string, resp authutils.TokenIntrospectionResponse) {
	t.Helper()
	_ = memoryStore.Set(token, resp, time.Until(resp.Expires))
}

func cacheDeleteToken(t *testing.T, token string) {
	t.Helper()
	_ = memoryStore.Delete(token)
}

func TestAuthenticationGinMiddleware_AuthServer_Succeeds_FromCache(t *testing.T) {
	token := gofakeit.UUID()
	resp := authutils.TokenIntrospectionResponse{
		Token:   token,
		Expires: time.Now().Add(30 * time.Minute),
	}

	cacheSetToken(t, token, resp)
	t.Cleanup(func() { cacheDeleteToken(t, token) })

	primary := &fakePrimary{ok: false}
	introspector := &fakeIntrospector{}

	mw := AuthenticationGinMiddleware(memoryStore, primary, introspector)
	router := testRouter(mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	withAuthHeader(req, token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthenticationGinMiddleware_Keycloak_Succeeds_WhenAuthServerFails(t *testing.T) {
	token := gofakeit.UUID()
	_ = memoryStore.Delete(token)

	primary := &fakePrimary{
		ok:     false,
		errMap: map[string]string{"error": "authserver failed"},
	}
	introspector := &fakeIntrospector{
		claims: &keycloak.TokenClaims{
			Subject: gofakeit.UUID(),
			Exp:     time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	mw := AuthenticationGinMiddleware(memoryStore, primary, introspector)
	router := testRouter(mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	withAuthHeader(req, token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthenticationGinMiddleware_Fails_WhenBothPathsFail(t *testing.T) {
	token := gofakeit.UUID()
	_ = memoryStore.Delete(token)

	primary := &fakePrimary{
		ok:     false,
		errMap: map[string]string{"error": "authserver failed"},
	}
	introspector := &fakeIntrospector{
		err: errors.New("introspection error"),
	}

	mw := AuthenticationGinMiddleware(memoryStore, primary, introspector)
	router := testRouter(mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	withAuthHeader(req, token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	assert.Equal(t, "Unauthorized. Authentication failed", body["message"])
}

func TestAuthenticationGinMiddleware_KeycloakExpiredToken_Fails(t *testing.T) {
	token := gofakeit.UUID()
	_ = memoryStore.Delete(token)

	primary := &fakePrimary{
		ok:     false,
		errMap: map[string]string{"error": "authserver failed"},
	}

	introspector := &fakeIntrospector{
		claims: &keycloak.TokenClaims{
			Subject: "kc-user-2",
			Exp:     time.Now().Add(-5 * time.Minute).Unix(),
		},
	}

	mw := AuthenticationGinMiddleware(memoryStore, primary, introspector)
	router := testRouter(mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	withAuthHeader(req, token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// No Authorization header present.
// HasValidCachedToken fails to extract token, returns error map.
// Primary path succeeds, middleware should allow request.
func TestAuthenticationGinMiddleware_ExtractBearer_MissingHeader_PrimarySucceeds(t *testing.T) {
	token := gofakeit.UUID()
	primary := &fakePrimary{
		ok: true,
		resp: &authutils.TokenIntrospectionResponse{
			Token:   token,
			Expires: time.Now().Add(15 * time.Minute),
		},
	}
	introspector := &fakeIntrospector{}

	mw := AuthenticationGinMiddleware(memoryStore, primary, introspector)
	router := testRouter(mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// No Authorization header present.
// HasValidCachedToken fails to extract token; both Authserver & Keycloak path also fails to extract token.
func TestAuthenticationGinMiddleware_ExtractBearer_MissingHeader_BothFail(t *testing.T) {
	primary := &fakePrimary{
		ok:     false,
		errMap: map[string]string{"error": "authserver failed"},
	}

	introspector := &fakeIntrospector{}

	mw := AuthenticationGinMiddleware(memoryStore, primary, introspector)
	router := testRouter(mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil) // no header
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// Malformed Authorization header. ExtractBearerToken should return an error on both cache path and keycloak path.
func TestAuthenticationGinMiddleware_ExtractBearer_MalformedHeader_BothFail(t *testing.T) {
	primary := &fakePrimary{
		ok:     false,
		errMap: map[string]string{"error": "authserver failed"},
	}

	introspector := &fakeIntrospector{}

	mw := AuthenticationGinMiddleware(memoryStore, primary, introspector)
	router := testRouter(mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	withMalformedAuthHeader(req, "Token missingBearerPrefix")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestExtractBearerToken_Helper(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	const token = "abc123"
	withAuthHeader(r, token)

	got, err := firebasetools.ExtractBearerToken(r)
	assert.NoError(t, err)
	assert.Equal(t, token, got)
}

func TestAuthenticationGinMiddleware_KeycloakNilClaims_Fails(t *testing.T) {
	token := gofakeit.UUID()
	_ = memoryStore.Delete(token)

	primary := &fakePrimary{
		ok:     false, // force authserver check to fail
		errMap: map[string]string{"error": "authserver failed"},
	}

	// nil claims, fail
	introspector := &fakeIntrospector{
		claims: nil,
		err:    nil,
	}

	mw := AuthenticationGinMiddleware(memoryStore, primary, introspector)
	router := testRouter(mw)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	withAuthHeader(req, token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	assert.Equal(t, "Unauthorized. Authentication failed", body["message"])
}
