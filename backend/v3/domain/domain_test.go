package domain_test

import (
	"os"
	"testing"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/cache/cachemock"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		domain.SetOrgRepository(repository.OrganizationRepository)
		domain.SetCache(cachemock.NewOrganizationCacheMock())
		return m.Run()
	}())
}
