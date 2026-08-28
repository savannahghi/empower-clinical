package urlshortener

import (
	"context"

	"github.com/savannahghi/silurlshortener"
)

// IServiceURLShortener holds the methods to interact with the SIL's URL Shortener service
type IServiceURLShortener interface {
	Shorten(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error)
}

// IURLShortenerClient defines the methods used to communicate with the SIL's url shortener library.
type IURLShortenerClient interface {
	ShortenURL(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error)
}

// URLShortenerServiceImpl implements SIL's url shortener service interface
type URLShortenerServiceImpl struct {
	client IURLShortenerClient
}

// NewURLShortenerService initializes SIL's URL shortener client
func NewURLShortenerService(client IURLShortenerClient) *URLShortenerServiceImpl {
	return &URLShortenerServiceImpl{
		client: client,
	}
}

// Shorten is used to shorten any given URL
func (s *URLShortenerServiceImpl) Shorten(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
	return s.client.ShortenURL(ctx, payload)
}
