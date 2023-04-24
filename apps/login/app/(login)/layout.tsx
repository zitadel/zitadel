import { ZitadelLogo } from "#/ui/ZitadelLogo";
import React from "react";

export default async function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="mx-auto flex max-w-[400px] flex-col items-center space-y-4 py-10">
      <div className="relative">
        <ZitadelLogo height={70} width={180} />
      </div>

      <div className="w-full">{children}</div>
      <div className="flex flex-row justify-between"></div>
    </div>
  );
}
