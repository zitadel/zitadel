package admin

import (
	"context"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/management"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	management_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/policy"
	v1_pb "github.com/zitadel/zitadel/pkg/grpc/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (s *Server) ImportData(ctx context.Context, req *admin_pb.ImportDataRequest) (*admin_pb.ImportDataResponse, error) {
	errors := make([]*admin_pb.ImportDataError, 0)
	success := &admin_pb.ImportDataSuccess{}

	appSecretGenerator, err := s.query.InitHashGenerator(ctx, domain.SecretGeneratorTypeAppSecret, s.passwordHashAlg)
	if err != nil {
		return nil, err
	}
	initCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	phoneCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	passwordlessInitCode, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}

	ctxData := authz.GetCtxData(ctx)
	orgs := make([]*admin_pb.DataOrg, 0)

	if req.GetDataOrgsv1() != nil {
		dataOrgs, err := s.dataOrgsV1ToDataOrgs(ctx, req.GetDataOrgsv1())
		if err != nil {
			return nil, err
		}
		orgs = dataOrgs.GetOrgs()
	} else {
		orgs = req.GetDataOrgs().GetOrgs()
	}

	for _, org := range orgs {
		_, err := s.command.AddOrgWithID(ctx, org.GetOrg().GetName(), ctxData.UserID, ctxData.ResourceOwner, org.GetOrgId(), []string{})
		if err != nil {
			errors = append(errors, &admin_pb.ImportDataError{Type: "org", Id: org.GetOrgId(), Message: err.Error()})
		}
		successOrg := &admin_pb.ImportDataSuccessOrg{
			OrgId:               org.GetOrgId(),
			ProjectIds:          []string{},
			OidcAppIds:          []string{},
			ApiAppIds:           []string{},
			HumanUserIds:        []string{},
			MachineUserIds:      []string{},
			ActionIds:           []string{},
			ProjectGrants:       []*admin_pb.ImportDataSuccessProjectGrant{},
			UserGrants:          []*admin_pb.ImportDataSuccessUserGrant{},
			OrgMembers:          []string{},
			ProjectMembers:      []*admin_pb.ImportDataSuccessProjectMember{},
			ProjectGrantMembers: []*admin_pb.ImportDataSuccessProjectGrantMember{},
		}

		domainPolicy := org.GetDomainPolicy()
		if org.DomainPolicy != nil {
			_, err := s.command.AddOrgDomainPolicy(ctx, org.GetOrgId(), DomainPolicyToDomain(domainPolicy.UserLoginMustBeDomain, domainPolicy.ValidateOrgDomains, domainPolicy.SmtpSenderAddressMatchesInstanceDomain))
			if err != nil {
				errors = append(errors, &admin_pb.ImportDataError{Type: "domain_policy", Id: org.GetOrgId(), Message: err.Error()})
			}
		}
		if org.Domains != nil {
			for _, domainR := range org.Domains {
				orgDomain := &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: org.GetOrgId(),
					},
					Domain: domainR.DomainName,
				}
				_, err := s.command.AddOrgDomain(ctx, orgDomain, []string{})
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "domain", Id: org.GetOrgId() + "_" + domainR.DomainName, Message: err.Error()})
					continue
				}
				if domainR.IsVerified {
					if _, err := s.command.ValidateOrgDomain(ctx, orgDomain, []string{}); err != nil {
						errors = append(errors, &admin_pb.ImportDataError{Type: "vaildate_domain", Id: org.GetOrgId() + "_" + domainR.DomainName, Message: err.Error()})
						continue
					}
				}
				if domainR.IsPrimary {
					if _, err := s.command.SetPrimaryOrgDomain(ctx, orgDomain); err != nil {
						errors = append(errors, &admin_pb.ImportDataError{Type: "primary_domain", Id: org.GetOrgId() + "_" + domainR.DomainName, Message: err.Error()})
						continue
					}
				}
				successOrg.Domains = append(successOrg.Domains, domainR.DomainName)
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
		if org.OidcIdps != nil {
			for _, idp := range org.OidcIdps {
				_, err := s.command.ImportIDPConfig(ctx, management.AddOIDCIDPRequestToDomain(idp.Idp), idp.IdpId, org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "oidc_idp", Id: idp.IdpId, Message: err.Error()})
					continue
				}
				successOrg.OidcIpds = append(successOrg.OidcIpds, idp.GetIdpId())
			}
		}
		if org.JwtIdps != nil {
			for _, idp := range org.JwtIdps {
				_, err := s.command.ImportIDPConfig(ctx, management.AddJWTIDPRequestToDomain(idp.Idp), idp.IdpId, org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "jwt_idp", Id: idp.IdpId, Message: err.Error()})
					continue
				}
				successOrg.JwtIdps = append(successOrg.JwtIdps, idp.GetIdpId())
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
		if org.LoginTexts != nil {
			for _, text := range org.GetLoginTexts() {
				_, err := s.command.SetOrgLoginText(ctx, org.GetOrgId(), management.SetLoginCustomTextToDomain(text))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "login_texts", Id: org.GetOrgId() + "_" + text.Language, Message: err.Error()})
				}
			}
		}
		if org.InitMessages != nil {
			for _, message := range org.GetInitMessages() {
				_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetInitCustomTextToDomain(message))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "init_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
				}
			}
		}
		if org.PasswordResetMessages != nil {
			for _, message := range org.GetPasswordResetMessages() {
				_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetPasswordResetCustomTextToDomain(message))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "password_reset_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
				}
			}
		}
		if org.VerifyEmailMessages != nil {
			for _, message := range org.GetVerifyEmailMessages() {
				_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetVerifyEmailCustomTextToDomain(message))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "verify_email_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
				}
			}
		}
		if org.VerifyPhoneMessages != nil {
			for _, message := range org.GetVerifyPhoneMessages() {
				_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetVerifyPhoneCustomTextToDomain(message))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "verify_phone_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
				}
			}
		}
		if org.DomainClaimedMessages != nil {
			for _, message := range org.GetDomainClaimedMessages() {
				_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetDomainClaimedCustomTextToDomain(message))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "domain_claimed_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
				}
			}
		}
		if org.PasswordlessRegistrationMessages != nil {
			for _, message := range org.GetPasswordlessRegistrationMessages() {
				_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetPasswordlessRegistrationCustomTextToDomain(message))
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "passwordless_registration_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
				}
			}
		}

		if org.HumanUsers != nil {
			for _, user := range org.GetHumanUsers() {
				human, passwordless := management.ImportHumanUserRequestToDomain(user.User)
				human.AggregateID = user.UserId
				_, _, err := s.command.ImportHuman(ctx, org.GetOrgId(), human, passwordless, initCodeGenerator, phoneCodeGenerator, passwordlessInitCode)
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
				successOrg.MachineUserIds = append(successOrg.MachineUserIds, user.GetUserId())
			}
		}
		if org.UserMetadata != nil {
			for _, userMetadata := range org.GetUserMetadata() {
				_, err := s.command.SetUserMetadata(ctx, &domain.Metadata{Key: userMetadata.GetKey(), Value: userMetadata.GetValue()}, userMetadata.GetId(), org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "user_metadata", Id: userMetadata.GetId() + "_" + userMetadata.GetKey(), Message: err.Error()})
					continue
				}
				successOrg.UserMetadata = append(successOrg.UserMetadata, &admin_pb.ImportDataSuccessUserMetadata{UserId: userMetadata.GetId(), Key: userMetadata.GetKey()})
			}
		}
		if org.UserLinks != nil {
			for _, userLinks := range org.GetUserLinks() {
				externalIDP := &domain.UserIDPLink{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: userLinks.UserId},
					IDPConfigID:    userLinks.IdpId,
					ExternalUserID: userLinks.ProvidedUserId,
					DisplayName:    userLinks.ProvidedUserName,
				}
				if err := s.command.AddUserIDPLink(ctx, userLinks.UserId, org.GetOrgId(), externalIDP); err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "user_link", Id: userLinks.UserId + "_" + userLinks.IdpId, Message: err.Error()})
					continue
				}
				successOrg.UserLinks = append(successOrg.UserLinks, &admin_pb.ImportDataSuccessUserLinks{UserId: userLinks.GetUserId(), IdpId: userLinks.GetIdpId(), ExternalUserId: userLinks.GetProvidedUserId(), DisplayName: userLinks.GetProvidedUserName()})
			}
		}
		if org.Projects != nil {
			for _, project := range org.GetProjects() {
				_, err := s.command.AddProjectWithID(ctx, management.ProjectCreateToDomain(project.GetProject()), org.GetOrgId(), project.GetProjectId())
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
		if org.ProjectRoles != nil {
			for _, role := range org.GetProjectRoles() {
				_, err := s.command.AddProjectRole(ctx, management.AddProjectRoleRequestToDomain(role), org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "project_role", Id: role.ProjectId + "_" + role.RoleKey, Message: err.Error()})
					continue
				}
				successOrg.ProjectRoles = append(successOrg.ActionIds, role.ProjectId+"_"+role.RoleKey)
			}
		}
		success.Orgs = append(success.Orgs, successOrg)
	}

	for _, org := range orgs {
		var successOrg *admin_pb.ImportDataSuccessOrg
		for _, oldOrd := range success.Orgs {
			if org.OrgId == oldOrd.OrgId {
				successOrg = oldOrd
			}
		}
		if org.TriggerActions != nil {
			for _, triggerAction := range org.GetTriggerActions() {
				_, err := s.command.SetTriggerActions(ctx, domain.FlowType(triggerAction.FlowType), domain.TriggerType(triggerAction.TriggerType), triggerAction.ActionIds, org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "trigger_action", Id: triggerAction.FlowType.String() + "_" + triggerAction.TriggerType.String(), Message: err.Error()})
					continue
				}
				successOrg.TriggerActions = append(successOrg.TriggerActions, &management_pb.SetTriggerActionsRequest{FlowType: triggerAction.FlowType, TriggerType: triggerAction.TriggerType, ActionIds: triggerAction.GetActionIds()})
			}
		}
		if org.ProjectGrants != nil {
			for _, grant := range org.GetProjectGrants() {
				_, err := s.command.AddProjectGrantWithID(ctx, management.AddProjectGrantRequestToDomain(grant.GetProjectGrant()), grant.GetGrantId(), org.GetOrgId())
				if err != nil {
					errors = append(errors, &admin_pb.ImportDataError{Type: "project_grant", Id: org.GetOrgId() + "_" + grant.GetProjectGrant().GetProjectId() + "_" + grant.GetProjectGrant().GetGrantedOrgId(), Message: err.Error()})
					continue
				}
				successOrg.ProjectGrants = append(successOrg.ProjectGrants, &admin_pb.ImportDataSuccessProjectGrant{GrantId: grant.GetGrantId(), ProjectId: grant.GetProjectGrant().GetProjectId(), OrgId: grant.GetProjectGrant().GetGrantedOrgId()})
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
	}

	if success != nil && success.Orgs != nil {
		for _, org := range orgs {
			var successOrg *admin_pb.ImportDataSuccessOrg
			for _, oldOrd := range success.Orgs {
				if org.OrgId == oldOrd.OrgId {
					successOrg = oldOrd
				}
			}
			if successOrg == nil {
				continue
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
					successOrg.ProjectGrantMembers = append(successOrg.ProjectGrantMembers, &admin_pb.ImportDataSuccessProjectGrantMember{ProjectId: member.GetProjectId(), GrantId: member.GetGrantId(), UserId: member.GetUserId()})
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
			}
		}
	}
	return &admin_pb.ImportDataResponse{
		Errors:  errors,
		Success: success,
	}, nil
}

