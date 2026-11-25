"use server";

import { answerRejection } from "@/lib/server/api-error";
import { getSessionCookieById } from "@/lib/server/cookies";
import { addIdpFlowProof } from "@/lib/server/fingerprint";
import { getPublicHost } from "@/lib/server/host";
import { getServiceConfig } from "@/lib/server/service-url";
import { ServiceConfig, startIdentityProviderFlow, startLDAPIdentityProviderFlow } from "@/lib/server/zitadel";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { redirect } from "next/navigation";

export type RedirectToIdpState =
  { error?: string | null; samlData?: { url: string; fields: Record<string, string> } } | undefined;

export async function redirectToIdp(prevState: RedirectToIdpState, formData: FormData): Promise<RedirectToIdpState> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const host = getPublicHost(_headers);

  const params = new URLSearchParams();

  const sessionId = formData.get("sessionId") as string;
  const requestId = formData.get("requestId") as string;
  const organization = formData.get("organization") as string;
  const idpId = formData.get("id") as string;
  const provider = formData.get("provider") as string;
  const postErrorRedirectUrl = formData.get("postErrorRedirectUrl") as string;
  const loginHint = formData.get("loginHint") as string;

  if (sessionId) {
    // the session to link must be one of this browser's sessions
    if (!(await getSessionCookieById({ sessionId }).catch(() => undefined))) {
      return { error: "Invalid session for linking" };
    }
    params.set("linkToSessionId", sessionId);
  }
  if (requestId) params.set("requestId", requestId);
  if (organization) params.set("organization", organization);
  if (postErrorRedirectUrl) params.set("postErrorRedirectUrl", postErrorRedirectUrl);

  // redirect to LDAP page where username and password is requested
  if (provider === "ldap") {
    params.set("idpId", idpId);
    redirect(`/idp/ldap?` + params.toString());
  }

  // binds the flow to this browser (and the session to link), checked when the callback is processed
  await addIdpFlowProof(params, sessionId);

  const t = await getTranslations("idp");

  let response;
  try {
    response = await startIDPFlow({
      serviceConfig,
      host,
      idpId,
      successUrl: `/idp/${provider}/process?` + params.toString(),
      failureUrl: `/idp/${provider}/failure?` + params.toString(),
      loginHint: loginHint || undefined,
    });
  } catch (error) {
    // e.g. the IdP was removed or is misconfigured
    return answerRejection(error, t("errors.couldNotStartIDPFlow"), { context: { idpId } });
  }

  if (!response) {
    return { error: t("errors.couldNotStartIDPFlow") };
  }

  if ("redirect" in response) {
    redirect(response.redirect);
  }

  return { samlData: response.samlData };
}

export type StartIDPFlowCommand = {
  serviceConfig: ServiceConfig;
  host: string;
  idpId: string;
  successUrl: string;
  failureUrl: string;
  loginHint?: string;
};

async function startIDPFlow(
  command: StartIDPFlowCommand,
): Promise<{ redirect: string } | { samlData: { url: string; fields: Record<string, string> } } | null> {
  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  const response = await startIdentityProviderFlow({
    serviceConfig: command.serviceConfig,
    idpId: command.idpId,
    urls: {
      successUrl: `${command.host.includes("localhost") ? "http://" : "https://"}${command.host}${basePath}${command.successUrl}`,
      failureUrl: `${command.host.includes("localhost") ? "http://" : "https://"}${command.host}${basePath}${command.failureUrl}`,
      loginHint: command.loginHint,
    },
  });

  if (!response || !response.url) {
    return null;
  }

  if (response.fields) {
    return { samlData: { url: response.url, fields: response.fields } };
  }

  return { redirect: response.url };
}

type createNewSessionForLDAPCommand = {
  username: string;
  password: string;
  idpId: string;
  link: boolean;
  requestId?: string;
  organization?: string;
  postErrorRedirectUrl?: string;
  linkToSessionId?: string;
};

export async function createNewSessionForLDAP(command: createNewSessionForLDAPCommand) {
  const _headers = await headers();

  const { serviceConfig } = getServiceConfig(_headers);

  const t = await getTranslations("ldap");

  if (!command.username || !command.password) {
    return { error: t("errors.couldNotVerify") };
  }

  // the session to link must be one of this browser's sessions, same as redirectToIdp
  if (command.linkToSessionId) {
    const sessionCookie = await getSessionCookieById({ sessionId: command.linkToSessionId }).catch(() => undefined);
    if (!sessionCookie) {
      return { error: "Invalid session for linking" };
    }
  }

  let response;
  try {
    response = await startLDAPIdentityProviderFlow({
      serviceConfig,
      idpId: command.idpId,
      username: command.username,
      password: command.password,
    });
  } catch (error) {
    // wrong credentials are answered with InvalidArgument (COMMAND-nzun2i)
    return answerRejection(error, t("errors.couldNotVerify"), { context: { idpId: command.idpId } });
  }

  if (!response || response.nextStep.case !== "idpIntent" || !response.nextStep.value) {
    return { error: t("errors.couldNotVerify") };
  }

  const { userId, idpIntentId, idpIntentToken } = response.nextStep.value;

  const params = new URLSearchParams({
    userId,
    id: idpIntentId,
    token: idpIntentToken,
  });
  await addIdpFlowProof(params, command.linkToSessionId);

  if (command.link) {
    params.set("link", "true");
  }
  if (command.requestId) params.set("requestId", command.requestId);
  if (command.organization) params.set("organization", command.organization);
  if (command.postErrorRedirectUrl) params.set("postErrorRedirectUrl", command.postErrorRedirectUrl);
  if (command.linkToSessionId) params.set("linkToSessionId", command.linkToSessionId);

  return {
    redirect: `/idp/ldap/process?` + params.toString(),
  };
}
