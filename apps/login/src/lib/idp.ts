import {
  AddHumanUserRequest,
  AddHumanUserRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { IDPInformation } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { create } from "@zitadel/client";

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
// TODO: extend this object from a other file which can be overwritten by customers like map = { ...PROVIDER_MAPPING, ...customerMap }
export type OIDC_USER = {
  User: {
    email: string;
    name?: string;
    given_name?: string;
    family_name?: string;
  };
};

export const PROVIDER_MAPPING: {
  [provider: string]: (rI: IDPInformation) => AddHumanUserRequest;
} = {
  [idpTypeToSlug(IdentityProviderType.GOOGLE)]: (idp: IDPInformation) => {
    const rawInfo = idp.rawInformation as OIDC_USER;
    console.log(rawInfo);

    return create(AddHumanUserRequestSchema, {
      username: idp.userName,
      email: {
        email: rawInfo.User?.email,
        verification: { case: "isVerified", value: true },
      },
      profile: {
        displayName: rawInfo.User?.name ?? "",
        givenName: rawInfo.User?.given_name ?? "",
        familyName: rawInfo.User?.family_name ?? "",
      },
      idpLinks: [
        {
          idpId: idp.idpId,
          userId: idp.userId,
          userName: idp.userName,
        },
      ],
    });
  },
  [idpTypeToSlug(IdentityProviderType.AZURE_AD)]: (idp: IDPInformation) => {
    const rawInfo = idp.rawInformation as {
      jobTitle: string;
      mail: string;
      mobilePhone: string;
      preferredLanguage: string;
      id: string;
      displayName?: string;
      givenName?: string;
      surname?: string;
      officeLocation?: string;
      userPrincipalName: string;
    };

    console.log(rawInfo, rawInfo.userPrincipalName);

    return create(AddHumanUserRequestSchema, {
      username: idp.userName,
      email: {
        email: rawInfo.mail || rawInfo.userPrincipalName || "",
        verification: { case: "isVerified", value: true },
      },
      profile: {
        displayName: rawInfo.displayName ?? "",
        givenName: rawInfo.givenName ?? "",
        familyName: rawInfo.surname ?? "",
      },
      idpLinks: [
        {
          idpId: idp.idpId,
          userId: idp.userId,
          userName: idp.userName,
        },
      ],
    });
  },
  [idpTypeToSlug(IdentityProviderType.GITHUB)]: (idp: IDPInformation) => {
    const rawInfo = idp.rawInformation as {
      email: string;
      name: string;
    };

    return create(AddHumanUserRequestSchema, {
      username: idp.userName,
      email: {
        email: rawInfo.email,
        verification: { case: "isVerified", value: true },
      },
      profile: {
        displayName: rawInfo.name ?? "",
        givenName: rawInfo.name ?? "",
        familyName: rawInfo.name ?? "",
      },
      idpLinks: [
        {
          idpId: idp.idpId,
          userId: idp.userId,
          userName: idp.userName,
        },
      ],
    });
  },
};
