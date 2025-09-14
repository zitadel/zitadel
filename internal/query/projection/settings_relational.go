package projection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"

	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	SettingsRelationalProjectionTable = "zitadel.settings"

	SettingInstanceIDCol = "instance_id"
	SettingsOrgIDCol     = "org_id"
	SettingsIDCol        = "id"
	SettingsTypeCol      = "type"
	SettingsSettingsCol  = "settings"
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
				// Login
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
				// {
				// 	Event:  org.OrgRemovedEventType,
				// 	Reduce: s.reduceOrgLoginPolicyRemoved,
				// },
				// label -----------------------------------------------------------
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
				// TODO
				// {
				// 	Event:  org.LabelPolicyActivatedEventType,
				// 	Reduce: s.reduceLabelActivated,
				// },
				{
					Event:  org.LabelPolicyLogoAddedEventType,
					Reduce: s.reduceLabelLogoAdded,
				},
				// {
				// 	Event:  org.LabelPolicyLogoRemovedEventType,
				// 	Reduce: s.reduceLogoRemoved,
				// },
				// {
				// 	Event:  org.LabelPolicyIconAddedEventType,
				// 	Reduce: s.reduceIconAdded,
				// },
				// {
				// 	Event:  org.LabelPolicyIconRemovedEventType,
				// 	Reduce: s.reduceIconRemoved,
				// },
				// {
				// 	Event:  org.LabelPolicyLogoDarkAddedEventType,
				// 	Reduce: s.reduceLogoAdded,
				// },
				// {
				// 	Event:  org.LabelPolicyLogoDarkRemovedEventType,
				// 	Reduce: s.reduceLogoRemoved,
				// },
				// {
				// 	Event:  org.LabelPolicyIconDarkAddedEventType,
				// 	Reduce: s.reduceIconAdded,
				// },
				// {
				// 	Event:  org.LabelPolicyIconDarkRemovedEventType,
				// 	Reduce: s.reduceIconRemoved,
				// },
				// {
				// 	Event:  org.LabelPolicyFontAddedEventType,
				// 	Reduce: s.reduceFontAdded,
				// },
				// {
				// 	Event:  org.LabelPolicyFontRemovedEventType,
				// 	Reduce: s.reduceFontRemoved,
				// },
				// {
				// 	Event:  org.LabelPolicyAssetsRemovedEventType,
				// 	Reduce: s.reduceAssetsRemoved,
				// },
				// {
				// 	Event:  org.OrgRemovedEventType,
				// 	Reduce: s.reduceOwnerRemoved,
				// },
				// Password Complexity -------------------------------------------------------
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
				// {
				// 	Event:  org.OrgRemovedEventType,
				// 	Reduce: s.reducePasswordComplexityOrgRemoved,
				// },
				// Password Policy
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
				// {
				// 	Event:  org.OrgRemovedEventType,
				// 	Reduce: s.reducePasswordPolicyOrgRemoved,
				// },
				// Lockout Policy -------------------------------------------------------
				{
					Event: org.LockoutPolicyAddedEventType,
					// Reduce: s.reduceOrgLockoutPolicyAdded,
					Reduce: s.reduceLockoutPolicyAdded,
				},
				{
					Event: org.LockoutPolicyChangedEventType,
					// Reduce: s.reduceOrgLockoutPolciyChanged,
					Reduce: s.reduceLockoutPolicyChanged,
				},
				{
					Event:  org.LockoutPolicyRemovedEventType,
					Reduce: s.reduceOrgLockoutPolicyRemoved,
				},
				// {
				// 	Event:  org.OrgRemovedEventType,
				// 	Reduce: s.reduceOrgLockoutPolicyOrgRemoved,
				// },
				// Domain Policy -------------------------------------------------------
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
				// {
				// 	Event:  org.OrgRemovedEventType,
				// 	Reduce: s.reduceOrgDomainPolicyOrgRemoved,
				// },
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				// Login
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
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: s.reduceInstanceRemoved,
				// },
				// Label
				{
					Event:  instance.LabelPolicyAddedEventType,
					Reduce: s.reduceLabelAdded,
				},
				{
					Event:  instance.LabelPolicyChangedEventType,
					Reduce: s.reduceLabelChanged,
				},
				// TODO
				// {
				// 	Event:  instance.LabelPolicyActivatedEventType,
				// 	Reduce: s.reduceLabelActivated,
				// },
				{
					Event:  instance.LabelPolicyLogoAddedEventType,
					Reduce: s.reduceLabelLogoAdded,
				},
				// {
				// 	Event:  instance.LabelPolicyLogoRemovedEventType,
				// 	Reduce: s.reduceLogoRemoved,
				// },
				// {
				// 	Event:  instance.LabelPolicyIconAddedEventType,
				// 	Reduce: s.reduceIconAdded,
				// },
				// {
				// 	Event:  instance.LabelPolicyIconRemovedEventType,
				// 	Reduce: s.reduceIconRemoved,
				// },
				// {
				// 	Event:  instance.LabelPolicyLogoDarkAddedEventType,
				// 	Reduce: s.reduceLogoAdded,
				// },
				// {
				// 	Event:  instance.LabelPolicyLogoDarkRemovedEventType,
				// 	Reduce: s.reduceLogoRemoved,
				// },
				// {
				// 	Event:  instance.LabelPolicyIconDarkAddedEventType,
				// 	Reduce: s.reduceIconAdded,
				// },
				// {
				// 	Event:  instance.LabelPolicyIconDarkRemovedEventType,
				// 	Reduce: s.reduceIconRemoved,
				// },
				// {
				// 	Event:  instance.LabelPolicyFontAddedEventType,
				// 	Reduce: s.reduceFontAdded,
				// },
				// {
				// 	Event:  instance.LabelPolicyFontRemovedEventType,
				// 	Reduce: s.reduceFontRemoved,
				// },
				// {
				// 	Event:  instance.LabelPolicyAssetsRemovedEventType,
				// 	Reduce: s.reduceAssetsRemoved,
				// },
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: reduceInstanceRemovedHelper(LabelPolicyInstanceIDCol),
				// },
				// Password Complexity -------------------------------------------------------
				{
					Event:  instance.PasswordComplexityPolicyAddedEventType,
					Reduce: s.reducePassedComplexityAdded,
				},
				{
					Event:  instance.PasswordComplexityPolicyChangedEventType,
					Reduce: s.reducePasswordComplexityChanged,
				},
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: s.reducePasswordComplexityInstanceRemoved,
				// },
				// Password Policy -------------------------------------------------------
				{
					Event:  instance.PasswordAgePolicyAddedEventType,
					Reduce: s.reducePasswordPolicyAdded,
				},
				{
					Event:  instance.PasswordAgePolicyChangedEventType,
					Reduce: s.reducePasswordPolicyChanged,
				},
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: s.reduceInstancePasswordPoicyRemoved,
				// },
				// Lockout Policy -------------------------------------------------------
				{
					Event:  instance.LockoutPolicyAddedEventType,
					Reduce: s.reduceLockoutPolicyAdded,
				},
				{
					Event:  instance.LockoutPolicyChangedEventType,
					Reduce: s.reduceLockoutPolicyChanged,
				},
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: s.reduceInstanceLockoutRemoved,
				// },
				// Domain Policy -------------------------------------------------------
				{
					Event:  instance.DomainPolicyAddedEventType,
					Reduce: s.reduceDomainPolicyAdded,
				},
				{
					Event:  instance.DomainPolicyChangedEventType,
					Reduce: s.reduceDomainPolicyChanged,
				},
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: s.reduceInstanceDomainPolicyRemoved,
				// },
				// Security Policy -------------------------------------------------------
				{
					Event:  instance.SecurityPolicySetEventType,
					Reduce: s.reduceSecurityPolicySet,
				},
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: s.reduceInstanceSecurityPoicyRemoved,
				// },
			},
		},
	}
}

