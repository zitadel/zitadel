/**
 * Mock Zitadel Connect/gRPC server for testing trace propagation.
 * Captures incoming headers and writes them to a file for test assertions.
 */
const http = require("http");
const fs = require("fs");
const path = require("path");

const PORT = process.env.PORT || 8080;
const OUTPUT_DIR = process.env.OUTPUT_DIR || "/tmp/otel";
const HEADERS_FILE = path.join(OUTPUT_DIR, "captured-headers.json");

// Store captured requests
const capturedRequests = [];

const server = http.createServer((req, res) => {
  // Capture trace headers
  const traceHeaders = {
    timestamp: new Date().toISOString(),
    method: req.method,
    url: req.url,
    traceparent: req.headers["traceparent"] || null,
    tracestate: req.headers["tracestate"] || null,
    baggage: req.headers["baggage"] || null,
    grpcAcceptEncoding: req.headers["grpc-accept-encoding"] || null,
    contentType: req.headers["content-type"] || null,
  };

  capturedRequests.push(traceHeaders);

  // Write to file for test assertions
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
  fs.writeFileSync(HEADERS_FILE, JSON.stringify(capturedRequests, null, 2));

  console.log(`[${traceHeaders.timestamp}] ${req.method} ${req.url}`);
  console.log(`  traceparent: ${traceHeaders.traceparent || "(none)"}`);

  // Return mock Connect/gRPC responses
  // The actual response doesn't matter much - we just need to not error
  const url = req.url || "";

  // Health check
  if (url.includes("/health") || url === "/") {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ status: "ok" }));
    return;
  }

  // Return captured headers (for debugging)
  if (url === "/captured-headers") {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify(capturedRequests, null, 2));
    return;
  }

  // Mock Connect/gRPC responses
  // These are minimal responses to prevent the login app from crashing
  res.writeHead(200, {
    "Content-Type": "application/json",
    "grpc-status": "0",
    "grpc-message": "OK",
  });

  // Return empty but valid protobuf-like responses
  if (url.includes("GetLoginSettings") || url.includes("SettingsService")) {
    res.end(JSON.stringify({ settings: {} }));
  } else if (url.includes("GetBrandingSettings")) {
    res.end(JSON.stringify({ settings: {} }));
  } else if (url.includes("SessionService") || url.includes("Session")) {
    res.end(JSON.stringify({ session: {} }));
  } else if (url.includes("UserService") || url.includes("User")) {
    res.end(JSON.stringify({ user: {} }));
  } else {
    res.end(JSON.stringify({}));
  }
});

server.listen(PORT, "0.0.0.0", () => {
  console.log(`Mock Zitadel server listening on port ${PORT}`);
  console.log(`Captured headers will be written to: ${HEADERS_FILE}`);
});

// Graceful shutdown
process.on("SIGTERM", () => {
  console.log("Shutting down mock server...");
  server.close(() => {
    console.log("Mock server stopped");
    process.exit(0);
  });
});
