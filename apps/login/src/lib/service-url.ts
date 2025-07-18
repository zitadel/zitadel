import { ReadonlyHeaders } from "next/dist/server/web/spec-extension/adapters/headers";
import { NextRequest } from "next/server";

/**
 * Extracts the service url and region from the headers if used in a multitenant context (host, x-zitadel-forward-host header)
 * or falls back to the ZITADEL_API_URL for a self hosting deployment
 * or falls back to the host header for a self hosting deployment using custom domains
 * @param headers
 * @returns the service url and region from the headers
 * @throws if the service url could not be determined
 *
 */
export function getServiceUrlFromHeaders(headers: ReadonlyHeaders): {
  serviceUrl: string;
} {
  let instanceUrl;

  const forwardedHost = headers.get("x-zitadel-forward-host");
  // use the forwarded host if available (multitenant), otherwise fall back to the host of the deployment itself
  if (forwardedHost) {
    instanceUrl = forwardedHost;
    instanceUrl = instanceUrl.startsWith("http://")
      ? instanceUrl
      : `https://${instanceUrl}`;
  } else if (process.env.ZITADEL_API_URL) {
    instanceUrl = process.env.ZITADEL_API_URL;
  } else {
    const host = headers.get("host");

    if (host) {
      const [hostname] = host.split(":");
      if (hostname !== "localhost") {
        instanceUrl = host.startsWith("http") ? host : `https://${host}`;
      }
    }
  }

  if (!instanceUrl) {
    throw new Error("Service URL could not be determined");
  }

  return {
    serviceUrl: instanceUrl,
  };
}

export function constructUrl(request: NextRequest, path: string) {
  const forwardedProto = request.headers.get("x-forwarded-proto")
    ? `${request.headers.get("x-forwarded-proto")}:`
    : request.nextUrl.protocol;

  const forwardedHost =
    request.headers.get("x-zitadel-forward-host") ??
    request.headers.get("x-forwarded-host") ??
    request.headers.get("host");
  const basePath = process.env.NEXT_PUBLIC_BASE_PATH || "";
  return new URL(`${basePath}${path}`, `${forwardedProto}//${forwardedHost}`);
}
