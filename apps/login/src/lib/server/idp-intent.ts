"use server";

import { getServiceConfig } from "@/lib/service-url";
import {
  retrieveIDPIntent,
  getIDPByID,
  addIDPLink,
  listUsers,
  addHuman,
  getLoginSettings,
  getOrgsByDomain,
  getActiveIdentityProviders,
  getUserByID,
  getDefaultOrg,
  updateHuman,
  ServiceConfig,
} from "@/lib/zitadel";
import { headers } from "next/headers";
import { Code, ConnectError, create } from "@zitadel/client";
import { AutoLinkingOption } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { OrganizationSchema } from "@zitadel/proto/zitadel/object/v2/object_pb";
import {
  AddHumanUserRequest,
  AddHumanUserRequestSchema,
  UpdateHumanUserRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getSession } from "@/lib/zitadel";
import { getSessionCookieById } from "@/lib/cookies";
import { createNewSessionFromIdpIntent } from "./idp";
import { getTranslations } from "next-intl/server";
import crypto from "crypto";
import { getFingerprintIdCookie } from "../fingerprint";

const ORG_SUFFIX_REGEX = /(?<=@)(.+)/;

async function resolveOrganizationForUser({
  organization,
  addHumanUser,
  serviceConfig,
}: {
  organization?: string;
  addHumanUser?: AddHumanUserRequest;
  serviceConfig: ServiceConfig;
}): Promise<string | undefined> {
  if (organization) return organization;

  if (addHumanUser?.username && ORG_SUFFIX_REGEX.test(addHumanUser.username)) {
    const matched = ORG_SUFFIX_REGEX.exec(addHumanUser.username);
    const suffix = matched?.[1] ?? "";

    const orgs = await getOrgsByDomain({ serviceConfig, domain: suffix });
    const orgToCheckForDiscovery = orgs.result && orgs.result.length === 1 ? orgs.result[0].id : undefined;

    if (orgToCheckForDiscovery) {
      const orgLoginSettings = await getLoginSettings({ serviceConfig, organization: orgToCheckForDiscovery });
      if (orgLoginSettings?.allowDomainDiscovery) {
        return orgToCheckForDiscovery;
      }
    }
  }

  // Fallback to default organization if no org was resolved through discovery
  const defaultOrg = await getDefaultOrg({ serviceConfig });
  return defaultOrg?.id;
}

/**
 * Validates if IDP linking is allowed for a user's organization.
 * Checks:
 * 1. Organization allows external IDP login (allowExternalIdp)
 * 2. The specific IDP is activated for the organization
 *
 */
export async function validateIDPLinkingPermissions({
  serviceConfig,
  userOrganizationId,
  idpId,
}: {
  serviceConfig: ServiceConfig;
  userOrganizationId: string;
  idpId: string;
}): Promise<boolean> {
  // Check organization login settings
  const loginSettings = await getLoginSettings({ serviceConfig, organization: userOrganizationId });

  if (!loginSettings?.allowExternalIdp) {
    return false;
  }

  // Check if the IDP is activated for the organization and allows linking
  const activeIDPs = await getActiveIdentityProviders({ serviceConfig, orgId: userOrganizationId, linking_allowed: true });

  const isIDPActive = activeIDPs.identityProviders?.some((idp) => idp.id === idpId);

  if (!isIDPActive) {
    return false;
  }

  return true;
}

type IDPIntentResult = Awaited<ReturnType<typeof retrieveIDPIntent>>;
type IDPConfig = Awaited<ReturnType<typeof getIDPByID>>;

interface IDPHandlerContext {
  serviceConfig: ServiceConfig;
  t: (key: string) => string;
  intent: IDPIntentResult;
  idp: NonNullable<IDPConfig>;
  options: NonNullable<NonNullable<IDPConfig>["config"]>["options"];
  params: {
    provider: string;
    id: string;
    token: string;
    requestId?: string;
    organization?: string;
    postErrorRedirectUrl?: string;
    sessionId?: string;
    linkFingerprint?: string;
  };
  buildRedirectParams: (additionalParams?: Record<string, string>, includeToken?: boolean) => string;
}

type IDPHandlerResult = { redirect?: string; error?: string } | null;

/**
 * CASE 1: Explicit Linking (via sessionId)
 * This happens when a logged-in user initiates an IDP flow to link it to their account.
 */
