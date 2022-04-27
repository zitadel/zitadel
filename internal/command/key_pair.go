package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	keypair "github.com/zitadel/zitadel/internal/repository/keypair"
)

const (
	oidcUser = "OIDC"
)

func (c *Commands) GenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	ctx = setOIDCCtx(ctx)
	privateCrypto, publicCrypto, err := crypto.GenerateEncryptedKeyPair(c.keySize, c.keyAlgorithm)
	if err != nil {
		return err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}

	privateKeyExp := time.Now().UTC().Add(c.privateKeyLifetime)
	publicKeyExp := time.Now().UTC().Add(c.publicKeyLifetime)

	keyPairWriteModel := NewKeyPairWriteModel(keyID, domain.IAMID)
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSigning,
		algorithm,
		privateCrypto, publicCrypto,
		privateKeyExp, publicKeyExp))
	return err
}

func setOIDCCtx(ctx context.Context) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: oidcUser, OrgID: domain.IAMID})
}
