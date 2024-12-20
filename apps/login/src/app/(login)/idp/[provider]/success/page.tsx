import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { IdpSignin } from "@/components/idp-signin";
import { idpTypeToIdentityProviderType, PROVIDER_MAPPING } from "@/lib/idp";
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
import { create } from "@zitadel/client";
import { AutoLinkingOption } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { OrganizationSchema } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import {
  AddHumanUserRequest,
  AddHumanUserRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getLocale, getTranslations } from "next-intl/server";

const ORG_SUFFIX_REGEX = /(?<=@)(.+)/;

async function loginFailed(branding?: BrandingSettings, error: string = "") {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("loginError.title")}</h1>
        <p className="ztdl-p">{t("loginError.description")}</p>
        {error && (
          <div className="w-full">
            {<Alert type={AlertType.ALERT}>{error}</Alert>}
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}

async function loginSuccess(
  userId: string,
  idpIntent: { idpIntentId: string; idpIntentToken: string },
  authRequestId?: string,
  branding?: BrandingSettings,
) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("loginSuccess.title")}</h1>
        <p className="ztdl-p">{t("loginSuccess.description")}</p>

        <IdpSignin
          userId={userId}
          idpIntent={idpIntent}
          authRequestId={authRequestId}
        />
      </div>
    </DynamicTheme>
  );
}

async function linkingSuccess(
  userId: string,
  idpIntent: { idpIntentId: string; idpIntentToken: string },
  authRequestId?: string,
  branding?: BrandingSettings,
) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("linkingSuccess.title")}</h1>
        <p className="ztdl-p">{t("linkingSuccess.description")}</p>

        <IdpSignin
          userId={userId}
          idpIntent={idpIntent}
          authRequestId={authRequestId}
        />
      </div>
    </DynamicTheme>
  );
}

async function linkingFailed(branding?: BrandingSettings) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("linkingError.title")}</h1>
        <div className="w-full">
          {
            <Alert type={AlertType.ALERT}>
              {t("linkingError.description")}
            </Alert>
          }
        </div>
      </div>
    </DynamicTheme>
  );
}

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const params = await props.params;
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });
  const { id, token, authRequestId, organization, link } = searchParams;
  const { provider } = params;

  const branding = await getBrandingSettings(organization);

  if (!provider || !id || !token) {
    return loginFailed(branding, "IDP context missing");
  }

  const intent = await retrieveIDPIntent(id, token);

  const { idpInformation, userId } = intent;

  // sign in user. If user should be linked continue
  if (userId && !link) {
    // TODO: update user if idp.options.isAutoUpdate is true

    return loginSuccess(
      userId,
      { idpIntentId: id, idpIntentToken: token },
      authRequestId,
      branding,
    );
  }

  if (!idpInformation) {
    return loginFailed(branding, "IDP information missing");
  }

  const idp = await getIDPByID(idpInformation.idpId);
  const options = idp?.config?.options;

  if (!idp) {
    throw new Error("IDP not found");
  }

  const providerType = idpTypeToIdentityProviderType(idp.type);

  // search for potential user via username, then link
  if (options?.isLinkingAllowed) {
    let foundUser;
    const email = PROVIDER_MAPPING[providerType](idpInformation).email?.email;

    if (options.autoLinking === AutoLinkingOption.EMAIL && email) {
      foundUser = await listUsers({ email }).then((response) => {
        return response.result ? response.result[0] : null;
      });
    } else if (options.autoLinking === AutoLinkingOption.USERNAME) {
      foundUser = await listUsers(
        options.autoLinking === AutoLinkingOption.USERNAME
          ? { userName: idpInformation.userName }
          : { email },
      ).then((response) => {
        return response.result ? response.result[0] : null;
      });
    } else {
      foundUser = await listUsers({
        userName: idpInformation.userName,
        email,
      }).then((response) => {
        return response.result ? response.result[0] : null;
      });
    }

    if (foundUser) {
      const idpLink = await addIDPLink(
        {
          id: idpInformation.idpId,
          userId: idpInformation.userId,
          userName: idpInformation.userName,
        },
        foundUser.userId,
      ).catch((error) => {
        return linkingFailed(branding);
      });

      if (idpLink) {
        return linkingSuccess(
          foundUser.userId,
          { idpIntentId: id, idpIntentToken: token },
          authRequestId,
          branding,
        );
      }
    }
  }

  // if link === true, do not create user
  if (options?.isCreationAllowed && options.isAutoCreation && !link) {
    let orgToRegisterOn: string | undefined = organization;

    let userData: AddHumanUserRequest =
      PROVIDER_MAPPING[providerType](idpInformation);

    if (
      !orgToRegisterOn &&
      userData.username && // username or email?
      ORG_SUFFIX_REGEX.test(userData.username)
    ) {
      const matched = ORG_SUFFIX_REGEX.exec(userData.username);
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

    if (orgToRegisterOn) {
      const organizationSchema = create(OrganizationSchema, {
        org: { case: "orgId", value: orgToRegisterOn },
      });

      userData = create(AddHumanUserRequestSchema, {
        ...userData,
        organization: organizationSchema,
      });
    }

    const newUser = await addHuman(userData);

    if (newUser) {
      return (
        <DynamicTheme branding={branding}>
          <div className="flex flex-col items-center space-y-4">
            <h1>{t("registerSuccess.title")}</h1>
            <p className="ztdl-p">{t("registerSuccess.description")}</p>
            <IdpSignin
              userId={newUser.userId}
              idpIntent={{ idpIntentId: id, idpIntentToken: token }}
              authRequestId={authRequestId}
            />
          </div>
        </DynamicTheme>
      );
    }
  }

  if (link) {
    return linkingFailed(branding);
  }

  // return login failed if no linking or creation is allowed and no user was found
  return loginFailed(branding, "No user found");
}
