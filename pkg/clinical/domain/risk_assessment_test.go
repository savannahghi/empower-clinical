package domain

import "testing"

func TestProbabilityRiskAssessmentEnum_Display(t *testing.T) {
	tests := []struct {
		name string
		p    ProbabilityRiskAssessmentEnum
		want string
	}{
		{
			name: "Happy case: correct display method for high risk",
			p:    HighRiskProbability,
			want: "High likelihood",
		},
		{
			name: "Happy case: correct display method for low risk",
			p:    LowRiskRiskProbability,
			want: "Low likelihood",
		},
		{
			name: "Happy case: correct display method for moderate risk",
			p:    ModerateRiskProbability,
			want: "Moderate likelihood",
		},
		{
			name: "Happy case: correct display method for certain risk",
			p:    CertainRiskProbability,
			want: "Certain",
		},
		{
			name: "Happy case: correct display method for negligible risk",
			p:    NegligibleRiskProbability,
			want: "Negligible likelihood",
		},
		{
			name: "Sad case: unknown code",
			p:    "",
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Display(); got != tt.want {
				t.Errorf("ProbabilityRiskAssessmentEnum.Display() = %v, want %v", got, tt.want)
			}
		})
	}
}
