package query

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	exec "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	executionTable = table{
		name:          projection.ExecutionTable,
		instanceIDCol: projection.ExecutionInstanceIDCol,
	}
	ExecutionColumnID = Column{
		name:  projection.ExecutionIDCol,
		table: executionTable,
	}
	ExecutionColumnCreationDate = Column{
		name:  projection.ExecutionCreationDateCol,
		table: executionTable,
	}
	ExecutionColumnChangeDate = Column{
		name:  projection.ExecutionChangeDateCol,
		table: executionTable,
	}
	ExecutionColumnInstanceID = Column{
		name:  projection.ExecutionInstanceIDCol,
		table: executionTable,
	}
	executionTargetsTable = table{
		name:          projection.ExecutionTable + "_" + projection.ExecutionTargetSuffix,
		instanceIDCol: projection.ExecutionTargetInstanceIDCol,
	}
	executionTargetsTableAlias       = executionTargetsTable.setAlias("execution_targets")
	ExecutionTargetsColumnInstanceID = Column{
		name:  projection.ExecutionTargetInstanceIDCol,
		table: executionTargetsTableAlias,
	}
	ExecutionTargetsColumnExecutionID = Column{
		name:  projection.ExecutionTargetExecutionIDCol,
		table: executionTargetsTableAlias,
	}
	executionTargetsListCol = Column{
		name:  "targets",
		table: executionTargetsTableAlias,
	}
)

var (
	//go:embed execution_targets.sql
	executionTargetsQuery string
	//go:embed targets_by_execution_id.sql
	TargetsByExecutionIDQuery string
	//go:embed targets_by_execution_ids.sql
	TargetsByExecutionIDsQuery string
)

type Executions struct {
	SearchResponse
	Executions []*Execution
}

func (e *Executions) SetState(s *State) {
	e.State = s
}

type Execution struct {
	domain.ObjectDetails

	Targets []*exec.Target
}

type ExecutionSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *ExecutionSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchExecutions(ctx context.Context, queries *ExecutionSearchQueries) (executions *Executions, err error) {
	eq := sq.Eq{
		ExecutionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareExecutionsQuery(ctx, q.client)
	return genericRowsQueryWithState[*Executions](ctx, q.client, executionTable, combineToWhereStmt(query, queries.toQuery, eq), scan)
}

func (q *Queries) GetExecutionByID(ctx context.Context, id string) (execution *Execution, err error) {
	eq := sq.Eq{
		ExecutionColumnID.identifier():         id,
		ExecutionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareExecutionQuery(ctx, q.client)
	return genericRowQuery[*Execution](ctx, q.client, query.Where(eq), scan)
}

func NewExecutionInIDsSearchQuery(values []string) (SearchQuery, error) {
	return NewInTextQuery(ExecutionColumnID, values)
}

func NewExecutionTypeSearchQuery(t domain.ExecutionType) (SearchQuery, error) {
	return NewTextQuery(ExecutionColumnID, t.String(), TextStartsWith)
}

func NewTargetSearchQuery(target string) (SearchQuery, error) {
	data, err := targetItemJSONB(domain.ExecutionTargetTypeTarget, target)
	if err != nil {
		return nil, err
	}
	return NewListContains(executionTargetsListCol, data)
}

func NewIncludeSearchQuery(include string) (SearchQuery, error) {
	data, err := targetItemJSONB(domain.ExecutionTargetTypeInclude, include)
	if err != nil {
		return nil, err
	}
	return NewListContains(executionTargetsListCol, data)
}

// marshall executionTargets into the same JSONB structure as in the SQL queries
func targetItemJSONB(t domain.ExecutionTargetType, targetItem string) ([]byte, error) {
	var target *executionTarget
	switch t {
	case domain.ExecutionTargetTypeTarget:
		target = &executionTarget{Target: targetItem}
	case domain.ExecutionTargetTypeInclude:
		target = &executionTarget{Include: targetItem}
	case domain.ExecutionTargetTypeUnspecified:
		return nil, nil
	default:
		return nil, nil
	}
	return json.Marshal([]*executionTarget{target})
}

// TargetsByExecutionID query list of targets for best match of a list of IDs,  for example:
// [ "request/zitadel.action.v3alpha.ActionService/GetTargetByID",
// "request/zitadel.action.v3alpha.ActionService",
// "request" ]
func (q *Queries) TargetsByExecutionID(ctx context.Context, ids []string) (execution []*ExecutionTarget, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	if instanceID == "" {
		return nil, nil
	}

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			execution, err = scanExecutionTargets(rows)
			return err
		},
		TargetsByExecutionIDQuery,
		instanceID,
		database.TextArray[string](ids),
	)
	for i := range execution {
		if err := execution[i].decryptSigningKey(q.targetEncryptionAlgorithm); err != nil {
			return nil, err
		}
	}
	return execution, err
}

