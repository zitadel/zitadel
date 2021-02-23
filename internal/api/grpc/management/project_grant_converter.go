package management

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/v2/domain"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func projectGrantFromDomain(grant *domain.ProjectGrant) *management.ProjectGrant {
	return &management.ProjectGrant{
		Id:           grant.GrantID,
		State:        projectGrantStateFromDomain(grant.State),
		CreationDate: timestamppb.New(grant.CreationDate),
		ChangeDate:   timestamppb.New(grant.ChangeDate),
		GrantedOrgId: grant.GrantedOrgID,
		RoleKeys:     grant.RoleKeys,
		Sequence:     grant.Sequence,
		ProjectId:    grant.AggregateID,
	}
}

func projectGrantFromModel(grant *proj_model.ProjectGrant) *management.ProjectGrant {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-8d73s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	return &management.ProjectGrant{
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

func projectGrantCreateToDomain(grant *management.ProjectGrantCreate) *domain.ProjectGrant {
	return &domain.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: grant.ProjectId,
		},
		GrantedOrgID: grant.GrantedOrgId,
		RoleKeys:     grant.RoleKeys,
	}
}

func projectGrantUpdateToDomain(grant *management.ProjectGrantUpdate) *domain.ProjectGrant {
	return &domain.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: grant.ProjectId,
		},
		GrantID:  grant.Id,
		RoleKeys: grant.RoleKeys,
	}
}

func projectGrantSearchRequestsToModel(request *management.ProjectGrantSearchRequest) *proj_model.ProjectGrantViewSearchRequest {
	return &proj_model.ProjectGrantViewSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: projectGrantSearchQueriesToModel(request.ProjectId, request.Queries),
	}
}

func projectGrantSearchQueriesToModel(projectId string, queries []*management.ProjectGrantSearchQuery) []*proj_model.ProjectGrantViewSearchQuery {
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

func projectGrantSearchQueryToModel(query *management.ProjectGrantSearchQuery) *proj_model.ProjectGrantViewSearchQuery {
	return &proj_model.ProjectGrantViewSearchQuery{
		Key:    projectGrantViewSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectGrantViewSearchKeyToModel(key management.ProjectGrantSearchKey) proj_model.ProjectGrantViewSearchKey {
	switch key {
	case management.ProjectGrantSearchKey_PROJECTGRANTSEARCHKEY_PROJECT_NAME:
		return proj_model.GrantedProjectSearchKeyProjectID
	case management.ProjectGrantSearchKey_PROJECTGRANTSEARCHKEY_ROLE_KEY:
		return proj_model.GrantedProjectSearchKeyRoleKeys
	default:
		return proj_model.GrantedProjectSearchKeyUnspecified
	}
}

func projectGrantSearchResponseFromModel(response *proj_model.ProjectGrantViewSearchResponse) *management.ProjectGrantSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-MCjs7").OnError(err).Debug("unable to parse timestamp")
	return &management.ProjectGrantSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            projectGrantsFromGrantedProjectModel(response.Result),
		ViewTimestamp:     timestamp,
		ProcessedSequence: response.Sequence,
	}
}

func projectGrantsFromGrantedProjectModel(projects []*proj_model.ProjectGrantView) []*management.ProjectGrantView {
	converted := make([]*management.ProjectGrantView, len(projects))
	for i, project := range projects {
		converted[i] = projectGrantFromGrantedProjectModel(project)
	}
	return converted
}

func projectGrantFromGrantedProjectModel(project *proj_model.ProjectGrantView) *management.ProjectGrantView {
	creationDate, err := ptypes.TimestampProto(project.CreationDate)
	logging.Log("GRPC-dlso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(project.ChangeDate)
	logging.Log("GRPC-sope3").OnError(err).Debug("unable to parse timestamp")

	return &management.ProjectGrantView{
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

func projectGrantStateFromDomain(state domain.ProjectGrantState) management.ProjectGrantState {
	switch state {
	case domain.ProjectGrantStateActive:
		return management.ProjectGrantState_PROJECTGRANTSTATE_ACTIVE
	case domain.ProjectGrantStateInactive:
		return management.ProjectGrantState_PROJECTGRANTSTATE_INACTIVE
	default:
		return management.ProjectGrantState_PROJECTGRANTSTATE_UNSPECIFIED
	}
}
func projectGrantStateFromModel(state proj_model.ProjectGrantState) management.ProjectGrantState {
	switch state {
	case proj_model.ProjectGrantStateActive:
		return management.ProjectGrantState_PROJECTGRANTSTATE_ACTIVE
	case proj_model.ProjectGrantStateInactive:
		return management.ProjectGrantState_PROJECTGRANTSTATE_INACTIVE
	default:
		return management.ProjectGrantState_PROJECTGRANTSTATE_UNSPECIFIED
	}
}

func projectGrantStateFromProjectStateModel(state proj_model.ProjectState) management.ProjectGrantState {
	switch state {
	case proj_model.ProjectStateActive:
		return management.ProjectGrantState_PROJECTGRANTSTATE_ACTIVE
	case proj_model.ProjectStateInactive:
		return management.ProjectGrantState_PROJECTGRANTSTATE_INACTIVE
	default:
		return management.ProjectGrantState_PROJECTGRANTSTATE_UNSPECIFIED
	}
}

func projectGrantsToIDs(projectGrants []*proj_model.ProjectGrantView) []string {
	converted := make([]string, len(projectGrants))
	for i, grant := range projectGrants {
		converted[i] = grant.GrantID
	}
	return converted
}
