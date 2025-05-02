import { useTranslations } from "next-intl";
import Link from "next/link";
import { Button, ButtonVariants } from "./button";

export function ConsentScreen({
  scope,
  nextUrl,
}: {
  scope?: string[];
  nextUrl: string;
}) {
  const t = useTranslations();

  return (
    <div className="w-full flex flex-col items-center space-y-4">
      <ul className="list-disc space-y-2 w-full">
        {scope?.map((s) => {
          const translationKey = `device.scope.${s}`;
          const description = t(translationKey, null);

          // Check if the key itself is returned and provide a fallback
          const resolvedDescription =
            description === translationKey
              ? "No description available."
              : description;

          return (
            <li
              key={s}
              className="grid grid-cols-4 w-full text-sm flex flex-row items-center bg-background-light-400 dark:bg-background-dark-400  border border-divider-light py-2 px-4 rounded-md transition-all"
            >
              <strong className="col-span-1">{s}</strong>
              <span className="col-span-3">{resolvedDescription}</span>
            </li>
          );
        })}
      </ul>

      <div className="mt-4 flex w-full flex-row items-center">
        <Button variant={ButtonVariants.Destructive} data-testid="deny-button">
          Deny
        </Button>
        <span className="flex-grow"></span>

        <Link href={nextUrl}>
          <Button
            data-testid="submit-button"
            type="submit"
            className="self-end"
            variant={ButtonVariants.Primary}
          >
            continue
          </Button>
        </Link>
      </div>
    </div>
  );
}
