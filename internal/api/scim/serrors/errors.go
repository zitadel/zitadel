package serrors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zitadel/logging"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	zhttp_middleware "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type scimErrorType string

type wrappedScimError struct {
	Parent   error
	ScimType scimErrorType
	Status   int
}

type ScimError struct {
	Schemas       []schemas.ScimSchemaType `json:"schemas"`
	ScimType      scimErrorType            `json:"scimType,omitempty"`
	Detail        string                   `json:"detail,omitempty"`
	StatusCode    int                      `json:"-"`
	Status        string                   `json:"status"`
	ZitadelDetail *ErrorDetail             `json:"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail,omitempty"`
}

type ErrorDetail struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

const (
	// ScimTypeInvalidValue A required value was missing,
	// or the value specified was not compatible with the operation,
	// or attribute type (see Section 2.2 of RFC7643),
	// or resource schema (see Section 4 of RFC7643).
	ScimTypeInvalidValue scimErrorType = "invalidValue"

	// ScimTypeInvalidSyntax The request body message structure was invalid or did
	// not conform to the request schema.
	ScimTypeInvalidSyntax scimErrorType = "invalidSyntax"

	// ScimTypeInvalidFilter The specified filter syntax as invalid, or the
	// specified attribute and filter comparison combination is not supported.
	ScimTypeInvalidFilter scimErrorType = "invalidFilter"

	// ScimTypeInvalidPath The "path" attribute was invalid or malformed.
	ScimTypeInvalidPath scimErrorType = "invalidPath"

	// ScimTypeNoTarget The specified "path" did not
	// yield an attribute or attribute value that could be operated on.
	// This occurs when the specified "path" value contains a filter that yields no match.
	ScimTypeNoTarget scimErrorType = "noTarget"

	// ScimTypeUniqueness One or more of the attribute values are already in use or are reserved.
	ScimTypeUniqueness scimErrorType = "uniqueness"
)

func ErrorHandler(translator *i18n.Translator) func(next zhttp_middleware.HandlerFuncWithError) http.Handler {
	return func(next zhttp_middleware.HandlerFuncWithError) http.Handler {
		var err error

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err = next(w, r); err == nil {
				return
			}

			scimErr := MapToScimError(r.Context(), translator, err)
			w.WriteHeader(scimErr.StatusCode)

			jsonErr := json.NewEncoder(w).Encode(scimErr)
			logging.OnError(jsonErr).Warn("Failed to marshal scim error response")
		})
	}
}

func ThrowInvalidValue(parent error) error {
	return &wrappedScimError{
		Parent:   parent,
		ScimType: ScimTypeInvalidValue,
	}
}

func ThrowInvalidSyntax(parent error) error {
	return &wrappedScimError{
		Parent:   parent,
		ScimType: ScimTypeInvalidSyntax,
	}
}

func ThrowInvalidFilter(parent error) error {
	return &wrappedScimError{
		Parent:   parent,
		ScimType: ScimTypeInvalidFilter,
	}
}

func ThrowInvalidPath(parent error) error {
	return &wrappedScimError{
		Parent:   parent,
		ScimType: ScimTypeInvalidPath,
	}
}

func ThrowNoTarget(parent error) error {
	return &wrappedScimError{
		Parent:   parent,
		ScimType: ScimTypeNoTarget,
	}
}

func ThrowPayloadTooLarge(parent error) error {
	return &wrappedScimError{
		Parent: parent,
		Status: http.StatusRequestEntityTooLarge,
	}
}

func IsScimOrZitadelError(err error) bool {
	return IsScimError(err) || zerrors.IsZitadelError(err)
}

func IsScimError(err error) bool {
	var scimErr *wrappedScimError
	return errors.As(err, &scimErr)
}

func (err *ScimError) Error() string {
	return fmt.Sprintf("SCIM Error: %s: %s", err.ScimType, err.Detail)
}

func (err *wrappedScimError) Error() string {
	return fmt.Sprintf("SCIM Error: %s: %s", err.ScimType, err.Parent.Error())
}

func MapToScimError(ctx context.Context, translator *i18n.Translator, err error) *ScimError {
	scimError := new(ScimError)
	if ok := errors.As(err, &scimError); ok {
		return scimError
	}

	scimWrappedError := new(wrappedScimError)
	if ok := errors.As(err, &scimWrappedError); ok {
		mappedErr := MapToScimError(ctx, translator, scimWrappedError.Parent)
		if scimWrappedError.ScimType != "" {
			mappedErr.ScimType = scimWrappedError.ScimType
		}

		if scimWrappedError.Status != 0 {
			mappedErr.Status = strconv.Itoa(scimWrappedError.Status)
			mappedErr.StatusCode = scimWrappedError.Status
		}

		return mappedErr
	}

	zitadelErr := new(zerrors.ZitadelError)
	if ok := errors.As(err, &zitadelErr); !ok {
		return &ScimError{
			Schemas:    []schemas.ScimSchemaType{schemas.IdError},
			Detail:     "Unknown internal server error",
			Status:     strconv.Itoa(http.StatusInternalServerError),
			StatusCode: http.StatusInternalServerError,
		}
	}

	statusCode, ok := http_util.ZitadelErrorToHTTPStatusCode(err)
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	localizedMsg := translator.LocalizeFromCtx(ctx, zitadelErr.GetMessage(), nil)
	return &ScimError{
		Schemas:    []schemas.ScimSchemaType{schemas.IdError, schemas.IdZitadelErrorDetail},
		ScimType:   mapErrorToScimErrorType(err),
		Detail:     localizedMsg,
		StatusCode: statusCode,
		Status:     strconv.Itoa(statusCode),
		ZitadelDetail: &ErrorDetail{
			ID:      zitadelErr.GetID(),
			Message: zitadelErr.GetMessage(),
		},
	}
}

func mapErrorToScimErrorType(err error) scimErrorType {
	switch {
	case zerrors.IsErrorInvalidArgument(err):
		return ScimTypeInvalidValue
	case zerrors.IsErrorAlreadyExists(err):
		return ScimTypeUniqueness
	default:
		return ""
	}
}
