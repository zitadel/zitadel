package admin

import (
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	member_pb "github.com/zitadel/zitadel/pkg/grpc/member"
)

func AddIAMMemberToCommand(req *admin_pb.AddIAMMemberRequest, instanceID string) *command.AddInstanceMember {
	return &command.AddInstanceMember{
		InstanceID: instanceID,
		UserID:     req.UserId,
		Roles:      req.Roles,
	}
}

func UpdateIAMMemberToCommand(req *admin_pb.UpdateIAMMemberRequest, instanceID string) *command.ChangeInstanceMember {
	return &command.ChangeInstanceMember{
		InstanceID: instanceID,
		UserID:     req.UserId,
		Roles:      req.Roles,
	}
}

func ListIAMMembersRequestToQuery(req *admin_pb.ListIAMMembersRequest) (*query.IAMMembersQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := member_grpc.MemberQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.IAMMembersQuery{
		MembersQuery: query.MembersQuery{
			SearchRequest: query.SearchRequest{
				Offset:        offset,
				Limit:         limit,
				Asc:           asc,
				SortingColumn: fieldNameToMemberColumn(req.SortingColumn),
			},
			Queries: queries,
		},
	}, nil
}

func fieldNameToMemberColumn(fieldName member_pb.MemberFieldColumnName) query.Column {
	switch fieldName {
	case member_pb.MemberFieldColumnName_MEMBER_FIELD_NAME_UNSPECIFIED:
		return query.InstanceMemberInstanceID
	case member_pb.MemberFieldColumnName_MEMBER_FIELD_NAME_USER_ID:
		return query.InstanceMemberUserID
	case member_pb.MemberFieldColumnName_MEMBER_FIELD_NAME_CREATION_DATE:
		return query.InstanceMemberCreationDate
	case member_pb.MemberFieldColumnName_MEMBER_FIELD_NAME_CHANGE_DATE:
		return query.InstanceMemberChangeDate
	case member_pb.MemberFieldColumnName_MEMBER_FIELD_NAME_USER_RESOURCE_OWNER:
		return query.InstanceMemberResourceOwner
	default:
		return query.Column{}
	}
}