type loginSettings struct {
	policy.LoginPolicyAddedEvent
	IsDefault *bool `json:"isDefault,omitempty"`
}

func (s *settingsRelationalProjection) reduceLoginPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LoginPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *instance.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		isDefault = true
	case *org.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		isDefault = false
		orgId = &policyEvent.Aggregate().ResourceOwner
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YYPxS", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyAddedEventType, instance.LoginPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		loginSettings := loginSettings{
			LoginPolicyAddedEvent: policyEvent,
			IsDefault:             &isDefault,
		}
		settings, err := json.Marshal(loginSettings)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
		setting := domain.Setting{
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			Type:       domain.SettingTypeLogin,
			Settings:   settings,
		}
		err = settingsRepo.Create(ctx, &setting)
		return err
	}), nil
}

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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLogin(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
		}

		if policyEvent.AllowRegister != nil {
			setting.Settings.AllowRegister = *policyEvent.AllowRegister
		}
		if policyEvent.AllowUserNamePassword != nil {
			setting.Settings.AllowUserNamePassword = *policyEvent.AllowUserNamePassword
		}
		if policyEvent.AllowExternalIDP != nil {
			setting.Settings.AllowExternalSetting = *policyEvent.AllowExternalIDP
		}
		if policyEvent.ForceMFA != nil {
			setting.Settings.ForceMFA = *policyEvent.ForceMFA
		}
		if policyEvent.ForceMFALocalOnly != nil {
			setting.Settings.ForceMFALocalOnly = *policyEvent.ForceMFALocalOnly
		}
		if policyEvent.PasswordlessType != nil {
			setting.Settings.PasswordlessType = domain.PasswordlessType(*policyEvent.PasswordlessType)
		}
		if policyEvent.HidePasswordReset != nil {
			setting.Settings.HidePasswordReset = *policyEvent.HidePasswordReset
		}
		if policyEvent.IgnoreUnknownUsernames != nil {
			setting.Settings.IgnoreUnknownUsernames = *policyEvent.IgnoreUnknownUsernames
		}
		if policyEvent.AllowDomainDiscovery != nil {
			setting.Settings.AllowDomainDiscovery = *policyEvent.AllowDomainDiscovery
		}
		if policyEvent.DisableLoginWithEmail != nil {
			setting.Settings.DisableLoginWithEmail = *policyEvent.DisableLoginWithEmail
		}
		if policyEvent.DisableLoginWithPhone != nil {
			setting.Settings.DisableLoginWithPhone = *policyEvent.DisableLoginWithPhone
		}
		if policyEvent.DefaultRedirectURI != nil {
			setting.Settings.DefaultRedirectURI = *policyEvent.DefaultRedirectURI
		}
		if policyEvent.PasswordCheckLifetime != nil {
			setting.Settings.PasswordCheckLifetime = *policyEvent.PasswordCheckLifetime
		}
		if policyEvent.ExternalLoginCheckLifetime != nil {
			setting.Settings.ExternalLoginCheckLifetime = *policyEvent.ExternalLoginCheckLifetime
		}
		if policyEvent.MFAInitSkipLifetime != nil {
			setting.Settings.MFAInitSkipLifetime = *policyEvent.MFAInitSkipLifetime
		}
		if policyEvent.SecondFactorCheckLifetime != nil {
			setting.Settings.SecondFactorCheckLifetime = *policyEvent.SecondFactorCheckLifetime
		}
		if policyEvent.MultiFactorCheckLifetime != nil {
			setting.Settings.MultiFactorCheckLifetime = *policyEvent.MultiFactorCheckLifetime
		}

		_, err = settingsRepo.UpdateLogin(ctx, setting)
		return err
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
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLogin(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
		}

		if slices.Contains(setting.Settings.MFAType, domain.MultiFactorType(policyEvent.MFAType)) {
			return nil
		}

		setting.Settings.MFAType = append(setting.Settings.MFAType, domain.MultiFactorType(policyEvent.MFAType))

		_, err = settingsRepo.UpdateLogin(ctx, setting)
		return err
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
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLogin(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
		}

		setting.Settings.MFAType = slices.DeleteFunc(setting.Settings.MFAType, func(mfaType domain.MultiFactorType) bool {
			return mfaType == domain.MultiFactorType(policyEvent.MFAType)
		})

		_, err = settingsRepo.UpdateLogin(ctx, setting)
		return err
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		_, err := settingsRepo.DeleteLogin(
			ctx,
			loginPolicyRemovedEvent.Aggregate().InstanceID,
			&loginPolicyRemovedEvent.Aggregate().ID)
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLogin(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-H7m9m", "error accessing login policy record")
		}

		if slices.Contains(setting.Settings.SecondFactorTypes, domain.SecondFactorType(policyEvent.MFAType)) {
			return nil
		}

		setting.Settings.SecondFactorTypes = append(setting.Settings.SecondFactorTypes, domain.SecondFactorType(policyEvent.MFAType))

		_, err = settingsRepo.UpdateLogin(ctx, setting)
		return err
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLogin(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-rsk9m", "error accessing login policy record")
		}

		setting.Settings.SecondFactorTypes = slices.DeleteFunc(setting.Settings.SecondFactorTypes, func(secondFactorType domain.SecondFactorType) bool {
			return secondFactorType == domain.SecondFactorType(policyEvent.MFAType)
		})

		_, err = settingsRepo.UpdateLogin(ctx, setting)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	removeInstanceEvent, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-28UlS", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewStatement(removeInstanceEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		_, err := settingsRepo.Delete(ctx, removeInstanceEvent.Aggregate().ID, nil, domain.SettingTypeLogin)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceOrgLoginPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	removeOrgEvent, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-T8NZa", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewStatement(removeOrgEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-arg9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		_, err := settingsRepo.Delete(ctx, removeOrgEvent.Aggregate().InstanceID, &removeOrgEvent.Aggregate().ID, domain.SettingTypeLogin)
		return err
	}), nil
}

