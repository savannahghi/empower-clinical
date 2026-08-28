package authutils

import (
	"context"
	"net/http"

	"github.com/savannahghi/authutils"
)

type OAuthClientService interface {
	HasValidSlade360BearerToken(ctx context.Context, r *http.Request) (bool, map[string]string, *authutils.TokenIntrospectionResponse)
	Authenticate() (*authutils.OAUTHResponse, error)
}

type OAuthClientImpl struct {
	Client OAuthClientService
}

func NewAuthClient(client OAuthClientService) *OAuthClientImpl {
	return &OAuthClientImpl{
		Client: client,
	}
}

func (o *OAuthClientImpl) Authenticate() (*authutils.OAUTHResponse, error) {
	return o.Client.Authenticate()
}

func (o *OAuthClientImpl) HasValidSlade360BearerToken(ctx context.Context, r *http.Request) (bool, map[string]string, *authutils.TokenIntrospectionResponse) {
	return o.Client.HasValidSlade360BearerToken(ctx, r)
}
