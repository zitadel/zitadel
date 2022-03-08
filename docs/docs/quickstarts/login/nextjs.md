---
title: Next.js
---

This guide shows you how to integrate ZITADEL with your [Next.js](https://nextjs.org/) app. 

It covers how to
- Authenticate as a user
- Retrieve user information from the OIDC endpoint.

> The template code is part of our zitadel-example repo. Take a look [here](https://github.com/caos/zitadel-examples/tree/main/nextjs).

## Get Started

1. To start, create a new NextJS app with `npx create-next-app`. This sets up everything automatically for you. 

1. Then, create a project with:

```bash
npx create-next-app --typescript
# or
yarn create next-app --typescript
```

## Install Authentication library

To keep the template as easy as possible, we use [next-auth](https://next-auth.js.org/) as our main authentication library. 

1. To install, run:

```bash
npm i next-auth
# or
yarn add next-auth
```

2. To run the app, type:

```bash
npm run dev
```

3. To check the result, open [http://localhost:3000](http://localhost:3000) 
in your browser.

## Configuration

NextAuth.js exposes a REST API that your client uses.

1. To setup your configuration, create a file called [...nextauth].tsx in `pages/api/auth`.

2. Paste the following snippet in.

```ts
import NextAuth from 'next-auth';

export const ZITADEL = {
    id: "zitadel",
    name: "zitadel",
    type: "oauth",
    version: "2.0",
    scope: "openid profile email",
    params: { response_type: "code", grant_type: "authorization_code" },
    authorizationParams: { grant_type: "authorization_code", response_type: "code" },
    accessTokenUrl: "https://api.zitadel.dev/oauth/v2/token",
    requestTokenUrl: "https://api.zitadel.dev/oauth/v2/token",
    authorizationUrl: "https://accounts.zitadel.dev/oauth/v2/authorize",
    profileUrl: "https://api.zitadel.dev/oauth/v2/userinfo",
    protection: "pkce",
    async profile(profile, tokens) {
        console.log(profile, tokens);
        return {
            id: profile.sub,
            name: profile.name,
            email: profile.email,
            image: profile.picture
        };
    },
    clientId: process.env.ZITADEL_CLIENT_ID,
    session: {
        jwt: true,
    },
};

export default NextAuth({
    providers: [
        ZITADEL
    ],
});
```

3. Replace the `https//api.zitadel.dev` endpoint with your endpoint.
If you use ZITADEL CLOUD tier, the endpoint is `https://api.zitadel.ch/`.
If you use a self-hosted ENTERPRISE tier, replace it with your own endpoint.

We recommend using the Authentication Code flow secured by PKCE for the Authentication flow.

To connect to ZITADEL:

1. Navigate to your [Console Projects](https://console.zitadel.ch/projects).
1. Create or select an existing project.
1. To add your app, select **WEB**, then **PKCE**. 
1. Add `http://localhost:3000/api/auth/callback/zitadel` as the redirect URL to your app. 

   For simplicity, we use the same default as the one that next-auth provides.
   You can change the redirect later if you want.
1. Hit **Create**, then in the detail view of your application, make sure to enable dev mode. Dev mode ensures that you can start an auth flow from a non https endpoint for testing.

> Note that we get a clientId but no clientSecret because it is not needed for our authentication flow.

## Environment

1. Create a file `.env` in the root of the project. 
1. Add the following keys to it.

```
NEXTAUTH_URL=http://localhost:3000
ZITADEL_CLIENT_ID=[yourClientId]
```

# User interface

Now we can start editing the homepage by modifying `pages/index.tsx`.

Add this snippet your file. This code gets your auth session from next-auth.
If you are authenticated, it renders a Logout button.
If you aren't, it shows a Signup button.

Note that the `signIn` method requires the id of the provider we provided earlier.
It also lets you add a callback URL.
If login is successful, 
Auth Next will redirect you to the specified route.

```ts
import { signIn, signOut, useSession } from 'next-auth/client';

export default function Page() {
    const [session, loading] = useSession();
    ...
    {!session && <>
        Not signed in <br />
        <button onClick={() => signIn('zitadel', { callbackUrl: 'http://localhost:3000/profile' })}>
            Sign in
        </button>
    </>}
    {session && <>
        Signed in as {session.user.email} <br />
        <button onClick={() => signOut()}>Sign out</button>
    </>}
    ...
}
```

### Session state

Sharing session states between pages has a multiple benefits:
- Improves performance,
- Reduces network traffic,
- Avoids component state changes while rendering.

To allow session state sharing, you can use the NextAuth.js Provider in `/pages/_app.tsx`.
Take a look at the template `_app.tsx`.

```ts
import { Provider } from 'next-auth/client';

function MyApp({ Component, pageProps }) {
    return <Provider
        session={pageProps.session} >
        <Component {...pageProps} />
    </Provider>;
}

export default MyApp;
```

To render the callback page, create a `profile.tsx` in `/pages`.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/deployment) for more details.
