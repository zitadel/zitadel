"use client";

import { forwardRef } from "react";
import { BaseButton, SignInWithIdentityProviderProps } from "./base-button";

export const SignInWithGeneric = forwardRef<
  HTMLButtonElement,
  SignInWithIdentityProviderProps
>(function SignInWithGeneric(props, ref) {
  const {
    children,
    name = "",
    className = "h-[50px] pl-20",
    ...restProps
  } = props;
  return (
    <BaseButton {...restProps} ref={ref} className={className}>
      {children ? children : <span>{name}</span>}
    </BaseButton>
  );
});