func (s *Server) dataOrgsV1ToDataOrgs(ctx context.Context, dataOrgs *v1_pb.ImportDataOrg) (*admin_pb.ImportDataOrg, error) {
	orgs := make([]*admin_pb.DataOrg, 0)
	for _, orgV1 := range dataOrgs.Orgs {
		org := &admin_pb.DataOrg{
			OrgId:                            orgV1.GetOrgId(),
			Org:                              orgV1.GetOrg(),
			DomainPolicy:                     nil,
			LabelPolicy:                      orgV1.GetLabelPolicy(),
			LockoutPolicy:                    orgV1.GetLockoutPolicy(),
			LoginPolicy:                      orgV1.GetLoginPolicy(),
			PasswordComplexityPolicy:         orgV1.GetPasswordComplexityPolicy(),
			PrivacyPolicy:                    orgV1.GetPrivacyPolicy(),
			Projects:                         orgV1.GetProjects(),
			ProjectRoles:                     orgV1.GetProjectRoles(),
			ApiApps:                          orgV1.GetApiApps(),
			OidcApps:                         orgV1.GetOidcApps(),
			HumanUsers:                       orgV1.GetHumanUsers(),
			MachineUsers:                     orgV1.GetMachineUsers(),
			TriggerActions:                   orgV1.GetTriggerActions(),
			Actions:                          orgV1.GetActions(),
			ProjectGrants:                    orgV1.GetProjectGrants(),
			UserGrants:                       orgV1.GetUserGrants(),
			OrgMembers:                       orgV1.GetOrgMembers(),
			ProjectMembers:                   orgV1.GetProjectMembers(),
			ProjectGrantMembers:              orgV1.GetProjectGrantMembers(),
			UserMetadata:                     orgV1.GetUserMetadata(),
			LoginTexts:                       orgV1.GetLoginTexts(),
			InitMessages:                     orgV1.GetInitMessages(),
			PasswordResetMessages:            orgV1.GetPasswordResetMessages(),
			VerifyEmailMessages:              orgV1.GetVerifyEmailMessages(),
			VerifyPhoneMessages:              orgV1.GetVerifyPhoneMessages(),
			DomainClaimedMessages:            orgV1.GetDomainClaimedMessages(),
			PasswordlessRegistrationMessages: orgV1.GetPasswordlessRegistrationMessages(),
			OidcIdps:                         orgV1.GetOidcIdps(),
			JwtIdps:                          orgV1.GetJwtIdps(),
			UserLinks:                        orgV1.GetUserLinks(),
			Domains:                          orgV1.GetDomains(),
		}
		if orgV1.IamPolicy != nil {
			defaultDomainPolicy, err := s.query.DefaultDomainPolicy(ctx)
			if err != nil {
				return nil, err
			}

			org.DomainPolicy = &admin_pb.AddCustomDomainPolicyRequest{
				UserLoginMustBeDomain:                  orgV1.IamPolicy.UserLoginMustBeDomain,
				ValidateOrgDomains:                     defaultDomainPolicy.ValidateOrgDomains,
				SmtpSenderAddressMatchesInstanceDomain: defaultDomainPolicy.SMTPSenderAddressMatchesInstanceDomain,
			}
		}
		if org.LoginPolicy != nil {
			defaultLoginPolicy, err := s.query.DefaultLoginPolicy(ctx)
			if err != nil {
				return nil, err
			}
			org.LoginPolicy.ExternalLoginCheckLifetime = durationpb.New(defaultLoginPolicy.ExternalLoginCheckLifetime)
			org.LoginPolicy.MultiFactorCheckLifetime = durationpb.New(defaultLoginPolicy.MultiFactorCheckLifetime)
			org.LoginPolicy.SecondFactorCheckLifetime = durationpb.New(defaultLoginPolicy.SecondFactorCheckLifetime)
			org.LoginPolicy.PasswordCheckLifetime = durationpb.New(defaultLoginPolicy.PasswordCheckLifetime)
			org.LoginPolicy.MfaInitSkipLifetime = durationpb.New(defaultLoginPolicy.MFAInitSkipLifetime)

			if orgV1.SecondFactors != nil {
				org.LoginPolicy.SecondFactors = make([]policy.SecondFactorType, len(orgV1.SecondFactors))
				for i, factor := range orgV1.SecondFactors {
					org.LoginPolicy.SecondFactors[i] = factor.GetType()
				}
			}
			if orgV1.MultiFactors != nil {
				org.LoginPolicy.MultiFactors = make([]policy.MultiFactorType, len(orgV1.MultiFactors))
				for i, factor := range orgV1.MultiFactors {
					org.LoginPolicy.MultiFactors[i] = factor.GetType()
				}
			}
			if orgV1.Idps != nil {
				org.LoginPolicy.Idps = make([]*management_pb.AddCustomLoginPolicyRequest_IDP, len(orgV1.Idps))
				for i, idpR := range orgV1.Idps {
					org.LoginPolicy.Idps[i] = &management_pb.AddCustomLoginPolicyRequest_IDP{
						IdpId:     idpR.GetIdpId(),
						OwnerType: idpR.GetOwnerType(),
					}
				}
			}
		}
		orgs = append(orgs, org)
	}

	return &admin_pb.ImportDataOrg{
		Orgs: orgs,
	}, nil
}
