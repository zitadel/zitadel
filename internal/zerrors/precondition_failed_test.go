package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestPreconditionFailed(t *testing.T) {
	parentErr := errors.New("parent error")
	id := "test_id"
	message := "test message"

	t.Run("ThrowPreconditionFailed", func(t *testing.T) {
		err := zerrors.ThrowPreconditionFailed(parentErr, id, message)
		assert.NotNil(t, err)

		zitadelErr, ok := zerrors.AsZitadelError(err)
		assert.True(t, ok)
		assert.Equal(t, zerrors.KindPreconditionFailed, zitadelErr.Kind)

		zitadelError := new(zerrors.ZitadelError)
		if errors.As(err, &zitadelError) {
			assert.Equal(t, parentErr, zitadelError.Unwrap())
			assert.Equal(t, id, zitadelError.ID)
			assert.Equal(t, message, zitadelError.Message)
		} else {
			t.Errorf("error is not of type ZitadelError")
		}
	})

	t.Run("ThrowPreconditionFailedf", func(t *testing.T) {
		format := "formatted %s"
		arg := "message"
		expectedMessage := "formatted message"

		err := zerrors.ThrowPreconditionFailedf(parentErr, id, format, arg)
		assert.NotNil(t, err)

		zitadelErr, ok := zerrors.AsZitadelError(err)
		assert.True(t, ok)
		assert.Equal(t, zerrors.KindPreconditionFailed, zitadelErr.Kind)

		zitadelError := new(zerrors.ZitadelError)
		if errors.As(err, &zitadelError) {
			assert.Equal(t, parentErr, zitadelError.Unwrap())
			assert.Equal(t, id, zitadelError.ID)
			assert.Equal(t, expectedMessage, zitadelError.Message)
		} else {
			t.Errorf("error is not of type ZitadelError")
		}
	})

	t.Run("IsPreconditionFailed", func(t *testing.T) {
		err := zerrors.ThrowPreconditionFailed(parentErr, id, message)
		isPreconditionFailed := zerrors.IsPreconditionFailed(err)
		assert.True(t, isPreconditionFailed)
	})
}
