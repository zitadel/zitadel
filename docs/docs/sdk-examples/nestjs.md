---
title: NestJS
sidebar_label: NestJS
framework_url: https://nestjs.com
auth_library: "@auth/core"
auth_library_url: https://www.npmjs.com/package/@auth/core
example_repo: https://github.com/zitadel/example-auth-nestjs
auth_flow: pkce
status: stable
---

## Overview

[NestJS](https://nestjs.com) is a progressive Node.js framework for building efficient, reliable, and scalable server-side applications using [TypeScript](https://www.typescriptlang.org). This example demonstrates how to integrate **Zitadel** using the **[OAuth 2.0 PKCE flow](/docs/concepts/pkce)** to authenticate users securely and maintain sessions across your NestJS application.

#### Auth library

This example uses **[@auth/core](https://www.npmjs.com/package/@auth/core)**, the core authentication library powering [Auth.js](https://authjs.dev), which implements the [OpenID Connect (OIDC)](https://openid.net/connect/) flow, manages [PKCE](/docs/concepts/pkce), performs [token exchange](https://zitadel.com/docs/apis/openidoauth/endpoints), and exposes helpers for [session state](https://zitadel.com/docs/guides/integrate/login/oidc/session-handling). The example integrates this through the **[@mridang/nestjs-auth](https://www.npmjs.com/package/@mridang/nestjs-auth)** [NestJS module](https://docs.nestjs.com/modules), which provides [decorators](https://docs.nestjs.com/custom-decorators), [guards](https://docs.nestjs.com/guards), and middleware for seamless authentication within the NestJS ecosystem.

---

## What this example demonstrates

This example implements a complete authentication flow using [PKCE](/docs/concepts/pkce) with [Zitadel](https://zitadel.com/docs) as the identity provider. Users begin on a public landing page and click a login button to authenticate through Zitadel's authorization server. After successful authentication, they're redirected to a protected profile page displaying their user information retrieved from the [ID token](https://zitadel.com/docs/apis/openidoauth/claims) and [access token](https://zitadel.com/docs/concepts/tokens).

The application leverages [NestJS controllers](https://docs.nestjs.com/controllers) to handle authentication routes and [@mridang/nestjs-auth guards](https://www.npmjs.com/package/@mridang/nestjs-auth) to protect routes requiring authentication. [Session management](https://zitadel.com/docs/guides/integrate/login/oidc/session-handling) is handled through encrypted [JWT-based sessions](https://authjs.dev/concepts/session-strategies#jwt-session) stored in secure, HTTP-only cookies. The example includes automatic [token refresh](https://zitadel.com/docs/apis/openidoauth/grant-types#refresh-token) functionality using [refresh tokens](https://oauth.net/2/refresh-tokens/), ensuring users maintain their sessions without interruption when access tokens expire.

The logout implementation demonstrates [federated logout](https://zitadel.com/docs/guides/integrate/login/oidc/logout) by redirecting users to Zitadel's end-session endpoint, terminating both the local application session and the Zitadel session. [CSRF protection](https://owasp.org/www-community/attacks/csrf) during logout is achieved through a state parameter validated in the callback. The example also showcases accessing Zitadel's [UserInfo endpoint](https://zitadel.com/docs/apis/openidoauth/endpoints#userinfo) to fetch real-time user data, including [custom claims](https://zitadel.com/docs/apis/openidoauth/claims), [roles](https://zitadel.com/docs/guides/manage/console/roles), and organization membership.

All protected routes are secured using the [@mridang/nestjs-auth global guard](https://www.npmjs.com/package/@mridang/nestjs-auth), with the `@Public()` decorator marking routes that don't require authentication. The application uses [Handlebars templates](https://handlebarsjs.com) for server-side rendering and [Tailwind CSS](https://tailwindcss.com) for styling, providing a complete reference implementation for NestJS developers.

---

## Getting started

### Prerequisites

Before running this example, you need to create and configure a PKCE application in the [Zitadel Console](https://zitadel.com/docs/guides/manage/console/overview). Follow the **[PKCE Application Setup guide](/docs/guides/integrate/login/oidc/pkce-setup)** to:

1. Create a new Web application in your Zitadel project
2. Configure it to use the [PKCE authentication method](/docs/concepts/pkce)
3. Set up your redirect URIs (e.g., `http://localhost:3000/auth/callback` for development)
4. Configure post-logout redirect URIs (e.g., `http://localhost:3000/auth/logout/callback`)
5. Copy your **Client ID** for use in the next steps

> **Note:** Make sure to enable **Dev Mode** in the Zitadel Console if you're using HTTP URLs during local development. For production, always use HTTPS URLs and disable Dev Mode.

### Run the example

Once you have your Zitadel application configured:

1. Clone the [repository](https://github.com/zitadel/example-auth-nestjs).
2. Create a `.env` file in the project root and configure it with the values from your [Zitadel application](https://zitadel.com/docs/guides/manage/console/overview). **Use the exact environment variable names from the repository:**
   ```
   NODE_ENV=development
   PORT=3000

   SESSION_SECRET=your-very-secret-and-strong-session-key
   SESSION_SALT=your-cryptographic-salt
   SESSION_DURATION=3600

   ZITADEL_DOMAIN=https://your-instance.zitadel.cloud
   ZITADEL_CLIENT_ID=your_client_id_from_console
   ZITADEL_CLIENT_SECRET=your_randomly_generated_secret
   ZITADEL_CALLBACK_URL=http://localhost:3000/auth/callback
   ZITADEL_POST_LOGIN_URL=/profile
   ZITADEL_POST_LOGOUT_URL=http://localhost:3000/auth/logout/callback
   ```
   Replace these values with:
	- Your actual Zitadel instance URL (the **Issuer**)
	- The **Client ID** you copied when creating the application
	- A randomly generated **Client Secret** (generate using: `node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"`)
	- A secure **Session Secret** (generate using the same command)
	- A cryptographic **Session Salt** for cookie encryption
	- The **redirect URI** you configured in the PKCE setup (must match exactly)
	- The **post-logout redirect URI** for the logout callback
3. Install dependencies using [npm](https://www.npmjs.com) and start the development server:
   ```bash
   npm install
   npm run dev
   ```
   The application will be running at `http://localhost:3000`. Visit the URL to verify the authentication flow end-to-end.

---

## Learn more and resources

* **Zitadel documentation:** [https://zitadel.com/docs](https://zitadel.com/docs)
* **PKCE concept:** [/docs/concepts/pkce](/docs/concepts/pkce)
* **PKCE application setup:** [/docs/guides/integrate/login/oidc/pkce-setup](/docs/guides/integrate/login/oidc/pkce-setup)
* **Federated logout:** [https://zitadel.com/docs/guides/integrate/login/oidc/logout](https://zitadel.com/docs/guides/integrate/login/oidc/logout)
* **OIDC integration guide:** [https://zitadel.com/docs/guides/integrate/login/oidc/](https://zitadel.com/docs/guides/integrate/login/oidc/)
* **Framework docs:** [https://nestjs.com](https://nestjs.com)
* **Auth library (@auth/core):** [https://www.npmjs.com/package/@auth/core](https://www.npmjs.com/package/@auth/core)
* **NestJS auth module:** [https://www.npmjs.com/package/@mridang/nestjs-auth](https://www.npmjs.com/package/@mridang/nestjs-auth)
* **Example repository:** [https://github.com/zitadel/example-auth-nestjs](https://github.com/zitadel/example-auth-nestjs)
* **All framework examples:** [/docs/examples](/docs/examples)
