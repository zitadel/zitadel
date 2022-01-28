package model

import (
	"testing"
)

func TestPasswordLockoutPolicyChanges(t *testing.T) {
	type args struct {
		existing *LockoutPolicy
		new      *LockoutPolicy
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
			name: "lockout policy all attributes change",
			args: args{
				existing: &LockoutPolicy{MaxPasswordAttempts: 365, ShowLockOutFailures: true},
				new:      &LockoutPolicy{MaxPasswordAttempts: 730, ShowLockOutFailures: false},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &LockoutPolicy{MaxPasswordAttempts: 10, ShowLockOutFailures: true},
				new:      &LockoutPolicy{MaxPasswordAttempts: 10, ShowLockOutFailures: true},
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
