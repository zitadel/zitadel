package org

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestAddOrg(t *testing.T) {
	type args struct {
		a    *org.Aggregate
		name string
	}

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "invalid domain",
			args: args{
				a:    org.NewAggregate("test", "test"),
				name: "",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "ORG-mruNY", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:    org.NewAggregate("test", "test"),
				name: "domain",
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
					org.NewOrgAddedEvent(context.Background(), &org.NewAggregate("test", "test").Aggregate, "domain"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparation.AssertValidation(t, AddOrg(tt.args.a, tt.args.name), tt.want)
		})
	}
}
