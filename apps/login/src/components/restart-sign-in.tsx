import { clsx } from "clsx";
import Link from "next/link";
import { textLinkClasses } from "./text-button";
import { Translated } from "./translated";

type Props = {
  loginName?: string;
  organization?: string;
  requestId?: string;
};

/**
 * Way out of a step whose session can no longer be resolved: back to the login name step,
 * keeping the auth request so the application's sign-in continues afterwards.
 */
export function RestartSignIn({ loginName, organization, requestId }: Props) {
  const params = new URLSearchParams();

  if (loginName) {
    params.set("loginName", loginName);
  }
  if (organization) {
    params.set("organization", organization);
  }
  if (requestId) {
    params.set("requestId", requestId);
  }

  return (
    <Link href={`/loginname?${params}`} className={clsx(textLinkClasses, "mt-4 inline-block")} data-testid="restart-sign-in">
      <Translated i18nKey="sessionResolution.restart" namespace="error" />
    </Link>
  );
}
