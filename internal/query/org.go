package query

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	domain_pkg "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/readmodel"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	orgsTable = table{
		name:          projection.OrgProjectionTable,
		instanceIDCol: projection.OrgColumnInstanceID,
	}
	OrgColumnID = Column{
		name:  projection.OrgColumnID,
		table: orgsTable,
	}
	OrgColumnCreationDate = Column{
		name:  projection.OrgColumnCreationDate,
		table: orgsTable,
	}
	OrgColumnChangeDate = Column{
		name:  projection.OrgColumnChangeDate,
		table: orgsTable,
	}
	OrgColumnResourceOwner = Column{
		name:  projection.OrgColumnResourceOwner,
		table: orgsTable,
	}
	OrgColumnInstanceID = Column{
		name:  projection.OrgColumnInstanceID,
		table: orgsTable,
	}
	OrgColumnState = Column{
		name:  projection.OrgColumnState,
		table: orgsTable,
	}
	OrgColumnSequence = Column{
		name:  projection.OrgColumnSequence,
		table: orgsTable,
	}
	OrgColumnName = Column{
		name:           projection.OrgColumnName,
		table:          orgsTable,
		isOrderByLower: true,
	}
	OrgColumnDomain = Column{
		name:  projection.OrgColumnDomain,
		table: orgsTable,
	}
)

type Orgs struct {
	SearchResponse
	Orgs []*Org
}

type Org struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain_pkg.OrgState
	Sequence      uint64

	Name   string
	Domain string
}

func orgsCheckPermission(ctx context.Context, orgs *Orgs, permissionCheck domain_pkg.PermissionCheck) {
	orgs.Orgs = slices.DeleteFunc(orgs.Orgs,
		func(org *Org) bool {
			if err := permissionCheck(ctx, domain_pkg.PermissionOrgRead, org.ID, org.ID); err != nil {
				return true
			}
			return false
		},
	)
}

type OrgSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *OrgSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) OrgByID(ctx context.Context, shouldTriggerBulk bool, id string) (org *Org, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org, ok := q.caches.org.Get(ctx, orgIndexByID, id); ok {
		return org, nil
	}
	defer func() {
		if err == nil && org != nil {
			q.caches.org.Set(ctx, org)
		}
	}()

	if !authz.GetInstance(ctx).Features().ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeOrgByID) {
		return q.oldOrgByID(ctx, shouldTriggerBulk, id)
	}

	foundOrg := readmodel.NewOrg(id)
	eventCount, err := q.eventStoreV4.Query(
		ctx,
		eventstore.NewQuery(
			authz.GetInstance(ctx).InstanceID(),
			foundOrg,
			eventstore.AppendFilters(foundOrg.Filter()...),
		),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-AWx52", "Errors.Query.SQLStatement")
	}

	if eventCount == 0 {
		return nil, zerrors.ThrowNotFound(nil, "QUERY-leq5Q", "Errors.Org.NotFound")
	}

	return &Org{
		ID:            foundOrg.ID,
		CreationDate:  foundOrg.CreationDate,
		ChangeDate:    foundOrg.ChangeDate,
		ResourceOwner: foundOrg.Owner,
		State:         domain_pkg.OrgState(foundOrg.State.State),
		Sequence:      uint64(foundOrg.Sequence),
		Name:          foundOrg.Name,
		Domain:        foundOrg.PrimaryDomain.Domain,
	}, nil
}

func (q *Queries) oldOrgByID(ctx context.Context, shouldTriggerBulk bool, id string) (org *Org, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerOrgProjection")
		ctx, err = projection.OrgProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareOrgQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		OrgColumnID.identifier():         id,
		OrgColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-AWx52", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		org, err = scan(row)
		return err
	}, query, args...)
	return org, err
}

func (q *Queries) OrgByPrimaryDomain(ctx context.Context, domain string) (org *Org, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	org, ok := q.caches.org.Get(ctx, orgIndexByPrimaryDomain, domain)
	if ok {
		return org, nil
	}

	stmt, scan := prepareOrgQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		OrgColumnDomain.identifier():     domain,
		OrgColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		OrgColumnState.identifier():      domain_pkg.OrgStateActive,
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-TYUCE", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		org, err = scan(row)
		return err
	}, query, args...)
	if err == nil {
		q.caches.org.Set(ctx, org)
	}
	return org, err
}

func (q *Queries) OrgByVerifiedDomain(ctx context.Context, domain string) (org *Org, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareOrgWithDomainsQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		OrgDomainDomainCol.identifier():     domain,
		OrgDomainIsVerifiedCol.identifier(): true,
		OrgColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-TYUCE", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		org, err = scan(row)
		return err
	}, query, args...)
	return org, err
}

func (q *Queries) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if name == "" && domain == "" {
		return false, zerrors.ThrowInvalidArgument(nil, "QUERY-DGqfd", "Errors.Query.InvalidRequest")
	}
	query, scan := prepareOrgUniqueQuery(ctx, q.client)
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{
				OrgColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
				OrgDomainIsVerifiedCol.identifier(): true,
			},
			sq.Or{
				sq.ILike{
					OrgDomainDomainCol.identifier(): domain,
				},
				sq.ILike{
					OrgColumnName.identifier(): name,
				},
			},
			sq.NotEq{
				OrgColumnState.identifier(): domain_pkg.OrgStateRemoved,
			},
		}).ToSql()
	if err != nil {
		return false, zerrors.ThrowInternal(err, "QUERY-Dgbe2", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		isUnique, err = scan(row)
		return err
	}, stmt, args...)
	return isUnique, err
}

