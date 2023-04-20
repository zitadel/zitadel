import { Avatar, AvatarSize } from "#/ui/Avatar";

type Props = {
  name: string;
};

export default function UserAvatar({ name }: Props) {
  return (
    <div className="flex h-full w-full flex-row items-center rounded-full border p-[1px] dark:border-white/20">
      <div>
        <Avatar size={AvatarSize.SMALL} name={name} loginName={name} />
      </div>
      <span className="ml-4 text-14px">{name}</span>
    </div>
  );
}
