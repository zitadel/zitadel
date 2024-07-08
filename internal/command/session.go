package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/activity"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SessionCommand func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error)

type SessionCommands struct {
	sessionCommands []SessionCommand

	sessionWriteModel *SessionWriteModel
	intentWriteModel  *IDPIntentWriteModel
	eventstore        *eventstore.Eventstore
	eventCommands     []eventstore.Command

	hasher      *crypto.Hasher
	intentAlg   crypto.EncryptionAlgorithm
	totpAlg     crypto.EncryptionAlgorithm
	otpAlg      crypto.EncryptionAlgorithm
	createCode  encryptedCodeWithDefaultFunc
	createToken func(sessionID string) (id string, token string, err error)
	now         func() time.Time
}

func (c *Commands) NewSessionCommands(cmds []SessionCommand, session *SessionWriteModel) *SessionCommands {
	return &SessionCommands{
		sessionCommands:   cmds,
		sessionWriteModel: session,
		eventstore:        c.eventstore,
		hasher:            c.userPasswordHasher,
		intentAlg:         c.idpConfigEncryption,
		totpAlg:           c.multifactors.OTP.CryptoMFA,
		otpAlg:            c.userEncryption,
		createCode:        c.newEncryptedCodeWithDefault,
		createToken:       c.sessionTokenCreator,
		now:               time.Now,
	}
}

// CheckUser defines a user check to be executed for a session update
func CheckUser(id string, resourceOwner string, preferredLanguage *language.Tag) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
		if cmd.sessionWriteModel.UserID != "" && id != "" && cmd.sessionWriteModel.UserID != id {
			return nil, zerrors.ThrowInvalidArgument(nil, "", "user change not possible")
		}
		return nil, cmd.UserChecked(ctx, id, resourceOwner, cmd.now(), preferredLanguage)
	}
}

// CheckPassword defines a password check to be executed for a session update
func CheckPassword(password string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
		commands, err := checkPassword(ctx, cmd.sessionWriteModel.UserID, password, cmd.eventstore, cmd.hasher, nil)
		if err != nil {
			return commands, err
		}
		cmd.eventCommands = append(cmd.eventCommands, commands...)
		cmd.PasswordChecked(ctx, cmd.now())
		return nil, nil
	}
}

// CheckIntent defines a check for a succeeded intent to be executed for a session update
func CheckIntent(intentID, token string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
		if cmd.sessionWriteModel.UserID == "" {
			return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sfw3r", "Errors.User.UserIDMissing")
		}
		if err := crypto.CheckToken(cmd.intentAlg, token, intentID); err != nil {
			return nil, err
		}
		cmd.intentWriteModel = NewIDPIntentWriteModel(intentID, "")
		err := cmd.eventstore.FilterToQueryReducer(ctx, cmd.intentWriteModel)
		if err != nil {
			return nil, err
		}
		if cmd.intentWriteModel.State != domain.IDPIntentStateSucceeded {
			return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Df4bw", "Errors.Intent.NotSucceeded")
		}
		if cmd.intentWriteModel.UserID != "" {
			if cmd.intentWriteModel.UserID != cmd.sessionWriteModel.UserID {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-O8xk3w", "Errors.Intent.OtherUser")
			}
		} else {
			linkWriteModel := NewUserIDPLinkWriteModel(cmd.sessionWriteModel.UserID, cmd.intentWriteModel.IDPID, cmd.intentWriteModel.IDPUserID, cmd.sessionWriteModel.UserResourceOwner)
			err := cmd.eventstore.FilterToQueryReducer(ctx, linkWriteModel)
			if err != nil {
				return nil, err
			}
			if linkWriteModel.State != domain.UserIDPLinkStateActive {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-O8xk3w", "Errors.Intent.OtherUser")
			}
		}
		cmd.IntentChecked(ctx, cmd.now())
		return nil, nil
	}
}

