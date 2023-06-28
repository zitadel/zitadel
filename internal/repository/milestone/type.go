//go:generate stringer -type Type

package milestone

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
