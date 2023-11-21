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
  - When a language is disallowed, all users who have set this language as their preferred language will be reset to the instances default language.
  - New users can only be created or updated with an allowed preferred language.
  - Custom Texts can also be created for disallowed languages as long as ZITADEL supports that language, so all texts can be customized before allowing a language.

Feature restrictions for an instance are intended to be configured by a user that is managed within that instance.
However, if you are self-hosting and need to control your virtual instances usage, [read about the APIs for limits and quotas](/self-hosting/manage/usage_control) that are intended to be used by system users.
