---
title:  Single Page Application
description: ...
---

### SPA Protocol and Flow recommendation

If your [client](administrate#Clients) is a single page application (SPA) we recommend that you use [Authorization Code](documentation#Authorization_Code) in combination with [Proof Key for Code Exchange](documentation#Proof_Key_for_Code_Exchange).

This flow has great support with most modern languages and frameworks and is the recommended default.

> In the OIDC and OAuth world this **client profile** is called "user-agent-based application"

### Typescript Example

#### Typescript Authentication Example

If you use a framework like Angular, Vue, React, ... you can use this code snippet here to integrate **ZITADEL** into you application

Library used for this example [https://github.com/IdentityModel/oidc-client-js](https://github.com/IdentityModel/oidc-client-js)

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