func CheckTOTP(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) (_ []eventstore.Command, err error) {
		commands, err := checkTOTP(
			ctx,
			cmd.sessionWriteModel.UserID,
			"",
			code,
			cmd.eventstore.FilterToQueryReducer,
			cmd.totpAlg,
			nil,
		)
		if err != nil {
			return commands, err
		}
		cmd.eventCommands = append(cmd.eventCommands, commands...)
		cmd.TOTPChecked(ctx, cmd.now())
		return nil, nil
	}
}

// Exec will execute the commands specified and returns an error on the first occurrence.
// In case of an error there might be specific commands returned, e.g. a failed pw check will have to be stored.
func (s *SessionCommands) Exec(ctx context.Context) ([]eventstore.Command, error) {
	for _, cmd := range s.sessionCommands {
		if cmds, err := cmd(ctx, s); err != nil {
			return cmds, err
		}
	}
	return nil, nil
}

func (s *SessionCommands) Start(ctx context.Context, userAgent *domain.UserAgent) {
	s.eventCommands = append(s.eventCommands, session.NewAddedEvent(ctx, s.sessionWriteModel.aggregate, userAgent))
}

func (s *SessionCommands) UserChecked(ctx context.Context, userID, resourceOwner string, checkedAt time.Time, preferredLanguage *language.Tag) error {
	s.eventCommands = append(s.eventCommands, session.NewUserCheckedEvent(ctx, s.sessionWriteModel.aggregate, userID, resourceOwner, checkedAt, preferredLanguage))
	// set the userID so other checks can use it
	s.sessionWriteModel.UserID = userID
	s.sessionWriteModel.UserResourceOwner = resourceOwner
	return nil
}

func (s *SessionCommands) PasswordChecked(ctx context.Context, checkedAt time.Time) {
	s.eventCommands = append(s.eventCommands, session.NewPasswordCheckedEvent(ctx, s.sessionWriteModel.aggregate, checkedAt))
}

func (s *SessionCommands) IntentChecked(ctx context.Context, checkedAt time.Time) {
	s.eventCommands = append(s.eventCommands, session.NewIntentCheckedEvent(ctx, s.sessionWriteModel.aggregate, checkedAt))
}

func (s *SessionCommands) WebAuthNChallenged(ctx context.Context, challenge string, allowedCrentialIDs [][]byte, userVerification domain.UserVerificationRequirement, rpid string) {
	s.eventCommands = append(s.eventCommands, session.NewWebAuthNChallengedEvent(ctx, s.sessionWriteModel.aggregate, challenge, allowedCrentialIDs, userVerification, rpid))
}

func (s *SessionCommands) WebAuthNChecked(ctx context.Context, checkedAt time.Time, tokenID string, signCount uint32, userVerified bool) {
	s.eventCommands = append(s.eventCommands,
		session.NewWebAuthNCheckedEvent(ctx, s.sessionWriteModel.aggregate, checkedAt, userVerified),
	)
	if s.sessionWriteModel.WebAuthNChallenge.UserVerification == domain.UserVerificationRequirementRequired {
		s.eventCommands = append(s.eventCommands,
			user.NewHumanPasswordlessSignCountChangedEvent(ctx, s.sessionWriteModel.aggregate, tokenID, signCount),
		)
	} else {
		s.eventCommands = append(s.eventCommands,
			user.NewHumanU2FSignCountChangedEvent(ctx, s.sessionWriteModel.aggregate, tokenID, signCount),
		)
	}
}

func (s *SessionCommands) TOTPChecked(ctx context.Context, checkedAt time.Time) {
	s.eventCommands = append(s.eventCommands, session.NewTOTPCheckedEvent(ctx, s.sessionWriteModel.aggregate, checkedAt))
}

func (s *SessionCommands) OTPSMSChallenged(ctx context.Context, code *crypto.CryptoValue, expiry time.Duration, returnCode bool) {
	s.eventCommands = append(s.eventCommands, session.NewOTPSMSChallengedEvent(ctx, s.sessionWriteModel.aggregate, code, expiry, returnCode))
}

