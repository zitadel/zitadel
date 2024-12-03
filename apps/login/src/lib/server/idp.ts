"use server";

import { startIdentityProviderFlow } from "@/lib/zitadel";
import { headers } from "next/headers";

export type StartIDPFlowCommand = {
  idpId: string;
  successUrl: string;
  failureUrl: string;
};

export async function startIDPFlow(command: StartIDPFlowCommand) {
  const host = (await headers()).get("host");

  return startIdentityProviderFlow({
    idpId: command.idpId,
    urls: {
      successUrl: `${host}${command.successUrl}`,
      failureUrl: `${host}${command.failureUrl}`,
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
