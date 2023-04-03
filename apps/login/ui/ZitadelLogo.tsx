import { ZitadelLogoDark } from './ZitadelLogoDark';
import { ZitadelLogoLight } from './ZitadelLogoLight';

export function ZitadelLogo() {
  return (
    <>
      <div className="hidden dark:flex">
        <ZitadelLogoLight />
      </div>
      <div className="flex dark:hidden">
        <ZitadelLogoDark />
      </div>
    </>
  );
}
