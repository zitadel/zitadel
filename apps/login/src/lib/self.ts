"use server";

import { createUserServiceClient } from "@zitadel/client/v2";
import { createServerTransport } from "@zitadel/node";
import { getSessionCookieById } from "./cookies";
import { getSession } from "./zitadel";

const transport = (token: string) =>
  createServerTransport(token, {
    baseUrl: process.env.ZITADEL_API_URL!,
  });

const myUserService = (sessionToken: string) => {
  return createUserServiceClient(transport(sessionToken));
};

export async function setMyPassword({
  sessionId,
  password,
}: {
  sessionId: string;
  password: string;
}) {
  const sessionCookie = await getSessionCookieById({ sessionId });

  const { session } = await getSession({
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  if (!session) {
    return { error: "Could not load session" };
  }

  const service = await myUserService(`${sessionCookie.token}`);

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
