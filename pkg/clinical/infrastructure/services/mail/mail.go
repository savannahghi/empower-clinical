package mail

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
)

// MailSender is an interface defining the contract for sending emails.
type MailSender interface {
	SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
	SendRawEmail(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error)
}

// MailSenderImpl is a concrete implementation of the MailSender interface.
type MailSenderImpl struct {
	client MailSender
}

// NewMailSender creates a new instance of MailSenderImpl with the provided sender
func NewMailSender(sender MailSender) *MailSenderImpl {
	return &MailSenderImpl{
		client: sender,
	}
}

// SendEmail sends an email using the configured MailSender client, in this case, it is amazon SES.
func (m MailSenderImpl) SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
	return m.client.SendEmail(ctx, params, optFns...)
}

// SendRawEmail more flexible than the SendEmail operation. When you use the SendRawEmail operation,
// you can specify the headers of the message as well as its content.
func (m MailSenderImpl) SendRawEmail(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error) {
	return m.client.SendRawEmail(ctx, params, optFns...)
}
