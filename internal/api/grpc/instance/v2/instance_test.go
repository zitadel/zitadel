package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestValidateParam(t *testing.T) {
	tt := []struct {
		param       string
		paramName   string
		expectedErr error
	}{
		{"", "instance_id", zerrors.ThrowInvalidArgument(nil, "instance_id", "instance_id must not be empty")},
		{" ", "instance_id", zerrors.ThrowInvalidArgument(nil, "instance_id", "instance_id must not be empty")},
		{"valid_id", "instance_id", nil},
	}

	for _, tc := range tt {
		t.Run(tc.param, func(t *testing.T) {
			err := validateParam(tc.param, tc.paramName)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
