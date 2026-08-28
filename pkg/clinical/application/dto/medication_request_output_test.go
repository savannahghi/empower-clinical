package dto_test

import (
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
)

func TestCreateMedicationRequestConnection(t *testing.T) {
	medID := gofakeit.UUID()
	facilityName := gofakeit.Name()
	type args struct {
		medicationRequests []*dto.MedicationRequestOutput
		pageInfo           dto.PageInfo
		total              int
	}
	tests := []struct {
		name string
		args args
		want dto.MedicationRequestConnection
	}{
		{
			name: "Happy case",
			args: args{
				medicationRequests: []*dto.MedicationRequestOutput{
					{
						ID:           medID,
						FacilityName: facilityName,
					},
				},
				total: 1,
			},
			want: dto.MedicationRequestConnection{
				TotalCount: 1,
				Edges: []dto.MedicationRequestEdge{
					{
						Node: dto.MedicationRequestOutput{
							ID:           medID,
							FacilityName: facilityName,
						},
						Cursor: medID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dto.CreateMedicationRequestConnection(tt.args.medicationRequests, tt.args.pageInfo, tt.args.total); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMedicationRequestConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
