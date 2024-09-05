"use server";

import { headers } from "next/headers";
import { idpTypeToSlug } from "../idp";
import {
  getActiveIdentityProviders,
  getLoginSettings,
  listAuthenticationMethodTypes,
  listUsers,
  startIdentityProviderFlow,
} from "../zitadel";
import { createSessionForUserIdAndUpdateCookie } from "../../utils/session";
import { redirect } from "next/navigation";

export type SendLoginnameCommand = {
  loginName: string;
  authRequestId?: string;
  organization?: string;
};

export async function sendLoginname(command: SendLoginnameCommand) {
  const users = await listUsers({
    userName: command.loginName,
    organizationId: command.organization,
  });

  if (users.details?.totalResult == BigInt(1) && users.result[0].userId) {
    const userId = users.result[0].userId;
    const session = await createSessionForUserIdAndUpdateCookie(
      userId,
      undefined,
      undefined,
      command.authRequestId,
    );

    if (!session.factors?.user?.id) {
      throw Error("Could not create session for user");
    }

    const methods = await listAuthenticationMethodTypes(
      session.factors?.user?.id,
    );

    return {
      authMethodTypes: methods.authMethodTypes,
      sessionId: session.id,
      factors: session.factors,
    };
  }

  const loginSettings = await getLoginSettings(command.organization);
  // TODO: check if allowDomainDiscovery has to be allowed too, to redirect to the register page
  // user not found, check if register is enabled on organization

  if (loginSettings?.allowRegister && !loginSettings?.allowUsernamePassword) {
    // TODO redirect to loginname page with idp hint
    const identityProviders = await getActiveIdentityProviders(
      command.organization,
    ).then((resp) => {
      return resp.identityProviders;
    });

    if (identityProviders.length === 1) {
      const host = headers().get("host");
      console.log("host", host);
      const identityProviderType = identityProviders[0].type;

      const provider = idpTypeToSlug(identityProviderType);

      const params = new URLSearchParams();

      if (command.authRequestId) {
        params.set("authRequestId", command.authRequestId);
      }

      if (command.organization) {
        params.set("organization", command.organization);
      }

      return startIdentityProviderFlow({
        idpId: identityProviders[0].id,
        urls: {
          successUrl:
            `${host}/idp/${provider}/success?` + new URLSearchParams(params),
          failureUrl:
            `${host}/idp/${provider}/failure?` + new URLSearchParams(params),
        },
      }).then((resp: any) => {
        if (resp.authUrl) {
          return redirect(resp.authUrl);
        }
      });
    } else {
      throw Error("Could not find user");
    }
  } else if (
    loginSettings?.allowRegister &&
    loginSettings?.allowUsernamePassword
  ) {
    const params = new URLSearchParams();

    if (command.organization) {
      params.set("organization", command.organization);
    }
    if (command.authRequestId) {
      params.set("authRequestId", command.authRequestId);
    }
    if (command.loginName) {
      params.set("loginName", command.loginName);
    }

    const registerUrl = new URL("/register?" + params);

    return redirect(registerUrl.toString());
  }

  throw Error("Could not find user");
}
