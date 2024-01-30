---
title: ZITADEL with Vue
sidebar_label: Vue
---

This integration guide demonstrates the recommended way to incorporate ZITADEL into your Vue application. 
It explains how to enable user login in your application and how to fetch data from the user info endpoint.

By the end of this guide, your application will have login functionality and will be able to access the current user's profile.

:::tip 
This documentation references our [example](https://github.com/zitadel/zitadel-vue) on GitHub.
It also uses the @zitadel/vue package with its default configuration.
:::

## Set up application and obtain keys

Before we begin developing our application, we need to perform a few configuration steps in the ZITADEL Console.
You'll need to provide some information about your app.
We recommend creating a new app to start from scratch.
Navigate to your project, then add a new application at the top of the page.
Select the **User Agent** application type and continue.
We recommend that you use [Proof Key for Code Exchange (PKCE)](/apis/openidoauth/grant-types#proof-key-for-code-exchange) for all single page applications.

![Create app in console](/img/vue/app-create.png)

### Redirect URIs

The redirect URIs field tells ZITADEL where it's allowed to redirect users after authentication. For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.
The post logout redirect sends your users back to a public route on your application after they have logged out.

:::tip
If you are following along with the [example](https://github.com/zitadel/zitadel-vue), set the dev mode switch to `true`.
Configure a redirect URIs to *http:/<span></span>/localhost:5173/auth/signinwin/zitadel* and a post redirect URI to *http:/<span></span>/localhost:5173/*.
:::

Continue and create the application.

### Refresh Token and Client ID

After successful creation of the app, make sure you tick the checkbox to enable refresh tokens.
Also copy the client ID, as you will need it to configure your Vue client.

![Tick refresh token checkbox](/img/vue/tick-refresh-token.png)

## Create a project role "admin" and assign it to your user

Also note the projects resource ID, as you will need it to configure your Vue client.

![Create project role "admin"](/img/vue/project-role.png)

![Assign the "admin" role to your user](/img/vue/project-authz.png)

## Vue setup

Now that you have configured your web application on the ZITADEL side, you can proceed with the integration of your Vue client.

### Install Vue dependencies

To conveniently connect with ZITADEL, you can install the [@zitadel/vue NPM package](https://www.npmjs.com/package/@zitadel/vue). Run the following command:

```bash
npm install --save @zitadel/vue
```

### Create and configure the auth service

The @zitadel/vue package provides a `createZITADELAuth()` function which sets some defaults and calls the underlying [vue-oidc-client packages](https://github.com/soukoku/vue-oidc-client) `createOidcAuth()` function.
You can overwrite all the defaults with the aguments you pass to `createZITADELAuth()`.

Export the object returned from `createZITADELAuth()`

```ts reference
https://github.com/zitadel/zitadel-vue/blob/main/src/services/zitadelAuth.ts
```

### Register the auth service in your global variables when bootstrapping Vue

```ts reference
https://github.com/zitadel/zitadel-vue/blob/main/src/main.ts
```

### Add three new views to your application

The restricted admin view will only be shown if the user is authenticated and has the role "admin" in the apps project in ZITADEL.

```ts reference
https://github.com/zitadel/zitadel-vue/blob/main/src/views/Admin.vue
```

The restricted login view is shown to all authenticated users.
It prints all the information it gets from the token and from the user info endpoint.

```ts reference
https://github.com/zitadel/zitadel-vue/blob/main/src/views/Login.vue
```

The public no access view is shown to authenticated users who navigate to a page they don't have access to based on their roles.

```ts reference
https://github.com/zitadel/zitadel-vue/blob/main/src/views/NoAccess.vue
```

### Add protected routes to your new pages as well as a Signout link

Note that we conditionally render the admin view or the no access view based on the user's roles.

```ts reference
https://github.com/zitadel/zitadel-vue/blob/main/src/router/index.ts
```

## Completion

Congratulations! You have successfully integrated your Vue application with ZITADEL!

If you get stuck, consider checking out the [ZITADEL Vue example application](https://github.com/zitadel/zitadel-vue).
This application includes all the functionalities mentioned in this quickstart.
You can start by cloning the repository and change the arguments to createZITADELAuth so they fit your requirements.
If you face issues, contact us or [raise an issue on GitHub](https://github.com/zitadel/zitadel-vue/issues).

![App in console](/img/vue/app-screen.png)

### What's next?

Now that you have enabled authentication, you are ready to call add authorization to your application using ZITADEL APIs.
To do this, [refer to the API docs](/apis/introduction) or check out [the ZITADEL Console code on GitHub](https://github.com/zitadel/zitadel) which uses gRPC to access data.

For more information on how to create an Vue application, you can refer to [Vue](https://vuejs.org/guide/quick-start.html).
If you want to learn more about the libraries wrapped by [@zitadel/vue](https://www.npmjs.com/package/@zitadel/vue), [read the docs for vue-oidc-client](https://github.com/soukoku/vue-oidc-client/wiki/V1-Docs).