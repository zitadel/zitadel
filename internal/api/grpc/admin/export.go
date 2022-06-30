package admin

import (
	"context"
	text_grpc "github.com/zitadel/zitadel/internal/api/grpc/text"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	action_pb "github.com/zitadel/zitadel/pkg/grpc/action"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	app_pb "github.com/zitadel/zitadel/pkg/grpc/app"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp"
	management_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	org_pb "github.com/zitadel/zitadel/pkg/grpc/org"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (s *Server) ExportData(ctx context.Context, req *admin_pb.ExportDataRequest) (_ *admin_pb.ExportDataResponse, err error) {

	orgSearchQuery := &query.OrgSearchQueries{}
	if len(req.OrgIds) > 0 {
		orgIDsSearchQuery, err := query.NewOrgIDsSearchQuery(req.OrgIds...)
		if err != nil {
			return nil, err
		}
		orgSearchQuery.Queries = []query.SearchQuery{orgIDsSearchQuery}
	}
	queriedOrgs, err := s.query.SearchOrgs(ctx, orgSearchQuery)
	if err != nil {
		return nil, err
	}

	orgs := make([]*admin_pb.DataOrg, len(queriedOrgs.Orgs))
	processedOrgs := make([]string, len(queriedOrgs.Orgs))
	processedProjects := make([]string, 0)
	processedGrants := make([]string, 0)
	processedUsers := make([]string, 0)
	processedActions := make([]string, 0)

	for i, queriedOrg := range queriedOrgs.Orgs {
		if req.ExcludedOrgIds != nil {
			found := false
			for _, excludedOrg := range req.ExcludedOrgIds {
				if excludedOrg == queriedOrg.ID {
					found = true
				}
			}
			if found {
				continue
			}
		}
		processedOrgs = append(processedOrgs, queriedOrg.ID)

		/******************************************************************************************************************
		Organization
		******************************************************************************************************************/
		org := &admin_pb.DataOrg{OrgId: queriedOrg.ID, Org: &management_pb.AddOrgRequest{Name: queriedOrg.Name}}
		orgs[i] = org
	}

	for _, org := range orgs {
		org.IamPolicy, err = s.getIAMPolicy(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		org.Domains, err = s.getDomains(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		org.OidcIdps, org.JwtIdps, err = s.getIDPs(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		org.LabelPolicy, err = s.getLabelPolicy(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		org.LoginPolicy, org.SecondFactors, org.MultiFactors, org.Idps, err = s.getLoginPolicy(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		org.UserLinks, err = s.getUserLinks(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		org.LockoutPolicy, err = s.getLockoutPolicy(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		org.PasswordComplexityPolicy, err = s.getPasswordComplexityPolicy(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		org.PrivacyPolicy, err = s.getPrivacyPolicy(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}

		langResp, err := s.GetSupportedLanguages(ctx, &admin_pb.GetSupportedLanguagesRequest{})
		if err != nil {
			return nil, err
		}

		org.LoginTexts, err = s.getCustomLoginTexts(ctx, org.GetOrgId(), langResp.Languages)
		if err != nil {
			return nil, err
		}

		org.InitMessages, err = s.getCustomInitMessageTexts(ctx, org.GetOrgId(), langResp.Languages)
		if err != nil {
			return nil, err
		}

		org.PasswordResetMessages, err = s.getCustomPasswordResetMessageTexts(ctx, org.GetOrgId(), langResp.Languages)
		if err != nil {
			return nil, err
		}

		org.VerifyEmailMessages, err = s.getCustomVerifyEmailMessageTexts(ctx, org.GetOrgId(), langResp.Languages)
		if err != nil {
			return nil, err
		}

		org.VerifyPhoneMessages, err = s.getCustomVerifyPhoneMessageTexts(ctx, org.GetOrgId(), langResp.Languages)
		if err != nil {
			return nil, err
		}

		org.DomainClaimedMessages, err = s.getCustomDomainClaimedMessageTexts(ctx, org.GetOrgId(), langResp.Languages)
		if err != nil {
			return nil, err
		}

		org.PasswordlessRegistrationMessages, err = s.getCustomPasswordlessRegistrationMessageTexts(ctx, org.GetOrgId(), langResp.Languages)
		if err != nil {
			return nil, err
		}

		/******************************************************************************************************************
		Users
		******************************************************************************************************************/
		org.HumanUsers, org.MachineUsers, org.UserMetadata, err = s.getUsers(ctx, org.GetOrgId(), req.WithPasswords, req.WithOtp)
		if err != nil {
			return nil, err
		}
		for _, processedUser := range org.HumanUsers {
			processedUsers = append(processedUsers, processedUser.UserId)
		}
		for _, processedUser := range org.MachineUsers {
			processedUsers = append(processedUsers, processedUser.UserId)
		}

		/******************************************************************************************************************
		Project and Applications
		******************************************************************************************************************/
		org.Projects, org.ProjectRoles, org.OidcApps, org.ApiApps, err = s.getProjectsAndApps(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}
		for _, processedProject := range org.Projects {
			processedProjects = append(processedProjects, processedProject.ProjectId)
		}

		/******************************************************************************************************************
		Actions
		******************************************************************************************************************/
		org.Actions, err = s.getActions(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}
		for _, processedAction := range org.Actions {
			processedActions = append(processedActions, processedAction.ActionId)
		}
	}

	for _, org := range orgs {
		/******************************************************************************************************************
		  Flows
		  ******************************************************************************************************************/
		org.TriggerActions, err = s.getTriggerActions(ctx, org.OrgId, processedActions)
		if err != nil {
			return nil, err
		}

		/******************************************************************************************************************
		  Grants
		  ******************************************************************************************************************/
		org.ProjectGrants, err = s.getNecessaryProjectGrantsForOrg(ctx, org.OrgId, processedOrgs, processedProjects)
		if err != nil {
			return nil, err
		}
		for _, processedGrant := range org.ProjectGrants {
			processedGrants = append(processedGrants, processedGrant.GrantId)
		}

		org.UserGrants, err = s.getNecessaryUserGrantsForOrg(ctx, org.OrgId, processedProjects, processedGrants, processedUsers)
		if err != nil {
			return nil, err
		}
	}

	for _, org := range orgs {
		/******************************************************************************************************************
		  Members
		  ******************************************************************************************************************/
		org.OrgMembers, err = s.getNecessaryOrgMembersForOrg(ctx, org.OrgId, processedUsers)
		if err != nil {
			return nil, err
		}

		org.ProjectMembers, err = s.getNecessaryProjectMembersForOrg(ctx, processedProjects, processedUsers)
		if err != nil {
			return nil, err
		}

		org.ProjectGrantMembers, err = s.getNecessaryProjectGrantMembersForOrg(ctx, org.OrgId, processedProjects, processedGrants, processedUsers)
		if err != nil {
			return nil, err
		}
	}

	return &admin_pb.ExportDataResponse{
		Orgs: orgs,
	}, nil
}

func (s *Server) getIAMPolicy(ctx context.Context, orgID string) (*admin_pb.AddCustomOrgIAMPolicyRequest, error) {
	queriedIAMPolicy, err := s.query.OrgIAMPolicyByOrg(ctx, false, orgID)
	if err != nil {
		return nil, err
	}
	if !queriedIAMPolicy.IsDefault {
		return &admin_pb.AddCustomOrgIAMPolicyRequest{
			OrgId:                 orgID,
			UserLoginMustBeDomain: queriedIAMPolicy.UserLoginMustBeDomain,
		}, nil
	}
	return nil, nil
}

func (s *Server) getDomains(ctx context.Context, orgID string) ([]*org_pb.Domain, error) {
	orgDomainOrgIDQuery, err := query.NewOrgDomainOrgIDSearchQuery(orgID)
	if err != nil {
		return nil, err
	}
	orgDomainsQuery, err := s.query.SearchOrgDomains(ctx, &query.OrgDomainSearchQueries{Queries: []query.SearchQuery{orgDomainOrgIDQuery}})
	if err != nil {
		return nil, err
	}
	orgDomains := make([]*org_pb.Domain, len(orgDomainsQuery.Domains))
	for i, orgDomain := range orgDomainsQuery.Domains {
		orgDomains[i] = &org_pb.Domain{
			OrgId:          orgDomain.OrgID,
			DomainName:     orgDomain.Domain,
			IsVerified:     orgDomain.IsVerified,
			IsPrimary:      orgDomain.IsPrimary,
			ValidationType: org_pb.DomainValidationType(orgDomain.ValidationType),
		}
	}
	return orgDomains, nil
}

func (s *Server) getIDPs(ctx context.Context, orgID string) ([]*admin_pb.DataOIDCIDP, []*admin_pb.DataJWTIDP, error) {
	ownerType, err := query.NewIDPOwnerTypeSearchQuery(domain.IdentityProviderTypeOrg)
	if err != nil {
		return nil, nil, err
	}
	idpQuery, err := query.NewIDPResourceOwnerSearchQuery(orgID)
	if err != nil {
		return nil, nil, err
	}
	idps, err := s.query.IDPs(ctx, &query.IDPSearchQueries{Queries: []query.SearchQuery{idpQuery, ownerType}})
	if err != nil {
		return nil, nil, err
	}
	oidcIdps := make([]*admin_pb.DataOIDCIDP, 0)
	jwtIdps := make([]*admin_pb.DataJWTIDP, 0)
	for _, idp := range idps.IDPs {
		if idp.OIDCIDP != nil {
			clientSecret, err := s.query.GetOIDCIDPClientSecret(ctx, false, orgID, idp.ID)
			if err != nil && !caos_errors.IsNotFound(err) {
				return nil, nil, err
			}
			oidcIdps = append(oidcIdps, &admin_pb.DataOIDCIDP{
				IdpId: idp.ID,
				Idp: &management_pb.AddOrgOIDCIDPRequest{
					Name:               idp.Name,
					StylingType:        idp_pb.IDPStylingType(idp.StylingType),
					ClientId:           idp.ClientID,
					ClientSecret:       clientSecret,
					Issuer:             idp.OIDCIDP.Issuer,
					Scopes:             idp.Scopes,
					DisplayNameMapping: idp_pb.OIDCMappingField(idp.DisplayNameMapping),
					UsernameMapping:    idp_pb.OIDCMappingField(idp.UsernameMapping),
					AutoRegister:       idp.AutoRegister,
				},
			})
		} else if idp.JWTIDP != nil {
			jwtIdps = append(jwtIdps, &admin_pb.DataJWTIDP{
				IdpId: idp.ID,
				Idp: &management_pb.AddOrgJWTIDPRequest{
					Name:         idp.Name,
					StylingType:  idp_pb.IDPStylingType(idp.StylingType),
					JwtEndpoint:  idp.JWTIDP.Endpoint,
					Issuer:       idp.JWTIDP.Issuer,
					KeysEndpoint: idp.KeysEndpoint,
					HeaderName:   idp.HeaderName,
					AutoRegister: idp.AutoRegister,
				},
			})
		}
	}
	return oidcIdps, jwtIdps, nil
}

func (s *Server) getLabelPolicy(ctx context.Context, orgID string) (*management_pb.AddCustomLabelPolicyRequest, error) {
	queriedLabel, err := s.query.ActiveLabelPolicyByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if !queriedLabel.IsDefault {
		return &management_pb.AddCustomLabelPolicyRequest{
			PrimaryColor:        queriedLabel.Light.PrimaryColor,
			HideLoginNameSuffix: queriedLabel.HideLoginNameSuffix,
			WarnColor:           queriedLabel.Light.WarnColor,
			BackgroundColor:     queriedLabel.Light.BackgroundColor,
			FontColor:           queriedLabel.Light.FontColor,
			PrimaryColorDark:    queriedLabel.Dark.PrimaryColor,
			BackgroundColorDark: queriedLabel.Dark.BackgroundColor,
			WarnColorDark:       queriedLabel.Dark.WarnColor,
			FontColorDark:       queriedLabel.Dark.FontColor,
			DisableWatermark:    queriedLabel.WatermarkDisabled,
		}, nil
	}
	return nil, nil
}

func (s *Server) getLoginPolicy(ctx context.Context, orgID string) (*management_pb.AddCustomLoginPolicyRequest, []*management_pb.AddSecondFactorToLoginPolicyRequest, []*management_pb.AddMultiFactorToLoginPolicyRequest, []*management_pb.AddIDPToLoginPolicyRequest, error) {
	queriedLogin, err := s.query.LoginPolicyByID(ctx, false, orgID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if !queriedLogin.IsDefault {

		secondFactors := make([]*management_pb.AddSecondFactorToLoginPolicyRequest, len(queriedLogin.SecondFactors))
		if queriedLogin.SecondFactors != nil {
			for i, factor := range queriedLogin.SecondFactors {
				secondFactors[i] = &management_pb.AddSecondFactorToLoginPolicyRequest{Type: policy_pb.SecondFactorType(factor)}
			}
		}

		multiFactors := make([]*management_pb.AddMultiFactorToLoginPolicyRequest, len(queriedLogin.MultiFactors))
		if queriedLogin.MultiFactors != nil {
			for i, factor := range queriedLogin.MultiFactors {
				multiFactors[i] = &management_pb.AddMultiFactorToLoginPolicyRequest{Type: policy_pb.MultiFactorType(factor)}
			}
		}

		queriedIdpLinks, err := s.query.IDPLoginPolicyLinks(ctx, orgID, &query.IDPLoginPolicyLinksSearchQuery{})
		if err != nil {
			return nil, nil, nil, nil, err
		}
		idpLinks := make([]*management_pb.AddIDPToLoginPolicyRequest, len(queriedIdpLinks.Links))
		if queriedIdpLinks != nil {
			for i, idpLink := range queriedIdpLinks.Links {
				idpLinks[i] = &management_pb.AddIDPToLoginPolicyRequest{
					IdpId:     idpLink.IDPID,
					OwnerType: idp_pb.IDPOwnerType(idpLink.IDPType),
				}
			}
		}

		return &management_pb.AddCustomLoginPolicyRequest{
			AllowUsernamePassword:  queriedLogin.AllowUsernamePassword,
			AllowRegister:          queriedLogin.AllowRegister,
			AllowExternalIdp:       queriedLogin.AllowExternalIDPs,
			ForceMfa:               queriedLogin.ForceMFA,
			PasswordlessType:       policy_pb.PasswordlessType(queriedLogin.PasswordlessType),
			HidePasswordReset:      queriedLogin.HidePasswordReset,
			IgnoreUnknownUsernames: queriedLogin.IgnoreUnknownUsernames,
			DefaultRedirectUri:     queriedLogin.DefaultRedirectURI,
		}, secondFactors, multiFactors, idpLinks, nil
	}

	return nil, nil, nil, nil, nil
}

func (s *Server) getUserLinks(ctx context.Context, orgID string) ([]*idp_pb.IDPUserLink, error) {
	userLinksResourceOwner, err := query.NewIDPUserLinksResourceOwnerSearchQuery(orgID)
	if err != nil {
		return nil, err
	}
	idpUserLinks, err := s.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: []query.SearchQuery{userLinksResourceOwner}})
	if err != nil && !caos_errors.IsNotFound(err) {
		return nil, err
	}
	userLinks := make([]*idp_pb.IDPUserLink, 0)
	if !caos_errors.IsNotFound(err) && idpUserLinks != nil {
		for _, idpUserLink := range idpUserLinks.Links {
			userLinks = append(userLinks, &idp_pb.IDPUserLink{
				UserId:           idpUserLink.UserID,
				IdpId:            idpUserLink.IDPID,
				IdpName:          idpUserLink.IDPName,
				ProvidedUserId:   idpUserLink.ProvidedUserID,
				ProvidedUserName: idpUserLink.ProvidedUsername,
				IdpType:          idp_pb.IDPType(idpUserLink.IDPType),
			})
		}
	}

	return userLinks, nil
}

func (s *Server) getLockoutPolicy(ctx context.Context, orgID string) (*management_pb.AddCustomLockoutPolicyRequest, error) {
	queriedLockout, err := s.query.LockoutPolicyByOrg(ctx, false, orgID)
	if err != nil {
		return nil, err
	}
	if !queriedLockout.IsDefault {
		return &management_pb.AddCustomLockoutPolicyRequest{
			MaxPasswordAttempts: uint32(queriedLockout.MaxPasswordAttempts),
		}, nil
	}
	return nil, nil
}

func (s *Server) getPasswordComplexityPolicy(ctx context.Context, orgID string) (*management_pb.AddCustomPasswordComplexityPolicyRequest, error) {
	queriedPasswordComplexity, err := s.query.PasswordComplexityPolicyByOrg(ctx, false, orgID)
	if err != nil {
		return nil, err
	}
	if !queriedPasswordComplexity.IsDefault {
		return &management_pb.AddCustomPasswordComplexityPolicyRequest{
			MinLength:    queriedPasswordComplexity.MinLength,
			HasUppercase: queriedPasswordComplexity.HasUppercase,
			HasLowercase: queriedPasswordComplexity.HasLowercase,
			HasNumber:    queriedPasswordComplexity.HasNumber,
			HasSymbol:    queriedPasswordComplexity.HasSymbol,
		}, nil
	}
	return nil, nil
}

func (s *Server) getPrivacyPolicy(ctx context.Context, orgID string) (*management_pb.AddCustomPrivacyPolicyRequest, error) {
	queriedPrivacy, err := s.query.PrivacyPolicyByOrg(ctx, false, orgID)
	if err != nil {
		return nil, err
	}
	if !queriedPrivacy.IsDefault {
		return &management_pb.AddCustomPrivacyPolicyRequest{
			TosLink:     queriedPrivacy.TOSLink,
			PrivacyLink: queriedPrivacy.PrivacyLink,
			HelpLink:    queriedPrivacy.HelpLink,
		}, nil
	}
	return nil, nil
}

func (s *Server) getUsers(ctx context.Context, org string, withPasswords bool, withOTP bool) ([]*admin_pb.DataHumanUser, []*admin_pb.DataMachineUser, []*management_pb.SetUserMetadataRequest, error) {
	orgSearch, err := query.NewUserResourceOwnerSearchQuery(org, query.TextEquals)
	if err != nil {
		return nil, nil, nil, err
	}
	users, err := s.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: []query.SearchQuery{orgSearch}})
	if err != nil && !caos_errors.IsNotFound(err) {
		return nil, nil, nil, err
	}
	humanUsers := make([]*admin_pb.DataHumanUser, 0)
	machineUsers := make([]*admin_pb.DataMachineUser, 0)
	userMetadata := make([]*management_pb.SetUserMetadataRequest, 0)
	if err != nil && caos_errors.IsNotFound(err) {
		return humanUsers, machineUsers, userMetadata, nil
	}
	for _, user := range users.Users {
		switch user.Type {
		case domain.UserTypeHuman:
			dataUser := &admin_pb.DataHumanUser{
				UserId: user.ID,
				User: &admin_pb.ExportHumanUser{
					UserName: user.Username,
					Profile: &admin_pb.ExportHumanUser_Profile{
						FirstName:         user.Human.FirstName,
						LastName:          user.Human.LastName,
						NickName:          user.Human.NickName,
						DisplayName:       user.Human.DisplayName,
						PreferredLanguage: user.Human.PreferredLanguage.String(),
						Gender:            user_pb.Gender(user.Human.Gender),
					},
				},
			}
			if user.Human.Email != "" {
				dataUser.User.Email = &admin_pb.ExportHumanUser_Email{
					Email:           user.Human.Email,
					IsEmailVerified: user.Human.IsEmailVerified,
				}
			}
			if user.Human.Phone != "" {
				dataUser.User.Phone = &admin_pb.ExportHumanUser_Phone{
					Phone:           user.Human.Phone,
					IsPhoneVerified: user.Human.IsPhoneVerified,
				}
			}
			if withPasswords {
				hashedPassword, hashAlgorithm, err := s.query.GetHumanPassword(ctx, org, user.ID)
				if err != nil && !caos_errors.IsNotFound(err) {
					return nil, nil, nil, err
				}
				if err == nil && hashedPassword != nil {
					dataUser.User.HashedPassword = &admin_pb.ExportHumanUser_HashedPassword{
						Value:     string(hashedPassword),
						Algorithm: hashAlgorithm,
					}
				}
			}
			if withOTP {
				code, err := s.query.GetHumanOTPSecret(ctx, user.ID, org)
				if err != nil && !caos_errors.IsNotFound(err) {
					return nil, nil, nil, err
				}
				if err == nil && code != "" {
					dataUser.User.OtpCode = code
				}
			}

			humanUsers = append(humanUsers, dataUser)
		case domain.UserTypeMachine:
			machineUsers = append(machineUsers, &admin_pb.DataMachineUser{
				UserId: user.ID,
				User: &management_pb.AddMachineUserRequest{
					UserName:    user.Username,
					Name:        user.Machine.Name,
					Description: user.Machine.Description,
				},
			})
		}

		metadataOrgSearch, err := query.NewUserMetadataResourceOwnerSearchQuery(org)
		if err != nil {
			return nil, nil, nil, err
		}
		metadataList, err := s.query.SearchUserMetadata(ctx, user.ID, false, &query.UserMetadataSearchQueries{Queries: []query.SearchQuery{metadataOrgSearch}})
		if err != nil {
			return nil, nil, nil, err
		}
		for _, metadata := range metadataList.Metadata {
			userMetadata = append(userMetadata, &management_pb.SetUserMetadataRequest{
				Id:    user.ID,
				Key:   metadata.Key,
				Value: metadata.Value,
			})
		}
	}
	return humanUsers, machineUsers, userMetadata, nil
}

func (s *Server) getTriggerActions(ctx context.Context, org string, processedActions []string) ([]*management_pb.SetTriggerActionsRequest, error) {
	flowTypes := []domain.FlowType{domain.FlowTypeExternalAuthentication}
	triggerActions := make([]*management_pb.SetTriggerActionsRequest, 0)

	for _, flowType := range flowTypes {
		flow, err := s.query.GetFlow(ctx, flowType, org)
		if err != nil {
			return nil, err
		}

		for triggerType, triggerAction := range flow.TriggerActions {
			actions := make([]string, 0)
			for _, action := range triggerAction {
				for _, actionID := range processedActions {
					if action.ID == actionID {
						actions = append(actions, action.ID)
					}
				}
			}

			triggerActions = append(triggerActions, &management_pb.SetTriggerActionsRequest{
				FlowType:    action_pb.FlowType(flowType),
				TriggerType: action_pb.TriggerType(triggerType),
				ActionIds:   actions,
			})
		}
	}
	return triggerActions, nil
}

func (s *Server) getActions(ctx context.Context, org string) ([]*admin_pb.DataAction, error) {
	actionSearch, err := query.NewActionResourceOwnerQuery(org)
	if err != nil {
		return nil, err
	}
	queriedActions, err := s.query.SearchActions(ctx, &query.ActionSearchQueries{Queries: []query.SearchQuery{actionSearch}})
	if err != nil && !caos_errors.IsNotFound(err) {
		return nil, err
	}
	actions := make([]*admin_pb.DataAction, len(queriedActions.Actions))
	if err != nil && caos_errors.IsNotFound(err) {
		return actions, nil
	}
	for i, action := range queriedActions.Actions {
		timeout := durationpb.New(action.Timeout)

		actions[i] = &admin_pb.DataAction{
			ActionId: action.ID,
			Action: &management_pb.CreateActionRequest{
				Name:          action.Name,
				Script:        action.Script,
				Timeout:       timeout,
				AllowedToFail: action.AllowedToFail,
			},
		}
	}

	return actions, nil
}

func (s *Server) getProjectsAndApps(ctx context.Context, org string) ([]*admin_pb.DataProject, []*management_pb.AddProjectRoleRequest, []*admin_pb.DataOIDCApplication, []*admin_pb.DataAPIApplication, error) {
	projectSearch, err := query.NewProjectResourceOwnerSearchQuery(org)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	queriedProjects, err := s.query.SearchProjects(ctx, &query.ProjectSearchQueries{Queries: []query.SearchQuery{projectSearch}})
	if err != nil && !caos_errors.IsNotFound(err) {
		return nil, nil, nil, nil, err
	}

	projects := make([]*admin_pb.DataProject, len(queriedProjects.Projects))
	orgProjectRoles := make([]*management_pb.AddProjectRoleRequest, 0)
	oidcApps := make([]*admin_pb.DataOIDCApplication, 0)
	apiApps := make([]*admin_pb.DataAPIApplication, 0)
	if err != nil && caos_errors.IsNotFound(err) {
		return projects, orgProjectRoles, oidcApps, apiApps, nil
	}
	for i, queriedProject := range queriedProjects.Projects {
		projects[i] = &admin_pb.DataProject{
			ProjectId: queriedProject.ID,
			Project: &management_pb.AddProjectRequest{
				Name:                   queriedProject.Name,
				ProjectRoleAssertion:   queriedProject.ProjectRoleAssertion,
				ProjectRoleCheck:       queriedProject.ProjectRoleCheck,
				HasProjectCheck:        queriedProject.HasProjectCheck,
				PrivateLabelingSetting: project_pb.PrivateLabelingSetting(queriedProject.PrivateLabelingSetting),
			},
		}

		projectRoleSearch, err := query.NewProjectRoleProjectIDSearchQuery(queriedProject.ID)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		queriedProjectRoles, err := s.query.SearchProjectRoles(ctx, false, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectRoleSearch}})
		if err != nil && !caos_errors.IsNotFound(err) {
			return nil, nil, nil, nil, err
		}
		if queriedProjectRoles != nil {
			for _, role := range queriedProjectRoles.ProjectRoles {
				orgProjectRoles = append(orgProjectRoles, &management_pb.AddProjectRoleRequest{
					ProjectId:   role.ProjectID,
					RoleKey:     role.Key,
					DisplayName: role.DisplayName,
					Group:       role.Group,
				})
			}
		}

		appSearch, err := query.NewAppProjectIDSearchQuery(queriedProject.ID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		apps, err := s.query.SearchApps(ctx, &query.AppSearchQueries{Queries: []query.SearchQuery{appSearch}})
		if err != nil && !caos_errors.IsNotFound(err) {
			return nil, nil, nil, nil, err
		}
		if apps != nil {
			for _, app := range apps.Apps {
				if app.OIDCConfig != nil {
					responseTypes := make([]app_pb.OIDCResponseType, 0)
					for _, ty := range app.OIDCConfig.ResponseTypes {
						responseTypes = append(responseTypes, app_pb.OIDCResponseType(ty))
					}

					grantTypes := make([]app_pb.OIDCGrantType, 0)
					for _, ty := range app.OIDCConfig.GrantTypes {
						grantTypes = append(grantTypes, app_pb.OIDCGrantType(ty))
					}

					oidcApps = append(oidcApps, &admin_pb.DataOIDCApplication{
						AppId: app.ID,
						App: &management_pb.AddOIDCAppRequest{
							ProjectId:                app.ProjectID,
							Name:                     app.Name,
							RedirectUris:             app.OIDCConfig.RedirectURIs,
							ResponseTypes:            responseTypes,
							GrantTypes:               grantTypes,
							AppType:                  app_pb.OIDCAppType(app.OIDCConfig.AppType),
							AuthMethodType:           app_pb.OIDCAuthMethodType(app.OIDCConfig.AuthMethodType),
							PostLogoutRedirectUris:   app.OIDCConfig.PostLogoutRedirectURIs,
							Version:                  app_pb.OIDCVersion(app.OIDCConfig.Version),
							DevMode:                  app.OIDCConfig.IsDevMode,
							AccessTokenType:          app_pb.OIDCTokenType(app.OIDCConfig.AccessTokenType),
							AccessTokenRoleAssertion: app.OIDCConfig.AssertAccessTokenRole,
							IdTokenRoleAssertion:     app.OIDCConfig.AssertIDTokenRole,
							IdTokenUserinfoAssertion: app.OIDCConfig.AssertIDTokenUserinfo,
							ClockSkew:                durationpb.New(app.OIDCConfig.ClockSkew),
							AdditionalOrigins:        app.OIDCConfig.AdditionalOrigins,
						},
					})
				}
				if app.APIConfig != nil {
					apiApps = append(apiApps, &admin_pb.DataAPIApplication{
						AppId: app.ID,
						App: &management_pb.AddAPIAppRequest{
							ProjectId:      app.ProjectID,
							Name:           app.Name,
							AuthMethodType: app_pb.APIAuthMethodType(app.APIConfig.AuthMethodType),
						},
					})
				}
			}
		}
	}
	return projects, orgProjectRoles, oidcApps, apiApps, nil
}

func (s *Server) getNecessaryProjectGrantMembersForOrg(ctx context.Context, org string, processedProjects []string, processedGrants []string, processedUsers []string) ([]*management_pb.AddProjectGrantMemberRequest, error) {
	projectMembers := make([]*management_pb.AddProjectGrantMemberRequest, 0)

	for _, projectID := range processedProjects {
		for _, grantID := range processedGrants {
			search, err := query.NewMemberResourceOwnerSearchQuery(org)
			if err != nil {
				return nil, err
			}

			queriedProjectMembers, err := s.query.ProjectGrantMembers(ctx, &query.ProjectGrantMembersQuery{ProjectID: projectID, OrgID: org, GrantID: grantID, MembersQuery: query.MembersQuery{Queries: []query.SearchQuery{search}}})
			if err != nil {
				return nil, err
			}
			for _, projectMember := range queriedProjectMembers.Members {
				for _, userID := range processedUsers {
					if userID == projectMember.UserID {

						projectMembers = append(projectMembers, &management_pb.AddProjectGrantMemberRequest{
							ProjectId: projectID,
							UserId:    userID,
							GrantId:   grantID,
							Roles:     projectMember.Roles,
						})
						break
					}
				}

			}
		}
	}
	return projectMembers, nil
}

func (s *Server) getNecessaryProjectMembersForOrg(ctx context.Context, processedProjects []string, processedUsers []string) ([]*management_pb.AddProjectMemberRequest, error) {
	projectMembers := make([]*management_pb.AddProjectMemberRequest, 0)

	for _, projectID := range processedProjects {
		queriedProjectMembers, err := s.query.ProjectMembers(ctx, &query.ProjectMembersQuery{ProjectID: projectID})
		if err != nil {
			return nil, err
		}
		for _, projectMember := range queriedProjectMembers.Members {
			for _, userID := range processedUsers {
				if userID == projectMember.UserID {
					projectMembers = append(projectMembers, &management_pb.AddProjectMemberRequest{
						ProjectId: projectID,
						UserId:    userID,
						Roles:     projectMember.Roles,
					})
					break
				}
			}
		}
	}
	return projectMembers, nil
}

func (s *Server) getNecessaryOrgMembersForOrg(ctx context.Context, org string, processedUsers []string) ([]*management_pb.AddOrgMemberRequest, error) {
	queriedOrgMembers, err := s.query.OrgMembers(ctx, &query.OrgMembersQuery{OrgID: org})
	if err != nil {
		return nil, err
	}
	orgMembers := make([]*management_pb.AddOrgMemberRequest, 0, len(queriedOrgMembers.Members))
	for _, orgMember := range queriedOrgMembers.Members {
		for _, userID := range processedUsers {
			if userID == orgMember.UserID {
				orgMembers = append(orgMembers, &management_pb.AddOrgMemberRequest{
					UserId: orgMember.UserID,
					Roles:  orgMember.Roles,
				})
				break
			}
		}
	}
	return orgMembers, nil
}

func (s *Server) getNecessaryProjectGrantsForOrg(ctx context.Context, org string, processedOrgs []string, processedProjects []string) ([]*admin_pb.DataProjectGrant, error) {

	projectGrantSearchOrg, err := query.NewProjectGrantResourceOwnerSearchQuery(org)
	if err != nil {
		return nil, err
	}
	queriedProjectGrants, err := s.query.SearchProjectGrants(ctx, &query.ProjectGrantSearchQueries{Queries: []query.SearchQuery{projectGrantSearchOrg}})
	if err != nil {
		return nil, err
	}
	projectGrants := make([]*admin_pb.DataProjectGrant, 0, len(queriedProjectGrants.ProjectGrants))
	for _, projectGrant := range queriedProjectGrants.ProjectGrants {
		for _, projectID := range processedProjects {
			if projectID == projectGrant.ProjectID {
				foundOrg := false
				for _, orgID := range processedOrgs {
					if orgID == projectGrant.GrantedOrgID {
						projectGrants = append(projectGrants, &admin_pb.DataProjectGrant{
							GrantId: projectGrant.GrantID,
							ProjectGrant: &management_pb.AddProjectGrantRequest{
								ProjectId:    projectGrant.ProjectID,
								GrantedOrgId: projectGrant.GrantedOrgID,
								RoleKeys:     projectGrant.GrantedRoleKeys,
							},
						})
						foundOrg = true
						break
					}
				}
				if foundOrg {
					break
				}
			}
		}
	}
	return projectGrants, nil
}

func (s *Server) getNecessaryUserGrantsForOrg(ctx context.Context, org string, processedProjects []string, processedGrants []string, processedUsers []string) ([]*management_pb.AddUserGrantRequest, error) {
	userGrantSearchOrg, err := query.NewUserGrantResourceOwnerSearchQuery(org)
	if err != nil {
		return nil, err
	}

	queriedUserGrants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{Queries: []query.SearchQuery{userGrantSearchOrg}})
	if err != nil {
		return nil, err
	}
	userGrants := make([]*management_pb.AddUserGrantRequest, 0, len(queriedUserGrants.UserGrants))
	for _, userGrant := range queriedUserGrants.UserGrants {
		for _, projectID := range processedProjects {
			if projectID == userGrant.ProjectID {
				//if usergrant is on a granted project
				if userGrant.GrantID != "" {
					for _, grantID := range processedGrants {
						if grantID == userGrant.GrantID {
							for _, userID := range processedUsers {
								if userID == userGrant.UserID {
									userGrants = append(userGrants, &management_pb.AddUserGrantRequest{
										UserId:         userGrant.UserID,
										ProjectId:      userGrant.ProjectID,
										ProjectGrantId: userGrant.GrantID,
										RoleKeys:       userGrant.Roles,
									})
								}
							}
						}
					}
				} else {
					for _, userID := range processedUsers {
						if userID == userGrant.UserID {
							userGrants = append(userGrants, &management_pb.AddUserGrantRequest{
								UserId:         userGrant.UserID,
								ProjectId:      userGrant.ProjectID,
								ProjectGrantId: userGrant.GrantID,
								RoleKeys:       userGrant.Roles,
							})
						}
					}
				}
			}
		}
	}
	return userGrants, nil
}
func (s *Server) getCustomLoginTexts(ctx context.Context, org string, languages []string) ([]*management_pb.SetCustomLoginTextsRequest, error) {
	customTexts := make([]*management_pb.SetCustomLoginTextsRequest, 0, len(languages))
	for _, lang := range languages {
		text, err := s.query.GetCustomLoginTexts(ctx, org, lang)
		if err != nil {
			return nil, err
		}
		if !text.IsDefault {
			customTexts = append(customTexts, &management_pb.SetCustomLoginTextsRequest{
				Language:                             lang,
				SelectAccountText:                    text_grpc.SelectAccountScreenToPb(text.SelectAccount),
				LoginText:                            text_grpc.LoginScreenTextToPb(text.Login),
				PasswordText:                         text_grpc.PasswordScreenTextToPb(text.Password),
				UsernameChangeText:                   text_grpc.UsernameChangeScreenTextToPb(text.UsernameChange),
				UsernameChangeDoneText:               text_grpc.UsernameChangeDoneScreenTextToPb(text.UsernameChangeDone),
				InitPasswordText:                     text_grpc.InitPasswordScreenTextToPb(text.InitPassword),
				InitPasswordDoneText:                 text_grpc.InitPasswordDoneScreenTextToPb(text.InitPasswordDone),
				EmailVerificationText:                text_grpc.EmailVerificationScreenTextToPb(text.EmailVerification),
				EmailVerificationDoneText:            text_grpc.EmailVerificationDoneScreenTextToPb(text.EmailVerificationDone),
				InitializeUserText:                   text_grpc.InitializeUserScreenTextToPb(text.InitUser),
				InitializeDoneText:                   text_grpc.InitializeUserDoneScreenTextToPb(text.InitUserDone),
				InitMfaPromptText:                    text_grpc.InitMFAPromptScreenTextToPb(text.InitMFAPrompt),
				InitMfaOtpText:                       text_grpc.InitMFAOTPScreenTextToPb(text.InitMFAOTP),
				InitMfaU2FText:                       text_grpc.InitMFAU2FScreenTextToPb(text.InitMFAU2F),
				InitMfaDoneText:                      text_grpc.InitMFADoneScreenTextToPb(text.InitMFADone),
				MfaProvidersText:                     text_grpc.MFAProvidersTextToPb(text.MFAProvider),
				VerifyMfaOtpText:                     text_grpc.VerifyMFAOTPScreenTextToPb(text.VerifyMFAOTP),
				VerifyMfaU2FText:                     text_grpc.VerifyMFAU2FScreenTextToPb(text.VerifyMFAU2F),
				PasswordlessText:                     text_grpc.PasswordlessScreenTextToPb(text.Passwordless),
				PasswordlessPromptText:               text_grpc.PasswordlessPromptScreenTextToPb(text.PasswordlessPrompt),
				PasswordlessRegistrationText:         text_grpc.PasswordlessRegistrationScreenTextToPb(text.PasswordlessRegistration),
				PasswordlessRegistrationDoneText:     text_grpc.PasswordlessRegistrationDoneScreenTextToPb(text.PasswordlessRegistrationDone),
				PasswordChangeText:                   text_grpc.PasswordChangeScreenTextToPb(text.PasswordChange),
				PasswordChangeDoneText:               text_grpc.PasswordChangeDoneScreenTextToPb(text.PasswordChangeDone),
				PasswordResetDoneText:                text_grpc.PasswordResetDoneScreenTextToPb(text.PasswordResetDone),
				RegistrationOptionText:               text_grpc.RegistrationOptionScreenTextToPb(text.RegisterOption),
				RegistrationUserText:                 text_grpc.RegistrationUserScreenTextToPb(text.RegistrationUser),
				ExternalRegistrationUserOverviewText: text_grpc.ExternalRegistrationUserOverviewScreenTextToPb(text.ExternalRegistrationUserOverview),
				RegistrationOrgText:                  text_grpc.RegistrationOrgScreenTextToPb(text.RegistrationOrg),
				LinkingUserDoneText:                  text_grpc.LinkingUserDoneScreenTextToPb(text.LinkingUsersDone),
				ExternalUserNotFoundText:             text_grpc.ExternalUserNotFoundScreenTextToPb(text.ExternalNotFoundOption),
				SuccessLoginText:                     text_grpc.SuccessLoginScreenTextToPb(text.LoginSuccess),
				LogoutText:                           text_grpc.LogoutDoneScreenTextToPb(text.LogoutDone),
				FooterText:                           text_grpc.FooterTextToPb(text.Footer),
			})
		}
	}

	return customTexts, nil
}

func (s *Server) getCustomInitMessageTexts(ctx context.Context, org string, languages []string) ([]*management_pb.SetCustomInitMessageTextRequest, error) {
	customTexts := make([]*management_pb.SetCustomInitMessageTextRequest, 0, len(languages))
	for _, lang := range languages {
		text, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, org, domain.InitCodeMessageType, lang)
		if err != nil {
			return nil, err
		}

		if !text.IsDefault {
			customTexts = append(customTexts, &management_pb.SetCustomInitMessageTextRequest{
				Language:   lang,
				Title:      text.Title,
				PreHeader:  text.PreHeader,
				Subject:    text.Subject,
				Greeting:   text.Greeting,
				Text:       text.Text,
				ButtonText: text.ButtonText,
				FooterText: text.Footer,
			})
		}
	}

	return customTexts, nil
}

func (s *Server) getCustomPasswordResetMessageTexts(ctx context.Context, org string, languages []string) ([]*management_pb.SetCustomPasswordResetMessageTextRequest, error) {
	customTexts := make([]*management_pb.SetCustomPasswordResetMessageTextRequest, 0, len(languages))
	for _, lang := range languages {
		text, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, org, domain.PasswordResetMessageType, lang)
		if err != nil {
			return nil, err
		}

		if !text.IsDefault {
			customTexts = append(customTexts, &management_pb.SetCustomPasswordResetMessageTextRequest{
				Language:   lang,
				Title:      text.Title,
				PreHeader:  text.PreHeader,
				Subject:    text.Subject,
				Greeting:   text.Greeting,
				Text:       text.Text,
				ButtonText: text.ButtonText,
				FooterText: text.Footer,
			})
		}
	}

	return customTexts, nil
}

func (s *Server) getCustomVerifyEmailMessageTexts(ctx context.Context, org string, languages []string) ([]*management_pb.SetCustomVerifyEmailMessageTextRequest, error) {
	customTexts := make([]*management_pb.SetCustomVerifyEmailMessageTextRequest, 0, len(languages))
	for _, lang := range languages {
		text, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, org, domain.VerifyEmailMessageType, lang)
		if err != nil {
			return nil, err
		}

		if !text.IsDefault {
			customTexts = append(customTexts, &management_pb.SetCustomVerifyEmailMessageTextRequest{
				Language:   lang,
				Title:      text.Title,
				PreHeader:  text.PreHeader,
				Subject:    text.Subject,
				Greeting:   text.Greeting,
				Text:       text.Text,
				ButtonText: text.ButtonText,
				FooterText: text.Footer,
			})
		}
	}

	return customTexts, nil
}

func (s *Server) getCustomVerifyPhoneMessageTexts(ctx context.Context, org string, languages []string) ([]*management_pb.SetCustomVerifyPhoneMessageTextRequest, error) {
	customTexts := make([]*management_pb.SetCustomVerifyPhoneMessageTextRequest, 0, len(languages))
	for _, lang := range languages {
		text, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, org, domain.VerifyPhoneMessageType, lang)
		if err != nil {
			return nil, err
		}

		if !text.IsDefault {
			customTexts = append(customTexts, &management_pb.SetCustomVerifyPhoneMessageTextRequest{
				Language:   lang,
				Title:      text.Title,
				PreHeader:  text.PreHeader,
				Subject:    text.Subject,
				Greeting:   text.Greeting,
				Text:       text.Text,
				ButtonText: text.ButtonText,
				FooterText: text.Footer,
			})
		}
	}

	return customTexts, nil
}

