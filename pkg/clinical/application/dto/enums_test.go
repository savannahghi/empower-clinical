package dto

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
)

func TestConsentState_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     SegmentationCategory
		wantW string
	}{
		{
			name:  "valid type s",
			c:     SegmentationCategoryLowRisk,
			wantW: strconv.Quote("CERVICAL_CANCER_LOW_RISK"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SegmentationCategory.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestSegmentationCategory_String(t *testing.T) {
	tests := []struct {
		name string
		e    SegmentationCategory
		want string
	}{
		{
			name: "CERVICAL_CANCER_TIPS",
			e:    SegmentationCategoryNoRisk,
			want: "CERVICAL_CANCER_TIPS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("SegmentationCategory.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSegmentationCategory_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    SegmentationCategory
		want bool
	}{
		{
			name: "valid type",
			e:    SegmentationCategoryHighRiskNegative,
			want: true,
		},
		{
			name: "invalid type",
			e:    SegmentationCategory("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("SegmentationCategory.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSegmentationCategory_UnmarshalGQL(t *testing.T) {
	value := SegmentationCategoryHighRiskNegative
	invalid := SegmentationCategory("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *SegmentationCategory
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "CERVICAL_CANCER_HIGH_RISK",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SegmentationCategory.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSegmentationCategory_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     SegmentationCategory
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     SegmentationCategoryHighRiskNegative,
			b:     w,
			wantW: strconv.Quote("CERVICAL_CANCER_HIGH_RISK"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SegmentationCategory.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestSAppointment_String(t *testing.T) {
	tests := []struct {
		name string
		e    AppointmentStatus
		want string
	}{
		{
			name: "noshow",
			e:    AppointmentStatusNoShow,
			want: "noshow",
		},
		{
			name: "checked-in",
			e:    AppointmentStatus(AppointmentStatusCheckedIn.String()),
			want: "checked-in",
		},
		{
			name: "entered-in-error",
			e:    AppointmentStatus(AppointmentStatusEnteredInError.String()),
			want: "entered-in-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AppointmentStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppointmentStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AppointmentStatus
		want bool
	}{
		{
			name: "valid type",
			e:    AppointmentStatusCheckedIn,
			want: true,
		},
		{
			name: "invalid type",
			e:    AppointmentStatus("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AppointmentStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppointmentStatus_UnmarshalGQL(t *testing.T) {
	value := AppointmentStatusBooked
	invalid := AppointmentStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AppointmentStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "booked",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("AppointmentStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppointmentStatus_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     AppointmentStatus
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     AppointmentStatusBooked,
			b:     w,
			wantW: strconv.Quote("booked"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AppointmentStatus.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestParticipationStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    ParticipationStatus
		want string
	}{
		{
			name: "accepted",
			e:    ParticipationStatusAccepted,
			want: "accepted",
		},
		{
			name: "needs-action",
			e:    ParticipationStatus(ParticipationStatusNeedsAction.String()),
			want: "needs-action",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ParticipationStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParticipationStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ParticipationStatus
		want bool
	}{
		{
			name: "valid type",
			e:    ParticipationStatusAccepted,
			want: true,
		},
		{
			name: "invalid type",
			e:    ParticipationStatus("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ParticipationStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParticipationStatus_UnmarshalGQL(t *testing.T) {
	value := ParticipationStatusAccepted
	invalid := ParticipationStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ParticipationStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "accepted",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ParticipationStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParticipationStatus_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     ParticipationStatus
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     ParticipationStatusAccepted,
			b:     w,
			wantW: strconv.Quote("accepted"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ParticipationStatus.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestParticipantRequired_String(t *testing.T) {
	tests := []struct {
		name string
		e    ParticipantRequired
		want string
	}{
		{
			name: "required",
			e:    ParticipantRequiredRequired,
			want: "required",
		},
		{
			name: "information-only",
			e:    ParticipantRequired(ParticipantRequiredInformationOnly.String()),
			want: "information-only",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ParticipantRequired.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParticipantRequired_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ParticipantRequired
		want bool
	}{
		{
			name: "valid type",
			e:    ParticipantRequiredRequired,
			want: true,
		},
		{
			name: "invalid type",
			e:    ParticipantRequired("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ParticipantRequired.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParticipantRequired_UnmarshalGQL(t *testing.T) {
	value := ParticipantRequiredRequired
	invalid := ParticipantRequired("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ParticipantRequired
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "required",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ParticipantRequired.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParticipantRequired_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     ParticipantRequired
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     ParticipantRequiredRequired,
			b:     w,
			wantW: strconv.Quote("required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ParticipantRequired.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestScreeningTypeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ScreeningTypeEnum
		want bool
	}{
		{
			name: "Valid screening type - breast cancer",
			e:    BreastCancerScreeningTypeEnum,
			want: true,
		},
		{
			name: "Invalid screening type status",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ScreeningTypeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningTypeEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    ScreeningTypeEnum
		want string
	}{
		{
			name: "Breast Cancer Screening",
			e:    BreastCancerScreeningTypeEnum,
			want: "BREAST_CANCER_SCREENING",
		},
		{
			name: "Cervical Cancer Screening",
			e:    CervicalCancerScreeningTypeEnum,
			want: "CERVICAL_CANCER_SCREENING",
		},
		{
			name: "Prostate Cancer Screening",
			e:    ProstateCancerScreeningTypeEnum,
			want: "PROSTATE_CANCER_SCREENING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ScreeningTypeEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningTypeEnum_UnmarshalGQL(t *testing.T) {
	value := BreastCancerScreeningTypeEnum
	invalidType := ScreeningTypeEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ScreeningTypeEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "BREAST_CANCER_SCREENING",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalidType,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidType,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ScreeningTypeEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScreeningTypeEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ScreeningTypeEnum
		wantW string
	}{
		{
			name:  "BREAST_CANCER_SCREENING",
			e:     BreastCancerScreeningTypeEnum,
			wantW: strconv.Quote("BREAST_CANCER_SCREENING"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ScreeningTypeEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestScreeningTypeEnum_Text(t *testing.T) {
	tests := []struct {
		name          string
		screeningType ScreeningTypeEnum
		want          string
	}{
		{
			name:          "Valid type - breast",
			screeningType: BreastCancerScreeningTypeEnum,
			want:          "Breast Cancer Screening",
		},
		{
			name:          "Valid type - cervical",
			screeningType: CervicalCancerScreeningTypeEnum,
			want:          "Cervical Cancer Screening",
		},
		{
			name:          "Valid type - prostate",
			screeningType: ProstateCancerScreeningTypeEnum,
			want:          "Prostate Cancer Screening",
		},
		{
			name:          "Invalid type",
			screeningType: ScreeningTypeEnum("invalid"),
			want:          "unknown screening type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.screeningType.Text(); got != tt.want {
				t.Errorf("ScreeningTypeEnum.Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVIAOutcomeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    VIAOutcomeEnum
		want bool
	}{
		{
			name: "Valid type - Negative",
			e:    VIAOutcomeNegative,
			want: true,
		},
		{
			name: "Valid type - Positive",
			e:    VIAOutcomePositive,
			want: true,
		},
		{
			name: "Invalid type",
			e:    VIAOutcomeEnum("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("VIAOutcomeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVIAOutcomeEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    VIAOutcomeEnum
		want string
	}{
		{
			name: "Negative",
			e:    VIAOutcomeNegative,
			want: "negative",
		},
		{
			name: "Positive",
			e:    VIAOutcomePositive,
			want: "positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("VIAOutcomeEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVIAOutcomeEnum_UnmarshalGQL(t *testing.T) {
	value := VIAOutcomeNegative
	invalidType := VIAOutcomeEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *VIAOutcomeEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "negative",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalidType,
			args: args{
				v: "INVALID",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidType,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("VIAOutcomeEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVIAOutcomeEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     VIAOutcomeEnum
		wantW string
	}{
		{
			name:  "negative",
			e:     VIAOutcomeNegative,
			wantW: strconv.Quote("negative"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("VIAOutcomeEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestObservationStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ObservationStatusEnum
		want bool
	}{
		{
			name: "Valid type - Registered",
			e:    ObservationStatusEnumRegistered,
			want: true,
		},
		{
			name: "Valid type - Final",
			e:    ObservationStatusEnumFinal,
			want: true,
		},
		{
			name: "Invalid type",
			e:    ObservationStatusEnum("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ObservationStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObservationStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    ObservationStatusEnum
		want string
	}{
		{
			name: "Registered",
			e:    ObservationStatusEnumRegistered,
			want: "registered",
		},
		{
			name: "Final",
			e:    ObservationStatusEnumFinal,
			want: "final",
		},
		{
			name: "Entered-in-error",
			e:    ObservationStatusEnumEnteredInError,
			want: "entered-in-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ObservationStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObservationStatusEnum_UnmarshalGQL(t *testing.T) {
	value := ObservationStatusEnumRegistered
	invalidType := ObservationStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ObservationStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "registered",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalidType,
			args: args{
				v: "INVALID",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidType,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ObservationStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObservationStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ObservationStatusEnum
		wantW string
	}{
		{
			name:  "registered",
			e:     ObservationStatusEnumRegistered,
			wantW: strconv.Quote("registered"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ObservationStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    TaskStatus
		want bool
	}{
		{
			name: "Valid task status - completed",
			e:    CompletedTasksStatus,
			want: true,
		},
		{
			name: "Invalid task status status",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("TaskStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    TaskStatus
		want string
	}{
		{
			name: "completed",
			e:    CompletedTasksStatus,
			want: "completed",
		},
		{
			name: "requested",
			e:    RequestedTasksStatus,
			want: "requested",
		},
		{
			name: "cancelled",
			e:    CancelledTasksStatus,
			want: "cancelled",
		},
		{
			name: "accepted",
			e:    AcceptedTasksStatus,
			want: "accepted",
		},
		{
			name: "ready",
			e:    ReadyTasksStatus,
			want: "ready",
		},
		{
			name: "received",
			e:    ReceivedTasksStatus,
			want: "received",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("TaskStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_UnmarshalGQL(t *testing.T) {
	value := CompletedTasksStatus
	invalidType := TaskStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *TaskStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "completed",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalidType,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidType,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("TaskStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskStatus_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     TaskStatus
		wantW string
	}{
		{
			name:  "completed",
			e:     CompletedTasksStatus,
			wantW: strconv.Quote("completed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if got := w.String(); got != tt.wantW {
				t.Errorf("TaskStatus.MarshalGQL() = %v, want %v", got, tt.wantW)
			}
		})
	}
}

func TestServiceRequestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ServiceRequestStatus
		want bool
	}{
		{
			name: "Valid service request status - completed",
			e:    ServiceRequestStatusCompleted,
			want: true,
		},
		{
			name: "Invalid service request status status",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ServiceRequest.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceRequestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    ServiceRequestStatus
		want string
	}{
		{
			name: "draft",
			e:    ServiceRequestStatusDraft,
			want: "draft",
		},
		{
			name: "active",
			e:    ServiceRequestStatusActive,
			want: "active",
		},
		{
			name: "on-hold",
			e:    ServiceRequestStatusOnHold,
			want: "on-hold",
		},
		{
			name: "revoked",
			e:    ServiceRequestStatusRevoked,
			want: "revoked",
		},
		{
			name: "completed",
			e:    ServiceRequestStatusCompleted,
			want: "completed",
		},
		{
			name: "entered-in-error",
			e:    ServiceRequestStatusEnteredInError,
			want: "entered-in-error",
		},
		{
			name: "unknown",
			e:    ServiceRequestStatusUnknown,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ServiceRequest.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceRequestStatus_UnmarshalGQL(t *testing.T) {
	value := ServiceRequestStatusCompleted
	invalidType := ServiceRequestStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ServiceRequestStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "completed",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalidType,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidType,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ServiceRequestStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceRequestStatus_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ServiceRequestStatus
		wantW string
	}{
		{
			name:  "completed",
			e:     ServiceRequestStatusCompleted,
			wantW: strconv.Quote("completed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if got := w.String(); got != tt.wantW {
				t.Errorf("ServiceRequestStatus.MarshalGQL() = %v, want %v", got, tt.wantW)
			}
		})
	}
}

func TestDayOfWeek_String(t *testing.T) {
	tests := []struct {
		name string
		e    DayOfWeek
		want string
	}{
		{
			name: "MONDAY",
			e:    DayOfWeekMonday,
			want: "Monday",
		},
		{
			name: "TUESDAY",
			e:    DayOfWeekTuesday,
			want: "Tuesday",
		},
		{
			name: "WEDNESDAY",
			e:    DayOfWeekWednesday,
			want: "Wednesday",
		},
		{
			name: "Thursday",
			e:    DayOfWeekThursday,
			want: "Thursday",
		},
		{
			name: "FRIDAY",
			e:    DayOfWeekFriday,
			want: "Friday",
		},
		{
			name: "SATURDAY",
			e:    DayOfWeekSaturday,
			want: "Saturday",
		},
		{
			name: "SUNDAY",
			e:    DayOfWeekSunday,
			want: "Sunday",
		},
		{
			name: "Invalid Day",
			e:    "Invalid Day",
			want: "<unknown>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Display(); got != tt.want {
				t.Errorf("DayOfWeek.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayOfWeek_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    DayOfWeek
		want bool
	}{
		{
			name: "valid type",
			e:    DayOfWeekFriday,
			want: true,
		},
		{
			name: "invalid type",
			e:    DayOfWeek("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("DayOfWeek.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayOfWeek_UnmarshalGQL(t *testing.T) {
	value := DayOfWeekSunday
	invalid := DayOfWeek("invalid")

	type args struct {
		v interface{}
	}

	tests := []struct {
		name    string
		e       *DayOfWeek
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "SUNDAY",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DayOfWeek.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDayOfWeek_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}

	tests := []struct {
		name  string
		e     DayOfWeek
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     DayOfWeekSunday,
			b:     w,
			wantW: strconv.Quote("Sunday"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DayOfWeek.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestTerminologies_String(t *testing.T) {
	tests := []struct {
		name string
		e    Terminologies
		want string
	}{
		{
			name: "CIEL",
			e:    TerminologiesCIEL,
			want: "CIEL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Terminologies.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTerminologies_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    Terminologies
		want bool
	}{
		{
			name: "valid type",
			e:    TerminologiesCIEL,
			want: true,
		},
		{
			name: "invalid type",
			e:    Terminologies("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Terminologies.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTerminologies_UnmarshalGQL(t *testing.T) {
	value := TerminologiesCIEL
	invalid := Terminologies("invalid")

	type args struct {
		v interface{}
	}

	tests := []struct {
		name    string
		e       *Terminologies
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "CIEL",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Terminologies.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTerminologies_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}

	tests := []struct {
		name  string
		e     Terminologies
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     TerminologiesCIEL,
			b:     w,
			wantW: strconv.Quote("CIEL"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Terminologies.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFacilityIdentifierType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    FacilityIdentifierType
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    FacilityIdentifierTypeMFLCode,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    FacilityIdentifierType("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("FacilityIdentifierType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFacilityIdentifierType_String(t *testing.T) {
	tests := []struct {
		name string
		c    FacilityIdentifierType
		want string
	}{
		{
			name: "Happy Case",
			c:    FacilityIdentifierTypeMFLCode,
			want: FacilityIdentifierTypeMFLCode.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("FacilityIdentifierType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFacilityIdentifierType_UnmarshalGQL(t *testing.T) {
	validValue := FacilityIdentifierTypeMFLCode
	invalidType := FacilityIdentifierType("INVALID")

	type args struct {
		v interface{}
	}

	tests := []struct {
		name    string
		c       *FacilityIdentifierType
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: FacilityIdentifierTypeMFLCode.String(),
			},
			c:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid type",
			args: args{
				v: "invalid type",
			},
			c:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid type(float)",
			args: args{
				v: 45.1,
			},
			c:       &validValue,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("FacilityIdentifierType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFacilityIdentifierType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     FacilityIdentifierType
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     FacilityIdentifierTypeMFLCode,
			wantW: strconv.Quote("MFL_CODE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FacilityIdentifierType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestCountry_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    Country
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    CountryKenya,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    Country("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("Country.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountry_String(t *testing.T) {
	tests := []struct {
		name string
		c    Country
		want string
	}{
		{
			name: "Happy Case",
			c:    CountryKenya,
			want: CountryKenya.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Country.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountry_UnmarshalGQL(t *testing.T) {
	validValue := CountryKenya
	invalidType := Country("INVALID")

	type args struct {
		v interface{}
	}

	tests := []struct {
		name    string
		c       *Country
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: CountryKenya.String(),
			},
			c:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid type",
			args: args{
				v: "invalid type",
			},
			c:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid type(float)",
			args: args{
				v: 45.1,
			},
			c:       &validValue,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Country.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCountry_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     Country
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     CountryKenya,
			wantW: strconv.Quote("KE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Country.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMedicationRequestPriority_String(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationRequestPriority
		want string
	}{
		{
			name: "stat",
			e:    MedicationRequestPriorityStat,
			want: "stat",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MedicationRequestPriority.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationRequestPriority_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationRequestPriority
		want bool
	}{
		{
			name: "valid type",
			e:    MedicationRequestPriorityRoutine,
			want: true,
		},
		{
			name: "invalid type",
			e:    MedicationRequestPriority("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MedicationRequestPriority.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationRequestPriority_UnmarshalGQL(t *testing.T) {
	value := MedicationRequestPriorityStat
	invalid := MedicationRequestPriority("invalid")

	type args struct {
		v interface{}
	}

	tests := []struct {
		name    string
		e       *MedicationRequestPriority
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "stat",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("MedicationRequestPriority.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMedicationRequestPriority_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}

	tests := []struct {
		name  string
		e     MedicationRequestPriority
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     MedicationRequestPriorityRoutine,
			b:     w,
			wantW: strconv.Quote("routine"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MedicationRequestPriority.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestObservationCategory_String(t *testing.T) {
	tests := []struct {
		name string
		e    ObservationCategory
		want string
	}{
		{
			name: "laboratory",
			e:    ObservationCategoryLab,
			want: "laboratory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ObservationCategory.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObservationCategory_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ObservationCategory
		want bool
	}{
		{
			name: "valid type",
			e:    ObservationCategoryExam,
			want: true,
		},
		{
			name: "invalid type",
			e:    ObservationCategory("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ObservationCategory.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObservationCategory_UnmarshalGQL(t *testing.T) {
	value := ObservationCategoryLab
	invalid := ObservationCategory("invalid")

	type args struct {
		v interface{}
	}

	tests := []struct {
		name    string
		e       *ObservationCategory
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "laboratory",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ObservationCategory.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObservationCategory_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}

	tests := []struct {
		name  string
		e     ObservationCategory
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     ObservationCategoryExam,
			b:     w,
			wantW: strconv.Quote("exam"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ObservationCategory.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestObservationCategory_Text(t *testing.T) {
	tests := []struct {
		name                string
		observationCategory ObservationCategory
		want                string
	}{
		{
			name:                "Valid type - vital_signs",
			observationCategory: ObservationCategory(ObservationCategoryVitalSigns),
			want:                "vital-signs",
		},
		{
			name:                "Valid type - exam",
			observationCategory: ObservationCategory(ObservationCategoryExam),
			want:                "exam",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.observationCategory.Text(); got != tt.want {
				t.Errorf("ObservationCategory.Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionSeverity_ToString(t *testing.T) {
	tests := []struct {
		name string
		cs   ConditionSeverity
		want string
	}{
		{
			name: "mild",
			cs:   ConditionSeverityMild,
			want: "Mild",
		},
		{
			name: "moderate",
			cs:   ConditionSeverityModerate,
			want: "Moderate",
		},
		{
			name: "severe",
			cs:   ConditionSeveritySevere,
			want: "Severe",
		},
		{
			name: "severee",
			cs:   "ConditionSeveritySevere",
			want: "Unknown Condition Severity: ConditionSeveritySevere",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.ToString(); got != tt.want {
				t.Errorf("ConditionSeverity.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionCategory_ToCategoryCode(t *testing.T) {
	tests := []struct {
		name string
		cs   ConditionCategory
		want string
	}{
		{
			name: "ENCOUNTER_DIAGNOSIS",
			cs:   ConditionCategoryDiagnosis,
			want: "encounter-diagnosis",
		},
		{
			name: "PROBLEM_LIST_ITEM",
			cs:   ConditionCategoryProblemList,
			want: "problem-list-item",
		},
		{
			name: "severee",
			cs:   "ConditionCategoryProblemList",
			want: "Unknown Condition Category: ConditionCategoryProblemList",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.ToCategoryCode(); got != tt.want {
				t.Errorf("ConditionCategory.ToCategoryCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionCategory_ToString(t *testing.T) {
	tests := []struct {
		name string
		cs   ConditionCategory
		want string
	}{
		{
			name: "ENCOUNTER_DIAGNOSIS",
			cs:   ConditionCategoryDiagnosis,
			want: "Encounter Diagnosis",
		},
		{
			name: "PROBLEM_LIST_ITEM",
			cs:   ConditionCategoryProblemList,
			want: "Problem List Item",
		},
		{
			name: "severee",
			cs:   "ConditionCategoryProblemList",
			want: "Unknown Condition Category: ConditionCategoryProblemList",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.ToString(); got != tt.want {
				t.Errorf("ConditionCategory.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpeciality_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		speciality SpecialityEnum
		want       bool
	}{
		{
			name:       "Happy case: Valid speciality",
			speciality: Radiologist,
			want:       Radiologist.IsValid(),
		},
		{
			name:       "Sad case: invalid speciality",
			speciality: "Mganga",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.speciality.IsValid(); got != tt.want {
				t.Errorf("peciality.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLabTestTypeEnum_Code(t *testing.T) {
	tests := []struct {
		name     string
		testType LabTestTypeEnum
		want     string
	}{
		{name: "Mammogram", testType: MammogramTest, want: common.MammogramTerminologyCode},
		{name: "Biopsy", testType: BiopsyTest, want: common.BiopsyTerminologySystem},
		{name: "MRI", testType: MRITest, want: common.MRITerminologySystem},
		{name: "Ultrasound", testType: UltrasoundTest, want: common.ChestUltrasoundTerminologySystem},
		{name: "CBE", testType: CBETest, want: common.BreastExaminationLOINCTerminologySystem},
		{name: "Pap smear", testType: PapSmearTest, want: common.PapSmearTerminologyCode},
		{name: "Whole blood", testType: WholeBloodTest, want: common.WholeBloodTerminologyCode},
		{name: "Prostatic serum antigen", testType: ProstaticSerumAntigenTest, want: common.ProstateCancerTerminologyCode},
		{name: "HPV PCR DNA", testType: HPV_PCR_DNA_Test, want: common.HPV_PCR_DNATerminologyCode},
		{name: "HPV oncoprotein", testType: HPV_ONCOPROTEIN_Test, want: common.HPV_OncoproteinTerminologyCode},
		{name: "IHC progesterone receptor", testType: IHCProgesteroneReceptorTest, want: common.IHCProgesteroneReceptorLOINCCode},
		{name: "IHC estrogen receptor", testType: IHCEstrogenReceptorTest, want: common.IHCEstrogenReceptorLOINCCode},
		{name: "IHC HER2", testType: IHCHER2ReceptorTest, want: common.HER2LOINCCode},
		{name: "IHC Ki67", testType: IHCKi67Test, want: common.Ki67LOINCCode},
		{name: "Unmapped test type returns empty", testType: "SOME_UNKNOWN_TEST", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.testType.Code(); got != tt.want {
				t.Errorf("LabTestTypeEnum.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLabTestTypeFromString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  LabTestTypeEnum
	}{
		{name: "Raw enum value", value: "IHC_HER2", want: IHCHER2ReceptorTest},
		{name: "Humanised display label", value: "IHC HER2", want: IHCHER2ReceptorTest},
		{name: "Display label with irregular casing", value: "ihc her2", want: IHCHER2ReceptorTest},
		{name: "Pap smear display (does not round-trip via casing)", value: "Pap Smear", want: PapSmearTest},
		{name: "Pap smear raw enum", value: "PAPSMEAR", want: PapSmearTest},
		{name: "Mammogram display", value: "Mammogram", want: MammogramTest},
		{name: "Unknown value returned verbatim", value: "SOME_UNKNOWN_TEST", want: LabTestTypeEnum("SOME_UNKNOWN_TEST")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LabTestTypeFromString(tt.value); got != tt.want {
				t.Errorf("LabTestTypeFromString(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestSpeciality_Code(t *testing.T) {
	tests := []struct {
		name string
		code SpecialityEnum
		want string
	}{
		{
			name: "Happy case: Valid code",
			code: Dentist,
			want: Dentist.Code(),
		},
		{
			name: "Happy case: Valid code",
			code: Anesthesiologist,
			want: Anesthesiologist.Code(),
		},
		{
			name: "Happy case: Valid code",
			code: Radiologist,
			want: Radiologist.Code(),
		},
		{
			name: "Happy case: Valid code",
			code: GeneralPractioner,
			want: GeneralPractioner.Code(),
		},
		{
			name: "Happy case: Valid code",
			code: Pediatrician,
			want: Pediatrician.Code(),
		},
		{
			name: "Happy case: Valid code",
			code: Surgeon,
			want: Surgeon.Code(),
		},
		{
			name: "Happy case: Valid code",
			code: Obstetrician,
			want: Obstetrician.Code(),
		},
		{
			name: "Sad case: Invalid code",
			code: "Mganga",
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Code(); got != tt.want {
				t.Errorf("Speciality.Code() = %v want %v", got, tt.want)
			}
		})
	}
}

func TestSpeciality_Display(t *testing.T) {
	tests := []struct {
		name string
		code SpecialityEnum
		want string
	}{
		{
			name: "Happy case: Valid code",
			code: Dentist,
			want: Dentist.Display(),
		},
		{
			name: "Happy case: Valid code",
			code: Anesthesiologist,
			want: Anesthesiologist.Display(),
		},
		{
			name: "Happy case: Valid code",
			code: Radiologist,
			want: Radiologist.Display(),
		},
		{
			name: "Happy case: Valid code",
			code: Pediatrician,
			want: Pediatrician.Display(),
		},
		{
			name: "Happy case: Valid code",
			code: Surgeon,
			want: Surgeon.Display(),
		},
		{
			name: "Happy case: Valid code",
			code: Obstetrician,
			want: Obstetrician.Display(),
		},
		{
			name: "Happy case: Valid code",
			code: GeneralPractioner,
			want: GeneralPractioner.Display(),
		},
		{
			name: "Sad case: Invalid code",
			code: "Mganga",
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Display(); got != tt.want {
				t.Errorf("Speciality.Display() = %v want %v", got, tt.want)
			}
		})
	}
}

func TestCompositionStatusEnum_ToCode(t *testing.T) {
	tests := []struct {
		name string
		cs   CompositionStatusEnum
		want string
	}{
		{
			name: "PRELIMINARY",
			cs:   "PRELIMINARY",
			want: "preliminary",
		},
		{
			name: "FINAL",
			cs:   "FINAL",
			want: "final",
		},
		{
			name: "AMENDED",
			cs:   "AMENDED",
			want: "amended",
		},
		{
			name: "ENTERED_IN_ERROR",
			cs:   "ENTERED_IN_ERROR",
			want: "entered-in-error",
		},
		{
			name: "UNKNOWN",
			cs:   "UNKNOWN",
			want: "Unknown CompositionStatusEnum: UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.ToCode(); got != tt.want {
				t.Errorf("CompositionStatusEnum.ToCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayOfWeek_Code(t *testing.T) {
	tests := []struct {
		name string
		m    DayOfWeek
		want string
	}{
		{
			name: "MONDAY",
			m:    DayOfWeekMonday,
			want: "mon",
		},
		{
			name: "TUESDAY",
			m:    DayOfWeekTuesday,
			want: "tue",
		},
		{
			name: "WEDNESDAY",
			m:    DayOfWeekWednesday,
			want: "wed",
		},
		{
			name: "THURSDAY",
			m:    DayOfWeekThursday,
			want: "thu",
		},
		{
			name: "FRIDAY",
			m:    DayOfWeekFriday,
			want: "fri",
		},
		{
			name: "SATURDAY",
			m:    DayOfWeekSaturday,
			want: "sat",
		},
		{
			name: "SUNDAY",
			m:    DayOfWeekSunday,
			want: "sun",
		},
		{
			name: "SUNDAYT",
			m:    "DayOfWeekSunday",
			want: "<unknown>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Code(); got != tt.want {
				t.Errorf("DayOfWeek.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayOfWeek_Display(t *testing.T) {
	tests := []struct {
		name string
		m    DayOfWeek
		want string
	}{
		{
			name: "MONDAY",
			m:    DayOfWeekMonday,
			want: "Monday",
		},
		{
			name: "TUESDAY",
			m:    DayOfWeekTuesday,
			want: "Tuesday",
		},
		{
			name: "WEDNESDAY",
			m:    DayOfWeekWednesday,
			want: "Wednesday",
		},
		{
			name: "THURSDAY",
			m:    DayOfWeekThursday,
			want: "Thursday",
		},
		{
			name: "FRIDAY",
			m:    DayOfWeekFriday,
			want: "Friday",
		},
		{
			name: "SATURDAY",
			m:    DayOfWeekSaturday,
			want: "Saturday",
		},
		{
			name: "SUNDAY",
			m:    DayOfWeekSunday,
			want: "Sunday",
		},
		{
			name: "SUNDAYy",
			m:    "DayOfWeekSunday",
			want: "<unknown>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Display(); got != tt.want {
				t.Errorf("DayOfWeek.Display() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPatientVitalSignsEnumEnum_IsValid(t *testing.T) {
	type args struct {
		input PatientVitalSignsEnum
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Happy case: valid patient vital sign",
			args: args{
				input: PatientMenstrualPeriod,
			},
			want: true,
		},
		{
			name: "Sad case: invalid patient vital sign",
			args: args{
				input: "MBI",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PatientVitalSignsEnum(tt.args.input).IsValid()
			if got != tt.want {
				t.Errorf("PatientVitalSignsEnum.IsValid() = %v, want = %v", got, tt.want)
			}
		})
	}
}
