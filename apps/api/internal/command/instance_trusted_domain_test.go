package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddTrustedDomain(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		trustedDomain string
	}
	type want struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty domain, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "",
			},
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-Stk21", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid domain (length), error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "my-very-endleeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeess-looooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo0ooooooooooooooooooooooong.domain.com",
			},
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-Stk21", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid domain (chars), error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "&.com",
			},
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-S3v3w", "Errors.Instance.Domain.InvalidCharacter"),
			},
		},
		{
			name: "domain already exists, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewTrustedDomainAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate, "domain.com"),
						),
					),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "domain.com",
			},
			want: want{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMA-hg42a", "Errors.Instance.Domain.AlreadyExists"),
			},
		},
		{
			name: "domain add ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewTrustedDomainAddedEvent(context.Background(),
							&instance.NewAggregate("instanceID").Aggregate, "domain.com"),
					),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "domain.com",
			},
			want: want{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.AddTrustedDomain(tt.args.ctx, tt.args.trustedDomain)
			assert.ErrorIs(t, err, tt.want.err)
			assertObjectDetails(t, tt.want.details, got)
		})
	}
}

func TestCommands_RemoveTrustedDomain(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		trustedDomain string
	}
	type want struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "domain empty string, error",
			fields: fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore { return &eventstore.Eventstore{} },
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: " ",
			},
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-ajAzwu", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "domain invalid character, error",
			fields: fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore { return &eventstore.Eventstore{} },
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "? ",
			},
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-lfs3Te", "Errors.Instance.Domain.InvalidCharacter"),
			},
		},
		{
			name: "domain length exceeded, error",
			fields: fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore { return &eventstore.Eventstore{} },
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "averylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongdomain",
			},
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-ajAzwu", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "domain does not exists, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "domain.com",
			},
			want: want{
				err: zerrors.ThrowNotFound(nil, "COMMA-de3z9", "Errors.Instance.Domain.NotFound"),
			},
		},
		{
			name: "domain remove ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewTrustedDomainAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate, "domain.com"),
						),
					),
					expectPush(
						instance.NewTrustedDomainRemovedEvent(context.Background(),
							&instance.NewAggregate("instanceID").Aggregate, "domain.com"),
					),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "domain.com",
			},
			want: want{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.RemoveTrustedDomain(tt.args.ctx, tt.args.trustedDomain)
			assert.ErrorIs(t, err, tt.want.err)
			assertObjectDetails(t, tt.want.details, got)
		})
	}
}
