package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/golang/protobuf/ptypes"
)

func projectGrantFromModel(grant *proj_model.ProjectGrant) *ProjectGrant {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-8d73s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	return &ProjectGrant{
		Id:           grant.GrantID,
		State:        projectGrantStateFromModel(grant.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		GrantedOrgId: grant.GrantedOrgID,
		RoleKeys:     grant.RoleKeys,
		Sequence:     grant.Sequence,
	}
}

func projectGrantCreateToModel(grant *ProjectGrantCreate) *proj_model.ProjectGrant {
	return &proj_model.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: grant.ProjectId,
		},
		GrantedOrgID: grant.GrantedOrgId,
		RoleKeys:     grant.RoleKeys,
	}
}

func projectGrantUpdateToModel(grant *ProjectGrantUpdate) *proj_model.ProjectGrant {
	return &proj_model.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: grant.ProjectId,
		},
		GrantID:  grant.Id,
		RoleKeys: grant.RoleKeys,
	}
}

func projectGrantSearchRequestsToModel(request *ProjectGrantSearchRequest) *proj_model.GrantedProjectSearchRequest {
	return &proj_model.GrantedProjectSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: projectGrantSearchQueriesToModel(request.ProjectId),
	}
}

func projectGrantSearchQueriesToModel(projectId string) []*proj_model.GrantedProjectSearchQuery {
	converted := make([]*proj_model.GrantedProjectSearchQuery, 0)
	return append(converted, &proj_model.GrantedProjectSearchQuery{
		Key:    proj_model.GRANTEDPROJECTSEARCHKEY_PROJECTID,
		Method: model.SEARCHMETHOD_EQUALS,
		Value:  projectId,
	})
}

func projectGrantSearchResponseFromModel(response *proj_model.GrantedProjectSearchResponse) *ProjectGrantSearchResponse {
	return &ProjectGrantSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectGrantsFromGrantedProjectModel(response.Result),
	}
}

func projectGrantsFromGrantedProjectModel(projects []*proj_model.GrantedProjectView) []*ProjectGrantView {
	converted := make([]*ProjectGrantView, len(projects))
	for i, project := range projects {
		converted = projectGrantFromGrantedProjectModel(project)
	}
	return converted
}

func projectGrantFromGrantedProjectModel(project *proj_model.GrantedProjectView) *ProjectGrantView {
	creationDate, err := ptypes.TimestampProto(project.CreationDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(project.ChangeDate)
	logging.Log("GRPC-sope3").OnError(err).Debug("unable to parse timestamp")

	return &ProjectGrantView{
		ProjectId:        project.ProjectID,
		State:            projectGrantStateFromProjectStateModel(project.State),
		CreationDate:     creationDate,
		ChangeDate:       changeDate,
		ProjectName:      project.Name,
		Sequence:         project.Sequence,
		GrantedOrgId:     project.OrgID,
		GrantedOrgName:   project.OrgName,
		GrantedOrgDomain: project.OrgDomain,
		Id:               project.GrantID,
	}
}

func projectGrantStateFromModel(state proj_model.ProjectGrantState) ProjectGrantState {
	switch state {
	case proj_model.PROJECTGRANTSTATE_ACTIVE:
		return ProjectGrantState_PROJECTGRANTSTATE_ACTIVE
	case proj_model.PROJECTGRANTSTATE_INACTIVE:
		return ProjectGrantState_PROJECTGRANTSTATE_INACTIVE
	default:
		return ProjectGrantState_PROJECTGRANTSTATE_UNSPECIFIED
	}
}

func projectGrantStateFromProjectStateModel(state proj_model.ProjectState) ProjectGrantState {
	switch state {
	case proj_model.PROJECTSTATE_ACTIVE:
		return ProjectGrantState_PROJECTGRANTSTATE_ACTIVE
	case proj_model.PROJECTSTATE_INACTIVE:
		return ProjectGrantState_PROJECTGRANTSTATE_INACTIVE
	default:
		return ProjectGrantState_PROJECTGRANTSTATE_UNSPECIFIED
	}
}
