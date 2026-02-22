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

func TestDomainAddedEvent(t *testing.T) {
	ctx := t.Context()
	instanceDomainRepo := repository.InstanceDomainRepository()
	instanceRepo := repository.InstanceRepository()
	beforeAdd := time.Now()
	err := instanceRepo.Create(ctx, pool, &domain.Instance{
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
	tstCmd := instance.NewDomainAddedEvent(ctx, &instance.NewAggregate("123").Aggregate, "test-domain.com", false)
	_, err = es.Push(ctx, tstCmd)
	require.NoError(t, err)
	retryDuration, tick := waitForAndTickWithMaxDuration(ctx, time.Second*30)
	afterAdd := time.Now()
	assert.EventuallyWithT(t, func(t *assert.CollectT) {
		domain, err := instanceDomainRepo.Get(ctx, pool,
			database.WithCondition(
				database.And(
					instanceDomainRepo.InstanceIDCondition("123"),
					instanceDomainRepo.DomainCondition(database.TextOperationEqual, "test-domain.com"),
					instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
				),
			),
		)
		require.NoError(t, err)
		assert.Equal(t, "test-domain.com", domain.Domain)
		assert.Equal(t, "123", domain.InstanceID)
		assert.False(t, *domain.IsPrimary)
		assert.WithinRange(t, domain.CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, domain.UpdatedAt, beforeAdd, afterAdd)
	}, retryDuration, tick)
}
