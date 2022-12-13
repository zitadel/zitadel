package query

import (
	"context"
	"database/sql"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/projection"
	projection_old "github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	InstancesFilterTableAlias = "f"
)

var (
	instanceTable = table{
		name:          projection_old.InstanceProjectionTable,
		instanceIDCol: projection_old.InstanceColumnID,
	}
	InstanceColumnID = Column{
		name:  projection_old.InstanceColumnID,
		table: instanceTable,
	}
	InstanceColumnName = Column{
		name:  projection_old.InstanceColumnName,
		table: instanceTable,
	}
	InstanceColumnCreationDate = Column{
		name:  projection_old.InstanceColumnCreationDate,
		table: instanceTable,
	}
	InstanceColumnChangeDate = Column{
		name:  projection_old.InstanceColumnChangeDate,
		table: instanceTable,
	}
	InstanceColumnSequence = Column{
		name:  projection_old.InstanceColumnSequence,
		table: instanceTable,
	}
	InstanceColumnDefaultOrgID = Column{
		name:  projection_old.InstanceColumnDefaultOrgID,
		table: instanceTable,
	}
	InstanceColumnProjectID = Column{
		name:  projection_old.InstanceColumnProjectID,
		table: instanceTable,
	}
	InstanceColumnConsoleID = Column{
		name:  projection_old.InstanceColumnConsoleID,
		table: instanceTable,
	}
	InstanceColumnConsoleAppID = Column{
		name:  projection_old.InstanceColumnConsoleAppID,
		table: instanceTable,
	}
	InstanceColumnDefaultLanguage = Column{
		name:  projection_old.InstanceColumnDefaultLanguage,
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

	filter, query, scan := prepareInstancesQuery()
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

	instance := projection.NewInstance(authz.GetInstance(ctx).InstanceID(), authz.GetInstance(ctx).RequestedDomain())
	events, err := q.eventstore.Filter(ctx, instance.SearchQuery(ctx))
	if err != nil {
		return nil, err
	}
	instance.Reduce(events)

	return mapInstance(instance), nil
}

func (q *Queries) InstanceByHost(ctx context.Context, host string) (_ authz.Instance, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// remove possible port
	domain := strings.Split(host, ":")[0]

	domainSearch := projection.NewSearchInstanceDomain(domain)
	events, err := q.eventstore.Filter(ctx, domainSearch.SearchQuery(ctx))
	if err != nil {
		return nil, err
	}
	domainSearch.Reduce(events)
	if domainSearch.InstanceID == "" {
		return nil, errors.ThrowNotFound(nil, "QUERY-VZHH2", "Errors.NotFound")
	}

	instance := projection.NewInstance(domainSearch.InstanceID, host)
	events, err = q.eventstore.Filter(ctx, instance.SearchQuery(ctx))
	if err != nil {
		return nil, err
	}
	instance.Reduce(events)

	return mapInstance(instance), nil
}

func mapInstance(instance *projection.Instance) *Instance {
	return &Instance{
		ID:           instance.ID,
		ChangeDate:   instance.ChangeDate,
		CreationDate: instance.CreationDate,
		Sequence:     instance.Sequence,
		Name:         instance.Name,
		host:         instance.Host,

		DefaultOrgID: instance.DefaultOrgID,
		IAMProjectID: instance.IAMProjectID,
		ConsoleID:    instance.ConsoleID,
		ConsoleAppID: instance.ConsoleAppID,
		DefaultLang:  instance.DefaultLang,
		Domains:      mapInstanceDomains(instance.Domains),
	}
}

func mapInstanceDomains(domains []*projection.InstanceDomain) []*InstanceDomain {
	instanceDomains := make([]*InstanceDomain, len(domains))

	for i, domain := range domains {
		instanceDomains[i] = &InstanceDomain{
			CreationDate: domain.CreationDate,
			ChangeDate:   domain.ChangeDate,
			Sequence:     domain.Sequence,
			Domain:       domain.Domain,
			InstanceID:   domain.InstanceID,
			IsGenerated:  domain.IsGenerated,
			IsPrimary:    domain.IsPrimary,
		}
	}

	return instanceDomains
}

func (q *Queries) GetDefaultLanguage(ctx context.Context) language.Tag {
	instance, err := q.Instance(ctx, false)
	if err != nil {
		return language.Und
	}
	return instance.DefaultLanguage()
}

func prepareInstancesQuery() (sq.SelectBuilder, func(sq.SelectBuilder) sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
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
				LeftJoin(join(InstanceDomainInstanceIDCol, instanceFilterIDColumn)).
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
