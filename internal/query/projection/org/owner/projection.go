package owner

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
	"golang.org/x/text/language"
)

type OrgOwner struct {
	OrgID             string
	OrgName           string
	OrgCreationDate   time.Time
	OwnerID           string
	OwnerLanguage     *language.Tag
	OwnerEmailAddress string
	OwnerFirstName    string
	OwnerLastName     string
	OwnerGender       domain.Gender
}

type OrgOwnerProjection struct {
	crdb.StatementHandler
}

const (
	orgTableSuffix     = "orgs"
	orgIDCol           = "id"
	orgNameCol         = "name"
	orgCreationDateCol = "creation_date"

	userTableSuffix  = "users"
	userOrgIDCol     = "org_id"
	userIDCol        = "owner_id"
	userLanguageCol  = "language"
	userEmailCol     = "email"
	userFirstNameCol = "first_name"
	userLastNameCol  = "last_name"
	userGenderCol    = "gender"
)

func NewOrgOwnerProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgOwnerProjection {
	p := &OrgOwnerProjection{}
	config.ProjectionName = "zitadel.projections.org_owners"
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OrgOwnerProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgAddedEventType,
					Reduce: p.reduceOrgAdded,
				},
				{
					Event:  org.OrgChangedEventType,
					Reduce: p.reduceOrgChanged,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
				},
				{
					Event:  org.MemberAddedEventType,
					Reduce: p.reduceMemberAdded,
				},
				{
					Event:  org.MemberChangedEventType,
					Reduce: p.reduceMemberChanged,
				},
				{
					Event:  org.MemberRemovedEventType,
					Reduce: p.reduceMemberRemoved,
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.HumanEmailChangedType,
					Reduce: p.reduceHumanEmailChanged,
				},
				{
					Event:  user.UserV1EmailChangedType,
					Reduce: p.reduceHumanEmailChanged,
				},
				{
					Event:  user.HumanProfileChangedType,
					Reduce: p.reduceHumanProfileChanged,
				},
				{
					Event:  user.UserV1ProfileChangedType,
					Reduce: p.reduceHumanProfileChanged,
				},
			},
		},
	}
}

func (p *OrgOwnerProjection) reduceMemberAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("PROJE-kL530", "seq", event.Sequence, "expected", org.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-OkiBV", "reduce.wrong.event.type")
	}

	if !isOrgOwner(e.Roles) {
		return crdb.NewNoOpStatement(e), nil
	}

	stmt, err := p.addOwner(e, e.Aggregate().ResourceOwner, e.UserID)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *OrgOwnerProjection) reduceMemberChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-kL530", "seq", event.Sequence, "expected", org.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-OkiBV", "reduce.wrong.event.type")
	}

	if !isOrgOwner(e.Roles) {
		return p.deleteOwner(e, e.Aggregate().ID, e.UserID), nil
	}

	stmt, err := p.addOwner(e, e.Aggregate().ResourceOwner, e.UserID)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *OrgOwnerProjection) reduceMemberRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-boIbP", "seq", event.Sequence, "expected", org.MemberRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-pk6TS", "reduce.wrong.event.type")
	}

	return p.deleteOwner(e, e.Aggregate().ID, e.UserID), nil
}

func (p *OrgOwnerProjection) reduceHumanEmailChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-IHFwh", "seq", event.Sequence, "expected", user.HumanEmailChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-jMlwT", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(userEmailCol, e.EmailAddress),
		},
		[]handler.Condition{
			handler.NewCond(userIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(userTableSuffix),
	), nil
}

func (p *OrgOwnerProjection) reduceHumanProfileChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanProfileChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-WqgUS", "seq", event.Sequence, "expected", user.HumanProfileChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Cdkkf", "reduce.wrong.event.type")
	}

	values := []handler.Column{}
	if e.FirstName != "" {
		values = append(values, handler.NewCol(userFirstNameCol, e.FirstName))
	}
	if e.LastName != "" {
		values = append(values, handler.NewCol(userLastNameCol, e.LastName))
	}
	if e.PreferredLanguage != nil {
		values = append(values, handler.NewCol(userLanguageCol, e.PreferredLanguage.String()))
	}
	if e.Gender != nil {
		values = append(values, handler.NewCol(userGenderCol, *e.Gender))
	}

	if len(values) == 0 {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewUpdateStatement(
		e,
		values,
		[]handler.Condition{
			handler.NewCond(userIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(userTableSuffix),
	), nil
}

func (p *OrgOwnerProjection) reduceOrgAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		logging.LogWithFields("PROJE-wbOrL", "seq", event.Sequence, "expected", org.OrgAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-pk6TS", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(orgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCol(orgNameCol, e.Name),
			handler.NewCol(orgCreationDateCol, e.CreationDate()),
		},
		crdb.WithTableSuffix(orgTableSuffix),
	), nil
}

