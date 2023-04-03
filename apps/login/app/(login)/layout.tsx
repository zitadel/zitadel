import { ColorService } from '#/utils/colors';
import Image from 'next/image';
import React, { useEffect } from 'react';

export default async function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const colorService = new ColorService();

  return (
    <div className="mx-auto flex max-w-[400px] flex-col items-center space-y-4 py-10">
      <div className="relative h-28 w-48">
        <Image
          fill
          priority
          sizes="100%"
          src="https://zitadel.com/zitadel-logo-light.svg"
          alt="Login logo"
        />
      </div>

      <div className="w-full">{children}</div>
      <div className="flex flex-row justify-between"></div>
    </div>
  );
}
