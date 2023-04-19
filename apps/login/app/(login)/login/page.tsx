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
      <h1>Welcome back!</h1>
      <p className="ztdl-p">Enter your login data.</p>

      <form className="w-full" onSubmit={() => submit()}>
        <div className="block">
          <TextInput title="loginname" label="Loginname" />
        </div>

        <div>
          <IdentityProviders />
        </div>
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
