package runtime

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/zitadel/zitadel/apps/cli/internal/config"

	// Register proto descriptors.
	_ "github.com/zitadel/zitadel/pkg/grpc/application/v2/applicationconnect"
	_ "github.com/zitadel/zitadel/pkg/grpc/user/v2/userconnect"
)

func init() {
	// Register minimal specs for testing, since we can't import gen/ (circular dep).
	RegisterAll([]CommandSpec{
		{
			Group:          "users",
			Verb:           "list",
			FullMethodName: "zitadel.user.v2.UserService/ListUsers",
			Short:          "List users",
			IsListMethod:   true,
			ListFieldName:  "result",
		},
		{
			Group:          "users",
			Verb:           "create",
			FullMethodName: "zitadel.user.v2.UserService/CreateUser",
			Short:          "Create user",
		},
		{
			Group:          "apps",
			Verb:           "create",
			FullMethodName: "zitadel.application.v2.ApplicationService/CreateApplication",
			Short:          "Create application",
		},
		{
			Group:          "apps",
			Verb:           "delete",
			FullMethodName: "zitadel.application.v2.ApplicationService/DeleteApplication",
			Short:          "Delete application",
		},
	})
}

// sendMCP sends JSON-RPC messages to an in-process MCP server and returns responses.
// WARNING: not safe for t.Parallel() — mutates os.Stdin/os.Stdout.
func sendMCP(t *testing.T, services []string, messages ...string) []json.RawMessage {
	t.Helper()

	input := strings.Join(messages, "\n") + "\n"

	server := NewMCPServer(func() *config.Config { return &config.Config{} }, services)

	// Redirect stdin to our input
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		_, _ = w.WriteString(input)
		_ = w.Close()
	}()

	// Capture stdout
	oldStdout := os.Stdout
	outR, outW, _ := os.Pipe()
	os.Stdout = outW
	defer func() { os.Stdout = oldStdout }()

	// Drain stdout concurrently to prevent pipe buffer deadlock
	// when MCP responses exceed 64KB.
	var outBuf bytes.Buffer
	outDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(&outBuf, outR)
		close(outDone)
	}()

	_ = server.Run()
	_ = outW.Close()
	<-outDone

	var responses []json.RawMessage
	scanner := bufio.NewScanner(bytes.NewReader(outBuf.Bytes()))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 0 {
			responses = append(responses, json.RawMessage(append([]byte(nil), line...)))
		}
	}
	_ = outR.Close()
	return responses
}

// parseResponse extracts a top-level field from a JSON-RPC response.
func parseResponse(t *testing.T, raw json.RawMessage) map[string]json.RawMessage {
	t.Helper()
	var resp map[string]json.RawMessage
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	return resp
}

func TestMCPInitialize(t *testing.T) {
	responses := sendMCP(t, nil,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`,
	)
	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(responses))
	}
	resp := parseResponse(t, responses[0])

	// Must have a result, not an error
	if _, ok := resp["error"]; ok {
		t.Fatalf("expected no error, got: %s", resp["error"])
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp["result"], &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if result["protocolVersion"] != "2024-11-05" {
		t.Errorf("expected protocolVersion=2024-11-05, got %v", result["protocolVersion"])
	}
}

func TestMCPToolsList(t *testing.T) {
	responses := sendMCP(t, nil,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
	)
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	resp := parseResponse(t, responses[1])
	if _, ok := resp["error"]; ok {
		t.Fatalf("expected no error, got: %s", resp["error"])
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(resp["result"], &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	var tools []map[string]interface{}
	if err := json.Unmarshal(result["tools"], &tools); err != nil {
		t.Fatalf("unmarshal tools: %v", err)
	}
	if len(tools) == 0 {
		t.Error("expected at least one tool")
	}

	// Every tool must have name, description, inputSchema
	for _, tool := range tools {
		if tool["name"] == nil || tool["name"] == "" {
			t.Error("tool missing name")
		}
		if tool["description"] == nil {
			t.Error("tool missing description")
		}
		if tool["inputSchema"] == nil {
			t.Error("tool missing inputSchema")
		}
	}
}

func TestMCPToolsListFiltered(t *testing.T) {
	responses := sendMCP(t, []string{"users"},
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
	)
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	resp := parseResponse(t, responses[1])

	var result map[string]json.RawMessage
	if err := json.Unmarshal(resp["result"], &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	var tools []map[string]interface{}
	if err := json.Unmarshal(result["tools"], &tools); err != nil {
		t.Fatalf("unmarshal tools: %v", err)
	}

	for _, tool := range tools {
		name := tool["name"].(string)
		if !strings.HasPrefix(name, "users_") {
			t.Errorf("expected only users_* tools when filtered, got %s", name)
		}
	}
}

func TestMCPToolsCallUnknown(t *testing.T) {
	responses := sendMCP(t, nil,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"nonexistent_tool","arguments":{}}}`,
	)
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	resp := parseResponse(t, responses[1])

	// Should be a JSON-RPC error, not a panic
	if _, ok := resp["error"]; !ok {
		t.Fatal("expected JSON-RPC error for unknown tool")
	}
}

func TestMCPToolsCallInvalidArgs(t *testing.T) {
	responses := sendMCP(t, nil,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"users_list","arguments":"not json"}}`,
	)
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	resp := parseResponse(t, responses[1])

	// Should be an error, not a panic
	var result map[string]interface{}
	if err := json.Unmarshal(resp["result"], &result); err == nil {
		if isError, ok := result["isError"]; ok && isError.(bool) {
			return // expected
		}
	}
	// Also accept a JSON-RPC error
	if _, ok := resp["error"]; ok {
		return // also acceptable
	}
	t.Fatal("expected error for invalid args")
}

func TestMCPInputSchemaHasEnumValues(t *testing.T) {
	responses := sendMCP(t, []string{"apps"},
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
	)
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	resp := parseResponse(t, responses[1])

	var result map[string]json.RawMessage
	if err := json.Unmarshal(resp["result"], &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	var tools []map[string]json.RawMessage
	if err := json.Unmarshal(result["tools"], &tools); err != nil {
		t.Fatalf("unmarshal tools: %v", err)
	}

	// Find apps_create tool
	for _, tool := range tools {
		var name string
		if err := json.Unmarshal(tool["name"], &name); err != nil {
			t.Fatalf("unmarshal tool name: %v", err)
		}
		if name == "apps_create" {
			schema := string(tool["inputSchema"])
			if !bytes.Contains(tool["inputSchema"], []byte("OIDC_APP_TYPE_WEB")) {
				t.Errorf("expected OIDC_APP_TYPE_WEB in inputSchema, got: %s", schema[:min(200, len(schema))])
			}
			return
		}
	}
	t.Error("expected apps_create tool in filtered results")
}
