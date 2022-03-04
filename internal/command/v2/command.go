package command

import "github.com/caos/zitadel/internal/eventstore"

type Command struct {
	es        *eventstore.Eventstore
	iamDomain string
}

func New(es *eventstore.Eventstore, iamDomain string) *Command {
	return &Command{
		es:        es,
		iamDomain: iamDomain,
	}
}
