package command

import (
	"context"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/static/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/features"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/static"
)

func TestCommandSide_SetOrgFeatures(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		iamDomain  string
		static     static.Storage
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		features      *domain.Features
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
			name: "resourceowner missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					LoginPolicyPasswordReset: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no change, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					LoginPolicyPasswordReset: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
					MetadataUser:             false,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "org does not exist, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					LoginPolicyPasswordReset: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
					MetadataUser:             false,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "set with default policies, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								false,
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
							instance.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"primary",
								"secondary",
								"warn",
								"font",
								"primary-dark",
								"secondary-dark",
								"warn-dark",
								"font-dark",
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.InitCodeMessageType,
								domain.MessageSubject,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.LoginCustomText,
								domain.LoginKeyExternalRegistrationUserOverviewTitle,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewPrivacyPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"toslink",
								"privacylink",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLockoutPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								5,
								false,
							),
						),
					),
					expectFilter(),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
							),
						},
					),
				),
				iamDomain: "iam-domain",
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					LoginPolicyPasswordReset: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
					CustomTextMessage:        false,
					CustomTextLogin:          false,
					PrivacyPolicy:            false,
					MetadataUser:             false,
					LockoutPolicy:            false,
					ActionsAllowed:           domain.ActionsNotAllowed,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set with default policies, custom domains, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								false,
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
							instance.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"primary",
								"secondary",
								"warn",
								"font",
								"primary-dark",
								"secondary-dark",
								"warn-dark",
								"font-dark",
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test1",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test1",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test2",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.InitCodeMessageType,
								domain.MessageSubject,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.LoginCustomText,
								domain.LoginKeyExternalRegistrationUserOverviewTitle,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewPrivacyPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"toslink",
								"privacylink",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLockoutPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								5,
								false,
							),
						),
					),
					expectFilter(),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewDomainRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "test1", true),
							),
							eventFromEventPusher(
								org.NewDomainRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "test2", false),
							),
							eventFromEventPusher(
								newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
							),
						},
						uniqueConstraintsFromEventConstraint(org.NewRemoveOrgDomainUniqueConstraint("test1")),
					),
				),
				iamDomain: "iam-domain",
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					LoginPolicyPasswordReset: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
					MetadataUser:             false,
					PrivacyPolicy:            false,
					LockoutPolicy:            false,
					ActionsAllowed:           domain.ActionsNotAllowed,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set with default policies, custom domains, default not primary, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								false,
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
							instance.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"primary",
								"secondary",
								"warn",
								"font",
								"primary-dark",
								"secondary-dark",
								"warn-dark",
								"font-dark",
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test1",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test1",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test2",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.InitCodeMessageType,
								domain.MessageSubject,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.LoginCustomText,
								domain.LoginKeyExternalRegistrationUserOverviewTitle,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewPrivacyPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"toslink",
								"privacylink",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLockoutPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								5,
								false,
							),
						),
					),
					expectFilter(),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewDomainPrimarySetEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "org1.iam-domain"),
							),
							eventFromEventPusher(
								org.NewDomainRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "test1", true),
							),
							eventFromEventPusher(
								org.NewDomainRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "test2", false),
							),
							eventFromEventPusher(
								newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
							),
						},
						uniqueConstraintsFromEventConstraint(org.NewRemoveOrgDomainUniqueConstraint("test1")),
					),
				),
				iamDomain: "iam-domain",
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					LoginPolicyPasswordReset: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
					MetadataUser:             false,
					PrivacyPolicy:            false,
					LockoutPolicy:            false,
					ActionsAllowed:           domain.ActionsNotAllowed,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set with default policies, custom domains, default not existing, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								false,
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
							instance.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"primary",
								"secondary",
								"warn",
								"font",
								"primary-dark",
								"secondary-dark",
								"warn-dark",
								"font-dark",
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test1",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test1",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"test2",
							),
						),
						eventFromEventPusher(
							org.NewDomainRemovedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain", true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.InitCodeMessageType,
								domain.MessageSubject,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.LoginCustomText,
								domain.LoginKeyExternalRegistrationUserOverviewTitle,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewPrivacyPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"toslink",
								"privacylink",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLockoutPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								5,
								false,
							),
						),
					),
					expectFilter(),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewDomainAddedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "org1.iam-domain"),
							),
							eventFromEventPusher(
								org.NewDomainPrimarySetEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "org1.iam-domain"),
							),
							eventFromEventPusher(
								org.NewDomainRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "test1", true),
							),
							eventFromEventPusher(
								org.NewDomainRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, "test2", false),
							),
							eventFromEventPusher(
								newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
							),
						},
						uniqueConstraintsFromEventConstraint(org.NewRemoveOrgDomainUniqueConstraint("test1")),
					),
				),
				iamDomain: "iam-domain",
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					LoginPolicyPasswordReset: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
					MetadataUser:             false,
					PrivacyPolicy:            false,
					LockoutPolicy:            false,
					ActionsAllowed:           domain.ActionsNotAllowed,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set with custom policies, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					//checkOrgExists
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
					),
					//NewOrgFeaturesWriteModel
					expectFilter(),
					//begin ensureOrgSettingsToFeatures
					//begin setAllowedLoginPolicy
					//orgLoginPolicyWriteModelByID
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
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
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
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
								time.Hour*50,
							),
						),
					),
					//getDefaultLoginPolicy
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
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
					//begin setDefaultAuthFactorsInCustomLoginPolicy
					//orgLoginPolicyAuthFactorsWriteModel
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicySecondFactorAddedEvent(context.Background(), &instance.NewAggregate().Aggregate, domain.SecondFactorTypeU2F),
						),
						eventFromEventPusher(
							instance.NewLoginPolicyMultiFactorAddedEvent(context.Background(), &instance.NewAggregate().Aggregate, domain.MultiFactorTypeU2FWithPIN),
						),
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.SecondFactorTypeOTP),
						),
					),
					//addSecondFactorToLoginPolicy
					expectFilter(),
					//removeSecondFactorFromLoginPolicy
					expectFilter(),
					//addMultiFactorToLoginPolicy
					expectFilter(),
					//end setDefaultAuthFactorsInCustomLoginPolicy
					//orgPasswordComplexityPolicyWriteModelByID
					expectFilter(
						eventFromEventPusher(
							instance.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								7,
								false,
								false,
								false,
								false,
							),
						),
					),
					//begin setDefaultAuthFactorsInCustomLoginPolicy
					//orgLabelPolicyWriteModelByID
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"primary",
								"secondary",
								"warn",
								"font",
								"primary-dark",
								"secondary-dark",
								"warn-dark",
								"font-dark",
								false,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewLabelPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"custom",
								"secondary",
								"warn",
								"font",
								"primary-dark",
								"secondary-dark",
								"warn-dark",
								"font-dark",
								false,
								false,
								false,
							),
						),
					),
					//removeLabelPolicy
					expectFilter(),
					//end setDefaultAuthFactorsInCustomLoginPolicy
					//removeCustomDomains
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.InitCodeMessageType,
								domain.MessageSubject,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.LoginCustomText,
								domain.LoginKeyExternalRegistrationUserOverviewTitle,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPrivacyPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"toslink",
								"privacylink",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								5,
								false,
							),
						),
					),
					expectFilter(),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLoginPolicySecondFactorRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.SecondFactorTypeOTP),
							),
							eventFromEventPusher(
								org.NewLoginPolicySecondFactorAddedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.SecondFactorTypeU2F),
							),
							eventFromEventPusher(
								org.NewLoginPolicyMultiFactorAddedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.MultiFactorTypeU2FWithPIN),
							),
							eventFromEventPusher(
								newLoginPolicyChangedEvent(context.Background(), "org1",
									true, true, true, true, true, domain.PasswordlessTypeAllowed,
									nil,
									nil,
									nil,
									nil,
									nil),
							),
							eventFromEventPusher(
								org.NewPasswordComplexityPolicyRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate),
							),
							eventFromEventPusher(
								org.NewLabelPolicyRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate),
							),
							eventFromEventPusher(
								org.NewCustomTextTemplateRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.InitCodeMessageType, language.English),
							),
							eventFromEventPusher(
								org.NewCustomTextTemplateRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, language.English),
							),
							eventFromEventPusher(
								org.NewPrivacyPolicyRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate),
							),
							eventFromEventPusher(
								org.NewLockoutPolicyRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate),
							),
							eventFromEventPusher(
								newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
							),
						},
					),
				),
				iamDomain: "iam-domain",
				static:    mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectsNoError(),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
					MetadataUser:             false,
					PrivacyPolicy:            false,
					LockoutPolicy:            false,
					ActionsAllowed:           domain.ActionsNotAllowed,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set with default policies, usermetadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
					),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								false,
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
							instance.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"primary",
								"secondary",
								"warn",
								"font",
								"primary-dark",
								"secondary-dark",
								"warn-dark",
								"font-dark",
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.InitCodeMessageType,
								domain.MessageSubject,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.LoginCustomText,
								domain.LoginKeyExternalRegistrationUserOverviewTitle,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewPrivacyPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"toslink",
								"privacylink",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLockoutPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								5,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewMetadataRemovedAllEvent(context.Background(), &user.NewAggregate("user1", "org1").Aggregate),
							),
							eventFromEventPusher(
								newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
							),
						},
					),
				),
				iamDomain: "iam-domain",
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					LoginPolicyPasswordReset: false,
					PasswordComplexityPolicy: false,
					LabelPolicyPrivateLabel:  false,
					LabelPolicyWatermark:     false,
					CustomDomain:             false,
					CustomTextMessage:        false,
					CustomTextLogin:          false,
					PrivacyPolicy:            false,
					MetadataUser:             false,
					LockoutPolicy:            false,
					ActionsAllowed:           domain.ActionsNotAllowed,
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
				iamDomain:  tt.fields.iamDomain,
				static:     tt.fields.static,
			}
			got, err := r.SetOrgFeatures(tt.args.ctx, tt.args.resourceOwner, tt.args.features)
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

