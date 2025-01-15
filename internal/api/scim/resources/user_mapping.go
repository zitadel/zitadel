package resources

import (
	"context"
	"strconv"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (h *UsersHandler) mapToAddHuman(ctx context.Context, scimUser *ScimUser) (*command.AddHuman, error) {
	human := &command.AddHuman{
		Username:    scimUser.UserName,
		NickName:    scimUser.NickName,
		DisplayName: scimUser.DisplayName,
	}

	if scimUser.Active != nil && !*scimUser.Active {
		human.SetInactive = true
	}

	if email := h.mapPrimaryEmail(scimUser); email != nil {
		human.Email = *email
	}

	if phone := h.mapPrimaryPhone(scimUser); phone != nil {
		human.Phone = *phone
	}

	md, err := h.mapMetadataToCommands(ctx, scimUser)
	if err != nil {
		return nil, err
	}
	human.Metadata = md

	if scimUser.Password != nil {
		human.Password = scimUser.Password.String()
		scimUser.Password = nil
	}

	if scimUser.Name != nil {
		human.FirstName = scimUser.Name.GivenName
		human.LastName = scimUser.Name.FamilyName

		// the direct mapping displayName => displayName has priority
		// over the formatted name assignment
		if human.DisplayName == "" {
			human.DisplayName = scimUser.Name.Formatted
		} else {
			// update user to match the actual stored value
			scimUser.Name.Formatted = human.DisplayName
		}
	}

	if err := domain.LanguageIsDefined(scimUser.PreferredLanguage); err != nil {
		human.PreferredLanguage = language.English
		scimUser.PreferredLanguage = language.English
	}

	return human, nil
}

func (h *UsersHandler) mapToChangeHuman(ctx context.Context, scimUser *ScimUser) (*command.ChangeHuman, error) {
	human := &command.ChangeHuman{
		ID:       scimUser.ID,
		Username: &scimUser.UserName,
		Profile: &command.Profile{
			NickName:    &scimUser.NickName,
			DisplayName: &scimUser.DisplayName,
		},
		Email: h.mapPrimaryEmail(scimUser),
		Phone: h.mapPrimaryPhone(scimUser),
	}

	if scimUser.Active != nil {
		if *scimUser.Active {
			human.State = gu.Ptr(domain.UserStateActive)
		} else {
			human.State = gu.Ptr(domain.UserStateInactive)
		}
	}

	md, mdRemovedKeys, err := h.mapMetadataToDomain(ctx, scimUser)
	if err != nil {
		return nil, err
	}
	human.Metadata = md
	human.MetadataKeysToRemove = mdRemovedKeys

	if scimUser.Password != nil {
		human.Password = &command.Password{
			Password: scimUser.Password.String(),
		}
		scimUser.Password = nil
	}

	if scimUser.Name != nil {
		human.Profile.FirstName = &scimUser.Name.GivenName
		human.Profile.LastName = &scimUser.Name.FamilyName

		// the direct mapping displayName => displayName has priority
		// over the formatted name assignment
		if *human.Profile.DisplayName == "" {
			human.Profile.DisplayName = &scimUser.Name.Formatted
		} else {
			// update user to match the actual stored value
			scimUser.Name.Formatted = *human.Profile.DisplayName
		}
	}

	if err := domain.LanguageIsDefined(scimUser.PreferredLanguage); err != nil {
		human.Profile.PreferredLanguage = &language.English
		scimUser.PreferredLanguage = language.English
	}

	return human, nil
}

func (h *UsersHandler) mapPrimaryEmail(scimUser *ScimUser) *command.Email {
	for _, email := range scimUser.Emails {
		if !email.Primary {
			continue
		}

		return &command.Email{
			Address:  domain.EmailAddress(email.Value),
			Verified: h.config.EmailVerified,
		}
	}

	return nil
}

func (h *UsersHandler) mapPrimaryPhone(scimUser *ScimUser) *command.Phone {
	for _, phone := range scimUser.PhoneNumbers {
		if !phone.Primary {
			continue
		}

		return &command.Phone{
			Number:   domain.PhoneNumber(phone.Value),
			Verified: h.config.PhoneVerified,
		}
	}

	return nil
}

func (h *UsersHandler) mapAddCommandToScimUser(ctx context.Context, user *ScimUser, addHuman *command.AddHuman) {
	user.ID = addHuman.Details.ID
	user.Resource = buildResource(ctx, h, addHuman.Details)
	user.Password = nil

	// ZITADEL supports only one (primary) phone number or email.
	// Therefore, only the primary one should be returned.
	// Note that the phone number might also be reformatted.
	if addHuman.Phone.Number != "" {
		user.PhoneNumbers = []*ScimPhoneNumber{
			{
				Value:   string(addHuman.Phone.Number),
				Primary: true,
			},
		}
	}

	if addHuman.Email.Address != "" {
		user.Emails = []*ScimEmail{
			{
				Value:   string(addHuman.Email.Address),
				Primary: true,
			},
		}
	}
}

func (h *UsersHandler) mapChangeCommandToScimUser(ctx context.Context, user *ScimUser, changeHuman *command.ChangeHuman) {
	user.ID = changeHuman.Details.ID
	user.Resource = buildResource(ctx, h, changeHuman.Details)
	user.Password = nil

	// ZITADEL supports only one (primary) phone number or email.
	// Therefore, only the primary one should be returned.
	// Note that the phone number might also be reformatted.
	if changeHuman.Phone != nil {
		user.PhoneNumbers = []*ScimPhoneNumber{
			{
				Value:   string(changeHuman.Phone.Number),
				Primary: true,
			},
		}
	}

	if changeHuman.Email != nil {
		user.Emails = []*ScimEmail{
			{
				Value:   string(changeHuman.Email.Address),
				Primary: true,
			},
		}
	}
}

func (h *UsersHandler) mapToScimUser(ctx context.Context, user *query.User, md map[metadata.ScopedKey][]byte) *ScimUser {
	scimUser := &ScimUser{
		Resource:     h.buildResourceForQuery(ctx, user),
		ID:           user.ID,
		ExternalID:   extractScalarMetadata(ctx, md, metadata.KeyExternalId),
		UserName:     user.Username,
		ProfileUrl:   extractHttpURLMetadata(ctx, md, metadata.KeyProfileUrl),
		Title:        extractScalarMetadata(ctx, md, metadata.KeyTitle),
		Locale:       extractScalarMetadata(ctx, md, metadata.KeyLocale),
		Timezone:     extractScalarMetadata(ctx, md, metadata.KeyTimezone),
		Active:       gu.Ptr(user.State.IsEnabled()),
		Ims:          make([]*ScimIms, 0),
		Addresses:    make([]*ScimAddress, 0),
		Photos:       make([]*ScimPhoto, 0),
		Entitlements: make([]*ScimEntitlement, 0),
		Roles:        make([]*ScimRole, 0),
	}

	if scimUser.Locale != "" {
		_, err := language.Parse(scimUser.Locale)
		if err != nil {
			logging.OnError(err).Warn("Failed to load locale of scim user")
			scimUser.Locale = ""
		}
	}

	if scimUser.Timezone != "" {
		_, err := time.LoadLocation(scimUser.Timezone)
		if err != nil {
			logging.OnError(err).Warn("Failed to load timezone of scim user")
			scimUser.Timezone = ""
		}
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyIms, &scimUser.Ims); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim ims metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyAddresses, &scimUser.Addresses); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim addresses metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyPhotos, &scimUser.Photos); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim photos metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyEntitlements, &scimUser.Entitlements); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim entitlements metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyRoles, &scimUser.Roles); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim roles metadata")
	}

	if user.Human != nil {
		mapHumanToScimUser(ctx, user.Human, scimUser, md)
	}

	return scimUser
}

