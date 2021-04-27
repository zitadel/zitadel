package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/usergrant"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandSide_AddProjectGrant(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		projectGrant  *domain.ProjectGrant
		resourceOwner string
	}
	type res struct {
		want *domain.ProjectGrant
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid usergrant, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantedOrgID: "grantedorg1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "granted org not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantedOrgID: "grantedorg1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project roles not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1", "grantedorg1").Aggregate,
								"granted org",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "usergrant for project, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1", "grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"grantedorg1",
								[]string{"key1"},
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(project.NewAddProjectGrantUniqueConstraint("grantedorg1", "project1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "projectgrant1"),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
					State:        domain.ProjectGrantStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			got, err := r.AddProjectGrant(tt.args.ctx, tt.args.projectGrant, tt.args.resourceOwner)
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

func TestCommandSide_ChangeProjectGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx                 context.Context
		projectGrant        *domain.ProjectGrant
		resourceOwner       string
		cascadeUserGrantIDs []string
	}
	type res struct {
		want *domain.ProjectGrant
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid projectgrant, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "projectgrant not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "project not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "granted org not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project roles not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1", "grantedorg1").Aggregate,
								"granted org",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant not changed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1", "grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant only added roles, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1", "grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantChangedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								[]string{"key1", "key2"},
							)),
						},
						nil,
					),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1", "key2"},
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1", "key2"},
					State:        domain.ProjectGrantStateActive,
				},
			},
		},
		{
			name: "projectgrant remove roles, usergrant not found, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1", "key2"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1", "grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"",
							),
						),
					),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantChangedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								[]string{"key1"},
							)),
						},
						nil,
					),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
				resourceOwner:       "org1",
				cascadeUserGrantIDs: []string{"usergrant1"},
			},
			res: res{
				want: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
					State:        domain.ProjectGrantStateActive,
				},
			},
		},
		{
			name: "projectgrant remove roles, usergrant not found, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1", "key2"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1", "grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"projectgrant1",
								[]string{"key1", "key2"}),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantChangedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								[]string{"key1"},
							)),
							eventFromEventPusher(usergrant.NewUserGrantCascadeChangedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								[]string{"key1"},
							)),
						},
						nil,
					),
				),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
				resourceOwner:       "org1",
				cascadeUserGrantIDs: []string{"usergrant1"},
			},
			res: res{
				want: &domain.ProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
					State:        domain.ProjectGrantStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeProjectGrant(tt.args.ctx, tt.args.projectGrant, tt.args.resourceOwner, tt.args.cascadeUserGrantIDs...)
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

func TestCommandSide_DeactivateProjectGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		projectID     string
		grantID       string
		resourceOwner string
	}
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
			name: "missing projectid, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing grantid, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, precondition failed error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "projectgrant already deactivated, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
						eventFromEventPusher(project.NewGrantDeactivateEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant deactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantDeactivateEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
							)),
						},
						nil,
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.DeactivateProjectGrant(tt.args.ctx, tt.args.projectID, tt.args.grantID, tt.args.resourceOwner)
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

func TestCommandSide_ReactivateProjectGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		projectID     string
		grantID       string
		resourceOwner string
	}
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
			name: "missing projectid, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing grantid, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, precondition failed error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "projectgrant not inactive, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant reactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
						eventFromEventPusher(project.NewGrantDeactivateEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantReactivatedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
							)),
						},
						nil,
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ReactivateProjectGrant(tt.args.ctx, tt.args.projectID, tt.args.grantID, tt.args.resourceOwner)
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

func TestCommandSide_RemoveProjectGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx                 context.Context
		projectID           string
		grantID             string
		resourceOwner       string
		cascadeUserGrantIDs []string
	}
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
			name: "missing projectid, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing grantid, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, precondition failed error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "projectgrant remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"grantedorg1",
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectGrantUniqueConstraint("grantedorg1", "project1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove, cascading usergrant not found, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"grantedorg1",
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectGrantUniqueConstraint("grantedorg1", "project1")),
					),
				),
			},
			args: args{
				ctx:                 context.Background(),
				projectID:           "project1",
				grantID:             "projectgrant1",
				resourceOwner:       "org1",
				cascadeUserGrantIDs: []string{"usergrant1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove with cascading usergrants, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(usergrant.NewUserGrantAddedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"projectgrant1",
							[]string{"key1"}))),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewGrantRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"grantedorg1",
							)),
							eventFromEventPusher(usergrant.NewUserGrantCascadeRemovedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"projectgrant1",
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectGrantUniqueConstraint("grantedorg1", "project1")),
						uniqueConstraintsFromEventConstraint(usergrant.NewRemoveUserGrantUniqueConstraint("org1", "user1", "project1", "projectgrant1")),
					),
				),
			},
			args: args{
				ctx:                 context.Background(),
				projectID:           "project1",
				grantID:             "projectgrant1",
				resourceOwner:       "org1",
				cascadeUserGrantIDs: []string{"usergrant1"},
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveProjectGrant(tt.args.ctx, tt.args.projectID, tt.args.grantID, tt.args.resourceOwner, tt.args.cascadeUserGrantIDs...)
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
