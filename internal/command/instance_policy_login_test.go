package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
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
				ctx: context.Background(),
				policy: &ChangeLoginPolicy{
					AllowRegister:    true,
					AllowExternalIDP: true,
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
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
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
					EnableRegistrationCaptcha:  false,
					EnableLoginCaptcha:         false,
					CaptchaType:                domain.CaptchaTypeDisabled,
					CaptchaSiteKey:             "",
					CaptchaSecretKey:           "",
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
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
							),
						),
					),
					expectPush(
						newDefaultLoginPolicyChangedEvent(context.Background(),
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
							time.Hour*10,
							time.Hour*20,
							time.Hour*30,
							time.Hour*40,
							time.Hour*50,
							true,
							true,
							domain.CaptchaTypeReCaptcha,
							"sitekey",
							"secretkey",
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				policy: &ChangeLoginPolicy{
					AllowRegister:              false,
					AllowUsernamePassword:      false,
					AllowExternalIDP:           false,
					ForceMFA:                   false,
					ForceMFALocalOnly:          false,
					HidePasswordReset:          false,
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
					EnableRegistrationCaptcha:  true,
					EnableLoginCaptcha:         true,
					CaptchaType:                domain.CaptchaTypeReCaptcha,
					CaptchaSiteKey:             "sitekey",
					CaptchaSecretKey:           "secretkey",
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
			got, err := r.ChangeDefaultLoginPolicy(tt.args.ctx, tt.args.policy)
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
				ctx: context.Background(),
				provider: &domain.IDPProvider{
					IDPConfigID: "config1",
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
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
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
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
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
				err: zerrors.IsErrorAlreadyExists,
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
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
						instance.NewIdentityProviderAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"config1"),
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
		ctx      context.Context
		provider *domain.IDPProvider
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
				err: zerrors.IsErrorInvalidArgument,
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
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
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
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
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
							instance.NewLoginPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
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
						instance.NewIdentityProviderRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"config1"),
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
								false,
								false,
								domain.CaptchaTypeDisabled,
								"",
								"",
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
						instance.NewIdentityProviderRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"config1"),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveIDPProviderFromDefaultLoginPolicy(tt.args.ctx, tt.args.provider)
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
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeTOTP,
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
						instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.SecondFactorTypeTOTP),
					),
				),
			},
			args: args{
				ctx:    authz.WithInstanceID(context.Background(), "INSTANCE"),
				factor: domain.SecondFactorTypeTOTP,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "add factor otp email, ok ",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.SecondFactorTypeOTPEmail),
					),
				),
			},
			args: args{
				ctx:    authz.WithInstanceID(context.Background(), "INSTANCE"),
				factor: domain.SecondFactorTypeOTPEmail,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "add factor otp sms, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.SecondFactorTypeOTPSMS),
					),
				),
			},
			args: args{
				ctx:    authz.WithInstanceID(context.Background(), "INSTANCE"),
				factor: domain.SecondFactorTypeOTPSMS,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "add factor totp, add otp sms, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
					expectPush(
						instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.SecondFactorTypeOTPSMS),
					),
				),
			},
			args: args{
				ctx:    authz.WithInstanceID(context.Background(), "INSTANCE"),
				factor: domain.SecondFactorTypeOTPSMS,
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
			got, err := r.AddSecondFactorToDefaultLoginPolicy(tt.args.ctx, tt.args.factor)
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
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeTOTP,
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "factor removed totp, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeTOTP,
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "factor removed otp email, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTPEmail,
							),
						),
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTPEmail,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTPEmail,
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "factor removed otp sms, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTPSMS,
							),
						),
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTPSMS,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTPSMS,
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove factor totp, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
					expectPush(
						instance.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.SecondFactorTypeTOTP),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeTOTP,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "remove factor email, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTPEmail,
							),
						),
					),
					expectPush(
						instance.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.SecondFactorTypeOTPEmail),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTPEmail,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "remove factor sms, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeOTPSMS,
							),
						),
					),
					expectPush(
						instance.NewLoginPolicySecondFactorRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.SecondFactorTypeOTPSMS),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTPSMS,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "factor added totp, removed otp sms, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.SecondFactorTypeTOTP,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.SecondFactorTypeOTPSMS,
			},
			res: res{
				err: zerrors.IsNotFound,
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
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddMultiFactorDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		factor domain.MultiFactorType
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
							instance.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.MultiFactorTypeU2FWithPIN,
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
						instance.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.MultiFactorTypeU2FWithPIN),
					),
				),
			},
			args: args{
				ctx:    authz.WithInstanceID(context.Background(), "INSTANCE"),
				factor: domain.MultiFactorTypeU2FWithPIN,
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
			got, err := r.AddMultiFactorToDefaultLoginPolicy(tt.args.ctx, tt.args.factor)
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

func TestCommandSide_RemoveMultiFactorDefaultLoginPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		factor domain.MultiFactorType
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
				ctx:    context.Background(),
				factor: domain.MultiFactorTypeU2FWithPIN,
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
				ctx:    context.Background(),
				factor: domain.MultiFactorTypeU2FWithPIN,
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
							instance.NewLoginPolicyMultiFactorAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								domain.MultiFactorTypeU2FWithPIN,
							),
						),
					),
					expectPush(
						instance.NewLoginPolicyMultiFactorRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.MultiFactorTypeU2FWithPIN),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				factor: domain.MultiFactorTypeU2FWithPIN,
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
			got, err := r.RemoveMultiFactorFromDefaultLoginPolicy(tt.args.ctx, tt.args.factor)
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

