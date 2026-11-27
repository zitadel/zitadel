import { beforeEach, describe, expect, test, vi } from "vitest";

vi.mock("next-intl/server", () => ({
  getRequestConfig: (fn: unknown) => fn,
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(),
  cookies: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://api.example.com" } })),
}));

vi.mock("@/lib/server/zitadel", () => ({
  getAllowedLanguages: vi.fn(),
  getHostedLoginTranslation: vi.fn(),
}));

type Config = { locale: string; messages: any };

async function loadConfig(requestHeaders: Record<string, string>): Promise<Config> {
  const { headers, cookies } = await import("next/headers");
  vi.mocked(headers).mockResolvedValue(new Headers(requestHeaders) as any);
  vi.mocked(cookies).mockResolvedValue({ get: () => undefined } as any);

  const getConfig = (await import("./request")).default as unknown as () => Promise<Config>;
  return getConfig();
}

describe("i18n request config", () => {
  let zitadel: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = vi.mocked(await import("@/lib/server/zitadel"));
    zitadel.getAllowedLanguages.mockResolvedValue({ allowedLanguages: [], defaultLanguage: "en" });
  });

  test("lets the organization texts win over the instance texts and both win over the local files", async () => {
    zitadel.getHostedLoginTranslation.mockImplementation(async ({ organization }: { organization?: string }) =>
      organization
        ? { register: { labels: { lastname: "Org last name" } } }
        : { register: { labels: { firstname: "Instance first name", lastname: "Instance last name" } } },
    );

    const { messages } = await loadConfig({ "x-zitadel-i18n-organization": "org-1" });

    expect(zitadel.getHostedLoginTranslation).toHaveBeenCalledWith(expect.objectContaining({ organization: undefined }));
    expect(zitadel.getHostedLoginTranslation).toHaveBeenCalledWith(expect.objectContaining({ organization: "org-1" }));
    expect(messages.register.labels).toEqual({
      firstname: "Instance first name",
      lastname: "Org last name",
      email: "E-mail",
    });
  });

  test("loads only the instance texts without an organization", async () => {
    zitadel.getHostedLoginTranslation.mockResolvedValue({ loginname: { title: "Instance title" } });

    const { messages } = await loadConfig({});

    expect(zitadel.getHostedLoginTranslation).toHaveBeenCalledTimes(1);
    expect(zitadel.getHostedLoginTranslation).toHaveBeenCalledWith(expect.objectContaining({ organization: undefined }));
    expect(messages.loginname.title).toBe("Instance title");
  });

  test("keeps the instance texts when the organization texts can't be loaded", async () => {
    zitadel.getHostedLoginTranslation.mockImplementation(async ({ organization }: { organization?: string }) => {
      if (organization) {
        throw new Error("unavailable");
      }
      return { loginname: { title: "Instance title" } };
    });

    const { messages } = await loadConfig({ "x-zitadel-i18n-organization": "org-1" });

    expect(messages.loginname.title).toBe("Instance title");
    expect(messages.loginname.submit).toBe("Continue");
  });
});
