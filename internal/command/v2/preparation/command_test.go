package preparation

import (
	"errors"
	"testing"
)

func Test_validate(t *testing.T) {
	type args struct {
		validations []Validation
	}
	type want struct {
		len int
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no validations",
		},
		{
			name: "validations return no command creators",
		},
		{
			name: "error in validation",
		},
		{
			name: "correct",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate(tt.args.validations)
			if !errors.Is(err, tt.want.err) {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if len(got) != tt.want.len {
				t.Errorf("validate() len = %v, want %v", len(got), tt.want.len)
			}
		})
	}
}
