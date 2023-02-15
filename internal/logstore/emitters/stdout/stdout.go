package stdout

import (
	"context"
	"encoding/json"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/logstore"
)

func NewStdoutEmitter() logstore.LogEmitter {
	return logstore.LogEmitterFunc(func(ctx context.Context, bulk []logstore.LogRecord) error {
		for idx := range bulk {
			bytes, err := json.Marshal(bulk[idx])
			if err != nil {
				return err
			}
			logging.WithFields("record", string(bytes)).Info("log record emitted")
		}
		return nil
	})
}
