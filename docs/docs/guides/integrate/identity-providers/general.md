---
title: General configurations for identity providers
sidebar_label: General
---

This guide shows you what you need to configure no matter which identity provider template you choose.


## Add custom login policy

The login policy can be configured on two levels. Once as default on the instance and this can be overwritten for each organization.
The only difference is where you configure it. Go either to the settings page of a specific organization or to the settings page of your instance.
Instance: $YOUR-DOMAIN/ui/console/settings?id=general
Organization: Choose the organization in the menu and go to $YOUR-DOMAIN/ui/console/org-settings?id=login

1. Go to the Settings
2. Modify your login policy in the menu "Login Behavior and Security"
3. Enable the attribute "External IDP allowed"

![Allow External IDP](/img/guides/zitadel_allow_external_idp.png)

### Trigger configuration on the login for a specific organization

Per default ZITADEL will always show the settings configured on your instance, because these are the default settings.
If you have overwritten the settings on an organization and you want to trigger them you have to send the organization scope in your request.

The organization scope does look like this: ```urn:zitadel:iam:org:id:{id}``` or you can read more about the reserved scopes [here](/apis/openidoauth/scopes#reserved-scopes)
Or use our [OIDC Playground](/apis/openidoauth/authrequest) to try out what happens with the login if you have different scopes.

## Configure identity provider

Now that you have allowed to use identity providers for your login you can configure the providers you need.
Go to the settings page of your instance or organization and choose "Identity Providers".

In the table you can see all the providers you have configured and below the different providers
![Identity Provider Overview](/img/guides/zitadel_identity_provider_overview.png)

To setup your specific providers go to the configuration guide you need:
- [Google](./google)
- [GitHub](./github)
- [GitLab](./gitlab)
- [LDAP](./ldap)