import { ColorShade, getColorHash } from "@/helpers/colors";
import { getComponentRoundness } from "@/lib/theme";
import { CSSProperties } from "react";

interface AvatarProps {
  name: string | null | undefined;
  loginName: string;
  imageUrl?: string;
  size?: "small" | "base" | "large";
  shadow?: boolean;
}

export function getInitials(name: string, loginName: string) {
  if (name) {
    const split = name.split(" ");
    return split[0].charAt(0) + (split[1] ? split[1].charAt(0) : "");
  }

  const username = loginName.split("@")[0];
  let separator = "_";
  if (username.includes("-")) {
    separator = "-";
  }
  if (username.includes(".")) {
    separator = ".";
  }
  const split = username.split(separator);
  return split[0].charAt(0) + (split[1] ? split[1].charAt(0) : "");
}

// Helper function to get avatar roundness from theme
function getAvatarRoundness(): string {
  return getComponentRoundness("avatar");
}

export function Avatar({ size = "base", name, loginName, imageUrl, shadow }: AvatarProps) {
  const credentials = getInitials(name ?? loginName, loginName);
  const avatarRoundness = getAvatarRoundness();

  const color: ColorShade = getColorHash(loginName);

  // Both palettes are exposed as theme-independent CSS variables so the server
  // and the client render identical markup. The `dark:` variants switch
  // between them once next-themes resolves the theme on the client, which
  // avoids a hydration mismatch when the system theme is light - an inline
  // style based on `resolvedTheme` would render the dark branch on the server
  // and the light branch during hydration.
  const avatarVars = {
    "--avatar-bg": color[200],
    "--avatar-fg": color[900],
    "--avatar-bg-dark": color[900],
    "--avatar-fg-dark": color[200],
  } as CSSProperties;

  return (
    <div
      className={`dark:group-focus:ring-offset-blue dark:text-blue bg-[var(--avatar-bg)] text-[var(--avatar-fg)] hover:bg-primary-light-400 group-focus:ring-primary-light-200 dark:bg-[var(--avatar-bg-dark)] dark:text-[var(--avatar-fg-dark)] hover:dark:bg-primary-dark-500 dark:group-focus:ring-primary-dark-400 pointer-events-none flex h-full w-full flex-shrink-0 cursor-default items-center justify-center transition-colors duration-200 group-focus:ring-2 group-focus:outline-none ${avatarRoundness} ${
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
      style={avatarVars}
    >
      {imageUrl ? (
        <img
          height={48}
          width={48}
          alt="avatar"
          className={`border-divider-light dark:border-divider-dark h-full w-full border ${avatarRoundness}`}
          src={imageUrl}
        />
      ) : (
        <span className={`uppercase ${size === "large" ? "text-xl" : "text-13px"}`}>{credentials}</span>
      )}
    </div>
  );
}
