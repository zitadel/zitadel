package projection

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	legacy_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	settings "github.com/zitadel/zitadel/internal/repository/organization_settings"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SettingsRelationalProjectionTable = "zitadel.settings"
)

type settingsRelationalProjection struct{}

func newSettingsRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(settingsRelationalProjection))
}

func (*settingsRelationalProjection) Name() string {
	return SettingsRelationalProjectionTable
}

func (s *settingsRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				// 		// Login
				{
					Event:  org.LoginPolicyAddedEventType,
					Reduce: s.reduceLoginPolicyAdded,
				},
				{
					Event:  org.LoginPolicyChangedEventType,
					Reduce: s.reduceLoginPolicyChanged,
				},
				{
					Event:  org.LoginPolicyRemovedEventType,
					Reduce: s.reduceLoginPolicyRemoved,
				},
				{
					Event:  org.LoginPolicyMultiFactorAddedEventType,
					Reduce: s.reduceMFAAdded,
				},
				{
					Event:  org.LoginPolicyMultiFactorRemovedEventType,
					Reduce: s.reduceMFARemoved,
				},
				{
					Event:  org.LoginPolicySecondFactorAddedEventType,
					Reduce: s.reduceSecondFactorAdded,
				},
				{
					Event:  org.LoginPolicySecondFactorRemovedEventType,
					Reduce: s.reduceSecondFactorRemoved,
				},
				// label
				{
					Event:  org.LabelPolicyAddedEventType,
					Reduce: s.reduceLabelAdded,
				},
				{
					Event:  org.LabelPolicyChangedEventType,
					Reduce: s.reduceLabelChanged,
				},
				{
					Event:  org.LabelPolicyRemovedEventType,
					Reduce: s.reduceLabelPolicyRemoved,
				},
				{
					Event:  org.LabelPolicyActivatedEventType,
					Reduce: s.reduceLabelActivated,
				},
				{
					Event:  org.LabelPolicyLogoAddedEventType,
					Reduce: s.reduceLogoAdded,
				},
				{
					Event:  org.LabelPolicyLogoRemovedEventType,
					Reduce: s.reduceLogoRemoved,
				},
				{
					Event:  org.LabelPolicyLogoDarkAddedEventType,
					Reduce: s.reduceLogoDarkAdded,
				},
				{
					Event:  org.LabelPolicyLogoDarkRemovedEventType,
					Reduce: s.reduceLogoDarkRemoved,
				},
				{
					Event:  org.LabelPolicyIconAddedEventType,
					Reduce: s.reduceIconAdded,
				},
				{
					Event:  org.LabelPolicyIconRemovedEventType,
					Reduce: s.reduceIconRemoved,
				},
				{
					Event:  org.LabelPolicyIconDarkAddedEventType,
					Reduce: s.reduceIconDarkAdded,
				},
				{
					Event:  org.LabelPolicyIconDarkRemovedEventType,
					Reduce: s.reduceIconDarkRemoved,
				},
				{
					Event:  org.LabelPolicyFontAddedEventType,
					Reduce: s.reduceFontAdded,
				},
				{
					Event:  org.LabelPolicyFontRemovedEventType,
					Reduce: s.reduceFontRemoved,
				},
				// 		// Password Complexity
				{
					Event:  org.PasswordComplexityPolicyAddedEventType,
					Reduce: s.reducePasswordComplexityAdded,
				},
				{
					Event:  org.PasswordComplexityPolicyChangedEventType,
					Reduce: s.reducePasswordComplexityChanged,
				},
				{
					Event:  org.PasswordComplexityPolicyRemovedEventType,
					Reduce: s.reducePasswordComplexityRemoved,
				},
				// 		// Password Policy
				{
					Event:  org.PasswordAgePolicyAddedEventType,
					Reduce: s.reducePasswordPolicyAdded,
				},
				{
					Event:  org.PasswordAgePolicyChangedEventType,
					Reduce: s.reducePasswordPolicyChanged,
				},
				{
					Event:  org.PasswordAgePolicyRemovedEventType,
					Reduce: s.reducePasswordPolicyRemoved,
				},
				// 		// Lockout Policy
				{
					Event:  org.LockoutPolicyAddedEventType,
					Reduce: s.reduceLockoutPolicyAdded,
				},
				{
					Event:  org.LockoutPolicyChangedEventType,
					Reduce: s.reduceLockoutPolicyChanged,
				},
				{
					Event:  org.LockoutPolicyRemovedEventType,
					Reduce: s.reduceOrgLockoutPolicyRemoved,
				},
				// 		// Domain Policy
				{
					Event:  org.DomainPolicyAddedEventType,
					Reduce: s.reduceDomainPolicyAdded,
				},
				{
					Event:  org.DomainPolicyChangedEventType,
					Reduce: s.reduceDomainPolicyChanged,
				},
				{
					Event:  org.DomainPolicyRemovedEventType,
					Reduce: s.reduceOrgDomainPolicyRemoved,
				},
				// Notification
				{
					Event:  org.NotificationPolicyAddedEventType,
					Reduce: s.reduceNotificationPolicyAdded,
				},
				{
					Event:  org.NotificationPolicyChangedEventType,
					Reduce: s.reduceNotificationPolicyChanged,
				},
				{
					Event:  org.NotificationPolicyRemovedEventType,
					Reduce: s.reduceOrgNotificationPolicyRemoved,
				},
				// Privacy
				{
					Event:  org.PrivacyPolicyAddedEventType,
					Reduce: s.reducePrivacyPolicyAdded,
				},
				{
					Event:  org.PrivacyPolicyChangedEventType,
					Reduce: s.reducePrivacyPolicyChanged,
				},
				{
					Event:  org.PrivacyPolicyRemovedEventType,
					Reduce: s.reduceOrgPrivacyPolicyRemoved,
				},
			},
		},
		// // settings
		{
			Aggregate: settings.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  settings.OrganizationSettingsSetEventType,
					Reduce: s.reduceOrganizationSettingsSet,
				},
				{
					Event:  settings.OrganizationSettingsRemovedEventType,
					Reduce: s.reduceOrganizationSettingsRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				// 		// Login
				{
					Event:  instance.LoginPolicyAddedEventType,
					Reduce: s.reduceLoginPolicyAdded,
				},
				{
					Event:  instance.LoginPolicyChangedEventType,
					Reduce: s.reduceLoginPolicyChanged,
				},
				{
					Event:  instance.LoginPolicyMultiFactorAddedEventType,
					Reduce: s.reduceMFAAdded,
				},
				{
					Event:  instance.LoginPolicyMultiFactorRemovedEventType,
					Reduce: s.reduceMFARemoved,
				},
				{
					Event:  instance.LoginPolicySecondFactorAddedEventType,
					Reduce: s.reduceSecondFactorAdded,
				},
				{
					Event:  instance.LoginPolicySecondFactorRemovedEventType,
					Reduce: s.reduceSecondFactorRemoved,
				},
				// 		// Label
				{
					Event:  instance.LabelPolicyAddedEventType,
					Reduce: s.reduceLabelAdded,
				},
				{
					Event:  instance.LabelPolicyChangedEventType,
					Reduce: s.reduceLabelChanged,
				},
				{
					Event:  instance.LabelPolicyActivatedEventType,
					Reduce: s.reduceLabelActivated,
				},
				{
					Event:  instance.LabelPolicyLogoAddedEventType,
					Reduce: s.reduceLogoAdded,
				},
				{
					Event:  instance.LabelPolicyLogoRemovedEventType,
					Reduce: s.reduceLogoRemoved,
				},
				{
					Event:  instance.LabelPolicyLogoDarkAddedEventType,
					Reduce: s.reduceLogoDarkAdded,
				},
				{
					Event:  instance.LabelPolicyLogoDarkRemovedEventType,
					Reduce: s.reduceLogoDarkRemoved,
				},
				{
					Event:  instance.LabelPolicyIconAddedEventType,
					Reduce: s.reduceIconAdded,
				},
				{
					Event:  instance.LabelPolicyIconRemovedEventType,
					Reduce: s.reduceIconRemoved,
				},
				{
					Event:  instance.LabelPolicyIconDarkAddedEventType,
					Reduce: s.reduceIconDarkAdded,
				},
				{
					Event:  instance.LabelPolicyIconDarkRemovedEventType,
					Reduce: s.reduceIconDarkRemoved,
				},
				{
					Event:  instance.LabelPolicyFontAddedEventType,
					Reduce: s.reduceFontAdded,
				},
				{
					Event:  instance.LabelPolicyFontRemovedEventType,
					Reduce: s.reduceFontRemoved,
				},
				// 		// Password Complexity
				{
					Event:  instance.PasswordComplexityPolicyAddedEventType,
					Reduce: s.reducePasswordComplexityAdded,
				},
				{
					Event:  instance.PasswordComplexityPolicyChangedEventType,
					Reduce: s.reducePasswordComplexityChanged,
				},
				// 		// Password Policy
				{
					Event:  instance.PasswordAgePolicyAddedEventType,
					Reduce: s.reducePasswordPolicyAdded,
				},
				{
					Event:  instance.PasswordAgePolicyChangedEventType,
					Reduce: s.reducePasswordPolicyChanged,
				},
				// 		// Lockout Policy
				{
					Event:  instance.LockoutPolicyAddedEventType,
					Reduce: s.reduceLockoutPolicyAdded,
				},
				{
					Event:  instance.LockoutPolicyChangedEventType,
					Reduce: s.reduceLockoutPolicyChanged,
				},
				// 		// Domain Policy
				{
					Event:  instance.DomainPolicyAddedEventType,
					Reduce: s.reduceDomainPolicyAdded,
				},
				{
					Event:  instance.DomainPolicyChangedEventType,
					Reduce: s.reduceDomainPolicyChanged,
				},
				// 		// Security Policy
				{
					Event:  instance.SecurityPolicySetEventType,
					Reduce: s.reduceSecurityPolicySet,
				},
				// 	Notification
				{
					Event:  instance.NotificationPolicyAddedEventType,
					Reduce: s.reduceNotificationPolicyAdded,
				},
				{
					Event:  instance.NotificationPolicyChangedEventType,
					Reduce: s.reduceNotificationPolicyChanged,
				},
				// Privacy policy
				{
					Event:  instance.PrivacyPolicyAddedEventType,
					Reduce: s.reducePrivacyPolicyAdded,
				},
				{
					Event:  instance.PrivacyPolicyChangedEventType,
					Reduce: s.reducePrivacyPolicyChanged,
				},
				// Secret Generator
				{
					Event:  instance.SecretGeneratorAddedEventType,
					Reduce: s.reduceSecretGeneratorAdded,
				},
				{
					Event:  instance.SecretGeneratorChangedEventType,
					Reduce: s.reduceSecretGeneratorChanged,
				},
				{
					Event:  instance.SecretGeneratorRemovedEventType,
					Reduce: s.reduceSecretGeneratorRemoved,
				},
			},
		},
	}
}

