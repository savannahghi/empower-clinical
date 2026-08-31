package dto

import (
	"testing"

	"github.com/brianvoe/gofakeit"
)

func TestCreateTaskConnection(t *testing.T) {
	type args struct {
		tasks    []*TaskOutput
		pageInfo PageInfo
		total    int
	}
	tests := []struct {
		name string
		args args
		want TaskOutputConnection
	}{
		{
			name: "Happy case: creates task connection",
			args: args{
				tasks: []*TaskOutput{
					{
						ID: gofakeit.UUID(),
					},
					{
						ID: gofakeit.UUID(),
					},
				},
				pageInfo: PageInfo{
					HasNextPage:     true,
					HasPreviousPage: true,
				},
				total: 10,
			},
			want: TaskOutputConnection{
				Edges: []TaskEdge{
					{
						Node: TaskOutput{ID: gofakeit.UUID()},
					},
					{
						Node: TaskOutput{ID: gofakeit.UUID()},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateTaskConnection(tt.args.tasks, tt.args.pageInfo, &tt.args.total)
			if len(got.Edges) != len(tt.want.Edges) {
				t.Errorf("CreateTaskConnection() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
