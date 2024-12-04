package query

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
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

func (q *Queries) ActiveInstances() []string {
	return q.caches.activeInstances.Keys()
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

func (q *Queries) InstanceByHost(ctx context.Context, instanceHost, publicHost string) (_ authz.Instance, err error) {
	var instance *authzInstance
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		if err != nil {
			err = fmt.Errorf("unable to get instance by host: instanceHost %s, publicHost %s: %w", instanceHost, publicHost, err)
		} else {
			q.caches.activeInstances.Add(instance.ID, true)
		}
		span.EndWithError(err)
	}()

	instanceDomain := strings.Split(instanceHost, ":")[0] // remove possible port
	publicDomain := strings.Split(publicHost, ":")[0]     // remove possible port

	instance, ok := q.caches.instance.Get(ctx, instanceIndexByHost, instanceDomain)
	if ok {
		return instance, instance.checkDomain(instanceDomain, publicDomain)
	}
	instance, scan := scanAuthzInstance()
	if err = q.client.QueryRowContext(ctx, scan, instanceByDomainQuery, instanceDomain); err != nil {
		return nil, err
	}
	q.caches.instance.Set(ctx, instance)

	return instance, instance.checkDomain(instanceDomain, publicDomain)
}

func (q *Queries) InstanceByID(ctx context.Context, id string) (_ authz.Instance, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	defer func() {
		if err != nil {
			return
		}
		q.caches.activeInstances.Add(id, true)
	}()

	instance, ok := q.caches.instance.Get(ctx, instanceIndexByID, id)
	if ok {
		return instance, nil
	}

	instance, scan := scanAuthzInstance()
	err = q.client.QueryRowContext(ctx, scan, instanceByIDQuery, id)
	logging.OnError(err).WithField("instance_id", id).Warn("instance by ID")
	if err == nil {
		q.caches.instance.Set(ctx, instance)
	}
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
	ID              string                     `json:"id,omitempty"`
	IAMProjectID    string                     `json:"iam_project_id,omitempty"`
	ConsoleID       string                     `json:"console_id,omitempty"`
	ConsoleAppID    string                     `json:"console_app_id,omitempty"`
	DefaultLang     language.Tag               `json:"default_lang,omitempty"`
	DefaultOrgID    string                     `json:"default_org_id,omitempty"`
	CSP             csp                        `json:"csp,omitempty"`
	Impersonation   bool                       `json:"impersonation,omitempty"`
	IsBlocked       *bool                      `json:"is_blocked,omitempty"`
	LogRetention    *time.Duration             `json:"log_retention,omitempty"`
	Feature         feature.Features           `json:"feature,omitempty"`
	ExternalDomains database.TextArray[string] `json:"external_domains,omitempty"`
	TrustedDomains  database.TextArray[string] `json:"trusted_domains,omitempty"`
}

type csp struct {
	EnableIframeEmbedding bool                       `json:"enable_iframe_embedding,omitempty"`
	AllowedOrigins        database.TextArray[string] `json:"allowed_origins,omitempty"`
}

func (i *authzInstance) InstanceID() string {
	return i.ID
}

func (i *authzInstance) ProjectID() string {
	return i.IAMProjectID
}

func (i *authzInstance) ConsoleClientID() string {
	return i.ConsoleID
}

func (i *authzInstance) ConsoleApplicationID() string {
	return i.ConsoleAppID
}

func (i *authzInstance) DefaultLanguage() language.Tag {
	return i.DefaultLang
}

func (i *authzInstance) DefaultOrganisationID() string {
	return i.DefaultOrgID
}

func (i *authzInstance) SecurityPolicyAllowedOrigins() []string {
	if !i.CSP.EnableIframeEmbedding {
		return nil
	}
	return i.CSP.AllowedOrigins
}

func (i *authzInstance) EnableImpersonation() bool {
	return i.Impersonation
}

func (i *authzInstance) Block() *bool {
	return i.IsBlocked
}

