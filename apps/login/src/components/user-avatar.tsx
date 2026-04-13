import { Avatar } from "@/components/avatar";
import { getComponentRoundness } from "@/lib/theme";
import { ChevronDownIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

// Helper function to get user avatar container roundness from theme
function getUserAvatarRoundness(): string {
  return getComponentRoundness("avatarContainer");
}

type Props = {
  loginName?: string;
  displayName?: string;
  showDropdown: boolean;
  searchParams?: Record<string | number | symbol, string | undefined>;
};

export function UserAvatar({ loginName, displayName, showDropdown, searchParams }: Props) {
  const params = new URLSearchParams({});
  const userAvatarRoundness = getUserAvatarRoundness();

  if (searchParams?.sessionId) {
    params.set("sessionId", searchParams.sessionId);
  }

  if (searchParams?.organization) {
    params.set("organization", searchParams.organization);
  }

  if (searchParams?.requestId) {
    params.set("requestId", searchParams.requestId);
  }

  if (searchParams?.loginName) {
    params.set("loginName", searchParams.loginName);
  }

  return (
    <div className={`flex h-full flex-row items-center border p-[1px] dark:border-white/20 ${userAvatarRoundness}`}>
      <div>
        <Avatar size="small" name={displayName ?? loginName ?? ""} loginName={loginName ?? ""} />
      </div>
      <span className="text-14px ml-4 max-w-[250px] overflow-hidden pr-4 text-ellipsis">{loginName}</span>
      <span className="flex-grow"></span>
      {showDropdown && (
        <Link
          href={"/accounts?" + params}
          className={`mr-1 ml-4 flex items-center justify-center p-1 transition-all hover:bg-black/10 dark:hover:bg-white/10 ${userAvatarRoundness}`}
        >
          <ChevronDownIcon className="h-4 w-4" />
        </Link>
      )}
    </div>
  );
}
