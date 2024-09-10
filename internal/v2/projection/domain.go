package projection

import (
	"github.com/zitadel/zitadel/internal/v2/domain"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
)

type Domain struct {
	id string

	Name      string
	Verified  bool
	IsPrimary bool
}

func (d *Domain) shouldReduce(event *v2_es.StorageEvent, domain string) bool {
	return event.Aggregate.ID == d.id && d.Name == domain
}

func (d *Domain) reduceAddedPayload(payload domain.AddedPayload) {
	d.Name = payload.Name
}

func (d *Domain) reduceVerifiedPayload(payload domain.VerifiedPayload) {
	d.Verified = true
}

func (d *Domain) reducePrimarySetPayload(payload domain.PrimarySetPayload) {
	d.IsPrimary = payload.Name == d.Name
}
