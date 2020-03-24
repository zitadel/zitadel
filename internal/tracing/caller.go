package tracing

import (
	"runtime"

	"github.com/caos/logging"
)

func GetCaller() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		logging.Log("HELPE-rWjfC").Debug("no caller")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		logging.Log("HELPE-25POw").Debug("caller was nil")
	}

	// Print the name of the function
	return caller.Name()
}
