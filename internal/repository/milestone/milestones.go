//go:generate stringer -type=Milestone

package milestone

type Milestone int

const (
	unknown Milestone = iota
	InstanceCreated
	AuthenticationSucceededOnInstance
	ProjectCreated
	ApplicationCreated
	AuthenticationSucceededOnApplication
	InstanceDeleted

	milestonesCount
)

func All() []Milestone {
	milestones := make([]Milestone, milestonesCount-1)
	for i := 1; i < int(milestonesCount); i++ {
		milestones[i] = Milestone(i)
	}
	return milestones
}
