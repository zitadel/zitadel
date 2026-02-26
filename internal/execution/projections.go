package execution

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/queue"
)

var (
	projections []*handler.Handler
)

func Register(
	ctx context.Context,
	workerConfig WorkerConfig,
	queue *queue.Queue,
	targetEncAlg crypto.EncryptionAlgorithm,
	activeSigningKey GetActiveSigningWebKey,
) {
	queue.ShouldStart()
	queue.AddWorkers(ctx, NewWorker(workerConfig, targetEncAlg, activeSigningKey, time.Now))
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
}
