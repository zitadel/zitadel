package model

import "testing"

func TestPhoneChanges(t *testing.T) {
	type args struct {
		existing *Phone
		new      *Phone
	}
	type res struct {
		changesLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "all fields changed",
			args: args{
				existing: &Phone{Phone: "Phone", IsPhoneVerified: true},
				new:      &Phone{Phone: "PhoneChanged", IsPhoneVerified: false},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no fields changed",
			args: args{
				existing: &Phone{Phone: "Phone", IsPhoneVerified: true},
				new:      &Phone{Phone: "Phone", IsPhoneVerified: false},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existing.Changes(tt.args.new)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}
