import { useZitadel } from "../context.js";

/**
 * Hook to access the current user information.
 */
export function useUser() {
  const { user } = useZitadel();
  return { user };
}
