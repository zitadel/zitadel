import { ColorShade, getColorHash } from "@/helpers/colors";
import { useTheme } from "next-themes";
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
      className={`dark:group-focus:ring-offset-blue dark:text-blue pointer-events-none flex h-[100px] w-[100px] cursor-default items-center justify-center rounded-full bg-primary-light-500 text-primary-light-contrast-500 transition-colors duration-200 hover:bg-primary-light-400 group-focus:outline-none group-focus:ring-2 group-focus:ring-primary-light-200 dark:bg-primary-dark-300 dark:text-primary-dark-contrast-300 hover:dark:bg-primary-dark-500 dark:group-focus:ring-primary-dark-400 ${
        shadow ? "shadow" : ""
      }`}
      style={resolvedTheme === "light" ? avatarStyleLight : avatarStyleDark}
    >
      {imageUrl ? (
        <img
          height={48}
          width={48}
          alt="avatar"
          className="h-full w-full rounded-full border border-divider-light dark:border-divider-dark"
          src={imageUrl}
        />
      ) : (
        <span className={`text-3xl uppercase`}>{credentials}</span>
      )}
    </div>
  );
}
