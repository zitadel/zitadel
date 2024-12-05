---
tags: 
  - FAQ
---

# Why is my key list empty?

Keys are currently managed on expiry. On freshly created instances or instances that did not serve traffic for some time, there might not be an active key and thus the list might be empty. As soon as a user successfully authenticates and is issued a token, a new key pair will be created an the key list will return a public key.

If your clients are not able to handle dynamic key rotation, you can enable an experimental flag to fully manage the OAuth2 / OpenID Connect keys yourself.
Please checkout the corresponding [guide](https://zitadel.com/docs/guides/integrate/login/oidc/webkeys) to do so.
