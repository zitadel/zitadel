package model

import "testing"

func TestEmailChanges(t *testing.T) {
	type args struct {
		existing *Email
		new      *Email
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
				existing: &Email{EmailAddress: "Email", IsEmailVerified: true},
				new:      &Email{EmailAddress: "EmailChanged", IsEmailVerified: false},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no fields changed",
			args: args{
				existing: &Email{EmailAddress: "Email", IsEmailVerified: true},
				new:      &Email{EmailAddress: "Email", IsEmailVerified: false},
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
