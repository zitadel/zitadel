import Image from "next/image";
type Props = {
  height?: number;
  width?: number;
};

export function ZitadelLogo({ height = 40, width = 147.5 }: Props) {
  return (
    <>
      <div className="hidden dark:flex">
        {/* <ZitadelLogoLight /> */}

        <Image
          height={height}
          width={width}
          src="/zitadel-logo-light.svg"
          alt="zitadel logo"
          priority={true}
        />
      </div>
      <div className="flex dark:hidden">
        <Image
          height={height}
          width={width}
          priority={true}
          src="/zitadel-logo-dark.svg"
          alt="zitadel logo"
        />
      </div>
    </>
  );
}
