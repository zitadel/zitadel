package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCommandSide_AddSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
		alg         crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx context.Context
		sms *twilio.TwilioConfig
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
							eventFromEventPusher(iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
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
				ctx: context.Background(),
				sms: &twilio.TwilioConfig{
					SID:        "sid",
					Token:      "token",
					SenderName: "senderName",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
				smsCrypto:   tt.fields.alg,
			}
			_, got, err := r.AddSMSConfigTwilio(tt.args.ctx, tt.args.sms)
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
		ctx context.Context
		id  string
		sms *twilio.TwilioConfig
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
				sms: &twilio.TwilioConfig{},
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
				ctx: context.Background(),
				sms: &twilio.TwilioConfig{},
				id:  "id",
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
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
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
				sms: &twilio.TwilioConfig{
					SID:        "sid",
					Token:      "token",
					SenderName: "senderName",
				},
				id: "providerid",
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
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
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
				sms: &twilio.TwilioConfig{
					SID:        "sid2",
					Token:      "token2",
					SenderName: "senderName2",
				},
				id: "providerid",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeSMSConfigTwilio(tt.args.ctx, tt.args.id, tt.args.sms)
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
		ctx context.Context
		id  string
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
				ctx: context.Background(),
				id:  "id",
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
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
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
								iam.NewSMSConfigTwilioActivatedEvent(
									context.Background(),
									&iam.NewAggregate().Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "providerid",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ActivateSMSConfigTwilio(tt.args.ctx, tt.args.id)
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
		ctx context.Context
		id  string
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
				ctx: context.Background(),
				id:  "id",
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
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
							),
						),
						eventFromEventPusher(
							iam.NewSMSConfigTwilioActivatedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewSMSConfigDeactivatedEvent(
									context.Background(),
									&iam.NewAggregate().Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "providerid",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.DeactivateSMSConfigTwilio(tt.args.ctx, tt.args.id)
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

func TestCommandSide_RemoveSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
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
				ctx: context.Background(),
				id:  "id",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "sms config twilio remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
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
								iam.NewSMSConfigRemovedEvent(
									context.Background(),
									&iam.NewAggregate().Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "providerid",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveSMSConfigTwilio(tt.args.ctx, tt.args.id)
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

func newSMSConfigTwilioChangedEvent(ctx context.Context, id, sid, senderName string) *iam.SMSConfigTwilioChangedEvent {
	changes := []iam.SMSConfigTwilioChanges{
		iam.ChangeSMSConfigTwilioSID(sid),
		iam.ChangeSMSConfigTwilioSenderName(senderName),
	}
	event, _ := iam.NewSMSConfigTwilioChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		id,
		changes,
	)
	return event
}
