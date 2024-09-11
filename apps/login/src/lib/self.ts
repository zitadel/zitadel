"use server";

import {
  createSessionServiceClient,
  createUserServiceClient,
} from "@zitadel/client/v2";
import { createServerTransport } from "@zitadel/node";
import { getSessionCookieById } from "./cookies";

const transport = (token: string) =>
  createServerTransport(token, {
    baseUrl: process.env.ZITADEL_API_URL!,
    httpVersion: "2",
  });

const sessionService = (sessionId: string) => {
  return getSessionCookieById({ sessionId }).then((session) => {
    return createSessionServiceClient(transport(session.token));
  });
};

const userService = (sessionId: string) => {
  return getSessionCookieById({ sessionId }).then((session) => {
    return createUserServiceClient(transport(session.token));
  });
};

export async function setPassword({
  sessionId,
  userId,
  password,
}: {
  sessionId: string;
  userId: string;
  password: string;
}) {
  return (await userService(sessionId)).setPassword(
    {
      userId,
      newPassword: { password, changeRequired: false },
    },
    {},
  );
}
