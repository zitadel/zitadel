package projection

import (
	"context"
	"database/sql"
	"encoding/json"
	"slices"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceLDAPIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.LDAPIDPAddedEvent
	switch e := event.(type) {
	case *org.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPAddedEventType, instance.LDAPIDPAddedEventType})
	}

	ldap := domain.LDAP{
		Servers:           idpEvent.Servers,
		StartTLS:          idpEvent.StartTLS,
		BaseDN:            idpEvent.BaseDN,
		BindDN:            idpEvent.BindDN,
		BindPassword:      idpEvent.BindPassword,
		UserBase:          idpEvent.UserBase,
		UserObjectClasses: idpEvent.UserObjectClasses,
		UserFilters:       idpEvent.UserFilters,
		Timeout:           idpEvent.Timeout,
		LDAPAttributes: domain.LDAPAttributes{
			IDAttribute:                idpEvent.IDAttribute,
			FirstNameAttribute:         idpEvent.FirstNameAttribute,
			LastNameAttribute:          idpEvent.LastNameAttribute,
			DisplayNameAttribute:       idpEvent.DisplayNameAttribute,
			NickNameAttribute:          idpEvent.NickNameAttribute,
			PreferredUsernameAttribute: idpEvent.PreferredUsernameAttribute,
			EmailAttribute:             idpEvent.EmailAttribute,
			EmailVerifiedAttribute:     idpEvent.EmailVerifiedAttribute,
			PhoneAttribute:             idpEvent.PhoneAttribute,
			PhoneVerifiedAttribute:     idpEvent.PhoneVerifiedAttribute,
			PreferredLanguageAttribute: idpEvent.PreferredLanguageAttribute,
			AvatarURLAttribute:         idpEvent.AvatarURLAttribute,
			ProfileAttribute:           idpEvent.ProfileAttribute,
		},
	}

	payloadJSON, err := json.Marshal(ldap)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-XCJ8w", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeLDAP),
			AllowCreation:     idpEvent.IsCreationAllowed,
			AllowAutoCreation: idpEvent.IsAutoCreation,
			AllowAutoUpdate:   idpEvent.IsAutoUpdate,
			AllowLinking:      idpEvent.IsLinkingAllowed,
			AutoLinkingField:  mapAutoLinkingField(idpEvent.AutoLinkingOption),
			Payload:           payloadJSON,
			CreatedAt:         event.CreatedAt(),
			UpdatedAt:         event.CreatedAt(),
		})
	}), nil
}

func (p *relationalTablesProjection) reduceLDAPIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.LDAPIDPChangedEvent
	switch e := event.(type) {
	case *org.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPChangedEventType, instance.LDAPIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		ldap, err := idpRepo.GetLDAP(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &ldap.LDAP
		payloadChanged := p.reduceLDAPIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			changes = append(changes, idpRepo.SetPayload(string(payloadJSON)))
		}

		changes = append(changes, idpRepo.SetUpdatedAt(gu.Ptr(event.CreatedAt())))
		_, err = idpRepo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(idpRepo, idpEvent.Aggregate().InstanceID, idpEvent.ID, orgId), changes...)
		return err

	}), nil
}

