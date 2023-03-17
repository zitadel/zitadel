package query

import (
	"context"
	"database/sql"
	errs "errors"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	InstancesFilterTableAlias = "f"
)

var (
	instanceTable = table{
		name:          projection.InstanceProjectionTable,
		instanceIDCol: projection.InstanceColumnID,
	}
	InstanceColumnID = Column{
		name:  projection.InstanceColumnID,
		table: instanceTable,
	}
	InstanceColumnName = Column{
		name:  projection.InstanceColumnName,
		table: instanceTable,
	}
	InstanceColumnCreationDate = Column{
		name:  projection.InstanceColumnCreationDate,
		table: instanceTable,
	}
	InstanceColumnChangeDate = Column{
		name:  projection.InstanceColumnChangeDate,
		table: instanceTable,
	}
	InstanceColumnSequence = Column{
		name:  projection.InstanceColumnSequence,
		table: instanceTable,
	}
	InstanceColumnDefaultOrgID = Column{
		name:  projection.InstanceColumnDefaultOrgID,
		table: instanceTable,
	}
	InstanceColumnProjectID = Column{
		name:  projection.InstanceColumnProjectID,
		table: instanceTable,
	}
	InstanceColumnConsoleID = Column{
		name:  projection.InstanceColumnConsoleID,
		table: instanceTable,
	}
	InstanceColumnConsoleAppID = Column{
		name:  projection.InstanceColumnConsoleAppID,
		table: instanceTable,
	}
	InstanceColumnDefaultLanguage = Column{
		name:  projection.InstanceColumnDefaultLanguage,
		table: instanceTable,
	}
)

type Instance struct {
	ID           string
	ChangeDate   time.Time
	CreationDate time.Time
	Sequence     uint64
	Name         string

	DefaultOrgID string
	IAMProjectID string
	ConsoleID    string
	ConsoleAppID string
	DefaultLang  language.Tag
	Domains      []*InstanceDomain
	host         string
	csp          csp
}

type csp struct {
	enabled        bool
	allowedOrigins database.StringArray
}

type Instances struct {
	SearchResponse
	Instances []*Instance
}

func (i *Instance) InstanceID() string {
	return i.ID
}

func (i *Instance) ProjectID() string {
	return i.IAMProjectID
}

func (i *Instance) ConsoleClientID() string {
	return i.ConsoleID
}

func (i *Instance) ConsoleApplicationID() string {
	return i.ConsoleAppID
}

func (i *Instance) RequestedDomain() string {
	return strings.Split(i.host, ":")[0]
}

func (i *Instance) RequestedHost() string {
	return i.host
}

func (i *Instance) DefaultLanguage() language.Tag {
	return i.DefaultLang
}

func (i *Instance) DefaultOrganisationID() string {
	return i.DefaultOrgID
}

func (i *Instance) SecurityPolicyAllowedOrigins() []string {
	if !i.csp.enabled {
		return nil
	}
	return i.csp.allowedOrigins
}

type InstanceSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewInstanceIDsListSearchQuery(ids ...string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(InstanceColumnID, list, ListIn)
}

func (q *InstanceSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchInstances(ctx context.Context, queries *InstanceSearchQueries) (instances *Instances, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	filter, query, scan := prepareInstancesQuery(ctx, q.client)
	stmt, args, err := query(queries.toQuery(filter)).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-M9fow", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-3j98f", "Errors.Internal")
	}
	instances, err = scan(rows)
	if err != nil {
		return nil, err
	}
	return instances, err
}

func (q *Queries) Instance(ctx context.Context, shouldTriggerBulk bool) (_ *Instance, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.InstanceProjection.Trigger(ctx)
	}

	stmt, scan := prepareInstanceDomainQuery(ctx, q.client, authz.GetInstance(ctx).RequestedDomain())
	query, args, err := stmt.Where(sq.Eq{
		InstanceColumnID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-d9ngs", "Errors.Query.SQLStatement")
	}

	row, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return scan(row)
}

