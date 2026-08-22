package block

import "testing"

func TestFirebaseImpersonationEnabled(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		allow       string
		want        bool
	}{
		{name: "exact local true", environment: "local", allow: "true", want: true},
		{name: "case and whitespace tolerated", environment: " LOCAL ", allow: " TRUE ", want: true},
		{name: "production denied", environment: "production", allow: "true", want: false},
		{name: "missing environment denied", environment: "", allow: "true", want: false},
		{name: "explicit false denied", environment: "local", allow: "false", want: false},
		{name: "missing allow denied", environment: "local", allow: "", want: false},
		{name: "numeric truthy denied", environment: "local", allow: "1", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firebaseImpersonationEnabled(tt.environment, tt.allow); got != tt.want {
				t.Fatalf("firebaseImpersonationEnabled(%q, %q) = %t, want %t", tt.environment, tt.allow, got, tt.want)
			}
		})
	}
}