func (q *Queries) ExistsOrg(ctx context.Context, id, domain string) (verifiedID string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var org *Org
	if id != "" {
		org, err = q.OrgByID(ctx, true, id)
	} else {
		org, err = q.OrgByVerifiedDomain(ctx, domain)
	}
	if err != nil {
		return "", err
	}
	return org.ID, nil
}

func (q *Queries) SearchOrgs(ctx context.Context, queries *OrgSearchQueries, permissionCheck domain_pkg.PermissionCheck) (*Orgs, error) {
	orgs, err := q.searchOrgs(ctx, queries)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil {
		orgsCheckPermission(ctx, orgs, permissionCheck)
	}
	return orgs, nil
}

func (q *Queries) searchOrgs(ctx context.Context, queries *OrgSearchQueries) (orgs *Orgs, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareOrgsQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).
		Where(sq.And{
			sq.Eq{
				OrgColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
			},
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-wQ3by", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		orgs, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-M6mYN", "Errors.Internal")
	}

	orgs.State, err = q.latestState(ctx, orgsTable)
	return orgs, err
}

func NewOrgIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(OrgColumnID, value, TextEquals)
}

func NewOrgVerifiedDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	domainQuery, err := NewTextQuery(OrgDomainDomainCol, value, method)
	if err != nil {
		return nil, err
	}
	verifiedQuery, err := NewBoolQuery(OrgDomainIsVerifiedCol, true)
	if err != nil {
		return nil, err
	}
	subSelect, err := NewSubSelect(OrgDomainOrgIDCol, []SearchQuery{domainQuery, verifiedQuery})
	if err != nil {
		return nil, err
	}
	return NewListQuery(
		OrgColumnID,
		subSelect,
		ListIn,
	)
}

func NewOrgNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(OrgColumnName, value, method)
}

func NewOrgStateSearchQuery(value domain_pkg.OrgState) (SearchQuery, error) {
	return NewNumberQuery(OrgColumnState, value, NumberEquals)
}

func NewOrgIDsSearchQuery(ids ...string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(OrgColumnID, list, ListIn)
}

func prepareOrgsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Orgs, error)) {
	return sq.Select(
			OrgColumnID.identifier(),
			OrgColumnCreationDate.identifier(),
			OrgColumnChangeDate.identifier(),
			OrgColumnResourceOwner.identifier(),
			OrgColumnState.identifier(),
			OrgColumnSequence.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),
			countColumn.identifier()).
			From(orgsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Orgs, error) {
			orgs := make([]*Org, 0)
			var count uint64
			for rows.Next() {
				org := new(Org)
				err := rows.Scan(
					&org.ID,
					&org.CreationDate,
					&org.ChangeDate,
					&org.ResourceOwner,
					&org.State,
					&org.Sequence,
					&org.Name,
					&org.Domain,
					&count,
				)
				if err != nil {
					return nil, err
				}
				orgs = append(orgs, org)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-QMXJv", "Errors.Query.CloseRows")
			}

			return &Orgs{
				Orgs: orgs,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareOrgQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Org, error)) {
	return sq.Select(
			OrgColumnID.identifier(),
			OrgColumnCreationDate.identifier(),
			OrgColumnChangeDate.identifier(),
			OrgColumnResourceOwner.identifier(),
			OrgColumnState.identifier(),
			OrgColumnSequence.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),
		).
			From(orgsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Org, error) {
			o := new(Org)
			err := row.Scan(
				&o.ID,
				&o.CreationDate,
				&o.ChangeDate,
				&o.ResourceOwner,
				&o.State,
				&o.Sequence,
				&o.Name,
				&o.Domain,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-iTTGJ", "Errors.Org.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-pWS5H", "Errors.Internal")
			}
			return o, nil
		}
}

func prepareOrgWithDomainsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Org, error)) {
	return sq.Select(
			OrgColumnID.identifier(),
			OrgColumnCreationDate.identifier(),
			OrgColumnChangeDate.identifier(),
			OrgColumnResourceOwner.identifier(),
			OrgColumnState.identifier(),
			OrgColumnSequence.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),
		).
			From(orgsTable.identifier()).
			LeftJoin(join(OrgDomainOrgIDCol, OrgColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Org, error) {
			o := new(Org)
			err := row.Scan(
				&o.ID,
				&o.CreationDate,
				&o.ChangeDate,
				&o.ResourceOwner,
				&o.State,
				&o.Sequence,
				&o.Name,
				&o.Domain,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-iTTGJ", "Errors.Org.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-pWS5H", "Errors.Internal")
			}
			return o, nil
		}
}

func prepareOrgUniqueQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (bool, error)) {
	return sq.Select(uniqueColumn.identifier()).
			From(orgsTable.identifier()).
			LeftJoin(join(OrgDomainOrgIDCol, OrgColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (isUnique bool, err error) {
			err = row.Scan(&isUnique)
			if err != nil {
				return false, zerrors.ThrowInternal(err, "QUERY-e6EiG", "Errors.Internal")
			}
			return isUnique, err
		}
}

type orgIndex int

//go:generate enumer -type orgIndex -linecomment
const (
	// Empty line comment ensures empty string for unspecified value
	orgIndexUnspecified orgIndex = iota //
	orgIndexByID
	orgIndexByPrimaryDomain
)

// Keys implements [cache.Entry]
func (o *Org) Keys(index orgIndex) []string {
	switch index {
	case orgIndexByID:
		return []string{o.ID}
	case orgIndexByPrimaryDomain:
		return []string{o.Domain}
	case orgIndexUnspecified:
	}
	return nil
}

func (c *Caches) registerOrgInvalidation() {
	invalidate := cacheInvalidationFunc(c.org, orgIndexByID, getAggregateID)
	projection.OrgProjection.RegisterCacheInvalidation(invalidate)
}
