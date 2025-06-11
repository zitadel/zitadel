"use server";

import { addHumanUser, createInviteCode } from "@/lib/zitadel";
import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { headers } from "next/headers";
import { getServiceUrlFromHeaders } from "../service-url";

type InviteUserCommand = {
  email: string;
  firstName: string;
  lastName: string;
  password?: string;
  organization: string;
  requestId?: string;
};

export type RegisterUserResponse = {
  userId: string;
  sessionId: string;
  factors: Factors | undefined;
};

export async function inviteUser(command: InviteUserCommand) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);
  const host = _headers.get("host");

  if (!host) {
    return { error: "Could not get domain" };
  }

  const human = await addHumanUser({
    serviceUrl,
    email: command.email,
    firstName: command.firstName,
    lastName: command.lastName,
    password: command.password ? command.password : undefined,
    organization: command.organization,
  });

  if (!human) {
    return { error: "Could not create user" };
  }

  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  const codeResponse = await createInviteCode({
    serviceUrl,
    urlTemplate: `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&invite=true`,
    userId: human.userId,
  });

  if (!codeResponse || !human) {
    return { error: "Could not create invite code" };
  }

  return human.userId;
}
