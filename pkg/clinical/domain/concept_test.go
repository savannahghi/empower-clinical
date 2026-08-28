package domain

import "testing"

func TestConcept_GetConceptDescription(t *testing.T) {
	type fields struct {
		DisplayName  string
		Descriptions []*Descriptions
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy case: get description",
			fields: fields{
				DisplayName: "Test Description",
				Descriptions: []*Descriptions{
					{
						Description: "Test Description",
					},
				},
			},
			want: "Test Description",
		},
		{
			name: "Happy case: get display name",
			fields: fields{
				DisplayName: "Test Display name",
			},
			want: "Test Display name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Concept{
				DisplayName:  tt.fields.DisplayName,
				Descriptions: tt.fields.Descriptions,
			}

			if got := d.GetConceptDisplay(); got != tt.want {
				t.Errorf("Concept.GetConceptDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}
