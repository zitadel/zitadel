package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestAlreadyExists(t *testing.T) {
	parentErr := errors.New("parent error")
	id := "test_id"
	message := "test message"

	t.Run("ThrowAlreadyExists", func(t *testing.T) {
		err := zerrors.ThrowAlreadyExists(parentErr, id, message)
		assert.NotNil(t, err)

		zitadelErr, ok := zerrors.AsZitadelError(err)
		assert.True(t, ok)
		assert.Equal(t, zerrors.KindAlreadyExists, zitadelErr.Kind)

		zitadelError := new(zerrors.ZitadelError)
		if errors.As(err, &zitadelError) {
			assert.Equal(t, parentErr, zitadelError.Unwrap())
			assert.Equal(t, id, zitadelError.ID)
			assert.Equal(t, message, zitadelError.Message)
		} else {
			t.Errorf("error is not of type ZitadelError")
		}
	})

	t.Run("ThrowAlreadyExistsf", func(t *testing.T) {
		format := "formatted %s"
		arg := "message"
		expectedMessage := "formatted message"

		err := zerrors.ThrowAlreadyExistsf(parentErr, id, format, arg)
		assert.NotNil(t, err)

		zitadelErr, ok := zerrors.AsZitadelError(err)
		assert.True(t, ok)
		assert.Equal(t, zerrors.KindAlreadyExists, zitadelErr.Kind)

		zitadelError := new(zerrors.ZitadelError)
		if errors.As(err, &zitadelError) {
			assert.Equal(t, parentErr, zitadelError.Unwrap())
			assert.Equal(t, id, zitadelError.ID)
			assert.Equal(t, expectedMessage, zitadelError.Message)
		} else {
			t.Errorf("error is not of type ZitadelError")
		}
	})

	t.Run("IsErrorAlreadyExists", func(t *testing.T) {
		err := zerrors.ThrowAlreadyExists(parentErr, id, message)
		isAlreadyExists := zerrors.IsErrorAlreadyExists(err)
		assert.True(t, isAlreadyExists)

		otherErr := errors.New("some other error")
		isAlreadyExists = zerrors.IsErrorAlreadyExists(otherErr)
		assert.False(t, isAlreadyExists)
	})
}
