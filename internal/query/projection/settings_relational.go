package projection

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	db_json "github.com/zitadel/zitadel/backend/v3/storage/database/json"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
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
				// 		{
				// 			Event:  org.LoginPolicyAddedEventType,
				// 			Reduce: s.reduceLoginPolicyAdded,
				// 		},
				// 		{
				// 			Event:  org.LoginPolicyChangedEventType,
				// 			Reduce: s.reduceLoginPolicyChanged,
				// 		},
				// 		{
				// 			Event:  org.LoginPolicyMultiFactorAddedEventType,
				// 			Reduce: s.reduceMFAAdded,
				// 		},
				// 		{
				// 			Event:  org.LoginPolicyMultiFactorRemovedEventType,
				// 			Reduce: s.reduceMFARemoved,
				// 		},
				// 		{
				// 			Event:  org.LoginPolicyRemovedEventType,
				// 			Reduce: s.reduceLoginPolicyRemoved,
				// 		},
				// 		{
				// 			Event:  org.LoginPolicySecondFactorAddedEventType,
				// 			Reduce: s.reduceSecondFactorAdded,
				// 		},
				// 		{
				// 			Event:  org.LoginPolicySecondFactorRemovedEventType,
				// 			Reduce: s.reduceSecondFactorRemoved,
				// 		},
				// 		// label
				{
					Event:  org.LabelPolicyAddedEventType,
					Reduce: s.reduceLabelAdded,
				},
				{
					Event:  org.LabelPolicyChangedEventType,
					Reduce: s.reduceLabelChanged,
				},
				// {
				// 	Event:  org.LabelPolicyRemovedEventType,
				// 	Reduce: s.reduceLabelPolicyRemoved,
				// },
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
				// 		{
				// 			Event:  org.PasswordComplexityPolicyAddedEventType,
				// 			Reduce: s.reducePassedComplexityAdded,
				// 		},
				// 		{
				// 			Event:  org.PasswordComplexityPolicyChangedEventType,
				// 			Reduce: s.reducePasswordComplexityChanged,
				// 		},
				// 		{
				// 			Event:  org.PasswordComplexityPolicyRemovedEventType,
				// 			Reduce: s.reducePasswordComplexityRemoved,
				// 		},
				// 		// Password Policy
				// 		{
				// 			Event:  org.PasswordAgePolicyAddedEventType,
				// 			Reduce: s.reducePasswordPolicyAdded,
				// 		},
				// 		{
				// 			Event:  org.PasswordAgePolicyChangedEventType,
				// 			Reduce: s.reducePasswordPolicyChanged,
				// 		},
				// 		{
				// 			Event:  org.PasswordAgePolicyRemovedEventType,
				// 			Reduce: s.reducePasswordPolicyRemoved,
				// 		},
				// 		// Lockout Policy
				// 		{
				// 			Event:  org.LockoutPolicyAddedEventType,
				// 			Reduce: s.reduceLockoutPolicyAdded,
				// 		},
				// 		{
				// 			Event:  org.LockoutPolicyChangedEventType,
				// 			Reduce: s.reduceLockoutPolicyChanged,
				// 		},
				// 		{
				// 			Event:  org.LockoutPolicyRemovedEventType,
				// 			Reduce: s.reduceOrgLockoutPolicyRemoved,
				// 		},
				// 		// Domain Policy
				// 		{
				// 			Event:  org.DomainPolicyAddedEventType,
				// 			Reduce: s.reduceDomainPolicyAdded,
				// 		},
				// 		{
				// 			Event:  org.DomainPolicyChangedEventType,
				// 			Reduce: s.reduceDomainPolicyChanged,
				// 		},
				// 		{
				// 			Event:  org.DomainPolicyRemovedEventType,
				// 			Reduce: s.reduceOrgDomainPolicyRemoved,
				// 		},
				// 		// Delete org
				// 		{
				// 			Event:  org.OrgRemovedEventType,
				// 			Reduce: s.reduceOrgRemoved,
				// 		},
				// 	},
				// },
				// // settings
				// {
				// 	Aggregate: settings.AggregateType,
				// 	EventReducers: []handler.EventReducer{
				// 		{
				// 			Event:  settings.OrganizationSettingsSetEventType,
				// 			Reduce: s.reduceOrganizationSettingsSet,
				// 		},
				// 		{
				// 			Event:  settings.OrganizationSettingsRemovedEventType,
				// 			Reduce: s.reduceOrganizationSettingsRemoved,
				// 		},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				// 		// Login
				// 		{
				// 			Event:  instance.LoginPolicyAddedEventType,
				// 			Reduce: s.reduceLoginPolicyAdded,
				// 		},
				// 		{
				// 			Event:  instance.LoginPolicyChangedEventType,
				// 			Reduce: s.reduceLoginPolicyChanged,
				// 		},
				// 		{
				// 			Event:  instance.LoginPolicyMultiFactorAddedEventType,
				// 			Reduce: s.reduceMFAAdded,
				// 		},
				// 		{
				// 			Event:  instance.LoginPolicyMultiFactorRemovedEventType,
				// 			Reduce: s.reduceMFARemoved,
				// 		},
				// 		{
				// 			Event:  instance.LoginPolicySecondFactorAddedEventType,
				// 			Reduce: s.reduceSecondFactorAdded,
				// 		},
				// 		{
				// 			Event:  instance.LoginPolicySecondFactorRemovedEventType,
				// 			Reduce: s.reduceSecondFactorRemoved,
				// 		},
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
				// 		{
				// 			Event:  instance.PasswordComplexityPolicyAddedEventType,
				// 			Reduce: s.reducePassedComplexityAdded,
				// 		},
				// 		{
				// 			Event:  instance.PasswordComplexityPolicyChangedEventType,
				// 			Reduce: s.reducePasswordComplexityChanged,
				// 		},
				// 		// Password Policy
				// 		{
				// 			Event:  instance.PasswordAgePolicyAddedEventType,
				// 			Reduce: s.reducePasswordPolicyAdded,
				// 		},
				// 		{
				// 			Event:  instance.PasswordAgePolicyChangedEventType,
				// 			Reduce: s.reducePasswordPolicyChanged,
				// 		},
				// 		// Lockout Policy
				// 		{
				// 			Event:  instance.LockoutPolicyAddedEventType,
				// 			Reduce: s.reduceLockoutPolicyAdded,
				// 		},
				// 		{
				// 			Event:  instance.LockoutPolicyChangedEventType,
				// 			Reduce: s.reduceLockoutPolicyChanged,
				// 		},
				// 		// Domain Policy
				// 		{
				// 			Event:  instance.DomainPolicyAddedEventType,
				// 			Reduce: s.reduceDomainPolicyAdded,
				// 		},
				// 		{
				// 			Event:  instance.DomainPolicyChangedEventType,
				// 			Reduce: s.reduceDomainPolicyChanged,
				// 		},
				// 		// Security Policy
				// 		{
				// 			Event:  instance.SecurityPolicySetEventType,
				// 			Reduce: s.reduceSecurityPolicySet,
				// 		},
				// 		// Delete Instance
				// 		{
				// 			Event:  instance.InstanceRemovedEventType,
				// 			Reduce: s.reduceInstanceRemoved,
			},
		},
		// },
	}
}

