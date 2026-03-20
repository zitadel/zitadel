package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CheckIDPIntentType struct {
	ID    string
	Token string
}

type IDPIntentCheckCommand struct {
	CheckIntent *CheckIDPIntentType

	SessionID  string
	InstanceID string
	EncAlgo    crypto.EncryptionAlgorithm

	FetchedUser        User
	IsCheckComplete    bool
	IntentLastVerified time.Time
}

// NewIDPIntentCheckCommand returns an IDPIntentCheckCommand initialized with the input values.
//
// If encryptionAlgo is nil, the default [crypto.EncryptionAlgorithm] will be used
func NewIDPIntentCheckCommand(request *CheckIDPIntentType, sessionID, instanceID string, encryptionAlgo crypto.EncryptionAlgorithm) *IDPIntentCheckCommand {
	idpCheckCommand := &IDPIntentCheckCommand{
		CheckIntent: request,
		SessionID:   sessionID,
		InstanceID:  instanceID,
		EncAlgo:     idpEncryptionAlgo,
	}

	if encryptionAlgo != nil {
		idpCheckCommand.EncAlgo = encryptionAlgo
	}

	return idpCheckCommand
}

// Events implements [Commander].
func (i *IDPIntentCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if i.CheckIntent == nil || !i.IsCheckComplete {
		return nil, nil
	}

	return []eventstore.Command{
		session.NewIntentCheckedEvent(ctx, &session.NewAggregate(i.SessionID, i.InstanceID).Aggregate, i.IntentLastVerified),
		idpintent.NewConsumedEvent(ctx, &idpintent.NewAggregate(i.CheckIntent.ID, "").Aggregate),
	}, nil
}

// Execute implements [Commander].
func (i *IDPIntentCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if i.CheckIntent == nil {
		return nil
	}

	intentRepo := opts.idpIntentRepo
	sessionRepo := opts.sessionRepo

	beginner, ok := opts.DB().(database.Beginner)
	if !ok {
		return zerrors.ThrowInternal(nil, "DOM-jiuIbh", "database doesn't implement database.Beginner")
	}

	tx, txErr := beginner.Begin(ctx, nil)
	if txErr != nil {
		return zerrors.ThrowInternal(txErr, "DOM-Msuwn2", "failed starting transaction")
	}

	defer func() {
		if endErr := tx.End(ctx, txErr); endErr != nil {
			err = endErr
		}
	}()

	removedRows, err := intentRepo.Delete(ctx, opts.DB(), intentRepo.PrimaryKeyCondition(i.InstanceID, i.CheckIntent.ID))
	if err != nil {
		txErr = zerrors.ThrowInternal(err, "DOM-j1s5Eu", "failed deleting IDP intent")
		return txErr
	}
	if removedRows != 1 {
		txErr = zerrors.ThrowInternal(NewRowsReturnedMismatchError(1, removedRows), "DOM-3CBpdB", "unexpected number of rows deleted")
		return txErr
	}

	idpFactor := &SessionFactorIdentityProviderIntent{LastVerifiedAt: time.Now()}
	updateCount, updateErr := sessionRepo.Update(ctx, opts.DB(),
		sessionRepo.PrimaryKeyCondition(i.InstanceID, i.SessionID),
		sessionRepo.SetFactor(idpFactor),
	)
	if updateErr != nil {
		txErr = zerrors.ThrowInternal(updateErr, "DOM-pec0al", "failed updating session")
		return txErr
	}

	if updateCount == 0 {
		txErr = zerrors.ThrowNotFound(nil, "DOM-CopO4e", "session not found")
		return txErr
	}
	if updateCount > 1 {
		txErr = zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-mlbibw", "unexpected number of rows updated")
		return txErr
	}

	i.IntentLastVerified = idpFactor.LastVerifiedAt
	i.IsCheckComplete = true

	return nil
}

// String implements [Commander].
func (i *IDPIntentCheckCommand) String() string {
	return "IDPIntentCheckCommand"
}

// Validate implements [Commander].
func (i *IDPIntentCheckCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if i.CheckIntent == nil {
		return nil
	}

	if i.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-5Y8pb4", "Errors.Missing.SessionID")
	}
	if i.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-Q4YFIq", "Errors.Missing.InstanceID")
	}

	sessionRepo := opts.sessionRepo
	idpIntentRepo := opts.idpIntentRepo
	userRepo := opts.userRepo

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.PrimaryKeyCondition(i.InstanceID, i.SessionID)))
	if err = handleGetError(err, "DOM-EhIgey", "session"); err != nil {
		return err
	}

	if session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-IJcVkV", "Errors.User.UserIDMissing")
	}

	if err := crypto.CheckToken(i.EncAlgo, i.CheckIntent.Token, i.CheckIntent.ID); err != nil {
		return err
	}

	intent, err := idpIntentRepo.Get(ctx, opts.DB(), database.WithCondition(idpIntentRepo.PrimaryKeyCondition(i.InstanceID, i.CheckIntent.ID)))
	if err = handleGetError(err, "DOM-5XkWJV", "intent"); err != nil {
		return err
	}

	if intent.State != IDPIntentStateSucceeded {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-0UKHku", "Errors.Intent.NotSucceeded")
	}

	if intent.ExpiresAt == nil || intent.ExpiresAt.Before(time.Now()) {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-kDR1XK", "Errors.Intent.Expired")
	}

	user, err := userRepo.Get(ctx, opts.DB(), database.WithCondition(userRepo.PrimaryKeyCondition(i.InstanceID, session.UserID)))
	if err = handleGetError(err, "DOM-Vnx2G9", "user"); err != nil {
		return err
	}
	if user.Human == nil {
		return zerrors.ThrowInternal(nil, "DOM-FkX5lZ", "user not human")
	}

	if intentUsr := intent.UserID; intentUsr != "" {
		if intentUsr != session.UserID {
			return zerrors.ThrowPreconditionFailed(nil, "DOM-FLdnLH", "Errors.Intent.OtherUser")
		}
		i.FetchedUser = *user
		return nil
	}

	var matchingLink *IdentityProviderLink
	for _, idpLink := range user.Human.IdentityProviderLinks {
		if idpLink.ProviderID == intent.IDPID {
			matchingLink = idpLink
			break
		}
	}
	if matchingLink == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-XuNkt7", "Errors.Intent.OtherUser")
	}

	i.FetchedUser = *user
	return nil
}

var _ Commander = (*IDPIntentCheckCommand)(nil)
