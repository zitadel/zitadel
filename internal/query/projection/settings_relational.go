package projection

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
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

	SettingIDCol          = "id"
	SettingInstanceIDCol  = "instance_id"
	SettingsOrgIDCol      = "org_id"
	SettingsTypeCol       = "type"
	SettingsLabelStateCol = "label_state"
	SettingsSettingsCol   = "settings"
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
					Event:  org.LoginPolicyMultiFactorAddedEventType,
					Reduce: s.reduceMFAAdded,
				},
				{
					Event:  org.LoginPolicyMultiFactorRemovedEventType,
					Reduce: s.reduceMFARemoved,
				},
				{
					Event:  org.LoginPolicyRemovedEventType,
					Reduce: s.reduceLoginPolicyRemoved,
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
					Reduce: s.reduceLabelLogoAdded,
				},
				{
					Event:  org.LabelPolicyLogoRemovedEventType,
					Reduce: s.reduceLogoRemoved,
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
					Event:  org.LabelPolicyLogoDarkAddedEventType,
					Reduce: s.reduceLabelLogoAdded,
				},
				{
					Event:  org.LabelPolicyLogoDarkRemovedEventType,
					Reduce: s.reduceLogoRemoved,
				},
				{
					Event:  org.LabelPolicyIconDarkAddedEventType,
					Reduce: s.reduceIconAdded,
				},
				{
					Event:  org.LabelPolicyIconDarkRemovedEventType,
					Reduce: s.reduceIconRemoved,
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
					Reduce: s.reducePassedComplexityAdded,
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
					Reduce: s.reduceLabelLogoAdded,
				},
				{
					Event:  instance.LabelPolicyLogoRemovedEventType,
					Reduce: s.reduceLogoRemoved,
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
					Event:  instance.LabelPolicyLogoDarkAddedEventType,
					Reduce: s.reduceLabelLogoAdded,
				},
				{
					Event:  instance.LabelPolicyLogoDarkRemovedEventType,
					Reduce: s.reduceLogoRemoved,
				},
				{
					Event:  instance.LabelPolicyIconDarkAddedEventType,
					Reduce: s.reduceIconAdded,
				},
				{
					Event:  instance.LabelPolicyIconDarkRemovedEventType,
					Reduce: s.reduceIconRemoved,
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
					Reduce: s.reducePassedComplexityAdded,
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
		settingsRepo := repository.LoginSettingsRepository()

		passwordlessType := domain.PasswordlessType(policyEvent.PasswordlessType)
		setting := domain.LoginSettings{
			Settings: domain.Settings{
				InstanceID:     policyEvent.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      &policyEvent.AllowUserNamePassword,
				AllowRegister:              &policyEvent.AllowRegister,
				AllowExternalIDP:           &policyEvent.AllowExternalIDP,
				ForceMFA:                   &policyEvent.ForceMFA,
				ForceMFALocalOnly:          &policyEvent.ForceMFALocalOnly,
				HidePasswordReset:          &policyEvent.HidePasswordReset,
				IgnoreUnknownUsernames:     &policyEvent.IgnoreUnknownUsernames,
				AllowDomainDiscovery:       &policyEvent.AllowDomainDiscovery,
				DisableLoginWithEmail:      &policyEvent.DisableLoginWithEmail,
				DisableLoginWithPhone:      &policyEvent.DisableLoginWithPhone,
				PasswordlessType:           &passwordlessType,
				DefaultRedirectURI:         &policyEvent.DefaultRedirectURI,
				PasswordCheckLifetime:      &policyEvent.PasswordCheckLifetime,
				ExternalLoginCheckLifetime: &policyEvent.ExternalLoginCheckLifetime,
				MFAInitSkipLifetime:        &policyEvent.MFAInitSkipLifetime,
				SecondFactorCheckLifetime:  &policyEvent.SecondFactorCheckLifetime,
				MultiFactorCheckLifetime:   &policyEvent.MultiFactorCheckLifetime,
				MFAType:                    nil,
				SecondFactorTypes:          nil,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &setting)
	}), nil
}

