import { CheckIcon } from "@heroicons/react/24/solid";
import { clsx } from "clsx";
import Link from "next/link";
import { ReactNode } from "react";
import { BadgeState, StateBadge } from "./state-badge";

const cardClasses = (alreadyAdded: boolean) =>
  clsx(
    "relative bg-background-light-400 dark:bg-background-dark-400 group block space-y-1.5 rounded-md px-5 py-3  border border-divider-light dark:border-divider-dark transition-all ",
    alreadyAdded
      ? "opacity-50 cursor-default"
      : "hover:shadow-lg hover:dark:bg-white/10",
  );

const LinkWrapper = ({
  alreadyAdded,
  children,
  link,
}: {
  alreadyAdded: boolean;
  children: ReactNode;
  link: string;
}) => {
  return !alreadyAdded ? (
    <Link href={link} className={cardClasses(alreadyAdded)}>
      {children}
    </Link>
  ) : (
    <div className={cardClasses(alreadyAdded)}>{children}</div>
  );
};

export const TOTP = (alreadyAdded: boolean, link: string) => {
  return (
    <LinkWrapper key={link} alreadyAdded={alreadyAdded} link={link}>
      <div
        className={clsx(
          "flex items-center font-medium",
          alreadyAdded ? "opacity-50" : "",
        )}
      >
        <svg
          className="mr-4 h-8 w-8 -translate-x-[2px] transform fill-current text-black dark:text-white"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
        >
          <title>timer-lock-outline</title>
          <path d="M11 8H13V14H11V8M13 19.92C12.67 19.97 12.34 20 12 20C8.13 20 5 16.87 5 13S8.13 6 12 6C14.82 6 17.24 7.67 18.35 10.06C18.56 10.04 18.78 10 19 10C19.55 10 20.07 10.11 20.57 10.28C20.23 9.22 19.71 8.24 19.03 7.39L20.45 5.97C20 5.46 19.55 5 19.04 4.56L17.62 6C16.07 4.74 14.12 4 12 4C7.03 4 3 8.03 3 13S7.03 22 12 22C12.42 22 12.83 21.96 13.24 21.91C13.09 21.53 13 21.12 13 20.7V19.92M15 1H9V3H15V1M23 17.3V20.8C23 21.4 22.4 22 21.7 22H16.2C15.6 22 15 21.4 15 20.7V17.2C15 16.6 15.6 16 16.2 16V14.5C16.2 13.1 17.6 12 19 12S21.8 13.1 21.8 14.5V16C22.4 16 23 16.6 23 17.3M20.5 14.5C20.5 13.7 19.8 13.2 19 13.2S17.5 13.7 17.5 14.5V16H20.5V14.5Z" />
        </svg>{" "}
        <span>Authenticator App</span>
      </div>
      {alreadyAdded && (
        <>
          <Setup />
        </>
      )}
    </LinkWrapper>
  );
};

export const U2F = (alreadyAdded: boolean, link: string) => {
  return (
    <LinkWrapper key={link} alreadyAdded={alreadyAdded} link={link}>
      <div
        className={clsx(
          "flex items-center font-medium",
          alreadyAdded ? "" : "",
        )}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
          stroke="currentColor"
          className="mr-4 h-8 w-8"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M7.864 4.243A7.5 7.5 0 0119.5 10.5c0 2.92-.556 5.709-1.568 8.268M5.742 6.364A7.465 7.465 0 004.5 10.5a7.464 7.464 0 01-1.15 3.993m1.989 3.559A11.209 11.209 0 008.25 10.5a3.75 3.75 0 117.5 0c0 .527-.021 1.049-.064 1.565M12 10.5a14.94 14.94 0 01-3.6 9.75m6.633-4.596a18.666 18.666 0 01-2.485 5.33"
          />
        </svg>
        <span>Universal Second Factor</span>
      </div>
      {alreadyAdded && (
        <>
          <Setup />
        </>
      )}
    </LinkWrapper>
  );
};

