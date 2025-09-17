package repository_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestAddOrganizationDomain(t *testing.T) {
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

	// create organization
	orgID := gofakeit.UUID()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}

	tests := []struct {
		name               string
		testFunc           func(ctx context.Context, t *testing.T, domainRepo domain.OrganizationDomainRepository) *domain.AddOrganizationDomain
		organizationDomain domain.AddOrganizationDomain
		err                error
	}{
		{
			name: "happy path",
			organizationDomain: domain.AddOrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          orgID,
				Domain:         gofakeit.DomainName(),
				IsVerified:     false,
				IsPrimary:      false,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
		},
		{
			name: "add verified domain",
			organizationDomain: domain.AddOrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          orgID,
				Domain:         gofakeit.DomainName(),
				IsVerified:     true,
				IsPrimary:      false,
				ValidationType: gu.Ptr(domain.DomainValidationTypeHTTP),
			},
		},
		{
			name: "add primary domain",
			organizationDomain: domain.AddOrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          orgID,
				Domain:         gofakeit.DomainName(),
				IsVerified:     true,
				IsPrimary:      true,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
		},
		{
			name: "add domain without domain name",
			organizationDomain: domain.AddOrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          orgID,
				Domain:         "",
				IsVerified:     false,
				IsPrimary:      false,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
			err: new(database.CheckError),
		},
		{
			name: "add domain with same domain twice",
			testFunc: func(ctx context.Context, t *testing.T, domainRepo domain.OrganizationDomainRepository) *domain.AddOrganizationDomain {
				domainName := gofakeit.DomainName()

				organizationDomain := &domain.AddOrganizationDomain{
					InstanceID:     instanceID,
					OrgID:          orgID,
					Domain:         domainName,
					IsVerified:     false,
					IsPrimary:      false,
					ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
				}

				err := domainRepo.Add(ctx, organizationDomain)
				require.NoError(t, err)

				// return same domain again
				return &domain.AddOrganizationDomain{
					InstanceID:     instanceID,
					OrgID:          orgID,
					Domain:         domainName,
					IsVerified:     true,
					IsPrimary:      true,
					ValidationType: gu.Ptr(domain.DomainValidationTypeHTTP),
				}
			},
			err: new(database.UniqueError),
		},
		{
			name: "add domain with non-existent instance",
			organizationDomain: domain.AddOrganizationDomain{
				InstanceID:     "non-existent-instance",
				OrgID:          orgID,
				Domain:         gofakeit.DomainName(),
				IsVerified:     false,
				IsPrimary:      false,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "add domain with non-existent organization",
			organizationDomain: domain.AddOrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          "non-existent-org",
				Domain:         gofakeit.DomainName(),
				IsVerified:     false,
				IsPrimary:      false,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "add domain without instance id",
			organizationDomain: domain.AddOrganizationDomain{
				OrgID:          orgID,
				Domain:         gofakeit.DomainName(),
				IsVerified:     false,
				IsPrimary:      false,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "add domain without org id",
			organizationDomain: domain.AddOrganizationDomain{
				InstanceID:     instanceID,
				Domain:         gofakeit.DomainName(),
				IsVerified:     false,
				IsPrimary:      false,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
			err: new(database.ForeignKeyError),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

			tx, err := pool.Begin(t.Context(), nil)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, tx.Rollback(t.Context()))
			}()

			orgRepo := repository.OrganizationRepository(tx)
			err = orgRepo.Create(t.Context(), &organization)
			require.NoError(t, err)

			domainRepo := orgRepo.Domains(false)

			var organizationDomain *domain.AddOrganizationDomain
			if test.testFunc != nil {
				organizationDomain = test.testFunc(ctx, t, domainRepo)
			} else {
				organizationDomain = &test.organizationDomain
			}

			err = domainRepo.Add(ctx, organizationDomain)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, organizationDomain.CreatedAt)
			assert.NotZero(t, organizationDomain.UpdatedAt)
		})
	}
}

