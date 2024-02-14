package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	projectsTable = table{
		name:          projection.ProjectProjectionTable,
		instanceIDCol: projection.ProjectColumnInstanceID,
	}
	ProjectColumnID = Column{
		name:  projection.ProjectColumnID,
		table: projectsTable,
	}
	ProjectColumnName = Column{
		name:  projection.ProjectColumnName,
		table: projectsTable,
	}
	ProjectColumnProjectRoleAssertion = Column{
		name:  projection.ProjectColumnProjectRoleAssertion,
		table: projectsTable,
	}
	ProjectColumnProjectRoleCheck = Column{
		name:  projection.ProjectColumnProjectRoleCheck,
		table: projectsTable,
	}
	ProjectColumnHasProjectCheck = Column{
		name:  projection.ProjectColumnHasProjectCheck,
		table: projectsTable,
	}
	ProjectColumnPrivateLabelingSetting = Column{
		name:  projection.ProjectColumnPrivateLabelingSetting,
		table: projectsTable,
	}
	ProjectColumnCreationDate = Column{
		name:  projection.ProjectColumnCreationDate,
		table: projectsTable,
	}
	ProjectColumnChangeDate = Column{
		name:  projection.ProjectColumnChangeDate,
		table: projectsTable,
	}
	ProjectColumnResourceOwner = Column{
		name:  projection.ProjectColumnResourceOwner,
		table: projectsTable,
	}
	ProjectColumnInstanceID = Column{
		name:  projection.ProjectColumnInstanceID,
		table: projectsTable,
	}
	ProjectColumnSequence = Column{
		name:  projection.ProjectColumnSequence,
		table: projectsTable,
	}
	ProjectColumnState = Column{
		name:  projection.ProjectColumnState,
		table: projectsTable,
	}
)

type Projects struct {
	SearchResponse
	Projects []*Project
}

type Project struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.ProjectState
	Sequence      uint64

	Name                   string
	ProjectRoleAssertion   bool
	ProjectRoleCheck       bool
	HasProjectCheck        bool
	PrivateLabelingSetting domain.PrivateLabelingSetting
}

type ProjectSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) ProjectByID(ctx context.Context, shouldTriggerBulk bool, id string) (project *Project, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerProjectProjection")
		ctx, err = projection.ProjectProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareProjectQuery(ctx, q.client)
	eq := sq.Eq{
		ProjectColumnID.identifier():         id,
		ProjectColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-2m00Q", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		project, err = scan(row)
		return err
	}, query, args...)
	return project, err
}

func (q *Queries) SearchProjects(ctx context.Context, queries *ProjectSearchQueries) (projects *Projects, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareProjectsQuery(ctx, q.client)
	eq := sq.Eq{ProjectColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-fn9ew", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		projects, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-2j00f", "Errors.Internal")
	}
	projects.State, err = q.latestState(ctx, projectsTable)
	return projects, err
}

func NewProjectNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnName, value, method)
}

func NewProjectIDSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(ProjectColumnID, list, ListIn)
}

func NewProjectResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnResourceOwner, value, TextEquals)
}

func (r *ProjectSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *ProjectSearchQueries) AppendPermissionQueries(permissions []string) error {
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetAllPermissionCtxIDs(permissions)
		query, err := NewProjectIDSearchQuery(ids)
		if err != nil {
			return err
		}
		r.Queries = append(r.Queries, query)
	}
	return nil
}

func (q *ProjectSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareProjectQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Project, error)) {
	return sq.Select(
			ProjectColumnID.identifier(),
			ProjectColumnCreationDate.identifier(),
			ProjectColumnChangeDate.identifier(),
			ProjectColumnResourceOwner.identifier(),
			ProjectColumnState.identifier(),
			ProjectColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectColumnProjectRoleAssertion.identifier(),
			ProjectColumnProjectRoleCheck.identifier(),
			ProjectColumnHasProjectCheck.identifier(),
			ProjectColumnPrivateLabelingSetting.identifier()).
			From(projectsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Project, error) {
			p := new(Project)
			err := row.Scan(
				&p.ID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.State,
				&p.Sequence,
				&p.Name,
				&p.ProjectRoleAssertion,
				&p.ProjectRoleCheck,
				&p.HasProjectCheck,
				&p.PrivateLabelingSetting,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-fk2fs", "Errors.Project.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-dj2FF", "Errors.Internal")
			}
			return p, nil
		}
}

func prepareProjectsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Projects, error)) {
	return sq.Select(
			ProjectColumnID.identifier(),
			ProjectColumnCreationDate.identifier(),
			ProjectColumnChangeDate.identifier(),
			ProjectColumnResourceOwner.identifier(),
			ProjectColumnState.identifier(),
			ProjectColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectColumnProjectRoleAssertion.identifier(),
			ProjectColumnProjectRoleCheck.identifier(),
			ProjectColumnHasProjectCheck.identifier(),
			ProjectColumnPrivateLabelingSetting.identifier(),
			countColumn.identifier()).
			From(projectsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Projects, error) {
			projects := make([]*Project, 0)
			var count uint64
			for rows.Next() {
				project := new(Project)
				err := rows.Scan(
					&project.ID,
					&project.CreationDate,
					&project.ChangeDate,
					&project.ResourceOwner,
					&project.State,
					&project.Sequence,
					&project.Name,
					&project.ProjectRoleAssertion,
					&project.ProjectRoleCheck,
					&project.HasProjectCheck,
					&project.PrivateLabelingSetting,
					&count,
				)
				if err != nil {
					return nil, err
				}
				projects = append(projects, project)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-QMXJv", "Errors.Query.CloseRows")
			}

			return &Projects{
				Projects: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
