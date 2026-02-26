import { createContext, useContext, type ReactNode } from "react";

export interface ZitadelContextValue {
  /** Whether the user is currently authenticated. */
  isAuthenticated: boolean;
  /** The current access token, if available. */
  accessToken: string | undefined;
  /** The current user information, if available. */
  user: Record<string, unknown> | undefined;
}

const ZitadelContext = createContext<ZitadelContextValue | undefined>(undefined);

export interface ZitadelProviderProps {
  children: ReactNode;
  value: ZitadelContextValue;
}

/**
 * Provides ZITADEL context to child components.
 */
export function ZitadelProvider({ children, value }: ZitadelProviderProps) {
  return <ZitadelContext.Provider value={value}>{children}</ZitadelContext.Provider>;
}

/**
 * Hook to access the ZITADEL context.
 * Must be used within a ZitadelProvider.
 */
export function useZitadel(): ZitadelContextValue {
  const context = useContext(ZitadelContext);
  if (context === undefined) {
    throw new Error("useZitadel must be used within a ZitadelProvider");
  }
  return context;
}
