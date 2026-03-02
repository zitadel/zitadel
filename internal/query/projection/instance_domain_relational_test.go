package projection

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestInstanceDomainReducers(t *testing.T) {
	handler := &instanceDomainRelationalProjection{}
	rawTx, tx := getTransactions(t)

	t.Cleanup(func() {
		require.NoError(t, rawTx.Rollback())
	})
	ctx := t.Context()

	instanceDomainRepo := repository.InstanceDomainRepository()
	instanceRepo := repository.InstanceRepository()
	err := instanceRepo.Create(ctx, tx, &domain.Instance{
		ID:              "123",
		Name:            "my instance",
		DefaultOrgID:    gofakeit.UUID(),
		IAMProjectID:    gofakeit.UUID(),
		ConsoleClientID: gofakeit.UUID(),
		ConsoleAppID:    gofakeit.UUID(),
		DefaultLanguage: gofakeit.LanguageAbbreviation(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	t.Run("'custom domain added' reducer should create custom domain", func(t *testing.T) {
		// Given
		domainName := "test-domain.com"
		domainAddedEvt := instance.NewDomainAddedEvent(ctx, &instance.NewAggregate("123").Aggregate, domainName, false)

		// Test
		res := callReduce(t, ctx, rawTx, handler, domainAddedEvt)
		require.True(t, res)

		// Verify
		require.NoError(t, err)

		domain, err := instanceDomainRepo.Get(ctx, tx,
			database.WithCondition(
				database.And(
					instanceDomainRepo.InstanceIDCondition("123"),
					instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
				),
			),
		)
		require.NoError(t, err)
		assert.Equal(t, domainName, domain.Domain)
		assert.Equal(t, "123", domain.InstanceID)

		require.NotNil(t, domain.IsPrimary)
		assert.False(t, *domain.IsPrimary)

		assert.NotZero(t, domain.CreatedAt)
		assert.NotZero(t, domain.UpdatedAt)
	})

	t.Run("'primary domain set' reducer should update existing domain", func(t *testing.T) {
		// Given
		domainName := "test-domain-primary-set.com"
		err := instanceDomainRepo.Add(ctx, tx, &domain.AddInstanceDomain{
			InstanceID:  "123",
			Domain:      domainName,
			Type:        domain.DomainTypeCustom,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsPrimary:   new(bool),
			IsGenerated: new(bool),
		})
		require.NoError(t, err)

		domainPrimarySetEvt := instance.NewDomainPrimarySetEvent(ctx, &instance.NewAggregate("123").Aggregate, domainName)

		// Test
		res := callReduce(t, ctx, rawTx, handler, domainPrimarySetEvt)
		require.True(t, res)

		// Verify
		require.NoError(t, err)

		domain, err := instanceDomainRepo.Get(ctx, tx,
			database.WithCondition(
				database.And(
					instanceDomainRepo.InstanceIDCondition("123"),
					instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
					instanceDomainRepo.IsPrimaryCondition(true),
				),
			),
		)
		require.NoError(t, err)
		assert.Equal(t, domainName, domain.Domain)
		assert.Equal(t, "123", domain.InstanceID)

		require.NotNil(t, domain.IsPrimary)
		assert.True(t, *domain.IsPrimary)

		assert.NotZero(t, domain.CreatedAt)
		assert.NotZero(t, domain.UpdatedAt)
	})

	t.Run("'custom domain removed' reducer should remove custom domain", func(t *testing.T) {
		// Given
		domainName := "test-domain-removed.com"
		err := instanceDomainRepo.Add(ctx, tx, &domain.AddInstanceDomain{
			InstanceID:  "123",
			Domain:      domainName,
			Type:        domain.DomainTypeCustom,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsPrimary:   new(bool),
			IsGenerated: new(bool),
		})
		require.NoError(t, err)

		domainRemovedEvt := instance.NewDomainRemovedEvent(ctx, &instance.NewAggregate("123").Aggregate, domainName)

		// Test
		res := callReduce(t, ctx, rawTx, handler, domainRemovedEvt)
		require.True(t, res)

		// Verify
		require.NoError(t, err)

		_, err = instanceDomainRepo.Get(ctx, tx,
			database.WithCondition(
				database.And(
					instanceDomainRepo.InstanceIDCondition("123"),
					instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
				),
			),
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, &database.NoRowFoundError{})
	})

	t.Run("'trusted domain added' reducer should create trusted domain", func(t *testing.T) {
		// Given
		domainName := "trusted-test-domain.com"
		trustedDomainAddedEvt := instance.NewTrustedDomainAddedEvent(ctx, &instance.NewAggregate("123").Aggregate, domainName)

		// Test
		res := callReduce(t, ctx, rawTx, handler, trustedDomainAddedEvt)
		require.True(t, res)

		// Verify
		require.NoError(t, err)

		domain, err := instanceDomainRepo.Get(ctx, tx,
			database.WithCondition(
				database.And(
					instanceDomainRepo.InstanceIDCondition("123"),
					instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					instanceDomainRepo.TypeCondition(domain.DomainTypeTrusted),
				),
			),
		)
		require.NoError(t, err)
		assert.Equal(t, domainName, domain.Domain)
		assert.Equal(t, "123", domain.InstanceID)
		assert.NotZero(t, domain.CreatedAt)
		assert.NotZero(t, domain.UpdatedAt)
	})

	t.Run("'trusted domain removed' reducer should remove trusted domain", func(t *testing.T) {
		// Given
		domainName := "test-trusted-domain-removed.com"
		err := instanceDomainRepo.Add(ctx, tx, &domain.AddInstanceDomain{
			InstanceID: "123",
			Domain:     domainName,
			Type:       domain.DomainTypeTrusted,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
		require.NoError(t, err)

		trustedDomainRemovedEvt := instance.NewTrustedDomainRemovedEvent(ctx, &instance.NewAggregate("123").Aggregate, domainName)

		// Test
		res := callReduce(t, ctx, rawTx, handler, trustedDomainRemovedEvt)
		require.True(t, res)

		// Verify
		require.NoError(t, err)

		_, err = instanceDomainRepo.Get(ctx, tx,
			database.WithCondition(
				database.And(
					instanceDomainRepo.InstanceIDCondition("123"),
					instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					instanceDomainRepo.TypeCondition(domain.DomainTypeTrusted),
				),
			),
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, &database.NoRowFoundError{})
	})
}
