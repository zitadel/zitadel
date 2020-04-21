package model

import (
	"net"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user_agent/model"
)

type UserAgent struct {
	es_models.ObjectRoot
	UserAgent      string         `json:"userAgent,omitempty"`
	AcceptLanguage string         `json:"acceptLanguage,omitempty"`
	RemoteIP       net.IP         `json:"remoteIP,omitempty"`
	State          int32          `json:"-"`
	UserSessions   []*UserSession `json:"-"`
}

func UserAgentFromModel(userAgent *model.UserAgent) *UserAgent {
	return &UserAgent{
		ObjectRoot:     userAgent.ObjectRoot,
		UserAgent:      userAgent.UserAgent,
		AcceptLanguage: userAgent.AcceptLanguage,
		RemoteIP:       userAgent.RemoteIP,
		State:          int32(userAgent.State),
		UserSessions:   UserSessionsFromModel(userAgent.UserSessions),
	}
}

func UserAgentToModel(userAgent *UserAgent) *model.UserAgent {
	return &model.UserAgent{
		ObjectRoot:     userAgent.ObjectRoot,
		UserAgent:      userAgent.UserAgent,
		AcceptLanguage: userAgent.AcceptLanguage,
		RemoteIP:       userAgent.RemoteIP,
		State:          model.UserAgentState(userAgent.State),
		UserSessions:   UserSessionsToModel(userAgent.UserSessions),
	}
}

func (p *UserAgent) Changes(changed *UserAgent) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func (p *UserAgent) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *UserAgent) AppendEvent(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case UserAgentAdded, UserAgentChanged:
		if err := json.Unmarshal(event.Data, p); err != nil {
			logging.Log("EVEN-46ss2").WithError(err).Error("could not unmarshal event data")
			return err
		}
		p.State = int32(model.UserAgentStateActive)
		return nil
	case UserAgentRevoked:
		p.State = int32(model.UserAgentStateRevoked)
		return nil
	case UserSessionAdded:
		return p.appendUserSessionAddedEvent(event)
	case UserSessionTerminated:
		return p.appendUserSessionTerminatedEvent(event)
	case UserNameCheckSucceeded:
		return p.appendUserNameCheckSucceededEvent(event)
	case UserNameCheckFailed:
		return p.appendUserNameCheckFailedEvent(event)
	case PasswordCheckSucceeded:
		return p.appendPasswordCheckSucceededEvent(event)
	case PasswordCheckFailed:
		return p.appendPasswordCheckFailedEvent(event)
	case MfaCheckSucceeded:
		return p.appendMfaCheckSucceededEvent(event)
	case MfaCheckFailed:
		return p.appendMfaCheckFailedEvent(event)
	case ReAuthRequested:
		return p.appendReAuthRequestedEvent(event)
	case AuthSessionAdded:
		return p.appendAuthSessionAddedEvent(event)
	case AuthSessionSet:
		return p.appendAuthSessionSetEvent(event)
	case TokenAdded:
		return p.appendTokenAddedEvent(event)
	}
	return nil
}

func (a *UserAgent) appendUserSessionAddedEvent(event *es_models.Event) error {
	s := new(UserSession)
	if err := s.getData(event); err != nil {
		return err
	}
	s.State = int32(model.UserSessionStateActive)
	a.UserSessions = append(a.UserSessions, s)
	return nil
}

func (p *UserAgent) appendUserSessionTerminatedEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendUserNameCheckSucceededEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendUserNameCheckFailedEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendPasswordCheckSucceededEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendPasswordCheckFailedEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendMfaCheckSucceededEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendMfaCheckFailedEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendReAuthRequestedEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendAuthSessionAddedEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendAuthSessionSetEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}

func (p *UserAgent) appendTokenAddedEvent(event *es_models.Event) error {
	p.State = model.UserAgentStateToInt(model.Inactive)
	return nil
}
