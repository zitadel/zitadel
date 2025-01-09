package scim

import (
	"errors"
	"strconv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AssertedScimError struct {
	Error *ScimError
}

func RequireScimError(t require.TestingT, httpStatus int, err error) AssertedScimError {
	require.Error(t, err)

	var scimErr *ScimError
	assert.True(t, errors.As(err, &scimErr))
	assert.Equal(t, strconv.Itoa(httpStatus), scimErr.Status)
	return AssertedScimError{scimErr} // wrap it, otherwise error handling is enforced
}
