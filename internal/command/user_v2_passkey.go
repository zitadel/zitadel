package command

import (
	"context"
	"io"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func (c *Commands) AddUserPasskeyCode(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	details, err := c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, "", false)
	return details.ObjectDetails, err
}

func (c *Commands) AddUserPasskeyCodeURLTemplate(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm, urlTmpl string) (*domain.ObjectDetails, error) {
	if err := domain.RenderPasskeyURLTemplate(io.Discard, urlTmpl, userID, resourceOwner, "codeID", "code"); err != nil {
		return nil, err
	}
	details, err := c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, urlTmpl, false)
	return details.ObjectDetails, err
}

func (c *Commands) AddUserPasskeyCodeReturn(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm) (*domain.PasskeyCodeDetails, error) {
	return c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, "", true)
}

func (c *Commands) addUserPasskeyCode(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm, urlTmpl string, returnCode bool) (*domain.PasskeyCodeDetails, error) {
	config, err := secretGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordlessInitCode)
	if err != nil {
		return nil, err
	}
	codeID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	gen := crypto.NewEncryptionGenerator(*config, alg)
	cryptoCode, code, err := crypto.NewCode(gen)
	if err != nil {
		return nil, err
	}

	var (
		wm    = NewHumanPasswordlessInitCodeWriteModel(userID, codeID, resourceOwner)
		aggr  = UserAggregateFromWriteModel(&wm.WriteModel)
		event = user.NewHumanPasswordlessInitCodeRequestedEvent(ctx, aggr, codeID, cryptoCode, gen.Expiry(), urlTmpl, returnCode)
	)

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(wm, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return &domain.PasskeyCodeDetails{
		ObjectDetails: writeModelToObjectDetails(&wm.WriteModel),
		CodeID:        codeID,
		Code:          code,
	}, nil
}
