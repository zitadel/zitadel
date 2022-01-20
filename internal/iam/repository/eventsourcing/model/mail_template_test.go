package model

import (
	"testing"
)

func TestMailTemplateChanges(t *testing.T) {
	type args struct {
		existing *MailTemplate
		new      *MailTemplate
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
			name: "mailtemplate all attributes change",
			args: args{
				existing: &MailTemplate{Template: []byte("<doctype html>")},
				new:      &MailTemplate{Template: []byte("<!doctype html>")},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &MailTemplate{Template: []byte("<!doctype html>")},
				new:      &MailTemplate{Template: []byte("<!doctype html>")},
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
