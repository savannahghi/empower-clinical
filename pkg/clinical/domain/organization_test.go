package domain

import (
	"testing"

	"github.com/brianvoe/gofakeit"
)

func TestFHIROrganization_GetReceivingFacilityContactDetails(t *testing.T) {
	emailContactSystem := "email"
	phoneContactSystem := "phone"
	phone := gofakeit.Phone()
	email := gofakeit.Email()
	type args struct {
		input *FHIROrganization
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: organization email contact exists",
			args: args{
				input: &FHIROrganization{
					Contact: []FHIROrganizationContact{
						{
							Telecom: []FHIRContactPoint{
								{
									System: (*ContactPointSystemEnum)(&emailContactSystem),
									Value:  &email,
								},
								{
									System: (*ContactPointSystemEnum)(&phoneContactSystem),
									Value:  &phone,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: empty contact details",
			args: args{
				input: &FHIROrganization{
					Contact: []FHIROrganizationContact{
						{
							Telecom: []FHIRContactPoint{},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.args.input.GetReceivingFacilityContactDetails()
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIROrganization.GetReceivingFacilityContactDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