type labelSettings struct {
	policy.LabelPolicyAddedEvent
	IsDefault *bool `json:"isDefault,omitempty"`
}

// label
func (s *settingsRelationalProjection) reduceLabelAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
		isDefault = false
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CSE7A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAddedEventType, instance.LabelPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		labelSettings := labelSettings{
			LabelPolicyAddedEvent: policyEvent,
			IsDefault:             &isDefault,
		}
		settings, err := json.Marshal(labelSettings)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
		setting := domain.Setting{
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			Type:       domain.SettingTypeLabel,
			Settings:   settings,
		}
		err = settingsRepo.Create(ctx, &setting)
		return err
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
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLabel(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
		}

		if policyEvent.PrimaryColor != nil {
			setting.Settings.PrimaryColor = *policyEvent.PrimaryColor
		}
		if policyEvent.BackgroundColor != nil {
			setting.Settings.BackgroundColor = *policyEvent.BackgroundColor
		}
		if policyEvent.WarnColor != nil {
			setting.Settings.WarnColor = *policyEvent.WarnColor
		}
		if policyEvent.FontColor != nil {
			setting.Settings.FontColor = *policyEvent.FontColor
		}
		if policyEvent.PrimaryColorDark != nil {
			setting.Settings.PrimaryColorDark = *policyEvent.PrimaryColorDark
		}
		if policyEvent.BackgroundColorDark != nil {
			setting.Settings.BackgroundColorDark = *policyEvent.BackgroundColorDark
		}
		if policyEvent.WarnColorDark != nil {
			setting.Settings.WarnColorDark = *policyEvent.WarnColorDark
		}
		if policyEvent.FontColorDark != nil {
			setting.Settings.FontColorDark = *policyEvent.FontColorDark
		}
		if policyEvent.HideLoginNameSuffix != nil {
			setting.Settings.HideLoginNameSuffix = *policyEvent.HideLoginNameSuffix
		}
		if policyEvent.ErrorMsgPopup != nil {
			setting.Settings.ErrorMsgPopup = *policyEvent.ErrorMsgPopup
		}
		if policyEvent.DisableWatermark != nil {
			setting.Settings.DisableWatermark = *policyEvent.DisableWatermark
		}
		if policyEvent.ThemeMode != nil {
			setting.Settings.ThemeMode = domain.LabelPolicyThemeMode(*policyEvent.ThemeMode)
		}
		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
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
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypeLabel)
		return err
	}), nil
}