// func (s *settingsRelationalProjection) reduceLoginPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.LoginPolicyAddedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *instance.LoginPolicyAddedEvent:
// 		policyEvent = e.LoginPolicyAddedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	case *org.LoginPolicyAddedEvent:
// 		policyEvent = e.LoginPolicyAddedEvent
// 		ownerType = domain.OwnerTypeOrganization
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YYPxS", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyAddedEventType, instance.LoginPolicyAddedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		settingJSON, err := json.Marshal(policyEvent)
// 		if err != nil {
// 			return err
// 		}

// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
// 		}
// 		settingsRepo := repository.SettingsRepository()
// 		setting := domain.Setting{
// 			InstanceID: policyEvent.Aggregate().InstanceID,
// 			OrgID:      orgId,
// 			Type:       domain.SettingTypeLogin,
// 			OwnerType:  ownerType,
// 			Settings:   settingJSON,
// 			CreatedAt:  policyEvent.CreationDate(),
// 			UpdatedAt:  &policyEvent.Creation,
// 		}
// 		err = settingsRepo.Create(ctx, v3_sql.SQLTx(tx), &setting)
// 		fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ADDDING setting.InstanceID = %+v\n", setting.InstanceID)
// 		return err
// 	}), nil
// }

// //nolint:gocognit
// func (s *settingsRelationalProjection) reduceLoginPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var ownerType domain.OwnerType
// 	var policyEvent policy.LoginPolicyChangedEvent
// 	switch e := event.(type) {
// 	case *instance.LoginPolicyChangedEvent:
// 		policyEvent = e.LoginPolicyChangedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	case *org.LoginPolicyChangedEvent:
// 		policyEvent = e.LoginPolicyChangedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 		ownerType = domain.OwnerTypeOrganization
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BHd86", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyChangedEventType, instance.LoginPolicyChangedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLk9y", "reduce.wrong.db.pool %T", ex)
// 		}

// 		loginRepo := repository.LoginRepository()

// 		setting := &domain.LoginSetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		if policyEvent.AllowRegister != nil {
// 			// setting.Settings.AllowRegister = *policyEvent.AllowRegister
// 			setting.Settings.AllowRegister = policyEvent.AllowRegister
// 		}
// 		if policyEvent.AllowUserNamePassword != nil {
// 			setting.Settings.AllowUserNamePassword = policyEvent.AllowUserNamePassword
// 		}
// 		if policyEvent.AllowExternalIDP != nil {
// 			setting.Settings.AllowExternalIDP = policyEvent.AllowExternalIDP
// 		}
// 		if policyEvent.ForceMFA != nil {
// 			// setting.Settings.ForceMFA = policyEvent.ForceMFA
// 			forceMFA := *policyEvent.ForceMFA
// 			setting.Settings.ForceMFA = &forceMFA
// 		}
// 		if policyEvent.ForceMFALocalOnly != nil {
// 			setting.Settings.ForceMFALocalOnly = policyEvent.ForceMFALocalOnly
// 		}
// 		if policyEvent.PasswordlessType != nil {
// 			passwordlessType := domain.PasswordlessType(*policyEvent.PasswordlessType)
// 			setting.Settings.PasswordlessType = &passwordlessType
// 		}
// 		if policyEvent.HidePasswordReset != nil {
// 			setting.Settings.HidePasswordReset = policyEvent.HidePasswordReset
// 		}
// 		if policyEvent.IgnoreUnknownUsernames != nil {
// 			setting.Settings.IgnoreUnknownUsernames = policyEvent.IgnoreUnknownUsernames
// 		}
// 		if policyEvent.AllowDomainDiscovery != nil {
// 			setting.Settings.AllowDomainDiscovery = policyEvent.AllowDomainDiscovery
// 		}
// 		if policyEvent.DisableLoginWithEmail != nil {
// 			setting.Settings.DisableLoginWithEmail = policyEvent.DisableLoginWithEmail
// 		}
// 		if policyEvent.DisableLoginWithPhone != nil {
// 			setting.Settings.DisableLoginWithPhone = policyEvent.DisableLoginWithPhone
// 		}
// 		if policyEvent.DefaultRedirectURI != nil {
// 			setting.Settings.DefaultRedirectURI = *policyEvent.DefaultRedirectURI
// 		}
// 		if policyEvent.PasswordCheckLifetime != nil {
// 			setting.Settings.PasswordCheckLifetime = *policyEvent.PasswordCheckLifetime
// 		}
// 		if policyEvent.ExternalLoginCheckLifetime != nil {
// 			setting.Settings.ExternalLoginCheckLifetime = *policyEvent.ExternalLoginCheckLifetime
// 		}
// 		if policyEvent.MFAInitSkipLifetime != nil {
// 			setting.Settings.MFAInitSkipLifetime = *policyEvent.MFAInitSkipLifetime
// 		}
// 		if policyEvent.SecondFactorCheckLifetime != nil {
// 			setting.Settings.SecondFactorCheckLifetime = *policyEvent.SecondFactorCheckLifetime
// 		}
// 		if policyEvent.MultiFactorCheckLifetime != nil {
// 			setting.Settings.MultiFactorCheckLifetime = *policyEvent.MultiFactorCheckLifetime
// 		}

// 		setting.UpdatedAt = &policyEvent.Creation

// 		err := loginRepo.Set(ctx, v3_sql.SQLTx(tx), setting, loginRepo.SetUpdatedAt(&policyEvent.Creation))

// 		fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> UPDATING setting.InstanceID = %+v\n", setting.InstanceID)
// 		fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> err = %+v\n", err)
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceMFAAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var ownerType domain.OwnerType
// 	var policyEvent policy.MultiFactorAddedEvent
// 	switch e := event.(type) {
// 	case *instance.LoginPolicyMultiFactorAddedEvent:
// 		policyEvent = e.MultiFactorAddedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	case *org.LoginPolicyMultiFactorAddedEvent:
// 		policyEvent = e.MultiFactorAddedEvent
// 		orgId = &policyEvent.Aggregate().ID
// 		ownerType = domain.OwnerTypeOrganization
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-WghuV", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyMultiFactorAddedEventType, instance.LoginPolicyMultiFactorAddedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLw7y", "reduce.wrong.db.pool %T", ex)
// 		}

// 		loginRepo := repository.LoginRepository()

// 		setting := &domain.LoginSetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		if slices.Contains(setting.Settings.MFAType, domain.MultiFactorType(policyEvent.MFAType)) {
// 			return nil
// 		}

// 		setting.Settings.MFAType = append(setting.Settings.MFAType, domain.MultiFactorType(policyEvent.MFAType))

// 		err := loginRepo.Set(ctx, v3_sql.SQLTx(tx), setting, loginRepo.SetUpdatedAt(&policyEvent.Creation))
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceMFARemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var ownerType domain.OwnerType
// 	var policyEvent policy.MultiFactorRemovedEvent
// 	switch e := event.(type) {
// 	case *instance.LoginPolicyMultiFactorRemovedEvent:
// 		policyEvent = e.MultiFactorRemovedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	case *org.LoginPolicyMultiFactorRemovedEvent:
// 		policyEvent = e.MultiFactorRemovedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 		ownerType = domain.OwnerTypeOrganization
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-cHU7u", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyMultiFactorRemovedEventType, instance.LoginPolicyMultiFactorRemovedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLi9y", "reduce.wrong.db.pool %T", ex)
// 		}