func (s *settingsRelationalProjection) reduceLoginPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LoginPolicyAddedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
	case *org.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YYPxS", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyAddedEventType, instance.LoginPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.LoginSettings()

		passwordlessType := domain.PasswordlessType(policyEvent.PasswordlessType)
		setting := domain.LoginSetting{
			Setting: domain.Setting{
				InstanceID:     policyEvent.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			AllowUsernamePassword:       policyEvent.AllowUserNamePassword,
			AllowRegister:               policyEvent.AllowRegister,
			AllowExternalIDP:            policyEvent.AllowExternalIDP,
			ForceMultiFactor:            policyEvent.ForceMFA,
			ForceMultiFactorLocalOnly:   policyEvent.ForceMFALocalOnly,
			HidePasswordReset:           policyEvent.HidePasswordReset,
			IgnoreUnknownUsernames:      policyEvent.IgnoreUnknownUsernames,
			AllowDomainDiscovery:        policyEvent.AllowDomainDiscovery,
			DisableLoginWithEmail:       policyEvent.DisableLoginWithEmail,
			DisableLoginWithPhone:       policyEvent.DisableLoginWithPhone,
			PasswordlessType:            passwordlessType,
			DefaultRedirectURI:          policyEvent.DefaultRedirectURI,
			PasswordCheckLifetime:       policyEvent.PasswordCheckLifetime,
			ExternalLoginCheckLifetime:  policyEvent.ExternalLoginCheckLifetime,
			MultiFactorInitSkipLifetime: policyEvent.MFAInitSkipLifetime,
			SecondFactorCheckLifetime:   policyEvent.SecondFactorCheckLifetime,
			MultiFactorCheckLifetime:    policyEvent.MultiFactorCheckLifetime,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &setting)
	}), nil
}

