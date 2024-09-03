package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore  func(*testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
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
				eventstore:    tt.fields.eventstore(t),
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
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.Config{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
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
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "sms config twilio change, ok",
			fields: fields{
				eventstore: expectEventstore(
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
						newSMSConfigTwilioChangedEvent(
							context.Background(),
							"providerid",
							"sid2",
							"senderName2",
						),
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ChangeSMSConfigTwilio(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.sms)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddSMSConfigHTTP(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
		alg         crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		instanceID string
		http       *AddSMSHTTP
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
			name: "add sms config http, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewSMSConfigHTTPAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
							"endpoint",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "providerid"),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				http: &AddSMSHTTP{
					Endpoint: "endpoint",
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
			_, got, err := r.AddSMSConfigHTTP(tt.args.ctx, tt.args.instanceID, tt.args.http)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeSMSConfigHTTP(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
		http       *ChangeSMSHTTP
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:  context.Background(),
				http: &ChangeSMSHTTP{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				http:       &ChangeSMSHTTP{},
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"endpoint",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				http: &ChangeSMSHTTP{
					Endpoint: "endpoint",
				},
				instanceID: "INSTANCE",
				id:         "providerid",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "sms config http change, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"endpoint",
							),
						),
					),
					expectPush(
						newSMSConfigHTTPChangedEvent(
							context.Background(),
							"providerid",
							"endpoint2",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				http: &ChangeSMSHTTP{
					Endpoint: "endpoint2",
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ChangeSMSConfigHTTP(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.http)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ActivateSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx         context.Context
		instanceID  string
		id          string
		activatedID string
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "sms config twilio activate, ok",
			fields: fields{
				eventstore: expectEventstore(
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
						instance.NewSMSConfigActivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
						),
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
		{
			name: "sms config http activate, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"endpoint",
							),
						),
					),
					expectPush(
						instance.NewSMSConfigActivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
						),
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
		{
			name: "sms config http activate,with already active, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"activeproviderid",
								"activeendpoint",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"activeproviderid",
							),
						),
					),
					expectPush(
						instance.NewSMSConfigDeactivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"activeproviderid",
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"endpoint",
							),
						),
					),
					expectPush(
						instance.NewSMSConfigActivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
						),
					),
				),
			},
			args: args{
				ctx:         context.Background(),
				instanceID:  "INSTANCE",
				id:          "providerid",
				activatedID: "activeproviderid",
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ActivateSMSConfig(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.activatedID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_DeactivateSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "sms config twilio deactivate, ok",
			fields: fields{
				eventstore: expectEventstore(
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
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
							),
						),
					),
					expectPush(
						instance.NewSMSConfigDeactivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
						),
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
		{
			name: "sms config http deactivate, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"endpoint",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
							),
						),
					),
					expectPush(
						instance.NewSMSConfigDeactivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
						),
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.DeactivateSMSConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveSMSConfig(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "sms config remove, twilio, ok",
			fields: fields{
				eventstore: expectEventstore(
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
						instance.NewSMSConfigRemovedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
						),
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
		{
			name: "sms config remove, http, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"endpoint",
							),
						),
					),
					expectPush(
						instance.NewSMSConfigRemovedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
						),
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.RemoveSMSConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
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

func newSMSConfigHTTPChangedEvent(ctx context.Context, id, endpoint string) *instance.SMSConfigHTTPChangedEvent {
	changes := []instance.SMSConfigHTTPChanges{
		instance.ChangeSMSConfigHTTPEndpoint(endpoint),
	}
	event, _ := instance.NewSMSConfigHTTPChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		id,
		changes,
	)
	return event
}
