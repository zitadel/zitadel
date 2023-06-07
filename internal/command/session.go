package command

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type SessionCommand func(ctx context.Context, cmd *SessionCommands) error

type SessionCommands struct {
	cmds []SessionCommand

	sessionWriteModel  *SessionWriteModel
	passwordWriteModel *HumanPasswordWriteModel
	eventstore         *eventstore.Eventstore
	userPasswordAlg    crypto.HashAlgorithm
	createToken        func(sessionID string) (id string, token string, err error)
	now                func() time.Time
}

func (c *Commands) NewSessionCommands(cmds []SessionCommand, session *SessionWriteModel) *SessionCommands {
	return &SessionCommands{
		cmds:              cmds,
		sessionWriteModel: session,
		eventstore:        c.eventstore,
		userPasswordAlg:   c.userPasswordAlg,
		createToken:       c.sessionTokenCreator,
		now:               time.Now,
	}
}

// CheckUser defines a user check to be executed for a session update
func CheckUser(id string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		if cmd.sessionWriteModel.UserID != "" && id != "" && cmd.sessionWriteModel.UserID != id {
			return caos_errs.ThrowInvalidArgument(nil, "", "user change not possible")
		}
		return cmd.sessionWriteModel.UserChecked(ctx, id, cmd.now())
	}
}

// CheckPassword defines a password check to be executed for a session update
func CheckPassword(password string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		if cmd.sessionWriteModel.UserID == "" {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sfw3f", "Errors.User.UserIDMissing")
		}
		cmd.passwordWriteModel = NewHumanPasswordWriteModel(cmd.sessionWriteModel.UserID, "")
		err := cmd.eventstore.FilterToQueryReducer(ctx, cmd.passwordWriteModel)
		if err != nil {
			return err
		}
		if cmd.passwordWriteModel.UserState == domain.UserStateUnspecified || cmd.passwordWriteModel.UserState == domain.UserStateDeleted {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Df4b3", "Errors.User.NotFound")
		}

		if cmd.passwordWriteModel.Secret == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-WEf3t", "Errors.User.Password.NotSet")
		}
		ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
		err = crypto.CompareHash(cmd.passwordWriteModel.Secret, []byte(password), cmd.userPasswordAlg)
		spanPasswordComparison.EndWithError(err)
		if err != nil {
			//TODO: maybe we want to reset the session in the future https://github.com/zitadel/zitadel/issues/5807
			return caos_errs.ThrowInvalidArgument(err, "COMMAND-SAF3g", "Errors.User.Password.Invalid")
		}
		cmd.sessionWriteModel.PasswordChecked(ctx, cmd.now())
		return nil
	}
}

// Exec will execute the commands specified and returns an error on the first occurrence
func (s *SessionCommands) Exec(ctx context.Context) error {
	for _, cmd := range s.cmds {
		if err := cmd(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionCommands) gethumanWriteModel(ctx context.Context) (*HumanWriteModel, error) {
	if s.sessionWriteModel.UserID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-eeR2e", "Errors.User.UserIDMissing")
	}
	humanWriteModel := NewHumanWriteModel(s.sessionWriteModel.UserID, s.sessionWriteModel.ResourceOwner)
	err := s.eventstore.FilterToQueryReducer(ctx, humanWriteModel)
	if err != nil {
		return nil, err
	}
	if humanWriteModel.UserState != domain.UserStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Df4b3", "Errors.ie4Ai.NotFound")
	}
	return humanWriteModel, nil
}

func (s *SessionCommands) commands(ctx context.Context) (string, []eventstore.Command, error) {
	if len(s.sessionWriteModel.commands) == 0 {
		return "", nil, nil
	}

	tokenID, token, err := s.createToken(s.sessionWriteModel.AggregateID)
	if err != nil {
		return "", nil, err
	}
	s.sessionWriteModel.SetToken(ctx, tokenID)
	return token, s.sessionWriteModel.commands, nil
}

func (c *Commands) CreateSession(ctx context.Context, cmds []SessionCommand, metadata map[string][]byte) (set *SessionChanged, err error) {
	sessionID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetCtxData(ctx).OrgID)
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return nil, err
	}
	cmd := c.NewSessionCommands(cmds, sessionWriteModel)
	cmd.sessionWriteModel.Start(ctx)
	return c.updateSession(ctx, cmd, metadata)
}

func (c *Commands) UpdateSession(ctx context.Context, sessionID, sessionToken string, cmds []SessionCommand, metadata map[string][]byte) (set *SessionChanged, err error) {
	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetCtxData(ctx).OrgID)
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return nil, err
	}
	if err := c.sessionPermission(ctx, sessionWriteModel, sessionToken, domain.PermissionSessionWrite); err != nil {
		return nil, err
	}
	cmd := c.NewSessionCommands(cmds, sessionWriteModel)
	return c.updateSession(ctx, cmd, metadata)
}

func (c *Commands) TerminateSession(ctx context.Context, sessionID, sessionToken string) (*domain.ObjectDetails, error) {
	sessionWriteModel := NewSessionWriteModel(sessionID, "")
	if err := c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel); err != nil {
		return nil, err
	}
	if err := c.sessionPermission(ctx, sessionWriteModel, sessionToken, domain.PermissionSessionDelete); err != nil {
		return nil, err
	}
	if sessionWriteModel.State != domain.SessionStateActive {
		return writeModelToObjectDetails(&sessionWriteModel.WriteModel), nil
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
	return writeModelToObjectDetails(&sessionWriteModel.WriteModel), nil
}

// updateSession execute the [SessionCommands] where new events will be created and as well as for metadata (changes)
func (c *Commands) updateSession(ctx context.Context, checks *SessionCommands, metadata map[string][]byte) (set *SessionChanged, err error) {
	if checks.sessionWriteModel.State == domain.SessionStateTerminated {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMAND-SAjeh", "Errors.Session.Terminated")
	}
	if err := checks.Exec(ctx); err != nil {
		// TODO: how to handle failed checks (e.g. pw wrong) https://github.com/zitadel/zitadel/issues/5807
		return nil, err
	}
	checks.sessionWriteModel.ChangeMetadata(ctx, metadata)
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
	return c.sessionTokenVerifier(ctx, sessionToken, sessionWriteModel.AggregateID, sessionWriteModel.TokenID)
}

func sessionTokenCreator(idGenerator id.Generator, sessionAlg crypto.EncryptionAlgorithm) func(sessionID string) (id string, token string, err error) {
	return func(sessionID string) (id string, token string, err error) {
		id, err = idGenerator.Next()
		if err != nil {
			return "", "", err
		}
		encrypted, err := sessionAlg.Encrypt([]byte(fmt.Sprintf(authz.SessionTokenFormat, sessionID, id)))
		if err != nil {
			return "", "", err
		}
		return id, base64.RawURLEncoding.EncodeToString(encrypted), nil
	}
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
