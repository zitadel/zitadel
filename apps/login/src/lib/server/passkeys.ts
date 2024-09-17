"use server";

import {
  createPasskeyRegistrationLink,
  getSession,
  registerPasskey,
  verifyPasskeyRegistration,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import {
  RegisterPasskeyResponse,
  VerifyPasskeyRegistrationRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";
import { userAgent } from "next/server";
import { getSessionCookieById } from "../cookies";

type VerifyPasskeyCommand = {
  passkeyId: string;
  passkeyName?: string;
  publicKeyCredential: any;
  sessionId: string;
};

type RegisterPasskeyCommand = {
  sessionId: string;
};

export async function registerPasskeyLink(
  command: RegisterPasskeyCommand,
): Promise<RegisterPasskeyResponse> {
  const { sessionId } = command;

  const sessionCookie = await getSessionCookieById({ sessionId });
  const session = await getSession({
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  const domain = headers().get("host");

  if (!domain) {
    throw new Error("Could not get domain");
  }

  const userId = session?.session?.factors?.user?.id;

  if (!userId) {
    throw new Error("Could not get session");
  }
  // TODO: add org context

  // use session token to add the passkey
  const registerLink = await createPasskeyRegistrationLink(
    userId,
    // sessionCookie.token,
  );

  if (!registerLink.code) {
    throw new Error("Missing code in response");
  }

  return registerPasskey(userId, registerLink.code, domain);
}

export async function verifyPasskey(command: VerifyPasskeyCommand) {
  // if no name is provided, try to generate one from the user agent
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
    throw new Error("Could not get session");
  }

  return verifyPasskeyRegistration(
    create(VerifyPasskeyRegistrationRequestSchema, {
      passkeyId: command.passkeyId,
      publicKeyCredential: command.publicKeyCredential,
      passkeyName,
      userId,
    }),
  );
}
