import * as React from "react";

export interface SignInWithIDPProps {
  children?: React.ReactNode;
  orgId?: string;
}

export function SignInWithIDP(props: SignInWithIDPProps) {
  return (
    <div className="ztdl-flex ztdl-flex-row border ztdl-border-divider-light dark:ztdl-border-divider-dark rounded-md px-4 text-sm">
      <div></div>
      {props.children}
    </div>
  );
}

SignInWithIDP.displayName = "SignInWithIDP";
