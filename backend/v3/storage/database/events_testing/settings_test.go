//go:build integration

package events_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
)

func TestServer_TestLoginSettingsReduces(t *testing.T) {
	// instanceID := Instance.ID()

	// orgID := Instance.DefaultOrg.Id

	settingsRepo := repository.SettingsRepository(pool)

	t.Run("test adding login settings reduces", func(t *testing.T) {
		// beforeCreate := time.Now()
		newInstance := integration.NewInstance(t.Context())

		setting, err := settingsRepo.List(
			t.Context(), settingsRepo.InstanceIDCondition(newInstance.ID()),
			settingsRepo.TypeCondition(domain.SettingTypeLogin))
		// afterCreate := time.Now()
		require.NoError(t, err)

		fmt.Printf("[DEBUGPRINT] [:1] setting = %+v\n", setting[0])
	})
}
