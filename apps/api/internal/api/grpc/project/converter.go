package project

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	proj_pb "github.com/zitadel/zitadel/pkg/grpc/project"
)

func ProjectViewsToPb(projects []*query.Project) []*proj_pb.Project {
	o := make([]*proj_pb.Project, len(projects))
	for i, org := range projects {
		o[i] = ProjectViewToPb(org)
	}
	return o
}

func ProjectViewToPb(project *query.Project) *proj_pb.Project {
	return &proj_pb.Project{
		Id:                     project.ID,
		State:                  projectStateToPb(project.State),
		Name:                   project.Name,
		PrivateLabelingSetting: privateLabelingSettingToPb(project.PrivateLabelingSetting),
		HasProjectCheck:        project.HasProjectCheck,
		ProjectRoleAssertion:   project.ProjectRoleAssertion,
		ProjectRoleCheck:       project.ProjectRoleCheck,
		Details: object.ToViewDetailsPb(
			project.Sequence,
			project.CreationDate,
			project.ChangeDate,
			project.ResourceOwner,
		),
	}
}

func GrantedProjectViewsToPb(projects []*query.ProjectGrant) []*proj_pb.GrantedProject {
	p := make([]*proj_pb.GrantedProject, len(projects))
	for i, project := range projects {
		p[i] = GrantedProjectViewToPb(project)
	}
	return p
}

func GrantedProjectViewToPb(project *query.ProjectGrant) *proj_pb.GrantedProject {
	return &proj_pb.GrantedProject{
		ProjectId:        project.ProjectID,
		GrantId:          project.GrantID,
		Details:          object.ToViewDetailsPb(project.Sequence, project.CreationDate, project.ChangeDate, project.ResourceOwner),
		ProjectName:      project.ProjectName,
		State:            projectGrantStateToPb(project.State),
		ProjectOwnerId:   project.ResourceOwner,
		ProjectOwnerName: project.ResourceOwnerName,
		GrantedOrgId:     project.GrantedOrgID,
		GrantedOrgName:   project.OrgName,
		GrantedRoleKeys:  project.GrantedRoleKeys,
	}
}
func ProjectQueriesToModel(queries []*proj_pb.ProjectQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = ProjectQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func ProjectQueryToModel(apiQuery *proj_pb.ProjectQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *proj_pb.ProjectQuery_NameQuery:
		return query.NewProjectNameSearchQuery(object.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	case *proj_pb.ProjectQuery_ProjectResourceOwnerQuery:
		return query.NewProjectResourceOwnerSearchQuery(q.ProjectResourceOwnerQuery.ResourceOwner)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func projectStateToPb(state domain.ProjectState) proj_pb.ProjectState {
	switch state {
	case domain.ProjectStateActive:
		return proj_pb.ProjectState_PROJECT_STATE_ACTIVE
	case domain.ProjectStateInactive:
		return proj_pb.ProjectState_PROJECT_STATE_INACTIVE
	default:
		return proj_pb.ProjectState_PROJECT_STATE_UNSPECIFIED
	}
}

func projectGrantStateToPb(state domain.ProjectGrantState) proj_pb.ProjectGrantState {
	switch state {
	case domain.ProjectGrantStateActive:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_ACTIVE
	case domain.ProjectGrantStateInactive:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_INACTIVE
	default:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_UNSPECIFIED
	}
}

func privateLabelingSettingToPb(setting domain.PrivateLabelingSetting) proj_pb.PrivateLabelingSetting {
	switch setting {
	case domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy:
		return proj_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY
	case domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy:
		return proj_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY
	default:
		return proj_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED
	}
}

func RoleQueriesToModel(queries []*proj_pb.RoleQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = RoleQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func RoleQueryToModel(apiQuery *proj_pb.RoleQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *proj_pb.RoleQuery_KeyQuery:
		return query.NewProjectRoleKeySearchQuery(object.TextMethodToQuery(q.KeyQuery.Method), q.KeyQuery.Key)
	case *proj_pb.RoleQuery_DisplayNameQuery:
		return query.NewProjectRoleDisplayNameSearchQuery(object.TextMethodToQuery(q.DisplayNameQuery.Method), q.DisplayNameQuery.DisplayName)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-fms0e", "List.Query.Invalid")
	}
}

func RoleViewsToPb(roles []*query.ProjectRole) []*proj_pb.Role {
	o := make([]*proj_pb.Role, len(roles))
	for i, org := range roles {
		o[i] = RoleViewToPb(org)
	}
	return o
}

func RoleViewToPb(role *query.ProjectRole) *proj_pb.Role {
	return &proj_pb.Role{
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
