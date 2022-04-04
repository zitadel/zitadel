package command

import (
	"context"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/caos/zitadel/internal/repository/user"

	"github.com/stretchr/testify/assert"
)

func TestCommandSide_AddDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.LoginPolicy
	}
	type res struct {
		want *domain.LoginPolicy
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
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								false,
								false,
								false,
								domain.PasswordlessTypeAllowed,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LoginPolicy{
					AllowRegister:         true,
					AllowUsernamePassword: true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewLoginPolicyAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									true,
									true,
									true,
									true,
									true,
									domain.PasswordlessTypeAllowed,
									time.Hour*1,
									time.Hour*2,
									time.Hour*3,
									time.Hour*4,
									time.Hour*5,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				policy: &domain.LoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					HidePasswordReset:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
				},
			},
			res: res{
				want: &domain.LoginPolicy{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					HidePasswordReset:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddDefaultLoginPolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ChangeDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.LoginPolicy
	}
	type res struct {
		want *domain.LoginPolicy
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
				ctx: context.Background(),
				policy: &domain.LoginPolicy{
					AllowRegister:    true,
					AllowExternalIDP: true,
				},
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
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
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
				ctx: context.Background(),
				policy: &domain.LoginPolicy{
					AllowRegister:              true,
					AllowUsernamePassword:      true,
					AllowExternalIDP:           true,
					ForceMFA:                   true,
					HidePasswordReset:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					PasswordCheckLifetime:      time.Hour * 1,
					ExternalLoginCheckLifetime: time.Hour * 2,
					MFAInitSkipLifetime:        time.Hour * 3,
					SecondFactorCheckLifetime:  time.Hour * 4,
					MultiFactorCheckLifetime:   time.Hour * 5,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								newDefaultLoginPolicyChangedEvent(context.Background(),
									false,
									false,
									false,
									false,
									false,
									domain.PasswordlessTypeNotAllowed,
									time.Hour*10,
									time.Hour*20,
									time.Hour*30,
									time.Hour*40,
									time.Hour*50),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				policy: &domain.LoginPolicy{
					AllowRegister:              false,
					AllowUsernamePassword:      false,
					AllowExternalIDP:           false,
					ForceMFA:                   false,
					HidePasswordReset:          false,
					PasswordlessType:           domain.PasswordlessTypeNotAllowed,
					PasswordCheckLifetime:      time.Hour * 10,
					ExternalLoginCheckLifetime: time.Hour * 20,
					MFAInitSkipLifetime:        time.Hour * 30,
					SecondFactorCheckLifetime:  time.Hour * 40,
					MultiFactorCheckLifetime:   time.Hour * 50,
				},
			},
			res: res{
				want: &domain.LoginPolicy{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					AllowRegister:              false,
					AllowUsernamePassword:      false,
					AllowExternalIDP:           false,
					ForceMFA:                   false,
					HidePasswordReset:          false,
					PasswordlessType:           domain.PasswordlessTypeNotAllowed,
					PasswordCheckLifetime:      time.Hour * 10,
					ExternalLoginCheckLifetime: time.Hour * 20,
					MFAInitSkipLifetime:        time.Hour * 30,
					SecondFactorCheckLifetime:  time.Hour * 40,
					MultiFactorCheckLifetime:   time.Hour * 50,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultLoginPolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_AddIDPProviderDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		provider *domain.IDPProvider
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
			name: "provider invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:      context.Background(),
				provider: &domain.IDPProvider{},
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
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "config not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
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
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "provider already exists, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
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
							instance.NewIDPConfigAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
							instance.NewIdentityProviderAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"config1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add provider, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewIDPConfigAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewIdentityProviderAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									"config1"),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				want: &domain.IDPProvider{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					IDPConfigID: "config1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddIDPProviderToDefaultLoginPolicy(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_RemoveIDPProviderDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx                 context.Context
		provider            *domain.IDPProvider
		cascadeExternalIDPs []*domain.UserIDPLink
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
			name: "provider invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:      context.Background(),
				provider: &domain.IDPProvider{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "login policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "provider not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
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
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "provider removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
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
							instance.NewIdentityProviderAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"config1",
							),
						),
						eventFromEventPusher(
							instance.NewIdentityProviderRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"config1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "remove provider, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
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
							instance.NewIdentityProviderAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"config1",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewIdentityProviderRemovedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									"config1"),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "remove provider external idp not found, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
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
							instance.NewIdentityProviderAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"config1",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewIdentityProviderRemovedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									"config1"),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
				cascadeExternalIDPs: []*domain.UserIDPLink{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user1",
						},
						IDPConfigID: "config1",
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "remove provider with external idps, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
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
							instance.NewIdentityProviderAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"config1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"config1", "", "externaluser1"),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewIdentityProviderRemovedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									"config1"),
							),
							eventFromEventPusher(
								user.NewUserIDPLinkCascadeRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"config1", "externaluser1")),
						},
						uniqueConstraintsFromEventConstraint(user.NewRemoveUserIDPLinkUniqueConstraint("config1", "externaluser1")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
				cascadeExternalIDPs: []*domain.UserIDPLink{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user1",
						},
						IDPConfigID:    "config1",
						ExternalUserID: "externaluser1",
					},
				},
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveIDPProviderFromDefaultLoginPolicy(tt.args.ctx, tt.args.provider, tt.args.cascadeExternalIDPs...)
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

