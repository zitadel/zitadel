package command

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
	}
	type args struct {
		ctx  context.Context
		user *CreateSchemaUser
	}
	type res struct {
		id              string
		returnCodeEmail bool
		returnCodePhone bool
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.ResourceOwnerMissing"))
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.UserSchema.User.Type.Missing"))
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
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO"))
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "TODO", "TODO"))
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
				id: "id1",
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"user created, full",
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
						schemauser.NewEmailChangedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							"test@example.com",
						),
						schemauser.NewEmailCodeAddedEvent(context.Background(),
							&schemauser.NewAggregate("id1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
							false,
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
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							false,
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
					Email: &Email{
						Address:     "test@example.com",
						Verified:    false,
						URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
					},
					Phone: &Phone{
						Number:   "+41791234567",
						Verified: false,
					},
				},
			},
			res{
				id: "id1",
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
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
						schemauser.NewEmailChangedEvent(context.Background(),
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
				id: "id1",
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
				idGenerator:     tt.fields.idGenerator,
				checkPermission: tt.fields.checkPermission,
			}
			err := c.CreateSchemaUser(tt.args.ctx, tt.args.user, GetMockSecretGenerator(t), GetMockSecretGenerator(t))
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, tt.args.user.ID)
				assertObjectDetails(t, tt.res.details, tt.args.user.Details)
			}

			if tt.res.returnCodePhone {
				assert.NotEmpty(t, tt.args.user.ReturnCodePhone)
			}
			if tt.res.returnCodeEmail {
				assert.NotEmpty(t, tt.args.user.ReturnCodeEmail)
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.IDMissing"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
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
						schemauser.NewDeletedEvent(context.Background(),
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