// //nolint:gocognit
func (s *settingsRelationalProjection) reduceLoginPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgID *string
	var policyEvent policy.LoginPolicyChangedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyChangedEvent:
		policyEvent = e.LoginPolicyChangedEvent
	case *org.LoginPolicyChangedEvent:
		policyEvent = e.LoginPolicyChangedEvent
		orgID = &policyEvent.Aggregate().ResourceOwner
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BHd86", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyChangedEventType, instance.LoginPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LoginSettings()

		changes := make([]database.Change, 0, 16)
		if policyEvent.AllowUserNamePassword != nil {
			changes = append(changes, settingsRepo.SetAllowUsernamePassword(*policyEvent.AllowUserNamePassword))
		}
		if policyEvent.AllowRegister != nil {
			changes = append(changes, settingsRepo.SetAllowRegister(*policyEvent.AllowRegister))
		}
		if policyEvent.AllowExternalIDP != nil {
			changes = append(changes, settingsRepo.SetAllowExternalIDP(*policyEvent.AllowExternalIDP))
		}
		if policyEvent.ForceMFA != nil {
			changes = append(changes, settingsRepo.SetForceMultiFactor(*policyEvent.ForceMFA))
		}
		if policyEvent.ForceMFALocalOnly != nil {
			changes = append(changes, settingsRepo.SetForceMultiFactorLocalOnly(*policyEvent.ForceMFALocalOnly))
		}
		if policyEvent.HidePasswordReset != nil {
			changes = append(changes, settingsRepo.SetHidePasswordReset(*policyEvent.HidePasswordReset))
		}
		if policyEvent.IgnoreUnknownUsernames != nil {
			changes = append(changes, settingsRepo.SetIgnoreUnknownUsernames(*policyEvent.IgnoreUnknownUsernames))
		}
		if policyEvent.AllowDomainDiscovery != nil {
			changes = append(changes, settingsRepo.SetAllowDomainDiscovery(*policyEvent.AllowDomainDiscovery))
		}
		if policyEvent.DisableLoginWithEmail != nil {
			changes = append(changes, settingsRepo.SetDisableLoginWithEmail(*policyEvent.DisableLoginWithEmail))
		}
		if policyEvent.DisableLoginWithPhone != nil {
			changes = append(changes, settingsRepo.SetDisableLoginWithPhone(*policyEvent.DisableLoginWithPhone))
		}
		if policyEvent.PasswordlessType != nil {
			changes = append(changes, settingsRepo.SetPasswordlessType(mapPasswordlessType(*policyEvent.PasswordlessType)))
		}
		if policyEvent.DefaultRedirectURI != nil {
			changes = append(changes, settingsRepo.SetDefaultRedirectURI(*policyEvent.DefaultRedirectURI))
		}
		if policyEvent.PasswordCheckLifetime != nil {
			changes = append(changes, settingsRepo.SetPasswordCheckLifetime(*policyEvent.PasswordCheckLifetime))
		}
		if policyEvent.ExternalLoginCheckLifetime != nil {
			changes = append(changes, settingsRepo.SetExternalLoginCheckLifetime(*policyEvent.ExternalLoginCheckLifetime))
		}
		if policyEvent.MFAInitSkipLifetime != nil {
			changes = append(changes, settingsRepo.SetMultiFactorInitSkipLifetime(*policyEvent.MFAInitSkipLifetime))
		}
		if policyEvent.SecondFactorCheckLifetime != nil {
			changes = append(changes, settingsRepo.SetSecondFactorCheckLifetime(*policyEvent.SecondFactorCheckLifetime))
		}
		if policyEvent.MultiFactorCheckLifetime != nil {
			changes = append(changes, settingsRepo.SetMultiFactorCheckLifetime(*policyEvent.MultiFactorCheckLifetime))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgID, changes...)
	}), nil
}

func mapPasswordlessType(passwordlessType legacy_domain.PasswordlessType) domain.PasswordlessType {
	switch passwordlessType {
	case legacy_domain.PasswordlessTypeAllowed:
		return domain.PasswordlessTypeAllowed
	case legacy_domain.PasswordlessTypeNotAllowed:
		fallthrough
	default:
		return domain.PasswordlessTypeNotAllowed
	}
}

func (s *settingsRelationalProjection) reduceMFAAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.MultiFactorAddedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyMultiFactorAddedEvent:
		policyEvent = e.MultiFactorAddedEvent
	case *org.LoginPolicyMultiFactorAddedEvent:
		policyEvent = e.MultiFactorAddedEvent
		orgId = &policyEvent.Aggregate().ID
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-WghuV", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyMultiFactorAddedEventType, instance.LoginPolicyMultiFactorAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLw7y", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.LoginSettings()
		return repo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, repo.AddMultiFactorType(mapMultiFactorType(policyEvent.MFAType)))
	}), nil
}

func mapMultiFactorType(mfaType legacy_domain.MultiFactorType) domain.MultiFactorType {
	switch mfaType {
	case legacy_domain.MultiFactorTypeU2FWithPIN:
		return domain.MultiFactorTypeU2FWithPIN
	case legacy_domain.MultiFactorTypeUnspecified:
		fallthrough
	default:
		return domain.MultiFactorTypeU2FWithPIN
	}
}

func (s *settingsRelationalProjection) reduceMFARemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.MultiFactorRemovedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyMultiFactorRemovedEvent:
		policyEvent = e.MultiFactorRemovedEvent
	case *org.LoginPolicyMultiFactorRemovedEvent:
		policyEvent = e.MultiFactorRemovedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-cHU7u", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyMultiFactorRemovedEventType, instance.LoginPolicyMultiFactorRemovedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLi9y", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.LoginSettings()
		return repo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, repo.RemoveMultiFactorType(mapMultiFactorType(policyEvent.MFAType)))
	}), nil
}

func (*settingsRelationalProjection) reduceLoginPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	loginPolicyRemovedEvent, ok := event.(*org.LoginPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-oRSvD", "reduce.wrong.event.type %s", org.LoginPolicyRemovedEventType)
	}
	return handler.NewStatement(loginPolicyRemovedEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-arg9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LoginSettings()
		_, err := settingsRepo.Delete(
			ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(loginPolicyRemovedEvent.Aggregate().InstanceID, &loginPolicyRemovedEvent.Aggregate().ID, domain.SettingTypeLogin, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceSecondFactorAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.SecondFactorAddedEvent
	switch e := event.(type) {
	case *instance.LoginPolicySecondFactorAddedEvent:
		policyEvent = e.SecondFactorAddedEvent
	case *org.LoginPolicySecondFactorAddedEvent:
		policyEvent = e.SecondFactorAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-apB2E", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicySecondFactorAddedEventType, instance.LoginPolicySecondFactorAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iLk4m", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LoginSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.AddSecondFactorType(mapSecondFactorType(policyEvent.MFAType)))
	}), nil
}

