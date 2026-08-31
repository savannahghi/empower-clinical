package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
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

// introspectKeyCloakToken verifies a Keycloak bearer token and returns its claims.
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

// AuthenticationGinMiddleware authenticates a request against Keycloak.
func AuthenticationGinMiddleware(introspector TokenClaimsIntrospector) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := introspectKeyCloakToken(c.Request.Context(), c.Request, introspector)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized. Authentication failed",
			})

			return
		}

		ctx := context.WithValue(c.Request.Context(), ContextKeyUserID, claims.Subject)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
