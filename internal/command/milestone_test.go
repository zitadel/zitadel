package command

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/connector/gomap"
	"github.com/zitadel/zitadel/internal/cache/connector/noop"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/milestone"
)

func TestCommands_GetMilestonesReached(t *testing.T) {
	cached := &MilestonesReached{
		InstanceID:                        "cached-id",
		InstanceCreated:                   true,
		AuthenticationSucceededOnInstance: true,
	}

	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	aggregate := milestone.NewAggregate(ctx)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *MilestonesReached
		wantErr error
	}{
		{
			name: "cached",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "cached-id"),
			},
			want: cached,
		},
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "no events, all false",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID: "instanceID",
			},
		},
		{
			name: "instance created",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.InstanceCreated)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:      "instanceID",
				InstanceCreated: true,
			},
		},
		{
			name: "instance auth",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:                        "instanceID",
				AuthenticationSucceededOnInstance: true,
			},
		},
		{
			name: "project created",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:     "instanceID",
				ProjectCreated: true,
			},
		},
		{
			name: "app created",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:         "instanceID",
				ApplicationCreated: true,
			},
		},
		{
			name: "app auth",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:                           "instanceID",
				AuthenticationSucceededOnApplication: true,
			},
		},
		{
			name: "instance deleted",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.InstanceDeleted)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:      "instanceID",
				InstanceDeleted: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := gomap.NewCache[milestoneIndex, string, *MilestonesReached](
				context.Background(),
				[]milestoneIndex{milestoneIndexInstanceID},
				cache.Config{Connector: cache.ConnectorMemory},
			)
			cache.Set(context.Background(), cached)

			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: cache,
				},
			}
			got, err := c.GetMilestonesReached(tt.args.ctx)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_milestonesCompleted(t *testing.T) {
	c := &Commands{
		caches: &Caches{
			milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
		},
	}
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	arg := &MilestonesReached{
		InstanceID:                           "instanceID",
		InstanceCreated:                      true,
		AuthenticationSucceededOnInstance:    true,
		ProjectCreated:                       true,
		ApplicationCreated:                   true,
		AuthenticationSucceededOnApplication: true,
		InstanceDeleted:                      false,
	}
	c.setCachedMilestonesReached(ctx, arg)
	got, ok := c.getCachedMilestonesReached(ctx)
	assert.True(t, ok)
	assert.Equal(t, arg, got)
}

func TestCommands_MilestonePushed(t *testing.T) {
	aggregate := milestone.NewInstanceAggregate("instanceID")
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		msType     milestone.Type
		endpoints  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "milestone pushed",
			fields: fields{
				eventstore: expectEventstore(
					expectPush(
						milestone.NewPushedEvent(
							context.Background(),
							aggregate,
							milestone.ApplicationCreated,
							[]string{"foo.com", "bar.com"},
							"example.com",
						),
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "instanceID",
				msType:     milestone.ApplicationCreated,
				endpoints:  []string{"foo.com", "bar.com"},
			},
			wantErr: nil,
		},
		{
			name: "pusher error",
			fields: fields{
				eventstore: expectEventstore(
					expectPushFailed(
						io.ErrClosedPipe,
						milestone.NewPushedEvent(
							context.Background(),
							aggregate,
							milestone.ApplicationCreated,
							[]string{"foo.com", "bar.com"},
							"example.com",
						),
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "instanceID",
				msType:     milestone.ApplicationCreated,
				endpoints:  []string{"foo.com", "bar.com"},
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:     tt.fields.eventstore(t),
				externalDomain: "example.com",
			}
			err := c.MilestonePushed(tt.args.ctx, tt.args.instanceID, tt.args.msType, tt.args.endpoints)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestOIDCSessionEvents_SetMilestones(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	ctx = authz.WithConsoleClientID(ctx, "console")
	aggregate := milestone.NewAggregate(ctx)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		clientID string
		isHuman  bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantEvents []eventstore.Command
		wantErr    error
	}{
		{
			name: "get error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  true,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "milestones already reached",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance)),
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication)),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  true,
			},
			wantErr: nil,
		},
		{
			name: "auth on instance",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "console",
				isHuman:  true,
			},
			wantEvents: []eventstore.Command{
				milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance),
			},
			wantErr: nil,
		},
		{
			name: "subsequent console login, no milestone",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance)),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "console",
				isHuman:  true,
			},
			wantErr: nil,
		},
		{
			name: "subsequent machine login, no milestone",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance)),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  false,
			},
			wantErr: nil,
		},
		{
			name: "auth on app",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance)),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  true,
			},
			wantEvents: []eventstore.Command{
				milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
				},
			}
			s := &OIDCSessionEvents{
				commands: c,
			}
			postCommit, err := s.SetMilestones(tt.args.ctx, tt.args.clientID, tt.args.isHuman)
			postCommit(tt.args.ctx)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantEvents, s.events)
		})
	}
}

func TestCommands_projectCreatedMilestone(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	systemCtx := authz.SetCtxData(ctx, authz.CtxData{
		SystemMemberships: authz.Memberships{
			&authz.Membership{
				MemberType: authz.MemberTypeSystem,
			},
		},
	})
	aggregate := milestone.NewAggregate(ctx)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantEvents []eventstore.Command
		wantErr    error
	}{
		{
			name: "system user",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: systemCtx,
			},
			wantErr: nil,
		},
		{
			name: "get error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "milestone already reached",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: nil,
		},
		{
			name: "milestone reached event",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantEvents: []eventstore.Command{
				milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
				},
			}
			var cmds []eventstore.Command
			postCommit, err := c.projectCreatedMilestone(tt.args.ctx, &cmds)
			postCommit(tt.args.ctx)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantEvents, cmds)
		})
	}
}

func TestCommands_applicationCreatedMilestone(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	systemCtx := authz.SetCtxData(ctx, authz.CtxData{
		SystemMemberships: authz.Memberships{
			&authz.Membership{
				MemberType: authz.MemberTypeSystem,
			},
		},
	})
	aggregate := milestone.NewAggregate(ctx)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantEvents []eventstore.Command
		wantErr    error
	}{
		{
			name: "system user",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: systemCtx,
			},
			wantErr: nil,
		},
		{
			name: "get error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "milestone already reached",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: nil,
		},
		{
			name: "milestone reached event",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantEvents: []eventstore.Command{
				milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
				},
			}
			var cmds []eventstore.Command
			postCommit, err := c.applicationCreatedMilestone(tt.args.ctx, &cmds)
			postCommit(tt.args.ctx)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantEvents, cmds)
		})
	}
}

func (c *Commands) setMilestonesCompletedForTest(instanceID string) {
	c.milestonesCompleted.Store(instanceID, struct{}{})
}
