---
title: Next.js
---

This is our Zitadel [Next.js](https://nextjs.org/) template. It shows how to authenticate as a user and retrieve user information from the OIDC endpoint.

> The template code is part of our zitadel-example repo. Take a look [here](https://github.com/zitadel/zitadel-examples/tree/main/nextjs).

## Getting Started

First, we start by creating a new NextJS app with `npx create-next-app`, which sets up everything automatically for you. To create a project, run:

```bash
npx create-next-app --typescript
# or
yarn create next-app --typescript
```

# Install Authentication library

To keep the template as easy as possible we use [next-auth](https://next-auth.js.org/) as our main authentication library. To install, run:

```bash
npm i next-auth
# or
yarn add next-auth
```

To run the app, type:

```bash
npm run dev
```

then open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

# Configuration

NextAuth.js exposes a REST API which is used by your client.
To setup your configuration, create a file called [...nextauth].tsx in `pages/api/auth`.

```ts
import NextAuth from "next-auth";

export const ZITADEL = {
  id: "zitadel",
  name: "zitadel",
  type: "oauth",
  version: "2",
  wellKnown: process.env.ZITADEL_ISSUER,
  authorization: {
    params: {
      scope: "openid email profile",
    },
  },
  idToken: true,
  checks: ["pkce", "state"],
  client: {
    token_endpoint_auth_method: "none",
  },
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
  clientId: process.env.ZITADEL_CLIENT_ID,
};

export default NextAuth({
  providers: [ZITADEL],
});
```

Replace the endpoints `https:/[your-domain]-[random-string].zitadel.cloud` with your instance or if your using a ZITADEL CLOUD tier or your own endpoint if your using a self hosted ENTERPRISE tier respectively.

We recommend using the Authentication Code flow secured by PKCE for the Authentication flow.
To be able to connect to ZITADEL, navigate to your Console Projects, create or select an existing project and add your app selecting WEB, then PKCE, and then add `http://localhost:3000/api/auth/callback/zitadel` as redirect url to your app.

For simplicity reasons we set the default to the one that next-auth provides us. You'll be able to change the redirect later if you want to.

Hit Create, then in the detail view of your application make sure to enable dev mode. Dev mode ensures that you can start an auth flow from a non https endpoint for testing.

> Note that we get a clientId but no clientSecret because it is not needed for our authentication flow.

## Environment

Create a file `.env` in the root of the project and add the following keys to it.

```
NEXTAUTH_URL=http://localhost:3000
ZITADEL_CLIENT_ID=[yourClientId]
ZITADEL_ISSUER=https:/[your-domain]-[random-string].zitadel.cloud
```

# User interface

Now we can start editing the homepage by modifying `pages/index.tsx`.
Add this snippet your file. This code gets your auth session from next-auth, renders a Logout button if your authenticated or shows a Signup button if your not.
Note that signIn method requires the id of the provider we provided earlier, and provides a possibilty to add a callback url, Auth Next will redirect you to the specified route if logged in successfully.

```ts
import { signIn, signOut, useSession } from 'next-auth/client';

export default function Page() {
    const { data: session } = useSession();
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

To allow session state to be shared between pages - which improves performance, reduces network traffic and avoids component state changes while rendering - you can use the NextAuth.js Provider in `/pages/_app.tsx`.
Take a loot at the template `_app.tsx`.

```ts
import { SessionProvider } from "next-auth/react";

function MyApp({ Component, pageProps }) {
  return (
    <SessionProvider session={pageProps.session}>
      <Component {...pageProps} />
    </SessionProvider>
  );
}

export default MyApp;
```

Last thing: create a `profile.tsx` in /pages which renders the callback page.

```ts
import Link from "next/link";

import styles from "../styles/Home.module.css";

export default function Profile() {
  return (
    <div className={styles.container}>
      <h1>Login successful</h1>
      <Link href="/">
        <button>Back to Home</button>
      </Link>
    </div>
  );
}
```
