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
	"time"

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
	//go:embed insert.sql
	insertStmt string
	//go:embed mig.sql
	create string
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
	argsPerRow = 9
)

func (e *EventStore) Push(ctx context.Context, commands ...Cmd) (events []*Event, err error) {
	args := make([]interface{}, 0, len(commands)*argsPerRow)
	placeholders := make([]string, len(commands))
	events = make([]*Event, len(commands))

	for i, cmd := range commands {
		payload, err := payloadAsJSON(cmd)
		if err != nil {
			return nil, err
		}
		events[i] = &Event{
			Type:      cmd.Type(),
			Version:   cmd.Version(),
			Payload:   payload,
			Editor:    *cmd.Editor(),
			Aggregate: *cmd.Aggregate(),
		}

		args = append(args,
			// TODO: replaced instance id with aggregate id to simulate multiple aggregates in a single instance
			authz.GetInstance(ctx).InstanceID(),
			events[i].Aggregate.Type,
			events[i].Aggregate.Owner,
			events[i].Aggregate.ID,
			events[i].Editor.UserID,
			events[i].Editor.Service,
			events[i].Type,
			events[i].Version,
			events[i].Payload,
		)

		placeholders[i] = "(" +
			"$" + strconv.Itoa(i*argsPerRow+1) +
			", $" + strconv.Itoa(i*argsPerRow+2) +
			", $" + strconv.Itoa(i*argsPerRow+3) +
			", $" + strconv.Itoa(i*argsPerRow+4) +
			", $" + strconv.Itoa(i*argsPerRow+5) +
			", $" + strconv.Itoa(i*argsPerRow+6) +
			", $" + strconv.Itoa(i*argsPerRow+7) +
			", $" + strconv.Itoa(i*argsPerRow+8) +
			", $" + strconv.Itoa(i*argsPerRow+9) +
			", now() + '" + fmt.Sprintf("%f", time.Duration(time.Microsecond*time.Duration(i)).Seconds()) + "s'" +
			")"
	}

	query := fmt.Sprintf(insertStmt, strings.Join(placeholders, ", "))
	err = crdb.ExecuteTx(ctx, e.db, nil, func(tx *sql.Tx) error {
		begin := time.Now()
		rows, err := tx.Query(query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for i := 0; rows.Next(); i++ {
			err = rows.Scan(&events[i].CreationDate)
			if err != nil {
				return err
			}
		}

		if remaining := begin.Sub(time.Now()).Microseconds() - int64(len(commands)); remaining > 0 {
			time.Sleep(time.Duration(remaining) * time.Microsecond)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return events, nil
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
