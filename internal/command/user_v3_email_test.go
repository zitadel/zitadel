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
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_ChangeSchemaUserEmail(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
		newCode         encrypedCodeFunc
	}
	type args struct {
		ctx  context.Context
		user *ChangeSchemaUserEmail
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
				user: &ChangeSchemaUserEmail{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-0oj2PquNGA", "Errors.IDMissing"))
				},
			},
		},
		{
			"no valid email, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserEmail{
					ID:    "user1",
					Email: &Email{Address: "noemail"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "EMAIL-599BI", "Errors.User.Email.Invalid"))
				},
			},
		},
		{
			"no valid template, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserEmail{
					ID:    "user1",
					Email: &Email{Address: "noemail", URLTemplate: "{{"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "EMAIL-599BI", "Errors.User.Email.Invalid"))
				},
			},
		},
		{
			"email update, user not found",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserEmail{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-nJ0TQFuRmP", "Errors.User.NotFound"))
				},
			},
		},
		{
			"email update, no permission",
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
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserEmail{
					ID:    "user1",
					Email: &Email{Address: "noemail@example.com"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"email update, email not changed",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"test@example.com",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserEmail{
					ID: "user1",
					Email: &Email{
						Address:     "test@example.com",
						ReturnCode:  true,
						URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
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
			"email update, email return",
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
					expectPush(
						schemauser.NewEmailUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"test@example.com",
						),
						schemauser.NewEmailCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailverify"),
							},
							time.Hour*1,
							"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("emailverify", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserEmail{
					ID: "user1",
					Email: &Email{
						Address:     "test@example.com",
						ReturnCode:  true,
						URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				returnCode: "emailverify",
			},
		},
		{
			"user updated, email to verify",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						)),
					expectPush(
						schemauser.NewEmailUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"test@example.com",
						),
						schemauser.NewEmailCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailverify"),
							},
							time.Hour*1,
							"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
							false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("emailverify", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserEmail{
					ID: "user1",
					Email: &Email{
						Address:     "test@example.com",
						URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
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
					expectFilter(
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
					),
					expectPush(
						schemauser.NewEmailUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							"test@example.com",
						),
						schemauser.NewEmailVerifiedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUserEmail{
					ID:    "user1",
					Email: &Email{Address: "test@example.com", Verified: true},
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
				eventstore:       tt.fields.eventstore(t),
				checkPermission:  tt.fields.checkPermission,
				newEncryptedCode: tt.fields.newCode,
				userEncryption:   crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			}
			details, err := c.ChangeSchemaUserEmail(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
				if tt.res.returnCode != "" {
					assert.NotNil(t, tt.args.user.ReturnCode)
					assert.Equal(t, tt.res.returnCode, *tt.args.user.ReturnCode)
				}
			}
		})
	}
}

func TestCommands_VerifySchemaUserEmail(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-y3n4Sdu8j5", "Errors.IDMissing"))
				},
			},
		},
		{
			"email verify, user not found",
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-qbGyMPvjvj", "Errors.User.NotFound"))
				},
			},
		},
		{
			"email verify, no code",
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
				id:  "user1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"email verify, already verified",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"test@example.com",
							),
						),
						eventFromEventPusher(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
							),
						),
						eventFromEventPusher(
							schemauser.NewEmailVerifiedEvent(
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
			"email update, no permission",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"test@example.com",
							),
						),
						eventFromEventPusher(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
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
			"email verify, wrong code",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"test@example.com",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
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
			"email verify, ok",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"test@example.com",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							schemauser.NewEmailVerifiedEvent(
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
				code: "emailverify",
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
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
				userEncryption:  crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			}
			details, err := c.VerifySchemaUserEmail(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.code)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_ResendSchemaUserEmailCode(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
		newCode         encrypedCodeFunc
	}
	type args struct {
		ctx  context.Context
		user *ResendSchemaUserEmailCode
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
				user: &ResendSchemaUserEmailCode{
					ID: "",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-KvPc5o9GeJ", "Errors.IDMissing"))
				},
			},
		},
		{
			"email code resend, user not found",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserEmailCode{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-EajeF6ypOV", "Errors.User.NotFound"))
				},
			},
		},
		{
			"email code resend, no code",
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
				user: &ResendSchemaUserEmailCode{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-QRkNTBwF8q", "Errors.User.Code.Empty"))
				},
			},
		},
		{
			"email code resend, already verified",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"test@example.com",
							),
						),
						eventFromEventPusher(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
							),
						),
						eventFromEventPusher(
							schemauser.NewEmailVerifiedEvent(
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
				user: &ResendSchemaUserEmailCode{
					ID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-QRkNTBwF8q", "Errors.User.Code.Empty"))
				},
			},
		},
		{
			"email code resend, no permission",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"test@example.com",
							),
						),
						eventFromEventPusher(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserEmailCode{
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
			"email code resend, ok",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"test@example.com",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify2"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("emailverify2", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserEmailCode{
					ID:          "user1",
					URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"email code resend, return, ok",
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
							schemauser.NewEmailUpdatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"test@example.com",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								false,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							schemauser.NewEmailCodeAddedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("emailverify2"),
								},
								time.Hour*1,
								"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
								true,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("emailverify2", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ResendSchemaUserEmailCode{
					ID:          "user1",
					URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
					ReturnCode:  true,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				returnCode: "emailverify2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:       tt.fields.eventstore(t),
				checkPermission:  tt.fields.checkPermission,
				newEncryptedCode: tt.fields.newCode,
				userEncryption:   crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			}
			details, err := c.ResendSchemaUserEmailCode(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
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
