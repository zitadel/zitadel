---
title: Authenticate service users and client applications
sidebar_label: Authenticate service users
sidebar_position: 1
---

This guide explains ZITADEL service users and their role in facilitating secure machine-to-machine communication within your applications.

## What are Service Users?

Service users in ZITADEL represent **non-human entities** within your system.
They are ideal for scenarios involving secure communication between applications, particularly when interacting with backend services or APIs.
Service users in combination with [Manager](/concepts/structure/managers) permissions are used to access ZITADEL's APIs, for example, to manage user resources.
Unlike regular human users, service users don't rely on traditional login methods (e.g., username/password) and require alternative authentication mechanisms.

## Benefits of using Service Users

### Enhanced security

* **Principle of Least Privilege:** Grant service users only the minimum permissions they need, minimizing potential damage in case of compromise.
* **Distinct Credentials:** Avoid embedding sensitive credentials like API keys directly in code. Service user credentials can be rotated independently.

### Segregated authorization

Manage authorization for service users separately from human users, providing an extra layer of control.

### API and backend access

Service users offer a secure way to authenticate and access various API endpoints and protected backend services.

You can [use service users to access ZITADEL APIs](../zitadel-apis/access-zitadel-apis), follow the guides to learn how to access the different ZITADEL APIs.
While you can define the scopes and required information in your requests for your applications API endpoints, when using the ZITADEL APIs, you must include the scope `urn:zitadel:iam:org:project:id:zitadel:aud` to gain access.

### Improved auditability

Actions performed by service users are clearly identifiable in logs, facilitating easier auditing and tracing.

Using the [Event API](../zitadel-apis/event-api) you can use these logs for further analysis or to integrate the logs with [external SOC / SIEM](../external-audit-log) systems.

## Authentication methods

ZITADEL supports two primary authentication methods for service users:

### Private key JWT authentication

#### How private key JWT authentication works

* Generate a private/public key pair associated with the service user.
* Sign JWTs with the private key.
* ZITADEL validates the signature using the service user's public key.
* JWTs can include expiration dates and scopes to control access.

Follow our guide on using [private key JWT client authentication](./private-key-jwt) to get started authenticating service users and clients.

#### Benefits of private key JWT authentication

* **Decentralized Verification:** No need for constant server calls, improving performance and scalability.
* **Flexibility and Control:** Define scopes and expiration within the JWT itself for granular access control.
* **Stateless:** The server doesn't need to maintain a session state, simplifying server implementation.

#### Drawbacks of private key JWT authentication

* **Complexity:** Slightly more complex to implement compared to other methods, requiring knowledge of JWT and digital signing.
* **Revocation:** Invalidating a JWT before its expiry can be challenging; blacklisting mechanisms might be required. 

#### Security considerations when using private key JWT authentication

* **Secure Key Storage:** The private key used for signing must be stored with the highest level of security. Compromise could allow attackers to forge tokens.
* **Short Expirations:** Implementing short expiration durations for JWTs helps limit the impact of stolen tokens.

### Client credentials grant

* Presents a client ID and client secret associated with the service user.
* Simpler than the JWT profile in specific scenarios.

Follow our guide on using [client credentials grant](./client-credentials) to get started authenticating service users and clients.

This method is still available in ZITADEL but is generally considered less secure than JWT due to:

* **Centralized Validation:** Relies on the server to verify credentials for every request, potentially impacting performance and requiring more server resources.
* **Credentials Exposure:** Leaked client ID and secret could be used by attackers to impersonate the service user until rotation occurs.

### Personal Access Tokens (PATs)

* **Ready-to-use tokens:** Generated for specific service users and can be directly included in the authorization header of API requests.
* **Currently available only for machine users** (service users) and not regular human users.

Follow our guide on using [personal access tokens](./personal-access-token) to get started authenticating service users and clients.

PAT offer some benefits, such as:

* **Ease of Use:** Ready-to-use tokens, eliminating the need for complex signing logic.

However, PATs also come with limitations:

* **Centralized Validation:** Similar to Client Credentials, relying on the server for verification could impact performance under high load.
* **Revocation:** Requires deleting the PAT directly, potentially causing downtime if not managed carefully.
* **Leakage:** PATs are long-lived tokens that can be readily used in API calls, if leaked the attacker can access all resources until the PAT is expired or deleted. Private key JWT and client credentials create a short-lived access token instead.

## Using Service Users

1. **Creation:** Access the ZITADEL management console and create a new service user. Assign a descriptive name that reflects its purpose. Follow our detailed guide on [how to create service users](../../manage/console/users).
2. **Credentials:** Choose your preferred authentication method (JWT or Client Credentials) and securely store the generated credentials (private key, client secret).
3. **Making API Calls:** When your service needs to make an API call:
    * **For JWT:** Generate and sign a JWT. Include it in the "Authorization" header of your API request.
    * **For Client Credentials:** Include the client ID and client secret in your API request.
    * **For PATs:** Include the PAT directly in the "Authorization" header of your API request.
4. ZITADEL Verifies the credentials and authorizes the service user to perform the requested action based on its granted permissions.

We have guides for the different authentication methods:

- [Private key JWT authentication](private-key-jwt)
- [Client credential authentication](client-credentials)
- [Personal access token authentication](personal-access-token)

## Important considerations

* **Secure Credentials:** Treat service user credentials (private keys, client secrets) with utmost care. Store them securely, similar to any other sensitive information like API keys or passwords.
* **Expiry Management:** Set appropriate expiration dates for JWTs and regularly rotate all credentials to maintain strong security practices.
* **Permission Granting:** Adhere to the principle of least privilege by granting only the specific permissions required for a service user's function.

## Choosing the right authentication method

For most service user scenarios in ZITADEL, [private key JWT authentication](./private-key-jwt.md) is the recommended choice due to its benefits in security, performance, and control.
However, [client credentials authentication](./client-credentials.md) might be considered in specific situations where simplicity and trust between servers are priorities.

## Further resources

* Read about the [different methods to authenticate service users](./authenticate-service-users)
* [Service User API reference](/docs/category/apis/resources/mgmt/user-machine)
* [OIDC JWT with private key](/docs/apis/openidoauth/authn-methods#jwt-with-private-key) authentication method reference
* [Access ZITADEL APIs](../zitadel-apis/access-zitadel-apis)
* Validate access tokens with [token introspection with private key jwt](../token-introspection/private-key-jwt.mdx)

import DocCardList from '@theme/DocCardList';

<DocCardList />
