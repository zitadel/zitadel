import { describe, expect, it } from "vitest";
import { instanceRolesFromClaim } from "./instance-roles";

const CLAIM = "urn:zitadel:iam:org:project:roles";
const rolesInfo = [{ organizationId: "org-1", organizationDomain: "support.example.com" }];

describe("instanceRolesFromClaim", () => {
  it("returns roles granted in a configured organization", () => {
    const raw = { [CLAIM]: { IAM_OWNER_VIEWER: { "org-1": "support.example.com" } } };
    expect(instanceRolesFromClaim(raw, rolesInfo)).toEqual(["IAM_OWNER_VIEWER"]);
  });

  it("returns multiple matching roles sorted deterministically", () => {
    const raw = {
      [CLAIM]: {
        IAM_OWNER_VIEWER: { "org-1": "support.example.com" },
        IAM_ORG_MANAGER: { "org-1": "support.example.com" },
      },
    };
    expect(instanceRolesFromClaim(raw, rolesInfo)).toEqual(["IAM_ORG_MANAGER", "IAM_OWNER_VIEWER"]);
  });

  it("ignores roles granted in an unconfigured organization", () => {
    const raw = { [CLAIM]: { IAM_OWNER_VIEWER: { "other-org": "other.example.com" } } };
    expect(instanceRolesFromClaim(raw, rolesInfo)).toEqual([]);
  });

  it("requires both organization id and domain to match", () => {
    const wrongDomain = { [CLAIM]: { IAM_OWNER_VIEWER: { "org-1": "evil.example.com" } } };
    const wrongOrg = { [CLAIM]: { IAM_OWNER_VIEWER: { "org-2": "support.example.com" } } };
    expect(instanceRolesFromClaim(wrongDomain, rolesInfo)).toEqual([]);
    expect(instanceRolesFromClaim(wrongOrg, rolesInfo)).toEqual([]);
  });

  it("ignores roles without the IAM_ prefix", () => {
    const raw = {
      [CLAIM]: {
        PROJECT_OWNER: { "org-1": "support.example.com" },
        SUPPORT_HERO: { "org-1": "support.example.com" },
      },
    };
    expect(instanceRolesFromClaim(raw, rolesInfo)).toEqual([]);
  });

  it("matches when any of multiple configured organizations grants the role", () => {
    const info = [
      { organizationId: "org-1", organizationDomain: "support.example.com" },
      { organizationId: "org-2", organizationDomain: "second.example.com" },
    ];
    const raw = { [CLAIM]: { IAM_OWNER_VIEWER: { "org-2": "second.example.com" } } };
    expect(instanceRolesFromClaim(raw, info)).toEqual(["IAM_OWNER_VIEWER"]);
  });

  it("returns empty for a missing or malformed claim", () => {
    expect(instanceRolesFromClaim(undefined, rolesInfo)).toEqual([]);
    expect(instanceRolesFromClaim(null, rolesInfo)).toEqual([]);
    expect(instanceRolesFromClaim({}, rolesInfo)).toEqual([]);
    expect(instanceRolesFromClaim({ [CLAIM]: "not-an-object" }, rolesInfo)).toEqual([]);
    expect(instanceRolesFromClaim({ [CLAIM]: ["IAM_OWNER_VIEWER"] }, rolesInfo)).toEqual([]);
    expect(instanceRolesFromClaim({ [CLAIM]: { IAM_OWNER_VIEWER: "org-1" } }, rolesInfo)).toEqual([]);
    expect(instanceRolesFromClaim({ [CLAIM]: { IAM_OWNER_VIEWER: { "org-1": 42 } } }, rolesInfo)).toEqual([]);
  });
});
