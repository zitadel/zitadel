package execution

import (
	"context"
	"net/http"
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
	httpClient *http.Client,
	queue *queue.Queue,
	targetEncAlg crypto.EncryptionAlgorithm,
	activeSigningKey GetActiveSigningWebKey,
) {
	queue.ShouldStart()
	queue.AddWorkers(ctx, NewWorker(workerConfig, targetEncAlg, activeSigningKey, time.Now, httpClient))
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
}