func (i *authzInstance) AuditLogRetention() *time.Duration {
	return i.LogRetention
}

func (i *authzInstance) Features() feature.Features {
	return i.Feature
}

var errPublicDomain = "public domain %q not trusted"

func (i *authzInstance) checkDomain(instanceDomain, publicDomain string) error {
	// in case public domain is empty, or the same as the instance domain, we do not need to check it
	if publicDomain == "" || instanceDomain == publicDomain {
		return nil
	}
	if !slices.Contains(i.TrustedDomains, publicDomain) {
		return zerrors.ThrowNotFound(fmt.Errorf(errPublicDomain, publicDomain), "QUERY-IuGh1", "Errors.IAM.NotFound")
	}
	return nil
}

// Keys implements [cache.Entry]
func (i *authzInstance) Keys(index instanceIndex) []string {
	switch index {
	case instanceIndexByID:
		return []string{i.ID}
	case instanceIndexByHost:
		return i.ExternalDomains
	case instanceIndexUnspecified:
	}
	return nil
}

func scanAuthzInstance() (*authzInstance, func(row *sql.Row) error) {
	instance := &authzInstance{}
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
			&instance.ID,
			&instance.DefaultOrgID,
			&instance.IAMProjectID,
			&instance.ConsoleID,
			&instance.ConsoleAppID,
			&lang,
			&enableIframeEmbedding,
			&instance.CSP.AllowedOrigins,
			&enableImpersonation,
			&auditLogRetention,
			&block,
			&features,
			&instance.ExternalDomains,
			&instance.TrustedDomains,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return zerrors.ThrowNotFound(nil, "QUERY-1kIjX", "Errors.IAM.NotFound")
		}
		if err != nil {
			return zerrors.ThrowInternal(err, "QUERY-d3fas", "Errors.Internal")
		}
		instance.DefaultLang = language.Make(lang)
		if auditLogRetention.Valid {
			instance.LogRetention = &auditLogRetention.Duration
		}
		if block.Valid {
			instance.IsBlocked = &block.Bool
		}
		instance.CSP.EnableIframeEmbedding = enableIframeEmbedding.Bool
		instance.Impersonation = enableImpersonation.Bool
		if len(features) == 0 {
			return nil
		}
		if err = json.Unmarshal(features, &instance.Feature); err != nil {
			return zerrors.ThrowInternal(err, "QUERY-Po8ki", "Errors.Internal")
		}
		return nil
	}
}

func (c *Caches) registerInstanceInvalidation() {
	invalidate := cacheInvalidationFunc(c.instance, instanceIndexByID, getAggregateID)
	projection.InstanceProjection.RegisterCacheInvalidation(invalidate)
	projection.InstanceDomainProjection.RegisterCacheInvalidation(invalidate)
	projection.InstanceFeatureProjection.RegisterCacheInvalidation(invalidate)
	projection.InstanceTrustedDomainProjection.RegisterCacheInvalidation(invalidate)
	projection.SecurityPolicyProjection.RegisterCacheInvalidation(invalidate)

	// These projections have their own aggregate ID, invalidate using resource owner.
	invalidate = cacheInvalidationFunc(c.instance, instanceIndexByID, getResourceOwner)
	projection.LimitsProjection.RegisterCacheInvalidation(invalidate)
	projection.RestrictionsProjection.RegisterCacheInvalidation(invalidate)

	// System feature update should invalidate all instances, so Truncate the cache.
	projection.SystemFeatureProjection.RegisterCacheInvalidation(func(ctx context.Context, _ []*eventstore.Aggregate) {
		err := c.instance.Truncate(ctx)
		logging.OnError(err).Warn("cache truncate failed")
	})
}

type instanceIndex int

//go:generate enumer -type instanceIndex -linecomment
const (
	// Empty line comment ensures empty string for unspecified value
	instanceIndexUnspecified instanceIndex = iota //
	instanceIndexByID
	instanceIndexByHost
)