// //nolint:gocognit
func (s *settingsRelationalProjection) reduceLoginPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LoginPolicyChangedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyChangedEvent:
		policyEvent = e.LoginPolicyChangedEvent
	case *org.LoginPolicyChangedEvent:
		policyEvent = e.LoginPolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BHd86", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyChangedEventType, instance.LoginPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LoginSettingsRepository()

		var passwordlessType *domain.PasswordlessType
		if policyEvent.PasswordlessType != nil {
			passwordlessType = gu.Ptr(domain.PasswordlessType(*policyEvent.PasswordlessType))
		}

		settings := domain.LoginSettings{
			Settings: domain.Settings{
				InstanceID:     policyEvent.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      policyEvent.Creation,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      policyEvent.AllowUserNamePassword,
				AllowRegister:              policyEvent.AllowRegister,
				AllowExternalIDP:           policyEvent.AllowExternalIDP,
				ForceMFA:                   policyEvent.ForceMFA,
				ForceMFALocalOnly:          policyEvent.ForceMFALocalOnly,
				HidePasswordReset:          policyEvent.HidePasswordReset,
				IgnoreUnknownUsernames:     policyEvent.IgnoreUnknownUsernames,
				AllowDomainDiscovery:       policyEvent.AllowDomainDiscovery,
				DisableLoginWithEmail:      policyEvent.DisableLoginWithEmail,
				DisableLoginWithPhone:      policyEvent.DisableLoginWithPhone,
				PasswordlessType:           passwordlessType,
				DefaultRedirectURI:         policyEvent.DefaultRedirectURI,
				PasswordCheckLifetime:      policyEvent.PasswordCheckLifetime,
				ExternalLoginCheckLifetime: policyEvent.ExternalLoginCheckLifetime,
				MFAInitSkipLifetime:        policyEvent.MFAInitSkipLifetime,
				SecondFactorCheckLifetime:  policyEvent.SecondFactorCheckLifetime,
				MultiFactorCheckLifetime:   policyEvent.MultiFactorCheckLifetime,
				MFAType:                    nil,
				SecondFactorTypes:          nil,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
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

		settingsRepo := repository.LoginSettingsRepository()
		return settingsRepo.SetColumns(ctx, v3_sql.SQLTx(tx),
			&domain.Settings{
				InstanceID:     policyEvent.Agg.InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      policyEvent.Creation,
			},
			settingsRepo.AddMFAType(domain.MultiFactorType(policyEvent.MFAType)),
		)
	}), nil
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

		settingsRepo := repository.LoginSettingsRepository()
		return settingsRepo.SetColumns(ctx, v3_sql.SQLTx(tx),
			&domain.Settings{
				InstanceID:     policyEvent.Agg.InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      policyEvent.Creation,
			},
			settingsRepo.RemoveMFAType(domain.MultiFactorType(policyEvent.MFAType)),
		)
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

		settingsRepo := repository.LoginSettingsRepository()
		_, err := settingsRepo.Delete(
			ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(loginPolicyRemovedEvent.Aggregate().InstanceID),
				settingsRepo.OrganizationIDCondition(&loginPolicyRemovedEvent.Aggregate().ID),
			))

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

		settingsRepo := repository.LoginSettingsRepository()
		return settingsRepo.SetColumns(ctx, v3_sql.SQLTx(tx),
			&domain.Settings{
				InstanceID:     policyEvent.Agg.InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      policyEvent.Creation,
			},
			settingsRepo.AddSecondFactorTypes(domain.SecondFactorType(policyEvent.MFAType)),
		)
	}), nil
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
		settingsRepo := repository.LoginSettingsRepository()
		return settingsRepo.SetColumns(ctx, v3_sql.SQLTx(tx),
			&domain.Settings{
				InstanceID:     policyEvent.Agg.InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      policyEvent.Creation,
			},
			settingsRepo.RemoveSecondFactorTypes(domain.SecondFactorType(policyEvent.MFAType)),
		)
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
		settingsRepo := repository.BrandingSettingsRepository()
		themeMode := domain.BrandingPolicyThemeMode(policyEvent.ThemeMode)
		settings := domain.BrandingSettings{
			Settings: domain.Settings{
				InstanceID:     policyEvent.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    &policyEvent.PrimaryColor,
				BackgroundColorLight: &policyEvent.BackgroundColor,
				WarnColorLight:       &policyEvent.WarnColor,
				FontColorLight:       &policyEvent.FontColor,
				PrimaryColorDark:     &policyEvent.PrimaryColorDark,
				BackgroundColorDark:  &policyEvent.BackgroundColorDark,
				WarnColorDark:        &policyEvent.WarnColorDark,
				FontColorDark:        &policyEvent.FontColorDark,
				HideLoginNameSuffix:  &policyEvent.HideLoginNameSuffix,
				ErrorMsgPopup:        &policyEvent.ErrorMsgPopup,
				DisableWatermark:     &policyEvent.DisableWatermark,
				ThemeMode:            &themeMode,
				LogoURLLight:         nil,
				IconURLLight:         nil,
				LogoURLDark:          nil,
				IconURLDark:          nil,
				FontURL:              nil,
			},
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
		settingsRepo := repository.BrandingSettingsRepository()
		var themeMode *domain.BrandingPolicyThemeMode
		if policyEvent.ThemeMode != nil {
			themeMode = gu.Ptr(domain.BrandingPolicyThemeMode(*policyEvent.ThemeMode))
		}
		settings := domain.BrandingSettings{
			Settings: domain.Settings{
				InstanceID:     policyEvent.Aggregate().InstanceID,
				OrganizationID: orgId,
				CreatedAt:      policyEvent.Creation,
				UpdatedAt:      policyEvent.Creation,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    policyEvent.PrimaryColor,
				BackgroundColorLight: policyEvent.BackgroundColor,
				WarnColorLight:       policyEvent.WarnColor,
				FontColorLight:       policyEvent.FontColor,
				PrimaryColorDark:     policyEvent.PrimaryColorDark,
				BackgroundColorDark:  policyEvent.BackgroundColorDark,
				WarnColorDark:        policyEvent.WarnColorDark,
				FontColorDark:        policyEvent.FontColorDark,
				HideLoginNameSuffix:  policyEvent.HideLoginNameSuffix,
				ErrorMsgPopup:        policyEvent.ErrorMsgPopup,
				DisableWatermark:     policyEvent.DisableWatermark,
				ThemeMode:            themeMode,
				LogoURLLight:         nil,
				IconURLLight:         nil,
				LogoURLDark:          nil,
				IconURLDark:          nil,
				FontURL:              nil,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
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

		settingsRepo := repository.BrandingSettingsRepository()

		orgId := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(policyEvent.Agg.InstanceID),
				settingsRepo.OrganizationIDCondition(orgId),
			))
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

		settingsRepo := repository.BrandingSettingsRepository()
		_, err := settingsRepo.Activate(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(event.Aggregate().InstanceID),
				settingsRepo.OrganizationIDCondition(orgId),
			),
			settingsRepo.SetUpdatedAt(&policyEvent.Creation),
		)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceLabelLogoAdded(event eventstore.Event) (*handler.Statement, error) {
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettingsRepository()

		updatedAt := event.CreatedAt()
		settings := domain.BrandingSettings{
			Settings: domain.Settings{
				InstanceID: event.Aggregate().InstanceID,
				Type:       domain.SettingTypeBranding,
				UpdatedAt:  updatedAt,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{},
		}

		switch e := event.(type) {
		case *org.LabelPolicyLogoAddedEvent:
			settings.OrganizationID = &e.Aggregate().ResourceOwner
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			settings.LogoURLLight = url
		case *instance.LabelPolicyLogoAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			settings.LogoURLLight = url
		case *org.LabelPolicyLogoDarkAddedEvent:
			settings.OrganizationID = &e.Aggregate().ResourceOwner
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			settings.LogoURLDark = url
		case *instance.LabelPolicyLogoDarkAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			settings.LogoURLDark = url
		}

		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reduceLogoRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var logoLight, logoDark *url.URL
	switch event.(type) {
	case *org.LabelPolicyLogoRemovedEvent:
		orgId = &event.Aggregate().ID
		logoLight = &url.URL{}
	case *instance.LabelPolicyLogoRemovedEvent:
		logoLight = &url.URL{}
	case *org.LabelPolicyLogoDarkRemovedEvent:
		orgId = &event.Aggregate().ID
		logoDark = &url.URL{}
	case *instance.LabelPolicyLogoDarkRemovedEvent:
		logoDark = &url.URL{}
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoRemovedEventType, instance.LabelPolicyLogoRemovedEventType, org.LabelPolicyLogoDarkRemovedEventType, instance.LabelPolicyLogoDarkRemovedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.BrandingSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.BrandingSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				LogoURLLight: logoLight,
				LogoURLDark:  logoDark,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reduceIconAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var iconLight, iconDark *url.URL
	switch e := event.(type) {
	case *org.LabelPolicyIconAddedEvent:
		orgId = &event.Aggregate().ID
		url, err := url.Parse(e.StoreKey)
		if err != nil {
			return nil, err
		}
		iconLight = url
	case *instance.LabelPolicyIconAddedEvent:
		url, err := url.Parse(e.StoreKey)
		if err != nil {
			return nil, err
		}
		iconLight = url
	case *org.LabelPolicyIconDarkAddedEvent:
		orgId = &event.Aggregate().ID
		url, err := url.Parse(e.StoreKey)
		if err != nil {
			return nil, err
		}
		iconDark = url
	case *instance.LabelPolicyIconDarkAddedEvent:
		url, err := url.Parse(e.StoreKey)
		if err != nil {
			return nil, err
		}
		iconDark = url
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-e2JFz", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconAddedEventType, instance.LabelPolicyIconAddedEventType, org.LabelPolicyIconDarkAddedEventType, instance.LabelPolicyIconDarkAddedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.BrandingSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.BrandingSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				IconURLLight: iconLight,
				IconURLDark:  iconDark,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reduceIconRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var iconLight, iconDark *url.URL
	switch event.(type) {
	case *org.LabelPolicyIconRemovedEvent:
		orgId = &event.Aggregate().ID
		iconLight = &url.URL{}
	case *instance.LabelPolicyIconRemovedEvent:
		iconLight = &url.URL{}
	case *org.LabelPolicyIconDarkRemovedEvent:
		orgId = &event.Aggregate().ID
		iconDark = &url.URL{}
	case *instance.LabelPolicyIconDarkRemovedEvent:
		iconDark = &url.URL{}
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gfgbY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconRemovedEventType, instance.LabelPolicyIconRemovedEventType, org.LabelPolicyIconDarkRemovedEventType, instance.LabelPolicyIconDarkRemovedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.BrandingSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				IconURLLight: iconLight,
				IconURLDark:  iconDark,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reduceFontAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var font *url.URL
	switch e := event.(type) {
	case *org.LabelPolicyFontAddedEvent:
		orgId = &event.Aggregate().ID
		url, err := url.Parse(e.StoreKey)
		if err != nil {
			return nil, err
		}
		font = url
	case *instance.LabelPolicyFontAddedEvent:
		url, err := url.Parse(e.StoreKey)
		if err != nil {
			return nil, err
		}
		font = url
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-65i9W", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontAddedEventType, instance.LabelPolicyFontAddedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.BrandingSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				FontURL: font,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reduceFontRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	switch event.(type) {
	case *org.LabelPolicyFontRemovedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyFontRemovedEvent:
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-xf32J", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontRemovedEventType, instance.LabelPolicyFontRemovedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.BrandingSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.BrandingSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				FontURL: &url.URL{},
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reducePassedComplexityAdded(event eventstore.Event) (*handler.Statement, error) {
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

		settingsRepo := repository.PasswordComplexitySettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.PasswordComplexitySettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    &policyEvent.MinLength,
				HasLowercase: &policyEvent.HasLowercase,
				HasUppercase: &policyEvent.HasUppercase,
				HasNumber:    &policyEvent.HasNumber,
				HasSymbol:    &policyEvent.HasSymbol,
			},
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

		settingsRepo := repository.PasswordComplexitySettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.PasswordComplexitySettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    policyEvent.MinLength,
				HasLowercase: policyEvent.HasLowercase,
				HasUppercase: policyEvent.HasUppercase,
				HasNumber:    policyEvent.HasNumber,
				HasSymbol:    policyEvent.HasSymbol,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
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

		settingsRepo := repository.PasswordComplexitySettingsRepository()

		orgId := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(policyEvent.Aggregate().InstanceID),
				settingsRepo.OrganizationIDCondition(orgId),
			))
		return err
	}), nil
}

func (p *settingsRelationalProjection) reducePasswordPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
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

		settingsRepo := repository.PasswordExpiryRepository()
		updatedAt := event.CreatedAt()
		settings := domain.PasswordExpirySettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: &policyEvent.ExpireWarnDays,
				MaxAgeDays:     &policyEvent.MaxAgeDays,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reducePasswordPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
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

		settingsRepo := repository.PasswordExpiryRepository()
		updatedAt := event.CreatedAt()
		settings := domain.PasswordExpirySettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: policyEvent.ExpireWarnDays,
				MaxAgeDays:     policyEvent.MaxAgeDays,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reducePasswordPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PasswordAgePolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-EtHWB", "reduce.wrong.event.type %s", org.PasswordAgePolicyRemovedEventType)
	}
	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.PasswordExpiryRepository()

		orgId := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(policyEvent.Aggregate().InstanceID),
				settingsRepo.OrganizationIDCondition(orgId),
			))
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceOrgLockoutPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.LockoutPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.LockoutPolicyRemovedEventType)
	}
	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.LockoutSettingsRepository()

		orgId := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(policyEvent.Aggregate().InstanceID),
				settingsRepo.OrganizationIDCondition(orgId),
			))
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceLockoutPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
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
		settingsRepo := repository.LockoutSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.LockoutSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: &policyEvent.MaxPasswordAttempts,
				MaxOTPAttempts:      &policyEvent.MaxOTPAttempts,
				ShowLockOutFailures: &policyEvent.ShowLockOutFailures,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reduceLockoutPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
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

		settingsRepo := repository.LockoutSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.LockoutSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: policyEvent.MaxPasswordAttempts,
				MaxOTPAttempts:      policyEvent.MaxOTPAttempts,
				ShowLockOutFailures: policyEvent.ShowLockOutFailures,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reduceDomainPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
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

		settingsRepo := repository.DomainSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.DomainSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                &policyEvent.UserLoginMustBeDomain,
				RequireOrgDomainVerification:           &policyEvent.ValidateOrgDomains,
				SMTPSenderAddressMatchesInstanceDomain: &policyEvent.SMTPSenderAddressMatchesInstanceDomain,
			},
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

		settingsRepo := repository.DomainSettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.DomainSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: orgId,
				UpdatedAt:      updatedAt,
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                policyEvent.UserLoginMustBeDomain,
				RequireOrgDomainVerification:           policyEvent.ValidateOrgDomains,
				SMTPSenderAddressMatchesInstanceDomain: policyEvent.SMTPSenderAddressMatchesInstanceDomain,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
	}), nil
}

