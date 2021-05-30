package crdb

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

//reduce implements handler.Reduce function
func (h *StatementHandler) reduce(event eventstore.EventReader) ([]handler.Statement, error) {
	reduce, ok := h.reduces[string(event.Aggregate().Typ)+"."+string(event.Type())]
	if !ok {
		logging.LogWithFields("CRDB-TXMc2", "aggregateType", event.Aggregate().Typ, "eventType", event.Type()).Panic("no reduce function found")
	}

	return reduce(event)
}
