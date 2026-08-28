package domain

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/scalarutils"
)

func TestFHIRTask_GetServiceRequestIDFromTask(t *testing.T) {
	serviceReqType := scalarutils.URI("Patient Referral (Service Request)")
	id := gofakeit.UUID()
	ref := fmt.Sprintf("ServiceRequest/%s", id)

	type fields struct {
		BasedOn []*FHIRReference
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Happy case: get referral ID from task",
			fields: fields{
				BasedOn: []*FHIRReference{
					{
						ID:        &id,
						Reference: &ref,
						Display:   ReferralServiceRequestType.String(),
						Type:      &serviceReqType,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &FHIRTask{
				BasedOn: tt.fields.BasedOn,
			}
			_, err := tr.GetServiceRequestIDFromTask()
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRTask.GetServiceRequestIDFromTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
