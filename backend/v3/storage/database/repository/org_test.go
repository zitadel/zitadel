package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBla is an example and can be removed later
func TestBla(t *testing.T) {
	var count int
	err := pool.QueryRow(context.Background(), "select count(*) from zitadel.instances").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
