import * as http2 from "node:http2";
import * as fs from "node:fs";

const OUTPUT_DIR = process.env.OUTPUT_DIR || "/output";

const server = http2.createSecureServer({
  key: fs.readFileSync("/certs/server.key"),
  cert: fs.readFileSync("/certs/server.crt"),
  allowHTTP1: true,
});

interface CapturedRequest {
  method: string;
  url: string;
  tlsConnected: boolean;
}

const requests: CapturedRequest[] = [];

server.on("stream", (stream, headers) => {
  requests.push({
    method: (headers[":method"] as string) || "UNKNOWN",
    url: (headers[":path"] as string) || "/",
    tlsConnected: true,
  });
  fs.writeFileSync(`${OUTPUT_DIR}/requests.json`, JSON.stringify(requests, null, 2));

  stream.respond({
    ":status": 200,
    "content-type": "application/json",
    "grpc-status": "0",
  });
  stream.end(JSON.stringify({ status: "ok" }));
});

server.listen(443, "0.0.0.0", () => {
  console.log("Mock TLS server listening on port 443");
});
