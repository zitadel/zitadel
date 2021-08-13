package project

import (
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/errors"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_pb "github.com/caos/zitadel/pkg/grpc/project"
)

func ProjectToPb(project *proj_model.ProjectView) *proj_pb.Project {
	return &proj_pb.Project{
		Id:                   project.ProjectID,
		Details:              object_grpc.ToViewDetailsPb(project.Sequence, project.CreationDate, project.ChangeDate, project.ResourceOwner),
		Name:                 project.Name,
		State:                projectStateToPb(project.State),
		ProjectRoleAssertion: project.ProjectRoleAssertion,
		ProjectRoleCheck:     project.ProjectRoleCheck,
		HasProjectCheck:      project.HasProjectCheck,
	}
}

func GrantedProjectToPb(project *proj_model.ProjectGrantView) *proj_pb.GrantedProject {
	return &proj_pb.GrantedProject{
		GrantId:          project.GrantID,
		ProjectId:        project.ProjectID,
		Details:          object_grpc.ToViewDetailsPb(project.Sequence, project.CreationDate, project.ChangeDate, project.ResourceOwner),
		ProjectName:      project.Name,
		State:            grantedProjectStateToPb(project.State),
		ProjectOwnerId:   project.ResourceOwner,
		ProjectOwnerName: project.ResourceOwnerName,
		GrantedOrgId:     project.OrgID,
		GrantedOrgName:   project.OrgName,
		GrantedRoleKeys:  project.GrantedRoleKeys,
	}
}

func ProjectsToPb(projects []*proj_model.ProjectView) []*proj_pb.Project {
	p := make([]*proj_pb.Project, len(projects))
	for i, project := range projects {
		p[i] = ProjectToPb(project)
	}
	return p
}

func GrantedProjectsToPb(projects []*proj_model.ProjectGrantView) []*proj_pb.GrantedProject {
	p := make([]*proj_pb.GrantedProject, len(projects))
	for i, project := range projects {
		p[i] = GrantedProjectToPb(project)
	}
	return p
}

func projectStateToPb(state proj_model.ProjectState) proj_pb.ProjectState {
	switch state {
	case proj_model.ProjectStateActive:
		return proj_pb.ProjectState_PROJECT_STATE_ACTIVE
	case proj_model.ProjectStateInactive:
		return proj_pb.ProjectState_PROJECT_STATE_INACTIVE
	default:
		return proj_pb.ProjectState_PROJECT_STATE_UNSPECIFIED
	}
}

func grantedProjectStateToPb(state proj_model.ProjectState) proj_pb.ProjectGrantState {
	switch state {
	case proj_model.ProjectStateActive:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_ACTIVE
	case proj_model.ProjectStateInactive:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_INACTIVE
	default:
		return proj_pb.ProjectGrantState_PROJECT_GRANT_STATE_UNSPECIFIED
	}
}

func ProjectQueriesToModel(queries []*proj_pb.ProjectQuery) (_ []*proj_model.ProjectViewSearchQuery, err error) {
	q := make([]*proj_model.ProjectViewSearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = ProjectQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func ProjectQueryToModel(query *proj_pb.ProjectQuery) (*proj_model.ProjectViewSearchQuery, error) {
	switch q := query.Query.(type) {
	case *proj_pb.ProjectQuery_NameQuery:
		return ProjectQueryNameToModel(q.NameQuery), nil
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-Ags42", "List.Query.Invalid")
	}
}

func ProjectQueryNameToModel(query *proj_pb.ProjectNameQuery) *proj_model.ProjectViewSearchQuery {
	return &proj_model.ProjectViewSearchQuery{
		Key:    proj_model.ProjectViewSearchKeyName,
		Method: object_grpc.TextMethodToModel(query.Method),
		Value:  query.Name,
	}
}

func GrantedProjectQueriesToModel(queries []*proj_pb.ProjectQuery) (_ []*proj_model.ProjectGrantViewSearchQuery, err error) {
	q := make([]*proj_model.ProjectGrantViewSearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = GrantedProjectQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func GrantedProjectQueryToModel(query *proj_pb.ProjectQuery) (*proj_model.ProjectGrantViewSearchQuery, error) {
	switch q := query.Query.(type) {
	case *proj_pb.ProjectQuery_NameQuery:
		return GrantedProjectQueryNameToModel(q.NameQuery), nil
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-Ags42", "List.Query.Invalid")
	}
}

func GrantedProjectQueryNameToModel(query *proj_pb.ProjectNameQuery) *proj_model.ProjectGrantViewSearchQuery {
	return &proj_model.ProjectGrantViewSearchQuery{
		Key:    proj_model.GrantedProjectSearchKeyName,
		Method: object_grpc.TextMethodToModel(query.Method),
		Value:  query.Name,
	}
}

func RoleQueriesToModel(queries []*proj_pb.RoleQuery) (_ []*proj_model.ProjectRoleSearchQuery, err error) {
	q := make([]*proj_model.ProjectRoleSearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = RoleQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func RoleQueryToModel(query *proj_pb.RoleQuery) (*proj_model.ProjectRoleSearchQuery, error) {
	switch q := query.Query.(type) {
	case *proj_pb.RoleQuery_KeyQuery:
		return RoleQueryKeyToModel(q.KeyQuery), nil
	case *proj_pb.RoleQuery_DisplayNameQuery:
		return RoleQueryDisplayNameToModel(q.DisplayNameQuery), nil
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-Ags42", "List.Query.Invalid")
	}
}

func RoleQueryKeyToModel(query *proj_pb.RoleKeyQuery) *proj_model.ProjectRoleSearchQuery {
	return &proj_model.ProjectRoleSearchQuery{
		Key:    proj_model.ProjectRoleSearchKeyKey,
		Method: object_grpc.TextMethodToModel(query.Method),
		Value:  query.Key,
	}
}

func RoleQueryDisplayNameToModel(query *proj_pb.RoleDisplayNameQuery) *proj_model.ProjectRoleSearchQuery {
	return &proj_model.ProjectRoleSearchQuery{
		Key:    proj_model.ProjectRoleSearchKeyDisplayName,
		Method: object_grpc.TextMethodToModel(query.Method),
		Value:  query.DisplayName,
	}
}

func RolesToPb(roles []*proj_model.ProjectRoleView) []*proj_pb.Role {
	r := make([]*proj_pb.Role, len(roles))
	for i, role := range roles {
		r[i] = RoleToPb(role)
	}
	return r
}

func RoleToPb(role *proj_model.ProjectRoleView) *proj_pb.Role {
	return &proj_pb.Role{
		Key:         role.Key,
		Details:     object_grpc.ToViewDetailsPb(role.Sequence, role.CreationDate, role.ChangeDate, role.ResourceOwner),
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}
