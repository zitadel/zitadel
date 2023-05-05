import { Button, ButtonVariants } from "#/ui/Button";
import { NextPage, NextPageContext } from "next";
import Link from "next/link";
import { useSearchParams } from "next/navigation";

type Props = {
  searchParams: { [key: string]: string | string[] | undefined };
};
export default async function Page({ searchParams }: Props) {
  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Register successful</h1>
      <p className="ztdl-p">You are registered.</p>

      {`userId: ${searchParams["userid"]}`}
      <Link href="/register">
        <Button variant={ButtonVariants.Primary}>go back</Button>
      </Link>
    </div>
  );
}
