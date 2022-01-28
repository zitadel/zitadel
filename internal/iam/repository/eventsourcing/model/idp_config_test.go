package model

import (
	"testing"
)

func TestIdpConfigChanges(t *testing.T) {
	type args struct {
		existing *IDPConfig
		new      *IDPConfig
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
			name: "idp config name changes",
			args: args{
				existing: &IDPConfig{IDPConfigID: "IDPConfigID", Name: "Name"},
				new:      &IDPConfig{IDPConfigID: "IDPConfigID", Name: "NameChanged"},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &IDPConfig{IDPConfigID: "IDPConfigID", Name: "Name"},
				new:      &IDPConfig{IDPConfigID: "IDPConfigID", Name: "Name"},
			},
			res: res{
				changesLen: 1,
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
