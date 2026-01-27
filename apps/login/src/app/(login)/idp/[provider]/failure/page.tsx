import { Alert, AlertType } from "@/components/alert";
import { ChooseAuthenticatorToLogin } from "@/components/choose-authenticator-to-login";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getServiceConfig } from "@/lib/service-url";
import { getBrandingSettings, getLoginSettings, getUserByID, listAuthenticationMethodTypes } from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const searchParams = await props.searchParams;

  const { organization, userId } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const branding = await getBrandingSettings({ serviceConfig, organization,
  });

  const loginSettings = await getLoginSettings({ serviceConfig, organization,
  });

  let authMethods: AuthenticationMethodType[] = [];
  let user: User | undefined = undefined;
  let human: HumanUser | undefined = undefined;

  const params = new URLSearchParams({});
  if (organization) {
    params.set("organization", organization);
  }
  if (userId) {
    params.set("userId", userId);
  }

  if (userId) {
    const userResponse = await getUserByID({ serviceConfig, userId,
    });
    if (userResponse) {
      user = userResponse.user;
      if (user?.type.case === "human") {
        human = user.type.value as HumanUser;
      }

      if (user?.preferredLoginName) {
        params.set("loginName", user.preferredLoginName);
      }
    }

    const authMethodsResponse = await listAuthenticationMethodTypes({ serviceConfig, userId,
    });
    if (authMethodsResponse.authMethodTypes) {
      authMethods = authMethodsResponse.authMethodTypes;
    }
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="loginError.title" namespace="idp" />
        </h1>
      </div>

      <div className="w-full">
        <Alert type={AlertType.ALERT}>
          <Translated i18nKey="loginError.description" namespace="idp" />
        </Alert>

        {userId && authMethods.length && (
          <>
            {user && human && (
              <UserAvatar
                loginName={user.preferredLoginName}
                displayName={human?.profile?.displayName}
                showDropdown={false}
              />
            )}

            <ChooseAuthenticatorToLogin
              authMethods={authMethods}
              loginSettings={loginSettings}
              params={params}
            ></ChooseAuthenticatorToLogin>
          </>
        )}
      </div>
    </DynamicTheme>
  );
}
