const BASE_DIRECTIVES: Record<string, string[]> = {
  "default-src": ["'self'"],
  "script-src": ["'self'", "'unsafe-inline'", "'unsafe-eval'"],
  "connect-src": ["'self'"],
  "style-src": ["'self'", "'unsafe-inline'"],
  "font-src": ["'self'"],
  "img-src": ["'self'"],
  "frame-ancestors": ["'none'"],
  "object-src": ["'none'"],
};

export interface CSPOptions {
  serviceUrl?: string;
  iframeOrigins?: string[] | null;
}

/* The API URL is how the server reaches ZITADEL, which in a container
 * deployment is an internal name such as `http://zitadel:8080`. Naming it in a
 * header the browser reads publishes the service's internal hostname to every
 * visitor, and is of no use to a browser that is not on that network.
 *
 * Best effort, on the shape of the hostname alone: a dotted name, an IPv6
 * literal, or `localhost` can be a real destination for a browser; a bare
 * single label is a container or service name. It errs towards keeping the
 * URL, so a reachable host is never dropped by mistake. */
function isBrowserReachable(serviceUrl: string): boolean {
  try {
    const { hostname } = new URL(serviceUrl);
    // URL keeps the brackets around an IPv6 literal, e.g. "[::1]".
    const isIpv6Literal = hostname.startsWith("[") && hostname.endsWith("]");
    return hostname.includes(".") || isIpv6Literal || hostname === "localhost";
  } catch {
    return false;
  }
}

export function buildCSP(options: CSPOptions = {}): string {
  const directives: Record<string, string[]> = { ...BASE_DIRECTIVES };

  // next/font inlines the fallback faces it generates as data: URIs, and a
  // policy without `data:` blocks every font on the page.
  directives["font-src"] = [...directives["font-src"], "data:"];

  if (options.serviceUrl && isBrowserReachable(options.serviceUrl)) {
    directives["img-src"] = [...directives["img-src"], options.serviceUrl];
    directives["font-src"] = [...directives["font-src"], options.serviceUrl];
  }

  if (options.iframeOrigins && options.iframeOrigins.length > 0) {
    directives["frame-ancestors"] = [...options.iframeOrigins];
  }

  return serializeCSP(directives);
}

function serializeCSP(directives: Record<string, string[]>): string {
  return Object.entries(directives)
    .map(([key, values]) => [key, ...values].join(" "))
    .join("; ");
}
