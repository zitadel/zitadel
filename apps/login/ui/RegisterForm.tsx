"use client";

import { PasswordComplexityPolicy, PrivacyPolicy } from "@zitadel/server";
import PasswordComplexity from "./PasswordComplexity";
import { useState } from "react";
import { Button, ButtonVariants } from "./Button";
import { TextInput } from "./Input";
import { PrivacyPolicyCheckboxes } from "./PrivacyPolicyCheckboxes";

type Props = {
  privacyPolicy: PrivacyPolicy;
  passwordComplexityPolicy: PasswordComplexityPolicy;
};

export default function RegisterForm({
  privacyPolicy,
  passwordComplexityPolicy,
}: Props) {
  const [tosAndPolicyAccepted, setTosAndPolicyAccepted] = useState(false);

  return (
    <form className="w-full">
      <div className="grid grid-cols-2 gap-4 mb-4">
        <div className="">
          <TextInput label="Firstname" />
        </div>
        <div className="">
          <TextInput label="Lastname" />
        </div>
        <div className="col-span-2">
          <TextInput label="Email" />
        </div>
        <div className="">
          <TextInput label="Password" />
        </div>
        <div className="">
          <TextInput label="Password Confirmation" />
        </div>
      </div>

      {passwordComplexityPolicy && (
        <PasswordComplexity
          passwordComplexityPolicy={passwordComplexityPolicy}
          password={""}
          equals={false}
        />
      )}

      {privacyPolicy && (
        <PrivacyPolicyCheckboxes
          privacyPolicy={privacyPolicy}
          onChange={setTosAndPolicyAccepted}
        />
      )}

      <div className="mt-8 flex w-full flex-row items-center justify-between">
        <Button type="button" variant={ButtonVariants.Secondary}>
          back
        </Button>
        <Button type="submit" variant={ButtonVariants.Primary}>
          continue
        </Button>
      </div>
    </form>
  );
}
