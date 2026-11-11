package otp

import (
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Duration wraps time.Duration so JSON bodies can specify values as either
// duration strings ("5m", "30s") or integer nanoseconds.
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		parsed, perr := time.ParseDuration(s)
		if perr != nil {
			return perr
		}
		*d = Duration(parsed)
		return nil
	}
	var n int64
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	*d = Duration(n)
	return nil
}

// PreOTPSMSCodeContext is the JSON payload delivered to a preotpsmscode action target
// before an OTP SMS code is generated.
type PreOTPSMSCodeContext struct {
	FunctionName         string                 `json:"function,omitempty"`
	User                 *query.User            `json:"user,omitempty"`
	Org                  *query.UserInfoOrg     `json:"org,omitempty"`
	RecipientPhoneNumber string                 `json:"recipient_phone_number,omitempty"`
	GeneratorConfig      *PublicGeneratorConfig `json:"generator_config,omitempty"`
	Response             *PreOTPSMSCodeResponse `json:"response,omitempty"`
}

// PreOTPEmailCodeContext is the JSON payload delivered to a preotpemailcode action target
// before an OTP Email code is generated.
type PreOTPEmailCodeContext struct {
	FunctionName          string                   `json:"function,omitempty"`
	User                  *query.User              `json:"user,omitempty"`
	Org                   *query.UserInfoOrg       `json:"org,omitempty"`
	RecipientEmailAddress string                   `json:"recipient_email_address,omitempty"`
	GeneratorConfig       *PublicGeneratorConfig   `json:"generator_config,omitempty"`
	Response              *PreOTPEmailCodeResponse `json:"response,omitempty"`
}

// PublicGeneratorConfig is a JSON-safe snapshot of the effective OTP secret generator
// configuration passed to the action target for context.
type PublicGeneratorConfig struct {
	Length              uint32   `json:"length,omitempty"`
	Expiry              Duration `json:"expiry,omitempty"`
	IncludeLowerLetters bool     `json:"include_lower_letters,omitempty"`
	IncludeUpperLetters bool     `json:"include_upper_letters,omitempty"`
	IncludeDigits       bool     `json:"include_digits,omitempty"`
	IncludeSymbols      bool     `json:"include_symbols,omitempty"`
}

// GenerationOverrides lets an action target override individual OTP generator
// parameters for a single request. A nil field means "inherit from the instance
// default". Expiry is intentionally not part of this struct; it is controlled
// by PreOTP*CodeResponse.Expiry because it applies whether the code is generated
// or supplied by the action.
type GenerationOverrides struct {
	Length              *uint32 `json:"length,omitempty"`
	IncludeLowerLetters *bool   `json:"include_lower_letters,omitempty"`
	IncludeUpperLetters *bool   `json:"include_upper_letters,omitempty"`
	IncludeDigits       *bool   `json:"include_digits,omitempty"`
	IncludeSymbols      *bool   `json:"include_symbols,omitempty"`
}

// PreOTPSMSCodeResponse is the JSON body returned by a preotpsmscode action target.
// Code and Generate are mutually exclusive: setting both is rejected. If neither
// is set the instance defaults are used.
type PreOTPSMSCodeResponse struct {
	Expiry   *Duration            `json:"expiry,omitempty"`
	Code     *string              `json:"code,omitempty"`
	Generate *GenerationOverrides `json:"generate,omitempty"`
}

// PreOTPEmailCodeResponse mirrors PreOTPSMSCodeResponse for the OTP Email channel.
type PreOTPEmailCodeResponse struct {
	Expiry   *Duration            `json:"expiry,omitempty"`
	Code     *string              `json:"code,omitempty"`
	Generate *GenerationOverrides `json:"generate,omitempty"`
}

func (c *PreOTPSMSCodeContext) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *PreOTPSMSCodeContext) SetHTTPResponseBody(resp []byte) error {
	if !json.Valid(resp) {
		return zerrors.ThrowPreconditionFailed(nil, "ACTION-p7q2w", "Errors.Execution.ResponseIsNotValidJSON")
	}
	if c.Response == nil {
		c.Response = &PreOTPSMSCodeResponse{}
	}
	if err := json.Unmarshal(resp, c.Response); err != nil {
		return err
	}
	// Code and Generate describe different things (a finished value vs generation parameters);
	// accepting both would require an unspecified precedence rule — reject instead.
	if c.Response.Code != nil && c.Response.Generate != nil {
		c.Response = &PreOTPSMSCodeResponse{}
		return zerrors.ThrowPreconditionFailed(nil, "ACTION-k3j9z", "Errors.Execution.Invalid")
	}
	return nil
}

func (c *PreOTPSMSCodeContext) GetContent() any {
	return c.Response
}

func (c *PreOTPEmailCodeContext) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *PreOTPEmailCodeContext) SetHTTPResponseBody(resp []byte) error {
	if !json.Valid(resp) {
		return zerrors.ThrowPreconditionFailed(nil, "ACTION-m8x4n", "Errors.Execution.ResponseIsNotValidJSON")
	}
	if c.Response == nil {
		c.Response = &PreOTPEmailCodeResponse{}
	}
	if err := json.Unmarshal(resp, c.Response); err != nil {
		return err
	}
	if c.Response.Code != nil && c.Response.Generate != nil {
		c.Response = &PreOTPEmailCodeResponse{}
		return zerrors.ThrowPreconditionFailed(nil, "ACTION-r5t8b", "Errors.Execution.Invalid")
	}
	return nil
}

func (c *PreOTPEmailCodeContext) GetContent() any {
	return c.Response
}
