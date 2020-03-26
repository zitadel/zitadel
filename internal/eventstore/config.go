package eventstore

import (
	"github.com/caos/eventstore-lib"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

type Config struct {
	Repository repository.Config
}

func Start(conf Config) App {
	repo := repository.Start(conf.Repository)
	return &app{
		eventstore: eventstore.Start(eventstore.Config{repo}),
	}
}
