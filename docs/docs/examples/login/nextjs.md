---
title: Next.js
---

This is our Zitadel [Next.js](https://nextjs.org/) template. It shows how to authenticate as a user and retrieve user information from the OIDC endpoint.

> The template code is part of our zitadel-nextjs repo. Take a look [here](https://github.com/zitadel/zitadel-nextjs).

## Getting Started

### Install dependencies

To install the dependencies type:

```bash
yarn install
```

then to run the app:

```bash
npm run dev
```

then open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Setup Application and Get Keys

Before we can start building our application, we have to do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app.
Navigate to your Project, then add a new application at the top of the page.
Select Web application type and continue.
We recommend you use [Authorization Code](/apis/openidoauth/grant-types#authorization-code) in combination with [Proof Key for Code Exchange (PKCE)](/apis/openidoauth/grant-types#proof-key-for-code-exchange) for all web applications.
As the requests from your application to ZITADEL are made on NextJS serverside, you can select `CODE` in the next step. This makes sure you still get a secret which is then used in combination with PKCE. Note that the secret never gets exposed on the browser and is therefore kept in a confidential environment.

![Create app in console](/img/nextjs/app-create.png)

### Redirect URIs

With the Redirect URIs field, you tell ZITADEL where it is allowed to redirect users to after authentication. For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.

> If you are following along with the [example](https://github.com/zitadel/zitadel-angular), set dev mode to `true` and the Redirect URIs to <http://localhost:300/api/auth/callback/zitadel>.

If you want to redirect the users back to a route on your application after they have logged out, add an optional redirect in the Post Logout URIs field.

Continue and create the application.

### Client ID

After successful app creation, a pop-up will appear, showing the app's client ID. Copy the client ID, as you will need it to configure your NextJS app.

## NextJS Setup

Now that you have your web application configured on the ZITADEL side, you can go ahead and integrate your NextJS app.

### Configuration

NextAuth.js exposes a REST API which is used by your client.
To setup your configuration, create a file called [...nextauth].tsx in `pages/api/auth`.
You can directly import the ZITADEL provider from [next-auth](https://next-auth.js.org/providers/zitadel).

```ts reference
https://github.com/zitadel/zitadel-nextjs/blob/main/pages/api/auth/%5B...nextauth%5D.tsx
```

You can overwrite the default callbacks, just append them to the ZITADEL provider.

```ts
...
ZitadelProvider({
    issuer: process.env.ZITADEL_ISSUER,
    clientId: process.env.ZITADEL_CLIENT_ID,
    clientSecret: process.env.ZITADEL_CLIENT_SECRET,
    async profile(profile) {
        return {
          id: profile.sub,
          name: profile.name,
          firstName: profile.given_name,
          lastName: profile.family_name,
          email: profile.email,
          loginName: profile.preferred_username,
          image: profile.picture,
        };
    },
}),
...
```

We recommend using the Authentication Code flow secured by PKCE for the Authentication flow.
To be able to connect to ZITADEL, navigate to your Console Projects, create or select an existing project and add your app selecting WEB, then PKCE, and then add `http://localhost:3000/api/auth/callback/zitadel` as redirect url to your app.

For simplicity reasons we set the default to the one that next-auth provides us. You'll be able to change the redirect later if you want to.

Hit Create, then in the detail view of your application make sure to enable dev mode. Dev mode ensures that you can start an auth flow from a non https endpoint for testing.

> Note that we get a clientId but no clientSecret because it is not needed for our authentication flow.

Now go to Token settings and check the checkbox for **User Info inside ID Token** to get your users name directly on authentication.

### Environment

Create a file `.env` in the root of the project and add the following keys to it.
You can find your Issuer Url on the application detail page in console.

```env reference
https://github.com/zitadel/zitadel-nextjs/blob/main/.env
```

next-auth requires a secret for all providers, so just define a random value here.

### User interface

Now we can start editing the homepage by modifying `pages/index.tsx`. On the homepage, your authenticated user or a Signin button is shown.

Add the following component to render the UI elements:

```ts reference

https://github.com/zitadel/zitadel-nextjs/blob/main/components/profile.tsx#L4-L38
```

Note that the signIn method requires the id of our provider which is in our case `zitadel`.

### Session state

To allow session state to be shared between pages - which improves performance, reduces network traffic and avoids component state changes while rendering - you can use the NextAuth.js Provider in `/pages/_app.tsx`.
Take a loot at the template `_app.tsx`.

```ts reference
https://github.com/zitadel/zitadel-nextjs/blob/main/pages/_app.tsx
```

Last thing: create a `profile.tsx` in /pages which renders the callback page.

```ts reference
https://github.com/zitadel/zitadel-nextjs/blob/main/pages/profile.tsx
```
