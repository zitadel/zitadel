package federatedlogout

type Index int

const (
	IndexUnspecified Index = iota
	IndexRequestID
)

type FederatedLogout struct {
	InstanceID            string
	FingerPrintID         string
	SessionID             string
	IDPID                 string
	UserID                string
	PostLogoutRedirectURI string
	State                 State
}

// Keys implements cache.Entry
func (c *FederatedLogout) Keys(i Index) []string {
	if i == IndexRequestID {
		return []string{Key(c.InstanceID, c.SessionID)}
	}
	return nil
}

func Key(instanceID, sessionID string) string {
	return instanceID + "-" + sessionID
}

type State int

const (
	StateCreated State = iota
	StateRedirected
)
