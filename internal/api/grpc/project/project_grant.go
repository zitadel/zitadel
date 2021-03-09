package project

import (
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/grpc/project"
	proj_pb "github.com/caos/zitadel/pkg/grpc/project"
)

func ProjectGrantQueriesToModel(queries []*project.ProjectGrantQuery) []*proj_model.ProjectGrantViewSearchQuery {
	q := make([]*proj_model.ProjectGrantViewSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = ProjectGrantQueryToModel(query)
	}
	return q
}

func ProjectGrantQueryToModel(query *project.ProjectGrantQuery) *proj_model.ProjectGrantViewSearchQuery {
	switch q := query.Query.(type) {
	case *proj_pb.ProjectGrantQuery_ProjectNameQuery:
		return ProjectGrantNameToProjectQuery(q.ProjectNameQuery)
	case *proj_pb.ProjectGrantQuery_RoleKeyQuery:
		return ProjectGrantRoleToProjectQuery(q.RoleKeyQuery)
	default:
		return nil
	}
}

func ProjectGrantNameToProjectQuery(query *proj_pb.GrantProjectNameQuery) *proj_model.ProjectGrantViewSearchQuery {
	return &proj_model.ProjectGrantViewSearchQuery{
		Key:    proj_model.GrantedProjectSearchKeyName,
		Method: object_grpc.TextMethodToModel(query.Method),
		Value:  query.Name,
	}
}

func ProjectGrantRoleToProjectQuery(query *proj_pb.GrantRoleKeyQuery) *proj_model.ProjectGrantViewSearchQuery {
	return &proj_model.ProjectGrantViewSearchQuery{
		Key:    proj_model.GrantedProjectSearchKeyRoleKeys,
		Method: object_grpc.TextMethodToModel(query.Method),
		Value:  query.RoleKey,
	}
}
