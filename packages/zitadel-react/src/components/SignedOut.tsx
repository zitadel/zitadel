import type { ReactNode } from "react";
import { useZitadel } from "../context.js";

/**
 * Renders children only when the user is NOT authenticated.
 */
export function SignedOut({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useZitadel();
  return isAuthenticated ? null : <>{children}</>;
}
