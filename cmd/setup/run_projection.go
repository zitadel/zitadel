package setup

import (
	"context"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type runProjections struct {
	es *eventstore.Eventstore

	Version string `json:"version"`
}

func (mig *runProjections) Check(lastRun map[string]interface{}) bool {
	currentVersion, _ := lastRun["version"].(string)
	return currentVersion != mig.Version
}

func (mig *runProjections) Execute(ctx context.Context) error {
	cmd, err := command.StartCommands(
		mig.es,
		mig.defaults,
		nil,
		nil,
		nil,
		mig.ExternalDomain,
		mig.ExternalSecure,
		mig.ExternalPort,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
		0,
		0,
		nil,
	)

	if err != nil {
		return err
	}
	return cmd.ChangeSystemConfig(ctx, mig.currentExternalDomain, mig.currentExternalPort, mig.currentExternalSecure)
}

func (mig *runProjections) String() string {
	return "run_projections"
}