// 		loginRepo := repository.LoginRepository()

// 		setting := &domain.LoginSetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		setting.Settings.MFAType = slices.DeleteFunc(setting.Settings.MFAType, func(mfaType domain.MultiFactorType) bool {
// 			return mfaType == domain.MultiFactorType(policyEvent.MFAType)
// 		})

// 		setting.UpdatedAt = &policyEvent.Creation

// 		err := loginRepo.Set(ctx, v3_sql.SQLTx(tx), setting, loginRepo.SetUpdatedAt(&policyEvent.Creation))
// 		return err
// 	}), nil
// }

// func (*settingsRelationalProjection) reduceLoginPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	loginPolicyRemovedEvent, ok := event.(*org.LoginPolicyRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-oRSvD", "reduce.wrong.event.type %s", org.LoginPolicyRemovedEventType)
// 	}
// 	return handler.NewStatement(loginPolicyRemovedEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-arg9y", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		_, err := settingsRepo.Delete(
// 			ctx, v3_sql.SQLTx(tx),
// 			loginPolicyRemovedEvent.Aggregate().InstanceID,
// 			&loginPolicyRemovedEvent.Aggregate().ID,
// 			domain.SettingTypeLogin)
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceSecondFactorAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var ownerType domain.OwnerType
// 	var policyEvent policy.SecondFactorAddedEvent
// 	switch e := event.(type) {
// 	case *instance.LoginPolicySecondFactorAddedEvent:
// 		policyEvent = e.SecondFactorAddedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	case *org.LoginPolicySecondFactorAddedEvent:
// 		policyEvent = e.SecondFactorAddedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 		ownerType = domain.OwnerTypeOrganization
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-apB2E", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicySecondFactorAddedEventType, instance.LoginPolicySecondFactorAddedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iLk4m", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		// setting, err := settingsRepo.GetLogin(ctx, v3_sql.SQLTx(tx), policyEvent.Agg.InstanceID, orgId)
// 		// if err != nil {
// 		// 	return zerrors.ThrowInternal(err, "HANDL-H7m9m", "error accessing login policy record")
// 		// }
// 		loginRepo := repository.LoginRepository()

// 		setting := &domain.LoginSetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		if slices.Contains(setting.Settings.SecondFactorTypes, domain.SecondFactorType(policyEvent.MFAType)) {
// 			return nil
// 		}

// 		setting.UpdatedAt = &policyEvent.Creation

// 		setting.Settings.SecondFactorTypes = append(setting.Settings.SecondFactorTypes, domain.SecondFactorType(policyEvent.MFAType))

// 		err := loginRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&policyEvent.Creation))
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceSecondFactorRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var ownerType domain.OwnerType
// 	var policyEvent policy.SecondFactorRemovedEvent
// 	switch e := event.(type) {
// 	case *instance.LoginPolicySecondFactorRemovedEvent:
// 		policyEvent = e.SecondFactorRemovedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	case *org.LoginPolicySecondFactorRemovedEvent:
// 		policyEvent = e.SecondFactorRemovedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 		ownerType = domain.OwnerTypeOrganization
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-bYpmA", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicySecondFactorRemovedEventType, instance.LoginPolicySecondFactorRemovedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rnd0y", "reduce.wrong.db.pool %T", ex)
// 		}

// 		// settingsRepo := repository.SettingsRepository()

// 		// setting, err := settingsRepo.GetLogin(ctx, v3_sql.SQLTx(tx), policyEvent.Agg.InstanceID, orgId)
// 		// if err != nil {
// 		// 	return zerrors.ThrowInternal(err, "HANDL-rsk9m", "error accessing login policy record")
// 		// }
// 		loginRepo := repository.LoginRepository()

// 		setting := &domain.LoginSetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		setting.Settings.SecondFactorTypes = slices.DeleteFunc(setting.Settings.SecondFactorTypes, func(secondFactorType domain.SecondFactorType) bool {
// 			return secondFactorType == domain.SecondFactorType(policyEvent.MFAType)
// 		})

// 		// _, err = settingsRepo.UpdateLogin(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&policyEvent.Creation))
// 		err := loginRepo.Set(ctx, v3_sql.SQLTx(tx), setting, loginRepo.SetUpdatedAt(&policyEvent.Creation))
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	removeInstanceEvent, ok := event.(*instance.InstanceRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-2ZUFS", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
// 	}

// 	return handler.NewStatement(removeInstanceEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rrdHy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		_, err := settingsRepo.DeleteSettingsForInstance(ctx, v3_sql.SQLTx(tx), removeInstanceEvent.Aggregate().InstanceID)
// 		return err
// 	}), nil
// }

// // label
func (s *settingsRelationalProjection) reduceLabelAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyAddedEvent
	var ownerType domain.OwnerType
	switch e := event.(type) {
	case *org.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
		ownerType = domain.OwnerTypeOrganization
		orgId = &policyEvent.Aggregate().ResourceOwner
	case *instance.LabelPolicyAddedEvent:
		policyEvent = e.LabelPolicyAddedEvent
		ownerType = domain.OwnerTypeInstance
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CSE7A", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyAddedEventType, instance.LabelPolicyAddedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		settings, err := json.Marshal(policyEvent)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		labelStatePreview := domain.LabelStatePreview
		settingsRepo := repository.SettingsRepository()
		setting := domain.Setting{
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			OwnerType:  ownerType,
			Type:       domain.SettingTypeLabel,
			LabelState: &labelStatePreview,
			Settings:   settings,
			CreatedAt:  policyEvent.CreationDate(),
			UpdatedAt:  &policyEvent.Creation,
		}
		err = settingsRepo.Create(ctx, v3_sql.SQLTx(tx), &setting)
		fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> err = %+v\n", err)
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceLabelChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var policyEvent policy.LabelPolicyChangedEvent
	var ownerType domain.OwnerType
	switch e := event.(type) {
	case *org.LabelPolicyChangedEvent:
		policyEvent = e.LabelPolicyChangedEvent
		orgId = &policyEvent.Aggregate().ResourceOwner
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyChangedEvent:
		policyEvent = e.LabelPolicyChangedEvent
		ownerType = domain.OwnerTypeInstance
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyChangedEventType, instance.LabelPolicyChangedEventType})
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-lhb9y", "reduce.wrong.db.pool %T", ex)
		}

		// settingsRepo := repository.SettingsRepository()

		// setting, err := settingsRepo.GetLabel(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, domain.LabelStatePreview)
		// if err != nil {
		// 	return zerrors.ThrowInternal(err, "HANDL-r879m", "error accessing login policy record")
		// }

		settingsRepo := repository.LabelRepository()

		change := make([]db_json.JSONFieldChange, 0, 9)

		if policyEvent.PrimaryColor != nil {
			change = append(change, settingsRepo.SetPrimaryColorField(*policyEvent.PrimaryColor))
		}
		if policyEvent.BackgroundColor != nil {
			change = append(change, settingsRepo.SetBackgroundColorField(*policyEvent.BackgroundColor))
		}
		if policyEvent.WarnColor != nil {
			change = append(change, settingsRepo.SetWarnColorField(*policyEvent.WarnColor))
		}
		if policyEvent.FontColor != nil {
			change = append(change, settingsRepo.SetFontColorField(*policyEvent.FontColor))
		}
		if policyEvent.PrimaryColorDark != nil {
			change = append(change, settingsRepo.SetPrimaryCcolorDarkField(*policyEvent.PrimaryColorDark))
		}
		if policyEvent.BackgroundColorDark != nil {
			change = append(change, settingsRepo.SetBackgroundColorDarkField(*policyEvent.BackgroundColorDark))
		}
		if policyEvent.WarnColorDark != nil {
			change = append(change, settingsRepo.SetWarnColorDarkField(*policyEvent.WarnColorDark))
		}
		if policyEvent.FontColorDark != nil {
			change = append(change, settingsRepo.SetFontColorDarkField(*policyEvent.FontColorDark))
		}
		if policyEvent.HideLoginNameSuffix != nil {
			change = append(change, settingsRepo.SetHideLoginNameSuffixField(*policyEvent.HideLoginNameSuffix))
		}
		if policyEvent.ErrorMsgPopup != nil {
			change = append(change, settingsRepo.SetErrorMsgPopupField(*policyEvent.ErrorMsgPopup))
		}
		if policyEvent.DisableWatermark != nil {
			change = append(change, settingsRepo.SetDisableWatermarkField(*policyEvent.DisableWatermark))
		}
		if policyEvent.ThemeMode != nil {
			change = append(change, settingsRepo.SetThemeModeField(domain.LabelPolicyThemeMode(*policyEvent.ThemeMode)))
		}

		_, err := settingsRepo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(policyEvent.Agg.InstanceID),
				settingsRepo.OrgIDCondition(orgId),
				settingsRepo.TypeCondition(domain.SettingTypeLabel),
				settingsRepo.OwnerTypeCondition(ownerType),
				settingsRepo.LabelStateCondition(domain.LabelStatePreview),
			),
			settingsRepo.SetLabelSettings(change...),
			settingsRepo.SetUpdatedAt(&policyEvent.Creation))

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
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-r7k0y", "reduce.wrong.db.pool %T", ex)
		}

		settingsRepo := repository.SettingsRepository()

		orgId := &policyEvent.Aggregate().ID

		// _, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx), policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypeLabel)
		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(policyEvent.Agg.InstanceID),
				settingsRepo.OrgIDCondition(orgId),
				settingsRepo.TypeCondition(domain.SettingTypeLabel),
				settingsRepo.OwnerTypeCondition(domain.OwnerTypeOrganization),
				settingsRepo.LabelStateCondition(domain.LabelStatePreview),
			))
		return err
	}), nil
}

