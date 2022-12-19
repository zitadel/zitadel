import Link from "next/link";
import * as React from "react";

type Props = {
  method: string;
  link: string;
  children: React.ReactNode;
};

export function Endpoint({ method, link, children }: Props) {
  return (
    <Link href={link} className="block text-black/80 dark:text-white/80">
      <span
        className={`mr-4 ${
          method === "POST"
            ? "text-amber-500 dark:text-amber-500"
            : method === "GET"
            ? "text-green-500"
            : "text-primary-light-500 dark:text-primary-dark-500"
        }`}
      >
        {method}
      </span>
      {children}
    </Link>
  );
}
