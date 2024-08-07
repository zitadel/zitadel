package main

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/integration"
)

func main() {
	logging.OnError(integration.InitTesterState(context.TODO())).Fatal("integration setup failed")
}
