export interface WebhookHandlerOptions {
  /** The shared webhook signing key. */
  signingKey: string;
  /** Callback invoked with the parsed webhook payload when signature is valid. */
  onEvent: (event: unknown) => void | Promise<void>;
}

/**
 * Creates a Next.js API route handler for ZITADEL webhooks.
 * Placeholder — to be implemented with Next.js Route Handler primitives.
 */
export function createZitadelWebhookHandler(_options: WebhookHandlerOptions) {
  // TODO: implement with NextRequest/NextResponse
  return async function handler() {
    // Placeholder
  };
}
