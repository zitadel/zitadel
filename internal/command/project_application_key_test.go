package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	permissionmock "github.com/zitadel/zitadel/internal/domain/mock"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddAPIApplicationKey(t *testing.T) {
	type fields struct {
		eventstore          func(*testing.T) *eventstore.Eventstore
		idGenerator         id.Generator
		keySize             int
		permissionCheckMock domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		key           *domain.ApplicationKey
		resourceOwner string
	}
	type res struct {
		want *domain.APIApp
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "no aggregateid, invalid argument error",
			fields: fields{
				eventstore:          expectEventstore(),
				permissionCheckMock: permissionmock.MockPermissionCheckOK(),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ApplicationID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no appid, invalid argument error",
			fields: fields{
				eventstore:          expectEventstore(),
				permissionCheckMock: permissionmock.MockPermissionCheckOK(),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore:          expectEventstore(expectFilter()),
				permissionCheckMock: permissionmock.MockPermissionCheckOK(),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					ApplicationID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "create key not allowed, precondition error 1",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								"secret",
								domain.APIAuthMethodTypeBasic),
						),
					),
				),
				idGenerator:         id_mock.NewIDGeneratorExpectIDs(t, "key1"),
				permissionCheckMock: permissionmock.MockPermissionCheckOK(),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					ApplicationID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "permission check failed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								"secret",
								domain.APIAuthMethodTypeBasic),
						),
					),
				),
				idGenerator:         id_mock.NewIDGeneratorExpectIDs(t, "key1"),
				keySize:             10,
				permissionCheckMock: permissionmock.MockPermissionCheckErr(zerrors.ThrowPermissionDenied(nil, "mock.err", "mock permission check failed")),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					ApplicationID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "create key not allowed, precondition error 2",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								"secret",
								domain.APIAuthMethodTypeBasic),
						),
					),
				),
				idGenerator:         id_mock.NewIDGeneratorExpectIDs(t, "key1"),
				keySize:             10,
				permissionCheckMock: permissionmock.MockPermissionCheckOK(),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					ApplicationID: "app1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				idGenerator:        tt.fields.idGenerator,
				applicationKeySize: tt.fields.keySize,
				checkPermission:    tt.fields.permissionCheckMock,
			}
			got, err := r.AddApplicationKey(tt.args.ctx, tt.args.key, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
