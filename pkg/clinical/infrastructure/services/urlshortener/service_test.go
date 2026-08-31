package urlshortener_test

import (
	"context"
	"errors"
	"testing"

	"github.com/savannahghi/silurlshortener"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/urlshortener"
	mockURLShortener "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/urlshortener/mock"
)

func TestURLShortenerServiceImpl_Shorten(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload *silurlshortener.ShortenURLPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: shorten URL",
			args: args{
				ctx: context.Background(),
				payload: &silurlshortener.ShortenURLPayload{
					LongURL:         "http://gooooooooooooooooo.com",
					ValidSince:      "",
					ValidUntil:      "",
					MaxVisits:       0,
					Tags:            []string{},
					Title:           "",
					Crawlable:       false,
					ForwardQuery:    false,
					CustomSlug:      "",
					PathPrefix:      "",
					FindIfExists:    false,
					Domain:          "",
					ShortCodeLength: 5,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to shorten URL",
			args: args{
				ctx: context.Background(),
				payload: &silurlshortener.ShortenURLPayload{
					LongURL:         "http://gooooooooooooooooo.com",
					ValidSince:      "",
					ValidUntil:      "",
					MaxVisits:       0,
					Tags:            []string{},
					Title:           "",
					Crawlable:       false,
					ForwardQuery:    false,
					CustomSlug:      "",
					PathPrefix:      "",
					FindIfExists:    false,
					Domain:          "",
					ShortCodeLength: 5,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockURLShortener.NewURLShortenerClientMock()
			us := urlshortener.NewURLShortenerService(fakeClient)

			if tt.name == "Sad case: unable to shorten URL" {
				fakeClient.MockShortenURLFn = func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
					return nil, errors.New("unable to shorten URL")
				}
			}

			_, err := us.Shorten(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLShortenerServiceImpl.Shorten() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
