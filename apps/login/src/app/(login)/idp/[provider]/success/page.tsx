import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { IdpSignin } from "@/components/idp-signin";
import { idpTypeToIdentityProviderType, PROVIDER_MAPPING } from "@/lib/idp";
import {
  addIDPLink,
  createUser,
  getBrandingSettings,
  getIDPByID,
  listUsers,
  retrieveIDPIntent,
} from "@/lib/zitadel";
import { AutoLinkingOption } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { getLocale, getTranslations } from "next-intl/server";

async function loginFailed(branding?: BrandingSettings) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("loginError.title")}</h1>
        <div className="w-full">
          {<Alert type={AlertType.ALERT}>{t("loginError.title")}</Alert>}
        </div>
      </div>
    </DynamicTheme>
  );
}
export default async function Page(
  props: {
    searchParams: Promise<Record<string | number | symbol, string | undefined>>;
    params: Promise<{ provider: string }>;
  }
) {
  const params = await props.params;
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });
  const { id, token, authRequestId, organization } = searchParams;
  const { provider } = params;

  const branding = await getBrandingSettings(organization);

  if (!provider || !id || !token) {
    return loginFailed(branding);
  }

  const intent = await retrieveIDPIntent(id, token);

  const { idpInformation, userId } = intent;

  if (userId) {
    // TODO: update user if idp.options.isAutoUpdate is true

    return (
      <DynamicTheme branding={branding}>
        <div className="flex flex-col items-center space-y-4">
          <h1>{t("loginSuccess.title")}</h1>
          <div>{t("loginSuccess.description")}</div>

          <IdpSignin
            userId={userId}
            idpIntent={{ idpIntentId: id, idpIntentToken: token }}
            authRequestId={authRequestId}
          />
        </div>
      </DynamicTheme>
    );
  }

  if (!idpInformation) {
    return loginFailed(branding);
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
      });

      if (idpLink) {
        return (
          // TODO: possibily login user now
          (<DynamicTheme branding={branding}>
            <div className="flex flex-col items-center space-y-4">
              <h1>{t("linkingSuccess.title")}</h1>
              <div>{t("linkingSuccess.description")}</div>
            </div>
          </DynamicTheme>)
        );
      }
    }
  }

  if (options?.isCreationAllowed && options.isAutoCreation) {
    const newUser = await createUser(providerType, idpInformation);

    if (newUser) {
      return (
        <DynamicTheme branding={branding}>
          <div className="flex flex-col items-center space-y-4">
            <h1>{t("registerSuccess.title")}</h1>
            <div>{t("registerSuccess.description")}</div>
          </div>
        </DynamicTheme>
      );
    }
  }

  // return login failed if no linking or creation is allowed and no user was found
  return loginFailed;
}
