package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommandSide_AddSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
		alg         crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		instanceID string
		sms        *twilio.Config
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
			name: "add sms config twilio, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"sid",
								"senderName",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								},
							),
							),
						},
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "providerid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				sms: &twilio.Config{
					SID:          "sid",
					Token:        "token",
					SenderNumber: "senderName",
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
				eventstore:    tt.fields.eventstore,
				idGenerator:   tt.fields.idGenerator,
				smsEncryption: tt.fields.alg,
			}
			_, got, err := r.AddSMSConfigTwilio(tt.args.ctx, tt.args.instanceID, tt.args.sms)
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

func TestCommandSide_ChangeSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
		sms        *twilio.Config
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
			name: "id empty, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.Config{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				sms:        &twilio.Config{},
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"sid",
								"senderName",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.Config{
					SID:          "sid",
					Token:        "token",
					SenderNumber: "senderName",
				},
				instanceID: "INSTANCE",
				id:         "providerid",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "sms config twilio change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"sid",
								"token",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newSMSConfigTwilioChangedEvent(
									context.Background(),
									"providerid",
									"sid2",
									"senderName2",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.Config{
					SID:          "sid2",
					Token:        "token2",
					SenderNumber: "senderName2",
				},
				instanceID: "INSTANCE",
				id:         "providerid",
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
			got, err := r.ChangeSMSConfigTwilio(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.sms)
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

func TestCommandSide_ActivateSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
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
			name: "id empty, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "sms config twilio activate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewSMSConfigTwilioActivatedEvent(
									context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "providerid",
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
			got, err := r.ActivateSMSConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
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

func TestCommandSide_DeactivateSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
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
			name: "id empty, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "sms config twilio deactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigTwilioActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewSMSConfigDeactivatedEvent(
									context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "providerid",
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
			got, err := r.DeactivateSMSConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
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

func TestCommandSide_RemoveSMSConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
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
			name: "id empty, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "sms config remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewSMSConfigRemovedEvent(
									context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "providerid",
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
			got, err := r.RemoveSMSConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
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

func newSMSConfigTwilioChangedEvent(ctx context.Context, id, sid, senderName string) *instance.SMSConfigTwilioChangedEvent {
	changes := []instance.SMSConfigTwilioChanges{
		instance.ChangeSMSConfigTwilioSID(sid),
		instance.ChangeSMSConfigTwilioSenderNumber(senderName),
	}
	event, _ := instance.NewSMSConfigTwilioChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		id,
		changes,
	)
	return event
}
