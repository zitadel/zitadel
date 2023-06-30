//go:generate stringer -type Type

package milestone

import (
	"fmt"
	"strings"
)

type Type int

const (
	unknown Type = iota

	InstanceCreated
	AuthenticationSucceededOnInstance
	ProjectCreated
	ApplicationCreated
	AuthenticationSucceededOnApplication
	InstanceDeleted

	typesCount
)

func AllTypes() []Type {
	types := make([]Type, typesCount-1)
	for i := Type(1); i < typesCount; i++ {
		types[i-1] = i
	}
	return types
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

func (t *Type) UnmarshalJSON(data []byte) error {
	*t = typeFromString(strings.Trim(string(data), `"`))
	return nil
}

func typeFromString(t string) Type {
	switch t {
	case InstanceCreated.String():
		return InstanceCreated
	case AuthenticationSucceededOnInstance.String():
		return AuthenticationSucceededOnInstance
	case ProjectCreated.String():
		return ProjectCreated
	case ApplicationCreated.String():
		return ApplicationCreated
	case AuthenticationSucceededOnApplication.String():
		return AuthenticationSucceededOnApplication
	case InstanceDeleted.String():
		return InstanceDeleted
	default:
		return unknown
	}
}
