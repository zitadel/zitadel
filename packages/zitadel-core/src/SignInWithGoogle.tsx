import * as React from "react";

export interface SignInWithGoogleProps {
  children: React.ReactNode;
}

export function SignInWithGoogle(props: SignInWithGoogleProps) {
  return <button>{props.children}</button>;
}

SignInWithGoogle.displayName = "SignInWithGoogle";
