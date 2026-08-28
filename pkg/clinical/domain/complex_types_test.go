package domain

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/go-cmp/cmp"
	"github.com/savannahghi/scalarutils"
)

func TestAddressTypeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AddressTypeEnum
		want bool
	}{
		{
			name: "valid postal address type",
			e:    AddressTypeEnumPostal,
			want: true,
		},
		{
			name: "valid physical address type",
			e:    AddressTypeEnumPhysical,
			want: true,
		},
		{
			name: "valid both address type",
			e:    AddressTypeEnumBoth,
			want: true,
		},
		{
			name: "invalid address type",
			e:    AddressTypeEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AddressTypeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressTypeEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    AddressTypeEnum
		want string
	}{
		{
			name: "postal address type",
			e:    AddressTypeEnumPostal,
			want: "postal",
		},
		{
			name: "physical address type",
			e:    AddressTypeEnumPhysical,
			want: "physical",
		},
		{
			name: "both address type",
			e:    AddressTypeEnumBoth,
			want: "both",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AddressTypeEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressTypeEnum_UnmarshalGQL(t *testing.T) {
	value := AddressTypeEnumPhysical
	invalidType := AddressTypeEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AddressTypeEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "physical",
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
				t.Errorf("AddressTypeEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddressTypeEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     AddressTypeEnum
		wantW string
	}{
		{
			name:  "postal address type",
			e:     AddressTypeEnumPostal,
			wantW: strconv.Quote("postal"),
		},
		{
			name:  "physical address type",
			e:     AddressTypeEnumPhysical,
			wantW: strconv.Quote("physical"),
		},
		{
			name:  "both address type",
			e:     AddressTypeEnumBoth,
			wantW: strconv.Quote("both"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AddressTypeEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestAddressUseEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AddressUseEnum
		want bool
	}{
		{
			name: "valid home address use",
			e:    AddressUseEnumHome,
			want: true,
		},
		{
			name: "valid work address use",
			e:    AddressUseEnumWork,
			want: true,
		},
		{
			name: "valid temporary address use",
			e:    AddressUseEnumTemp,
			want: true,
		},
		{
			name: "valid old address use",
			e:    AddressUseEnumOld,
			want: true,
		},
		{
			name: "valid billing address use",
			e:    AddressUseEnumBilling,
			want: true,
		},
		{
			name: "invalid address use",
			e:    AddressUseEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AddressUseEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressUseEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    AddressUseEnum
		want string
	}{
		{
			name: "home address use",
			e:    AddressUseEnumHome,
			want: "home",
		},
		{
			name: "work address use",
			e:    AddressUseEnumWork,
			want: "work",
		},
		{
			name: "temp address use",
			e:    AddressUseEnumTemp,
			want: "temp",
		},
		{
			name: "old address use",
			e:    AddressUseEnumOld,
			want: "old",
		},
		{
			name: "billing address use",
			e:    AddressUseEnumBilling,
			want: "billing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AddressUseEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressUseEnum_UnmarshalGQL(t *testing.T) {
	value := AddressUseEnumWork
	invalidType := AddressUseEnum("invalid")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AddressUseEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "work",
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
				t.Errorf("AddressUseEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddressUseEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     AddressUseEnum
		wantW string
	}{
		{
			name:  "home address use",
			e:     AddressUseEnumHome,
			wantW: strconv.Quote("home"),
		},
		{
			name:  "temporary address use",
			e:     AddressUseEnumTemp,
			wantW: strconv.Quote("temp"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AddressUseEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestAgeComparatorEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AgeComparatorEnum
		want bool
	}{
		{
			name: "valid less than",
			e:    AgeComparatorEnumLessThan,
			want: true,
		},
		{
			name: "valid less than or equal to",
			e:    AgeComparatorEnumLessThanOrEqualTo,
			want: true,
		},
		{
			name: "valid greater than or equal to",
			e:    AgeComparatorEnumGreaterThanOrEqualTo,
			want: true,
		},
		{
			name: "valid greater than",
			e:    AgeComparatorEnumGreaterThan,
			want: true,
		},
		{
			name: "invalid",
			e:    AgeComparatorEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AgeComparatorEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgeComparatorEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    AgeComparatorEnum
		want string
	}{
		{
			name: "less than",
			e:    AgeComparatorEnumLessThan,
			want: "<",
		},
		{
			name: "less than or equal to",
			e:    AgeComparatorEnumLessThanOrEqualTo,
			want: "<=",
		},
		{
			name: "greater than or equal to",
			e:    AgeComparatorEnumGreaterThanOrEqualTo,
			want: ">=",
		},
		{
			name: "greater than",
			e:    AgeComparatorEnumGreaterThan,
			want: ">",
		},
		{
			name: "invalid",
			e:    AgeComparatorEnum("invalid"),
			want: "invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AgeComparatorEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgeComparatorEnum_UnmarshalGQL(t *testing.T) {
	value := AgeComparatorEnumLessThan
	invalidType := AgeComparatorEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AgeComparatorEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "less_than",
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
				t.Errorf("AgeComparatorEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgeComparatorEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     AgeComparatorEnum
		wantW string
	}{
		{
			name:  "less than",
			e:     AgeComparatorEnumLessThan,
			wantW: strconv.Quote("<"),
		},
		{
			name:  "less than or equal to",
			e:     AgeComparatorEnumLessThanOrEqualTo,
			wantW: strconv.Quote("<="),
		},
		{
			name:  "greater than or equal to",
			e:     AgeComparatorEnumGreaterThanOrEqualTo,
			wantW: strconv.Quote(">="),
		},
		{
			name:  "greater than",
			e:     AgeComparatorEnumGreaterThan,
			wantW: strconv.Quote(">"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AgeComparatorEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestContactPointUseEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ContactPointUseEnum
		want bool
	}{
		{
			name: "Happy Case - valid home contact point",
			e:    ContactPointUseEnumHome,
			want: true,
		},
		{
			name: "Happy Case - valid work contact point",
			e:    ContactPointUseEnumWork,
			want: true,
		},
		{
			name: "Happy Case - valid temp contact point",
			e:    ContactPointUseEnumTemp,
			want: true,
		},
		{
			name: "Happy Case - valid old contact point",
			e:    ContactPointUseEnumOld,
			want: true,
		},
		{
			name: "Happy Case - valid mobile contact point",
			e:    ContactPointUseEnumMobile,
			want: true,
		},
		{
			name: "Invalid contact point",
			e:    ContactPointUseEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ContactPointUseEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactPointSystemEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ContactPointSystemEnum
		want bool
	}{
		{
			name: "Happy Case - valid email contact point",
			e:    ContactPointSystemEnumEmail,
			want: true,
		},
		{
			name: "Happy Case - valid fax contact point",
			e:    ContactPointSystemEnumFax,
			want: true,
		},
		{
			name: "Happy Case - valid other contact point",
			e:    ContactPointSystemEnumOther,
			want: true,
		},
		{
			name: "Happy Case - valid pager contact point",
			e:    ContactPointSystemEnumPager,
			want: true,
		},
		{
			name: "Happy Case - valid phone contact point",
			e:    ContactPointSystemEnumPhone,
			want: true,
		},
		{
			name: "Invalid System contact point",
			e:    ContactPointSystemEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ContactPointSystemEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactPointSystemEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    ContactPointSystemEnum
		want string
	}{
		{
			name: "phone",
			e:    ContactPointSystemEnumPhone,
			want: "phone",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ContactPointSystemEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactPointSystemEnum_UnmarshalGQL(t *testing.T) {
	value := ContactPointSystemEnumEmail
	invalidContact := ContactPointSystemEnum("invalid")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ContactPointSystemEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid contact system",
			e:    &value,
			args: args{
				v: "email",
			},
			wantErr: false,
		},
		{
			name: "invalid contact",
			e:    &invalidContact,
			args: args{
				v: "this is not a valid contact",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidContact,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ContactPointSystemEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContactPointSystemEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ContactPointSystemEnum
		wantW string
	}{
		{
			name:  "email",
			e:     ContactPointSystemEnumEmail,
			wantW: strconv.Quote("email"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ContactPointSystemEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestContactPointUseEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    ContactPointUseEnum
		want string
	}{
		{
			name: "Home",
			e:    ContactPointUseEnumHome,
			want: "home",
		},
		{
			name: "Work",
			e:    ContactPointUseEnumWork,
			want: "work",
		},
		{
			name: "Temp",
			e:    ContactPointUseEnumTemp,
			want: "temp",
		},
		{
			name: "Old",
			e:    ContactPointUseEnumOld,
			want: "old",
		},
		{
			name: "Mobile",
			e:    ContactPointUseEnumMobile,
			want: "mobile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ContactPointUseEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactPointUseEnum_UnmarshalGQL(t *testing.T) {
	value := ContactPointUseEnumHome
	invalidType := ContactPointUseEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ContactPointUseEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "home",
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
				t.Errorf("ContactPointUseEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContactPointUseEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ContactPointUseEnum
		wantW string
	}{
		{
			name:  "Home contact point",
			e:     ContactPointUseEnumHome,
			wantW: strconv.Quote("home"),
		},
		{
			name:  "Work contact point",
			e:     ContactPointUseEnumWork,
			wantW: strconv.Quote("work"),
		},
		{
			name:  "Temp contact point",
			e:     ContactPointUseEnumTemp,
			wantW: strconv.Quote("temp"),
		},
		{
			name:  "Old contact point",
			e:     ContactPointUseEnumOld,
			wantW: strconv.Quote("old"),
		},
		{
			name:  "Mobile contact point",
			e:     ContactPointUseEnumMobile,
			wantW: strconv.Quote("mobile"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ContactPointUseEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestDurationComparatorEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    DurationComparatorEnum
		want bool
	}{
		{
			name: "valid less than",
			e:    DurationComparatorEnumLessThan,
			want: true,
		},
		{
			name: "valid less than or equal to",
			e:    DurationComparatorEnumLessThanOrEqualTo,
			want: true,
		},
		{
			name: "valid greater than or equal to",
			e:    DurationComparatorEnumGreaterThanOrEqualTo,
			want: true,
		},
		{
			name: "valid greater than",
			e:    DurationComparatorEnumGreaterThan,
			want: true,
		},
		{
			name: "invalid",
			e:    DurationComparatorEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("DurationComparatorEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDurationComparatorEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    DurationComparatorEnum
		want string
	}{
		{
			name: "less than",
			e:    DurationComparatorEnumLessThan,
			want: "<",
		},
		{
			name: "less than or equal to",
			e:    DurationComparatorEnumLessThanOrEqualTo,
			want: "<=",
		},
		{
			name: "greater than",
			e:    DurationComparatorEnumGreaterThan,
			want: ">",
		},
		{
			name: "greater than or equal to",
			e:    DurationComparatorEnumGreaterThanOrEqualTo,
			want: ">=",
		},
		{
			name: "invalid value",
			e:    DurationComparatorEnum("invalid"),
			want: "invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("DurationComparatorEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDurationComparatorEnum_UnmarshalGQL(t *testing.T) {
	value := DurationComparatorEnumLessThan
	invalidType := DurationComparatorEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *DurationComparatorEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "less_than",
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
				t.Errorf("DurationComparatorEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDurationComparatorEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     DurationComparatorEnum
		wantW string
	}{
		{
			name:  "less than",
			e:     DurationComparatorEnumLessThan,
			wantW: strconv.Quote("<"),
		},
		{
			name:  "less than or equal to",
			e:     DurationComparatorEnumLessThanOrEqualTo,
			wantW: strconv.Quote("<="),
		},
		{
			name:  "greater than",
			e:     DurationComparatorEnumGreaterThan,
			wantW: strconv.Quote(">"),
		},
		{
			name:  "greater than or equal to",
			e:     DurationComparatorEnumGreaterThanOrEqualTo,
			wantW: strconv.Quote(">="),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DurationComparatorEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestHumanNameUseEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    HumanNameUseEnum
		want bool
	}{
		{
			name: "valid enum value",
			e:    HumanNameUseEnumUsual,
			want: true,
		},
		{
			name: "invalid enum value",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("HumanNameUseEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanNameUseEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    HumanNameUseEnum
		want string
	}{
		{
			name: "usual",
			e:    HumanNameUseEnumUsual,
			want: "usual",
		},
		{
			name: "official",
			e:    HumanNameUseEnumOfficial,
			want: "official",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("HumanNameUseEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanNameUseEnum_UnmarshalGQL(t *testing.T) {
	validValue := HumanNameUseEnumUsual
	invalidValue := HumanNameUseEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *HumanNameUseEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid value",
			e:    &validValue,
			args: args{
				v: "usual",
			},
			wantErr: false,
		},
		{
			name: "invalid value",
			e:    &invalidValue,
			args: args{
				v: "invalid",
			},
			wantErr: true,
		},
		{
			name: "non-string value",
			e:    &invalidValue,
			args: args{
				v: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("HumanNameUseEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHumanNameUseEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     HumanNameUseEnum
		wantW string
	}{
		{
			name:  "anonymous",
			e:     HumanNameUseEnumAnonymous,
			wantW: strconv.Quote("anonymous"),
		},
		{
			name:  "maiden",
			e:     HumanNameUseEnumMaiden,
			wantW: `"maiden"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("HumanNameUseEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestIdentifierUseEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    IdentifierUseEnum
		want bool
	}{
		{
			name: "Happy Case - Usual use",
			e:    IdentifierUseEnumUsual,
			want: true,
		},
		{
			name: "Happy Case - Official use",
			e:    IdentifierUseEnumOfficial,
			want: true,
		},
		{
			name: "Happy Case - Temp use",
			e:    IdentifierUseEnumTemp,
			want: true,
		},
		{
			name: "Happy Case - Secondary use",
			e:    IdentifierUseEnumSecondary,
			want: true,
		},
		{
			name: "Happy Case - Old use",
			e:    IdentifierUseEnumOld,
			want: true,
		},
		{
			name: "Sad Case - Invalid use",
			e:    IdentifierUseEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IdentifierUseEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifierUseEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    IdentifierUseEnum
		want string
	}{
		{
			name: "Happy case",
			e:    IdentifierUseEnumOfficial,
			want: "official",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("IdentifierUseEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifierUseEnum_UnmarshalGQL(t *testing.T) {
	validValue := IdentifierUseEnumOfficial
	invalidValue := IdentifierUseEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *IdentifierUseEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid value",
			e:    &validValue,
			args: args{
				v: "usual",
			},
			wantErr: false,
		},
		{
			name: "invalid value",
			e:    &invalidValue,
			args: args{
				v: "invalid",
			},
			wantErr: true,
		},
		{
			name: "non-string value",
			e:    &invalidValue,
			args: args{
				v: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("IdentifierUseEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIdentifierUseEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     IdentifierUseEnum
		wantW string
	}{
		{
			name:  "official",
			e:     IdentifierUseEnumOfficial,
			wantW: strconv.Quote("official"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("IdentifierUseEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestNarrativeStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    NarrativeStatusEnum
		want bool
	}{
		{
			name: "Happy Case - Generated status",
			e:    NarrativeStatusEnumGenerated,
			want: true,
		},
		{
			name: "Happy Case - Extensions status",
			e:    NarrativeStatusEnumExtensions,
			want: true,
		},
		{
			name: "Happy Case - Additional status",
			e:    NarrativeStatusEnumAdditional,
			want: true,
		},
		{
			name: "Happy Case - Empty status",
			e:    NarrativeStatusEnumEmpty,
			want: true,
		},
		{
			name: "Sad Case - Invalid status",
			e:    NarrativeStatusEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("NarrativeStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNarrativeStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    NarrativeStatusEnum
		want string
	}{
		{
			name: "Generated",
			e:    NarrativeStatusEnumGenerated,
			want: "generated",
		},
		{
			name: "Extensions",
			e:    NarrativeStatusEnumExtensions,
			want: "extensions",
		},
		{
			name: "Additional",
			e:    NarrativeStatusEnumAdditional,
			want: "additional",
		},
		{
			name: "Empty",
			e:    NarrativeStatusEnumEmpty,
			want: "empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("NarrativeStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNarrativeStatusEnum_UnmarshalGQL(t *testing.T) {
	validType := NarrativeStatusEnumGenerated
	invalidType := NarrativeStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *NarrativeStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - valid type",
			e:    &validType,
			args: args{
				v: "generated",
			},
			wantErr: false,
		},
		{
			name: "invalid value",
			e:    &invalidType,
			args: args{
				v: "invalid",
			},
			wantErr: true,
		},
		{
			name: "non-string value",
			e:    &invalidType,
			args: args{
				v: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("NarrativeStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNarrativeStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     NarrativeStatusEnum
		wantW string
	}{
		{
			name:  "additional",
			e:     NarrativeStatusEnumAdditional,
			wantW: strconv.Quote("additional"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NarrativeStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestQuantityComparatorEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    QuantityComparatorEnum
		want bool
	}{
		{
			name: "Happy case - valid less_than comparator",
			e:    QuantityComparatorEnumLessThan,
			want: true,
		},
		{
			name: "Happy case - valid less_than_or_equal_to comparator",
			e:    QuantityComparatorEnumLessThanOrEqualTo,
			want: true,
		},
		{
			name: "Happy case - valid greater_than_or_equal_to comparator",
			e:    QuantityComparatorEnumGreaterThanOrEqualTo,
			want: true,
		},
		{
			name: "Happy case - valid greater_than comparator",
			e:    QuantityComparatorEnumGreaterThan,
			want: true,
		},
		{
			name: "Error case - invalid string",
			e:    QuantityComparatorEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("QuantityComparatorEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuantityComparatorEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    QuantityComparatorEnum
		want string
	}{
		{
			name: "less_than",
			e:    QuantityComparatorEnumLessThan,
			want: "<",
		},
		{
			name: "less_than_or_equal_to",
			e:    QuantityComparatorEnumLessThanOrEqualTo,
			want: "<=",
		},
		{
			name: "greater_than_or_equal_to",
			e:    QuantityComparatorEnumGreaterThanOrEqualTo,
			want: ">=",
		},
		{
			name: "greater_than",
			e:    QuantityComparatorEnumGreaterThan,
			want: ">",
		},
		{
			name: "invalid",
			e:    QuantityComparatorEnum("invalid"),
			want: "invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("QuantityComparatorEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuantityComparatorEnum_UnmarshalGQL(t *testing.T) {
	validType := QuantityComparatorEnumLessThan
	invalidType := QuantityComparatorEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *QuantityComparatorEnum
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - valid less_than comparator",
			e:    &validType,
			args: args{
				v: "less_than",
			},
			wantErr: false,
		},
		{
			name: "Happy case - valid greater_than_or_equal_to comparator",
			e:    &validType,
			args: args{
				v: "greater_than_or_equal_to",
			},
			wantErr: false,
		},
		{
			name: "invalid value",
			e:    &invalidType,
			args: args{
				v: "invalid",
			},
			wantErr: true,
		},
		{
			name: "non-string value",
			e:    &invalidType,
			args: args{
				v: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("QuantityComparatorEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuantityComparatorEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     QuantityComparatorEnum
		wantW string
	}{
		{
			name:  "less than",
			e:     QuantityComparatorEnumLessThan,
			wantW: strconv.Quote("<"),
		},
		{
			name:  "less than or equal to",
			e:     QuantityComparatorEnumLessThanOrEqualTo,
			wantW: strconv.Quote("<="),
		},
		{
			name:  "greater than or equal to",
			e:     QuantityComparatorEnumGreaterThanOrEqualTo,
			wantW: strconv.Quote(">="),
		},
		{
			name:  "greater than",
			e:     QuantityComparatorEnumGreaterThan,
			wantW: strconv.Quote(">"),
		},
		{
			name:  "invalid",
			e:     QuantityComparatorEnum("invalid"),
			wantW: strconv.Quote("invalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("QuantityComparatorEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestUnitsOfTimeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    UnitsOfTimeEnum
		want bool
	}{
		{
			name: "valid_s",
			e:    UnitsOfTimeEnumS,
			want: true,
		},
		{
			name: "valid_min",
			e:    UnitsOfTimeEnumMin,
			want: true,
		},
		{
			name: "valid_h",
			e:    UnitsOfTimeEnumH,
			want: true,
		},
		{
			name: "valid_d",
			e:    UnitsOfTimeEnumD,
			want: true,
		},
		{
			name: "valid_wk",
			e:    UnitsOfTimeEnumWk,
			want: true,
		},
		{
			name: "valid_mo",
			e:    UnitsOfTimeEnumMo,
			want: true,
		},
		{
			name: "valid_a",
			e:    UnitsOfTimeEnumA,
			want: true,
		},
		{
			name: "invalid",
			e:    UnitsOfTimeEnum("foo"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("UnitsOfTimeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitsOfTimeEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    UnitsOfTimeEnum
		want string
	}{
		{
			name: "s",
			e:    UnitsOfTimeEnumS,
			want: "s",
		},
		{
			name: "min",
			e:    UnitsOfTimeEnumMin,
			want: "min",
		},
		{
			name: "h",
			e:    UnitsOfTimeEnumH,
			want: "h",
		},
		{
			name: "d",
			e:    UnitsOfTimeEnumD,
			want: "d",
		},
		{
			name: "wk",
			e:    UnitsOfTimeEnumWk,
			want: "wk",
		},
		{
			name: "mo",
			e:    UnitsOfTimeEnumMo,
			want: "mo",
		},
		{
			name: "a",
			e:    UnitsOfTimeEnumA,
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("UnitsOfTimeEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitsOfTimeEnum_UnmarshalGQL(t *testing.T) {
	validType := UnitsOfTimeEnumD
	invalidType := UnitsOfTimeEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *UnitsOfTimeEnum
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - valid type",
			e:    &validType,
			args: args{
				v: "d",
			},
			wantErr: false,
		},
		{
			name: "invalid value",
			e:    &invalidType,
			args: args{
				v: "invalid",
			},
			wantErr: true,
		},
		{
			name: "non-string value",
			e:    &invalidType,
			args: args{
				v: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UnitsOfTimeEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnitsOfTimeEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     UnitsOfTimeEnum
		wantW string
	}{
		{
			name:  "Valid enum value",
			e:     UnitsOfTimeEnumA,
			wantW: strconv.Quote("a"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UnitsOfTimeEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestTimingRepeatWhenEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    TimingRepeatWhenEnum
		want bool
	}{
		{
			name: "Valid enum value",
			e:    TimingRepeatWhenEnumMorn,
			want: true,
		},
		{
			name: "Valid enum value",
			e:    TimingRepeatWhenEnumCv,
			want: true,
		},
		{
			name: "Invalid enum value",
			e:    TimingRepeatWhenEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("TimingRepeatWhenEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimingRepeatWhenEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    TimingRepeatWhenEnum
		want string
	}{
		{
			name: "MORN",
			e:    TimingRepeatWhenEnumMorn,
			want: "MORN",
		},
		{
			name: "MORN_early",
			e:    TimingRepeatWhenEnumMornEarly,
			want: "MORN_early",
		},
		{
			name: "MORN_late",
			e:    TimingRepeatWhenEnumMornLate,
			want: "MORN_late",
		},
		{
			name: "NOON",
			e:    TimingRepeatWhenEnumNoon,
			want: "NOON",
		},
		{
			name: "AFT",
			e:    TimingRepeatWhenEnumAft,
			want: "AFT",
		},
		{
			name: "AFT_early",
			e:    TimingRepeatWhenEnumAftEarly,
			want: "AFT_early",
		},
		{
			name: "AFT_late",
			e:    TimingRepeatWhenEnumAftLate,
			want: "AFT_late",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("TimingRepeatWhenEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimingRepeatWhenEnum_UnmarshalGQL(t *testing.T) {
	validType := TimingRepeatWhenEnumWake
	invalidType := TimingRepeatWhenEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *TimingRepeatWhenEnum
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - valid type",
			e:    &validType,
			args: args{
				v: "WAKE",
			},
			wantErr: false,
		},
		{
			name: "invalid value",
			e:    &invalidType,
			args: args{
				v: "invalid",
			},
			wantErr: true,
		},
		{
			name: "non-string value",
			e:    &invalidType,
			args: args{
				v: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("TimingRepeatWhenEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimingRepeatWhenEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     TimingRepeatWhenEnum
		wantW string
	}{
		{
			name:  "Valid enum value",
			e:     TimingRepeatWhenEnumAc,
			wantW: strconv.Quote("AC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("TimingRepeatWhenEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestAllergyIntoleranceCategoryEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AllergyIntoleranceCategoryEnum
		want bool
	}{
		{
			name: "Valid food",
			e:    AllergyIntoleranceCategoryEnumFood,
			want: true,
		},
		{
			name: "Valid medication",
			e:    AllergyIntoleranceCategoryEnumMedication,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AllergyIntoleranceCategoryEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllergyIntoleranceCategoryEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    AllergyIntoleranceCategoryEnum
		want string
	}{
		{
			name: "Happy Case - Return valid string",
			e:    AllergyIntoleranceCategoryEnumFood,
			want: "food",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AllergyIntoleranceCategoryEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllergyIntoleranceCategoryEnum_UnmarshalGQL(t *testing.T) {
	value := AllergyIntoleranceCategoryEnumEnvironment
	invalidType := AllergyIntoleranceCategoryEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AllergyIntoleranceCategoryEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "environment",
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
				t.Errorf("AllergyIntoleranceCategoryEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAllergyIntoleranceCategoryEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     AllergyIntoleranceCategoryEnum
		wantW string
	}{
		{
			name:  "biologic",
			e:     AllergyIntoleranceCategoryEnumBiologic,
			wantW: strconv.Quote("biologic"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AllergyIntoleranceCategoryEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestAllergyIntoleranceCriticalityEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AllergyIntoleranceCriticalityEnum
		want bool
	}{
		{
			name: "Valid low",
			e:    AllergyIntoleranceCriticalityEnumLow,
			want: true,
		},
		{
			name: "Valid high",
			e:    AllergyIntoleranceCriticalityEnumHigh,
			want: true,
		},
		{
			name: "Invalid",
			e:    AllergyIntoleranceCriticalityEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AllergyIntoleranceCriticalityEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllergyIntoleranceCriticalityEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    AllergyIntoleranceCriticalityEnum
		want string
	}{
		{
			name: "Happy Case - Return valid string",
			e:    AllergyIntoleranceCriticalityEnumLow,
			want: "low",
		},
		{
			name: "Happy Case - Return valid string",
			e:    AllergyIntoleranceCriticalityEnumUnableToAssess,
			want: "unable-to-assess",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AllergyIntoleranceCriticalityEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllergyIntoleranceCriticalityEnum_UnmarshalGQL(t *testing.T) {
	value := AllergyIntoleranceCriticalityEnumHigh
	invalidType := AllergyIntoleranceCriticalityEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AllergyIntoleranceCriticalityEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "high",
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
				t.Errorf("AllergyIntoleranceCriticalityEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAllergyIntoleranceCriticalityEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     AllergyIntoleranceCriticalityEnum
		wantW string
	}{
		{
			name:  "high",
			e:     AllergyIntoleranceCriticalityEnumHigh,
			wantW: strconv.Quote("high"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AllergyIntoleranceCriticalityEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestAllergyIntoleranceTypeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AllergyIntoleranceTypeEnum
		want bool
	}{
		{
			name: "Valid allergy",
			e:    AllergyIntoleranceTypeEnumAllergy,
			want: true,
		},
		{
			name: "Valid intolerance",
			e:    AllergyIntoleranceTypeEnumIntolerance,
			want: true,
		},
		{
			name: "Invalid type",
			e:    AllergyIntoleranceTypeEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AllergyIntoleranceTypeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllergyIntoleranceTypeEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    AllergyIntoleranceTypeEnum
		want string
	}{
		{
			name: "Happy Case - Return valid string",
			e:    AllergyIntoleranceTypeEnumAllergy,
			want: "allergy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AllergyIntoleranceTypeEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllergyIntoleranceTypeEnum_UnmarshalGQL(t *testing.T) {
	value := AllergyIntoleranceTypeEnumIntolerance
	invalidType := AllergyIntoleranceTypeEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AllergyIntoleranceTypeEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "intolerance",
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
				t.Errorf("AllergyIntoleranceTypeEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAllergyIntoleranceTypeEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     AllergyIntoleranceTypeEnum
		wantW string
	}{
		{
			name:  "allergy",
			e:     AllergyIntoleranceTypeEnumAllergy,
			wantW: strconv.Quote("allergy"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AllergyIntoleranceTypeEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestAllergyIntoleranceReactionSeverityEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AllergyIntoleranceReactionSeverityEnum
		want bool
	}{
		{
			name: "valid_mild",
			e:    AllergyIntoleranceReactionSeverityEnumMild,
			want: true,
		},
		{
			name: "valid_moderate",
			e:    AllergyIntoleranceReactionSeverityEnumModerate,
			want: true,
		},
		{
			name: "valid_severe",
			e:    AllergyIntoleranceReactionSeverityEnumSevere,
			want: true,
		},
		{
			name: "invalid_enum_value",
			e:    AllergyIntoleranceReactionSeverityEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AllergyIntoleranceReactionSeverityEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllergyIntoleranceReactionSeverityEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    AllergyIntoleranceReactionSeverityEnum
		want string
	}{
		{
			name: "mild",
			e:    AllergyIntoleranceReactionSeverityEnumMild,
			want: "mild",
		},
		{
			name: "moderate",
			e:    AllergyIntoleranceReactionSeverityEnumModerate,
			want: "moderate",
		},
		{
			name: "severe",
			e:    AllergyIntoleranceReactionSeverityEnumSevere,
			want: "severe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AllergyIntoleranceReactionSeverityEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllergyIntoleranceReactionSeverityEnum_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AllergyIntoleranceReactionSeverityEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid_mild",
			e:    new(AllergyIntoleranceReactionSeverityEnum),
			args: args{
				v: "mild",
			},
			wantErr: false,
		},
		{
			name: "invalid_enum_value",
			e:    new(AllergyIntoleranceReactionSeverityEnum),
			args: args{
				v: "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid_type",
			e:    new(AllergyIntoleranceReactionSeverityEnum),
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("AllergyIntoleranceReactionSeverityEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAllergyIntoleranceReactionSeverityEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     AllergyIntoleranceReactionSeverityEnum
		wantW string
	}{
		{
			name:  "mild",
			e:     AllergyIntoleranceReactionSeverityEnumMild,
			wantW: strconv.Quote("mild"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AllergyIntoleranceReactionSeverityEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestCompositionStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    CompositionStatusEnum
		want bool
	}{
		{
			name: "valid preliminary status",
			e:    CompositionStatusEnumPreliminary,
			want: true,
		},
		{
			name: "valid final status",
			e:    CompositionStatusEnumFinal,
			want: true,
		},
		{
			name: "valid amended status",
			e:    CompositionStatusEnumAmended,
			want: true,
		},
		{
			name: "valid entered_in_error status",
			e:    CompositionStatusEnumEnteredInError,
			want: true,
		},
		{
			name: "invalid status",
			e:    CompositionStatusEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("CompositionStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompositionStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    CompositionStatusEnum
		want string
	}{
		{
			name: "preliminary status",
			e:    CompositionStatusEnumPreliminary,
			want: "preliminary",
		},
		{
			name: "final status",
			e:    CompositionStatusEnumFinal,
			want: "final",
		},
		{
			name: "amended status",
			e:    CompositionStatusEnumAmended,
			want: "amended",
		},
		{
			name: "entered_in_error status",
			e:    CompositionStatusEnumEnteredInError,
			want: "entered-in-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("CompositionStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompositionStatusEnum_UnmarshalGQL(t *testing.T) {
	value := CompositionStatusEnumFinal
	invalidType := CompositionStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *CompositionStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "final",
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
				t.Errorf("CompositionStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompositionStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     CompositionStatusEnum
		wantW string
	}{
		{
			name:  "amended",
			e:     CompositionStatusEnumAmended,
			wantW: strconv.Quote("amended"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("CompositionStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestCompositionAttesterModeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    CompositionAttesterModeEnum
		want bool
	}{
		{
			name: "valid mode",
			e:    CompositionAttesterModeEnumPersonal,
			want: true,
		},
		{
			name: "invalid mode",
			e:    CompositionAttesterModeEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("CompositionAttesterModeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompositionAttesterModeEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    CompositionAttesterModeEnum
		want string
	}{
		{
			name: "valid mode",
			e:    CompositionAttesterModeEnumPersonal,
			want: "personal",
		},
		{
			name: "invalid mode",
			e:    CompositionAttesterModeEnum("invalid"),
			want: "invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("CompositionAttesterModeEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompositionAttesterModeEnum_UnmarshalGQL(t *testing.T) {
	value := CompositionAttesterModeEnumOfficial
	invalidMode := CompositionAttesterModeEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *CompositionAttesterModeEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid mode",
			e:    &value,
			args: args{
				v: "official",
			},
			wantErr: false,
		},
		{
			name: "invalid mode",
			e:    &invalidMode,
			args: args{
				v: "this is not a valid mode",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidMode,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("CompositionAttesterModeEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompositionAttesterModeEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     CompositionAttesterModeEnum
		wantW string
	}{
		{
			name:  "legal",
			e:     CompositionAttesterModeEnumLegal,
			wantW: strconv.Quote("legal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("CompositionAttesterModeEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestEncounterStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    EncounterStatusEnum
		want bool
	}{
		{
			name: "valid status",
			e:    EncounterStatusEnumPlanned,
			want: true,
		},
		{
			name: "invalid status",
			e:    EncounterStatusEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("EncounterStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncounterStatusEnum_IsFinal(t *testing.T) {
	tests := []struct {
		name string
		e    EncounterStatusEnum
		want bool
	}{
		{
			name: "valid status",
			e:    EncounterStatusEnumCompleted,
			want: true,
		},
		{
			name: "invalid status",
			e:    EncounterStatusEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsFinal(); got != tt.want {
				t.Errorf("EncounterStatusEnum.IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncounterStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    EncounterStatusEnum
		want string
	}{
		{
			name: "in progress status",
			e:    EncounterStatusEnumInProgress,
			want: "in-progress",
		},
		{
			name: "entered in error status",
			e:    EncounterStatusEnumEnteredInError,
			want: "entered-in-error",
		},
		{
			name: "other status",
			e:    EncounterStatusEnumDiscontinued,
			want: "discontinued",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("EncounterStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncounterStatusEnum_UnmarshalGQL(t *testing.T) {
	value := EncounterStatusEnumDiscontinued
	invalidStatus := EncounterStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *EncounterStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid status",
			e:    &value,
			args: args{
				v: "discontinued",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			e:    &invalidStatus,
			args: args{
				v: "this is not a valid status",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidStatus,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("EncounterStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncounterStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     EncounterStatusEnum
		wantW string
	}{
		{
			name:  "in progress",
			e:     EncounterStatusEnumInProgress,
			wantW: strconv.Quote("in-progress"),
		},
		{
			name:  "unknown",
			e:     EncounterStatusEnumUnknown,
			wantW: strconv.Quote("unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("EncounterStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestEncounterLocationStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    EncounterLocationStatusEnum
		want bool
	}{
		{
			name: "valid status",
			e:    EncounterLocationStatusEnumPlanned,
			want: true,
		},
		{
			name: "invalid status",
			e:    EncounterLocationStatusEnum("invalid"),
			want: false,
		},
		{
			name: "reserved status",
			e:    EncounterLocationStatusEnumReserved,
			want: true,
		},
		{
			name: "active status",
			e:    EncounterLocationStatusEnumActive,
			want: true,
		},
		{
			name: "completed status",
			e:    EncounterLocationStatusEnumCompleted,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("EncounterLocationStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncounterLocationStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    EncounterLocationStatusEnum
		want string
	}{
		{
			name: "planned status",
			e:    EncounterLocationStatusEnumPlanned,
			want: "planned",
		},
		{
			name: "reserved status",
			e:    EncounterLocationStatusEnumReserved,
			want: "reserved",
		},
		{
			name: "active status",
			e:    EncounterLocationStatusEnumActive,
			want: "active",
		},
		{
			name: "completed status",
			e:    EncounterLocationStatusEnumCompleted,
			want: "completed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("EncounterLocationStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncounterLocationStatusEnum_UnmarshalGQL(t *testing.T) {
	value := EncounterLocationStatusEnumPlanned
	invalidStatus := EncounterLocationStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *EncounterLocationStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid status",
			e:    &value,
			args: args{
				v: "planned",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			e:    &invalidStatus,
			args: args{
				v: "this is not a valid status",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidStatus,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("EncounterLocationStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncounterLocationStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     EncounterLocationStatusEnum
		wantW string
	}{
		{
			name:  "completed",
			e:     EncounterLocationStatusEnumCompleted,
			wantW: strconv.Quote("completed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("EncounterLocationStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestEncounterStatusHistoryStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    EncounterStatusHistoryStatusEnum
		want bool
	}{
		{
			name: "valid status",
			e:    EncounterStatusHistoryStatusEnumPlanned,
			want: true,
		},
		{
			name: "invalid status",
			e:    EncounterStatusHistoryStatusEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("EncounterStatusHistoryStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncounterStatusHistoryStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    EncounterStatusHistoryStatusEnum
		want string
	}{
		{
			name: "valid status",
			e:    EncounterStatusHistoryStatusEnumPlanned,
			want: "planned",
		},
		{
			name: "valid status",
			e:    EncounterStatusHistoryStatusEnumInProgress,
			want: "in-progress",
		},
		{
			name: "valid status",
			e:    EncounterStatusHistoryStatusEnumEnteredInError,
			want: "entered-in-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("EncounterStatusHistoryStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncounterStatusHistoryStatusEnum_UnmarshalGQL(t *testing.T) {
	value := EncounterStatusHistoryStatusEnumPlanned
	invalidStatus := EncounterStatusHistoryStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *EncounterStatusHistoryStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid status",
			e:    &value,
			args: args{
				v: "planned",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			e:    &invalidStatus,
			args: args{
				v: "this is not a valid status",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidStatus,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("EncounterStatusHistoryStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncounterStatusHistoryStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     EncounterStatusHistoryStatusEnum
		wantW string
	}{
		{
			name:  "planned",
			e:     EncounterStatusHistoryStatusEnumPlanned,
			wantW: strconv.Quote("planned"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("EncounterStatusHistoryStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestEpisodeOfCareEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    EpisodeOfCareStatusEnum
		want bool
	}{
		{
			name: "valid status",
			e:    EpisodeOfCareStatusEnumPlanned,
			want: true,
		},
		{
			name: "invalid status",
			e:    EpisodeOfCareStatusEnum("invalid"),
			want: false,
		},
		{
			name: "waitlist status",
			e:    EpisodeOfCareStatusEnumWaitlist,
			want: true,
		},
		{
			name: "active status",
			e:    EpisodeOfCareStatusEnumActive,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("EpisodeOfCareStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpisodeOfCareEnum_IsFinal(t *testing.T) {
	tests := []struct {
		name string
		e    EpisodeOfCareStatusEnum
		want bool
	}{
		{
			name: "finished",
			e:    EpisodeOfCareStatusEnumFinished,
			want: true,
		},
		{
			name: "invalid status",
			e:    EpisodeOfCareStatusEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsFinal(); got != tt.want {
				t.Errorf("EpisodeOfCareStatusEnum.isFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpisodeOfCareEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    EpisodeOfCareStatusEnum
		want string
	}{
		{
			name: "entered-in-error",
			e:    EpisodeOfCareStatusEnumEnteredInError,
			want: "entered-in-error",
		},
		{
			name: "active",
			e:    EpisodeOfCareStatusEnumActive,
			want: "active",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("EpisodeOfCareEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpisodeOfCareStatusEnum_UnmarshalGQL(t *testing.T) {
	value := EpisodeOfCareStatusEnumActive
	invalidStatus := EpisodeOfCareStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *EpisodeOfCareStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid status",
			e:    &value,
			args: args{
				v: "active",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			e:    &invalidStatus,
			args: args{
				v: "this is not a valid status",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidStatus,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("EpisodeOfCareStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEpisodeOfCareStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     EpisodeOfCareStatusEnum
		wantW string
	}{
		{
			name:  "active",
			e:     EpisodeOfCareStatusEnumActive,
			wantW: strconv.Quote("active"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("EpisodeOfCareStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestEpisodeOfCareStatusHistoryStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    EpisodeOfCareStatusHistoryStatusEnum
		want bool
	}{
		{
			name: "valid status",
			e:    EpisodeOfCareStatusHistoryStatusEnumPlanned,
			want: true,
		},
		{
			name: "invalid status",
			e:    EpisodeOfCareStatusHistoryStatusEnum("invalid"),
			want: false,
		},
		{
			name: "waitlist status",
			e:    EpisodeOfCareStatusHistoryStatusEnumWaitlist,
			want: true,
		},
		{
			name: "active status",
			e:    EpisodeOfCareStatusHistoryStatusEnumActive,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("EpisodeOfCareStatusHistoryStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpisodeOfCareStatusHistoryStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    EpisodeOfCareStatusHistoryStatusEnum
		want string
	}{
		{
			name: "valid status",
			e:    EpisodeOfCareStatusHistoryStatusEnumPlanned,
			want: "planned",
		},
		{
			name: "valid status",
			e:    EpisodeOfCareStatusHistoryStatusEnumEnteredInError,
			want: "entered-in-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("EpisodeOfCareStatusHistoryStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpisodeOfCareStatusHistoryStatusEnum_UnmarshalGQL(t *testing.T) {
	value := EpisodeOfCareStatusHistoryStatusEnumActive
	invalidStatus := EpisodeOfCareStatusHistoryStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *EpisodeOfCareStatusHistoryStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid status",
			e:    &value,
			args: args{
				v: "active",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			e:    &invalidStatus,
			args: args{
				v: "this is not a valid status",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidStatus,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("EpisodeOfCareStatusHistoryStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEpisodeOfCareStatusHistoryStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     EpisodeOfCareStatusHistoryStatusEnum
		wantW string
	}{
		{
			name:  "active",
			e:     EpisodeOfCareStatusHistoryStatusEnumActive,
			wantW: strconv.Quote("active"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("EpisodeOfCareStatusHistoryStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
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
			name: "cancelled",
			e:    ObservationStatusEnumCancelled,
			want: true,
		},
		{
			name: "entered in error",
			e:    ObservationStatusEnumEnteredInError,
			want: true,
		},
		{
			name: "unknown",
			e:    ObservationStatusEnumUnknown,
			want: true,
		},
		{
			name: "invalid",
			e:    ObservationStatusEnum("invalid"),
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
			name: "corrected",
			e:    ObservationStatusEnumCorrected,
			want: "corrected",
		},
		{
			name: "cancelled",
			e:    ObservationStatusEnumCancelled,
			want: "cancelled",
		},
		{
			name: "entered in error",
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
	value := ObservationStatusEnumFinal
	invalidStatus := ObservationStatusEnum("invalid")
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
			name: "valid status",
			e:    &value,
			args: args{
				v: "final",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			e:    &invalidStatus,
			args: args{
				v: "this is not a valid status",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidStatus,
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
			name:  "final",
			e:     ObservationStatusEnumFinal,
			wantW: strconv.Quote("final"),
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

func TestPatientGenderEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    PatientGenderEnum
		want bool
	}{
		{
			name: "valid - female",
			e:    PatientGenderEnumFemale,
			want: true,
		},
		{
			name: "valid - other",
			e:    PatientGenderEnumOther,
			want: true,
		},
		{
			name: "valid - unknown",
			e:    PatientGenderEnumUnknown,
			want: true,
		},
		{
			name: "invalid - empty",
			e:    "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("PatientGenderEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientGenderEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    PatientGenderEnum
		want string
	}{
		{
			name: "male",
			e:    PatientGenderEnumMale,
			want: "male",
		},
		{
			name: "female",
			e:    PatientGenderEnumFemale,
			want: "female",
		},
		{
			name: "other",
			e:    PatientGenderEnumOther,
			want: "other",
		},
		{
			name: "unknown",
			e:    PatientGenderEnumUnknown,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("PatientGenderEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientGenderEnum_UnmarshalGQL(t *testing.T) {
	value := PatientGenderEnumMale
	invalidStatus := PatientGenderEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *PatientGenderEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid status",
			e:    &value,
			args: args{
				v: "male",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			e:    &invalidStatus,
			args: args{
				v: "this is not a valid status",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalidStatus,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("PatientGenderEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatientGenderEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     PatientGenderEnum
		wantW string
	}{
		{
			name:  "male",
			e:     PatientGenderEnumMale,
			wantW: strconv.Quote("male"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PatientGenderEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPatientContactGenderEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    PatientContactGenderEnum
		want bool
	}{
		{
			name: "Valid: male",
			e:    PatientContactGenderEnumMale,
			want: true,
		},
		{
			name: "Valid: female",
			e:    PatientContactGenderEnumFemale,
			want: true,
		},
		{
			name: "Valid: other",
			e:    PatientContactGenderEnumOther,
			want: true,
		},
		{
			name: "Valid: unknown",
			e:    PatientContactGenderEnumUnknown,
			want: true,
		},
		{
			name: "Invalid",
			e:    PatientContactGenderEnum("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("PatientContactGenderEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientContactGenderEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    PatientContactGenderEnum
		want string
	}{
		{
			name: "male",
			e:    PatientContactGenderEnumMale,
			want: "male",
		},
		{
			name: "female",
			e:    PatientContactGenderEnumFemale,
			want: "female",
		},
		{
			name: "other",
			e:    PatientContactGenderEnumOther,
			want: "other",
		},
		{
			name: "unknown",
			e:    PatientContactGenderEnumUnknown,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("PatientContactGenderEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientContactGenderEnum_UnmarshalGQL(t *testing.T) {
	value := PatientContactGenderEnumMale
	invalidType := PatientContactGenderEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *PatientContactGenderEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "male",
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
				t.Errorf("PatientContactGenderEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatientContactGenderEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     PatientContactGenderEnum
		wantW string
	}{
		{
			name:  "male",
			e:     PatientContactGenderEnumMale,
			wantW: strconv.Quote("male"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PatientContactGenderEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPatientLinkTypeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    PatientLinkTypeEnum
		want bool
	}{
		{
			name: "Valid patient link type - replaced_by",
			e:    PatientLinkTypeEnumReplacedBy,
			want: true,
		},
		{
			name: "Invalid patient link type",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("PatientLinkTypeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientLinkTypeEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    PatientLinkTypeEnum
		want string
	}{
		{
			name: "Patient link type - replaced_by",
			e:    PatientLinkTypeEnumReplacedBy,
			want: "replaced-by",
		},
		{
			name: "Patient link type - refer",
			e:    PatientLinkTypeEnumRefer,
			want: "refer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("PatientLinkTypeEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientLinkTypeEnum_UnmarshalGQL(t *testing.T) {
	value := PatientLinkTypeEnumRefer
	invalidType := PatientLinkTypeEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *PatientLinkTypeEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "refer",
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
				t.Errorf("PatientLinkTypeEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatientLinkTypeEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     PatientLinkTypeEnum
		wantW string
	}{
		{
			name:  "refer",
			e:     PatientLinkTypeEnumRefer,
			wantW: strconv.Quote("refer"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PatientLinkTypeEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMedicationStatementStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationStatementStatusEnum
		want bool
	}{
		{
			name: "Valid medication statement status - active",
			e:    MedicationStatementStatusEnumActive,
			want: true,
		},
		{
			name: "Invalid medication statement status",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MedicationStatementStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationStatementStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationStatementStatusEnum
		want string
	}{
		{
			name: "Medication statement status - active",
			e:    MedicationStatementStatusEnumActive,
			want: "active",
		},
		{
			name: "Medication statement status - intended",
			e:    MedicationStatementStatusEnumIntended,
			want: "intended",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MedicationStatementStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationStatementStatusEnum_UnmarshalGQL(t *testing.T) {
	value := MedicationStatementStatusEnumActive
	invalidType := MedicationStatementStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *MedicationStatementStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "active",
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
				t.Errorf("MedicationStatementStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMedicationStatementStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     MedicationStatementStatusEnum
		wantW string
	}{
		{
			name:  "active",
			e:     MedicationStatementStatusEnumActive,
			wantW: strconv.Quote("active"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MedicationStatementStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMedicationStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationStatusEnum
		want bool
	}{
		{
			name: "Valid medication status - active",
			e:    MedicationStatusEnumActive,
			want: true,
		},
		{
			name: "Valid medication status - inactive",
			e:    MedicationStatusEnumInActive,
			want: true,
		},
		{
			name: "Valid medication status - entered-in-error",
			e:    MedicationStatusEnumEnteredInError,
			want: true,
		},
		{
			name: "Invalid medication status",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MedicationStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationStatusEnum
		want string
	}{
		{
			name: "Medication status - active",
			e:    MedicationStatusEnumActive,
			want: "active",
		},
		{
			name: "Medication status - inactive",
			e:    MedicationStatusEnumInActive,
			want: "inactive",
		},
		{
			name: "Medication status - entered-in-error",
			e:    MedicationStatusEnumEnteredInError,
			want: "entered-in-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MedicationStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationStatusEnum_UnmarshalGQL(t *testing.T) {
	value := MedicationStatusEnumActive
	invalidType := MedicationStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *MedicationStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "active",
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
				t.Errorf("MedicationStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMedicationStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     MedicationStatusEnum
		wantW string
	}{
		{
			name:  "active",
			e:     MedicationStatusEnumActive,
			wantW: strconv.Quote("active"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MedicationStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMediaStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MediaStatusEnum
		want bool
	}{
		{
			name: "Valid media status - completed",
			e:    MediaStatusCompleted,
			want: true,
		},
		{
			name: "Invalid media status",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MediaStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMediaStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    MediaStatusEnum
		want string
	}{
		{
			name: "Media status - completed",
			e:    MediaStatusCompleted,
			want: "completed",
		},
		{
			name: "Media status - in-progress",
			e:    MediaStatusInProgress,
			want: "in-progress",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MediaStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMediaStatusEnum_UnmarshalGQL(t *testing.T) {
	value := MediaStatusCompleted
	invalidType := MediaStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *MediaStatusEnum
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
				t.Errorf("MediaStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMediaStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     MediaStatusEnum
		wantW string
	}{
		{
			name:  "completed",
			e:     MediaStatusCompleted,
			wantW: strconv.Quote("completed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MediaStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestDiagnosticReportStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    DiagnosticReportStatusEnum
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    DiagnosticReportStatusAmended,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    DiagnosticReportStatusEnum("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("DiagnosticReportStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiagnosticReportStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		c    DiagnosticReportStatusEnum
		want string
	}{
		{
			name: "Happy Case",
			c:    DiagnosticReportStatusAmended,
			want: DiagnosticReportStatusAmended,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("DiagnosticReportStatusEnum = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiagnosticReportStatusEnum_UnmarshalGQL(t *testing.T) {
	validValue := DiagnosticReportStatusUnknown
	invalidType := DiagnosticReportStatusEnum("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *DiagnosticReportStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: DiagnosticReportStatusAmended,
			},
			c:       (*DiagnosticReportStatusEnum)(&validValue),
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
			c:       (*DiagnosticReportStatusEnum)(&validValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DiagnosticReportStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConsentState_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     DiagnosticReportStatusEnum
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     DiagnosticReportStatusAmended,
			wantW: strconv.Quote("amended"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DiagnosticReportStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestDocumentReferenceStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    DocumentReferenceStatusEnum
		want bool
	}{
		{
			name: "valid current status",
			e:    DocumentReferenceStatusEnumCurrent,
			want: true,
		},
		{
			name: "valid superseded status",
			e:    DocumentReferenceStatusEnumSuperseded,
			want: true,
		},
		{
			name: "valid entered-in-error status",
			e:    DocumentReferenceStatusEnumEnteredInError,
			want: true,
		},
		{
			name: "invalid status",
			e:    DocumentReferenceStatusEnum("unknown"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("DocumentReferenceStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentReferenceStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    DocumentReferenceStatusEnum
		want string
	}{
		{
			name: "current status",
			e:    DocumentReferenceStatusEnumCurrent,
			want: "current",
		},
		{
			name: "superseded status",
			e:    DocumentReferenceStatusEnumSuperseded,
			want: "superseded",
		},
		{
			name: "entered-in-error status",
			e:    DocumentReferenceStatusEnumEnteredInError,
			want: "entered-in-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("DocumentReferenceStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentRelationshipTypeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    DocumentRelationshipTypeEnum
		want bool
	}{
		{
			name: "valid replaces relationship",
			e:    DocumentRelationshipTypeEnumReplaces,
			want: true,
		},
		{
			name: "valid transforms relationship",
			e:    DocumentRelationshipTypeEnumTransforms,
			want: true,
		},
		{
			name: "valid signs relationship",
			e:    DocumentRelationshipTypeEnumSigns,
			want: true,
		},
		{
			name: "valid appends relationship",
			e:    DocumentRelationshipTypeEnumAppends,
			want: true,
		},
		{
			name: "invalid relationship",
			e:    DocumentRelationshipTypeEnum("unknown"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("DocumentRelationshipTypeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentRelationshipTypeEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    DocumentRelationshipTypeEnum
		want string
	}{
		{
			name: "replaces relationship",
			e:    DocumentRelationshipTypeEnumReplaces,
			want: "replaces",
		},
		{
			name: "transforms relationship",
			e:    DocumentRelationshipTypeEnumTransforms,
			want: "transforms",
		},
		{
			name: "signs relationship",
			e:    DocumentRelationshipTypeEnumSigns,
			want: "signs",
		},
		{
			name: "appends relationship",
			e:    DocumentRelationshipTypeEnumAppends,
			want: "appends",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("DocumentRelationshipTypeEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ConsentStatusEnum
		want bool
	}{
		{
			name: "Valid type - Active",
			e:    ConsentStatusActive,
			want: true,
		},
		{
			name: "Valid type - Inactive",
			e:    ConsentStatusInactive,
			want: true,
		},
		{
			name: "Invalid type",
			e:    ConsentStatusEnum("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ConsentStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    ConsentStatusEnum
		want string
	}{
		{
			name: "Active",
			e:    ConsentStatusActive,
			want: "active",
		},
		{
			name: "Inactive",
			e:    ConsentStatusInactive,
			want: "inactive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ConsentStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentStatusEnum_UnmarshalGQL(t *testing.T) {
	value := ConsentStatusActive
	invalidType := ConsentStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ConsentStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "active",
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
				t.Errorf("ConsentStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConsentStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ConsentStatusEnum
		wantW string
	}{
		{
			name:  "active",
			e:     ConsentStatusActive,
			wantW: strconv.Quote("active"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ConsentStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestConsentDecisionEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ConsentDecisionEnum
		want bool
	}{
		{
			name: "Valid type - Deny",
			e:    ConsentDecisionDeny,
			want: true,
		},
		{
			name: "Valid type - Permit",
			e:    ConsentDecisionPermit,
			want: true,
		},
		{
			name: "Invalid type",
			e:    ConsentDecisionEnum("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ConsentDecisionEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentDecisionEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    ConsentDecisionEnum
		want string
	}{
		{
			name: "Deny",
			e:    ConsentDecisionDeny,
			want: "deny",
		},
		{
			name: "Permit",
			e:    ConsentDecisionPermit,
			want: "permit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ConsentDecisionEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentDecisionEnum_UnmarshalGQL(t *testing.T) {
	value := ConsentDecisionDeny
	invalidType := ConsentDecisionEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ConsentDecisionEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "deny",
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
				t.Errorf("ConsentDecisionEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConsentDecisionEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ConsentDecisionEnum
		wantW string
	}{
		{
			name:  "deny",
			e:     ConsentDecisionDeny,
			wantW: strconv.Quote("deny"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ConsentDecisionEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestConsentDataMeaningEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ConsentDataMeaningEnum
		want bool
	}{
		{
			name: "Valid type - Instance",
			e:    ConsentDataMeaningInstance,
			want: true,
		},
		{
			name: "Valid type - Related",
			e:    ConsentDataMeaningRelated,
			want: true,
		},
		{
			name: "Invalid type",
			e:    ConsentDataMeaningEnum("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ConsentDataMeaningEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentDataMeaningEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    ConsentDataMeaningEnum
		want string
	}{
		{
			name: "Instance",
			e:    ConsentDataMeaningInstance,
			want: "instance",
		},
		{
			name: "Related",
			e:    ConsentDataMeaningRelated,
			want: "related",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ConsentDataMeaningEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentDataMeaningEnum_UnmarshalGQL(t *testing.T) {
	value := ConsentDataMeaningInstance
	invalidType := ConsentDataMeaningEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ConsentDataMeaningEnum
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "instance",
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
				t.Errorf("ConsentDataMeaningEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConsentDataMeaningEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ConsentDataMeaningEnum
		wantW string
	}{
		{
			name:  "instance",
			e:     ConsentDataMeaningInstance,
			wantW: strconv.Quote("instance"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ConsentDataMeaningEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestQuestionnaireResponseStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    QuestionnaireResponseStatusEnum
		want bool
	}{
		{
			name: "Valid type - Completed",
			e:    QuestionnaireResponseStatusEnumCompleted,
			want: true,
		},
		{
			name: "Invalid type",
			e:    QuestionnaireResponseStatusEnum("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("QuestionnaireResponseStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireResponseStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    QuestionnaireResponseStatusEnum
		want string
	}{
		{
			name: "Completed",
			e:    QuestionnaireResponseStatusEnumCompleted,
			want: "completed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("QuestionnaireResponseStatusEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireResponseStatusEnum_UnmarshalGQL(t *testing.T) {
	value := QuestionnaireResponseStatusEnumCompleted
	invalidType := QuestionnaireResponseStatusEnum("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *QuestionnaireResponseStatusEnum
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
				t.Errorf("QuestionnaireResponseStatusEnum.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestionnaireResponseStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     QuestionnaireResponseStatusEnum
		wantW string
	}{
		{
			name:  "completed",
			e:     QuestionnaireResponseStatusEnumCompleted,
			wantW: strconv.Quote("completed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("QuestionnaireResponseStatusEnum.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMedicationStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationRequestStatus
		want string
	}{
		{
			name: "completed",
			e:    CompletedMedicationStatus,
			want: "completed",
		},
		{
			name: "cancelled",
			e:    CancelledMedicationStatus,
			want: "cancelled",
		},
		{
			name: "active",
			e:    ActiveMedicationStatus,
			want: "active",
		},
		{
			name: "on_hold",
			e:    OnHoldMedicationStatus,
			want: "on-hold",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MedicationRequestStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationStatus_UnmarshalGQL(t *testing.T) {
	value := CompletedMedicationStatus
	invalidType := MedicationRequestStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *MedicationRequestStatus
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
				t.Errorf("MedicationRequestStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMedicationStatus_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     MedicationRequestStatus
		wantW string
	}{
		{
			name:  "completed",
			e:     CompletedMedicationStatus,
			wantW: strconv.Quote("completed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if got := w.String(); got != tt.wantW {
				t.Errorf("MedicationRequestStatus.MarshalGQL() = %v, want %v", got, tt.wantW)
			}
		})
	}
}

func TestMedicationStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationRequestStatus
		want bool
	}{
		{
			name: "Valid medication request status - completed",
			e:    CompletedMedicationStatus,
			want: true,
		},
		{
			name: "Invalid medication request status",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MedicationRequestStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationIntent_String(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationRequestIntent
		want string
	}{
		{
			name: "proposal",
			e:    ProposalMedicationIntent,
			want: "proposal",
		},
		{
			name: "plan",
			e:    PlanMedicationIntent,
			want: "plan",
		},
		{
			name: "order",
			e:    OrderMedicationIntent,
			want: "order",
		},
		{
			name: "original-order",
			e:    OriginalOrderMedicationIntent,
			want: "original-order",
		},
		{
			name: "reflex-order",
			e:    ReflexOrderMedicationIntent,
			want: "reflex-order",
		},
		{
			name: "filler-order",
			e:    FillerOrderMedicationIntent,
			want: "filler-order",
		},
		{
			name: "instance-order",
			e:    InstanceOrderMedicationIntent,
			want: "instance-order",
		},
		{
			name: "option",
			e:    OptionMedicationIntent,
			want: "option",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MedicationRequestIntent.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationIntent_UnmarshalGQL(t *testing.T) {
	value := ProposalMedicationIntent
	invalidType := MedicationRequestIntent("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *MedicationRequestIntent
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "proposal",
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
				t.Errorf("MedicationRequestIntent.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMedicationIntent_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     MedicationRequestIntent
		wantW string
	}{
		{
			name:  "proposal",
			e:     ProposalMedicationIntent,
			wantW: strconv.Quote("proposal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if got := w.String(); got != tt.wantW {
				t.Errorf("MedicationRequestIntent.MarshalGQL() = %v, want %v", got, tt.wantW)
			}
		})
	}
}

func TestMedicationIntent_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MedicationRequestIntent
		want bool
	}{
		{
			name: "Valid medication request intent - proposal",
			e:    ProposalMedicationIntent,
			want: true,
		},
		{
			name: "Invalid medication request intent",
			e:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MedicationRequestIntent.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueType_String(t *testing.T) {
	tests := []struct {
		name string
		e    IssueType
		want string
	}{
		{
			name: "invalid",
			e:    IssueTypeInvalid,
			want: "invalid",
		},
		{
			name: "structure",
			e:    IssueTypeStructure,
			want: "structure",
		},
		{
			name: "required",
			e:    IssueTypeRequired,
			want: "required",
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("IssueType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueType_UnmarshalGQL(t *testing.T) {
	value := IssueTypeInvalid
	invalidType := IssueType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *IssueType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "invalid",
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
				t.Errorf("IssueType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIssueType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     IssueType
		wantW string
	}{
		{
			name:  "invalid",
			e:     IssueTypeInvalid,
			wantW: strconv.Quote("invalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if got := w.String(); got != tt.wantW {
				t.Errorf("IssueType.MarshalGQL() = %v, want %v", got, tt.wantW)
			}
		})
	}
}

func TestIssueType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    IssueType
		want bool
	}{
		{
			name: "Valid issue type - invalid",
			e:    IssueTypeInvalid,
			want: true,
		},
		{
			name: "Invalid issue type",
			e:    "not-a-valid-type",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IssueType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueSeverity_String(t *testing.T) {
	tests := []struct {
		name string
		e    IssueSeverity
		want string
	}{
		{
			name: "fatal",
			e:    IssueSeverityFatal,
			want: "fatal",
		},
		{
			name: "error",
			e:    IssueSeverityError,
			want: "error",
		},
		{
			name: "warning",
			e:    IssueSeverityWarning,
			want: "warning",
		},
		{
			name: "information",
			e:    IssueSeverityInformation,
			want: "information",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("IssueSeverity.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueSeverity_UnmarshalGQL(t *testing.T) {
	value := IssueSeverityFatal
	invalidType := IssueSeverity("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *IssueSeverity
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "fatal",
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
				t.Errorf("IssueSeverity.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIssueSeverity_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     IssueSeverity
		wantW string
	}{
		{
			name:  "fatal",
			e:     IssueSeverityFatal,
			wantW: strconv.Quote("fatal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if got := w.String(); got != tt.wantW {
				t.Errorf("IssueSeverity.MarshalGQL() = %v, want %v", got, tt.wantW)
			}
		})
	}
}

func TestIssueSeverity_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    IssueSeverity
		want bool
	}{
		{
			name: "Valid issue severity - fatal",
			e:    IssueSeverityFatal,
			want: true,
		},
		{
			name: "Invalid issue severity",
			e:    "not-a-valid-severity",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IssueSeverity.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcedure_IsValid(t *testing.T) {
	tests := []struct {
		name string
		p    ProcedureStatusEnum
		want bool
	}{
		{
			name: "Happy case: A valid procedure enum type",
			p:    ProcedureStatusCompleted,
			want: true,
		},
		{
			name: "Sad case: A valid procedure enum type",
			p:    "ProcedureStatusCompleted",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsValid(); got != tt.want {
				t.Errorf("ProcedureStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcedure_Code(t *testing.T) {
	tests := []struct {
		name string
		code ProcedureStatusEnum
		want string
	}{
		{
			name: "Happy Case: Valid code",
			code: ProcedureStatusCompleted,
			want: ProcedureStatusCompleted.Code(),
		},
		{
			name: "Happy Case: Valid code",
			code: ProcedureStatusInProgress,
			want: ProcedureStatusInProgress.Code(),
		},
		{
			name: "Happy Case: Valid code",
			code: ProcedureStatusNotDone,
			want: ProcedureStatusNotDone.Code(),
		},
		{
			name: "Happy Case: Valid code",
			code: ProcedureStatusEnteredInError,
			want: ProcedureStatusEnteredInError.Code(),
		},
		{
			name: "Happy Case: Valid code",
			code: ProcedureStatusOnHold,
			want: ProcedureStatusOnHold.Code(),
		},
		{
			name: "Happy Case: Valid code",
			code: ProcedureStatusUnknown,
			want: ProcedureStatusUnknown.Code(),
		},
		{
			name: "Happy Case: Valid code",
			code: ProcedureStatusPreparation,
			want: ProcedureStatusPreparation.Code(),
		},
		{
			name: "Happy Case: Valid code",
			code: ProcedureStatusStopped,
			want: ProcedureStatusStopped.Code(),
		},
		{
			name: "Sad Case: Invalid code",
			code: "ProcedureStatusStopped",
			want: "<unknown>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Code(); got != tt.want {
				t.Errorf("ProcedureStatusEnum.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcedure_Display(t *testing.T) {
	tests := []struct {
		name string
		code ProcedureStatusEnum
		want string
	}{
		{
			name: "Happy case: A valid procedure enum type",
			code: ProcedureStatusCompleted,
			want: ProcedureStatusCompleted.Display(),
		},
		{
			name: "Happy case: A valid procedure enum type",
			code: ProcedureStatusPreparation,
			want: ProcedureStatusPreparation.Display(),
		},
		{
			name: "Happy case: A valid procedure enum type",
			code: ProcedureStatusEnteredInError,
			want: ProcedureStatusEnteredInError.Display(),
		},
		{
			name: "Happy case: A valid procedure enum type",
			code: ProcedureStatusInProgress,
			want: ProcedureStatusInProgress.Display(),
		},
		{
			name: "Happy case: A valid procedure enum type",
			code: ProcedureStatusNotDone,
			want: ProcedureStatusNotDone.Display(),
		},
		{
			name: "Happy case: A valid procedure enum type",
			code: ProcedureStatusOnHold,
			want: ProcedureStatusOnHold.Display(),
		},
		{
			name: "Happy case: A valid procedure enum type",
			code: ProcedureStatusStopped,
			want: ProcedureStatusStopped.Display(),
		},
		{
			name: "Sad case: A valid procedure enum type",
			code: "ProcedureStatusCompleted",
			want: "<unknown>",
		},
		{
			name: "Happy case: A valid procedure enum type",
			code: ProcedureStatusUnknown,
			want: ProcedureStatusUnknown.Display(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Display(); got != tt.want {
				t.Errorf("ProcedureStatusEnum.Display() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParticipationEnum_Display(t *testing.T) {
	tests := []struct {
		name string
		d    ParticipantStatusEnum
		want string
	}{
		{
			name: "Happy case: Valid display",
			d:    Accepted,
			want: Accepted.Display(),
		},
		{
			name: "Happy case: Valid display",
			d:    Tenative,
			want: Tenative.Display(),
		},
		{
			name: "Happy case: Valid display",
			d:    NeedsAction,
			want: NeedsAction.Display(),
		},
		{
			name: "Happy case: Valid display",
			d:    Declined,
			want: Declined.Display(),
		},
		{
			name: "Sad case: Invalid display",
			d:    "Declined",
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Display(); got != tt.want {
				t.Errorf("ParticipantStatusEnum.Code()  = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParticipationEnum_Code(t *testing.T) {
	tests := []struct {
		name string
		d    ParticipantStatusEnum
		want string
	}{
		{
			name: "Happy case: Valid display",
			d:    Accepted,
			want: Accepted.Code(),
		},
		{
			name: "Happy case: Valid display",
			d:    Tenative,
			want: Tenative.Code(),
		},
		{
			name: "Happy case: Valid display",
			d:    NeedsAction,
			want: NeedsAction.Code(),
		},
		{
			name: "Happy case: Valid display",
			d:    Declined,
			want: Declined.Code(),
		},
		{
			name: "Sad case: Invalid display",
			d:    "Declined",
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Code(); got != tt.want {
				t.Errorf("ParticipantStatusEnum.Code()  = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParticipationEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		p    ParticipantStatusEnum
		want bool
	}{
		{
			name: "Happy case: Valid code",
			p:    Accepted,
			want: Accepted.IsValid(),
		},
		{
			name: "Sad case: Invalid code",
			p:    "Unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsValid(); got != tt.want {
				t.Errorf("ParticipantStatusEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriorityEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		p    PriorityEnum
		want bool
	}{
		{
			name: "Happy case: Valid code",
			p:    ASAP,
			want: Accepted.IsValid(),
		},
		{
			name: "Sad case: Invalid code",
			p:    "unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsValid(); got != tt.want {
				t.Errorf("PriorityEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriorityEnum_Code(t *testing.T) {
	tests := []struct {
		name string
		d    PriorityEnum
		want string
	}{
		{
			name: "Happy case: Valid Code",
			d:    ASAP,
			want: ASAP.Code(),
		},
		{
			name: "Happy case: Valid Code",
			d:    Emergency,
			want: Emergency.Code(),
		},
		{
			name: "Happy case: Valid Code",
			d:    Routine,
			want: Routine.Code(),
		},
		{
			name: "Happy case: Valid Code",
			d:    Urgent,
			want: Urgent.Code(),
		},
		{
			name: "Sad case: Invalid Code",
			d:    "Declined",
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Code(); got != tt.want {
				t.Errorf("PriorityEnum.Code()  = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriorityEnum_Display(t *testing.T) {
	tests := []struct {
		name string
		d    PriorityEnum
		want string
	}{
		{
			name: "Happy case: Valid display",
			d:    ASAP,
			want: ASAP.Display(),
		},
		{
			name: "Happy case: Valid display",
			d:    Emergency,
			want: Emergency.Display(),
		},
		{
			name: "Happy case: Valid display",
			d:    Routine,
			want: Routine.Display(),
		},
		{
			name: "Happy case: Valid display",
			d:    Urgent,
			want: Urgent.Display(),
		},
		{
			name: "Sad case: Invalid display",
			d:    "Declined",
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Display(); got != tt.want {
				t.Errorf("PriorityEnum.Display()  = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppointmentTypeEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		a    AppointmentTypeEnum
		want bool
	}{
		{
			name: "Happy case: Valid code",
			a:    RoutineAppointment,
			want: true,
		},
		{
			name: "Sad case: Invalid code",
			a:    "unknown",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsValid(); got != tt.want {
				t.Errorf("AppointmentTypeEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppointmentTypeEnum_Code(t *testing.T) {
	tests := []struct {
		name string
		a    AppointmentTypeEnum
		want string
	}{
		{
			name: "Happy case: Valid Code",
			a:    RoutineAppointment,
			want: string(RoutineAppointment),
		},
		{
			name: "Happy case: Valid Code",
			a:    CheckUpAppointment,
			want: string(CheckUpAppointment),
		},
		{
			name: "Happy case: Valid Code",
			a:    EmergencyAppointment,
			want: string(EmergencyAppointment),
		},
		{
			name: "Happy case: Valid Code",
			a:    WalkInAppointment,
			want: string(WalkInAppointment),
		},
		{
			name: "Sad case: Invalid Code",
			a:    "",
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Code(); got != tt.want {
				t.Errorf("AppointmentTypeEnum.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppointmentTypeEnum_Display(t *testing.T) {
	tests := []struct {
		name string
		a    AppointmentTypeEnum
		want string
	}{
		{
			name: "Happy case: Valid display",
			a:    RoutineAppointment,
			want: RoutineAppointment.Display(),
		},
		{
			name: "Happy case: Valid display",
			a:    CheckUpAppointment,
			want: CheckUpAppointment.Display(),
		},
		{
			name: "Happy case: Valid display",
			a:    EmergencyAppointment,
			want: EmergencyAppointment.Display(),
		},
		{
			name: "Happy case: Valid display",
			a:    WalkInAppointment,
			want: WalkInAppointment.Display(),
		},
		{
			name: "Sad case: Invalid display",
			a:    "",
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Display(); got != tt.want {
				t.Errorf("AppointmentTypeEnum.Display() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeableConcept_ConditionCode(t *testing.T) {
	code := "referral"
	system := ("http://hl7.org/fhir/ValueSet/request-status")
	id := gofakeit.UUID()
	tests := []struct {
		name  string
		input FHIRCodeableConcept
		want  *FHIRCodeableReferenceInput
	}{
		{
			name: "Happy case: get code from service request resource",
			input: FHIRCodeableConcept{
				Coding: []*FHIRCoding{
					{
						Code:    (*scalarutils.Code)(&code),
						System:  (*scalarutils.URI)(&system),
						Display: "Referral",
						ID:      &id,
					},
				},
			},
			want: &FHIRCodeableReferenceInput{
				Concept: &FHIRCodeableConceptInput{
					Coding: []*FHIRCodingInput{
						{
							Code:    scalarutils.Code(code),
							System:  (*scalarutils.URI)(&system),
							Display: "Referral",
							ID:      id,
						},
					},
				},
			},
		},
		{
			name: "Sad case: missing code from service request resource",
			input: FHIRCodeableConcept{
				Coding: []*FHIRCoding{
					{
						Code:    nil,
						System:  nil,
						Display: "",
						ID:      nil,
					},
					{
						Code:    (*scalarutils.Code)(&code),
						System:  (*scalarutils.URI)(&system),
						Display: "Referral",
						ID:      &id,
					},
				},
			},
			want: &FHIRCodeableReferenceInput{
				Concept: &FHIRCodeableConceptInput{
					Coding: []*FHIRCodingInput{
						{
							Code:    scalarutils.Code(code),
							System:  (*scalarutils.URI)(&system),
							Display: "Referral",
							ID:      id,
						},
					},
				},
			},
		},
		{
			name: "Sad case: get code from service request fails",
			input: FHIRCodeableConcept{
				Coding: []*FHIRCoding{},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.input.ConditionCode()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("TestCodeableConcept.ConditionCode() got = %s", diff)
			}
		})
	}
}
