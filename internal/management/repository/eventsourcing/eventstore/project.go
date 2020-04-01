package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_es "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type ProjectEventstore struct {
	es_int.Eventstore
}

type ProjectConfig struct {
	es_int.Eventstore
}

func StartProject(conf ProjectConfig) (*ProjectEventstore, error) {
	return &ProjectEventstore{Eventstore: conf.Eventstore}, nil
}

func (es *ProjectEventstore) ProjectByID(ctx context.Context, project *proj_model.Project) error {
	filter := proj_es.ProjectByIDQuery(project.ID, project.Sequence)
	events, err := es.Eventstore.FilterEvents(ctx, filter)
	if err != nil {
		return err
	}
	foundProject, err := proj_es.FromEvents(nil, events...)

	*project = *proj_es.ProjectToModel(foundProject)
	return err
}

func (es *ProjectEventstore) CreateProject(ctx context.Context, name string) (id string, err error) {
	project := proj_model.Project{Name: name}

	projectAggregate, err := proj_es.ProjectCreateEvents(ctx, es.Eventstore.AggregateCreator(), &project)
	if err != nil {
		return "", err
	}
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return "", err
	}
	return projectAggregate.ID, nil
}

func (es *ProjectEventstore) UpdateProject(ctx context.Context, changes *proj_model.ProjectChange) (sequence uint64, err error) {
	projectAggregate := proj_es.ProjectUpdateEvents(changes)
	save := pkg.PrepareSave(ctx, orgAggregate, uniqueAggregates...)
	err = save(es.Client)
	err = es.PushAggregates(ctx, projectAggregate)
	if err != nil {
		return 0, err
	}
	return projectAggregate, nil
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
