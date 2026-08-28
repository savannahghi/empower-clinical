package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
)

func TestValidateEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid  email address",
			args: args{
				email: firebasetools.TestUserEmail,
			},
			wantErr: false,
		},
		{
			name: "invalid email address",
			args: args{
				email: "hey@notavalidemail",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateEmail(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReportErrorToSentry(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy case",
			args: args{
				err: fmt.Errorf("test error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReportErrorToSentry(tt.args.err)
		})
	}
}

func TestAddPubSubNamespace(t *testing.T) {
	environment := fmt.Sprintf("service-test-%s-v2", serverutils.GetRunningEnvironment())
	type args struct {
		topicName   string
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy case: add pubsub namespace",
			args: args{
				topicName:   "test",
				serviceName: "service",
			},
			want: environment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddPubSubNamespace(tt.args.topicName, tt.args.serviceName); got != tt.want {
				t.Errorf("AddPubSubNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertDateStringToDateScalar(t *testing.T) {
	layout := time.RFC3339
	type args struct {
		layout     string
		dateString string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: convert date string to date",
			args: args{
				layout:     layout,
				dateString: "2024-05-11T09:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to convert date string to date",
			args: args{
				layout:     layout,
				dateString: "2024-0",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConvertDateStringToDateScalar(tt.args.layout, tt.args.dateString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertDateStringToDateScalar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNarrativeGenerator(t *testing.T) {
	status := "additional"
	type args struct {
		text   string
		status *string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy case: generate narrative",
			args: args{
				text:   "test",
				status: &status,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NarrativeGenerator(tt.args.text, tt.args.status)
		})
	}
}

func TestNewCustomError(t *testing.T) {
	type args struct {
		err     error
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy case: custom error",
			args: args{
				err:     fmt.Errorf("error"),
				message: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewCustomError(tt.args.err, tt.args.message)
		})
	}
}
