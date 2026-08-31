package domain

import "testing"

func TestDaysOfWeek_Code(t *testing.T) {
	tests := []struct {
		name string
		code DaysOfWeek
		want string
	}{
		{
			name: "Monday",
			code: DaysOfWeekMon,
			want: "mon",
		},
		{
			name: "Tuesday",
			code: DaysOfWeekTue,
			want: "tue",
		},
		{
			name: "Wednesday",
			code: DaysOfWeekWed,
			want: "wed",
		},
		{
			name: "Thursday",
			code: DaysOfWeekThu,
			want: "thu",
		},
		{
			name: "Friday",
			code: DaysOfWeekFri,
			want: "fri",
		},
		{
			name: "Saturday",
			code: DaysOfWeekSat,
			want: "sat",
		},
		{
			name: "Sunday",
			code: DaysOfWeekSun,
			want: "sun",
		},
		{
			name: "Unknown",
			code: DaysOfWeek(99), // An invalid value
			want: "<unknown>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Code(); got != tt.want {
				t.Errorf("DaysOfWeek.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRLocationMode(t *testing.T) {
	tests := []struct {
		name             string
		locationModeCode FHIRLocationMode
		wantCode         string
		wantDisplay      string
		wantDefinition   string
		input            string
		wantJSON         string
		wantErr          bool
	}{
		{
			name:             "Instance",
			locationModeCode: FHIRLocationModeInstance,
			wantCode:         "instance",
			wantDisplay:      "Instance",
			wantDefinition:   "The Location resource represents a specific instance of a location (e.g. Operating Theatre 1A).",
			input:            `"instance"`,
			wantJSON:         `"instance"`,
			wantErr:          false,
		},
		{
			name:             "Kind",
			locationModeCode: FHIRLocationModeKind,
			wantCode:         "kind",
			wantDisplay:      "Kind",
			wantDefinition:   "The Location represents a class of locations (e.g. Any Operating Theatre) although this class of locations could be constrained within a specific boundary (such as organization, or parent location, address etc.).",
			input:            `"kind"`,
			wantJSON:         `"kind"`,
			wantErr:          false,
		},
		{
			name:             "Unknown",
			locationModeCode: FHIRLocationMode("invalid"),
			wantCode:         "<unknown>",
			wantDisplay:      "<unknown>",
			wantDefinition:   "<unknown>",
			input:            `"invalid"`,
			wantJSON:         `"\u003cunknown\u003e"`,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.locationModeCode.Code(); got != tt.wantCode {
				t.Errorf("FHIRLocationMode.Code() = %v, want %v", got, tt.wantCode)
			}

			if got := tt.locationModeCode.Display(); got != tt.wantDisplay {
				t.Errorf("FHIRLocationMode.Display() = %v, want %v", got, tt.wantDisplay)
			}

			if got := tt.locationModeCode.Definition(); got != tt.wantDefinition {
				t.Errorf("FHIRLocationMode.Definition() = %v, want %v", got, tt.wantDefinition)
			}

			jsonBytes, err := tt.locationModeCode.MarshalJSON()
			if err != nil {
				t.Errorf("FHIRLocationMode.MarshalJSON() error = %v", err)
			}

			if got := string(jsonBytes); got != tt.wantJSON {
				t.Errorf("FHIRLocationMode.MarshalJSON() = %v, want %v", got, tt.wantJSON)
			}

			var unmarshalledCode FHIRLocationMode
			err = unmarshalledCode.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Errorf("FHIRLocationMode.UnmarshalJSON() expected error, but got none")
				}

			} else {
				if err != nil {
					t.Errorf("FHIRLocationMode.UnmarshalJSON() error = %v", err)
				}
				if unmarshalledCode != tt.locationModeCode {
					t.Errorf("FHIRLocationMode.UnmarshalJSON() = %v, want %v", unmarshalledCode, tt.locationModeCode)
				}
			}
		})
	}
}

func TestFHIRLocationStatus(t *testing.T) {
	tests := []struct {
		name           string
		locationStatus FHIRLocationStatus
		wantCode       string
		wantDisplay    string
		wantDefinition string
		jsonInput      string
		wantJSON       string
		wantError      bool
	}{
		{
			name:           "Active",
			locationStatus: FHIRLocationStatusActive,
			wantCode:       "active",
			wantDisplay:    "Active",
			wantDefinition: "The location is operational.",
			jsonInput:      `"active"`,
			wantJSON:       `"active"`,
			wantError:      false,
		},
		{
			name:           "Suspended",
			locationStatus: FHIRLocationStatusSuspended,
			wantCode:       "suspended",
			wantDisplay:    "Suspended",
			wantDefinition: "The location is temporarily closed.",
			jsonInput:      `"suspended"`,
			wantJSON:       `"suspended"`,
			wantError:      false,
		},
		{
			name:           "Inactive",
			locationStatus: FHIRLocationStatusInactive,
			wantCode:       "inactive",
			wantDisplay:    "Inactive",
			wantDefinition: "The location is no longer used.",
			jsonInput:      `"inactive"`,
			wantJSON:       `"inactive"`,
			wantError:      false,
		},
		{
			name:           "Unknown",
			locationStatus: FHIRLocationStatus("invalid"),
			wantCode:       "<unknown>",
			wantDisplay:    "<unknown>",
			wantDefinition: "<unknown>",
			jsonInput:      `"invalid"`,
			wantJSON:       `"\u003cunknown\u003e"`,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.locationStatus.Code(); got != tt.wantCode {
				t.Errorf("FHIRLocationStatus.Code() = %v, want %v", got, tt.wantCode)
			}

			if got := tt.locationStatus.Display(); got != tt.wantDisplay {
				t.Errorf("FHIRLocationStatus.Display() = %v, want %v", got, tt.wantDisplay)
			}

			if got := tt.locationStatus.Definition(); got != tt.wantDefinition {
				t.Errorf("FHIRLocationStatus.Definition() = %v, want %v", got, tt.wantDefinition)
			}

			jsonBytes, err := tt.locationStatus.MarshalJSON()
			if err != nil {
				t.Errorf("FHIRLocationStatus.MarshalJSON() error = %v", err)
			}
			if got := string(jsonBytes); got != tt.wantJSON {
				t.Errorf("FHIRLocationStatus.MarshalJSON() = %v, want %v", got, tt.wantJSON)
			}

			var unmarshalledCode FHIRLocationStatus
			err = unmarshalledCode.UnmarshalJSON([]byte(tt.jsonInput))
			if tt.wantError {
				if err == nil {
					t.Errorf("FHIRLocationStatus.UnmarshalJSON() expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("FHIRLocationStatus.UnmarshalJSON() error = %v", err)
				}
				if unmarshalledCode != tt.locationStatus {
					t.Errorf("FHIRLocationStatus.UnmarshalJSON() = %v, want %v", unmarshalledCode, tt.locationStatus)
				}
			}
		})
	}
}
