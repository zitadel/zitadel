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
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestLimits_SetLimits(t *testing.T) {
	type fields func(*testing.T) (*eventstore.Eventstore, id_generator.Generator)
	type args struct {
		ctx       context.Context
		setLimits *SetLimits
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
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
			name: "update limits audit log retention, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
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
										&limits.NewAggregate("limits1", "instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
			name: "update limits unblock, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Minute)),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimits: &SetLimits{
					Block: gu.Ptr(false),
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
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
							eventFromEventPusher(
								limits.NewResetEvent(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1").Aggregate,
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits2", "instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
			_eventstore, _id_generator := tt.fields(t)
			r := &Commands{eventstore: _eventstore}
			id_generator.SetGenerator(_id_generator)
			got, err := r.SetLimits(tt.args.ctx, tt.args.setLimits)
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

func TestLimits_SetLimitsBulk(t *testing.T) {
	type fields func(*testing.T) (*eventstore.Eventstore, id_generator.Generator)
	type args struct {
		ctx           context.Context
		setLimitsBulk []*SetInstanceLimitsBulk
	}
	type res struct {
		want       *domain.ObjectDetails
		wantTarget []*domain.ObjectDetails
		err        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "create limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimitsBulk: []*SetInstanceLimitsBulk{{
					InstanceID: "instance1",
					SetLimits: SetLimits{
						AuditLogRetention: gu.Ptr(time.Hour),
					},
				}},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "instance1",
				}},
			},
		},
		{
			name: "update limits audit log retention, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
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
										&limits.NewAggregate("limits1", "instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimitsBulk: []*SetInstanceLimitsBulk{{
					InstanceID: "instance1",
					SetLimits: SetLimits{
						AuditLogRetention: gu.Ptr(time.Hour),
					},
				}},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "instance1",
				}},
			},
		},
		{
			name: "update limits unblock, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Minute)),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimitsBulk: []*SetInstanceLimitsBulk{{
					InstanceID: "instance1",
					SetLimits: SetLimits{
						Block: gu.Ptr(false),
					},
				}},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "instance1",
				}},
			},
		},
		{
			name: "set limits after resetting limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
							eventFromEventPusher(
								limits.NewResetEvent(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1").Aggregate,
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits2", "instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimitsBulk: []*SetInstanceLimitsBulk{{
					InstanceID: "instance1",
					SetLimits: SetLimits{
						AuditLogRetention: gu.Ptr(time.Hour),
					},
				},
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "instance1",
				}},
			},
		},
		{
			name: "set many limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("add-block", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("switch-to-block", "instance2").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("already-blocked", "instance3").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("unblock", "instance4").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("relimit", "instance5").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusher(
								limits.NewResetEvent(
									context.Background(),
									&limits.NewAggregate("relimit", "instance5").Aggregate,
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance0",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("create-limits", "instance0").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("add-block", "instance1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusherWithInstanceID(
								"instance2",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("switch-to-block", "instance2").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusherWithInstanceID(
								"instance4",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("unblock", "instance4").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
							eventFromEventPusherWithInstanceID(
								"instance5",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("relimit", "instance5").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "create-limits", "relimit")
			},
			args: args{
				ctx: context.Background(),
				setLimitsBulk: []*SetInstanceLimitsBulk{{
					InstanceID: "instance0",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}, {
					InstanceID: "instance1",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}, {
					InstanceID: "instance2",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}, {
					InstanceID: "instance3",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}, {
					InstanceID: "instance4",
					SetLimits: SetLimits{
						Block: gu.Ptr(false),
					},
				}, {
					InstanceID: "instance5",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "instance0",
				}, {
					ResourceOwner: "instance1",
				}, {
					ResourceOwner: "instance2",
				}, {
					ResourceOwner: "instance3",
				}, {
					ResourceOwner: "instance4",
				}, {
					ResourceOwner: "instance5",
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_eventstore, _id_generator := tt.fields(t)
			r := &Commands{eventstore: _eventstore}
			id_generator.SetGenerator(_id_generator)
			gotDetails, gotTargetDetails, err := r.SetInstanceLimitsBulk(tt.args.ctx, tt.args.setLimitsBulk)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, gotDetails)
				assert.Equal(t, tt.res.wantTarget, gotTargetDetails)
			}
		})
	}
}

func TestLimits_ResetLimits(t *testing.T) {
	type fields func(*testing.T) *eventstore.Eventstore
	type args struct {
		ctx context.Context
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
									&limits.NewAggregate("limits1", "instance1").Aggregate,
									limits.SetEventType,
								),
								limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
							),
						),
						eventFromEventPusher(
							limits.NewResetEvent(context.Background(),
								&limits.NewAggregate("limits1", "instance1").Aggregate,
							),
						),
					),
				)
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
									&limits.NewAggregate("limits1", "instance1").Aggregate,
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
								&limits.NewAggregate("limits1", "instance1").Aggregate,
							),
						),
					),
				)
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
			got, err := r.ResetLimits(tt.args.ctx)
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
