package repository

import (
	"context"
	"time"

	"github.com/caos/utils/logging"
	"github.com/cockroachdb/cockroach-go/crdb"

	"database/sql"

	lib_models "github.com/caos/eventstore-lib/pkg/models"
	caos_errs "github.com/caos/utils/errors"
	"github.com/caos/utils/tracing"
	"github.com/caos/zitadel/internal/eventstore/models"
)

const (
	eventsTable = "eventstore.events"
)

type Event struct {
	ID               string    `sql:"column:id;default:NULL"`
	CreationDate     time.Time `sql:"column:creation_date;default:now()"`
	Typ              string    `sql:"column:event_type;default:NULL"`
	Sequence         uint64    `sql:"column:event_sequence;default:nextval('eventstore.event_seq')"`
	PreviousSequence uint64    `sql:"column:previous_sequence;default:NULL"`
	Data             []byte    `sql:"column:event_data;default:NULL"`
	ModifierService  string    `sql:"column:modifier_service;default:NULL"`
	ModifierTenant   string    `sql:"column:modifier_tenant;default:NULL"`
	ModiferUser      string    `sql:"column:modifier_user;default:NULL"`
	ResourceOwner    string    `sql:"column:resource_owner;default:NULL"`
	AggregateType    string    `sql:"column:aggregate_type"`
	AggregateID      string    `sql:"column:aggregate_id"`
}

func eventFromApp(event lib_models.Event, aggregate aggregateRoot) *Event {
	e := event.(*models.Event)

	data := e.Data()
	if data == nil || len(data) == 0 {
		//json decoder failes with EOF if json text is empty
		data = []byte("{}")
	}

	return &Event{
		CreationDate:     e.CreationDate(),
		Typ:              e.Type(),
		PreviousSequence: e.PreviousSequence,
		Data:             data,
		ModifierService:  e.ModifierService,
		ModifierTenant:   e.ModifierTenant,
		ModiferUser:      e.ModifierUser,
		ResourceOwner:    e.ResourceOwner,
		AggregateType:    aggregate.Type(),
		AggregateID:      aggregate.ID(),
	}
}

func updateAppEvent(event lib_models.Event, savedEvent *Event) {
	e, ok := event.(*models.Event)
	if !ok {
		logging.Log("SQL-h7K6x").Warn("cannot update event with saved parameters")
		return
	}
	e.Sequence = savedEvent.Sequence
	e.PreviousSequence = savedEvent.PreviousSequence
	e.AggregateType = savedEvent.AggregateType
	e.AggregateID = savedEvent.AggregateID
}

func eventToApp(event *Event) *models.Event {
	e := models.NewEvent(event.ID)
	e.SetCreationDate(event.CreationDate)
	e.SetType(event.Typ)
	e.SetData(event.Data)
	e.ModifierService = event.ModifierService
	e.ModifierTenant = event.ModifierTenant
	e.ModifierUser = event.ModiferUser
	e.ResourceOwner = event.ResourceOwner
	e.Sequence = event.Sequence
	e.PreviousSequence = event.PreviousSequence
	e.AggregateType = event.AggregateType
	e.AggregateID = event.AggregateID

	return e
}

func (db *SQL) PushEvents(ctx context.Context, aggregates ...lib_models.Aggregate) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	err = crdb.ExecuteTx(ctx, db.sqlClient, nil, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare("insert into eventstore.events " +
			"(event_type, aggregate_type, aggregate_id, creation_date, event_data, modifier_user, modifier_service, modifier_tenant, resource_owner, previous_sequence) " +
			"select $1, $2, $3, coalesce($4, now()), $5, $6, $7, $8, $9, " +
			// case is to set the highest sequence or NULL in previous_sequence
			"case (select exists(select event_sequence from eventstore.events where aggregate_type = $10 AND aggregate_id = $11)) " +
			"WHEN true then (select event_sequence from eventstore.events where aggregate_type = $12 AND aggregate_id = $13 order by event_sequence desc limit 1) " +
			"ELSE NULL " +
			"end " +
			"where (" +
			// exactly one event of requested aggregate must have a >= sequence (last inserted event)
			"(select count(id) from eventstore.events where event_sequence >= $14 AND aggregate_type = $15 AND aggregate_id = $16) = 1 OR " +
			// previous sequence = 0, no events must exist for the requested aggregate
			"((select count(id) from eventstore.events where aggregate_type = $17 and aggregate_id = $18) = 0 AND $19 = 0)) " +
			"RETURNING id, event_sequence, creation_date")
		if err != nil {
			tx.Rollback()
			logging.Log("SQL-9ctx5").WithError(err).Warn("prepare failed")
			return caos_errs.ThrowInternal(err, "SQL-juCgA", "prepare failed")
		}

		for _, aggregate := range aggregates {
			previousSequence := aggregate.LatestSequence()
			for _, event := range aggregate.Events().GetAll() {
				e := eventFromApp(event, aggregate)
				e.PreviousSequence = previousSequence

				rows, err := stmt.Query(e.Typ, e.AggregateType, e.AggregateID, e.CreationDate, e.Data, e.ModiferUser, e.ModifierService, e.ModifierTenant, e.ResourceOwner,
					e.AggregateType, e.AggregateID,
					e.AggregateType, e.AggregateID,
					e.PreviousSequence, e.AggregateType, e.AggregateID,
					e.AggregateType, e.AggregateID, e.PreviousSequence)

				if err != nil {
					logging.Log("SQL-EXA0q").WithError(err).Info("query failed")
					tx.Rollback()
					return caos_errs.ThrowInternal(err, "SQL-SBP37", "unable to create event")
				}
				defer rows.Close()

				rowInserted := false
				for rows.Next() {
					rowInserted = true
					err = rows.Scan(&e.ID, &e.Sequence, &e.CreationDate)
					logging.Log("SQL-rAvLD").OnError(err).Info("unable to scan result into event")
				}

				if !rowInserted {
					tx.Rollback()
					return caos_errs.ThrowAlreadyExists(nil, "SQL-GKcAa", "wrong sequence")
				}

				updateAppEvent(event, e)
				previousSequence = e.Sequence
			}
		}
		return nil
	})

	if _, ok := err.(*caos_errs.CaosError); !ok && err != nil {
		err = caos_errs.ThrowInternal(err, "SQL-DjgtG", "unable to store events")
	}

	return err
}
