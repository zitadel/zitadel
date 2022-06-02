package admin

import (
	"context"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/management"
	"github.com/zitadel/zitadel/internal/domain"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ImportData(ctx context.Context, req *admin_pb.ImportDataRequest) (*admin_pb.ImportDataResponse, error) {
	errors := make([]*admin_pb.ImportDataError, 0)
	success := &admin_pb.ImportDataSuccess{}

	appSecretGenerator, err := s.query.InitHashGenerator(ctx, domain.SecretGeneratorTypeAppSecret, s.passwordHashAlg)
	if err != nil {
		return nil, err
	}

	ctxData := authz.GetCtxData(ctx)

	for _, org := range req.GetOrgs() {
		_, err := s.command.AddOrgWithID(ctx, org.GetOrg().GetName(), ctxData.UserID, ctxData.ResourceOwner, org.GetOrgId(), []string{})
		if err != nil {
			errors = append(errors, &admin_pb.ImportDataError{Type: "org", Id: org.GetOrgId(), Message: err.Error()})
			continue
		}
		successOrg := &admin_pb.ImportDataSuccessOrg{OrgId: org.GetOrgId()}

		domainPolicy := org.GetDomainPolicy()
		if org.DomainPolicy != nil {
			_, err := s.command.AddOrgDomainPolicy(ctx, org.GetOrgId(), DomainPolicyToDomain(domainPolicy.UserLoginMustBeDomain, domainPolicy.ValidateOrgDomains, domainPolicy.SmtpSenderAddressMatchesInstanceDomain))
			if err != nil {
				errors = append(errors, &admin_pb.ImportDataError{Type: "domain_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
		}
		if org.LabelPolicy != nil {
			_, err = s.command.AddLabelPolicy(ctx, org.GetOrgId(), management.AddLabelPolicyToDomain(org.GetLabelPolicy()))
			if err != nil {
				errors = append(errors, &admin_pb.ImportDataError{Type: "label_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
		}
		if org.LockoutPolicy != nil {
			_, err = s.command.AddLockoutPolicy(ctx, org.GetOrgId(), management.AddLockoutPolicyToDomain(org.GetLockoutPolicy()))
			if err != nil {
				errors = append(errors, &admin_pb.ImportDataError{Type: "lockout_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
		}
		if org.LoginPolicy != nil {
			_, err = s.command.AddLoginPolicy(ctx, org.GetOrgId(), management.AddLoginPolicyToDomain(org.GetLoginPolicy()))
			if err != nil {
				errors = append(errors, &admin_pb.ImportDataError{Type: "login_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
		}
		if org.PasswordComplexityPolicy != nil {
			_, err = s.command.AddPasswordComplexityPolicy(ctx, org.GetOrgId(), management.AddPasswordComplexityPolicyToDomain(org.GetPasswordComplexityPolicy()))
			if err != nil {
				errors = append(errors, &admin_pb.ImportDataError{Type: "password_complexity_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
		}
		if org.PrivacyPolicy != nil {
			_, err = s.command.AddPrivacyPolicy(ctx, org.GetOrgId(), management.AddPrivacyPolicyToDomain(org.GetPrivacyPolicy()))
			if err != nil {
				errors = append(errors, &admin_pb.ImportDataError{Type: "privacy_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
		}
		if org.HumanUsers != nil {
			for _, user := range org.GetHumanUsers() {
				_, err := s.command.AddHumanWithID(ctx, org.GetOrgId(), user.GetUserId(), management.AddHumanUserRequestToAddHuman(user.GetUser()))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "human_user", Id: user.GetUserId(), Message: err.Error()})
					continue
				}
				successOrg.HumanUserIds = append(successOrg.HumanUserIds, user.GetUserId())
			}
		}
		if org.MachineUsers != nil {
			for _, user := range org.GetMachineUsers() {
				_, err := s.command.AddMachineWithID(ctx, org.GetOrgId(), user.GetUserId(), management.AddMachineUserRequestToDomain(user.GetUser()))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "machine_user", Id: user.GetUserId(), Message: err.Error()})
					continue
				}
				successOrg.MachinUserIds = append(successOrg.MachinUserIds, user.GetUserId())
			}
		}
		if org.Projects != nil {
			for _, project := range org.GetProjects() {
				_, err := s.command.AddProjectWithID(ctx, management.ProjectCreateToDomain(project.GetProject()), org.GetOrgId(), project.GetOwnerId(), project.GetProjectId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "project", Id: project.GetProjectId(), Message: err.Error()})
					continue
				}
				successOrg.ProjectIds = append(successOrg.ProjectIds, project.GetProjectId())
			}
		}
		if org.OidcApps != nil {
			for _, app := range org.GetOidcApps() {
				_, err := s.command.AddOIDCApplicationWithID(ctx, management.AddOIDCAppRequestToDomain(app.App), org.GetOrgId(), app.GetAppId(), appSecretGenerator)
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "oidc_app", Id: app.GetAppId(), Message: err.Error()})
					continue
				}
				successOrg.OidcAppIds = append(successOrg.OidcAppIds, app.GetAppId())
			}
		}
		if org.ApiApps != nil {
			for _, app := range org.GetApiApps() {
				_, err := s.command.AddAPIApplicationWithID(ctx, management.AddAPIAppRequestToDomain(app.GetApp()), org.GetOrgId(), app.GetAppId(), appSecretGenerator)
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "api_app", Id: app.GetAppId(), Message: err.Error()})
					continue
				}
				successOrg.ApiAppIds = append(successOrg.ApiAppIds, app.GetAppId())
			}
		}
		if org.Actions != nil {
			for _, action := range org.GetActions() {
				_, _, err := s.command.AddActionWithID(ctx, management.CreateActionRequestToDomain(action.GetAction()), org.GetOrgId(), action.GetActionId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "action", Id: action.GetActionId(), Message: err.Error()})
					continue
				}
				successOrg.ActionIds = append(successOrg.ActionIds, action.ActionId)
			}
		}
		if org.ProjectGrants != nil {
			for _, grant := range org.GetProjectGrants() {
				_, err := s.command.AddProjectGrant(ctx, management.AddProjectGrantRequestToDomain(grant), org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "project_grant", Id: org.GetOrgId() + "_" + grant.GetProjectId() + "_" + grant.GetGrantedOrgId(), Message: err.Error()})
					continue
				}
				successOrg.ProjectGrants = append(successOrg.ProjectGrants, &admin_pb.ImportDataSuccessProjectGrant{ProjectId: grant.GetProjectId(), OrgId: grant.GetGrantedOrgId()})
			}
		}
		if org.UserGrants != nil {
			for _, grant := range org.GetUserGrants() {
				_, err := s.command.AddUserGrant(ctx, management.AddUserGrantRequestToDomain(grant), org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "user_grant", Id: org.GetOrgId() + "_" + grant.GetProjectId() + "_" + grant.GetUserId(), Message: err.Error()})
					continue
				}
				successOrg.UserGrants = append(successOrg.UserGrants, &admin_pb.ImportDataSuccessUserGrant{ProjectId: grant.GetProjectId(), UserId: grant.GetUserId()})
			}
		}
		if org.OrgMembers != nil {
			for _, member := range org.GetOrgMembers() {
				_, err := s.command.AddOrgMember(ctx, org.GetOrgId(), member.GetUserId(), member.GetRoles()...)
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "org_member", Id: org.GetOrgId() + "_" + member.GetUserId(), Message: err.Error()})
					continue
				}
				successOrg.OrgMembers = append(successOrg.OrgMembers, member.GetUserId())
			}
		}
		if org.ProjectGrantMembers != nil {
			for _, member := range org.GetProjectGrantMembers() {
				_, err := s.command.AddProjectGrantMember(ctx, management.AddProjectGrantMemberRequestToDomain(member))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "project_grant_member", Id: org.GetOrgId() + "_" + member.GetProjectId() + "_" + member.GetGrantId() + "_" + member.GetUserId(), Message: err.Error()})
					continue
				}
				successOrg.ProjectGrantMembers = append(successOrg.ProjectGrantMembers, &admin_pb.ImportDataSuccessProjectGrantMember{ProjectId: member.GetProjectId(), GrantId: member.GetProjectId(), UserId: member.GetUserId()})
			}
		}
		if org.ProjectMembers != nil {
			for _, member := range org.GetProjectMembers() {
				_, err := s.command.AddProjectMember(ctx, management.AddProjectMemberRequestToDomain(member), org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "project_member", Id: org.GetOrgId() + "_" + member.GetProjectId() + "_" + member.GetUserId(), Message: err.Error()})
					continue
				}
				successOrg.ProjectMembers = append(successOrg.ProjectMembers, &admin_pb.ImportDataSuccessProjectMember{ProjectId: member.GetProjectId(), UserId: member.GetUserId()})
			}

			success.Orgs = append(success.Orgs, successOrg)
		}
	}

	return &admin_pb.ImportDataResponse{
		Errors:  errors,
		Success: success,
	}, nil
}
