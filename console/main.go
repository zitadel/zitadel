package main

import (
	"context"

	"github.com/caos/zitadel/console/config"
	"github.com/caos/zitadel/console/service"
	"github.com/caos/utils/logging"
)

func main() {
	conf, err := config.ReadConfig("$PROJECT_PATH/config/startup.yaml")
	logging.Log("MAIN-o8xk9U").OnError(err).Panic("error reading config")

	ctx := context.Background()
	//console.Console()
	service.Start(ctx, conf)
	<-ctx.Done()
}
