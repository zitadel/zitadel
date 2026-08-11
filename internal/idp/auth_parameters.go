package idp

import (
	"maps"
	"slices"
	"strings"
)

// reservedAuthorizationParameters are the parameters of an upstream authorization request
// which are fully managed by ZITADEL. They can neither be configured statically on a provider
// nor be forwarded from the original authorization request.
var reservedAuthorizationParameters = []string{
	"client_id",
	"client_secret",
	"code",
	"code_challenge",
	"code_challenge_method",
	"nonce",
	"redirect_uri",
	"response_type",
	"scope",
	"state",
}

// IsReservedAuthorizationParameter reports whether the parameter is managed by ZITADEL
// and therefore must not be set or forwarded by configuration.
func IsReservedAuthorizationParameter(key string) bool {
	return slices.Contains(reservedAuthorizationParameters, strings.ToLower(strings.TrimSpace(key)))
}

// ReservedAuthorizationParameters returns the parameters managed by ZITADEL.
func ReservedAuthorizationParameters() []string {
	return slices.Clone(reservedAuthorizationParameters)
}

// AuthorizationParameters allows to pass parameters of the original authorization request
// to BeginAuth. A provider only forwards those which it is explicitly configured to forward.
type AuthorizationParameters map[string]string

func (p AuthorizationParameters) setValue() {}

// ResolveAuthorizationParameters merges the parameters which are added to the upstream
// authorization request on top of the ones ZITADEL generates itself.
//
// Parameters of the original authorization request are only taken into account if their name is
// contained in forward. The statically configured parameters take precedence over the forwarded ones.
// Reserved parameters are dropped from both sources.
//
// A returned empty value means the parameter must not be sent at all, which allows to disable
// defaults such as `prompt=select_account`.
func ResolveAuthorizationParameters(static map[string]string, forward []string, params []Parameter) map[string]string {
	resolved := make(map[string]string, len(static)+len(forward))
	if len(forward) > 0 {
		for _, param := range params {
			requested, ok := param.(AuthorizationParameters)
			if !ok {
				continue
			}
			for key, value := range requested {
				key = normalizeAuthorizationParameter(key)
				if !slices.ContainsFunc(forward, func(allowed string) bool {
					return normalizeAuthorizationParameter(allowed) == key
				}) {
					continue
				}
				resolved[key] = value
			}
		}
	}
	for key, value := range static {
		resolved[normalizeAuthorizationParameter(key)] = value
	}
	maps.DeleteFunc(resolved, func(key, _ string) bool {
		return key == "" || IsReservedAuthorizationParameter(key)
	})
	return resolved
}

func normalizeAuthorizationParameter(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}
