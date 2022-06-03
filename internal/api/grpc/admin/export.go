package admin

import (
	"context"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	app_pb "github.com/zitadel/zitadel/pkg/grpc/app"
	management_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
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
	processedUsers := make([]string, 0)

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

		/******************************************************************************************************************
		Users
		******************************************************************************************************************/
		humanUsers, machineUsers, err := s.getUsers(ctx, queriedOrg.ID)
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
		orgProjects, oidcApps, apiApps, err := s.getProjectsAndApps(ctx, queriedOrg.ID)
		if err != nil {
			return nil, err
		}
		org.Projects = orgProjects
		org.OidcApps = oidcApps
		org.ApiApps = apiApps

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

		orgs = append(orgs, org)
	}

	for _, org := range orgs {

		/******************************************************************************************************************
		Grants
		******************************************************************************************************************/
		projectGrants, err := s.getNecessaryProjectGrantsForOrg(ctx, org.OrgId, processedOrgs, processedProjects)
		if err != nil {
			return nil, err
		}
		org.ProjectGrants = projectGrants

		userGrants, err := s.getNecessaryUserGrantsForOrg(ctx, org.OrgId, processedProjects, processedUsers)
		if err != nil {
			return nil, err
		}
		org.UserGrants = userGrants

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

		projectGrantMembers, err := s.getNecessaryProjectGrantMembersForOrg(ctx, org.OrgId, processedProjects, processedUsers)
		if err != nil {
			return nil, err
		}
		org.ProjectGrantMembers = projectGrantMembers
	}

	return &admin_pb.ExportDataResponse{
		Orgs: orgs,
	}, nil
}

func (s *Server) getUsers(ctx context.Context, org string) ([]*admin_pb.DataHumanUser, []*admin_pb.DataMachineUser, error) {
	orgSearch, err := query.NewUserResourceOwnerSearchQuery(org, query.TextEquals)
	if err != nil {
		return nil, nil, err
	}
	users, err := s.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: []query.SearchQuery{orgSearch}})
	if err != nil {
		return nil, nil, err
	}
	humanUsers := make([]*admin_pb.DataHumanUser, 0)
	machineUsers := make([]*admin_pb.DataMachineUser, 0)
	for _, user := range users.Users {
		switch user.Type {
		case domain.UserTypeHuman:
			dataUser := &admin_pb.DataHumanUser{
				UserId: user.ID,
				User: &management_pb.AddHumanUserRequest{
					UserName: user.Username,
					Profile: &management_pb.AddHumanUserRequest_Profile{
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
				dataUser.User.Email = &management_pb.AddHumanUserRequest_Email{
					Email:           user.Human.Email,
					IsEmailVerified: user.Human.IsEmailVerified,
				}
			}
			if user.Human.Phone != "" {
				dataUser.User.Phone = &management_pb.AddHumanUserRequest_Phone{
					Phone:           user.Human.Phone,
					IsPhoneVerified: user.Human.IsPhoneVerified,
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
	}
	return humanUsers, machineUsers, nil
}

func (s *Server) getActions(ctx context.Context, org string) ([]*admin_pb.DataAction, error) {
	actionSearch, err := query.NewActionResourceOwnerQuery(org)
	if err != nil {
		return nil, err
	}
	queriedActions, err := s.query.SearchActions(ctx, &query.ActionSearchQueries{Queries: []query.SearchQuery{actionSearch}})
	if err != nil {
		return nil, err
	}
	actions := make([]*admin_pb.DataAction, 0)
	for _, action := range queriedActions.Actions {
		timeout := durationpb.New(action.Timeout)

		actions = append(actions, &admin_pb.DataAction{
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

func (s *Server) getProjectsAndApps(ctx context.Context, org string) ([]*admin_pb.DataProject, []*admin_pb.DataOIDCApplication, []*admin_pb.DataAPIApplication, error) {
	projectSearch, err := query.NewProjectResourceOwnerSearchQuery(org)
	if err != nil {
		return nil, nil, nil, err
	}
	projects, err := s.query.SearchProjects(ctx, &query.ProjectSearchQueries{Queries: []query.SearchQuery{projectSearch}})
	if err != nil {
		return nil, nil, nil, err
	}

	orgProjects := make([]*admin_pb.DataProject, 0)
	oidcApps := make([]*admin_pb.DataOIDCApplication, 0)
	apiApps := make([]*admin_pb.DataAPIApplication, 0)
	for _, project := range projects.Projects {
		setting := project_pb.PrivateLabelingSetting(project.PrivateLabelingSetting)
		orgProjects = append(orgProjects, &admin_pb.DataProject{
			ProjectId: project.ID,
			Project: &management_pb.AddProjectRequest{
				Name:                   project.Name,
				ProjectRoleAssertion:   project.ProjectRoleAssertion,
				ProjectRoleCheck:       project.ProjectRoleCheck,
				HasProjectCheck:        project.HasProjectCheck,
				PrivateLabelingSetting: setting,
			},
		})

		appSearch, err := query.NewAppProjectIDSearchQuery(project.ID)
		if err != nil {
			return nil, nil, nil, err
		}
		apps, err := s.query.SearchApps(ctx, &query.AppSearchQueries{Queries: []query.SearchQuery{appSearch}})
		if err != nil {
			return nil, nil, nil, err
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
						ClockSkew:                duration,
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
	return orgProjects, oidcApps, apiApps, nil
}

func (s *Server) getNecessaryProjectGrantMembersForOrg(ctx context.Context, org string, processedProjects []string, processedUsers []string) ([]*management_pb.AddProjectGrantMemberRequest, error) {
	projectMembers := make([]*management_pb.AddProjectGrantMemberRequest, 0)

	for _, projectID := range processedProjects {
		queriedProjectMembers, err := s.query.ProjectGrantMembers(ctx, &query.ProjectGrantMembersQuery{ProjectID: projectID, OrgID: org})
		if err != nil {
			return nil, err
		}
		for _, projectMember := range queriedProjectMembers.Members {
			for _, userID := range processedUsers {
				if userID == projectMember.UserID {
					projectMembers = append(projectMembers, &management_pb.AddProjectGrantMemberRequest{
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

func (s *Server) getNecessaryProjectGrantsForOrg(ctx context.Context, org string, processedOrgs []string, processedProjects []string) ([]*management_pb.AddProjectGrantRequest, error) {
	projectGrants := make([]*management_pb.AddProjectGrantRequest, 0)
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
						projectGrants = append(projectGrants, &management_pb.AddProjectGrantRequest{
							ProjectId:    projectGrant.ProjectID,
							GrantedOrgId: projectGrant.GrantedOrgID,
							RoleKeys:     projectGrant.GrantedRoleKeys,
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

func (s *Server) getNecessaryUserGrantsForOrg(ctx context.Context, org string, processedProjects []string, processedUsers []string) ([]*management_pb.AddUserGrantRequest, error) {
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
				foundUser := false
				for _, userID := range processedUsers {
					if userID == userGrant.UserID {
						userGrants = append(userGrants, &management_pb.AddUserGrantRequest{
							UserId:    userGrant.UserID,
							ProjectId: userGrant.ProjectID,
							RoleKeys:  userGrant.Roles,
						})
						foundUser = true
						break
					}
				}
				if foundUser {
					break
				}
			}
		}
	}
	return userGrants, nil
}
