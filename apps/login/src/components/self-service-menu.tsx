import Link from "next/link";

export function SelfServiceMenu({ sessionId }: { sessionId: string }) {
  const list: any[] = [];

  // if (!!config.selfservice.change_password.enabled) {
  //   list.push({
  //     link:
  //       `/me/change-password?` +
  //       new URLSearchParams({
  //         sessionId: sessionId,
  //       }),
  //     name: "Change password",
  //   });
  // }

  return (
    <div className="w-full flex flex-col space-y-2">
      {list.map((menuitem, index) => {
        return (
          <SelfServiceItem
            link={menuitem.link}
            key={"self-service-" + index}
            name={menuitem.name}
          />
        );
      })}
    </div>
  );
}

const SelfServiceItem = ({ name, link }: { name: string; link: string }) => {
  return (
    <Link
      prefetch={false}
      href={link}
      className="w-full group flex flex-row items-center bg-background-light-400 dark:bg-background-dark-400  border border-divider-light hover:shadow-lg dark:hover:bg-white/10 py-2 px-4 rounded-md transition-all"
    >
      {name}
    </Link>
  );
};
