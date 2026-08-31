package domain

import "testing"

func TestMedicationDispenseStatusCodes(t *testing.T) {
	tests := []struct {
		name string
		code MedicationDispenseStatusCodes
		want string
	}{
		{
			name: "preparation",
			code: MedicationDispenseStatusCodesPreparation,
			want: "preparation",
		},
		{
			name: "in-progress",
			code: MedicationDispenseStatusCodesInProgress,
			want: "in-progress",
		},
		{
			name: "cancelled",
			code: MedicationDispenseStatusCodesCancelled,
			want: "cancelled",
		},
		{
			name: "on-hold",
			code: MedicationDispenseStatusCodesOnHold,
			want: "on-hold",
		},
		{
			name: "completed",
			code: MedicationDispenseStatusCodesCompleted,
			want: "completed",
		},
		{
			name: "entered-in-error",
			code: MedicationDispenseStatusCodesEnteredInError,
			want: "entered-in-error",
		},
		{
			name: "stopped",
			code: MedicationDispenseStatusCodesStopped,
			want: "stopped",
		},
		{
			name: "declined",
			code: MedicationDispenseStatusCodesDeclined,
			want: "declined",
		},
		{
			name: "INVALID",
			code: "MedicationDispenseStatusCodesDeclined",
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Code(); got != tt.want {
				t.Errorf("MedicationDispenseStatusCodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedicationDispenseStatusDisplay(t *testing.T) {
	tests := []struct {
		name string
		code MedicationDispenseStatusCodes
		want string
	}{
		{
			name: "Preparation",
			code: MedicationDispenseStatusCodesPreparation,
			want: "Preparation",
		},
		{
			name: "In Progress",
			code: MedicationDispenseStatusCodesInProgress,
			want: "In Progress",
		},
		{
			name: "Cancelled",
			code: MedicationDispenseStatusCodesCancelled,
			want: "Cancelled",
		},
		{
			name: "On Hold",
			code: MedicationDispenseStatusCodesOnHold,
			want: "On Hold",
		},
		{
			name: "Completed",
			code: MedicationDispenseStatusCodesCompleted,
			want: "Completed",
		},
		{
			name: "Entered in Error",
			code: MedicationDispenseStatusCodesEnteredInError,
			want: "Entered in Error",
		},
		{
			name: "Stopped",
			code: MedicationDispenseStatusCodesStopped,
			want: "Stopped",
		},
		{
			name: "Declined",
			code: MedicationDispenseStatusCodesDeclined,
			want: "Declined",
		},
		{
			name: "Unknown",
			code: MedicationDispenseStatusCodesUnknown,
			want: "Unknown",
		},
		{
			name: "Invalid",
			code: "MedicationDispenseStatusCodesDeclined",
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Display(); got != tt.want {
				t.Errorf("MedicationDispenseStatusDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}
