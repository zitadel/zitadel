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
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func groupStateToPb(state domain.GroupState) group_pb.GroupState {
	switch state {
	case domain.GroupStateActive:
		return group_pb.GroupState_GROUP_STATE_ACTIVE
	case domain.GroupStateInactive:
		return group_pb.GroupState_GROUP_STATE_INACTIVE
	default:
		return group_pb.GroupState_GROUP_STATE_UNSPECIFIED
	}
}

func projectGrantStateToPb(state domain.ProjectGrantState) group_pb.ProjectGrantState {
	switch state {
	case domain.ProjectGrantStateActive:
		return group_pb.ProjectGrantState_PROJECT_GRANT_STATE_ACTIVE
	case domain.ProjectGrantStateInactive:
		return group_pb.ProjectGrantState_PROJECT_GRANT_STATE_INACTIVE
	default:
		return group_pb.ProjectGrantState_PROJECT_GRANT_STATE_UNSPECIFIED
	}
}

func privateLabelingSettingToPb(setting domain.PrivateLabelingSetting) group_pb.PrivateLabelingSetting {
	switch setting {
	case domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy:
		return group_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY
	case domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy:
		return group_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY
	default:
		return group_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED
	}
}

func RoleQueriesToModel(queries []*group_pb.RoleQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = RoleQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func RoleQueryToModel(apiQuery *group_pb.RoleQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *group_pb.RoleQuery_KeyQuery:
		return query.NewProjectRoleKeySearchQuery(object.TextMethodToQuery(q.KeyQuery.Method), q.KeyQuery.Key)
	case *group_pb.RoleQuery_DisplayNameQuery:
		return query.NewProjectRoleDisplayNameSearchQuery(object.TextMethodToQuery(q.DisplayNameQuery.Method), q.DisplayNameQuery.DisplayName)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-fms0e", "List.Query.Invalid")
	}
}

func RoleViewsToPb(roles []*query.ProjectRole) []*group_pb.Role {
	o := make([]*group_pb.Role, len(roles))
	for i, org := range roles {
		o[i] = RoleViewToPb(org)
	}
	return o
}

func RoleViewToPb(role *query.ProjectRole) *group_pb.Role {
	return &group_pb.Role{
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
		Details: object.ToViewDetailsPb(

			role.Sequence,
			role.CreationDate,
			role.ChangeDate,
			role.ResourceOwner,
		),
	}
}
