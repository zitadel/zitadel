import { useZitadel } from "../context.js";

/**
 * Hook to access the current session state.
 */
export function useSession() {
  const { isAuthenticated, accessToken } = useZitadel();
  return { isAuthenticated, accessToken };
}
