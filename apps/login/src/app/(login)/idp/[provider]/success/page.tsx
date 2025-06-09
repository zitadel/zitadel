import { DynamicTheme } from "@/components/dynamic-theme";
import { IdpSignin } from "@/components/idp-signin";
import { linkingFailed } from "@/components/idps/pages/linking-failed";
import { linkingSuccess } from "@/components/idps/pages/linking-success";
import { loginFailed } from "@/components/idps/pages/login-failed";
import { loginSuccess } from "@/components/idps/pages/login-success";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import {
  addHuman,
  addIDPLink,
  getBrandingSettings,
  getIDPByID,
  getLoginSettings,
  getOrgsByDomain,
  listUsers,
  retrieveIDPIntent,
} from "@/lib/zitadel";
import { ConnectError, create } from "@zitadel/client";
import { AutoLinkingOption } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { OrganizationSchema } from "@zitadel/proto/zitadel/object/v2/object_pb";
import {
  AddHumanUserRequest,
  AddHumanUserRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

const ORG_SUFFIX_REGEX = /(?<=@)(.+)/;

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const params = await props.params;
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });
  const { id, token, requestId, organization, link } = searchParams;
  const { provider } = params;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

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

  // sign in user. If user should be linked continue
  if (userId && !link) {
    // TODO: update user if idp.options.isAutoUpdate is true

    return loginSuccess(
      userId,
      { idpIntentId: id, idpIntentToken: token },
      requestId,
      branding,
    );
  }

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
      return linkingSuccess(
        userId,
        { idpIntentId: id, idpIntentToken: token },
        requestId,
        branding,
      );
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
        return linkingSuccess(
          foundUser.userId,
          { idpIntentId: id, idpIntentToken: token },
          requestId,
          branding,
        );
      }
    }
  }

  if (options?.isAutoCreation) {
    let orgToRegisterOn: string | undefined = organization;
    let newUser;

    if (
      !orgToRegisterOn &&
      addHumanUser?.username && // username or email?
      ORG_SUFFIX_REGEX.test(addHumanUser.username)
    ) {
      const matched = ORG_SUFFIX_REGEX.exec(addHumanUser.username);
      const suffix = matched?.[1] ?? "";

      // this just returns orgs where the suffix is set as primary domain
      const orgs = await getOrgsByDomain({
        serviceUrl,
        domain: suffix,
      });
      const orgToCheckForDiscovery =
        orgs.result && orgs.result.length === 1 ? orgs.result[0].id : undefined;

      const orgLoginSettings = await getLoginSettings({
        serviceUrl,
        organization: orgToCheckForDiscovery,
      });
      if (orgLoginSettings?.allowDomainDiscovery) {
        orgToRegisterOn = orgToCheckForDiscovery;
      }
    }

    if (addHumanUser) {
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
        addHumanUserWithOrganization = create(
          AddHumanUserRequestSchema,
          addHumanUser,
        );
      }

      try {
        newUser = await addHuman({
          serviceUrl,
          request: addHumanUserWithOrganization,
        });
      } catch (error: unknown) {
        console.error(
          "An error occurred while creating the user:",
          error,
          addHumanUser,
        );
        return loginFailed(
          branding,
          (error as ConnectError).message
            ? (error as ConnectError).message
            : "Could not create user",
        );
      }
    }

    if (newUser) {
      return (
        <DynamicTheme branding={branding}>
          <div className="flex flex-col items-center space-y-4">
            <h1>{t("registerSuccess.title")}</h1>
            <p className="ztdl-p">{t("registerSuccess.description")}</p>
            <IdpSignin
              userId={newUser.userId}
              idpIntent={{ idpIntentId: id, idpIntentToken: token }}
              requestId={requestId}
            />
          </div>
        </DynamicTheme>
      );
    }
  }

  // return login failed if no linking or creation is allowed and no user was found
  return loginFailed(branding, "No user found");
}
