package query

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	InstancesFilterTableAlias = "f"
)

var (
	instanceTable = table{
		name:          projection.InstanceProjectionTable,
		instanceIDCol: projection.InstanceColumnID,
	}
	limitsTable = table{
		name:          projection.LimitsProjectionTable,
		instanceIDCol: projection.LimitsColumnInstanceID,
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
	LimitsColumnInstanceID = Column{
		name:  projection.LimitsColumnInstanceID,
		table: limitsTable,
	}
	LimitsColumnAuditLogRetention = Column{
		name:  projection.LimitsColumnAuditLogRetention,
		table: limitsTable,
	}
	LimitsColumnBlock = Column{
		name:  projection.LimitsColumnBlock,
		table: limitsTable,
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
}

type Instances struct {
	SearchResponse
	Instances []*Instance
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

func NewInstanceDomainsListSearchQuery(domains ...string) (SearchQuery, error) {
	list := make([]interface{}, len(domains))
	for i, value := range domains {
		list[i] = value
	}

	return NewListQuery(InstanceDomainDomainCol, list, ListIn)
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
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-M9fow", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		instances, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-3j98f", "Errors.Internal")
	}
	return instances, err
}

func (q *Queries) Instance(ctx context.Context, shouldTriggerBulk bool) (instance *Instance, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerInstanceProjection")
		ctx, err = projection.InstanceProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareInstanceDomainQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		InstanceColumnID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-d9ngs", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		instance, err = scan(rows)
		return err
	}, query, args...)
	return instance, err
}

var (
	//go:embed instance_by_domain.sql
	instanceByDomainQuery string

	//go:embed instance_by_id.sql
	instanceByIDQuery string
)

func (q *Queries) InstanceByHost(ctx context.Context, host string) (_ authz.Instance, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		if err != nil {
			err = fmt.Errorf("unable to get instance by host %s: %w", host, err)
		}
		span.EndWithError(err)
	}()

	domain := strings.Split(host, ":")[0] // remove possible port
	instance, scan := scanAuthzInstance(host, domain)
	err = q.client.QueryRowContext(ctx, scan, instanceByDomainQuery, domain)
	logging.OnError(err).WithField("host", host).WithField("domain", domain).Warn("instance by host")
	return instance, err
}

func (q *Queries) InstanceByID(ctx context.Context) (_ authz.Instance, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	instance, scan := scanAuthzInstance("", "")
	err = q.client.QueryRowContext(ctx, scan, instanceByIDQuery, instanceID)
	logging.OnError(err).WithField("instance_id", instanceID).Warn("instance by ID")
	return instance, err
}

func (q *Queries) GetDefaultLanguage(ctx context.Context) language.Tag {
	instance, err := q.Instance(ctx, false)
	if err != nil {
		return language.Und
	}
	return instance.DefaultLang
}

func prepareInstancesQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(sq.SelectBuilder) sq.SelectBuilder, func(*sql.Rows) (*Instances, error)) {
	instanceFilterTable := instanceTable.setAlias(InstancesFilterTableAlias)
	instanceFilterIDColumn := InstanceColumnID.setTable(instanceFilterTable)
	instanceFilterCountColumn := InstancesFilterTableAlias + ".count"
	return sq.Select(
			InstanceColumnID.identifier(),
			countColumn.identifier(),
		).Distinct().From(instanceTable.identifier()).
			LeftJoin(join(InstanceDomainInstanceIDCol, InstanceColumnID)),
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
				return nil, zerrors.ThrowInternal(err, "QUERY-8nlWW", "Errors.Query.CloseRows")
			}

			return &Instances{
				Instances: instances,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareInstanceDomainQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Instance, error)) {
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
					return nil, zerrors.ThrowInternal(err, "QUERY-d9nw", "Errors.Internal")
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
				return nil, zerrors.ThrowNotFound(nil, "QUERY-n0wng", "Errors.IAM.NotFound")
			}
			instance.DefaultLang = language.Make(lang)
			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-Dfbe2", "Errors.Query.CloseRows")
			}
			return instance, nil
		}
}

type authzInstance struct {
	id                  string
	iamProjectID        string
	consoleID           string
	consoleAppID        string
	host                string
	domain              string
	defaultLang         language.Tag
	defaultOrgID        string
	csp                 csp
	enableImpersonation bool
	block               *bool
	auditLogRetention   *time.Duration
	features            feature.Features
}

type csp struct {
	enableIframeEmbedding bool
	allowedOrigins        database.TextArray[string]
}

func (i *authzInstance) InstanceID() string {
	return i.id
}

func (i *authzInstance) ProjectID() string {
	return i.iamProjectID
}

func (i *authzInstance) ConsoleClientID() string {
	return i.consoleID
}

func (i *authzInstance) ConsoleApplicationID() string {
	return i.consoleAppID
}

func (i *authzInstance) RequestedDomain() string {
	return strings.Split(i.host, ":")[0]
}

func (i *authzInstance) RequestedHost() string {
	return i.host
}

func (i *authzInstance) DefaultLanguage() language.Tag {
	return i.defaultLang
}

func (i *authzInstance) DefaultOrganisationID() string {
	return i.defaultOrgID
}

func (i *authzInstance) SecurityPolicyAllowedOrigins() []string {
	if !i.csp.enableIframeEmbedding {
		return nil
	}
	return i.csp.allowedOrigins
}

func (i *authzInstance) EnableImpersonation() bool {
	return i.enableImpersonation
}

func (i *authzInstance) Block() *bool {
	return i.block
}

func (i *authzInstance) AuditLogRetention() *time.Duration {
	return i.auditLogRetention
}

func (i *authzInstance) Features() feature.Features {
	return i.features
}

func scanAuthzInstance(host, domain string) (*authzInstance, func(row *sql.Row) error) {
	instance := &authzInstance{
		host:   host,
		domain: domain,
	}
	return instance, func(row *sql.Row) error {
		var (
			lang                  string
			enableIframeEmbedding sql.NullBool
			enableImpersonation   sql.NullBool
			auditLogRetention     database.NullDuration
			block                 sql.NullBool
			features              []byte
		)
		err := row.Scan(
			&instance.id,
			&instance.defaultOrgID,
			&instance.iamProjectID,
			&instance.consoleID,
			&instance.consoleAppID,
			&lang,
			&enableIframeEmbedding,
			&instance.csp.allowedOrigins,
			&enableImpersonation,
			&auditLogRetention,
			&block,
			&features,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return zerrors.ThrowNotFound(nil, "QUERY-1kIjX", "Errors.IAM.NotFound")
		}
		if err != nil {
			return zerrors.ThrowInternal(err, "QUERY-d3fas", "Errors.Internal")
		}
		instance.defaultLang = language.Make(lang)
		if auditLogRetention.Valid {
			instance.auditLogRetention = &auditLogRetention.Duration
		}
		if block.Valid {
			instance.block = &block.Bool
		}
		instance.csp.enableIframeEmbedding = enableIframeEmbedding.Bool
		instance.enableImpersonation = enableImpersonation.Bool
		if len(features) == 0 {
			return nil
		}
		if err = json.Unmarshal(features, &instance.features); err != nil {
			return zerrors.ThrowInternal(err, "QUERY-Po8ki", "Errors.Internal")
		}
		return nil
	}
}