func mapHumanToScimUser(ctx context.Context, human *query.Human, user *ScimUser, md map[metadata.ScopedKey][]byte) {
	user.DisplayName = human.DisplayName
	user.NickName = human.NickName
	user.PreferredLanguage = human.PreferredLanguage
	user.Name = &ScimUserName{
		Formatted:       human.DisplayName,
		FamilyName:      human.LastName,
		GivenName:       human.FirstName,
		MiddleName:      extractScalarMetadata(ctx, md, metadata.KeyMiddleName),
		HonorificPrefix: extractScalarMetadata(ctx, md, metadata.KeyHonorificPrefix),
		HonorificSuffix: extractScalarMetadata(ctx, md, metadata.KeyHonorificSuffix),
	}

	if string(human.Email) != "" {
		user.Emails = []*ScimEmail{
			{
				Value:   string(human.Email),
				Primary: true,
			},
		}
	}

	if string(human.Phone) != "" {
		user.PhoneNumbers = []*ScimPhoneNumber{
			{
				Value:   string(human.Phone),
				Primary: true,
			},
		}
	}
}

func (h *UsersHandler) buildResourceForQuery(ctx context.Context, user *query.User) *Resource {
	return &Resource{
		Schemas: []schemas.ScimSchemaType{schemas.IdUser},
		Meta: &ResourceMeta{
			ResourceType: schemas.UserResourceType,
			Created:      user.CreationDate.UTC(),
			LastModified: user.ChangeDate.UTC(),
			Version:      strconv.FormatUint(user.Sequence, 10),
			Location:     buildLocation(ctx, h, user.ID),
		},
	}
}

func cascadingMemberships(memberships []*query.Membership) []*command.CascadingMembership {
	cascades := make([]*command.CascadingMembership, len(memberships))
	for i, membership := range memberships {
		cascades[i] = &command.CascadingMembership{
			UserID:        membership.UserID,
			ResourceOwner: membership.ResourceOwner,
			IAM:           cascadingIAMMembership(membership.IAM),
			Org:           cascadingOrgMembership(membership.Org),
			Project:       cascadingProjectMembership(membership.Project),
			ProjectGrant:  cascadingProjectGrantMembership(membership.ProjectGrant),
		}
	}
	return cascades
}

func cascadingIAMMembership(membership *query.IAMMembership) *command.CascadingIAMMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingIAMMembership{IAMID: membership.IAMID}
}

func cascadingOrgMembership(membership *query.OrgMembership) *command.CascadingOrgMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingOrgMembership{OrgID: membership.OrgID}
}

func cascadingProjectMembership(membership *query.ProjectMembership) *command.CascadingProjectMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingProjectMembership{ProjectID: membership.ProjectID}
}

func cascadingProjectGrantMembership(membership *query.ProjectGrantMembership) *command.CascadingProjectGrantMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingProjectGrantMembership{ProjectID: membership.ProjectID, GrantID: membership.GrantID}
}

func userGrantsToIDs(userGrants []*query.UserGrant) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}
