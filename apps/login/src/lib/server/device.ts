"use server";

import { authorizeOrDenyDeviceAuthorization, getDeviceAuthorizationRequest } from "@/lib/zitadel";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { getSessionCookieById } from "../cookies";
import { checkDeviceCodeCheckCookie, deleteDeviceCodeCheckCookie } from "../device-code-check";
import { getServiceConfig } from "../service-url";
import { isSessionValid, loadSession } from "../session";
import { sendLoginname } from "./loginname";

type DeviceAuthorizationCommand = {
  userCode: string;
  sessionId: string;
  organization?: string;
};

// Shared by approve and deny: the user code must have been entered in this browser,
// and the session must come from this browser's `sessions` cookie and still exist.
async function resolveConsentSession({ userCode, sessionId, organization }: DeviceAuthorizationCommand) {
  if (!(await checkDeviceCodeCheckCookie(userCode))) {
    return { redirect: "/device?" + new URLSearchParams({ user_code: userCode }) };
  }

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const t = await getTranslations("error");

  const cookie = await getSessionCookieById({ sessionId, organization });

  if (!cookie) {
    return { error: t("unknownContext") };
  }

  const session = await loadSession(serviceConfig, cookie);

  if (!session?.factors?.user) {
    return { error: t("sessionExpired") };
  }

  return { serviceConfig, t, cookie, session, user: session.factors.user };
}

/**
 * Approves a device authorization request for the user of the current session.
 *
 * The session token is never taken from the caller: it is read from the `sessions` cookie
 * server-side. Together with the Next.js server action origin check this makes the approval
 * an explicit, same-origin user action instead of a side effect of rendering a page.
 * The `deviceCodeCheck` cookie makes sure the user code was entered in this browser.
 */
export async function approveDeviceAuthorization(command: DeviceAuthorizationCommand) {
  const consent = await resolveConsentSession(command);

  if (!consent.session) {
    return consent;
  }

  const { userCode, organization } = command;
  const { serviceConfig, t, cookie, session, user } = consent;

  // Validate the session (including MFA) before handing out a token to the device.
  const valid = await isSessionValid({ serviceConfig, session });

  if (!valid) {
    // Send the user through the login funnel, the same way continueWithSession does.
    const res = await sendLoginname({
      loginName: user.loginName,
      organization: user.organizationId,
      requestId: `device_${userCode}`,
    });

    if (res && "redirect" in res && res.redirect) {
      return { redirect: res.redirect };
    }

    return { error: t("sessionExpired") };
  }

  const { deviceAuthorizationRequest } = await getDeviceAuthorizationRequest({ serviceConfig, userCode });

  if (!deviceAuthorizationRequest) {
    return { error: t("noDeviceRequest") };
  }

  await authorizeOrDenyDeviceAuthorization({
    serviceConfig,
    deviceAuthorizationId: deviceAuthorizationRequest.id,
    session: { sessionId: cookie.id, sessionToken: cookie.token },
  });

  await deleteDeviceCodeCheckCookie();

  const params = new URLSearchParams({ requestId: `device_${userCode}`, sessionId: cookie.id });

  if (organization) {
    params.append("organization", organization);
  }

  return { redirect: `/signedin?` + params };
}

export async function denyDeviceAuthorization(command: DeviceAuthorizationCommand) {
  const consent = await resolveConsentSession(command);

  if (!consent.session) {
    return consent;
  }

  const { userCode } = command;
  const { serviceConfig, t } = consent;

  const { deviceAuthorizationRequest } = await getDeviceAuthorizationRequest({ serviceConfig, userCode });

  if (!deviceAuthorizationRequest) {
    return { error: t("noDeviceRequest") };
  }

  // without a session, the device auth request is denied
  const response = await authorizeOrDenyDeviceAuthorization({
    serviceConfig,
    deviceAuthorizationId: deviceAuthorizationRequest.id,
  });

  await deleteDeviceCodeCheckCookie();

  return response;
}
