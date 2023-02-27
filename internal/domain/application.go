package domain

import "net/url"

type Application interface {
	GetAppID() string
	GetApplicationName() string
	GetApplicationExternalURL() string
	GetApplicationIsVisibleToEndUser() bool
	GetApplicationLightIcon() string
	GetApplicationDarkIcon() string
	GetState() AppState
	//GetSequence() uint64
	//GetChangeDate() time.Time
	//GetResourceOwner() string
}

type AppState int32

const (
	AppStateUnspecified AppState = iota
	AppStateActive
	AppStateInactive
	AppStateRemoved
)

func (a AppState) Exists() bool {
	return !(a == AppStateUnspecified || a == AppStateRemoved)
}

type ChangeApp struct {
	AppID   string
	AppName string
	ExternalURL string
	IsVisibleToEndUser bool	
	LightIconURL string
	DarkIconURL string
	State   AppState
}

func (a *ChangeApp) GetAppID() string {
	return a.AppID
}

func (a *ChangeApp) GetApplicationName() string {
	return a.AppName
}

func (a *ChangeApp) GetApplicationExternalURL() string {
	return a.ExternalURL
}

func (a *ChangeApp) GetApplicationIsVisibleToEndUser() bool {
	return a.IsVisibleToEndUser
}

func (a *ChangeApp) GetApplicationLightIcon() string {
	return a.LightIconURL
}

func (a *ChangeApp) GetApplicationDarkIcon() string {
	return a.DarkIconURL
}

func (a *ChangeApp) GetState() AppState {
	return a.State
}

func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
    return err == nil && u.Scheme != "" && u.Host != ""
}
