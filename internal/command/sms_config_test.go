package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
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
		ctx context.Context
		sms *AddTwilioConfig
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
			name: "add sms config twilio, missing resourceowner",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				sms: &AddTwilioConfig{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-ZLrZhKSKq0", "Errors.ResourceOwnerMissing"))
				},
			},
		},
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
							"description",
							"sid",
							"senderName",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("token"),
							},
							"",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "providerid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				sms: &AddTwilioConfig{
					ResourceOwner:    "INSTANCE",
					Description:      "description",
					SID:              "sid",
					Token:            "token",
					SenderNumber:     "senderName",
					VerifyServiceSID: "",
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
			err := r.AddSMSConfigTwilio(tt.args.ctx, tt.args.sms)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.sms.Details)
			}
		})
	}
}

func TestCommandSide_ChangeSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		sms *ChangeTwilioConfig
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
			name: "resourceowner empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				sms: &ChangeTwilioConfig{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-RHXryJwmFG", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				sms: &ChangeTwilioConfig{
					ResourceOwner: "INSTANCE",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-gMr93iNhTR", "Errors.IDMissing"))
				},
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
				ctx: context.Background(),
				sms: &ChangeTwilioConfig{
					ResourceOwner: "INSTANCE",
					ID:            "id",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-MUY0IFAf8O", "Errors.SMSConfig.NotFound"))
				},
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"description",
								"sid",
								"senderName",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								},
								"",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &ChangeTwilioConfig{
					ResourceOwner:    "INSTANCE",
					ID:               "providerid",
					SID:              gu.Ptr("sid"),
					Token:            gu.Ptr("token"),
					SenderNumber:     gu.Ptr("senderName"),
					VerifyServiceSID: gu.Ptr(""),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
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
								"description",
								"sid",
								"token",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								},
								"verifyServiceSid",
							),
						),
					),
					expectPush(
						newSMSConfigTwilioChangedEvent(
							context.Background(),
							"providerid",
							"sid2",
							"senderName2",
							"description2",
							"verifyServiceSid2",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &ChangeTwilioConfig{
					ResourceOwner:    "INSTANCE",
					ID:               "providerid",
					Description:      gu.Ptr("description2"),
					SID:              gu.Ptr("sid2"),
					Token:            gu.Ptr("token2"),
					SenderNumber:     gu.Ptr("senderName2"),
					VerifyServiceSID: gu.Ptr("verifyServiceSid2"),
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
				eventstore: tt.fields.eventstore(t),
			}
			err := r.ChangeSMSConfigTwilio(tt.args.ctx, tt.args.sms)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.sms.Details)
			}
		})
	}
}

func TestCommandSide_AddSMSConfigHTTP(t *testing.T) {
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		idGenerator                 id.Generator
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
		alg                         crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx  context.Context
		http *AddSMSHTTP
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
			name: "add sms config http, resource owner missing",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:  context.Background(),
				http: &AddSMSHTTP{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-huy99qWjX4", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "add sms config http, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewSMSConfigHTTPAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"providerid",
							"description",
							"endpoint",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
						),
					),
				),
				idGenerator:                 id_mock.NewIDGeneratorExpectIDs(t, "providerid"),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args: args{
				ctx: context.Background(),
				http: &AddSMSHTTP{
					ResourceOwner: "INSTANCE",
					Description:   "description",
					Endpoint:      "endpoint",
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
				eventstore:                  tt.fields.eventstore(t),
				idGenerator:                 tt.fields.idGenerator,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				smsEncryption:               tt.fields.alg,
			}
			err := r.AddSMSConfigHTTP(tt.args.ctx, tt.args.http)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.http.Details)
			}
		})
	}
}

