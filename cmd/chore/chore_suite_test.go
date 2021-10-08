package chore_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChore(t *testing.T) {
	os.Setenv("ZITADEL_E2E_TAG", "operator-e2e")
	os.Setenv("ZITADEL_E2E_DBUSER", "test")

	RegisterFailHandler(Fail)
	RunSpecs(t, "Chore Suite")
}
