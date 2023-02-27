package object

import (
	"net"
	"time"

	"github.com/dop251/goja"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/domain"
)

// AuthRequestField accepts the domain.AuthRequest by value, so it's not mutated
func AuthRequestField(authRequest *domain.AuthRequest) func(c *actions.FieldConfig) interface{} {
	return func(c *actions.FieldConfig) interface{} {
		return AuthRequestFromDomain(c, authRequest)
	}
}

func AuthRequestFromDomain(c *actions.FieldConfig, request *domain.AuthRequest) goja.Value {
	var maxAuthAge *time.Duration
	if request.MaxAuthAge != nil {
		maxAuthAgeCopy := *request.MaxAuthAge
		maxAuthAge = &maxAuthAgeCopy
	}

	return c.Runtime.ToValue(&authRequest{
		Id:                       request.ID,
		AgentId:                  request.AgentID,
		CreationDate:             request.CreationDate,
		ChangeDate:               request.ChangeDate,
		BrowserInfo:              browserInfoFromDomain(request.BrowserInfo),
		ApplicationId:            request.ApplicationID,
		CallbackUri:              request.CallbackURI,
		TransferState:            request.TransferState,
		Prompt:                   request.Prompt,
		UiLocales:                request.UiLocales,
		LoginHint:                request.LoginHint,
		MaxAuthAge:               maxAuthAge,
		InstanceId:               request.InstanceID,
		Request:                  requestFromDomain(request.Request),
		UserId:                   request.UserID,
		UserName:                 request.UserName,
		LoginName:                request.LoginName,
		DisplayName:              request.DisplayName,
		ResourceOwner:            request.UserOrgID,
		RequestedOrgId:           request.RequestedOrgID,
		RequestedOrgName:         request.RequestedOrgName,
		RequestedPrimaryDomain:   request.RequestedPrimaryDomain,
		RequestedOrgDomain:       request.RequestedOrgDomain,
		ApplicationResourceOwner: request.ApplicationResourceOwner,
		PrivateLabelingSetting:   request.PrivateLabelingSetting,
		SelectedIdpConfigId:      request.SelectedIDPConfigID,
		LinkingUsers:             externalUsersFromDomain(request.LinkingUsers),
		PasswordVerified:         request.PasswordVerified,
		MfasVerified:             request.MFAsVerified,
		Audience:                 request.Audience,
		AuthTime:                 request.AuthTime,
	})
}

type authRequest struct {
	Id            string
	AgentId       string
	CreationDate  time.Time
	ChangeDate    time.Time
	BrowserInfo   *browserInfo
	ApplicationId string
	CallbackUri   string
	TransferState string
	Prompt        []domain.Prompt
	UiLocales     []string
	LoginHint     string
	MaxAuthAge    *time.Duration
	InstanceId    string
	Request       *request
	UserId        string
	UserName      string
	LoginName     string
	DisplayName   string
	// UserOrgID string
	ResourceOwner string
	// requested by scope
	RequestedOrgId string
	// requested by scope
	RequestedOrgName string
	// requested by scope
	RequestedPrimaryDomain string
	// requested by scope
	RequestedOrgDomain bool
	// client
	ApplicationResourceOwner string
	PrivateLabelingSetting   domain.PrivateLabelingSetting
	SelectedIdpConfigId      string
	LinkingUsers             []*externalUser
	PasswordVerified         bool
	MfasVerified             []domain.MFAType
	Audience                 []string
	AuthTime                 time.Time
}

func browserInfoFromDomain(info *domain.BrowserInfo) *browserInfo {
	if info == nil {
		return nil
	}
	return &browserInfo{
		UserAgent:      info.UserAgent,
		AcceptLanguage: info.AcceptLanguage,
		RemoteIp:       info.RemoteIP,
	}
}

func requestFromDomain(req domain.Request) *request {
	r := new(request)

	if oidcRequest, ok := req.(*domain.AuthRequestOIDC); ok {
		r.Oidc = OIDCRequest{Scopes: oidcRequest.Scopes}
	}

	return r
}

type request struct {
	Oidc OIDCRequest
}

type OIDCRequest struct {
	Scopes []string
}

type browserInfo struct {
	UserAgent      string
	AcceptLanguage string
	RemoteIp       net.IP
}
