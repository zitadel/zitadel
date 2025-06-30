package command

import (
	"context"
	"io"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_SetSystemFeatures(t *testing.T) {
	aggregate := feature_v2.NewAggregate("SYSTEM", "SYSTEM")

	type args struct {
		ctx context.Context
		f   *SystemFeatures
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
			args: args{context.Background(), &SystemFeatures{
				LoginDefaultOrg: gu.Ptr(true),
			}},
			wantErr: io.ErrClosedPipe,
		},
		{
			name:       "all nil, No Change",
			eventstore: expectEventstore(),
			args:       args{context.Background(), &SystemFeatures{}},
			wantErr:    zerrors.ThrowInternal(nil, "COMMAND-Oop8a", "Errors.NoChangesFound"),
		},
		{
			name: "set LoginDefaultOrg",
			eventstore: expectEventstore(
				expectFilter(),
				expectPush(
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, true,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				LoginDefaultOrg: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
		{
			name: "set UserSchema",
			eventstore: expectEventstore(
				expectFilter(),
				expectPush(
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemUserSchemaEventType, true,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				UserSchema: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
		{
			name: "push error",
			eventstore: expectEventstore(
				expectFilter(),
				expectPushFailed(io.ErrClosedPipe,
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemEnableBackChannelLogout, true,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				EnableBackChannelLogout: gu.Ptr(true),
			}},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "set all",
			eventstore: expectEventstore(
				expectFilter(),
				expectPush(
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemUserSchemaEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemOIDCSingleV1SessionTerminationEventType, true,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				LoginDefaultOrg:                gu.Ptr(true),
				UserSchema:                     gu.Ptr(true),
				OIDCSingleV1SessionTermination: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
		{
			name: "set only updated",
			eventstore: expectEventstore(
				// throw in some set events, reset and set again.
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewResetEvent(
						context.Background(), aggregate,
						feature_v2.SystemResetEventType,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, false,
					)),
				),
				expectPush(
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemUserSchemaEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemOIDCSingleV1SessionTerminationEventType, false,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				LoginDefaultOrg:                gu.Ptr(true),
				UserSchema:                     gu.Ptr(true),
				OIDCSingleV1SessionTermination: gu.Ptr(false),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.eventstore(t),
			}
			got, err := c.SetSystemFeatures(tt.args.ctx, tt.args.f)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_ResetSystemFeatures(t *testing.T) {
	aggregate := feature_v2.NewAggregate("SYSTEM", "SYSTEM")
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
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, true,
					)),
				),
				expectPushFailed(io.ErrClosedPipe,
					feature_v2.NewResetEvent(context.Background(), aggregate, feature_v2.SystemResetEventType),
				),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, true,
					)),
				),
				expectPush(
					feature_v2.NewResetEvent(context.Background(), aggregate, feature_v2.SystemResetEventType),
				),
			),
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
		{
			name: "no change after previous reset",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewResetEvent(
						context.Background(), aggregate,
						feature_v2.SystemResetEventType,
					)),
				),
			),
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
		{
			name: "no change without previous events",
			eventstore: expectEventstore(
				expectFilter(),
			),
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.eventstore(t),
			}
			got, err := c.ResetSystemFeatures(context.Background())
			require.ErrorIs(t, err, tt.wantErr)
			assertObjectDetails(t, tt.want, got)
		})
	}
}
