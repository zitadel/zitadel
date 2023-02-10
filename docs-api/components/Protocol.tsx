import { useRouter } from "next/router";
import * as React from "react";

type Props = {
  showDefault?: boolean;
  showOnProtocol?: string;
  children: React.ReactNode;
};

export function Protocol({ showOnProtocol, showDefault, children }: Props) {
  const router = useRouter();

  const { protocol } = router.query;

  return (
    (protocol ? showOnProtocol === protocol : !!showDefault) && (
      <div className="my-4 bg-white dark:bg-background-dark-400 border border-border-light dark:border-border-dark rounded-md w-full">
        <div className="px-4">{children}</div>
      </div>
    )
  );
}
