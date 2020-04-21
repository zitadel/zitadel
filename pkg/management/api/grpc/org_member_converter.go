package grpc

import org_model "github.com/caos/zitadel/internal/org/model"

func addOrgMemberToModel(member *AddOrgMemberRequest) *org_model.OrgMember {
	memberModel := org_model.NewOrgMember(member.OrgId, member.UserId)
	memberModel.Roles = member.Roles

	return memberModel
}

func changeOrgMemberToModel(member *ChangeOrgMemberRequest) *org_model.OrgMember {
	memberModel := org_model.NewOrgMember(member.OrgId, member.UserId)
	memberModel.Roles = member.Roles

	return memberModel
}
