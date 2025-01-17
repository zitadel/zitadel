"use server";

import { createServerTransport } from "@zitadel/client/node";
import { createUserServiceClient } from "@zitadel/client/v2";
import { headers } from "next/headers";
import { getInstanceUrl } from "./api";
import { getSessionCookieById } from "./cookies";
import { getSession } from "./zitadel";

const transport = async (host: string, token: string) => {
  let instanceUrl;
  try {
    instanceUrl = await getInstanceUrl(host);
  } catch (error) {
    console.error(
      `Could not get instance url for ${host}, fallback to ZITADEL_API_URL`,
      error,
    );
    instanceUrl = process.env.ZITADEL_API_URL;
  }

  return createServerTransport(token, {
    baseUrl: instanceUrl,
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
  const host = (await headers()).get("host");

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

  const sessionCookie = await getSessionCookieById({ sessionId });

  const { session } = await getSession({
    host,
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  if (!session) {
    return { error: "Could not load session" };
  }

  const service = await myUserService(host, `${sessionCookie.token}`);

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
