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
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type IDPIntentCheckCommand struct {
	CheckIntent *session_grpc.CheckIDPIntent

	sessionID  string
	instanceID string
	encAlgo    crypto.EncryptionAlgorithm

	// fetchedUser     User // todo: commenting it out to please the linter
	isCheckComplete bool
}

// NewIDPIntentCheckCommand returns an IDPIntentCheckCommand initialized with the input values.
//
// If encryptionAlgo is nil, the default [crypto.EncryptionAlgorithm] will be used
func NewIDPIntentCheckCommand(sessionID, instanceID string, request *session_grpc.CheckIDPIntent, encryptionAlgo crypto.EncryptionAlgorithm) *IDPIntentCheckCommand {
	idpCheckCommand := &IDPIntentCheckCommand{
		CheckIntent: request,
		sessionID:   sessionID,
		instanceID:  instanceID,
	}

	idpCheckCommand.encAlgo = idpEncryptionAlgo
	if encryptionAlgo != nil {
		idpCheckCommand.encAlgo = encryptionAlgo
	}

	return idpCheckCommand
}

// RequiresTransaction implements [Transactional].
func (i *IDPIntentCheckCommand) RequiresTransaction() {}

// Events implements [Commander].
func (i *IDPIntentCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if i.CheckIntent == nil || !i.isCheckComplete {
		return nil, nil
	}

	return []eventstore.Command{
		session.NewIntentCheckedEvent(ctx, &session.NewAggregate(i.sessionID, i.instanceID).Aggregate, time.Now()),
		idpintent.NewConsumedEvent(ctx, &idpintent.NewAggregate(i.CheckIntent.GetIdpIntentId(), "").Aggregate),
	}, nil
}

// Execute implements [Commander].
func (i *IDPIntentCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if i.CheckIntent == nil {
		return nil
	}

	// TODO(IAM-Marco): Implement when intent repo is available
	// Should implement this idpintent.NewConsumedEvent(ctx, IDPIntentAggregateFromWriteModel(&s.intentWriteModel.WriteModel))
	// intentRepo := opts.intentRepo
	// rowCount, err := intentRepo.Update(ctx, opts.DB(), ???)

	i.isCheckComplete = true

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

	sessionRepo := opts.sessionRepo
	// TODO(IAM-Marco): Uncomment when IDP intents are available
	// idpIntentRepo := opts.idpIntentRepo
	// userRepo := opts.userRepo.LoadIdentityProviderLinks()

	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(i.sessionID)))
	if err := handleGetError(err, "DOM-EhIgey", "session"); err != nil {
		return nil
	}

	if session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-IJcVkV", "Errors.User.UserIDMissing")
	}

	if err := crypto.CheckToken(i.encAlgo, i.CheckIntent.GetIdpIntentToken(), i.CheckIntent.GetIdpIntentId()); err != nil {
		return err
	}

	// TODO(IAM-Marco): Uncomment when IDP intents are available
	// intent, err := idpIntentRepo.Get(ctx, opts.DB(), database.WithCondition(idpIntentRepo.IDCondition(i.CheckIntent.GetIdpIntentId())))
	// if err := handleGetError(err, "DOM-5XkWJV", "intent"); err != nil {
	// 	return err
	// }

	// if intent.State != IntentStateSucceeded {
	// 	return zerrors.ThrowPreconditionFailed(nil, "DOM-0UKHku", "Errors.Intent.NotSucceeded")
	// }

	// if intent.ExpiresAt.Before(time.Now()) {
	// 	zerrors.ThrowPreconditionFailed(nil, "DOM-kDR1XK", "Errors.Intent.Expired")
	// }

	// if intentUsr := intent.UserID; intentUsr != "" {
	// 	if intentUsr != session.UserID {
	// 		return zerrors.ThrowPreconditionFailed(nil, "DOM-FLdnLH", "Errors.Intent.OtherUser")
	// 	}
	// 	return nil
	// }

	// user, err := userRepo.Get(ctx, opts.DB(), database.WithCondition(userRepo.IDCondition(session.UserID)))
	// if err := handleGetError(err, "DOM-Vnx2G9", "intent"); err != nil {
	// 	return err
	// }
	// if user.Human == nil {
	// 	return zerrors.ThrowInternal(nil, "DOM-FkX5lZ", "user not human")
	// }

	// i.fetchedUser = *user

	// var matchingLink *IdentityProviderLink
	// for _, idpLink := range user.Human.IdentityProviderLinks {
	// 	if idpLink.ProviderID == intent.ProviderID {
	// 		matchingLink = idpLink
	// 		break
	// 	}
	// }
	// if matchingLink == nil {
	// 	return zerrors.ThrowPreconditionFailed(nil, "DOM-XuNkt7", "Errors.Intent.OtherUser")
	// }

	return nil
}

var _ Commander = (*IDPIntentCheckCommand)(nil)
var _ Transactional = (*IDPIntentCheckCommand)(nil)
