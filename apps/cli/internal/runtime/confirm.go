package runtime

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// destructiveVerbs are prefixes that indicate a destructive operation.
var destructiveVerbs = []string{
	"delete",
	"deactivate",
	"remove",
	"revoke",
	"reset",
	"clear",
}

// isDestructiveVerb returns true if the verb indicates a destructive operation.
func isDestructiveVerb(verb string) bool {
	for _, prefix := range destructiveVerbs {
		if verb == prefix || strings.HasPrefix(verb, prefix+"-") {
			return true
		}
	}
	return false
}

// confirmDestructive prompts the user to confirm a destructive operation.
// It extracts the primary resource identifier from the request (first positional
// arg or first ID-like field) and asks the user to type it back.
//
// Confirmation is skipped when:
//   - --yes/-y flag is set
//   - stdin is not a TTY (non-interactive; returns an error instead)
func confirmDestructive(cmd *cobra.Command, spec CommandSpec, req *dynamicpb.Message, reqDesc protoreflect.MessageDescriptor, yesFlag bool) error {
	if yesFlag {
		return nil
	}

	// Extract the resource identifier for the confirmation prompt.
	resourceID := extractResourceID(req, reqDesc, spec)

	// Non-interactive: refuse to run destructive commands without --yes.
	if !isTerminal(os.Stdin) {
		if resourceID != "" {
			return fmt.Errorf("destructive operation on %q requires confirmation; use --yes to skip in non-interactive mode", resourceID)
		}
		return fmt.Errorf("destructive operation requires confirmation; use --yes to skip in non-interactive mode")
	}

	// Interactive: prompt the user.
	action := spec.Verb
	resource := spec.Group

	fmt.Fprintf(os.Stderr, "\n⚠️  You are about to %s %s", action, resource)
	if resourceID != "" {
		fmt.Fprintf(os.Stderr, " %q", resourceID)
	}
	fmt.Fprintln(os.Stderr, ". This action cannot be undone.")

	if resourceID != "" {
		fmt.Fprintf(os.Stderr, "   Type the resource ID to confirm: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("aborted: no confirmation received")
		}
		input := strings.TrimSpace(scanner.Text())
		if input != resourceID {
			return fmt.Errorf("aborted: expected %q, got %q", resourceID, input)
		}
	} else {
		fmt.Fprintf(os.Stderr, "   Type 'yes' to confirm: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("aborted: no confirmation received")
		}
		input := strings.TrimSpace(scanner.Text())
		if input != "yes" {
			return fmt.Errorf("aborted: expected 'yes', got %q", input)
		}
	}

	return nil
}

// extractResourceID finds the primary resource identifier from the request message.
// It checks positional args first, then falls back to fields ending in "_id".
func extractResourceID(req *dynamicpb.Message, reqDesc protoreflect.MessageDescriptor, spec CommandSpec) string {
	// Check positional args first — they're the primary identifiers.
	for _, pa := range spec.PositionalArgs {
		fd := reqDesc.Fields().ByName(protoreflect.Name(pa.ProtoFieldName))
		if fd != nil && fd.Kind() == protoreflect.StringKind {
			val := req.Get(fd).String()
			if val != "" {
				return val
			}
		}
	}

	// Fall back to the first *_id field that has a value.
	for i := 0; i < reqDesc.Fields().Len(); i++ {
		fd := reqDesc.Fields().Get(i)
		name := string(fd.Name())
		if fd.Kind() == protoreflect.StringKind && strings.HasSuffix(name, "_id") {
			val := req.Get(fd).String()
			if val != "" {
				return val
			}
		}
	}

	return ""
}

// isTerminal reports whether f is connected to a terminal.
func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
