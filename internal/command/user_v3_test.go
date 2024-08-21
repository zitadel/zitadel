package command

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
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
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		resourceowner string
		userSchema    *CreateSchemaUser
	}
	type res struct {
		id      string
		details *domain.ObjectDetails
		err     error
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
				eventstore: expectEventstore(),
			},
			args{
				ctx:        authz.NewMockContext("instanceID", "", ""),
				userSchema: &CreateSchemaUser{},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.ResourceOwnerMissing"),
			},
		},
		{
			"no type, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           authz.NewMockContext("instanceID", "", ""),
				resourceowner: "instanceID",
				userSchema:    &CreateSchemaUser{},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.UserSchema.User.Type.Missing"),
			},
		},
		{
			"no revision, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           authz.NewMockContext("instanceID", "", ""),
				resourceowner: "instanceID",
				userSchema: &CreateSchemaUser{
					SchemaType: "type",
				},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "TODO", "Errors.UserSchema.User.Revision.Missing"),
			},
		},
		{
			"schema not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx:           authz.NewMockContext("instanceID", "", ""),
				resourceowner: "instanceID",
				userSchema: &CreateSchemaUser{
					SchemaType:     "type",
					SchemaRevision: 1,
				},
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "TODO", "TODO"),
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
			},
			args{
				ctx:           authz.NewMockContext("instanceID", "", ""),
				resourceowner: "instanceID",
				userSchema: &CreateSchemaUser{
					SchemaType:     "type",
					SchemaRevision: 1,
				},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "TODO", "TODO"),
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
							&schemauser.NewAggregate("id1", "instanceID").Aggregate,
							"type",
							1,
							"",
							"",
							json.RawMessage(`{
						"name": "user"
					}`),
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx:           authz.NewMockContext("instanceID", "", ""),
				resourceowner: "instanceID",
				userSchema: &CreateSchemaUser{
					SchemaType:     "type",
					SchemaRevision: 1,
					Data: json.RawMessage(`{
						"name": "user"
					}`),
				},
			},
			res{
				id: "id1",
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore(t),
				idGenerator: tt.fields.idGenerator,
			}
			gotID, gotDetails, err := c.CreateSchemaUser(tt.args.ctx, tt.args.resourceowner, tt.args.userSchema)
			assert.Equal(t, tt.res.id, gotID)
			assertObjectDetails(t, tt.res.details, gotDetails)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}
