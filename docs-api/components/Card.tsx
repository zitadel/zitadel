import * as React from "react";
import LanguageSwitcher from "./LanguageSwitcher";
import ProtocolSwitcher from "./ProtocolSwitcher";

type Props = {
  title: string;
  hasLanguageToggle: boolean;
  hasProtocolToggle: boolean;
  language?: string;
  protocol?: string;
  children: React.ReactNode;
};

export function Card({
  title,
  hasLanguageToggle,
  hasProtocolToggle,
  children,
  language,
  protocol,
}: Props) {
  console.log(protocol);
  return (
    <div className="my-4 bg-white dark:bg-background-dark-400 border border-border-light dark:border-border-dark rounded-md w-full">
      <div className="py-2 px-4 bg-black/10 text-sm dark:bg-white/10 flex flex-row items-center justify-between rounded-t-md">
        {title}
        {hasLanguageToggle && <LanguageSwitcher initial={language} />}
        {hasProtocolToggle && <ProtocolSwitcher initial={protocol} />}
      </div>
      <div className="px-4">{children}</div>
    </div>
  );
}
