---
title: React
---

This Integration guide shows you the recommended way to integrate **ZITADEL** into your React Application.
It demonstrates how to add a user login to your application and fetch some data from the user info endpoint.

At the end of the guide you should have an application able to login a user and read the user profile.

## Setup Application and get Keys

Before we can start building our application we have to do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your [Project](https://console.zitadel.ch/projects) and add a new application at the top of the page.
Select User Agent and continue. More about the different app types you can finde [here](https://docs.zitadel.ch/docs/guides/usage/oauth-recommended-flows#different-client-profiles)
We recommend that you use [Authorization Code](../apis/openidoauth/grant-types#authorization-code) in combination with [Proof Key for Code Exchange](../apis/openidoauth/grant-types#proof-key-for-code-exchange) for all web applications.

### Redirect URLs

A redirect URL is a URL in your application where ZITADEL redirects the user after they have authenticated. Set your url to the domain the app will be deployed to or use `http://localhost:3000/` because this will be the default of npm.

> You should set dev mode to `true` as this will enable unsecure http for the moment.

If you want to redirect the users back to a route on your application after they have logged out, add an optional redirect in the post redirectURI field.

Continue and Create the application.

### Client ID

After successful app creation a popup will appear showing you your clientID.
Copy your client ID as it will be needed in the next step.

## React Setup

### Create React app

Create a new React app with the following command

```bash
npx create-react-app my-app
```

### Install oidc client

You need to install an oauth / oidc client to connect with ZITADEL. Run the following command:

```bash
npm install oidc-react
```

This library helps integrating ZITADEL Authentication in your React Application.

### Create and configure Auth Module

With the installed oidc pakage you will need an AuthProvider which should contain the OIDC configuration.

The oidc configuration should contain **openid**, **profile** and **email** as scope and **code** as responseType.
In the code below the authority is already set to the issuer of zitadel.ch you can find this in the ZITADEL Console on you application.
Replace the clientId value 'YOUR-CLIENT-ID' with the generated client id of you application in ZITADEL Console.


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

Start your react application with the following command

```bash
npm start
```

Your browser should automatically open the app site or just go to `http://localhost:3000/`.
On opening the app in the browser you will be redirected to the login of zitadel.ch
After successfully authenticating your user, you will get back to you application.
It should show a popup which says: **You logged in {FirstName} {LastName}**

## Completion

You have successfully integrated ZITADEL in your React Application!

### Whats next?

Now you can proceed implementing our APIs to include Authorization. You can findour API Docs [here](../apis/apis)

For more information about creating a React application we refer to [React](https://reactjs.org/docs/getting-started.html) and for more information about the used oauth/oidc library consider reading their docs at [oidc-react](https://www.npmjs.com/package/oidc-react).
