import { PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType as Method } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { beforeEach, expect, test, vi } from "vitest";
import { getLoginMethods } from "./login-methods";
import { getActiveIdentityProviders, listIDPLinks } from "./zitadel";
vi.mock("./zitadel", () => ({ getActiveIdentityProviders: vi.fn(), listIDPLinks: vi.fn() }));
const options = {
  serviceConfig: {} as any,
  userId: "user",
  organization: "org",
  authMethods: [Method.PASSWORD, Method.PASSKEY, Method.IDP],
  loginSettings: { allowLocalAuthentication: true, allowExternalIdp: true, passkeysType: PasskeysType.ALLOWED } as any,
};
beforeEach(() => {
  vi.clearAllMocks();
  vi.mocked(getActiveIdentityProviders).mockResolvedValue({
    identityProviders: [{ id: "google", name: "Google" }, { id: "unlinked" }],
  } as any);
  vi.mocked(listIDPLinks).mockResolvedValue({
    result: [{ idpId: "google", userName: "external@example.org" }, { idpId: "disabled" }],
  } as any);
});
test("offers both local methods and only linked active providers with their external login hint", async () => {
  const result = await getLoginMethods(options);
  expect(result).toMatchObject({
    password: true,
    passkey: true,
    count: 3,
    identityProviders: [{ id: "google", loginHint: "external@example.org" }],
  });
  expect(getActiveIdentityProviders).toHaveBeenCalledWith({ serviceConfig: {}, orgId: "org" });
  expect(listIDPLinks).toHaveBeenCalledWith({ serviceConfig: {}, userId: "user" });
});
test("respects local and external login policies", async () => {
  const result = await getLoginMethods({
    ...options,
    loginSettings: { allowLocalAuthentication: false, allowExternalIdp: false } as any,
  });
  expect(result.count).toBe(0);
  expect(getActiveIdentityProviders).not.toHaveBeenCalled();
  expect(listIDPLinks).not.toHaveBeenCalled();
});
test("does not offer disabled passkeys or methods not configured on the user", async () => {
  const result = await getLoginMethods({
    ...options,
    authMethods: [Method.PASSWORD],
    loginSettings: { ...options.loginSettings, passkeysType: PasskeysType.NOT_ALLOWED },
  });
  expect(result).toEqual({ password: true, passkey: false, identityProviders: [], count: 1 });
  expect(listIDPLinks).not.toHaveBeenCalled();
});
test("does not hide provider lookup failures", async () => {
  vi.mocked(listIDPLinks).mockRejectedValue(new Error("unavailable"));
  await expect(getLoginMethods(options)).rejects.toThrow("unavailable");
});

test("protected choices are based only on policy and do not inspect user links", async () => {
  const result = await getLoginMethods({ ...options, userId: undefined, authMethods: undefined });
  expect(result.count).toBe(4);
  expect(result.identityProviders.map((p) => p.id)).toEqual(["google", "unlinked"]);
  expect(result.identityProviders.every((p) => p.loginHint === undefined)).toBe(true);
  expect(listIDPLinks).not.toHaveBeenCalled();
});
