import * as https from "node:https";
import * as http from "node:http";

const scheme = process.env.ZITADEL_TLS_ENABLED === "true" ? "https" : "http";
const port = process.env.PORT || "3000";
const url = process.argv[2] || `${scheme}://localhost:${port}/ui/v2/login/healthy`;

const get = scheme === "https" ? https.get : http.get;

try {
  const res = await new Promise((resolve, reject) => {
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
