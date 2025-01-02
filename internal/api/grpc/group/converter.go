package group

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	group_pb "github.com/zitadel/zitadel/pkg/grpc/group"
)

func GroupViewsToPb(groups []*query.Group) []*group_pb.Group {
	o := make([]*group_pb.Group, len(groups))
	for i, org := range groups {
		o[i] = GroupViewToPb(org)
	}
	return o
}

func GroupViewToPb(group *query.Group) *group_pb.Group {
	return &group_pb.Group{
		Id:          group.ID,
		State:       groupStateToPb(group.State),
		Name:        group.Name,
		Description: group.Description,
		Details: object.ToViewDetailsPb(
			group.Sequence,
			group.CreationDate,
			group.ChangeDate,
			group.ResourceOwner,
		),
	}
}

// func GrantedProjectViewsToPb(projects []*query.GroupGrant) []*group_pb.GrantedGroup {
// 	p := make([]*group_pb.Granted, len(projects))
// 	for i, project := range projects {
// 		p[i] = GrantedProjectViewToPb(project)
// 	}
// 	return p
// }

//	func GrantedProjectViewToPb(project *query.ProjectGrant) *group_pb.GrantedProject {
//		return &group_pb.GrantedProject{
//			ProjectId:        project.ProjectID,
//			GrantId:          project.GrantID,
//			Details:          object.ToViewDetailsPb(project.Sequence, project.CreationDate, project.ChangeDate, project.ResourceOwner),
//			ProjectName:      project.ProjectName,
//			State:            projectGrantStateToPb(project.State),
//			ProjectOwnerId:   project.ResourceOwner,
//			ProjectOwnerName: project.ResourceOwnerName,
//			GrantedOrgId:     project.GrantedOrgID,
//			GrantedOrgName:   project.OrgName,
//			GrantedRoleKeys:  project.GrantedRoleKeys,
//		}
//	}
func GroupQueriesToModel(queries []*group_pb.GroupQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = GroupQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func GroupQueryToModel(apiQuery *group_pb.GroupQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *group_pb.GroupQuery_NameQuery:
		return query.NewGroupNameSearchQuery(object.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	case *group_pb.GroupQuery_GroupResourceOwnerQuery:
		return query.NewGroupResourceOwnerSearchQuery(q.GroupResourceOwnerQuery.ResourceOwner)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-uQ9nC", "List.Query.Invalid")
	}
}

func groupStateToPb(state domain.GroupState) group_pb.GroupState {
	switch state {
	case domain.GroupStateActive:
		return group_pb.GroupState_GROUP_STATE_ACTIVE
	case domain.GroupStateInactive:
		return group_pb.GroupState_GROUP_STATE_INACTIVE
	case domain.GroupStateUnspecified:
		return group_pb.GroupState_GROUP_STATE_UNSPECIFIED
	default:
		return group_pb.GroupState_GROUP_STATE_UNSPECIFIED
	}
}

func groupGrantStateToPb(state domain.GroupGrantState) group_pb.GroupGrantState {
	switch state {
	case domain.GroupGrantStateActive:
		return group_pb.GroupGrantState_GROUP_GRANT_STATE_ACTIVE
	case domain.GroupGrantStateInactive:
		return group_pb.GroupGrantState_GROUP_GRANT_STATE_INACTIVE
	default:
		return group_pb.GroupGrantState_GROUP_GRANT_STATE_UNSPECIFIED
	}
}
