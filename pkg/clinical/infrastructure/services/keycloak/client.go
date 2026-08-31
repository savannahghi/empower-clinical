package keycloak

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/keycloak")

type Config struct {
	ClientID     string
	ClientSecret string
	Timeout      time.Duration
	Logger       *slog.Logger
	Realm        string
}

type LoginInfo struct {
	AccessToken string
	ExpiresIn   int
	TokenType   string
}

// TokenClaims represents the claims extracted from a JWT token.
type TokenClaims struct {
	Subject    string `json:"sub"`
	TenantID   string `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
	TenantType string `json:"tenant_type"`
	Exp        int64  `json:"exp"`
	Iat        int64  `json:"iat"`
}

type gocloaker interface {
	RetrospectToken(
		ctx context.Context,
		accessToken,
		clientID,
		clientSecret,
		realm string,
	) (*gocloak.IntroSpectTokenResult, error)
}

type Client struct {
	Gocloak gocloaker
	Config  Config
}

func (c *Config) Validate() error {
	if c.ClientID == "" {
		return errors.New("client ID is required")
	}

	if c.Logger == nil {
		return errors.New("logger is required")
	}

	if c.ClientSecret == "" {
		return errors.New("client secret is required")
	}

	if c.Realm == "" {
		return errors.New("keycloak realm is empty")
	}

	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	return nil
}

// New creates a new client that wraps github.com/Nerzal/gocloak.
func New(config Config, gocloak gocloaker) (*Client, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &Client{
		Config:  config,
		Gocloak: gocloak,
	}, nil
}

// IntrospectTokenWithClaims verifies the JWT and extracts tenant claims using Keycloak introspection.
func (cl *Client) IntrospectTokenWithClaims(
	ctx context.Context,
	accessToken string,
) (*TokenClaims, error) {
	ctx, span := tracer.Start(ctx, "IntrospectTokenWithClaims")
	defer span.End()

	result, err := cl.Gocloak.RetrospectToken(
		ctx,
		accessToken,
		cl.Config.ClientID,
		cl.Config.ClientSecret,
		cl.Config.Realm,
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return nil, fmt.Errorf("failed to introspect token: %w", err)
	}

	if result.Active == nil || !*result.Active {
		return nil, errors.New("token is not active")
	}

	// Since RetrospectToken doesn't return all claims, we need to parse the JWT locally
	// This is safe because we've already validated the token with Keycloak
	parser := jwt.Parser{}

	token, _, err := parser.ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims from token")
	}

	tokenClaims := &TokenClaims{
		Subject: getStringClaim(claims, "sub"),
	}

	if result.Exp != nil {
		tokenClaims.Exp = int64(*result.Exp)
	}

	if result.Iat != nil {
		tokenClaims.Iat = int64(*result.Iat)
	}

	return tokenClaims, nil
}

// getStringClaim safely extracts a string claim from the token.
func getStringClaim(claims jwt.MapClaims, key string) string {
	if value, ok := claims[key].(string); ok {
		return value
	}

	return ""
}
