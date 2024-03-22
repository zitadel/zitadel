---
title: Let Users Login with Preferred Identity Provider
sidebar_label: Login users with SSO
sidebar_position: 1
---

## External Identity Providers and SSO authentication

An **Identity Provider (IdP)** is a service that creates and maintains identity information and then provides authentication services to your applications. Incorporating SSO (Single Sign-On) authentication allows users to access multiple services with a single set of credentials, enhancing user convenience and security. As you develop your application, you want to give users the freedom to choose an Identity Provider they use to sign in to your application, despite having your own IdP (e.g., ZITADEL). Many users already have trusted accounts with popular IdPs like Google, GitHub, or LinkedIn, making these services effective identity brokers in the SSO landscape. Or it could be an IdP they’re already using in the workplace, such as EntraID or Auth0, which supports SSO authentication. These can be identified as *External Identity Providers*. Allowing them to use these accounts for signing into your application means they don’t need to create new usernames and passwords, thereby significantly enhancing their convenience and the likelihood they’ll engage with your product.

Using external IdPs can also improve security and user trust in your application. These major IdPs have robust security measures in place, reducing the burden on your application to manage secure passwords. Furthermore, by leveraging external IdPs as part of an SSO framework, you're aligning with a user-centric approach, prioritizing their preferences and simplifying their access to your services.

However, integrating these external IdPs with your system requires a clear understanding of the connection process with your IdP, acting as an  identity broker  to facilitate  SSO authentication . This ensures that you can offer this flexibility without compromising on security or user experience. Integrating external IdPs alongside your own offers the best of both worlds: the ease and security of established accounts through SSO and the tailored experience of your application.


## Where ZITADEL fits in

ZITADEL positions itself as the central hub or the identity broker in the interaction between your application and various external IdPs. It handles the orchestration of authentication requests and manages the seamless flow of identity verification between your application and the chosen IdP. This is achieved through a process known as federation, where ZITADEL acts as a federated identity provider, integrating IdPs via authentication protocols like OpenID Connect and SAML.


## Adding external identity providers to your application

With ZITADEL, you can enhance your application's accessibility by integrating social login IdPs such as Google or other IdPs such as EntraID (formerly known as AzureAD) and Auth0. For organizations with bespoke identity solutions, ZITADEL supports integration with custom-built IdPs that adhere to OpenID Connect or SAML protocols. By default, ZITADEL can serve as the primary user store for your applications. This centralizes user management and simplifies the authentication process, allowing users to sign in with their email and password, or via the external IdPs you've integrated.

ZITADEL excels in B2B scenarios by offering the flexibility to configure different IdPs for distinct customers, enhancing ease of use and customization. For instance, you could have Customer A utilize EntraID for authentication, while Customer B uses Okta. This versatility enables you to tailor authentication solutions to meet the specific needs of each customer, streamlining their access while maintaining a centralized management system.


## The advantages of using ZITADEL

The benefits of integrating ZITADEL for managing external IdPs are multiple:

- **No need for custom authentication code**: Your application only needs to interact with ZITADEL, which handles all communications with external IdPs. This abstraction saves you from the complexity of direct integration with multiple IdPs.


- **Unified protocol handling**: ZITADEL abstracts away the specific protocols used by different IdPs. Your application communicates with ZITADEL using a standard protocol (e.g., OpenID Connect), while ZITADEL takes care of the rest.


- **Centralized user management**: All user profiles are managed within ZITADEL, allowing for a unified view of user identities regardless of the IdP used for authentication.


- **Dynamic profile synchronization**: When users update their profiles on an external IdP, those changes are reflected in ZITADEL at the next login, ensuring that user data remains current.


- **Simplified account linking**: ZITADEL can link identities from multiple IdPs to a single user profile, facilitating a cohesive user experience across different authentication methods.



## The user journey 

- **Integration with external identity providers**: ZITADEL supports integrating a variety of external identity providers, including social logins like Google or GitHub, as well as custom IdPs that use OpenID Connect or SAML protocols.


- **Initiation of sign-in process**: Users select to sign in with an external IdP within your application.


- **Redirection to ZITADEL**: The application redirects the user to ZITADEL, which guides them to their chosen external IdP for authentication.


- **User authentication**: After successfully logging in at the external IdP site, users are redirected back to the application through ZITADEL, bringing along authentication tokens and profile information.


- **Account linking for new users**: If the user's identity from the external IdP does not exist in ZITADEL, they're presented with two options before moving forward:
  - **Create a new account**: Choosing this option creates a new ZITADEL account linked to the external IdP.
  - **Linking with an existing local account**: Users have the option to link their new external identity to an existing local account in ZITADEL, enabling future logins with either their local account or the external IdP.


- **Profile pre-filling from external IdP**: ZITADEL uses information from the external IdP to pre-fill the user's profile, simplifying the account creation or linking process.


- **User option to update profile**: Users can review and, if necessary, update their pre-filled profile information.


- **Session creation and access granting**: After the account is created or linked and the profile is set, the application grants access by creating a session for the user based on their authenticated identity.


## Setting up external identity providers in ZITADEL


In ZITADEL, you have the flexibility to link an external Identity Provider (IdP) to your entire instance, making it the default option for all organizations within your instance, or to connect it exclusively to a specific organization. This setup allows organization members to leverage similar capabilities in self-service if permitted.


### Adjusting the custom login policy


The login policy can be set as a default at the instance level and can be customized for each organization. The configuration process varies slightly depending on your focus:


