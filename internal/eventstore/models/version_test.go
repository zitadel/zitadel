package models

import "testing"

func TestVersion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       Version
		wantErr bool
	}{
		{
			"correct version",
			"v1.23.23",
			false,
		},
		{
			"no v prefix",
			"1.2.2",
			true,
		},
		{
			"letters in version",
			"v1.as.3",
			true,
		},
		{
			"no version",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Version.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
