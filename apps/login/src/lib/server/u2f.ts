"use server";

import { getSession, registerU2F, verifyU2FRegistration } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { VerifyU2FRegistrationRequestSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";
import { userAgent } from "next/server";
import { getSessionCookieById } from "../cookies";

type RegisterU2FCommand = {
  sessionId: string;
};

type VerifyU2FCommand = {
  u2fId: string;
  passkeyName?: string;
  publicKeyCredential: any;
  sessionId: string;
};

export async function addU2F(command: RegisterU2FCommand) {
  const sessionCookie = await getSessionCookieById({
    sessionId: command.sessionId,
  });

  if (!sessionCookie) {
    return { error: "Could not get session" };
  }

  const session = await getSession({
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  const domain = headers().get("host");

  if (!domain) {
    return { error: "Could not get domain" };
  }

  const userId = session?.session?.factors?.user?.id;

  if (!session || !userId) {
    return { error: "Could not get session" };
  }

  return registerU2F(
    userId,
    domain,
    // sessionCookie.token
  );
}

export async function verifyU2F(command: VerifyU2FCommand) {
  let passkeyName = command.passkeyName;
  if (!!!passkeyName) {
    const headersList = headers();
    const userAgentStructure = { headers: headersList };
    const { browser, device, os } = userAgent(userAgentStructure);

    passkeyName = `${device.vendor ?? ""} ${device.model ?? ""}${
      device.vendor || device.model ? ", " : ""
    }${os.name}${os.name ? ", " : ""}${browser.name}`;
  }
  const sessionCookie = await getSessionCookieById({
    sessionId: command.sessionId,
  });

  const session = await getSession({
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  const userId = session?.session?.factors?.user?.id;

  if (!userId) {
    return { error: "Could not get session" };
  }

  const req = create(VerifyU2FRegistrationRequestSchema, {
    u2fId: command.u2fId,
    publicKeyCredential: command.publicKeyCredential,
    tokenName: passkeyName,
    userId,
  });

  return verifyU2FRegistration(req);
}
