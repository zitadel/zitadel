package command

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	webauthn_helper "github.com/zitadel/zitadel/internal/webauthn"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_RegisterUserPasskey(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	ctx = authz.WithRequestedDomain(ctx, "example.com")

	webauthnConfig := &webauthn_helper.Config{
		DisplayName:    "test",
		ExternalSecure: true,
	}
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
		rpID          string
		authenticator domain.AuthenticatorAttachment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.WebAuthNRegistrationDetails
		wantErr error
	}{
		{
			name: "wrong user",
			args: args{
				userID:        "foo",
				resourceOwner: "org1",
				authenticator: domain.AuthenticatorAttachmentCrossPlattform,
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTH-Bohd2", "Errors.User.UserIDWrong"),
		},
		{
			name: "get human passwordless error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				authenticator: domain.AuthenticatorAttachmentCrossPlattform,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "id generator error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(), // getHumanPasswordlessTokens
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(ctx,
							userAgg,
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
					)),
					expectFilter(eventFromEventPusher(
						org.NewOrgAddedEvent(ctx,
							&org.NewAggregate("org1").Aggregate,
							"org1",
						),
					)),
					expectFilter(eventFromEventPusher(
						org.NewDomainPolicyAddedEvent(ctx,
							&org.NewAggregate("org1").Aggregate,
							false, false, false,
						),
					)),
				),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				authenticator: domain.AuthenticatorAttachmentCrossPlattform,
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:     tt.fields.eventstore,
				webauthnConfig: webauthnConfig,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			_, err := c.RegisterUserPasskey(ctx, tt.args.userID, tt.args.resourceOwner, tt.args.rpID, tt.args.authenticator)
			require.ErrorIs(t, err, tt.wantErr)
			// successful case can't be tested due to random challenge.
		})
	}
}

func TestCommands_RegisterUserPasskeyWithCode(t *testing.T) {
	ctx := authz.WithRequestedDomain(context.Background(), "example.com")
	webauthnConfig := &webauthn_helper.Config{
		DisplayName:    "test",
		ExternalSecure: true,
	}
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	es := eventstoreExpect(t,
		expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
	)
	code, err := newEncryptedCode(ctx, es.Filter, domain.SecretGeneratorTypePasswordlessInitCode, alg) //nolint:staticcheck
	require.NoError(t, err)
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
		rpID          string
		authenticator domain.AuthenticatorAttachment
		codeID        string
		code          string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "code verification error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								userAgg, "123", code.Crypted, time.Minute, "", false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "123"),
						),
					),
					expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
					expectPush(
						user.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, "123"),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				authenticator: domain.AuthenticatorAttachmentCrossPlattform,
				codeID:        "123",
				code:          "wrong",
			},
			wantErr: zerrors.ThrowInvalidArgument(err, "COMMAND-Eeb2a", "Errors.User.Code.Invalid"),
		},
		{
			name: "code verification ok, get human passwordless error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								userAgg, "123", code.Crypted, time.Minute, "", false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "123"),
						),
					),
					expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				authenticator: domain.AuthenticatorAttachmentCrossPlattform,
				codeID:        "123",
				code:          code.Plain,
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:     tt.fields.eventstore,
				webauthnConfig: webauthnConfig,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			_, err := c.RegisterUserPasskeyWithCode(ctx, tt.args.userID, tt.args.resourceOwner, tt.args.authenticator, tt.args.codeID, tt.args.code, tt.args.rpID, alg)
			require.ErrorIs(t, err, tt.wantErr)
			// successful case can't be tested due to random challenge.
		})
	}
}

