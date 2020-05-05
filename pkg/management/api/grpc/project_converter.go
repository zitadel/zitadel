package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/golang/protobuf/ptypes"
)

func projectFromModel(project *proj_model.Project) *Project {
	creationDate, err := ptypes.TimestampProto(project.CreationDate)
	logging.Log("GRPC-iejs3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(project.ChangeDate)
	logging.Log("GRPC-di7rw").OnError(err).Debug("unable to parse timestamp")

	return &Project{
		Id:           project.AggregateID,
		State:        projectStateFromModel(project.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Name:         project.Name,
		Sequence:     project.Sequence,
	}
}

func grantedProjectSearchResponseFromModel(response *proj_model.GrantedProjectSearchResponse) *GrantedProjectSearchResponse {
	return &GrantedProjectSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      grantedProjectsFromModel(response.Result),
	}
}

func grantedProjectsFromModel(projects []*proj_model.GrantedProject) []*GrantedProject {
	converted := make([]*GrantedProject, 0)
	for _, q := range projects {
		converted = append(converted, grantedProjectFromModel(q))
	}
	return converted
}

func grantedProjectFromModel(project *proj_model.GrantedProject) *GrantedProject {
	creationDate, err := ptypes.TimestampProto(project.CreationDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(project.ChangeDate)
	logging.Log("GRPC-sope3").OnError(err).Debug("unable to parse timestamp")

	return &GrantedProject{
		Id:            project.ProjectID,
		State:         projectStateFromModel(project.State),
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		Name:          project.Name,
		Sequence:      project.Sequence,
		ResourceOwner: project.ResourceOwner,
		OrgId:         project.OrgID,
		OrgName:       project.OrgName,
		OrgDomain:     project.OrgDomain,
		GrantId:       project.GrantID,
	}
}

func projectRoleSearchResponseFromModel(response *proj_model.ProjectRoleSearchResponse) *ProjectRoleSearchResponse {
	return &ProjectRoleSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectRoleViewsFromModel(response.Result),
	}
}

func projectRoleViewsFromModel(roles []*proj_model.ProjectRoleView) []*ProjectRoleView {
	converted := make([]*ProjectRoleView, 0)
	for _, q := range roles {
		converted = append(converted, projectRoleViewFromModel(q))
	}
	return converted
}

func projectRoleViewFromModel(role *proj_model.ProjectRoleView) *ProjectRoleView {
	creationDate, err := ptypes.TimestampProto(role.CreationDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	return &ProjectRoleView{
		ProjectId:    role.ProjectID,
		CreationDate: creationDate,
		Key:          role.Key,
		Group:        role.Group,
		DisplayName:  role.DisplayName,
		Sequence:     role.Sequence,
	}
}

func projectStateFromModel(state proj_model.ProjectState) ProjectState {
	switch state {
	case proj_model.PROJECTSTATE_ACTIVE:
		return ProjectState_PROJECTSTATE_ACTIVE
	case proj_model.PROJECTSTATE_INACTIVE:
		return ProjectState_PROJECTSTATE_INACTIVE
	default:
		return ProjectState_PROJECTSTATE_UNSPECIFIED
	}
}

func projectTypeFromModel(projecttype proj_model.ProjectType) ProjectType {
	switch projecttype {
	case proj_model.PROJECTTYPE_OWNED:
		return ProjectType_PROJECTTYPE_OWNED
	case proj_model.PROJECTTYPE_GRANTED:
		return ProjectType_PROJECTTYPE_GRANTED
	default:
		return ProjectType_PROJECTTYPE_UNSPECIFIED
	}
}

func projectUpdateToModel(project *ProjectUpdateRequest) *proj_model.Project {
	return &proj_model.Project{
		ObjectRoot: models.ObjectRoot{
			AggregateID: project.Id,
		},
		Name: project.Name,
	}
}

func projectRoleFromModel(role *proj_model.ProjectRole) *ProjectRole {
	creationDate, err := ptypes.TimestampProto(role.CreationDate)
	logging.Log("GRPC-due83").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(role.ChangeDate)
	logging.Log("GRPC-id93s").OnError(err).Debug("unable to parse timestamp")

	return &ProjectRole{
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     role.Sequence,
		Key:          role.Key,
		DisplayName:  role.DisplayName,
		Group:        role.Group,
	}
}

func projectRoleAddToModel(role *ProjectRoleAdd) *proj_model.ProjectRole {
	return &proj_model.ProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID: role.Id,
		},
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}

func projectRoleChangeToModel(role *ProjectRoleChange) *proj_model.ProjectRole {
	return &proj_model.ProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID: role.Id,
		},
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}

