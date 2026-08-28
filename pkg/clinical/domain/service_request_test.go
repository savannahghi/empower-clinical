package domain

import (
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/scalarutils"
)

func TestFHIRServiceRequest_GetReceivingFacilityDetails(t *testing.T) {
	type fields struct {
		ID                      *string
		Text                    *FHIRNarrative
		Identifier              []*FHIRIdentifier
		InstantiatesCanonical   *scalarutils.Canonical
		InstantiatesURI         *scalarutils.Instant
		BasedOn                 []*FHIRReference
		Replaces                []*FHIRReference
		Requisition             *FHIRIdentifier
		Status                  ServiceRequestStatusEnum
		Intent                  ServiceRequestIntentEnum
		Category                []*FHIRCodeableConcept
		Priority                ServiceRequestPriorityEnum
		DoNotPerform            *bool
		Code                    *FHIRCodeableReference
		OrderDetail             []*FHIRCodeableConcept
		QuantityQuantity        *FHIRQuantity
		QuantityRatio           *FHIRRatio
		QuantityRange           *FHIRRange
		Subject                 *FHIRReference
		Encounter               *FHIRReference
		OccurrenceDateTime      *scalarutils.Date
		OccurrencePeriod        *FHIRPeriod
		OccurrenceTiming        *FHIRTiming
		AsNeededBoolean         *bool
		AsNeededCodeableConcept *scalarutils.Code
		AuthoredOn              *scalarutils.DateTime
		Requester               *FHIRReference
		PerformerType           *FHIRCodeableConcept
		Performer               []*FHIRReference
		LocationCode            *scalarutils.Code
		LocationReference       []*FHIRReference
		ReasonCode              *scalarutils.Code
		ReasonReference         []*FHIRReference
		Insurance               []*FHIRReference
		SupportingInfo          []*FHIRReference
		Specimen                []*FHIRReference
		BodySite                []*FHIRCodeableConcept
		Note                    []*FHIRAnnotation
		PatientInstruction      *string
		RelevantHistory         []*FHIRReference
		Meta                    *FHIRMeta
		Extension               []*FHIRExtension
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ReceivingFacility
		wantErr bool
	}{
		{
			name: "Happy Case: Get receiving facility details",
			fields: fields{
				Extension: []*FHIRExtension{
					{
						URL: "http://savannahghi.org/fhir/StructureDefinition/referred-facility",
						Extension: []Extension{
							{
								URL:         "facilityName",
								ValueString: "Nairobi Hospital",
							},
							{
								URL:         "facilityCounty",
								ValueString: "Nairobi",
							},
							{
								URL:         "facilityContact",
								ValueString: "0711223344",
							},
						},
					},
				},
			},
			want: &ReceivingFacility{
				FacilityName:    "Nairobi Hospital",
				FacilityCounty:  "Nairobi",
				FacilityContact: "0711223344",
			},
			wantErr: false,
		},
		{
			name:    "Sad Case: nil extension",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FHIRServiceRequest{
				ID:                      tt.fields.ID,
				Text:                    tt.fields.Text,
				Identifier:              tt.fields.Identifier,
				InstantiatesCanonical:   tt.fields.InstantiatesCanonical,
				InstantiatesURI:         tt.fields.InstantiatesURI,
				BasedOn:                 tt.fields.BasedOn,
				Replaces:                tt.fields.Replaces,
				Requisition:             tt.fields.Requisition,
				Status:                  tt.fields.Status,
				Intent:                  tt.fields.Intent,
				Category:                tt.fields.Category,
				Priority:                tt.fields.Priority,
				DoNotPerform:            tt.fields.DoNotPerform,
				Code:                    tt.fields.Code,
				OrderDetail:             tt.fields.OrderDetail,
				QuantityQuantity:        tt.fields.QuantityQuantity,
				QuantityRatio:           tt.fields.QuantityRatio,
				QuantityRange:           tt.fields.QuantityRange,
				Subject:                 tt.fields.Subject,
				Encounter:               tt.fields.Encounter,
				OccurrenceDateTime:      tt.fields.OccurrenceDateTime,
				OccurrencePeriod:        tt.fields.OccurrencePeriod,
				OccurrenceTiming:        tt.fields.OccurrenceTiming,
				AsNeededBoolean:         tt.fields.AsNeededBoolean,
				AsNeededCodeableConcept: tt.fields.AsNeededCodeableConcept,
				AuthoredOn:              tt.fields.AuthoredOn,
				Requester:               tt.fields.Requester,
				PerformerType:           tt.fields.PerformerType,
				Performer:               tt.fields.Performer,
				LocationCode:            tt.fields.LocationCode,
				LocationReference:       tt.fields.LocationReference,
				ReasonCode:              tt.fields.ReasonCode,
				ReasonReference:         tt.fields.ReasonReference,
				Insurance:               tt.fields.Insurance,
				SupportingInfo:          tt.fields.SupportingInfo,
				Specimen:                tt.fields.Specimen,
				BodySite:                tt.fields.BodySite,
				Note:                    tt.fields.Note,
				PatientInstruction:      tt.fields.PatientInstruction,
				RelevantHistory:         tt.fields.RelevantHistory,
				Meta:                    tt.fields.Meta,
				Extension:               tt.fields.Extension,
			}
			got, err := f.GetReceivingFacilityDetails()
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRServiceRequest.GetReceivingFacilityDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRServiceRequest.GetReceivingFacilityDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRServiceRequest_GetPatientReferralTest(t *testing.T) {
	code := scalarutils.Code("TEST")

	type fields struct {
		Code *FHIRCodeableReference
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy case: get correct display",
			fields: fields{
				Code: &FHIRCodeableReference{
					Concept: &FHIRCodeableConcept{
						ID: new(string),
						Coding: []*FHIRCoding{
							{
								Code:    &code,
								Display: "TEST",
							},
						},
						Text: "TEST",
					},
				},
			},
			want: "TEST",
		},
		{
			name:   "Sad case: empty patient referral test",
			fields: fields{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FHIRServiceRequest{
				Code: tt.fields.Code,
			}
			if got := f.GetPatientReferralTest(); got != tt.want {
				t.Errorf("FHIRServiceRequest.GetPatientReferralTest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRServiceRequest_GetPractitionersNotes(t *testing.T) {
	practionersNote := scalarutils.Markdown("The patient is doing great")

	type fields struct {
		Note []*FHIRAnnotation
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy case: successfully get practioners notes",
			fields: fields{
				Note: []*FHIRAnnotation{
					{
						Text: &practionersNote,
					},
				},
			},
			want: "The patient is doing great",
		},
		{
			name: "Sad case: no practitioner's notes",
			fields: fields{
				Note: []*FHIRAnnotation{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f *FHIRServiceRequest
			if tt.fields.Note != nil {
				f = &FHIRServiceRequest{
					Note: tt.fields.Note,
				}
			}
			if got := f.GetPractitionersNotes(); got != tt.want {
				t.Errorf("FHIRServiceRequest.GetPractitionersNotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRServiceRequest_GetFacilityName(t *testing.T) {
	facilitySystem := scalarutils.URI("http://mycarehub/tenant-identification/facility")

	code := scalarutils.Code("159623")
	type fields struct {
		Meta *FHIRMeta
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy case: Successfully getting facility name",
			fields: fields{
				Meta: &FHIRMeta{
					Tag: []FHIRCoding{
						{
							System:  &facilitySystem,
							Code:    &code,
							Display: "Tumaini Clinic",
						},
						{
							System:  &facilitySystem,
							Display: "Tumaini Clinic",
						},
					},
				},
			},
			want: "Tumaini Clinic",
		},
		{
			name: "Sad case: Does not get the facility name",
			fields: fields{
				Meta: &FHIRMeta{
					Tag: nil,
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f *FHIRServiceRequest
			if tt.fields.Meta != nil {
				f = &FHIRServiceRequest{
					Meta: tt.fields.Meta,
				}
			}
			if got := f.GetFacilityName(); got != tt.want {
				t.Errorf("FHIRServiceRequest.GetFacilityName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRServiceRequest_GetPerformerID(t *testing.T) {
	id := gofakeit.UUID()
	type args struct {
		input *FHIRServiceRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy Case: Gets performer ID",
			args: args{
				input: &FHIRServiceRequest{
					Performer: []*FHIRReference{
						{
							ID:      &id,
							Display: id,
						},
					},
				},
			},
			want: id,
		},
		{
			name: "Sad Case: No performer ID",
			args: args{
				input: &FHIRServiceRequest{
					Performer: []*FHIRReference{
						{
							Display: "",
						},
					},
				},
			},
			want: "",
		},
		{
			name: "Sad Case: No performer",
			args: args{
				input: &FHIRServiceRequest{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.input.GetPerformerID()
			if got != tt.want {
				t.Errorf("FHIRServiceRequest.GetPerformerID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRServiceRequest_GetPatientReferralReason(t *testing.T) {
	type fields struct {
		ID   *string
		Code *FHIRCodeableReference
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Happy case: get referral reason",
			fields: fields{
				Code: &FHIRCodeableReference{
					ID: new(string),
					Concept: &FHIRCodeableConcept{
						Text: "Test",
					},
					Reference: &FHIRReference{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FHIRServiceRequest{
				ID:   tt.fields.ID,
				Code: tt.fields.Code,
			}

			_ = f.GetPatientReferralReason()
		})
	}
}

func TestFHIRServiceRequest_GetCoding(t *testing.T) {
	code := "123"

	type fields struct {
		Code *FHIRCodeableReference
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  string
	}{
		{
			name: "Happy case: get codings",
			fields: fields{
				Code: &FHIRCodeableReference{
					Concept: &FHIRCodeableConcept{
						Coding: []*FHIRCoding{
							{
								Code:    (*scalarutils.Code)(&code),
								Display: "Test code",
							},
						},
					},
				},
			},
			want:  "123",
			want1: "Test code",
		},
		{
			name: "Happy case: no codings",
			fields: fields{
				Code: &FHIRCodeableReference{
					Concept: &FHIRCodeableConcept{},
				},
			},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FHIRServiceRequest{
				Code: tt.fields.Code,
			}

			got, got1 := s.GetCoding()
			if got != tt.want {
				t.Errorf("FHIRServiceRequest.GetCoding() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FHIRServiceRequest.GetCoding() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestServiceRequestCategoryType_Display(t *testing.T) {
	type args struct {
		input ServiceRequestCategoryType
		want  string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy case: Lab procedure display",
			args: args{
				input: LaboratoryProcedureCategoryType,
				want:  "Laboratory procedure",
			},
		},
		{
			name: "Happy case: counselling",
			args: args{
				input: CounsellingCategoryType,
				want:  "Counselling",
			},
		},
		{
			name: "Happy case: referral display",
			args: args{
				input: ReferralCategoryType,
				want:  "Referral",
			},
		},
		{
			name: "Happy case: Lab procedure display",
			args: args{
				input: EducationCategoryType,
				want:  "Education",
			},
		},
		{
			name: "Sad case: Unknown",
			args: args{
				input: "EducationCategType",
				want:  "Unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.input.Display(); got != tt.args.want {
				t.Errorf("ServiceRequestCategoryType.Display() = %v, want %v", got, tt.args.want)
			}
		})
	}
}

func TestFHIRServiceRequest_GetCode(t *testing.T) {
	code := "referall"
	tests := []struct {
		name   string
		input  *FHIRServiceRequest
		output string
	}{
		{
			name:   "Happy case: found code from ServiceRequest",
			output: code,
			input: &FHIRServiceRequest{
				Code: &FHIRCodeableReference{
					Concept: &FHIRCodeableConcept{
						Coding: []*FHIRCoding{
							{
								Code: (*scalarutils.Code)(&code),
							},
						},
					},
				},
			},
		},
		{
			name:   "Sad case: no code found from ServiceRequest",
			output: "",
			input: &FHIRServiceRequest{
				Code: &FHIRCodeableReference{
					Concept: &FHIRCodeableConcept{
						Coding: []*FHIRCoding{
							{
								Code: nil,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.GetCode(); got != tt.output {
				t.Errorf("ServiceRequest.GetCode() got = %v expected = %v", tt.input, tt.output)
			}
		})
	}
}
