"use server";

import { createSessionAndUpdateCookie } from "@/lib/server/cookie";
import { addHumanUser } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";
import {
  ChecksJson,
  ChecksSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { redirect } from "next/navigation";

type RegisterUserCommand = {
  email: string;
  firstName: string;
  lastName: string;
  password?: string;
  organization?: string;
  authRequestId?: string;
};

export type RegisterUserResponse = {
  userId: string;
  sessionId: string;
  factors: Factors | undefined;
};
export async function registerUser(command: RegisterUserCommand) {
  const human = await addHumanUser({
    email: command.email,
    firstName: command.firstName,
    lastName: command.lastName,
    password: command.password ? command.password : undefined,
    organization: command.organization,
  });

  if (!human) {
    return { error: "Could not create user" };
  }

  let checkPayload: any = {
    user: { search: { case: "userId", value: human.userId } },
  };

  if (command.password) {
    checkPayload = {
      ...checkPayload,
      password: { password: command.password },
    } as ChecksJson;
  }

  const checks = create(ChecksSchema, checkPayload);

  const session = await createSessionAndUpdateCookie(
    checks,
    undefined,
    command.authRequestId,
  ).then((session) => {
    return {
      userId: human.userId,
      sessionId: session.id,
      factors: session.factors,
    };
  });

  if (!session || !session.factors?.user) {
    return { error: "Could not create session" };
  }

  if (!command.password) {
    const params = new URLSearchParams({
      loginName: session.factors.user.loginName,
      organization: session.factors.user.organizationId,
    });

    if (command.authRequestId) {
      params.append("authRequestId", command.authRequestId);
    }

    return redirect("/passkey/set?" + params);
  } else {
    const params = new URLSearchParams({
      loginName: session.factors.user.loginName,
      organization: session.factors.user.organizationId,
    });

    if (command.authRequestId && session.userId) {
      params.append("authRequest", command.authRequestId);
      params.append("userId", session.userId);

      return redirect("/login?" + params);
    } else {
      return redirect("/signedin?" + params);
    }
  }
}
