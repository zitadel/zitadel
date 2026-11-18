package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// HandleError prints the error to the given writer, formatting it as JSON if
// requested or if stdout is piped.
func HandleError(err error, outputFlag string, tip string, w io.Writer) {
	errStr := err.Error()
	code := ""
	hint := fmt.Sprintf("Run '%s -h' for help.", tip)

	// Extract connect error code from "[CODE] message" format.
	if strings.HasPrefix(errStr, "[") {
		if idx := strings.Index(errStr, "] "); idx > 0 {
			code = errStr[1:idx]
			errStr = errStr[idx+2:]
		}
	}

	// Add actionable hints based on error codes
	if code == "unauthenticated" {
		hint = "Your token may have expired or is invalid. Try running 'zitadel-cli login' again."
	} else if code == "permission_denied" {
		hint = "You don't have permission to perform this action. Check your project/organization grants."
	} else if strings.Contains(errStr, "no active context") {
		hint = "Run 'zitadel-cli login' to set up a context."
	}

	if outputFlag == "json" || IsStdoutPiped() {
		je := struct {
			Error string `json:"error"`
			Code  string `json:"code,omitempty"`
			Hint  string `json:"hint"`
		}{
			Error: errStr,
			Code:  code,
			Hint:  hint,
		}
		data, _ := json.MarshalIndent(je, "", "  ")
		fmt.Fprintln(w, string(data))
	} else {
		if code != "" {
			fmt.Fprintf(w, "Error [%s]: %v\n\n💡 Hint: %s\n", code, errStr, hint)
		} else {
			fmt.Fprintf(w, "Error: %v\n\n💡 Hint: %s\n", errStr, hint)
		}
	}
}
