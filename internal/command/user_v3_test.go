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
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateSchemaUser(t *testing.T) {
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
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		idGenerator                 id.Generator
		checkPermission             domain.PermissionCheck
		newCode                     encrypedCodeFunc
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
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
			"user created, additional property",
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
						"additional": "property"
					}`),
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
						"additional": "property"
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
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
						schemauser.NewPhoneUpdatedEvent(context.Background(),
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
							"",
						),
					),
				),
				idGenerator:                 mock.ExpectID(t, "id1"),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
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
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
						schemauser.NewPhoneUpdatedEvent(context.Background(),
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
							"",
						),
					),
				),
				idGenerator:                 mock.ExpectID(t, "id1"),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
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
			"user created, phone to verify (external)",
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
						schemauser.NewCreatedEvent(
							context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"type",
							1,
							json.RawMessage(`{
						"name": "user"
					}`),
						),
						schemauser.NewPhoneUpdatedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"+41791234567",
						),
						schemauser.NewPhoneCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							nil,
							0,
							false,
							"id",
						),
					),
				),
				idGenerator:                 mock.ExpectID(t, "id1"),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
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
						schemauser.NewPhoneUpdatedEvent(context.Background(),
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
				eventstore:                  tt.fields.eventstore(t),
				idGenerator:                 tt.fields.idGenerator,
				checkPermission:             tt.fields.checkPermission,
				newEncryptedCode:            tt.fields.newCode,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				userEncryption:              crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			}
			details, err := c.CreateSchemaUser(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}

			if tt.res.returnCodePhone != "" {
				assert.NotNil(t, tt.args.user.ReturnCodePhone)
				assert.Equal(t, tt.res.returnCodePhone, *tt.args.user.ReturnCodePhone)
			}
			if tt.res.returnCodeEmail != "" {
				assert.NotNil(t, tt.args.user.ReturnCodeEmail)
				assert.Equal(t, tt.res.returnCodeEmail, *tt.args.user.ReturnCodeEmail)
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
			orgID  string
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
			got, err := r.DeleteSchemaUser(tt.args.ctx, tt.args.orgID, tt.args.userID)
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

func TestCommandSide_LockSchemaUser(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-Eu8I2VAfjF", "Errors.IDMissing"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-G4LOrnjY7q", "Errors.User.NotFound"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-G4LOrnjY7q", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user locked, precondition error",
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
							schemauser.NewLockedEvent(context.Background(),
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-G4LOrnjY7q", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "lock user, ok",
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
						schemauser.NewLockedEvent(context.Background(),
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
			name: "lock user, no permission",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.LockSchemaUser(tt.args.ctx, tt.args.orgID, tt.args.userID)
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

func TestCommandSide_UnlockSchemaUser(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-krXtYscQZh", "Errors.IDMissing"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-gpBv46Lh9m", "Errors.User.NotFound"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-gpBv46Lh9m", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user not locked, precondition error",
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-gpBv46Lh9m", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "unlock user, ok",
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
							schemauser.NewLockedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						schemauser.NewUnlockedEvent(context.Background(),
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
			name: "unlock user, no permission",
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
							schemauser.NewLockedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.UnlockSchemaUser(tt.args.ctx, tt.args.orgID, tt.args.userID)
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

func TestCommandSide_DeactivateSchemaUser(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-pjJhge86ZV", "Errors.IDMissing"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-Ob6lR5iFTe", "Errors.User.NotFound"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-Ob6lR5iFTe", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user not active, precondition error",
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
							schemauser.NewDeactivatedEvent(context.Background(),
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-Ob6lR5iFTe", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "deactivate user, ok",
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
						schemauser.NewDeactivatedEvent(context.Background(),
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
			name: "deactivate user, no permission",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.DeactivateSchemaUser(tt.args.ctx, tt.args.orgID, tt.args.userID)
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

func TestCommandSide_ReactivateSchemaUser(t *testing.T) {
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-17XupGvxBJ", "Errors.IDMissing"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-rQjbBr4J3j", "Errors.User.NotFound"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-rQjbBr4J3j", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user not inactive, precondition error",
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-rQjbBr4J3j", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "activate user, ok",
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
							schemauser.NewDeactivatedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						schemauser.NewActivatedEvent(context.Background(),
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
			name: "activate user, no permission",
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
							schemauser.NewDeactivatedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.ActivateSchemaUser(tt.args.ctx, tt.args.orgID, tt.args.userID)
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

func TestCommands_ChangeSchemaUser(t *testing.T) {
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
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		checkPermission             domain.PermissionCheck
		newCode                     encrypedCodeFunc
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		ctx  context.Context
		user *ChangeSchemaUser
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
			"no userID, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-gEJR1QOGHb", "Errors.IDMissing"))
				},
			},
		},
		{
			"schema not existing, error",
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
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "type",
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-VLDTtxT3If", "Errors.UserSchema.NotExists"))
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
				user: &ChangeSchemaUser{
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
			"no valid phone, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID:    "user1",
					Phone: &Phone{Number: "invalid"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "PHONE-so0wa", "Errors.User.Phone.Invalid"))
				},
			},
		},
		{
			"user update, no permission",
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
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						Data: json.RawMessage(`{
						"name": "user"
					}`),
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"user updated, same schema",
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
					expectPush(
						schemauser.NewUpdatedEvent(
							context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							[]schemauser.Changes{
								schemauser.ChangeData(
									json.RawMessage(`{
						"name": "user2"
					}`),
								),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						Data: json.RawMessage(`{
						"name": "user2"
					}`),
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
			"user updated, changed schema",
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
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id2", "instanceID").Aggregate,
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
					expectPush(
						schemauser.NewUpdatedEvent(
							context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							[]schemauser.Changes{
								schemauser.ChangeSchemaID("id2"),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "id2",
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
			"user updated, new schema",
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
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id2", "instanceID").Aggregate,
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
					expectPush(
						schemauser.NewUpdatedEvent(
							context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							[]schemauser.Changes{
								schemauser.ChangeSchemaID("id2"),
								schemauser.ChangeData(
									json.RawMessage(`{
						"name": "user2"
					}`),
								),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "id2",
						Data: json.RawMessage(`{
						"name": "user2"
					}`),
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
			"user updated, same schema revision",
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
						"name1": "user1"
					}`),
							),
						),
					),
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
									"name1": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
						eventFromEventPusher(
							schema.NewUpdatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								[]schema.Changes{
									schema.IncreaseRevision(1),
									schema.ChangeSchema(json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name2": {
										"type": "string"
									}
								}
							}`)),
								},
							),
						),
					),
					expectPush(
						schemauser.NewUpdatedEvent(
							context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							[]schemauser.Changes{
								schemauser.ChangeSchemaRevision(2),
								schemauser.ChangeData(
									json.RawMessage(`{
						"name2": "user2"
					}`),
								),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						Data: json.RawMessage(`{
						"name2": "user2"
					}`),
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
			"user updated, new schema and revision",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schemauser.NewCreatedEvent(
								context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"id1",
								2,
								json.RawMessage(`{
						"name1": "user1"
					}`),
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id2", "instanceID").Aggregate,
								"type",
								json.RawMessage(`{
								"$schema": "urn:zitadel:schema:v1",
								"type": "object",
								"properties": {
									"name2": {
										"type": "string"
									}
								}
							}`),
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectPush(
						schemauser.NewUpdatedEvent(
							context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							[]schemauser.Changes{
								schemauser.ChangeSchemaID("id2"),
								schemauser.ChangeSchemaRevision(1),
								schemauser.ChangeData(
									json.RawMessage(`{
						"name2": "user2"
					}`),
								),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "id2",
						Data: json.RawMessage(`{
						"name2": "user2"
					}`),
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
			"user update, no field permission as admin",
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
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "id1",
						Data: json.RawMessage(`{
						"name": "user"
					}`),
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid"))
				},
			},
		},
		{
			"user update, no field permission as user",
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
					), expectFilter(
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
			},
			args{
				ctx: authz.NewMockContext("instanceID", "org1", "user1"),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "type",
						Data: json.RawMessage(`{
						"name": "user"
					}`),
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid"))
				},
			},
		},
		{
			"user update, invalid data type",
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
			},
			args{
				ctx: authz.NewMockContext("instanceID", "org1", "user1"),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "type",
						Data: json.RawMessage(`{
						"name": 1
					}`),
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid"))
				},
			},
		},
		{
			"user update, additional property",
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
					expectPush(
						schemauser.NewUpdatedEvent(
							context.Background(),
							&schemauser.NewAggregate("user1", "org1").Aggregate,
							[]schemauser.Changes{
								schemauser.ChangeData(
									json.RawMessage(`{
						"name": "user1",
						"additional": "property"
					}`),
								),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "id1",
						Data: json.RawMessage(`{
						"name": "user1",
						"additional": "property"
					}`),
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
			"user update, invalid data attribute name",
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
			},
			args{
				ctx: authz.NewMockContext("instanceID", "org1", "user1"),
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "type",
						Data: json.RawMessage(`{
						"invalid": "user"
					}`),
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid"))
				},
			},
		},
		{
			"user update, email not changed",
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
				user: &ChangeSchemaUser{
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
			"user update, email return",
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
				user: &ChangeSchemaUser{
					ID: "user1",
					SchemaUser: &SchemaUser{
						SchemaID: "id1",
					},
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
				returnCodeEmail: "emailverify",
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
				user: &ChangeSchemaUser{
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
			"user updated, phone no change",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
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
						eventFromEventPusher(
							schemauser.NewPhoneUpdatedEvent(context.Background(),
								&schemauser.NewAggregate("user1", "org1").Aggregate,
								"+41791234567",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
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
			"user updated, phone return",
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
				defaultSecretGenerators:     defaultGenerators,
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
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
				returnCodePhone: "phoneverify",
			},
		},
		{
			"user updated, phone to verify",
			fields{
				eventstore: expectEventstore(
					expectFilter(
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
							false,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
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
					expectFilter(
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
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneverify", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &ChangeSchemaUser{
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
			"user updated, full verified",
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
				user: &ChangeSchemaUser{
					ID:    "user1",
					Email: &Email{Address: "test@example.com", Verified: true},
					Phone: &Phone{Number: "+41791234567", Verified: true},
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
				newEncryptedCode:            tt.fields.newCode,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				userEncryption:              crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			}
			details, err := c.ChangeSchemaUser(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}

			if tt.res.returnCodePhone != "" {
				assert.NotNil(t, tt.args.user.ReturnCodePhone)
				assert.Equal(t, tt.res.returnCodePhone, *tt.args.user.ReturnCodePhone)
			}
			if tt.res.returnCodeEmail != "" {
				assert.NotNil(t, tt.args.user.ReturnCodeEmail)
				assert.Equal(t, tt.res.returnCodeEmail, *tt.args.user.ReturnCodeEmail)
			}
		})
	}
}
