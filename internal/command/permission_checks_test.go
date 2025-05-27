package command

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CheckPermission(t *testing.T) {
	type fields struct {
		eventstore            func(*testing.T) *eventstore.Eventstore
		domainPermissionCheck func(*testing.T) domain.PermissionCheck
	}
	type args struct {
		ctx                        context.Context
		permission                 string
		aggregateType              eventstore.AggregateType
		resourceOwner, aggregateID string
	}
	type want struct {
		err func(error) bool
	}
	ctx := context.Background()
	filterErr := errors.New("filter error")
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "resource owner is given, no query",
			fields: fields{
				domainPermissionCheck: mockDomainPermissionCheck(
					ctx,
					"permission",
					"resourceOwner",
					"aggregateID"),
			},
			args: args{
				ctx:           ctx,
				permission:    "permission",
				resourceOwner: "resourceOwner",
				aggregateID:   "aggregateID",
			},
		},
		{
			name: "resource owner is empty, query for resource owner",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(&repository.Event{
						AggregateID:   "aggregateID",
						ResourceOwner: sql.NullString{String: "resourceOwner"},
					}),
				),
				domainPermissionCheck: mockDomainPermissionCheck(ctx, "permission", "resourceOwner", "aggregateID"),
			},
			args: args{
				ctx:           ctx,
				permission:    "permission",
				resourceOwner: "",
				aggregateID:   "aggregateID",
			},
		},
		{
			name: "resource owner is empty, query for resource owner, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(filterErr),
				),
			},
			args: args{
				ctx:           ctx,
				permission:    "permission",
				resourceOwner: "",
				aggregateID:   "aggregateID",
			},
			want: want{
				err: func(err error) bool {
					return errors.Is(err, filterErr)
				},
			},
		},
		{
			name: "resource owner is empty, query for resource owner, no events",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           ctx,
				permission:    "permission",
				resourceOwner: "",
				aggregateID:   "aggregateID",
			},
			want: want{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no aggregateID, internal error",
			args: args{
				ctx: ctx,
			},
			want: want{
				err: zerrors.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				checkPermission: func(ctx context.Context, permission, orgID, resourceID string) (err error) {
					assert.Failf(t, "Domain permission check should not be called", "Called c.checkPermission(%v,%v,%v,%v)", ctx, permission, orgID, resourceID)
					return nil
				},
				eventstore: expectEventstore()(t),
			}
			if tt.fields.domainPermissionCheck != nil {
				c.checkPermission = tt.fields.domainPermissionCheck(t)
			}
			if tt.fields.eventstore != nil {
				c.eventstore = tt.fields.eventstore(t)
			}
			err := c.newPermissionCheck(tt.args.ctx, tt.args.permission, tt.args.aggregateType)(tt.args.resourceOwner, tt.args.aggregateID)
			if tt.want.err != nil {
				assert.True(t, tt.want.err(err))
			}
		})
	}
}

func TestCommands_CheckPermissionUserWrite(t *testing.T) {
	type fields struct {
		domainPermissionCheck func(*testing.T) domain.PermissionCheck
	}
	type args struct {
		ctx                        context.Context
		resourceOwner, aggregateID string
	}
	type want struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "self, no permission check",
			args: args{
				ctx: authz.SetCtxData(context.Background(), authz.CtxData{
					UserID: "aggregateID",
				}),
				resourceOwner: "resourceOwner",
				aggregateID:   "aggregateID",
			},
		},
		{
			name: "not self, permission check",
			fields: fields{
				domainPermissionCheck: mockDomainPermissionCheck(
					context.Background(),
					"user.write",
					"resourceOwner",
					"foreignAggregateID"),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "resourceOwner",
				aggregateID:   "foreignAggregateID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				checkPermission: func(ctx context.Context, permission, orgID, resourceID string) (err error) {
					assert.Failf(t, "Domain permission check should not be called", "Called c.checkPermission(%v,%v,%v,%v)", ctx, permission, orgID, resourceID)
					return nil
				},
			}
			if tt.fields.domainPermissionCheck != nil {
				c.checkPermission = tt.fields.domainPermissionCheck(t)
			}
			err := c.NewPermissionCheckUserWrite(tt.args.ctx)(tt.args.resourceOwner, tt.args.aggregateID)
			if tt.want.err != nil {
				assert.True(t, tt.want.err(err))
			}
		})
	}
}

func TestCommands_CheckPermissionUserDelete(t *testing.T) {
	type fields struct {
		domainPermissionCheck func(*testing.T) domain.PermissionCheck
	}
	type args struct {
		ctx                        context.Context
		resourceOwner, aggregateID string
	}
	type want struct {
		err func(error) bool
	}
	userCtx := authz.SetCtxData(context.Background(), authz.CtxData{
		UserID: "aggregateID",
	})
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "self, no permission check",
			args: args{
				ctx:           userCtx,
				resourceOwner: "resourceOwner",
				aggregateID:   "aggregateID",
			},
		},
		{
			name: "not self, permission check",
			fields: fields{
				domainPermissionCheck: mockDomainPermissionCheck(
					context.Background(),
					"user.delete",
					"resourceOwner",
					"foreignAggregateID"),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "resourceOwner",
				aggregateID:   "foreignAggregateID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				checkPermission: func(ctx context.Context, permission, orgID, resourceID string) (err error) {
					assert.Failf(t, "Domain permission check should not be called", "Called c.checkPermission(%v,%v,%v,%v)", ctx, permission, orgID, resourceID)
					return nil
				},
			}
			if tt.fields.domainPermissionCheck != nil {
				c.checkPermission = tt.fields.domainPermissionCheck(t)
			}
			err := c.checkPermissionDeleteUser(tt.args.ctx, tt.args.resourceOwner, tt.args.aggregateID)
			if tt.want.err != nil {
				assert.True(t, tt.want.err(err))
			}
		})
	}
}

func mockDomainPermissionCheck(expectCtx context.Context, expectPermission, expectResourceOwner, expectResourceID string) func(t *testing.T) domain.PermissionCheck {
	return func(t *testing.T) domain.PermissionCheck {
		return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
			assert.Equal(t, expectCtx, ctx)
			assert.Equal(t, expectPermission, permission)
			assert.Equal(t, expectResourceOwner, orgID)
			assert.Equal(t, expectResourceID, resourceID)
			return nil
		}
	}
}
