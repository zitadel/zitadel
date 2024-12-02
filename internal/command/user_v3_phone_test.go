package command

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/senders/mock"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func filterSMSProviderActivated() expect {
	return expectFilter(
		eventFromEventPusher(
			instance.NewSMSConfigActivatedEvent(
				context.Background(),
				&instance.NewAggregate("instanceID").Aggregate,
				"id",
			),
		),
	)
}

func filterSMSProviderExternalExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			instance.NewSMSConfigTwilioAddedEvent(
				context.Background(),
				&instance.NewAggregate("instanceID").Aggregate,
				"id",
				"",
				"sid",
				"senderNumber",
				&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
				"verifyServiceSid",
			),
		),
		eventFromEventPusher(
			instance.NewSMSConfigActivatedEvent(
				context.Background(),
				&instance.NewAggregate("instanceID").Aggregate,
				"id",
			),
		),
	)
}

func filterSMSProviderInternalExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			instance.NewSMSConfigTwilioAddedEvent(
				context.Background(),
				&instance.NewAggregate("instanceID").Aggregate,
				"id",
				"",
				"sid",
				"senderNumber",
				&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
				"",
			),
		),
		eventFromEventPusher(
			instance.NewSMSConfigActivatedEvent(
				context.Background(),
				&instance.NewAggregate("instanceID").Aggregate,
				"id",
			),
		),
	)
}

func filterSchemaUserPhoneExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			schemauser.NewCreatedEvent(
				context.Background(),
				&schemauser.NewAggregate("user1", "org1").Aggregate,
				"id1",
				1,
				json.RawMessage(`{
						"name": "user1"
					}`),
			),
		),
		eventFromEventPusher(
			schemauser.NewPhoneUpdatedEvent(
				context.Background(),
				&schema.NewAggregate("id1", "instanceID").Aggregate,
				"+41791234567",
			),
		),
	)
}

func filterSchemaUserPhoneCodeExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			schemauser.NewCreatedEvent(
				context.Background(),
				&schemauser.NewAggregate("user1", "org1").Aggregate,
				"id1",
				1,
				json.RawMessage(`{
						"name": "user1"
					}`),
			),
		),
		eventFromEventPusher(
			schemauser.NewPhoneUpdatedEvent(
				context.Background(),
				&schemauser.NewAggregate("user1", "org1").Aggregate,
				"+41791234567",
			),
		),
		eventFromEventPusherWithCreationDateNow(
			schemauser.NewPhoneCodeAddedEvent(
				context.Background(),
				&schemauser.NewAggregate("user1", "org1").Aggregate,
				&crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("phoneverify"),
				},
				time.Hour*1,
				false,
				"",
			),
		),
	)
}

func filterSchemaUserPhoneCodeExternalExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			schemauser.NewCreatedEvent(
				context.Background(),
				&schemauser.NewAggregate("user1", "org1").Aggregate,
				"id1",
				1,
				json.RawMessage(`{
						"name": "user1"
					}`),
			),
		),
		eventFromEventPusher(
			schemauser.NewPhoneUpdatedEvent(
				context.Background(),
				&schemauser.NewAggregate("user1", "org1").Aggregate,
				"+41791234567",
			),
		),
		eventFromEventPusher(
			schemauser.NewPhoneCodeAddedEvent(
				context.Background(),
				&schemauser.NewAggregate("user1", "org1").Aggregate,
				nil,
				0,
				false,
				"id",
			),
		),
		eventFromEventPusherWithCreationDateNow(
			schemauser.NewPhoneCodeSentEvent(context.Background(),
				&authenticator.NewAggregate("user1", "org1").Aggregate,
				&senders.CodeGeneratorInfo{
					ID:             "id",
					VerificationID: "verificationID",
				},
			),
		),
	)
}

