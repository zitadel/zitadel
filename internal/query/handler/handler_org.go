package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
)

const (
	idColName            = "id"
	creationDateColName  = "creation_date"
	changeDateColName    = "change_date"
	resourceOwnerColName = "resource_owner"
	stateColName         = "org_state"
	sequenceColName      = "sequence"
	domainColName        = "domain"
	nameColName          = "name"
)

type OrgHandler struct {
	handler.ReadModelHandler
	crdb.StatementHandler

	TableName string
}

func NewOrgHandler(
	ctx context.Context,
	es *eventstore.Eventstore,
	client *sql.DB,
) *OrgHandler {
	h := &OrgHandler{
		ReadModelHandler: *handler.NewReadModelHandler(es, 1*time.Minute),
		StatementHandler: crdb.NewStatementHandler(
			client,
			"read_models.orgs",
			"read_models.current_sequences",
			"read_models.locks",
		),
		TableName: "read_models.orgs",
	}
	go h.ReadModelHandler.Process(
		ctx,
		h.reduce,
		h.StatementHandler.Update,
		h.StatementHandler.Lock,
		h.StatementHandler.Unlock,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "orgs"),
	)

	h.ReadModelHandler.Handler.Subscribe("orgs")

	return h
}

func (h *OrgHandler) reduce(event eventstore.EventReader) ([]handler.Statement, error) {
	stmts := []handler.Statement{}

	switch e := event.(type) {
	case *org.OrgAddedEvent:
		stmts = append(stmts, h.orgAddedStmts(e)...)
	case *org.OrgChangedEvent:
		stmts = append(stmts, h.orgChangedStmts(e)...)
	case *org.OrgDeactivatedEvent:
		stmts = append(stmts, h.orgDeactivatedStmts(e)...)
	case *org.OrgReactivatedEvent:
		stmts = append(stmts, h.orgReactivatedStmts(e)...)
	case *org.DomainPrimarySetEvent:
		stmts = append(stmts, h.orgPrimaryDomainStmts(e)...)
	}

	return stmts, nil
}

func (h *OrgHandler) orgAddedStmts(event *org.OrgAddedEvent) []handler.Statement {
	return []handler.Statement{
		handler.NewCreateStatement(h.TableName, []handler.Column{
			{
				Name:  idColName,
				Value: event.Aggregate().ID,
			},
			{
				Name:  creationDateColName,
				Value: event.CreationDate(),
			},
			{
				Name:  changeDateColName,
				Value: event.CreationDate(),
			},
			{
				Name:  resourceOwnerColName,
				Value: event.Aggregate().ResourceOwner,
			},
			{
				Name:  sequenceColName,
				Value: event.Sequence(),
			},
			{
				Name:  nameColName,
				Value: event.Name,
			},
			{
				Name:  stateColName,
				Value: domain.OrgStateActive,
			},
		}, event.PreviousSequence()),
	}
}

func (h *OrgHandler) orgChangedStmts(event *org.OrgChangedEvent) []handler.Statement {
	values := []handler.Column{
		{
			Name:  changeDateColName,
			Value: event.CreationDate(),
		},
		{
			Name:  sequenceColName,
			Value: event.Sequence(),
		},
	}
	if event.Name != "" {
		values = append(values, handler.Column{
			Name:  nameColName,
			Value: event.Name,
		})
	}
	return []handler.Statement{
		handler.NewUpdateStatement(
			h.TableName,
			[]handler.Column{
				{
					Name:  idColName,
					Value: event.Aggregate().ID,
				},
			},
			values,
			event.PreviousSequence(),
		),
	}
}

func (h *OrgHandler) orgReactivatedStmts(event *org.OrgReactivatedEvent) []handler.Statement {
	return []handler.Statement{
		handler.NewUpdateStatement(
			h.TableName,
			[]handler.Column{
				{
					Name:  idColName,
					Value: event.Aggregate().ID,
				},
			},
			[]handler.Column{
				{
					Name:  changeDateColName,
					Value: event.CreationDate(),
				},
				{
					Name:  sequenceColName,
					Value: event.Sequence(),
				},
				{
					Name:  stateColName,
					Value: domain.OrgStateActive,
				},
			},
			event.PreviousSequence(),
		),
	}
}

func (h *OrgHandler) orgDeactivatedStmts(event *org.OrgDeactivatedEvent) []handler.Statement {
	return []handler.Statement{
		handler.NewUpdateStatement(
			h.TableName,
			[]handler.Column{
				{
					Name:  idColName,
					Value: event.Aggregate().ID,
				},
			},
			[]handler.Column{
				{
					Name:  changeDateColName,
					Value: event.CreationDate(),
				},
				{
					Name:  sequenceColName,
					Value: event.Sequence(),
				},
				{
					Name:  stateColName,
					Value: domain.OrgStateInactive,
				},
			},
			event.PreviousSequence(),
		),
	}
}

func (h *OrgHandler) orgPrimaryDomainStmts(event *org.DomainPrimarySetEvent) []handler.Statement {
	return []handler.Statement{
		handler.NewUpdateStatement(
			h.TableName,
			[]handler.Column{
				{
					Name:  idColName,
					Value: event.Aggregate().ID,
				},
			},
			[]handler.Column{
				{
					Name:  changeDateColName,
					Value: event.CreationDate(),
				},
				{
					Name:  sequenceColName,
					Value: event.Sequence(),
				},
				{
					Name:  nameColName,
					Value: event.Domain,
				},
			},
			event.PreviousSequence(),
		),
	}
}
