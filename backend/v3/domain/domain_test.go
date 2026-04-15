package domain_test

import (
	"os"
	"testing"

	"github.com/zitadel/passwap"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/crypto"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		domain.SetPasswordHasher(&crypto.Hasher{Swapper: &passwap.Swapper{}})
		return m.Run()
	}())
}
