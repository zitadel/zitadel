package projection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"slices"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	settings "github.com/zitadel/zitadel/internal/repository/organization_settings"
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
				// {
				// 	Event:  org.LabelPolicyActivatedEventType,
				// 	Reduce: s.reduceLabelActivated,
				// },
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
				{
					Event:  org.LabelPolicyAssetsRemovedEventType,
					Reduce: s.reduceAssetsRemoved,
				},
				// Password Complexity
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
				// Lockout Policy
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
				// Domain Policy
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
				// Delete org
				{
					Event:  org.OrgRemovedEventType,
					Reduce: s.reduceOrgRemoved,
				},
			},
		},
		// settings
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
				{
					Event:  instance.LabelPolicyAssetsRemovedEventType,
					Reduce: s.reduceAssetsRemoved,
				},
				// Password Complexity
				{
					Event:  instance.PasswordComplexityPolicyAddedEventType,
					Reduce: s.reducePassedComplexityAdded,
				},
				{
					Event:  instance.PasswordComplexityPolicyChangedEventType,
					Reduce: s.reducePasswordComplexityChanged,
				},
				// Password Policy
				{
					Event:  instance.PasswordAgePolicyAddedEventType,
					Reduce: s.reducePasswordPolicyAdded,
				},
				{
					Event:  instance.PasswordAgePolicyChangedEventType,
					Reduce: s.reducePasswordPolicyChanged,
				},
				// Lockout Policy
				{
					Event:  instance.LockoutPolicyAddedEventType,
					Reduce: s.reduceLockoutPolicyAdded,
				},
				{
					Event:  instance.LockoutPolicyChangedEventType,
					Reduce: s.reduceLockoutPolicyChanged,
				},
				// Domain Policy
				{
					Event:  instance.DomainPolicyAddedEventType,
					Reduce: s.reduceDomainPolicyAdded,
				},
				{
					Event:  instance.DomainPolicyChangedEventType,
					Reduce: s.reduceDomainPolicyChanged,
				},
				// Security Policy
				{
					Event:  instance.SecurityPolicySetEventType,
					Reduce: s.reduceSecurityPolicySet,
				},
				// Delete Instance
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: s.reduceInstanceRemoved,
				},
			},
		},
	}
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
		type settingStruct struct {
			policy.LoginPolicyAddedEvent
			IsDefault *bool `json:"isDefault,omitempty"`
		}

		loginPolicySetting := settingStruct{
			LoginPolicyAddedEvent: policyEvent,
			IsDefault:             &isDefault,
		}
		settingJSON, err := json.Marshal(loginPolicySetting)
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
			Settings:   settingJSON,
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

		_, err := settingsRepo.DeleteSettingsForInstance(ctx, removeInstanceEvent.Aggregate().InstanceID)
		return err
	}), nil
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
		type settingStruct struct {
			policy.LabelPolicyAddedEvent
			IsDefault *bool `json:"isDefault,omitempty"`
		}

		labelSetting := settingStruct{
			LabelPolicyAddedEvent: policyEvent,
			IsDefault:             &isDefault,
		}
		settings, err := json.Marshal(labelSetting)
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
		orgId = &e.Aggregate().ID
	case *instance.LabelPolicyLogoAddedEvent:
	// ok
	case *org.LabelPolicyLogoDarkAddedEvent:
		orgId = &e.Aggregate().ID
	case *instance.LabelPolicyLogoDarkAddedEvent:
	// ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-4UbiP", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoAddedEventType, instance.LabelPolicyLogoAddedEventType, org.LabelPolicyLogoDarkAddedEventType, instance.LabelPolicyLogoDarkAddedEventType})
	}

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

		switch e := event.(type) {
		case *org.LabelPolicyLogoAddedEvent:
			setting.Settings.LabelPolicyLightLogoURL = &e.StoreKey
		case *instance.LabelPolicyLogoAddedEvent:
			setting.Settings.LabelPolicyLightLogoURL = &e.StoreKey
		case *org.LabelPolicyLogoDarkAddedEvent:
			setting.Settings.LabelPolicyDarkLogoURL = &e.StoreKey
		case *instance.LabelPolicyLogoDarkAddedEvent:
			setting.Settings.LabelPolicyDarkLogoURL = &e.StoreKey
		}

		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceLogoRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	switch event.(type) {
	case *org.LabelPolicyLogoRemovedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyLogoRemovedEvent:
	case *org.LabelPolicyLogoDarkRemovedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyLogoDarkRemovedEvent:
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoRemovedEventType, instance.LabelPolicyLogoRemovedEventType, org.LabelPolicyLogoDarkRemovedEventType, instance.LabelPolicyLogoDarkRemovedEventType})
	}

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

		switch event.(type) {
		case *org.LabelPolicyLogoRemovedEvent:
			setting.Settings.LabelPolicyLightLogoURL = nil
		case *instance.LabelPolicyLogoRemovedEvent:
			setting.Settings.LabelPolicyLightLogoURL = nil
		case *org.LabelPolicyLogoDarkRemovedEvent:
			setting.Settings.LabelPolicyDarkLogoURL = nil
		case *instance.LabelPolicyLogoDarkRemovedEvent:
			setting.Settings.LabelPolicyDarkLogoURL = nil
		}

		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceIconAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	switch event.(type) {
	case *org.LabelPolicyIconAddedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyIconAddedEvent:
	case *org.LabelPolicyIconDarkAddedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyIconDarkAddedEvent:
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-e2JFz", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconAddedEventType, instance.LabelPolicyIconAddedEventType, org.LabelPolicyIconDarkAddedEventType, instance.LabelPolicyIconDarkAddedEventType})
	}

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

		switch e := event.(type) {
		case *org.LabelPolicyIconAddedEvent:
			setting.Settings.LabelPolicyLightIconURL = &e.StoreKey
		case *instance.LabelPolicyIconAddedEvent:
			setting.Settings.LabelPolicyLightIconURL = &e.StoreKey
		case *org.LabelPolicyIconDarkAddedEvent:
			setting.Settings.LabelPolicyDarkIconURL = &e.StoreKey
		case *instance.LabelPolicyIconDarkAddedEvent:
			setting.Settings.LabelPolicyDarkIconURL = &e.StoreKey
		}

		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceIconRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	switch event.(type) {
	case *org.LabelPolicyIconRemovedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyIconRemovedEvent:
	case *org.LabelPolicyIconDarkRemovedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyIconDarkRemovedEvent:
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gfgbY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconRemovedEventType, instance.LabelPolicyIconRemovedEventType, org.LabelPolicyIconDarkRemovedEventType, instance.LabelPolicyIconDarkRemovedEventType})
	}

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

		switch event.(type) {
		case *org.LabelPolicyIconRemovedEvent:
			setting.Settings.LabelPolicyLightIconURL = nil
		case *instance.LabelPolicyIconRemovedEvent:
			setting.Settings.LabelPolicyLightIconURL = nil
		case *org.LabelPolicyIconDarkRemovedEvent:
			setting.Settings.LabelPolicyDarkIconURL = nil
		case *instance.LabelPolicyIconDarkRemovedEvent:
			setting.Settings.LabelPolicyDarkIconURL = nil
		}

		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceFontAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	switch event.(type) {
	case *org.LabelPolicyFontAddedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyFontAddedEvent:
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-65i9W", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontAddedEventType, instance.LabelPolicyFontAddedEventType})
	}

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

		switch e := event.(type) {
		case *org.LabelPolicyFontAddedEvent:
			setting.Settings.LabelPolicyFontURL = &e.StoreKey
		case *instance.LabelPolicyFontAddedEvent:
			setting.Settings.LabelPolicyFontURL = &e.StoreKey
		}

		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
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
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		setting, err := settingsRepo.GetLabel(ctx, event.Aggregate().InstanceID, orgId)
		if err != nil {
			return zerrors.ThrowInternal(err, "HANDL-r7k9m", "error accessing login policy record")
		}

		setting.Settings.LabelPolicyFontURL = nil

		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceAssetsRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	switch event.(type) {
	case *org.LabelPolicyAssetsRemovedEvent:
		orgId = &event.Aggregate().ID
	case *instance.LabelPolicyAssetsRemovedEvent:
	// ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qi39A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAssetsRemovedEventType, instance.LabelPolicyAssetsRemovedEventType})
	}
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

		setting.Settings.LabelPolicyLightLogoURL = nil
		setting.Settings.LabelPolicyDarkLogoURL = nil
		setting.Settings.LabelPolicyLightIconURL = nil
		setting.Settings.LabelPolicyDarkIconURL = nil
		setting.Settings.LabelPolicyFontURL = nil

		_, err = settingsRepo.UpdateLabel(ctx, setting)
		return err
	}), nil
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
		type setting struct {
			policy.PasswordComplexityPolicyAddedEvent
			IsDefault *bool `json:"isDefault,omitempty"`
		}

		passwordComplexitySetting := setting{
			PasswordComplexityPolicyAddedEvent: policyEvent,
			IsDefault:                          &isDefault,
		}
		settingJSON, err := json.Marshal(passwordComplexitySetting)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
		newSetting := domain.Setting{
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			Type:       domain.SettingTypePasswordComplexity,
			Settings:   settingJSON,
		}
		err = settingsRepo.Create(ctx, &newSetting)
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
		type settingStruct struct {
			policy.PasswordAgePolicyAddedEvent
			IsDefault *bool `json:"isDefault,omitempty"`
		}

		passwordAgeSetting := settingStruct{
			PasswordAgePolicyAddedEvent: policyEvent,
			IsDefault:                   &isDefault,
		}

		settings, err := json.Marshal(passwordAgeSetting)
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

