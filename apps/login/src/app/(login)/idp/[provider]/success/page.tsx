import { idpTypeToIdentityProviderType, PROVIDER_MAPPING } from "@/lib/idp";
import {
  addIDPLink,
  createUser,
  getBrandingSettings,
  getIDPByID,
  listUsers,
  retrieveIDPIntent,
} from "@/lib/zitadel";
import Alert, { AlertType } from "@/ui/Alert";
import DynamicTheme from "@/ui/DynamicTheme";
import IdpSignin from "@/ui/IdpSignin";
import { AutoLinkingOption } from "@zitadel/proto/zitadel/idp/v2/idp_pb";

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: { provider: string };
}) {
  const { id, token, authRequestId, organization } = searchParams;
  const { provider } = params;

  const branding = await getBrandingSettings(organization);
  if (provider && id && token) {
    return retrieveIDPIntent(id, token)
      .then(async (resp) => {
        const { idpInformation, userId } = resp;

        if (userId) {
          // TODO: update user if idp.options.isAutoUpdate is true

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
        }

        if (idpInformation) {
          const idp = await getIDPByID(idpInformation.idpId);
          const options = idp?.config?.options;

          if (!idp) {
            throw new Error("IDP not found");
          }

          const providerType = idpTypeToIdentityProviderType(idp.type);

          // search for potential user via username, then link
          if (options?.isLinkingAllowed) {
            let foundUser;
            const email =
              PROVIDER_MAPPING[providerType](idpInformation).email?.email;

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
                      <h1>Linking failed</h1>
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

              if (idpLink) {
                return (
                  // TODO: possibily login user now
                  <DynamicTheme branding={branding}>
                    <div className="flex flex-col items-center space-y-4">
                      <h1>Account successfully linked</h1>
                      <div>Your account has successfully been linked!</div>
                    </div>
                  </DynamicTheme>
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
                    <h1>Register successful</h1>
                    <div>You have successfully been registered!</div>
                  </div>
                </DynamicTheme>
              );
            }
          }

          // return login failed if no linking or creation is allowed and no user was found
          return (
            <DynamicTheme branding={branding}>
              <div className="flex flex-col items-center space-y-4">
                <h1>Login failed</h1>
                <div className="w-full">
                  {
                    <Alert type={AlertType.ALERT}>
                      User could not be logged in
                    </Alert>
                  }
                </div>
              </div>
            </DynamicTheme>
          );
        } else {
          return (
            <DynamicTheme branding={branding}>
              <div className="flex flex-col items-center space-y-4">
                <h1>Login failed</h1>
                <div className="w-full">
                  {
                    <Alert type={AlertType.ALERT}>
                      Could not get user information
                    </Alert>
                  }
                </div>
              </div>
            </DynamicTheme>
          );
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
