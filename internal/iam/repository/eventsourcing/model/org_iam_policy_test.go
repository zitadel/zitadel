package model

import (
	"testing"
)

func TestOrgIAMPolicyChanges(t *testing.T) {
	type args struct {
		existing *DomainPolicy
		new      *DomainPolicy
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
			name: "org iam policy all attributes change",
			args: args{
				existing: &DomainPolicy{UserLoginMustBeDomain: true},
				new:      &DomainPolicy{UserLoginMustBeDomain: false},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &DomainPolicy{UserLoginMustBeDomain: true},
				new:      &DomainPolicy{UserLoginMustBeDomain: true},
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
