package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ProjectMetadataList struct {
	SearchResponse
	Metadata []*ProjectMetadata
}

type ProjectMetadata struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	ProjectID     string
	ResourceOwner string
	Sequence      uint64
	Key           string
	Value         []byte
}

type ProjectMetadataSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

var (
	projectMetadataTable = table{
		name:          projection.ProjectMetadataProjectionTable,
		instanceIDCol: projection.ProjectMetadataColumnInstanceID,
	}
	ProjectMetadataProjectIDCol = Column{
		name:  projection.ProjectMetadataColumnProjectID,
		table: projectMetadataTable,
	}
	ProjectMetadataCreationDateCol = Column{
		name:  projection.ProjectMetadataColumnCreationDate,
		table: projectMetadataTable,
	}
	ProjectMetadataChangeDateCol = Column{
		name:  projection.ProjectMetadataColumnChangeDate,
		table: projectMetadataTable,
	}
	ProjectMetadataResourceOwnerCol = Column{
		name:  projection.ProjectMetadataColumnResourceOwner,
		table: projectMetadataTable,
	}
	ProjectMetadataInstanceIDCol = Column{
		name:  projection.ProjectMetadataColumnInstanceID,
		table: projectMetadataTable,
	}
	ProjectMetadataSequenceCol = Column{
		name:  projection.ProjectMetadataColumnSequence,
		table: projectMetadataTable,
	}
	ProjectMetadataKeyCol = Column{
		name:  projection.ProjectMetadataColumnKey,
		table: projectMetadataTable,
	}
	ProjectMetadataValueCol = Column{
		name:  projection.ProjectMetadataColumnValue,
		table: projectMetadataTable,
	}
)

func (q *Queries) GetProjectMetadataByKey(ctx context.Context, shouldTriggerBulk bool, projectID string, key string, queries ...SearchQuery) (metadata *ProjectMetadata, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerProjectMetadataProjection")
		ctx, err = projection.ProjectMetadataProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	query, scan := prepareProjectMetadataQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		ProjectMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
		ProjectMetadataProjectIDCol.identifier():  projectID,
		ProjectMetadataKeyCol.identifier():        key,
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Ak4gcW", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		metadata, err = scan(row)
		return err
	}, stmt, args...)
	return metadata, err
}

func (q *Queries) SearchProjectMetadata(ctx context.Context, shouldTriggerBulk bool, projectID string, queries *ProjectMetadataSearchQueries) (metadata *ProjectMetadataList, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerProjectMetadataProjection")
		ctx, err = projection.ProjectMetadataProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}
	eq := sq.Eq{
		ProjectMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
		ProjectMetadataProjectIDCol.identifier():  projectID,
	}
	query, scan := prepareProjectMetadataListQuery()
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-zWKyX9", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		metadata, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-7C4yTM", "Errors.Internal")
	}

	metadata.State, err = q.latestState(ctx, projectMetadataTable)
	return metadata, err
}

func (q *ProjectMetadataSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (r *ProjectMetadataSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectMetadataResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func NewProjectMetadataResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectMetadataResourceOwnerCol, value, TextEquals)
}

func NewProjectMetadataKeySearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(ProjectMetadataKeyCol, value, comparison)
}

func prepareProjectMetadataQuery() (sq.SelectBuilder, func(*sql.Row) (*ProjectMetadata, error)) {
	return sq.Select(
			ProjectMetadataCreationDateCol.identifier(),
			ProjectMetadataChangeDateCol.identifier(),
			ProjectMetadataResourceOwnerCol.identifier(),
			ProjectMetadataSequenceCol.identifier(),
			ProjectMetadataKeyCol.identifier(),
			ProjectMetadataValueCol.identifier(),
		).
			From(projectMetadataTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*ProjectMetadata, error) {
			m := new(ProjectMetadata)
			err := row.Scan(
				&m.CreationDate,
				&m.ChangeDate,
				&m.ResourceOwner,
				&m.Sequence,
				&m.Key,
				&m.Value,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-ZLAE05", "Errors.Metadata.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-4ObU3Q", "Errors.Internal")
			}
			return m, nil
		}
}

func prepareProjectMetadataListQuery() (sq.SelectBuilder, func(*sql.Rows) (*ProjectMetadataList, error)) {
	return sq.Select(
			ProjectMetadataCreationDateCol.identifier(),
			ProjectMetadataChangeDateCol.identifier(),
			ProjectMetadataResourceOwnerCol.identifier(),
			ProjectMetadataSequenceCol.identifier(),
			ProjectMetadataKeyCol.identifier(),
			ProjectMetadataValueCol.identifier(),
			countColumn.identifier()).
			From(projectMetadataTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*ProjectMetadataList, error) {
			metadata := make([]*ProjectMetadata, 0)
			var count uint64
			for rows.Next() {
				m := new(ProjectMetadata)
				err := rows.Scan(
					&m.CreationDate,
					&m.ChangeDate,
					&m.ResourceOwner,
					&m.Sequence,
					&m.Key,
					&m.Value,
					&count,
				)
				if err != nil {
					return nil, err
				}

				metadata = append(metadata, m)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-XUeGQL", "Errors.Query.CloseRows")
			}

			return &ProjectMetadataList{
				Metadata: metadata,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