func mapSecondFactorType(mfaType legacy_domain.SecondFactorType) domain.SecondFactorType {
	switch mfaType {
	case legacy_domain.SecondFactorTypeTOTP:
		return domain.SecondFactorTypeTOTP
	case legacy_domain.SecondFactorTypeU2F:
		return domain.SecondFactorTypeU2F
	case legacy_domain.SecondFactorTypeOTPEmail:
		return domain.SecondFactorTypeOTPEmail
	case legacy_domain.SecondFactorTypeOTPSMS:
		return domain.SecondFactorTypeOTPSMS
	case legacy_domain.SecondFactorTypeRecoveryCodes:
		return domain.SecondFactorTypeRecoveryCodes
	case legacy_domain.SecondFactorTypeUnspecified:
		fallthrough
	default:
		return domain.SecondFactorTypeUnspecified
	}
}

func (s *settingsRelationalProjection) reduceSecondFactorRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.SecondFactorRemovedEvent
	switch e := event.(type) {
	case *instance.LoginPolicySecondFactorRemovedEvent:
		policyEvent = e.SecondFactorRemovedEvent
	case *org.LoginPolicySecondFactorRemovedEvent:
		policyEvent = e.SecondFactorRemovedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-bYpmA", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicySecondFactorRemovedEventType, instance.LoginPolicySecondFactorRemovedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rnd0y", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.LoginSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.RemoveSecondFactorType(mapSecondFactorType(policyEvent.MFAType)))
	}), nil
}

// label
func (s *settingsRelationalProjection) reduceLabelAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyAddedEvent
	switch e := event.(type) {
	case *org.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CSE7A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAddedEventType, instance.LabelPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.BrandingSettings()
		themeMode := domain.BrandingPolicyThemeMode(policyEvent.ThemeMode)
		settings := domain.BrandingSetting{
			Setting: domain.Setting{
				InstanceID:     policyEvent.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			PrimaryColorLight:    policyEvent.PrimaryColor,
			BackgroundColorLight: policyEvent.BackgroundColor,
			WarnColorLight:       policyEvent.WarnColor,
			FontColorLight:       policyEvent.FontColor,
			PrimaryColorDark:     policyEvent.PrimaryColorDark,
			BackgroundColorDark:  policyEvent.BackgroundColorDark,
			WarnColorDark:        policyEvent.WarnColorDark,
			FontColorDark:        policyEvent.FontColorDark,
			HideLoginNameSuffix:  policyEvent.HideLoginNameSuffix,
			ErrorMessagePopup:    policyEvent.ErrorMsgPopup,
			DisableWatermark:     policyEvent.DisableWatermark,
			ThemeMode:            themeMode,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (s *settingsRelationalProjection) reduceLabelChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyChangedEvent
	switch e := event.(type) {
	case *org.LabelPolicyChangedEvent:
		policyEvent = e.LabelPolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.LabelPolicyChangedEvent:
		policyEvent = e.LabelPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyChangedEventType, instance.LabelPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-lhb9y", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.BrandingSettings()

		changes := make([]database.Change, 0, 12)
		if policyEvent.PrimaryColor != nil {
			changes = append(changes, settingsRepo.SetPrimaryColorLight(*policyEvent.PrimaryColor))
		}
		if policyEvent.BackgroundColor != nil {
			changes = append(changes, settingsRepo.SetBackgroundColorLight(*policyEvent.BackgroundColor))
		}
		if policyEvent.WarnColor != nil {
			changes = append(changes, settingsRepo.SetWarnColorLight(*policyEvent.WarnColor))
		}
		if policyEvent.FontColor != nil {
			changes = append(changes, settingsRepo.SetFontColorLight(*policyEvent.FontColor))
		}
		if policyEvent.PrimaryColorDark != nil {
			changes = append(changes, settingsRepo.SetPrimaryColorDark(*policyEvent.PrimaryColorDark))
		}
		if policyEvent.BackgroundColorDark != nil {
			changes = append(changes, settingsRepo.SetBackgroundColorDark(*policyEvent.BackgroundColorDark))
		}
		if policyEvent.WarnColorDark != nil {
			changes = append(changes, settingsRepo.SetWarnColorDark(*policyEvent.WarnColorDark))
		}
		if policyEvent.FontColorDark != nil {
			changes = append(changes, settingsRepo.SetFontColorDark(*policyEvent.FontColorDark))
		}
		if policyEvent.HideLoginNameSuffix != nil {
			changes = append(changes, settingsRepo.SetHideLoginNameSuffix(*policyEvent.HideLoginNameSuffix))
		}
		if policyEvent.ErrorMsgPopup != nil {
			changes = append(changes, settingsRepo.SetErrorMessagePopup(*policyEvent.ErrorMsgPopup))
		}
		if policyEvent.DisableWatermark != nil {
			changes = append(changes, settingsRepo.SetDisableWatermark(*policyEvent.DisableWatermark))
		}
		if policyEvent.ThemeMode != nil {
			changes = append(changes, settingsRepo.SetThemeMode(domain.BrandingPolicyThemeMode(*policyEvent.ThemeMode)))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, changes...)
	}), nil
}

func (s *settingsRelationalProjection) reduceLabelPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.LabelPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-ATMBz", "reduce.wrong.event.type %s", org.LabelPolicyRemovedEventType)
	}

	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-r7k0y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettings()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(policyEvent.Aggregate().InstanceID, &policyEvent.Aggregate().ID, domain.SettingTypeBranding, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceLabelActivated(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyActivatedEvent
	switch e := event.(type) {
	case *org.LabelPolicyActivatedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyActivatedEvent
	case *instance.LabelPolicyActivatedEvent:
		policyEvent = e.LabelPolicyActivatedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-7Kd8U", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyActivatedEventType, instance.LabelPolicyActivatedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-r7k0y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettings()
		_, err := settingsRepo.ActivateAt(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(policyEvent.Aggregate().InstanceID, orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
			policyEvent.Creation,
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceLogoAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyLogoAddedEvent
	switch e := event.(type) {
	case *org.LabelPolicyLogoAddedEvent:
		orgId = &e.Aggregate().ResourceOwner
		policyEvent = e.LabelPolicyLogoAddedEvent
	case *instance.LabelPolicyLogoAddedEvent:
		policyEvent = e.LabelPolicyLogoAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoAddedEventType, instance.LabelPolicyLogoAddedEventType})
	}
	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		url, err := url.Parse(policyEvent.StoreKey)
		if err != nil {
			return err
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetLogoURLLight(url))
	}), nil
}

