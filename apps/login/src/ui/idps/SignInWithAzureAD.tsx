"use client";

import { ReactNode, forwardRef } from "react";
import { IdpButtonClasses, SignInWithIdentityProviderProps } from "./classes";

export const SignInWithAzureAD = forwardRef<
  HTMLButtonElement,
  SignInWithIdentityProviderProps
>(
  ({ children, className = "", name = "", ...props }, ref): ReactNode => (
    <button
      type="button"
      ref={ref}
      className={`${IdpButtonClasses} ${className}`}
      {...props}
    >
      <div className="h-12 p-[10px] w-12 flex items-center justify-center">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="21"
          height="21"
          viewBox="0 0 21 21"
          className="w-full h-full"
        >
          <path fill="#f25022" d="M1 1H10V10H1z"></path>
          <path fill="#00a4ef" d="M1 11H10V20H1z"></path>
          <path fill="#7fba00" d="M11 1H20V10H11z"></path>
          <path fill="#ffb900" d="M11 11H20V20H11z"></path>
        </svg>
      </div>
      {children ? (
        children
      ) : (
        <span className="ml-4">{name ? name : "Sign in with AzureAD"}</span>
      )}
    </button>
  ),
);

SignInWithAzureAD.displayName = "SignInWithAzureAD";