func (s *Server) getCustomDomainClaimedMessageTexts(ctx context.Context, org string, languages []string) ([]*management_pb.SetCustomDomainClaimedMessageTextRequest, error) {
	customTexts := make([]*management_pb.SetCustomDomainClaimedMessageTextRequest, 0, len(languages))
	for _, lang := range languages {
		text, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, org, domain.DomainClaimedMessageType, lang)
		if err != nil {
			return nil, err
		}

		if !text.IsDefault {
			customTexts = append(customTexts, &management_pb.SetCustomDomainClaimedMessageTextRequest{
				Language:   lang,
				Title:      text.Title,
				PreHeader:  text.PreHeader,
				Subject:    text.Subject,
				Greeting:   text.Greeting,
				Text:       text.Text,
				ButtonText: text.ButtonText,
				FooterText: text.Footer,
			})
		}
	}

	return customTexts, nil
}

func (s *Server) getCustomPasswordlessRegistrationMessageTexts(ctx context.Context, org string, languages []string) ([]*management_pb.SetCustomPasswordlessRegistrationMessageTextRequest, error) {
	customTexts := make([]*management_pb.SetCustomPasswordlessRegistrationMessageTextRequest, 0, len(languages))
	for _, lang := range languages {
		text, err := s.query.CustomMessageTextByTypeAndLanguage(ctx, org, domain.DomainClaimedMessageType, lang)
		if err != nil {
			return nil, err
		}

		if !text.IsDefault {
			customTexts = append(customTexts, &management_pb.SetCustomPasswordlessRegistrationMessageTextRequest{
				Language:   lang,
				Title:      text.Title,
				PreHeader:  text.PreHeader,
				Subject:    text.Subject,
				Greeting:   text.Greeting,
				Text:       text.Text,
				ButtonText: text.ButtonText,
				FooterText: text.Footer,
			})
		}
	}

	return customTexts, nil
}