// TODO
// func (p *settingsRelationalProjection) reduceLabelActivated(event eventstore.Event) (*handler.Statement, error) {
// 	switch event.(type) {
// 	case *org.LabelPolicyActivatedEvent, *instance.LabelPolicyActivatedEvent:
// 		// everything ok
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-7Kd8U", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyActivatedEventType, instance.LabelPolicyActivatedEventType})
// 	}
// 	return handler.NewStatement(&event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
// 		orgId := &event.Aggregate().ID

// 		setting, err := settingsRepo.GetLabel(ctx, event.Agg.InstanceID, orgId)
// 		if err != nil {
// 			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
// 		}

// 		setting.
// 		_, err = settingsRepo.UpdateLabel(ctx, setting)
// 		return err
// 	}), nil
// }

func (p *settingsRelationalProjection) reduceLabelLogoAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	switch e := event.(type) {
	case *org.LabelPolicyLogoAddedEvent:
		orgId = &e.Aggregate().ResourceOwner
	case *instance.LabelPolicyLogoAddedEvent:
	// ok
	case *org.LabelPolicyLogoDarkAddedEvent:
		orgId = &e.Aggregate().ResourceOwner
	case *instance.LabelPolicyLogoDarkAddedEvent:
	// ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-4UbiP", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoAddedEventType, instance.LabelPolicyLogoAddedEventType, org.LabelPolicyLogoDarkAddedEventType, instance.LabelPolicyLogoDarkAddedEventType})
	}

	fmt.Println("[DEBUGPRINT] [settings_org_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> LOGO")

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLabel(ctx, event.Aggregate().InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
		}

		fmt.Println("[DEBUGPRINT] [settings_org_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> LOGO")
		switch e := event.(type) {
		case *org.LabelPolicyLogoAddedEvent:
			setting.Settings.LabelPolicyLightLogoURL = e.StoreKey
		case *instance.LabelPolicyLogoAddedEvent:
			setting.Settings.LabelPolicyLightLogoURL = e.StoreKey
		case *org.LabelPolicyLogoDarkAddedEvent:
			setting.Settings.LabelPolicyDarkLogoURL = e.StoreKey
		case *instance.LabelPolicyLogoDarkAddedEvent:
			setting.Settings.LabelPolicyDarkLogoURL = e.StoreKey
		}

		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
	}), nil
}

