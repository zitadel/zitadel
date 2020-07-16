package eventstore

import svcacc_event "github.com/caos/zitadel/internal/service_account/repository/eventsourcing"

type ServiceAccountRepo struct {
	ServiceAccountEvents *svcacc_event.ServiceAccountEventstore
}
