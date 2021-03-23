package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_AddLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore    *eventstore.Eventstore
		tokenVerifier *authz.TokenVerifier
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "loginpolicy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.LoginPolicy{
					AllowRegister:         true,
					AllowUsernamePassword: true,
					AllowExternalIDP:      true,
					ForceMFA:              true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "loginpolicy not allowed, permission denied error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								false,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
				),
				tokenVerifier: GetMockVerifier(t),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.LoginPolicy{
					AllowRegister:         true,
					AllowUsernamePassword: true,
					AllowExternalIDP:      true,
					ForceMFA:              true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
				},
			},
			res: res{
				err: caos_errs.IsPermissionDenied,
			},
		},
		{
			name: "add policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								false,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLoginPolicyAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									true,
									true,
									true,
									true,
									domain.PasswordlessTypeAllowed,
								),
							),
						},
					),
				),
				tokenVerifier: GetMockVerifier(t, domain.FeatureLoginPolicyUsernameLogin),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.LoginPolicy{
					AllowRegister:         true,
					AllowUsernamePassword: true,
					AllowExternalIDP:      true,
					ForceMFA:              true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
				},
			},
			res: res{
				want: &domain.LoginPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					AllowRegister:         true,
					AllowUsernamePassword: true,
					AllowExternalIDP:      true,
					ForceMFA:              true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:    tt.fields.eventstore,
				tokenVerifier: tt.fields.tokenVerifier,
			}
			got, err := r.AddLoginPolicy(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_ChangeLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore    *eventstore.Eventstore
		tokenVerifier *authz.TokenVerifier
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LoginPolicy{
					AllowRegister:         true,
					AllowUsernamePassword: true,
					AllowExternalIDP:      true,
					ForceMFA:              true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
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
				policy: &domain.LoginPolicy{
					AllowRegister:         true,
					AllowUsernamePassword: true,
					AllowExternalIDP:      true,
					ForceMFA:              true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "not allowed, permission denied error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								false,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
				),
				tokenVerifier: GetMockVerifier(t),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.LoginPolicy{
					AllowRegister:         true,
					AllowUsernamePassword: true,
					AllowExternalIDP:      true,
					ForceMFA:              true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
				},
			},
			res: res{
				err: caos_errs.IsPermissionDenied,
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
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								false,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
				),
				tokenVerifier: GetMockVerifier(t, domain.FeatureLoginPolicyUsernameLogin),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.LoginPolicy{
					AllowRegister:         true,
					AllowUsernamePassword: true,
					AllowExternalIDP:      true,
					ForceMFA:              true,
					PasswordlessType:      domain.PasswordlessTypeAllowed,
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
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newLoginPolicyChangedEvent(context.Background(), "org1", false, false, false, false, domain.PasswordlessTypeNotAllowed),
							),
						},
					),
				),
				tokenVerifier: GetMockVerifier(t, domain.FeatureLoginPolicyUsernameLogin),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.LoginPolicy{
					AllowRegister:         false,
					AllowUsernamePassword: false,
					AllowExternalIDP:      false,
					ForceMFA:              false,
					PasswordlessType:      domain.PasswordlessTypeNotAllowed,
				},
			},
			res: res{
				want: &domain.LoginPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					AllowRegister:         false,
					AllowUsernamePassword: false,
					AllowExternalIDP:      false,
					ForceMFA:              false,
					PasswordlessType:      domain.PasswordlessTypeNotAllowed,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:    tt.fields.eventstore,
				tokenVerifier: tt.fields.tokenVerifier,
			}
			got, err := r.ChangeLoginPolicy(tt.args.ctx, tt.args.orgID, tt.args.policy)
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
			name: "remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLoginPolicyRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate),
							),
						},
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
				assert.Equal(t, tt.res.want, got)
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "config not existing, precondition error",
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
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1", "or1").Aggregate,
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
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add provider, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
							),
						),
					),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewIdentityProviderAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"config1",
									domain.IdentityProviderTypeOrg),
							),
						},
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
			got, err := r.AddIDPProviderToLoginPolicy(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
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
		ctx                 context.Context
		resourceOwner       string
		provider            *domain.IDPProvider
		cascadeExternalIDPs []*domain.ExternalIDP
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "provider not existing, not found error",
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
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								domain.IdentityProviderTypeOrg,
							),
						),
						eventFromEventPusher(
							org.NewIdentityProviderRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								domain.IdentityProviderTypeOrg,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewIdentityProviderRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"config1"),
							),
						},
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
			name: "remove provider external idp not found, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								domain.IdentityProviderTypeOrg,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewIdentityProviderRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"config1"),
							),
						},
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
				cascadeExternalIDPs: []*domain.ExternalIDP{
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
					ResourceOwner: "org1",
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
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
						eventFromEventPusher(
							org.NewIdentityProviderAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								domain.IdentityProviderTypeOrg,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanExternalIDPAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"config1", "", "externaluser1"),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewIdentityProviderRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"config1"),
							),
							eventFromEventPusher(
								user.NewHumanExternalIDPCascadeRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"config1", "externaluser1")),
						},
						uniqueConstraintsFromEventConstraint(user.NewRemoveExternalIDPUniqueConstraint("config1", "externaluser1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
				},
				cascadeExternalIDPs: []*domain.ExternalIDP{
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
			got, err := r.RemoveIDPProviderFromLoginPolicy(tt.args.ctx, tt.args.resourceOwner, tt.args.provider, tt.args.cascadeExternalIDPs...)
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
				err: caos_errs.IsErrorInvalidArgument,
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
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								domain.SecondFactorTypeOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTP,
				resourceOwner: "org1",
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
								org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									domain.SecondFactorTypeOTP),
							),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTP,
				resourceOwner: "org1",
			},
			res: res{
				want: domain.SecondFactorTypeOTP,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddSecondFactorToLoginPolicy(tt.args.ctx, tt.args.factor, tt.args.resourceOwner)
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
				factor: domain.SecondFactorTypeOTP,
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTP,
				resourceOwner: "org1",
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
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								domain.SecondFactorTypeOTP,
							),
						),
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								domain.SecondFactorTypeOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTP,
				resourceOwner: "org1",
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
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								domain.SecondFactorTypeOTP,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									domain.SecondFactorTypeOTP),
							),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				factor:        domain.SecondFactorTypeOTP,
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
			err := r.RemoveSecondFactorFromLoginPolicy(tt.args.ctx, tt.args.factor, tt.args.resourceOwner)
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
				err: caos_errs.IsErrorInvalidArgument,
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
							org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									domain.MultiFactorTypeU2FWithPIN),
							),
						},
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
			got, err := r.AddMultiFactorToLoginPolicy(tt.args.ctx, tt.args.factor, tt.args.resourceOwner)
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
				err: caos_errs.IsErrorInvalidArgument,
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
				ctx:           context.Background(),
				factor:        domain.MultiFactorTypeU2FWithPIN,
				resourceOwner: "org1",
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
							org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
						eventFromEventPusher(
							org.NewLoginPolicyMultiFactorRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
							org.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLoginPolicyMultiFactorRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									domain.MultiFactorTypeU2FWithPIN),
							),
						},
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
			err := r.RemoveMultiFactorFromLoginPolicy(tt.args.ctx, tt.args.factor, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newLoginPolicyChangedEvent(ctx context.Context, orgID string, usernamePassword, register, externalIDP, mfa bool, passwordlessType domain.PasswordlessType) *org.LoginPolicyChangedEvent {
	event, _ := org.NewLoginPolicyChangedEvent(ctx,
		&org.NewAggregate(orgID, orgID).Aggregate,
		[]policy.LoginPolicyChanges{
			policy.ChangeAllowUserNamePassword(usernamePassword),
			policy.ChangeAllowRegister(register),
			policy.ChangeAllowExternalIDP(externalIDP),
			policy.ChangeForceMFA(mfa),
			policy.ChangePasswordlessType(passwordlessType),
		},
	)
	return event
}