func (s *settingsRelationalProjection) reduceLabelActivated(event eventstore.Event) (*handler.Statement, error) {
	var orgCond handler.NamespacedCondition
	switch event.(type) {
	case *org.LabelPolicyActivatedEvent:
		orgId := &event.Aggregate().ID
		orgCond = handler.NewNamespacedCondition(SettingsOrgIDCol, orgId)
	case *instance.LabelPolicyActivatedEvent:
		orgCond = handler.NewIsNotNulNSlCond(SettingsOrgIDCol)
		// everything ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-7Kd8U", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyActivatedEventType, instance.LabelPolicyActivatedEventType})
	}

	return handler.NewCopyStatement(
		event,
		[]handler.Column{
			handler.NewCol(SettingInstanceIDCol, nil),
			handler.NewCol(SettingsOrgIDCol, nil),
			handler.NewCol(SettingsTypeCol, nil),
			handler.NewCol(SettingsLabelStateCol, nil),
		},
		[]handler.Condition{
			handler.NewCond(SettingsTypeCol, domain.SettingTypeLabel),
		},
		[]handler.Column{
			handler.NewCol(SettingInstanceIDCol, nil),
			handler.NewCol(SettingsOrgIDCol, nil),
			handler.NewCol(SettingsTypeCol, nil),
			handler.NewCol(SettingsLabelStateCol, domain.LabelStateActivated),
			handler.NewCol(SettingsSettingsCol, nil),
			handler.NewCol(UpdatedAt, event.CreatedAt()),
			handler.NewCol(CreatedAt, event.CreatedAt()),
		},
		[]handler.Column{
			handler.NewCol(SettingInstanceIDCol, nil),
			handler.NewCol(SettingsOrgIDCol, nil),
			handler.NewCol(SettingsTypeCol, nil),
			handler.NewCol(SettingsLabelStateCol, nil),
			handler.NewCol(SettingsSettingsCol, nil),
			handler.NewCol(UpdatedAt, nil),
			handler.NewCol(CreatedAt, nil),
		},
		[]handler.NamespacedCondition{
			handler.NewNamespacedCondition(SettingsTypeCol, domain.SettingTypeLabel),
			handler.NewNamespacedCondition(SettingInstanceIDCol, event.Aggregate().InstanceID),
			orgCond,
			handler.NewNamespacedCondition(SettingsLabelStateCol, domain.LabelStatePreview),
		}), nil
}

