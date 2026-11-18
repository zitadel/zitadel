import "server-only";

import { create } from "@zitadel/client";
import { UserAgent, UserAgentSchema } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { isIP } from "net";
import { ReadonlyHeaders } from "next/dist/server/web/spec-extension/adapters/headers";
import { cookies, headers } from "next/headers";
import { userAgent } from "next/server";
import { v4 as uuidv4 } from "uuid";
import { createIdpFlowProof } from "./idp-flow-proof";
import { createLogger } from "./logger";
import { getServiceConfig } from "./service-url";
import { getSecuritySettings } from "./zitadel";

const logger = createLogger("fingerprint");

const BINDING_ID_COOKIE = "bindingId";

// sets an httpOnly browser cookie (1 year by default) that also works when the login is embedded in an iframe
export async function setBrowserCookie(name: string, value: string, maxAge: number = 31536000) {
  const cookiesList = await cookies();

  let iFrameEnabled = false;
  try {
    const { serviceConfig } = getServiceConfig(await headers());
    const securitySettings = await getSecuritySettings({ serviceConfig });
    iFrameEnabled = !!securitySettings?.embeddedIframe?.enabled;
  } catch (error) {
    logger.warn(`Could not load security settings, setting the ${name} cookie without iframe support`, { error });
  }

  return cookiesList.set({
    name,
    value,
    httpOnly: true,
    path: "/",
    maxAge,
    // "none" is required for iframe embedding (with secure flag), "lax" allows cookies during top-level navigation
    sameSite: iFrameEnabled ? "none" : "lax",
    // "none" is only accepted by browsers together with "secure", otherwise the cookie is silently dropped
    secure: process.env.NODE_ENV === "production" || iFrameEnabled,
  });
}

// returns the value of the cookie, or sets it to a new random value
async function getOrSetBrowserCookie(name: string): Promise<string> {
  const cookie = (await cookies()).get(name);
  if (cookie) {
    return cookie.value;
  }

  const value = uuidv4();
  await setBrowserCookie(name, value);
  return value;
}

export async function getOrSetFingerprintId(): Promise<string> {
  return getOrSetBrowserCookie("fingerprintId");
}

export async function getBindingIdCookie() {
  const cookiesList = await cookies();
  return cookiesList.get(BINDING_ID_COOKIE);
}

/**
 * Returns the random binding id of this browser, which the signed proofs of the login (see
 * signed-proof.ts) are bound to. Unlike the fingerprintId it is never sent to the API, so it can't
 * be read back through a session and stays a secret of this browser.
 */
export async function getOrSetBindingId(): Promise<string> {
  return getOrSetBrowserCookie(BINDING_ID_COOKIE);
}

// every external IdP flow start calls this, so no start can forget the browser binding
export async function addIdpFlowProof(params: URLSearchParams, linkToSessionId?: string): Promise<void> {
  params.set("flowFingerprint", createIdpFlowProof(await getOrSetBindingId(), linkToSessionId));
}

const DEFAULT_TRUSTED_PROXY_HOPS = 0;

/**
 * Number of proxies between the client and the login that append their own address to
 * `X-Forwarded-For`. The default of 0 fits a single reverse proxy and Vercel, which write the
 * client address as the last entry; the GCP load balancer appends its own address, so it needs 1.
 */
function trustedProxyHops(): number {
  const configured = process.env.TRUSTED_PROXY_HOPS?.trim();
  if (!configured) {
    return DEFAULT_TRUSTED_PROXY_HOPS;
  }
  const hops = Number(configured);
  if (!Number.isInteger(hops) || hops < 0) {
    logger.warn("Ignoring invalid TRUSTED_PROXY_HOPS, falling back to the default", { configured });
    return DEFAULT_TRUSTED_PROXY_HOPS;
  }
  return hops;
}

// drops the port and IPv6 brackets, and anything that is not an address at all ("unknown", RFC 7239 obfuscated names)
function toIp(value: string): string {
  const trimmed = value.trim();
  if (isIP(trimmed)) {
    return trimmed; // IPv4, or IPv6 without brackets (which may not carry a port)
  }
  const bracketed = /^\[(.+)\](?::\d+)?$/.exec(trimmed);
  const withoutPort = bracketed ? bracketed[1] : trimmed.replace(/:\d+$/, "");
  return isIP(withoutPort) ? withoutPort : "";
}

/**
 * The client address, as far as it can be trusted. Every `X-Forwarded-For` entry left of the
 * trusted proxies was written by whoever sent the request, so it is spoofable and never used.
 * Returns "" when no address can be trusted: the backend then omits the IP from the session event
 * instead of recording an attacker-supplied one.
 */
function getClientIp(requestHeaders: ReadonlyHeaders): string {
  const forwardedFor = requestHeaders.get("x-forwarded-for");
  if (!forwardedFor) {
    // nothing in front of the login, or a front that doesn't set the header
    return toIp(requestHeaders.get("remoteAddress") ?? "");
  }
  // shortcut: a hop count rather than a list of trusted proxy addresses, which is all a fixed
  // deployment needs; match on addresses once a front can vary the number of hops per request
  return toIp(forwardedFor.split(",").at(-1 - trustedProxyHops()) ?? "");
}

export async function getUserAgent(): Promise<UserAgent> {
  const _headers = await headers();

  const fingerprintId = await getOrSetFingerprintId();

  const { device, engine, os, browser } = userAgent({ headers: _headers });

  const userAgentHeader = _headers.get("user-agent");

  const userAgentHeaderValues = userAgentHeader?.split(",");

  const deviceDescription = `${device?.type ? `${device.type},` : ""} ${device?.vendor ? `${device.vendor},` : ""} ${device.model ? `${device.model},` : ""} `;
  const osDescription = `${os?.name ? `${os.name},` : ""} ${os?.version ? `${os.version},` : ""} `;
  const engineDescription = `${engine?.name ? `${engine.name},` : ""} ${engine?.version ? `${engine.version},` : ""} `;
  const browserDescription = `${browser?.name ? `${browser.name},` : ""} ${browser.version ? `${browser.version},` : ""} `;

  const userAgentData: UserAgent = create(UserAgentSchema, {
    ip: getClientIp(_headers),
    header: { "user-agent": { values: userAgentHeaderValues } },
    description: `${browserDescription}, ${deviceDescription}, ${engineDescription}, ${osDescription}`,
    fingerprintId: fingerprintId,
  });

  return userAgentData;
}
