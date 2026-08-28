package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"github.com/sirupsen/logrus"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	auth "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/authutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/keycloak"
)

type (
	contextKey string
)

const (
	ContextKeyUserID contextKey = "user_id"
)

type TokenClaimsIntrospector interface {
	IntrospectTokenWithClaims(ctx context.Context, accessToken string) (*keycloak.TokenClaims, error)
}

// authCheckFn is a function type for authorization and authentication checks.
// there can be several e.g an authentication check runs first then an authorization.
// check runs next if the authentication passes etc.
type authCheckFn = func(
	ctx context.Context,
	r *http.Request,
) (bool, map[string]string, *authutils.TokenIntrospectionResponse)

// HasValidCachedToken returns an authentication check function for verifying
// the validity of a token stored in the provided cache store.
//
// Parameters:
// - cacheStore (persist.CacheStore): The cache store used for storing tokens.
func HasValidCachedToken(cacheStore persist.CacheStore) authCheckFn {
	return func(_ context.Context, r *http.Request) (bool, map[string]string, *authutils.TokenIntrospectionResponse) {
		token, err := firebasetools.ExtractBearerToken(r)
		if err != nil {
			utils.ReportErrorToSentry(err)
			return false, serverutils.ErrorMap(err), nil
		}

		tokenResponse := authutils.TokenIntrospectionResponse{}

		err = cacheStore.Get(token, &tokenResponse)
		if err != nil {
			utils.ReportErrorToSentry(err)

			if errors.Is(err, persist.ErrCacheMiss) {
				return false, serverutils.ErrorMap(errors.New("supplied access token not in cache")), nil
			}

			return false, serverutils.ErrorMap(err), nil
		}

		return true, nil, &tokenResponse
	}
}

// introspectAuthServerToken verifies the validity of authserver token
// Returns the token response on success & error on failure.
func introspectAuthServerToken(
	ctx context.Context,
	r *http.Request,
	cacheStore persist.CacheStore,
	primary auth.OAuthClientService,
) (*authutils.TokenIntrospectionResponse, error) {
	checkFuncs := []authCheckFn{
		HasValidCachedToken(cacheStore),
		primary.HasValidSlade360BearerToken,
	}

	var lastErr error

	for _, check := range checkFuncs {
		ok, errMap, tok := check(ctx, r)
		if ok && tok != nil {
			if !tok.Expires.IsZero() {
				if err := cacheStore.Set(tok.Token, *tok, time.Until(tok.Expires)); err != nil {
					logrus.WithError(err).Warn("failed to set token in cache")
				}
			}

			return tok, nil
		}

		// Skip token-length validation errors — let the next check (or Keycloak) handle it
		if errMap != nil {
			if msg, exists := errMap["error"]; exists && strings.Contains(msg, "255 characters") {
				continue
			}
		}

		lastErr = fmt.Errorf("%v", errMap)
	}

	if lastErr == nil {
		lastErr = errors.New("authutils chain failed")
	}

	return nil, lastErr
}

// introspectKeyCloakToken is used to verify KeyCloak auth token
// Returns claims on success & an error on failure.
func introspectKeyCloakToken(
	parentCtx context.Context,
	r *http.Request,
	introspector TokenClaimsIntrospector,
) (*keycloak.TokenClaims, error) {
	token, err := firebasetools.ExtractBearerToken(r)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("extract bearer token: %w", err)
	}

	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	claims, err := introspector.IntrospectTokenWithClaims(ctx, token)
	if err != nil {
		utils.ReportErrorToSentry(err)

		return nil, err
	}

	if claims == nil {
		errMsg := fmt.Errorf("nil claims returned from token introspection")
		utils.ReportErrorToSentry(errMsg)

		return nil, errMsg
	}

	if claims.Exp > 0 && time.Unix(claims.Exp, 0).Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// AuthenticationGinMiddleware tries AuthenticationGinMiddleware *first* (authutils path).
// If it fails, it tries the Keycloak Authentication path.
// Only aborts if both fail.
func AuthenticationGinMiddleware(
	cacheStore persist.CacheStore,
	primary auth.OAuthClientService,
	introspector TokenClaimsIntrospector,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// AuthenticationGinMiddleware first
		if tok, err := introspectAuthServerToken(c.Request.Context(), c.Request, cacheStore, primary); err == nil {
			ctx := context.WithValue(c.Request.Context(), authutils.AuthTokenContextKey, tok)

			c.Request = c.Request.WithContext(ctx)
			c.Next()

			return
		} else {
			utils.ReportErrorToSentry(err)
			logrus.WithError(err).Debug("authserver token introspection failed")
		}

		// Keycloak path if GinAuth fails
		if claims, err := introspectKeyCloakToken(c.Request.Context(), c.Request, introspector); err == nil {
			ctx := context.WithValue(c.Request.Context(), ContextKeyUserID, claims.Subject)

			c.Request = c.Request.WithContext(ctx)
			c.Next()

			return
		} else {
			utils.ReportErrorToSentry(err)
			logrus.WithError(err).Debug("keycloak auth token introspection failed")
		}

		// Both failed, abort!
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized. Authentication failed",
		})
	}
}
