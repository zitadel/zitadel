"use client";

import { useRouter } from "next/navigation";
import { Button, ButtonVariants } from "./button";
import { Translated } from "./translated";

type Props = {
  organization?: string;
  requestId?: string;
};

// Last block of the sign in page, after every sign in option, for people who have no account yet.
export function RegisterPrompt({ organization, requestId }: Props) {
  const router = useRouter();

  function goToRegister() {
    const registerParams = new URLSearchParams();
    if (organization) {
      registerParams.append("organization", organization);
    }
    if (requestId) {
      registerParams.append("requestId", requestId);
    }

    router.push("/register?" + registerParams);
  }

  return (
    <div className="border-divider-light dark:border-divider-dark mt-4 flex flex-col items-center border-t pt-4">
      <p className="ztdl-p mb-2 text-center" data-testid="first-time-here">
        <Translated i18nKey="firstTimeHere" namespace="loginname" />
      </p>
      <Button
        type="button"
        className="text-primary-light-500 dark:text-primary-dark-500 border-primary-light-500 dark:border-primary-dark-500 w-full justify-center font-semibold"
        variant={ButtonVariants.Secondary}
        onClick={goToRegister}
        data-testid="register-button"
      >
        <Translated i18nKey="register" namespace="loginname" />
      </Button>
    </div>
  );
}
