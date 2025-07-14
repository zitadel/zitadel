package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	v2domain "github.com/zitadel/zitadel/internal/v2/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	domainsTable = "zitadel.domains"
)

// DomainRepository implements both InstanceDomainRepository and OrganizationDomainRepository
type DomainRepository struct {
	client *database.DB
}

// NewDomainRepository creates a new domain repository
func NewDomainRepository(client *database.DB) *DomainRepository {
	return &DomainRepository{
		client: client,
	}
}

// Add creates a new instance domain (always verified)
func (r *DomainRepository) AddInstanceDomain(ctx context.Context, instanceID, domain string) (*v2domain.Domain, error) {
	query := sq.Insert(domainsTable).
		Columns("instance_id", "domain", "is_verified", "is_primary", "validation_type").
		Values(instanceID, domain, true, false, domain.OrgDomainValidationTypeUnspecified).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(sq.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-1n8sd", "Errors.Internal")
	}

	var id string
	var createdAt, updatedAt time.Time
	err = r.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&id, &createdAt, &updatedAt)
	}, stmt, args...)
	if err != nil {
		if database.IsConstraintError(err) {
			return nil, v2domain.ErrDomainAlreadyExists
		}
		return nil, zerrors.ThrowInternal(err, "DOMAIN-2n8sd", "Errors.Internal")
	}

	return &v2domain.Domain{
		ID:             id,
		InstanceID:     instanceID,
		OrganizationID: nil,
		Domain:         domain,
		IsVerified:     true,
		IsPrimary:      false,
		ValidationType: domain.OrgDomainValidationTypeUnspecified,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

// SetInstanceDomainPrimary sets a domain as primary for an instance
func (r *DomainRepository) SetInstanceDomainPrimary(ctx context.Context, instanceID, domain string) error {
	tx, err := r.client.BeginTx(ctx, nil)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-3n8sd", "Errors.Internal")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// First, unset all primary domains for this instance
	unsetQuery := sq.Update(domainsTable).
		Set("is_primary", false).
		Set("updated_at", "now()").
		Where(sq.Eq{"instance_id": instanceID, "is_primary": true, "deleted_at": nil}).
		Where(sq.Eq{"org_id": nil}). // instance domains only
		PlaceholderFormat(sq.Dollar)

	stmt, args, err := unsetQuery.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-4n8sd", "Errors.Internal")
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-5n8sd", "Errors.Internal")
	}

	// Then set the specified domain as primary
	setPrimaryQuery := sq.Update(domainsTable).
		Set("is_primary", true).
		Set("updated_at", "now()").
		Where(sq.Eq{"instance_id": instanceID, "domain": domain, "deleted_at": nil}).
		Where(sq.Eq{"org_id": nil}). // instance domains only
		PlaceholderFormat(sq.Dollar)

	stmt, args, err = setPrimaryQuery.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-6n8sd", "Errors.Internal")
	}

	result, err := tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-7n8sd", "Errors.Internal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-8n8sd", "Errors.Internal")
	}
	if rowsAffected == 0 {
		return v2domain.ErrDomainNotFound
	}

	return tx.Commit()
}

// RemoveInstanceDomain removes an instance domain
func (r *DomainRepository) RemoveInstanceDomain(ctx context.Context, instanceID, domain string) error {
	query := sq.Update(domainsTable).
		Set("deleted_at", "now()").
		Set("updated_at", "now()").
		Where(sq.Eq{"instance_id": instanceID, "domain": domain, "deleted_at": nil}).
		Where(sq.Eq{"org_id": nil}). // instance domains only
		PlaceholderFormat(sq.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-9n8sd", "Errors.Internal")
	}

	result, err := r.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-10n8sd", "Errors.Internal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-11n8sd", "Errors.Internal")
	}
	if rowsAffected == 0 {
		return v2domain.ErrDomainNotFound
	}

	return nil
}