func grantedProjectSearchRequestsToModel(project *GrantedProjectSearchRequest) *proj_model.GrantedProjectSearchRequest {
	return &proj_model.GrantedProjectSearchRequest{
		Offset:  project.Offset,
		Limit:   project.Limit,
		Queries: grantedProjectSearchQueriesToModel(project.Queries),
	}
}

func grantedProjectSearchQueriesToModel(queries []*GrantedProjectSearchQuery) []*proj_model.GrantedProjectSearchQuery {
	converted := make([]*proj_model.GrantedProjectSearchQuery, 0)
	for _, q := range queries {
		converted = append(converted, grantedProjectSearchQueryToModel(q))
	}
	return converted
}

func grantedProjectSearchQueryToModel(query *GrantedProjectSearchQuery) *proj_model.GrantedProjectSearchQuery {
	return &proj_model.GrantedProjectSearchQuery{
		Key:    projectSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectSearchKeyToModel(key GrantedProjectSearchKey) proj_model.GrantedProjectSearchKey {
	switch key {
	case GrantedProjectSearchKey_PROJECTSEARCHKEY_PROJECT_NAME:
		return proj_model.GRANTEDPROJECTSEARCHKEY_NAME
	default:
		return proj_model.GRANTEDPROJECTSEARCHKEY_UNSPECIFIED
	}
}

func projectRoleSearchRequestsToModel(role *ProjectRoleSearchRequest) *proj_model.ProjectRoleSearchRequest {
	return &proj_model.ProjectRoleSearchRequest{
		Offset:  role.Offset,
		Limit:   role.Limit,
		Queries: projectRoleSearchQueriesToModel(role.Queries),
	}
}

func projectRoleSearchQueriesToModel(queries []*ProjectRoleSearchQuery) []*proj_model.ProjectRoleSearchQuery {
	converted := make([]*proj_model.ProjectRoleSearchQuery, 0)
	for _, q := range queries {
		converted = append(converted, projectRoleSearchQueryToModel(q))
	}
	return converted
}

func projectRoleSearchQueryToModel(query *ProjectRoleSearchQuery) *proj_model.ProjectRoleSearchQuery {
	return &proj_model.ProjectRoleSearchQuery{
		Key:    projectRoleSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectRoleSearchKeyToModel(key ProjectRoleSearchKey) proj_model.ProjectRoleSearchKey {
	switch key {
	case ProjectRoleSearchKey_PROJECTROLESEARCHKEY_KEY:
		return proj_model.PROJECTROLESEARCHKEY_KEY
	case ProjectRoleSearchKey_PROJECTROLESEARCHKEY_DISPLAY_NAME:
		return proj_model.PROJECTROLESEARCHKEY_DISPLAY_NAME
	default:
		return proj_model.PROJECTROLESEARCHKEY_UNSPECIFIED
	}
}
func searchMethodToModel(method SearchMethod) model.SearchMethod {
	switch method {
	case SearchMethod_SEARCHMETHOD_EQUALS:
		return model.SEARCHMETHOD_EQUALS
	case SearchMethod_SEARCHMETHOD_CONTAINS:
		return model.SEARCHMETHOD_CONTAINS
	case SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return model.SEARCHMETHOD_STARTS_WITH
	case SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return model.SEARCHMETHOD_EQUALS_IGNORE_CASE
	case SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return model.SEARCHMETHOD_CONTAINS_IGNORE_CASE
	case SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return model.SEARCHMETHOD_STARTS_WITH_IGNORE_CASE
	default:
		return model.SEARCHMETHOD_EQUALS
	}
}
