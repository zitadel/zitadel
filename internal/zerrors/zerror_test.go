package zerrors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestErrorMethod(t *testing.T) {
	err := zerrors.ThrowError(nil, "id", "msg")
	expected := "ID=id Message=msg"
	assert.Equal(t, expected, err.Error())

	err = zerrors.ThrowError(err, "subID", "subMsg")
	subExptected := "ID=subID Message=subMsg Parent=(ID=id Message=msg)"
	assert.Equal(t, subExptected, err.Error())
}