func (s *settingsRelationalProjection) reduceLogoDarkAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyLogoDarkAddedEvent
	switch e := event.(type) {
	case *org.LabelPolicyLogoDarkAddedEvent:
		orgId = &e.Aggregate().ResourceOwner
		policyEvent = e.LabelPolicyLogoDarkAddedEvent
	case *instance.LabelPolicyLogoDarkAddedEvent:
		policyEvent = e.LabelPolicyLogoDarkAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoDarkAddedEventType, instance.LabelPolicyLogoDarkAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		url, err := url.Parse(policyEvent.StoreKey)
		if err != nil {
			return err
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetLogoURLDark(url))
	}), nil
}

func (s *settingsRelationalProjection) reduceLogoRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyLogoRemovedEvent
	switch e := event.(type) {
	case *org.LabelPolicyLogoRemovedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyLogoRemovedEvent
	case *instance.LabelPolicyLogoRemovedEvent:
		policyEvent = e.LabelPolicyLogoRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoRemovedEventType, instance.LabelPolicyLogoRemovedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetLogoURLLight(nil))
	}), nil
}

func (s *settingsRelationalProjection) reduceLogoDarkRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyLogoDarkRemovedEvent
	switch e := event.(type) {
	case *org.LabelPolicyLogoDarkRemovedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyLogoDarkRemovedEvent
	case *instance.LabelPolicyLogoDarkRemovedEvent:
		policyEvent = e.LabelPolicyLogoDarkRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoDarkRemovedEventType, instance.LabelPolicyLogoDarkRemovedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetLogoURLDark(nil))
	}), nil
}

func (s *settingsRelationalProjection) reduceIconAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyIconAddedEvent
	switch e := event.(type) {
	case *org.LabelPolicyIconAddedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyIconAddedEvent
	case *instance.LabelPolicyIconAddedEvent:
		policyEvent = e.LabelPolicyIconAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-e2JFz", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconDarkAddedEventType, instance.LabelPolicyIconDarkAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		url, err := url.Parse(policyEvent.StoreKey)
		if err != nil {
			return err
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetIconURLLight(url))
	}), nil
}

func (s *settingsRelationalProjection) reduceIconDarkAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyIconDarkAddedEvent
	switch e := event.(type) {
	case *org.LabelPolicyIconDarkAddedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyIconDarkAddedEvent
	case *instance.LabelPolicyIconDarkAddedEvent:
		policyEvent = e.LabelPolicyIconDarkAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-e2JFz", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconDarkAddedEventType, instance.LabelPolicyIconDarkAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		url, err := url.Parse(policyEvent.StoreKey)
		if err != nil {
			return err
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetIconURLDark(url))
	}), nil
}

func (s *settingsRelationalProjection) reduceIconRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyIconRemovedEvent
	switch e := event.(type) {
	case *org.LabelPolicyIconRemovedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyIconRemovedEvent
	case *instance.LabelPolicyIconRemovedEvent:
		policyEvent = e.LabelPolicyIconRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gfgbY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconRemovedEventType, instance.LabelPolicyIconRemovedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetIconURLLight(nil))
	}), nil
}

func (s *settingsRelationalProjection) reduceIconDarkRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyIconDarkRemovedEvent
	switch e := event.(type) {
	case *org.LabelPolicyIconDarkRemovedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyIconDarkRemovedEvent
	case *instance.LabelPolicyIconDarkRemovedEvent:
		policyEvent = e.LabelPolicyIconDarkRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gfgbY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconDarkRemovedEventType, instance.LabelPolicyIconDarkRemovedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetIconURLDark(nil))
	}), nil
}

func (s *settingsRelationalProjection) reduceFontAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyFontAddedEvent
	switch e := event.(type) {
	case *org.LabelPolicyFontAddedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyFontAddedEvent
	case *instance.LabelPolicyFontAddedEvent:
		policyEvent = e.LabelPolicyFontAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-65i9W", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontAddedEventType, instance.LabelPolicyFontAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		url, err := url.Parse(policyEvent.StoreKey)
		if err != nil {
			return err
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetFontURL(url))
	}), nil
}

func (s *settingsRelationalProjection) reduceFontRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyFontRemovedEvent
	switch e := event.(type) {
	case *org.LabelPolicyFontRemovedEvent:
		orgId = &event.Aggregate().ID
		policyEvent = e.LabelPolicyFontRemovedEvent
	case *instance.LabelPolicyFontRemovedEvent:
		policyEvent = e.LabelPolicyFontRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-xf32J", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontRemovedEventType, instance.LabelPolicyFontRemovedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettings()
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, settingsRepo.SetFontURL(nil))
	}), nil
}

