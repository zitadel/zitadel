package command

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_ChangeUserPhone(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID string
		phone  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyPhoneCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID: "user1",
				phone:  "",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing phone",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyPhoneCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID: "user1",
				phone:  "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyPhoneCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID: "user1",
				phone:  "+41791234567",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Phone.NotChanged"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			_, err := c.ChangeUserPhone(context.Background(), tt.args.userID, tt.args.phone, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_changeUserPhoneWithGenerator
		})
	}
}

func TestCommands_ChangeUserPhoneReturnCode(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID string
		phone  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyPhoneCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID: "user1",
				phone:  "+41791234567",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing phone",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyEmailCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID: "user1",
				phone:  "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			_, err := c.ChangeUserPhoneReturnCode(context.Background(), tt.args.userID, tt.args.phone, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_changeUserPhoneWithGenerator
		})
	}
}

func TestCommands_ResendUserPhoneCode(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyPhoneCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID: "user1",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "no code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyPhoneCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID: "user1",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "PHONE-5xrra88eq8", "Errors.User.Code.Empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			_, err := c.ResendUserPhoneCode(context.Background(), tt.args.userID, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_resendUserPhoneCodeWithGenerator
		})
	}
}

func TestCommands_ResendUserPhoneCodeReturnCode(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyPhoneCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID: "user1",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "no code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyEmailCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID: "user1",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "PHONE-5xrra88eq8", "Errors.User.Code.Empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			_, err := c.ResendUserPhoneCodeReturnCode(context.Background(), tt.args.userID, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_resendUserPhoneCodeWithGenerator
		})
	}
}

func TestCommands_ChangeUserPhoneVerified(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID string
		phone  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Phone
		wantErr error
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID: "",
				phone:  "+41791234567",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-xP292j", "Errors.User.Phone.IDMissing"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID: "user1",
				phone:  "+41791234567",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing phone",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID: "user1",
				phone:  "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty"),
		},
		{
			name: "phone changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"+41791234568",
						),
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID: "user1",
				phone:  "+41791234568",
			},
			want: &domain.Phone{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				PhoneNumber:     domain.PhoneNumber("+41791234568"),
				IsPhoneVerified: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.ChangeUserPhoneVerified(context.Background(), tt.args.userID, tt.args.phone)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_changeUserPhoneWithGenerator(t *testing.T) {
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
		checkPermission             domain.PermissionCheck
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		userEncryption              crypto.EncryptionAlgorithm
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		userID     string
		phone      string
		returnCode bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Phone
		wantErr error
	}{
		{
			name: "missing user",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				userID:     "",
				phone:      "+41791234567",
				returnCode: false,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-xP292j", "Errors.User.Phone.IDMissing"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID:     "user1",
				phone:      "+41791234567",
				returnCode: false,
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing phone",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:     "user1",
				phone:      "",
				returnCode: false,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:     "user1",
				phone:      "+41791234567",
				returnCode: false,
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Phone.NotChanged"),
		},
		{
			name: "phone changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
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
							"+41791234568",
						),
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							false,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("a", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				userID:     "user1",
				phone:      "+41791234568",
				returnCode: false,
			},
			want: &domain.Phone{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				PhoneNumber:     "+41791234568",
				IsPhoneVerified: false,
			},
		},
		{
			name: "phone changed, return code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
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
							"+41791234568",
						),
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							true,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("a", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				userID:     "user1",
				phone:      "+41791234568",
				returnCode: true,
			},
			want: &domain.Phone{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				PhoneNumber:     "+41791234568",
				IsPhoneVerified: false,
				PlainCode:       gu.Ptr("a"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				checkPermission:             tt.fields.checkPermission,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				userEncryption:              tt.fields.userEncryption,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
			}
			got, err := c.changeUserPhoneWithGenerator(context.Background(), tt.args.userID, tt.args.phone, GetMockSecretGenerator(t), tt.args.returnCode)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_resendUserPhoneCodeWithGenerator(t *testing.T) {
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
		checkPermission             domain.PermissionCheck
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		userID     string
		returnCode bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Phone
		wantErr error
	}{
		{
			name: "missing user",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				userID:     "",
				returnCode: false,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-xP292j", "Errors.User.Phone.IDMissing"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID:     "user1",
				returnCode: false,
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "no code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:     "user1",
				returnCode: false,
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "PHONE-5xrra88eq8", "Errors.User.Code.Empty"),
		},
		{
			name: "resend",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEventV2(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								true,
								"",
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
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							false,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("a", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				userID:     "user1",
				returnCode: false,
			},
			want: &domain.Phone{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				PhoneNumber:     "+41791234567",
				IsPhoneVerified: false,
			},
		},
		{
			name: "resend (external)",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEventV2(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil,
								0,
								true,
								"id",
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
					),
					expectPush(
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							nil,
							0,
							false,
							"id",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("a", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				userID:     "user1",
				returnCode: false,
			},
			want: &domain.Phone{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				PhoneNumber:     "+41791234567",
				IsPhoneVerified: false,
			},
		},
		{
			name: "return code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								event := user.NewHumanAddedEvent(context.Background(),
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
								)
								event.AddPhoneData("+41791234567")
								return event
							}(),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEventV2(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								true,
								"",
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
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							true,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("a", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				userID:     "user1",
				returnCode: true,
			},
			want: &domain.Phone{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				PhoneNumber:     "+41791234567",
				IsPhoneVerified: false,
				PlainCode:       gu.Ptr("a"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				checkPermission:             tt.fields.checkPermission,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
			}
			got, err := c.resendUserPhoneCodeWithGenerator(context.Background(), tt.args.userID, GetMockSecretGenerator(t), tt.args.returnCode)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}
