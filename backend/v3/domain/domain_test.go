package domain_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		return m.Run()
	}())
}
