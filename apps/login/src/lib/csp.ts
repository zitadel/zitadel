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
  imageSources?: string[];
  iframeOrigins?: string[] | null;
}

interface ImageSourceOptions {
  serviceUrl?: string;
  publicHost?: string;
  instanceHost?: string;
  customRequestHeaders?: string;
}

export function buildCSP(options: CSPOptions = {}): string {
  const directives: Record<string, string[]> = { ...BASE_DIRECTIVES };

  const imageSources: string[] = [];
  const rawSources = [options.serviceUrl, ...(options.imageSources ?? [])];
  for (const source of rawSources) {
    imageSources.push(...normalizeImageSources(source));
  }

  if (imageSources.length > 0) {
    directives["img-src"] = dedupeSources([...directives["img-src"], ...imageSources]);
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

function normalizeImageSources(source?: string): string[] {
  if (!source) {
    return [];
  }

  const trimmed = source.trim();
  if (!trimmed) {
    return [];
  }

  let parsed: URL;

  if (/^https?:\/\//i.test(trimmed)) {
    try {
      parsed = new URL(trimmed);
    } catch {
      return [];
    }
  } else {
    const protocol = trimmed.includes("localhost") || trimmed.startsWith("127.0.0.1") ? "http://" : "https://";

    try {
      parsed = new URL(`${protocol}${trimmed}`);
    } catch {
      return [];
    }
  }

  const normalizedSources = [parsed.origin];

  if (parsed.port) {
    normalizedSources.push(`${parsed.protocol}//${parsed.hostname}`);
  }

  return dedupeSources(normalizedSources);
}

function dedupeSources(sources: string[]): string[] {
  const unique: string[] = [];

  for (const source of sources) {
    if (!unique.includes(source)) {
      unique.push(source);
    }
  }

  return unique;
}

export function resolveImageSources({
  serviceUrl,
  publicHost,
  instanceHost,
  customRequestHeaders,
}: ImageSourceOptions): string[] {
  const customHosts = extractHostsFromCustomHeaders(customRequestHeaders);
  return dedupeSources([...customHosts, publicHost || "", instanceHost || "", serviceUrl || ""].filter(Boolean));
}

function extractHostsFromCustomHeaders(customRequestHeaders?: string): string[] {
  if (!customRequestHeaders) {
    return [];
  }

  const allowedHeaderNames = new Set(["x-zitadel-public-host", "x-zitadel-instance-host", "x-zitadel-forward-host"]);
  const hosts: string[] = [];

  for (const header of customRequestHeaders.split(",")) {
    const separatorIndex = header.indexOf(":");
    if (separatorIndex <= 0) {
      continue;
    }

    const key = header.slice(0, separatorIndex).trim().toLowerCase();
    const value = header.slice(separatorIndex + 1).trim();
    if (!value || !allowedHeaderNames.has(key)) {
      continue;
    }

    hosts.push(value);
  }

  return dedupeSources(hosts);
}
