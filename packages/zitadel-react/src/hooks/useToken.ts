import { useZitadel } from "../context.js";

/**
 * Hook to access the current access token.
 */
export function useToken() {
  const { accessToken } = useZitadel();
  return { accessToken };
}
