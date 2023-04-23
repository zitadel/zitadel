package domain

import (
	"strconv"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

// DeviceAuth describes a Device Authorization request.
// It is used as input and output model in the command and query packages.
type DeviceAuth struct {
	models.ObjectRoot

	ClientID   string
	DeviceCode string
	UserCode   string
	Expires    time.Time
	Scopes     []string
	Subject    string
	State      DeviceAuthState
}

// DeviceAuthState describes the step the
// the device authorization process is in.
// We generate the Stringer implemntation for pretier
// log output.
//
//go:generate stringer -type=DeviceAuthState -linecomment
type DeviceAuthState uint

const (
	DeviceAuthStateUndefined DeviceAuthState = iota // undefined
	DeviceAuthStateInitiated                        // initiated
	DeviceAuthStateApproved                         // approved
	DeviceAuthStateDenied                           // denied
	DeviceAuthStateExpired                          // expired
	DeviceAuthStateRemoved                          // removed
)

// Exists returns true when not Undefined and
// any status lower than Removed.
func (s DeviceAuthState) Exists() bool {
	return s > DeviceAuthStateUndefined && s < DeviceAuthStateRemoved
}

// Done returns true when DeviceAuthState is Approved.
// This implements the OIDC interface requirement of "Done"
func (s DeviceAuthState) Done() bool {
	return s == DeviceAuthStateApproved
}

// Denied returns true when DeviceAuthState is Denied, Expired or Removed.
// This implements the OIDC interface requirement of "Denied".
func (s DeviceAuthState) Denied() bool {
	return s >= DeviceAuthStateDenied
}

func (s DeviceAuthState) GoString() string {
	return strconv.Itoa(int(s))
}

// DeviceAuthCanceled is a subset of DeviceAuthState, allowed to
// be used in the deviceauth.CanceledEvent.
// The string type is used to make the eventstore more readable
// on the reason of cancelation.
type DeviceAuthCanceled string

const (
	DeviceAuthCanceledDenied  = "denied"
	DeviceAuthCanceledExpired = "expired"
)

func (c DeviceAuthCanceled) State() DeviceAuthState {
	switch c {
	case DeviceAuthCanceledDenied:
		return DeviceAuthStateDenied
	case DeviceAuthCanceledExpired:
		return DeviceAuthStateExpired
	default:
		return DeviceAuthStateUndefined
	}
}