func TestCommands_verifyUserPasskeyCode(t *testing.T) {
	ctx := authz.WithRequestedDomain(context.Background(), "example.com")
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	es := eventstoreExpect(t,
		expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
	)
	code, err := newEncryptedCode(ctx, es.Filter, domain.SecretGeneratorTypePasswordlessInitCode, alg) //nolint:staticcheck
	require.NoError(t, err)
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		codeID        string
		code          string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.HumanPasswordlessInitCodeCheckSucceededEvent
		wantErr error
	}{
		{
			name: "filter error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				codeID:        "123",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "code verification error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								userAgg, "123", code.Crypted, time.Minute, "", false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "123"),
						),
					),
					expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
					expectPush(
						user.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, "123"),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				codeID:        "123",
				code:          "wrong",
			},
			wantErr: zerrors.ThrowInvalidArgument(err, "COMMAND-Eeb2a", "Errors.User.Code.Invalid"),
		},
		{
			name: "success",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								userAgg, "123", code.Crypted, time.Minute, "", false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "123"),
						),
					),
					expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				codeID:        "123",
				code:          code.Plain,
			},
			want: user.NewHumanPasswordlessInitCodeCheckSucceededEvent(ctx, userAgg, "123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.verifyUserPasskeyCode(ctx, tt.args.userID, tt.args.resourceOwner, tt.args.codeID, tt.args.code, alg)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want, got(ctx, userAgg))
			}
		})
	}
}

