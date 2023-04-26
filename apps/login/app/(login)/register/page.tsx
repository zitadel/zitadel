"use client";

import { Button, ButtonVariants } from "#/ui/Button";
import IdentityProviders from "#/ui/IdentityProviders";
import { TextInput } from "#/ui/Input";
import { useRouter } from "next/navigation";

export default function Page() {
  const router = useRouter();

  function submit() {
    router.push("/password");
  }
  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Register</h1>
      <p className="ztdl-p">Create your ZITADEL account.</p>

      <form className="" onSubmit={() => submit()}>
        <div className="grid grid-cols-2 gap-4">
          <div className="">
            <TextInput label="Firstname" />
          </div>
          <div className="">
            <TextInput label="Lastname" />
          </div>
          <div className="">
            <TextInput label="Email" />
          </div>
          <div className="">
            <TextInput label="Password" />
          </div>
          <div className="">
            <TextInput label="Password Confirmation" />
          </div>
        </div>

        <PrivacyPolicyCheckboxes />

        <div className="mt-8 flex w-full flex-row items-center justify-between">
          <Button type="button" variant={ButtonVariants.Secondary}>
            back
          </Button>
          <Button
            type="submit"
            variant={ButtonVariants.Primary}
            onClick={() => submit()}
          >
            continue
          </Button>
        </div>
      </form>
    </div>
  );
}
