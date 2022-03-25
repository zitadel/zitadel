package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestAddDomain(t *testing.T) {
	type args struct {
		a      *org.Aggregate
		domain string
	}

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "invalid domain",
			args: args{
				a:      org.NewAggregate("test", "test"),
				domain: "",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "ORG-r3h4J", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:      org.NewAggregate("test", "test"),
				domain: "domain",
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
					org.NewDomainAddedEvent(context.Background(), &org.NewAggregate("test", "test").Aggregate, "domain"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparation.AssertValidation(t, AddOrgDomain(tt.args.a, tt.args.domain), nil, tt.want)
		})
	}
}

func TestVerifyDomain(t *testing.T) {
	type args struct {
		a      *org.Aggregate
		domain string
	}

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "invalid domain",
			args: args{
				a:      org.NewAggregate("test", "test"),
				domain: "",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "ORG-yqlVQ", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:      org.NewAggregate("test", "test"),
				domain: "domain",
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
					org.NewDomainVerifiedEvent(context.Background(), &org.NewAggregate("test", "test").Aggregate, "domain"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparation.AssertValidation(t, VerifyOrgDomain(tt.args.a, tt.args.domain), nil, tt.want)
		})
	}
}

func TestSetDomainPrimary(t *testing.T) {
	type args struct {
		a      *org.Aggregate
		domain string
	}

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "invalid domain",
			args: args{
				a:      org.NewAggregate("test", "test"),
				domain: "",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "ORG-gmNqY", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:      org.NewAggregate("test", "test"),
				domain: "domain",
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
					org.NewDomainPrimarySetEvent(context.Background(), &org.NewAggregate("test", "test").Aggregate, "domain"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparation.AssertValidation(t, SetPrimaryOrgDomain(tt.args.a, tt.args.domain), nil, tt.want)
		})
	}
}
