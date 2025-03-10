import FingerprintJS from "@fingerprintjs/fingerprintjs";
import { create } from "@zitadel/client";
import {
  UserAgent,
  UserAgentSchema,
} from "@zitadel/proto/zitadel/session/v2/session_pb";
import { cookies, headers } from "next/headers";
import { userAgent } from "next/server";

export async function getFingerprintId() {
  const fp = await FingerprintJS.load();
  const result = await fp.get();
  return result.visitorId;
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

  const { device } = userAgent({ headers: _headers });

  const userAgentHeader = _headers.get("user-agent");

  const userAgentHeaderValues = userAgentHeader?.split(",");

  const userAgentData: UserAgent = create(UserAgentSchema, {
    ip: _headers.get("x-forwarded-for") ?? _headers.get("remoteAddress") ?? "",
    header: { "user-agent": { values: userAgentHeaderValues } },
    description: `${device?.type ?? "unknown type"}, ${device?.vendor ?? "unknown vendor"} ${device?.model ?? "unknown model"}`,
    fingerprintId: fingerprintId,
  });

  return userAgentData;
}
