package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestHandleError_JSONFormat(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		outputFlag string
		tip        string
		wantCode   string
		wantError  string
		wantHint   string
	}{
		{
			name:       "plain error as JSON",
			err:        fmt.Errorf("something went wrong"),
			outputFlag: "json",
			tip:        "zitadel-cli users list",
			wantCode:   "",
			wantError:  "something went wrong",
			wantHint:   "Run 'zitadel-cli users list -h' for help.",
		},
		{
			name:       "connect error with code",
			err:        fmt.Errorf("[not_found] user not found"),
			outputFlag: "json",
			tip:        "zitadel-cli users get-by-id",
			wantCode:   "not_found",
			wantError:  "user not found",
			wantHint:   "Run 'zitadel-cli users get-by-id -h' for help.",
		},
		{
			name:       "unauthenticated error gets special hint",
			err:        fmt.Errorf("[unauthenticated] token expired"),
			outputFlag: "json",
			tip:        "zitadel-cli users list",
			wantCode:   "unauthenticated",
			wantError:  "token expired",
			wantHint:   "Your token may have expired or is invalid. Try running 'zitadel-cli login' again.",
		},
		{
			name:       "permission denied error gets special hint",
			err:        fmt.Errorf("[permission_denied] insufficient rights"),
			outputFlag: "json",
			tip:        "zitadel-cli orgs delete",
			wantCode:   "permission_denied",
			wantError:  "insufficient rights",
			wantHint:   "You don't have permission to perform this action. Check your project/organization grants.",
		},
		{
			name:       "no active context error",
			err:        fmt.Errorf("no active context configured"),
			outputFlag: "json",
			tip:        "zitadel-cli users list",
			wantCode:   "",
			wantError:  "no active context configured",
			wantHint:   "Run 'zitadel-cli login' to set up a context.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			HandleError(tt.err, tt.outputFlag, tt.tip, &buf)

			// Parse the JSON output.
			var result struct {
				Error string `json:"error"`
				Code  string `json:"code"`
				Hint  string `json:"hint"`
			}
			output := strings.TrimSpace(buf.String())
			if err := json.Unmarshal([]byte(output), &result); err != nil {
				t.Fatalf("output is not valid JSON: %v\nraw: %s", err, output)
			}

			if result.Error != tt.wantError {
				t.Errorf("error = %q, want %q", result.Error, tt.wantError)
			}
			if result.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", result.Code, tt.wantCode)
			}
			if result.Hint != tt.wantHint {
				t.Errorf("hint = %q, want %q", result.Hint, tt.wantHint)
			}
		})
	}
}

func TestHandleError_TextFormat(t *testing.T) {
	// Note: during test execution, stdout is piped to the test runner,
	// so IsStdoutPiped() returns true and HandleError uses JSON format
	// regardless of outputFlag. We test this by verifying the JSON
	// output is still valid even when outputFlag is not "json".
	var buf bytes.Buffer
	HandleError(fmt.Errorf("[not_found] user not found"), "table", "zitadel-cli users get-by-id", &buf)

	output := buf.String()
	// In a test environment (piped stdout), we get JSON.
	// Verify the output is parseable and contains the expected fields.
	if strings.Contains(output, "\"error\"") {
		// JSON path
		var result struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}
		if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if result.Code != "not_found" {
			t.Errorf("expected code not_found, got %q", result.Code)
		}
	} else {
		// Text path (when running with a real TTY)
		if !strings.Contains(output, "Error [not_found]") {
			t.Errorf("expected error code in text output, got: %s", output)
		}
		if !strings.Contains(output, "💡 Hint:") {
			t.Errorf("expected hint in text output, got: %s", output)
		}
	}
}