func (s *settingsRelationalProjection) reducePasswordComplexityAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.PasswordComplexityPolicyAddedEvent
	switch e := event.(type) {
	case *org.PasswordComplexityPolicyAddedEvent:
		policyEvent = e.PasswordComplexityPolicyAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.PasswordComplexityPolicyAddedEvent:
		policyEvent = e.PasswordComplexityPolicyAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-KTHmJ", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordComplexityPolicyAddedEventType, instance.PasswordComplexityPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.PasswordComplexitySettings()
		settings := domain.PasswordComplexitySetting{
			Setting: domain.Setting{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			MinLength:    policyEvent.MinLength,
			HasLowercase: policyEvent.HasLowercase,
			HasUppercase: policyEvent.HasUppercase,
			HasNumber:    policyEvent.HasNumber,
			HasSymbol:    policyEvent.HasSymbol,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (s *settingsRelationalProjection) reducePasswordComplexityChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.PasswordComplexityPolicyChangedEvent
	switch e := event.(type) {
	case *org.PasswordComplexityPolicyChangedEvent:
		policyEvent = e.PasswordComplexityPolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.PasswordComplexityPolicyChangedEvent:
		policyEvent = e.PasswordComplexityPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-cf3Xb", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordComplexityPolicyChangedEventType, instance.PasswordComplexityPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLrfy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.PasswordComplexitySettings()

		changes := make([]database.Change, 0, 5)
		if policyEvent.MinLength != nil {
			changes = append(changes, settingsRepo.SetMinLength(*policyEvent.MinLength))
		}
		if policyEvent.HasLowercase != nil {
			changes = append(changes, settingsRepo.SetHasLowercase(*policyEvent.HasLowercase))
		}
		if policyEvent.HasUppercase != nil {
			changes = append(changes, settingsRepo.SetHasUppercase(*policyEvent.HasUppercase))
		}
		if policyEvent.HasNumber != nil {
			changes = append(changes, settingsRepo.SetHasNumber(*policyEvent.HasNumber))
		}
		if policyEvent.HasSymbol != nil {
			changes = append(changes, settingsRepo.SetHasSymbol(*policyEvent.HasSymbol))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, changes...)
	}), nil
}

func (s *settingsRelationalProjection) reducePasswordComplexityRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PasswordComplexityPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-wttCd", "reduce.wrong.event.type %s", org.PasswordComplexityPolicyRemovedEventType)
	}

	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.PasswordComplexitySettings()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(policyEvent.Aggregate().InstanceID, &policyEvent.Aggregate().ID, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reducePasswordPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.PasswordAgePolicyAddedEvent
	switch e := event.(type) {
	case *org.PasswordAgePolicyAddedEvent:
		policyEvent = e.PasswordAgePolicyAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.PasswordAgePolicyAddedEvent:
		policyEvent = e.PasswordAgePolicyAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CJqF0", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyAddedEventType, instance.PasswordAgePolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.PasswordExpirySettings()
		settings := domain.PasswordExpirySetting{
			Setting: domain.Setting{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			ExpireWarnDays: policyEvent.ExpireWarnDays,
			MaxAgeDays:     policyEvent.MaxAgeDays,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (s *settingsRelationalProjection) reducePasswordPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.PasswordAgePolicyChangedEvent
	switch e := event.(type) {
	case *org.PasswordAgePolicyChangedEvent:
		policyEvent = e.PasswordAgePolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.PasswordAgePolicyChangedEvent:
		policyEvent = e.PasswordAgePolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-i7FZt", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyChangedEventType, instance.PasswordAgePolicyChangedEventType})
	}
	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-Mlk6y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.PasswordExpirySettings()

		changes := make([]database.Change, 0, 2)
		if policyEvent.ExpireWarnDays != nil {
			changes = append(changes, settingsRepo.SetExpireWarnDays(*policyEvent.ExpireWarnDays))
		}
		if policyEvent.MaxAgeDays != nil {
			changes = append(changes, settingsRepo.SetMaxAgeDays(*policyEvent.MaxAgeDays))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, changes...)
	}), nil
}

func (s *settingsRelationalProjection) reducePasswordPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PasswordAgePolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-EtHWB", "reduce.wrong.event.type %s", org.PasswordAgePolicyRemovedEventType)
	}
	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.PasswordExpirySettings()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(policyEvent.Aggregate().InstanceID, &policyEvent.Aggregate().ID, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceOrgLockoutPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.LockoutPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.LockoutPolicyRemovedEventType)
	}
	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LockoutSettings()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(policyEvent.Aggregate().InstanceID, &policyEvent.Aggregate().ID, domain.SettingTypeLockout, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceLockoutPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LockoutPolicyAddedEvent
	switch e := event.(type) {
	case *org.LockoutPolicyAddedEvent:
		policyEvent = e.LockoutPolicyAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.LockoutPolicyAddedEvent:
		policyEvent = e.LockoutPolicyAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-d8mZO", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyAddedEventType, instance.LockoutPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hnNE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.LockoutSettings()
		settings := domain.LockoutSetting{
			Setting: domain.Setting{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			MaxPasswordAttempts: policyEvent.MaxPasswordAttempts,
			MaxOTPAttempts:      policyEvent.MaxOTPAttempts,
			ShowLockOutFailures: policyEvent.ShowLockOutFailures,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (s *settingsRelationalProjection) reduceLockoutPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LockoutPolicyChangedEvent
	switch e := event.(type) {
	case *org.LockoutPolicyChangedEvent:
		policyEvent = e.LockoutPolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.LockoutPolicyChangedEvent:
		policyEvent = e.LockoutPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gT3BQ", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyChangedEventType, instance.LockoutPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rbsxy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LockoutSettings()

		changes := make([]database.Change, 0, 3)
		if policyEvent.MaxPasswordAttempts != nil {
			changes = append(changes, settingsRepo.SetMaxPasswordAttempts(*policyEvent.MaxPasswordAttempts))
		}
		if policyEvent.MaxOTPAttempts != nil {
			changes = append(changes, settingsRepo.SetMaxOTPAttempts(*policyEvent.MaxOTPAttempts))
		}
		if policyEvent.ShowLockOutFailures != nil {
			changes = append(changes, settingsRepo.SetShowLockOutFailures(*policyEvent.ShowLockOutFailures))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, changes...)
	}), nil
}

func (s *settingsRelationalProjection) reduceDomainPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.DomainPolicyAddedEvent
	switch e := event.(type) {
	case *org.DomainPolicyAddedEvent:
		policyEvent = e.DomainPolicyAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.DomainPolicyAddedEvent:
		policyEvent = e.DomainPolicyAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-8se7M", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyAddedEventType, instance.DomainPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-chduE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.DomainSettings()
		settings := domain.DomainSetting{
			Setting: domain.Setting{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			LoginNameIncludesDomain:                policyEvent.UserLoginMustBeDomain,
			RequireOrgDomainVerification:           policyEvent.ValidateOrgDomains,
			SMTPSenderAddressMatchesInstanceDomain: policyEvent.SMTPSenderAddressMatchesInstanceDomain,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (s *settingsRelationalProjection) reduceDomainPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.DomainPolicyChangedEvent
	switch e := event.(type) {
	case *org.DomainPolicyChangedEvent:
		policyEvent = e.DomainPolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.DomainPolicyChangedEvent:
		policyEvent = e.DomainPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyChangedEventType, instance.DomainPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rbsxy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.DomainSettings()

		changes := make([]database.Change, 0, 3)
		if policyEvent.UserLoginMustBeDomain != nil {
			changes = append(changes, settingsRepo.SetLoginNameIncludesDomain(*policyEvent.UserLoginMustBeDomain))
		}
		if policyEvent.ValidateOrgDomains != nil {
			changes = append(changes, settingsRepo.SetRequireOrgDomainVerification(*policyEvent.ValidateOrgDomains))
		}
		if policyEvent.SMTPSenderAddressMatchesInstanceDomain != nil {
			changes = append(changes, settingsRepo.SetSMTPSenderAddressMatchesInstanceDomain(*policyEvent.SMTPSenderAddressMatchesInstanceDomain))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, changes...)
	}), nil
}

func (s *settingsRelationalProjection) reduceOrgDomainPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.DomainPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.LockoutPolicyRemovedEventType)
	}
	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.DomainSettings()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(policyEvent.Aggregate().InstanceID, &policyEvent.Aggregate().ID, domain.SettingTypeDomain, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceNotificationPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.NotificationPolicyAddedEvent
	switch e := event.(type) {
	case *org.NotificationPolicyAddedEvent:
		policyEvent = e.NotificationPolicyAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.NotificationPolicyAddedEvent:
		policyEvent = e.NotificationPolicyAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-8se7M", "reduce.wrong.event.type %v", []eventstore.EventType{org.NotificationPolicyAddedEventType, instance.NotificationPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-chduE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.NotificationSettings()
		settings := domain.NotificationSetting{
			Setting: domain.Setting{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			PasswordChange: policyEvent.PasswordChange,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (s *settingsRelationalProjection) reduceNotificationPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.NotificationPolicyChangedEvent
	switch e := event.(type) {
	case *org.NotificationPolicyChangedEvent:
		policyEvent = e.NotificationPolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.NotificationPolicyChangedEvent:
		policyEvent = e.NotificationPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.NotificationPolicyChangedEventType, instance.NotificationPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rbsxy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.NotificationSettings()

		changes := make([]database.Change, 0, 1)
		if policyEvent.PasswordChange != nil {
			changes = append(changes, settingsRepo.SetPasswordChange(*policyEvent.PasswordChange))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, changes...)
	}), nil
}

func (s *settingsRelationalProjection) reduceOrgNotificationPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.NotificationPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.NotificationPolicyRemovedEventType)
	}
	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.NotificationSettings()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(policyEvent.Aggregate().InstanceID, &policyEvent.Aggregate().ID, domain.SettingTypeNotification, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reducePrivacyPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.PrivacyPolicyAddedEvent
	switch e := event.(type) {
	case *org.PrivacyPolicyAddedEvent:
		policyEvent = e.PrivacyPolicyAddedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.PrivacyPolicyAddedEvent:
		policyEvent = e.PrivacyPolicyAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-8se7M", "reduce.wrong.event.type %v", []eventstore.EventType{org.PrivacyPolicyAddedEventType, instance.PrivacyPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-chduE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LegalAndSupportSettings()
		email := string(policyEvent.SupportEmail)
		settings := domain.LegalAndSupportSetting{
			Setting: domain.Setting{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			TOSLink:           policyEvent.TOSLink,
			PrivacyPolicyLink: policyEvent.PrivacyLink,
			HelpLink:          policyEvent.HelpLink,
			SupportEmail:      email,
			DocsLink:          policyEvent.DocsLink,
			CustomLink:        policyEvent.CustomLink,
			CustomLinkText:    policyEvent.CustomLinkText,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (s *settingsRelationalProjection) reducePrivacyPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.PrivacyPolicyChangedEvent
	switch e := event.(type) {
	case *org.PrivacyPolicyChangedEvent:
		policyEvent = e.PrivacyPolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.PrivacyPolicyChangedEvent:
		policyEvent = e.PrivacyPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.PrivacyPolicyChangedEventType, instance.PrivacyPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rbsxy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LegalAndSupportSettings()

		changes := make([]database.Change, 0, 7)
		if policyEvent.TOSLink != nil {
			changes = append(changes, settingsRepo.SetTOSLink(*policyEvent.TOSLink))
		}
		if policyEvent.PrivacyLink != nil {
			changes = append(changes, settingsRepo.SetPrivacyPolicyLink(*policyEvent.PrivacyLink))
		}
		if policyEvent.HelpLink != nil {
			changes = append(changes, settingsRepo.SetHelpLink(*policyEvent.HelpLink))
		}
		if policyEvent.SupportEmail != nil {
			changes = append(changes, settingsRepo.SetSupportEmail(string(*policyEvent.SupportEmail)))
		}
		if policyEvent.DocsLink != nil {
			changes = append(changes, settingsRepo.SetDocsLink(*policyEvent.DocsLink))
		}
		if policyEvent.CustomLink != nil {
			changes = append(changes, settingsRepo.SetCustomLink(*policyEvent.CustomLink))
		}
		if policyEvent.CustomLinkText != nil {
			changes = append(changes, settingsRepo.SetCustomLinkText(*policyEvent.CustomLinkText))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, changes...)
	}), nil
}

func (s *settingsRelationalProjection) reduceOrgPrivacyPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PrivacyPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.PrivacyPolicyRemovedEventType)
	}
	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LegalAndSupportSettings()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(policyEvent.Aggregate().InstanceID, &policyEvent.Aggregate().ID, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceSecurityPolicySet(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*instance.SecurityPolicySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-83U8p", "reduce.wrong.event.type %s", instance.SecurityPolicySetEventType)
	}

	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-lhPul", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SecuritySettings()

		changes := make([]database.Change, 0, 3)
		if policyEvent.EnableIframeEmbedding != nil {
			changes = append(changes, settingsRepo.SetEnableIframeEmbedding(*policyEvent.EnableIframeEmbedding))
		}
		if policyEvent.AllowedOrigins != nil {
			changes = append(changes, settingsRepo.SetAllowedOrigins(*policyEvent.AllowedOrigins))
		}
		if policyEvent.EnableImpersonation != nil {
			changes = append(changes, settingsRepo.SetEnableImpersonation(*policyEvent.EnableImpersonation))
		}

		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, nil, changes...)
	}), nil
}

