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
	feature_v1 "github.com/zitadel/zitadel/internal/repository/feature"
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
			name: "filter error",
			eventstore: expectEventstore(
				expectFilterError(io.ErrClosedPipe),
			),
			args: args{ctx, &InstanceFeatures{
				LoginDefaultOrg: gu.Ptr(true),
			}},
			wantErr: io.ErrClosedPipe,
		},
		{
			name:       "all nil, No Change",
			eventstore: expectEventstore(),
			args:       args{ctx, &InstanceFeatures{}},
			wantErr:    zerrors.ThrowInvalidArgument(nil, "COMMAND-Vigh1", "Errors.NoChangesFound"),
		},
		{
			name: "set LoginDefaultOrg",
			eventstore: expectEventstore(
				expectFilter(),
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, true,
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
			name: "set LoginDefaultOrg, update from v1",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v1.NewSetEvent[feature_v1.Boolean](
						ctx, &eventstore.Aggregate{
							ID:            "instance1",
							ResourceOwner: "instance1",
						},
						feature_v1.DefaultLoginInstanceEventType,
						feature_v1.Boolean{
							Boolean: false,
						},
					)),
				),
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, true,
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
				expectFilter(),
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
				expectFilter(),
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
			name: "set UserSchema",
			eventstore: expectEventstore(
				expectFilter(),
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceUserSchemaEventType, true,
					),
				),
			),
			args: args{ctx, &InstanceFeatures{
				UserSchema: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "set Actions",
			eventstore: expectEventstore(
				expectFilter(),
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceActionsEventType, true,
					),
				),
			),
			args: args{ctx, &InstanceFeatures{
				Actions: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "push error",
			eventstore: expectEventstore(
				expectFilter(),
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
				expectFilter(),
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceTriggerIntrospectionProjectionsEventType, false,
					),
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLegacyIntrospectionEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceUserSchemaEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceActionsEventType, true,
					),
				),
			),
			args: args{ctx, &InstanceFeatures{
				LoginDefaultOrg:                 gu.Ptr(true),
				TriggerIntrospectionProjections: gu.Ptr(false),
				LegacyIntrospection:             gu.Ptr(true),
				UserSchema:                      gu.Ptr(true),
				Actions:                         gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "set only updated",
			eventstore: expectEventstore(
				// throw in some set events, reset and set again.
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceTriggerIntrospectionProjectionsEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewResetEvent(
						ctx, aggregate,
						feature_v2.InstanceResetEventType,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLegacyIntrospectionEventType, true,
					)),
				),
				expectPush(
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceTriggerIntrospectionProjectionsEventType, false,
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
			name: "filter error",
			eventstore: expectEventstore(
				expectFilterError(io.ErrClosedPipe),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "push error",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, true,
					)),
				),
				expectPushFailed(io.ErrClosedPipe,
					feature_v2.NewResetEvent(ctx, aggregate, feature_v2.InstanceResetEventType),
				),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, true,
					)),
				),
				expectPush(
					feature_v2.NewResetEvent(ctx, aggregate, feature_v2.InstanceResetEventType),
				),
			),
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "no change after previous reset",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewResetEvent(
						ctx, aggregate,
						feature_v2.InstanceResetEventType,
					)),
				),
			),
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "no change without previous events",
			eventstore: expectEventstore(
				expectFilter(),
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
