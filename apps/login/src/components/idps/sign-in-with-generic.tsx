"use client";

import { ReactNode, forwardRef } from "react";
import { IdpButtonClasses, SignInWithIdentityProviderProps } from "./classes";

export const SignInWithGeneric = forwardRef<
  HTMLButtonElement,
  SignInWithIdentityProviderProps
>(
  (
    { children, className = "h-[50px] pl-20", name = "", ...props },
    ref,
  ): ReactNode => (
    <button
      type="button"
      ref={ref}
      className={`${IdpButtonClasses} ${className}`}
      {...props}
    >
      {children ? children : <span className="">{name}</span>}
    </button>
  ),
);

SignInWithGeneric.displayName = "SignInWithGeneric";
