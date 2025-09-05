package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddDefaultDomainPolicy(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx                                    context.Context
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
			name: "domain policy already existing, already exists error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy,ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewDomainPolicyAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							true,
							true,
							true,
						),
					),
				),
			},
			args: args{
				ctx:                                    authz.WithInstanceID(context.Background(), "INSTANCE"),
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.AddDefaultDomainPolicy(tt.args.ctx, tt.args.userLoginMustBeDomain, tt.args.validateOrgDomains, tt.args.smtpSenderAddressMatchesInstanceDomain)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeDefaultDomainPolicy(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx                                    context.Context
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
			name: "domain policy not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					// domainPolicyOrgs
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								"org2",
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								false,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3",
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								false,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyRemovedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
							),
						),
					),
					// domainPolicyUsernames for each org
					// org1
					expectFilterOrganizationSettings("org1", false, false),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1.com",
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
								"user1@org1.com",
								false,
							),
						),
					),
					// org3
					expectFilterOrganizationSettings("org3", false, false),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3.com",
							),
						),
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org3").Aggregate,
								"user1",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user1@org3.com",
								false,
							),
						),
					),
					expectPush(
						newDefaultDomainPolicyChangedEvent(context.Background(), false, false, false),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"user1@org1.com",
							false,
							false,
							user.UsernameChangedEventWithPolicyChange(true),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org3").Aggregate,
							"user1",
							"user1@org3.com",
							false,
							false,
							user.UsernameChangedEventWithPolicyChange(true),
						),
					),
				),
			},
			args: args{
				ctx:                                    authz.WithInstanceID(context.Background(), "INSTANCE"),
				userLoginMustBeDomain:                  false,
				validateOrgDomains:                     false,
				smtpSenderAddressMatchesInstanceDomain: false,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "change, ok, org removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					// domainPolicyOrgs
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								"org2",
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								false,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3",
							),
						),
						eventFromEventPusher(
							org.NewOrgRemovedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org",
								[]string{},
								false,
								[]string{},
								[]*domain.UserIDPLink{},
								[]string{},
							),
						),
					),
					// domainPolicyUsernames for each org
					// org1
					expectFilterOrganizationSettings("org1", false, false),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1.com",
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
								"user1@org1.com",
								false,
							),
						),
					),
					expectPush(
						newDefaultDomainPolicyChangedEvent(context.Background(), false, false, false),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"user1@org1.com",
							false,
							false,
							user.UsernameChangedEventWithPolicyChange(true),
						),
					),
				),
			},
			args: args{
				ctx:                                    authz.WithInstanceID(context.Background(), "INSTANCE"),
				userLoginMustBeDomain:                  false,
				validateOrgDomains:                     false,
				smtpSenderAddressMatchesInstanceDomain: false,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "change, organization scoped usernames, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					// domainPolicyOrgs
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								"org2",
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								false,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3",
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								false,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyRemovedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
							),
						),
					),
					// domainPolicyUsernames for each org
					// org1
					expectFilterOrganizationSettings("org1", true, true),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1.com",
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
								"user1@org1.com",
								false,
							),
						),
					),
					// org3
					expectFilterOrganizationSettings("org3", true, true),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3.com",
							),
						),
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org3").Aggregate,
								"user1",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user1@org3.com",
								false,
							),
						),
					),
					expectPush(
						newDefaultDomainPolicyChangedEvent(context.Background(), false, false, false),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"user1@org1.com",
							false,
							true,
							user.UsernameChangedEventWithPolicyChange(true),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org3").Aggregate,
							"user1",
							"user1@org3.com",
							false,
							true,
							user.UsernameChangedEventWithPolicyChange(true),
						),
					),
				),
			},
			args: args{
				ctx:                                    authz.WithInstanceID(context.Background(), "INSTANCE"),
				userLoginMustBeDomain:                  false,
				validateOrgDomains:                     false,
				smtpSenderAddressMatchesInstanceDomain: false,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "change, organization scoped usernames, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								false,
								true,
								true,
							),
						),
					),
					// domainPolicyOrgs
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								"org2",
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								true,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3",
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								true,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewDomainPolicyRemovedEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
							),
						),
					),
					// domainPolicyUsernames for each org
					// org1
					expectFilterOrganizationSettings("org1", true, true),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org1.com",
							),
						),
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"user1@org1.com",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user1@org1.com",
								false,
							),
						),
					),
					// org3
					expectFilterOrganizationSettings("org3", true, true),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org3").Aggregate,
								"org3.com",
							),
						),
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org3").Aggregate,
								"user1@org3.com",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user1@org3.com",
								true,
							),
						),
					),
					expectPush(
						newDefaultDomainPolicyChangedEvent(context.Background(), true, false, false),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1@org1.com",
							"user1@org1.com",
							true,
							true,
							user.UsernameChangedEventWithPolicyChange(false),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org3").Aggregate,
							"user1@org3.com",
							"user1@org3.com",
							true,
							true,
							user.UsernameChangedEventWithPolicyChange(true),
						),
					),
				),
			},
			args: args{
				ctx:                                    authz.WithInstanceID(context.Background(), "INSTANCE"),
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     false,
				smtpSenderAddressMatchesInstanceDomain: false,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ChangeDefaultDomainPolicy(tt.args.ctx, tt.args.userLoginMustBeDomain, tt.args.validateOrgDomains, tt.args.smtpSenderAddressMatchesInstanceDomain)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func newDefaultDomainPolicyChangedEvent(ctx context.Context, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain bool) *instance.DomainPolicyChangedEvent {
	event, _ := instance.NewDomainPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.DomainPolicyChanges{
			policy.ChangeUserLoginMustBeDomain(userLoginMustBeDomain),
			policy.ChangeValidateOrgDomains(validateOrgDomains),
			policy.ChangeSMTPSenderAddressMatchesInstanceDomain(smtpSenderAddressMatchesInstanceDomain),
		},
	)
	return event
}