func (s *settingsRelationalProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-IoW0x", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgID := &e.Aggregate().ID

		_, err := settingsRepo.DeleteSettingsForOrg(ctx, orgID)
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
		type settingStruct struct {
			policy.LockoutPolicyAddedEvent
			IsDefault *bool `json:"isDefault,omitempty"`
		}

		loginSettings := settingStruct{
			LockoutPolicyAddedEvent: policyEvent,
			IsDefault:               &isDefault,
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
		type settingStruct struct {
			policy.DomainPolicyAddedEvent
			IsDefault *bool `json:"isDefault,omitempty"`
		}
		loginSettings := settingStruct{
			DomainPolicyAddedEvent: policyEvent,
			IsDefault:              &isDefault,
		}
		settingJSON, err := json.Marshal(loginSettings)
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
			Settings:   settingJSON,
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

func (s *settingsRelationalProjection) reduceOrganizationSettingsSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*settings.OrganizationSettingsSetEvent](event)
	if err != nil {
		return nil, err
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

		orgID := &e.Aggregate().ID

		existingSetting, err := settingsRepo.GetOrg(ctx, e.Agg.InstanceID, nil)
		if err != nil {
			if errors.Is(err, &database.NoRowFoundError{}) {
				setting := domain.Setting{
					InstanceID: e.Aggregate().InstanceID,
					OrgID:      orgID,
					Type:       domain.SettingTypeOrganization,
					Settings:   settings,
				}
				err = settingsRepo.Create(ctx, &setting)
				return err

			} else {
				return zerrors.ThrowInternal(err, "HANDL-rSkxt", "error accessing login policy record")
			}
		}

		existingSetting.Settings.OrganizationScopedUsernames = e.OrganizationScopedUsernames

		_, err = settingsRepo.UpdateOrg(ctx, existingSetting)
		return err
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

		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))

		orgId := &event.Aggregate().ID

		_, err := settingsRepo.Delete(ctx, e.Aggregate().InstanceID, orgId, domain.SettingTypeOrganization)
		return err
	}), nil
}
