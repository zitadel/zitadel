import { create } from "@zitadel/client";
import { IDPType } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { IDPInformation } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import {
  AddHumanUserRequest,
  AddHumanUserRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";

// This maps the IdentityProviderType to a slug which is used in the /success and /failure routes
export function idpTypeToSlug(idpType: IdentityProviderType) {
  switch (idpType) {
    case IdentityProviderType.GITHUB:
      return "github";
    case IdentityProviderType.GITHUB_ES:
      return "github_es";
    case IdentityProviderType.GITLAB:
      return "gitlab";
    case IdentityProviderType.GITLAB_SELF_HOSTED:
      return "gitlab_es";
    case IdentityProviderType.APPLE:
      return "apple";
    case IdentityProviderType.GOOGLE:
      return "google";
    case IdentityProviderType.AZURE_AD:
      return "azure";
    case IdentityProviderType.SAML:
      return "saml";
    case IdentityProviderType.OAUTH:
      return "oauth";
    case IdentityProviderType.OIDC:
      return "oidc";
    default:
      throw new Error("Unknown identity provider type");
  }
}

// TODO: this is ugly but needed atm as the getIDPByID returns a IDPType and not a IdentityProviderType
export function idpTypeToIdentityProviderType(
  idpType: IDPType,
): IdentityProviderType {
  switch (idpType) {
    case IDPType.IDP_TYPE_GITHUB:
      return IdentityProviderType.GITHUB;

    case IDPType.IDP_TYPE_GITHUB_ES:
      return IdentityProviderType.GITHUB_ES;

    case IDPType.IDP_TYPE_GITLAB:
      return IdentityProviderType.GITLAB;

    case IDPType.IDP_TYPE_GITLAB_SELF_HOSTED:
      return IdentityProviderType.GITLAB_SELF_HOSTED;

    case IDPType.IDP_TYPE_APPLE:
      return IdentityProviderType.APPLE;

    case IDPType.IDP_TYPE_GOOGLE:
      return IdentityProviderType.GOOGLE;

    case IDPType.IDP_TYPE_AZURE_AD:
      return IdentityProviderType.AZURE_AD;

    case IDPType.IDP_TYPE_SAML:
      return IdentityProviderType.SAML;

    case IDPType.IDP_TYPE_OAUTH:
      return IdentityProviderType.OAUTH;

    case IDPType.IDP_TYPE_OIDC:
      return IdentityProviderType.OIDC;

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

const GITLAB_MAPPING = (idp: IDPInformation) => {
  const rawInfo = idp.rawInformation as {
    name: string;
    email: string;
    email_verified: boolean;
  };

  return create(AddHumanUserRequestSchema, {
    username: idp.userName,
    email: {
      email: rawInfo.email,
      verification: { case: "isVerified", value: rawInfo.email_verified },
    },
    profile: {
      displayName: rawInfo.name || idp.userName || "",
      givenName: "",
      familyName: "",
    },
    idpLinks: [
      {
        idpId: idp.idpId,
        userId: idp.userId,
        userName: idp.userName,
      },
    ],
  });
};

const OIDC_MAPPING = (idp: IDPInformation) => {
  const rawInfo = idp.rawInformation as OIDC_USER;

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
};

const GITHUB_MAPPING = (idp: IDPInformation) => {
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
};

export const PROVIDER_MAPPING: {
  [provider: number]: (rI: IDPInformation) => AddHumanUserRequest;
} = {
  [IdentityProviderType.GOOGLE]: (idp: IDPInformation) => {
    const rawInfo = idp.rawInformation as OIDC_USER;

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
  [IdentityProviderType.GITLAB]: GITLAB_MAPPING,
  [IdentityProviderType.GITLAB_SELF_HOSTED]: GITLAB_MAPPING,
  [IdentityProviderType.OIDC]: OIDC_MAPPING,
  // check
  [IdentityProviderType.OAUTH]: OIDC_MAPPING,
  [IdentityProviderType.AZURE_AD]: (idp: IDPInformation) => {
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
  [IdentityProviderType.GITHUB]: GITHUB_MAPPING,
  [IdentityProviderType.GITHUB_ES]: GITHUB_MAPPING,
  [IdentityProviderType.APPLE]: (idp: IDPInformation) => {
    const rawInfo = idp.rawInformation as {
      name?: string;
      firstName?: string;
      lastName?: string;
      email?: string;
    };

    return create(AddHumanUserRequestSchema, {
      username: idp.userName,
      email: {
        email: rawInfo.email ?? "",
        verification: { case: "isVerified", value: true },
      },
      profile: {
        displayName: rawInfo.name ?? "",
        givenName: rawInfo.firstName ?? "",
        familyName: rawInfo.lastName ?? "",
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
