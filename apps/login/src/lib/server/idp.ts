"use server";

import { startIdentityProviderFlow } from "@/lib/zitadel";

export type StartIDPFlowOptions = {
  idpId: string;
  successUrl: string;
  failureUrl: string;
};
export async function startIDPFlow(options: StartIDPFlowOptions) {
  const { idpId, successUrl, failureUrl } = options;

  return startIdentityProviderFlow({
    idpId,
    urls: {
      successUrl,
      failureUrl,
    },
  });
}