// GetInstanceDomain returns a single instance domain matching criteria
func (r *DomainRepository) GetInstanceDomain(ctx context.Context, criteria v2domain.DomainSearchCriteria) (*v2domain.Domain, error) {
	query := r.buildSelectQuery().
		Where(sq.Eq{"org_id": nil}). // instance domains only
		Where(sq.Eq{"deleted_at": nil})

	query = r.applyCriteria(query, criteria)

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-12n8sd", "Errors.Internal")
	}

	var domains []*v2domain.Domain
	err = r.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			domain := &v2domain.Domain{}
			var orgID sql.NullString
			var deletedAt sql.NullTime
			var validationType sql.NullInt16

			err := rows.Scan(
				&domain.ID,
				&domain.InstanceID,
				&orgID,
				&domain.Domain,
				&domain.IsVerified,
				&domain.IsPrimary,
				&validationType,
				&domain.CreatedAt,
				&domain.UpdatedAt,
				&deletedAt,
			)
			if err != nil {
				return err
			}

			if orgID.Valid {
				domain.OrganizationID = &orgID.String
			}
			if deletedAt.Valid {
				domain.DeletedAt = &deletedAt.Time
			}
			if validationType.Valid {
				domain.ValidationType = domain.OrgDomainValidationType(validationType.Int16)
			}

			domains = append(domains, domain)
		}
		return nil
	}, stmt, args...)

	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-13n8sd", "Errors.Internal")
	}

	if len(domains) == 0 {
		return nil, v2domain.ErrDomainNotFound
	}
	if len(domains) > 1 {
		return nil, v2domain.ErrMultipleDomainsFound
	}

	return domains[0], nil
}

// ListInstanceDomains returns instance domains matching criteria with pagination
func (r *DomainRepository) ListInstanceDomains(ctx context.Context, criteria v2domain.DomainSearchCriteria, pagination v2domain.DomainPagination) ([]*v2domain.Domain, int64, error) {
	// First get the count
	countQuery := sq.Select("COUNT(*)").
		From(domainsTable).
		Where(sq.Eq{"org_id": nil}). // instance domains only
		Where(sq.Eq{"deleted_at": nil})

	countQuery = r.applyCriteria(countQuery, criteria)

	countStmt, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, zerrors.ThrowInternal(err, "DOMAIN-14n8sd", "Errors.Internal")
	}

	var count int64
	err = r.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&count)
	}, countStmt, countArgs...)
	if err != nil {
		return nil, 0, zerrors.ThrowInternal(err, "DOMAIN-15n8sd", "Errors.Internal")
	}

	// Then get the actual results
	query := r.buildSelectQuery().
		Where(sq.Eq{"org_id": nil}). // instance domains only
		Where(sq.Eq{"deleted_at": nil})

	query = r.applyCriteria(query, criteria)
	query = r.applyPagination(query, pagination)

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, 0, zerrors.ThrowInternal(err, "DOMAIN-16n8sd", "Errors.Internal")
	}

	var domains []*v2domain.Domain
	err = r.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			domain := &v2domain.Domain{}
			var orgID sql.NullString
			var deletedAt sql.NullTime
			var validationType sql.NullInt16

			err := rows.Scan(
				&domain.ID,
				&domain.InstanceID,
				&orgID,
				&domain.Domain,
				&domain.IsVerified,
				&domain.IsPrimary,
				&validationType,
				&domain.CreatedAt,
				&domain.UpdatedAt,
				&deletedAt,
			)
			if err != nil {
				return err
			}

			if orgID.Valid {
				domain.OrganizationID = &orgID.String
			}
			if deletedAt.Valid {
				domain.DeletedAt = &deletedAt.Time
			}
			if validationType.Valid {
				domain.ValidationType = domain.OrgDomainValidationType(validationType.Int16)
			}

			domains = append(domains, domain)
		}
		return nil
	}, stmt, args...)

	if err != nil {
		return nil, 0, zerrors.ThrowInternal(err, "DOMAIN-17n8sd", "Errors.Internal")
	}

	return domains, count, nil
}

