"use server";

import { createUserServiceClient } from "@zitadel/client/v2";
import { headers } from "next/headers";
import { getSessionCookieById } from "./cookies";
import { getServiceConfig } from "./service-url";
import { createServerTransport, getSession, ServiceConfig } from "./zitadel";

const myUserService = async (serviceConfig: ServiceConfig, sessionToken: string) => {
  const transportPromise = await createServerTransport(sessionToken, serviceConfig);
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
  const { serviceConfig } = getServiceConfig(_headers);

  const sessionCookie = await getSessionCookieById({ sessionId });

  const { session } = await getSession({ serviceConfig, sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  if (!session) {
    return { error: "Could not load session" };
  }

  const service = await myUserService(serviceConfig, `${sessionCookie.token}`);

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