func (p *settingsRelationalProjection) reduceLabelLogoAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var ownerType domain.OwnerType
	switch e := event.(type) {
	case *org.LabelPolicyLogoAddedEvent:
		orgId = &e.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyLogoAddedEvent:
		ownerType = domain.OwnerTypeInstance
	case *org.LabelPolicyLogoDarkAddedEvent:
		orgId = &e.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyLogoDarkAddedEvent:
		ownerType = domain.OwnerTypeInstance
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-4UbiP", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoAddedEventType, instance.LabelPolicyLogoAddedEventType, org.LabelPolicyLogoDarkAddedEventType, instance.LabelPolicyLogoDarkAddedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		// settingsRepo := repository.SettingsRepository()

		// setting, err := settingsRepo.GetLabel(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, domain.LabelStatePreview)
		// if err != nil {
		// 	return zerrors.ThrowInternal(err, "HANDL-y7dDm", "error accessing login policy record")
		// }

		settingsRepo := repository.LabelRepository()

		// setting := &domain.LabelSetting{
		// 	Setting: &domain.Setting{
		// 		InstanceID: event.Aggregate().InstanceID,
		// 		OrgID:      orgId,
		// 		OwnerType:  ownerType,
		// 	},
		// }

		var change db_json.JSONFieldChange

		switch e := event.(type) {
		case *org.LabelPolicyLogoAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			change = settingsRepo.SetLabelPolicyLightLogoURL(url)
		case *instance.LabelPolicyLogoAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			change = settingsRepo.SetLabelPolicyLightLogoURL(url)
		case *org.LabelPolicyLogoDarkAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			change = settingsRepo.SetLabelPolicyDarkLogoURL(url)
		case *instance.LabelPolicyLogoDarkAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			change = settingsRepo.SetLabelPolicyDarkLogoURL(url)
		}

		CreatedAt := event.CreatedAt()

		// _, err := settingsRepo.Update(ctx, v3_sql.SQLTx(tx),
		u, err := settingsRepo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(event.Aggregate().InstanceID),
				settingsRepo.OrgIDCondition(orgId),
				settingsRepo.TypeCondition(domain.SettingTypeLabel),
				settingsRepo.OwnerTypeCondition(ownerType),
				settingsRepo.LabelStateCondition(domain.LabelStatePreview),
			),
			settingsRepo.SetLabelSettings(change),
			settingsRepo.SetUpdatedAt(&CreatedAt))
		fmt.Printf("[DEBUGPRINT] [settings_relational.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> u = %+v\n", u)
		fmt.Printf("[DEBUGPRINT] [settings_relational.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> LOGO err = %+v\n", err)

		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceLogoRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var ownerType domain.OwnerType
	switch event.(type) {
	case *org.LabelPolicyLogoRemovedEvent:
		orgId = &event.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyLogoRemovedEvent:
		ownerType = domain.OwnerTypeInstance
	case *org.LabelPolicyLogoDarkRemovedEvent:
		orgId = &event.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyLogoDarkRemovedEvent:
		ownerType = domain.OwnerTypeInstance
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-kg8H4", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyLogoRemovedEventType, instance.LabelPolicyLogoRemovedEventType, org.LabelPolicyLogoDarkRemovedEventType, instance.LabelPolicyLogoDarkRemovedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		// settingsRepo := repository.SettingsRepository()

		// setting, err := settingsRepo.GetLabel(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, domain.LabelStatePreview)
		// if err != nil {
		// 	return zerrors.ThrowInternal(err, "HANDL-d7L9s", "error accessing login policy record")
		// }
		settingsRepo := repository.LabelRepository()

		// setting := &domain.LabelSetting{
		// 	Setting: &domain.Setting{
		// 		InstanceID: event.Aggregate().InstanceID,
		// 		OrgID:      orgId,
		// 		OwnerType:  ownerType,
		// 	},
		// }

		var change db_json.JSONFieldChange

		switch event.(type) {
		case *org.LabelPolicyLogoRemovedEvent:
			change = settingsRepo.SetLabelPolicyLightLogoURL(nil)
		case *instance.LabelPolicyLogoRemovedEvent:
			change = settingsRepo.SetLabelPolicyLightLogoURL(nil)
		case *org.LabelPolicyLogoDarkRemovedEvent:
			change = settingsRepo.SetLabelPolicyDarkLogoURL(nil)
		case *instance.LabelPolicyLogoDarkRemovedEvent:
			change = settingsRepo.SetLabelPolicyDarkLogoURL(nil)
		}

		CreatedAt := event.CreatedAt()

		// _, err = settingsRepo.UpdateLabel(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		// err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		// return err

		_, err := settingsRepo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(event.Aggregate().InstanceID),
				settingsRepo.OrgIDCondition(orgId),
				settingsRepo.TypeCondition(domain.SettingTypeLabel),
				settingsRepo.OwnerTypeCondition(ownerType),
				settingsRepo.LabelStateCondition(domain.LabelStatePreview),
			),
			settingsRepo.SetLabelSettings(change),
			settingsRepo.SetUpdatedAt(&CreatedAt))
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceIconAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var ownerType domain.OwnerType
	switch event.(type) {
	case *org.LabelPolicyIconAddedEvent:
		orgId = &event.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyIconAddedEvent:
		ownerType = domain.OwnerTypeInstance
	case *org.LabelPolicyIconDarkAddedEvent:
		orgId = &event.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyIconDarkAddedEvent:
		ownerType = domain.OwnerTypeInstance
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-e2JFz", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconAddedEventType, instance.LabelPolicyIconAddedEventType, org.LabelPolicyIconDarkAddedEventType, instance.LabelPolicyIconDarkAddedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		// settingsRepo := repository.SettingsRepository()

		// setting, err := settingsRepo.GetLabel(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, domain.LabelStatePreview)
		// if err != nil {
		// 	return zerrors.ThrowInternal(err, "HANDL-s7a9m", "error accessing login policy record")
		// }

		settingsRepo := repository.LabelRepository()

		// setting := &domain.LabelSetting{
		// 	Setting: &domain.Setting{
		// 		InstanceID: event.Aggregate().InstanceID,
		// 		OrgID:      orgId,
		// 		OwnerType:  ownerType,
		// 	},
		// }

		var change db_json.JSONFieldChange

		switch e := event.(type) {
		case *org.LabelPolicyIconAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			// setting.Settings.LabelPolicyLightIconURL = url
			change = settingsRepo.SetLabelPolicyLightIconURL(url)
		case *instance.LabelPolicyIconAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			// setting.Settings.LabelPolicyLightIconURL = url
			// setting.Settings.LabelPolicyLightIconURL = url
			change = settingsRepo.SetLabelPolicyLightIconURL(url)
		case *org.LabelPolicyIconDarkAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			// setting.Settings.LabelPolicyDarkIconURL = url
			// setting.Settings.LabelPolicyLightIconURL = url
			change = settingsRepo.SetLabelPolicyDarkIconURL(url)
		case *instance.LabelPolicyIconDarkAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			// setting.Settings.LabelPolicyDarkIconURL = url
			change = settingsRepo.SetLabelPolicyDarkIconURL(url)
		}

		CreatedAt := event.CreatedAt()

		// _, err = settingsRepo.UpdateLabel(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		// err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		_, err := settingsRepo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(event.Aggregate().InstanceID),
				settingsRepo.OrgIDCondition(orgId),
				settingsRepo.TypeCondition(domain.SettingTypeLabel),
				settingsRepo.OwnerTypeCondition(ownerType),
				settingsRepo.LabelStateCondition(domain.LabelStatePreview),
			),
			settingsRepo.SetLabelSettings(change),
			settingsRepo.SetUpdatedAt(&CreatedAt))
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceIconRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var ownerType domain.OwnerType
	switch event.(type) {
	case *org.LabelPolicyIconRemovedEvent:
		orgId = &event.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyIconRemovedEvent:
		ownerType = domain.OwnerTypeInstance
	case *org.LabelPolicyIconDarkRemovedEvent:
		orgId = &event.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyIconDarkRemovedEvent:
		ownerType = domain.OwnerTypeInstance
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gfgbY", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyIconRemovedEventType, instance.LabelPolicyIconRemovedEventType, org.LabelPolicyIconDarkRemovedEventType, instance.LabelPolicyIconDarkRemovedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		// settingsRepo := repository.SettingsRepository()

		// setting, err := settingsRepo.GetLabel(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, domain.LabelStatePreview)
		// if err != nil {
		// 	return zerrors.ThrowInternal(err, "HANDL-B7L9m", "error accessing login policy record")
		// }

		settingsRepo := repository.LabelRepository()

		// setting := &domain.LabelSetting{
		// 	Setting: &domain.Setting{
		// 		InstanceID: event.Aggregate().InstanceID,
		// 		OrgID:      orgId,
		// 		OwnerType:  ownerType,
		// 	},
		// }

		var change db_json.JSONFieldChange

		switch event.(type) {
		case *org.LabelPolicyIconRemovedEvent:
			// setting.Settings.LabelPolicyLightIconURL = nil
			change = settingsRepo.SetLabelPolicyLightIconURL(nil)
		case *instance.LabelPolicyIconRemovedEvent:
			change = settingsRepo.SetLabelPolicyLightIconURL(nil)
		case *org.LabelPolicyIconDarkRemovedEvent:
			change = settingsRepo.SetLabelPolicyDarkIconURL(nil)
		case *instance.LabelPolicyIconDarkRemovedEvent:
			change = settingsRepo.SetLabelPolicyDarkIconURL(nil)
		}

		CreatedAt := event.CreatedAt()

		// _, err = settingsRepo.UpdateLabel(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		// err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		_, err := settingsRepo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(event.Aggregate().InstanceID),
				settingsRepo.OrgIDCondition(orgId),
				settingsRepo.TypeCondition(domain.SettingTypeLabel),
				settingsRepo.OwnerTypeCondition(ownerType),
				settingsRepo.LabelStateCondition(domain.LabelStatePreview),
			),
			settingsRepo.SetLabelSettings(change),
			settingsRepo.SetUpdatedAt(&CreatedAt))
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceFontAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var ownerType domain.OwnerType
	switch event.(type) {
	case *org.LabelPolicyFontAddedEvent:
		orgId = &event.Aggregate().ID
		ownerType = domain.OwnerTypeOrganization
	case *instance.LabelPolicyFontAddedEvent:
		ownerType = domain.OwnerTypeInstance
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-65i9W", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontAddedEventType, instance.LabelPolicyFontAddedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		// settingsRepo := repository.SettingsRepository()

		// setting, err := settingsRepo.GetLabel(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, domain.LabelStatePreview)
		// if err != nil {
		// 	return zerrors.ThrowInternal(err, "HANDL-H7S7m", "error accessing login policy record")
		// }

		settingsRepo := repository.LabelRepository()

		// setting := &domain.LabelSetting{
		// 	Setting: &domain.Setting{
		// 		InstanceID: event.Aggregate().InstanceID,
		// 		OrgID:      orgId,
		// 		OwnerType:  ownerType,
		// 	},
		// }

		var change db_json.JSONFieldChange
		switch e := event.(type) {
		case *org.LabelPolicyFontAddedEvent:
			// setting.Settings.LabelPolicyFontURL = &e.StoreKey
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			change = settingsRepo.SetLabelPolicyFontURL(url)
		case *instance.LabelPolicyFontAddedEvent:
			url, err := url.Parse(e.StoreKey)
			if err != nil {
				return err
			}
			change = settingsRepo.SetLabelPolicyFontURL(url)
			// setting.Settings.LabelPolicyFontURL = &e.StoreKey
		}

		CreatedAt := event.CreatedAt()

		// _, err = settingsRepo.UpdateLabel(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		// err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		_, err := settingsRepo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(event.Aggregate().InstanceID),
				settingsRepo.OrgIDCondition(orgId),
				settingsRepo.TypeCondition(domain.SettingTypeLabel),
				settingsRepo.OwnerTypeCondition(ownerType),
				settingsRepo.LabelStateCondition(domain.LabelStatePreview),
			),
			settingsRepo.SetLabelSettings(change),
			settingsRepo.SetUpdatedAt(&CreatedAt))
		return err
	}), nil
}

