package eventsourcing

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/changes/model"
	chg_model "github.com/caos/zitadel/internal/changes/model"
	chg_type "github.com/caos/zitadel/internal/changes/types"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_model "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/protobuf/ptypes"
)

func (es *ChangesEventstore) Changes(ctx context.Context, aggregateType es_model.AggregateType, id string, secId string, lastSequence uint64, limit uint64) (*chg_model.Changes, error) {
	aggregateTypeQuery := aggregateType
	if aggregateType == chg_model.Application {
		aggregateTypeQuery = chg_model.Project
	}
	query := ChangesQuery(id, lastSequence, aggregateTypeQuery)

	events, err := es.Eventstore.FilterEvents(context.Background(), query)
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-FpQqK", "no objects found")
	}
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-328b1", "unable to get current user")
	}

	result := make([]*model.Change, 0)

	for _, u := range events {
		creationDate, err := ptypes.TimestampProto(u.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &model.Change{
			ChangeDate: creationDate,
			EventType:  u.Type.String(),
			Modifier:   u.EditorUser,
			Sequence:   u.Sequence,
		}
		appendChanges := true

		switch u.AggregateType {
		case chg_model.User:
			userDummy := chg_type.User{}
			if u.Data != nil {
				if err := json.Unmarshal(u.Data, &userDummy); err != nil {
					log.Println("Error getting data!", err.Error())
				}
			}
			change.Data = userDummy
		case chg_model.Project:
			if aggregateType == chg_model.Project {
				logging.Log("Project").Debugln("Project")
				projectDummy := chg_type.Project{}
				appDummy := chg_type.App{}
				change.Data = projectDummy
				if u.Data != nil {
					if strings.Contains(change.EventType, "application") {
						if err := json.Unmarshal(u.Data, &appDummy); err != nil {
							log.Println("Error getting data!", err.Error())
						}
						change.Data = appDummy
					} else {
						if err := json.Unmarshal(u.Data, &projectDummy); err != nil {
							log.Println("Error getting data!", err.Error())
						}
						change.Data = projectDummy
					}
				}
			} else if aggregateType == chg_model.Application {
				if change.EventType == "project.application.added" ||
					change.EventType == "project.application.changed" ||
					change.EventType == "project.application.config.oidc.added" ||
					change.EventType == "project.application.config.oidc.changed" {
					appDummy := chg_type.App{}
					if u.Data != nil {
						if err := json.Unmarshal(u.Data, &appDummy); err != nil {
							log.Println("Error getting data!", err.Error())
						}
					}
					change.Data = appDummy
					if appDummy.AppId != secId {
						appendChanges = false
					}
				} else {
					appendChanges = false
				}
			}
		case chg_model.Org:
			orgDummy := chg_type.Org{}
			if u.Data != nil {
				if err := json.Unmarshal(u.Data, &orgDummy); err != nil {
					log.Println("Error getting data!", err.Error())
				}
			}
			change.Data = orgDummy
		}
		if appendChanges {
			result = append(result, change)
			if lastSequence < u.Sequence {
				lastSequence = u.Sequence
			}
		}
	}

	changes := &model.Changes{
		Changes:      result,
		LastSequence: lastSequence,
	}

	return changes, nil
}

func ChangesQuery(objectID string, latestSequence uint64, aggregateType es_model.AggregateType) *es_model.SearchQuery {
	query := es_model.NewSearchQuery().
		AggregateTypeFilter(aggregateType).
		LatestSequenceFilter(latestSequence).
		AggregateIDFilter(objectID)
	return query
}
