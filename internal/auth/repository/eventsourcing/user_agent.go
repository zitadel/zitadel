package eventsourcing

import (
	"context"

	proj_model "github.com/caos/zitadel/internal/project/model"
	user_agent_model "github.com/caos/zitadel/internal/user_agent/model"
	user_agent_event "github.com/caos/zitadel/internal/user_agent/repository/eventsourcing"
)

type UserAgentRepo struct {
	UserAgentEvents *user_agent_event.UserAgentEventstore
	//view      *view.View
}

func (repo *UserAgentRepo) UserAgentByID(ctx context.Context, id string) (*user_agent_model.UserAgent, err error) {
	return repo.UserAgentEvents.UserAgentByID(ctx, id)
}

func (repo *UserAgentRepo) CreateUserAgent(ctx context.Context, name string) (*proj_model.Project, error) {
	project := &proj_model.Project{Name: name}
	return repo.UserAgentEvents.CreateProject(ctx, project)
}