func (s *settingsRelationalProjection) reduceOrganizationSettingsSet(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, err := assertEvent[*settings.OrganizationSettingsSetEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-lhPul", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.OrganizationSettings()
		orgID := &event.Aggregate().ID
		return settingsRepo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgID, settingsRepo.SetOrganizationScopedUsernames(policyEvent.OrganizationScopedUsernames))
	}), nil
}

func (s *settingsRelationalProjection) reduceOrganizationSettingsRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*settings.OrganizationSettingsRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rHiHb", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.OrganizationSettings()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			settingsRepo.UniqueCondition(e.Aggregate().InstanceID, &e.Aggregate().ID, domain.SettingTypeOrganization, domain.SettingStateActive),
		)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceSecretGeneratorAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecretGeneratorAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-yTuWKA", "reduce.wrong.event.type %s", instance.SecretGeneratorAddedEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-jNhP6P", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SecretGeneratorSettings()
		setting, err := setSecretGeneratorSettingsAttrs(
			e.GeneratorType,
			e.Length,
			e.IncludeLowerLetters,
			e.IncludeUpperLetters,
			e.IncludeDigits,
			e.IncludeSymbols,
			&e.Expiry,
		)
		if err != nil {
			return err
		}
		setting.Setting = domain.Setting{
			InstanceID:     event.Aggregate().InstanceID,
			OrganizationID: nil,
			CreatedAt:      e.Creation,
			UpdatedAt:      e.Creation,
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &setting)
	}), nil
}

