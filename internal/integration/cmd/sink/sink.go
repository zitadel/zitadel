package main

import (
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/integration/sink"
)

func main() {
	err := sink.ListenAndServe()
	logging.OnError(err).Fatal("running sink failed")
}
