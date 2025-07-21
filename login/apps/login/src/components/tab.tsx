"use client";

import type { Item } from "@/components/tab-group";
import { clsx } from "clsx";
import Link from "next/link";
import { useSelectedLayoutSegment } from "next/navigation";

export const Tab = ({
  path,
  item: { slug, text },
}: {
  path: string;
  item: Item;
}) => {
  const segment = useSelectedLayoutSegment();
  const href = slug ? path + "/" + slug : path;
  const isActive =
    // Example home pages e.g. `/layouts`
    (!slug && segment === null) ||
    // Nested pages e.g. `/layouts/electronics`
    segment === slug;

  return (
    <Link
      href={href}
      className={clsx("mr-2 mt-2 rounded-lg px-3 py-1 text-sm font-medium", {
        "bg-gray-700 text-gray-100 hover:bg-gray-500 hover:text-white":
          !isActive,
        "bg-blue-500 text-white": isActive,
      })}
    >
      {text}
    </Link>
  );
};
