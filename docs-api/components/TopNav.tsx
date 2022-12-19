import React from "react";
import Link from "next/link";

export function TopNav() {
  return (
    <nav className="text-sm flex flex-row justify-end py-6">
      <Link href="https://zitadel.com/docs" className="flex flex-row px-4">
        Docs
      </Link>
      <Link href="https://zitadel.com/contact" className="flex flex-row px-4">
        Support
      </Link>
      <Link href="https://zitadel.cloud" className="flex flex-row px-4">
        Customer portal
      </Link>
    </nav>
  );
}
