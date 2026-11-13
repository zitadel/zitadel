package routing

import "testing"

func TestSelect(t *testing.T) {
	tests := []struct {
		name                 string
		configured           []string
		active               []string
		externalLoginAllowed bool
		want                 Decision
	}{
		{name: "trusted single provider", configured: []string{"entra"}, active: []string{"entra", "okta"}, externalLoginAllowed: true, want: Decision{ProviderID: "entra", Allowed: true}},
		{name: "no mapping", active: []string{"entra"}, externalLoginAllowed: true},
		{name: "ambiguous mapping", configured: []string{"entra", "okta"}, active: []string{"entra", "okta"}, externalLoginAllowed: true},
		{name: "disabled provider", configured: []string{"entra"}, active: []string{"okta"}, externalLoginAllowed: true},
		{name: "external login disabled", configured: []string{"entra"}, active: []string{"entra"}, externalLoginAllowed: false},
		{name: "empty provider", configured: []string{""}, active: []string{""}, externalLoginAllowed: true},
		{name: "duplicate configured provider remains ambiguous", configured: []string{"entra", "entra"}, active: []string{"entra"}, externalLoginAllowed: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Select(tt.configured, tt.active, tt.externalLoginAllowed)
			if got != tt.want {
				t.Fatalf("Select() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
