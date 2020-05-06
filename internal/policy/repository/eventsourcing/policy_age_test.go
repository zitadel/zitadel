package eventsourcing

import (
	"testing"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

func TestGetPasswordAgePolicyQuery(t *testing.T) {
	type args struct {
		recourceOwner string
		sequence      uint64
	}
	type res struct {
		filterLen int
		wantErr   bool
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "Get password age policy query ok",
			args: args{
				recourceOwner: "org",
				sequence:      14,
			},
			res: res{
				filterLen: 3,
			},
		},
		{
			name: "Get password age policy query, no org",
			args: args{
				sequence: 1,
			},
			res: res{
				filterLen: 3,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := PasswordAgePolicyQuery(tt.args.recourceOwner, tt.args.sequence)
			if !tt.res.wantErr && query == nil {
				t.Errorf("query should not be nil")
			}
			if !tt.res.wantErr && len(query.Filters) != tt.res.filterLen {
				t.Errorf("got wrong filter len: expected: %v, actual: %v ", tt.res.filterLen, len(query.Filters))
			}
		})
	}
}
