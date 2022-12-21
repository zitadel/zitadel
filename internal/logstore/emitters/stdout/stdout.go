package stdout

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/internal/logstore"
)

func NewStdoutEmitter() logstore.LogEmitter {
	return logstore.LogEmitterFunc(func(ctx context.Context, bulk []logstore.LogRecord) error {
		for idx := range bulk {
			bytes, err := json.Marshal(bulk[idx])
			if err != nil {
				return err
			}
			fmt.Println(string(bytes))
		}
		return nil
	})
}
