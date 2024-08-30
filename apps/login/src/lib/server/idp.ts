"use server";

import { startIdentityProviderFlow } from "@/lib/zitadel";

export type StartIDPFlowCommand = {
  idpId: string;
  successUrl: string;
  failureUrl: string;
};

export async function startIDPFlow(command: StartIDPFlowCommand) {
  const { idpId, successUrl, failureUrl } = command;

  return startIdentityProviderFlow({
    idpId,
    urls: {
      successUrl,
      failureUrl,
    },
  });
}
