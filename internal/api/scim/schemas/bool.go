package schemas

import (
	"encoding/json"
	"strings"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/zerrors"
)

// RelaxedBool a bool which is more relaxed when it comes to json (un)marshaling.
// This ensures compatibility with some bugged scim providers,
// such as Microsoft Entry, which sends booleans as "True" or "False".
// See also https://learn.microsoft.com/en-us/entra/identity/app-provisioning/application-provisioning-config-problem-scim-compatibility.
type RelaxedBool bool

func NewRelaxedBool(value bool) *RelaxedBool {
	return gu.Ptr(RelaxedBool(value))
}

func (b *RelaxedBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(*b))
}

func (b *RelaxedBool) UnmarshalJSON(bytes []byte) error {
	str := strings.ToLower(string(bytes))
	switch {
	case str == "true" || str == "\"true\"":
		*b = true
	case str == "false" || str == "\"false\"":
		*b = false
	default:
		return zerrors.ThrowInvalidArgumentf(nil, "SCIM-BOO1", "bool expected, got %v", str)
	}
	return nil
}

func (b *RelaxedBool) Bool() bool {
	return bool(*b)
}
