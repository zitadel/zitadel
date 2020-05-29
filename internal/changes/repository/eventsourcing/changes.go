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

func (es *ChangesEventstore) Changes(ctx context.Context, objectType es_model.AggregateType, objectID string, lastSequence uint64, limit uint64) (*chg_model.Changes, error) {

	query := ChangesQuery(objectID, lastSequence, objectType)

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
		result = append(result, change)
		if lastSequence < u.Sequence {
			lastSequence = u.Sequence
		}

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
		case chg_model.Application:
			type omitempty struct {
				ClientId string `json:"clientId,omitempty"`
			}
			type app struct {
				ClientId  string    `json:"clientId,omitempty"`
				Name      string    `json:"name,omitempty"`
				Omitempty omitempty `json:"omitempty,omitempty"`
				ProjectId string    `json:"projectId,omitempty"`
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
	}

	changes := &model.Changes{
		Changes:      result,
		LastSequence: lastSequence,
	}

	return changes, nil
}

func ChangesQuery(recourceOwner string, latestSequence uint64, aggregateType es_model.AggregateType) *es_model.SearchQuery {
	return es_model.NewSearchQuery().
		AggregateTypeFilter(aggregateType).
		LatestSequenceFilter(latestSequence).
		AggregateIDFilter(recourceOwner)
}
