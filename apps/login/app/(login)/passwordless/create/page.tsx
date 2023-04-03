'use client';
import { Button, ButtonVariants } from '#/ui/Button';
import { TextInput } from '#/ui/Input';
import { useRouter } from 'next/navigation';

export default function Page() {
  const router = useRouter();

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Password</h1>
      <p className="ztdl-p mb-6 block">Enter your password.</p>

      <div className="flex w-full flex-row items-center rounded-full border p-[1px] dark:border-white/20">
        {/* <Image
          height={20}
          width={20}
          className="avatar-img"
          src=""
          alt="user-avatar"
        /> */}
        <div className="h-8 w-8 rounded-full bg-primary-dark-800"></div>
        <span className="ml-4 text-14px">max@zitadel.cloud</span>
      </div>

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
