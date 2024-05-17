import { ProviderSlug } from "@/lib/demos";
import { getBrandingSettings, server } from "@/lib/zitadel";
import Alert, { AlertType } from "@/ui/Alert";
import DynamicTheme from "@/ui/DynamicTheme";
import IdpSignin from "@/ui/IdpSignin";
import {
  AddHumanUserRequest,
  IDPInformation,
  RetrieveIdentityProviderIntentResponse,
  user,
  IDPLink,
} from "@zitadel/server";
import { ClientError } from "nice-grpc";

const PROVIDER_MAPPING: {
  [provider: string]: (rI: IDPInformation) => Partial<AddHumanUserRequest>;
} = {
  [ProviderSlug.GOOGLE]: (idp: IDPInformation) => {
    const idpLink: IDPLink = {
      idpId: idp.idpId,
      userId: idp.userId,
      userName: idp.userName,
    };
    const req: Partial<AddHumanUserRequest> = {
      username: idp.userName,
      email: {
        email: idp.rawInformation?.User?.email,
        isVerified: true,
      },
      // organisation: Organisation | undefined;
      profile: {
        displayName: idp.rawInformation?.User?.name ?? "",
        givenName: idp.rawInformation?.User?.given_name ?? "",
        familyName: idp.rawInformation?.User?.family_name ?? "",
      },
      idpLinks: [idpLink],
    };
    return req;
  },
  [ProviderSlug.GITHUB]: (idp: IDPInformation) => {
    const idpLink: IDPLink = {
      idpId: idp.idpId,
      userId: idp.userId,
      userName: idp.userName,
    };
    const req: Partial<AddHumanUserRequest> = {
      username: idp.userName,
      email: {
        email: idp.rawInformation?.email,
        isVerified: true,
      },
      // organisation: Organisation | undefined;
      profile: {
        displayName: idp.rawInformation?.name ?? "",
        givenName: idp.rawInformation?.name ?? "",
        familyName: idp.rawInformation?.name ?? "",
      },
      idpLinks: [idpLink],
    };
    return req;
  },
};

function retrieveIDPIntent(
  id: string,
  token: string,
): Promise<RetrieveIdentityProviderIntentResponse> {
  const userService = user.getUser(server);
  return userService.retrieveIdentityProviderIntent(
    { idpIntentId: id, idpIntentToken: token },
    {},
  );
}

function createUser(
  provider: ProviderSlug,
  info: IDPInformation,
): Promise<string> {
  const userData = PROVIDER_MAPPING[provider](info);
  const userService = user.getUser(server);
  return userService.addHumanUser(userData, {}).then((resp) => resp.userId);
}

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: { provider: ProviderSlug };
}) {
  const { id, token, authRequestId, organization } = searchParams;
  const { provider } = params;

  const branding = await getBrandingSettings(server, organization);

  if (provider && id && token) {
    return retrieveIDPIntent(id, token)
      .then((resp) => {
        const { idpInformation, userId } = resp;

        if (idpInformation) {
          // handle login
          if (userId) {
            return (
              <DynamicTheme branding={branding}>
                <div className="flex flex-col items-center space-y-4">
                  <h1>Login successful</h1>
                  <div>You have successfully been loggedIn!</div>

                  <IdpSignin
                    userId={userId}
                    idpIntent={{ idpIntentId: id, idpIntentToken: token }}
                    authRequestId={authRequestId}
                  />
                </div>
              </DynamicTheme>
            );
          } else {
            // handle register
            return createUser(provider, idpInformation)
              .then((userId) => {
                return (
                  <DynamicTheme branding={branding}>
                    <div className="flex flex-col items-center space-y-4">
                      <h1>Register successful</h1>
                      <div>You have successfully been registered!</div>
                    </div>
                  </DynamicTheme>
                );
              })
              .catch((error: ClientError) => {
                return (
                  <DynamicTheme branding={branding}>
                    <div className="flex flex-col items-center space-y-4">
                      <h1>Register failed</h1>
                      <div className="w-full">
                        {
                          <Alert type={AlertType.ALERT}>
                            {JSON.stringify(error.details)}
                          </Alert>
                        }
                      </div>
                    </div>
                  </DynamicTheme>
                );
              });
          }
        } else {
          throw new Error("Could not get user information.");
        }
      })
      .catch((error) => {
        return (
          <DynamicTheme branding={branding}>
            <div className="flex flex-col items-center space-y-4">
              <h1>An error occurred</h1>
              <div className="w-full">
                {
                  <Alert type={AlertType.ALERT}>
                    {JSON.stringify(error.message)}
                  </Alert>
                }
              </div>
            </div>
          </DynamicTheme>
        );
      });
  } else {
    return (
      <DynamicTheme branding={branding}>
        <div className="flex flex-col items-center space-y-4">
          <div className="flex flex-col items-center space-y-4">
            <h1>Register</h1>
            <p className="ztdl-p">No id and token received!</p>
          </div>
        </div>
      </DynamicTheme>
    );
  }
}