// TargetsByExecutionIDs query list of targets for best matches of 2 separate lists of IDs, combined for performance, for example:
// [ "request/zitadel.action.v3alpha.ActionService/GetTargetByID",
// "request/zitadel.action.v3alpha.ActionService",
// "request" ]
// and
// [ "response/zitadel.action.v3alpha.ActionService/GetTargetByID",
// "response/zitadel.action.v3alpha.ActionService",
// "response" ]
func (q *Queries) TargetsByExecutionIDs(ctx context.Context, ids1, ids2 []string) (execution []*ExecutionTarget, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	if instanceID == "" {
		return nil, nil
	}

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			execution, err = scanExecutionTargets(rows)
			return err
		},
		TargetsByExecutionIDsQuery,
		instanceID,
		database.TextArray[string](ids1),
		database.TextArray[string](ids2),
	)
	for i := range execution {
		if err := execution[i].decryptSigningKey(q.targetEncryptionAlgorithm); err != nil {
			return nil, err
		}
	}
	return execution, err
}

func prepareExecutionQuery(context.Context, prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) (*Execution, error)) {
	return sq.Select(
			ExecutionColumnInstanceID.identifier(),
			ExecutionColumnID.identifier(),
			ExecutionColumnCreationDate.identifier(),
			ExecutionColumnChangeDate.identifier(),
			executionTargetsListCol.identifier(),
		).From(executionTable.identifier()).
			Join("(" + executionTargetsQuery + ") AS " + executionTargetsTableAlias.alias + " ON " +
				ExecutionTargetsColumnInstanceID.identifier() + " = " + ExecutionColumnInstanceID.identifier() + " AND " +
				ExecutionTargetsColumnExecutionID.identifier() + " = " + ExecutionColumnID.identifier(),
			).
			PlaceholderFormat(sq.Dollar),
		scanExecution
}

func prepareExecutionsQuery(context.Context, prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (*Executions, error)) {
	return sq.Select(
			ExecutionColumnInstanceID.identifier(),
			ExecutionColumnID.identifier(),
			ExecutionColumnCreationDate.identifier(),
			ExecutionColumnChangeDate.identifier(),
			executionTargetsListCol.identifier(),
			countColumn.identifier(),
		).From(executionTable.identifier()).
			Join("(" + executionTargetsQuery + ") AS " + executionTargetsTableAlias.alias + " ON " +
				ExecutionTargetsColumnInstanceID.identifier() + " = " + ExecutionColumnInstanceID.identifier() + " AND " +
				ExecutionTargetsColumnExecutionID.identifier() + " = " + ExecutionColumnID.identifier(),
			).
			PlaceholderFormat(sq.Dollar),
		scanExecutions
}

type executionTarget struct {
	Position int    `json:"position,omitempty"`
	Include  string `json:"include,omitempty"`
	Target   string `json:"target,omitempty"`
}

func scanExecution(row *sql.Row) (*Execution, error) {
	execution := new(Execution)
	targets := make([]byte, 0)

	err := row.Scan(
		&execution.ResourceOwner,
		&execution.ID,
		&execution.CreationDate,
		&execution.EventDate,
		&targets,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "QUERY-qzn1xycesh", "Errors.Execution.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-f8sjvm4tb8", "Errors.Internal")
	}

	executionTargets := make([]*executionTarget, 0)
	if err := json.Unmarshal(targets, &executionTargets); err != nil {
		return nil, err
	}

	execution.Targets = make([]*exec.Target, len(executionTargets))
	for i := range executionTargets {
		if executionTargets[i].Target != "" {
			execution.Targets[i] = &exec.Target{Type: domain.ExecutionTargetTypeTarget, Target: executionTargets[i].Target}
		}
		if executionTargets[i].Include != "" {
			execution.Targets[i] = &exec.Target{Type: domain.ExecutionTargetTypeInclude, Target: executionTargets[i].Include}
		}
	}

	return execution, nil
}

