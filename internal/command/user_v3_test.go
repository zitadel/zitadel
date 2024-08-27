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
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateSchemaUser(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		idGenerator     id.Generator
		checkPermission domain.PermissionCheck
		newCode         encrypedCodeFunc
	}
	type args struct {
		ctx  context.Context
		user *CreateSchemaUser
	}
	type res struct {
		returnCodeEmail string
		returnCodePhone string
		details         *domain.ObjectDetails
		err             func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no resourceOwner, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-urEJKa1tJM", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			"no schemaID, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner: "org1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-TFo06JgnF2", "Errors.UserSchema.ID.Missing"))
				},
			},
		},
		{
			"schema not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-N9QOuN4F7o", "Errors.UserSchema.NotExists"))
				},
			},
		},
		{
			"no data, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "id1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-7o3ZGxtXUz", "Errors.User.Invalid"))
				},
			},
		},
		{
			"user create, no permission",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
				idGenerator:     mock.ExpectID(t, "id1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"user created",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectFilter(),
					expectPush(
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "id1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "id1",
				},
			},
		},
		{
			"user create, no field permission as admin",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
									 	"urn:zitadel:schema:permission": {
											"owner": "r"
										},
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "id1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid"))
				},
			},
		},
		{
			"user create, no field permission as user",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
									 	"urn:zitadel:schema:permission": {
											"self": "r"
										},
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "org1", "id1"),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid"))
				},
			},
		},
		{
			"user create, invalid data type",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "id1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "org1", "user1"),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": 1
					}`),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid"))
				},
			},
		},
		{
			"user create, invalid data attribute name",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								},
       							"additionalProperties": false
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "id1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "org1", "user1"),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"invalid": "user"
					}`),
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid"))
				},
			},
		},
		{
			"user created, email return",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectFilter(),
					expectPush(
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
						schemauser.NewEmailUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"test@example.com",
						),
						schemauser.NewEmailCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
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
				idGenerator:     mock.ExpectID(t, "id1"),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("emailverify", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
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
					ID:            "id1",
				},
				returnCodeEmail: "emailverify",
			},
		},
		{
			"user created, email to verify",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectFilter(),
					expectPush(
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
						schemauser.NewEmailUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"test@example.com",
						),
						schemauser.NewEmailCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
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
				idGenerator:     mock.ExpectID(t, "id1"),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("emailverify", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
					Email: &Email{
						Address:     "test@example.com",
						URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "id1",
				},
			},
		},
		{
			"user created, phone return",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectFilter(),
					expectPush(
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
						schemauser.NewPhoneChangedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"+41791234567",
						),
						schemauser.NewPhoneCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneverify"),
							},
							time.Hour*1,
							true,
						),
					),
				),
				idGenerator:     mock.ExpectID(t, "id1"),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("phoneverify", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
					Phone: &Phone{
						Number:     "+41791234567",
						ReturnCode: true,
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "id1",
				},
				returnCodePhone: "phoneverify",
			},
		},
		{
			"user created, phone to verify",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectFilter(),
					expectPush(
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
						schemauser.NewPhoneChangedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"+41791234567",
						),
						schemauser.NewPhoneCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneverify"),
							},
							time.Hour*1,
							false,
						),
					),
				),
				idGenerator:     mock.ExpectID(t, "id1"),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("phoneverify", time.Hour),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
					Phone: &Phone{
						Number: "+41791234567",
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "id1",
				},
			},
		},
		{
			"user created, full verified",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectFilter(),
					expectPush(
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
						schemauser.NewEmailUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"test@example.com",
						),
						schemauser.NewEmailVerifiedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
						),
						schemauser.NewPhoneChangedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"+41791234567",
						),
						schemauser.NewPhoneVerifiedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
						),
					),
				),
				idGenerator:     mock.ExpectID(t, "id1"),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &CreateSchemaUser{
					ResourceOwner:  "org1",
					SchemaID:       "type",
					schemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
					Email: &Email{Address: "test@example.com", Verified: true},
					Phone: &Phone{Number: "+41791234567", Verified: true},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "id1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:       tt.fields.eventstore(t),
				idGenerator:      tt.fields.idGenerator,
				checkPermission:  tt.fields.checkPermission,
				newEncryptedCode: tt.fields.newCode,
			}
			err := c.CreateSchemaUser(tt.args.ctx, tt.args.user, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, tt.args.user.Details)
			}

			if tt.res.returnCodePhone != "" {
				assert.Equal(t, tt.res.returnCodePhone, tt.args.user.ReturnCodePhone)
			}
			if tt.res.returnCodeEmail != "" {
				assert.Equal(t, tt.res.returnCodeEmail, tt.args.user.ReturnCodeEmail)
			}
		})
	}
}

func TestCommandSide_DeleteSchemaUser(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx    context.Context
			userID string
		}
	)
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
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-Vs4wJCME7T", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"schema",
								1,
								json.RawMessage(`{
						"name": "user"
					}`),
							),
						),
						eventFromEventPusher(
							schemauser.NewDeletedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "remove user, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"schema",
								1,
								json.RawMessage(`{
						"name": "user"
					}`),
							),
						),
					),
					expectPush(
						schemauser.NewDeletedEvent(context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "remove user, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"schema",
								1,
								json.RawMessage(`{
						"name": "user"
					}`),
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "remove user, self",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"schema",
								1,
								json.RawMessage(`{
						"name": "user"
					}`),
							),
						),
					),
					expectPush(
						schemauser.NewDeletedEvent(authz.NewMockContext("instanceID", "org1", "user1"),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    authz.NewMockContext("instanceID", "org1", "user1"),
				userID: "user1",
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
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.DeleteSchemaUser(tt.args.ctx, tt.args.userID)
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
