import * as React from "react";
import { ZitadelServer, settings, getActiveIdentityProvidersResponse } from "@zitadel/server";

export interface SignInWithIDPProps {
  children?: React.ReactNode;
  server: ZitadelServer;
  orgId?: string;
}

function getIDPs(
    server: ZitadelServer,
    orgId?: string,
  ): Promise<GetActiveIdentityProvidersResponse | undefined> {
    const settingsService = settings.getSettings(server);
    return settingsService
      .getActiveIdentityProviders(orgId ? {ctx: {orgId}}: {ctx: {instance: true}}, {})
      .then((resp: getActiveIdentityProvidersResponse) => {
        return resp.settings;
      });


export function SignInWithIDP(props: SignInWithIDPProps) {
  return (
    <div className="ztdl-flex ztdl-flex-row border ztdl-border-divider-light dark:ztdl-border-divider-dark rounded-md px-4 text-sm">
      <div></div>
      {props.children}
    </div>
  );
}

SignInWithIDP.displayName = "SignInWithIDP";