func (p *OrgOwnerProjection) reduceOrgChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-piy2b", "seq", event.Sequence, "expected", org.OrgChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-MGbru", "reduce.wrong.event.type")
	}

	values := []handler.Column{}
	if e.Name != "" {
		values = append(values, handler.NewCol(orgNameCol, e.Name))
	}

	if len(values) == 0 {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewUpdateStatement(
		e,
		values,
		[]handler.Condition{
			handler.NewCond(orgIDCol, e.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(orgTableSuffix),
	), nil
}

func (p *OrgOwnerProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-F1mHQ", "seq", event.Sequence, "expected", org.OrgRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-9ZR2w", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(e,
		//delete org in org table
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(orgIDCol, e.Aggregate().ResourceOwner),
			},
			crdb.WithTableSuffix(orgTableSuffix),
		),
		// delete users of the org
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(userOrgIDCol, e.Aggregate().ResourceOwner),
			},
			crdb.WithTableSuffix(userTableSuffix),
		),
	), nil
}

func isOrgOwner(roles []string) bool {
	for _, role := range roles {
		if role == "ORG_OWNER" {
			return true
		}
	}
	return false
}

func (p *OrgOwnerProjection) deleteOwner(event eventstore.Event, orgID, ownerID string) *handler.Statement {
	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(userOrgIDCol, orgID),
			handler.NewCond(userIDCol, ownerID),
		},
		crdb.WithTableSuffix(userTableSuffix),
	)
}

func (p *OrgOwnerProjection) addOwner(event eventstore.Event, orgID, userID string) (*handler.Statement, error) {
	events, err := p.Eventstore.Filter(context.Background(),
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			AddQuery().
			AggregateTypes(user.AggregateType).
			EventTypes(
				user.HumanAddedType,
				user.UserV1AddedType,
				user.HumanRegisteredType,
				user.UserV1RegisteredType,
				user.HumanEmailChangedType,
				user.UserV1EmailChangedType,
				user.HumanProfileChangedType,
				user.UserV1ProfileChangedType,
				user.MachineAddedEventType,
				user.MachineChangedEventType).
			AggregateIDs(userID).
			SequenceLess(event.Sequence()).
			Builder())
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		logging.LogWithFields("mqd3w", "user", userID, "org", orgID, "seq", event.Sequence()).Warn("no events for user found")
		return nil, errors.ThrowInternal(nil, "PROJE-Qk7Tv", "unable to find user events")
	}

	owner := &OrgOwner{
		OrgID:   orgID,
		OwnerID: userID,
	}

	p.reduce(owner, events)

	values := []handler.Column{
		handler.NewCol(userOrgIDCol, owner.OrgID),
		handler.NewCol(userIDCol, owner.OwnerID),
		handler.NewCol(userEmailCol, owner.OwnerEmailAddress),
		handler.NewCol(userFirstNameCol, owner.OwnerFirstName),
		handler.NewCol(userLastNameCol, owner.OwnerLastName),
		handler.NewCol(userGenderCol, owner.OwnerGender),
	}

	if owner.OwnerLanguage != nil {
		values = append(values, handler.NewCol(userLanguageCol, owner.OwnerLanguage.String()))
	}

	return crdb.NewUpsertStatement(
		event,
		values,
		crdb.WithTableSuffix(userTableSuffix),
	), nil
}

func (p *OrgOwnerProjection) reduce(owner *OrgOwner, events []eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			owner.OwnerLanguage = &e.PreferredLanguage
			owner.OwnerEmailAddress = e.EmailAddress
			owner.OwnerFirstName = e.FirstName
			owner.OwnerLastName = e.LastName
			owner.OwnerGender = e.Gender
		case *user.HumanRegisteredEvent:
			owner.OwnerLanguage = &e.PreferredLanguage
			owner.OwnerEmailAddress = e.EmailAddress
			owner.OwnerFirstName = e.FirstName
			owner.OwnerLastName = e.LastName
			owner.OwnerGender = e.Gender
		case *user.HumanEmailChangedEvent:
			owner.OwnerEmailAddress = e.EmailAddress
		case *user.HumanProfileChangedEvent:
			if e.PreferredLanguage != nil {
				owner.OwnerLanguage = e.PreferredLanguage
			}
			if e.FirstName != "" {
				owner.OwnerFirstName = e.FirstName
			}
			if e.LastName != "" {
				owner.OwnerLastName = e.LastName
			}
			if e.Gender != nil {
				owner.OwnerGender = *e.Gender
			}
		case *user.MachineAddedEvent:
			owner.OwnerFirstName = "machine"
			owner.OwnerLastName = e.Name
			owner.OwnerEmailAddress = e.UserName
		case *user.MachineChangedEvent:
			if e.Name != nil {
				owner.OwnerLastName = *e.Name
			}
		default:
			// This happens only on implementation errors
			logging.LogWithFields("PROJE-sKNsR", "eventType", event.Type()).Panic("unexpected event type")
		}
	}
}
