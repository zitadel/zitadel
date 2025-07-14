package readmodel

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/domain"
	v2domain "github.com/zitadel/zitadel/internal/v2/domain"
	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Database interfaces for the repository
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type TX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Rollback() error
	Commit() error
}

const (
	domainsTable             = "zitadel.domains"
	domainsIDCol             = "id"
	domainsInstanceIDCol     = "instance_id"
	domainsOrgIDCol          = "org_id"
	domainsDomainCol         = "domain"
	domainsIsVerifiedCol     = "is_verified"
	domainsIsPrimaryCol      = "is_primary"
	domainsValidationTypeCol = "validation_type"
	domainsCreatedAtCol      = "created_at"
	domainsUpdatedAtCol      = "updated_at"
	domainsDeletedAtCol      = "deleted_at"
)

// DomainRepository implements both InstanceDomainRepository and OrganizationDomainRepository
type DomainRepository struct {
	client DB
}

// NewDomainRepository creates a new DomainRepository
func NewDomainRepository(client DB) *DomainRepository {
	return &DomainRepository{
		client: client,
	}
}

// Add creates a new instance domain (always verified)
func (r *DomainRepository) AddInstanceDomain(ctx context.Context, instanceID, domain string) (*v2domain.Domain, error) {
	now := time.Now()
	query := squirrel.Insert(domainsTable).
		Columns(
			domainsInstanceIDCol,
			domainsDomainCol,
			domainsIsVerifiedCol,
			domainsIsPrimaryCol,
			domainsCreatedAtCol,
			domainsUpdatedAtCol,
		).
		Values(instanceID, domain, true, false, now, now).
		Suffix("RETURNING " + domainsIDCol).
		PlaceholderFormat(squirrel.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-1n8fK", "Errors.Internal")
	}

	var id string
	err = r.client.QueryRowContext(ctx, stmt, args...).Scan(&id)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-3m9sL", "Errors.Internal")
	}

	return &v2domain.Domain{
		ID:             id,
		InstanceID:     instanceID,
		OrganizationID: nil,
		Domain:         domain,
		IsVerified:     true,
		IsPrimary:      false,
		ValidationType: nil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// SetPrimary sets the primary domain for an instance
func (r *DomainRepository) SetInstanceDomainPrimary(ctx context.Context, instanceID, domain string) error {
	return r.withTransaction(ctx, func(tx TX) error {
		// First, unset any existing primary domain for this instance
		unsetQuery := squirrel.Update(domainsTable).
			Set(domainsIsPrimaryCol, false).
			Set(domainsUpdatedAtCol, time.Now()).
			Where(squirrel.Eq{
				domainsInstanceIDCol: instanceID,
				domainsOrgIDCol:      nil,
				domainsIsPrimaryCol:  true,
			}).
			Where(squirrel.Expr(domainsDeletedAtCol + " IS NULL")).
			PlaceholderFormat(squirrel.Dollar)

		stmt, args, err := unsetQuery.ToSql()
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-4k2nL", "Errors.Internal")
		}

		_, err = tx.ExecContext(ctx, stmt, args...)
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-5n3mK", "Errors.Internal")
		}

		// Then set the new primary domain
		setPrimaryQuery := squirrel.Update(domainsTable).
			Set(domainsIsPrimaryCol, true).
			Set(domainsUpdatedAtCol, time.Now()).
			Where(squirrel.Eq{
				domainsInstanceIDCol: instanceID,
				domainsOrgIDCol:      nil,
				domainsDomainCol:     domain,
			}).
			Where(squirrel.Expr(domainsDeletedAtCol + " IS NULL")).
			PlaceholderFormat(squirrel.Dollar)

		stmt, args, err = setPrimaryQuery.ToSql()
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-6o4pM", "Errors.Internal")
		}

		result, err := tx.ExecContext(ctx, stmt, args...)
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-7p5qN", "Errors.Internal")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-8q6rO", "Errors.Internal")
		}

		if rowsAffected == 0 {
			return zerrors.ThrowNotFound(nil, "DOMAIN-9r7sP", "Errors.Domain.NotFound")
		}

		return nil
	})
}

