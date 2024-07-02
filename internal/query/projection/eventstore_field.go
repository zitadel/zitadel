package projection

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type fillFields struct {
	eventstore handler.EventStore
}

type fillProjectGrantFields struct{}

func (*fillProjectGrantFields) Name() string {
	return "project_grant_fields"
}

func newFillProjectGrantFields(config handler.Config) *handler.FieldHandler {
	return handler.NewFieldHandler(
		&config,
		"project_grant_fields",
		map[eventstore.AggregateType][]eventstore.EventType{
			org.AggregateType:     nil,
			project.AggregateType: nil,
		},
	)
}