func (p *settingsRelationalProjection) reduceFontRemoved(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var ownerType domain.OwnerType
	switch event.(type) {
	case *org.LabelPolicyFontRemovedEvent:
		orgId = &event.Aggregate().ID
		ownerType = domain.OwnerTypeInstance
	case *instance.LabelPolicyFontRemovedEvent:
		ownerType = domain.OwnerTypeInstance
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-xf32J", "reduce.wrong.event.type %v", []eventstore.EventType{org.LabelPolicyFontRemovedEventType, instance.LabelPolicyFontRemovedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
		}
		// settingsRepo := repository.SettingsRepository()

		// setting, err := settingsRepo.GetLabel(ctx, v3_sql.SQLTx(tx), event.Aggregate().InstanceID, orgId, domain.LabelStatePreview)
		// if err != nil {
		// 	return zerrors.ThrowInternal(err, "HANDL-77kMm", "error accessing login policy record")
		// }
		settingsRepo := repository.LabelRepository()

		// setting := &domain.LabelSetting{
		// 	Setting: &domain.Setting{
		// 		InstanceID: event.Aggregate().InstanceID,
		// 		OrgID:      orgId,
		// 		OwnerType:  ownerType,
		// 	},
		// }

		// setting.Settings.LabelPolicyFontURL = nil
		change := settingsRepo.SetLabelPolicyFontURL(nil)

		CreatedAt := event.CreatedAt()

		// _, err = settingsRepo.UpdateLabel(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		// err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
		_, err := settingsRepo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				settingsRepo.InstanceIDCondition(event.Aggregate().InstanceID),
				settingsRepo.OrgIDCondition(orgId),
				settingsRepo.TypeCondition(domain.SettingTypeLabel),
				settingsRepo.OwnerTypeCondition(ownerType),
				settingsRepo.LabelStateCondition(domain.LabelStatePreview),
			),
			settingsRepo.SetLabelSettings(change),
			settingsRepo.SetUpdatedAt(&CreatedAt))
		return err
	}), nil
}

// func (p *settingsRelationalProjection) reducePassedComplexityAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.PasswordComplexityPolicyAddedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *org.PasswordComplexityPolicyAddedEvent:
// 		policyEvent = e.PasswordComplexityPolicyAddedEvent
// 		ownerType = domain.OwnerTypeOrganization
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 	case *instance.PasswordComplexityPolicyAddedEvent:
// 		policyEvent = e.PasswordComplexityPolicyAddedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-KTHmJ", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordComplexityPolicyAddedEventType, instance.PasswordComplexityPolicyAddedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		settingJSON, err := json.Marshal(policyEvent)
// 		if err != nil {
// 			return err
// 		}

// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
// 		}
// 		settingsRepo := repository.SettingsRepository()
// 		newSetting := domain.Setting{
// 			InstanceID: policyEvent.Aggregate().InstanceID,
// 			OrgID:      orgId,
// 			Type:       domain.SettingTypePasswordComplexity,
// 			OwnerType:  ownerType,
// 			Settings:   settingJSON,
// 			CreatedAt:  policyEvent.CreationDate(),
// 			UpdatedAt:  &policyEvent.Creation,
// 		}
// 		err = settingsRepo.Create(ctx, v3_sql.SQLTx(tx), &newSetting)
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reducePasswordComplexityChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.PasswordComplexityPolicyChangedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *org.PasswordComplexityPolicyChangedEvent:
// 		policyEvent = e.PasswordComplexityPolicyChangedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 		ownerType = domain.OwnerTypeOrganization
// 	case *instance.PasswordComplexityPolicyChangedEvent:
// 		policyEvent = e.PasswordComplexityPolicyChangedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-cf3Xb", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordComplexityPolicyChangedEventType, instance.PasswordComplexityPolicyChangedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rLrfy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.PasswordComplexityRepository()

// 		setting := &domain.PasswordComplexitySetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		if policyEvent.MinLength != nil {
// 			setting.Settings.MinLength = *policyEvent.MinLength
// 		}
// 		if policyEvent.HasLowercase != nil {
// 			setting.Settings.HasLowercase = *policyEvent.HasLowercase
// 		}
// 		if policyEvent.HasUppercase != nil {
// 			setting.Settings.HasUppercase = *policyEvent.HasUppercase
// 		}
// 		if policyEvent.HasSymbol != nil {
// 			setting.Settings.HasSymbol = *policyEvent.HasSymbol
// 		}
// 		if policyEvent.HasNumber != nil {
// 			setting.Settings.HasNumber = *policyEvent.HasNumber
// 		}

// 		CreatedAt := event.CreatedAt()

