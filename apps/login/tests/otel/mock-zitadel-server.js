/**
 * Mock Zitadel Connect/gRPC server for testing trace propagation.
 * Captures incoming headers and writes them to a file for test assertions.
 *
 * Uses HTTP/2 with cleartext (h2c) to support gRPC protocol.
 */
const http2 = require("http2");
const http = require("http");
const fs = require("fs");
const path = require("path");

const PORT = process.env.PORT || 8080;
const OUTPUT_DIR = process.env.OUTPUT_DIR || "/tmp/otel";
const HEADERS_FILE = path.join(OUTPUT_DIR, "captured-headers.json");

// Store captured requests
const capturedRequests = [];

// Handle HTTP/2 request
function handleRequest(stream, headers, isHttp2 = true) {
  const method = isHttp2 ? headers[":method"] : headers.method;
  const url = isHttp2 ? headers[":path"] : headers.url;

  // Capture trace headers
  const traceHeaders = {
    timestamp: new Date().toISOString(),
    method: method,
    url: url,
    traceparent: headers["traceparent"] || null,
    tracestate: headers["tracestate"] || null,
    baggage: headers["baggage"] || null,
    grpcAcceptEncoding: headers["grpc-accept-encoding"] || null,
    contentType: headers["content-type"] || null,
  };

  capturedRequests.push(traceHeaders);

  // Write to file for test assertions
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
  fs.writeFileSync(HEADERS_FILE, JSON.stringify(capturedRequests, null, 2));

  console.log(`[${traceHeaders.timestamp}] ${method} ${url}`);
  console.log(`  traceparent: ${traceHeaders.traceparent || "(none)"}`);

  // Helper to send response
  const sendResponse = (statusCode, contentType, body) => {
    if (isHttp2) {
      stream.respond({
        ":status": statusCode,
        "content-type": contentType,
        "grpc-status": "0",
        "grpc-message": "OK",
      });
      stream.end(body);
    } else {
      stream.writeHead(statusCode, {
        "Content-Type": contentType,
        "grpc-status": "0",
        "grpc-message": "OK",
      });
      stream.end(body);
    }
  };

  // Health check
  if (url.includes("/health") || url === "/") {
    sendResponse(200, "application/json", JSON.stringify({ status: "ok" }));
    return;
  }

  // Return captured headers (for debugging)
  if (url === "/captured-headers") {
    sendResponse(200, "application/json", JSON.stringify(capturedRequests, null, 2));
    return;
  }

  // Mock Connect/gRPC responses
  // These are minimal responses to prevent the login app from crashing
  if (url.includes("GetLoginSettings") || url.includes("SettingsService")) {
    sendResponse(200, "application/json", JSON.stringify({ settings: {} }));
  } else if (url.includes("GetBrandingSettings")) {
    sendResponse(200, "application/json", JSON.stringify({ settings: {} }));
  } else if (url.includes("SessionService") || url.includes("Session")) {
    sendResponse(200, "application/json", JSON.stringify({ session: {} }));
  } else if (url.includes("UserService") || url.includes("User")) {
    sendResponse(200, "application/json", JSON.stringify({ user: {} }));
  } else if (url.includes("OrganizationService") || url.includes("Organization")) {
    sendResponse(200, "application/json", JSON.stringify({ result: [] }));
  } else if (url.includes("IdentityProviderService") || url.includes("IdentityProvider")) {
    sendResponse(200, "application/json", JSON.stringify({ identityProviders: [] }));
  } else {
    sendResponse(200, "application/json", JSON.stringify({}));
  }
}

// Create HTTP/2 server (cleartext for testing)
const http2Server = http2.createServer();

http2Server.on("stream", (stream, headers) => {
  handleRequest(stream, headers, true);
});

http2Server.on("error", (err) => console.error("HTTP/2 error:", err));

// Also create HTTP/1.1 server for health checks from wget
const httpServer = http.createServer((req, res) => {
  // Create a headers object compatible with handleRequest
  const headers = {
    method: req.method,
    url: req.url,
    ...req.headers,
  };
  handleRequest(res, headers, false);
});

// Start HTTP/2 server on main port
http2Server.listen(PORT, "0.0.0.0", () => {
  console.log(`Mock Zitadel HTTP/2 server listening on port ${PORT}`);
  console.log(`Captured headers will be written to: ${HEADERS_FILE}`);
});

// Start HTTP/1.1 server on a different port for health checks
const HEALTH_PORT = 7432;
httpServer.listen(HEALTH_PORT, "0.0.0.0", () => {
  console.log(`Mock Zitadel HTTP/1.1 health server listening on port ${HEALTH_PORT}`);
});

// Graceful shutdown
process.on("SIGTERM", () => {
  console.log("Shutting down mock servers...");
  http2Server.close(() => {
    httpServer.close(() => {
      console.log("Mock servers stopped");
      process.exit(0);
    });
  });
});
