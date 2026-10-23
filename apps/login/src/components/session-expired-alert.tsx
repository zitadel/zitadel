import Link from "next/link";
import { Alert } from "./alert";
import { Button, ButtonVariants } from "./button";
import { Translated } from "./translated";

type Props = {
  hasContext: boolean;
  loginName?: string;
  requestId?: string;
  organization?: string;
};

// The session a page needs is gone: say so and offer a way back to the start that keeps the auth request alive.
export function SessionExpiredAlert({ hasContext, loginName, requestId, organization }: Props) {
  const params = new URLSearchParams();
  if (loginName) {
    params.set("loginName", loginName);
  }
  if (requestId) {
    params.set("requestId", requestId);
  }
  if (organization) {
    params.set("organization", organization);
  }
  const query = params.toString();

  return (
    <div className="flex flex-col items-center space-y-4 py-4">
      <Alert>
        <Translated i18nKey={hasContext ? "sessionExpired" : "unknownContext"} namespace="error" />
      </Alert>
      <Link href={query ? `/loginname?${query}` : "/loginname"} data-testid="sign-in-again-link">
        <Button type="button" variant={ButtonVariants.Secondary}>
          <Translated i18nKey="signInAgain" namespace="error" />
        </Button>
      </Link>
    </div>
  );
}
