package command

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/v2/domain"
	keypair "github.com/caos/zitadel/internal/v2/repository/keypair"
	"time"
)

const (
	oidcUser = "OIDC"
)

func (r *CommandSide) GenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	ctx = setOIDCCtx(ctx)
	privateCrypto, publicCrypto, err := crypto.GenerateEncryptedKeyPair(r.keySize, r.keyAlgorithm)
	if err != nil {
		return err
	}
	keyID, err := r.idGenerator.Next()
	if err != nil {
		return err
	}

	privateKeyExp := time.Now().UTC().Add(r.privateKeyLifetime)
	publicKeyExp := time.Now().UTC().Add(r.publicKeyLifetime)

	keyPairWriteModel := NewKeyPairWriteModel(keyID, domain.IAMID)
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, keypair.NewAddedEvent(
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
