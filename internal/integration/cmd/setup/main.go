package main

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/integration"
)

func main() {
	logging.OnError(integration.InitFirstInstance(context.TODO())).Fatal("integration setup failed")
}
