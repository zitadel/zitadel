"use server";

import { createServerTransport } from "@zitadel/client/node";
import { createUserServiceClient } from "@zitadel/client/v2";
import { headers } from "next/headers";
import { getSessionCookieById } from "./cookies";
import { getApiUrlOfHeaders } from "./service";
import { getSession } from "./zitadel";

const transport = async (host: string, token: string) => {
  return createServerTransport(token, {
    baseUrl: host,
  });
};

const myUserService = async (host: string, sessionToken: string) => {
  const transportPromise = await transport(host, sessionToken);
  return createUserServiceClient(transportPromise);
};

export async function setMyPassword({
  sessionId,
  password,
}: {
  sessionId: string;
  password: string;
}) {
  const _headers = await headers();
  const instanceUrl = getApiUrlOfHeaders(_headers);

  if (!instanceUrl) {
    throw new Error("No host found");
  }

  const sessionCookie = await getSessionCookieById({ sessionId });

  const { session } = await getSession({
    host: instanceUrl,
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  if (!session) {
    return { error: "Could not load session" };
  }

  const service = await myUserService(instanceUrl, `${sessionCookie.token}`);

  if (!session?.factors?.user?.id) {
    return { error: "No user id found in session" };
  }

  return service
    .setPassword(
      {
        userId: session.factors.user.id,
        newPassword: { password, changeRequired: false },
      },
      {},
    )
    .catch((error) => {
      console.log(error);
      if (error.code === 7) {
        return { error: "Session is not valid." };
      }
      throw error;
    });
}
