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
	checks    []SessionCheck

	sessionWriteModel  *SessionWriteModel
	passwordWriteModel *HumanPasswordWriteModel
	eventstore         *eventstore.Eventstore
	userPasswordAlg    crypto.HashAlgorithm
	sessionAlg         crypto.EncryptionAlgorithm
}

func (c *Commands) NewSessionChecks() *SessionChecks {
	return &SessionChecks{
		eventstore:      c.eventstore,
		userPasswordAlg: c.userPasswordAlg,
		sessionAlg:      c.keyAlgorithm,
	}
}

// CheckUser defines a user check to be executed for a session update
func (s *SessionChecks) CheckUser(id string) {
	s.userCheck = func(ctx context.Context) error {
		// TODO: check here?
		if s.sessionWriteModel.UserID != "" && id != "" && s.sessionWriteModel.UserID != id {
			return caos_errs.ThrowInvalidArgument(nil, "", "user change not possible")
		}
		return s.sessionWriteModel.UserChecked(ctx, id, time.Now())
	}
}

// CheckPassword defines a password check to be executed for a session update
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
			//TODO: reset session?
			return caos_errs.ThrowInvalidArgument(err, "COMMAND-SAF3g", "Errors.User.Password.Invalid")
		}
		s.sessionWriteModel.PasswordChecked(ctx, time.Now())
		return nil
	})
}

// Check will execute the checks specified and return an error on the first occurrence
func (s *SessionChecks) Check(ctx context.Context) error {
	// do the user check first, so the user is set for other checks
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

func (s *SessionChecks) commands(ctx context.Context) (string, []eventstore.Command, error) {
	if len(s.sessionWriteModel.commands) == 0 {
		return "", nil, nil
	}

	// TODO: change config and algorithm (and encrypt instead?)
	// TODO: or maybe just change to generate id and encrypt it
	//token, plain, err := crypto.NewCode(crypto.NewHashGenerator(crypto.GeneratorConfig{
	token, plain, err := crypto.NewCode(crypto.NewEncryptionGenerator(crypto.GeneratorConfig{
		Length:              64,
		Expiry:              0,
		IncludeLowerLetters: true,
		IncludeUpperLetters: true,
		IncludeDigits:       true,
		IncludeSymbols:      false,
		//}, s.userPasswordAlg))
	}, s.sessionAlg))
	if err != nil {
		return "", nil, err
	}
	s.sessionWriteModel.SetToken(ctx, token)
	//if s.sessionWriteModel.State == domain.SessionStateUnspecified {
	//	s.sessionWriteModel.commands = append([]eventstore.Command{session.NewAddedEvent(ctx, s.sessionWriteModel.aggregate)}, s.sessionWriteModel.commands...)
	//}
	return plain, s.sessionWriteModel.commands, nil
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
	checks.sessionWriteModel.Start(ctx)
	return c.updateSession(ctx, checks, metadata)
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
	return c.updateSession(ctx, checks, metadata)
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

// updateSession execute the [SessionChecks] where new events will be created and as well as for metadata (changes)
func (c *Commands) updateSession(ctx context.Context, checks *SessionChecks, metadata map[string][]byte) (set *SessionChanged, err error) {
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
	sessionToken, cmds, err := checks.commands(ctx)
	if err != nil {
		return nil, err
	}
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
	changed := sessionWriteModelToSessionChanged(checks.sessionWriteModel)
	changed.NewToken = sessionToken
	return changed, nil
}

// sessionPermission will check that the provided sessionToken is correct or
// if empty, check that the caller is granted the necessary permission
func (c *Commands) sessionPermission(ctx context.Context, sessionWriteModel *SessionWriteModel, sessionToken, permission string) (err error) {
	if sessionToken == "" {
		return c.checkPermission(ctx, permission, authz.GetCtxData(ctx).OrgID, sessionWriteModel.AggregateID)
	}
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	//err = crypto.CompareHash(sessionWriteModel.Token, []byte(sessionToken), c.userPasswordAlg)
	var token string
	token, err = crypto.DecryptString(sessionWriteModel.Token, c.keyAlgorithm)
	spanPasswordComparison.EndWithError(err)
	if err != nil || token != sessionToken {
		return caos_errs.ThrowPermissionDenied(err, "COMMAND-sGr42", "Invalid token") //TODO: i18n
	}
	return nil
}

type SessionChanged struct {
	*domain.ObjectDetails
	ID       string
	NewToken string
}

func sessionWriteModelToSessionChanged(wm *SessionWriteModel) *SessionChanged {
	return &SessionChanged{
		ObjectDetails: &domain.ObjectDetails{
			Sequence:      wm.ProcessedSequence,
			EventDate:     wm.ChangeDate,
			ResourceOwner: wm.ResourceOwner,
		},
		ID: wm.AggregateID,
	}
}
