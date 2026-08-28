package domain

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/scalarutils"
)

func TestRelationshipTypeDisplay(t *testing.T) {
	type args struct {
		val RelationshipType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy Case - Relationship type C",
			args: args{
				val: RelationshipTypeC,
			},
			want: "Emergency Contact",
		},
		{
			name: "Happy Case - Relationship type E",
			args: args{
				val: RelationshipTypeE,
			},
			want: "Employer",
		},
		{
			name: "Happy Case - Relationship type F",
			args: args{
				val: RelationshipTypeF,
			},
			want: "Federal Agency",
		},
		{
			name: "Happy Case - Relationship type I",
			args: args{
				val: RelationshipTypeI,
			},
			want: "Insurance Company",
		},
		{
			name: "Happy Case - Relationship type N",
			args: args{
				val: RelationshipTypeN,
			},
			want: "Next-of-Kin",
		},
		{
			name: "Happy Case - Relationship type O",
			args: args{
				val: RelationshipTypeO,
			},
			want: "Other",
		},
		{
			name: "Happy Case - Relationship type S",
			args: args{
				val: RelationshipTypeS,
			},
			want: "State Agency",
		},
		{
			name: "Happy Case - Relationship type U",
			args: args{
				val: RelationshipTypeU,
			},
			want: "Unknown",
		},
		{
			name: "Happy Case - Relationship type U",
			args: args{
				val: RelationshipType("random"),
			},
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RelationshipTypeDisplay(tt.args.val); got != tt.want {
				t.Errorf("RelationshipTypeDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaritalStatusDisplay(t *testing.T) {
	type args struct {
		val MaritalStatus
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy Case - Marital status is Annulled",
			args: args{
				val: MaritalStatusA,
			},
			want: "Annulled",
		},
		{
			name: "Marital status is Divorced",
			args: args{
				val: MaritalStatusD,
			},
			want: "Divorced",
		},
		{
			name: "Marital status is Interlocutory",
			args: args{
				val: MaritalStatusI,
			},
			want: "Interlocutory",
		},
		{
			name: "Marital status is Legally Separated",
			args: args{
				val: MaritalStatusL,
			},
			want: "Legally Separated",
		},
		{
			name: "Marital status is Married",
			args: args{
				val: MaritalStatusM,
			},
			want: "Married",
		},
		{
			name: "Marital status is Polygamous",
			args: args{
				val: MaritalStatusP,
			},
			want: "Polygamous",
		},
		{
			name: "Marital status is Never Married",
			args: args{
				val: MaritalStatusS,
			},
			want: "Never Married",
		},
		{
			name: "Marital status is Domestic Partner",
			args: args{
				val: MaritalStatusT,
			},
			want: "Domestic Partner",
		},
		{
			name: "Marital status is unmarried",
			args: args{
				val: MaritalStatusU,
			},
			want: "unmarried",
		},
		{
			name: "Marital status is Widowed",
			args: args{
				val: MaritalStatusW,
			},
			want: "Widowed",
		},
		{
			name: "Marital status is unknown",
			args: args{
				val: MaritalStatusUnk,
			},
			want: "unknown",
		},
		{
			name: "Marital status is outside the range",
			args: args{
				val: MaritalStatus("random"),
			},
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaritalStatusDisplay(tt.args.val); got != tt.want {
				t.Errorf("MaritalStatusDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    RelationshipType
		want bool
	}{
		{
			name: "Happy Case - Valid relationship type",
			e:    RelationshipTypeC,
			want: true,
		},
		{
			name: "Valid Relationship Type - E",
			e:    RelationshipTypeE,
			want: true,
		},
		{
			name: "Valid Relationship Type - F",
			e:    RelationshipTypeF,
			want: true,
		},
		{
			name: "Valid Relationship Type - I",
			e:    RelationshipTypeI,
			want: true,
		},
		{
			name: "Valid Relationship Type - N",
			e:    RelationshipTypeN,
			want: true,
		},
		{
			name: "Valid Relationship Type - O",
			e:    RelationshipTypeO,
			want: true,
		},
		{
			name: "Valid Relationship Type - S",
			e:    RelationshipTypeS,
			want: true,
		},
		{
			name: "Invalid Relationship Type - U",
			e:    RelationshipType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("RelationshipType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipType_String(t *testing.T) {
	tests := []struct {
		name string
		e    RelationshipType
		want string
	}{
		{
			name: "Happy case - return valid string",
			e:    RelationshipTypeC,
			want: "C",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("RelationshipType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipType_UnmarshalGQL(t *testing.T) {
	value := RelationshipTypeC
	invalidType := RelationshipType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *RelationshipType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "C",
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
				t.Errorf("RelationshipType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRelationshipType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     RelationshipType
		wantW string
	}{
		{
			name:  "valid type enums",
			e:     RelationshipTypeC,
			wantW: strconv.Quote("C"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("RelationshipType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMaritalStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MaritalStatus
		want bool
	}{
		{
			name: "Happy case - valid marital status",
			e:    MaritalStatusA,
			want: true,
		},
		{
			name: "Invalid marital status",
			e:    MaritalStatus("X"),
			want: false,
		},
		{
			name: "Valid MaritalStatusD",
			e:    MaritalStatusD,
			want: true,
		},
		{
			name: "Valid MaritalStatusI",
			e:    MaritalStatusI,
			want: true,
		},
		{
			name: "Valid MaritalStatusL",
			e:    MaritalStatusL,
			want: true,
		},
		{
			name: "Valid MaritalStatusM",
			e:    MaritalStatusM,
			want: true,
		},
		{
			name: "Valid MaritalStatusP",
			e:    MaritalStatusP,
			want: true,
		},
		{
			name: "Valid MaritalStatusS",
			e:    MaritalStatusS,
			want: true,
		},
		{
			name: "Valid MaritalStatusT",
			e:    MaritalStatusT,
			want: true,
		},
		{
			name: "Valid MaritalStatusU",
			e:    MaritalStatusU,
			want: true,
		},
		{
			name: "Valid MaritalStatusW",
			e:    MaritalStatusW,
			want: true,
		},
		{
			name: "Valid MaritalStatusUnk",
			e:    MaritalStatusUnk,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MaritalStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaritalStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    MaritalStatus
		want string
	}{
		{
			name: "Happy case - valid string",
			e:    MaritalStatusA,
			want: "A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MaritalStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaritalStatus_UnmarshalGQL(t *testing.T) {
	value := MaritalStatusA
	invalidType := MaritalStatus("ZZ")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *MaritalStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "A",
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
				t.Errorf("MaritalStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMaritalStatus_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     MaritalStatus
		wantW string
	}{
		{
			name:  "valid type enums",
			e:     MaritalStatusA,
			wantW: strconv.Quote("A"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MaritalStatus.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestIDDocumentType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    IDDocumentType
		want bool
	}{
		{
			name: "Valid document type",
			e:    IDDocumentTypeAlienID,
			want: true,
		},
		{
			name: "Valid document type",
			e:    IDDocumentTypePassport,
			want: true,
		},
		{
			name: "Valid document type",
			e:    IDDocumentTypeNationalID,
			want: true,
		},
		{
			name: "invalid document type",
			e:    IDDocumentType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IDDocumentType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIDDocumentType_String(t *testing.T) {
	tests := []struct {
		name string
		e    IDDocumentType
		want string
	}{
		{
			name: "Happy case - return valid string",
			e:    IDDocumentTypePassport,
			want: "passport",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("IDDocumentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIDDocumentType_UnmarshalGQL(t *testing.T) {
	value := IDDocumentTypePassport
	invalidType := IDDocumentType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *IDDocumentType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "passport",
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
				t.Errorf("IDDocumentType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIDDocumentType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     IDDocumentType
		wantW string
	}{
		{
			name:  "valid type enums",
			e:     IDDocumentTypePassport,
			wantW: strconv.Quote("passport"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("IDDocumentType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFHIRPatient_Names(t *testing.T) {
	userName := gofakeit.Name()
	type fields struct {
		ID                   *string
		Text                 *FHIRNarrative
		Identifier           []*FHIRIdentifier
		Active               *bool
		Name                 []*FHIRHumanName
		Telecom              []*FHIRContactPoint
		Gender               *PatientGenderEnum
		BirthDate            *scalarutils.Date
		DeceasedBoolean      *bool
		DeceasedDateTime     *scalarutils.Date
		Address              []*FHIRAddress
		MaritalStatus        *FHIRCodeableConcept
		MultipleBirthBoolean *bool
		MultipleBirthInteger *string
		Photo                []*FHIRAttachment
		Contact              []*FHIRPatientContact
		Communication        []*FHIRPatientCommunication
		GeneralPractitioner  []*FHIRReference
		ManagingOrganization *FHIRReference
		Link                 []*FHIRPatientLink
		Meta                 *FHIRMeta
		Extension            []*FHIRExtension
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy Case - Empty name",
			fields: fields{
				Name: []*FHIRHumanName{},
			},
			want: "",
		},
		{
			name:   "Happy Case - Nil name",
			fields: fields{},
			want:   "",
		},
		{
			name: "Happy case - non-empty name",
			fields: fields{
				Name: []*FHIRHumanName{
					{
						Text: userName,
					},
				},
			},
			want: userName,
		},
		{
			name: "Happy case - empty name",
			fields: fields{
				Name: []*FHIRHumanName{
					{
						Text: "",
					},
				},
			},
			want: "",
		},
		{
			name: "Happy case - nil name",
			fields: fields{
				Name: []*FHIRHumanName{nil},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FHIRPatient{
				ID:                   tt.fields.ID,
				Text:                 tt.fields.Text,
				Identifier:           tt.fields.Identifier,
				Active:               tt.fields.Active,
				Name:                 tt.fields.Name,
				Telecom:              tt.fields.Telecom,
				Gender:               tt.fields.Gender,
				BirthDate:            tt.fields.BirthDate,
				DeceasedBoolean:      tt.fields.DeceasedBoolean,
				DeceasedDateTime:     tt.fields.DeceasedDateTime,
				Address:              tt.fields.Address,
				MaritalStatus:        tt.fields.MaritalStatus,
				MultipleBirthBoolean: tt.fields.MultipleBirthBoolean,
				MultipleBirthInteger: tt.fields.MultipleBirthInteger,
				Photo:                tt.fields.Photo,
				Contact:              tt.fields.Contact,
				Communication:        tt.fields.Communication,
				GeneralPractitioner:  tt.fields.GeneralPractitioner,
				ManagingOrganization: tt.fields.ManagingOrganization,
				Link:                 tt.fields.Link,
				Meta:                 tt.fields.Meta,
				Extension:            tt.fields.Extension,
			}
			if got := p.Names(); got != tt.want {
				t.Errorf("FHIRPatient.Names() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDummy_SetID(t *testing.T) {
	type fields struct {
		ID string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Happy case - successfully set ID",
			fields: fields{
				ID: "12345",
			},
			args: args{
				id: "12345",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dummy{
				ID: tt.fields.ID,
			}
			d.SetID(tt.args.id)
		})
	}
}

func TestDummy_IsNode(t *testing.T) {
	type fields struct {
		ID string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Happy case - is node",
			fields: fields{
				ID: "12345",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dummy{
				ID: tt.fields.ID,
			}
			d.IsNode()
		})
	}
}

func TestDummy_IsEntity(t *testing.T) {
	type fields struct {
		ID string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Happy case - is entity",
			fields: fields{
				ID: "12345",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Dummy{
				ID: tt.fields.ID,
			}
			d.IsEntity()
		})
	}
}
