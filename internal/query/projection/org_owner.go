package projection

// import (
// 	"context"
// 	"time"

// 	"github.com/caos/logging"
// 	"github.com/caos/zitadel/internal/domain"
// 	"github.com/caos/zitadel/internal/errors"
// 	"github.com/caos/zitadel/internal/eventstore"
// 	"github.com/caos/zitadel/internal/eventstore/handler"
// 	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
// 	"github.com/caos/zitadel/internal/repository/org"
// 	"github.com/caos/zitadel/internal/repository/user"
// 	"golang.org/x/text/language"
// )

// type OrgOwner struct {
// 	OrgID             string        `col:"org_id"`
// 	OrgName           string        `col:"org_name"`
// 	OrgCreationDate   time.Time     `col:"org_creation_date"`
// 	OwnerID           string        `col:"owner_id"`
// 	OwnerLanguage     *language.Tag `col:"owner_language"`
// 	OwnerEmailAddress string        `col:"owner_email"`
// 	OwnerFirstName    string        `col:"owner_first_name"`
// 	OwnerLastName     string        `col:"owner_last_name"`
// 	OwnerGender       domain.Gender `col:"owner_gender"`
// }

// type OrgOwnerProjection struct {
// 	crdb.StatementHandler
// }

// func NewOrgOwnerProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgOwnerProjection {
// 	p := &OrgOwnerProjection{}
// 	config.ProjectionName = "projections.org_owners"
// 	config.Reducers = p.reducers()
// 	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
// 	return p
// }

// func (p *OrgOwnerProjection) reducers() []handler.AggregateReducer {
// 	return []handler.AggregateReducer{
// 		{
// 			Aggregate: org.AggregateType,
// 			EventRedusers: []handler.EventReducer{
// 				{
// 					Event:  org.MemberAddedEventType,
// 					Reduce: p.reduceMemberAdded,
// 				},
// 				{
// 					Event:  org.MemberChangedEventType,
// 					Reduce: p.reduceMemberChanged,
// 				},
// 				{
// 					Event:  org.MemberRemovedEventType,
// 					Reduce: p.reduceMemberRemoved,
// 				},
// 				{
// 					Event:  org.OrgChangedEventType,
// 					Reduce: p.reduceOrgChanged,
// 				},
// 				{
// 					Event:  org.OrgRemovedEventType,
// 					Reduce: p.reduceOrgRemoved,
// 				},
// 			},
// 		},
// 		{
// 			Aggregate: user.AggregateType,
// 			EventRedusers: []handler.EventReducer{
// 				{
// 					Event:  user.HumanEmailChangedType,
// 					Reduce: p.reduceHumanEmailChanged,
// 				},
// 				{
// 					Event:  user.UserV1EmailChangedType,
// 					Reduce: p.reduceHumanEmailChanged,
// 				},
// 				{
// 					Event:  user.HumanProfileChangedType,
// 					Reduce: p.reduceHumanProfileChanged,
// 				},
// 				{
// 					Event:  user.UserV1ProfileChangedType,
// 					Reduce: p.reduceHumanProfileChanged,
// 				},
// 			},
// 		},
// 	}
// }

// const (
// 	orgOwnerOrgID           = "org_id"
// 	orgOwnerOrgName         = "org_name"
// 	orgOwnerOrgCreationDate = "org_creation_date"
// 	orgOwnerOwnerID         = "owner_id"
// 	orgOwnerLanguage        = "owner_language"
// 	orgOwnerEmail           = "owner_email"
// 	orgOwnerFirstName       = "owner_first_name"
// 	orgOwnerLastName        = "owner_last_name"
// 	orgOwnerGender          = "owner_gender"
// )

// func (p *OrgOwnerProjection) reduceMemberAdded(event eventstore.EventReader) ([]handler.Statement, error) {
// 	e, ok := event.(*org.MemberAddedEvent)
// 	if !ok {
// 		logging.LogWithFields("PROJE-kL530", "seq", event.Sequence, "expected", org.MemberAddedEventType).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "PROJE-OkiBV", "reduce.wrong.event.type")
// 	}

// 	if !isOrgOwner(e.Roles) {
// 		return []handler.Statement{crdb.NewNoOpStatement(e.Aggregate().Typ, e.Sequence(), e.PreviousAggregateTypeSequence())}, nil
// 	}

// 	stmt, err := p.addOwner(e, e.Aggregate().ResourceOwner, e.UserID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return []handler.Statement{stmt}, nil
// }

// func (p *OrgOwnerProjection) reduceMemberChanged(event eventstore.EventReader) ([]handler.Statement, error) {
// 	e, ok := event.(*org.MemberChangedEvent)
// 	if !ok {
// 		logging.LogWithFields("PROJE-kL530", "seq", event.Sequence, "expected", org.MemberAddedEventType).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "PROJE-OkiBV", "reduce.wrong.event.type")
// 	}

// 	if !isOrgOwner(e.Roles) {
// 		return []handler.Statement{p.deleteOwner(e, e.Aggregate().ID, e.UserID)}, nil
// 	}