// func (p *settingsRelationalProjection) reduceLogoRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var col string
// 	switch event.(type) {
// 	case *org.LabelPolicyLogoRemovedEvent:
// 		col = LabelPolicyLightLogoURLCol
// 	case *instance.LabelPolicyLogoRemovedEvent:
// 		col = LabelPolicyLightLogoURLCol
// 	case *org.LabelPolicyLogoDarkRemovedEvent:
// 		col = LabelPolicyDarkLogoURLCol
// 	case *instance.LabelPolicyLogoDarkRemovedEvent:
// 		col = LabelPolicyDarkLogoURLCol
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoRemovedEventType, instance.LabelPolicyLogoRemovedEventType, org.LabelPolicyLogoDarkRemovedEventType, instance.LabelPolicyLogoDarkRemovedEventType})
// 	}

// 	return handler.NewUpdateStatement(
// 		event,
// 		[]handler.Column{
// 			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
// 			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
// 			handler.NewCol(col, nil),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
// 			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
// 			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
// 		}), nil
// }

// func (p *settingsRelationalProjection) reduceIconAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var storeKey handler.Column
// 	switch e := event.(type) {
// 	case *org.LabelPolicyIconAddedEvent:
// 		storeKey = handler.NewCol(LabelPolicyLightIconURLCol, e.StoreKey)
// 	case *instance.LabelPolicyIconAddedEvent:
// 		storeKey = handler.NewCol(LabelPolicyLightIconURLCol, e.StoreKey)
// 	case *org.LabelPolicyIconDarkAddedEvent:
// 		storeKey = handler.NewCol(LabelPolicyDarkIconURLCol, e.StoreKey)
// 	case *instance.LabelPolicyIconDarkAddedEvent:
// 		storeKey = handler.NewCol(LabelPolicyDarkIconURLCol, e.StoreKey)
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-e2JFz", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconAddedEventType, instance.LabelPolicyIconAddedEventType, org.LabelPolicyIconDarkAddedEventType, instance.LabelPolicyIconDarkAddedEventType})
// 	}

// 	return handler.NewUpdateStatement(
// 		event,
// 		[]handler.Column{
// 			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
// 			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
// 			storeKey,
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
// 			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
// 			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
// 		}), nil
// }

// func (p *settingsRelationalProjection) reduceIconRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var col string
// 	switch event.(type) {
// 	case *org.LabelPolicyIconRemovedEvent:
// 		col = LabelPolicyLightIconURLCol
// 	case *instance.LabelPolicyIconRemovedEvent:
// 		col = LabelPolicyLightIconURLCol
// 	case *org.LabelPolicyIconDarkRemovedEvent:
// 		col = LabelPolicyDarkIconURLCol
// 	case *instance.LabelPolicyIconDarkRemovedEvent:
// 		col = LabelPolicyDarkIconURLCol
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gfgbY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconRemovedEventType, instance.LabelPolicyIconRemovedEventType, org.LabelPolicyIconDarkRemovedEventType, instance.LabelPolicyIconDarkRemovedEventType})
// 	}

// 	return handler.NewUpdateStatement(
// 		event,
// 		[]handler.Column{
// 			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
// 			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
// 			handler.NewCol(col, nil),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
// 			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
// 			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
// 		}), nil
// }

// func (p *settingsRelationalProjection) reduceFontAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var storeKey handler.Column
// 	switch e := event.(type) {
// 	case *org.LabelPolicyFontAddedEvent:
// 		storeKey = handler.NewCol(LabelPolicyFontURLCol, e.StoreKey)
// 	case *instance.LabelPolicyFontAddedEvent:
// 		storeKey = handler.NewCol(LabelPolicyFontURLCol, e.StoreKey)
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-65i9W", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontAddedEventType, instance.LabelPolicyFontAddedEventType})
// 	}

// 	return handler.NewUpdateStatement(
// 		event,
// 		[]handler.Column{
// 			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
// 			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
// 			storeKey,
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
// 			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
// 			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
// 		}), nil
// }

// func (p *settingsRelationalProjection) reduceFontRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var col string
// 	switch event.(type) {
// 	case *org.LabelPolicyFontRemovedEvent:
// 		col = LabelPolicyFontURLCol
// 	case *instance.LabelPolicyFontRemovedEvent:
// 		col = LabelPolicyFontURLCol
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-xf32J", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontRemovedEventType, instance.LabelPolicyFontRemovedEventType})
// 	}

