package domain

type Application interface {
	GetAppID() string
	GetApplicationName() string
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
	State   AppState
}

func (a *ChangeApp) GetAppID() string {
	return a.AppID
}

func (a *ChangeApp) GetApplicationName() string {
	return a.AppName
}

func (a *ChangeApp) GetState() AppState {
	return a.State
}

type LoginVersion int32

const (
	LoginVersionUnspecified LoginVersion = iota
	LoginVersion1
	LoginVersion2
)
