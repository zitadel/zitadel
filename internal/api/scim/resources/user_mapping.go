package resources

import (
	"context"
	"strconv"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
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

	if email, err := h.mapPrimaryEmail(scimUser); err != nil {
		return nil, err
	} else {
		human.Email = email
	}

	if phone := h.mapPrimaryPhone(scimUser); phone != nil {
		human.Phone = *phone
	}

	md, err := h.mapMetadataToCommands(ctx, scimUser)
	if err != nil {
		return nil, err
	}
	human.Metadata = md

	// Okta sends a random password during SCIM provisioning
	// irrespective of whether the Sync Password option is enabled or disabled on Okta.
	// This password does not comply with Zitadel's password complexity, and
	// the following workaround ignores the random password as it does not add any value.
	ignorePasswordOnCreate := metadata.GetScimContextData(ctx).IgnorePasswordOnCreate
	if scimUser.Password != nil && !ignorePasswordOnCreate {
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
		ID:            scimUser.ID,
		ResourceOwner: authz.GetCtxData(ctx).OrgID,
		Username:      &scimUser.UserName,
		Profile: &command.Profile{
			NickName:    &scimUser.NickName,
			DisplayName: &scimUser.DisplayName,
		},
		Phone: h.mapPrimaryPhone(scimUser),
	}

	if human.Phone == nil {
		human.Phone = &command.Phone{Remove: true}
	}

	if email, err := h.mapPrimaryEmail(scimUser); err != nil {
		return nil, err
	} else {
		human.Email = &email
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

		if scimUser.Name.GivenName == "" || scimUser.Name.FamilyName == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "SCIM-USN1", "The name of a user is mandatory")
		}
	} else {
		return nil, zerrors.ThrowInvalidArgument(nil, "SCIM-USN2", "The name of a user is mandatory")
	}

	if err := domain.LanguageIsDefined(scimUser.PreferredLanguage); err != nil {
		human.Profile.PreferredLanguage = &language.English
		scimUser.PreferredLanguage = language.English
	}

	return human, nil
}