func executionTargetsUnmarshal(data []byte) ([]*exec.Target, error) {
	executionTargets := make([]*executionTarget, 0)
	if err := json.Unmarshal(data, &executionTargets); err != nil {
		return nil, err
	}

	targets := make([]*exec.Target, len(executionTargets))
	// position starts with 1
	for _, item := range executionTargets {
		if item.Target != "" {
			targets[item.Position-1] = &exec.Target{Type: domain.ExecutionTargetTypeTarget, Target: item.Target}
		}
		if item.Include != "" {
			targets[item.Position-1] = &exec.Target{Type: domain.ExecutionTargetTypeInclude, Target: item.Include}
		}
	}
	return targets, nil
}

func scanExecutions(rows *sql.Rows) (*Executions, error) {
	executions := make([]*Execution, 0)
	var count uint64

	for rows.Next() {
		execution := new(Execution)
		targets := make([]byte, 0)

		err := rows.Scan(
			&execution.ResourceOwner,
			&execution.ID,
			&execution.CreationDate,
			&execution.EventDate,
			&targets,
			&count,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, zerrors.ThrowNotFound(err, "QUERY-tbrmno85vp", "Errors.Execution.NotFound")
			}
			return nil, zerrors.ThrowInternal(err, "QUERY-tyw2ydsj84", "Errors.Internal")
		}

		execution.Targets, err = executionTargetsUnmarshal(targets)
		if err != nil {
			return nil, err
		}
		executions = append(executions, execution)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-yhka3fs3mw", "Errors.Query.CloseRows")
	}

	return &Executions{
		Executions: executions,
		SearchResponse: SearchResponse{
			Count: count,
		},
	}, nil
}

type ExecutionTarget struct {
	InstanceID       string
	ExecutionID      string
	TargetID         string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
	signingKey       *crypto.CryptoValue
	SigningKey       string
}

func (e *ExecutionTarget) GetExecutionID() string {
	return e.ExecutionID
}
func (e *ExecutionTarget) GetTargetID() string {
	return e.TargetID
}
func (e *ExecutionTarget) IsInterruptOnError() bool {
	return e.InterruptOnError
}
func (e *ExecutionTarget) GetEndpoint() string {
	return e.Endpoint
}
func (e *ExecutionTarget) GetTargetType() domain.TargetType {
	return e.TargetType
}
func (e *ExecutionTarget) GetTimeout() time.Duration {
	return e.Timeout
}
func (e *ExecutionTarget) GetSigningKey() string {
	return e.SigningKey
}

func (t *ExecutionTarget) decryptSigningKey(alg crypto.EncryptionAlgorithm) error {
	if t.signingKey == nil {
		return nil
	}
	keyValue, err := crypto.DecryptString(t.signingKey, alg)
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-bxevy3YXwy", "Errors.Internal")
	}
	t.SigningKey = keyValue
	return nil
}

func scanExecutionTargets(rows *sql.Rows) ([]*ExecutionTarget, error) {
	targets := make([]*ExecutionTarget, 0)
	for rows.Next() {
		target := new(ExecutionTarget)

		var (
			instanceID       = &sql.NullString{}
			executionID      = &sql.NullString{}
			targetID         = &sql.NullString{}
			targetType       = &sql.NullInt32{}
			endpoint         = &sql.NullString{}
			timeout          = &sql.NullInt64{}
			interruptOnError = &sql.NullBool{}
			signingKey       = &crypto.CryptoValue{}
		)

		err := rows.Scan(
			executionID,
			instanceID,
			targetID,
			targetType,
			endpoint,
			timeout,
			interruptOnError,
			signingKey,
		)

		if err != nil {
			return nil, err
		}

		target.InstanceID = instanceID.String
		target.ExecutionID = executionID.String
		target.TargetID = targetID.String
		target.TargetType = domain.TargetType(targetType.Int32)
		target.Endpoint = endpoint.String
		target.Timeout = time.Duration(timeout.Int64)
		target.InterruptOnError = interruptOnError.Bool
		target.signingKey = signingKey

		targets = append(targets, target)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-37ardr0pki", "Errors.Query.CloseRows")
	}

	return targets, nil
}
