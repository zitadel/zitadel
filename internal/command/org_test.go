package command

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/member"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestCommandSide_AddOrg(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		iamDomain    string
		zitadelRoles []authz.RoleMapping
	}
	type args struct {
		ctx           context.Context
		name          string
		userID        string
		resourceOwner string
	}
	type res struct {
		want *domain.Org
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid org, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user removed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilterOrgDomainNotFound(),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
						eventFromEventPusher(
							user.NewUserRemovedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								true,
							),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "org2"),
			},
			args: args{
				ctx:           context.Background(),
				name:          "Org",
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "push failed unique constraint, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilterOrgDomainNotFound(),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilterOrgMemberNotFound(),
					expectPushFailed(caos_errs.ThrowAlreadyExists(nil, "id", "internal"),
						[]*repository.Event{
							eventFromEventPusher(org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"Org")),
							eventFromEventPusher(org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"org.iam-domain")),
							eventFromEventPusher(org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"org.iam-domain")),
							eventFromEventPusher(org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"org.iam-domain")),
							eventFromEventPusher(org.NewMemberAddedEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"user1", domain.RoleOrgOwner)),
						},
						uniqueConstraintsFromEventConstraint(org.NewAddOrgNameUniqueConstraint("Org")),
						uniqueConstraintsFromEventConstraint(org.NewAddOrgDomainUniqueConstraint("org.iam-domain")),
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("org2", "user1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "org2"),
				iamDomain:   "iam-domain",
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "ORG_OWNER",
					},
				},
			},
			args: args{
				ctx:           context.Background(),
				name:          "Org",
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "push failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilterOrgDomainNotFound(),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilterOrgMemberNotFound(),
					expectPushFailed(caos_errs.ThrowInternal(nil, "id", "internal"),
						[]*repository.Event{
							eventFromEventPusher(org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"Org")),
							eventFromEventPusher(org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"org.iam-domain")),
							eventFromEventPusher(org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"org.iam-domain")),
							eventFromEventPusher(org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"org.iam-domain")),
							eventFromEventPusher(org.NewMemberAddedEvent(
								context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"user1", domain.RoleOrgOwner)),
						},
						uniqueConstraintsFromEventConstraint(org.NewAddOrgNameUniqueConstraint("Org")),
						uniqueConstraintsFromEventConstraint(org.NewAddOrgDomainUniqueConstraint("org.iam-domain")),
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("org2", "user1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "org2"),
				iamDomain:   "iam-domain",
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "ORG_OWNER",
					},
				},
			},
			args: args{
				ctx:           context.Background(),
				name:          "Org",
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsInternal,
			},
		},
		{
			name: "add org, no error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilterOrgDomainNotFound(),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilterOrgMemberNotFound(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"Org",
							)),
							eventFromEventPusher(org.NewDomainAddedEvent(context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate, "org.iam-domain",
							)),
							eventFromEventPusher(org.NewDomainVerifiedEvent(context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"org.iam-domain",
							)),
							eventFromEventPusher(org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"org.iam-domain",
							)),
							eventFromEventPusher(org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								"user1",
								domain.RoleOrgOwner,
							)),
						},
						uniqueConstraintsFromEventConstraint(org.NewAddOrgNameUniqueConstraint("Org")),
						uniqueConstraintsFromEventConstraint(org.NewAddOrgDomainUniqueConstraint("org.iam-domain")),
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("org2", "user1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "org2"),
				iamDomain:   "iam-domain",
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "ORG_OWNER",
					},
				},
			},
			args: args{
				ctx:           context.Background(),
				name:          "Org",
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Org{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org2",
						ResourceOwner: "org2",
					},
					Name:          "Org",
					State:         domain.OrgStateActive,
					PrimaryDomain: "org.iam-domain",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:   tt.fields.eventstore,
				idGenerator:  tt.fields.idGenerator,
				iamDomain:    tt.fields.iamDomain,
				zitadelRoles: tt.fields.zitadelRoles,
			}
			got, err := r.AddOrg(tt.args.ctx, tt.args.name, tt.args.userID, tt.args.resourceOwner)
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

func TestCommandSide_DeactivateOrg(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
		iamDomain   string
	}
	type args struct {
		ctx   context.Context
		orgID string
	}
	type res struct {
		want *domain.Org
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "org not found, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "org already inactive, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org"),
						),
						eventFromEventPusher(
							org.NewOrgDeactivatedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "push failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org"),
						),
					),
					expectPushFailed(
						caos_errs.ThrowInternal(nil, "id", "message"),
						[]*repository.Event{
							eventFromEventPusher(org.NewOrgDeactivatedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate)),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsInternal,
			},
		},
		{
			name: "deactivate org",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org"),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewOrgDeactivatedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
							)),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			_, err := r.DeactivateOrg(tt.args.ctx, tt.args.orgID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
func TestCommandSide_ReactivateOrg(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
		iamDomain   string
	}
	type args struct {
		ctx   context.Context
		orgID string
	}
	type res struct {
		want *domain.Org
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "org not found, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "org already active, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org"),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "push failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org"),
						),
						eventFromEventPusher(
							org.NewOrgDeactivatedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
							),
						),
					),
					expectPushFailed(
						caos_errs.ThrowInternal(nil, "id", "message"),
						[]*repository.Event{
							eventFromEventPusher(org.NewOrgReactivatedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
							)),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsInternal,
			},
		},
		{
			name: "reactivate org",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org"),
						),
						eventFromEventPusher(
							org.NewOrgDeactivatedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewOrgReactivatedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate)),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			_, err := r.ReactivateOrg(tt.args.ctx, tt.args.orgID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
