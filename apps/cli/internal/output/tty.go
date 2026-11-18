package output

import "os"

// IsStdoutPiped returns true if the standard output is piped or redirected.
func IsStdoutPiped() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}
