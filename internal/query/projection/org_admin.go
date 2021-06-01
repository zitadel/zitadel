package projection

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
	"golang.org/x/text/language"
)

type OrgAdmin struct {
	OrgID             string        `col:"org_id"`
	OrgName           string        `col:"org_name"`
	OrgCreationDate   time.Time     `col:"org_creation_date"`
	OwnerID           string        `col:"owner_id"`
	OwnerLanguage     *language.Tag `col:"owner_language"`
	OwnerEmailAddress string        `col:"owner_email"`
	OwnerFirstName    string        `col:"owner_first_name"`
	OwnerLastName     string        `col:"owner_last_name"`
}

type OrgAdminProjection struct {
	crdb.StatementHandler
}

func NewOrgAdminProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgAdminProjection {
	p := &OrgAdminProjection{}
	config.ProjectionName = "projections.org_admins"
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OrgAdminProjection) reducers() []handler.EventReducer {
	return []handler.EventReducer{
		{
			Aggregate: "org",
			Event:     org.MemberAddedEventType,
			Reduce:    p.reduceMemberAdded,
		},
		{
			Aggregate: "org",
			Event:     org.MemberChangedEventType,
			Reduce:    p.reduceMemberChanged,
		},
		{
			Aggregate: "org",
			Event:     org.MemberRemovedEventType,
			Reduce:    p.reduceMemberRemoved,
		},
		{
			Aggregate: "org",
			Event:     org.OrgChangedEventType,
			Reduce:    p.reduceOrgChanged,
		},
		{
			Aggregate: "org",
			Event:     org.OrgRemovedEventType,
			Reduce:    p.reduceOrgRemoved,
		},
		{
			Aggregate: "user",
			Event:     user.HumanEmailChangedType,
			Reduce:    p.reduceHumanEmailChanged,
		},
		{
			Aggregate: "user",
			Event:     user.HumanProfileChangedType,
			Reduce:    p.reduceHumanProfileChanged,
		},
	}
}

const (
	orgAdminOrgID           = "org_id"
	orgAdminOrgName         = "org_name"
	orgAdminOrgCreationDate = "org_creation_date"
	orgAdminOwnerID         = "owner_id"
	orgAdminOwnerLanguage   = "owner_language"
	orgAdminOwnerEmail      = "owner_email"
	orgAdminOwnerFirstName  = "owner_first_name"
	orgAdminOwnerLastName   = "owner_last_name"
)

func (p *OrgAdminProjection) reduceMemberAdded(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("PROJE-kL530", "seq", event.Sequence, "expected", org.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-OkiBV", "reduce.wrong.event.type")
	}

	if !isOrgOwner(e.Roles) {
		return []handler.Statement{crdb.NewNoOpStatement(e.Sequence(), e.PreviousSequence())}, nil
	}

	stmt, err := p.addAdmin(e.Aggregate().ResourceOwner, e.UserID, e.Sequence(), e.PreviousSequence())
	if err != nil {
		return nil, err
	}

	return []handler.Statement{stmt}, nil
}

func (p *OrgAdminProjection) reduceMemberChanged(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-kL530", "seq", event.Sequence, "expected", org.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-OkiBV", "reduce.wrong.event.type")
	}

	if !isOrgOwner(e.Roles) {
		return []handler.Statement{p.deleteAdmin(e.Aggregate().ID, e.UserID, e.Sequence(), e.PreviousSequence())}, nil
	}

	stmt, err := p.addAdmin(e.Aggregate().ResourceOwner, e.UserID, e.Sequence(), e.PreviousSequence())
	if err != nil {
		return nil, err
	}

	return []handler.Statement{stmt}, nil
}

func (p *OrgAdminProjection) reduceMemberRemoved(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-boIbP", "seq", event.Sequence, "expected", org.MemberRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-pk6TS", "reduce.wrong.event.type")
	}

	return []handler.Statement{p.deleteAdmin(e.Aggregate().ID, e.UserID, e.Sequence(), e.PreviousSequence())}, nil
}