func (s *SessionCommands) OTPSMSChecked(ctx context.Context, checkedAt time.Time) {
	s.eventCommands = append(s.eventCommands, session.NewOTPSMSCheckedEvent(ctx, s.sessionWriteModel.aggregate, checkedAt))
}

func (s *SessionCommands) OTPEmailChallenged(ctx context.Context, code *crypto.CryptoValue, expiry time.Duration, returnCode bool, urlTmpl string) {
	s.eventCommands = append(s.eventCommands, session.NewOTPEmailChallengedEvent(ctx, s.sessionWriteModel.aggregate, code, expiry, returnCode, urlTmpl))
}

func (s *SessionCommands) OTPEmailChecked(ctx context.Context, checkedAt time.Time) {
	s.eventCommands = append(s.eventCommands, session.NewOTPEmailCheckedEvent(ctx, s.sessionWriteModel.aggregate, checkedAt))
}

func (s *SessionCommands) SetToken(ctx context.Context, tokenID string) {
	// trigger activity log for session for user
	activity.Trigger(ctx, s.sessionWriteModel.UserResourceOwner, s.sessionWriteModel.UserID, activity.SessionAPI, s.eventstore.FilterToQueryReducer)
	s.eventCommands = append(s.eventCommands, session.NewTokenSetEvent(ctx, s.sessionWriteModel.aggregate, tokenID))
}

func (s *SessionCommands) ChangeMetadata(ctx context.Context, metadata map[string][]byte) {
	var changed bool
	for key, value := range metadata {
		currentValue, exists := s.sessionWriteModel.Metadata[key]

		if len(value) != 0 {
			// if a value is provided, and it's not equal, change it
			if !bytes.Equal(currentValue, value) {
				s.sessionWriteModel.Metadata[key] = value
				changed = true
			}
		} else {
			// if there's no / an empty value, we only need to remove it on existing entries
			if exists {
				delete(s.sessionWriteModel.Metadata, key)
				changed = true
			}
		}
	}
	if changed {
		s.eventCommands = append(s.eventCommands, session.NewMetadataSetEvent(ctx, s.sessionWriteModel.aggregate, s.sessionWriteModel.Metadata))
	}
}

func (s *SessionCommands) SetLifetime(ctx context.Context, lifetime time.Duration) error {
	if lifetime < 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-asEG4", "Errors.Session.PositiveLifetime")
	}
	if lifetime == 0 {
		return nil
	}
	s.eventCommands = append(s.eventCommands, session.NewLifetimeSetEvent(ctx, s.sessionWriteModel.aggregate, lifetime))
	return nil
}

func (s *SessionCommands) gethumanWriteModel(ctx context.Context) (*HumanWriteModel, error) {
	if s.sessionWriteModel.UserID == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-eeR2e", "Errors.User.UserIDMissing")
	}
	humanWriteModel := NewHumanWriteModel(s.sessionWriteModel.UserID, s.sessionWriteModel.UserResourceOwner)
	err := s.eventstore.FilterToQueryReducer(ctx, humanWriteModel)
	if err != nil {
		return nil, err
	}
	if humanWriteModel.UserState != domain.UserStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Df4b3", "Errors.User.NotFound")
	}
	return humanWriteModel, nil
}

func (s *SessionCommands) commands(ctx context.Context) (string, []eventstore.Command, error) {
	if len(s.eventCommands) == 0 {
		return "", nil, nil
	}

	tokenID, token, err := s.createToken(s.sessionWriteModel.AggregateID)
	if err != nil {
		return "", nil, err
	}
	s.SetToken(ctx, tokenID)
	return token, s.eventCommands, nil
}

func (c *Commands) CreateSession(ctx context.Context, cmds []SessionCommand, metadata map[string][]byte, userAgent *domain.UserAgent, lifetime time.Duration) (set *SessionChanged, err error) {
	sessionID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}
	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetInstance(ctx).InstanceID())
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return nil, err
	}
	cmd := c.NewSessionCommands(cmds, sessionWriteModel)
	cmd.Start(ctx, userAgent)
	return c.updateSession(ctx, cmd, metadata, lifetime)
}

