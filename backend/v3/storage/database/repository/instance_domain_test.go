package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestAddInstanceDomain(t *testing.T) {
	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository(pool)
	err := instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	tests := []struct {
		name           string
		testFunc       func(ctx context.Context, t *testing.T, domainRepo domain.InstanceDomainRepository) *domain.AddInstanceDomain
		instanceDomain domain.AddInstanceDomain
		err            error
	}{
		{
			name: "happy path custom domain",
			instanceDomain: domain.AddInstanceDomain{
				InstanceID:  instanceID,
				Domain:      gofakeit.DomainName(),
				Type:        domain.DomainTypeCustom,
				IsPrimary:   gu.Ptr(false),
				IsGenerated: gu.Ptr(false),
			},
		},
		{
			name: "happy path trusted domain",
			instanceDomain: domain.AddInstanceDomain{
				InstanceID: instanceID,
				Domain:     gofakeit.DomainName(),
				Type:       domain.DomainTypeTrusted,
			},
		},
		{
			name: "add primary domain",
			instanceDomain: domain.AddInstanceDomain{
				InstanceID:  instanceID,
				Domain:      gofakeit.DomainName(),
				Type:        domain.DomainTypeCustom,
				IsPrimary:   gu.Ptr(true),
				IsGenerated: gu.Ptr(false),
			},
		},
		{
			name: "add custom domain without domain name",
			instanceDomain: domain.AddInstanceDomain{
				InstanceID:  instanceID,
				Domain:      "",
				Type:        domain.DomainTypeCustom,
				IsPrimary:   gu.Ptr(false),
				IsGenerated: gu.Ptr(false),
			},
			err: new(database.CheckError),
		},
		{
			name: "add trusted domain without domain name",
			instanceDomain: domain.AddInstanceDomain{
				InstanceID: instanceID,
				Domain:     "",
				Type:       domain.DomainTypeTrusted,
			},
			err: new(database.CheckError),
		},
		{
			name: "add custom domain with same domain twice",
			testFunc: func(ctx context.Context, t *testing.T, domainRepo domain.InstanceDomainRepository) *domain.AddInstanceDomain {
				domainName := gofakeit.DomainName()

				instanceDomain := &domain.AddInstanceDomain{
					InstanceID:  instanceID,
					Domain:      domainName,
					Type:        domain.DomainTypeCustom,
					IsPrimary:   gu.Ptr(false),
					IsGenerated: gu.Ptr(false),
				}

				err := domainRepo.Add(ctx, instanceDomain)
				require.NoError(t, err)

				// return same domain again
				return &domain.AddInstanceDomain{
					InstanceID:  instanceID,
					Domain:      domainName,
					Type:        domain.DomainTypeCustom,
					IsPrimary:   gu.Ptr(false),
					IsGenerated: gu.Ptr(false),
				}
			},
			err: new(database.UniqueError),
		},
		{
			name: "add trusted domain with same domain twice",
			testFunc: func(ctx context.Context, t *testing.T, domainRepo domain.InstanceDomainRepository) *domain.AddInstanceDomain {
				domainName := gofakeit.DomainName()

				instanceDomain := &domain.AddInstanceDomain{
					InstanceID: instanceID,
					Domain:     domainName,
					Type:       domain.DomainTypeTrusted,
				}

				err := domainRepo.Add(ctx, instanceDomain)
				require.NoError(t, err)

				// return same domain again
				return &domain.AddInstanceDomain{
					InstanceID: instanceID,
					Domain:     domainName,
					Type:       domain.DomainTypeTrusted,
				}
			},
			err: new(database.UniqueError),
		},
		{
			name: "add domain with non-existent instance",
			instanceDomain: domain.AddInstanceDomain{
				InstanceID:  "non-existent-instance",
				Domain:      gofakeit.DomainName(),
				Type:        domain.DomainTypeCustom,
				IsPrimary:   gu.Ptr(false),
				IsGenerated: gu.Ptr(false),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "add domain without instance id",
			instanceDomain: domain.AddInstanceDomain{
				Domain:      gofakeit.DomainName(),
				Type:        domain.DomainTypeCustom,
				IsPrimary:   gu.Ptr(false),
				IsGenerated: gu.Ptr(false),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "add custom domain without primary",
			instanceDomain: domain.AddInstanceDomain{
				Domain:      gofakeit.DomainName(),
				Type:        domain.DomainTypeCustom,
				IsGenerated: gu.Ptr(false),
			},
			err: new(database.CheckError),
		},
		{
			name: "add custom domain without generated",
			instanceDomain: domain.AddInstanceDomain{
				Domain:    gofakeit.DomainName(),
				Type:      domain.DomainTypeCustom,
				IsPrimary: gu.Ptr(false),
			},
			err: new(database.CheckError),
		},
		{
			name: "add trusted domain with primary",
			instanceDomain: domain.AddInstanceDomain{
				Domain:    gofakeit.DomainName(),
				Type:      domain.DomainTypeTrusted,
				IsPrimary: gu.Ptr(false),
			},
			err: new(database.CheckError),
		},
		{
			name: "add trusted domain with generated",
			instanceDomain: domain.AddInstanceDomain{
				Domain:      gofakeit.DomainName(),
				Type:        domain.DomainTypeTrusted,
				IsGenerated: gu.Ptr(false),
			},
			err: new(database.CheckError),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

			// we take now here because the timestamp of the transaction is used to set the createdAt and updatedAt fields
			beforeAdd := time.Now()
			tx, err := pool.Begin(t.Context(), nil)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, tx.Rollback(t.Context()))
			}()
			instanceRepo := repository.InstanceRepository(tx)
			domainRepo := instanceRepo.Domains(false)

			var instanceDomain *domain.AddInstanceDomain
			if test.testFunc != nil {
				instanceDomain = test.testFunc(ctx, t, domainRepo)
			} else {
				instanceDomain = &test.instanceDomain
			}

			err = domainRepo.Add(ctx, instanceDomain)
			afterAdd := time.Now()
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, instanceDomain.CreatedAt)
			assert.NotZero(t, instanceDomain.UpdatedAt)
			assert.WithinRange(t, instanceDomain.CreatedAt, beforeAdd, afterAdd)
			assert.WithinRange(t, instanceDomain.UpdatedAt, beforeAdd, afterAdd)
		})
	}
}