async function resolveUserIdFromSession({
  sessionId,
  serviceConfig,
  provider,
}: {
  sessionId: string;
  serviceConfig: ServiceConfig;
  provider: string;
}) {
  try {
    const sessionCookie = await getSessionCookieById({ sessionId });
    if (!sessionCookie) {
      console.warn("[IDP Process] Session for linking not found or invalid");
      return { redirect: `/idp/${provider}/linking-failed?error=session_invalid` };
    }

    const sessionResp = await getSession({
      serviceConfig,
      sessionId: sessionCookie.id,
      sessionToken: sessionCookie.token,
    });
    const session = sessionResp.session;

    if (!session?.factors?.user?.id) {
      console.warn("[IDP Process] Session found but no userId associated for linking.");
      return { redirect: `/idp/${provider}/linking-failed?error=session_invalid` };
    }

    return { userId: session.factors.user.id };
  } catch (error) {
    console.warn("[IDP Process] Error retrieving session for linking:", error);
    return { redirect: `/idp/${provider}/linking-failed?error=session_invalid` };
  }
}

async function handleExplicitLinking(ctx: IDPHandlerContext): Promise<IDPHandlerResult> {
  const { sessionId, linkFingerprint, provider } = ctx.params;
  const { userId } = ctx.intent;
  const { options, serviceConfig, intent, t, buildRedirectParams } = ctx;

  if (sessionId && !userId) {
    // Intent should not have a userId if it is a linking intent
    // 1. Security Check: Verify Fingerprint
    const fingerprintCookie = await getFingerprintIdCookie();

    if (!linkFingerprint || !fingerprintCookie?.value) {
      console.warn("[IDP Process] Missing fingerprint information for linking verification");
      return { redirect: `/idp/${provider}/linking-failed?error=session_mismatch` };
    }

    const expectedHash = crypto
      .createHash("sha256")
      .update(sessionId + fingerprintCookie.value)
      .digest("hex");

    if (linkFingerprint !== expectedHash) {
      console.warn("[IDP Process] Session linking fingerprint mismatch");
      return { redirect: `/idp/${provider}/linking-failed?error=session_mismatch` };
    }

    // 2. Retrieve Session & Resolve User
    const { userId: resolvedUserId, redirect: sessionRedirect } = await resolveUserIdFromSession({
      sessionId,
      serviceConfig,
      provider,
    });

    if (sessionRedirect || !resolvedUserId) {
      return { redirect: sessionRedirect || `/idp/${provider}/linking-failed?error=session_invalid` };
    }

    console.log("[IDP Process] Resolved userId from session link:", resolvedUserId);

    // 3. Perform Linking Logic
    if (!options?.isLinkingAllowed) {
      console.error("[IDP Process] Linking not allowed by IDP configuration");
      const params = buildRedirectParams();
      return { redirect: `/idp/${provider}/linking-failed?${params}&error=linking_not_allowed` };
    }

    try {
      const targetUser = await getUserByID({ serviceConfig, userId: resolvedUserId });

      if (!targetUser || !targetUser.details?.resourceOwner) {
        console.error("[IDP Process] User not found or missing organization");
        const params = buildRedirectParams();
        return { redirect: `/idp/${provider}/linking-failed?${params}&error=user_not_found` };
      }

      const isAllowed = await validateIDPLinkingPermissions({
        serviceConfig,
        userOrganizationId: targetUser.details.resourceOwner,
        idpId: intent.idpInformation!.idpId,
      });

      if (!isAllowed) {
        console.error("[IDP Process] IDP linking validation failed");
        const params = buildRedirectParams();
        return { redirect: `/idp/${provider}/linking-failed?${params}&error=validation_failed` };
      }

      await addIDPLink({
        serviceConfig,
        idp: {
          id: intent.idpInformation!.idpId,
          userId: intent.idpInformation!.userId,
          userName: intent.idpInformation!.userName,
        },
        userId: resolvedUserId,
      });
      console.log("[IDP Process] IDP linked successfully, creating session");

      const sessionResult = await createNewSessionFromIdpIntent({
        userId: resolvedUserId,
        idpIntent: {
          idpIntentId: ctx.params.id,
          idpIntentToken: ctx.params.token,
        },
        requestId: ctx.params.requestId,
        organization: ctx.params.organization,
      });

      if ("error" in sessionResult && sessionResult.error) {
        console.error("[IDP Process] Error creating session:", sessionResult.error);
        return { error: sessionResult.error };
      }

      if ("redirect" in sessionResult && sessionResult.redirect) {
        console.log("[IDP Process] Session created, redirecting to:", sessionResult.redirect);
        return { redirect: sessionResult.redirect };
      }

      return { error: t("errors.sessionCreationFailed") };
    } catch (error) {
      console.error("[IDP Process] Error linking IDP:", error);
      const errorMessage = error instanceof Error ? error.message : t("errors.unknownError");
      let params = buildRedirectParams({ error: errorMessage });
      if (error instanceof ConnectError && error.code === Code.AlreadyExists) {
        params = buildRedirectParams({ error: "external_idp_taken" });
      }
      return { redirect: `/idp/${provider}/linking-failed?${params}` };
    }
  }

  return null;
}

