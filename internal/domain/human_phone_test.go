package domain

import (
	"testing"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

func TestFormatPhoneNumber(t *testing.T) {
	type args struct {
		phone *Phone
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
				phone: &Phone{
					PhoneNumber: "PhoneNumber",
				},
			},
			errFunc: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "format phone 071...",
			args: args{
				phone: &Phone{
					PhoneNumber: "0711234567",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format phone 0041...",
			args: args{
				phone: &Phone{
					PhoneNumber: "0041711234567",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format phone 071 xxx xx xx",
			args: args{
				phone: &Phone{
					PhoneNumber: "071 123 45 67",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format phone +4171 xxx xx xx",
			args: args{
				phone: &Phone{
					PhoneNumber: "+4171 123 45 67",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format phone 004171 xxx xx xx",
			args: args{
				phone: &Phone{
					PhoneNumber: "004171 123 45 67",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format non swiss phone 004371 xxx xx xx",
			args: args{
				phone: &Phone{
					PhoneNumber: "004371 123 45 67",
				},
			},
			result: &Phone{
				PhoneNumber: "+43711234567",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.phone.formatPhone()

			if tt.errFunc == nil && tt.result.PhoneNumber != tt.args.phone.PhoneNumber {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.args.phone.PhoneNumber, tt.result.PhoneNumber)
			}
			if tt.errFunc != nil && !tt.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
