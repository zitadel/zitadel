import { ServiceConfig } from "./zitadel";

function applyServicePortToHost(host: string, serviceBaseUrl: string): string {
  let serviceUrl: URL;
  try {
    serviceUrl = new URL(serviceBaseUrl);
  } catch {
    return host;
  }

  const normalizedHost = host.replace(/^https?:\/\//i, "").split("/")[0];
  if (!normalizedHost) {
    return host;
  }

  if (!serviceUrl.port) {
    return normalizedHost;
  }

  if (normalizedHost.includes(":")) {
    return normalizedHost;
  }

  if (normalizedHost !== serviceUrl.hostname) {
    return normalizedHost;
  }

  return `${normalizedHost}:${serviceUrl.port}`;
}

export function getLogoutTokenVerificationOptions(
  serviceConfig: ServiceConfig,
  customRequestHeaders: string | undefined = process.env.CUSTOM_REQUEST_HEADERS,
): { instanceHost?: string; publicHost?: string } | undefined {
  // In self-host mode, `publicHost` can be the login UI host. Don't use it for JWKS unless it is explicitly configured.
  let instanceHost = serviceConfig.instanceHost;
  let publicHost = serviceConfig.instanceHost ? serviceConfig.publicHost : undefined;

  if (customRequestHeaders) {
    for (const header of customRequestHeaders.split(",")) {
      const separatorIndex = header.indexOf(":");
      if (separatorIndex <= 0) {
        continue;
      }

      const key = header.slice(0, separatorIndex).trim().toLowerCase();
      const rawValue = header.slice(separatorIndex + 1).trim();
      const value = rawValue ? applyServicePortToHost(rawValue, serviceConfig.baseUrl) : undefined;

      if (key === "x-zitadel-instance-host") {
        instanceHost = value;
      }
      if (key === "x-zitadel-public-host") {
        publicHost = value;
      }
      if (key === "x-zitadel-forward-host" && value) {
        if (!instanceHost) {
          instanceHost = value;
        }
        if (!publicHost) {
          publicHost = value;
        }
      }
    }
  }

  if (!instanceHost && !publicHost) {
    return undefined;
  }

  return {
    ...(instanceHost ? { instanceHost } : {}),
    ...(publicHost ? { publicHost } : {}),
  };
}
