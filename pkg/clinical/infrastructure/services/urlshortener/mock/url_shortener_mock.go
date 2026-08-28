package mock

import (
	"context"
	"time"

	"github.com/savannahghi/silurlshortener"
)

// URLShortenerClientMock mocks the SIL's URL shortener client library implementations
type URLShortenerClientMock struct {
	MockShortenURLFn func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error)
}

// NewURLShortenerClientMock initializes client mocks
func NewURLShortenerClientMock() *URLShortenerClientMock {
	return &URLShortenerClientMock{
		MockShortenURLFn: func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
			return &silurlshortener.ShortenURLResponse{
				ShortURL:      "https://e.google/silurlshort",
				ShortCode:     "silurlshort",
				LongURL:       "",
				DateCreated:   time.Time{},
				Tags:          []string{},
				Meta:          silurlshortener.Meta{},
				Domain:        "",
				Title:         "",
				Crawlable:     false,
				ForwardQuery:  false,
				VisitsSummary: silurlshortener.VisitsSummary{},
			}, nil
		},
	}
}

// ShortenURL mocks the implementation of mocking SIL url shortener
func (sh *URLShortenerClientMock) ShortenURL(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
	return sh.MockShortenURLFn(ctx, payload)
}
