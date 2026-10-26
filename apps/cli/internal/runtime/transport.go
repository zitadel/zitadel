package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// callConnect performs a generic unary ConnectRPC call using HTTP POST + protojson.
//
// Connect protocol: POST to baseURL/package.Service/Method with Content-Type: application/json.
// The body is the protojson-encoded request message. The response is protojson-decoded.
func callConnect(ctx context.Context, httpClient *http.Client, baseURL string, methodDesc protoreflect.MethodDescriptor, req *dynamicpb.Message) (proto.Message, error) {
	// Build the URL: baseURL/full.package.ServiceName/MethodName
	serviceFQN := methodDesc.Parent().(protoreflect.ServiceDescriptor).FullName()
	methodName := methodDesc.Name()
	url := strings.TrimRight(baseURL, "/") + "/" + string(serviceFQN) + "/" + string(methodName)

	// Marshal request to JSON.
	reqBody, err := protojson.MarshalOptions{}.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	// Build HTTP request.
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute.
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling %s: %w", url, err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, parseConnectError(httpResp.StatusCode, respBody)
	}

	// Unmarshal into a dynamic message for the response type.
	respMsg := dynamicpb.NewMessage(methodDesc.Output())
	if err := protojson.Unmarshal(respBody, respMsg); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return respMsg, nil
}

// connectErrorBody is the JSON structure returned by ConnectRPC on error.
type connectErrorBody struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details,omitempty"`
}

// parseConnectError parses a ConnectRPC error response body.
func parseConnectError(statusCode int, body []byte) error {
	var ce connectErrorBody
	if err := json.Unmarshal(body, &ce); err == nil && ce.Code != "" {
		return fmt.Errorf("[%s] %s", ce.Code, ce.Message)
	}
	return fmt.Errorf("HTTP %d: %s", statusCode, string(body))
}
