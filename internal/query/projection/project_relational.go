package projection

import (
	"context"

	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ProjectRelationalTable                       = "zitadel.projects"
	ProjectRelationalShouldAssertRoleCol         = "should_assert_role"
	ProjectRelationalIsAuthorizationRequiredCol  = "is_authorization_required"
	ProjectRelationalIsProjectAccessRequiredCol  = "is_project_access_required"
	ProjectRelationalUsedLabelingSettingOwnerCol = "used_labeling_setting_owner"
)

type projectRelationalProjection struct{}

func (*projectRelationalProjection) Name() string {
	return ProjectRelationalTable
}

func newProjectRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(projectRelationalProjection))
}

func (p *projectRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.ProjectAddedType,
					Reduce: p.reduceProjectAdded,
				},
				{
					Event:  project.ProjectChangedType,
					Reduce: p.reduceProjectChanged,
				},
				{
					Event:  project.ProjectDeactivatedType,
					Reduce: p.reduceProjectDeactivated,
				},
				{
					Event:  project.ProjectReactivatedType,
					Reduce: p.reduceProjectReactivated,
				},
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
	}
}

func (p *projectRelationalProjection) reduceProjectAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.ProjectAddedType)
	}

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(InstanceID, e.Aggregate().InstanceID),
			handler.NewCol(OrganizationID, e.Aggregate().ResourceOwner),
			handler.NewCol(ID, e.Aggregate().ID),
			handler.NewCol(CreatedAt, e.CreationDate()),
			handler.NewCol(UpdatedAt, e.CreationDate()),
			handler.NewCol(ProjectColumnName, e.Name),
			handler.NewCol(State, repoDomain.ProjectStateActive),
			handler.NewCol(ProjectRelationalShouldAssertRoleCol, e.ProjectRoleAssertion),
			handler.NewCol(ProjectRelationalIsAuthorizationRequiredCol, e.ProjectRoleCheck),
			handler.NewCol(ProjectRelationalIsProjectAccessRequiredCol, e.HasProjectCheck),
			handler.NewCol(ProjectRelationalUsedLabelingSettingOwnerCol, e.PrivateLabelingSetting),
		},
	), nil
}

func (p *projectRelationalProjection) reduceProjectChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectChangeEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.ProjectChangedType)
	}
	if e.Name == nil && e.HasProjectCheck == nil && e.ProjectRoleAssertion == nil && e.ProjectRoleCheck == nil && e.PrivateLabelingSetting == nil {
		return handler.NewNoOpStatement(e), nil
	}

	columns := make([]handler.Column, 0, 6)
	columns = append(columns,
		handler.NewCol(UpdatedAt, e.CreationDate()),
	)
	if e.Name != nil {
		columns = append(columns, handler.NewCol(ProjectColumnName, *e.Name))
	}
	if e.ProjectRoleAssertion != nil {
		columns = append(columns, handler.NewCol(ProjectRelationalShouldAssertRoleCol, *e.ProjectRoleAssertion))
	}
	if e.ProjectRoleCheck != nil {
		columns = append(columns, handler.NewCol(ProjectRelationalIsAuthorizationRequiredCol, *e.ProjectRoleCheck))
	}
	if e.HasProjectCheck != nil {
		columns = append(columns, handler.NewCol(ProjectRelationalIsProjectAccessRequiredCol, *e.HasProjectCheck))
	}
	if e.PrivateLabelingSetting != nil {
		columns = append(columns, handler.NewCol(ProjectRelationalUsedLabelingSettingOwnerCol, *e.PrivateLabelingSetting))
	}
	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(e.Aggregate().InstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ID, e.Aggregate().ID),
		},
	), nil
}

func (p *projectRelationalProjection) reduceProjectDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.ProjectDeactivatedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UpdatedAt, e.CreationDate()),
			handler.NewCol(State, repoDomain.ProjectStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ID, e.Aggregate().ID),
		},
	), nil
}

func (p *projectRelationalProjection) reduceProjectReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-oof4U", "reduce.wrong.event.type %s", project.ProjectDeactivatedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UpdatedAt, e.CreationDate()),
			handler.NewCol(State, repoDomain.ProjectStateActive),
		},
		[]handler.Condition{
			handler.NewCond(InstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ID, e.Aggregate().ID),
		},
	), nil
}

func (p *projectRelationalProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Xae7w", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ID, e.Aggregate().ID),
		},
	), nil
}
