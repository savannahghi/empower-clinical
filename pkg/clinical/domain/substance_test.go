package domain

import "testing"

func TestFHIRSubstanceStatus(t *testing.T) {
	tests := []struct {
		name            string
		substanceStatus FHIRSubstanceStatus
		wantCode        string
		wantDisplay     string
		jsonInput       string
		wantJSON        string
		wantError       bool
	}{
		{
			name:            "Active",
			substanceStatus: FHIRSubstanceStatusActive,
			wantCode:        "active",
			wantDisplay:     "Active",
			jsonInput:       `"active"`,
			wantJSON:        `"active"`,
			wantError:       false,
		},
		{
			name:            "Inactive",
			substanceStatus: FHIRSubstanceStatusInactive,
			wantCode:        "inactive",
			wantDisplay:     "Inactive",
			jsonInput:       `"inactive"`,
			wantJSON:        `"inactive"`,
			wantError:       false,
		},
		{
			name:            "EnteredInError",
			substanceStatus: FHIRSubstanceStatusEnteredInError,
			wantCode:        "entered-in-error",
			wantDisplay:     "Entered in Error",
			jsonInput:       `"entered-in-error"`,
			wantJSON:        `"entered-in-error"`,
			wantError:       false,
		},
		{
			name:            "Unknown",
			substanceStatus: FHIRSubstanceStatus("invalid"),
			wantCode:        "<unknown>",
			wantDisplay:     "<unknown>",
			jsonInput:       `"invalid"`,
			wantJSON:        `"\u003cunknown\u003e"`,
			wantError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.substanceStatus.Code(); got != tt.wantCode {
				t.Errorf("FHIRSubstanceStatus.Code() = %v, want %v", got, tt.wantCode)
			}

			if got := tt.substanceStatus.Display(); got != tt.wantDisplay {
				t.Errorf("FHIRSubstanceStatus.Display() = %v, want %v", got, tt.wantDisplay)
			}

			jsonBytes, err := tt.substanceStatus.MarshalJSON()
			if err != nil {
				t.Errorf("FHIRSubstanceStatus.MarshalJSON() error = %v", err)
			}

			if got := string(jsonBytes); got != tt.wantJSON {
				t.Errorf("FHIRSubstanceStatus.MarshalJSON() = %v, want %v", got, tt.wantJSON)
			}

			var unmarshalledCode FHIRSubstanceStatus
			err = unmarshalledCode.UnmarshalJSON([]byte(tt.jsonInput))
			if tt.wantError {
				if err == nil {
					t.Errorf("FHIRSubstanceStatus.UnmarshalJSON() expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("FHIRSubstanceStatus.UnmarshalJSON() error = %v", err)
				}
				if unmarshalledCode != tt.substanceStatus {
					t.Errorf("FHIRSubstanceStatus.UnmarshalJSON() = %v, want %v", unmarshalledCode, tt.substanceStatus)
				}
			}
		})
	}
}
