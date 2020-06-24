package grpc

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/management/grpc"

	"github.com/golang/protobuf/ptypes"
)

func projectGrantFromModel(grant *proj_model.ProjectGrant) *grpc.ProjectGrant {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-8d73s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	return &grpc.ProjectGrant{
		Id:           grant.GrantID,
		State:        projectGrantStateFromModel(grant.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		GrantedOrgId: grant.GrantedOrgID,
		RoleKeys:     grant.RoleKeys,
		Sequence:     grant.Sequence,
		ProjectId:    grant.AggregateID,
	}
}

func projectGrantCreateToModel(grant *grpc.ProjectGrantCreate) *proj_model.ProjectGrant {
	return &proj_model.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: grant.ProjectId,
		},
		GrantedOrgID: grant.GrantedOrgId,
		RoleKeys:     grant.RoleKeys,
	}
}

func projectGrantUpdateToModel(grant *grpc.ProjectGrantUpdate) *proj_model.ProjectGrant {
	return &proj_model.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: grant.ProjectId,
		},
		GrantID:  grant.Id,
		RoleKeys: grant.RoleKeys,
	}
}

func projectGrantSearchRequestsToModel(request *grpc.ProjectGrantSearchRequest) *proj_model.ProjectGrantViewSearchRequest {
	return &proj_model.ProjectGrantViewSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: projectGrantSearchQueriesToModel(request.ProjectId, request.Queries),
	}
}

func projectGrantSearchQueriesToModel(projectId string, queries []*grpc.ProjectGrantSearchQuery) []*proj_model.ProjectGrantViewSearchQuery {
	converted := make([]*proj_model.ProjectGrantViewSearchQuery, 0)
	converted = append(converted, &proj_model.ProjectGrantViewSearchQuery{
		Key:    proj_model.GrantedProjectSearchKeyProjectID,
		Method: model.SearchMethodEquals,
		Value:  projectId,
	})
	for i, query := range queries {
		converted[i] = projectGrantSearchQueryToModel(query)
	}
	return converted
}

func projectGrantSearchQueryToModel(query *grpc.ProjectGrantSearchQuery) *proj_model.ProjectGrantViewSearchQuery {
	return &proj_model.ProjectGrantViewSearchQuery{
		Key:    projectGrantViewSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectGrantViewSearchKeyToModel(key grpc.ProjectGrantSearchKey) proj_model.ProjectGrantViewSearchKey {
	switch key {
	case grpc.ProjectGrantSearchKey_PROJECTGRANTSEARCHKEY_PROJECT_NAME:
		return proj_model.GrantedProjectSearchKeyProjectID
	case grpc.ProjectGrantSearchKey_PROJECTGRANTSEARCHKEY_ROLE_KEY:
		return proj_model.GrantedProjectSearchKeyRoleKeys
	default:
		return proj_model.GrantedProjectSearchKeyUnspecified
	}
}

func projectGrantSearchResponseFromModel(response *proj_model.ProjectGrantViewSearchResponse) *grpc.ProjectGrantSearchResponse {
	return &grpc.ProjectGrantSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectGrantsFromGrantedProjectModel(response.Result),
	}
}

func projectGrantsFromGrantedProjectModel(projects []*proj_model.ProjectGrantView) []*grpc.ProjectGrantView {
	converted := make([]*grpc.ProjectGrantView, len(projects))
	for i, project := range projects {
		converted[i] = projectGrantFromGrantedProjectModel(project)
	}
	return converted
}

func projectGrantFromGrantedProjectModel(project *proj_model.ProjectGrantView) *grpc.ProjectGrantView {
	creationDate, err := ptypes.TimestampProto(project.CreationDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(project.ChangeDate)
	logging.Log("GRPC-sope3").OnError(err).Debug("unable to parse timestamp")

	return &grpc.ProjectGrantView{
		ProjectId:         project.ProjectID,
		State:             projectGrantStateFromProjectStateModel(project.State),
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		ProjectName:       project.Name,
		Sequence:          project.Sequence,
		GrantedOrgId:      project.OrgID,
		GrantedOrgName:    project.OrgName,
		Id:                project.GrantID,
		RoleKeys:          project.GrantedRoleKeys,
		ResourceOwner:     project.ResourceOwner,
		ResourceOwnerName: project.ResourceOwnerName,
	}
}

func projectGrantStateFromModel(state proj_model.ProjectGrantState) grpc.ProjectGrantState {
	switch state {
	case proj_model.ProjectGrantStateActive:
		return grpc.ProjectGrantState_PROJECTGRANTSTATE_ACTIVE
	case proj_model.ProjectGrantStateInactive:
		return grpc.ProjectGrantState_PROJECTGRANTSTATE_INACTIVE
	default:
		return grpc.ProjectGrantState_PROJECTGRANTSTATE_UNSPECIFIED
	}
}

func projectGrantStateFromProjectStateModel(state proj_model.ProjectState) grpc.ProjectGrantState {
	switch state {
	case proj_model.ProjectStateActive:
		return grpc.ProjectGrantState_PROJECTGRANTSTATE_ACTIVE
	case proj_model.ProjectStateInactive:
		return grpc.ProjectGrantState_PROJECTGRANTSTATE_INACTIVE
	default:
		return grpc.ProjectGrantState_PROJECTGRANTSTATE_UNSPECIFIED
	}
}
