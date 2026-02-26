import { timingSafeEqual, createHmac } from "node:crypto";

export interface WebhookVerifyOptions {
  /** The raw request body as a string or Buffer. */
  body: string | Buffer;
  /** The signature header value from the incoming request. */
  signature: string;
  /** The shared webhook signing key. */
  signingKey: string;
}

/**
 * Verifies a ZITADEL webhook signature using HMAC-SHA256.
 * Returns true if the signature is valid.
 */
export function verifyZitadelWebhook(options: WebhookVerifyOptions): boolean {
  const { body, signature, signingKey } = options;
  const expected = createHmac("sha256", signingKey)
    .update(body)
    .digest("hex");

  const sigBuffer = Buffer.from(signature, "hex");
  const expectedBuffer = Buffer.from(expected, "hex");

  if (sigBuffer.length !== expectedBuffer.length) {
    return false;
  }

  return timingSafeEqual(sigBuffer, expectedBuffer);
}

export interface WebhookHandlerOptions {
  /** The shared webhook signing key. */
  signingKey: string;
  /** Callback invoked with the parsed webhook payload when signature is valid. */
  onEvent: (event: unknown) => void | Promise<void>;
}

/**
 * Creates a webhook handler that verifies the signature and parses the body.
 */
export function createWebhookHandler(options: WebhookHandlerOptions) {
  return async (request: { body: string; headers: Record<string, string> }) => {
    const signature = request.headers["x-zitadel-signature"] ?? "";
    const isValid = verifyZitadelWebhook({
      body: request.body,
      signature,
      signingKey: options.signingKey,
    });

    if (!isValid) {
      return { status: 401, body: "Invalid signature" };
    }

    const event = JSON.parse(request.body);
    await options.onEvent(event);
    return { status: 200, body: "OK" };
  };
}
