"use server";

import { createUserServiceClient } from "@zitadel/client/v2";
import { headers } from "next/headers";
import { getSessionCookieById } from "./cookies";
import { getServiceUrlFromHeaders } from "./service-url";
import { createServerTransport, getSession } from "./zitadel";

const myUserService = async (serviceUrl: string, sessionToken: string) => {
  const transportPromise = await createServerTransport(
    sessionToken,
    serviceUrl,
  );
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
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const sessionCookie = await getSessionCookieById({ sessionId });

  const { session } = await getSession({
    serviceUrl,
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  if (!session) {
    return { error: "Could not load session" };
  }

  const service = await myUserService(serviceUrl, `${sessionCookie.token}`);

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
