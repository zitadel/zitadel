package crdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
)

var (
	errSeqNotUpdated = errors.ThrowInternal(nil, "CRDB-79GWt", "current sequence not updated")
)

type StatementHandlerConfig struct {
	handler.ProjectionHandlerConfig

	Client            *sql.DB
	SequenceTable     string
	LockTable         string
	FailedEventsTable string
	MaxFailureCount   uint
	BulkLimit         uint64

	Reducers  []handler.AggregateReducer
	InitCheck *handler.Check
}

type StatementHandler struct {
	*handler.ProjectionHandler
	Locker

	client                  *sql.DB
	sequenceTable           string
	currentSequenceStmt     string
	updateSequencesBaseStmt string
	maxFailureCount         uint
	failureCountStmt        string
	setFailureCountStmt     string

	aggregates []eventstore.AggregateType
	reduces    map[eventstore.EventType]handler.Reduce

	bulkLimit uint64
}

func NewStatementHandler(
	ctx context.Context,
	config StatementHandlerConfig,
) StatementHandler {
	aggregateTypes := make([]eventstore.AggregateType, 0, len(config.Reducers))
	reduces := make(map[eventstore.EventType]handler.Reduce, len(config.Reducers))
	for _, aggReducer := range config.Reducers {
		aggregateTypes = append(aggregateTypes, aggReducer.Aggregate)
		for _, eventReducer := range aggReducer.EventRedusers {
			reduces[eventReducer.Event] = eventReducer.Reduce
		}
	}

	h := StatementHandler{
		client:                  config.Client,
		sequenceTable:           config.SequenceTable,
		maxFailureCount:         config.MaxFailureCount,
		currentSequenceStmt:     fmt.Sprintf(latestEventStmtFormat, config.SequenceTable),
		updateSequencesBaseStmt: fmt.Sprintf(updateEventIDStmtFormat, config.SequenceTable),
		failureCountStmt:        fmt.Sprintf(failureCountStmtFormat, config.FailedEventsTable),
		setFailureCountStmt:     fmt.Sprintf(setFailureCountStmtFormat, config.FailedEventsTable),
		aggregates:              aggregateTypes,
		reduces:                 reduces,
		bulkLimit:               config.BulkLimit,
		Locker:                  NewLocker(config.Client, config.LockTable, config.ProjectionName),
	}

	initialized := make(chan bool)
	h.ProjectionHandler = handler.NewProjectionHandler(ctx, config.ProjectionHandlerConfig, h.reduce, h.Update, h.SearchQuery, h.Lock, h.Unlock, initialized)

	err := h.Init(ctx, initialized, config.InitCheck)
	logging.OnError(err).WithField("projection", config.ProjectionName).Fatal("unable to initialize projections")

	h.Subscribe(h.aggregates...)

	return h
}

func (h *StatementHandler) SearchQuery(ctx context.Context, instanceIDs []string) (*eventstore.SearchQueryBuilder, uint64, error) {
	sequences, err := h.currentSequences(ctx, h.client.QueryContext, instanceIDs)
	if err != nil {
		return nil, 0, err
	}

	queryBuilder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).Limit(h.bulkLimit)

	for _, aggregateType := range h.aggregates {
		for _, instanceID := range instanceIDs {
			var creationDate time.Time
			for _, sequence := range sequences[aggregateType] {
				if sequence.instanceID == instanceID {
					creationDate = sequence.creationDate
					break
				}
			}
			queryBuilder.
				AddQuery().
				AggregateTypes(aggregateType).
				CreationDateAfter(creationDate).
				InstanceID(instanceID)
		}
	}

	return queryBuilder, h.bulkLimit, nil
}

