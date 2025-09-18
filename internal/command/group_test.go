package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateGroup(t *testing.T) {
	t.Parallel()

	pushErr := errors.New("push error")
	filterErr := errors.New("filter error")
	idGeneratorErr := errors.New("id generator error")

	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx         context.Context
		group       *domain.Group
		aggregateID func() (string, error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "invalid group name, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:           " ",
					Description:    "example group",
					OrganizationID: "org1",
				},
			},
			wantErr: zerrors.IsErrorInvalidArgument,
		},
		{
			name: "missing organization id, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:        "example",
					Description: "example group",
				},
			},
			wantErr: zerrors.IsErrorInvalidArgument,
		},
		{
			name: "org not found, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:           "example",
					Description:    "example group",
					OrganizationID: "org1",
				},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "failed to generate group id, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:           "example",
					Description:    "example group",
					OrganizationID: "org1",
				},
				aggregateID: func() (string, error) {
					return "", idGeneratorErr
				},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, idGeneratorErr)
			},
		},
		{
			name: "group already exists, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("1234", "org1").Aggregate,
								"example",
								"example group",
								"org1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "1234",
					},
					Name:           "example",
					Description:    "example group",
					OrganizationID: "org1",
				},
			},
			wantErr: zerrors.IsErrorAlreadyExists,
		},
		{
			name: "failed to get org write model, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(filterErr),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:           "example",
					Description:    "example group",
					OrganizationID: "org1",
				},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-j1mH8l", "Errors.Org.NotFound"))
			},
		},
		{
			name: "failed to get group write model, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilterError(filterErr),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:           "example",
					Description:    "example group",
					OrganizationID: "org1",
				},
				aggregateID: func() (string, error) {
					return "12345", nil
				},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, filterErr)
			},
		},
		{
			name: "failed to push group added event, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(),
					expectPushFailed(
						pushErr,
						group.NewGroupAddedEvent(context.Background(),
							&group.NewAggregate("12345", "").Aggregate,
							"example",
							"example group",
							"org1",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:           "example",
					Description:    "example group",
					OrganizationID: "org1",
				},
				aggregateID: func() (string, error) {
					return "12345", nil
				},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, pushErr)
			},
		},
		{
			name: "group without user provided id, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(),
					expectPush(
						group.NewGroupAddedEvent(context.Background(),
							&group.NewAggregate("12345", "").Aggregate,
							"example",
							"example group",
							"org1",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:           "example",
					Description:    "example group",
					OrganizationID: "org1",
				},
				aggregateID: func() (string, error) {
					return "12345", nil
				}, // mock value generated by the id generator
			},
			want: &domain.ObjectDetails{
				ID: "12345",
			},
		},
		{
			name: "group with user provided id, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(),
					expectPush(
						group.NewGroupAddedEvent(context.Background(),
							&group.NewAggregate("9090", "").Aggregate,
							"example",
							"example group",
							"org1",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "9090",
					},
					Name:           "example",
					Description:    "example group",
					OrganizationID: "org1",
				},
			},
			want: &domain.ObjectDetails{
				ID: "9090",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup mock id generator
			// it's not defined directly in the tests because,
			// when run individually, it affects other tests that don't need this mock.
			if tt.args.aggregateID != nil {
				aggregateID, err := tt.args.aggregateID()
				if err != nil {
					tt.fields.idGenerator = mock.NewIDGeneratorExpectError(t, err)
				} else {
					tt.fields.idGenerator = mock.ExpectID(t, aggregateID)
				}
			}

			c := &Commands{
				eventstore:  tt.fields.eventstore(t),
				idGenerator: tt.fields.idGenerator,
			}

			got, err := c.CreateGroup(tt.args.ctx, tt.args.group)
			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotEmpty(t, got.ID)
				assertObjectDetails(t, tt.want, got)
				return
			}
			require.True(t, tt.wantErr(err))
		})
	}
}

func TestCommands_UpdateGroup(t *testing.T) {
	t.Parallel()

	filterErr := errors.New("filter error")
	pushErr := errors.New("push error")

	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx   context.Context
		group *domain.Group
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "invalid group name, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:           " ",
					Description:    "example group",
					OrganizationID: "org1",
				},
			},
			wantErr: zerrors.IsErrorInvalidArgument,
		},
		{
			name: "group not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "1234",
					},
					Name:           "updated name",
					Description:    "updated description",
					OrganizationID: "org1",
				},
			},
			wantErr: zerrors.IsNotFound,
		},
		{
			name: "failed to get group write model, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(filterErr),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:        "updated group name",
					Description: "updated group description",
				},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, filterErr)
			},
		},
		{
			name: "failed to push group changed event, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("1234", "org1").Aggregate,
								"group1",
								"group1 description",
								"org1",
							),
						),
					),
					expectPushFailed(
						pushErr,
						group.NewGroupChangedEvent(context.Background(),
							&group.NewAggregate("1234", "org1").Aggregate,
							"updated group name",
							"updated group description",
							"org1",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name:        "updated group name",
					Description: "updated group description",
				},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, pushErr)
			},
		},
		{
			name: "update group, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("1234", "org1").Aggregate,
								"group1",
								"group1 description",
								"org1",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							group.NewGroupChangedEvent(context.Background(),
								&group.NewAggregate("1234", "org1").Aggregate,
								"groupXX",
								"updated description",
								"org1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "1234",
					},
					Name:        "groupXX",
					Description: "updated description",
				},
			},
			want: &domain.ObjectDetails{
				ID:            "1234",
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.UpdateGroup(tt.args.ctx, tt.args.group)
			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotEmpty(t, got.ID)
				assertObjectDetails(t, tt.want, got)
				return
			}
			require.True(t, tt.wantErr(err))
		})
	}
}

func TestCommands_DeleteGroup(t *testing.T) {
	t.Parallel()

	filterErr := errors.New("filter error")
	pushErr := errors.New("push error")

	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx     context.Context
		groupID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "missing group ID, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "  ",
			},
			wantErr: zerrors.IsErrorInvalidArgument,
		},
		{
			name: "failed to get group write model, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(filterErr),
				),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "1234",
			},
			wantErr: func(err error) bool {
				return errors.Is(err, filterErr)
			},
		},
		{
			name: "group not found, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "1234",
			},
			want: &domain.ObjectDetails{
				ID: "1234",
			},
		},
		{
			name: "failed to push group delete event, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("1234", "").Aggregate,
								"group1",
								"group1 description",
								"org1",
							),
						),
					),
					expectPushFailed(
						pushErr,
						group.NewGroupRemovedEvent(context.Background(),
							&group.NewAggregate("1234", "").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "1234",
			},
			wantErr: func(err error) bool {
				return errors.Is(err, pushErr)
			},
		},
		{
			name: "delete group, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("1234", "").Aggregate,
								"group1",
								"group1 description",
								"org1",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							group.NewGroupRemovedEvent(context.Background(),
								&group.NewAggregate("1234", "").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "1234",
			},
			want: &domain.ObjectDetails{
				ID: "1234",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.DeleteGroup(tt.args.ctx, tt.args.groupID)
			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotEmpty(t, got.ID)
				assertObjectDetails(t, tt.want, got)
				return
			}
			require.True(t, tt.wantErr(err))
		})
	}
}