// Remove soft deletes an instance domain
func (r *DomainRepository) RemoveInstanceDomain(ctx context.Context, instanceID, domain string) error {
	query := squirrel.Update(domainsTable).
		Set(domainsDeletedAtCol, time.Now()).
		Set(domainsUpdatedAtCol, time.Now()).
		Where(squirrel.Eq{
			domainsInstanceIDCol: instanceID,
			domainsOrgIDCol:      nil,
			domainsDomainCol:     domain,
		}).
		Where(squirrel.Expr(domainsDeletedAtCol + " IS NULL")).
		PlaceholderFormat(squirrel.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-As8tQ", "Errors.Internal")
	}

	result, err := r.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-Bt9uR", "Errors.Internal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-Cu0vS", "Errors.Internal")
	}

	if rowsAffected == 0 {
		return zerrors.ThrowNotFound(nil, "DOMAIN-Dv1wT", "Errors.Domain.NotFound")
	}

	return nil
}

// Add creates a new organization domain
func (r *DomainRepository) AddOrganizationDomain(ctx context.Context, instanceID, organizationID, domain string, validationType domain.OrgDomainValidationType) (*v2domain.Domain, error) {
	now := time.Now()
	query := squirrel.Insert(domainsTable).
		Columns(
			domainsInstanceIDCol,
			domainsOrgIDCol,
			domainsDomainCol,
			domainsIsVerifiedCol,
			domainsIsPrimaryCol,
			domainsValidationTypeCol,
			domainsCreatedAtCol,
			domainsUpdatedAtCol,
		).
		Values(instanceID, organizationID, domain, false, false, int(validationType), now, now).
		Suffix("RETURNING " + domainsIDCol).
		PlaceholderFormat(squirrel.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-Ew2xU", "Errors.Internal")
	}

	var id string
	err = r.client.QueryRowContext(ctx, stmt, args...).Scan(&id)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-Fx3yV", "Errors.Internal")
	}

	return &v2domain.Domain{
		ID:             id,
		InstanceID:     instanceID,
		OrganizationID: &organizationID,
		Domain:         domain,
		IsVerified:     false,
		IsPrimary:      false,
		ValidationType: &validationType,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// SetVerified marks an organization domain as verified
func (r *DomainRepository) SetOrganizationDomainVerified(ctx context.Context, instanceID, organizationID, domain string) error {
	query := squirrel.Update(domainsTable).
		Set(domainsIsVerifiedCol, true).
		Set(domainsUpdatedAtCol, time.Now()).
		Where(squirrel.Eq{
			domainsInstanceIDCol: instanceID,
			domainsOrgIDCol:      organizationID,
			domainsDomainCol:     domain,
		}).
		Where(squirrel.Expr(domainsDeletedAtCol + " IS NULL")).
		PlaceholderFormat(squirrel.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-Gy4zW", "Errors.Internal")
	}

	result, err := r.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-Hz5aX", "Errors.Internal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-I06bY", "Errors.Internal")
	}

	if rowsAffected == 0 {
		return zerrors.ThrowNotFound(nil, "DOMAIN-J17cZ", "Errors.Domain.NotFound")
	}

	return nil
}

// SetPrimary sets the primary domain for an organization
func (r *DomainRepository) SetOrganizationDomainPrimary(ctx context.Context, instanceID, organizationID, domain string) error {
	return r.withTransaction(ctx, func(tx TX) error {
		// First, unset any existing primary domain for this organization
		unsetQuery := squirrel.Update(domainsTable).
			Set(domainsIsPrimaryCol, false).
			Set(domainsUpdatedAtCol, time.Now()).
			Where(squirrel.Eq{
				domainsInstanceIDCol: instanceID,
				domainsOrgIDCol:      organizationID,
				domainsIsPrimaryCol:  true,
			}).
			Where(squirrel.Expr(domainsDeletedAtCol + " IS NULL")).
			PlaceholderFormat(squirrel.Dollar)

		stmt, args, err := unsetQuery.ToSql()
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-K28d0", "Errors.Internal")
		}

		_, err = tx.ExecContext(ctx, stmt, args...)
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-L39e1", "Errors.Internal")
		}

		// Then set the new primary domain
		setPrimaryQuery := squirrel.Update(domainsTable).
			Set(domainsIsPrimaryCol, true).
			Set(domainsUpdatedAtCol, time.Now()).
			Where(squirrel.Eq{
				domainsInstanceIDCol: instanceID,
				domainsOrgIDCol:      organizationID,
				domainsDomainCol:     domain,
			}).
			Where(squirrel.Expr(domainsDeletedAtCol + " IS NULL")).
			PlaceholderFormat(squirrel.Dollar)

		stmt, args, err = setPrimaryQuery.ToSql()
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-M40f2", "Errors.Internal")
		}

		result, err := tx.ExecContext(ctx, stmt, args...)
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-N51g3", "Errors.Internal")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return zerrors.ThrowInternal(err, "DOMAIN-O62h4", "Errors.Internal")
		}

		if rowsAffected == 0 {
			return zerrors.ThrowNotFound(nil, "DOMAIN-P73i5", "Errors.Domain.NotFound")
		}

		return nil
	})
}