// 		err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&CreatedAt))
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reducePasswordComplexityRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	policyEvent, ok := event.(*org.PasswordComplexityPolicyRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-wttCd", "reduce.wrong.event.type %s", org.PasswordComplexityPolicyRemovedEventType)
// 	}

// 	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		orgID := &policyEvent.Aggregate().ID

// 		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx), policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypePasswordComplexity)
// 		return err
// 	}), nil
// }

// func (p *settingsRelationalProjection) reducePasswordPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.PasswordAgePolicyAddedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *org.PasswordAgePolicyAddedEvent:
// 		policyEvent = e.PasswordAgePolicyAddedEvent
// 		ownerType = domain.OwnerTypeOrganization
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 	case *instance.PasswordAgePolicyAddedEvent:
// 		policyEvent = e.PasswordAgePolicyAddedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CJqF0", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyAddedEventType, instance.PasswordAgePolicyAddedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		settings, err := json.Marshal(policyEvent)
// 		if err != nil {
// 			return err
// 		}

// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hONE", "reduce.wrong.db.pool %T", ex)
// 		}
// 		settingsRepo := repository.SettingsRepository()
// 		setting := domain.Setting{
// 			InstanceID: policyEvent.Aggregate().InstanceID,
// 			OrgID:      orgId,
// 			Type:       domain.SettingTypePasswordExpiry,
// 			OwnerType:  ownerType,
// 			Settings:   settings,
// 			CreatedAt:  policyEvent.CreationDate(),
// 			UpdatedAt:  &policyEvent.Creation,
// 		}
// 		err = settingsRepo.Create(ctx, v3_sql.SQLTx(tx), &setting)
// 		return err
// 	}), nil
// }

// func (p *settingsRelationalProjection) reducePasswordPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.PasswordAgePolicyChangedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *org.PasswordAgePolicyChangedEvent:
// 		policyEvent = e.PasswordAgePolicyChangedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 		ownerType = domain.OwnerTypeOrganization
// 	case *instance.PasswordAgePolicyChangedEvent:
// 		policyEvent = e.PasswordAgePolicyChangedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-i7FZt", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyChangedEventType, instance.PasswordAgePolicyChangedEventType})
// 	}
// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-Mlk6y", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.PasswordExpiryRepository()

// 		// setting, err := settingsRepo.GetPasswordExpiry(ctx, v3_sql.SQLTx(tx), policyEvent.Agg.InstanceID, orgId)
// 		// if err != nil {
// 		// 	return zerrors.ThrowInternal(err, "HANDL-z7k3m", "error accessing login policy record")
// 		// }

// 		setting := &domain.PasswordExpirySetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		if policyEvent.ExpireWarnDays != nil {
// 			setting.Settings.ExpireWarnDays = *policyEvent.ExpireWarnDays
// 		}
// 		if policyEvent.MaxAgeDays != nil {
// 			setting.Settings.MaxAgeDays = *policyEvent.MaxAgeDays
// 		}

// 		err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&policyEvent.Creation))
// 		return err
// 	}), nil
// }

// func (p *settingsRelationalProjection) reducePasswordPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	policyEvent, ok := event.(*org.PasswordAgePolicyRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-EtHWB", "reduce.wrong.event.type %s", org.PasswordAgePolicyRemovedEventType)
// 	}
// 	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		orgID := &policyEvent.Aggregate().ID

// 		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx), policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypePasswordExpiry)
// 		return err
// 	}), nil
// }

// func (p *settingsRelationalProjection) reduceOrgLockoutPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	policyEvent, ok := event.(*org.LockoutPolicyRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.LockoutPolicyRemovedEventType)
// 	}
// 	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		orgID := &policyEvent.Aggregate().ID

// 		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx), policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypeLockout)
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*org.OrgRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-IoW0x", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
// 	}

// 	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		orgID := e.Aggregate().ID

// 		_, err := settingsRepo.DeleteSettingsForOrg(ctx, v3_sql.SQLTx(tx), orgID)
// 		return err
// 	}), nil
// }

// func (p *settingsRelationalProjection) reduceLockoutPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.LockoutPolicyAddedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *org.LockoutPolicyAddedEvent:
// 		policyEvent = e.LockoutPolicyAddedEvent
// 		ownerType = domain.OwnerTypeOrganization
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 	case *instance.LockoutPolicyAddedEvent:
// 		policyEvent = e.LockoutPolicyAddedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-d8mZO", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyAddedEventType, instance.LockoutPolicyAddedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		settings, err := json.Marshal(policyEvent)
// 		if err != nil {
// 			return err
// 		}

// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5hnNE", "reduce.wrong.db.pool %T", ex)
// 		}
// 		settingsRepo := repository.SettingsRepository()
// 		setting := domain.Setting{
// 			InstanceID: policyEvent.Aggregate().InstanceID,
// 			OrgID:      orgId,
// 			Type:       domain.SettingTypeLockout,
// 			OwnerType:  ownerType,
// 			Settings:   settings,
// 			CreatedAt:  policyEvent.CreationDate(),
// 			UpdatedAt:  &policyEvent.Creation,
// 		}
// 		err = settingsRepo.Create(ctx, v3_sql.SQLTx(tx), &setting)
// 		return err
// 	}), nil
// }

// func (p *settingsRelationalProjection) reduceLockoutPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.LockoutPolicyChangedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *org.LockoutPolicyChangedEvent:
// 		policyEvent = e.LockoutPolicyChangedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 		ownerType = domain.OwnerTypeOrganization
// 	case *instance.LockoutPolicyChangedEvent:
// 		policyEvent = e.LockoutPolicyChangedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gT3BQ", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyChangedEventType, instance.LockoutPolicyChangedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rbsxy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.LockoutRepository()

// 		setting := &domain.LockoutSetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		// setting, err := settingsRepo.GetLockout(ctx, v3_sql.SQLTx(tx), policyEvent.Agg.InstanceID, orgId)
// 		// if err != nil {
// 		// 	return zerrors.ThrowInternal(err, "HANDL-rPkxm", "error accessing login policy record")
// 		// }

// 		if policyEvent.MaxPasswordAttempts != nil {
// 			setting.Settings.MaxPasswordAttempts = *policyEvent.MaxPasswordAttempts
// 		}
// 		if policyEvent.MaxOTPAttempts != nil {
// 			setting.Settings.MaxOTPAttempts = *policyEvent.MaxOTPAttempts
// 		}
// 		if policyEvent.ShowLockOutFailures != nil {
// 			setting.Settings.ShowLockOutFailures = *policyEvent.ShowLockOutFailures
// 		}

// 		err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&policyEvent.Creation))
// 		return err
// 	}), nil
// }

// func (p *settingsRelationalProjection) reduceDomainPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.DomainPolicyAddedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *org.DomainPolicyAddedEvent:
// 		policyEvent = e.DomainPolicyAddedEvent
// 		ownerType = domain.OwnerTypeOrganization
// 		ownerType = domain.OwnerTypeOrganization
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 	case *instance.DomainPolicyAddedEvent:
// 		policyEvent = e.DomainPolicyAddedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-8se7M", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyAddedEventType, instance.DomainPolicyAddedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		settingJSON, err := json.Marshal(policyEvent)
// 		if err != nil {
// 			return err
// 		}

// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-chduE", "reduce.wrong.db.pool %T", ex)
// 		}
// 		settingsRepo := repository.SettingsRepository()
// 		setting := domain.Setting{
// 			InstanceID: policyEvent.Aggregate().InstanceID,
// 			OrgID:      orgId,
// 			Type:       domain.SettingTypeDomain,
// 			OwnerType:  ownerType,
// 			Settings:   settingJSON,
// 			CreatedAt:  policyEvent.CreationDate(),
// 			UpdatedAt:  &policyEvent.Creation,
// 		}
// 		err = settingsRepo.Create(ctx, v3_sql.SQLTx(tx), &setting)
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceDomainPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
// 	var orgId *string
// 	var policyEvent policy.DomainPolicyChangedEvent
// 	var ownerType domain.OwnerType
// 	switch e := event.(type) {
// 	case *org.DomainPolicyChangedEvent:
// 		policyEvent = e.DomainPolicyChangedEvent
// 		orgId = &policyEvent.Aggregate().ResourceOwner
// 		ownerType = domain.OwnerTypeOrganization
// 	case *instance.DomainPolicyChangedEvent:
// 		policyEvent = e.DomainPolicyChangedEvent
// 		ownerType = domain.OwnerTypeInstance
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyChangedEventType, instance.DomainPolicyChangedEventType})
// 	}

