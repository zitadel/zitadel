package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestErrorMethod(t *testing.T) {
	err := caos_errs.ThrowError(nil, "id", "msg")
	expected := "ID=id Message=msg"
	assert.Equal(t, expected, err.Error())

	err = caos_errs.ThrowError(err, "subID", "subMsg")
	subExptected := "ID=subID Message=subMsg Parent=(ID=id Message=msg)"
	assert.Equal(t, subExptected, err.Error())
}
