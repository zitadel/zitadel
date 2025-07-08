"use client";

import { ColorShade, getColorHash } from "@/helpers/colors";
import { useTheme } from "next-themes";
import Image from "next/image";

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
      const initials =
        split[0].charAt(0) + (split[1] ? split[1].charAt(0) : "");
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

export function Avatar({
  size = "base",
  name,
  loginName,
  imageUrl,
  shadow,
}: AvatarProps) {
  const { resolvedTheme } = useTheme();
  const credentials = getInitials(name ?? loginName, loginName);

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
      className={`w-full h-full flex-shrink-0 flex justify-center items-center cursor-default pointer-events-none group-focus:outline-none group-focus:ring-2 transition-colors duration-200 dark:group-focus:ring-offset-blue bg-primary-light-500 text-primary-light-contrast-500 hover:bg-primary-light-400 hover:dark:bg-primary-dark-500 group-focus:ring-primary-light-200 dark:group-focus:ring-primary-dark-400 dark:bg-primary-dark-300 dark:text-primary-dark-contrast-300 dark:text-blue rounded-full ${
        shadow ? "shadow" : ""
      } ${
        size === "large"
          ? "h-20 w-20 font-normal"
          : size === "base"
            ? "w-[38px] h-[38px] font-bold"
            : size === "small"
              ? "!w-[32px] !h-[32px] font-bold text-[13px]"
              : "w-12 h-12"
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
        <span
          className={`uppercase ${size === "large" ? "text-xl" : "text-13px"}`}
        >
          {credentials}
        </span>
      )}
    </div>
  );
}