// Organization domain methods
func (r *DomainRepository) AddOrganizationDomain(ctx context.Context, instanceID, organizationID, domain string, validationType domain.OrgDomainValidationType) (*v2domain.Domain, error) {
	query := sq.Insert(domainsTable).
		Columns("instance_id", "org_id", "domain", "is_verified", "is_primary", "validation_type").
		Values(instanceID, organizationID, domain, false, false, int(validationType)).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(sq.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-18n8sd", "Errors.Internal")
	}

	var id string
	var createdAt, updatedAt time.Time
	err = r.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&id, &createdAt, &updatedAt)
	}, stmt, args...)
	if err != nil {
		if database.IsConstraintError(err) {
			return nil, v2domain.ErrDomainAlreadyExists
		}
		return nil, zerrors.ThrowInternal(err, "DOMAIN-19n8sd", "Errors.Internal")
	}

	return &v2domain.Domain{
		ID:             id,
		InstanceID:     instanceID,
		OrganizationID: &organizationID,
		Domain:         domain,
		IsVerified:     false,
		IsPrimary:      false,
		ValidationType: validationType,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

// SetOrganizationDomainVerified marks a domain as verified
func (r *DomainRepository) SetOrganizationDomainVerified(ctx context.Context, instanceID, organizationID, domain string) error {
	query := sq.Update(domainsTable).
		Set("is_verified", true).
		Set("updated_at", "now()").
		Where(sq.Eq{"instance_id": instanceID, "org_id": organizationID, "domain": domain, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-20n8sd", "Errors.Internal")
	}

	result, err := r.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-21n8sd", "Errors.Internal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-22n8sd", "Errors.Internal")
	}
	if rowsAffected == 0 {
		return v2domain.ErrDomainNotFound
	}

	return nil
}

// SetOrganizationDomainPrimary sets a domain as primary for an organization
func (r *DomainRepository) SetOrganizationDomainPrimary(ctx context.Context, instanceID, organizationID, domain string) error {
	tx, err := r.client.BeginTx(ctx, nil)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-23n8sd", "Errors.Internal")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// First, unset all primary domains for this organization
	unsetQuery := sq.Update(domainsTable).
		Set("is_primary", false).
		Set("updated_at", "now()").
		Where(sq.Eq{"instance_id": instanceID, "org_id": organizationID, "is_primary": true, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	stmt, args, err := unsetQuery.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-24n8sd", "Errors.Internal")
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-25n8sd", "Errors.Internal")
	}

	// Then set the specified domain as primary
	setPrimaryQuery := sq.Update(domainsTable).
		Set("is_primary", true).
		Set("updated_at", "now()").
		Where(sq.Eq{"instance_id": instanceID, "org_id": organizationID, "domain": domain, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	stmt, args, err = setPrimaryQuery.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-26n8sd", "Errors.Internal")
	}

	result, err := tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-27n8sd", "Errors.Internal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-28n8sd", "Errors.Internal")
	}
	if rowsAffected == 0 {
		return v2domain.ErrDomainNotFound
	}

	return tx.Commit()
}

// RemoveOrganizationDomain removes an organization domain
func (r *DomainRepository) RemoveOrganizationDomain(ctx context.Context, instanceID, organizationID, domain string) error {
	query := sq.Update(domainsTable).
		Set("deleted_at", "now()").
		Set("updated_at", "now()").
		Where(sq.Eq{"instance_id": instanceID, "org_id": organizationID, "domain": domain, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	stmt, args, err := query.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-29n8sd", "Errors.Internal")
	}

	result, err := r.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-30n8sd", "Errors.Internal")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOMAIN-31n8sd", "Errors.Internal")
	}
	if rowsAffected == 0 {
		return v2domain.ErrDomainNotFound
	}

	return nil
}

// GetOrganizationDomain returns a single organization domain matching criteria
func (r *DomainRepository) GetOrganizationDomain(ctx context.Context, criteria v2domain.DomainSearchCriteria) (*v2domain.Domain, error) {
	query := r.buildSelectQuery().
		Where(sq.NotEq{"org_id": nil}). // organization domains only
		Where(sq.Eq{"deleted_at": nil})

	query = r.applyCriteria(query, criteria)

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-32n8sd", "Errors.Internal")
	}

	var domains []*v2domain.Domain
	err = r.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			domain := &v2domain.Domain{}
			var orgID sql.NullString
			var deletedAt sql.NullTime
			var validationType sql.NullInt16

			err := rows.Scan(
				&domain.ID,
				&domain.InstanceID,
				&orgID,
				&domain.Domain,
				&domain.IsVerified,
				&domain.IsPrimary,
				&validationType,
				&domain.CreatedAt,
				&domain.UpdatedAt,
				&deletedAt,
			)
			if err != nil {
				return err
			}

			if orgID.Valid {
				domain.OrganizationID = &orgID.String
			}
			if deletedAt.Valid {
				domain.DeletedAt = &deletedAt.Time
			}
			if validationType.Valid {
				domain.ValidationType = domain.OrgDomainValidationType(validationType.Int16)
			}

			domains = append(domains, domain)
		}
		return nil
	}, stmt, args...)

	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DOMAIN-33n8sd", "Errors.Internal")
	}

	if len(domains) == 0 {
		return nil, v2domain.ErrDomainNotFound
	}
	if len(domains) > 1 {
		return nil, v2domain.ErrMultipleDomainsFound
	}

	return domains[0], nil
}

