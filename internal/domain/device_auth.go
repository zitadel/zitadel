package domain

import (
	"strconv"
)

// DeviceAuthState describes the step the
// the device authorization process is in.
// We generate the Stringer implementation for prettier
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
	DeviceAuthStateDone                             // done

	deviceAuthStateCount // invalid
)

// Exists returns true when not Undefined and
// any status lower than deviceAuthStateCount.
func (s DeviceAuthState) Exists() bool {
	return s > DeviceAuthStateUndefined && s < deviceAuthStateCount
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
