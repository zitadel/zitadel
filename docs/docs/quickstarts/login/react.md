---
title: React
---

This guide shows you how to integrate ZITADEL with into your React application.

It covers how to:
- Add a user login to your application
- Fetch some data from the user info endpoint.

## Setup Application and get Keys

Before you build your application, you'll need to do a few configuration steps in the ZITADEL Console.
You will need to provide some information about your app.
We recommend creating a new app to start from scratch.
To do so:

1. Navigate to your [Project](https://console.zitadel.ch/projects).
2. At the top of the page, add a new application.
3. Select User Agent and continue.

For all applications, we recommend combining an
[Authorization Code](../../apis/openidoauth/grant-types#authorization-code)
with a [Proof Key for Code Exchange (PKCE)](../../apis/openidoauth/grant-types#proof-key-for-code-exchange).

[Read more about the different app types](https://docs.zitadel.ch/docs/guides/authorization/oauth-recommended-flows#different-client-profiles).

### Redirect URLs

After users authenticate, ZITADEL redirects them to a redirect URL.

Set your redirect URL to the domain where the app will be deployed.
You can also use the npm default `http://localhost:3000/`.

> If redirecting to `localhost`, set dev mode to `true`.
> This will enable unsecure http for the moment.

After users log out, you can redirect them back to a route on your application.
To configure this, add an optional redirect in the Post Logout URIs field.

**Continue** and **Create** the application.

### Client ID

After you create your app, a popup shows you your clientID and secret.
Copy your client ID.
You'll use it in the next step.

## React Setup

### Create React app

Create a new React app with the following command:

```bash
npx create-react-app my-app
```

### Install an oidc client

You need to install an oauth / oidc client to connect with ZITADEL. Run the following command:

```bash
npm install oidc-react
```

This library helps integrate ZITADEL Authentication into your React Application.

### Create and configure Auth Module

With the installed oidc pakage, you will need an AuthProvider.
This should contain the OIDC configuration.

The oidc configuration should have the following values:
   * For `scope`, set `openid`, `profile` and `email`.
   * For `responseType`, use `code`

In the following code, the authority is already set to the issuer of zitadel.ch.
You can find this in the ZITADEL Console for your application.
Replace the clientId value `<YOUR-CLIENT-ID>` with the generated client ID of your application in ZITADEL Console.


```ts

import React from 'react';
import { AuthProvider } from 'oidc-react';
import './App.css';
const oidcConfig = {
    onSignIn: async (response: any) => {
        alert('You logged in :' + response.profile.given_name + ' ' + response.profile.family_name);
        window.location.hash = '';
    },
    authority: 'https://issuer.zitadel.ch',
    clientId:
        'YOUR-CLIENT-ID',
    responseType: 'code',
    redirectUri: 'http://localhost:3000/',
    scope: 'openid profile email'
};

function App() {
    return (
        <AuthProvider {...oidcConfig}>
        <div className="App">
        <header className="App-header">
            <p>Hello World</p>
    </header>
    </div>
    </AuthProvider>
);
}

export default App;
```

### Run application

Start your React application with the following command:

```bash
npm start
```

Your browser should automatically open the app site.
You can also just go to `http://localhost:3000/`.

When you open the app in your browser, you will be redirected to the zitadel.ch login.
After successfully authenticating your user, you will return to your application.
It should show a popup that says: **You logged in {FirstName} {LastName}**

## Completion

You have successfully integrated ZITADEL in your React Application!

### What's next?

Now you can proceed implementing our APIs to include Authorization. You can find our API Docs [here](../../apis/introduction)

For more information about creating a React application we refer to [React](https://reactjs.org/docs/getting-started.html) and for more information about the used oauth/oidc library consider reading their docs at [oidc-react](https://www.npmjs.com/package/oidc-react).
