package projection

import "github.com/zitadel/zitadel/internal/v2/eventstore"

type projection struct {
	instance string
	position float64
	sequence uint32
}

func (p projection) shouldReduce(event eventstore.Event) bool {
	return (p.instance == "" || p.instance == event.Aggregate().Instance) && p.position <= event.Position()
}
