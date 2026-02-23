if (process.env.ZITADEL_TLS_ENABLED === "true") {
  process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
}

const scheme = process.env.ZITADEL_TLS_ENABLED === "true" ? "https" : "http";
const port = process.env.PORT || "3000";
const url = process.argv[2] || `${scheme}://localhost:${port}/ui/v2/login/healthy`;

try {
  const res = await fetch(url);
  if (!res.ok) process.exit(1);
  process.exit(0);
} catch (e) {
  process.exit(1);
}
