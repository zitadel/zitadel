package chore_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChore(t *testing.T) {
	os.Setenv("ZITADEL_E2E_TAG", "operator-e2e")
	os.Setenv("ZITADEL_E2E_SHA", "SHA-40aed4b9f11a")
	os.Setenv("ZITADEL_E2E_DBUSER", "test")
	os.Setenv("ZITADEL_E2E_BACKUPSAJSON", "./artifacts/sa.json")
	os.Setenv("ZITADEL_E2E_BACKUPAKID", "./artifacts/akid")
	os.Setenv("ZITADEL_E2E_BACKUPSAK", "./artifacts/sak")

	RegisterFailHandler(Fail)
	RunSpecs(t, "Chore Suite")
}
