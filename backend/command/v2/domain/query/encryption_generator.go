package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
)

type encryptionConfigReceiver interface {
	GetEncryptionConfig(ctx context.Context) (*crypto.GeneratorConfig, error)
}

type encryptionGenerator struct {
	receiver  encryptionConfigReceiver
	algorithm crypto.EncryptionAlgorithm

	res crypto.Generator
}

func QueryEncryptionGenerator(receiver encryptionConfigReceiver, algorithm crypto.EncryptionAlgorithm) *encryptionGenerator {
	return &encryptionGenerator{
		receiver:  receiver,
		algorithm: algorithm,
	}
}

func (q *encryptionGenerator) Execute(ctx context.Context) error {
	config, err := q.receiver.GetEncryptionConfig(ctx)
	if err != nil {
		return err
	}
	q.res = crypto.NewEncryptionGenerator(*config, q.algorithm)
	return nil
}

func (q *encryptionGenerator) Name() string {
	return "query.encryption_generator"
}

func (q *encryptionGenerator) Result() crypto.Generator {
	return q.res
}
