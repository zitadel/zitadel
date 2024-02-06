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
			name:       "all nil, No Change",
			eventstore: expectEventstore(),
			args:       args{context.Background(), &SystemFeatures{}},
			wantErr:    zerrors.ThrowInvalidArgument(nil, "COMMAND-Dul8I", "Errors.NoChangesFound"),
		},
		{
			name: "set LoginDefaultOrg",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemDefaultLoginInstanceEventType, true,
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
			name: "set TriggerIntrospectionProjections",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				TriggerIntrospectionProjections: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
		{
			name: "set LegacyIntrospection",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, true,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				LegacyIntrospection: gu.Ptr(true),
			}},
			want: &domain.ObjectDetails{
				ResourceOwner: "SYSTEM",
			},
		},
		{
			name: "push error",
			eventstore: expectEventstore(
				expectPushFailed(io.ErrClosedPipe,
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, true,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				LegacyIntrospection: gu.Ptr(true),
			}},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "set all",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemDefaultLoginInstanceEventType, true,
					),
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, false,
					),
					feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, true,
					),
				),
			),
			args: args{context.Background(), &SystemFeatures{
				LoginDefaultOrg:                 gu.Ptr(true),
				TriggerIntrospectionProjections: gu.Ptr(false),
				LegacyIntrospection:             gu.Ptr(true),
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
			name: "push error",
			eventstore: expectEventstore(
				expectPushFailed(io.ErrClosedPipe,
					feature_v2.NewResetEvent(context.Background(), aggregate, feature_v2.SystemResetEventType),
				),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			eventstore: expectEventstore(
				expectPush(
					feature_v2.NewResetEvent(context.Background(), aggregate, feature_v2.SystemResetEventType),
				),
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
			assert.Equal(t, tt.want, got)
		})
	}
}
