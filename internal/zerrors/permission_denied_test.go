package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestPermissionDenied(t *testing.T) {
	parentErr := errors.New("parent error")
	id := "PERMISSION_DENIED"
	message := "you do not have permission"

	t.Run("ThrowPermissionDenied", func(t *testing.T) {
		err := zerrors.ThrowPermissionDenied(parentErr, id, message)
		assert.NotNil(t, err)

		zitadelErr, ok := zerrors.AsZitadelError(err)
		assert.True(t, ok)
		assert.Equal(t, zerrors.KindPermissionDenied, zitadelErr.Kind)

		zitadelError := new(zerrors.ZitadelError)
		if errors.As(err, &zitadelError) {
			assert.Equal(t, parentErr, zitadelError.Unwrap())
			assert.Equal(t, id, zitadelError.ID)
			assert.Equal(t, message, zitadelError.Message)
		} else {
			t.Errorf("error is not of type ZitadelError")
		}
	})

	t.Run("ThrowPermissionDeniedf", func(t *testing.T) {
		format := "formatted %s"
		arg := "message"
		expectedMessage := "formatted message"

		err := zerrors.ThrowPermissionDeniedf(parentErr, id, format, arg)
		assert.NotNil(t, err)

		zitadelErr, ok := zerrors.AsZitadelError(err)
		assert.True(t, ok)
		assert.Equal(t, zerrors.KindPermissionDenied, zitadelErr.Kind)

		zitadelError := new(zerrors.ZitadelError)
		if errors.As(err, &zitadelError) {
			assert.Equal(t, parentErr, zitadelError.Unwrap())
			assert.Equal(t, id, zitadelError.ID)
			assert.Equal(t, expectedMessage, zitadelError.Message)
		} else {
			t.Errorf("error is not of type ZitadelError")
		}
	})

	t.Run("IsPermissionDenied", func(t *testing.T) {
		err := zerrors.ThrowPermissionDenied(parentErr, id, message)
		isPermissionDenied := zerrors.IsPermissionDenied(err)
		assert.True(t, isPermissionDenied)
	})
}
