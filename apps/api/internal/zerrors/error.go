package zerrors

import (
	"strings"
)

// Error is a stdlib error extension.
// It contains parameters to identify errors through all application layers
type Error interface {
	GetParent() error
	GetMessage() string
	SetMessage(string)
	GetID() string
}

// Contains compares the error message with needle
func Contains(err error, needle string) bool {
	return err != nil && strings.Contains(err.Error(), needle)
}
