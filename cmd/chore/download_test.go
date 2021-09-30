package chore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChore_Download(t *testing.T) {
	assert.NoError(t, downloadZitadelctl("./artifacts/zitadelctl-darwin-amd64", "operator-e2e"))
}
