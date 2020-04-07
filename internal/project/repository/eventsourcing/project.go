package eventsourcing

import (
	"context"
	"strconv"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/sony/sonyflake"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

const (
	projectVersion = "v1"
)

type Project struct {
	es_models.ObjectRoot
	Name    string           `json:"name,omitempty"`
	State   int32            `json:"-"`
	Members []*ProjectMember `json:"-"`
}

type ProjectMember struct {
	es_models.ObjectRoot
	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func (p *Project) Changes(changed *Project) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func ProjectFromModel(project *model.Project) *Project {
	members := ProjectMembersFromModel(project.Members)
	return &Project{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ObjectRoot.ID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		Name:    project.Name,
		State:   model.ProjectStateToInt(project.State),
		Members: members,
	}
}

func ProjectToModel(project *Project) *model.Project {
	members := ProjectMembersToModel(project.Members)
	return &model.Project{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ID,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
			Sequence:     project.Sequence,
		},
		Name:    project.Name,
		State:   model.ProjectStateFromInt(project.State),
		Members: members,
	}
}

func ProjectMembersToModel(members []*ProjectMember) []*model.ProjectMember {
	convertedMembers := make([]*model.ProjectMember, len(members))
	for i, m := range members {
		convertedMembers[i] = ProjectMemberToModel(m)
	}
	return convertedMembers
}

func ProjectMembersFromModel(members []*model.ProjectMember) []*ProjectMember {
	convertedMembers := make([]*ProjectMember, len(members))
	for i, m := range members {
		convertedMembers[i] = ProjectMemberFromModel(m)
	}
	return convertedMembers
}

func ProjectMemberFromModel(member *model.ProjectMember) *ProjectMember {
	return &ProjectMember{
		ObjectRoot: es_models.ObjectRoot{
			ID:           member.ObjectRoot.ID,
			Sequence:     member.Sequence,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
		},
		UserID: member.UserID,
		Roles:  member.Roles,
	}
}

func ProjectMemberToModel(member *ProjectMember) *model.ProjectMember {
	return &model.ProjectMember{
		ObjectRoot: es_models.ObjectRoot{
			ID:           member.ID,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
			Sequence:     member.Sequence,
		},
		UserID: member.UserID,
		Roles:  member.Roles,
	}
}

func ProjectByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "id should be filled")
	}
	return ProjectQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func ProjectQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate).
		LatestSequenceFilter(latestSequence)
}

func ProjectAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, id string, sequence uint64) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, id, model.ProjectAggregate, projectVersion, sequence)
}

func ProjectCreateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, project *Project) (*es_models.Aggregate, error) {
	if project == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "project should not be nil")
	}
	var err error
	id, err := idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	project.ID = strconv.FormatUint(id, 10)

	agg, err := ProjectAggregate(ctx, aggCreator, project.ID, project.Sequence)
	if err != nil {
		return nil, err
	}

	return agg.AppendEvent(model.ProjectAdded, project)
}

func ProjectUpdateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project, new *Project) (*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
	}
	if new == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	changes := existing.Changes(new)
	return agg.AppendEvent(model.ProjectChanged, changes)
}

func ProjectDeactivateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project) (*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ueh45", "existing project should not be nil")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ProjectDeactivated, nil)
}

func ProjectReactivateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project) (*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-37dur", "existing project should not be nil")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ProjectReactivated, nil)
}

func ProjectMemberAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) (*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-di38f", "existing project should not be nil")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ProjectMemberAdded, member)
}

func ProjectMemberChangedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) (*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sle3d", "existing project should not be nil")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ProjectMemberChanged, member)
}

func ProjectMemberRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Project, member *ProjectMember) (*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-slo9e", "existing project should not be nil")
	}
	agg, err := ProjectAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	return agg.AppendEvent(model.ProjectMemberRemoved, member)
}
