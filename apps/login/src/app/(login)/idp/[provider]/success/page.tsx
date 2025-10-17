import { DynamicTheme } from "@/components/dynamic-theme";
import { IdpSignin } from "@/components/idp-signin";
import { completeIDP } from "@/components/idps/pages/complete-idp";
import { linkingFailed } from "@/components/idps/pages/linking-failed";
import { linkingSuccess } from "@/components/idps/pages/linking-success";
import { loginFailed } from "@/components/idps/pages/login-failed";
import { loginSuccess } from "@/components/idps/pages/login-success";
import { Translated } from "@/components/translated";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import {
  addHuman,
  addIDPLink,
  getBrandingSettings,
  getDefaultOrg,
  getIDPByID,
  getLoginSettings,
  getOrgsByDomain,
  listUsers,
  retrieveIDPIntent,
  updateHuman,
} from "@/lib/zitadel";
import { ConnectError, create } from "@zitadel/client";
import { redirect } from "next/navigation";
import { AutoLinkingOption } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { OrganizationSchema } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import {
  AddHumanUserRequest,
  AddHumanUserRequestSchema,
  UpdateHumanUserRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";

const ORG_SUFFIX_REGEX = /(?<=@)(.+)/;

async function resolveOrganizationForUser({
  organization,
  addHumanUser,
  serviceUrl,
}: {
  organization?: string;
  addHumanUser?: { username?: string };
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
  return undefined;
}

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const params = await props.params;
  const searchParams = await props.searchParams;
  let { id, token, requestId, organization, link, postErrorRedirectUrl } = searchParams;
  const { provider } = params;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  let branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  if (!organization) {
    const org: Organization | null = await getDefaultOrg({
      serviceUrl,
    });
    if (org) {
      organization = org.id;
    }
  }

  if (!provider || !id || !token) {
    return loginFailed(branding, "IDP context missing");
  }

  const intent = await retrieveIDPIntent({
    serviceUrl,
    id,
    token,
  });

  const { idpInformation, userId } = intent;
  let { addHumanUser } = intent;

  if (!idpInformation) {
    return loginFailed(branding, "IDP information missing");
  }

  const idp = await getIDPByID({
    serviceUrl,
    id: idpInformation.idpId,
  });

  const options = idp?.config?.options;

  if (!idp) {
    throw new Error("IDP not found");
  }

  // sign in user. If user should be linked continue
  if (userId && !link) {
    // if auto update is enabled, we will update the user with the new information
    if (options?.isAutoUpdate && addHumanUser) {
      try {
        await updateHuman({
          serviceUrl,
          request: create(UpdateHumanUserRequestSchema, {
            userId: userId,
            profile: addHumanUser.profile,
            email: addHumanUser.email,
            phone: addHumanUser.phone,
          }),
        });
      } catch (error: unknown) {
        // Log the error and continue with the login process
        console.warn("An error occurred while updating the user:", error);
      }
    }

    return loginSuccess(userId, { idpIntentId: id, idpIntentToken: token }, requestId, branding);
  }

  if (link) {
    if (!options?.isLinkingAllowed) {
      // linking was probably disallowed since the invitation was created
      return linkingFailed(branding, "Linking is no longer allowed");
    }

    let idpLink;
    try {
      idpLink = await addIDPLink({
        serviceUrl,
        idp: {
          id: idpInformation.idpId,
          userId: idpInformation.userId,
          userName: idpInformation.userName,
        },
        userId,
      });
    } catch (error) {
      console.error(error);
      return linkingFailed(branding);
    }

    if (!idpLink) {
      return linkingFailed(branding);
    } else {
      return linkingSuccess(userId, { idpIntentId: id, idpIntentToken: token }, requestId, branding);
    }
  }

  // search for potential user via username, then link
  if (options?.autoLinking) {
    let foundUser;
    const email = addHumanUser?.email?.email;

    if (options.autoLinking === AutoLinkingOption.EMAIL && email) {
      foundUser = await listUsers({ serviceUrl, email }).then((response) => {
        return response.result ? response.result[0] : null;
      });
    } else if (options.autoLinking === AutoLinkingOption.USERNAME) {
      foundUser = await listUsers(
        options.autoLinking === AutoLinkingOption.USERNAME
          ? { serviceUrl, userName: idpInformation.userName }
          : { serviceUrl, email },
      ).then((response) => {
        return response.result ? response.result[0] : null;
      });
    } else {
      foundUser = await listUsers({
        serviceUrl,
        userName: idpInformation.userName,
        email,
      }).then((response) => {
        return response.result ? response.result[0] : null;
      });
    }

    if (foundUser) {
      let idpLink;
      try {
        idpLink = await addIDPLink({
          serviceUrl,
          idp: {
            id: idpInformation.idpId,
            userId: idpInformation.userId,
            userName: idpInformation.userName,
          },
          userId: foundUser.userId,
        });
      } catch (error) {
        console.error(error);
        return linkingFailed(branding);
      }

      if (!idpLink) {
        return linkingFailed(branding);
      } else {
        return linkingSuccess(foundUser.userId, { idpIntentId: id, idpIntentToken: token }, requestId, branding);
      }
    }
  }

  let newUser;
  // automatic creation of a user is allowed and data is complete
  if (options?.isAutoCreation && addHumanUser) {
    const orgToRegisterOn = await resolveOrganizationForUser({
      organization,
      addHumanUser,
      serviceUrl,
    });

    let addHumanUserWithOrganization: AddHumanUserRequest;
    if (orgToRegisterOn) {
      const organizationSchema = create(OrganizationSchema, {
        org: { case: "orgId", value: orgToRegisterOn },
      });

      addHumanUserWithOrganization = create(AddHumanUserRequestSchema, {
        ...addHumanUser,
        organization: organizationSchema,
      });
    } else {
      addHumanUserWithOrganization = create(AddHumanUserRequestSchema, addHumanUser);
    }

    try {
      newUser = await addHuman({
        serviceUrl,
        request: addHumanUserWithOrganization,
      });
    } catch (error: unknown) {
      console.error("An error occurred while creating the user:", error, addHumanUser);
      return loginFailed(
        branding,
        (error as ConnectError).message ? (error as ConnectError).message : "Could not create user",
      );
    }
  } else if (options?.isCreationAllowed) {
    // if no user was found, we will create a new user manually / redirect to the registration page
    const orgToRegisterOn = await resolveOrganizationForUser({
      organization,
      addHumanUser,
      serviceUrl,
    });

    if (orgToRegisterOn) {
      branding = await getBrandingSettings({
        serviceUrl,
        organization: orgToRegisterOn,
      });
    }

    if (!orgToRegisterOn) {
      // Redirect to registration-failed page - couldn't determine organization for registration
      const queryParams = new URLSearchParams();
      if (requestId) queryParams.set("requestId", requestId);
      if (organization) queryParams.set("organization", organization);
      if (postErrorRedirectUrl) queryParams.set("postErrorRedirectUrl", postErrorRedirectUrl);
      redirect(`/idp/${provider}/registration-failed?${queryParams.toString()}`);
    }

    return completeIDP({
      branding,
      idpIntent: { idpIntentId: id, idpIntentToken: token },
      addHumanUser,
      organization: orgToRegisterOn,
      requestId,
      idpUserId: idpInformation?.userId,
      idpId: idpInformation?.idpId,
      idpUserName: idpInformation?.userName,
    });
  }

  if (newUser) {
    return (
      <DynamicTheme branding={branding}>
        <div className="flex flex-col space-y-4">
          <h1>
            <Translated i18nKey="registerSuccess.title" namespace="idp" />
          </h1>
          <p className="ztdl-p">
            <Translated i18nKey="registerSuccess.description" namespace="idp" />
          </p>
        </div>

        <div className="w-full">
          <IdpSignin userId={newUser.userId} idpIntent={{ idpIntentId: id, idpIntentToken: token }} requestId={requestId} />
        </div>
      </DynamicTheme>
    );
  }

  // Redirect to account-not-found page with postErrorRedirectUrl
  // This provides a graceful fallback when no user was found and creation/linking is not allowed
  const queryParams = new URLSearchParams();
  if (requestId) queryParams.set("requestId", requestId);
  if (organization) queryParams.set("organization", organization);
  if (postErrorRedirectUrl) queryParams.set("postErrorRedirectUrl", postErrorRedirectUrl);
  redirect(`/idp/${provider}/account-not-found?${queryParams.toString()}`);
}
