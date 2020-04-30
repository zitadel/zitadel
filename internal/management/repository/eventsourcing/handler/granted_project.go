package handler

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"time"
)

type GrantedProject struct {
	handler
	eventstore eventstore.Eventstore
}

const (
	grantedProjectTable = "management.granted_projects"
)

func (p *GrantedProject) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *GrantedProject) ViewModel() string {
	return grantedProjectTable
}

func (p *GrantedProject) EventQuery() (*models.SearchQuery, error) {
	return nil, nil
}

func (p *GrantedProject) Process(event *models.Event) (err error) {
	return nil
}

func (p *GrantedProject) OnError(event *models.Event, err error) error {
	return nil
}
