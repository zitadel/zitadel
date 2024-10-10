package query

import (
	"cmp"
	"context"
	"slices"
	"strconv"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type MilestoneReadModel struct {
	eventstore.ReadModel `json:"-"`
	startPosition        float64
	systemUsers          map[string]*authz.SystemAPIUser
	Milestones           []*Milestone
	IgnoredClientIDs     []string
}

func NewMilestoneReadModel(ctx context.Context, systemUsers map[string]*authz.SystemAPIUser) *MilestoneReadModel {
	instanceID := authz.GetInstance(ctx).InstanceID()
	milestones := make([]*Milestone, len(milestone.TypeValues()))
	for i, typ := range milestone.TypeValues() {
		milestones[i] = &Milestone{
			InstanceID: instanceID,
			Type:       typ,
		}
	}
	return &MilestoneReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID:   "",
			InstanceID:    instanceID,
			ResourceOwner: instanceID,
		},
		Milestones: milestones,
	}
}

func (m *MilestoneReadModel) Reduce() error {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *instance.InstanceAddedEvent:
			m.reduceReached(milestone.InstanceCreated, event)
		case *instance.DomainPrimarySetEvent:
			for _, milestone := range m.Milestones {
				milestone.PrimaryDomain = e.Domain
			}
		case *project.ProjectAddedEvent:
			m.reduceReachedIfUserEvent(milestone.ProjectCreated, event)
		case *project.ApplicationAddedEvent:
			m.reduceReachedIfUserEvent(milestone.ApplicationCreated, event)
		case *project.OIDCConfigAddedEvent:
			m.reduceAppConfigAdded(event, e.ClientID)
		case *project.APIConfigAddedEvent:
			m.reduceAppConfigAdded(event, e.ClientID)
		case *oidcsession.AddedEvent:
			m.reduceReached(milestone.AuthenticationSucceededOnInstance, event)
			// Ignore authentications without session, for example JWT profile,
			if e.SessionID == "" {
				m.reduceReached(milestone.AuthenticationSucceededOnApplication, event)
			}
		case *instance.InstanceRemovedEvent:
			m.reduceReached(milestone.InstanceDeleted, event)
		case *milestone.PushedEvent:

		}

	}
	return m.ReadModel.Reduce()
}

func (m *MilestoneReadModel) Query() *eventstore.SearchQueryBuilder {
	var (
		instanceEvents   []eventstore.EventType
		projectEvents    []eventstore.EventType
		oidcsessionEvent eventstore.EventType
	)
	if m.notReached(milestone.InstanceCreated) {
		instanceEvents = append(instanceEvents, instance.InstanceAddedEventType)
	}
	if m.emptyDomain() {
		instanceEvents = append(instanceEvents, instance.InstanceDomainPrimarySetEventType)
	}
	if m.notReached(milestone.AuthenticationSucceededOnInstance) {
		oidcsessionEvent = oidcsession.AddedType
	}
	if m.notReached(milestone.ProjectCreated) {
		projectEvents = append(projectEvents, project.ProjectAddedType)
	}
	if m.notReached(milestone.ApplicationCreated) {
		projectEvents = append(projectEvents, project.ApplicationAddedType)
	}
	if m.notReached(milestone.AuthenticationSucceededOnApplication) {
		projectEvents = append(projectEvents, project.OIDCConfigAddedType, project.APIConfigAddedType)
		oidcsessionEvent = oidcsession.AddedType
	}
	if m.notReached(milestone.InstanceDeleted) {
		instanceEvents = append(instanceEvents, instance.InstanceRemovedEventType)
	}

	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent)
	builder.InstanceID(m.InstanceID)
	if len(instanceEvents) > 0 {
		builder = builder.AddQuery().
			AggregateTypes(instance.AggregateType).
			EventTypes(instanceEvents...).
			PositionAfter(m.Position).
			Builder()
	}
	if len(projectEvents) > 0 {
		builder = builder.AddQuery().
			AggregateTypes(project.AggregateType).
			EventTypes(projectEvents...).
			PositionAfter(m.Position).
			Builder()
	}
	if oidcsessionEvent != "" {
		builder = builder.AddQuery().
			AggregateTypes(oidcsession.AggregateType).
			EventTypes(oidcsessionEvent).
			PositionAfter(m.Position).
			Builder()
	}
	if m.unPushed() {
		builder = builder.AddQuery().
			AggregateTypes(milestone.AggregateType).
			EventTypes(milestone.PushedEventType).
			PositionAfter(m.Position).
			Builder()
	}
	return builder
}