// Remove soft deletes an organization domain
func (r *DomainRepository) RemoveOrganizationDomain(ctx context.Context, instanceID, organizationID, domain string) error {
	query := squirrel.Update(domainsTable).
		Set(domainsDeletedAtCol, time.Now()).
		Set(domainsUpdatedAtCol, time.Now()).
		Where(squirrel.Eq{
			domainsInstanceIDCol: instanceID,
			domainsOrgIDCol:      organizationID,
			domainsDomainCol:     domain,
		}).
		Where(squirrel.Expr(domainsDeletedAtCol + " IS NULL")).
		PlaceholderFormat(squirrel.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-Q84j6", "Errors.Internal")
	}

	result, err := r.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-R95k7", "Errors.Internal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-S06l8", "Errors.Internal")
	}

	if rowsAffected == 0 {
		return zerrors.ThrowNotFound(nil, "DOMAIN-T17m9", "Errors.Domain.NotFound")
	}

	return nil
}

// Get returns a single domain matching the criteria
func (r *DomainRepository) Get(ctx context.Context, criteria v2domain.DomainSearchCriteria) (*v2domain.Domain, error) {
	query := r.buildSelectQuery(criteria, v2domain.DomainPagination{Limit: 2}) // Limit to 2 to detect multiple results

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-U28n0", "Errors.Internal")
	}

	rows, err := r.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-V39o1", "Errors.Internal")
	}
	defer rows.Close()

	var domains []*v2domain.Domain
	for rows.Next() {
		domain, err := r.scanDomain(rows)
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-W40p2", "Errors.Internal")
	}

	if len(domains) == 0 {
		return nil, zerrors.ThrowNotFound(nil, "DOMAIN-X51q3", "Errors.Domain.NotFound")
	}

	if len(domains) > 1 {
		return nil, zerrors.ThrowInvalidArgument(nil, "DOMAIN-Y62r4", "Errors.Domain.MultipleFound")
	}

	return domains[0], nil
}

// List returns a list of domains matching the criteria with pagination
func (r *DomainRepository) List(ctx context.Context, criteria v2domain.DomainSearchCriteria, pagination v2domain.DomainPagination) (*v2domain.DomainList, error) {
	// First get the total count
	countQuery := r.buildCountQuery(criteria)
	stmt, args, err := countQuery.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-Z73s5", "Errors.Internal")
	}

	var totalCount uint64
	err = r.client.QueryRowContext(ctx, stmt, args...).Scan(&totalCount)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-A84t6", "Errors.Internal")
	}

	// Then get the actual data
	query := r.buildSelectQuery(criteria, pagination)
	stmt, args, err = query.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-B95u7", "Errors.Internal")
	}

	rows, err := r.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-C06v8", "Errors.Internal")
	}
	defer rows.Close()

	var domains []*v2domain.Domain
	for rows.Next() {
		domain, err := r.scanDomain(rows)
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-D17w9", "Errors.Internal")
	}

	return &v2domain.DomainList{
		Domains:    domains,
		TotalCount: totalCount,
	}, nil
}

func (r *DomainRepository) buildSelectQuery(criteria v2domain.DomainSearchCriteria, pagination v2domain.DomainPagination) squirrel.SelectBuilder {
	query := squirrel.Select(
		domainsIDCol,
		domainsInstanceIDCol,
		domainsOrgIDCol,
		domainsDomainCol,
		domainsIsVerifiedCol,
		domainsIsPrimaryCol,
		domainsValidationTypeCol,
		domainsCreatedAtCol,
		domainsUpdatedAtCol,
		domainsDeletedAtCol,
	).From(domainsTable).
		PlaceholderFormat(squirrel.Dollar)

	query = r.applySearchCriteria(query, criteria)
	query = r.applyPagination(query, pagination)

	return query
}

