package command

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestQuotaWriteModel_NewChanges(t *testing.T) {
	type fields struct {
		from          time.Time
		resetInterval time.Duration
		amount        uint64
		limit         bool
		notifications []*quota.SetEventNotification
	}
	type args struct {
		idGenerator   id_generator.Generator
		createNew     bool
		amount        uint64
		from          time.Time
		resetInterval time.Duration
		limit         bool
		notifications []*QuotaNotification
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantEvent   quota.SetEvent
		wantChanges int
		wantErr     assert.ErrorAssertionFunc
	}{{
		name: "change reset interval",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: make([]*quota.SetEventNotification, 0),
		},
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Minute,
			limit:         true,
			notifications: make([]*QuotaNotification, 0),
		},
		wantChanges: 1,
		wantEvent: quota.SetEvent{
			ResetInterval: durationPtr(time.Minute),
		},
	}, {
		name: "change reset interval and amount",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: make([]*quota.SetEventNotification, 0),
		},
		args: args{
			amount:        10,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Minute,
			limit:         true,
			notifications: make([]*QuotaNotification, 0),
		},
		wantChanges: 2,
		wantEvent: quota.SetEvent{
			ResetInterval: durationPtr(time.Minute),
			Amount:        uint64Ptr(10),
		},
	}, {
		name: "change nothing",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*quota.SetEventNotification{},
		},
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*QuotaNotification{},
		},
	}, {
		name: "change limit to zero value",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: make([]*quota.SetEventNotification, 0),
		},
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         false,
			notifications: make([]*QuotaNotification, 0),
		},
		wantChanges: 1,
		wantEvent:   quota.SetEvent{Limit: boolPtr(false)},
	}, {
		name: "change amount to zero value",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: make([]*quota.SetEventNotification, 0),
		},
		args: args{
			amount:        0,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: make([]*QuotaNotification, 0),
		},
		wantChanges: 1,
		wantEvent:   quota.SetEvent{Amount: uint64Ptr(0)},
	}, {
		name: "change from to zero value",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: make([]*quota.SetEventNotification, 0),
		},
		args: args{
			amount:        5,
			from:          time.Time{},
			resetInterval: time.Hour,
			limit:         true,
			notifications: make([]*QuotaNotification, 0),
		},
		wantChanges: 1,
		wantEvent:   quota.SetEvent{From: &time.Time{}},
	}, {
		name: "add notification",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
		},
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*QuotaNotification{{
				Percent: 10,
				Repeat:  false,
				CallURL: "https://call.url",
			}},
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "notification1"),
		},
		wantChanges: 1,
		wantEvent: quota.SetEvent{Notifications: &[]*quota.SetEventNotification{{
			ID:      "notification1",
			Percent: 10,
			Repeat:  false,
			CallURL: "https://call.url",
		}}},
	}, {
		name: "change nothing with notification",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*quota.SetEventNotification{{
				ID:      "notification1",
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
		},
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*QuotaNotification{{
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
			idGenerator: id_mock.NewIDGenerator(t),
		},
	}, {
		name: "don't change notification order",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*quota.SetEventNotification{{
				ID:      "notification1",
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}, {
				ID:      "notification2",
				Percent: 10,
				Repeat:  false,
				CallURL: "https://call.url",
			}},
		},
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*QuotaNotification{{
				Percent: 10,
				Repeat:  false,
				CallURL: "https://call.url",
			}, {
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "newnotification1", "newnotification2"),
		},
	}, {
		name: "change notification to zero value",
		fields: fields{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*quota.SetEventNotification{{
				ID:      "notification1",
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
		},
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*QuotaNotification{},
		},
		wantChanges: 1,
		wantEvent:   quota.SetEvent{Notifications: &[]*quota.SetEventNotification{}},
	}, {
		name: "validate no duplicate notifications",
		args: args{
			notifications: []*QuotaNotification{{
				Percent: 10,
				Repeat:  false,
				CallURL: "https://call.url",
			}, {
				Percent: 10,
				Repeat:  false,
				CallURL: "https://call.url",
			}},
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "notification1", "notification2"),
		},
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			return zerrors.IsErrorInvalidArgument(err)
		},
	}, {
		name: "deduplicate existing notifications",
		fields: fields{
			notifications: []*quota.SetEventNotification{{
				ID:      "existingnotification1",
				Percent: 10,
				Repeat:  false,
				CallURL: "https://call.url",
			}, {
				ID:      "existingnotification2",
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}, {
				ID:      "existingnotification3",
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
		},
		args: args{
			notifications: []*QuotaNotification{{
				Percent: 10,
				Repeat:  false,
				CallURL: "https://call.url",
			}, {
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "notification1", "notification2"),
		},
		wantChanges: 1,
		wantEvent: quota.SetEvent{
			Notifications: &[]*quota.SetEventNotification{{
				ID:      "notification1",
				Percent: 10,
				Repeat:  false,
				CallURL: "https://call.url",
			}, {
				ID:      "notification2",
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}}},
	}, {
		name: "create new with all values",
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Hour,
			limit:         true,
			notifications: []*QuotaNotification{{
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "notification1"),
			createNew:   true,
		},
		wantChanges: 5,
		wantEvent: quota.SetEvent{
			From:          timePtr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			ResetInterval: durationPtr(time.Hour),
			Amount:        uint64Ptr(5),
			Limit:         boolPtr(true),
			Notifications: &[]*quota.SetEventNotification{{
				ID:      "notification1",
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
		},
	}, {
		name:        "create new with zero values",
		args:        args{createNew: true},
		wantChanges: 5,
		wantEvent: quota.SetEvent{
			From:          &time.Time{},
			ResetInterval: durationPtr(0),
			Amount:        uint64Ptr(0),
			Limit:         boolPtr(false),
			Notifications: &[]*quota.SetEventNotification{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &quotaWriteModel{
				from:          tt.fields.from,
				resetInterval: tt.fields.resetInterval,
				amount:        tt.fields.amount,
				limit:         tt.fields.limit,
				notifications: tt.fields.notifications,
			}
			id_generator.SetGenerator(tt.args.idGenerator)
			gotChanges, err := wm.NewChanges(tt.args.createNew, tt.args.amount, tt.args.from, tt.args.resetInterval, tt.args.limit, tt.args.notifications...)
			assert.Len(t, gotChanges, tt.wantChanges)
			if tt.wantErr != nil {
				tt.wantErr(t, err, fmt.Sprintf("NewChanges(%v, %v, %v, %v, %v, %v)", tt.args.createNew, tt.args.amount, tt.args.from, tt.args.resetInterval, tt.args.limit, tt.args.notifications))
				return
			}
			assert.NoError(t, err)
			marshalled, err := json.Marshal(quota.NewSetEvent(
				eventstore.NewBaseEventForPush(
					context.Background(),
					&quota.NewAggregate("quota1", "instance1").Aggregate,
					quota.SetEventType,
				),
				quota.Unimplemented,
				gotChanges...,
			))
			assert.NoError(t, err)
			unmarshalled := new(quota.SetEvent)
			assert.NoError(t, json.Unmarshal(marshalled, unmarshalled))
			assert.Equalf(t, tt.wantEvent, *unmarshalled, "NewChanges(%v, %v, %v, %v, %v, %v)", tt.args.createNew, tt.args.amount, tt.args.from, tt.args.resetInterval, tt.args.limit, tt.args.notifications)
		})
	}
}

func uint64Ptr(i uint64) *uint64                 { return &i }
func boolPtr(b bool) *bool                       { return &b }
func durationPtr(d time.Duration) *time.Duration { return &d }
func timePtr(t time.Time) *time.Time             { return &t }
