package eventsourcing

import "github.com/caos/zitadel/internal/eventstore"

type OrgEventstore struct {
	eventstore.Eventstore
}

type OrgConfig struct {
	eventstore.Eventstore
}

func StartOrg(conf OrgConfig) (*OrgEventstore, error) {
	return &OrgEventstore{Eventstore: conf.Eventstore}, nil
}
