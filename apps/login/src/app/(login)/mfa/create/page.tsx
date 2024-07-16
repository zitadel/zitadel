"use client";
import { Button, ButtonVariants } from "@/ui/Button";
import { TextInput } from "@/ui/Input";
import UserAvatar from "@/ui/UserAvatar";
import { useRouter } from "next/navigation";

export default function Page() {
  const router = useRouter();

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Password</h1>
      <p className="ztdl-p mb-6 block">Enter your password.</p>

      <UserAvatar
        showDropdown
        displayName="Max Peintner"
        loginName="max@zitadel.com"
      ></UserAvatar>

      <div className="w-full">
        <TextInput type="password" label="Password" />
      </div>
      <div className="flex w-full flex-row items-center justify-between">
        <Button
          onClick={() => router.back()}
          variant={ButtonVariants.Secondary}
        >
          back
        </Button>
        <Button variant={ButtonVariants.Primary}>continue</Button>
      </div>
    </div>
  );
}
