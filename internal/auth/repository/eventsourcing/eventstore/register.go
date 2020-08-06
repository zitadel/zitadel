package eventstore

import (
	auth_model "github.com/caos/zitadel/internal/auth/model"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

type Register struct {
	*model.Org
	*usr_es.User
}

func (r *Register) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		var err error
		switch event.AggregateType {
		case model.OrgAggregate:
			err = r.Org.AppendEvent(event)
		case usr_es.UserAggregate:
			err = r.User.AppendEvent(event)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterToModel(register *Register) *auth_model.RegisterOrg {
	return &auth_model.RegisterOrg{
		Org:  model.OrgToModel(register.Org),
		User: usr_es.UserToModel(register.User),
	}
}
