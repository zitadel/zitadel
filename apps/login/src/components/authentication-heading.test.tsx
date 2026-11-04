import { getAuthRequest } from "@/lib/zitadel";
import { cleanup, render } from "@testing-library/react";
import { create } from "@zitadel/client";
import { Prompt } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { GetAuthRequestResponseSchema } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { AuthenticationHeading } from "./authentication-heading";

vi.mock("@/lib/zitadel", () => ({ getAuthRequest: vi.fn() }));
vi.mock("next-intl", () => ({ useTranslations: (namespace: string) => (key: string) => `${namespace}.${key}` }));

const serviceConfig = { baseUrl: "https://idp.example.com" };

describe("AuthenticationHeading", () => {
  beforeEach(() => vi.resetAllMocks());
  afterEach(cleanup);

  test.each(["loginname", "password"] as const)("preserves ordinary %s headings", async (namespace) => {
    vi.mocked(getAuthRequest).mockResolvedValue(create(GetAuthRequestResponseSchema, { authRequest: { prompt: [] } }));
    const view = render(await AuthenticationHeading({ namespace, requestId: "oidc_V2_123", serviceConfig }));
    expect(view.getByRole("heading", { level: 1 })).toHaveTextContent(
      namespace === "password" ? "password.verify.title" : "loginname.title",
    );
  });

  test.each(["loginname", "password"] as const)("explains fresh authentication on the %s step", async (namespace) => {
    vi.mocked(getAuthRequest).mockResolvedValue(
      create(GetAuthRequestResponseSchema, { authRequest: { prompt: [Prompt.LOGIN] } }),
    );
    const view = render(await AuthenticationHeading({ namespace, requestId: "oidc_V2_123", serviceConfig }));
    expect(view.getByRole("heading", { level: 1 })).toHaveTextContent("reauthenticate.title");
    expect(view.getByText("reauthenticate.description")).toBeInTheDocument();
    expect(getAuthRequest).toHaveBeenCalledExactlyOnceWith({ serviceConfig, authRequestId: "V2_123" });
  });

  test("does not require an OIDC request for a direct or SAML login", async () => {
    await AuthenticationHeading({ namespace: "loginname", serviceConfig });
    await AuthenticationHeading({ namespace: "loginname", serviceConfig, requestId: "saml_123" });
    expect(getAuthRequest).not.toHaveBeenCalled();
  });

  test("lets the shared error boundary handle an invalid OIDC request", async () => {
    vi.mocked(getAuthRequest).mockRejectedValue(new Error("Request expired"));
    await expect(AuthenticationHeading({ namespace: "loginname", serviceConfig, requestId: "oidc_V2_123" })).rejects.toThrow(
      "Request expired",
    );
  });
});
