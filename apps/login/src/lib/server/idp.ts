"use server";

import { startIdentityProviderFlow } from "@/lib/zitadel";
import { headers } from "next/headers";

export type StartIDPFlowCommand = {
  idpId: string;
  successUrl: string;
  failureUrl: string;
};

export async function startIDPFlow(command: StartIDPFlowCommand) {
  const origin = (await headers()).get("origin");

  if (!origin) {
    return { error: "Could not get origin" };
  }

  return startIdentityProviderFlow({
    idpId: command.idpId,
    urls: {
      successUrl: `${origin}${command.successUrl}`,
      failureUrl: `${origin}${command.failureUrl}`,
    },
  }).then((response) => {
    if (
      response &&
      response.nextStep.case === "authUrl" &&
      response?.nextStep.value
    ) {
      return { redirect: response.nextStep.value };
    }
  });
}
