package runtime

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/zitadel/zitadel/apps/cli/internal/auth"
	"github.com/zitadel/zitadel/apps/cli/internal/client"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

// MCPServer runs a JSON-RPC 2.0 server over stdio, exposing CLI commands as MCP tools.
type MCPServer struct {
	getCfg  func() *config.Config
	filter  map[string]bool // if non-empty, only expose these groups
	specs   []CommandSpec
	scanner *bufio.Scanner
}

// NewMCPServer creates a new MCP server.
func NewMCPServer(getCfg func() *config.Config, filterGroups []string) *MCPServer {
	s := &MCPServer{
		getCfg: getCfg,
		specs:  AllSpecs(),
	}
	if len(filterGroups) > 0 {
		s.filter = make(map[string]bool)
		for _, g := range filterGroups {
			s.filter[g] = true
		}
	}
	return s
}

// JSON-RPC types
type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP types
type mcpToolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type mcpToolResult struct {
	Content []mcpContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

type mcpContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type mcpCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// Run starts the MCP server loop, reading JSON-RPC from stdin and writing to stdout.
func (s *MCPServer) Run() error {
	s.scanner = bufio.NewScanner(os.Stdin)
	s.scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	for s.scanner.Scan() {
		line := s.scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req jsonrpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.writeError(nil, -32700, "parse error")
			continue
		}

		resp := s.handleRequest(req)
		if resp != nil {
			s.writeResponse(*resp)
		}
	}
	return s.scanner.Err()
}

func (s *MCPServer) handleRequest(req jsonrpcRequest) *jsonrpcResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	case "notifications/initialized":
		return nil // no response needed for notifications
	default:
		return &jsonrpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &jsonrpcError{Code: -32601, Message: fmt.Sprintf("method not found: %s", req.Method)},
		}
	}
}

func (s *MCPServer) handleInitialize(req jsonrpcRequest) *jsonrpcResponse {
	return &jsonrpcResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "zitadel-cli",
				"version": "1.0.0",
			},
		},
	}
}

func (s *MCPServer) handleToolsList(req jsonrpcRequest) *jsonrpcResponse {
	var tools []mcpToolDef
	for _, spec := range s.filteredSpecs() {
		schema := s.buildInputSchema(spec)
		tools = append(tools, mcpToolDef{
			Name:        spec.Group + "_" + strings.ReplaceAll(spec.Verb, " ", "_"),
			Description: spec.Short,
			InputSchema: schema,
		})
	}
	return &jsonrpcResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"tools": tools},
	}
}

func (s *MCPServer) handleToolsCall(req jsonrpcRequest) *jsonrpcResponse {
	var params mcpCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &jsonrpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &jsonrpcError{Code: -32602, Message: "invalid params"},
		}
	}

	// Find the spec by tool name.
	var spec *CommandSpec
	for _, s := range s.filteredSpecs() {
		name := s.Group + "_" + strings.ReplaceAll(s.Verb, " ", "_")
		if name == params.Name {
			spec = &s
			break
		}
	}
	if spec == nil {
		return &jsonrpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &jsonrpcError{Code: -32602, Message: fmt.Sprintf("unknown tool: %s", params.Name)},
		}
	}

	// Execute the tool.
	result := s.executeTool(*spec, params.Arguments)
	return &jsonrpcResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *MCPServer) executeTool(spec CommandSpec, argsJSON json.RawMessage) mcpToolResult {
	// Resolve proto descriptors.
	methodDesc, reqDesc, err := resolveMethod(spec.FullMethodName)
	if err != nil {
		return mcpToolResult{
			Content: []mcpContent{{Type: "text", Text: fmt.Sprintf("error: %v", err)}},
			IsError: true,
		}
	}

	// Build request from JSON arguments.
	req := dynamicpb.NewMessage(reqDesc)
	if len(argsJSON) > 0 && string(argsJSON) != "{}" {
		if err := protojson.Unmarshal(argsJSON, req); err != nil {
			return mcpToolResult{
				Content: []mcpContent{{Type: "text", Text: fmt.Sprintf("error parsing arguments: %v", err)}},
				IsError: true,
			}
		}
	}

	// Get client config.
	cfg := s.getCfg()
	actx, _, err := config.ActiveCtx(cfg)
	if err != nil {
		return mcpToolResult{
			Content: []mcpContent{{Type: "text", Text: fmt.Sprintf("error: no active context. Run 'zitadel-cli login' first: %v", err)}},
			IsError: true,
		}
	}

	tokenSource, err := auth.TokenSource(context.Background(), actx)
	if err != nil {
		return mcpToolResult{
			Content: []mcpContent{{Type: "text", Text: fmt.Sprintf("auth error: %v", err)}},
			IsError: true,
		}
	}

	httpClient := client.New(tokenSource)
	baseURL := client.InstanceURL(actx.Instance)

	// Make the API call.
	respMsg, err := callConnect(context.Background(), httpClient, baseURL, methodDesc, req)
	if err != nil {
		return mcpToolResult{
			Content: []mcpContent{{Type: "text", Text: fmt.Sprintf("API error: %v", err)}},
			IsError: true,
		}
	}

	// Marshal response as JSON.
	respJSON, err := protojson.MarshalOptions{Indent: "  "}.Marshal(respMsg)
	if err != nil {
		return mcpToolResult{
			Content: []mcpContent{{Type: "text", Text: fmt.Sprintf("error marshalling response: %v", err)}},
			IsError: true,
		}
	}

	return mcpToolResult{
		Content: []mcpContent{{Type: "text", Text: string(respJSON)}},
	}
}

func (s *MCPServer) buildInputSchema(spec CommandSpec) json.RawMessage {
	_, reqDesc, err := resolveMethod(spec.FullMethodName)
	if err != nil {
		return json.RawMessage(`{"type":"object"}`)
	}
	schema := messageToSchema(reqDesc, 3)
	data, _ := json.Marshal(schema)
	return data
}

func (s *MCPServer) filteredSpecs() []CommandSpec {
	if len(s.filter) == 0 {
		return s.specs
	}
	var filtered []CommandSpec
	for _, spec := range s.specs {
		if s.filter[spec.Group] {
			filtered = append(filtered, spec)
		}
	}
	return filtered
}

func (s *MCPServer) writeResponse(resp jsonrpcResponse) {
	data, _ := json.Marshal(resp)
	fmt.Fprintln(os.Stdout, string(data))
}

func (s *MCPServer) writeError(id json.RawMessage, code int, message string) {
	resp := jsonrpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &jsonrpcError{Code: code, Message: message},
	}
	s.writeResponse(resp)
}

// Ensure io is used (for future streaming support).
var _ = io.EOF
