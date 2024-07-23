import { ProviderSlug } from "@/lib/demos";
import {
  addIDPLink,
  createUser,
  getBrandingSettings,
  retrieveIDPIntent,
} from "@/lib/zitadel";
import Alert, { AlertType } from "@/ui/Alert";
import DynamicTheme from "@/ui/DynamicTheme";
import IdpSignin from "@/ui/IdpSignin";

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: { provider: ProviderSlug };
}) {
  const { id, token, authRequestId, organization } = searchParams;
  const { provider } = params;

  const branding = await getBrandingSettings(organization);
  if (provider && id && token) {
    return retrieveIDPIntent(id, token)
      .then((resp) => {
        const { idpInformation, userId } = resp;

        if (idpInformation) {
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
              .catch((error) => {
                if (error.code === 6) {
                  return addIDPLink(
                    {
                      id: idpInformation.idpId,
                      userId: idpInformation.userId,
                      userName: idpInformation.userName,
                    },
                    userId,
                  ).then(() => {
                    return (
                      <DynamicTheme branding={branding}>
                        <div className="flex flex-col items-center space-y-4">
                          <h1>Account successfully linked</h1>
                          <div>Your account has successfully been linked!</div>
                        </div>
                      </DynamicTheme>
                    );
                  });
                } else {
                  return (
                    <DynamicTheme branding={branding}>
                      <div className="flex flex-col items-center space-y-4">
                        <h1>Register failed</h1>
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
                }
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
