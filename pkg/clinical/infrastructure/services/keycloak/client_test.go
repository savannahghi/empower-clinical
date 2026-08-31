package keycloak_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/keycloak"
	kcmock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/keycloak/mock"
)

func TestNew(t *testing.T) {
	config := keycloak.Config{
		ClientID: os.Getenv("KEYCLOAK_CLIENT_ID"),
		Logger:   slog.Default(),
	}
	tests := []struct {
		name string
		want *keycloak.Client
	}{
		{
			name: "Happy case: successfully created a keycloack instance",
		},
	}
	for _, tt := range tests {
		t.Setenv("KEYCLOAK_HOST", "http://localhost")
		t.Setenv("KEYCLOAK_CLIENT_ID", "ClientID")
		t.Setenv("KEYCLOAK_CLIENT_SECRET", "ClientSecret")
		t.Setenv("KEYCLOAK_REALM", "realm")
		t.Run(tt.name, func(t *testing.T) {
			keycloak.New(config, &kcmock.Mockgocloaker{})
		})
	}
}

func TestClient_IntrospectTokenWithClaims(t *testing.T) {
	type args struct {
		ctx         context.Context
		accessToken string
	}
	tests := []struct {
		name       string
		args       args
		config     keycloak.Config
		setupMocks func(*kcmock.Mockgocloaker)
		want       *keycloak.TokenClaims
		wantErr    bool
	}{
		{
			name: "Happy case: successfully introspected token with claims",
			args: args{
				ctx:         context.Background(),
				accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjMiLCJ0ZW5hbnRfaWQiOiJ0ZW5hbnQtMTIzIiwidGVuYW50X25hbWUiOiJUZXN0IEhvc3BpdGFsIiwidGVuYW50X3R5cGUiOiJIT1NQSVRBTCIsImV4cCI6MTczNTY4MDAwMCwiaWF0IjoxNzM1Njc5OTAwfQ.Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8", //gitleaks:allow
			},
			config: keycloak.Config{
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Realm:        "test-realm",
				Logger:       slog.Default(),
			},
			setupMocks: func(mockGocloak *kcmock.Mockgocloaker) {
				active := true
				exp := int(1735680000)
				iat := int(1735679900)
				mockGocloak.EXPECT().RetrospectToken(mock.Anything, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjMiLCJ0ZW5hbnRfaWQiOiJ0ZW5hbnQtMTIzIiwidGVuYW50X25hbWUiOiJUZXN0IEhvc3BpdGFsIiwidGVuYW50X3R5cGUiOiJIT1NQSVRBTCIsImV4cCI6MTczNTY4MDAwMCwiaWF0IjoxNzM1Njc5OTAwfQ.Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8Ej8", "test-client", "test-secret", "test-realm"). //gitleaks:allow
																																																Return(&gocloak.IntroSpectTokenResult{
						Active: &active,
						Exp:    &exp,
						Iat:    &iat,
					}, nil)
			},
			want: &keycloak.TokenClaims{
				Subject: "123",
				Exp:     1735680000,
				Iat:     1735679900,
			},
			wantErr: false,
		},
		{
			name: "Sad case: token introspection fails",
			args: args{
				ctx:         context.Background(),
				accessToken: "invalid-token",
			},
			config: keycloak.Config{
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Realm:        "test-realm",
				Logger:       slog.Default(),
			},
			setupMocks: func(mockGocloak *kcmock.Mockgocloaker) {
				mockGocloak.EXPECT().RetrospectToken(mock.Anything, "invalid-token", "test-client", "test-secret", "test-realm").
					Return(nil, errors.New("introspection failed"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: token is not active",
			args: args{
				ctx:         context.Background(),
				accessToken: "expired-token",
			},
			config: keycloak.Config{
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Realm:        "test-realm",
				Logger:       slog.Default(),
			},
			setupMocks: func(mockGocloak *kcmock.Mockgocloaker) {
				active := false
				mockGocloak.EXPECT().RetrospectToken(mock.Anything, "expired-token", "test-client", "test-secret", "test-realm").
					Return(&gocloak.IntroSpectTokenResult{
						Active: &active,
					}, nil)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGocloak := kcmock.NewMockgocloaker(t)

			if tt.setupMocks != nil {
				tt.setupMocks(mockGocloak)
			}

			cl := &keycloak.Client{
				Config:  tt.config,
				Gocloak: mockGocloak,
			}

			got, err := cl.IntrospectTokenWithClaims(tt.args.ctx, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.IntrospectTokenWithClaims() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.IntrospectTokenWithClaims() = %v, want %v", got, tt.want)
			}

			mockGocloak.AssertExpectations(t)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		ClientID     string
		ClientSecret string
		Timeout      time.Duration
		Logger       *slog.Logger
		Realm        string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Sad case: Client ID missing",
			fields: fields{
				ClientSecret: "middleware",
				Timeout:      10 * time.Second,
				Logger:       slog.Default(),
				Realm:        "master",
			},
			wantErr: true,
		},
		{
			name: "Sad case: Client secret missing",
			fields: fields{
				ClientID: "middleware",
				Timeout:  10 * time.Second,
				Logger:   slog.Default(),
				Realm:    "master",
			},
			wantErr: true,
		},
		{
			name: "Sad case: logger missing",
			fields: fields{
				ClientID:     "middleware",
				ClientSecret: "middleware",
				Timeout:      10 * time.Second,
				Realm:        "master",
			},
			wantErr: true,
		},
		{
			name: "Happy case: config has all fields present",
			fields: fields{
				ClientID:     "middleware",
				ClientSecret: "middleware",
				Timeout:      10 * time.Second,
				Logger:       slog.Default(),
				Realm:        "master",
			},
			wantErr: false,
		},
		{
			name: "Sad case: keycloak realm is not set",
			fields: fields{
				ClientID:     "middleware",
				ClientSecret: "middleware",
				Timeout:      10 * time.Second,
				Logger:       slog.Default(),
			},
			wantErr: true,
		},
		{
			name: "Should not fail if timeout is not set",
			fields: fields{
				ClientID:     "middleware",
				ClientSecret: "middleware",
				Logger:       slog.Default(),
				Realm:        "master",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &keycloak.Config{
				ClientID:     tt.fields.ClientID,
				ClientSecret: tt.fields.ClientSecret,
				Timeout:      tt.fields.Timeout,
				Logger:       tt.fields.Logger,
				Realm:        tt.fields.Realm,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
