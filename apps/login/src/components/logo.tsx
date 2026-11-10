type Props = {
  darkSrc?: string;
  lightSrc?: string;
  height?: number;
  width?: number;
};

export function Logo({ lightSrc, darkSrc, height = 40, width = 147.5 }: Props) {
  return (
    <>
      {darkSrc && (
        <div className="hidden dark:flex">
          <img height={height} width={width} src={darkSrc} alt="logo" />
        </div>
      )}
      {lightSrc && (
        <div className="flex dark:hidden">
          <img height={height} width={width} src={lightSrc} alt="logo" />
        </div>
      )}
    </>
  );
}
