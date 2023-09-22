package command

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/quota"
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
		idGenerator   id.Generator
		createNew     bool
		amount        uint64
		from          time.Time
		resetInterval time.Duration
		limit         bool
		notifications []*QuotaNotification
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantEventPayloadJSON string
		wantErr              assert.ErrorAssertionFunc
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
			resetInterval: time.Millisecond,
			limit:         true,
			notifications: make([]*QuotaNotification, 0),
		},
		wantEventPayloadJSON: `{"unit":0,"interval":1000000}`,
		wantErr:              assert.NoError,
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
			resetInterval: time.Millisecond,
			limit:         true,
			notifications: make([]*QuotaNotification, 0),
		},
		wantEventPayloadJSON: `{"unit":0,"interval":1000000,"amount":10}`,
		wantErr:              assert.NoError,
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
		wantEventPayloadJSON: `{"unit":0}`,
		wantErr:              assert.NoError,
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
		wantEventPayloadJSON: `{"unit":0,"limit":false}`,
		wantErr:              assert.NoError,
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
		wantEventPayloadJSON: `{"unit":0,"amount":0}`,
		wantErr:              assert.NoError,
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
		wantEventPayloadJSON: `{"unit":0,"from":"0001-01-01T00:00:00Z"}`,
		wantErr:              assert.NoError,
	}, {
		name: "add notification",
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
				Percent: 20,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "notification1"),
		},
		wantEventPayloadJSON: `{"unit":0,"notifications":[{"id":"notification1","percent":20,"repeat":true,"callUrl":"https://call.url"}]}`,
		wantErr:              assert.NoError,
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
		wantEventPayloadJSON: `{"unit":0}`,
		wantErr:              assert.NoError,
	}, {
		name: "change nothing but notification order",
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
		wantEventPayloadJSON: `{"unit":0}`,
		wantErr:              assert.NoError,
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
		wantEventPayloadJSON: `{"unit":0,"notifications":[]}`,
		wantErr:              assert.NoError,
	}, {
		name: "create new without notification",
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
		wantEventPayloadJSON: `{"unit":0,"notifications":[]}`,
		wantErr:              assert.NoError,
	}, {
		name: "create new with all values values",
		args: args{
			amount:        5,
			from:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			resetInterval: time.Millisecond,
			limit:         true,
			notifications: []*QuotaNotification{{
				Percent: 10,
				Repeat:  true,
				CallURL: "https://call.url",
			}},
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "notification1"),
			createNew:   true,
		},
		wantEventPayloadJSON: `{"unit":0,"from":"2020-01-01T00:00:00Z","interval":1000000,"amount":5,"limit":true,"notifications":[{"id":"notification1","percent":10,"repeat":true,"callUrl":"https://call.url"}]}`,
		wantErr:              assert.NoError,
	}, {
		name:                 "create new with zero values",
		args:                 args{createNew: true},
		wantEventPayloadJSON: `{"unit":0,"from":"0001-01-01T00:00:00Z","interval":0,"amount":0,"limit":false,"notifications":[]}`,
		wantErr:              assert.NoError,
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
			gotChanges, err := wm.NewChanges(tt.args.idGenerator, tt.args.createNew, tt.args.amount, tt.args.from, tt.args.resetInterval, tt.args.limit, tt.args.notifications...)
			if !tt.wantErr(t, err, fmt.Sprintf("NewChanges(%v, %v, %v, %v, %v, %v)", tt.args.createNew, tt.args.amount, tt.args.from, tt.args.resetInterval, tt.args.limit, tt.args.notifications)) {
				return
			}
			bytes, err := json.Marshal(quota.NewSetEvent(
				eventstore.NewBaseEventForPush(
					context.Background(),
					&quota.NewAggregate("quota1", "instance1").Aggregate,
					quota.SetEventType,
				),
				quota.Unimplemented,
				gotChanges...,
			))
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantEventPayloadJSON, string(bytes), "NewChanges(%v, %v, %v, %v, %v, %v)", tt.args.createNew, tt.args.amount, tt.args.from, tt.args.resetInterval, tt.args.limit, tt.args.notifications)
		})
	}
}