func TestCommands_ChangeSchemaUserPhone(t *testing.T) {
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		checkPermission             domain.PermissionCheck
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
	}
	type args struct {
		ctx  context.Context
		user *ChangeSchemaUserPhone
	}
	type res struct {
		returnCode string
		details    *domain.ObjectDetails
		err        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no userID, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-DkQ9aurv5u", "Errors.IDMissing"))
				},
			},
		},
		{
			"no valid phone, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{
					ID:    "user1",
					Phone: &Phone{Number: "nonumber"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "PHONE-so0wa", "Errors.User.Phone.Invalid"))
				},
			},
		}, {
			"phone update, user not found",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound"))
				},
			},
		},
		{
			"phone update, no permission",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{
					ID:    "user1",
					Phone: &Phone{Number: "+41791234567"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"phone update, phone not changed",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPhoneExisting(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{
					ID: "user1",
					Phone: &Phone{
						Number:     "+41791234567",
						ReturnCode: true,
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"phone update, phone return",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSMSProviderActivated(),
					filterSMSProviderInternalExisting(),
					expectPush(
						schemauser.NewPhoneUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"+41791234567",
						),
						schemauser.NewPhoneCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneverify"),
							},
							time.Hour*1,
							true,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{
					ID: "user1",
					Phone: &Phone{
						Number:     "+41791234567",
						ReturnCode: true,
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				returnCode: "phoneverify",
			},
		},
		{
			"user updated, phone to verify",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSMSProviderActivated(),
					filterSMSProviderInternalExisting(),
					expectPush(
						schemauser.NewPhoneUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"+41791234567",
						),
						schemauser.NewPhoneCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneverify"),
							}, time.Hour*1,
							false,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{
					ID: "user1",
					Phone: &Phone{
						Number: "+41791234567",
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"user updated, phone to verify (external)",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSMSProviderActivated(),
					filterSMSProviderExternalExisting(),
					expectPush(
						schemauser.NewPhoneUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"+41791234567",
						),
						schemauser.NewPhoneCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							nil,
							0,
							false,
							"id",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{
					ID: "user1",
					Phone: &Phone{
						Number: "+41791234567",
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"user updated, verified",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					expectPush(
						schemauser.NewPhoneUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"+41791234567",
						),
						schemauser.NewPhoneVerifiedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserPhone{
					ID: "user1",
					Phone: &Phone{
						Number:   "+41791234567",
						Verified: true,
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				checkPermission:             tt.fields.checkPermission,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				userEncryption:              crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				defaultSecretGenerators: &SecretGenerators{
					PhoneVerificationCode: &crypto.GeneratorConfig{
						Length:              8,
						Expiry:              time.Hour,
						IncludeLowerLetters: true,
						IncludeUpperLetters: true,
						IncludeDigits:       true,
						IncludeSymbols:      true,
					},
				},
			}
			details, err := c.ChangeSchemaUserPhone(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}

			if tt.res.returnCode != "" {
				assert.NotNil(t, tt.args.user.ReturnCode)
				assert.Equal(t, tt.res.returnCode, *tt.args.user.ReturnCode)
			}
		})
	}
}

func TestCommands_VerifySchemaUserPhone(t *testing.T) {
	type fields struct {
		eventstore        func(t *testing.T) *eventstore.Eventstore
		checkPermission   domain.PermissionCheck
		phoneCodeVerifier func(ctx context.Context, id string) (senders.CodeGenerator, error)
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		code          string
	}
	type res struct {
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no userID, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-R4LKY44Ke3", "Errors.IDMissing"))
				},
			},
		},
		{
			"phone verify, user not found",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "user1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound"))
				},
			},
		},
		{
			"phone verify, no code",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "user1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"phone verify, already verified",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"id1",
								1,
								json.RawMessage(`{
						"name": "user1"
					}`),
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"+41791234567",
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("phoneverify"),
								},
								time.Hour*1,
								false,
								"",
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneVerifiedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "user1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"phone update, no permission",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"id1",
								1,
								json.RawMessage(`{
						"name": "user1"
					}`),
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"+41791234567",
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("phoneverify"),
								},
								time.Hour*1,
								false,
								"",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "user1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"phone verify, wrong code",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"id1",
								1,
								json.RawMessage(`{
						"name": "user1"
					}`),
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"+41791234567",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							schemauser.NewPhoneCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("phoneverify"),
								},
								time.Hour*1,
								false,
								"",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "user1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "CODE-woT0xc", "Errors.User.Code.Invalid"))
				},
			},
		},
		{
			"phone verify, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPhoneCodeExisting(),
					expectPush(
						eventFromEventPusher(
							schemauser.NewPhoneVerifiedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				id:   "user1",
				code: "phoneverify",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"phone verify, external, not copnfigured",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPhoneCodeExternalExisting(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				id:   "user1",
				code: "phoneverify",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-S8kTrxy0aH", "Errors.User.Code.NotConfigured"))
				},
			},
		},
		{
			"phone verify, external, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPhoneCodeExternalExisting(),
					expectPush(
						eventFromEventPusher(
							schemauser.NewPhoneVerifiedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				phoneCodeVerifier: func(ctx context.Context, id string) (senders.CodeGenerator, error) {
					sender := mock.NewMockCodeGenerator(gomock.NewController(t))
					sender.EXPECT().VerifyCode("verificationID", "code").Return(nil)
					return sender, nil
				},
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				id:   "user1",
				code: "code",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:        tt.fields.eventstore(t),
				checkPermission:   tt.fields.checkPermission,
				userEncryption:    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				phoneCodeVerifier: tt.fields.phoneCodeVerifier,
			}
			details, err := c.VerifySchemaUserPhone(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.code)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_ResendSchemaUserPhoneCode(t *testing.T) {
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		checkPermission             domain.PermissionCheck
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
	}
	type args struct {
		ctx  context.Context
		user *ResendSchemaUserPhoneCode
	}
	type res struct {
		returnCode string
		details    *domain.ObjectDetails
		err        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no userID, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserPhoneCode{
					ID: "",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-zmxIFR2nMo", "Errors.IDMissing"))
				},
			},
		},
		{
			"phone code resend, user not found",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserPhoneCode{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound"))
				},
			},
		},
		{
			"phone code resend, no code",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"id1",
								1,
								json.RawMessage(`{
						"name": "user1"
					}`),
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserPhoneCode{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-fEsHdqECzb", "Errors.User.Code.Empty"))
				},
			},
		},
		{
			"phone code resend, already verified",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"id1",
								1,
								json.RawMessage(`{
						"name": "user1"
					}`),
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"+41791234567",
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("phoneverify"),
								},
								time.Hour*1,
								false,
								"",
							),
						),
						eventFromEventPusher(
							schemauser.NewPhoneVerifiedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserPhoneCode{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-fEsHdqECzb", "Errors.User.Code.Empty"))
				},
			},
		},
		{
			"phone code resend, no permission",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPhoneCodeExisting(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserPhoneCode{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"phone code resend, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPhoneCodeExisting(),
					filterSMSProviderActivated(),
					filterSMSProviderInternalExisting(),
					expectPush(
						eventFromEventPusher(
							schemauser.NewPhoneCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("phoneverify2"),
								},
								time.Hour*1,
								false,
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify2", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserPhoneCode{
					ID: "user1",
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"phone code resend, ok (external)",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPhoneCodeExisting(),
					filterSMSProviderActivated(),
					filterSMSProviderExternalExisting(),
					expectPush(
						eventFromEventPusher(
							schemauser.NewPhoneCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								nil,
								0,
								false,
								"id",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserPhoneCode{
					ID: "user1",
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"phone code resend, return, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPhoneCodeExisting(),
					filterSMSProviderActivated(),
					filterSMSProviderInternalExisting(),
					expectPush(
						eventFromEventPusher(
							schemauser.NewPhoneCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("phoneverify2"),
								},
								time.Hour*1,
								true,
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify2", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserPhoneCode{
					ID:         "user1",
					ReturnCode: true,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				returnCode: "phoneverify2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				checkPermission:             tt.fields.checkPermission,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				userEncryption:              crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				defaultSecretGenerators: &SecretGenerators{
					PhoneVerificationCode: &crypto.GeneratorConfig{
						Length:              8,
						Expiry:              time.Hour,
						IncludeLowerLetters: true,
						IncludeUpperLetters: true,
						IncludeDigits:       true,
						IncludeSymbols:      true,
					},
				},
			}
			details, err := c.ResendSchemaUserPhoneCode(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
				if tt.res.returnCode != "" {
					assert.NotNil(t, tt.args.user.PlainCode)
					assert.Equal(t, tt.res.returnCode, *tt.args.user.PlainCode)
				}
			}
		})
	}
}