//nolint:gocognit
func (p *relationalTablesProjection) reduceLDAPIDPChangedColumns(payload *domain.LDAP, idpEvent *idp.LDAPIDPChangedEvent) bool {
	payloadChanged := false
	if payload == nil || idpEvent == nil {
		return payloadChanged
	}

	if idpEvent.Servers != nil && !slices.Equal(idpEvent.Servers, payload.Servers) {
		payloadChanged = true
		payload.Servers = idpEvent.Servers
	}
	if idpEvent.StartTLS != nil && *idpEvent.StartTLS != payload.StartTLS {
		payloadChanged = true
		payload.StartTLS = *idpEvent.StartTLS
	}
	if idpEvent.BaseDN != nil && *idpEvent.BaseDN != payload.BaseDN {
		payloadChanged = true
		payload.BaseDN = *idpEvent.BaseDN
	}
	if idpEvent.BindDN != nil && *idpEvent.BindDN != payload.BindDN {
		payloadChanged = true
		payload.BindDN = *idpEvent.BindDN
	}
	if idpEvent.BindPassword != nil {
		payloadChanged = true
		payload.BindPassword = idpEvent.BindPassword
	}
	if idpEvent.UserBase != nil && *idpEvent.UserBase != payload.UserBase {
		payloadChanged = true
		payload.UserBase = *idpEvent.UserBase
	}
	if idpEvent.UserObjectClasses != nil && !slices.Equal(idpEvent.UserObjectClasses, payload.UserObjectClasses) {
		payloadChanged = true
		payload.UserObjectClasses = idpEvent.UserObjectClasses
	}
	if idpEvent.UserFilters != nil && !slices.Equal(idpEvent.UserFilters, payload.UserFilters) {
		payloadChanged = true
		payload.UserFilters = idpEvent.UserFilters
	}
	if idpEvent.Timeout != nil && *idpEvent.Timeout != payload.Timeout {
		payloadChanged = true
		payload.Timeout = *idpEvent.Timeout
	}
	if idpEvent.RootCA != nil && !slices.Equal(idpEvent.RootCA, payload.RootCA) {
		payloadChanged = true
		payload.RootCA = idpEvent.RootCA
	}
	if idpEvent.IDAttribute != nil && *idpEvent.IDAttribute != payload.IDAttribute {
		payloadChanged = true
		payload.IDAttribute = *idpEvent.IDAttribute
	}
	if idpEvent.FirstNameAttribute != nil && *idpEvent.FirstNameAttribute != payload.FirstNameAttribute {
		payloadChanged = true
		payload.FirstNameAttribute = *idpEvent.FirstNameAttribute
	}
	if idpEvent.LastNameAttribute != nil && *idpEvent.LastNameAttribute != payload.LastNameAttribute {
		payloadChanged = true
		payload.LastNameAttribute = *idpEvent.LastNameAttribute
	}
	if idpEvent.DisplayNameAttribute != nil && *idpEvent.DisplayNameAttribute != payload.DisplayNameAttribute {
		payloadChanged = true
		payload.DisplayNameAttribute = *idpEvent.DisplayNameAttribute
	}
	if idpEvent.NickNameAttribute != nil && *idpEvent.NickNameAttribute != payload.NickNameAttribute {
		payloadChanged = true
		payload.NickNameAttribute = *idpEvent.NickNameAttribute
	}
	if idpEvent.PreferredUsernameAttribute != nil && *idpEvent.PreferredUsernameAttribute != payload.PreferredUsernameAttribute {
		payloadChanged = true
		payload.PreferredUsernameAttribute = *idpEvent.PreferredUsernameAttribute
	}
	if idpEvent.EmailAttribute != nil && *idpEvent.EmailAttribute != payload.EmailAttribute {
		payloadChanged = true
		payload.EmailAttribute = *idpEvent.EmailAttribute
	}
	if idpEvent.EmailVerifiedAttribute != nil && *idpEvent.EmailVerifiedAttribute != payload.EmailVerifiedAttribute {
		payloadChanged = true
		payload.EmailVerifiedAttribute = *idpEvent.EmailVerifiedAttribute
	}
	if idpEvent.PhoneAttribute != nil && *idpEvent.PhoneAttribute != payload.PhoneAttribute {
		payloadChanged = true
		payload.PhoneAttribute = *idpEvent.PhoneAttribute
	}
	if idpEvent.PhoneVerifiedAttribute != nil && *idpEvent.PhoneVerifiedAttribute != payload.PhoneVerifiedAttribute {
		payloadChanged = true
		payload.PhoneVerifiedAttribute = *idpEvent.PhoneVerifiedAttribute
	}
	if idpEvent.PreferredLanguageAttribute != nil && *idpEvent.PreferredLanguageAttribute != payload.PreferredLanguageAttribute {
		payloadChanged = true
		payload.PreferredLanguageAttribute = *idpEvent.PreferredLanguageAttribute
	}
	if idpEvent.AvatarURLAttribute != nil && *idpEvent.AvatarURLAttribute != payload.AvatarURLAttribute {
		payloadChanged = true
		payload.AvatarURLAttribute = *idpEvent.AvatarURLAttribute
	}
	if idpEvent.ProfileAttribute != nil && *idpEvent.ProfileAttribute != payload.ProfileAttribute {
		payloadChanged = true
		payload.ProfileAttribute = *idpEvent.ProfileAttribute
	}
	return payloadChanged
}
