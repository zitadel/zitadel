import { Callout } from 'fumadocs-ui/components/callout';
import Link from 'next/link';

export function PreventLockout() {
  return (
    <Callout type="warn" title="Prevent Settings Misconfiguration Lockouts">
      Login policy settings misconfigurations that occur during the testing phase can easily lead to a{' '}
      <Link href="/legal/policies/account-lockout-policy">lockout</Link>. To ensure you don't lose access to your instance:
      <ol className="list-decimal list-inside mt-2">
        <li>
          <strong>Generate a backup PAT:</strong> Create a{' '}
          <Link href="/guides/integrate/service-accounts/personal-access-token">
            Service Account Personal Access Token
          </Link>{' '}
          with the <Link href="/guides/manage/console/administrators"><code>IAM_OWNER</code></Link> role to revert any login UI misconfigurations using the API.
        </li>
        <li>
          <strong>Add a second Instance Administrator:</strong> Always designate at least one{' '}
          <Link href="/guides/manage/console/administrators">
            <strong>second instance administrator</strong>
          </Link>.
        </li>
      </ol>
    </Callout>
  );
}
