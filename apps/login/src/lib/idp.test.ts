import { describe, it, expect } from "vitest";
import { idpTypeToSlug, idpTypeToIdentityProviderType } from "./idp";
import { IDPType } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";

describe("idp type mapping utilities", () => {
  describe("idpTypeToSlug", () => {
    it("should map GITHUB to 'github'", () => {
      const result = idpTypeToSlug(IdentityProviderType.GITHUB);
      expect(result).toBe("github");
    });

    it("should map GITHUB_ES to 'github_es'", () => {
      const result = idpTypeToSlug(IdentityProviderType.GITHUB_ES);
      expect(result).toBe("github_es");
    });

    it("should map GITLAB to 'gitlab'", () => {
      const result = idpTypeToSlug(IdentityProviderType.GITLAB);
      expect(result).toBe("gitlab");
    });

    it("should map GITLAB_SELF_HOSTED to 'gitlab_es'", () => {
      const result = idpTypeToSlug(IdentityProviderType.GITLAB_SELF_HOSTED);
      expect(result).toBe("gitlab_es");
    });

    it("should map APPLE to 'apple'", () => {
      const result = idpTypeToSlug(IdentityProviderType.APPLE);
      expect(result).toBe("apple");
    });

    it("should map GOOGLE to 'google'", () => {
      const result = idpTypeToSlug(IdentityProviderType.GOOGLE);
      expect(result).toBe("google");
    });

    it("should map AZURE_AD to 'azure'", () => {
      const result = idpTypeToSlug(IdentityProviderType.AZURE_AD);
      expect(result).toBe("azure");
    });

    it("should map SAML to 'saml'", () => {
      const result = idpTypeToSlug(IdentityProviderType.SAML);
      expect(result).toBe("saml");
    });

    it("should map OAUTH to 'oauth'", () => {
      const result = idpTypeToSlug(IdentityProviderType.OAUTH);
      expect(result).toBe("oauth");
    });

    it("should map OIDC to 'oidc'", () => {
      const result = idpTypeToSlug(IdentityProviderType.OIDC);
      expect(result).toBe("oidc");
    });

    it("should map LDAP to 'ldap'", () => {
      const result = idpTypeToSlug(IdentityProviderType.LDAP);
      expect(result).toBe("ldap");
    });

    it("should map JWT to 'jwt'", () => {
      const result = idpTypeToSlug(IdentityProviderType.JWT);
      expect(result).toBe("jwt");
    });

    it("should throw error for unknown identity provider type", () => {
      // Using a value that doesn't match any known type
      const unknownType = 9999 as IdentityProviderType;

      expect(() => idpTypeToSlug(unknownType)).toThrow("Unknown identity provider type");
    });

    it("should throw error for undefined type", () => {
      const undefinedType = undefined as any as IdentityProviderType;
      expect(() => idpTypeToSlug(undefinedType)).toThrow();
    });

    it("should handle all defined IdentityProviderType enum values", () => {
      // This test ensures all enum values are handled
      const handledTypes = [
        IdentityProviderType.GITHUB,
        IdentityProviderType.GITHUB_ES,
        IdentityProviderType.GITLAB,
        IdentityProviderType.GITLAB_SELF_HOSTED,
        IdentityProviderType.APPLE,
        IdentityProviderType.GOOGLE,
        IdentityProviderType.AZURE_AD,
        IdentityProviderType.SAML,
        IdentityProviderType.OAUTH,
        IdentityProviderType.OIDC,
        IdentityProviderType.LDAP,
        IdentityProviderType.JWT,
      ];

      handledTypes.forEach((type) => {
        expect(() => idpTypeToSlug(type)).not.toThrow();
      });
    });

    it("should return consistent slugs for the same input", () => {
      const type = IdentityProviderType.GITHUB;
      const result1 = idpTypeToSlug(type);
      const result2 = idpTypeToSlug(type);

      expect(result1).toBe(result2);
    });

    it("should return lowercase slugs", () => {
      const types = [
        IdentityProviderType.GITHUB,
        IdentityProviderType.GOOGLE,
        IdentityProviderType.AZURE_AD,
        IdentityProviderType.SAML,
      ];

      types.forEach((type) => {
        const slug = idpTypeToSlug(type);
        expect(slug).toBe(slug.toLowerCase());
      });
    });
  });

  describe("idpTypeToIdentityProviderType", () => {
    it("should map IDP_TYPE_GITHUB to GITHUB", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_GITHUB);
      expect(result).toBe(IdentityProviderType.GITHUB);
    });

    it("should map IDP_TYPE_GITHUB_ES to GITHUB_ES", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_GITHUB_ES);
      expect(result).toBe(IdentityProviderType.GITHUB_ES);
    });

    it("should map IDP_TYPE_GITLAB to GITLAB", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_GITLAB);
      expect(result).toBe(IdentityProviderType.GITLAB);
    });

    it("should map IDP_TYPE_GITLAB_SELF_HOSTED to GITLAB_SELF_HOSTED", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_GITLAB_SELF_HOSTED);
      expect(result).toBe(IdentityProviderType.GITLAB_SELF_HOSTED);
    });

    it("should map IDP_TYPE_APPLE to APPLE", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_APPLE);
      expect(result).toBe(IdentityProviderType.APPLE);
    });

    it("should map IDP_TYPE_GOOGLE to GOOGLE", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_GOOGLE);
      expect(result).toBe(IdentityProviderType.GOOGLE);
    });

    it("should map IDP_TYPE_AZURE_AD to AZURE_AD", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_AZURE_AD);
      expect(result).toBe(IdentityProviderType.AZURE_AD);
    });

    it("should map IDP_TYPE_SAML to SAML", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_SAML);
      expect(result).toBe(IdentityProviderType.SAML);
    });

    it("should map IDP_TYPE_OAUTH to OAUTH", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_OAUTH);
      expect(result).toBe(IdentityProviderType.OAUTH);
    });

    it("should map IDP_TYPE_OIDC to OIDC", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_OIDC);
      expect(result).toBe(IdentityProviderType.OIDC);
    });

    it("should map IDP_TYPE_JWT to JWT", () => {
      const result = idpTypeToIdentityProviderType(IDPType.IDP_TYPE_JWT);
      expect(result).toBe(IdentityProviderType.JWT);
    });

    it("should throw error for unknown IDP type", () => {
      // Using a value that doesn't match any known type
      const unknownType = 9999 as IDPType;

      expect(() => idpTypeToIdentityProviderType(unknownType)).toThrow("Unknown identity provider type");
    });

    it("should throw error for undefined type", () => {
      const undefinedType = undefined as any as IDPType;
      expect(() => idpTypeToIdentityProviderType(undefinedType)).toThrow();
    });

    it("should handle all defined IDPType enum values", () => {
      const handledTypes = [
        IDPType.IDP_TYPE_GITHUB,
        IDPType.IDP_TYPE_GITHUB_ES,
        IDPType.IDP_TYPE_GITLAB,
        IDPType.IDP_TYPE_GITLAB_SELF_HOSTED,
        IDPType.IDP_TYPE_APPLE,
        IDPType.IDP_TYPE_GOOGLE,
        IDPType.IDP_TYPE_AZURE_AD,
        IDPType.IDP_TYPE_SAML,
        IDPType.IDP_TYPE_OAUTH,
        IDPType.IDP_TYPE_OIDC,
        IDPType.IDP_TYPE_JWT,
      ];

      handledTypes.forEach((type) => {
        expect(() => idpTypeToIdentityProviderType(type)).not.toThrow();
      });
    });

    it("should return consistent results for the same input", () => {
      const type = IDPType.IDP_TYPE_GITHUB;
      const result1 = idpTypeToIdentityProviderType(type);
      const result2 = idpTypeToIdentityProviderType(type);

      expect(result1).toBe(result2);
    });

    it("should return valid IdentityProviderType enum values", () => {
      const types = [IDPType.IDP_TYPE_GITHUB, IDPType.IDP_TYPE_GOOGLE, IDPType.IDP_TYPE_AZURE_AD];

      types.forEach((type) => {
        const result = idpTypeToIdentityProviderType(type);
        expect(Object.values(IdentityProviderType)).toContain(result);
      });
    });
  });

  describe("round-trip conversions", () => {
    it("should successfully convert IDPType -> IdentityProviderType -> slug", () => {
      const idpType = IDPType.IDP_TYPE_GITHUB;
      const identityProviderType = idpTypeToIdentityProviderType(idpType);
      const slug = idpTypeToSlug(identityProviderType);

      expect(slug).toBe("github");
    });

    it("should handle enterprise providers in round-trip", () => {
      const idpType = IDPType.IDP_TYPE_AZURE_AD;
      const identityProviderType = idpTypeToIdentityProviderType(idpType);
      const slug = idpTypeToSlug(identityProviderType);

      expect(slug).toBe("azure");
    });

    it("should handle self-hosted variants", () => {
      const githubES = IDPType.IDP_TYPE_GITHUB_ES;
      const gitlabSH = IDPType.IDP_TYPE_GITLAB_SELF_HOSTED;

      const githubESType = idpTypeToIdentityProviderType(githubES);
      const gitlabSHType = idpTypeToIdentityProviderType(gitlabSH);

      const githubESSlug = idpTypeToSlug(githubESType);
      const gitlabSHSlug = idpTypeToSlug(gitlabSHType);

      expect(githubESSlug).toBe("github_es");
      expect(gitlabSHSlug).toBe("gitlab_es");
    });

    it("should handle protocol-based providers", () => {
      const providers = [
        { idpType: IDPType.IDP_TYPE_OIDC, expectedSlug: "oidc" },
        { idpType: IDPType.IDP_TYPE_OAUTH, expectedSlug: "oauth" },
        { idpType: IDPType.IDP_TYPE_SAML, expectedSlug: "saml" },
        { idpType: IDPType.IDP_TYPE_JWT, expectedSlug: "jwt" },
      ];

      providers.forEach(({ idpType, expectedSlug }) => {
        const identityProviderType = idpTypeToIdentityProviderType(idpType);
        const slug = idpTypeToSlug(identityProviderType);
        expect(slug).toBe(expectedSlug);
      });
    });

    it("should handle all supported providers in complete round-trip", () => {
      const allProviders = [
        { idpType: IDPType.IDP_TYPE_GITHUB, slug: "github" },
        { idpType: IDPType.IDP_TYPE_GITHUB_ES, slug: "github_es" },
        { idpType: IDPType.IDP_TYPE_GITLAB, slug: "gitlab" },
        { idpType: IDPType.IDP_TYPE_GITLAB_SELF_HOSTED, slug: "gitlab_es" },
        { idpType: IDPType.IDP_TYPE_APPLE, slug: "apple" },
        { idpType: IDPType.IDP_TYPE_GOOGLE, slug: "google" },
        { idpType: IDPType.IDP_TYPE_AZURE_AD, slug: "azure" },
        { idpType: IDPType.IDP_TYPE_SAML, slug: "saml" },
        { idpType: IDPType.IDP_TYPE_OAUTH, slug: "oauth" },
        { idpType: IDPType.IDP_TYPE_OIDC, slug: "oidc" },
        { idpType: IDPType.IDP_TYPE_JWT, slug: "jwt" },
      ];

      allProviders.forEach(({ idpType, slug }) => {
        const identityProviderType = idpTypeToIdentityProviderType(idpType);
        const actualSlug = idpTypeToSlug(identityProviderType);
        expect(actualSlug).toBe(slug);
      });
    });
  });

  describe("slug uniqueness", () => {
    it("should generate unique slugs for different providers", () => {
      const types = [
        IdentityProviderType.GITHUB,
        IdentityProviderType.GITLAB,
        IdentityProviderType.GOOGLE,
        IdentityProviderType.APPLE,
        IdentityProviderType.AZURE_AD,
        IdentityProviderType.SAML,
        IdentityProviderType.OIDC,
        IdentityProviderType.OAUTH,
        IdentityProviderType.LDAP,
        IdentityProviderType.JWT,
      ];

      const slugs = types.map((type) => idpTypeToSlug(type));
      const uniqueSlugs = new Set(slugs);

      expect(uniqueSlugs.size).toBe(slugs.length);
    });

    it("should use consistent naming convention for enterprise variants", () => {
      const githubES = idpTypeToSlug(IdentityProviderType.GITHUB_ES);
      const gitlabSH = idpTypeToSlug(IdentityProviderType.GITLAB_SELF_HOSTED);

      // Both use _es suffix for enterprise/self-hosted
      expect(githubES).toContain("_es");
      expect(gitlabSH).toContain("_es");
    });
  });

  describe("error handling", () => {
    it("should provide clear error messages", () => {
      const unknownType = 9999 as IdentityProviderType;

      expect(() => idpTypeToSlug(unknownType)).toThrow(/Unknown identity provider type/);
    });

    it("should handle undefined gracefully with error", () => {
      expect(() => idpTypeToSlug(undefined as any)).toThrow();
      expect(() => idpTypeToIdentityProviderType(undefined as any)).toThrow();
    });

    it("should handle null gracefully with error", () => {
      expect(() => idpTypeToSlug(null as any)).toThrow();
      expect(() => idpTypeToIdentityProviderType(null as any)).toThrow();
    });
  });
});
