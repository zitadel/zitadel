package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestQuota_AddQuota(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		ctx      context.Context
		setQuota *SetQuota
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
			name: "already existing",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							quota.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
									quota.SetEventType,
								),
								QuotaRequestsAllAuthenticated.Enum(),
								quota.ChangeFrom(time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)),
								quota.ChangeResetInterval(30*24*time.Hour),
								quota.ChangeAmount(1000),
								quota.ChangeLimit(false),
								quota.ChangeNotifications(make([]*quota.SetEventNotification, 0)),
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					ResetInterval: 30 * 24 * time.Hour,
					Amount:        1000,
					Limit:         true,
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "create quota, validation fail",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          "unimplemented",
					From:          time.Time{},
					ResetInterval: 0,
					Amount:        0,
					Limit:         false,
					Notifications: nil,
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-OTeSh", ""))
				},
			},
		},
		{
			name: "create quota, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						quota.NewSetEvent(
							eventstore.NewBaseEventForPush(
								context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								quota.SetEventType,
							),
							QuotaRequestsAllAuthenticated.Enum(),
							quota.ChangeFrom(time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)),
							quota.ChangeResetInterval(30*24*time.Hour),
							quota.ChangeAmount(1000),
							quota.ChangeLimit(true),
							quota.ChangeNotifications(make([]*quota.SetEventNotification, 0)),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					ResetInterval: 30 * 24 * time.Hour,
					Amount:        1000,
					Limit:         true,
					Notifications: nil,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "recreate quota, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
									quota.SetEventType,
								),
								QuotaRequestsAllAuthenticated.Enum(),
								quota.ChangeFrom(time.Now()),
								quota.ChangeResetInterval(30*24*time.Hour),
								quota.ChangeAmount(1000),
								quota.ChangeLimit(true),
								quota.ChangeNotifications(make([]*quota.SetEventNotification, 0)),
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewRemovedEvent(context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								QuotaRequestsAllAuthenticated.Enum(),
							),
						),
					),
					expectPush(
						quota.NewSetEvent(
							eventstore.NewBaseEventForPush(
								context.Background(),
								&quota.NewAggregate("quota2", "INSTANCE").Aggregate,
								quota.SetEventType,
							),
							QuotaRequestsAllAuthenticated.Enum(),
							quota.ChangeFrom(time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)),
							quota.ChangeResetInterval(30*24*time.Hour),
							quota.ChangeAmount(1000),
							quota.ChangeLimit(true),
							quota.ChangeNotifications(make([]*quota.SetEventNotification, 0)),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota2"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					ResetInterval: 30 * 24 * time.Hour,
					Amount:        1000,
					Limit:         true,
					Notifications: nil,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "create quota with notifications, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						quota.NewSetEvent(
							eventstore.NewBaseEventForPush(
								context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								quota.SetEventType,
							),
							QuotaRequestsAllAuthenticated.Enum(),
							quota.ChangeFrom(time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)),
							quota.ChangeResetInterval(30*24*time.Hour),
							quota.ChangeAmount(1000),
							quota.ChangeLimit(true),
							quota.ChangeNotifications(
								[]*quota.SetEventNotification{{
									ID:      "notification1",
									Percent: 20,
									Repeat:  false,
									CallURL: "https://url.com",
								}},
							),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota1", "notification1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					ResetInterval: 30 * 24 * time.Hour,
					Amount:        1000,
					Limit:         true,
					Notifications: QuotaNotifications{
						{
							Percent: 20,
							Repeat:  false,
							CallURL: "https://url.com",
						},
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := r.AddQuota(tt.args.ctx, tt.args.setQuota)
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

func TestQuota_SetQuota(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		ctx      context.Context
		setQuota *SetQuota
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
			name: "already existing",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							quota.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
									quota.SetEventType,
								),
								QuotaRequestsAllAuthenticated.Enum(),
								quota.ChangeFrom(time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)),
								quota.ChangeResetInterval(30*24*time.Hour),
								quota.ChangeAmount(1000),
								quota.ChangeLimit(true),
								quota.ChangeNotifications(make([]*quota.SetEventNotification, 0)),
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					ResetInterval: 30 * 24 * time.Hour,
					Amount:        1000,
					Limit:         true,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "create quota, validation fail",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          "unimplemented",
					From:          time.Time{},
					ResetInterval: 0,
					Amount:        0,
					Limit:         false,
					Notifications: nil,
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-OTeSh", ""))
				},
			},
		},
		{
			name: "create quota, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						quota.NewSetEvent(
							eventstore.NewBaseEventForPush(
								context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								quota.SetEventType,
							),
							QuotaRequestsAllAuthenticated.Enum(),
							quota.ChangeFrom(time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)),
							quota.ChangeResetInterval(30*24*time.Hour),
							quota.ChangeAmount(1000),
							quota.ChangeLimit(true),
							quota.ChangeNotifications(make([]*quota.SetEventNotification, 0)),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					ResetInterval: 30 * 24 * time.Hour,
					Amount:        1000,
					Limit:         true,
					Notifications: nil,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "recreate quota, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
									quota.SetEventType,
								),
								QuotaRequestsAllAuthenticated.Enum(),
								quota.ChangeFrom(time.Now()),
								quota.ChangeResetInterval(30*24*time.Hour),
								quota.ChangeAmount(1000),
								quota.ChangeLimit(true),
								quota.ChangeNotifications(make([]*quota.SetEventNotification, 0)),
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewRemovedEvent(context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								QuotaRequestsAllAuthenticated.Enum(),
							),
						),
					),
					expectPush(
						quota.NewSetEvent(
							eventstore.NewBaseEventForPush(
								context.Background(),
								&quota.NewAggregate("quota2", "INSTANCE").Aggregate,
								quota.SetEventType,
							),
							QuotaRequestsAllAuthenticated.Enum(),
							quota.ChangeFrom(time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)),
							quota.ChangeResetInterval(30*24*time.Hour),
							quota.ChangeAmount(1000),
							quota.ChangeLimit(true),
							quota.ChangeNotifications(make([]*quota.SetEventNotification, 0)),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota2"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					ResetInterval: 30 * 24 * time.Hour,
					Amount:        1000,
					Limit:         true,
					Notifications: nil,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "create quota with notifications, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						quota.NewSetEvent(
							eventstore.NewBaseEventForPush(
								context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								quota.SetEventType,
							),
							QuotaRequestsAllAuthenticated.Enum(),
							quota.ChangeFrom(time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)),
							quota.ChangeResetInterval(30*24*time.Hour),
							quota.ChangeAmount(1000),
							quota.ChangeLimit(true),
							quota.ChangeNotifications(
								[]*quota.SetEventNotification{{
									ID:      "notification1",
									Percent: 20,
									Repeat:  false,
									CallURL: "https://url.com",
								}},
							),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "quota1", "notification1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				setQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					ResetInterval: 30 * 24 * time.Hour,
					Amount:        1000,
					Limit:         true,
					Notifications: QuotaNotifications{
						{
							Percent: 20,
							Repeat:  false,
							CallURL: "https://url.com",
						},
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := r.SetQuota(tt.args.ctx, tt.args.setQuota)
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

func TestQuota_RemoveQuota(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx  context.Context
		unit QuotaUnit
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
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				unit: QuotaRequestsAllAuthenticated,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-WDfFf", ""))
				},
			},
		},
		{
			name: "already removed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
									quota.SetEventType,
								),
								QuotaRequestsAllAuthenticated.Enum(),
								quota.ChangeFrom(time.Now()),
								quota.ChangeResetInterval(30*24*time.Hour),
								quota.ChangeAmount(1000),
								quota.ChangeLimit(true),
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewRemovedEvent(context.Background(),
								&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
								QuotaRequestsAllAuthenticated.Enum(),
							),
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				unit: QuotaRequestsAllAuthenticated,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-WDfFf", ""))
				},
			},
		},
		{
			name: "remove quota, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							quota.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
									quota.SetEventType,
								),
								QuotaRequestsAllAuthenticated.Enum(),
								quota.ChangeFrom(time.Now()),
								quota.ChangeResetInterval(30*24*time.Hour),
								quota.ChangeAmount(1000),
								quota.ChangeLimit(false),
							),
						),
					),
					expectPush(
						quota.NewRemovedEvent(context.Background(),
							&quota.NewAggregate("quota1", "INSTANCE").Aggregate,
							QuotaRequestsAllAuthenticated.Enum(),
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				unit: QuotaRequestsAllAuthenticated,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveQuota(tt.args.ctx, tt.args.unit)
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

func TestQuota_QuotaNotification_validate(t *testing.T) {
	type args struct {
		quotaNotification *QuotaNotification
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "notification url parse failed",
			args: args{
				quotaNotification: &QuotaNotification{
					Percent: 20,
					Repeat:  false,
					CallURL: "%",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-bZ0Fj", ""))
				},
			},
		},
		{
			name: "notification url parse empty schema",
			args: args{
				quotaNotification: &QuotaNotification{
					Percent: 20,
					Repeat:  false,
					CallURL: "localhost:8080",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-HAYmN", ""))
				},
			},
		},
		{
			name: "notification url parse empty host",
			args: args{
				quotaNotification: &QuotaNotification{
					Percent: 20,
					Repeat:  false,
					CallURL: "https://",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-HAYmN", ""))
				},
			},
		},
		{
			name: "notification url parse percent 0",
			args: args{
				quotaNotification: &QuotaNotification{
					Percent: 0,
					Repeat:  false,
					CallURL: "https://localhost:8080",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-pBfjq", ""))
				},
			},
		},
		{
			name: "notification, ok",
			args: args{
				quotaNotification: &QuotaNotification{
					Percent: 20,
					Repeat:  false,
					CallURL: "https://localhost:8080",
				},
			},
			res: res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.quotaNotification.validate()
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestQuota_SetQuota_validate(t *testing.T) {
	type args struct {
		addQuota *SetQuota
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "notification url parse failed",
			args: args{
				addQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Now(),
					ResetInterval: time.Minute * 10,
					Amount:        100,
					Limit:         true,
					Notifications: QuotaNotifications{
						{
							Percent: 20,
							Repeat:  false,
							CallURL: "%",
						},
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-bZ0Fj", ""))
				},
			},
		},
		{
			name: "unit unimplemented",
			args: args{
				addQuota: &SetQuota{
					Unit:          "unimplemented",
					From:          time.Now(),
					ResetInterval: time.Minute * 10,
					Amount:        100,
					Limit:         true,
					Notifications: nil,
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-OTeSh", ""))
				},
			},
		},
		{
			name: "reset interval under 1 min",
			args: args{
				addQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Now(),
					ResetInterval: time.Second * 10,
					Amount:        100,
					Limit:         true,
					Notifications: nil,
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "QUOTA-R5otd", ""))
				},
			},
		},
		{
			name: "validate, ok",
			args: args{
				addQuota: &SetQuota{
					Unit:          QuotaRequestsAllAuthenticated,
					From:          time.Now(),
					ResetInterval: time.Minute * 10,
					Amount:        100,
					Limit:         false,
					Notifications: nil,
				},
			},
			res: res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.addQuota.validate()
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
