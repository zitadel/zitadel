package command

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func filterSchemaUserExisting() expect {
	return expectFilter(
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
	)
}

func filterSchemaExisting() expect {
	return expectFilter(
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
	)
}

func filterUsernameExisting(isOrgSpecifc bool) expect {
	return expectFilter(
		eventFromEventPusher(
			authenticator.NewUsernameCreatedEvent(
				context.Background(),
				&authenticator.NewAggregate("username1", "org1").Aggregate,
				"user1",
				isOrgSpecifc,
				"username",
			),
		),
	)
}

func TestCommands_AddUsername(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		idGenerator     id.Generator
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx  context.Context
		user *AddUsername
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
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				user: &AddUsername{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing"))
				},
			},
		},
		{
			"user not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddUsername{
					UserID: "notexisting",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSchemaExisting(),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
				idGenerator:     mock.ExpectID(t, "username1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddUsername{
					UserID: "user1",
					Username: &Username{
						Username: "user1",
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
			"userschema not existing, error",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddUsername{
					UserID: "user1",
					Username: &Username{
						Username: "user1",
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-VLDTtxT3If", "Errors.UserSchema.NotExists"))
				},
			},
		},
		{
			"username added, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSchemaExisting(),
					expectFilter(),
					expectPush(
						authenticator.NewUsernameCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("username1", "org1").Aggregate,
							"user1",
							false,
							"username",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "username1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddUsername{
					UserID: "user1",
					Username: &Username{
						Username:      "username",
						IsOrgSpecific: false,
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
			"username added, isOrgSpecific, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSchemaExisting(),
					expectFilter(),
					expectPush(
						authenticator.NewUsernameCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("username1", "org1").Aggregate,
							"user1",
							true,
							"username",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "username1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddUsername{
					UserID: "user1",
					Username: &Username{
						Username:      "username",
						IsOrgSpecific: true,
					},
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
				eventstore:      tt.fields.eventstore(t),
				idGenerator:     tt.fields.idGenerator,
				checkPermission: tt.fields.checkPermission,
			}
			details, err := c.AddUsername(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_DeleteUsername(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		userID        string
		id            string
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
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-J6ybG5WZiy", "Errors.IDMissing"))
				},
			},
		},
		{
			"no ID, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    authz.NewMockContext("instanceID", "", ""),
				userID: "user1",
				id:     "",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing"))
				},
			},
		},
		{
			"username not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    authz.NewMockContext("instanceID", "", ""),
				userID: "user1",
				id:     "notexisting",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-uEii8L6Awp", "Errors.User.NotFound"))
				},
			},
		},
		{
			"username already removed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							authenticator.NewUsernameCreatedEvent(
								context.Background(),
								&authenticator.NewAggregate("username1", "org1").Aggregate,
								"user1",
								true,
								"username",
							),
						),
						eventFromEventPusher(
							authenticator.NewUsernameDeletedEvent(
								context.Background(),
								&authenticator.NewAggregate("username1", "org1").Aggregate,
								true,
								"username",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    authz.NewMockContext("instanceID", "", ""),
				userID: "user1",
				id:     "notexisting",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-uEii8L6Awp", "Errors.User.NotFound"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					filterUsernameExisting(false),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:    authz.NewMockContext("instanceID", "", ""),
				userID: "user1",
				id:     "username1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"username removed, ok",
			fields{
				eventstore: expectEventstore(
					filterUsernameExisting(false),
					expectPush(
						authenticator.NewUsernameDeletedEvent(
							context.Background(),
							&authenticator.NewAggregate("username1", "org1").Aggregate,
							false,
							"username",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    authz.NewMockContext("instanceID", "", ""),
				userID: "user1",
				id:     "username1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"username removed, isOrgSpecific, ok",
			fields{
				eventstore: expectEventstore(
					filterUsernameExisting(true),
					expectPush(
						authenticator.NewUsernameDeletedEvent(
							context.Background(),
							&authenticator.NewAggregate("username1", "org1").Aggregate,
							true,
							"username",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    authz.NewMockContext("instanceID", "", ""),
				userID: "user1",
				id:     "username1",
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
			}
			details, err := c.DeleteUsername(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}
		})
	}
}