func (m *MilestoneReadModel) reduceReached(typ milestone.Type, event eventstore.Event) {
	milestone, ok := m.getMilestone(typ)
	if ok && milestone.ReachedDate.IsZero() {
		milestone.ReachedDate = event.CreatedAt()
	}
}

func (m *MilestoneReadModel) reduceReachedIfUserEvent(typ milestone.Type, event eventstore.Event) {
	if !m.isSystemEvent(event) {
		m.reduceReached(typ, event)
	}
}

func (m *MilestoneReadModel) reduceAppConfigAdded(event eventstore.Event, clientID string) {
	if !m.isSystemEvent(event) {
		return
	}
	milestone, ok := m.getMilestone(milestone.AuthenticationSucceededOnApplication)
	if ok && milestone.ReachedDate.IsZero() {
		m.IgnoredClientIDs = append(m.IgnoredClientIDs, clientID)
	}
}

func (m *MilestoneReadModel) isSystemEvent(event eventstore.Event) bool {
	if userId, err := strconv.Atoi(event.Creator()); err == nil && userId > 0 {
		return false
	}

	// check if it is a hard coded event creator
	for _, creator := range []string{"", "system", "OIDC", "LOGIN", "SYSTEM"} {
		if creator == event.Creator() {
			return true
		}
	}

	_, ok := m.systemUsers[event.Creator()]
	return ok
}

func (m *MilestoneReadModel) reducePushed(e *milestone.PushedEvent) {
	milestone, ok := m.getMilestone(e.MilestoneType)
	if ok {
		milestone.PushedDate = e.CreatedAt()
	}
}

func (m *MilestoneReadModel) allDone() bool {
	return !slices.ContainsFunc(m.Milestones, func(milestone *Milestone) bool {
		// contains unfinished milestone
		return milestone.PrimaryDomain == "" ||
			milestone.ReachedDate.IsZero() ||
			milestone.PushedDate.IsZero()
	})
}

func (m *MilestoneReadModel) emptyDomain() bool {
	return slices.ContainsFunc(m.Milestones, func(milestone *Milestone) bool {
		// contains milestone without domain
		return milestone.PrimaryDomain == ""
	})
}

func (m *MilestoneReadModel) notReached(typ milestone.Type) bool {
	milestone, ok := m.getMilestone(typ)
	return ok && milestone.ReachedDate.IsZero()
}

func (m *MilestoneReadModel) unPushed() bool {
	return slices.ContainsFunc(m.Milestones, func(milestone *Milestone) bool {
		// contains un-pushed milestone
		return !milestone.ReachedDate.IsZero() && milestone.PushedDate.IsZero()
	})
}

func (m *MilestoneReadModel) getMilestone(typ milestone.Type) (*Milestone, bool) {
	i, ok := slices.BinarySearchFunc(m.Milestones, typ, func(m *Milestone, typ milestone.Type) int {
		return cmp.Compare(m.Type, typ)
	})
	if ok {
		return m.Milestones[i], ok
	}
	return nil, ok
}

func (q *Queries) GetMilestoneSnapshot(ctx context.Context, systemUsers map[string]*authz.SystemAPIUser) (*MilestoneReadModel, error) {
	model := NewMilestoneReadModel(ctx, systemUsers)
	err := eventstore.SnapshotFromReadModel(&model.ReadModel, model).Populate(ctx, q.eventstore)
	if err != nil {
		return nil, err
	}
	model.startPosition = model.Position
	return model, nil
}