func newDefaultLoginPolicyChangedEvent(ctx context.Context, allowRegister, allowUsernamePassword, allowExternalIDP, forceMFA, forceMFALocalOnly,
	hidePasswordReset, ignoreUnknownUsernames, allowDomainDiscovery, disableLoginWithEmail, disableLoginWithPhone bool,
	passwordlessType domain.PasswordlessType,
	redirectURI string,
	passwordLifetime, externalLoginLifetime, mfaInitSkipLifetime, secondFactorLifetime, multiFactorLifetime time.Duration,
	enableRegistrationCaptcha, enableLoginCaptcha bool, captchaType domain.CaptchaType, captchaSiteKey, captchaSecretKey string) *instance.LoginPolicyChangedEvent {
	event, _ := instance.NewLoginPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.LoginPolicyChanges{
			policy.ChangeAllowRegister(allowRegister),
			policy.ChangeAllowExternalIDP(allowExternalIDP),
			policy.ChangeForceMFA(forceMFA),
			policy.ChangeForceMFALocalOnly(forceMFALocalOnly),
			policy.ChangeAllowUserNamePassword(allowUsernamePassword),
			policy.ChangeHidePasswordReset(hidePasswordReset),
			policy.ChangeIgnoreUnknownUsernames(ignoreUnknownUsernames),
			policy.ChangeAllowDomainDiscovery(allowDomainDiscovery),
			policy.ChangeDisableLoginWithEmail(disableLoginWithEmail),
			policy.ChangeDisableLoginWithPhone(disableLoginWithPhone),
			policy.ChangePasswordlessType(passwordlessType),
			policy.ChangeDefaultRedirectURI(redirectURI),
			policy.ChangePasswordCheckLifetime(passwordLifetime),
			policy.ChangeExternalLoginCheckLifetime(externalLoginLifetime),
			policy.ChangeMFAInitSkipLifetime(mfaInitSkipLifetime),
			policy.ChangeSecondFactorCheckLifetime(secondFactorLifetime),
			policy.ChangeMultiFactorCheckLifetime(multiFactorLifetime),
			policy.ChangeEnableRegistrationCaptcha(enableRegistrationCaptcha),
			policy.ChangeEnableLoginCaptcha(enableLoginCaptcha),
			policy.ChangeCaptchaType(captchaType),
			policy.ChangeCaptchaSiteKey(captchaSiteKey),
			policy.ChangeCaptchaSecretKey(captchaSecretKey),
		},
	)
	return event
}
