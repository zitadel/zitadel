package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestLimits_SetLimits(t *testing.T) {
	type fields func(*testing.T) (*eventstore.Eventstore, id.Generator)
	type args struct {
		ctx           context.Context
		resourceOwner string
		setLimits     *SetLimits
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "create limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "limits1")
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
				setLimits: &SetLimits{
					AuditLogRetention: gu.Ptr(time.Hour),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
		{
			name: "update limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Minute)),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
				setLimits: &SetLimits{
					AuditLogRetention: gu.Ptr(time.Hour),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
		{
			name: "set limits after resetting limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
							eventFromEventPusher(
								limits.NewResetEvent(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits2", "instance1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "limits2")
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
				setLimits: &SetLimits{
					AuditLogRetention: gu.Ptr(time.Hour),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(Commands)
			r.eventstore, r.idGenerator = tt.fields(t)
			got, err := r.SetLimits(tt.args.ctx, tt.args.resourceOwner, tt.args.setLimits)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestLimits_ResetLimits(t *testing.T) {
	type fields func(*testing.T) *eventstore.Eventstore
	type args struct {
		ctx           context.Context
		resourceOwner string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "not found",
			fields: func(tt *testing.T) *eventstore.Eventstore {
				return eventstoreExpect(
					tt,
					expectFilter(),
				)
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-9JToT", "Errors.Limits.NotFound"))
				},
			},
		},
		{
			name: "already removed",
			fields: func(tt *testing.T) *eventstore.Eventstore {
				return eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							limits.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
									limits.SetEventType,
								),
								limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
							),
						),
						eventFromEventPusher(
							limits.NewResetEvent(context.Background(),
								&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
							),
						),
					),
				)
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-9JToT", "Errors.Limits.NotFound"))
				},
			},
		},
		{
			name: "reset limits, ok",
			fields: func(tt *testing.T) *eventstore.Eventstore {
				return eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							limits.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
									limits.SetEventType,
								),
								limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
							),
						),
					),
					expectPush(
						eventFromEventPusherWithInstanceID(
							"instance1",
							limits.NewResetEvent(context.Background(),
								&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
							),
						),
					),
				)
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields(t),
			}
			got, err := r.ResetLimits(tt.args.ctx, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