/**
 * CASE 2: User exists and should sign in
 */
async function handleUserExists(ctx: IDPHandlerContext): Promise<IDPHandlerResult> {
  const { sessionId } = ctx.params;
  const { userId, updateHumanUser } = ctx.intent;
  const { options, serviceConfig, t } = ctx;

  if (userId && !sessionId) {
    // Auto-update user if enabled
    if (options?.isAutoUpdate && updateHumanUser) {
      try {
        console.log("[IDP Process] Auto-updating user profile");
        await updateHuman({
          serviceConfig,
          request: create(UpdateHumanUserRequestSchema, {
            userId,
            profile: updateHumanUser.profile,
            email: updateHumanUser.email,
            phone: updateHumanUser.phone,
          }),
        });
      } catch (error) {
        console.warn("[IDP Process] Failed to auto-update user:", error);
        // Continue with login even if update fails
      }
    }

    // Create session and handle redirect
    console.log("[IDP Process] Creating session for existing user");
    const sessionResult = await createNewSessionFromIdpIntent({
      userId,
      idpIntent: {
        idpIntentId: ctx.params.id,
        idpIntentToken: ctx.params.token,
      },
      requestId: ctx.params.requestId,
      organization: ctx.params.organization,
    });

    if ("error" in sessionResult && sessionResult.error) {
      console.error("[IDP Process] Error creating session:", sessionResult.error);
      return { error: sessionResult.error };
    }

    if ("redirect" in sessionResult && sessionResult.redirect) {
      console.log("[IDP Process] Session created, redirecting to:", sessionResult.redirect);
      return { redirect: sessionResult.redirect };
    }

    return { error: t("errors.sessionCreationFailed") };
  }

  return null;
}

/**
 * CASE 3: Auto-linking (search for user and link)
 */
async function handleAutoLinking(ctx: IDPHandlerContext): Promise<IDPHandlerResult> {
  const { options, intent, serviceConfig, buildRedirectParams, t } = ctx;
  const { addHumanUser, idpInformation } = intent;
  const { organization, provider } = ctx.params;

  if (options?.autoLinking) {
    let foundUser;
    const email = addHumanUser?.email?.email;

    if (options.autoLinking === AutoLinkingOption.EMAIL && email) {
      foundUser = await listUsers({ serviceConfig, email, organizationId: organization }).then((response) => {
        return response.result ? response.result[0] : null;
      });
    } else if (options.autoLinking === AutoLinkingOption.USERNAME) {
      foundUser = await listUsers({
        serviceConfig,
        userName: idpInformation!.userName,
        organizationId: organization,
      }).then((response) => {
        return response.result ? response.result[0] : null;
      });
    } else {
      foundUser = await listUsers({
        serviceConfig,
        userName: idpInformation!.userName,
        email,
        organizationId: organization,
      }).then((response) => {
        return response.result ? response.result[0] : null;
      });
    }

    if (foundUser) {
      try {
        if (!foundUser.details?.resourceOwner) {
          console.error("[IDP Process] Found user missing organization information");
          const params = buildRedirectParams();
          return { redirect: `/idp/${provider}/linking-failed?${params}&error=missing_organization` };
        }

        // Validate IDP linking permissions
        const isAllowed = await validateIDPLinkingPermissions({
          serviceConfig,
          userOrganizationId: foundUser.details.resourceOwner,
          idpId: idpInformation!.idpId,
        });

        if (!isAllowed) {
          console.error("[IDP Process] Auto-linking validation failed");
          const params = buildRedirectParams();
          return { redirect: `/idp/${provider}/linking-failed?${params}&error=validation_failed` };
        }

        await addIDPLink({
          serviceConfig,
          idp: {
            id: idpInformation!.idpId,
            userId: idpInformation!.userId,
            userName: idpInformation!.userName,
          },
          userId: foundUser.userId,
        });
        console.log("[IDP Process] User auto-linked successfully, creating session");

        // Create session after auto-linking
        const sessionResult = await createNewSessionFromIdpIntent({
          userId: foundUser.userId,
          idpIntent: {
            idpIntentId: ctx.params.id,
            idpIntentToken: ctx.params.token,
          },
          requestId: ctx.params.requestId,
          organization: ctx.params.organization,
        });

        if ("error" in sessionResult && sessionResult.error) {
          console.error("[IDP Process] Error creating session:", sessionResult.error);
          return { error: sessionResult.error };
        }

        if ("redirect" in sessionResult && sessionResult.redirect) {
          console.log("[IDP Process] Session created, redirecting to:", sessionResult.redirect);
          return { redirect: sessionResult.redirect };
        }

        return { error: t("errors.sessionCreationFailed") };
      } catch (error) {
        console.error("[IDP Process] Error auto-linking user:", error);
        const errorMessage = error instanceof Error ? error.message : t("errors.unknownError");
        const params = buildRedirectParams({ error: errorMessage });
        return { redirect: `/idp/${provider}/linking-failed?${params}` };
      }
    }
  }

  return null;
}

