package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func filterPublicKeyExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			authenticator.NewPublicKeyCreatedEvent(
				context.Background(),
				&authenticator.NewAggregate("pk1", "org1").Aggregate,
				"user1",
				time.Time{},
				[]byte("something"),
			),
		),
	)
}

func TestCommands_AddPublicKey(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		idGenerator     id.Generator
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx  context.Context
		user *AddPublicKey
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
				user: &AddPublicKey{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-14sGR7lTaj", "Errors.IDMissing"))
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
				user: &AddPublicKey{
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
				idGenerator:     mock.ExpectID(t, "pk1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddPublicKey{
					UserID: "user1",
					PublicKey: &PublicKey{
						PublicKey: []byte("something"),
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
				user: &AddPublicKey{
					UserID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-VLDTtxT3If", "Errors.UserSchema.NotExists"))
				},
			},
		},
		{
			"publickey added, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSchemaExisting(),
					expectFilter(),
					expectPush(
						authenticator.NewPublicKeyCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("pk1", "org1").Aggregate,
							"user1",
							time.Time{},
							[]byte("something"),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "pk1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddPublicKey{
					UserID: "user1",
					PublicKey: &PublicKey{
						PublicKey: []byte("something"),
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
			"publickey added, expirationdate, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSchemaExisting(),
					expectFilter(),
					expectPush(
						authenticator.NewPublicKeyCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("pk1", "org1").Aggregate,
							"user1",
							time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
							[]byte("something"),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "pk1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddPublicKey{
					UserID: "user1",
					PublicKey: &PublicKey{
						PublicKey:      []byte("something"),
						ExpirationDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
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
			details, err := c.AddPublicKey(tt.args.ctx, tt.args.user)
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

func TestCommands_DeletePublicKey(t *testing.T) {
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-hzqeAXW1qP", "Errors.IDMissing"))
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-BNNYJz6Yxt", "Errors.IDMissing"))
				},
			},
		},
		{
			"pk not existing, error",
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
				},
			},
		},
		{
			"pk already removed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							authenticator.NewPublicKeyCreatedEvent(
								context.Background(),
								&authenticator.NewAggregate("pk1", "org1").Aggregate,
								"user1",
								time.Time{},
								[]byte("something"),
							),
						),
						eventFromEventPusher(
							authenticator.NewPublicKeyDeletedEvent(
								context.Background(),
								&authenticator.NewAggregate("pk1", "org1").Aggregate,
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					filterPublicKeyExisting(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:    authz.NewMockContext("instanceID", "", ""),
				userID: "user1",
				id:     "pk1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"pk removed, ok",
			fields{
				eventstore: expectEventstore(
					filterPublicKeyExisting(),
					expectPush(
						authenticator.NewPublicKeyDeletedEvent(
							context.Background(),
							&authenticator.NewAggregate("pk1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    authz.NewMockContext("instanceID", "", ""),
				userID: "user1",
				id:     "pk1",
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
			details, err := c.DeletePublicKey(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.id)
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
