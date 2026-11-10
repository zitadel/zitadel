type Props = {
  height?: number;
  width?: number;
};

export function ZitadelLogo({ height = 40, width = 147.5 }: Props) {
  return (
    <>
      <div className="hidden dark:flex">
        {/* <ZitadelLogoLight /> */}

        <img height={height} width={width} src="/zitadel-logo-light.svg" alt="zitadel logo" />
      </div>
      <div className="flex dark:hidden">
        <img height={height} width={width} src="/zitadel-logo-dark.svg" alt="zitadel logo" />
      </div>
    </>
  );
}
