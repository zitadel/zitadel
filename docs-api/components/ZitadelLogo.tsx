import ZitadelLogoDark from "./ZitadelLogoDark";
import ZitadelLogoLight from "./ZitadelLogoLight";

export function ZitadelLogo() {
  return (
    <>
      <div className="hidden w-fit h-10 dark:flex relative">
        <ZitadelLogoLight />
        <span className="font-semibold absolute -right-6 bottom-0 text-sm text-pink-500 dark:text-pink-500">
          API
        </span>
      </div>
      <div className="flex h-10 dark:hidden relative">
        <ZitadelLogoDark />
        <span className="font-semibold absolute -right-6 bottom-0 text-sm text-pink-500 dark:text-pink-500">
          API
        </span>
      </div>
    </>
  );
}
