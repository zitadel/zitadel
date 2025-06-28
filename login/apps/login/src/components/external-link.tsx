import { ArrowRightIcon } from "@heroicons/react/24/solid";
import { ReactNode } from "react";

export const ExternalLink = ({
  children,
  href,
}: {
  children: ReactNode;
  href: string;
}) => {
  return (
    <a
      href={href}
      className="inline-flex space-x-2 rounded-lg bg-gray-700 px-3 py-1 text-sm font-medium text-gray-100 hover:bg-gray-500 hover:text-white"
    >
      <div>{children}</div>

      <ArrowRightIcon className="block w-4" />
    </a>
  );
};
