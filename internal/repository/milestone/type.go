//go:generate enumer -type Type -json

package milestone

type Type int

const (
	InstanceCreated Type = iota + 1
	AuthenticationSucceededOnInstance
	ProjectCreated
	ApplicationCreated
	AuthenticationSucceededOnApplication
	InstanceDeleted
)