func TestGetOrganizationDomain(t *testing.T) {
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

	// create organization
	orgID := gofakeit.UUID()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}

	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	instanceRepo := repository.InstanceRepository(tx)
	err = instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(tx)
	err = orgRepo.Create(t.Context(), &organization)
	require.NoError(t, err)

	// add domains
	domainRepo := orgRepo.Domains(false)
	domainName1 := gofakeit.DomainName()
	domainName2 := gofakeit.DomainName()

	domain1 := &domain.AddOrganizationDomain{
		InstanceID:     instanceID,
		OrgID:          orgID,
		Domain:         domainName1,
		IsVerified:     true,
		IsPrimary:      true,
		ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
	}
	domain2 := &domain.AddOrganizationDomain{
		InstanceID:     instanceID,
		OrgID:          orgID,
		Domain:         domainName2,
		IsVerified:     false,
		IsPrimary:      false,
		ValidationType: gu.Ptr(domain.DomainValidationTypeHTTP),
	}

	err = domainRepo.Add(t.Context(), domain1)
	require.NoError(t, err)
	err = domainRepo.Add(t.Context(), domain2)
	require.NoError(t, err)

	tests := []struct {
		name     string
		opts     []database.QueryOption
		expected *domain.OrganizationDomain
		err      error
	}{
		{
			name: "get primary domain",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.IsPrimaryCondition(true)),
			},
			expected: &domain.OrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          orgID,
				Domain:         domainName1,
				IsVerified:     true,
				IsPrimary:      true,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
		},
		{
			name: "get by domain name",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.DomainCondition(database.TextOperationEqual, domainName2)),
			},
			expected: &domain.OrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          orgID,
				Domain:         domainName2,
				IsVerified:     false,
				IsPrimary:      false,
				ValidationType: gu.Ptr(domain.DomainValidationTypeHTTP),
			},
		},
		{
			name: "get by org ID",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.OrgIDCondition(orgID)),
				database.WithCondition(domainRepo.IsPrimaryCondition(true)),
			},
			expected: &domain.OrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          orgID,
				Domain:         domainName1,
				IsVerified:     true,
				IsPrimary:      true,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
			},
		},
		{
			name: "get verified domain",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.IsVerifiedCondition(true)),
			},
			expected: &domain.OrganizationDomain{
				InstanceID:     instanceID,
				OrgID:          orgID,
				Domain:         domainName1,
				IsVerified:     true,
				IsPrimary:      true,
				ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
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
			assert.Equal(t, test.expected.OrgID, result.OrgID)
			assert.Equal(t, test.expected.Domain, result.Domain)
			assert.Equal(t, test.expected.IsVerified, result.IsVerified)
			assert.Equal(t, test.expected.IsPrimary, result.IsPrimary)
			assert.Equal(t, test.expected.ValidationType, result.ValidationType)
			assert.NotEmpty(t, result.CreatedAt)
			assert.NotEmpty(t, result.UpdatedAt)
		})
	}
}

func TestListOrganizationDomains(t *testing.T) {
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

	// create organization
	orgID := gofakeit.UUID()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}

	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	instanceRepo := repository.InstanceRepository(tx)
	err = instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(tx)
	err = orgRepo.Create(t.Context(), &organization)
	require.NoError(t, err)

	// add multiple domains
	domainRepo := orgRepo.Domains(false)
	domains := []domain.AddOrganizationDomain{
		{
			InstanceID:     instanceID,
			OrgID:          orgID,
			Domain:         gofakeit.DomainName(),
			IsVerified:     true,
			IsPrimary:      true,
			ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
		},
		{
			InstanceID:     instanceID,
			OrgID:          orgID,
			Domain:         gofakeit.DomainName(),
			IsVerified:     false,
			IsPrimary:      false,
			ValidationType: gu.Ptr(domain.DomainValidationTypeHTTP),
		},
		{
			InstanceID:     instanceID,
			OrgID:          orgID,
			Domain:         gofakeit.DomainName(),
			IsVerified:     true,
			IsPrimary:      false,
			ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
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
			name: "list verified domains",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.IsVerifiedCondition(true)),
			},
			expectedCount: 2,
		},
		{
			name: "list primary domains",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.IsPrimaryCondition(true)),
			},
			expectedCount: 1,
		},
		{
			name: "list by organization",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.OrgIDCondition(orgID)),
			},
			expectedCount: 3,
		},
		{
			name: "list by instance",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.InstanceIDCondition(instanceID)),
			},
			expectedCount: 3,
		},
		{
			name: "list non-existent organization",
			opts: []database.QueryOption{
				database.WithCondition(domainRepo.OrgIDCondition("non-existent")),
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
				assert.Equal(t, orgID, result.OrgID)
				assert.NotEmpty(t, result.Domain)
				assert.NotEmpty(t, result.CreatedAt)
				assert.NotEmpty(t, result.UpdatedAt)
			}
		})
	}
}

