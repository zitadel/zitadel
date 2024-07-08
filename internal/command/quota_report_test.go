package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

func TestQuotaReport_ReportQuotaUsage(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx              context.Context
		dueNotifications []*quota.NotificationDueEvent
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "no due events",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
			},
			res: res{},
		},
		{
			name: "due event already reported",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewNotificationDueEvent(context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								QuotaRequestsAllAuthenticated.Enum(),
								"id",
								"url",
								time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
								1000,
								200,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				dueNotifications: []*quota.NotificationDueEvent{
					quota.NewNotificationDueEvent(
						context.Background(),
						&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
						QuotaRequestsAllAuthenticated.Enum(),
						"id",
						"url",
						time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
						1000,
						250,
					),
				},
			},
			res: res{},
		},
		{
			name: "due event not reported",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						quota.NewNotificationDueEvent(context.Background(),
							&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
							QuotaRequestsAllAuthenticated.Enum(),
							"id",
							"url",
							time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
							1000,
							250,
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				dueNotifications: []*quota.NotificationDueEvent{
					quota.NewNotificationDueEvent(
						context.Background(),
						&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
						QuotaRequestsAllAuthenticated.Enum(),
						"id",
						"url",
						time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
						1000,
						250,
					),
				},
			},
			res: res{},
		},
		{
			name: "due events",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewNotificationDueEvent(context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								QuotaRequestsAllAuthenticated.Enum(),
								"id2",
								"url",
								time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
								1000,
								250,
							),
						),
					),
					expectFilter(),
					expectPush(
						quota.NewNotificationDueEvent(context.Background(),
							&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
							QuotaRequestsAllAuthenticated.Enum(),
							"id1",
							"url",
							time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
							1000,
							250,
						),
						quota.NewNotificationDueEvent(context.Background(),
							&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
							QuotaRequestsAllAuthenticated.Enum(),
							"id3",
							"url",
							time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
							1000,
							250,
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				dueNotifications: []*quota.NotificationDueEvent{
					quota.NewNotificationDueEvent(
						context.Background(),
						&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
						QuotaRequestsAllAuthenticated.Enum(),
						"id1",
						"url",
						time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
						1000,
						250,
					),
					quota.NewNotificationDueEvent(
						context.Background(),
						&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
						QuotaRequestsAllAuthenticated.Enum(),
						"id2",
						"url",
						time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
						1000,
						250,
					),
					quota.NewNotificationDueEvent(
						context.Background(),
						&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
						QuotaRequestsAllAuthenticated.Enum(),
						"id3",
						"url",
						time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
						1000,
						250,
					),
				},
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := r.ReportQuotaUsage(tt.args.ctx, tt.args.dueNotifications)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestQuotaReport_UsageNotificationSent(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		ctx             context.Context
		dueNotification *quota.NotificationDueEvent
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "usage notification sent, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPush(
						quota.NewNotifiedEvent(
							context.Background(),
							"quota1",
							quota.NewNotificationDueEvent(
								context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								QuotaRequestsAllAuthenticated.Enum(),
								"id1",
								"url",
								time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
								1000,
								250,
							),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				dueNotification: quota.NewNotificationDueEvent(
					context.Background(),
					&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
					QuotaRequestsAllAuthenticated.Enum(),
					"id1",
					"url",
					time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					1000,
					250,
				),
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			err := r.UsageNotificationSent(tt.args.ctx, tt.args.dueNotification)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
