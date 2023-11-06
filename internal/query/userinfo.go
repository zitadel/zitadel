package query

import (
	"context"
	"encoding/base64"
	"slices"
	"strings"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"golang.org/x/text/language"
)

func (q *Queries) GetOIDCUserInfo(ctx context.Context, userID string, scope, roleAudience []string) (_ *OIDCUserInfo, err error) {
	if slices.Contains(scope, domain.ScopeProjectsRoles) {
		roleAudience = domain.AddAudScopeToAudience(ctx, roleAudience, scope)
		// TODO: we need to get the project roles and user roles.
	}

	user := newOidcUserinfoReadModel(userID, scope)
	if err = q.eventstore.FilterToQueryReducer(ctx, user); err != nil {
		return nil, err
	}

	if hasOrgScope(scope) {
		org := newoidcUserinfoOrganizationReadModel(user.ResourceOwner)
		if err = q.eventstore.FilterToQueryReducer(ctx, org); err != nil {
			return nil, err
		}

		user.OrgID = org.AggregateID
		user.OrgName = org.Name
		user.OrgPrimaryDomain = org.PrimaryDomain
	}

	return &user.OIDCUserInfo, nil
}

func hasOrgScope(scope []string) bool {
	return slices.ContainsFunc(scope, func(s string) bool {
		return s == domain.ScopeResourceOwner || strings.HasPrefix(s, domain.OrgIDScope)
	})
}

type OIDCUserInfo struct {
	ID                string
	UserName          string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	PreferredLanguage language.Tag
	Gender            domain.Gender
	Avatar            string
	UpdatedAt         time.Time

	Email           domain.EmailAddress
	IsEmailVerified bool

	Phone           domain.PhoneNumber
	IsPhoneVerified bool

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string

	UserState domain.UserState
	UserType  domain.UserType

	OrgID            string
	OrgName          string
	OrgPrimaryDomain string

	Metadata map[string]string
}

type oidcUserInfoReadmodel struct {
	eventstore.ReadModel
	scope []string // Scope is used to determine events
	OIDCUserInfo
}

func newOidcUserinfoReadModel(userID string, scope []string) *oidcUserInfoReadmodel {
	return &oidcUserInfoReadmodel{
		ReadModel: eventstore.ReadModel{
			AggregateID: userID,
		},
		scope: scope,
		OIDCUserInfo: OIDCUserInfo{
			ID: userID,
		},
	}
}

func (rm *oidcUserInfoReadmodel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		AllowTimeTravel().
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(rm.AggregateID).
		EventTypes(rm.scopeToEventTypes()...).
		Builder()
}

// scopeToEventTypes sets required user events to obtain get the correct userinfo.
// Events such as UserLocked, UserDeactivated and UserRemoved are not checked,
// as access tokens should already be revoked.
func (rm *oidcUserInfoReadmodel) scopeToEventTypes() []eventstore.EventType {
	types := make([]eventstore.EventType, 0, len(rm.scope))
	types = append(types, user.HumanAddedType, user.MachineAddedEventType)

	for _, scope := range rm.scope {
		switch scope {
		case domain.ScopeEmail:
			types = append(types, user.HumanEmailChangedType, user.HumanEmailVerifiedType)
		case domain.ScopeProfile:
			types = append(types, user.HumanProfileChangedType, user.HumanAvatarAddedType, user.HumanAvatarRemovedType)
		case domain.ScopePhone:
			types = append(types, user.HumanPhoneChangedType, user.HumanPhoneVerifiedType, user.HumanPhoneRemovedType)
		case domain.ScopeAddress:
			types = append(types, user.HumanAddressChangedType)
		case domain.ScopeUserMetaData:
			types = append(types, user.MetadataSetType, user.MetadataRemovedType, user.MetadataRemovedAllType)
		}
	}
	return slices.Compact(types)
}