func (c *Commands) UpdateSession(ctx context.Context, sessionID string, cmds []SessionCommand, metadata map[string][]byte, lifetime time.Duration) (set *SessionChanged, err error) {
	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetInstance(ctx).InstanceID())
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return nil, err
	}
	cmd := c.NewSessionCommands(cmds, sessionWriteModel)
	return c.updateSession(ctx, cmd, metadata, lifetime)
}

func (c *Commands) TerminateSession(ctx context.Context, sessionID string, sessionToken string) (*domain.ObjectDetails, error) {
	return c.terminateSession(ctx, sessionID, sessionToken, true)
}

func (c *Commands) TerminateSessionWithoutTokenCheck(ctx context.Context, sessionID string) (*domain.ObjectDetails, error) {
	return c.terminateSession(ctx, sessionID, "", false)
}

func (c *Commands) terminateSession(ctx context.Context, sessionID, sessionToken string, mustCheckToken bool) (*domain.ObjectDetails, error) {
	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetInstance(ctx).InstanceID())
	if err := c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel); err != nil {
		return nil, err
	}
	if mustCheckToken {
		if err := c.checkSessionTerminationPermission(ctx, sessionWriteModel, sessionToken); err != nil {
			return nil, err
		}
	}
	if sessionWriteModel.CheckIsActive() != nil {
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
func (c *Commands) updateSession(ctx context.Context, checks *SessionCommands, metadata map[string][]byte, lifetime time.Duration) (set *SessionChanged, err error) {
	if err = checks.sessionWriteModel.CheckNotInvalidated(); err != nil {
		return nil, err
	}
	if cmds, err := checks.Exec(ctx); err != nil {
		if len(cmds) > 0 {
			_, pushErr := c.eventstore.Push(ctx, cmds...)
			logging.OnError(pushErr).Error("unable to store check failures")
		}
		return nil, err
	}
	checks.ChangeMetadata(ctx, metadata)
	err = checks.SetLifetime(ctx, lifetime)
	if err != nil {
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

// checkSessionTerminationPermission will check that the provided sessionToken is correct or
// if empty, check that the caller is either terminating the own session or
// is granted the "session.delete" permission on the resource owner of the authenticated user.
func (c *Commands) checkSessionTerminationPermission(ctx context.Context, model *SessionWriteModel, token string) error {
	if token != "" {
		return c.sessionTokenVerifier(ctx, token, model.AggregateID, model.TokenID)
	}
	if model.UserID != "" && model.UserID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	userResourceOwner, err := c.sessionUserResourceOwner(ctx, model)
	if err != nil {
		return err
	}
	return c.checkPermission(ctx, domain.PermissionSessionDelete, userResourceOwner, model.UserID)
}

// sessionUserResourceOwner will return the resourceOwner of the session form the [SessionWriteModel] or by additionally calling the eventstore,
// because before 2.42.0, the resourceOwner of a session used to be the organisation of the creator.
// Further the (checked) users organisation id was not stored.
// To be able to check the permission, we need to get the user's resourceOwner in this case.
func (c *Commands) sessionUserResourceOwner(ctx context.Context, model *SessionWriteModel) (string, error) {
	if model.UserID == "" || model.UserResourceOwner != "" {
		return model.UserResourceOwner, nil
	}
	r := NewResourceOwnerModel(authz.GetInstance(ctx).InstanceID(), user.AggregateType, model.UserID)
	err := c.eventstore.FilterToQueryReducer(ctx, r)
	if err != nil {
		return "", err
	}
	return r.resourceOwner, nil
}

func sessionTokenCreator(sessionAlg crypto.EncryptionAlgorithm) func(sessionID string) (id string, token string, err error) {
	return func(sessionID string) (id string, token string, err error) {
		id, err = id_generator.Next()
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
