package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	duration10 = time.Hour * 10
	duration20 = time.Hour * 20
	duration30 = time.Hour * 30
	duration40 = time.Hour * 40
	duration50 = time.Hour * 50
)

func TestCommandSide_AddLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
		policy *AddLoginPolicy
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
			name: "loginpolicy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								false,
								false,
								domain.PasswordlessTypeAllowed,
								"https://example.com/redirect",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &AddLoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "https://example.com/redirect",
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						org.NewLoginPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							domain.PasswordlessTypeAllowed,
							"https://example.com/redirect",
							time.Hour*1,
							time.Hour*2,
							time.Hour*3,
							time.Hour*4,
							time.Hour*5,
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &AddLoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "https://example.com/redirect",
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add policy with invalid factors, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &AddLoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "https://example.com/redirect",
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
					SecondFactors:              []domain.SecondFactorType{domain.SecondFactorTypeUnspecified},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add policy factors,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						org.NewLoginPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							domain.PasswordlessTypeAllowed,
							"https://example.com/redirect",
							time.Hour*1,
							time.Hour*2,
							time.Hour*3,
							time.Hour*4,
							time.Hour*5,
						),
						org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.SecondFactorTypeTOTP,
						),
						org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.MultiFactorTypeU2FWithPIN,
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &AddLoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "https://example.com/redirect",
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
					SecondFactors:              []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					MultiFactors:               []domain.MultiFactorType{domain.MultiFactorTypeU2FWithPIN},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add policy with unknown idp, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(), // reduce login policy
					expectFilter(), // check if is org idp
					expectFilter(), // check if is instance idp
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &AddLoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "https://example.com/redirect",
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
					IDPProviders: []*AddLoginPolicyIDP{
						{
							Type:     domain.IdentityProviderTypeSystem,
							ConfigID: "invalid",
						},
					},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "add policy instance idp, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewIDPConfigAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								true,
							),
						),
					),
					expectPush(
						org.NewLoginPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							domain.PasswordlessTypeAllowed,
							"https://example.com/redirect",
							time.Hour*1,
							time.Hour*2,
							time.Hour*3,
							time.Hour*4,
							time.Hour*5,
						),
						org.NewIdentityProviderAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"config1",
							domain.IdentityProviderTypeSystem,
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &AddLoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "https://example.com/redirect",
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
					IDPProviders: []*AddLoginPolicyIDP{
						{
							Type:     domain.IdentityProviderTypeSystem,
							ConfigID: "config1",
						},
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add policy org idp, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("ORG").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								true,
							),
						),
					),
					expectPush(
						org.NewLoginPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							true,
							domain.PasswordlessTypeAllowed,
							"https://example.com/redirect",
							time.Hour*1,
							time.Hour*2,
							time.Hour*3,
							time.Hour*4,
							time.Hour*5,
						),
						org.NewIdentityProviderAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"config1",
							domain.IdentityProviderTypeOrg,
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &AddLoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "https://example.com/redirect",
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
					IDPProviders: []*AddLoginPolicyIDP{
						{
							Type:     domain.IdentityProviderTypeOrg,
							ConfigID: "config1",
						},
					},
				},
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
			got, err := r.AddLoginPolicy(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_ChangeLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
		policy *ChangeLoginPolicy
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
			name: "loginpolicy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &ChangeLoginPolicy{
					AllowRegister:          true,
					AllowUsernamePassword:  true,
					AllowExternalIDP:       true,
					ForceMFA:               true,
					ForceMFALocalOnly:      true,
					IgnoreUnknownUsernames: true,
					AllowDomainDiscovery:   true,
					DisableLoginWithEmail:  true,
					DisableLoginWithPhone:  true,
					PasswordlessType:       domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:     "https://example.com/redirect",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"https://example.com/redirect",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &ChangeLoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					DefaultRedirectURI:         "https://example.com/redirect",
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"https://example.com/redirect",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectPush(
						newLoginPolicyChangedEvent(context.Background(),
							"org1",
							false,
							false,
							false,
							false,
							false,
							false,
							false,
							false,
							false,
							false,
							domain.PasswordlessTypeNotAllowed,
							"",
							&duration10,
							&duration20,
							&duration30,
							&duration40,
							&duration50,
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &ChangeLoginPolicy{
					AllowRegister:              false,
					AllowUsernamePassword:      false,
					AllowExternalIDP:           false,
					ForceMFA:                   false,
					ForceMFALocalOnly:          false,
					IgnoreUnknownUsernames:     false,
					AllowDomainDiscovery:       false,
					DisableLoginWithEmail:      false,
					DisableLoginWithPhone:      false,
					PasswordlessType:           domain.PasswordlessTypeNotAllowed,
					DefaultRedirectURI:         "",
					PasswordCheckLifetime:      time.Hour * 10,
					ExternalLoginCheckLifetime: time.Hour * 20,
					MFAInitSkipLifetime:        time.Hour * 30,
					SecondFactorCheckLifetime:  time.Hour * 40,
					MultiFactorCheckLifetime:   time.Hour * 50,
				},
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
			got, err := r.ChangeLoginPolicy(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_RemoveLoginPolicy(t *testing.T) {
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectPush(
						org.NewLoginPolicyRemovedEvent(context.Background(),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveLoginPolicy(tt.args.ctx, tt.args.orgID)
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

func TestCommandSide_AddIDPProviderLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		provider      *domain.IDPProvider
		resourceOwner string
	}
	type res struct {
		want *domain.IDPProvider
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "provider invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      &domain.IDPProvider{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "config not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "provider already exists, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								domain.IdentityProviderTypeOrg,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add provider, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
								true,
							),
						),
					),
					expectFilter(),
					expectPush(
						org.NewIdentityProviderAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"config1",
							domain.IdentityProviderTypeOrg),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.IDPProvider{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					IDPConfigID: "config1",
					Type:        domain.IdentityProviderTypeOrg,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddIDPToLoginPolicy(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
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

func TestCommandSide_RemoveIDPProviderLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      *domain.IDPProvider
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
			name: "resourceowner missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "provider invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      &domain.IDPProvider{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "login policy not exist, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "provider not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "provider removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								domain.IdentityProviderTypeOrg,
							),
						),
						eventFromEventPusher(
							org.NewIdentityProviderRemovedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove provider, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								domain.IdentityProviderTypeOrg,
							),
						),
					),
					expectPush(
						org.NewIdentityProviderRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"config1"),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "remove provider from login policy, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								domain.IdentityProviderTypeOrg,
							),
						),
					),
					expectPush(
						org.NewIdentityProviderRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"config1"),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
					Name:        "name",
					Type:        domain.IdentityProviderTypeOrg,
				},
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
			got, err := r.RemoveIDPFromLoginPolicy(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
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

func TestCommandSide_AddSecondFactorLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		factor        domain.SecondFactorType
		resourceOwner string
	}
	type res struct {
		want domain.SecondFactorType
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeU2F,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeUnspecified,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor already exists, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeTOTP,
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add factor totp, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.SecondFactorTypeTOTP),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeTOTP,
				resourceOwner: "org1",
			},
			res: res{
				want: domain.SecondFactorTypeTOTP,
			},
		},
		{
			name: "add factor otp email, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.SecondFactorTypeOTPEmail),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTPEmail,
				resourceOwner: "org1",
			},
			res: res{
				want: domain.SecondFactorTypeOTPEmail,
			},
		},
		{
			name: "add factor otp sms, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.SecondFactorTypeOTPSMS),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTPSMS,
				resourceOwner: "org1",
			},
			res: res{
				want: domain.SecondFactorTypeOTPSMS,
			},
		},
		{
			name: "add factor totp, add otp sms, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
					expectPush(
						org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.SecondFactorTypeOTPSMS),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTPSMS,
				resourceOwner: "org1",
			},
			res: res{
				want: domain.SecondFactorTypeOTPSMS,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, _, err := r.AddSecondFactorToLoginPolicy(tt.args.ctx, tt.args.factor, tt.args.resourceOwner)
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

func TestCommandSide_RemoveSecondFactoroginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		factor        domain.SecondFactorType
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
			name: "resourceowner missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeTOTP,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeUnspecified,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeTOTP,
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "factor totp removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeTOTP,
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "factor otp email removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeOTPEmail,
							),
						),
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeOTPEmail,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTPEmail,
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "factor otp sms removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeOTPSMS,
							),
						),
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeOTPSMS,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTPSMS,
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "add factor totp, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
					expectPush(
						org.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.SecondFactorTypeTOTP),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeTOTP,
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add factor otp email, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeOTPEmail,
							),
						),
					),
					expectPush(
						org.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.SecondFactorTypeOTPEmail),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTPEmail,
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add factor otp sms, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.SecondFactorTypeOTPSMS,
							),
						),
					),
					expectPush(
						org.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.SecondFactorTypeOTPSMS),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTPSMS,
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
			_, err := r.RemoveSecondFactorFromLoginPolicy(tt.args.ctx, tt.args.factor, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_AddMultiFactorLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		factor        domain.MultiFactorType
		resourceOwner string
	}
	type res struct {
		want domain.MultiFactorType
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.MultiFactorTypeUnspecified,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor already exists, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.MultiFactorTypeU2FWithPIN,
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add factor, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.MultiFactorTypeU2FWithPIN),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.MultiFactorTypeU2FWithPIN,
				resourceOwner: "org1",
			},
			res: res{
				want: domain.MultiFactorTypeU2FWithPIN,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, _, err := r.AddMultiFactorToLoginPolicy(tt.args.ctx, tt.args.factor, tt.args.resourceOwner)
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

func TestCommandSide_RemoveMultiFactorLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		factor        domain.MultiFactorType
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
			name: "resourceowner missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.MultiFactorTypeUnspecified,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.MultiFactorTypeU2FWithPIN,
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "factor removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
						eventFromEventPusher(
							org.NewLoginPolicyMultiFactorRemovedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.MultiFactorTypeU2FWithPIN,
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "add factor, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
					expectPush(
						org.NewLoginPolicyMultiFactorRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.MultiFactorTypeU2FWithPIN),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.MultiFactorTypeU2FWithPIN,
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
			_, err := r.RemoveMultiFactorFromLoginPolicy(tt.args.ctx, tt.args.factor, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newLoginPolicyChangedEvent(ctx context.Context, orgID string,
	usernamePassword, register, externalIDP, mfa, mfaLocalOnly, passwordReset, ignoreUnknownUsernames, allowDomainDiscovery, disableLoginWithEmail, disableLoginWithPhone bool,
	passwordlessType domain.PasswordlessType,
	redirectURI string,
	passwordLifetime, externalLoginLifetime, mfaInitSkipLifetime, secondFactorLifetime, multiFactorLifetime *time.Duration) *org.LoginPolicyChangedEvent {
	changes := []policy.LoginPolicyChanges{
		policy.ChangeAllowUserNamePassword(usernamePassword),
		policy.ChangeAllowRegister(register),
		policy.ChangeAllowExternalIDP(externalIDP),
		policy.ChangeForceMFA(mfa),
		policy.ChangeForceMFALocalOnly(mfaLocalOnly),
		policy.ChangeHidePasswordReset(passwordReset),
		policy.ChangeIgnoreUnknownUsernames(ignoreUnknownUsernames),
		policy.ChangeAllowDomainDiscovery(allowDomainDiscovery),
		policy.ChangePasswordlessType(passwordlessType),
		policy.ChangeDefaultRedirectURI(redirectURI),
		policy.ChangeDisableLoginWithEmail(disableLoginWithEmail),
		policy.ChangeDisableLoginWithPhone(disableLoginWithPhone),
	}
	if passwordLifetime != nil {
		changes = append(changes, policy.ChangePasswordCheckLifetime(*passwordLifetime))
	}
	if externalLoginLifetime != nil {
		changes = append(changes, policy.ChangeExternalLoginCheckLifetime(*externalLoginLifetime))
	}
	if mfaInitSkipLifetime != nil {
		changes = append(changes, policy.ChangeMFAInitSkipLifetime(*mfaInitSkipLifetime))
	}
	if secondFactorLifetime != nil {
		changes = append(changes, policy.ChangeSecondFactorCheckLifetime(*secondFactorLifetime))
	}
	if multiFactorLifetime != nil {
		changes = append(changes, policy.ChangeMultiFactorCheckLifetime(*multiFactorLifetime))
	}
	event, _ := org.NewLoginPolicyChangedEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		changes,
	)
	return event
}