func (rm *oidcUserInfoReadmodel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			rm.UserName = e.UserName
			rm.FirstName = e.FirstName
			rm.LastName = e.LastName
			rm.NickName = e.NickName
			rm.Name = e.DisplayName
			rm.PreferredLanguage = e.PreferredLanguage
			rm.Gender = e.Gender
			rm.Email = e.EmailAddress
			rm.Phone = e.PhoneNumber
			rm.Country = e.Country
			rm.Locality = e.Locality
			rm.PostalCode = e.PostalCode
			rm.Region = e.Region
			rm.StreetAddress = e.StreetAddress
			rm.UpdatedAt = e.Creation
		case *user.MachineAddedEvent:
			rm.UserName = e.UserName
			rm.Name = e.Name
			rm.UpdatedAt = e.Creation
		case *user.HumanEmailChangedEvent:
			rm.Email = e.EmailAddress
			rm.IsEmailVerified = false
			rm.UpdatedAt = e.Creation
		case *user.HumanEmailVerifiedEvent:
			rm.IsEmailVerified = e.IsEmailVerified
			rm.UpdatedAt = e.Creation
		case *user.HumanProfileChangedEvent:
			rm.FirstName = e.FirstName
			rm.LastName = e.LastName
			rm.NickName = gu.Value(e.NickName)
			rm.Name = gu.Value(e.DisplayName)
			rm.PreferredLanguage = gu.Value(e.PreferredLanguage)
			rm.Gender = gu.Value(e.Gender)
			rm.UpdatedAt = e.Creation
		case *user.HumanAvatarAddedEvent:
			rm.Avatar = e.StoreKey
			rm.UpdatedAt = e.Creation
		case *user.HumanAvatarRemovedEvent:
			rm.Avatar = ""
			rm.UpdatedAt = e.Creation
		case *user.HumanPhoneChangedEvent:
			rm.Phone = e.PhoneNumber
			rm.IsPhoneVerified = false
			rm.UpdatedAt = e.Creation
		case *user.HumanPhoneVerifiedEvent:
			rm.IsEmailVerified = e.IsPhoneVerified
			rm.UpdatedAt = e.Creation
		case *user.HumanPhoneRemovedEvent:
			rm.Phone = ""
			rm.IsPhoneVerified = false
			rm.UpdatedAt = e.Creation
		case *user.HumanAddressChangedEvent:
			rm.Country = gu.Value(e.Country)
			rm.Locality = gu.Value(e.Locality)
			rm.PostalCode = gu.Value(e.PostalCode)
			rm.Region = gu.Value(e.Region)
			rm.StreetAddress = gu.Value(e.StreetAddress)
			rm.UpdatedAt = e.Creation
		case *user.MetadataSetEvent:
			rm.Metadata[e.Key] = base64.RawURLEncoding.EncodeToString(e.Value)
			rm.UpdatedAt = e.Creation
		case *user.MetadataRemovedEvent:
			delete(rm.Metadata, e.Key)
			rm.UpdatedAt = e.Creation
		case *user.MetadataRemovedAllEvent:
			for key := range rm.Metadata {
				delete(rm.Metadata, key)
			}
			rm.UpdatedAt = e.Creation
		}
	}

	return rm.ReadModel.Reduce()
}

type oidcUserinfoOrganizationReadModel struct {
	eventstore.ReadModel

	Name          string
	PrimaryDomain string
}

func newoidcUserinfoOrganizationReadModel(orgID string) *oidcUserinfoOrganizationReadModel {
	return &oidcUserinfoOrganizationReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID: orgID,
		},
	}
}

func (rm *oidcUserinfoOrganizationReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		AllowTimeTravel().
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(rm.AggregateID).
		EventTypes(org.OrgAddedEventType, org.OrgChangedEventType, org.OrgDomainPrimarySetEventType).
		Builder()
}

func (rm *oidcUserinfoOrganizationReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *org.OrgAddedEvent:
			rm.Name = e.Name
		case *org.OrgChangedEvent:
			rm.Name = e.Name
		case *org.DomainPrimarySetEvent:
			rm.PrimaryDomain = e.Domain
		}
	}

	return rm.ReadModel.Reduce()
}