func TestCommandSide_AddSecondFactorDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		factor domain.SecondFactorType
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor already exists, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTP,
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add factor, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									domain.SecondFactorTypeOTP),
							),
						},
					),
				),
			},
			args: args{
				ctx:    authz.WithInstanceID(context.Background(), "INSTANCE"),
				factor: domain.SecondFactorTypeOTP,
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
				eventstore: tt.fields.eventstore,
			}
			_, got, err := r.AddSecondFactorToDefaultLoginPolicy(tt.args.ctx, tt.args.factor)
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

func TestCommandSide_RemoveSecondFactorDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		factor domain.SecondFactorType
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
				err: caos_errs.IsErrorInvalidArgument,
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
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTP,
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "factor removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTP,
							),
						),
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTP,
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "add factor, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTP,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									domain.SecondFactorTypeOTP),
							),
						},
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTP,
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveSecondFactorFromDefaultLoginPolicy(tt.args.ctx, tt.args.factor)
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

func TestCommandSide_AddMultiFactorDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		factor     domain.MultiFactorType
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
			name: "factor invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				factor:     domain.MultiFactorTypeUnspecified,
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "factor already exists, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				factor:     domain.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add factor, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									domain.MultiFactorTypeU2FWithPIN),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				factor:     domain.MultiFactorTypeU2FWithPIN,
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
				eventstore: tt.fields.eventstore,
			}
			_, got, err := r.AddMultiFactorToDefaultLoginPolicy(tt.args.ctx, tt.args.instanceID, tt.args.factor)
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

func TestCommandSide_RemoveMultiFactorDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		factor     domain.MultiFactorType
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
			name: "factor invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				factor:     domain.MultiFactorTypeUnspecified,
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				factor:     domain.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "factor removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
						eventFromEventPusher(
							instance.NewLoginPolicyMultiFactorRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				factor:     domain.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "add factor, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewLoginPolicyMultiFactorRemovedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									domain.MultiFactorTypeU2FWithPIN),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				factor:     domain.MultiFactorTypeU2FWithPIN,
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveMultiFactorFromDefaultLoginPolicy(tt.args.ctx, tt.args.instanceID, tt.args.factor)
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

func newDefaultLoginPolicyChangedEvent(ctx context.Context, allowRegister, allowUsernamePassword, allowExternalIDP, forceMFA, hidePasswordReset bool,
	passwordlessType domain.PasswordlessType,
	passwordLifetime, externalLoginLifetime, mfaInitSkipLifetime, secondFactorLifetime, multiFactorLifetime time.Duration) *instance.LoginPolicyChangedEvent {
	event, _ := instance.NewLoginPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.LoginPolicyChanges{
			policy.ChangeAllowRegister(allowRegister),
			policy.ChangeAllowExternalIDP(allowExternalIDP),
			policy.ChangeForceMFA(forceMFA),
			policy.ChangeAllowUserNamePassword(allowUsernamePassword),
			policy.ChangeHidePasswordReset(hidePasswordReset),
			policy.ChangePasswordlessType(passwordlessType),
			policy.ChangePasswordCheckLifetime(passwordLifetime),
			policy.ChangeExternalLoginCheckLifetime(externalLoginLifetime),
			policy.ChangeMFAInitSkipLifetime(mfaInitSkipLifetime),
			policy.ChangeSecondFactorCheckLifetime(secondFactorLifetime),
			policy.ChangeMultiFactorCheckLifetime(multiFactorLifetime),
		},
	)
	return event
}
