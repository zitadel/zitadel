package serrors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

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
}

type scimError struct {
	Schemas       []schemas.ScimSchemaType `json:"schemas"`
	ScimType      scimErrorType            `json:"scimType,omitempty"`
	Detail        string                   `json:"detail,omitempty"`
	StatusCode    int                      `json:"-"`
	Status        string                   `json:"status"`
	ZitadelDetail *errorDetail             `json:"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail,omitempty"`
}

type errorDetail struct {
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
)

var translator *i18n.Translator

func ErrorHandler(next zhttp_middleware.HandlerFuncWithError) http.Handler {
	var err error
	translator, err = i18n.NewZitadelTranslator(language.English)
	logging.OnError(err).Panic("unable to get translator")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err = next(w, r); err == nil {
			return
		}

		scimErr := mapToScimJsonError(r.Context(), err)
		w.WriteHeader(scimErr.StatusCode)

		jsonErr := json.NewEncoder(w).Encode(scimErr)
		logging.OnError(jsonErr).Warn("Failed to marshal scim error response")
	})
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

func (err *scimError) Error() string {
	return fmt.Sprintf("SCIM Error: %s: %s", err.ScimType, err.Detail)
}

func (err *wrappedScimError) Error() string {
	return fmt.Sprintf("SCIM Error: %s: %s", err.ScimType, err.Parent.Error())
}

func mapToScimJsonError(ctx context.Context, err error) *scimError {
	scimErr := new(wrappedScimError)
	if ok := errors.As(err, &scimErr); ok {
		mappedErr := mapToScimJsonError(ctx, scimErr.Parent)
		mappedErr.ScimType = scimErr.ScimType
		return mappedErr
	}

	zitadelErr := new(zerrors.ZitadelError)
	if ok := errors.As(err, &zitadelErr); !ok {
		return &scimError{
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
	return &scimError{
		Schemas:    []schemas.ScimSchemaType{schemas.IdError, schemas.IdZitadelErrorDetail},
		ScimType:   mapErrorToScimErrorType(err),
		Detail:     localizedMsg,
		StatusCode: statusCode,
		Status:     strconv.Itoa(statusCode),
		ZitadelDetail: &errorDetail{
			ID:      zitadelErr.GetID(),
			Message: zitadelErr.GetMessage(),
		},
	}
}

func mapErrorToScimErrorType(err error) scimErrorType {
	switch {
	case zerrors.IsErrorInvalidArgument(err):
		return ScimTypeInvalidValue
	default:
		return ""
	}
}
