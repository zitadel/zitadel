import { TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { createServiceForHost } from "./service";
import { getOrgsByDomain, listUsers, searchUsers } from "./zitadel";

vi.mock("./service", () => ({
  createServiceForHost: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

const serviceConfig = {} as never;
const emptyResponse = { details: { totalResult: BigInt(0) }, result: [] };

type Request = { queries: any[] };

/** Captures the requests that the login app sends to the API. */
function mockService() {
  const listUsersMock = vi.fn(async (_request: Request) => emptyResponse);
  const listOrganizationsMock = vi.fn(async (_request: Request) => emptyResponse);
  vi.mocked(createServiceForHost).mockResolvedValue({
    listUsers: listUsersMock,
    listOrganizations: listOrganizationsMock,
  } as never);
  return { listUsersMock, listOrganizationsMock };
}

/** Finds a query by its oneof case in a (possibly nested) query list. */
function findQuery(queries: any[], queryCase: string): any {
  for (const entry of queries) {
    if (entry.query?.case === queryCase) {
      return entry.query.value;
    }
    if (entry.query?.case === "orQuery" || entry.query?.case === "andQuery") {
      const nested = findQuery(entry.query.value.queries, queryCase);
      if (nested) {
        return nested;
      }
    }
  }
  return undefined;
}

beforeEach(() => {
  vi.clearAllMocks();
});

// A user cannot be stored twice with different casing: the username unique
// constraint is lowercased when it is written. Looking users up case sensitively
// therefore reports accounts as missing that provably exist — which is what breaks
// IdP auto-linking when the provider returns a different casing than was stored.
describe("listUsers", () => {
  it("should look up a login name ignoring case", async () => {
    const { listUsersMock } = mockService();

    await listUsers({ serviceConfig, loginName: "user@example.com" });

    const query = findQuery(listUsersMock.mock.calls[0][0].queries, "loginNameQuery");
    expect(query.loginName).toBe("user@example.com");
    expect(query.method).toBe(TextQueryMethod.EQUALS_IGNORE_CASE);
  });

  it("should look up an email ignoring case", async () => {
    const { listUsersMock } = mockService();

    await listUsers({ serviceConfig, email: "BREcloud@example.com" });

    const query = findQuery(listUsersMock.mock.calls[0][0].queries, "emailQuery");
    expect(query.emailAddress).toBe("BREcloud@example.com");
    expect(query.method).toBe(TextQueryMethod.EQUALS_IGNORE_CASE);
  });

  it("should look up a username ignoring case", async () => {
    const { listUsersMock } = mockService();

    await listUsers({ serviceConfig, userName: "BREcloud" });

    const query = findQuery(listUsersMock.mock.calls[0][0].queries, "userNameQuery");
    expect(query.method).toBe(TextQueryMethod.EQUALS_IGNORE_CASE);
  });

  // Phone numbers have no casing, so there is nothing to fold.
  it("should match a phone number exactly", async () => {
    const { listUsersMock } = mockService();

    await listUsers({ serviceConfig, phone: "+41791234567" });

    const query = findQuery(listUsersMock.mock.calls[0][0].queries, "phoneQuery");
    expect(query.method).toBe(TextQueryMethod.EQUALS);
  });
});

describe("searchUsers", () => {
  const loginSettings = {} as LoginSettings;

  it("should trim the search value before looking the user up", async () => {
    const { listUsersMock } = mockService();

    await searchUsers({ serviceConfig, searchValue: "  user@example.com \n", loginSettings });

    const query = findQuery(listUsersMock.mock.calls[0][0].queries, "loginNameQuery");
    expect(query.loginName).toBe("user@example.com");
  });

  it("should trim before concatenating the organization suffix", async () => {
    const { listUsersMock } = mockService();

    await searchUsers({ serviceConfig, searchValue: " user ", loginSettings, suffix: "example.com" });

    const query = findQuery(listUsersMock.mock.calls[0][0].queries, "loginNameQuery");
    expect(query.loginName).toBe("user@example.com");
  });

  it("should keep looking up login name and email ignoring case", async () => {
    const { listUsersMock } = mockService();

    await searchUsers({ serviceConfig, searchValue: "User@Example.com", loginSettings });

    expect(findQuery(listUsersMock.mock.calls[0][0].queries, "loginNameQuery").method).toBe(
      TextQueryMethod.EQUALS_IGNORE_CASE,
    );
    // no user was found by login name, so the email/phone fallback ran
    expect(findQuery(listUsersMock.mock.calls[1][0].queries, "emailQuery").method).toBe(TextQueryMethod.EQUALS_IGNORE_CASE);
  });
});

// Domains are case insensitive, and this one is the suffix of a user supplied login
// name, so it arrives in whatever casing was typed.
describe("getOrgsByDomain", () => {
  it("should look the domain up ignoring case", async () => {
    const { listOrganizationsMock } = mockService();

    await getOrgsByDomain({ serviceConfig, domain: "Example.com" });

    const query = findQuery(listOrganizationsMock.mock.calls[0][0].queries, "domainQuery");
    expect(query.domain).toBe("Example.com");
    expect(query.method).toBe(TextQueryMethod.EQUALS_IGNORE_CASE);
  });
});
