package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/queue"
)

var (
	projections []*handler.Handler
)

func Register(
	workerConfig WorkerConfig,
	queue *queue.Queue,
	targetEncAlg crypto.EncryptionAlgorithm,
) {
	queue.ShouldStart()
	queue.AddWorkers(NewWorker(workerConfig, targetEncAlg))
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
}
