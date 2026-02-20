
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
  iframeOrigins?: string[];
}

export function buildCSP(options: CSPOptions = {}): string {
  const directives: Record<string, string[]> = { ...BASE_DIRECTIVES };

  if (options.serviceUrl) {
    directives["img-src"] = [...directives["img-src"], options.serviceUrl];
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
