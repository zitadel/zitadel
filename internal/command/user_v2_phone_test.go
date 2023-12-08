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
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommands_ChangeUserPhone(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
		newCode         cryptoCodeFunc
	}
	type args struct {
		userID        string
		resourceOwner string
		phone         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Phone
		wantErr error
	}{
		{
			name: "missing permission",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "",
			},
			wantErr: caos_errs.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing phone",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234567",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Phone.NotChanged"),
		},
		{
			name: "phone changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							time.Hour*1,
							false,
						),
					),
				),
				newCode:         mockCode("phoneCode", time.Hour),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234568",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
				newCode:         tt.fields.newCode,
			}
			got, err := c.ChangeUserPhone(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.phone, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			// successful cases are tested in TestCommands_changeUserPhoneWithGenerator
		})
	}
}

func TestCommands_ChangeUserPhoneReturnCode(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
		newCode         cryptoCodeFunc
	}
	type args struct {
		userID        string
		resourceOwner string
		phone         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Phone
		wantErr error
	}{
		{
			name: "missing permission",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234567",
			},
			wantErr: caos_errs.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing phone",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty"),
		},
		{
			name: "phone changed, return code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							time.Hour*1,
							true,
						),
					),
				),
				newCode:         mockCode("phoneCode", time.Hour),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234568",
			},
			want: &domain.Phone{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				PhoneNumber:     "+41791234568",
				IsPhoneVerified: false,
				PlainCode:       gu.Ptr("phoneCode"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
				newCode:         tt.fields.newCode,
			}
			got, err := c.ChangeUserPhoneReturnCode(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.phone, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_ChangeUserPhoneVerified(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID        string
		resourceOwner string
		phone         string
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
				eventstore:      eventstoreExpect(t),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID:        "",
				resourceOwner: "org1",
				phone:         "+41791234567",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "COMMAND-xP292j", "Errors.User.Phone.IDMissing"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234567",
			},
			wantErr: caos_errs.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing phone",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty"),
		},
		{
			name: "phone changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234568",
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
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.ChangeUserPhoneVerified(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.phone)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_changeUserPhoneWithCode(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
		newCode         cryptoCodeFunc
	}
	type args struct {
		userID        string
		resourceOwner string
		phone         string
		returnCode    bool
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
				eventstore: eventstoreExpect(t),
			},
			args: args{
				userID:        "",
				resourceOwner: "org1",
				phone:         "+41791234567",
				returnCode:    false,
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "COMMAND-xP292j", "Errors.User.Phone.IDMissing"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234567",
				returnCode:    false,
			},
			wantErr: caos_errs.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing phone",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "",
				returnCode:    false,
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234567",
				returnCode:    false,
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Phone.NotChanged"),
		},
		{
			name: "phone changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							time.Hour*1,
							false,
						),
					),
				),
				newCode:         mockCode("phoneCode", time.Hour),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234568",
				returnCode:    false,
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
				eventstore: eventstoreExpect(
					t,
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
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							time.Hour*1,
							true,
						),
					),
				),
				newCode:         mockCode("phoneCode", time.Hour),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				phone:         "+41791234568",
				returnCode:    true,
			},
			want: &domain.Phone{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				PhoneNumber:     "+41791234568",
				IsPhoneVerified: false,
				PlainCode:       gu.Ptr("phoneCode"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
				newCode:         tt.fields.newCode,
			}
			got, err := c.changeUserPhoneWithCode(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.phone, crypto.CreateMockEncryptionAlg(gomock.NewController(t)), tt.args.returnCode)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}
