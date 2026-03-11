import { beforeEach, describe, expect, test, vi } from "vitest";

const handlerMock = vi.fn(async () => ({ status: 200, body: "OK" }));
const createWebhookHandlerMock = vi.fn(() => handlerMock);

vi.mock("@zitadel/zitadel-js/actions/webhook", () => ({
  createWebhookHandler: createWebhookHandlerMock,
}));

describe("createZitadelWebhookHandler", () => {
  beforeEach(() => {
    vi.resetModules();
    vi.clearAllMocks();

    delete process.env.ZITADEL_WEBHOOK_PAYLOAD_TYPE;
    delete process.env.ZITADEL_WEBHOOK_SECRET;
    delete process.env.ZITADEL_WEBHOOK_JWKS_ENDPOINT;
    delete process.env.ZITADEL_WEBHOOK_JWE_PRIVATE_KEY;
  });

  test("throws when JSON payload config is missing", async () => {
    const mod = await import("./webhook.js");

    expect(() =>
      mod.createZitadelWebhookHandler({
        payloadType: "json",
        onEvent: async () => {},
      }),
    ).toThrow("ZITADEL_WEBHOOK_SECRET");
  });

  test("throws when JWT payload config is missing", async () => {
    const mod = await import("./webhook.js");

    expect(() =>
      mod.createZitadelWebhookHandler({
        payloadType: "jwt",
        onEvent: async () => {},
      }),
    ).toThrow("ZITADEL_WEBHOOK_JWKS_ENDPOINT");
  });

  test("throws when JWE private key config is missing", async () => {
    process.env.ZITADEL_WEBHOOK_JWKS_ENDPOINT =
      "https://zitadel.example.com/oauth/v2/keys";
    const mod = await import("./webhook.js");

    expect(() =>
      mod.createZitadelWebhookHandler({
        payloadType: "jwe",
        onEvent: async () => {},
      }),
    ).toThrow("ZITADEL_WEBHOOK_JWE_PRIVATE_KEY");
  });

  test("uses environment payload type and forwards request to handler", async () => {
    process.env.ZITADEL_WEBHOOK_PAYLOAD_TYPE = "jwt";
    process.env.ZITADEL_WEBHOOK_JWKS_ENDPOINT =
      "https://zitadel.example.com/oauth/v2/keys";
    const mod = await import("./webhook.js");

    handlerMock.mockResolvedValueOnce({ status: 202, body: "accepted" });

    const POST = mod.createZitadelWebhookHandler({
      onEvent: async () => {},
    });

    const response = await POST(
      new Request("https://example.com/api/webhook", {
        method: "POST",
        headers: { "x-zitadel-signature": "abc" },
        body: '{"hello":"world"}',
      }),
    );

    expect(createWebhookHandlerMock).toHaveBeenCalledWith(
      expect.objectContaining({
        payloadType: "jwt",
        jwksEndpoint: "https://zitadel.example.com/oauth/v2/keys",
      }),
    );
    expect(handlerMock).toHaveBeenCalledWith(
      expect.objectContaining({
        body: '{"hello":"world"}',
      }),
    );
    expect(response.status).toBe(202);
    await expect(response.text()).resolves.toBe("accepted");
  });
});
