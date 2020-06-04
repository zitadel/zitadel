package eventsourcing

import (
	"context"
	"encoding/json"
	"log"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/changes/model"
	chg_model "github.com/caos/zitadel/internal/changes/model"
	"github.com/caos/zitadel/internal/errors"
	es_model "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/protobuf/ptypes"
)

func (es *ChangesEventstore) Changes(ctx context.Context, aggregateType es_model.AggregateType, id string, lastSequence uint64, limit uint64) (*chg_model.Changes, error) {
	aggregateTypeQuery := aggregateType
	idQuery := id
	if aggregateType == chg_model.Application {
		aggregateTypeQuery = chg_model.Project
		idQuery = ""
	}
	query := ChangesQuery(idQuery, lastSequence, aggregateTypeQuery)

	events, err := es.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-SUOQ8z").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-Uw6zSS", "unable to get current user")
	}

	result := make([]*model.Change, 0)

	for _, u := range events {
		creationDate, err := ptypes.TimestampProto(u.CreationDate)
		logging.Log("GRPC-8duwe").OnError(err).Debug("unable to parse timestamp")
		change := &model.Change{
			ChangeDate: creationDate,
			EventType:  u.Type.String(),
			Modifier:   u.EditorUser,
			Sequence:   u.Sequence,
		}
		appendChanges := true

		switch u.AggregateType {
		case chg_model.User:
			type user struct {
				FirstName    string `json:"firstName,omitempty"`
				LastName     string `json:"lastName,omitempty"`
				EMailAddress string `json:"email,omitempty"`
				Phone        string `json:"phone,omitempty"`
			}
			userDummy := user{}
			if u.Data != nil {
				if err := json.Unmarshal(u.Data, &userDummy); err != nil {
					log.Println("Error getting data!", err.Error())
				}
			}
			change.Data = userDummy
		case chg_model.Project:
			if aggregateType == chg_model.Project {
				logging.Log("Project").Debugln("Project")
				type project struct {
					Name string `json:"name,omitempty"`
				}
				projectDummy := project{}
				if u.Data != nil {
					if err := json.Unmarshal(u.Data, &projectDummy); err != nil {
						log.Println("Error getting data!", err.Error())
					}
				}
				change.Data = projectDummy
			} else if aggregateType == chg_model.Application {
				if change.EventType == "project.application.added" ||
					change.EventType == "project.application.changed" ||
					change.EventType == "project.application.config.oidc.added" ||
					change.EventType == "project.application.config.oidc.changed" {
					type omitempty struct {
						ClientId string `json:"clientId,omitempty"`
					}
					type app struct {
						ClientId  string    `json:"clientId,omitempty"`
						Name      string    `json:"name,omitempty"`
						Omitempty omitempty `json:"omitempty,omitempty"`
						AppId     string    `json:"AppId,omitempty"`
						AppType   int       `json:"AppType,omitempty"`
					}
					appDummy := app{}
					omitemptyDummy := omitempty{}
					appDummy.Omitempty = omitemptyDummy
					if u.Data != nil {
						if err := json.Unmarshal(u.Data, &appDummy); err != nil {
							log.Println("Error getting data!", err.Error())
						}
					}
					change.Data = appDummy
					if appDummy.AppId != id {
						appendChanges = false
					}
				} else {
					appendChanges = false
				}
			}
		case chg_model.Org:
			type org struct {
				Domain string   `json:"domain,omitempty"`
				Name   string   `json:"name,omitempty"`
				Roles  []string `json:"roles,omitempty"`
				UserId string   `json:"userId,omitempty"`
			}
			orgDummy := org{}
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
		LatestSequenceFilter(latestSequence)
	if objectID != "" {
		query = query.AggregateIDFilter(objectID)
	}
	return query
}