// ListOrganizationDomains returns organization domains matching criteria with pagination
func (r *DomainRepository) ListOrganizationDomains(ctx context.Context, criteria v2domain.DomainSearchCriteria, pagination v2domain.DomainPagination) ([]*v2domain.Domain, int64, error) {
	// First get the count
	countQuery := sq.Select("COUNT(*)").
		From(domainsTable).
		Where(sq.NotEq{"org_id": nil}). // organization domains only
		Where(sq.Eq{"deleted_at": nil})

	countQuery = r.applyCriteria(countQuery, criteria)

	countStmt, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, zerrors.ThrowInternal(err, "DOMAIN-34n8sd", "Errors.Internal")
	}

	var count int64
	err = r.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&count)
	}, countStmt, countArgs...)
	if err != nil {
		return nil, 0, zerrors.ThrowInternal(err, "DOMAIN-35n8sd", "Errors.Internal")
	}

	// Then get the actual results
	query := r.buildSelectQuery().
		Where(sq.NotEq{"org_id": nil}). // organization domains only
		Where(sq.Eq{"deleted_at": nil})

	query = r.applyCriteria(query, criteria)
	query = r.applyPagination(query, pagination)

	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, 0, zerrors.ThrowInternal(err, "DOMAIN-36n8sd", "Errors.Internal")
	}

	var domains []*v2domain.Domain
	err = r.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			domain := &v2domain.Domain{}
			var orgID sql.NullString
			var deletedAt sql.NullTime
			var validationType sql.NullInt16

			err := rows.Scan(
				&domain.ID,
				&domain.InstanceID,
				&orgID,
				&domain.Domain,
				&domain.IsVerified,
				&domain.IsPrimary,
				&validationType,
				&domain.CreatedAt,
				&domain.UpdatedAt,
				&deletedAt,
			)
			if err != nil {
				return err
			}

			if orgID.Valid {
				domain.OrganizationID = &orgID.String
			}
			if deletedAt.Valid {
				domain.DeletedAt = &deletedAt.Time
			}
			if validationType.Valid {
				domain.ValidationType = domain.OrgDomainValidationType(validationType.Int16)
			}

			domains = append(domains, domain)
		}
		return nil
	}, stmt, args...)

	if err != nil {
		return nil, 0, zerrors.ThrowInternal(err, "DOMAIN-37n8sd", "Errors.Internal")
	}

	return domains, count, nil
}

// Helper methods

func (r *DomainRepository) buildSelectQuery() sq.SelectBuilder {
	return sq.Select(
		"id",
		"instance_id",
		"org_id",
		"domain",
		"is_verified",
		"is_primary",
		"validation_type",
		"created_at",
		"updated_at",
		"deleted_at",
	).From(domainsTable).PlaceholderFormat(sq.Dollar)
}

func (r *DomainRepository) applyCriteria(query sq.SelectBuilder, criteria v2domain.DomainSearchCriteria) sq.SelectBuilder {
	if criteria.ID != nil {
		query = query.Where(sq.Eq{"id": *criteria.ID})
	}
	if criteria.InstanceID != nil {
		query = query.Where(sq.Eq{"instance_id": *criteria.InstanceID})
	}
	if criteria.OrganizationID != nil {
		query = query.Where(sq.Eq{"org_id": *criteria.OrganizationID})
	}
	if criteria.Domain != nil {
		query = query.Where(sq.Eq{"domain": *criteria.Domain})
	}
	if criteria.IsVerified != nil {
		query = query.Where(sq.Eq{"is_verified": *criteria.IsVerified})
	}
	if criteria.IsPrimary != nil {
		query = query.Where(sq.Eq{"is_primary": *criteria.IsPrimary})
	}
	return query
}

func (r *DomainRepository) applyPagination(query sq.SelectBuilder, pagination v2domain.DomainPagination) sq.SelectBuilder {
	// Apply sorting
	orderBy := "created_at"
	if pagination.SortBy != "" {
		switch pagination.SortBy {
		case "created_at", "updated_at", "domain":
			orderBy = pagination.SortBy
		}
	}

	order := "ASC"
	if strings.ToUpper(pagination.SortOrder) == "DESC" {
		order = "DESC"
	}

	query = query.OrderBy(fmt.Sprintf("%s %s", orderBy, order))

	// Apply pagination
	if pagination.Limit > 0 {
		query = query.Limit(uint64(pagination.Limit))
	}
	if pagination.Offset > 0 {
		query = query.Offset(uint64(pagination.Offset))
	}

	return query
}