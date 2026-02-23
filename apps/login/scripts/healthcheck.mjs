import * as https from "node:https";
import * as http from "node:http";

const scheme = process.env.ZITADEL_TLS_ENABLED === "true" ? "https" : "http";
const port = process.env.PORT || "3000";
const url = new URL(process.argv[2] || `/ui/v2/login/healthy`, `${scheme}://localhost:${port}`);

const get = scheme === "https" ? https.get : http.get;

try {
  const res = await new Promise((resolve, reject) => {
    // Safe: localhost-only probe, self-signed certs expected
    get(url, { rejectUnauthorized: false }, (res) => {
      res.resume();
      resolve(res);
    }).on("error", reject);
  });
  process.exit(res.statusCode >= 200 && res.statusCode < 400 ? 0 : 1);
} catch (e) {
  console.error("Healthcheck failed:", e);
  process.exit(1);
}