// 	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rbsxy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.DomainRepository()

// 		setting := &domain.DomainSetting{
// 			Setting: &domain.Setting{
// 				InstanceID: policyEvent.Agg.InstanceID,
// 				OrgID:      orgId,
// 				OwnerType:  ownerType,
// 			},
// 		}

// 		if policyEvent.UserLoginMustBeDomain != nil {
// 			setting.Settings.UserLoginMustBeDomain = *policyEvent.UserLoginMustBeDomain
// 		}
// 		if policyEvent.ValidateOrgDomains != nil {
// 			setting.Settings.ValidateOrgDomains = *policyEvent.ValidateOrgDomains
// 		}
// 		if policyEvent.SMTPSenderAddressMatchesInstanceDomain != nil {
// 			setting.Settings.SMTPSenderAddressMatchesInstanceDomain = *policyEvent.SMTPSenderAddressMatchesInstanceDomain
// 		}

// 		err := settingsRepo.Set(ctx, v3_sql.SQLTx(tx), setting, settingsRepo.SetUpdatedAt(&policyEvent.Creation))
// 		return err
// 	}), nil
// }

// func (p *settingsRelationalProjection) reduceOrgDomainPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	policyEvent, ok := event.(*org.DomainPolicyRemovedEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.LockoutPolicyRemovedEventType)
// 	}
// 	return handler.NewStatement(policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UrdHy", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		orgID := &policyEvent.Aggregate().ID

// 		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx), policyEvent.Aggregate().InstanceID, orgID, domain.SettingTypeDomain)
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceSecurityPolicySet(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*instance.SecurityPolicySetEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-83U8p", "reduce.wrong.event.type %s", instance.SecurityPolicySetEventType)
// 	}

// 	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-lhPul", "reduce.wrong.db.pool %T", ex)
// 		}
// 		settingsRepo := repository.SettingsRepository()

// 		existingSetting, err := settingsRepo.GetSecurity(ctx, v3_sql.SQLTx(tx), e.Agg.InstanceID, nil)
// 		if err != nil && !errors.Is(err, new(database.NoRowFoundError)) {
// 			return zerrors.ThrowInternal(err, "HANDL-rSkxt", "error accessing login policy record")
// 		}
// 		if errors.Is(err, new(database.NoRowFoundError)) {
// 			setting := new(domain.SecuritySettings)
// 			if e.EnableIframeEmbedding != nil {
// 				setting.EnableIframeEmbedding = *e.EnableIframeEmbedding
// 			}
// 			if e.Enabled != nil {
// 				setting.Enabled = *e.Enabled
// 			}
// 			if e.AllowedOrigins != nil {
// 				setting.AllowedOrigins = *e.AllowedOrigins
// 			}
// 			if e.EnableImpersonation != nil {
// 				setting.EnableImpersonation = *e.EnableImpersonation
// 			}
// 			payload, err := json.Marshal(setting)
// 			if err != nil {
// 				return err
// 			}
// 			return settingsRepo.Create(ctx, v3_sql.SQLTx(tx), &domain.Setting{
// 				InstanceID: e.Aggregate().InstanceID,
// 				Type:       domain.SettingTypeSecurity,
// 				Settings:   payload,
// 				CreatedAt:  e.CreatedAt(),
// 				UpdatedAt:  &e.Creation,
// 			})
// 		}

// 		if e.EnableIframeEmbedding != nil {
// 			existingSetting.Settings.EnableIframeEmbedding = *e.EnableIframeEmbedding
// 		} else if e.Enabled != nil {
// 			existingSetting.Settings.Enabled = *e.Enabled
// 		}
// 		if e.AllowedOrigins != nil {
// 			existingSetting.Settings.AllowedOrigins = *e.AllowedOrigins
// 		}
// 		if e.EnableImpersonation != nil {
// 			existingSetting.Settings.EnableImpersonation = *e.EnableImpersonation
// 		}

// 		CreatedAt := event.CreatedAt()

// 		_, err = settingsRepo.UpdateSecurity(ctx, v3_sql.SQLTx(tx), existingSetting, settingsRepo.SetUpdatedAt(&CreatedAt))
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceOrganizationSettingsSet(event eventstore.Event) (*handler.Statement, error) {
// 	e, err := assertEvent[*settings.OrganizationSettingsSetEvent](event)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		settings, err := json.Marshal(e)
// 		if err != nil {
// 			return err
// 		}

// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-chluS", "reduce.wrong.db.pool %T", ex)
// 		}
// 		settingsRepo := repository.SettingsRepository()

// 		orgID := &e.Aggregate().ID

// 		existingSetting, err := settingsRepo.GetOrg(ctx, v3_sql.SQLTx(tx), e.Agg.InstanceID, nil)
// 		if err != nil {
// 			if errors.Is(err, &database.NoRowFoundError{}) {
// 				setting := domain.Setting{
// 					InstanceID: e.Aggregate().InstanceID,
// 					OrgID:      orgID,
// 					Type:       domain.SettingTypeOrganization,
// 					Settings:   settings,
// 					CreatedAt:  e.CreatedAt(),
// 					UpdatedAt:  &e.Creation,
// 				}
// 				err = settingsRepo.Create(ctx, v3_sql.SQLTx(tx), &setting)
// 				return err

// 			} else {
// 				return zerrors.ThrowInternal(err, "HANDL-uhk0t", "error accessing login policy record")
// 			}
// 		}

// 		existingSetting.Settings.OrganizationScopedUsernames = e.OrganizationScopedUsernames

// 		CreatedAt := event.CreatedAt()

// 		_, err = settingsRepo.UpdateOrg(ctx, v3_sql.SQLTx(tx), existingSetting, settingsRepo.SetUpdatedAt(&CreatedAt))
// 		return err
// 	}), nil
// }

// func (s *settingsRelationalProjection) reduceOrganizationSettingsRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	e, err := assertEvent[*settings.OrganizationSettingsRemovedEvent](event)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
// 		tx, ok := ex.(*sql.Tx)
// 		if !ok {
// 			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rHiHb", "reduce.wrong.db.pool %T", ex)
// 		}

// 		settingsRepo := repository.SettingsRepository()

// 		orgId := &event.Aggregate().ID

// 		_, err := settingsRepo.Delete(ctx, v3_sql.SQLTx(tx), e.Aggregate().InstanceID, orgId, domain.SettingTypeOrganization)
// 		return err
// 	}), nil
// }
