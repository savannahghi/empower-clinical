package dto_test

import (
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
)

func TestCreateReferralDetailConnection(t *testing.T) {
	referralID := gofakeit.UUID()
	name := gofakeit.Name()

	type args struct {
		referralDetail []*dto.ReferralDetail
		pageInfo       dto.PageInfo
		total          int
	}

	tests := []struct {
		name string
		args args
		want dto.ReferralDetailConnection
	}{
		{
			name: "Happy case",
			args: args{
				referralDetail: []*dto.ReferralDetail{
					{
						ID:          referralID,
						PatientName: name,
						PatientID:   referralID,
					},
				},
				pageInfo: dto.PageInfo{
					HasNextPage:     true,
					HasPreviousPage: true,
				},
				total: 3,
			},
			want: dto.ReferralDetailConnection{
				TotalCount: 3,
				PageInfo: dto.PageInfo{
					HasNextPage:     true,
					HasPreviousPage: true,
				},
				Edges: []dto.ReferralDetailEdge{
					{
						Node: dto.ReferralDetail{
							ID:          referralID,
							PatientName: name,
							PatientID:   referralID,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dto.CreateReferralDetailConnection(tt.args.referralDetail, tt.args.pageInfo, tt.args.total); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateReferralDetailConnection() = got %v want %v", got, tt.want)

			}
		})
	}
}