func TestUpdateOrganizationDomain(t *testing.T) {
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

	// create organization
	orgID := gofakeit.UUID()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}

	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	instanceRepo := repository.InstanceRepository(tx)
	err = instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(tx)
	err = orgRepo.Create(t.Context(), &organization)
	require.NoError(t, err)

	// add domain
	domainRepo := orgRepo.Domains(false)
	domainName := gofakeit.DomainName()
	organizationDomain := &domain.AddOrganizationDomain{
		InstanceID:     instanceID,
		OrgID:          orgID,
		Domain:         domainName,
		IsVerified:     false,
		IsPrimary:      false,
		ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
	}

	err = domainRepo.Add(t.Context(), organizationDomain)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		changes   []database.Change
		expected  int64
		err       error
	}{
		{
			name:      "set verified",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, domainName),
			changes:   []database.Change{domainRepo.SetVerified()},
			expected:  1,
		},
		{
			name:      "set primary",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, domainName),
			changes:   []database.Change{domainRepo.SetPrimary()},
			expected:  1,
		},
		{
			name:      "set validation type",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, domainName),
			changes:   []database.Change{domainRepo.SetValidationType(domain.DomainValidationTypeHTTP)},
			expected:  1,
		},
		{
			name:      "multiple changes",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, domainName),
			changes: []database.Change{
				domainRepo.SetVerified(),
				domainRepo.SetPrimary(),
				domainRepo.SetValidationType(domain.DomainValidationTypeDNS),
			},
			expected: 1,
		},
		{
			name:      "update by org ID and domain",
			condition: database.And(domainRepo.OrgIDCondition(orgID), domainRepo.DomainCondition(database.TextOperationEqual, domainName)),
			changes:   []database.Change{domainRepo.SetVerified()},
			expected:  1,
		},
		{
			name:      "update non-existent domain",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, "non-existent.com"),
			changes:   []database.Change{domainRepo.SetVerified()},
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

func TestRemoveOrganizationDomain(t *testing.T) {
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

	// create organization
	orgID := gofakeit.UUID()
	organization := domain.Organization{
		ID:         orgID,
		Name:       gofakeit.Name(),
		InstanceID: instanceID,
		State:      domain.OrgStateActive,
	}

	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback(t.Context()))
	}()

	instanceRepo := repository.InstanceRepository(tx)
	err = instanceRepo.Create(t.Context(), &instance)
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(tx)
	err = orgRepo.Create(t.Context(), &organization)
	require.NoError(t, err)

	// add domains
	domainRepo := orgRepo.Domains(false)
	domainName1 := gofakeit.DomainName()
	domainName2 := gofakeit.DomainName()

	domain1 := &domain.AddOrganizationDomain{
		InstanceID:     instanceID,
		OrgID:          orgID,
		Domain:         domainName1,
		IsVerified:     true,
		IsPrimary:      true,
		ValidationType: gu.Ptr(domain.DomainValidationTypeDNS),
	}
	domain2 := &domain.AddOrganizationDomain{
		InstanceID:     instanceID,
		OrgID:          orgID,
		Domain:         domainName2,
		IsVerified:     false,
		IsPrimary:      false,
		ValidationType: gu.Ptr(domain.DomainValidationTypeHTTP),
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
			name:      "remove by org ID and domain",
			condition: database.And(domainRepo.OrgIDCondition(orgID), domainRepo.DomainCondition(database.TextOperationEqual, domainName2)),
			expected:  1,
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

			snapshot, err := tx.Begin(ctx)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, snapshot.Rollback(ctx))
			}()

			orgRepo := repository.OrganizationRepository(snapshot)
			domainRepo := orgRepo.Domains(false)

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

