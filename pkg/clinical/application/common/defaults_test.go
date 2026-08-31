package common

import (
	"testing"
	"time"

	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/assert"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

func TestDefaultIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		wantEmpty bool
	}{
		{
			name:      "Happy Case: set default identifier",
			wantEmpty: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultIdentifier()
			if !tt.wantEmpty {
				assert.NotEmpty(t, got, "DefaultIdentifier() = empty struct, want non-empty struct")
				assert.NotEmpty(t, got.Period, "DefaultIdentifier() = %v, Period field is empty", got)
			}
			if got != nil && tt.wantEmpty {
				t.Errorf("DefaultIdentifier() = %v, want %v", got == nil, tt.wantEmpty)
			}

		})
	}
}

func TestDefaultPeriod(t *testing.T) {
	tests := []struct {
		name      string
		wantEmpty bool
	}{
		{
			name:      "Happy case: set period",
			wantEmpty: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultPeriod()
			if !tt.wantEmpty {
				assert.NotEmpty(t, got, "DefaultPeriod() = empty struct, want non-empty struct")
			}
			if got != nil && tt.wantEmpty {
				t.Errorf("DefaultPeriod() = %v, want %v", got == nil, tt.wantEmpty)
			}
		})
	}
}

func TestDefaultPeriodInput(t *testing.T) {
	featureTime := time.Now().Add(time.Hour * NinentyYears)
	startDate := scalarutils.DateTime(time.Now().Format(timeFormatStr))
	endDate := scalarutils.DateTime(featureTime.Format(timeFormatStr))
	type args struct {
		output *domain.FHIRPeriod
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy case: Default period input",
			args: args{
				output: &domain.FHIRPeriod{
					Start: startDate,
					End:   endDate,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultPeriodInput()
			if *got.Start != tt.args.output.Start {
				t.Errorf("DefaultPeriodInput()")
				return
			}
		})
	}
}
