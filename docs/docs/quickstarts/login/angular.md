---
title: Angular
---

This guide shows you how to integrate ZITADEL into your Angular application.

It covers how to:
- Add a user login to your application
- Fetch some data from the user info endpoint.

> This documentation refers to our [example](https://github.com/zitadel/zitadel-examples/tree/main/angular) in GitHub.
> Note that we've written the ZITADEL console in Angular.
> You can also use that as a reference.

## Set up application and get keys

Before you build your application, you'll need to head to the ZITADEL console and add some information about your application.
To start, we recommend creating a new app from scratch.
To do so:

1. Navigate to your [project](https://console.zitadel.ch/projects).
1. At the top of the page, add a new application.
1. Select **Web application type** and continue.

For all applications, we recommend combining
[Authorization Code](../../apis/openidoauth/grant-types#authorization-code)
with [Proof Key for Code Exchange (PKCE)](../../apis/openidoauth/grant-types#proof-key-for-code-exchange).

[Read more about the different app types](https://docs.zitadel.ch/docs/guides/authorization/oauth-recommended-flows#different-client-profiles).

![Create app in console](/img/angular/app-create-light.png)

### Add redirect URIs

In the Redirect URIs field, tell ZITADEL where to redirect users after authentication.
For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.

> If you are following along with the [example](https://github.com/zitadel/zitadel-examples/tree/main/angular), set dev mode to `true` and the Redirect URIs to <http://localhost:4200/auth/callback>.

After users log out, you can redirect them back to a route on your application.
To do so, add an optional redirect in the Post Logout URIs field.

Continue and create the application.

### Copy client ID

After you create your app, a pop-up will show the app's client ID.
Copy the client ID, as you will need it to configure your Angular client.

## Angular setup

Now that you have your web application configured on the ZITADEL side, you can integrate your Angular client.

### Install Angular dependencies

To connect with ZITADEL, you need to install an OAuth / OIDC client.
To do so, run this command:

```bash
npm install angular-oauth2-oidc
```

### Create and configure AuthModule

1. Add `OAuthModule` and `HTTPClientModule` to your Angular imports in `AppModule`

2. In the `AuthConfig` object, add the following values:
   * For `scope`, set `openid`, `profile` and `email`.
   * For `responseType`, use `code`
   * Set `oidc` to `true`.

3. Provide the `AuthConfig` object in the providers' section.


```ts
...
import { AuthConfig, OAuthModule } from 'angular-oauth2-oidc';
import { HttpClientModule } from '@angular/common/http';

const authConfig: AuthConfig = {
    scope: 'openid profile email',
    responseType: 'code',
    oidc: true,
    clientId: 'YOUR-CLIENT-ID', // replace with your appid
    issuer: 'https://issuer.zitadel.ch',
    redirectUri: 'http://localhost:4200/auth/callback',
    postLogoutRedirectUri: 'http://localhost:4200/signedout', // optional
    requireHttps: false // required for running locally
};

@NgModule({
...
    imports: [
        OAuthModule.forRoot(),
        HttpClientModule,
...
    providers: [
        {
            provide: AuthConfig,
            useValue: authConfig
        }
...
```

4. Create an authentication service to provide the functions to authenticate your user.
   You can use Angular’s schematics to do so:

   ```bash
   ng g service services/authentication
   ```

5. Copy the following code to your service.
   This code provides a function `authenticate()`, which redirects the user to ZITADEL.
   After successful login, ZITADEL redirects the user back to the redirect URI configured in `AuthModule` and the ZITADEL console.
   Make sure both correspond, otherwise ZITADEL throws an error.

```ts
import { Injectable } from '@angular/core';
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, from, Observable } from 'rxjs';

import { StatehandlerService } from './statehandler.service';

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {
    private _authenticated: boolean = false;
    private readonly _authenticationChanged: BehaviorSubject<
        boolean
    > = new BehaviorSubject(this.authenticated);

    constructor(
        private oauthService: OAuthService,
        private authConfig: AuthConfig,
        private statehandler: StatehandlerService,
    ) { }

    public get authenticated(): boolean {
        return this._authenticated;
    }

    public get authenticationChanged(): Observable<boolean> {
        return this._authenticationChanged;
    }

    public getOIDCUser(): Observable<any> {
        return from(this.oauthService.loadUserProfile());
    }

    public async authenticate(
        setState: boolean = true,
    ): Promise<boolean> {
        this.oauthService.configure(this.authConfig);

        this.oauthService.strictDiscoveryDocumentValidation = false;
        await this.oauthService.loadDiscoveryDocumentAndTryLogin();

        this._authenticated = this.oauthService.hasValidAccessToken();

        if (!this.oauthService.hasValidIdToken() || !this.authenticated) {
            const newState = setState ? await this.statehandler.createState().toPromise() : undefined;
            this.oauthService.initCodeFlow(newState);
        }
        this._authenticationChanged.next(this.authenticated);

        return this.authenticated;
    }

    public signout(): void {
        this.oauthService.logOut();
        this._authenticated = false;
        this._authenticationChanged.next(false);
    }
}
```

Our example includes a `StatehandlerService` to redirect the users back to the route where they started.
If you don't need such behavior, you can omit the following line from the `authenticate()` method above.

```ts
...
const newState = setState ? await this.statehandler.createState().toPromise() : undefined;
...
```

If you decide to use the `StatehandlerService`, provide it in the `AppModule`.
Make sure it gets initialized first using Angular’s `APP_INITIALIZER`.
You can find the service implementation in the [example](https://github.com/zitadel/zitadel-examples/tree/main/angular).

```ts

const stateHandlerFn = (stateHandler: StatehandlerService) => {
    return () => {
        return stateHandler.initStateHandler();
    };
};

...
providers: [
        {
            provide: APP_INITIALIZER,
            useFactory: stateHandlerFn,
            multi: true,
            deps: [StatehandlerService],
        },
        {
            provide: StatehandlerProcessorService,
            useClass: StatehandlerProcessorServiceImpl,
        },
        {
            provide: StatehandlerService,
            useClass: StatehandlerServiceImpl,
        },
]
...
```

### Add login to your application

To log a user in, you probably need a _component_ or _guard_.

- A component could provide a button, starting the login flow on click.

- A guard starts a login flow when a user tries to access a protected route
  without a stored access token.

How you use these components depends on your application.
In most cases, you need both.

Generate a component like this:

```bash
ng g component components/login
```

Inject the `AuthenticationService` and call `authenticate()` on some click event.

Same for the guard:

```bash
ng g guard guards/auth
```

This code shows the `AuthGuard` used in the ZITADEL console.

```ts
import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, UrlTree } from '@angular/router';
import { Observable } from 'rxjs';
import { AuthenticationService } from '../services/authentication.service';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate {

  constructor(private auth: AuthenticationService) { }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
        if (!this.auth.authenticated) {
            return this.auth.authenticate();
        }
        return this.auth.authenticated;
    }
}
```

Add the guard to your `RouterModule` similar to this:

```ts
...
const routes: Routes = [
    {
        path: '',
        loadChildren: () => import('./pages/home/home.module').then(m => m.HomeModule),
        canActivate: [AuthGuard],
    },
...
```

> Note: Make sure you redirect the user from your callback URL to a guarded page, so `authenticate()` is called again and the access token is stored.

```ts
...
    {
        path: 'auth/callback',
        redirectTo: 'user',
    },
....
```

### Add logout to your application

To log the current user out, call `auth.signout()`.

To redirect your users after logout, configure a logout redirect URI.

```ts
import { AuthenticationService } from 'src/app/services/authentication.service';

export class SomeComponentWithLogout {
    constructor(private authService: AuthenticationService){}

    public signout(): Promise<void> {
        return this.authService.signout();
    }
}
```

### Show user information

To fetch user data, call ZITADEL's user info endpoint.
This data contains sensitive information and artifacts about the current user's identity, and the scopes you defined in your `AuthConfig`.

Our `AuthenticationService` already includes a method called `getOIDCUser()`.
You can call it wherever you need this information.

```ts
import { AuthenticationService } from 'src/app/services/authentication.service';

public user$: Observable<any>;

constructor(private auth: AuthenticationService) {
    this.user$ = this.auth.getOIDCUser();
}
```

And in your HTML file:

```html
<div *ngIf="user$ | async as user">
    <p>{{user | json}}</p>
</div>
```

## Completion

You have successfully integrated your Angular application with ZITADEL!

If you get stuck, check out our [example](https://github.com/zitadel/zitadel-examples/tree/main/angular) application.
It includes all the mentioned functionality of this quickstart.
You can start by cloning the repository and replacing the `AuthConfig` in the `AppModule` with your own configuration.

If you run into issues, contact us or raise an issue on [GitHub](https://github.com/zitadel/zitadel).

![App in console](/img/angular/app-screen.png)

### What's next?

Now that you have enabled authentication, it's time to add authorization to your application using ZITADEL APIs.
Refer to the [docs](../../apis/introduction) or check out our ZITADEL console code on [GitHub](https://github.com/zitadel/zitadel), which uses gRPC to access data.

For more information about creating an Angular application, refer to [Angular](https://angular.io/start) and for more information about the OAuth/OIDC library used above, consider reading their docs at [angular-oauth2-oidc](https://github.com/manfredsteyer/angular-oauth2-oidc).
