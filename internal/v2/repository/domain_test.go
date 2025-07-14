package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	v2domain "github.com/zitadel/zitadel/internal/v2/domain"
)

func TestDomainRepository_AddInstanceDomain(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name        string
		instanceID  string
		domain      string
		expectError bool
		setupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			name:       "success",
			instanceID: "instance1",
			domain:     "example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow("1", time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO zitadel\.domains`).
					WithArgs("instance1", "example.com", true, false, domain.OrgDomainValidationTypeUnspecified).
					WillReturnRows(rows)
			},
		},
		{
			name:        "duplicate domain",
			instanceID:  "instance1",
			domain:      "example.com",
			expectError: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO zitadel\.domains`).
					WithArgs("instance1", "example.com", true, false, domain.OrgDomainValidationTypeUnspecified).
					WillReturnError(&database.Error{SQLCode: database.CodeUniqueViolation})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			dbClient := &database.DB{DB: db}
			repo := NewDomainRepository(dbClient)

			tt.setupMock(mock)

			result, err := repo.AddInstanceDomain(ctx, tt.instanceID, tt.domain)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.instanceID, result.InstanceID)
				assert.Equal(t, tt.domain, result.Domain)
				assert.True(t, result.IsVerified)
				assert.False(t, result.IsPrimary)
				assert.Nil(t, result.OrganizationID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDomainRepository_AddOrganizationDomain(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name           string
		instanceID     string
		organizationID string
		domain         string
		validationType domain.OrgDomainValidationType
		expectError    bool
		setupMock      func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			instanceID:     "instance1",
			organizationID: "org1",
			domain:         "org.example.com",
			validationType: domain.OrgDomainValidationTypeDNS,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow("1", time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO zitadel\.domains`).
					WithArgs("instance1", "org1", "org.example.com", false, false, int(domain.OrgDomainValidationTypeDNS)).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			dbClient := &database.DB{DB: db}
			repo := NewDomainRepository(dbClient)

			tt.setupMock(mock)

			result, err := repo.AddOrganizationDomain(ctx, tt.instanceID, tt.organizationID, tt.domain, tt.validationType)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.instanceID, result.InstanceID)
				assert.Equal(t, tt.organizationID, *result.OrganizationID)
				assert.Equal(t, tt.domain, result.Domain)
				assert.False(t, result.IsVerified)
				assert.False(t, result.IsPrimary)
				assert.Equal(t, tt.validationType, result.ValidationType)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDomainRepository_SetInstanceDomainPrimary(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name        string
		instanceID  string
		domain      string
		expectError bool
		setupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			name:       "success",
			instanceID: "instance1",
			domain:     "example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				
				// Unset existing primary domains
				mock.ExpectExec(`UPDATE zitadel\.domains SET is_primary`).
					WithArgs("instance1").
					WillReturnResult(sqlmock.NewResult(0, 1))
				
				// Set new primary domain
				mock.ExpectExec(`UPDATE zitadel\.domains SET is_primary`).
					WithArgs("instance1", "example.com").
					WillReturnResult(sqlmock.NewResult(0, 1))
				
				mock.ExpectCommit()
			},
		},
		{
			name:        "domain not found",
			instanceID:  "instance1",
			domain:      "notfound.com",
			expectError: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				
				// Unset existing primary domains
				mock.ExpectExec(`UPDATE zitadel\.domains SET is_primary`).
					WithArgs("instance1").
					WillReturnResult(sqlmock.NewResult(0, 0))
				
				// Set new primary domain - not found
				mock.ExpectExec(`UPDATE zitadel\.domains SET is_primary`).
					WithArgs("instance1", "notfound.com").
					WillReturnResult(sqlmock.NewResult(0, 0))
				
				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			dbClient := &database.DB{DB: db}
			repo := NewDomainRepository(dbClient)

			tt.setupMock(mock)

			err = repo.SetInstanceDomainPrimary(ctx, tt.instanceID, tt.domain)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, v2domain.ErrDomainNotFound, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDomainRepository_GetInstanceDomain(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name        string
		criteria    v2domain.DomainSearchCriteria
		expectError bool
		setupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success",
			criteria: v2domain.DomainSearchCriteria{
				InstanceID: stringPtr("instance1"),
				Domain:     stringPtr("example.com"),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "instance_id", "org_id", "domain", "is_verified", "is_primary", 
					"validation_type", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					"1", "instance1", nil, "example.com", true, true, 
					nil, time.Now(), time.Now(), nil,
				)
				mock.ExpectQuery(`SELECT .+ FROM zitadel\.domains`).
					WithArgs("instance1", "example.com").
					WillReturnRows(rows)
			},
		},
		{
			name: "not found",
			criteria: v2domain.DomainSearchCriteria{
				InstanceID: stringPtr("instance1"),
				Domain:     stringPtr("notfound.com"),
			},
			expectError: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "instance_id", "org_id", "domain", "is_verified", "is_primary", 
					"validation_type", "created_at", "updated_at", "deleted_at",
				})
				mock.ExpectQuery(`SELECT .+ FROM zitadel\.domains`).
					WithArgs("instance1", "notfound.com").
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			dbClient := &database.DB{DB: db}
			repo := NewDomainRepository(dbClient)

			tt.setupMock(mock)

			result, err := repo.GetInstanceDomain(ctx, tt.criteria)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "instance1", result.InstanceID)
				assert.Equal(t, "example.com", result.Domain)
				assert.True(t, result.IsVerified)
				assert.True(t, result.IsPrimary)
				assert.Nil(t, result.OrganizationID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Helper function for creating string pointers
func stringPtr(s string) *string {
	return &s
}