/**
 * Mock Zitadel ConnectRPC server for testing trace propagation.
 * Captures incoming headers and writes them to a file for test assertions.
 */
const http = require("http");
const fs = require("fs");
const path = require("path");

const PORT = process.env.PORT || 8080;
const OUTPUT_DIR = process.env.OUTPUT_DIR || "/tmp/otel";
const HEADERS_FILE = path.join(OUTPUT_DIR, "captured-headers.json");

/**
 * Safely stringify a value to JSON with XSS-preventing character escapes.
 *
 * Escapes `<`, `>`, `&`, and `/` to their Unicode equivalents to prevent
 * script injection when the output is reflected in HTTP responses or
 * embedded in HTML contexts.
 *
 * @param {unknown} value - The value to stringify
 * @returns {string} JSON string with dangerous characters escaped
 */
function safeJsonStringify(value) {
  return JSON.stringify(value, null, 2)
    .replace(/</g, "\\u003c")
    .replace(/>/g, "\\u003e")
    .replace(/&/g, "\\u0026")
    .replace(/\//g, "\\/");
}

// Store captured requests
const capturedRequests = [];

const server = http.createServer((req, res) => {
  const { method, url } = req;

  // Capture trace headers
  const traceHeaders = {
    timestamp: new Date().toISOString(),
    method,
    url,
    traceparent: req.headers["traceparent"] || null,
    tracestate: req.headers["tracestate"] || null,
    baggage: req.headers["baggage"] || null,
    contentType: req.headers["content-type"] || null,
  };

  capturedRequests.push(traceHeaders);

  // Write to file for test assertions
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
  fs.writeFileSync(HEADERS_FILE, JSON.stringify(capturedRequests, null, 2));

  console.log(`[${traceHeaders.timestamp}] ${method} ${url}`);
  console.log(`  traceparent: ${traceHeaders.traceparent || "(none)"}`);

  const sendResponse = (statusCode, contentType, body) => {
    res.writeHead(statusCode, { "Content-Type": contentType });
    res.end(body);
  };

  // Health check
  if (url.includes("/health") || url === "/") {
    sendResponse(200, "application/json", safeJsonStringify({ status: "ok" }));
    return;
  }

  // Return captured headers (for debugging)
  if (url === "/captured-headers") {
    sendResponse(200, "application/json", safeJsonStringify(capturedRequests));
    return;
  }

  // Mock ConnectRPC responses
  // These are minimal responses to prevent the login app from crashing
  if (url.includes("GetLoginSettings") || url.includes("SettingsService")) {
    sendResponse(200, "application/json", safeJsonStringify({ settings: {} }));
  } else if (url.includes("GetBrandingSettings")) {
    sendResponse(200, "application/json", safeJsonStringify({ settings: {} }));
  } else if (url.includes("SessionService") || url.includes("Session")) {
    sendResponse(200, "application/json", safeJsonStringify({ session: {} }));
  } else if (url.includes("UserService") || url.includes("User")) {
    sendResponse(200, "application/json", safeJsonStringify({ user: {} }));
  } else if (url.includes("OrganizationService") || url.includes("Organization")) {
    sendResponse(200, "application/json", safeJsonStringify({ result: [] }));
  } else if (url.includes("IdentityProviderService") || url.includes("IdentityProvider")) {
    sendResponse(200, "application/json", safeJsonStringify({ identityProviders: [] }));
  } else {
    sendResponse(200, "application/json", safeJsonStringify({}));
  }
});

server.on("error", (err) => console.error("Server error:", err));

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
