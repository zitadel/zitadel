import { ColorShade, getColorHash } from "@/helpers/colors";
import { useTheme } from "next-themes";
import Image from "next/image";
import { getInitials } from "./avatar";

interface AvatarProps {
  appName: string;
  imageUrl?: string;
  shadow?: boolean;
}

export function AppAvatar({ appName, imageUrl, shadow }: AvatarProps) {
  const { resolvedTheme } = useTheme();
  const credentials = getInitials(appName, appName);

  const color: ColorShade = getColorHash(appName);

  const avatarStyleDark = {
    backgroundColor: color[900],
    color: color[200],
  };

  const avatarStyleLight = {
    backgroundColor: color[200],
    color: color[900],
  };

  return (
    <div
      className={`w-[100px] h-[100px] flex justify-center items-center cursor-default pointer-events-none group-focus:outline-none group-focus:ring-2 transition-colors duration-200 dark:group-focus:ring-offset-blue bg-primary-light-500 text-primary-light-contrast-500 hover:bg-primary-light-400 hover:dark:bg-primary-dark-500 group-focus:ring-primary-light-200 dark:group-focus:ring-primary-dark-400 dark:bg-primary-dark-300 dark:text-primary-dark-contrast-300 dark:text-blue rounded-full ${
        shadow ? "shadow" : ""
      }`}
      style={resolvedTheme === "light" ? avatarStyleLight : avatarStyleDark}
    >
      {imageUrl ? (
        <Image
          height={48}
          width={48}
          alt="avatar"
          className="w-full h-full border border-divider-light dark:border-divider-dark rounded-full"
          src={imageUrl}
        />
      ) : (
        <span className={`uppercase text-3xl`}>{credentials}</span>
      )}
    </div>
  );
}