/**
 * CASE 4: Auto-creation of user
 */
async function handleAutoCreation(ctx: IDPHandlerContext): Promise<IDPHandlerResult> {
  const { options, intent, serviceConfig, buildRedirectParams, t } = ctx;
  const { addHumanUser } = intent;
  const { organization, provider } = ctx.params;

  if (options?.isAutoCreation && addHumanUser) {
    const orgToRegisterOn = await resolveOrganizationForUser({
      organization,
      addHumanUser,
      serviceConfig,
    });

    if (!orgToRegisterOn) {
      console.error("[IDP Process] Could not determine organization for auto-creation (no default org available)");
      const params = buildRedirectParams();
      return { redirect: `/idp/${provider}/failure?${params}&error=no_organization_context` };
    }

    const organizationSchema = create(OrganizationSchema, {
      org: { case: "orgId", value: orgToRegisterOn },
    });

    const addHumanUserWithOrganization = create(AddHumanUserRequestSchema, {
      ...addHumanUser,
      organization: organizationSchema,
    });

    try {
      const newUser = await addHuman({ serviceConfig, request: addHumanUserWithOrganization });
      console.log("[IDP Process] User auto-created successfully, creating session");

      // Create session for newly created user
      const sessionResult = await createNewSessionFromIdpIntent({
        userId: newUser.userId,
        idpIntent: {
          idpIntentId: ctx.params.id,
          idpIntentToken: ctx.params.token,
        },
        requestId: ctx.params.requestId,
        organization: ctx.params.organization,
      });

      if ("error" in sessionResult && sessionResult.error) {
        console.error("[IDP Process] Error creating session:", sessionResult.error);
        return { error: sessionResult.error };
      }

      if ("redirect" in sessionResult && sessionResult.redirect) {
        console.log("[IDP Process] Session created, redirecting to:", sessionResult.redirect);
        return { redirect: sessionResult.redirect };
      }

      return { error: t("errors.sessionCreationFailed") };
    } catch (error: unknown) {
      console.error("[IDP Process] Error auto-creating user:", error);
      const params = buildRedirectParams();
      return { redirect: `/idp/${provider}/failure?${params}&error=user_creation_failed` };
    }
  }

  return null;
}

/**
 * CASE 5: Manual user creation allowed
 */
async function handleManualCreation(ctx: IDPHandlerContext): Promise<IDPHandlerResult> {
  const { options, intent, serviceConfig, buildRedirectParams } = ctx;
  const { addHumanUser, idpInformation } = intent;
  const { organization, provider } = ctx.params;

  if (options?.isCreationAllowed && addHumanUser) {
    const orgToRegisterOn = await resolveOrganizationForUser({
      organization,
      addHumanUser,
      serviceConfig,
    });

    if (!orgToRegisterOn) {
      console.error("[IDP Process] Could not determine organization for registration (no default org available)");
      const params = buildRedirectParams();
      return { redirect: `/idp/${provider}/registration-failed?${params}` };
    }

    // Store user data for manual registration form
    // Note: includeToken=true because the session hasn't been created yet
    // The token will be needed when registerUserAndLinkToIDP creates the session
    const params = buildRedirectParams(
      {
        organization: orgToRegisterOn,
        idpId: idpInformation!.idpId,
        idpUserId: idpInformation!.userId || "",
        idpUserName: idpInformation!.userName || "",
        // User data for pre-filling form
        givenName: addHumanUser.profile?.givenName || "",
        familyName: addHumanUser.profile?.familyName || "",
        email: addHumanUser.email?.email || "",
      },
      true,
    ); // includeToken=true
    return { redirect: `/idp/${provider}/complete-registration?${params}` };
  }

  return null;
}

