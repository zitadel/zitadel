export interface MiddlewareOptions {
  /** Paths that require authentication. */
  protectedPaths?: string[];
  /** The sign-in URL to redirect unauthenticated users to. */
  signInUrl?: string;
}

/**
 * Creates a Next.js middleware function that protects routes.
 * Placeholder — to be implemented with Next.js middleware primitives.
 */
export function createZitadelMiddleware(_options?: MiddlewareOptions) {
  // TODO: implement with NextResponse and NextRequest
  return async function middleware() {
    // Placeholder
  };
}
