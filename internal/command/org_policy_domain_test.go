package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddDomainPolicy(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:                                    context.Background(),
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy already existing, already exists error",
			fields: fields{
				eventstore: expectEventstore(
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
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy, no userLoginMustBeDomain change, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", false, false),
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
							"user1",
							"user1",
							true,
							false,
							user.UsernameChangedEventWithPolicyChange(false),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user2", "org1").Aggregate,
							"user@test.com",
							"user@test.com",
							true,
							false,
							user.UsernameChangedEventWithPolicyChange(false),
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
			name: "add policy, userLoginMustBeDomain changed, org scoped usernames, ok",
			fields: fields{
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", true, true),
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
							"user1",
							"user1",
							true,
							true,
							user.UsernameChangedEventWithPolicyChange(false),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user2", "org1").Aggregate,
							"user@test.com",
							"user@test.com",
							true,
							true,
							user.UsernameChangedEventWithPolicyChange(false),
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
			name: "add policy, userLoginMustBeDomain removed, org scoped usernames, ok",
			fields: fields{
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", true, true),
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
								"user1",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.English,
								domain.GenderUnspecified,
								"user1@org.com",
								true,
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
								true,
							),
						),
					),
					expectPush(
						org.NewDomainPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							false,
							true,
							true,
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"user1@org.com",
							false,
							true,
							user.UsernameChangedEventWithPolicyChange(true),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user2", "org1").Aggregate,
							"user@test.com",
							"user@test.com@org.com",
							false,
							true,
							user.UsernameChangedEventWithPolicyChange(true),
						),
					),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  false,
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.AddOrgDomainPolicy(tt.args.ctx, tt.args.orgID, tt.args.userLoginMustBeDomain, tt.args.validateOrgDomains, tt.args.smtpSenderAddressMatchesInstanceDomain)
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

func TestCommandSide_ChangeDomainPolicy(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:                                    context.Background(),
				userLoginMustBeDomain:                  true,
				validateOrgDomains:                     true,
				smtpSenderAddressMatchesInstanceDomain: true,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
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
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change, no userLoginMustBeDomain change, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", false, false),
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
							false,
							user.UsernameChangedEventWithPolicyChange(true),
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
			name: "change, userLoginMustBeDomain changed, org scoped usernames, ok",
			fields: fields{
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", true, true),
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
							policy.ChangeUserLoginMustBeDomain(true),
							policy.ChangeValidateOrgDomains(false),
							policy.ChangeSMTPSenderAddressMatchesInstanceDomain(false),
						),
						user.NewUsernameChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"user1",
							true,
							true,
							user.UsernameChangedEventWithPolicyChange(false),
						),
					),
				),
			},
			args: args{
				ctx:                                    context.Background(),
				orgID:                                  "org1",
				userLoginMustBeDomain:                  true,
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
			name: "change, userLoginMustBeDomain removed, org scoped usernames, ok",
			fields: fields{
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", true, true),
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
								true,
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
							true,
							user.UsernameChangedEventWithPolicyChange(true),
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ChangeOrgDomainPolicy(tt.args.ctx, tt.args.orgID, tt.args.userLoginMustBeDomain, tt.args.validateOrgDomains, tt.args.smtpSenderAddressMatchesInstanceDomain)
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

func TestCommandSide_RemoveDomainPolicy(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove, no userLoginMustBeDomain change, ok",
			fields: fields{
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", false, false),
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
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", false, false),
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
							false,
							user.UsernameChangedEventWithPolicyChange(true),
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
		{
			name: "remove, userLoginMustBeDomain removed, org scoped usernames, ok",
			fields: fields{
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", true, true),
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
							true,
							user.UsernameChangedEventWithPolicyChange(true),
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
		{
			name: "remove, userLoginMustBeDomain changed, org scoped usernames, ok",
			fields: fields{
				eventstore: expectEventstore(
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
					expectFilterOrganizationSettings("org1", true, true),
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
								true,
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
							"user1",
							true,
							true,
							user.UsernameChangedEventWithPolicyChange(false),
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.RemoveOrgDomainPolicy(tt.args.ctx, tt.args.orgID)
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

func newDomainPolicyChangedEvent(ctx context.Context, orgID string, changes ...policy.DomainPolicyChanges) *org.DomainPolicyChangedEvent {
	event, _ := org.NewDomainPolicyChangedEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		changes,
	)
	return event
}
