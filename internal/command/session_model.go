package command

import (
	"bytes"
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
)

type PasskeyChallengeModel struct {
	Challenge          string
	AllowedCrentialIDs [][]byte
	UserVerification   domain.UserVerificationRequirement
}

func (p *PasskeyChallengeModel) WebAuthNLogin(human *domain.Human, credentialAssertionData []byte) (*domain.WebAuthNLogin, error) {
	if p == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Ioqu5", "Errors.Session.Passkey.NoChallenge")
	}
	return &domain.WebAuthNLogin{
		ObjectRoot:              human.ObjectRoot,
		CredentialAssertionData: credentialAssertionData,
		Challenge:               p.Challenge,
		AllowedCredentialIDs:    p.AllowedCrentialIDs,
		UserVerification:        p.UserVerification,
	}, nil
}

type SessionWriteModel struct {
	eventstore.WriteModel

	TokenID           string
	UserID            string
	UserCheckedAt     time.Time
	PasswordCheckedAt time.Time
	PasskeyCheckedAt  time.Time
	Metadata          map[string][]byte
	State             domain.SessionState

	PasskeyChallenge *PasskeyChallengeModel

	commands  []eventstore.Command
	aggregate *eventstore.Aggregate
}

func NewSessionWriteModel(sessionID string, resourceOwner string) *SessionWriteModel {
	return &SessionWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   sessionID,
			ResourceOwner: resourceOwner,
		},
		Metadata:  make(map[string][]byte),
		aggregate: &session.NewAggregate(sessionID, resourceOwner).Aggregate,
	}
}

func (wm *SessionWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *session.AddedEvent:
			wm.reduceAdded(e)
		case *session.UserCheckedEvent:
			wm.reduceUserChecked(e)
		case *session.PasswordCheckedEvent:
			wm.reducePasswordChecked(e)
		case *session.PasskeyChallengedEvent:
			wm.reducePasskeyChallenged(e)
		case *session.PasskeyCheckedEvent:
			wm.reducePasskeyChecked(e)
		case *session.TokenSetEvent:
			wm.reduceTokenSet(e)
		case *session.TerminateEvent:
			wm.reduceTerminate()
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *SessionWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(session.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			session.AddedType,
			session.UserCheckedType,
			session.PasswordCheckedType,
			session.PasskeyChallengedType,
			session.PasskeyCheckedType,
			session.TokenSetType,
			session.MetadataSetType,
			session.TerminateType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *SessionWriteModel) reduceAdded(e *session.AddedEvent) {
	wm.State = domain.SessionStateActive
}

func (wm *SessionWriteModel) reduceUserChecked(e *session.UserCheckedEvent) {
	wm.UserID = e.UserID
	wm.UserCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reducePasswordChecked(e *session.PasswordCheckedEvent) {
	wm.PasswordCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reducePasskeyChallenged(e *session.PasskeyChallengedEvent) {
	wm.PasskeyChallenge = &PasskeyChallengeModel{
		Challenge:          e.Challenge,
		AllowedCrentialIDs: e.AllowedCrentialIDs,
		UserVerification:   e.UserVerification,
	}
}

func (wm *SessionWriteModel) reducePasskeyChecked(e *session.PasskeyCheckedEvent) {
	wm.PasskeyChallenge = nil
	wm.PasskeyCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reduceTokenSet(e *session.TokenSetEvent) {
	wm.TokenID = e.TokenID
}

func (wm *SessionWriteModel) reduceTerminate() {
	wm.State = domain.SessionStateTerminated
}

func (wm *SessionWriteModel) Start(ctx context.Context) {
	wm.commands = append(wm.commands, session.NewAddedEvent(ctx, wm.aggregate))
}

func (wm *SessionWriteModel) UserChecked(ctx context.Context, userID string, checkedAt time.Time) error {
	wm.commands = append(wm.commands, session.NewUserCheckedEvent(ctx, wm.aggregate, userID, checkedAt))
	// set the userID so other checks can use it
	wm.UserID = userID
	return nil
}

func (wm *SessionWriteModel) PasswordChecked(ctx context.Context, checkedAt time.Time) {
	wm.commands = append(wm.commands, session.NewPasswordCheckedEvent(ctx, wm.aggregate, checkedAt))
}

func (wm *SessionWriteModel) PasskeyChallenged(ctx context.Context, challenge string, allowedCrentialIDs [][]byte, userVerification domain.UserVerificationRequirement) {
	wm.commands = append(wm.commands, session.NewPasskeyChallengedEvent(ctx, wm.aggregate, challenge, allowedCrentialIDs, userVerification))
}

func (wm *SessionWriteModel) PasskeyChecked(ctx context.Context, checkedAt time.Time, tokenID string, signCount uint32) {
	wm.commands = append(wm.commands,
		session.NewPasskeyCheckedEvent(ctx, wm.aggregate, checkedAt),
		usr_repo.NewHumanPasswordlessSignCountChangedEvent(ctx, wm.aggregate, tokenID, signCount),
	)
}

func (wm *SessionWriteModel) SetToken(ctx context.Context, tokenID string) {
	wm.commands = append(wm.commands, session.NewTokenSetEvent(ctx, wm.aggregate, tokenID))
}

func (wm *SessionWriteModel) ChangeMetadata(ctx context.Context, metadata map[string][]byte) {
	var changed bool
	for key, value := range metadata {
		currentValue, exists := wm.Metadata[key]

		if len(value) != 0 {
			// if a value is provided, and it's not equal, change it
			if !bytes.Equal(currentValue, value) {
				wm.Metadata[key] = value
				changed = true
			}
		} else {
			// if there's no / an empty value, we only need to remove it on existing entries
			if exists {
				delete(wm.Metadata, key)
				changed = true
			}
		}
	}
	if changed {
		wm.commands = append(wm.commands, session.NewMetadataSetEvent(ctx, wm.aggregate, wm.Metadata))
	}
}