func (p *settingsRelationalProjection) reduceOrgDomainPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.DomainPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.LockoutPolicyRemovedEventType)
	}
	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.DomainSettingsRepository()

		orgId := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(policyEvent.Aggregate().InstanceID),
				settingsRepo.OrganizationIDCondition(orgId),
			))
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

		var allowedOrigins []string
		if policyEvent.AllowedOrigins != nil {
			allowedOrigins = *policyEvent.AllowedOrigins
		}

		settingsRepo := repository.SecuritySettingsRepository()
		updatedAt := event.CreatedAt()
		settings := domain.SecuritySettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: nil,
				UpdatedAt:      updatedAt,
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: policyEvent.EnableIframeEmbedding,
				AllowedOrigins:        allowedOrigins,
				EnableImpersonation:   policyEvent.EnableImpersonation,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
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

		settingsRepo := repository.OrganizationSettingRepository()
		updatedAt := event.CreatedAt()
		settings := domain.OrganizationSettings{
			Settings: domain.Settings{
				InstanceID:     event.Aggregate().InstanceID,
				OrganizationID: &event.Aggregate().ID,
				UpdatedAt:      updatedAt,
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: &policyEvent.OrganizationScopedUsernames,
			},
		}
		return settingsRepo.Set(ctx, v3_sql.SQLTx(tx), &settings)
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

		settingsRepo := repository.OrganizationSettingRepository()
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(e.Agg.InstanceID),
				settingsRepo.OrganizationIDCondition(&event.Aggregate().ID),
			))
		return err
	}), nil
}
