import { IDPType } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";

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
    case IdentityProviderType.LDAP:
      return "ldap";
    case IdentityProviderType.JWT:
      return "jwt";
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

    case IDPType.IDP_TYPE_JWT:
      return IdentityProviderType.JWT;

    default:
      throw new Error("Unknown identity provider type");
  }
}
