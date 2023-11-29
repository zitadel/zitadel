---
title: Feature Restrictions
---

New self-hosted and [ZITADEL Cloud instances](https://zitadel.com/signin) are unrestricted by default.
Self-hosters can change this default using the DefaultInstance.Restrictions configuration section.
Users with the role IAM_OWNER can change the restrictions of their instance using the [Feature Restrictions Admin API](/category/apis/resources/admin/feature-restrictions).
Currently, the following restrictions are available:

- *Disallow public organization registrations* - If restricted, only users with the role IAM_OWNERS can create new organizations. The endpoint */ui/login/register/org* returns HTTP status 404 on GET requests, and 409 on POST requests.
- *AllowedLanguages* - The following rules apply if languages are restricted:
  - Only allowed languages are listed in the OIDC discovery endpoint */.well-kown/openid-configuration*.
  - Login UI texts are only rendered in allowed languages.
  - Notification message texts are only rendered in allowed languages.
  - Custom Texts can be created for disallowed languages as long as ZITADEL supports that language. Therefore, all texts can be customized before allowing a language.

Feature restrictions for an instance are intended to be configured by a user that is managed within that instance.
However, if you are self-hosting and need to control your virtual instances usage, [read about the APIs for limits and quotas](/self-hosting/manage/usage_control) that are intended to be used by system users.
