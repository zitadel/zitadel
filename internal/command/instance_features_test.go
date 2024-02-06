package command

import (
	"context"
	"io"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_SetInstanceFeatures(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	aggregate := feature_v2.NewAggregate("instance1", "instance1")

	type args struct {
		ctx context.Context
		f   *InstanceFeatures
	}
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		args       args
		want       *domain.ObjectDetails
		wantErr    error
	}{
		{
			name:       "all nil, No Change",
			eventstore: expectEventstore(),
			args:       args{ctx, &InstanceFeatures{}},
			wantErr:    zerrors.ThrowInvalidArgument(nil, "COMMAND-Gie6U", "Errors.NoChangesFound"),
		},
		{
			name: "set LoginDefaultOrg",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceDefaultLoginInstanceEventType, true,
					),
				),
			),
			args: args{ctx, &InstanceFeatures{
				LoginDefaultOrg: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "set TriggerIntrospectionProjections",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceTriggerIntrospectionProjectionsEventType, true,
					),
				),
			),
			args: args{ctx, &InstanceFeatures{
				TriggerIntrospectionProjections: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "set LegacyIntrospection",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLegacyIntrospectionEventType, true,
					),
				),
			),
			args: args{ctx, &InstanceFeatures{
				LegacyIntrospection: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "push error",
			eventstore: expectEventstore(
				expectPushFailed(io.ErrClosedPipe,
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLegacyIntrospectionEventType, true,
					),
				),
			),
			args: args{ctx, &InstanceFeatures{
				LegacyIntrospection: gu.Ptr(true),
			}},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "set all",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceDefaultLoginInstanceEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceTriggerIntrospectionProjectionsEventType, false,
					),
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLegacyIntrospectionEventType, true,
					),
				),
			),
			args: args{ctx, &InstanceFeatures{
				LoginDefaultOrg:                 gu.Ptr(true),
				TriggerIntrospectionProjections: gu.Ptr(false),
				LegacyIntrospection:             gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.eventstore(t),
			}
			got, err := c.SetInstanceFeatures(tt.args.ctx, tt.args.f)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_ResetInstanceFeatures(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	aggregate := feature_v2.NewAggregate("instance1", "instance1")
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		want       *domain.ObjectDetails
		wantErr    error
	}{
		{
			name: "push error",
			eventstore: expectEventstore(
				expectPushFailed(io.ErrClosedPipe,
					feature_v2.NewResetEvent(ctx, aggregate, feature_v2.InstanceResetEventType),
				),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewResetEvent(ctx, aggregate, feature_v2.InstanceResetEventType),
				),
			),
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.eventstore(t),
			}
			got, err := c.ResetInstanceFeatures(ctx)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
