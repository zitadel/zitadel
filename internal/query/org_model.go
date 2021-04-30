package query

// import (
// 	"time"

// 	"github.com/caos/zitadel/internal/domain"
// 	"github.com/caos/zitadel/internal/eventstore"
// 	"github.com/caos/zitadel/internal/repository/org"
// )

// type OrgReadModel struct {
// 	eventstore.ReadModel

// 	stmts []eventstore.Statement
// }

// func NewOrgReadModel(orgID string) *OrgReadModel {
// 	return &OrgReadModel{
// 		ReadModel: eventstore.ReadModel{
// 			AggregateID:   orgID,
// 			ResourceOwner: orgID,
// 		},
// 	}
// }

// func (rm *OrgReadModel) Reduce() (err error) {
// 	for _, event := range rm.Events {
// 		switch e := event.(type) {
// 		case *org.OrgAddedEvent:
// 			rm.stmts = append(rm.stmts, &eventstore.CreateStatement{
// 				Values: []eventstore.Column{
// 					idColumn(e.Aggregate().ID),
// 					creationDateColumn(e.CreationDate()),
// 					changeDateColumn(e.CreationDate()),
// 					resourceOwnerColumn(e.Aggregate().ResourceOwner),
// 					stateColumn(domain.OrgStateActive),
// 					sequenceColumn(e.Sequence()),
// 					nameColumn(e.Name),
// 				},
// 			})
// 		case *org.OrgDeactivatedEvent:
// 			rm.stmts = append(rm.stmts, &eventstore.UpdateStatement{
// 				PK: []eventstore.Column{
// 					idColumn(e.Aggregate().ID),
// 				},
// 				Values: []eventstore.Column{
// 					changeDateColumn(e.CreationDate()),
// 					stateColumn(domain.OrgStateInactive),
// 				},
// 			})
// 		case *org.OrgReactivatedEvent:
// 			rm.stmts = append(rm.stmts, &eventstore.UpdateStatement{
// 				PK: []eventstore.Column{
// 					idColumn(e.Aggregate().ID),
// 				},
// 				Values: []eventstore.Column{
// 					changeDateColumn(e.CreationDate()),
// 					stateColumn(domain.OrgStateActive),
// 				},
// 			})
// 		case *org.OrgChangedEvent:
// 			rm.stmts = append(rm.stmts, &eventstore.UpdateStatement{
// 				PK: []eventstore.Column{
// 					idColumn(e.Aggregate().ID),
// 				},
// 				Values: []eventstore.Column{
// 					changeDateColumn(e.CreationDate()),
// 					nameColumn(e.Name),
// 				},
// 			})
// 		case *org.DomainPrimarySetEvent:
// 			rm.stmts = append(rm.stmts, &eventstore.UpdateStatement{
// 				PK: []eventstore.Column{
// 					idColumn(e.Aggregate().ID),
// 				},
// 				Values: []eventstore.Column{
// 					changeDateColumn(e.CreationDate()),
// 					domainColumn(e.Domain),
// 				},
// 			})
// 		}
// 	}
// 	return nil
// }

// func (rm *OrgReadModel) Statements() []eventstore.Statement {
// 	return rm.stmts
// }

// func idColumn(id string) eventstore.Column {
// 	return eventstore.Column{
// 		Name:  "id",
// 		Value: id,
// 	}
// }
// func creationDateColumn(creationDate time.Time) eventstore.Column {
// 	return eventstore.Column{
// 		Name:  "creation_date",
// 		Value: creationDate,
// 	}
// }
// func changeDateColumn(changeDate time.Time) eventstore.Column {
// 	return eventstore.Column{
// 		Name:  "change_date",
// 		Value: changeDate,
// 	}
// }
// func resourceOwnerColumn(resourceOwner string) eventstore.Column {
// 	return eventstore.Column{
// 		Name:  "resource_owner",
// 		Value: resourceOwner,
// 	}
// }
// func stateColumn(state domain.OrgState) eventstore.Column {
// 	return eventstore.Column{
// 		Name:  "org_state",
// 		Value: state,
// 	}
// }
// func sequenceColumn(sequence uint64) eventstore.Column {
// 	return eventstore.Column{
// 		Name:  "sequence",
// 		Value: sequence,
// 	}
// }
// func domainColumn(domain string) eventstore.Column {
// 	return eventstore.Column{
// 		Name:  "domain",
// 		Value: domain,
// 	}
// }
// func nameColumn(name string) eventstore.Column {
// 	return eventstore.Column{
// 		Name:  "name",
// 		Value: name,
// 	}
// }
