package mock

import (
	"context"
	"time"

	"github.com/savannahghi/silurlshortener"
)

// URLShortenerServiceMock mocks SIL's url shortener service implementations
type URLShortenerServiceMock struct {
	MockShortenFn func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error)
}

// NewSMSClientMock initializes our client mocks
func NewURLShortenerServiceMock() *URLShortenerServiceMock {
	return &URLShortenerServiceMock{
		MockShortenFn: func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
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

// Shorten mocks the implementation of shortening the URLs
func (m *URLShortenerServiceMock) Shorten(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
	return m.MockShortenFn(ctx, payload)
}