- **For default settings**, navigate to: `$YOUR-DOMAIN/ui/console/instance?id=general`
- **For specific organization settings**, select the organization from the menu and visit: `$YOUR-DOMAIN/ui/console/org-settings?id=login`


Once in the settings:
- Access the **Login Behavior and Security** section to modify your login policy. Here, ensure you enable the option for **External IDP Allowed**.


![Allow External IDP](/img/guides/zitadel_allow_external_idp.png)



### Configuring IdP Providers


Access the settings page of your instance or the specific organization and select **Identity Providers**.


The ZITADEL Console will display a list of all the IdPs you've configured, along with available provider templates. Selecting any listed IdP will guide you through the process of configuring that specific Identity Provider.

![Identity Provider Overview](/img/guides/zitadel_identity_provider_overview.png)


## Available guides

In the guides below, some of which utilize the Generic OIDC or SAML templates for configuration, you'll learn how to configure and set up your preferred external Identity Provider (IdP) in ZITADEL. 

- [Google](./google)
- [Entra ID (OIDC)](./azure-ad-oidc)
- [Entra ID SAML](./azure-ad-saml)
- [GitHub](./github)
- [GitLab](./gitlab)
- [Apple](./apple)
- [LDAP](./ldap) 
- [Local OpenLDAP](./openldap.mdx)
- [OKTA generic OIDC](./okta-oidc)
- [OKTA SAML](./okta-saml)
- [Keycloak generic OIDC](./keycloak)
- [MockSAML](./mocksaml)
- [JWT IdP](./jwt_idp)


### Configuring IdPs without predefined templates

If ZITADEL doesn't offer a specific template for your Identity Provider (IdP) and your IdP is fully compliant with OpenID Connect (OIDC), you have the option to use the generic OIDC provider configuration.

For those utilizing a SAML Service Provider, the SAML Service Provider option is available. You can learn how to set up a SAML Service Provider with our [MockSAML example](https://zitadel.com/docs/guides/integrate/identity-providers/mocksaml).

Should you wish to transition from a generic OIDC provider to Entra ID (formerly Azure Active Directory) or Google, consider following this [guide](https://zitadel.com/docs/guides/integrate/identity-providers/migrate).



## Key settings on the templates

When configuring external IdP templates in ZITADEL, several common settings enable customized integration to suit your application's authentication flow. Below is a generic explanation of these settings:

- **Scopes**: Specifies the permissions your application requests from the user's account on the external IdP. Common scopes include `openid`, `profile`, and `email`, essential for accessing basic user information for authentication and account management in ZITADEL. You can specify additional scopes based on your application's requirements and the information needed from the external IdP.

- **Automatic creation**: When enabled, this allows ZITADEL to automatically create a new user account if someone logs in with their external IdP credentials and no corresponding account exists in ZITADEL. This facilitates a smooth user onboarding experience by eliminating the need for manual account creation.

- **Automatic update**: This feature, when activated, allows ZITADEL to automatically update a user's profile information whenever changes are detected in the user's account on the external IdP. For example, if a user changes their last name in their Google or Microsoft account, ZITADEL will reflect this update in the user's account upon their next login.

- **Account creation allowed**: Determines whether new user accounts can be created in ZITADEL through the external IdP authentication process. Enabling this setting is crucial for allowing users who are new to your application to register and create accounts seamlessly via their existing external IdP accounts.

- **Account linking allowed**: Enables existing ZITADEL accounts to be linked with identities from external IdPs. It requires that a linkable ZITADEL account already exists for the user attempting to log in with an external IdP. Account linking is beneficial for users who wish to associate multiple login methods with their ZITADEL account, providing flexibility and convenience in how they access your application.



## Configure external IdPs at the organization level or on the default settings

Deciding whether to configure an external Identity Provider (IdP) at the organization level or in the default settings in ZITADEL depends on the scope of access and management you intend to provide. Here’s when to choose each option:

### Create an external IdP on an organization

- **Targeted access control**: When you want to allow authentication through the external IdP specifically for users of a particular organization. This is useful for businesses that manage multiple organizations within ZITADEL and require distinct authentication strategies for each.

- **Customized authentication flow**: If different organizations have unique requirements or policies around user authentication, configuring an IdP at the organization level allows for customized authentication flows that cater to the specific needs of each organization.

- **Delegated administration**: Enabling organizations to manage their own IdPs empowers them to administer their authentication mechanisms independently. This is beneficial in scenarios where organizations have the autonomy to choose their preferred IdPs or when they manage their user base directly.

### Create an external IdP as default

- **Unified authentication strategy**: When you aim to provide a consistent authentication experience across all organizations within your ZITADEL instance. Configuring an IdP in the default settings applies the same authentication mechanism universally, simplifying management and user experience.

- **Centralized management**: Setting up an IdP in the default settings is ideal for scenarios where a single administrative body oversees user authentication across all organizations. This approach centralizes the management of external IdPs, making it easier to maintain and update authentication policies.

- **Broad access needs**: If all users, regardless of their organization, require access to external IdPs for authentication, configuring the IdP in the default settings ensures that these options are available universally. This is particularly useful for platforms that serve a wide range of users with common access requirements.


## References

- [Identity brokering in ZITADEL](https://zitadel.com/docs/concepts/features/identity-brokering)
- [The ZITADEL API reference for managing external IdPs](https://zitadel.com/docs/category/apis/resources/admin/identity-providers)
- [Handle external logins in a custom login UI](https://zitadel.com/docs/guides/integrate/login-ui/external-login)