func (p *OrgAdminProjection) reduceOrgChanged(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-piy2b", "seq", event.Sequence, "expected", org.OrgChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-MGbru", "reduce.wrong.event.type")
	}

	values := []handler.Column{}
	if e.Name != "" {
		values = append(values, handler.NewCol(orgAdminOrgName, e.Name))
	}

	if len(values) == 0 {
		return []handler.Statement{crdb.NewNoOpStatement(e.Sequence(), e.PreviousSequence())}, nil
	}

	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				handler.NewCol(orgAdminOrgID, e.Aggregate().ResourceOwner),
			},
			values,
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func (p *OrgAdminProjection) reduceOrgRemoved(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-F1mHQ", "seq", event.Sequence, "expected", org.OrgRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-9ZR2w", "reduce.wrong.event.type")
	}

	return []handler.Statement{
		crdb.NewDeleteStatement(
			[]handler.Column{
				handler.NewCol(orgAdminOrgID, e.Aggregate().ResourceOwner),
			},
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func (p *OrgAdminProjection) reduceHumanEmailChanged(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*user.HumanEmailChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-IHFwh", "seq", event.Sequence, "expected", user.HumanEmailChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-jMlwT", "reduce.wrong.event.type")
	}

	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				handler.NewCol(orgAdminOwnerID, e.Aggregate().ID),
			},
			[]handler.Column{handler.NewCol(orgAdminOwnerEmail, e.EmailAddress)},
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func (p *OrgAdminProjection) reduceHumanProfileChanged(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*user.HumanProfileChangedEvent)
	if !ok {
		logging.LogWithFields("PROJE-WqgUS", "seq", event.Sequence, "expected", user.HumanProfileChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Cdkkf", "reduce.wrong.event.type")
	}

	values := []handler.Column{}
	if e.FirstName != "" {
		values = append(values, handler.NewCol(orgAdminOwnerFirstName, e.FirstName))
	}
	if e.LastName != "" {
		values = append(values, handler.NewCol(orgAdminOwnerLastName, e.LastName))
	}
	if !e.PreferredLanguage.IsRoot() {
		values = append(values, handler.NewCol(orgAdminOwnerLanguage, e.PreferredLanguage))
	}

	if len(values) == 0 {
		return []handler.Statement{crdb.NewNoOpStatement(e.Sequence(), e.PreviousSequence())}, nil
	}

	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				handler.NewCol(orgAdminOwnerID, e.Aggregate().ID),
			},
			values,
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func isOrgOwner(roles []string) bool {
	for _, role := range roles {
		if role == "ORG_OWNER" {
			return true
		}
	}
	return false
}

func (p *OrgAdminProjection) deleteAdmin(orgID, ownerID string, sequence, previousSequence uint64) handler.Statement {
	return crdb.NewDeleteStatement([]handler.Column{
		handler.NewCol(orgAdminOrgID, orgID),
		handler.NewCol(orgAdminOwnerID, ownerID),
	}, sequence, previousSequence)
}

func (p *OrgAdminProjection) addAdmin(orgID, userID string, sequence, previousSequence uint64) (handler.Statement, error) {
	events, err := p.Eventstore.FilterEvents(context.Background(),
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType, user.AggregateType).
			EventTypes(org.OrgAddedEventType, org.OrgChangedEventType, org.OrgDomainRemovedEventType,
				user.HumanAddedType, user.HumanEmailChangedType, user.HumanProfileChangedType).
			AggregateIDs(orgID, userID).
			SequenceLess(sequence))
	if err != nil {
		return handler.Statement{}, err
	}

	if len(events) == 0 {
		return handler.Statement{}, errors.ThrowInternal(nil, "PROJE-Qk7Tv", "unable to find org events")
	}

	admin := &OrgAdmin{
		OrgID:   orgID,
		OwnerID: userID,
	}

	p.reduce(admin, events)

	return crdb.NewCreateStatement([]handler.Column{
		handler.NewCol(orgAdminOrgID, admin.OrgID),
		handler.NewCol(orgAdminOrgName, admin.OrgName),
		handler.NewCol(orgAdminOrgCreationDate, admin.OrgCreationDate),
		handler.NewCol(orgAdminOwnerID, admin.OwnerID),
		handler.NewCol(orgAdminOwnerLanguage, admin.OwnerLanguage.String()),
		handler.NewCol(orgAdminOwnerEmail, admin.OwnerEmailAddress),
		handler.NewCol(orgAdminOwnerFirstName, admin.OwnerFirstName),
		handler.NewCol(orgAdminOwnerLastName, admin.OwnerLastName),
	}, sequence, previousSequence), nil
}

func (p *OrgAdminProjection) reduce(admin *OrgAdmin, events []eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			admin.OwnerLanguage = &e.PreferredLanguage
			admin.OwnerEmailAddress = e.EmailAddress
			admin.OwnerFirstName = e.FirstName
			admin.OwnerLastName = e.LastName
		case *user.HumanEmailChangedEvent:
			admin.OwnerEmailAddress = e.EmailAddress
		case *user.HumanProfileChangedEvent:
			if e.PreferredLanguage != nil {
				admin.OwnerLanguage = e.PreferredLanguage
			}
			if e.FirstName != "" {
				admin.OwnerFirstName = e.FirstName
			}
			if e.LastName != "" {
				admin.OwnerLastName = e.LastName
			}
		case *org.OrgAddedEvent:
			admin.OrgName = e.Name
			admin.OrgCreationDate = e.CreationDate()
		case *org.OrgChangedEvent:
			if e.Name != "" {
				admin.OrgName = e.Name
			}
		default:
			// This happens only on implementation errors
			logging.LogWithFields("PROJE-sKNsR", "eventType", event.Type()).Panic("unexpected event type")
		}
	}
}
