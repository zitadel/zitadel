"use server";

import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";
import { redirect } from "next/navigation";
import { idpTypeToIdentityProviderType, idpTypeToSlug } from "../idp";

import {
  getActiveIdentityProviders,
  getIDPByID,
  getLoginSettings,
  getOrgsByDomain,
  listAuthenticationMethodTypes,
  listIDPLinks,
  listUsers,
  startIdentityProviderFlow,
} from "../zitadel";
import { createSessionAndUpdateCookie } from "./cookie";

export type SendLoginnameCommand = {
  loginName: string;
  authRequestId?: string;
  organization?: string;
};

const ORG_SUFFIX_REGEX = /(?<=@)(.+)/;

export async function sendLoginname(command: SendLoginnameCommand) {
  const users = await listUsers({
    loginName: command.loginName,
    organizationId: command.organization,
  });

  const loginSettings = await getLoginSettings(command.organization);

  const redirectUserToSingleIDPIfAvailable = async () => {
    const identityProviders = await getActiveIdentityProviders(
      command.organization,
    ).then((resp) => {
      return resp.identityProviders;
    });

    if (identityProviders.length === 1) {
      const host = (await headers()).get("host");
      const identityProviderType = identityProviders[0].type;

      const provider = idpTypeToSlug(identityProviderType);

      const params = new URLSearchParams();

      if (command.authRequestId) {
        params.set("authRequestId", command.authRequestId);
      }

      if (command.organization) {
        params.set("organization", command.organization);
      }

      const resp = await startIdentityProviderFlow({
        idpId: identityProviders[0].id,
        urls: {
          successUrl:
            `${host}/idp/${provider}/success?` + new URLSearchParams(params),
          failureUrl:
            `${host}/idp/${provider}/failure?` + new URLSearchParams(params),
        },
      });

      if (resp?.nextStep.case === "authUrl") {
        return redirect(resp.nextStep.value);
      }
    }
  };

  const redirectUserToIDP = async (userId: string) => {
    const identityProviders = await listIDPLinks(userId).then((resp) => {
      return resp.result;
    });

    if (identityProviders.length === 1) {
      const host = (await headers()).get("host");
      const identityProviderId = identityProviders[0].idpId;

      const idp = await getIDPByID(identityProviderId);
      const idpType = idp?.type;

      if (!idp || !idpType) {
        throw new Error("Could not find identity provider");
      }

      const identityProviderType = idpTypeToIdentityProviderType(idpType);
      const provider = idpTypeToSlug(identityProviderType);

      const params = new URLSearchParams();

      if (command.authRequestId) {
        params.set("authRequestId", command.authRequestId);
      }

      if (command.organization) {
        params.set("organization", command.organization);
      }

      const resp = await startIdentityProviderFlow({
        idpId: idp.id,
        urls: {
          successUrl:
            `${host}/idp/${provider}/success?` + new URLSearchParams(params),
          failureUrl:
            `${host}/idp/${provider}/failure?` + new URLSearchParams(params),
        },
      });

      if (resp?.nextStep.case === "authUrl") {
        return redirect(resp.nextStep.value);
      }
    }
  };

  if (users.details?.totalResult == BigInt(1) && users.result[0].userId) {
    const userId = users.result[0].userId;

    const checks = create(ChecksSchema, {
      user: { search: { case: "userId", value: userId } },
    });

    const session = await createSessionAndUpdateCookie(
      checks,
      undefined,
      command.authRequestId,
    );

    if (!session.factors?.user?.id) {
      return { error: "Could not create session for user" };
    }

    if (users.result[0].state === UserState.INITIAL) {
      const params = new URLSearchParams({
        loginName: session.factors?.user?.loginName,
        initial: "true", // this does not require a code to be set
      });

      if (command.organization || session.factors?.user?.organizationId) {
        params.append(
          "organization",
          command.organization ?? session.factors?.user?.organizationId,
        );
      }

      if (command.authRequestId) {
        params.append("authRequestid", command.authRequestId);
      }

      return redirect("/password/set?" + params);
    }

    const methods = await listAuthenticationMethodTypes(
      session.factors?.user?.id,
    );

    if (!methods.authMethodTypes || !methods.authMethodTypes.length) {
      if (
        users.result[0].type.case === "human" &&
        users.result[0].type.value.email &&
        !users.result[0].type.value.email.isVerified
      ) {
        const paramsVerify = new URLSearchParams({
          loginName: session.factors?.user?.loginName,
          userId: session.factors?.user?.id, // verify needs user id
          invite: "true", // TODO: check - set this to true as we dont expect old email verification method here
        });

        if (command.organization || session.factors?.user?.organizationId) {
          paramsVerify.append(
            "organization",
            command.organization ?? session.factors?.user?.organizationId,
          );
        }

        if (command.authRequestId) {
          paramsVerify.append("authRequestId", command.authRequestId);
        }

        redirect("/verify?" + paramsVerify);
      }

      const paramsAuthenticatorSetup = new URLSearchParams({
        loginName: session.factors?.user?.loginName,
        userId: session.factors?.user?.id, // verify needs user id
      });

      if (command.organization || session.factors?.user?.organizationId) {
        paramsAuthenticatorSetup.append(
          "organization",
          command.organization ?? session.factors?.user?.organizationId,
        );
      }

      if (command.authRequestId) {
        paramsAuthenticatorSetup.append("authRequestId", command.authRequestId);
      }

      redirect("/authenticator/set?" + paramsAuthenticatorSetup);
    }

    if (methods.authMethodTypes.length == 1) {
      const method = methods.authMethodTypes[0];
      switch (method) {
        case AuthenticationMethodType.PASSWORD: // user has only password as auth method
          const paramsPassword: any = {
            loginName: session.factors?.user?.loginName,
          };

          // TODO: does this have to be checked in loginSettings.allowDomainDiscovery

          if (command.organization || session.factors?.user?.organizationId) {
            paramsPassword.organization =
              command.organization ?? session.factors?.user?.organizationId;
          }

          if (command.authRequestId) {
            paramsPassword.authRequestId = command.authRequestId;
          }

          return redirect("/password?" + new URLSearchParams(paramsPassword));
        case AuthenticationMethodType.PASSKEY: // AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSKEY
          const paramsPasskey: any = { loginName: command.loginName };
          if (command.authRequestId) {
            paramsPasskey.authRequestId = command.authRequestId;
          }

          if (command.organization || session.factors?.user?.organizationId) {
            paramsPasskey.organization =
              command.organization ?? session.factors?.user?.organizationId;
          }

          return redirect("/passkey?" + new URLSearchParams(paramsPasskey));
      }
    } else {
      // prefer passkey in favor of other methods
      if (methods.authMethodTypes.includes(AuthenticationMethodType.PASSKEY)) {
        const passkeyParams: any = {
          loginName: command.loginName,
          altPassword: `${methods.authMethodTypes.includes(1)}`, // show alternative password option
        };

        if (command.authRequestId) {
          passkeyParams.authRequestId = command.authRequestId;
        }

        if (command.organization || session.factors?.user?.organizationId) {
          passkeyParams.organization =
            command.organization ?? session.factors?.user?.organizationId;
        }

        return redirect("/passkey?" + new URLSearchParams(passkeyParams));
      } else if (
        methods.authMethodTypes.includes(AuthenticationMethodType.IDP)
      ) {
        await redirectUserToIDP(userId);
      } else if (
        methods.authMethodTypes.includes(AuthenticationMethodType.PASSWORD)
      ) {
        // user has no passkey setup and login settings allow passkeys
        const paramsPasswordDefault: any = { loginName: command.loginName };

        if (command.authRequestId) {
          paramsPasswordDefault.authRequestId = command.authRequestId;
        }

        if (command.organization || session.factors?.user?.organizationId) {
          paramsPasswordDefault.organization =
            command.organization ?? session.factors?.user?.organizationId;
        }

        return redirect(
          "/password?" + new URLSearchParams(paramsPasswordDefault),
        );
      }
    }
  }

  // user not found, check if register is enabled on organization
  if (loginSettings?.allowRegister && !loginSettings?.allowUsernamePassword) {
    // TODO: do we need to handle login hints for IDPs here?
    await redirectUserToSingleIDPIfAvailable();

    return { error: "Could not find user" };
  } else if (
    loginSettings?.allowRegister &&
    loginSettings?.allowUsernamePassword
  ) {
    let orgToRegisterOn: string | undefined = command.organization;

    if (
      !loginSettings?.ignoreUnknownUsernames &&
      !orgToRegisterOn &&
      command.loginName &&
      ORG_SUFFIX_REGEX.test(command.loginName)
    ) {
      const matched = ORG_SUFFIX_REGEX.exec(command.loginName);
      const suffix = matched?.[1] ?? "";

      // this just returns orgs where the suffix is set as primary domain
      const orgs = await getOrgsByDomain(suffix);
      const orgToCheckForDiscovery =
        orgs.result && orgs.result.length === 1 ? orgs.result[0].id : undefined;

      const orgLoginSettings = await getLoginSettings(orgToCheckForDiscovery);
      if (orgLoginSettings?.allowDomainDiscovery) {
        orgToRegisterOn = orgToCheckForDiscovery;
      }
    }

    // do not register user if ignoreUnknownUsernames is set
    if (orgToRegisterOn && !loginSettings?.ignoreUnknownUsernames) {
      const params = new URLSearchParams({ organization: orgToRegisterOn });

      if (command.authRequestId) {
        params.set("authRequestId", command.authRequestId);
      }
      if (command.loginName) {
        params.set("loginName", command.loginName);
      }

      return redirect("/register?" + params);
    }
  }

  if (loginSettings?.ignoreUnknownUsernames) {
    const paramsPasswordDefault = new URLSearchParams({
      loginName: command.loginName,
    });

    if (command.authRequestId) {
      paramsPasswordDefault.append("authRequestId", command.authRequestId);
    }

    if (command.organization) {
      paramsPasswordDefault.append("organization", command.organization);
    }

    return redirect("/password?" + paramsPasswordDefault);
  }

  // fallbackToPassword

  return { error: "Could not find user" };
}
