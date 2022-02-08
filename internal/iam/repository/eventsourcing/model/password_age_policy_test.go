package model

import (
	"testing"
)

func TestPasswordAgePolicyChanges(t *testing.T) {
	type args struct {
		existing *PasswordAgePolicy
		new      *PasswordAgePolicy
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
			name: "age policy all attributes change",
			args: args{
				existing: &PasswordAgePolicy{MaxAgeDays: 365, ExpireWarnDays: 5},
				new:      &PasswordAgePolicy{MaxAgeDays: 730, ExpireWarnDays: 10},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &PasswordAgePolicy{MaxAgeDays: 10, ExpireWarnDays: 10},
				new:      &PasswordAgePolicy{MaxAgeDays: 10, ExpireWarnDays: 10},
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