// 	return handler.NewUpdateStatement(
// 		event,
// 		[]handler.Column{
// 			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
// 			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
// 			handler.NewCol(col, nil),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
// 			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
// 			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
// 		}), nil
// }

// func (p *settingsRelationalProjection) reduceAssetsRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	switch event.(type) {
// 	case *org.LabelPolicyAssetsRemovedEvent, *instance.LabelPolicyAssetsRemovedEvent:
// 		// ok
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qi39A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAssetsRemovedEventType, instance.LabelPolicyAssetsRemovedEventType})
// 	}

// 	return handler.NewUpdateStatement(
// 		event,
// 		[]handler.Column{
// 			handler.NewCol(LabelPolicyChangeDateCol, event.CreatedAt()),
// 			handler.NewCol(LabelPolicySequenceCol, event.Sequence()),
// 			handler.NewCol(LabelPolicyLightLogoURLCol, nil),
// 			handler.NewCol(LabelPolicyLightIconURLCol, nil),
// 			handler.NewCol(LabelPolicyDarkLogoURLCol, nil),
// 			handler.NewCol(LabelPolicyDarkIconURLCol, nil),
// 			handler.NewCol(LabelPolicyFontURLCol, nil),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(LabelPolicyIDCol, event.Aggregate().ID),
// 			handler.NewCond(LabelPolicyStateCol, domain.LabelPolicyStatePreview),
// 			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
// 		}), nil
// }

// func (p *settingsRelationalProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*org.OrgRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Su6pX", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
// 	}

// 	return handler.NewDeleteStatement(
// 		e,
// 		[]handler.Condition{
// 			handler.NewCond(LabelPolicyInstanceIDCol, e.Aggregate().InstanceID),
// 			handler.NewCond(LabelPolicyResourceOwnerCol, e.Aggregate().ID),
// 		},
// 	), nil
// }

type setting struct {
	any
	IsDefault *bool `json:"isDefault,omitempty"`
}

func (p *settingsRelationalProjection) reducePassedComplexityAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.PasswordComplexityPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.PasswordComplexityPolicyAddedEvent:
		policyEvent = e.PasswordComplexityPolicyAddedEvent
		isDefault = false
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.PasswordComplexityPolicyAddedEvent:
		policyEvent = e.PasswordComplexityPolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-KTHmJ", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordComplexityPolicyAddedEventType, instance.PasswordComplexityPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		loginSettings := setting{
			any:       policyEvent,
			IsDefault: &isDefault,
		}
		settings, err := json.Marshal(loginSettings)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
		setting := domain.Setting{
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			Type:       domain.SettingTypePasswordComplexity,
			Settings:   settings,
		}
		err = settingsRepo.Create(ctx, &setting)
		return err
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
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetPasswordComplexity(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
		}

		if policyEvent.MinLength != nil {
			setting.Settings.MinLength = *policyEvent.MinLength
		}
		if policyEvent.HasLowercase != nil {
			setting.Settings.HasLowercase = *policyEvent.HasLowercase
		}
		if policyEvent.HasUppercase != nil {
			setting.Settings.HasUppercase = *policyEvent.HasUppercase
		}
		if policyEvent.HasSymbol != nil {
			setting.Settings.HasSymbol = *policyEvent.HasSymbol
		}
		if policyEvent.HasNumber != nil {
			setting.Settings.HasNumber = *policyEvent.HasNumber
		}

		_, err = settingsRepo.UpdatePasswordComplexity(ctx, setting)
		return err
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypePasswordComplexity)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reducePasswordComplexityOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-pGTz9", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &e.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, e.Aggregate().InstanceID, orgID, domain.SettingTypePasswordComplexity)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reducePasswordComplexityInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	removeInstanceEvent, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-28UlS", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewStatement(removeInstanceEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		_, err := settingsRepo.Delete(ctx, removeInstanceEvent.Aggregate().InstanceID, nil, domain.SettingTypePasswordComplexity)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reducePasswordPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.PasswordAgePolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.PasswordAgePolicyAddedEvent:
		policyEvent = e.PasswordAgePolicyAddedEvent
		isDefault = false
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.PasswordAgePolicyAddedEvent:
		policyEvent = e.PasswordAgePolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CJqF0", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyAddedEventType, instance.PasswordAgePolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		loginSettings := setting{
			any:       policyEvent,
			IsDefault: &isDefault,
		}
		settings, err := json.Marshal(loginSettings)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
		setting := domain.Setting{
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			Type:       domain.SettingTypePasswordExpiry,
			Settings:   settings,
		}
		err = settingsRepo.Create(ctx, &setting)
		return err
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
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetPasswordExpiry(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
		}

		if policyEvent.ExpireWarnDays != nil {
			setting.Settings.ExpireWarnDays = *policyEvent.ExpireWarnDays
		}
		if policyEvent.MaxAgeDays != nil {
			setting.Settings.MaxAgeDays = *policyEvent.MaxAgeDays
		}

		_, err = settingsRepo.UpdatePasswordExpiry(ctx, setting)
		return err
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypePasswordExpiry)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reducePasswordPolicyOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-edLs2", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &e.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, e.Aggregate().InstanceID, orgID, domain.SettingTypePasswordExpiry)
		return err
	}), nil
}

