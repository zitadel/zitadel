package chore_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chore Suite")
}
