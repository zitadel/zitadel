"use server";

import { getSessionCookieById } from "@/lib/server/cookies";
import { addIdpFlowProof } from "@/lib/server/fingerprint";
import { getPublicHost } from "@/lib/server/host";
import { getServiceConfig } from "@/lib/server/service-url";
import { ServiceConfig, startIdentityProviderFlow, startLDAPIdentityProviderFlow } from "@/lib/server/zitadel";
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

  const response = await startIDPFlow({
    serviceConfig,
    host,
    idpId,
    successUrl: `/idp/${provider}/process?` + params.toString(),
    failureUrl: `/idp/${provider}/failure?` + params.toString(),
    loginHint: loginHint || undefined,
  });

  if (!response) {
    return { error: "Could not start IDP flow" };
  }

  if (response && "error" in response && response?.error) {
    return { error: response.error };
  }

  if (response && "redirect" in response && response?.redirect) {
    redirect(response.redirect);
  }

  if (response && "samlData" in response && response?.samlData) {
    return { samlData: response.samlData };
  }

  return { error: "Unexpected response from IDP flow" };
}

export type StartIDPFlowCommand = {
  serviceConfig: ServiceConfig;
  host: string;
  idpId: string;
  successUrl: string;
  failureUrl: string;
  loginHint?: string;
};

async function startIDPFlow(command: StartIDPFlowCommand) {
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
    return { error: "Could not start IDP flow" };
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

  if (!command.username || !command.password) {
    return { error: "No username or password provided" };
  }

  // the session to link must be one of this browser's sessions, same as redirectToIdp
  if (command.linkToSessionId) {
    const sessionCookie = await getSessionCookieById({ sessionId: command.linkToSessionId }).catch(() => undefined);
    if (!sessionCookie) {
      return { error: "Invalid session for linking" };
    }
  }

  const response = await startLDAPIdentityProviderFlow({
    serviceConfig,
    idpId: command.idpId,
    username: command.username,
    password: command.password,
  });

  if (!response || response.nextStep.case !== "idpIntent" || !response.nextStep.value) {
    return { error: "Could not start LDAP identity provider flow" };
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