// func (s *settingsRelationalProjection) reduceOrgLockoutPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.LockoutPolicyAddedEvent
// 	var isDefault bool
// 	switch e := event.(type) {
// 	case *org.LockoutPolicyAddedEvent:
// 		policyEvent = e.LockoutPolicyAddedEvent
// 		isDefault = false
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 	case *instance.LockoutPolicyAddedEvent:
// 		policyEvent = e.LockoutPolicyAddedEvent
// 		isDefault = true
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-d8mZO", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyAddedEventType, instance.LockoutPolicyAddedEventType})
// 	}
// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		loginSettings := setting{
// 			any:       policyEvent,
// 			IsDefault: &isDefault,
// 		}
// 		settings, err := json.Marshal(loginSettings)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Printf("[DEBUGPRINT] [settings_org_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> policyEvent = %+v\n", policyEvent)

// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5QPNE", "reduce.wrong.db.pool %T", ex)
// 		}
// 		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
// 		setting := domain.Setting{
// 			InstanceID: policyEvent.Aggregate().InstanceID,
// 			OrgID:      orgId,
// 			Type:       domain.SettingTypeLockout,
// 			Settings:   settings,
// 		}
// 		err = settingsRepo.Create(ctx, &setting)
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceOrgLockoutPolciyChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.LockoutPolicyChangedEvent
// 	switch e := event.(type) {
// 	case *org.LockoutPolicyChangedEvent:
// 		policyEvent = e.LockoutPolicyChangedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 	case *instance.LockoutPolicyChangedEvent:
// 		policyEvent = e.LockoutPolicyChangedEvent
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-pT3mQ", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyChangedEventType, instance.LockoutPolicyChangedEventType})
// 	}
// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

// 		setting, err := settingsRepo.GetLockout(ctx, policyEvent.Agg.InstanceID, orgId)
// 		if err != nil {
// 			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
// 		}

// 		if policyEvent.MaxPasswordAttempts != nil {
// 			setting.Settings.MaxPasswordAttempts = *policyEvent.MaxPasswordAttempts
// 		}
// 		if policyEvent.MaxOTPAttempts != nil {
// 			setting.Settings.MaxOTPAttempts = *policyEvent.MaxOTPAttempts
// 		}
// 		if policyEvent.ShowLockOutFailures != nil {
// 			setting.Settings.ShowLockOutFailures = *policyEvent.ShowLockOutFailures
// 		}

// 		_, err = settingsRepo.UpdateLockout(ctx, setting)
// 		return err
// 	}), nil
// }

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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypeLockout)
		return err
	}), nil
}

// func (s *settingsRelationalProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*org.OrgRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-IoW0x", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
// 	}

// 	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

// 		orgID := &e.Aggregate().ID

// 		_, err := settingsRepo.Delete(ctx, e.Aggregate().InstanceID, orgID, domain.SettingTypeLockout)
// 		return err
// 	}), nil
// }

