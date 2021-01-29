---
title:  Configure Zitadel
description: ...
---

### Setup Application and get Keys
We recommend that you use [Authorization Code](architecture#Authorization_Code) in combination with [Proof Key for Code Exchange](architecture#Proof_Key_for_Code_Exchange) for all web applications.
This flow has great support with most modern languages and frameworks and is the recommended default.

> In the OIDC and OAuth world this **client profile** is called "user-agent-based application"

Go to [ZITADEL Console Projects](https://console.zitadel.ch/projects), select your project, and click on `new application`. 
Enter a name for your new Application, select Web as Type, Basis as Auth Method, and then proceed.

### Redirect URLs

A redirect URL is a URL in your application where ZITADEL redirects the user after they have authenticated. Set your url, add an optional redirect after logout and then proceed and create.

> If you are following along with the sample project you downloaded from our Templates, you should set the Allowed Callback URL to http://localhost:4200.

Copy your client ID as it will be needed in the next step.

### Install Angular dependencies

You need to install an oauth client to connect with ZITADEL. Run the following command:

```bash
npm install angular-oauth2-oidc
```

This library helps integrating ZITADEL Authentication in your Angular Application.

#### Configure Auth Module

This example shows how to ...

```ts
import { UserManager, WebStorageStateStore, User } from 'oidc-client';

export default class AuthService {
    private userManager: UserManager;

    constructor() {
        const ZITADEL_ISSUER_DOMAIN: string = "https://issuer.zitadel.ch";

        const settings: any = {
            userStore: new WebStorageStateStore({ store: window.localStorage }),
            authority: ZITADEL_ISSUER_DOMAIN,
            client_id: 'YOUR_ZITADEL_CLIENT_ID',
            redirect_uri: 'http://localhost:44444/callback.html',
            response_type: 'code',
            scope: 'openid profile',
            post_logout_redirect_uri: 'http://localhost:44444/',
        };

        this.userManager = new UserManager(settings);
    }

    public getUser(): Promise<User | null> {
        return this.userManager.getUser();
    }

    public login(): Promise<void> {
        return this.userManager.signinRedirect();
    }

    public logout(): Promise<void> {
        return this.userManager.signoutRedirect();
    }

    public getAccessToken(): Promise<string> {
        return this.userManager.getUser().then((data: any) => {
            return data.access_token;
        });
    }
}
```

### Add Login in your application

The library gives you tools to quickly login to ZITADEL. To start a login attempt you need to create Login Component containing a button to start the authentication flow. `loginWithRedirect` redirects your user to ZITADEL.ch for authentication. Upon successfull Authentication, ZITADEL will redirect the user back to your previously defined Redirect URL. 

```ts
import { UserManager, WebStorageStateStore, User } from 'oidc-client';

export default class AuthService {
    private userManager: UserManager;

    constructor() {
        const ZITADEL_ISSUER_DOMAIN: string = "https://issuer.zitadel.ch";

        const settings: any = {
            userStore: new WebStorageStateStore({ store: window.localStorage }),
            authority: ZITADEL_ISSUER_DOMAIN,
            client_id: 'YOUR_ZITADEL_CLIENT_ID',
            redirect_uri: 'http://localhost:44444/callback.html',
            response_type: 'code',
            scope: 'openid profile',
            post_logout_redirect_uri: 'http://localhost:44444/',
        };

        this.userManager = new UserManager(settings);
    }

    public getUser(): Promise<User | null> {
        return this.userManager.getUser();
    }

    public login(): Promise<void> {
        return this.userManager.signinRedirect();
    }

    public logout(): Promise<void> {
        return this.userManager.signoutRedirect();
    }

    public getAccessToken(): Promise<string> {
        return this.userManager.getUser().then((data: any) => {
            return data.access_token;
        });
    }
}
```


### Add Logout in your application

The library also provider a useful function for logging out your users. Just call `auth.logout` to log out your user. Note that you can also configure your Logout Redirect URL if you want your Users to be redirected after logout.
```ts
import { UserManager, WebStorageStateStore, User } from 'oidc-client';

export default class AuthService {
    public logout(): Promise<void> {
        return this.userManager.signoutRedirect();
    }
}
```

### Show User Information

To fetch user data, ZITADELS user info endpoint has to be called. This data contains sensitive information and artifacts related to your users identity and the scopes you defined in your Auth Config. 
...

