package mail_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/mail"
	mailMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/mail/mock"
)

type mockHandler struct {
	mail *mailMock.MockMailSender
}

func TestMailSenderImpl_SendEmail(t *testing.T) {
	type args struct {
		ctx    context.Context
		params *ses.SendEmailInput
		optFns []func(*ses.Options)
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *ses.SendEmailOutput
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully send email",
			setup: func(mh *mockHandler) args {
				mh.mail.EXPECT().SendEmail(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
						return &ses.SendEmailOutput{}, nil
					})
				return args{
					ctx:    context.TODO(),
					params: &ses.SendEmailInput{},
					optFns: []func(*ses.Options){},
				}
			},
			want:    &ses.SendEmailOutput{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mailMock.NewMockMailSender(t)
			m := mail.NewMailSender(mock)
			args := tt.setup(&mockHandler{mail: mock})

			got, err := m.SendEmail(args.ctx, args.params, args.optFns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MailSenderImpl.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MailSenderImpl.SendEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMailSenderImpl_SendRawEmail(t *testing.T) {
	type args struct {
		ctx    context.Context
		params *ses.SendRawEmailInput
		optFns []func(*ses.Options)
	}
	tests := []struct {
		name    string
		setup   func(mh *mockHandler) args
		want    *ses.SendRawEmailOutput
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully send email",
			setup: func(mh *mockHandler) args {
				mh.mail.EXPECT().SendRawEmail(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error) {
						return &ses.SendRawEmailOutput{}, nil
					})
				return args{
					ctx:    context.TODO(),
					params: &ses.SendRawEmailInput{},
					optFns: []func(*ses.Options){},
				}
			},
			want:    &ses.SendRawEmailOutput{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mailMock.NewMockMailSender(t)
			m := mail.NewMailSender(mock)
			args := tt.setup(&mockHandler{mail: mock})

			got, err := m.SendRawEmail(args.ctx, args.params, args.optFns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MailSenderImpl.SendRawEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MailSenderImpl.SendRawEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
