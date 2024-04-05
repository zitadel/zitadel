package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddAccessAndRefreshToken(
	ctx context.Context,
	orgID,
	agentID,
	clientID,
	userID,
	refreshToken string,
	audience,
	scopes,
	authMethodsReferences []string,
	accessLifetime,
	refreshIdleExpiration,
	refreshExpiration time.Duration,
	authTime time.Time,
	reason domain.TokenReason,
	actor *domain.TokenActor,
) (accessToken *domain.Token, newRefreshToken string, err error) {
	if refreshToken == "" {
		return c.AddNewRefreshTokenAndAccessToken(ctx, userID, orgID, agentID, clientID, audience, scopes, authMethodsReferences, refreshExpiration, accessLifetime, refreshIdleExpiration, authTime, reason, actor)
	}
	return c.RenewRefreshTokenAndAccessToken(ctx, userID, orgID, refreshToken, agentID, clientID, audience, scopes, refreshIdleExpiration, accessLifetime, actor)
}

func (c *Commands) AddNewRefreshTokenAndAccessToken(
	ctx context.Context,
	userID,
	orgID,
	agentID,
	clientID string,
	audience,
	scopes,
	authMethodsReferences []string,
	refreshExpiration,
	accessLifetime,
	refreshIdleExpiration time.Duration,
	authTime time.Time,
	reason domain.TokenReason,
	actor *domain.TokenActor,
) (accessToken *domain.Token, newRefreshToken string, err error) {
	if userID == "" || clientID == "" {
		return nil, "", zerrors.ThrowInvalidArgument(nil, "COMMAND-adg4r", "Errors.IDMissing")
	}
	userWriteModel := NewUserWriteModel(userID, orgID)
	refreshTokenID, err := c.idGenerator.Next()
	if err != nil {
		return nil, "", err
	}
	cmds, accessToken, err := c.addUserToken(ctx, userWriteModel, agentID, clientID, refreshTokenID, audience, scopes, authMethodsReferences, accessLifetime, authTime, reason, actor)
	if err != nil {
		return nil, "", err
	}
	refreshTokenEvent, newRefreshToken, err := c.addRefreshToken(ctx, accessToken, authMethodsReferences, authTime, refreshIdleExpiration, refreshExpiration, actor)
	if err != nil {
		return nil, "", err
	}
	cmds = append(cmds, refreshTokenEvent)
	_, err = c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, "", err
	}
	return accessToken, newRefreshToken, nil
}

func (c *Commands) RenewRefreshTokenAndAccessToken(
	ctx context.Context,
	userID,
	orgID,
	refreshToken,
	agentID,
	clientID string,
	audience,
	scopes []string,
	idleExpiration,
	accessLifetime time.Duration,
	actor *domain.TokenActor,
) (accessToken *domain.Token, newRefreshToken string, err error) {
	renewed, err := c.renewRefreshToken(ctx, userID, orgID, refreshToken, idleExpiration)
	if err != nil {
		return nil, "", err
	}
	userWriteModel := NewUserWriteModel(userID, orgID)
	cmds, accessToken, err := c.addUserToken(ctx, userWriteModel, agentID, clientID, renewed.tokenID, audience, scopes, renewed.authMethodsReferences, accessLifetime, renewed.authTime, domain.TokenReasonRefresh, actor)
	if err != nil {
		return nil, "", err
	}
	_, err = c.eventstore.Push(ctx, append(cmds, renewed.event)...)
	if err != nil {
		return nil, "", err
	}
	return accessToken, renewed.token, nil
}

func (c *Commands) RevokeRefreshToken(ctx context.Context, userID, orgID, tokenID string) (*domain.ObjectDetails, error) {
	removeEvent, refreshTokenWriteModel, err := c.removeRefreshToken(ctx, userID, orgID, tokenID)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(refreshTokenWriteModel, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&refreshTokenWriteModel.WriteModel), nil
}

