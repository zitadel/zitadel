package eventsourcing

import (
	"encoding/json"
	"net"
	"time"

	"github.com/caos/logging"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user_agent/model"
)

func UserAgentFromEvents(userAgent *UserAgent, events ...*es_models.Event) (*UserAgent, error) {
	if userAgent == nil {
		userAgent = &UserAgent{}
	}

	return userAgent, userAgent.AppendEvents(events...)
}

func (p *UserAgent) Changes(changed *UserAgent) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}
