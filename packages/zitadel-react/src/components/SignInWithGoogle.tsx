import * as React from "react";

export interface SignInWithGoogleProps {
  children?: React.ReactNode;
}

export function SignInWithGoogle(props: SignInWithGoogleProps) {
  return (
    <div className="ui-flex ui-flex-row ui-items-center ui-bg-white ui-text-black dark:ui-bg-transparent dark:ui-text-white rounded-md p-4 text-sm">
      <img
        className="h-8 w-8"
        src="./public/google.png"
        alt="google"
        height={24}
        width={24}
      />
      <span className="ui-ml-4">Sign in with Google</span>
    </div>
  );
}

SignInWithGoogle.displayName = "SignInWithGoogle";
