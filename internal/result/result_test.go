package result

import "testing"

func TestStatus_String(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{Success, "success"},
		{Failure, "failure"},
		{Error, "error"},
		{Status(99), "unknown"},
	}

	for _, tt := range tests {
		got := tt.status.String()
		if got != tt.want {
			t.Errorf("Status(%d).String() = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestResult_Passed(t *testing.T) {
	tests := []struct {
		status Status
		want   bool
	}{
		{Success, true},
		{Failure, false},
		{Error, false},
	}

	for _, tt := range tests {
		r := Result{Status: tt.status}
		if got := r.Passed(); got != tt.want {
			t.Errorf("Result{Status: %s}.Passed() = %v, want %v", tt.status, got, tt.want)
		}
	}
}
