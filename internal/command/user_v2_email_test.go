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

func TestCommands_ChangeUserEmail(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
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
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing email",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Email.NotChanged"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
			}
			_, err := c.ChangeUserEmail(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_changeUserEmailWithGenerator
		})
	}
}

func TestCommands_ChangeUserEmailURLTemplate(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
		urlTmpl       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "invalid template",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
				urlTmpl:       "{{",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "permission missing",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
				urlTmpl:       "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
				urlTmpl:       "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Email.NotChanged"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
			}
			_, err := c.ChangeUserEmailURLTemplate(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email, crypto.CreateMockEncryptionAlg(gomock.NewController(t)), tt.args.urlTmpl)
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_changeUserEmailWithGenerator
		})
	}
}

func TestCommands_ChangeUserEmailReturnCode(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
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
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing email",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
			}
			_, err := c.ChangeUserEmailReturnCode(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_changeUserEmailWithGenerator
		})
	}
}

func TestCommands_ChangeUserEmailVerified(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Email
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
				email:         "email@test.ch",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing email",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"),
		},
		{
			name: "email changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"email-changed@test.ch",
						),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
			},
			want: &domain.Email{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				EmailAddress:    "email-changed@test.ch",
				IsEmailVerified: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.ChangeUserEmailVerified(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_changeUserEmailWithGenerator(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
		returnCode    bool
		urlTmpl       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Email
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
				email:         "email@test.ch",
				returnCode:    false,
				urlTmpl:       "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
				returnCode:    false,
				urlTmpl:       "",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "missing email",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
				returnCode:    false,
				urlTmpl:       "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
				returnCode:    false,
				urlTmpl:       "",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Email.NotChanged"),
		},
		{
			name: "email changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"email-changed@test.ch",
						),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							"", false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
				returnCode:    false,
				urlTmpl:       "",
			},
			want: &domain.Email{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				EmailAddress:    "email-changed@test.ch",
				IsEmailVerified: false,
			},
		},
		{
			name: "email changed, return code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"email-changed@test.ch",
						),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							"", true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
				returnCode:    true,
				urlTmpl:       "",
			},
			want: &domain.Email{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				EmailAddress:    "email-changed@test.ch",
				IsEmailVerified: false,
				PlainCode:       gu.Ptr("a"),
			},
		},
		{
			name: "email changed, URL template",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"email-changed@test.ch",
						),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}", false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
				returnCode:    false,
				urlTmpl:       "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
			},
			want: &domain.Email{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				EmailAddress:    "email-changed@test.ch",
				IsEmailVerified: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.changeUserEmailWithGenerator(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email, GetMockSecretGenerator(t), tt.args.returnCode, tt.args.urlTmpl)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_VerifyUserEmail(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		code          string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyEmailCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
				),
			},
			args: args{
				userID:        "",
				resourceOwner: "org1",
				code:          "a",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing"),
		},
		{
			name: "missing code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				code:          "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-Fia4a", "Errors.User.Code.Empty"),
		},
		{
			name: "wrong code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							user.NewHumanEmailCodeAddedEventV2(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"", false,
							),
						),
					),
					expectPush(
						user.NewHumanEmailVerificationFailedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				code:          "wrong",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-eis9R", "Errors.User.Code.Invalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			_, err := c.VerifyUserEmail(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.code, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_verifyUserEmailWithGenerator
		})
	}
}

func TestCommands_verifyUserEmailWithGenerator(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		code          string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr error
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				userID:        "",
				resourceOwner: "org1",
				code:          "a",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing"),
		},
		{
			name: "missing code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				userID:        "user1",
				resourceOwner: "org1",
				code:          "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-Fia4a", "Errors.User.Code.Empty"),
		},
		{
			name: "good code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							user.NewHumanEmailCodeAddedEventV2(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"", false,
							),
						),
					),
					expectPush(
						user.NewHumanEmailVerificationFailedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				code:          "wrong",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-eis9R", "Errors.User.Code.Invalid"),
		},
		{
			name: "wrong code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							user.NewHumanEmailCodeAddedEventV2(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"", false,
							),
						),
					),
					expectPush(
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				code:          "a",
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.verifyUserEmailWithGenerator(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.code, GetMockSecretGenerator(t))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_NewUserEmailEvents(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				userID:        "",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing"),
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t, expectFilter()),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-ieJ2e", "Errors.User.Email.NotFound"),
		},
		{
			name: "user not initialized",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil, time.Hour*1,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-uz0Uu", "Errors.User.NotInitialised"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			_, err := c.NewUserEmailEvents(context.Background(), tt.args.userID, tt.args.resourceOwner)
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_changeUserEmailWithGenerator
		})
	}
}