func (s *settingsRelationalProjection) reduceOrgLockoutPolicyOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	removeOrgEvent, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-28UlS", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewStatement(removeOrgEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rV8Hy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		_, err := settingsRepo.Delete(ctx, removeOrgEvent.Aggregate().InstanceID, &removeOrgEvent.Aggregate().ID, domain.SettingTypeLockout)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceLockoutPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LockoutPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.LockoutPolicyAddedEvent:
		policyEvent = e.LockoutPolicyAddedEvent
		isDefault = false
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.LockoutPolicyAddedEvent:
		policyEvent = e.LockoutPolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-d8mZO", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyAddedEventType, instance.LockoutPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		loginSettings := setting{
			any:       policyEvent,
			IsDefault: &isDefault,
		}
		settings, err := json.Marshal(loginSettings)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hnNE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
		setting := domain.Setting{
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			Type:       domain.SettingTypeLockout,
			Settings:   settings,
		}
		err = settingsRepo.Create(ctx, &setting)
		return err
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLockout(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-rPkxm", "error accessing login policy record")
		}

		if policyEvent.MaxPasswordAttempts != nil {
			setting.Settings.MaxPasswordAttempts = *policyEvent.MaxPasswordAttempts
		}
		if policyEvent.MaxOTPAttempts != nil {
			setting.Settings.MaxOTPAttempts = *policyEvent.MaxOTPAttempts
		}
		if policyEvent.ShowLockOutFailures != nil {
			setting.Settings.ShowLockOutFailures = *policyEvent.ShowLockOutFailures
		}

		_, err = settingsRepo.UpdateLockout(ctx, setting)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceInstanceLockoutRemoved(event eventstore.Event) (*handler.Statement, error) {
	removeInstanceEvent, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-28gl9", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewStatement(removeInstanceEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-ZrdHz", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		_, err := settingsRepo.Delete(ctx, removeInstanceEvent.Aggregate().InstanceID, nil, domain.SettingTypeLockout)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceDomainPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.DomainPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.DomainPolicyAddedEvent:
		policyEvent = e.DomainPolicyAddedEvent
		isDefault = false
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.DomainPolicyAddedEvent:
		policyEvent = e.DomainPolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-8se7M", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyAddedEventType, instance.DomainPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		loginSettings := setting{
			any:       policyEvent,
			IsDefault: &isDefault,
		}
		settings, err := json.Marshal(loginSettings)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-chduE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
		setting := domain.Setting{
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			Type:       domain.SettingTypeDomain,
			Settings:   settings,
		}
		err = settingsRepo.Create(ctx, &setting)
		return err
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetDomain(ctx, policyEvent.Agg.InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-rPkxm", "error accessing login policy record")
		}

		if policyEvent.UserLoginMustBeDomain != nil {
			setting.Settings.UserLoginMustBeDomain = *policyEvent.UserLoginMustBeDomain
		}
		if policyEvent.ValidateOrgDomains != nil {
			setting.Settings.ValidateOrgDomains = *policyEvent.ValidateOrgDomains
		}
		if policyEvent.SMTPSenderAddressMatchesInstanceDomain != nil {
			setting.Settings.SMTPSenderAddressMatchesInstanceDomain = *policyEvent.SMTPSenderAddressMatchesInstanceDomain
		}

		_, err = settingsRepo.UpdateDomain(ctx, setting)
		return err
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &policyEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypeDomain)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceOrgDomainPolicyOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	removeOrgEvent, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-28UlS", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewStatement(removeOrgEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rV8Hy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &removeOrgEvent.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, removeOrgEvent.Aggregate().InstanceID, orgID, domain.SettingTypeDomain)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceInstanceDomainPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	removeInstanceEvent, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-78jl9", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewStatement(removeInstanceEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-zrdJz", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		_, err := settingsRepo.Delete(ctx, removeInstanceEvent.Aggregate().InstanceID, nil, domain.SettingTypeDomain)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceSecurityPolicySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecurityPolicySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-83U8p", "reduce.wrong.event.type %s", instance.SecurityPolicySetEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		settings, err := json.Marshal(e)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-chluS", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		existingSetting, err := settingsRepo.GetSecurity(ctx, e.Agg.InstanceID, nil)
		if err != nil {
			if errors.Is(err, &database.NoRowFoundError{}) {
				setting := domain.Setting{
					InstanceID: e.Aggregate().InstanceID,
					Type:       domain.SettingTypeSecurity,
					Settings:   settings,
				}
				err = settingsRepo.Create(ctx, &setting)
				return err

			} else {
				return zerrors.ThrowInternal(err, "HANDL-rSkxt", "error accessing login policy record")
			}
		}

		if e.EnableIframeEmbedding != nil {
			existingSetting.Settings.EnableIframeEmbedding = *e.EnableIframeEmbedding
		} else if e.Enabled != nil {
			existingSetting.Settings.Enabled = *e.Enabled
		}
		if e.AllowedOrigins != nil {
			existingSetting.Settings.AllowedOrigins = *e.AllowedOrigins
		}
		if e.EnableImpersonation != nil {
			existingSetting.Settings.EnableImpersonation = *e.EnableImpersonation
		}
		_, err = settingsRepo.UpdateSecurity(ctx, existingSetting)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceInstanceSecurityPoicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	removeInstanceEvent, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-b88lS", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewStatement(removeInstanceEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rHiHb", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		_, err := settingsRepo.Delete(ctx, removeInstanceEvent.Aggregate().InstanceID, nil, domain.SettingTypeSecurity)
		return err
	}), nil
}
