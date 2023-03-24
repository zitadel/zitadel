# Security Policy

## Introduction

At ZITADEL we are extremely grateful for security aware people who disclose vulnerabilities to us and the open source community.
All reports will be investigated by our team and we will work with you closely to validate and fix vulnerabilities reported to us.

We require that you keep vulnerabilities confidential until we are able to address them, since public disclosure of security vulnerabilities could put the ZITADEL community at risk.

## Scope

The scope of this policy applies to all security issues that concern our Product in form of Software in our [open source repositories](https://github.com/zitadel).

Out of scope are all websites and services operated by ZITADEL (CAOS Ltd.).
Please refer to the separate [vulnerability disclosure policy](https://zitadel.com/docs/legal/vulnerability-disclosure-policy).

### Supported Versions

Supported are releases that are newer and not older than 6 months from our stable release
https://github.com/zitadel/zitadel/blob/main/release-channels.yaml#L1

## Reporting a vulnerability

To file an incident, please disclose it by e-mail to [security@zitadel.com](mailto:security@zitadel.com) including the following details of the vulnerability:

- Target: ZITADEL, Website (zitadel.com), ZITADEL Cloud (zitadel.cloud), Other (please describe)
- Type: For example DoS, authentication bypass, information disclosure, broken authorization, ...
- Description: Provide a detailed explanation of the issue, steps to reproduce, and assumptions you have made
- URL / Location (optional): The URL of the vulnerability
- Contact details (optional): In case we should contact you on a different channel

At the moment GPG encryption is no yet supported, however you may sign your message at will.

Your email will be acknowledged within 48 hours.
We will follow-up within the next 3 business days indicating next steps in handling your report.

If you haven't received a response within 48 hours, or you didn't get a reply from our security team within the last 5 days, please contact [support@zitadel.com](mailto:support@zitadel.com).

Please inform us in your report whether we should mention your contribution.
We will not publish this information by default to protect your privacy.

### When should I NOT report a vulnerability

- Disclosure of known public files or directories, e.g. robots.txt, files under .well-known, or files that are included in our public repositories (eg, go.mod)
- DoS of users when [Lockout Policy is enabled](https://zitadel.com/docs/guides/manage/console/instance-settings#lockout)
- You need help applying security related settings

## Disclosure Process

Our security team will follow the disclosure process: 

1. We will acknowledge the receipt of your vulnerability report
2. Our security team will try to verify, reproduce, and determine the impact of your report
3. A member of our team will respond to either confirm or reject your report, including an explanation
4. Code will be audited to assess if the report uncovers similar issues
5. Fixes are prepared for the latest release
6. On the date that the fixes are applied, we will create a CVE and publish a [security advisory](https://github.com/zitadel/zitadel/security/advisories). Affected users of our Product, Services, or Website will be informed of the fix and required actions.

We think it is crucial to publish advisories `ASAP` as mitigations are ready. But due to the unknown nature of the disclosures the time frame can range from 7 to 90 days.
