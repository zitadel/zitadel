package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cleanStaticQueries(t *testing.T) {
	query := `select
	foo,
	bar
from table;`
	want := "select foo, bar from table;"
	cleanStaticQueries(&query)
	assert.Equal(t, want, query)
}
