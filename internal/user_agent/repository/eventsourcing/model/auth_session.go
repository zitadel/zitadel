package model

//
//import (
//	es_models "github.com/caos/zitadel/internal/eventstore/models"
//	"github.com/caos/zitadel/internal/user_agent/model"
//)
//
//type AuthSession struct {
//	es_models.ObjectRoot
//
//	SessionID string
//	//Type                  AuthSessionType
//	BrowserInfo           *BrowserInfo
//	ApplicationID         string   //clientID
//	CallbackURI           string   //redirectURi
//	TransferState         string   //state //oidc only?
//	Prompt                int32    //name?
//	LevelOfAssurance      string   //acr
//	RequestedPossibleLOAs []string //acr_values
//	RequestedUiLocales    []string //language.Tag?
//	LoginHint             string
//	PreselectedUserID     string
//	MaxAuthAge            uint32
//	ProjectApplicationIDs []string //aud?
//	//OIDC                  *AuthSessionOIDC
//	Request     Request
//	UserSession *UserSession
//	Token       *Token
//}
//
//func AuthSessionFromModel(authSession *model.AuthSession) *AuthSession {
//	return &AuthSession{
//		ObjectRoot:            authSession.ObjectRoot,
//		SessionID:             authSession.SessionID,
//		BrowserInfo:           BrowserInfoFromModel(authSession.BrowserInfo),
//		ApplicationID:         authSession.ApplicationID,
//		CallbackURI:           authSession.CallbackURI,
//		TransferState:         authSession.TransferState,
//		Prompt:                int32(authSession.Prompt),
//		LevelOfAssurance:      authSession.LevelOfAssurance,
//		RequestedPossibleLOAs: authSession.RequestedPossibleLOAs,
//		RequestedUiLocales:    authSession.RequestedUiLocales,
//		LoginHint:             authSession.LoginHint,
//		PreselectedUserID:     authSession.PreselectedUserID,
//		MaxAuthAge:            authSession.MaxAuthAge,
//		ProjectApplicationIDs: authSession.ProjectApplicationIDs,
//		Request:               RequestFromModel(authSession.Request),
//		UserSession:           UserSessionFromModel(authSession.UserSession),
//	}
//}
//
//func AuthSessionToModel(authSession *AuthSession) *model.AuthSession {
//	return &model.AuthSession{
//		ObjectRoot:            authSession.ObjectRoot,
//		SessionID:             authSession.SessionID,
//		BrowserInfo:           BrowserInfoToModel(authSession.BrowserInfo),
//		ApplicationID:         authSession.ApplicationID,
//		CallbackURI:           authSession.CallbackURI,
//		TransferState:         authSession.TransferState,
//		Prompt:                model.Prompt(authSession.Prompt),
//		RequestedPossibleLOAs: authSession.RequestedPossibleLOAs,
//		RequestedUiLocales:    authSession.RequestedUiLocales,
//		LoginHint:             authSession.LoginHint,
//		PreselectedUserID:     authSession.PreselectedUserID,
//		MaxAuthAge:            authSession.MaxAuthAge,
//		ProjectApplicationIDs: authSession.ProjectApplicationIDs,
//		Request:               RequestToModel(authSession.Request),
//		UserSession:           UserSessionToModel(authSession.UserSession),
//	}
//}
//
//func GetAuthSession(sessions []*AuthSession, id string) (int, *AuthSession) {
//	for i, s := range sessions {
//		if s.SessionID == id {
//			return i, s
//		}
//	}
//	return -1, nil
//}
//
//func (p *UserAgent) appendAuthSessionAddedEvent(event *es_models.Event) error {
//	p.State = model.UserAgentStateToInt(model.Inactive)
//	return nil
//}
//
//func (p *UserAgent) appendAuthSessionSetEvent(event *es_models.Event) error {
//	p.State = model.UserAgentStateToInt(model.Inactive)
//	return nil
//}
