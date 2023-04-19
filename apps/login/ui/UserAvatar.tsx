type Props = {
  name: string;
};

export default function UserAvatar({ name }: Props) {
  return (
    <div className="flex w-full flex-row items-center rounded-full border p-[1px] dark:border-white/20">
      {/* <Image
          height={20}
          width={20}
          className="avatar-img"
          src=""
          alt="user-avatar"
        /> */}
      <div className="h-8 w-8 rounded-full bg-primary-light-700 dark:bg-primary-dark-800"></div>
      <span className="ml-4 text-14px">{name}</span>
    </div>
  );
}