export const EMAIL = (alreadyAdded: boolean, link: string) => {
  return (
    <LinkWrapper key={link} alreadyAdded={alreadyAdded} link={link}>
      <div
        className={clsx(
          "flex items-center font-medium",
          alreadyAdded ? "" : "",
        )}
      >
        <svg
          className="mr-4 h-8 w-8"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={1.5}
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75"
          />
        </svg>

        <span>Code via Email</span>
      </div>
      {alreadyAdded && (
        <>
          <Setup />
        </>
      )}
    </LinkWrapper>
  );
};

export const SMS = (alreadyAdded: boolean, link: string) => {
  return (
    <LinkWrapper key={link} alreadyAdded={alreadyAdded} link={link}>
      <div
        className={clsx(
          "flex items-center font-medium",
          alreadyAdded ? "" : "",
        )}
      >
        <svg
          className="mr-4 h-8 w-8"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M10.5 1.5H8.25A2.25 2.25 0 006 3.75v16.5a2.25 2.25 0 002.25 2.25h7.5A2.25 2.25 0 0018 20.25V3.75a2.25 2.25 0 00-2.25-2.25H13.5m-3 0V3h3V1.5m-3 0h3m-3 18.75h3"
          />
        </svg>
        <span>Code via SMS</span>
      </div>
      {alreadyAdded && (
        <>
          <Setup />
        </>
      )}
    </LinkWrapper>
  );
};

export const PASSKEYS = (alreadyAdded: boolean, link: string) => {
  return (
    <LinkWrapper key={link} alreadyAdded={alreadyAdded} link={link}>
      <div
        className={clsx(
          "flex items-center font-medium",
          alreadyAdded ? "" : "",
        )}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
          stroke="currentColor"
          className="mr-4 h-8 w-8"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M7.864 4.243A7.5 7.5 0 0119.5 10.5c0 2.92-.556 5.709-1.568 8.268M5.742 6.364A7.465 7.465 0 004.5 10.5a7.464 7.464 0 01-1.15 3.993m1.989 3.559A11.209 11.209 0 008.25 10.5a3.75 3.75 0 117.5 0c0 .527-.021 1.049-.064 1.565M12 10.5a14.94 14.94 0 01-3.6 9.75m6.633-4.596a18.666 18.666 0 01-2.485 5.33"
          />
        </svg>
        <span>Passkeys</span>
      </div>
      {alreadyAdded && (
        <>
          <Setup />
        </>
      )}
    </LinkWrapper>
  );
};

export const PASSWORD = (alreadyAdded: boolean, link: string) => {
  return (
    <LinkWrapper key={link} alreadyAdded={alreadyAdded} link={link}>
      <div
        className={clsx(
          "flex items-center font-medium",
          alreadyAdded ? "" : "",
        )}
      >
        <svg
          className="mr-4 h-7 w-8 fill-current"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
        >
          <title>form-textbox-password</title>
          <path d="M17,7H22V17H17V19A1,1 0 0,0 18,20H20V22H17.5C16.95,22 16,21.55 16,21C16,21.55 15.05,22 14.5,22H12V20H14A1,1 0 0,0 15,19V5A1,1 0 0,0 14,4H12V2H14.5C15.05,2 16,2.45 16,3C16,2.45 16.95,2 17.5,2H20V4H18A1,1 0 0,0 17,5V7M2,7H13V9H4V15H13V17H2V7M20,15V9H17V15H20M8.5,12A1.5,1.5 0 0,0 7,10.5A1.5,1.5 0 0,0 5.5,12A1.5,1.5 0 0,0 7,13.5A1.5,1.5 0 0,0 8.5,12M13,10.89C12.39,10.33 11.44,10.38 10.88,11C10.32,11.6 10.37,12.55 11,13.11C11.55,13.63 12.43,13.63 13,13.11V10.89Z" />
        </svg>
        <span>Password</span>
      </div>
      {alreadyAdded && (
        <>
          <Setup />
        </>
      )}
    </LinkWrapper>
  );
};

function Setup() {
  return (
    <div className="absolute right-2 top-0 transform">
      <StateBadge evenPadding={true} state={BadgeState.Success}>
        <CheckIcon className="h-4 w-4" />
      </StateBadge>
    </div>
  );
}
