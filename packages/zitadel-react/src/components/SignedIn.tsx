import type { ReactNode } from "react";
import { useZitadel } from "../context.js";

/**
 * Renders children only when the user is authenticated.
 */
export function SignedIn({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useZitadel();
  return isAuthenticated ? <>{children}</> : null;
}
