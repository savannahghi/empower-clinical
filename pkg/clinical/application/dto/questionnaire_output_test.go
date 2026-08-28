package dto_test

import (
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
)

func TestCreateQuestionnaireConnection(t *testing.T) {
	url := gofakeit.URL()
	questionnaireID := gofakeit.UUID()
	type args struct {
		questionnaires []*dto.Questionnaire
		pageInfo       dto.PageInfo
		total          int
	}
	tests := []struct {
		name    string
		args    args
		want    dto.QuestionnaireConnection
		wantErr bool
	}{
		{
			name: "Happy case: successfully creates questionnaire connection",
			args: args{
				questionnaires: []*dto.Questionnaire{
					{
						ID:  questionnaireID,
						URL: (*scalarutils.URI)(&url),
					},
					{
						ID:  questionnaireID,
						URL: (*scalarutils.URI)(&url),
					},
				},
				pageInfo: dto.PageInfo{
					HasNextPage:     true,
					HasPreviousPage: true,
				},
				total: 10,
			},
			want: dto.QuestionnaireConnection{
				PageInfo: dto.PageInfo{
					HasNextPage:     true,
					HasPreviousPage: true,
				},
				TotalCount: 10,
				Edges: []dto.QuestionnaireEdge{
					{
						Node: dto.Questionnaire{
							ID:  questionnaireID,
							URL: (*scalarutils.URI)(&url),
						},
						Cursor: questionnaireID,
					},
					{
						Node: dto.Questionnaire{
							ID:  questionnaireID,
							URL: (*scalarutils.URI)(&url),
						},
						Cursor: questionnaireID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dto.CreateQuestionnaireConnection(tt.args.questionnaires, tt.args.pageInfo, tt.args.total); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateQuestionnaireConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