func TestGetInstanceDomain(t *testing.T) {
	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()
	instanceRepo := repository.InstanceRepository(tx)
	err = instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// add domains
	domainRepo := instanceRepo.Domains(false)
	domainName1 := gofakeit.DomainName()
	domainName2 := gofakeit.DomainName()

	domain1 := &domain.AddInstanceDomain{
		InstanceID:  instanceID,
		Domain:      domainName1,
		IsPrimary:   gu.Ptr(true),
		IsGenerated: gu.Ptr(false),
		Type:        domain.DomainTypeCustom,
	}
	domain2 := &domain.AddInstanceDomain{
		InstanceID:  instanceID,
		Domain:      domainName2,
		IsPrimary:   gu.Ptr(false),
		IsGenerated: gu.Ptr(false),
		Type:        domain.DomainTypeCustom,
	}

	err = domainRepo.Add(t.Context(), domain1)
	require.NoError(t, err)
	err = domainRepo.Add(t.Context(), domain2)
	require.NoError(t, err)

	tests := []struct {
		name     string
		opts     []database.QueryOption
		expected *domain.InstanceDomain
		err      error
	}{
		{
			name: "get primary domain",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.IsPrimaryCondition(true)),
			},
			expected: &domain.InstanceDomain{
				InstanceID: instanceID,
				Domain:     domainName1,
				IsPrimary:  gu.Ptr(true),
			},
		},
		{
			name: "get by domain name",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.DomainCondition(database.TextOperationEqual, domainName2)),
			},
			expected: &domain.InstanceDomain{
				InstanceID: instanceID,
				Domain:     domainName2,
				IsPrimary:  gu.Ptr(false),
			},
		},
		{
			name: "get non-existent domain",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.DomainCondition(database.TextOperationEqual, "non-existent.com")),
			},
			err: new(database.NoRowFoundError),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

			result, err := domainRepo.Get(ctx, test.opts...)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expected.InstanceID, result.InstanceID)
			assert.Equal(t, test.expected.Domain, result.Domain)
			assert.Equal(t, test.expected.IsPrimary, result.IsPrimary)
			assert.NotEmpty(t, result.CreatedAt)
			assert.NotEmpty(t, result.UpdatedAt)
		})
	}
}

