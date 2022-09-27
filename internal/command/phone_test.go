package command

import (
	"testing"

	"github.com/zitadel/zitadel/internal/errors"
)

func TestFormatPhoneNumber(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name    string
		args    args
		result  *Phone
		errFunc func(err error) bool
	}{
		{
			name: "invalid phone number",
			args: args{
				number: "PhoneNumber",
			},
			errFunc: errors.IsErrorInvalidArgument,
		},
		{
			name: "format phone +4171 xxx xx xx",
			args: args{
				number: "+4171 123 45 67",
			},
			result: &Phone{
				Number: "+41711234567",
			},
		},
		{
			name: "format non swiss phone +4371 xxx xx xx",
			args: args{
				number: "+4371 123 45 67",
			},
			result: &Phone{
				Number: "+43711234567",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := FormatPhoneNumber(tt.args.number)

			if tt.errFunc == nil && tt.result.Number != formatted {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.args.number, formatted)
			}
			if tt.errFunc != nil && !tt.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