func (h *UsersHandler) mapPrimaryEmail(scimUser *ScimUser) (command.Email, error) {
	for _, email := range scimUser.Emails {
		if !email.Primary {
			continue
		}

		return command.Email{
			Address:  domain.EmailAddress(email.Value),
			Verified: h.config.EmailVerified,
		}, nil
	}

	// if no primary email was found, the first email will be used
	for _, email := range scimUser.Emails {
		email.Primary = true
		return command.Email{
			Address:  domain.EmailAddress(email.Value),
			Verified: h.config.EmailVerified,
		}, nil
	}

	return command.Email{}, zerrors.ThrowInvalidArgument(nil, "SCIM-EM19", "Errors.User.Email.Empty")
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

	// if no primary phone was found, the first phone will be used
	for _, phone := range scimUser.PhoneNumbers {
		phone.Primary = true
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

func (h *UsersHandler) mapToScimUsers(ctx context.Context, users []*query.User, md map[string]map[metadata.ScopedKey][]byte) []*ScimUser {
	result := make([]*ScimUser, len(users))
	for i, user := range users {
		userMetadata, ok := md[user.ID]
		if !ok {
			userMetadata = make(map[metadata.ScopedKey][]byte)
		}

		result[i] = h.mapToScimUser(ctx, user, userMetadata)
	}

	return result
}

func (h *UsersHandler) mapToScimUser(ctx context.Context, user *query.User, md map[metadata.ScopedKey][]byte) *ScimUser {
	scimUser := &ScimUser{
		Resource:          h.buildResourceForQuery(ctx, user),
		ID:                user.ID,
		UserName:          user.Username,
		DisplayName:       user.Human.DisplayName,
		NickName:          user.Human.NickName,
		PreferredLanguage: user.Human.PreferredLanguage,
		Name: &ScimUserName{
			Formatted:  user.Human.DisplayName,
			FamilyName: user.Human.LastName,
			GivenName:  user.Human.FirstName,
		},
		Active: schemas.NewRelaxedBool(user.State.IsEnabled()),
	}

	if string(user.Human.Email) != "" {
		scimUser.Emails = []*ScimEmail{
			{
				Value:   string(user.Human.Email),
				Primary: true,
			},
		}
	}

	if string(user.Human.Phone) != "" {
		scimUser.PhoneNumbers = []*ScimPhoneNumber{
			{
				Value:   string(user.Human.Phone),
				Primary: true,
			},
		}
	}

	h.mapAndValidateMetadata(ctx, scimUser, md)
	return scimUser
}

func (h *UsersHandler) mapWriteModelToScimUser(ctx context.Context, user *command.UserV2WriteModel) *ScimUser {
	scimUser := &ScimUser{
		Resource:          h.buildResourceForWriteModel(ctx, user),
		ID:                user.AggregateID,
		UserName:          user.UserName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage,
		Name: &ScimUserName{
			Formatted:  user.DisplayName,
			FamilyName: user.LastName,
			GivenName:  user.FirstName,
		},
		Active: schemas.NewRelaxedBool(user.UserState.IsEnabled()),
	}

	if string(user.Email) != "" {
		scimUser.Emails = []*ScimEmail{
			{
				Value:   string(user.Email),
				Primary: true,
			},
		}
	}

	if string(user.Phone) != "" {
		scimUser.PhoneNumbers = []*ScimPhoneNumber{
			{
				Value:   string(user.Phone),
				Primary: true,
			},
		}
	}

	md := metadata.MapToScopedKeyMap(user.Metadata)
	h.mapAndValidateMetadata(ctx, scimUser, md)
	return scimUser
}

func (h *UsersHandler) mapAndValidateMetadata(ctx context.Context, user *ScimUser, md map[metadata.ScopedKey][]byte) {
	user.ExternalID = extractScalarMetadata(ctx, md, metadata.KeyExternalId)
	user.ProfileUrl = extractHttpURLMetadata(ctx, md, metadata.KeyProfileUrl)
	user.Title = extractScalarMetadata(ctx, md, metadata.KeyTitle)
	user.Locale = extractScalarMetadata(ctx, md, metadata.KeyLocale)
	user.Timezone = extractScalarMetadata(ctx, md, metadata.KeyTimezone)
	user.Name.MiddleName = extractScalarMetadata(ctx, md, metadata.KeyMiddleName)
	user.Name.HonorificPrefix = extractScalarMetadata(ctx, md, metadata.KeyHonorificPrefix)
	user.Name.HonorificSuffix = extractScalarMetadata(ctx, md, metadata.KeyHonorificSuffix)

	if user.Locale != "" {
		_, err := language.Parse(user.Locale)
		if err != nil {
			logging.OnError(err).Warn("Failed to load locale of scim user")
			user.Locale = ""
		}
	}

	if user.Timezone != "" {
		_, err := time.LoadLocation(user.Timezone)
		if err != nil {
			logging.OnError(err).Warn("Failed to load timezone of scim user")
			user.Timezone = ""
		}
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyIms, &user.Ims); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim ims metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyAddresses, &user.Addresses); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim addresses metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyPhotos, &user.Photos); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim photos metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyEntitlements, &user.Entitlements); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim entitlements metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyRoles, &user.Roles); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim roles metadata")
	}

	if err := extractJsonMetadata(ctx, md, metadata.KeyEmails, &user.Emails); err != nil {
		logging.OnError(err).Warn("Could not deserialize scim emails metadata")
	}
}

func (h *UsersHandler) buildResourceForQuery(ctx context.Context, user *query.User) *schemas.Resource {
	return &schemas.Resource{
		ID:      user.ID,
		Schemas: []schemas.ScimSchemaType{schemas.IdUser},
		Meta: &schemas.ResourceMeta{
			ResourceType: schemas.UserResourceType,
			Created:      gu.Ptr(user.CreationDate.UTC()),
			LastModified: gu.Ptr(user.ChangeDate.UTC()),
			Version:      strconv.FormatUint(user.Sequence, 10),
			Location:     schemas.BuildLocationForResource(ctx, h.schema.PluralName, user.ID),
		},
	}
}

func (h *UsersHandler) buildResourceForWriteModel(ctx context.Context, user *command.UserV2WriteModel) *schemas.Resource {
	return &schemas.Resource{
		Schemas: []schemas.ScimSchemaType{schemas.IdUser},
		Meta: &schemas.ResourceMeta{
			ResourceType: schemas.UserResourceType,
			Created:      gu.Ptr(user.CreationDate.UTC()),
			LastModified: gu.Ptr(user.ChangeDate.UTC()),
			Version:      strconv.FormatUint(user.ProcessedSequence, 10),
			Location:     schemas.BuildLocationForResource(ctx, h.schema.PluralName, user.AggregateID),
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

func usersToIDs(users []*query.User) []string {
	ids := make([]string, len(users))
	for i, user := range users {
		ids[i] = user.ID
	}
	return ids
}