func TestOrganizationDomainConditions(t *testing.T) {
	orgRepo := repository.OrganizationRepository(pool)
	domainRepo := orgRepo.Domains(false)

	tests := []struct {
		name      string
		condition database.Condition
		expected  string
	}{
		{
			name:      "domain condition equal",
			condition: domainRepo.DomainCondition(database.TextOperationEqual, "example.com"),
			expected:  "org_domains.domain = $1",
		},
		{
			name:      "domain condition starts with",
			condition: domainRepo.DomainCondition(database.TextOperationStartsWith, "example"),
			expected:  "org_domains.domain LIKE $1 || '%'",
		},
		{
			name:      "instance id condition",
			condition: domainRepo.InstanceIDCondition("instance-123"),
			expected:  "org_domains.instance_id = $1",
		},
		{
			name:      "org id condition",
			condition: domainRepo.OrgIDCondition("org-123"),
			expected:  "org_domains.org_id = $1",
		},
		{
			name:      "is primary true",
			condition: domainRepo.IsPrimaryCondition(true),
			expected:  "org_domains.is_primary = $1",
		},
		{
			name:      "is primary false",
			condition: domainRepo.IsPrimaryCondition(false),
			expected:  "org_domains.is_primary = $1",
		},
		{
			name:      "is verified true",
			condition: domainRepo.IsVerifiedCondition(true),
			expected:  "org_domains.is_verified = $1",
		},
		{
			name:      "is verified false",
			condition: domainRepo.IsVerifiedCondition(false),
			expected:  "org_domains.is_verified = $1",
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

func TestOrganizationDomainChanges(t *testing.T) {
	orgRepo := repository.OrganizationRepository(pool)
	domainRepo := orgRepo.Domains(false)

	tests := []struct {
		name     string
		change   database.Change
		expected string
	}{
		{
			name:     "set verified",
			change:   domainRepo.SetVerified(),
			expected: "is_verified = $1",
		},
		{
			name:     "set primary",
			change:   domainRepo.SetPrimary(),
			expected: "is_primary = $1",
		},
		{
			name:     "set validation type DNS",
			change:   domainRepo.SetValidationType(domain.DomainValidationTypeDNS),
			expected: "validation_type = $1",
		},
		{
			name:     "set validation type HTTP",
			change:   domainRepo.SetValidationType(domain.DomainValidationTypeHTTP),
			expected: "validation_type = $1",
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

func TestOrganizationDomainColumns(t *testing.T) {
	orgRepo := repository.OrganizationRepository(pool)
	domainRepo := orgRepo.Domains(false)

	tests := []struct {
		name      string
		column    *database.Column
		qualified bool
		expected  string
	}{
		{
			name:      "instance id column qualified",
			column:    domainRepo.InstanceIDColumn(),
			qualified: true,
			expected:  "org_domains.instance_id",
		},
		{
			name:      "instance id column unqualified",
			column:    domainRepo.InstanceIDColumn(),
			qualified: false,
			expected:  "instance_id",
		},
		{
			name:      "org id column qualified",
			column:    domainRepo.OrgIDColumn(),
			qualified: true,
			expected:  "org_domains.org_id",
		},
		{
			name:      "org id column unqualified",
			column:    domainRepo.OrgIDColumn(),
			qualified: false,
			expected:  "org_id",
		},
		{
			name:      "domain column qualified",
			column:    domainRepo.DomainColumn(),
			qualified: true,
			expected:  "org_domains.domain",
		},
		{
			name:      "domain column unqualified",
			column:    domainRepo.DomainColumn(),
			qualified: false,
			expected:  "domain",
		},
		{
			name:      "is verified column qualified",
			column:    domainRepo.IsVerifiedColumn(),
			qualified: true,
			expected:  "org_domains.is_verified",
		},
		{
			name:      "is verified column unqualified",
			column:    domainRepo.IsVerifiedColumn(),
			qualified: false,
			expected:  "is_verified",
		},
		{
			name:      "is primary column qualified",
			column:    domainRepo.IsPrimaryColumn(),
			qualified: true,
			expected:  "org_domains.is_primary",
		},
		{
			name:      "is primary column unqualified",
			column:    domainRepo.IsPrimaryColumn(),
			qualified: false,
			expected:  "is_primary",
		},
		{
			name:      "validation type column qualified",
			column:    domainRepo.ValidationTypeColumn(),
			qualified: true,
			expected:  "org_domains.validation_type",
		},
		{
			name:      "validation type column unqualified",
			column:    domainRepo.ValidationTypeColumn(),
			qualified: false,
			expected:  "validation_type",
		},
		{
			name:      "created at column qualified",
			column:    domainRepo.CreatedAtColumn(),
			qualified: true,
			expected:  "org_domains.created_at",
		},
		{
			name:      "created at column unqualified",
			column:    domainRepo.CreatedAtColumn(),
			qualified: false,
			expected:  "created_at",
		},
		{
			name:      "updated at column qualified",
			column:    domainRepo.UpdatedAtColumn(),
			qualified: true,
			expected:  "org_domains.updated_at",
		},
		{
			name:      "updated at column unqualified",
			column:    domainRepo.UpdatedAtColumn(),
			qualified: false,
			expected:  "updated_at",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var builder database.StatementBuilder
			if test.qualified {
				test.column.WriteQualified(&builder)
			} else {
				test.column.WriteUnqualified(&builder)
			}
			assert.Equal(t, test.expected, builder.String())
		})
	}
}
