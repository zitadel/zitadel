package instrumentation

import (
	"log/slog"
	"runtime"
)

type Caller struct {
	Function string
	File     string
	Line     int
}

// GetCaller returns the caller information
// skipping the given number of stack frames relative to the caller of this function.
func GetCaller(skip int) (Caller, bool) {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		slog.Default().Debug("no caller info", "id", "TRACE-3r8bE")
		return Caller{}, false
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		slog.Default().Debug("caller was nil", "id", "TRACE-25POw")
		return Caller{}, false
	}
	return Caller{
		Function: f.Name(),
		File:     file,
		Line:     line,
	}, true
}

const CallerUnknown = "unknown"

// GetCallingFunc returns the caller function name,
// skipping the given number of stack frames relative to the caller of this function.
// If the caller cannot be determined, [CallerUnknown] is returned.
func GetCallingFunc(skip int) string {
	caller, ok := GetCaller(skip + 1)
	if !ok {
		return CallerUnknown
	}
	return caller.Function
}
