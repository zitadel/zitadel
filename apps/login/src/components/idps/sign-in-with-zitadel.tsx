"use client";

import { forwardRef } from "react";
import { Translated } from "../translated";
import { BaseButton, SignInWithIdentityProviderProps } from "./base-button";

export const SignInWithZitadel = forwardRef<HTMLButtonElement, SignInWithIdentityProviderProps>(
  function SignInWithZitadel(props, ref) {
    const { children, name, ...restProps } = props;

    return (
      <BaseButton {...restProps} ref={ref}>
        <div className="flex h-12 w-12 items-center justify-center">
          <div className="hidden h-6 w-6 dark:flex">
            <img className="h-full w-full object-contain" src="/logo/zitadel-logo-solo-darkdesign.svg" alt="" />
          </div>
          <div className="flex h-6 w-6 dark:hidden">
            <img className="h-full w-full object-contain" src="/logo/zitadel-logo-solo-lightdesign.svg" alt="" />
          </div>
        </div>
        {children ? (
          children
        ) : (
          <span className="ml-4">{name ? name : <Translated i18nKey="signInWithZitadel" namespace="idp" />}</span>
        )}
      </BaseButton>
    );
  },
);
