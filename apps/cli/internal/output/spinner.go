package output

import (
	"time"

	"github.com/briandowns/spinner"
)

var (
	activeSpinner *spinner.Spinner
)

// StartSpinner starts a loading indicator if stdout is a TTY.
func StartSpinner(message string) {
	if IsStdoutPiped() {
		return
	}
	if activeSpinner == nil {
		activeSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	}
	activeSpinner.Suffix = " " + message
	activeSpinner.Start()
}

// StopSpinner stops the active loading indicator.
func StopSpinner() {
	if activeSpinner != nil && activeSpinner.Active() {
		activeSpinner.Stop()
	}
}