func TestCommandSide_RemoveOrgFeatures(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		iamDomain  string
	}
	type args struct {
		ctx           context.Context
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
			name: "resourceowner missing, error",
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
			name: "no features set, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "remove with default policies, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour)),
					),
					expectFilter(
						eventFromEventPusher(
							newIAMFeaturesSetEvent(context.Background(), "Default", domain.FeaturesStateActive, time.Hour)),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLoginPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								false,
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
							instance.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"primary",
								"secondary",
								"warn",
								"font",
								"primary-dark",
								"secondary-dark",
								"warn-dark",
								"font-dark",
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(
								context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"org1.iam-domain",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.InitCodeMessageType,
								domain.MessageSubject,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewCustomTextSetEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								domain.LoginCustomText,
								domain.LoginKeyExternalRegistrationUserOverviewTitle,
								"text",
								language.English,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewPrivacyPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								"toslink",
								"privacylink",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewLockoutPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								5,
								false,
							),
						),
					),
					expectFilter(),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewFeaturesRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate),
							),
						},
					),
				),
				iamDomain: "iam-domain",
			},
			args: args{
				ctx:           context.Background(),
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
				iamDomain:  tt.fields.iamDomain,
			}
			got, err := r.RemoveOrgFeatures(tt.args.ctx, tt.args.resourceOwner)
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

func newIAMFeaturesSetEvent(ctx context.Context, tierName string, state domain.FeaturesState, auditLog time.Duration) *instance.FeaturesSetEvent {
	event, _ := instance.NewFeaturesSetEvent(
		ctx,
		&instance.NewAggregate().Aggregate,
		[]features.FeaturesChanges{
			features.ChangeTierName(tierName),
			features.ChangeState(state),
			features.ChangeAuditLogRetention(auditLog),
		},
	)
	return event
}

func newFeaturesSetEvent(ctx context.Context, orgID string, tierName string, state domain.FeaturesState, auditLog time.Duration) *org.FeaturesSetEvent {
	event, _ := org.NewFeaturesSetEvent(
		ctx,
		&org.NewAggregate(orgID, orgID).Aggregate,
		[]features.FeaturesChanges{
			features.ChangeTierName(tierName),
			features.ChangeState(state),
			features.ChangeAuditLogRetention(auditLog),
		},
	)
	return event
}
