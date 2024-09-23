package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func filterSchemaUserPasswordExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			authenticator.NewPasswordCreatedEvent(
				context.Background(),
				&authenticator.NewAggregate("user1", "org1").Aggregate,
				"user1",
				"encoded",
				false,
			),
		),
	)
}

func filterPasswordComplexityPolicyExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
				&org.NewAggregate("org1").Aggregate,
				1,
				false,
				false,
				false,
				false,
			),
		),
	)
}

func TestCommands_SetSchemaUserPassword(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		userPasswordHasher *crypto.Hasher
		checkPermission    domain.PermissionCheck
	}
	type args struct {
		ctx  context.Context
		user *SetSchemaUserPassword
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
				user: &SetSchemaUserPassword{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing"))
				},
			},
		},
		{
			"no password, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty"))
				},
			},
		},
		{
			"user not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:   "notexisting",
					Password: "password",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:   "user1",
					Password: "password",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"password added, ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					filterSchemaUserExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password",
							false,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:         "user1",
					Password:       "password",
					ChangeRequired: false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password",
							false,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:         "user1",
					Password:       "password",
					ChangeRequired: false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, changeRequired, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password",
							true,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:         "user1",
					Password:       "password",
					ChangeRequired: true,
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
				eventstore:         tt.fields.eventstore(t),
				checkPermission:    tt.fields.checkPermission,
				userPasswordHasher: tt.fields.userPasswordHasher,
			}
			details, err := c.SetSchemaUserPassword(tt.args.ctx, tt.args.user)
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

func TestCommands_DeleteSchemaUserPassword(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
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
			"no ID, error",
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing"))
				},
			},
		},
		{
			"password not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "notexisting",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
				},
			},
		},
		{
			"password already removed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							authenticator.NewPasswordCreatedEvent(
								context.Background(),
								&authenticator.NewAggregate("user1", "org1").Aggregate,
								"id1",
								"hash",
								false,
							),
						),
						eventFromEventPusher(
							authenticator.NewPasswordDeletedEvent(
								context.Background(),
								&authenticator.NewAggregate("user1", "org1").Aggregate,
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
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
			"password removed, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					expectPush(
						authenticator.NewPasswordDeletedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			details, err := c.DeleteSchemaUserPassword(tt.args.ctx, tt.args.resourceOwner, tt.args.id)
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
