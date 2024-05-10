---
title: Account linking
---

ZITADEL supports linking of user accounts from different external identity providers such as social logins or enterprise IdPs.
A user can authenticate from any of their accounts and still be recognized by your app and associated with the same user profile.

In ZITADEL, users have one user account to simplify access management and provide a consistent audit trail.
A user account can have a link to multiple external identities.

### Advantages

- Users can login with multiple identity providers without managing separate profiles
- Already registered users can link existing profiles
- Provide backup authentication methods in case an IdP is unavailable
- Unified audit trail across multiple identities

### How it works

ZITADEL gives you great flexibility to configure account linking for an organization and based on the external identity provider.

When using external identity providers (ie. social login, enterprise SSO), a user account will be created in ZITADEL.
The [external identity](../structure/users#federated-users) will be linked to the ZITADEL account.

If login with "Username / Password" (ie. local account) is enabled and you have configured external IDPs, the user can decide if they want to login with an external IDP or the local account.

When only one external identity provider is configured and login with "Username / Password" is disabled, then the user is immediately redirected to the external identity provider.

In cases when a local account already exists and a user logs in with an external identity provider, you can instruct ZITADEL to link the external identity to the local account based on the username or email address.

### Automatic account linking

You can link accounts with the same email or username and prompt users to link them.
On an [identity provider template settings](/docs/guides/integrate/identity-providers/introduction#key-settings-on-the-templates), you must enable "Account linking allowed".

Automatic account linking is beneficial for users who wish to associate multiple login methods with their ZITADEL account, providing flexibility and convenience in how they access your application.
