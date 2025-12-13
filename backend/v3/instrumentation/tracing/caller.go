package tracing

import (
	"runtime"

	"github.com/zitadel/logging"
)

func GetCaller() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		logging.WithFields("logID", "TRACE-rWjfC").Debug("no caller")
		return ""
	}
	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		logging.WithFields("logID", "TRACE-25POw").Debug("caller was nil")
		return ""
	}
	return caller.Name()
}