func (s *settingsRelationalProjection) reduceSecretGeneratorChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecretGeneratorChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-HhaZbQ", "reduce.wrong.event.type %s", instance.SecretGeneratorChangedEventType)
	}

	repo := repository.SecretGeneratorSettings()

	changes := make([]database.Change, 0, 6)
	if e.Length != nil {
		changes = append(changes, repo.SetLength(*e.Length))
	}
	if e.Expiry != nil {
		changes = append(changes, repo.SetExpiry(*e.Expiry))
	}
	if e.IncludeLowerLetters != nil {
		changes = append(changes, repo.SetIncludeLowerLetters(*e.IncludeLowerLetters))
	}
	if e.IncludeUpperLetters != nil {
		changes = append(changes, repo.SetIncludeUpperLetters(*e.IncludeUpperLetters))
	}
	if e.IncludeDigits != nil {
		changes = append(changes, repo.SetIncludeDigits(*e.IncludeDigits))
	}
	if e.IncludeSymbols != nil {
		changes = append(changes, repo.SetIncludeSymbols(*e.IncludeSymbols))
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-3SidEv", "reduce.wrong.db.pool %T", ex)
		}

		return repo.Ensure(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, nil, mapSecretGeneratorTypeChanges(repo, e.GeneratorType, changes))
	}), nil
}

func mapSecretGeneratorTypeChanges(repo domain.SecretGeneratorSettingsRepository, generatorType legacy_domain.SecretGeneratorType, changes []database.Change) database.Change {
	switch generatorType {
	case legacy_domain.SecretGeneratorTypeInitCode:
		return repo.SetInitializeUserCodeSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeVerifyEmailCode:
		return repo.SetEmailVerificationCodeSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeVerifyPhoneCode:
		return repo.SetPhoneVerificationCodeSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeVerifyDomain:
		return repo.SetDomainVerificationSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypePasswordResetCode:
		return repo.SetPasswordVerificationCodeSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypePasswordlessInitCode:
		return repo.SetPasswordlessInitCodeSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeAppSecret:
		return repo.SetAppSecretSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeOTPSMS:
		return repo.SetOTPSMSSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeOTPEmail:
		return repo.SetOTPEmailSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeInviteCode:
		return repo.SetInviteCodeSecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeSigningKey:
		return repo.SetSigningKeySecretGenerator(changes...)
	case legacy_domain.SecretGeneratorTypeUnspecified:
		fallthrough
	default:
		panic(fmt.Sprintf("unknown secret generator type %s", generatorType))
	}
}

func (s *settingsRelationalProjection) reduceSecretGeneratorRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecretGeneratorRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-fmiIf", "reduce.wrong.event.type %s", instance.SecretGeneratorRemovedEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-xOUC9O", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.SecretGeneratorSettings()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			repo.UniqueCondition(e.Aggregate().InstanceID, nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
		)
		return err
	}), nil
}

func setSecretGeneratorSettingsAttrs(generatorType legacy_domain.SecretGeneratorType, length uint, includeLowerLetters, includeUpperLetters, includeDigits, includeSymbols bool, expiry *time.Duration) (domain.SecretGeneratorSetting, error) {
	var secretGeneratorSettingsAttrs domain.SecretGeneratorSetting
	attrs := domain.SecretGeneratorAttrs{
		Length:              length,
		IncludeLowerLetters: includeLowerLetters,
		IncludeUpperLetters: includeUpperLetters,
		IncludeDigits:       includeDigits,
		IncludeSymbols:      includeSymbols,
	}
	attrsWithExpiry := domain.SecretGeneratorAttrsWithExpiry{
		SecretGeneratorAttrs: attrs,
		Expiry:               expiry,
	}
	switch generatorType {
	case legacy_domain.SecretGeneratorTypeAppSecret:
		secretGeneratorSettingsAttrs.ClientSecret = &domain.ClientSecretAttributes{
			SecretGeneratorAttrsWithExpiry: attrsWithExpiry,
		}
	case legacy_domain.SecretGeneratorTypeInitCode:
		secretGeneratorSettingsAttrs.InitializeUserCode = &domain.InitializeUserCodeAttributes{
			SecretGeneratorAttrsWithExpiry: attrsWithExpiry,
		}
	case legacy_domain.SecretGeneratorTypeVerifyEmailCode:
		secretGeneratorSettingsAttrs.EmailVerificationCode = &domain.EmailVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: attrsWithExpiry,
		}
	case legacy_domain.SecretGeneratorTypeVerifyPhoneCode:
		secretGeneratorSettingsAttrs.PhoneVerificationCode = &domain.PhoneVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: attrsWithExpiry,
		}
	case legacy_domain.SecretGeneratorTypePasswordlessInitCode:
		secretGeneratorSettingsAttrs.PasswordlessInitCode = &domain.PasswordlessInitCodeAttributes{
			SecretGeneratorAttrsWithExpiry: attrsWithExpiry,
		}
	case legacy_domain.SecretGeneratorTypePasswordResetCode:
		secretGeneratorSettingsAttrs.PasswordVerificationCode = &domain.PasswordVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: attrsWithExpiry,
		}
	case legacy_domain.SecretGeneratorTypeVerifyDomain:
		secretGeneratorSettingsAttrs.DomainVerification = &domain.DomainVerificationAttributes{
			SecretGeneratorAttrs: attrs,
		}
	case legacy_domain.SecretGeneratorTypeOTPSMS:
		secretGeneratorSettingsAttrs.OTPSMS = &domain.OTPSMSAttributes{
			SecretGeneratorAttrsWithExpiry: attrsWithExpiry,
		}
	case legacy_domain.SecretGeneratorTypeOTPEmail:
		secretGeneratorSettingsAttrs.OTPEmail = &domain.OTPEmailAttributes{
			SecretGeneratorAttrsWithExpiry: attrsWithExpiry,
		}
	case legacy_domain.SecretGeneratorTypeInviteCode, legacy_domain.SecretGeneratorTypeSigningKey:
		// do nothing as these secret generators are not persisted in the settings
	case legacy_domain.SecretGeneratorTypeUnspecified:
		return domain.SecretGeneratorSetting{}, zerrors.ThrowInvalidArgumentf(nil, "HANDL-2n3fK", "unspecified secret generator type")
	default:
		return domain.SecretGeneratorSetting{}, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9mG7f", "unknown secret generator type %s", generatorType)
	}
	return secretGeneratorSettingsAttrs, nil
}