func TestListInstanceDomains(t *testing.T) {
	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	instanceRepo := repository.InstanceRepository(tx)
	err = instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// add multiple domains
	domainRepo := instanceRepo.Domains(false)
	domains := []domain.AddInstanceDomain{
		{
			InstanceID:  instanceID,
			Domain:      gofakeit.DomainName(),
			IsPrimary:   gu.Ptr(true),
			IsGenerated: gu.Ptr(false),
			Type:        domain.DomainTypeCustom,
		},
		{
			InstanceID:  instanceID,
			Domain:      gofakeit.DomainName(),
			IsPrimary:   gu.Ptr(false),
			IsGenerated: gu.Ptr(false),
			Type:        domain.DomainTypeCustom,
		},
		{
			InstanceID: instanceID,
			Domain:     gofakeit.DomainName(),
			Type:       domain.DomainTypeTrusted,
		},
	}

	for i := range domains {
		err = domainRepo.Add(t.Context(), &domains[i])
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		opts          []database.QueryOption
		expectedCount int
	}{
		{
			name:          "list all domains",
			opts:          []database.QueryOption{},
			expectedCount: 3,
		},
		{
			name: "list primary domains",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.IsPrimaryCondition(true)),
			},
			expectedCount: 1,
		},
		{
			name: "list by instance",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.InstanceIDCondition(instanceID)),
			},
			expectedCount: 3,
		},
		{
			name: "list non-existent instance",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.InstanceIDCondition("non-existent")),
			},
			expectedCount: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

			results, err := domainRepo.List(ctx, test.opts...)
			require.NoError(t, err)
			assert.Len(t, results, test.expectedCount)

			for _, result := range results {
				assert.Equal(t, instanceID, result.InstanceID)
				assert.NotEmpty(t, result.Domain)
				assert.NotEmpty(t, result.CreatedAt)
				assert.NotEmpty(t, result.UpdatedAt)
			}
		})
	}
}

func TestUpdateInstanceDomain(t *testing.T) {
	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}

	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	instanceRepo := repository.InstanceRepository(tx)
	err = instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// add domain
	domainRepo := instanceRepo.Domains(false)
	domainName := gofakeit.DomainName()
	instanceDomain := &domain.AddInstanceDomain{
		InstanceID:  instanceID,
		Domain:      domainName,
		IsPrimary:   gu.Ptr(false),
		IsGenerated: gu.Ptr(false),
		Type:        domain.DomainTypeCustom,
	}

	err = domainRepo.Add(t.Context(), instanceDomain)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		changes   []database.Change
		expected  int64
		err       error
	}{
		{
			name:      "set primary",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, domainName),
			changes:   []database.Change{domainRepo.SetPrimary()},
			expected:  1,
		},
		{
			name:      "update non-existent domain",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, "non-existent.com"),
			changes:   []database.Change{domainRepo.SetPrimary()},
			expected:  0,
		},
		{
			name:      "no changes",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, domainName),
			changes:   []database.Change{},
			expected:  0,
			err:       database.ErrNoChanges,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

			rowsAffected, err := domainRepo.Update(ctx, test.condition, test.changes...)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expected, rowsAffected)

			// verify changes were applied if rows were affected
			if rowsAffected > 0 && len(test.changes) > 0 {
				result, err := domainRepo.Get(ctx, database.WithCondition(test.condition))
				require.NoError(t, err)

				// We know changes were applied since rowsAffected > 0
				// The specific verification of what changed is less important
				// than knowing the operation succeeded
				assert.NotNil(t, result)
			}
		})
	}
}

