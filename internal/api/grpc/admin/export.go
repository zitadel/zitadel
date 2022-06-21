package admin

import (
	"context"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	action_pb "github.com/zitadel/zitadel/pkg/grpc/action"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	app_pb "github.com/zitadel/zitadel/pkg/grpc/app"
	management_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
	v1_pb "github.com/zitadel/zitadel/pkg/grpc/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (s *Server) ExportData(ctx context.Context, req *admin_pb.ExportDataRequest) (*admin_pb.ExportDataResponse, error) {
	queriedOrgs, err := s.query.SearchOrgs(ctx, &query.OrgSearchQueries{})
	if err != nil {
		return nil, err
	}

	orgs := make([]*admin_pb.DataOrg, 0)
	processedOrgs := make([]string, 0)
	processedProjects := make([]string, 0)
	processedGrants := make([]string, 0)
	processedUsers := make([]string, 0)
	processedActions := make([]string, 0)

	for _, queriedOrg := range queriedOrgs.Orgs {
		if req.OrgIds != nil || len(req.OrgIds) > 0 {
			found := false
			for _, searchingOrg := range req.OrgIds {
				if queriedOrg.ID == searchingOrg {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		processedOrgs = append(processedOrgs, queriedOrg.ID)

		/******************************************************************************************************************
		Organization
		******************************************************************************************************************/
		queriedOrg, err := s.query.OrgByID(ctx, queriedOrg.ID)
		if err != nil {
			return nil, err
		}
		org := &admin_pb.DataOrg{OrgId: queriedOrg.ID, Org: &management_pb.AddOrgRequest{Name: queriedOrg.Name}}

		queriedDomain, err := s.query.DomainPolicyByOrg(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}
		if !queriedDomain.IsDefault {
			org.DomainPolicy = &admin_pb.AddCustomDomainPolicyRequest{
				OrgId:                                  org.GetOrgId(),
				UserLoginMustBeDomain:                  queriedDomain.UserLoginMustBeDomain,
				ValidateOrgDomains:                     queriedDomain.ValidateOrgDomains,
				SmtpSenderAddressMatchesInstanceDomain: queriedDomain.SMTPSenderAddressMatchesInstanceDomain,
			}
		}

		queriedLabel, err := s.query.ActiveLabelPolicyByOrg(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}
		if !queriedLabel.IsDefault {
			org.LabelPolicy = &management_pb.AddCustomLabelPolicyRequest{
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
			}
		}

		queriedLogin, err := s.query.LoginPolicyByID(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}
		if !queriedLogin.IsDefault {
			pwCheck := durationpb.New(queriedLogin.PasswordCheckLifetime)
			externalLogin := durationpb.New(queriedLogin.ExternalLoginCheckLifetime)
			mfaInitSkip := durationpb.New(queriedLogin.MFAInitSkipLifetime)
			secondFactor := durationpb.New(queriedLogin.SecondFactorCheckLifetime)
			multiFactor := durationpb.New(queriedLogin.MultiFactorCheckLifetime)

			secondFactors := []policy_pb.SecondFactorType{}
			for _, factor := range queriedLogin.SecondFactors {
				secondFactors = append(secondFactors, policy_pb.SecondFactorType(factor))
			}

			multiFactors := []policy_pb.MultiFactorType{}
			for _, factor := range queriedLogin.MultiFactors {
				multiFactors = append(multiFactors, policy_pb.MultiFactorType(factor))
			}

			org.LoginPolicy = &management_pb.AddCustomLoginPolicyRequest{
				AllowUsernamePassword:      queriedLogin.AllowUsernamePassword,
				AllowRegister:              queriedLogin.AllowRegister,
				AllowExternalIdp:           queriedLogin.AllowExternalIDPs,
				ForceMfa:                   queriedLogin.ForceMFA,
				PasswordlessType:           policy_pb.PasswordlessType(queriedLogin.PasswordlessType),
				HidePasswordReset:          queriedLogin.HidePasswordReset,
				IgnoreUnknownUsernames:     queriedLogin.IgnoreUnknownUsernames,
				DefaultRedirectUri:         queriedLogin.DefaultRedirectURI,
				PasswordCheckLifetime:      pwCheck,
				ExternalLoginCheckLifetime: externalLogin,
				MfaInitSkipLifetime:        mfaInitSkip,
				SecondFactorCheckLifetime:  secondFactor,
				MultiFactorCheckLifetime:   multiFactor,
				SecondFactors:              secondFactors,
				MultiFactors:               multiFactors,
				// TODO ???
				//Idps:                       queriedLogin.id,
			}
		}

		queriedLockout, err := s.query.LockoutPolicyByOrg(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}
		if !queriedLockout.IsDefault {
			org.LockoutPolicy = &management_pb.AddCustomLockoutPolicyRequest{
				MaxPasswordAttempts: uint32(queriedLockout.MaxPasswordAttempts),
			}
		}

		queriedPasswordComplexity, err := s.query.PasswordComplexityPolicyByOrg(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}
		if !queriedPasswordComplexity.IsDefault {
			org.PasswordComplexityPolicy = &management_pb.AddCustomPasswordComplexityPolicyRequest{
				MinLength:    queriedPasswordComplexity.MinLength,
				HasUppercase: queriedPasswordComplexity.HasUppercase,
				HasLowercase: queriedPasswordComplexity.HasLowercase,
				HasNumber:    queriedPasswordComplexity.HasNumber,
				HasSymbol:    queriedPasswordComplexity.HasSymbol,
			}
		}

		queriedPrivacy, err := s.query.PrivacyPolicyByOrg(ctx, org.GetOrgId())
		if err != nil {
			return nil, err
		}
		if !queriedPrivacy.IsDefault {
			org.PrivacyPolicy = &management_pb.AddCustomPrivacyPolicyRequest{
				TosLink:     queriedPrivacy.TOSLink,
				PrivacyLink: queriedPrivacy.PrivacyLink,
				HelpLink:    queriedPrivacy.HelpLink,
			}
		}

		/******************************************************************************************************************
		Users
		******************************************************************************************************************/
		humanUsers, machineUsers, err := s.getUsers(ctx, queriedOrg.ID, req.WithPasswords)
		if err != nil {
			return nil, err
		}
		org.HumanUsers = humanUsers
		org.MachineUsers = machineUsers
		for _, processedUser := range humanUsers {
			processedUsers = append(processedUsers, processedUser.UserId)
		}
		for _, processedUser := range machineUsers {
			processedUsers = append(processedUsers, processedUser.UserId)
		}

		/******************************************************************************************************************
		Project and Applications
		******************************************************************************************************************/
		orgProjects, orgProjectRoles, oidcApps, apiApps, err := s.getProjectsAndApps(ctx, queriedOrg.ID)
		if err != nil {
			return nil, err
		}
		org.Projects = orgProjects
		org.OidcApps = oidcApps
		org.ApiApps = apiApps
		org.ProjectRoles = orgProjectRoles

		for _, processedProject := range orgProjects {
			processedProjects = append(processedProjects, processedProject.ProjectId)
		}

		/******************************************************************************************************************
		Actions
		******************************************************************************************************************/
		actions, err := s.getActions(ctx, queriedOrg.ID)
		if err != nil {
			return nil, err
		}
		org.Actions = actions
		for _, processedAction := range actions {
			processedActions = append(processedActions, processedAction.ActionId)
		}

		orgs = append(orgs, org)
	}

	for _, org := range orgs {
		/******************************************************************************************************************
		Flows
		******************************************************************************************************************/
		triggerActions, err := s.getTriggerActions(ctx, org.OrgId, processedActions)
		if err != nil {
			return nil, err
		}
		org.TriggerActions = triggerActions

		/******************************************************************************************************************
		Grants
		******************************************************************************************************************/
		projectGrants, err := s.getNecessaryProjectGrantsForOrg(ctx, org.OrgId, processedOrgs, processedProjects)
		if err != nil {
			return nil, err
		}
		org.ProjectGrants = projectGrants
		for _, processedGrant := range projectGrants {
			processedGrants = append(processedGrants, processedGrant.GrantId)
		}

		userGrants, err := s.getNecessaryUserGrantsForOrg(ctx, org.OrgId, processedProjects, processedGrants, processedUsers)
		if err != nil {
			return nil, err
		}
		org.UserGrants = userGrants
	}

	for _, org := range orgs {
		/******************************************************************************************************************
		Members
		******************************************************************************************************************/
		orgMembers, err := s.getNecessaryOrgMembersForOrg(ctx, org.OrgId, processedUsers)
		if err != nil {
			return nil, err
		}
		org.OrgMembers = orgMembers

		projectMembers, err := s.getNecessaryProjectMembersForOrg(ctx, processedProjects, processedUsers)
		if err != nil {
			return nil, err
		}
		org.ProjectMembers = projectMembers

		projectGrantMembers, err := s.getNecessaryProjectGrantMembersForOrg(ctx, org.OrgId, processedProjects, processedGrants, processedUsers)
		if err != nil {
			return nil, err
		}
		org.ProjectGrantMembers = projectGrantMembers
	}

	return &admin_pb.ExportDataResponse{
		Orgs: orgs,
	}, nil
}

func (s *Server) getUsers(ctx context.Context, org string, withPasswords bool) ([]*v1_pb.DataHumanUser, []*v1_pb.DataMachineUser, error) {
	orgSearch, err := query.NewUserResourceOwnerSearchQuery(org, query.TextEquals)
	if err != nil {
		return nil, nil, err
	}
	users, err := s.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: []query.SearchQuery{orgSearch}})
	if err != nil {
		return nil, nil, err
	}
	humanUsers := make([]*v1_pb.DataHumanUser, 0)
	machineUsers := make([]*v1_pb.DataMachineUser, 0)
	for _, user := range users.Users {
		switch user.Type {
		case domain.UserTypeHuman:
			dataUser := &v1_pb.DataHumanUser{
				UserId: user.ID,
				User: &management_pb.ImportHumanUserRequest{
					UserName: user.Username,
					Profile: &management_pb.ImportHumanUserRequest_Profile{
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
				dataUser.User.Email = &management_pb.ImportHumanUserRequest_Email{
					Email:           user.Human.Email,
					IsEmailVerified: user.Human.IsEmailVerified,
				}
			}
			if user.Human.Phone != "" {
				dataUser.User.Phone = &management_pb.ImportHumanUserRequest_Phone{
					Phone:           user.Human.Phone,
					IsPhoneVerified: user.Human.IsPhoneVerified,
				}
			}
			if withPasswords {
				hashedPassword, hashAlgorithm, err := s.query.GetHumanPassword(ctx, org, user.ID)
				if err != nil {
					return nil, nil, err
				}
				if hashedPassword != nil {
					dataUser.User.HashedPassword = &management_pb.ImportHumanUserRequest_HashedPassword{
						Value:     string(hashedPassword),
						Algorithm: hashAlgorithm,
					}
				}
			}

			humanUsers = append(humanUsers, dataUser)
		case domain.UserTypeMachine:
			machineUsers = append(machineUsers, &v1_pb.DataMachineUser{
				UserId: user.ID,
				User: &management_pb.AddMachineUserRequest{
					UserName:    user.Username,
					Name:        user.Machine.Name,
					Description: user.Machine.Description,
				},
			})
		}
	}
	return humanUsers, machineUsers, nil
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

func (s *Server) getActions(ctx context.Context, org string) ([]*v1_pb.DataAction, error) {
	actionSearch, err := query.NewActionResourceOwnerQuery(org)
	if err != nil {
		return nil, err
	}
	queriedActions, err := s.query.SearchActions(ctx, &query.ActionSearchQueries{Queries: []query.SearchQuery{actionSearch}})
	if err != nil {
		return nil, err
	}
	actions := make([]*v1_pb.DataAction, 0)
	for _, action := range queriedActions.Actions {
		timeout := durationpb.New(action.Timeout)

		actions = append(actions, &v1_pb.DataAction{
			ActionId: action.ID,
			Action: &management_pb.CreateActionRequest{
				Name:          action.Name,
				Script:        action.Script,
				Timeout:       timeout,
				AllowedToFail: action.AllowedToFail,
			},
		})
	}

	return actions, nil
}

func (s *Server) getProjectsAndApps(ctx context.Context, org string) ([]*v1_pb.DataProject, []*management_pb.AddProjectRoleRequest, []*v1_pb.DataOIDCApplication, []*v1_pb.DataAPIApplication, error) {
	projectSearch, err := query.NewProjectResourceOwnerSearchQuery(org)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	projects, err := s.query.SearchProjects(ctx, &query.ProjectSearchQueries{Queries: []query.SearchQuery{projectSearch}})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	orgProjects := make([]*v1_pb.DataProject, 0)
	orgProjectRoles := make([]*management_pb.AddProjectRoleRequest, 0)
	oidcApps := make([]*v1_pb.DataOIDCApplication, 0)
	apiApps := make([]*v1_pb.DataAPIApplication, 0)
	for _, project := range projects.Projects {
		setting := project_pb.PrivateLabelingSetting(project.PrivateLabelingSetting)
		orgProjects = append(orgProjects, &v1_pb.DataProject{
			ProjectId: project.ID,
			Project: &management_pb.AddProjectRequest{
				Name:                   project.Name,
				ProjectRoleAssertion:   project.ProjectRoleAssertion,
				ProjectRoleCheck:       project.ProjectRoleCheck,
				HasProjectCheck:        project.HasProjectCheck,
				PrivateLabelingSetting: setting,
			},
		})

		projectRoleSearch, err := query.NewProjectRoleProjectIDSearchQuery(project.ID)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		queriedProjectRoles, err := s.query.SearchProjectRoles(ctx, &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectRoleSearch}})
		if err != nil {
			return nil, nil, nil, nil, err
		}
		for _, role := range queriedProjectRoles.ProjectRoles {
			orgProjectRoles = append(orgProjectRoles, &management_pb.AddProjectRoleRequest{
				ProjectId:   role.ProjectID,
				RoleKey:     role.Key,
				DisplayName: role.DisplayName,
				Group:       role.Group,
			})
		}

		appSearch, err := query.NewAppProjectIDSearchQuery(project.ID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		apps, err := s.query.SearchApps(ctx, &query.AppSearchQueries{Queries: []query.SearchQuery{appSearch}})
		if err != nil {
			return nil, nil, nil, nil, err
		}
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
				duration := durationpb.New(app.OIDCConfig.ClockSkew)

				oidcApps = append(oidcApps, &v1_pb.DataOIDCApplication{
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
						ClockSkew:                duration,
						AdditionalOrigins:        app.OIDCConfig.AdditionalOrigins,
					},
				})
			}
			if app.APIConfig != nil {
				apiApps = append(apiApps, &v1_pb.DataAPIApplication{
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
	return orgProjects, orgProjectRoles, oidcApps, apiApps, nil
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
	orgMembers := make([]*management_pb.AddOrgMemberRequest, 0)
	queriedOrgMembers, err := s.query.OrgMembers(ctx, &query.OrgMembersQuery{OrgID: org})
	if err != nil {
		return nil, err
	}
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

func (s *Server) getNecessaryProjectGrantsForOrg(ctx context.Context, org string, processedOrgs []string, processedProjects []string) ([]*v1_pb.DataProjectGrant, error) {
	projectGrants := make([]*v1_pb.DataProjectGrant, 0)
	projectGrantSearchOrg, err := query.NewProjectGrantResourceOwnerSearchQuery(org)
	if err != nil {
		return nil, err
	}
	queriedProjectGrants, err := s.query.SearchProjectGrants(ctx, &query.ProjectGrantSearchQueries{Queries: []query.SearchQuery{projectGrantSearchOrg}})
	if err != nil {
		return nil, err
	}
	for _, projectGrant := range queriedProjectGrants.ProjectGrants {
		for _, projectID := range processedProjects {
			if projectID == projectGrant.ProjectID {
				foundOrg := false
				for _, orgID := range processedOrgs {
					if orgID == projectGrant.GrantedOrgID {
						projectGrants = append(projectGrants, &v1_pb.DataProjectGrant{
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
	userGrants := make([]*management_pb.AddUserGrantRequest, 0)

	userGrantSearchOrg, err := query.NewUserGrantResourceOwnerSearchQuery(org)
	if err != nil {
		return nil, err
	}

	queriedUserGrants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{Queries: []query.SearchQuery{userGrantSearchOrg}})
	if err != nil {
		return nil, err
	}
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