func TestCommandSide_ChangeSMSConfigHTTP(t *testing.T) {
	type fields struct {
		eventstore                  func(*testing.T) *eventstore.Eventstore
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		ctx  context.Context
		http *ChangeSMSHTTP
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
			name: "resourceowner empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:  context.Background(),
				http: &ChangeSMSHTTP{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-M622CFQnwK", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				http: &ChangeSMSHTTP{
					ResourceOwner: "INSTANCE",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-phyb2e4Kll", "Errors.IDMissing"))
				},
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
				ctx: context.Background(),
				http: &ChangeSMSHTTP{
					ResourceOwner: "INSTANCE",
					ID:            "id",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-6NW4I5Kqzj", "Errors.SMSConfig.NotFound"))
				},
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"description",
								"endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				http: &ChangeSMSHTTP{
					ResourceOwner: "INSTANCE",
					ID:            "providerid",
					Endpoint:      gu.Ptr("endpoint"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
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
								"description",
								"endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
							),
						),
					),
					expectPush(
						newSMSConfigHTTPChangedEvent(
							context.Background(),
							"providerid",
							"endpoint2",
							"description2",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("87654321"),
							},
						),
					),
				),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("87654321", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args: args{
				ctx: context.Background(),
				http: &ChangeSMSHTTP{
					ResourceOwner:        "INSTANCE",
					ID:                   "providerid",
					Description:          gu.Ptr("description2"),
					Endpoint:             gu.Ptr("endpoint2"),
					ExpirationSigningKey: true,
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
				eventstore:                  tt.fields.eventstore(t),
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
			}
			err := r.ChangeSMSConfigHTTP(tt.args.ctx, tt.args.http)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.http.Details)
			}
		})
	}
}

func TestCommandSide_ActivateSMSConfig(t *testing.T) {
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
			name: "resourceowner empty, invalid error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-EFgoOg997V", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, invalid error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-jJ6TVqzvjp", "Errors.IDMissing"))
				},
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
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-9ULtp9PH5E", "Errors.SMSConfig.NotFound"))
				},
			},
		},
		{
			name: "sms existing, already active",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"description",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
								"",
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
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "providerid",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-B25GFeIvRi", "Errors.SMSConfig.AlreadyActive"))
				},
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
								"description",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
								"",
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
								"description",
								"endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ActivateSMSConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
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

func TestCommandSide_DeactivateSMSConfig(t *testing.T) {
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
	}{{
		name: "resourceowner empty, invalid error",
		fields: fields{
			eventstore: expectEventstore(),
		},
		args: args{
			ctx: context.Background(),
		},
		res: res{
			err: func(err error) bool {
				return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-V9NWOZj8Gi", "Errors.ResourceOwnerMissing"))
			},
		},
	},
		{
			name: "id empty, invalid error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-xs1ah1v1CL", "Errors.IDMissing"))
				},
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
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-La91dGNhbM", "Errors.SMSConfig.NotFound"))
				},
			},
		},
		{
			name: "sms config twilio deactivate, already deactivated",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
								"description",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
								"",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigDeactivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"providerid",
							),
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
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-OSZAEkYvk7", "Errors.SMSConfig.AlreadyDeactivated"))
				},
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
								"description",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
								"",
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
								"description",
								"endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
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
			name: "resourceowner empty, invalid error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-cw0NSJsn1v", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, invalid error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-Qrz7lvdC4c", "Errors.IDMissing"))
				},
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
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-povEVHPCkV", "Errors.SMSConfig.NotFound"))
				},
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
								"description",
								"sid",
								"sender-name",
								&crypto.CryptoValue{},
								"",
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
								"description",
								"endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
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

func newSMSConfigTwilioChangedEvent(ctx context.Context, id, sid, senderName, description, verifyServiceSid string) *instance.SMSConfigTwilioChangedEvent {
	changes := []instance.SMSConfigTwilioChanges{
		instance.ChangeSMSConfigTwilioSID(sid),
		instance.ChangeSMSConfigTwilioSenderNumber(senderName),
		instance.ChangeSMSConfigTwilioDescription(description),
		instance.ChangeSMSConfigTwilioVerifyServiceSID(verifyServiceSid),
	}
	event, _ := instance.NewSMSConfigTwilioChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		id,
		changes,
	)
	return event
}

func newSMSConfigHTTPChangedEvent(ctx context.Context, id, endpoint, description string, signingKey *crypto.CryptoValue) *instance.SMSConfigHTTPChangedEvent {
	changes := []instance.SMSConfigHTTPChanges{
		instance.ChangeSMSConfigHTTPEndpoint(endpoint),
		instance.ChangeSMSConfigHTTPDescription(description),
		instance.ChangeSMSConfigHTTPSigningKey(signingKey),
	}
	event, _ := instance.NewSMSConfigHTTPChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		id,
		changes,
	)
	return event
}