// Update implements handler.Update
func (h *StatementHandler) Update(ctx context.Context, stmts []*handler.Statement, reduce handler.Reduce) (index int, err error) {
	if len(stmts) == 0 {
		return -1, nil
	}
	instanceIDs := make([]string, 0, len(stmts))
	for _, stmt := range stmts {
		instanceIDs = appendToInstanceIDs(instanceIDs, stmt.InstanceID)
	}
	tx, err := h.client.BeginTx(ctx, nil)
	if err != nil {
		return -1, errors.ThrowInternal(err, "CRDB-e89Gq", "begin failed")
	}

	sequences, err := h.currentSequences(ctx, tx.QueryContext, instanceIDs)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	//checks for events between create statement and current sequence
	// because there could be events between current sequence and a creation event
	// and we cannot check via stmt.PreviousSequence
	// if stmts[0].PreviousEventDate.IsZero() {
	previousStmts, err := h.fetchPreviousStmts(ctx, tx, stmts[0].CreationDate, stmts[0].InstanceID, sequences, reduce)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	stmts = append(previousStmts, stmts...)
	// }

	lastSuccessfulIdx := h.executeStmts(tx, &stmts, sequences)

	if lastSuccessfulIdx >= 0 {
		err = h.updateCurrentSequences(tx, sequences)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	if lastSuccessfulIdx < len(stmts)-1 {
		return lastSuccessfulIdx, handler.ErrSomeStmtsFailed
	}

	return lastSuccessfulIdx, nil
}

func (h *StatementHandler) fetchPreviousStmts(ctx context.Context, tx *sql.Tx, stmtDate time.Time, instanceID string, sequences events, reduce handler.Reduce) (previousStmts []*handler.Statement, err error) {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).SetTx(tx)
	queriesAdded := false
	for _, aggregateType := range h.aggregates {
		for _, sequence := range sequences[aggregateType] {
			if sequence.creationDate.After(stmtDate) && instanceID == sequence.instanceID {
				// if stmtSeq. <= sequence.eventID && instanceID == sequence.instanceID {
				continue
			}

			query.
				AddQuery().
				AggregateTypes(aggregateType).
				CreationDateAfter(sequence.creationDate).
				// CreationDateBefore(stmtDate).
				InstanceID(sequence.instanceID)

			queriesAdded = true
		}
	}

	if !queriesAdded {
		return nil, nil
	}

	events, err := h.Eventstore.Filter(ctx, query)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		stmt, err := reduce(event)
		if err != nil {
			return nil, err
		}
		previousStmts = append(previousStmts, stmt)
	}
	return previousStmts, nil
}

func (h *StatementHandler) executeStmts(
	tx *sql.Tx,
	stmts *[]*handler.Statement,
	sequences events,
) int {

	lastSuccessfulIdx := -1
stmts:
	for i := 0; i < len(*stmts); i++ {
		stmt := (*stmts)[i]
		for _, sequence := range sequences[stmt.AggregateType] {
			if sequence.creationDate.After(stmt.CreationDate) && stmt.InstanceID == sequence.instanceID {
				// if stmt.Sequence <= sequence.eventID && stmt.InstanceID == sequence.instanceID {
				logging.WithFields("currentSequence", sequence).Debug("statement dropped")
				if i < len(*stmts)-1 {
					copy((*stmts)[i:], (*stmts)[i+1:])
				}
				*stmts = (*stmts)[:len(*stmts)-1]
				i--
				continue stmts
			}
			// if !stmt.PreviousEventDate.IsZero() && !stmt.PreviousEventDate.Equal(sequence.creationDate) && stmt.InstanceID == sequence.instanceID {
			// 	// if stmt.PreviousSequence > 0 && stmt.PreviousSequence != sequence.eventID && stmt.InstanceID == sequence.instanceID {
			// 	logging.WithFields("projection", h.ProjectionName, "aggregateType", stmt.AggregateType, "creationDate", stmt.CreationDate.String(), "prevCreationDate", stmt.PreviousEventDate, "creationDate", sequence.creationDate).Warn("sequences do not match")
			// 	break stmts
			// }
		}
		err := h.executeStmt(tx, stmt)
		if err == nil {
			updateSequences(sequences, stmt)
			lastSuccessfulIdx = i
			continue
		}

		shouldContinue := h.handleFailedStmt(tx, stmt, err)
		if !shouldContinue {
			break
		}

		updateSequences(sequences, stmt)
		lastSuccessfulIdx = i
		continue
	}
	return lastSuccessfulIdx
}

// executeStmt handles sql statements
// an error is returned if the statement could not be inserted properly
func (h *StatementHandler) executeStmt(tx *sql.Tx, stmt *handler.Statement) error {
	if stmt.IsNoop() {
		return nil
	}
	_, err := tx.Exec("SAVEPOINT push_stmt")
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-i1wp6", "unable to create savepoint")
	}
	err = stmt.Execute(tx, h.ProjectionName)
	if err != nil {
		_, rollbackErr := tx.Exec("ROLLBACK TO SAVEPOINT push_stmt")
		if rollbackErr != nil {
			return errors.ThrowInternal(rollbackErr, "CRDB-zzp3P", "rollback to savepoint failed")
		}
		return errors.ThrowInternal(err, "CRDB-oRkaN", "unable execute stmt")
	}
	_, err = tx.Exec("RELEASE push_stmt")
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-qWgwT", "unable to release savepoint")
	}
	return nil
}

func updateSequences(sequences events, stmt *handler.Statement) {
	for _, sequence := range sequences[stmt.AggregateType] {
		if sequence.instanceID == stmt.InstanceID {
			sequence.creationDate = stmt.CreationDate
			return
		}
	}
	sequences[stmt.AggregateType] = append(sequences[stmt.AggregateType], &instanceEvents{
		instanceID:   stmt.InstanceID,
		creationDate: stmt.CreationDate,
	})
}

func appendToInstanceIDs(instances []string, id string) []string {
	for _, instance := range instances {
		if instance == id {
			return instances
		}
	}
	return append(instances, id)
}
