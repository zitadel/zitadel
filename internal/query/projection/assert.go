package projection

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

func assertEvent[T eventstore.Event](event eventstore.Event) (T, error) {
	e, ok := event.(T)
	if !ok {
		return e, zerrors.ThrowInvalidArgumentf(nil, "HANDL-1m9fS", "reduce.wrong.event.type %T", event)
	}
	return e, nil
}