func TestCommands_pushUserPasskey(t *testing.T) {
	ctx := authz.WithRequestedDomain(authz.NewMockContext("instance1", "org1", "user1"), "example.com")
	webauthnConfig := &webauthn_helper.Config{
		DisplayName:    "test",
		ExternalSecure: true,
	}
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	prep := []expect{
		expectFilter(), // getHumanPasswordlessTokens
		expectFilter(eventFromEventPusher(
			user.NewHumanAddedEvent(ctx,
				userAgg,
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
		)),
		expectFilter(eventFromEventPusher(
			org.NewOrgAddedEvent(ctx,
				&org.NewAggregate("org1").Aggregate,
				"org1",
			),
		)),
		expectFilter(eventFromEventPusher(
			org.NewDomainPolicyAddedEvent(ctx,
				&org.NewAggregate("org1").Aggregate,
				false, false, false,
			),
		)),
		expectFilter(eventFromEventPusher(
			user.NewHumanWebAuthNAddedEvent(eventstore.NewBaseEventForPush(
				ctx, &org.NewAggregate("org1").Aggregate, user.HumanPasswordlessTokenAddedType,
			), "111", "challenge", "rpID"),
		)),
	}

	type args struct {
		events []eventCallback
	}
	tests := []struct {
		name       string
		expectPush func(challenge string) expect
		args       args
		wantErr    error
	}{
		{
			name: "push error",
			expectPush: func(challenge string) expect {
				return expectPushFailed(io.ErrClosedPipe,
					user.NewHumanPasswordlessAddedEvent(ctx,
						userAgg, "123", challenge, "rpID",
					),
				)
			},
			args:    args{},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			expectPush: func(challenge string) expect {
				return expectPush(
					user.NewHumanPasswordlessAddedEvent(ctx,
						userAgg, "123", challenge, "rpID",
					),
				)
			},
			args: args{},
		},
		{
			name: "initcode succeeded event",
			expectPush: func(challenge string) expect {
				return expectPush(
					user.NewHumanPasswordlessAddedEvent(ctx,
						userAgg, "123", challenge, "rpID",
					),
					user.NewHumanPasswordlessInitCodeCheckSucceededEvent(ctx, userAgg, "123"),
				)
			},
			args: args{
				events: []eventCallback{func(ctx context.Context, userAgg *eventstore.Aggregate) eventstore.Command {
					return user.NewHumanPasswordlessInitCodeCheckSucceededEvent(ctx, userAgg, "123")
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:     eventstoreExpect(t, prep...),
				webauthnConfig: webauthnConfig,
			}
			id_generator.SetGenerator(id_mock.NewIDGeneratorExpectIDs(t, "123"))
			wm, userAgg, webAuthN, err := c.createUserPasskey(ctx, "user1", "org1", "rpID", domain.AuthenticatorAttachmentCrossPlattform)
			require.NoError(t, err)

			c.eventstore = eventstoreExpect(t, tt.expectPush(webAuthN.Challenge))

			got, err := c.pushUserPasskey(ctx, wm, userAgg, webAuthN, tt.args.events...)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.NotEmpty(t, got.PublicKeyCredentialCreationOptions)
				assert.Equal(t, "123", got.ID)
				assert.Equal(t, "org1", got.ObjectDetails.ResourceOwner)
			}
		})
	}
}

func TestCommands_AddUserPasskeyCode(t *testing.T) {
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	type fields struct {
		newCode     encrypedCodeFunc
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr error
	}{
		{
			name: "id generator error",
			fields: fields{
				newCode:     mockEncryptedCode("passkey1", time.Hour),
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				newCode: mockEncryptedCode("passkey1", time.Minute),
				eventstore: expectEventstore(
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							userAgg,
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
					)),
					expectPush(
						user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
							userAgg,
							"123", &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("passkey1"),
							}, time.Minute, "", false,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				newEncryptedCode: tt.fields.newCode,
				eventstore:       tt.fields.eventstore(t),
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.AddUserPasskeyCode(context.Background(), tt.args.userID, tt.args.resourceOwner, alg)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_AddUserPasskeyCodeURLTemplate(t *testing.T) {
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	type fields struct {
		newCode     encrypedCodeFunc
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
		urlTmpl       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr error
	}{
		{
			name: "template error",
			fields: fields{
				newCode:    newEncryptedCode,
				eventstore: eventstoreExpect(t),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				urlTmpl:       "{{",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "id generator error",
			fields: fields{
				newCode:     newEncryptedCode,
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				urlTmpl:       "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				newCode: mockEncryptedCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							userAgg,
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
					)),
					expectPush(
						user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
							userAgg,
							"123", &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("passkey1"),
							},
							time.Minute,
							"https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}",
							false,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				urlTmpl:       "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}",
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				newEncryptedCode: tt.fields.newCode,
				eventstore:       tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.AddUserPasskeyCodeURLTemplate(context.Background(), tt.args.userID, tt.args.resourceOwner, alg, tt.args.urlTmpl)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_AddUserPasskeyCodeReturn(t *testing.T) {
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	type fields struct {
		newCode     encrypedCodeFunc
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.PasskeyCodeDetails
		wantErr error
	}{
		{
			name: "id generator error",
			fields: fields{
				newCode:     newEncryptedCode,
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				newCode: mockEncryptedCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							userAgg,
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
					)),
					expectPush(
						user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
							userAgg,
							"123", &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("passkey1"),
							}, time.Minute, "", true,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			want: &domain.PasskeyCodeDetails{
				ObjectDetails: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				CodeID: "123",
				Code:   "passkey1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				newEncryptedCode: tt.fields.newCode,
				eventstore:       tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.AddUserPasskeyCodeReturn(context.Background(), tt.args.userID, tt.args.resourceOwner, alg)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_addUserPasskeyCode(t *testing.T) {
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	type fields struct {
		newCode     encrypedCodeFunc
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.PasskeyCodeDetails
		wantErr error
	}{
		{
			name: "id generator error",
			fields: fields{
				newCode:     newEncryptedCode,
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "crypto error",
			fields: fields{
				newCode:     newEncryptedCode,
				eventstore:  eventstoreExpect(t, expectFilterError(io.ErrClosedPipe)),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "filter query error",
			fields: fields{
				newCode: newEncryptedCode,
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
					expectFilterError(io.ErrClosedPipe),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "push error",
			fields: fields{
				newCode: mockEncryptedCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							userAgg,
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
					)),
					expectPushFailed(io.ErrClosedPipe,
						user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"123", &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("passkey1"),
							}, time.Minute, "", false,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				newCode: mockEncryptedCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							userAgg,
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
					)),
					expectPush(
						user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
							userAgg,
							"123", &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("passkey1"),
							}, time.Minute, "", false,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			want: &domain.PasskeyCodeDetails{
				ObjectDetails: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				CodeID: "123",
				Code:   "passkey1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				newEncryptedCode: tt.fields.newCode,
				eventstore:       tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.addUserPasskeyCode(context.Background(), tt.args.userID, tt.args.resourceOwner, alg, "", false)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
