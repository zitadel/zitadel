export interface WebhookHandlerOptions {
  /** The shared webhook signing key. */
  signingKey: string;
  /** Callback invoked with the parsed webhook payload when signature is valid. */
  onEvent: (event: unknown) => void | Promise<void>;
}

/**
 * Creates a webhook handler for use in Angular server-side routes.
 * Placeholder — to be implemented.
 */
export function createWebhookHandler(_options: WebhookHandlerOptions) {
  return async (_request: unknown) => {
    // TODO: implement webhook verification
    return { status: 200, body: "OK" };
  };
}
