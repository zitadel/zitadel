package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type SessionCheck func(ctx context.Context) error

type SessionChecks struct {
	userCheck SessionCheck
	userID    string
	checks    []SessionCheck

	sessionWriteModel  *SessionWriteModel
	passwordWriteModel *HumanPasswordWriteModel
	eventstore         *eventstore.Eventstore
	userPasswordAlg    crypto.HashAlgorithm
}

func (c *Commands) NewSessionChecks() *SessionChecks {
	return &SessionChecks{
		eventstore:      c.eventstore,
		userPasswordAlg: c.userPasswordAlg,
	}
}

func (s *SessionChecks) CheckUser(id string) {
	s.userCheck = func(ctx context.Context) error {
		if s.sessionWriteModel.UserID != "" && s.userID != "" && s.sessionWriteModel.UserID != s.userID {
			return caos_errs.ThrowInvalidArgument(nil, "", "user change not possible")
		}
		s.sessionWriteModel.UserChecked(ctx, id, time.Now())
		return nil
	}
}

func (s *SessionChecks) CheckPassword(password string) {
	s.checks = append(s.checks, func(ctx context.Context) error {
		if s.sessionWriteModel.UserID == "" {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sfw3f", "Errors.User.UserIDMissing")
		}
		s.passwordWriteModel = NewHumanPasswordWriteModel(s.sessionWriteModel.UserID, "")
		err := s.eventstore.FilterToQueryReducer(ctx, s.passwordWriteModel)
		if err != nil {
			return err
		}
		if s.passwordWriteModel.UserState == domain.UserStateUnspecified || s.passwordWriteModel.UserState == domain.UserStateDeleted {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Df4b3", "Errors.User.NotFound")
		}

		if s.passwordWriteModel.Secret == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-WEf3t", "Errors.User.Password.NotSet")
		}
		ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
		err = crypto.CompareHash(s.passwordWriteModel.Secret, []byte(password), s.userPasswordAlg)
		spanPasswordComparison.EndWithError(err)
		if err != nil {
			//TODO: reset?
			return caos_errs.ThrowInvalidArgument(err, "COMMAND-SAF3g", "Errors.User.Password.Invalid")
		}
		s.sessionWriteModel.PasswordChecked(ctx, time.Now())
		return nil
	})
}

func (s *SessionChecks) Check(ctx context.Context) error {
	if s.userCheck != nil {
		if err := s.userCheck(ctx); err != nil {
			return err
		}
	}
	for _, check := range s.checks {
		if err := check(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionChecks) events(ctx context.Context) []eventstore.Command {
	if s.sessionWriteModel.State == domain.SessionStateActive && s.sessionWriteModel.event == nil {
		return nil
	}
	return []eventstore.Command{s.sessionWriteModel.SetToken(ctx)}
}

func (c *Commands) CreateSession(ctx context.Context, checks *SessionChecks, metadata map[string][]byte) (set *SessionChanged, err error) {
	sessionID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	checks.sessionWriteModel = NewSessionWriteModel(sessionID, authz.GetCtxData(ctx).OrgID)
	err = c.eventstore.FilterToQueryReducer(ctx, checks.sessionWriteModel)
	if err != nil {
		return nil, err
	}
	return c.changeSession(ctx, checks, metadata)
}

func (c *Commands) UpdateSession(ctx context.Context, sessionID, sessionToken string, checks *SessionChecks, metadata map[string][]byte) (set *SessionChanged, err error) {
	checks.sessionWriteModel = NewSessionWriteModel(sessionID, authz.GetCtxData(ctx).OrgID)
	err = c.eventstore.FilterToQueryReducer(ctx, checks.sessionWriteModel)
	if err != nil {
		return nil, err
	}
	if err := c.sessionPermission(ctx, checks.sessionWriteModel, sessionToken, permissionSessionWrite); err != nil {
		return nil, err
	}
	return c.changeSession(ctx, checks, metadata)
}

func (c *Commands) TerminateSession(ctx context.Context, sessionID, sessionToken string) (*SessionChanged, error) {
	sessionWriteModel := NewSessionWriteModel(sessionID, "")
	if err := c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel); err != nil {
		return nil, err
	}
	if err := c.sessionPermission(ctx, sessionWriteModel, sessionToken, permissionSessionDelete); err != nil {
		return nil, err
	}
	if sessionWriteModel.State != domain.SessionStateActive {
		return sessionWriteModelToSessionChanged(sessionWriteModel), nil
	}
	terminate := session.NewTerminateEvent(ctx, &session.NewAggregate(sessionWriteModel.AggregateID, sessionWriteModel.ResourceOwner).Aggregate)
	pushedEvents, err := c.eventstore.Push(ctx, terminate)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(sessionWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return sessionWriteModelToSessionChanged(sessionWriteModel), nil
}

func (c *Commands) changeSession(ctx context.Context, checks *SessionChecks, metadata map[string][]byte) (set *SessionChanged, err error) {
	if checks.sessionWriteModel.State == domain.SessionStateTerminated {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMAND-SAjeh", "Errors.Session.Terminated") //TODO: i18n
	}
	if err := checks.Check(ctx); err != nil {
		//TODO: how to handle failed checks (e.g. pw wrong)
		//if e := checks.sessionWriteModel.event; e != nil {
		//	_, err := c.eventstore.Push(ctx, e)
		//	logging.OnError(err).Error("could not push event check failed events")
		//}
		return nil, err
	}
	if err := checks.sessionWriteModel.ChangeMetadata(ctx, metadata); err != nil {
		return nil, err
	}
	cmds := checks.events(ctx)
	if len(cmds) == 0 {
		return sessionWriteModelToSessionChanged(checks.sessionWriteModel), nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(checks.sessionWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return sessionWriteModelToSessionChanged(checks.sessionWriteModel), nil
}

func (c *Commands) sessionPermission(ctx context.Context, sessionWriteModel *SessionWriteModel, sessionToken, permission string) error {
	if sessionToken == "" {
		return c.checkPermission(ctx, permission, authz.GetCtxData(ctx).OrgID, sessionWriteModel.AggregateID)
	}
	if sessionWriteModel.Token == sessionToken {
		return nil
	}
	return caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Invalid token") //TODO: i18n
}

type SessionChanged struct {
	*domain.ObjectDetails
	ID    string
	Token string
}

func sessionWriteModelToSessionChanged(wm *SessionWriteModel) *SessionChanged {
	return &SessionChanged{
		ObjectDetails: &domain.ObjectDetails{
			Sequence:      wm.ProcessedSequence,
			EventDate:     wm.ChangeDate,
			ResourceOwner: wm.ResourceOwner,
		},
		ID:    wm.AggregateID,
		Token: wm.Token,
	}
}