func (c *Commands) RevokeRefreshTokens(ctx context.Context, userID, orgID string, tokenIDs []string) (err error) {
	if len(tokenIDs) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Gfj42", "Errors.IDMissing")
	}
	events := make([]eventstore.Command, len(tokenIDs))
	for i, tokenID := range tokenIDs {
		event, _, err := c.removeRefreshToken(ctx, userID, orgID, tokenID)
		if err != nil {
			return err
		}
		events[i] = event
	}
	_, err = c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) addRefreshToken(ctx context.Context, accessToken *domain.Token, authMethodsReferences []string, authTime time.Time, idleExpiration, expiration time.Duration, actor *domain.TokenActor) (*user.HumanRefreshTokenAddedEvent, string, error) {
	refreshToken, err := domain.NewRefreshToken(accessToken.AggregateID, accessToken.RefreshTokenID, c.keyAlgorithm)
	if err != nil {
		return nil, "", err
	}
	refreshTokenWriteModel := NewHumanRefreshTokenWriteModel(accessToken.AggregateID, accessToken.ResourceOwner, accessToken.RefreshTokenID)
	userAgg := UserAggregateFromWriteModel(&refreshTokenWriteModel.WriteModel)
	return user.NewHumanRefreshTokenAddedEvent(ctx, userAgg, accessToken.RefreshTokenID, accessToken.ApplicationID, accessToken.UserAgentID,
			accessToken.PreferredLanguage, accessToken.Audience, accessToken.Scopes, authMethodsReferences, authTime, idleExpiration, expiration, actor),
		refreshToken, nil
}

type renewedRefreshToken struct {
	event                 *user.HumanRefreshTokenRenewedEvent
	authTime              time.Time
	authMethodsReferences []string
	tokenID               string
	token                 string
}

func (c *Commands) renewRefreshToken(ctx context.Context, userID, orgID, refreshToken string, idleExpiration time.Duration) (*renewedRefreshToken, error) {
	if refreshToken == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-DHrr3", "Errors.IDMissing")
	}

	tokenUserID, tokenID, token, err := domain.FromRefreshToken(refreshToken, c.keyAlgorithm)
	if err != nil {
		return nil, err
	}
	if tokenUserID != userID {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Ht2g2", "Errors.User.RefreshToken.Invalid")
	}
	refreshTokenWriteModel := NewHumanRefreshTokenWriteModel(userID, orgID, tokenID)
	err = c.eventstore.FilterToQueryReducer(ctx, refreshTokenWriteModel)
	if err != nil {
		return nil, err
	}
	if refreshTokenWriteModel.UserState != domain.UserStateActive {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-BHnhs", "Errors.User.RefreshToken.Invalid")
	}
	if refreshTokenWriteModel.RefreshToken != token ||
		refreshTokenWriteModel.IdleExpiration.Before(time.Now()) ||
		refreshTokenWriteModel.Expiration.Before(time.Now()) {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Vr43e", "Errors.User.RefreshToken.Invalid")
	}

	newToken, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := domain.RefreshToken(userID, tokenID, newToken, c.keyAlgorithm)
	if err != nil {
		return nil, err
	}
	userAgg := UserAggregateFromWriteModel(&refreshTokenWriteModel.WriteModel)
	return &renewedRefreshToken{
		event:                 user.NewHumanRefreshTokenRenewedEvent(ctx, userAgg, tokenID, newToken, idleExpiration),
		authTime:              refreshTokenWriteModel.AuthTime,
		authMethodsReferences: refreshTokenWriteModel.AuthMethodsReferences,
		tokenID:               tokenID,
		token:                 newRefreshToken,
	}, nil
}

func (c *Commands) removeRefreshToken(ctx context.Context, userID, orgID, tokenID string) (*user.HumanRefreshTokenRemovedEvent, *HumanRefreshTokenWriteModel, error) {
	if userID == "" || orgID == "" || tokenID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-GVDgf", "Errors.IDMissing")
	}
	refreshTokenWriteModel := NewHumanRefreshTokenWriteModel(userID, orgID, tokenID)
	err := c.eventstore.FilterToQueryReducer(ctx, refreshTokenWriteModel)
	if err != nil {
		return nil, nil, err
	}
	if refreshTokenWriteModel.UserState != domain.UserStateActive {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-BHt2w", "Errors.User.RefreshToken.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&refreshTokenWriteModel.WriteModel)
	return user.NewHumanRefreshTokenRemovedEvent(ctx, userAgg, tokenID), refreshTokenWriteModel, nil
}
