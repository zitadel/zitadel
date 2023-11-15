package projection

import (
	"github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func assertEvent[T eventstore.Event](event eventstore.Event) (T, error) {
	e, ok := event.(T)
	if !ok {
		return e, errors.ThrowInvalidArgumentf(nil, "HANDL-1m9fS", "reduce.wrong.event.type %T", event)
	}
	return e, nil
}
