package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/zitadel/zitadel/internal/api/authz"
)

type EventStore struct {
	db *sql.DB
}

type Config struct {
	Client *sql.DB
}

func Start(config *Config) (*EventStore, error) {
	es := &EventStore{
		config.Client,
	}

	_, err := es.db.Exec(create)
	if err != nil {
		return nil, err
	}

	return es, nil
}

var (
	//go:embed insert_events.sql
	insertStmt string
	//go:embed mig.sql
	create string
	//go:embed query_sequences.sql
	sequenceQuery string
)

// Cmd is an abstraction of [Command]
// It requires to extend [Command]
type Cmd interface {
	// Type represents the action made
	Type() string
	// Version describes the object definition at the time the event was created.
	// The version counts up as soon as the action definition changes
	Version() uint8
	// Payload represents the data added or changed
	Payload() any
	// Editor represents who created the event
	Editor() *Editor
	// Aggregate represents a definition of an object
	Aggregate() *Aggregate
	// internal verifies that the [Command] object is extended
	internal()
}

const (
	eventArgsPerRow     = 10
	aggregateArgsPerRow = 3
)

func (e *EventStore) Push(ctx context.Context, commands ...Cmd) (events []*Event, err error) {
	eventArgs := make([]interface{}, 0, len(commands)*eventArgsPerRow)
	eventPlaceholders := make([]string, len(commands))
	events = make([]*Event, len(commands))

	err = crdb.ExecuteTx(ctx, e.db, nil, func(tx *sql.Tx) error {
		aggregates, err := e.queryAggregates(ctx, tx, commands)
		if err != nil {
			return err
		}
		for i, cmd := range commands {
			payload, err := payloadAsJSON(cmd)
			if err != nil {
				return err
			}
			aggregate := searchAggregate(aggregates, authz.GetInstance(ctx).InstanceID(), cmd.Aggregate().ID)
			aggregate.sequence++
			events[i] = &Event{
				Type:    cmd.Type(),
				Version: cmd.Version(),
				Payload: payload,
				Editor:  *cmd.Editor(),
				Aggregate: Aggregate{
					ID:         authz.GetInstance(ctx).InstanceID(),
					Type:       cmd.Aggregate().Type,
					Owner:      cmd.Aggregate().Owner,
					InstanceID: cmd.Aggregate().ID,
					sequence:   aggregate.sequence,
				},
			}

			eventArgs = append(eventArgs,
				events[i].Aggregate.ID,
				events[i].Aggregate.Type,
				events[i].Aggregate.Owner,
				events[i].Aggregate.InstanceID,
				events[i].Editor.UserID,
				events[i].Editor.Service,
				events[i].Type,
				events[i].Version,
				events[i].Payload,
				events[i].Aggregate.sequence,
			)

			eventPlaceholders[i] = "(" +
				"$" + strconv.Itoa(i*eventArgsPerRow+1) +
				", $" + strconv.Itoa(i*eventArgsPerRow+2) +
				", $" + strconv.Itoa(i*eventArgsPerRow+3) +
				", $" + strconv.Itoa(i*eventArgsPerRow+4) +
				", $" + strconv.Itoa(i*eventArgsPerRow+5) +
				", $" + strconv.Itoa(i*eventArgsPerRow+6) +
				", $" + strconv.Itoa(i*eventArgsPerRow+7) +
				", $" + strconv.Itoa(i*eventArgsPerRow+8) +
				", $" + strconv.Itoa(i*eventArgsPerRow+9) +
				", $" + strconv.Itoa(i*eventArgsPerRow+10) +
				")"
		}

		eventsStmt := fmt.Sprintf(insertStmt, strings.Join(eventPlaceholders, ", "))
		_, err = tx.Exec(eventsStmt, eventArgs...)

		return err
	})

	if err != nil {
		return nil, err
	}
	return events, nil
}

func (es *EventStore) queryAggregates(ctx context.Context, tx *sql.Tx, cmds []Cmd) ([]*Aggregate, error) {
	aggregates := make([]*Aggregate, 0, len(cmds))

	for _, cmd := range cmds {
		if searchAggregate(aggregates, cmd.Aggregate().InstanceID, cmd.Aggregate().ID) != nil {
			continue
		}
		aggregates = append(aggregates, &Aggregate{
			ID:         authz.GetInstance(ctx).InstanceID(),
			Type:       cmd.Aggregate().Type,
			Owner:      cmd.Aggregate().Owner,
			InstanceID: cmd.Aggregate().ID,
		})
	}

	args := make([]any, 0, len(aggregates)*2)
	conditions := make([]string, len(aggregates))
	for i, agg := range aggregates {
		conditions[i] = "(aggregate_id = $" + strconv.Itoa(i*2+1) + " AND instance_id = $" + strconv.Itoa(i*2+2) + ")"
		args = append(args, agg.ID, agg.InstanceID)
	}

	query := fmt.Sprintf(sequenceQuery, strings.Join(conditions, " OR "))
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var (
			aggID, instanceID string
			sequence          uint64
		)

		err := rows.Scan(&aggID, &instanceID, &sequence)
		if err != nil {
			return nil, err
		}

		agg := searchAggregate(aggregates, aggID, instanceID)
		agg.sequence = sequence
	}

	return aggregates, nil
}

func searchAggregate(aggregates []*Aggregate, id, instance string) *Aggregate {
	for _, aggregate := range aggregates {
		if aggregate.ID == id && aggregate.InstanceID == instance {
			return aggregate
		}
	}
	return nil
}

func (e *EventStore) Filter(ctx context.Context, filter *Filter) ([]*Event, error) {
	return nil, nil
}

// payloadAsJSON is used to encode the payload to json
func payloadAsJSON(c Cmd) ([]byte, error) {
	if c.Payload() == nil {
		return nil, nil
	}

	if payload, ok := c.Payload().([]byte); ok {
		if json.Valid(payload) {
			return payload, nil
		}
		return nil, errors.New("payload is not valid json")
	}

	return json.Marshal(c.Payload())
}
