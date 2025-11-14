"use server";

import { getServiceUrlFromHeaders } from "@/lib/service-url";
import {
  retrieveIDPIntent,
  getIDPByID,
  updateHuman,
  addIDPLink,
  listUsers,
  addHuman,
  getLoginSettings,
  getOrgsByDomain,
  getActiveIdentityProviders,
  getUserByID,
  getDefaultOrg,
} from "@/lib/zitadel";
import { headers } from "next/headers";
import { create } from "@zitadel/client";
import { AutoLinkingOption } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { OrganizationSchema } from "@zitadel/proto/zitadel/object/v2/object_pb";
import {
  AddHumanUserRequest,
  AddHumanUserRequestSchema,
  UpdateHumanUserRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { createNewSessionFromIdpIntent } from "./idp";
import { getTranslations } from "next-intl/server";

const ORG_SUFFIX_REGEX = /(?<=@)(.+)/;

async function resolveOrganizationForUser({
  organization,
  addHumanUser,
  serviceUrl,
}: {
  organization?: string;
  addHumanUser?: AddHumanUserRequest;
  serviceUrl: string;
}): Promise<string | undefined> {
  if (organization) return organization;

  if (addHumanUser?.username && ORG_SUFFIX_REGEX.test(addHumanUser.username)) {
    const matched = ORG_SUFFIX_REGEX.exec(addHumanUser.username);
    const suffix = matched?.[1] ?? "";

    const orgs = await getOrgsByDomain({
      serviceUrl,
      domain: suffix,
    });
    const orgToCheckForDiscovery = orgs.result && orgs.result.length === 1 ? orgs.result[0].id : undefined;

    if (orgToCheckForDiscovery) {
      const orgLoginSettings = await getLoginSettings({
        serviceUrl,
        organization: orgToCheckForDiscovery,
      });
      if (orgLoginSettings?.allowDomainDiscovery) {
        return orgToCheckForDiscovery;
      }
    }
  }

  // Fallback to default organization if no org was resolved through discovery
  const defaultOrg = await getDefaultOrg({ serviceUrl });
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
  serviceUrl,
  userOrganizationId,
  idpId,
}: {
  serviceUrl: string;
  userOrganizationId: string;
  idpId: string;
}): Promise<boolean> {
  // Check organization login settings
  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: userOrganizationId,
  });

  if (!loginSettings?.allowExternalIdp) {
    return false;
  }

  // Check if the IDP is activated for the organization and allows linking
  const activeIDPs = await getActiveIdentityProviders({
    serviceUrl,
    orgId: userOrganizationId,
    linking_allowed: true,
  });

  const isIDPActive = activeIDPs.identityProviders?.some((idp) => idp.id === idpId);

  if (!isIDPActive) {
    return false;
  }

  return true;
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
  link,
  postErrorRedirectUrl,
}: {
  provider: string;
  id: string;
  token: string;
  requestId?: string;
  organization?: string;
  link?: string;
  postErrorRedirectUrl?: string;
}): Promise<{ redirect?: string; error?: string }> {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

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
    const intent = await retrieveIDPIntent({
      serviceUrl,
      id,
      token,
    });

    console.log("[IDP Process] Intent retrieved successfully, processing business logic");

    const { idpInformation, userId, addHumanUser, updateHumanUser } = intent;

    if (!idpInformation) {
      console.error("[IDP Process] IDP information missing");
      return { redirect: `/idp/${provider}/failure?error=missing_idp_info` };
    }

    // Get IDP configuration
    const idp = await getIDPByID({
      serviceUrl,
      id: idpInformation.idpId,
    });

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

      if (additionalParams) {
        Object.entries(additionalParams).forEach(([key, value]) => {
          if (value) params.set(key, value);
        });
      }

      return params.toString();
    };

    // ============================================
    // CASE 1: User exists and should sign in
    // ============================================
    if (userId && !link) {
      // Auto-update user if enabled
      if (options?.isAutoUpdate && updateHumanUser) {
        try {
          await updateHuman({
            serviceUrl,
            request: create(UpdateHumanUserRequestSchema, {
              userId: userId,
              profile: updateHumanUser.profile,
              email: updateHumanUser.email,
              phone: updateHumanUser.phone,
            }),
          });
          console.log("[IDP Process] User auto-updated successfully");
        } catch (error: unknown) {
          console.warn("[IDP Process] Error auto-updating user:", error);
        }
      }

      // Create session and handle redirect
      console.log("[IDP Process] Creating session for existing user");
      const sessionResult = await createNewSessionFromIdpIntent({
        userId,
        idpIntent: {
          idpIntentId: id,
          idpIntentToken: token,
        },
        requestId,
        organization,
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

    // ============================================
    // CASE 2: Link IDP to existing user
    // ============================================
    if (link && userId) {
      if (!options?.isLinkingAllowed) {
        console.error("[IDP Process] Linking not allowed by IDP configuration");
        const params = buildRedirectParams();
        return { redirect: `/idp/${provider}/linking-failed?${params}&error=linking_not_allowed` };
      }

      try {
        // Get user to retrieve their organization
        const targetUser = await getUserByID({ serviceUrl, userId });

        if (!targetUser || !targetUser.details?.resourceOwner) {
          console.error("[IDP Process] User not found or missing organization");
          const params = buildRedirectParams();
          return { redirect: `/idp/${provider}/linking-failed?${params}&error=user_not_found` };
        }

        // Validate IDP linking permissions
        const isAllowed = await validateIDPLinkingPermissions({
          serviceUrl,
          userOrganizationId: targetUser.details.resourceOwner,
          idpId: idpInformation.idpId,
        });

        if (!isAllowed) {
          console.error("[IDP Process] IDP linking validation failed");
          const params = buildRedirectParams();
          return { redirect: `/idp/${provider}/linking-failed?${params}&error=validation_failed` };
        }

        await addIDPLink({
          serviceUrl,
          idp: {
            id: idpInformation.idpId,
            userId: idpInformation.userId,
            userName: idpInformation.userName,
          },
          userId,
        });
        console.log("[IDP Process] IDP linked successfully, creating session");

        // Create session after linking
        const sessionResult = await createNewSessionFromIdpIntent({
          userId,
          idpIntent: {
            idpIntentId: id,
            idpIntentToken: token,
          },
          requestId,
          organization,
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
        const params = buildRedirectParams();
        return { redirect: `/idp/${provider}/linking-failed?${params}` };
      }
    }

    // ============================================
    // CASE 3: Auto-linking (search for user and link)
    // ============================================
    if (options?.autoLinking) {
      let foundUser;
      const email = addHumanUser?.email?.email;

      if (options.autoLinking === AutoLinkingOption.EMAIL && email) {
        foundUser = await listUsers({ serviceUrl, email, organizationId: organization }).then((response) => {
          return response.result ? response.result[0] : null;
        });
      } else if (options.autoLinking === AutoLinkingOption.USERNAME) {
        foundUser = await listUsers({
          serviceUrl,
          userName: idpInformation.userName,
          organizationId: organization,
        }).then((response) => {
          return response.result ? response.result[0] : null;
        });
      } else {
        foundUser = await listUsers({
          serviceUrl,
          userName: idpInformation.userName,
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
            serviceUrl,
            userOrganizationId: foundUser.details.resourceOwner,
            idpId: idpInformation.idpId,
          });

          if (!isAllowed) {
            console.error("[IDP Process] Auto-linking validation failed");
            const params = buildRedirectParams();
            return { redirect: `/idp/${provider}/linking-failed?${params}&error=validation_failed` };
          }

          await addIDPLink({
            serviceUrl,
            idp: {
              id: idpInformation.idpId,
              userId: idpInformation.userId,
              userName: idpInformation.userName,
            },
            userId: foundUser.userId,
          });
          console.log("[IDP Process] User auto-linked successfully, creating session");

          // Create session after auto-linking
          const sessionResult = await createNewSessionFromIdpIntent({
            userId: foundUser.userId,
            idpIntent: {
              idpIntentId: id,
              idpIntentToken: token,
            },
            requestId,
            organization,
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
          const params = buildRedirectParams();
          return { redirect: `/idp/${provider}/linking-failed?${params}` };
        }
      }
    }

    // ============================================
    // CASE 4: Auto-creation of user
    // ============================================
    if (options?.isAutoCreation && addHumanUser) {
      const orgToRegisterOn = await resolveOrganizationForUser({
        organization,
        addHumanUser,
        serviceUrl,
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
        const newUser = await addHuman({
          serviceUrl,
          request: addHumanUserWithOrganization,
        });
        console.log("[IDP Process] User auto-created successfully, creating session");

        // Create session for newly created user
        const sessionResult = await createNewSessionFromIdpIntent({
          userId: newUser.userId,
          idpIntent: {
            idpIntentId: id,
            idpIntentToken: token,
          },
          requestId,
          organization,
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

    // ============================================
    // CASE 5: Manual user creation allowed
    // ============================================
    if (options?.isCreationAllowed && addHumanUser) {
      const orgToRegisterOn = await resolveOrganizationForUser({
        organization,
        addHumanUser,
        serviceUrl,
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
          idpId: idpInformation.idpId,
          idpUserId: idpInformation.userId || "",
          idpUserName: idpInformation.userName || "",
          // User data for pre-filling form
          givenName: addHumanUser.profile?.givenName || "",
          familyName: addHumanUser.profile?.familyName || "",
          email: addHumanUser.email?.email || "",
        },
        true,
      ); // includeToken=true
      return { redirect: `/idp/${provider}/complete-registration?${params}` };
    }

    // ============================================
    // CASE 6: No user found and creation not allowed
    // ============================================
    console.log("[IDP Process] No matching user and creation not allowed");
    const params = buildRedirectParams();
    return { redirect: `/idp/${provider}/account-not-found?${params}` };
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