/**
 * CASE 6: No user found and creation not allowed
 */
async function handleNoUserFound(ctx: IDPHandlerContext): Promise<IDPHandlerResult> {
  const { buildRedirectParams } = ctx;
  console.log("[IDP Process] No matching user and creation not allowed");
  const params = buildRedirectParams();
  return { redirect: `/idp/${ctx.params.provider}/account-not-found?${params}` };
}

/**
 * Server action to process IDP callback and handle ALL business logic.
 * This action:
 * 1. Consumes the single-use token once
 * 2. Performs all IDP-related operations (auto-update, auto-linking, auto-creation)
 * 3. Returns redirect URL or error for client-side navigation
 */
export async function processIDPCallback({
  provider,
  id,
  token,
  requestId,
  organization,
  postErrorRedirectUrl,
  sessionId,
  linkFingerprint,
}: {
  provider: string;
  id: string;
  token: string;
  requestId?: string;
  organization?: string;
  postErrorRedirectUrl?: string;
  sessionId?: string;
  linkFingerprint?: string;
}): Promise<{ redirect?: string; error?: string }> {
  // ... (headers and config retrieval) ...
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const t = await getTranslations("idp");

  // Validate required parameters
  if (!provider || !id || !token) {
    console.error("[IDP Process] Missing required parameters:", { provider, id, hasToken: !!token });
    const errorParams = new URLSearchParams();
    if (requestId) errorParams.set("requestId", requestId);
    if (organization) errorParams.set("organization", organization);
    if (postErrorRedirectUrl) errorParams.set("postErrorRedirectUrl", postErrorRedirectUrl);

    return { redirect: `/idp/${provider}/failure?${errorParams.toString()}` };
  }

  try {
    console.log("[IDP Process] Retrieving IDP intent (single call):", {
      id,
      tokenPreview: token.substring(0, 10) + "...",
      timestamp: new Date().toISOString(),
    });

    // Consume the single-use token ONCE
    const intent = await retrieveIDPIntent({ serviceConfig, id, token });

    console.log("[IDP Process] Intent retrieved successfully, processing business logic");

    const { idpInformation } = intent;

    // Verify we have IDP info early on
    if (!idpInformation) {
      console.error("[IDP Process] IDP information missing");
      return { redirect: `/idp/${provider}/failure?error=missing_idp_info` };
    }

    // Get IDP configuration
    const idp = await getIDPByID({ serviceConfig, id: idpInformation.idpId });

    if (!idp) {
      return { error: t("errors.idpNotFound") };
    }

    const options = idp?.config?.options;

    // Build base redirect params
    const buildRedirectParams = (additionalParams?: Record<string, string>, includeToken: boolean = false) => {
      const params = new URLSearchParams();
      params.set("id", id);
      if (includeToken) params.set("token", token);
      if (requestId) params.set("requestId", requestId);
      if (organization) params.set("organization", organization);
      if (postErrorRedirectUrl) params.set("postErrorRedirectUrl", postErrorRedirectUrl);
      if (sessionId) params.set("linkToSessionId", sessionId); // Include sessionId in redirects

      if (additionalParams) {
        Object.entries(additionalParams).forEach(([key, value]) => {
          if (value) params.set(key, value);
        });
      }

      return params.toString();
    };

    const ctx: IDPHandlerContext = {
      serviceConfig,
      t,
      intent,
      idp,
      options,
      params: {
        provider,
        id,
        token,
        requestId,
        organization,
        postErrorRedirectUrl,
        sessionId,
        linkFingerprint,
      },
      buildRedirectParams,
    };

    const handlers = [
      handleExplicitLinking,
      handleUserExists,
      handleAutoLinking,
      handleAutoCreation,
      handleManualCreation,
      handleNoUserFound,
    ];

    for (const handler of handlers) {
      const result = await handler(ctx);
      if (result) {
        return result;
      }
    }

    // Should theoretically be unreachable if handleNoUserFound covers the rest
    return { error: t("errors.unknown") };
  } catch (error: unknown) {
    console.error("[IDP Process] Error processing intent:", error);

    const errorParams = new URLSearchParams();
    if (requestId) errorParams.set("requestId", requestId);
    if (organization) errorParams.set("organization", organization);
    if (postErrorRedirectUrl) errorParams.set("postErrorRedirectUrl", postErrorRedirectUrl);
    errorParams.set("error", error instanceof Error ? error.message : t("errors.unknownError"));

    return { redirect: `/idp/${provider}/failure?${errorParams.toString()}` };
  }
}
