package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
)

func TestPatientLink_GetID(t *testing.T) {
	id := uuid.New().String()
	type fields struct {
		ID        string
		PatientID string
		OpaqueID  string
		Expires   time.Time
		Deleted   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy case",
			fields: fields{
				ID: id,
			},
			want: id,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := &PatientLink{
				ID:        tt.fields.ID,
				PatientID: tt.fields.PatientID,
				OpaqueID:  tt.fields.OpaqueID,
				Expires:   tt.fields.Expires,
				Deleted:   tt.fields.Deleted,
			}
			if got := pl.GetID(); got != tt.want {
				t.Errorf("PatientLink.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientLink_SetID(t *testing.T) {
	id := uuid.New().String()
	type fields struct {
		ID        string
		PatientID string
		OpaqueID  string
		Expires   time.Time
		Deleted   bool
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
			name: "Happy case",
			fields: fields{
				ID: id,
			},
			args: args{
				id: id,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := &PatientLink{
				ID:        tt.fields.ID,
				PatientID: tt.fields.PatientID,
				OpaqueID:  tt.fields.OpaqueID,
				Expires:   tt.fields.Expires,
				Deleted:   tt.fields.Deleted,
			}
			pl.SetID(tt.args.id)
		})
	}
}

func TestFHIRPatient_GetIDs(t *testing.T) {
	healthID := scalarutils.URI("HEALTH_ID")
	nationalID := scalarutils.URI("NATIONAL_ID")
	idDocument := scalarutils.URI("healthcloud.iddocument")
	healthIDValue := scalarutils.Code("1234567")
	nationalIDValue := scalarutils.Code("5678901")
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
		name    string
		fields  fields
		want    *Identifiers
		wantErr bool
	}{
		{
			name: "Happy Case: Get all identifiers",
			fields: fields{
				Identifier: []*FHIRIdentifier{
					{
						Use: "official",
						Type: FHIRCodeableConcept{
							ID: new(string),
							Coding: []*FHIRCoding{
								{
									System:  &healthID,
									Code:    &healthIDValue,
									Display: string(healthIDValue),
								},
							},
							Text: "",
						},
						System: &idDocument,
						Value:  string(healthIDValue),
					},
					{
						Use: "official",
						Type: FHIRCodeableConcept{
							ID: new(string),
							Coding: []*FHIRCoding{
								{
									System:  &nationalID,
									Code:    &nationalIDValue,
									Display: string(nationalIDValue),
								},
							},
							Text: "",
						},
						System: &idDocument,
						Value:  string(nationalIDValue),
					},
				},
			},
			want: &Identifiers{
				HealthID:   string(healthIDValue),
				NationalID: string(nationalIDValue),
			},
			wantErr: false,
		},
		{
			name:    "Sad Case: Nil identifiers",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FHIRPatient{
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
			got, err := p.GetIDs()
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRPatient.GetIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRPatient.GetIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRPatient_GetPhoneNumbers(t *testing.T) {
	phone := ContactPointSystemEnumPhone
	email := ContactPointSystemEnumEmail
	phoneNumber := gofakeit.Phone()
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
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get the contacts",
			fields: fields{
				Telecom: []*FHIRContactPoint{
					{
						System: &phone,
						Value:  &phoneNumber,
					},
				},
			},
			want:    []string{phoneNumber},
			wantErr: false,
		},
		{
			name:    "Sad Case: nil contacts",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "Sad Case: no phone contacts",
			fields: fields{
				Telecom: []*FHIRContactPoint{
					{
						System: &email,
						Value:  &phoneNumber,
					},
				},
			},
			wantErr: true,
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
			got, err := p.GetPhoneNumbers()
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRPatient.GetPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRPatient.GetPhoneNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRPatient_GetPatientTelecom(t *testing.T) {
	system := "phone"
	value := "12345"
	type fields struct {
		ID      *string
		Telecom []*FHIRContactPoint
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy case: find patient telecom",
			fields: fields{
				Telecom: []*FHIRContactPoint{
					{
						System: (*ContactPointSystemEnum)(&system),
						Value:  &value,
					},
				},
			},
			want: "12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FHIRPatient{
				ID:      tt.fields.ID,
				Telecom: tt.fields.Telecom,
			}
			if got := f.GetPatientTelecom(); got != tt.want {
				t.Errorf("FHIRPatient.GetPatientTelecom() = %v, want %v", got, tt.want)
			}
		})
	}
}
