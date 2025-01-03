import Image from "next/image";

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
          <Image
            height={height}
            width={width}
            src={darkSrc}
            alt="logo"
            priority={true}
          />
        </div>
      )}
      {lightSrc && (
        <div className="flex dark:hidden">
          <Image
            height={height}
            width={width}
            priority={true}
            src={lightSrc}
            alt="logo"
          />
        </div>
      )}
    </>
  );
}
