---
title: ZITADEL with React
sidebar_label: React
---

This integration guide demonstrates the recommended way to incorporate ZITADEL into your React application.
It explains how to enable user login in your application and how to fetch data from the user info endpoint.

By the end of this guide, your application will have login functionality and will be able to access the current user's profile.

:::tip
This documentation references our [example](https://github.com/zitadel/zitadel-react) on GitHub.
It also uses the @zitadel/react package with its default configuration.
:::

## Set up application and obtain keys

Before we begin developing our application, we need to perform a few configuration steps in the ZITADEL Console.
You'll need to provide some information about your app.
We recommend creating a new app to start from scratch.
Navigate to your project, then add a new application at the top of the page.
Select the **User Agent** application type and continue.
We recommend that you use [Proof Key for Code Exchange (PKCE)](/apis/openidoauth/grant-types#proof-key-for-code-exchange) for all single page applications.

![Create app in console](/img/react/app-create.png)

### Redirect URIs

The redirect URIs field tells ZITADEL where it's allowed to redirect users after authentication. For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.
The post logout redirect sends your users back to a public route on your application after they have logged out.

:::tip
If you are following along with the [example](https://github.com/zitadel/zitadel-react), set the dev mode switch to `true`.
Configure a redirect URIs to \http://localhost:3000/callback and a post redirect URI to \http://localhost:3000/.
:::

Continue and create the application.

### Copy Client ID

After successful creation of the app, make sure copy the client ID, as you will need it to configure your React client.

## Create a project role "admin" and assign it to your user

Also note the projects resource ID, as you will need it to configure your React client.

![Create project role "admin"](/img/react/project-role.png)

![Assign the "admin" role to your user](/img/react/project-authz.png)

If you want to read your users roles from user info endpoint, make sure to enable the checkbox in your project.

## React setup

Now that you have configured your web application on the ZITADEL side, you can proceed with the integration of your React client.

### Install React dependencies

To conveniently connect with ZITADEL, you can install the [@zitadel/react NPM package](https://www.npmjs.com/package/@zitadel/react). Run the following command:

```bash
yarn add @zitadel/react
```

### Create and configure the auth service

The @zitadel/react package provides a `createZitadelAuth()` function which sets some defaults and initializes the underlying [oidc-client-ts](https://github.com/authts/oidc-client-ts) `UserManager` class.
You can overwrite all the defaults with the aguments you pass to `createZitadelAuth()`.

Export the object returned from `createZitadelAuth()`

### Initialize user manager

```ts reference
https://github.com/zitadel/zitadel-react/blob/main/src/App.tsx
```

### Add two new components to your application

First, add the component which prompts the user to login.

```ts reference
https://github.com/zitadel/zitadel-react/blob/main/src/components/Login.tsx
```

Then create the component for the page where the users will be redirected.
It loads the user info endpoint once the code flow completes and prints all the information.

```ts reference
https://github.com/zitadel/zitadel-react/blob/main/src/components/Callback.tsx
```

You can now read a user's role to show protected areas of the application.

### Run

Finally, you can start your application by running the following:

```
yarn start
```

## Completion

Congratulations! You have successfully integrated your React application with ZITADEL!

If you get stuck, consider checking out the [ZITADEL React example application](https://github.com/zitadel/zitadel-react).
This application includes all the functionalities mentioned in this quickstart.
You can start by cloning the repository and changing the arguments to `createZitadelAuth` to fit your requirements.
If you face issues, contact us or [raise an issue on GitHub](https://github.com/zitadel/zitadel-react/issues).

![App in console](/img/react/app-screen.png)

### What's next?

Now that you have enabled authentication, you are ready to add authorization to your application by using ZITADEL APIs.
To do this, [refer to the API docs](/apis/introduction) or check out [the ZITADEL Console code on GitHub](https://github.com/zitadel/zitadel) which uses gRPC to access data.

For more information on how to create a React application, you can refer to [Create React App](https://github.com/facebook/create-react-app).
If you want to learn more about the libraries wrapped by [@zitadel/react](https://www.npmjs.com/package/@zitadel/react), read the docs for [oidc-client-ts](https://github.com/authts/oidc-client-ts).
