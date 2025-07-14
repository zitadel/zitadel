package readmodel

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	v2domain "github.com/zitadel/zitadel/internal/v2/domain"
	"github.com/zitadel/zitadel/internal/v2/database"
)

func TestDomainRepository_AddInstanceDomain(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDomainRepository(db)

	instanceID := "test-instance-id"
	domainName := "test.example.com"
	expectedID := "domain-id-123"

	mock.ExpectQuery(`INSERT INTO zitadel\.domains`).
		WithArgs(instanceID, domainName, true, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	result, err := repo.AddInstanceDomain(context.Background(), instanceID, domainName)

	require.NoError(t, err)
	assert.Equal(t, expectedID, result.ID)
	assert.Equal(t, instanceID, result.InstanceID)
	assert.Nil(t, result.OrganizationID)
	assert.Equal(t, domainName, result.Domain)
	assert.True(t, result.IsVerified)
	assert.False(t, result.IsPrimary)
	assert.Nil(t, result.ValidationType)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainRepository_AddOrganizationDomain(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDomainRepository(db)

	instanceID := "test-instance-id"
	orgID := "test-org-id"
	domainName := "test.example.com"
	validationType := domain.OrgDomainValidationTypeHTTP
	expectedID := "domain-id-456"

	mock.ExpectQuery(`INSERT INTO zitadel\.domains`).
		WithArgs(instanceID, orgID, domainName, false, false, int(validationType), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	result, err := repo.AddOrganizationDomain(context.Background(), instanceID, orgID, domainName, validationType)

	require.NoError(t, err)
	assert.Equal(t, expectedID, result.ID)
	assert.Equal(t, instanceID, result.InstanceID)
	assert.Equal(t, orgID, *result.OrganizationID)
	assert.Equal(t, domainName, result.Domain)
	assert.False(t, result.IsVerified)
	assert.False(t, result.IsPrimary)
	assert.Equal(t, validationType, *result.ValidationType)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainRepository_SetInstanceDomainPrimary(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDomainRepository(db)

	instanceID := "test-instance-id"
	domainName := "test.example.com"

	// Mock transaction begin
	mock.ExpectBegin()

	// Mock unset existing primary
	mock.ExpectExec(`UPDATE zitadel\.domains SET.*is_primary.*=.*false`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Mock set new primary
	mock.ExpectExec(`UPDATE zitadel\.domains SET.*is_primary.*=.*true`).
		WithArgs(sqlmock.AnyArg(), true, instanceID, nil, domainName).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Mock transaction commit
	mock.ExpectCommit()

	err = repo.SetInstanceDomainPrimary(context.Background(), instanceID, domainName)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDomainRepository(db)

	instanceID := "test-instance-id"
	domainName := "test.example.com"
	now := time.Now()

	criteria := v2domain.DomainSearchCriteria{
		InstanceID: &instanceID,
		Domain:     &domainName,
	}

	rows := sqlmock.NewRows([]string{
		"id", "instance_id", "org_id", "domain", "is_verified", "is_primary", "validation_type", "created_at", "updated_at", "deleted_at",
	}).AddRow("domain-123", instanceID, nil, domainName, true, false, nil, now, now, nil)

	mock.ExpectQuery(`SELECT .* FROM zitadel\.domains`).
		WithArgs(domainName, instanceID).
		WillReturnRows(rows)

	result, err := repo.Get(context.Background(), criteria)

	require.NoError(t, err)
	assert.Equal(t, "domain-123", result.ID)
	assert.Equal(t, instanceID, result.InstanceID)
	assert.Nil(t, result.OrganizationID)
	assert.Equal(t, domainName, result.Domain)
	assert.True(t, result.IsVerified)
	assert.False(t, result.IsPrimary)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDomainRepository(db)

	instanceID := "test-instance-id"
	now := time.Now()

	criteria := v2domain.DomainSearchCriteria{
		InstanceID: &instanceID,
	}

	pagination := v2domain.DomainPagination{
		Limit:  10,
		Offset: 0,
		SortBy: v2domain.DomainSortFieldDomain,
		Order:  database.SortOrderAsc,
	}

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM zitadel\.domains`).
		WithArgs(instanceID).
		WillReturnRows(countRows)

	// Mock data query
	rows := sqlmock.NewRows([]string{
		"id", "instance_id", "org_id", "domain", "is_verified", "is_primary", "validation_type", "created_at", "updated_at", "deleted_at",
	}).
		AddRow("domain-instance", instanceID, nil, "instance.example.com", true, true, nil, now, now, nil).
		AddRow("domain-org", instanceID, "org-id", "org.example.com", false, false, int(domain.OrgDomainValidationTypeHTTP), now, now, nil)

	mock.ExpectQuery(`SELECT .* FROM zitadel\.domains.*ORDER BY domain ASC.*LIMIT 10`).
		WithArgs(instanceID).
		WillReturnRows(rows)

	result, err := repo.List(context.Background(), criteria, pagination)

	require.NoError(t, err)
	assert.Equal(t, uint64(2), result.TotalCount)
	assert.Len(t, result.Domains, 2)

	// Check first domain (instance domain)
	assert.Equal(t, "domain-instance", result.Domains[0].ID)
	assert.Equal(t, instanceID, result.Domains[0].InstanceID)
	assert.Nil(t, result.Domains[0].OrganizationID)
	assert.Equal(t, "instance.example.com", result.Domains[0].Domain)
	assert.True(t, result.Domains[0].IsVerified)
	assert.True(t, result.Domains[0].IsPrimary)

	// Check second domain (org domain)
	assert.Equal(t, "domain-org", result.Domains[1].ID)
	assert.Equal(t, instanceID, result.Domains[1].InstanceID)
	assert.Equal(t, "org-id", *result.Domains[1].OrganizationID)
	assert.Equal(t, "org.example.com", result.Domains[1].Domain)
	assert.False(t, result.Domains[1].IsVerified)
	assert.False(t, result.Domains[1].IsPrimary)

	assert.NoError(t, mock.ExpectationsWereMet())
}