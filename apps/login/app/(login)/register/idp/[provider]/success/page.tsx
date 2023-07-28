import { ProviderSlug } from "#/lib/demos";
import { addHumanUser, server } from "#/lib/zitadel";
import {
  AddHumanUserRequest,
  AddHumanUserResponse,
  IDPInformation,
  Provider,
  RetrieveIdentityProviderInformationResponse,
  user,
  IDPLink,
} from "@zitadel/server";

const PROVIDER_MAPPING: {
  [provider: string]: (rI: IDPInformation) => Partial<AddHumanUserRequest>;
} = {
  [ProviderSlug.GOOGLE]: (idp: IDPInformation) => {
    console.log("idp", idp);
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
        firstName: idp.rawInformation?.User?.given_name ?? "",
        lastName: idp.rawInformation?.User?.family_name ?? "",
      },
      idpLinks: [idpLink],
    };
    return req;
  },
  [ProviderSlug.GITHUB]: (idp: IDPInformation) => {
    console.log("idp", idp);
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
        firstName: idp.rawInformation?.name ?? "",
        lastName: idp.rawInformation?.name ?? "",
      },
      idpLinks: [idpLink],
    };
    return req;
  },
};

function retrieveIDP(
  id: string,
  token: string
): Promise<IDPInformation | undefined> {
  const userService = user.getUser(server);
  console.log("req");
  return userService
    .retrieveIdentityProviderInformation({ intentId: id, token: token }, {})
    .then((resp: RetrieveIdentityProviderInformationResponse) => {
      return resp.idpInformation;
    });
}

function createUser(
  provider: ProviderSlug,
  info: IDPInformation
): Promise<string> {
  const userData = (PROVIDER_MAPPING as any)[provider](info);
  console.log(userData);
  const userService = user.getUser(server);
  console.log(userData.profile);
  return userService.addHumanUser(userData, {}).then((resp) => resp.userId);
}

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: { provider: ProviderSlug };
}) {
  const { id, token } = searchParams;
  const { provider } = params;

  if (provider && id && token) {
    const information = await retrieveIDP(id, token);
    let user;
    if (information) {
      user = await createUser(provider, information);
    }

    return (
      <div className="flex flex-col items-center space-y-4">
        <h1>Register successful</h1>
        <p className="ztdl-p">Your account has successfully been created.</p>
        {user && <div>{JSON.stringify(user)}</div>}
      </div>
    );
  } else {
    return (
      <div className="flex flex-col items-center space-y-4">
        <h1>Register successful</h1>
        <p className="ztdl-p">No id and token received!</p>
      </div>
    );
  }
}