// 	stmt, err := p.addOwner(e, e.Aggregate().ResourceOwner, e.UserID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return []handler.Statement{stmt}, nil
// }

// func (p *OrgOwnerProjection) reduceMemberRemoved(event eventstore.EventReader) ([]handler.Statement, error) {
// 	e, ok := event.(*org.MemberRemovedEvent)
// 	if !ok {
// 		logging.LogWithFields("PROJE-boIbP", "seq", event.Sequence, "expected", org.MemberRemovedEventType).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "PROJE-pk6TS", "reduce.wrong.event.type")
// 	}

// 	return []handler.Statement{p.deleteOwner(e, e.Aggregate().ID, e.UserID)}, nil
// }

// func (p *OrgOwnerProjection) reduceOrgChanged(event eventstore.EventReader) ([]handler.Statement, error) {
// 	e, ok := event.(*org.OrgChangedEvent)
// 	if !ok {
// 		logging.LogWithFields("PROJE-piy2b", "seq", event.Sequence, "expected", org.OrgChangedEventType).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "PROJE-MGbru", "reduce.wrong.event.type")
// 	}

// 	values := []handler.Column{}
// 	if e.Name != "" {
// 		values = append(values, handler.NewCol(orgOwnerOrgName, e.Name))
// 	}

// 	if len(values) == 0 {
// 		return []handler.Statement{crdb.NewNoOpStatement(e.Aggregate().Typ, e.Sequence(), e.PreviousAggregateTypeSequence())}, nil
// 	}

// 	return []handler.Statement{
// 		crdb.NewUpdateStatement(
// 			e.Aggregate().Typ,
// 			e.Sequence(),
// 			e.PreviousAggregateTypeSequence(),
// 			values,
// 			[]handler.Column{
// 				handler.NewCol(orgOwnerOrgID, e.Aggregate().ResourceOwner),
// 			},
// 		),
// 	}, nil
// }

// func (p *OrgOwnerProjection) reduceOrgRemoved(event eventstore.EventReader) ([]handler.Statement, error) {
// 	e, ok := event.(*org.OrgChangedEvent)
// 	if !ok {
// 		logging.LogWithFields("PROJE-F1mHQ", "seq", event.Sequence, "expected", org.OrgRemovedEventType).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "PROJE-9ZR2w", "reduce.wrong.event.type")
// 	}

// 	return []handler.Statement{
// 		crdb.NewDeleteStatement(
// 			e.Aggregate().Typ,
// 			e.Sequence(),
// 			e.PreviousAggregateTypeSequence(),
// 			[]handler.Column{
// 				handler.NewCol(orgOwnerOrgID, e.Aggregate().ResourceOwner),
// 			},
// 		),
// 	}, nil
// }

// func (p *OrgOwnerProjection) reduceHumanEmailChanged(event eventstore.EventReader) ([]handler.Statement, error) {
// 	e, ok := event.(*user.HumanEmailChangedEvent)
// 	if !ok {
// 		logging.LogWithFields("PROJE-IHFwh", "seq", event.Sequence, "expected", user.HumanEmailChangedType).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "PROJE-jMlwT", "reduce.wrong.event.type")
// 	}

// 	return []handler.Statement{
// 		crdb.NewUpdateStatement(
// 			e.Aggregate().Typ,
// 			e.Sequence(),
// 			e.PreviousAggregateTypeSequence(),
// 			[]handler.Column{
// 				handler.NewCol(orgOwnerEmail, e.EmailAddress),
// 			},
// 			[]handler.Column{
// 				handler.NewCol(orgOwnerOwnerID, e.Aggregate().ID),
// 			},
// 		),
// 	}, nil
// }

// func (p *OrgOwnerProjection) reduceHumanProfileChanged(event eventstore.EventReader) ([]handler.Statement, error) {
// 	e, ok := event.(*user.HumanProfileChangedEvent)
// 	if !ok {
// 		logging.LogWithFields("PROJE-WqgUS", "seq", event.Sequence, "expected", user.HumanProfileChangedType).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Cdkkf", "reduce.wrong.event.type")
// 	}

// 	values := []handler.Column{}
// 	if e.FirstName != "" {
// 		values = append(values, handler.NewCol(orgOwnerFirstName, e.FirstName))
// 	}
// 	if e.LastName != "" {
// 		values = append(values, handler.NewCol(orgOwnerLastName, e.LastName))
// 	}
// 	if e.PreferredLanguage != nil {
// 		values = append(values, handler.NewCol(orgOwnerLanguage, e.PreferredLanguage.String()))
// 	}
// 	if e.Gender != nil {
// 		values = append(values, handler.NewCol(orgOwnerGender, *e.Gender))
// 	}

// 	if len(values) == 0 {
// 		return []handler.Statement{crdb.NewNoOpStatement(e.Aggregate().Typ, e.Sequence(), e.PreviousAggregateTypeSequence())}, nil
// 	}

// 	return []handler.Statement{
// 		crdb.NewUpdateStatement(
// 			e.Aggregate().Typ,
// 			e.Sequence(),
// 			e.PreviousAggregateTypeSequence(),
// 			values,
// 			[]handler.Column{
// 				handler.NewCol(orgOwnerOwnerID, e.Aggregate().ID),
// 			},
// 		),
// 	}, nil
// }

