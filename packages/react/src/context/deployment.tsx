"use client";

import React, { createContext, useContext } from "react";

type DeploymentMode = "self-hosted" | "cloud";

interface DeploymentContextType {
  /** Current deployment mode */
  mode: DeploymentMode;
  /** Whether this is a ZITADEL Cloud deployment */
  isCloud: boolean;
  /** Whether this is a self-hosted deployment */
  isSelfHosted: boolean;
}

const DeploymentContext = createContext<DeploymentContextType>({
  mode: "self-hosted",
  isCloud: false,
  isSelfHosted: true,
});

/**
 * DeploymentProvider reads NEXT_PUBLIC_DEPLOYMENT_MODE from env to
 * determine whether the console is running in cloud or self-hosted mode.
 *
 * Cloud mode enables: instance management, billing, team, support features.
 * Self-hosted mode hides those and focuses on single-instance management.
 */
export function DeploymentProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const envMode = process.env.NEXT_PUBLIC_DEPLOYMENT_MODE as
    | DeploymentMode
    | undefined;
  const mode: DeploymentMode = envMode === "cloud" ? "cloud" : "self-hosted";

  return (
    <DeploymentContext.Provider
      value={{
        mode,
        isCloud: mode === "cloud",
        isSelfHosted: mode === "self-hosted",
      }}
    >
      {children}
    </DeploymentContext.Provider>
  );
}

/**
 * Hook to check the current deployment mode.
 *
 * Usage:
 *   const { isCloud, isSelfHosted } = useDeployment();
 *   if (isCloud) { showBillingNav(); }
 */
export function useDeployment() {
  return useContext(DeploymentContext);
}
