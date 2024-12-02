package command

import (
	"context"
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
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func filterPATExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			authenticator.NewPATCreatedEvent(
				context.Background(),
				&authenticator.NewAggregate("pk1", "org1").Aggregate,
				"user1",
				time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
				[]string{"first", "second", "third"},
			),
		),
	)
}

func TestCommands_AddPAT(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		idGenerator     id.Generator
		checkPermission domain.PermissionCheck
		tokenAlg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx  context.Context
		user *AddPAT
	}
	type res struct {
		details *domain.ObjectDetails
		token   string
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
				user: &AddPAT{
					PAT: &PAT{},
				},
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
				user: &AddPAT{
					UserID: "notexisting",
					PAT:    &PAT{},
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
				user: &AddPAT{
					UserID: "user1",
					PAT:    &PAT{},
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
				user: &AddPAT{
					UserID: "user1",
					PAT:    &PAT{},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-VLDTtxT3If", "Errors.UserSchema.NotExists"))
				},
			},
		},
		{
			"pat added, expirationDate before now",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddPAT{
					UserID: "user1",
					PAT: &PAT{
						Scopes:         []string{"first", "second", "third"},
						ExpirationDate: time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC),
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "DOMAIN-dv3t5", "Errors.AuthNKey.ExpireBeforeNow"))
				},
			},
		},
		{
			"pat added, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSchemaExisting(),
					expectFilter(),
					expectPush(
						authenticator.NewPATCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("pk1", "org1").Aggregate,
							"user1",
							time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
							[]string{"first", "second", "third"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "pk1"),
				tokenAlg:        crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddPAT{
					UserID: "user1",
					PAT: &PAT{
						Scopes: []string{"first", "second", "third"},
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				token: "pat_cGsxOnVzZXIx",
			},
		},
		{
			"pat added, expirationdate, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserExisting(),
					filterSchemaExisting(),
					expectFilter(),
					expectPush(
						authenticator.NewPATCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("pk1", "org1").Aggregate,
							"user1",
							time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
							[]string{"first", "second", "third"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     mock.ExpectID(t, "pk1"),
				tokenAlg:        crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddPAT{
					UserID: "user1",
					PAT: &PAT{
						ExpirationDate: time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
						Scopes:         []string{"first", "second", "third"},
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				token: "pat_cGsxOnVzZXIx",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				idGenerator:     tt.fields.idGenerator,
				checkPermission: tt.fields.checkPermission,
				keyAlgorithm:    tt.fields.tokenAlg,
			}
			details, err := c.AddPAT(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
				assert.Equal(t, tt.res.token, tt.args.user.PAT.Token)
			}
		})
	}
}

func TestCommands_DeletePAT(t *testing.T) {
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-ur4kxtxIhW", "Errors.User.NotFound"))
				},
			},
		},
		{
			"pk already removed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							authenticator.NewPATCreatedEvent(
								context.Background(),
								&authenticator.NewAggregate("pk1", "org1").Aggregate,
								"user1",
								time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
								[]string{"first", "second", "third"},
							),
						),
						eventFromEventPusher(
							authenticator.NewPATDeletedEvent(
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-ur4kxtxIhW", "Errors.User.NotFound"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					filterPATExisting(),
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
					filterPATExisting(),
					expectPush(
						authenticator.NewPATDeletedEvent(
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
			details, err := c.DeletePAT(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.id)
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
