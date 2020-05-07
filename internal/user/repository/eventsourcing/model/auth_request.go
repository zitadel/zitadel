package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

type AuthRequest struct {
	es_models.ObjectRoot

	UserAgentID string
	*BrowserInfo
}

func AuthRequestFromModel(request *model.AuthRequest) *AuthRequest {
	return &AuthRequest{
		ObjectRoot:  request.ObjectRoot,
		BrowserInfo: BrowserInfoFromModel(request.BrowserInfo),
	}
}

func AuthRequestToModel(request *AuthRequest) *model.AuthRequest {
	return &model.AuthRequest{
		ObjectRoot:  request.ObjectRoot,
		BrowserInfo: BrowserInfoToModel(request.BrowserInfo),
	}
}

//
//func (u *User) appendUserAuthRequestChangedEvent(event *es_models.Event) error {
//	if u.AuthRequest == nil {
//		u.AuthRequest = new(AuthRequest)
//	}
//	return u.AuthRequest.setData(event)
//}
//
//func (a *AuthRequest) setData(event *es_models.Event) error {
//	a.ObjectRoot.AppendEvent(event)
//	if err := json.Unmarshal(event.Data, a); err != nil {
//		logging.Log("EVEN-clos0").WithError(err).Error("could not unmarshal event data")
//		return caos_errs.ThrowInternal(err, "MODEL-so92s", "could not unmarshal event")
//	}
//	return nil
//}

//
//func (a *AuthRequest) AddPossibleStep(step NextStep) {
//	a.possibleSteps = append(a.possibleSteps, step)
//}
