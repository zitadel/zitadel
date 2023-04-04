import * as React from "react";

export interface SignInWithGoogleProps {
  children?: React.ReactNode;
}

export function SignInWithGoogle(props: SignInWithGoogleProps) {
  return (
    <div className="flex flex-row items-center bg-white text-gray-500 dark:bg-transparent dark:text-white rounded-md p-4 text-sm">
      <img
        className="h-8 w-8"
        src="idp/google.png"
        alt="google"
        height={24}
        width={24}
      />
      Sign in with Google
    </div>
  );
}

SignInWithGoogle.displayName = "SignInWithGoogle";
