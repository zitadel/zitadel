package query

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgReadModel struct {
	eventstore.ReadModel

	Name          string
	State         domain.OrgState
	PrimaryDomain string
}

func NewOrgReadModel(orgID string) *OrgReadModel {
	return &OrgReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
	}
}

func (rm *OrgReadModel) Reduce() (stmt []string, err error) {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *org.OrgAddedEvent:
			//create stmt
			rm.Name = e.Name
			rm.State = domain.OrgStateActive
		case *org.OrgDeactivatedEvent:
			//update stmt
			rm.State = domain.OrgStateInactive
		case *org.OrgReactivatedEvent:
			//update stmt
			rm.State = domain.OrgStateActive
		case *org.OrgChangedEvent:
			//update stmt
			rm.Name = e.Name
		case *org.DomainPrimarySetEvent:
			//update stmt
			rm.PrimaryDomain = e.Domain
		}
	}
	return nil, nil
}
