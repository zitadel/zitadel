package system

import (
	"context"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/admin"
	"github.com/zitadel/zitadel/internal/api/grpc/management"
	"github.com/zitadel/zitadel/internal/domain"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func (s *Server) ImportInstanceData(ctx context.Context, req *system_pb.ImportInstanceDataRequest) (*system_pb.ImportInstanceDataResponse, error) {
	ctx = authz.WithInstanceID(ctx, req.InstanceId)
	errors := make([]*system_pb.ImportInstanceDataError, 0)

	appSecretGenerator, err := s.query.InitHashGenerator(ctx, domain.SecretGeneratorTypeAppSecret, s.passwordHashAlg)
	if err != nil {
		return nil, err
	}

	resourceOwner := ""

	for _, org := range req.GetOrgs() {
		_, err := s.command.AddOrg(ctx, org.GetOrg().GetName(), org.GetOwnerId(), resourceOwner, []string{})
		if err != nil {
			errors = append(errors, &system_pb.ImportInstanceDataError{Type: "org", Id: org.GetOrgId(), Message: err.Error()})
		} else {
			domainPolicy := org.GetDomainPolicy()
			_, err := s.command.AddOrgDomainPolicy(ctx, org.GetOrgId(), admin.DomainPolicyToDomain(domainPolicy.UserLoginMustBeDomain, domainPolicy.ValidateOrgDomains, domainPolicy.SmtpSenderAddressMatchesInstanceDomain))
			if err != nil {
				errors = append(errors, &system_pb.ImportInstanceDataError{Type: "domain_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
			_, err = s.command.AddLabelPolicy(ctx, org.GetOrgId(), management.AddLabelPolicyToDomain(org.GetLabelPolicy()))
			if err != nil {
				errors = append(errors, &system_pb.ImportInstanceDataError{Type: "label_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
			_, err = s.command.AddLockoutPolicy(ctx, org.GetOrgId(), management.AddLockoutPolicyToDomain(org.GetLockoutPolicy()))
			if err != nil {
				errors = append(errors, &system_pb.ImportInstanceDataError{Type: "lockout_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
			_, err = s.command.AddLoginPolicy(ctx, org.GetOrgId(), management.AddLoginPolicyToDomain(org.GetLoginPolicy()))
			if err != nil {
				errors = append(errors, &system_pb.ImportInstanceDataError{Type: "login_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
			_, err = s.command.AddPasswordComplexityPolicy(ctx, org.GetOrgId(), management.AddPasswordComplexityPolicyToDomain(org.GetPasswordComplexityPolicy()))
			if err != nil {
				errors = append(errors, &system_pb.ImportInstanceDataError{Type: "password_complexity_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
			_, err = s.command.AddPrivacyPolicy(ctx, org.GetOrgId(), management.AddPrivacyPolicyToDomain(org.GetPrivacyPolicy()))
			if err != nil {
				errors = append(errors, &system_pb.ImportInstanceDataError{Type: "privacy_policy", Id: org.GetOrgId(), Message: err.Error()})
			}

			for _, user := range org.GetHumanUsers() {
				_, err := s.command.AddHumanWithID(ctx, org.GetOrgId(), user.GetUserId(), management.AddHumanUserRequestToAddHuman(user.GetUser()))
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "human_user", Id: user.GetUserId(), Message: err.Error()})
				}
			}
			for _, user := range org.GetMachineUsers() {
				_, err := s.command.AddMachineWithID(ctx, org.GetOrgId(), user.GetUserId(), management.AddMachineUserRequestToDomain(user.GetUser()))
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "machine_user", Id: user.GetUserId(), Message: err.Error()})
				}
			}
			for _, project := range org.GetProjects() {
				_, err := s.command.AddProjectWithID(ctx, management.ProjectCreateToDomain(project.GetProject()), org.GetOrgId(), project.GetOwnerId(), project.GetProjectId())
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "project", Id: project.GetProjectId(), Message: err.Error()})
				}
			}
			for _, app := range org.GetOidcApps() {
				_, err := s.command.AddOIDCApplicationWithID(ctx, management.AddOIDCAppRequestToDomain(app.App), org.GetOrgId(), app.GetAppId(), appSecretGenerator)
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "oidc_app", Id: app.GetAppId(), Message: err.Error()})
				}
			}
			for _, app := range org.GetApiApps() {
				_, err := s.command.AddAPIApplicationWithID(ctx, management.AddAPIAppRequestToDomain(app.GetApp()), org.GetOrgId(), app.GetAppId(), appSecretGenerator)
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "api_app", Id: app.GetAppId(), Message: err.Error()})
				}
			}
			for _, action := range org.GetActions() {
				_, _, err := s.command.AddActionWithID(ctx, management.CreateActionRequestToDomain(action.GetAction()), org.GetOrgId(), action.GetActionId())
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "action", Id: action.GetActionId(), Message: err.Error()})
				}
			}
			for _, grant := range org.GetProjectGrants() {
				_, err := s.command.AddProjectGrant(ctx, management.AddProjectGrantRequestToDomain(grant), org.GetOrgId())
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "project_grant", Id: org.GetOrgId() + "_" + grant.GetProjectId() + "_" + grant.GetGrantedOrgId(), Message: err.Error()})
				}
			}
			for _, grant := range org.GetUserGrants() {
				_, err := s.command.AddUserGrant(ctx, management.AddUserGrantRequestToDomain(grant), org.GetOrgId())
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "user_grant", Id: org.GetOrgId() + "_" + grant.GetProjectId() + "_" + grant.GetUserId(), Message: err.Error()})
				}
			}
			for _, member := range org.GetOrgMembers() {
				_, err := s.command.AddOrgMember(ctx, org.GetOrgId(), member.GetUserId(), member.GetRoles()...)
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "org_member", Id: org.GetOrgId() + "_" + member.GetUserId(), Message: err.Error()})
				}
			}
			for _, member := range org.GetProjectGrantMembers() {
				_, err := s.command.AddProjectGrantMember(ctx, management.AddProjectGrantMemberRequestToDomain(member))
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "project_grant_member", Id: org.GetOrgId() + "_" + member.GetProjectId() + "_" + member.GetGrantId() + "_" + member.GetUserId(), Message: err.Error()})
				}
			}
			for _, member := range org.GetProjectMembers() {
				_, err := s.command.AddProjectMember(ctx, management.AddProjectMemberRequestToDomain(member), org.GetOrgId())
				if err != nil {
					errors = append(errors, &system_pb.ImportInstanceDataError{Type: "project_member", Id: org.GetOrgId() + "_" + member.GetProjectId() + "_" + member.GetUserId(), Message: err.Error()})
				}
			}
		}
	}

	return &system_pb.ImportInstanceDataResponse{}, nil
}
