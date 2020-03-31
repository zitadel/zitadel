package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es "github.com/caos/zitadel/internal/eventstore"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_es "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type ProjectEventstore struct {
	es.App
	aggregateCreator es.AggregateCreator
}

type ProjectConfig struct {
	es.App
	aggregateCreator es.AggregateCreator
}

func StartProject(conf ProjectConfig) (*ProjectEventstore, error) {
	return &ProjectEventstore{App: conf.App, aggregateCreator: conf.aggregateCreator}, nil
}

func (es *ProjectEventstore) ProjectByID(ctx context.Context, project *proj_model.Project) error {
	filter := proj_es.ProjectByIDQuery(project.ID, project.Sequence)
	events, err := es.App.FilterEvents(ctx, filter)
	if err != nil {
		return err
	}
	foundProject, err := proj_es.FromEvents(nil, events...)

	*project = *proj_es.ProjectToModel(foundProject)
	return err
}

func (es *ProjectEventstore) CreateProject(ctx context.Context, name string) (id string, err error) {
	project := proj_model.Project{Name: name}

	projectAggregate, err := proj_es.ProjectCreateEvents(ctx, es.aggregateCreator, &project)
	if err != nil {
		return "", err
	}
	_, err := es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return "", err
	}
	return projectAggregate.ID, nil
}

func (es *ProjectEventstore) UpdateProject(ctx context.Context, changes *proj_model.ProjectChange) (sequence uint64, err error) {
	orgAggregate := proj_es.ProjectUpdateEvents(changes)
	save := pkg.PrepareSave(ctx, orgAggregate, uniqueAggregates...)
	err = save(es.Client)
	if err != nil {
		return changes.Sequence, err
	}
	return orgAggregate.Sequence, nil
}

//
//func (es *ProjectEventstore) DeactivateOrg(ctx context.Context, orgID string, sequence uint64) (uint64, error) {
//	aggregate := org_es.OrgDeactivateEvents(orgID, sequence)
//
//	save := pkg.PrepareSave(ctx, aggregate)
//	err := save(es.Client)
//	return aggregate.Sequence, err
//}
//
//func (es *ProjectEventstore) ReactivateOrg(ctx context.Context, orgID string, sequence uint64) (uint64, error) {
//	aggregate := org_es.OrgReactivateEvents(orgID, sequence)
//
//	save := pkg.PrepareSave(ctx, aggregate)
//	err := save(es.Client)
//	return aggregate.Sequence, err
//}
