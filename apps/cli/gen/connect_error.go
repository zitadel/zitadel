package gen

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
)

// FormatConnectError extracts the connect error code and message from
// a connect.Error, producing a structured "[CODE] message" format.
// Non-connect errors are returned unchanged.
func FormatConnectError(err error) error {
	var ce *connect.Error
	if errors.As(err, &ce) {
		return fmt.Errorf("[%s] %s", ce.Code(), ce.Message())
	}
	return err
}
