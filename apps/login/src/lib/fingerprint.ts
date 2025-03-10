import { create } from "@zitadel/client";
import {
  UserAgent,
  UserAgentSchema,
} from "@zitadel/proto/zitadel/session/v2/session_pb";
import { cookies, headers } from "next/headers";
import { userAgent } from "next/server";
import { v4 as uuidv4 } from "uuid";

export async function getFingerprintId() {
  return uuidv4();
}

export async function setFingerprintIdCookie(fingerprintId: string) {
  const cookiesList = await cookies();

  return cookiesList.set({
    name: "fingerprintId",
    value: fingerprintId,
    httpOnly: true,
    path: "/",
    maxAge: 31536000, // 1 year
  });
}

export async function getFingerprintIdCookie() {
  const cookiesList = await cookies();
  return cookiesList.get("fingerprintId");
}

export async function getOrSetFingerprintId(): Promise<string> {
  const cookie = await getFingerprintIdCookie();
  if (cookie) {
    return cookie.value;
  }

  const fingerprintId = await getFingerprintId();
  await setFingerprintIdCookie(fingerprintId);
  return fingerprintId;
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
    ip: _headers.get("x-forwarded-for") ?? _headers.get("remoteAddress") ?? "",
    header: { "user-agent": { values: userAgentHeaderValues } },
    description: `${browserDescription}, ${deviceDescription}, ${engineDescription}, ${osDescription}`,
    fingerprintId: fingerprintId,
  });

  return userAgentData;
}
