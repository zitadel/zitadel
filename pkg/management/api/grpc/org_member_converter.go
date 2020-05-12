package grpc

import (
	"github.com/caos/logging"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/golang/protobuf/ptypes"
)

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

func orgMemberFromModel(member *org_model.OrgMember) *OrgMember {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-jC5wY").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-Nc2jJ").OnError(err).Debug("date parse failed")

	return &OrgMember{
		UserId:       member.UserID,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Roles:        member.Roles,
		Sequence:     member.Sequence,
	}
}
