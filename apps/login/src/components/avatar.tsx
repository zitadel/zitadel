"use client";

import { ColorShade, getColorHash } from "@/helpers/colors";
import { useTheme } from "next-themes";
import { getComponentRoundness } from "@/lib/theme";

interface AvatarProps {
  name: string | null | undefined;
  loginName: string;
  imageUrl?: string;
  size?: "small" | "base" | "large";
  shadow?: boolean;
}

export function getInitials(name: string, loginName: string) {
  let credentials = "";
  if (name) {
    const split = name.split(" ");
    if (split) {
      const initials = split[0].charAt(0) + (split[1] ? split[1].charAt(0) : "");
      credentials = initials;
    } else {
      credentials = name.charAt(0);
    }
  } else {
    const username = loginName.split("@")[0];
    let separator = "_";
    if (username.includes("-")) {
      separator = "-";
    }
    if (username.includes(".")) {
      separator = ".";
    }
    const split = username.split(separator);
    const initials = split[0].charAt(0) + (split[1] ? split[1].charAt(0) : "");
    credentials = initials;
  }

  return credentials;
}

// Helper function to get avatar roundness from theme
function getAvatarRoundness(): string {
  return getComponentRoundness("avatar");
}

export function Avatar({ size = "base", name, loginName, imageUrl, shadow }: AvatarProps) {
  const { resolvedTheme } = useTheme();
  const credentials = getInitials(name ?? loginName, loginName);
  const avatarRoundness = getAvatarRoundness();

  const color: ColorShade = getColorHash(loginName);

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
      className={`dark:group-focus:ring-offset-blue dark:text-blue pointer-events-none flex h-full w-full flex-shrink-0 cursor-default items-center justify-center bg-primary-light-500 text-primary-light-contrast-500 transition-colors duration-200 hover:bg-primary-light-400 group-focus:outline-none group-focus:ring-2 group-focus:ring-primary-light-200 dark:bg-primary-dark-300 dark:text-primary-dark-contrast-300 hover:dark:bg-primary-dark-500 dark:group-focus:ring-primary-dark-400 ${avatarRoundness} ${
        shadow ? "shadow" : ""
      } ${
        size === "large"
          ? "h-20 w-20 font-normal"
          : size === "base"
            ? "h-[38px] w-[38px] font-bold"
            : size === "small"
              ? "!h-[32px] !w-[32px] text-[13px] font-bold"
              : "h-12 w-12"
      }`}
      style={resolvedTheme === "light" ? avatarStyleLight : avatarStyleDark}
    >
      {imageUrl ? (
        <img
          height={48}
          width={48}
          alt="avatar"
          className={`h-full w-full border border-divider-light dark:border-divider-dark ${avatarRoundness}`}
          src={imageUrl}
        />
      ) : (
        <span className={`uppercase ${size === "large" ? "text-xl" : "text-13px"}`}>{credentials}</span>
      )}
    </div>
  );
}
