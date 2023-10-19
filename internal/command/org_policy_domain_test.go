package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommandSide_AddDomainPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx                                    context.Context
		orgID                                  string
		userLoginMustBeDomain                  bool
		validateOrgDomains                     bool
		smtpSenderAddressMatchesInstanceDomain bool
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:                                    context.Background(),
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy, no userLoginMustBeDomain change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								true,
								false,
								false,
							),
						),
					),
					expectPush(
						org.NewDomainPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							true,
							true,
							true,
						),
					),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add policy, userLoginMustBeDomain changed, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"test.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainRemovedEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"test.com",
								true,
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org.com",
							),
						),
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"user1@org.com",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user1@org.com",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user2", "org1").Aggregate,
								"user@test.com",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user@test.com",
								false,
							),
						),
					),
					expectPush(
						org.NewDomainPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							true,
							true,
							true,
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1@org.com",
							"user1",
							true,
							user.UsernameChangedEventWithPolicyChange(),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user2", "org1").Aggregate,
							"user@test.com",
							"user@test.com",
							true,
							user.UsernameChangedEventWithPolicyChange(),
						),
					),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
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
			got, err := r.AddOrgDomainPolicy(tt.args.ctx, tt.args.orgID, tt.args.userLoginMustBeDomain, tt.args.validateOrgDomains, tt.args.smtpSenderAddressMatchesInstanceDomain)
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

func TestCommandSide_ChangeDomainPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx                                    context.Context
		orgID                                  string
		userLoginMustBeDomain                  bool
		validateOrgDomains                     bool
		smtpSenderAddressMatchesInstanceDomain bool
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:                                    context.Background(),
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "change, no userLoginMustBeDomain change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								false,
								true,
								true,
							),
						),
					),
					expectPush(
						newDomainPolicyChangedEvent(context.Background(), "org1",
							policy.ChangeValidateOrgDomains(false),
							policy.ChangeSMTPSenderAddressMatchesInstanceDomain(false),
						),
					),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  false,
				validateOrgDomains:                     false,
				smtpSenderAddressMatchesInstanceDomain: false,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change, userLoginMustBeDomain changed, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org.com",
							),
						),
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"user1",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user1@org.com",
								false,
							),
						),
					),
					expectPush(
						newDomainPolicyChangedEvent(context.Background(), "org1",
							policy.ChangeUserLoginMustBeDomain(false),
							policy.ChangeValidateOrgDomains(false),
							policy.ChangeSMTPSenderAddressMatchesInstanceDomain(false),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"user1@org.com",
							false,
							user.UsernameChangedEventWithPolicyChange(),
						),
					),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  false,
				validateOrgDomains:                     false,
				smtpSenderAddressMatchesInstanceDomain: false,
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
			got, err := r.ChangeOrgDomainPolicy(tt.args.ctx, tt.args.orgID, tt.args.userLoginMustBeDomain, tt.args.validateOrgDomains, tt.args.smtpSenderAddressMatchesInstanceDomain)
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

func TestCommandSide_RemoveDomainPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx   context.Context
		orgID string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy not existing, not found error",
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
			name: "remove, no userLoginMustBeDomain change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						org.NewDomainPolicyRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "remove, userLoginMustBeDomain changed, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								false,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org.com",
							),
						),
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"user1",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user1@org.com",
								false,
							),
						),
					),
					expectPush(
						org.NewDomainPolicyRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"user1@org.com",
							false,
							user.UsernameChangedEventWithPolicyChange(),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
			got, err := r.RemoveOrgDomainPolicy(tt.args.ctx, tt.args.orgID)
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

func newDomainPolicyChangedEvent(ctx context.Context, orgID string, changes ...policy.DomainPolicyChanges) *org.DomainPolicyChangedEvent {
	event, _ := org.NewDomainPolicyChangedEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		changes,
	)
	return event
}
