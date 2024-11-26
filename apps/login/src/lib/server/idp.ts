"use server";

import { startIdentityProviderFlow } from "@/lib/zitadel";

export type StartIDPFlowCommand = {
  idpId: string;
  successUrl: string;
  failureUrl: string;
};

export async function startIDPFlow(command: StartIDPFlowCommand) {
  return startIdentityProviderFlow({
    idpId: command.idpId,
    urls: {
      successUrl: command.successUrl,
      failureUrl: command.failureUrl,
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