// func isOrgOwner(roles []string) bool {
// 	for _, role := range roles {
// 		if role == "ORG_OWNER" {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (p *OrgOwnerProjection) deleteOwner(event eventstore.EventReader, orgID, ownerID string) handler.Statement {
// 	return crdb.NewDeleteStatement(
// 		event.Aggregate().Typ,
// 		event.Sequence(),
// 		event.PreviousAggregateTypeSequence(),
// 		[]handler.Column{
// 			handler.NewCol(orgOwnerOrgID, orgID),
// 			handler.NewCol(orgOwnerOwnerID, ownerID),
// 		})
// }

// func (p *OrgOwnerProjection) addOwner(event eventstore.EventReader, orgID, userID string) (handler.Statement, error) {
// 	events, err := p.Eventstore.FilterEvents(context.Background(),
// 		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
// 			AddQuery().
// 			AggregateTypes(org.AggregateType).
// 			EventTypes(
// 				org.OrgAddedEventType,
// 				org.OrgChangedEventType).
// 			AggregateIDs(orgID).
// 			Or().
// 			AggregateTypes(user.AggregateType).
// 			EventTypes(
// 				user.HumanAddedType,
// 				user.UserV1AddedType,
// 				user.HumanRegisteredType,
// 				user.UserV1RegisteredType,
// 				user.HumanEmailChangedType,
// 				user.UserV1EmailChangedType,
// 				user.HumanProfileChangedType,
// 				user.UserV1ProfileChangedType).
// 			AggregateIDs(userID).
// 			SequenceLess(event.Sequence()).
// 			Builder())
// 	if err != nil {
// 		return handler.Statement{}, err
// 	}

// 	if len(events) == 0 {
// 		return handler.Statement{}, errors.ThrowInternal(nil, "PROJE-Qk7Tv", "unable to find org events")
// 	}

// 	owner := &OrgOwner{
// 		OrgID:   orgID,
// 		OwnerID: userID,
// 	}

// 	p.reduce(owner, events)

// 	values := []handler.Column{
// 		handler.NewCol(orgOwnerOrgID, owner.OrgID),
// 		handler.NewCol(orgOwnerOrgName, owner.OrgName),
// 		handler.NewCol(orgOwnerOrgCreationDate, owner.OrgCreationDate),
// 		handler.NewCol(orgOwnerOwnerID, owner.OwnerID),
// 		handler.NewCol(orgOwnerEmail, owner.OwnerEmailAddress),
// 		handler.NewCol(orgOwnerFirstName, owner.OwnerFirstName),
// 		handler.NewCol(orgOwnerLastName, owner.OwnerLastName),
// 		handler.NewCol(orgOwnerGender, owner.OwnerGender),
// 	}

// 	if owner.OwnerLanguage != nil {
// 		values = append(values, handler.NewCol(orgOwnerLanguage, owner.OwnerLanguage.String()))
// 	}

// 	return crdb.NewUpsertStatement(
// 		event.Aggregate().Typ,
// 		event.Sequence(),
// 		event.PreviousAggregateTypeSequence(),
// 		values), nil
// }

// func (p *OrgOwnerProjection) reduce(admin *OrgOwner, events []eventstore.EventReader) {
// 	for _, event := range events {
// 		switch e := event.(type) {
// 		case *user.HumanAddedEvent:
// 			admin.OwnerLanguage = &e.PreferredLanguage
// 			admin.OwnerEmailAddress = e.EmailAddress
// 			admin.OwnerFirstName = e.FirstName
// 			admin.OwnerLastName = e.LastName
// 			admin.OwnerGender = e.Gender
// 		case *user.HumanRegisteredEvent:
// 			admin.OwnerLanguage = &e.PreferredLanguage
// 			admin.OwnerEmailAddress = e.EmailAddress
// 			admin.OwnerFirstName = e.FirstName
// 			admin.OwnerLastName = e.LastName
// 			admin.OwnerGender = e.Gender
// 		case *user.HumanEmailChangedEvent:
// 			admin.OwnerEmailAddress = e.EmailAddress
// 		case *user.HumanProfileChangedEvent:
// 			if e.PreferredLanguage != nil {
// 				admin.OwnerLanguage = e.PreferredLanguage
// 			}
// 			if e.FirstName != "" {
// 				admin.OwnerFirstName = e.FirstName
// 			}
// 			if e.LastName != "" {
// 				admin.OwnerLastName = e.LastName
// 			}
// 			if e.Gender != nil {
// 				admin.OwnerGender = *e.Gender
// 			}
// 		case *org.OrgAddedEvent:
// 			admin.OrgName = e.Name
// 			admin.OrgCreationDate = e.CreationDate()
// 		case *org.OrgChangedEvent:
// 			if e.Name != "" {
// 				admin.OrgName = e.Name
// 			}
// 		default:
// 			// This happens only on implementation errors
// 			logging.LogWithFields("PROJE-sKNsR", "eventType", event.Type()).Panic("unexpected event type")
// 		}
// 	}
// }
