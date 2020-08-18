package eventstore

import (
	admin_model "github.com/caos/zitadel/internal/admin/model"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

type Setup struct {
	*model.Org
	*usr_es.Human
}

func (s *Setup) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		var err error
		switch event.AggregateType {
		case model.OrgAggregate:
			err = s.Org.AppendEvent(event)
		case usr_es.UserAggregate:
			err = s.User.AppendEvent(event)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func SetupToModel(setup *Setup) *admin_model.SetupOrg {
	return &admin_model.SetupOrg{
		Org:  model.OrgToModel(setup.Org),
		User: usr_es.HumanToModel(setup.User),
	}
}
