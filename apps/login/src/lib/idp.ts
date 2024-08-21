import { AddHumanUserRequest } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { IDPInformation, IDPLink } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { PartialMessage } from "@zitadel/client";

// This maps the IdentityProviderType to a slug which is used in the /success and /failure routes
export function idpTypeToSlug(idpType: IdentityProviderType) {
  switch (idpType) {
    case IdentityProviderType.GITHUB:
      return "github";
    case IdentityProviderType.GOOGLE:
      return "google";
    case IdentityProviderType.AZURE_AD:
      return "azure";
    case IdentityProviderType.SAML:
      return "saml";
    case IdentityProviderType.OIDC:
      return "oidc";
    default:
      throw new Error("Unknown identity provider type");
  }
}

// this maps the IDPInformation to the AddHumanUserRequest which is used when creating a user or linking a user (email)
export const PROVIDER_MAPPING: {
  [provider: string]: (
    rI: IDPInformation,
  ) => PartialMessage<AddHumanUserRequest>;
} = {
  [idpTypeToSlug(IdentityProviderType.GOOGLE)]: (idp: IDPInformation) => {
    const rawInfo = idp.rawInformation?.toJson() as {
      User: {
        email: string;
        name?: string;
        given_name?: string;
        family_name?: string;
      };
    };

    const idpLink: PartialMessage<IDPLink> = {
      idpId: idp.idpId,
      userId: idp.userId,
      userName: idp.userName,
    };

    const req: PartialMessage<AddHumanUserRequest> = {
      username: idp.userName,
      email: {
        email: rawInfo.User?.email,
        verification: { case: "isVerified", value: true },
      },
      // organisation: Organisation | undefined;
      profile: {
        displayName: rawInfo.User?.name ?? "",
        givenName: rawInfo.User?.given_name ?? "",
        familyName: rawInfo.User?.family_name ?? "",
      },
      idpLinks: [idpLink],
    };

    return req;
  },
  [idpTypeToSlug(IdentityProviderType.AZURE_AD)]: (idp: IDPInformation) => {
    const rawInfo = idp.rawInformation?.toJson() as {
      mail: string;
      displayName?: string;
      givenName?: string;
      surname?: string;
    };

    const idpLink: PartialMessage<IDPLink> = {
      idpId: idp.idpId,
      userId: idp.userId,
      userName: idp.userName,
    };

    const req: PartialMessage<AddHumanUserRequest> = {
      username: idp.userName,
      email: {
        email: rawInfo?.mail,
        verification: { case: "isVerified", value: true },
      },
      // organisation: Organisation | undefined;
      profile: {
        displayName: rawInfo?.displayName ?? "",
        givenName: rawInfo?.givenName ?? "",
        familyName: rawInfo?.surname ?? "",
      },
      idpLinks: [idpLink],
    };

    return req;
  },
  [idpTypeToSlug(IdentityProviderType.GITHUB)]: (idp: IDPInformation) => {
    const rawInfo = idp.rawInformation?.toJson() as {
      email: string;
      name: string;
    };

    const idpLink: PartialMessage<IDPLink> = {
      idpId: idp.idpId,
      userId: idp.userId,
      userName: idp.userName,
    };

    const req: PartialMessage<AddHumanUserRequest> = {
      username: idp.userName,
      email: {
        email: rawInfo?.email,
        verification: { case: "isVerified", value: true },
      },
      // organisation: Organisation | undefined;
      profile: {
        displayName: rawInfo?.name ?? "",
        givenName: rawInfo?.name ?? "",
        familyName: rawInfo?.name ?? "",
      },
      idpLinks: [idpLink],
    };

    return req;
  },
};