func (q *Queries) InstanceByHost(ctx context.Context, host string) (_ authz.Instance, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAuthzInstanceQuery(ctx, q.client, host)
	host = strings.Split(host, ":")[0] //remove possible port
	query, args, err := stmt.Where(sq.Eq{
		InstanceDomainDomainCol.identifier(): host,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-SAfg2", "Errors.Query.SQLStatement")
	}

	row, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return scan(row)
}

func (q *Queries) InstanceByID(ctx context.Context) (_ authz.Instance, err error) {
	return q.Instance(ctx, true)
}

func (q *Queries) GetDefaultLanguage(ctx context.Context) language.Tag {
	instance, err := q.Instance(ctx, false)
	if err != nil {
		return language.Und
	}
	return instance.DefaultLanguage()
}

func prepareInstanceQuery(ctx context.Context, db prepareDatabase, host string) (sq.SelectBuilder, func(*sql.Row) (*Instance, error)) {
	return sq.Select(
			InstanceColumnID.identifier(),
			InstanceColumnCreationDate.identifier(),
			InstanceColumnChangeDate.identifier(),
			InstanceColumnSequence.identifier(),
			InstanceColumnDefaultOrgID.identifier(),
			InstanceColumnProjectID.identifier(),
			InstanceColumnConsoleID.identifier(),
			InstanceColumnConsoleAppID.identifier(),
			InstanceColumnDefaultLanguage.identifier(),
		).
			From(instanceTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Instance, error) {
			instance := &Instance{host: host}
			lang := ""
			err := row.Scan(
				&instance.ID,
				&instance.CreationDate,
				&instance.ChangeDate,
				&instance.Sequence,
				&instance.DefaultOrgID,
				&instance.IAMProjectID,
				&instance.ConsoleID,
				&instance.ConsoleAppID,
				&lang,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-5m09s", "Errors.IAM.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-3j9sf", "Errors.Internal")
			}
			instance.DefaultLang = language.Make(lang)
			return instance, nil
		}
}

func prepareInstancesQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(sq.SelectBuilder) sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
	instanceFilterTable := instanceTable.setAlias(InstancesFilterTableAlias)
	instanceFilterIDColumn := InstanceColumnID.setTable(instanceFilterTable)
	instanceFilterCountColumn := InstancesFilterTableAlias + ".count"
	return sq.Select(
			InstanceColumnID.identifier(),
			countColumn.identifier(),
		).From(instanceTable.identifier()),
		func(builder sq.SelectBuilder) sq.SelectBuilder {
			return sq.Select(
				instanceFilterCountColumn,
				instanceFilterIDColumn.identifier(),
				InstanceColumnCreationDate.identifier(),
				InstanceColumnChangeDate.identifier(),
				InstanceColumnSequence.identifier(),
				InstanceColumnName.identifier(),
				InstanceColumnDefaultOrgID.identifier(),
				InstanceColumnProjectID.identifier(),
				InstanceColumnConsoleID.identifier(),
				InstanceColumnConsoleAppID.identifier(),
				InstanceColumnDefaultLanguage.identifier(),
				InstanceDomainDomainCol.identifier(),
				InstanceDomainIsPrimaryCol.identifier(),
				InstanceDomainIsGeneratedCol.identifier(),
				InstanceDomainCreationDateCol.identifier(),
				InstanceDomainChangeDateCol.identifier(),
				InstanceDomainSequenceCol.identifier(),
			).FromSelect(builder, InstancesFilterTableAlias).
				LeftJoin(join(InstanceColumnID, instanceFilterIDColumn)).
				LeftJoin(join(InstanceDomainInstanceIDCol, instanceFilterIDColumn) + db.Timetravel(call.Took(ctx))).
				PlaceholderFormat(sq.Dollar)
		},
		func(rows *sql.Rows) (*Instances, error) {
			instances := make([]*Instance, 0)
			var lastInstance *Instance
			var count uint64
			for rows.Next() {
				instance := new(Instance)
				lang := ""
				var (
					domain       sql.NullString
					isPrimary    sql.NullBool
					isGenerated  sql.NullBool
					changeDate   sql.NullTime
					creationDate sql.NullTime
					sequence     sql.NullInt64
				)
				err := rows.Scan(
					&count,
					&instance.ID,
					&instance.CreationDate,
					&instance.ChangeDate,
					&instance.Sequence,
					&instance.Name,
					&instance.DefaultOrgID,
					&instance.IAMProjectID,
					&instance.ConsoleID,
					&instance.ConsoleAppID,
					&lang,
					&domain,
					&isPrimary,
					&isGenerated,
					&changeDate,
					&creationDate,
					&sequence,
				)
				if err != nil {
					return nil, err
				}
				if instance.ID == "" || !domain.Valid {
					continue
				}
				instance.DefaultLang = language.Make(lang)
				instanceDomain := &InstanceDomain{
					CreationDate: creationDate.Time,
					ChangeDate:   changeDate.Time,
					Sequence:     uint64(sequence.Int64),
					Domain:       domain.String,
					IsPrimary:    isPrimary.Bool,
					IsGenerated:  isGenerated.Bool,
					InstanceID:   instance.ID,
				}
				if lastInstance != nil && instance.ID == lastInstance.ID {
					lastInstance.Domains = append(lastInstance.Domains, instanceDomain)
					continue
				}
				lastInstance = instance
				instance.Domains = append(instance.Domains, instanceDomain)
				instances = append(instances, instance)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-8nlWW", "Errors.Query.CloseRows")
			}

			return &Instances{
				Instances: instances,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareInstanceDomainQuery(ctx context.Context, db prepareDatabase, host string) (sq.SelectBuilder, func(*sql.Rows) (*Instance, error)) {
	return sq.Select(
			InstanceColumnID.identifier(),
			InstanceColumnCreationDate.identifier(),
			InstanceColumnChangeDate.identifier(),
			InstanceColumnSequence.identifier(),
			InstanceColumnName.identifier(),
			InstanceColumnDefaultOrgID.identifier(),
			InstanceColumnProjectID.identifier(),
			InstanceColumnConsoleID.identifier(),
			InstanceColumnConsoleAppID.identifier(),
			InstanceColumnDefaultLanguage.identifier(),
			InstanceDomainDomainCol.identifier(),
			InstanceDomainIsPrimaryCol.identifier(),
			InstanceDomainIsGeneratedCol.identifier(),
			InstanceDomainCreationDateCol.identifier(),
			InstanceDomainChangeDateCol.identifier(),
			InstanceDomainSequenceCol.identifier(),
		).
			From(instanceTable.identifier()).
			LeftJoin(join(InstanceDomainInstanceIDCol, InstanceColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Instance, error) {
			instance := &Instance{
				host:    host,
				Domains: make([]*InstanceDomain, 0),
			}
			lang := ""
			for rows.Next() {
				var (
					domain       sql.NullString
					isPrimary    sql.NullBool
					isGenerated  sql.NullBool
					changeDate   sql.NullTime
					creationDate sql.NullTime
					sequence     sql.NullInt64
				)
				err := rows.Scan(
					&instance.ID,
					&instance.CreationDate,
					&instance.ChangeDate,
					&instance.Sequence,
					&instance.Name,
					&instance.DefaultOrgID,
					&instance.IAMProjectID,
					&instance.ConsoleID,
					&instance.ConsoleAppID,
					&lang,
					&domain,
					&isPrimary,
					&isGenerated,
					&changeDate,
					&creationDate,
					&sequence,
				)
				if err != nil {
					return nil, errors.ThrowInternal(err, "QUERY-d9nw", "Errors.Internal")
				}
				if !domain.Valid {
					continue
				}
				instance.Domains = append(instance.Domains, &InstanceDomain{
					CreationDate: creationDate.Time,
					ChangeDate:   changeDate.Time,
					Sequence:     uint64(sequence.Int64),
					Domain:       domain.String,
					IsPrimary:    isPrimary.Bool,
					IsGenerated:  isGenerated.Bool,
					InstanceID:   instance.ID,
				})
			}
			if instance.ID == "" {
				return nil, errors.ThrowNotFound(nil, "QUERY-n0wng", "Errors.IAM.NotFound")
			}
			instance.DefaultLang = language.Make(lang)
			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-Dfbe2", "Errors.Query.CloseRows")
			}
			return instance, nil
		}
}

func prepareAuthzInstanceQuery(ctx context.Context, db prepareDatabase, host string) (sq.SelectBuilder, func(*sql.Rows) (*Instance, error)) {
	return sq.Select(
			InstanceColumnID.identifier(),
			InstanceColumnCreationDate.identifier(),
			InstanceColumnChangeDate.identifier(),
			InstanceColumnSequence.identifier(),
			InstanceColumnName.identifier(),
			InstanceColumnDefaultOrgID.identifier(),
			InstanceColumnProjectID.identifier(),
			InstanceColumnConsoleID.identifier(),
			InstanceColumnConsoleAppID.identifier(),
			InstanceColumnDefaultLanguage.identifier(),
			InstanceDomainDomainCol.identifier(),
			InstanceDomainIsPrimaryCol.identifier(),
			InstanceDomainIsGeneratedCol.identifier(),
			InstanceDomainCreationDateCol.identifier(),
			InstanceDomainChangeDateCol.identifier(),
			InstanceDomainSequenceCol.identifier(),
			SecurityPolicyColumnEnabled.identifier(),
			SecurityPolicyColumnAllowedOrigins.identifier(),
		).
			From(instanceTable.identifier()).
			LeftJoin(join(InstanceDomainInstanceIDCol, InstanceColumnID)).
			LeftJoin(join(SecurityPolicyColumnInstanceID, InstanceColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Instance, error) {
			instance := &Instance{
				host:    host,
				Domains: make([]*InstanceDomain, 0),
			}
			lang := ""
			for rows.Next() {
				var (
					domain                sql.NullString
					isPrimary             sql.NullBool
					isGenerated           sql.NullBool
					changeDate            sql.NullTime
					creationDate          sql.NullTime
					sequence              sql.NullInt64
					securityPolicyEnabled sql.NullBool
				)
				err := rows.Scan(
					&instance.ID,
					&instance.CreationDate,
					&instance.ChangeDate,
					&instance.Sequence,
					&instance.Name,
					&instance.DefaultOrgID,
					&instance.IAMProjectID,
					&instance.ConsoleID,
					&instance.ConsoleAppID,
					&lang,
					&domain,
					&isPrimary,
					&isGenerated,
					&changeDate,
					&creationDate,
					&sequence,
					&securityPolicyEnabled,
					&instance.csp.allowedOrigins,
				)
				if err != nil {
					return nil, errors.ThrowInternal(err, "QUERY-d3fas", "Errors.Internal")
				}
				if !domain.Valid {
					continue
				}
				instance.Domains = append(instance.Domains, &InstanceDomain{
					CreationDate: creationDate.Time,
					ChangeDate:   changeDate.Time,
					Sequence:     uint64(sequence.Int64),
					Domain:       domain.String,
					IsPrimary:    isPrimary.Bool,
					IsGenerated:  isGenerated.Bool,
					InstanceID:   instance.ID,
				})
				instance.csp.enabled = securityPolicyEnabled.Bool
			}
			if instance.ID == "" {
				return nil, errors.ThrowNotFound(nil, "QUERY-n0wng", "Errors.IAM.NotFound")
			}
			instance.DefaultLang = language.Make(lang)
			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-Dfbe2", "Errors.Query.CloseRows")
			}
			return instance, nil
		}
}
