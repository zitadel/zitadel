package main

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/integration"
)

const tmpPath = "tmp/"

func main() {
	logging.OnError(integration.InitTesterState(context.TODO(), tmpPath)).Fatal("integration setup failed")
}