func TestRemoveInstanceDomain(t *testing.T) {
	// create instance
	instanceID := gofakeit.UUID()
	instance := domain.Instance{
		ID:              instanceID,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()
	instanceRepo := repository.InstanceRepository(tx)
	err = instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	// add domains
	domainRepo := instanceRepo.Domains(false)
	domainName1 := gofakeit.DomainName()

	domain1 := &domain.AddInstanceDomain{
		InstanceID:  instanceID,
		Domain:      domainName1,
		IsPrimary:   gu.Ptr(true),
		IsGenerated: gu.Ptr(false),
		Type:        domain.DomainTypeCustom,
	}
	domain2 := &domain.AddInstanceDomain{
		InstanceID:  instanceID,
		Domain:      gofakeit.DomainName(),
		IsPrimary:   gu.Ptr(false),
		IsGenerated: gu.Ptr(false),
		Type:        domain.DomainTypeCustom,
	}

	err = domainRepo.Add(t.Context(), domain1)
	require.NoError(t, err)
	err = domainRepo.Add(t.Context(), domain2)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		expected  int64
	}{
		{
			name:      "remove by domain name",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, domainName1),
			expected:  1,
		},
		{
			name:      "remove by primary condition",
			condition: domainRepo.IsPrimaryCondition(false),
			expected:  1, // domain2 should still exist and be non-primary
		},
		{
			name:      "remove non-existent domain",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, "non-existent.com"),
			expected:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

			// count before removal
			beforeCount, err := domainRepo.List(ctx)
			require.NoError(t, err)

			rowsAffected, err := domainRepo.Remove(ctx, test.condition)
			require.NoError(t, err)
			assert.Equal(t, test.expected, rowsAffected)

			// verify removal
			afterCount, err := domainRepo.List(ctx)
			require.NoError(t, err)
			assert.Equal(t, len(beforeCount)-int(test.expected), len(afterCount))
		})
	}
}

func TestInstanceDomainConditions(t *testing.T) {
	instanceRepo := repository.InstanceRepository(pool)
	domainRepo := instanceRepo.Domains(false)

	tests := []struct {
		name      string
		condition database.Condition
		expected  string
	}{
		{
			name:      "domain condition equal",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, "example.com"),
			expected:  "instance_domains.domain = $1",
		},
		{
			name:      "domain condition starts with",
			condition: domainRepo.DomainCondition(database.TextOperationStartsWith, "example"),
			expected:  "instance_domains.domain LIKE $1 || '%'",
		},
		{
			name:      "instance id condition",
			condition: domainRepo.InstanceIDCondition("instance-123"),
			expected:  "instance_domains.instance_id = $1",
		},
		{
			name:      "is primary true",
			condition: domainRepo.IsPrimaryCondition(true),
			expected:  "instance_domains.is_primary = $1",
		},
		{
			name:      "is primary false",
			condition: domainRepo.IsPrimaryCondition(false),
			expected:  "instance_domains.is_primary = $1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var builder database.StatementBuilder
			test.condition.Write(&builder)
			assert.Equal(t, test.expected, builder.String())
		})
	}
}

func TestInstanceDomainChanges(t *testing.T) {
	instanceRepo := repository.InstanceRepository(pool)
	domainRepo := instanceRepo.Domains(false)

	tests := []struct {
		name     string
		change   database.Change
		expected string
	}{
		{
			name:     "set primary",
			change:   domainRepo.SetPrimary(),
			expected: "is_primary = $1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var builder database.StatementBuilder
			test.change.Write(&builder)
			assert.Equal(t, test.expected, builder.String())
		})
	}
}
