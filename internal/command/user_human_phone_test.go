package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/senders/mock"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeHumanPhone(t *testing.T) {
	defaultGenerators := &SecretGenerators{
		OTPSMS: &crypto.GeneratorConfig{
			Length:              8,
			Expiry:              time.Hour,
			IncludeLowerLetters: true,
			IncludeUpperLetters: true,
			IncludeDigits:       true,
			IncludeSymbols:      true,
		},
	}
	type fields struct {
		eventstore                  func(*testing.T) *eventstore.Eventstore
		userEncryption              crypto.EncryptionAlgorithm
		defaultSecretGenerators     *SecretGenerators
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
	}
	type args struct {
		ctx             context.Context
		email           *domain.Phone
		resourceOwner   string
		secretGenerator crypto.Generator
	}
	type res struct {
		want *domain.Phone
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid phone, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber: "+41711234567",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "phone not changed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+41711234567",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber: "+41711234567",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "verified phone changed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+41711234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"+41719876543",
						),
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber:     "+41719876543",
					IsPhoneVerified: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber:     "+41719876543",
					IsPhoneVerified: true,
				},
			},
		},
		{
			name: "phone changed to verified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+41711234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber:     "+41711234567",
					IsPhoneVerified: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber:     "+41711234567",
					IsPhoneVerified: true,
				},
			},
		},
		{
			name: "phone changed to verified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+41711234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber:     "+41711234567",
					IsPhoneVerified: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber:     "+41711234567",
					IsPhoneVerified: true,
				},
			},
		},
		{
			name: "phone changed with code, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
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
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"+41711234567",
						),
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							"",
						),
					),
				),
				userEncryption:              crypto.NewMockEncryptionAlgorithm(gomock.NewController(t)),
				defaultSecretGenerators:     defaultGenerators,
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("a", time.Hour),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber: "+41711234567",
				},
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber: "+41711234567",
				},
			},
		},
		{
			name: "phone changed with code (external), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
								"",
								"sid",
								"senderNumber",
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
								"verifyServiceSID",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"+41711234567",
						),
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							nil,
							0,
							"id",
						),
					),
				),
				userEncryption:          crypto.NewMockEncryptionAlgorithm(gomock.NewController(t)),
				defaultSecretGenerators: defaultGenerators,
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber: "+41711234567",
				},
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber: "+41711234567",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				userEncryption:              tt.fields.userEncryption,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
			}
			got, err := r.ChangeHumanPhone(tt.args.ctx, tt.args.email, tt.args.resourceOwner, tt.args.secretGenerator)
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

func TestCommandSide_VerifyHumanPhone(t *testing.T) {
	type fields struct {
		eventstore        func(*testing.T) *eventstore.Eventstore
		phoneCodeVerifier func(ctx context.Context, id string) (senders.CodeGenerator, error)
	}
	type args struct {
		ctx             context.Context
		userID          string
		code            string
		resourceOwner   string
		secretGenerator crypto.Generator
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
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				code:          "aa",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "code missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				code:          "aa",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "code not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				code:          "aa",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "invalid code, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneVerificationFailedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:             context.Background(),
				userID:          "user1",
				code:            "test",
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "valid code, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPhoneCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:             context.Background(),
				userID:          "user1",
				code:            "a",
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "valid code (external), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPhoneCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil,
								0,
								"id",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPhoneCodeSentEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&senders.CodeGeneratorInfo{
									ID:             "id",
									VerificationID: "verificationID",
								},
							),
						),
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				phoneCodeVerifier: func(ctx context.Context, id string) (senders.CodeGenerator, error) {
					sender := mock.NewMockCodeGenerator(gomock.NewController(t))
					sender.EXPECT().VerifyCode("verificationID", "a")
					return sender, nil
				},
			},
			args: args{
				ctx:             context.Background(),
				userID:          "user1",
				code:            "a",
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:        tt.fields.eventstore(t),
				phoneCodeVerifier: tt.fields.phoneCodeVerifier,
			}
			got, err := r.VerifyHumanPhone(tt.args.ctx, tt.args.userID, tt.args.code, tt.args.resourceOwner, tt.args.secretGenerator)
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

func TestCommandSide_CreateVerificationCodeHumanPhone(t *testing.T) {
	defaultGenerators := &SecretGenerators{
		OTPSMS: &crypto.GeneratorConfig{
			Length:              8,
			Expiry:              time.Hour,
			IncludeLowerLetters: true,
			IncludeUpperLetters: true,
			IncludeDigits:       true,
			IncludeSymbols:      true,
		},
	}
	type fields struct {
		eventstore                  func(*testing.T) *eventstore.Eventstore
		userEncryption              crypto.EncryptionAlgorithm
		defaultSecretGenerators     *SecretGenerators
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
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
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "phone already verified, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "new code, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
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
					),
					expectPush(
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
							time.Hour*1,
							"",
						),
					),
				),
				userEncryption:              crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				defaultSecretGenerators:     defaultGenerators,
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "new code (external), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
								"",
								"sid",
								"senderNumber",
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
								"verifyServiceSID",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							nil,
							0,
							"id",
						),
					),
				),
				userEncryption:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				defaultSecretGenerators: defaultGenerators,
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				userEncryption:              tt.fields.userEncryption,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
			}
			got, err := r.CreateHumanPhoneVerificationCode(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
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

func TestCommandSide_PhoneVerificationCodeSent(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		generatorInfo *senders.CodeGeneratorInfo
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
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "code sent, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneCodeSentEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&senders.CodeGeneratorInfo{},
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				generatorInfo: &senders.CodeGeneratorInfo{},
			},
			res: res{},
		},
		{
			name: "code sent (external), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneCodeSentEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&senders.CodeGeneratorInfo{
								ID:             "generatorID",
								VerificationID: "verificationID",
							},
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				generatorInfo: &senders.CodeGeneratorInfo{
					ID:             "generatorID",
					VerificationID: "verificationID",
				},
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := r.HumanPhoneVerificationCodeSent(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.generatorInfo)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_RemoveHumanPhone(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
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
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "phone not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove phone, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.RemoveHumanPhone(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
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
