import { Avatar } from "@/ui/Avatar";
import { ChevronDownIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

type Props = {
  loginName?: string;
  displayName?: string;
  showDropdown: boolean;
  searchParams?: Record<string | number | symbol, string | undefined>;
};

export default function UserAvatar({
  loginName,
  displayName,
  showDropdown,
  searchParams,
}: Props) {
  const params = new URLSearchParams({});

  if (searchParams?.sessionId) {
    params.set("sessionId", searchParams.sessionId);
  }

  if (searchParams?.organization) {
    params.set("organization", searchParams.organization);
  }

  if (searchParams?.authRequestId) {
    params.set("authRequestId", searchParams.authRequestId);
  }

  if (searchParams?.loginName) {
    params.set("loginName", searchParams.loginName);
  }

  return (
    <div className="flex h-full flex-row items-center rounded-full border p-[1px] dark:border-white/20">
      <div>
        <Avatar
          size="small"
          name={displayName ?? loginName ?? ""}
          loginName={loginName ?? ""}
        />
      </div>
      <span className="ml-4 text-14px">{loginName}</span>
      <span className="flex-grow"></span>
      {showDropdown && (
        <Link
          href={"/accounts?" + params}
          className="ml-4 flex items-center justify-center p-1 hover:bg-black/10 dark:hover:bg-white/10 rounded-full mr-1 transition-all"
        >
          <ChevronDownIcon className="h-4 w-4" />
        </Link>
      )}
    </div>
  );
}