func (r *DomainRepository) buildCountQuery(criteria v2domain.DomainSearchCriteria) squirrel.SelectBuilder {
	query := squirrel.Select("COUNT(*)").
		From(domainsTable).
		PlaceholderFormat(squirrel.Dollar)

	return r.applySearchCriteria(query, criteria)
}

func (r *DomainRepository) applySearchCriteria(query squirrel.SelectBuilder, criteria v2domain.DomainSearchCriteria) squirrel.SelectBuilder {
	// Always exclude soft-deleted domains
	query = query.Where(squirrel.Expr(domainsDeletedAtCol + " IS NULL"))

	if criteria.ID != nil {
		query = query.Where(squirrel.Eq{domainsIDCol: *criteria.ID})
	}

	if criteria.Domain != nil {
		query = query.Where(squirrel.Eq{domainsDomainCol: *criteria.Domain})
	}

	if criteria.InstanceID != nil {
		query = query.Where(squirrel.Eq{domainsInstanceIDCol: *criteria.InstanceID})
	}

	if criteria.OrganizationID != nil {
		query = query.Where(squirrel.Eq{domainsOrgIDCol: *criteria.OrganizationID})
	}

	if criteria.IsVerified != nil {
		query = query.Where(squirrel.Eq{domainsIsVerifiedCol: *criteria.IsVerified})
	}

	if criteria.IsPrimary != nil {
		query = query.Where(squirrel.Eq{domainsIsPrimaryCol: *criteria.IsPrimary})
	}

	return query
}

func (r *DomainRepository) applyPagination(query squirrel.SelectBuilder, pagination v2domain.DomainPagination) squirrel.SelectBuilder {
	// Apply sorting
	var orderBy string
	switch pagination.SortBy {
	case v2domain.DomainSortFieldCreatedAt:
		orderBy = domainsCreatedAtCol
	case v2domain.DomainSortFieldUpdatedAt:
		orderBy = domainsUpdatedAtCol
	case v2domain.DomainSortFieldDomain:
		orderBy = domainsDomainCol
	default:
		orderBy = domainsCreatedAtCol
	}

	if pagination.Order == database.SortOrderDesc {
		orderBy += " DESC"
	} else {
		orderBy += " ASC"
	}

	query = query.OrderBy(orderBy)

	// Apply pagination
	if pagination.Limit > 0 {
		query = query.Limit(uint64(pagination.Limit))
	}

	if pagination.Offset > 0 {
		query = query.Offset(uint64(pagination.Offset))
	}

	return query
}

func (r *DomainRepository) scanDomain(rows *sql.Rows) (*v2domain.Domain, error) {
	var domainRecord v2domain.Domain
	var orgID sql.NullString
	var validationType sql.NullInt32
	var deletedAt sql.NullTime

	err := rows.Scan(
		&domainRecord.ID,
		&domainRecord.InstanceID,
		&orgID,
		&domainRecord.Domain,
		&domainRecord.IsVerified,
		&domainRecord.IsPrimary,
		&validationType,
		&domainRecord.CreatedAt,
		&domainRecord.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-E28x0", "Errors.Internal")
	}

	if orgID.Valid {
		domainRecord.OrganizationID = &orgID.String
	}

	if validationType.Valid {
		validationTypeValue := domain.OrgDomainValidationType(validationType.Int32)
		domainRecord.ValidationType = &validationTypeValue
	}

	if deletedAt.Valid {
		domainRecord.DeletedAt = &deletedAt.Time
	}

	return &domainRecord, nil
}

func (r *DomainRepository) withTransaction(ctx context.Context, fn func(TX) error) error {
	tx, err := r.client.BeginTx(ctx, nil)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-F39y1", "Errors.Internal")
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return zerrors.ThrowInternal(rbErr, "DOMAIN-G40z2", "Errors.Internal")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-H51a3", "Errors.Internal")
	}

	return nil
}