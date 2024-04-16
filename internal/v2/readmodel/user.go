package readmodel

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type User struct {
	readModel
	projection.User
	LoginNames *projection.LoginNames
	Human      *projection.Human
	Machine    *projection.Machine
}

func (u *User) PreferredLoginName() string {
	if u.LoginNames == nil {
		return ""
	}
	for _, loginName := range u.LoginNames.LoginNames {
		if loginName.IsPrimary {
			return loginName.Name
		}
	}

	return ""
}

func NewUser(id string) *User {
	return &User{
		User:    *projection.NewUserProjection(id),
		Human:   projection.NewHumanProjection(id),
		Machine: projection.NewMachineProjection(id),
	}
}

func (u *User) Query(ctx context.Context, querier eventstore.Querier, opts ...QueryOpt) error {
	queryOpts := make([]eventstore.QueryOpt, 0, len(opts)+1)
	queryOpts = append(queryOpts, eventstore.AppendFilters(u.User.Filter()...), eventstore.AppendFilters(u.Human.Filter()...), eventstore.AppendFilters(u.Machine.Filter()...))
	for _, opt := range opts {
		queryOpts = opt(queryOpts)
	}

	eventCount, err := querier.Query(
		ctx,
		eventstore.NewQuery(
			authz.GetInstance(ctx).InstanceID(),
			u,
			queryOpts...,
		),
	)
	if err != nil {
		return err
	}
	if eventCount == 0 {
		return zerrors.ThrowNotFound(nil, "READM-TWPSk", "Errors.User.NotFound")
	}
	u.LoginNames = projection.NewLoginNamesWithOwner(u.ID, u.Instance, u.Owner)

	queryOpts = make([]eventstore.QueryOpt, 0, len(opts)+1)
	queryOpts = append(queryOpts, eventstore.AppendFilters(u.LoginNames.Filter()...))
	for _, opt := range opts {
		queryOpts = opt(queryOpts)
	}

	_, err = querier.Query(
		ctx,
		eventstore.NewQuery(
			authz.GetInstance(ctx).InstanceID(),
			u.LoginNames,
			queryOpts...,
		),
	)

	if err != nil {
		return err
	}

	u.LoginNames.Generate()
	return nil
}

func (u *User) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) (err error) {

eventLoop:
	for _, event := range events {
		switch event.Type {
		case "user.human.added", "user.added", "user.human.selfregistered":
			u.Machine = nil
			break eventLoop
		case "user.machine.added":
			u.Human = nil
			break eventLoop
		}
	}
	if err = u.User.Reduce(events...); err != nil {
		return err
	}
	u.reduce(events[len(events)-1])
	if u.Type == domain.UserTypeHuman {
		err = u.Human.Reduce(events...)
	} else if u.Type == domain.UserTypeMachine {
		err = u.Machine.Reduce(events...)
	}

	return err
}
