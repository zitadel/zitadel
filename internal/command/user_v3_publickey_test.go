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
				time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
				[]byte("something"),
			),
		),
	)
}

var publicKeyExample = []byte("-----BEGIN PUBLIC KEY-----\nMIICITANBgkqhkiG9w0BAQEFAAOCAg4AMIICCQKCAgB5tWxwCGRloCqvpgI2ZXPl\nxQ+WZbQPuHTqAxwbXbsKOJoAAq16iHmzriLKpqVDxRUXqTH3cY0P0A1IZbCBB2gG\nyq3Lk08sR5ute+MEQ+QibX2qpk+mccRr+eP6B1otcyBWxRhZ/YtWphDpZ4GCb4oN\nAzTIebU0ztlu1OOnDDSEEhwScu2LhG40bx4hVU8XNgIqEjxiR61J89vfZpCmn0Rl\nsqYvmX9sqtqPokdsKl3LPItRyDAJMG0uhwwGKsHffDNeLDZN1OCZE/ZS7USarJQH\nbtGeqFQKsCL33xsKbNL+QjnAhqHW09bMdwofJvlwYLfL0rGJQr5aVCaERAfKAOE6\npy0nVkEJsRLxvdx/ZbTtZdCBk/LiznkE1xp9J02obQ+kWHtdUYxM1OSJqPRGQpbS\nZTxurdBQ43gRjO07iWNV9CB0i6QN2GtDBmHVb48i6aPdA++uJqnPYzy46FWA3KMA\nSlxiZ1RDcGH+fN9uklC2cwAurctAxed3Me2RYGdxl813udeV4Ef3qaiV2dix/pKA\nvN1KIfPTpTdULCDBLjtaAYflJ2WYXHeWMJMMC4oJc3bcKpA4mWjZibZ3pSGX/STQ\nXwHUtKsGlrVBSeqjjILVpH+2G0rusrqkGOlPKN+qOIsnwJf9x47v+xEw1slqdDWm\n+x3gc+8m9oowCcq20OeNTQIDAQAB\n-----END PUBLIC KEY-----")

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
					PublicKey: &PublicKey{
						PublicKey: publicKeyExample,
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound"))
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
					PublicKey: &PublicKey{
						PublicKey: publicKeyExample,
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
						PublicKey: publicKeyExample,
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
			"publickey added, no public key format",
			fields{
				eventstore: expectEventstore(),
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-WdWlhUSVqK", "Errors.User.Machine.Key.Invalid"))
				},
			},
		},
		{
			"publickey added, expirationDate before now",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &AddPublicKey{
					UserID: "user1",
					PublicKey: &PublicKey{
						PublicKey:      publicKeyExample,
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
							time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
							publicKeyExample,
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
						PublicKey: publicKeyExample,
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
							time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
							publicKeyExample,
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
						PublicKey:      publicKeyExample,
						ExpirationDate: time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-CqNteIqtCt", "Errors.User.NotFound"))
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
								time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
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
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-CqNteIqtCt", "Errors.User.NotFound"))
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
