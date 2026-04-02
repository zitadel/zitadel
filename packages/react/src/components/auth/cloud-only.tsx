"use client";

import { useDeployment } from "../../context/deployment";

/**
 * Wrapper component that only renders children when deploying in cloud mode.
 * Used to hide features like instance management, billing, team, and support
 * that are only available in ZITADEL Cloud.
 *
 * Usage:
 *   <CloudOnly>
 *     <BillingSection />
 *   </CloudOnly>
 */
export function CloudOnly({ children }: { children: React.ReactNode }) {
  const { isCloud } = useDeployment();

  if (!isCloud) {
    return null;
  }

  return <>{children}</>;
}

/**
 * Wrapper component that only renders children when deploying in self-hosted mode.
 */
export function SelfHostedOnly({ children }: { children: React.ReactNode }) {
  const { isSelfHosted } = useDeployment();

  if (!isSelfHosted) {
    return null;
  }

  return <>{children}</>;
}
