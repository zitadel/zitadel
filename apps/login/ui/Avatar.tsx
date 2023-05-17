import { ColorShade, getColorHash } from "#/utils/colors";
import { useTheme } from "next-themes";
import { FC } from "react";

export enum AvatarSize {
  SMALL = "small",
  BASE = "base",
  LARGE = "large",
}

interface AvatarProps {
  name: string | null | undefined;
  loginName: string;
  imageUrl?: string;
  size?: AvatarSize;
  shadow?: boolean;
}

export const Avatar: FC<AvatarProps> = ({
  size = AvatarSize.BASE,
  name,
  loginName,
  imageUrl,
  shadow,
}) => {
  //   const { resolvedTheme } = useTheme();
  let credentials = "";

  console.log(name, loginName);
  if (name) {
    const split = name.split(" ");
    if (split) {
      const initials =
        split[0].charAt(0) + (split[1] ? split[1].charAt(0) : "");
      credentials = initials;
    } else {
      return name.charAt(0);
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

  const color: ColorShade = getColorHash(loginName);

  //   const avatarStyleDark = {
  //     backgroundColor: color[900],
  //     color: color[200],
  //   };

  //   const avatarStyleLight = {
  //     backgroundColor: color[200],
  //     color: color[900],
  //   };

  return (
    <div
      className={`w-full h-full flex-shrink-0 flex justify-center items-center cursor-default pointer-events-none group-focus:outline-none group-focus:ring-2 transition-colors duration-200 dark:group-focus:ring-offset-blue bg-primary-light-500 text-primary-light-contrast-500 hover:bg-primary-light-400 hover:dark:bg-primary-dark-500 group-focus:ring-primary-light-200 dark:group-focus:ring-primary-dark-400 dark:bg-primary-dark-300 dark:text-primary-dark-contrast-300 dark:text-blue rounded-full ${
        shadow ? "shadow" : ""
      } ${
        size === AvatarSize.LARGE
          ? "h-20 w-20 font-normal"
          : size === AvatarSize.BASE
          ? "w-[38px] h-[38px] font-bold"
          : size === AvatarSize.SMALL
          ? "w-[32px] h-[32px] font-bold"
          : ""
      }`}
      //   style={resolvedTheme === "light" ? avatarStyleLight : avatarStyleDark}
    >
      {imageUrl ? (
        <img
          className="border border-divider-light dark:border-divider-dark rounded-full w-12 h-12"
          src={imageUrl}
        />
      ) : (
        <span
          className={`uppercase ${
            size === AvatarSize.LARGE ? "text-xl" : "text-13px"
          }`}
        >
          {credentials}
        </span>
      )}
    </div>
  );
};
