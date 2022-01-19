---
title: Angular
---

This integration guide shows you the recommended way to integrate ZITADEL into your Angular application.
It shows how to add user login to your application and fetch some data from the user info endpoint.  

At the end of the guide, your application has login functionality and has access to the current user's profile.  

> This documentation refers to our [example](https://github.com/caos/zitadel-examples/tree/main/angular) in GitHub. Note that we've written ZITADEL Console in Angular, so you can also use that as a reference.  

## Setup Application and Get Keys

Before we can start building our application, we have to do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your [Project](https://console.zitadel.ch/projects), then add a new application at the top of the page.
Select Web application type and continue.
We recommend you use [Authorization Code](../../apis/openidoauth/grant-types#authorization-code) in combination with [Proof Key for Code Exchange (PKCE)](../../apis/openidoauth/grant-types#proof-key-for-code-exchange) for all web applications.  

![Create app in console](/img/angular/app-create-light.png)

### Redirect URIs

With the Redirect URIs field, you tell ZITADEL where it is allowed to redirect users to after authentication. For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.  

> If you are following along with the [example](https://github.com/caos/zitadel-examples/tree/main/angular), set dev mode to `true` and the Redirect URIs to <http://localhost:4200/auth/callback>.  

If you want to redirect the users back to a route on your application after they have logged out, add an optional redirect in the Post Logout URIs field.  

Continue and create the application.

### Client ID and Secret

After successful app creation, a pop-up will appear, showing the app's client ID. Copy the client ID, as you will need it to configure your Angular client.

## Angular Setup

Now that you have your web application configured on the ZITADEL side, you can go ahead and integrate your Angular client.

### Install Angular Dependencies

You need to install an OAuth / OIDC client to connect with ZITADEL. Run the following command:

```bash
npm install angular-oauth2-oidc
```

### Create and Configure Auth Module

Add _OAuthModule_ to your Angular imports in _AppModule_ and provide the _AuthConfig_ in the providers' section. Also, ensure you import the _HTTPClientModule_.

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

Set _openid_, _profile_ and _email_ as scope, _code_ as responseType, and oidc to _true_. Then create an authentication service to provide the functions to authenticate your user.

You can use Angular’s schematics to do so:

```bash
ng g service services/authentication
```

Copy the following code to your service. This code provides a function `authenticate()` which redirects the user to ZITADEL. After successful login, ZITADEL redirects the user back to the redirect URI configured in _AuthModule_ and ZITADEL Console. Make sure both correspond, otherwise ZITADEL throws an error.

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

Our example includes a _StatehandlerService_ to redirect the user back to the route where he initially came from.
If you don't need such behavior, you can omit the following line from the `authenticate()` method above.

```ts
...
const newState = setState ? await this.statehandler.createState().toPromise() : undefined;
...
```

If you decide to use the _StatehandlerService_, provide it in the `app.module`. Make sure it gets initialized first using Angular’s `APP_INITIALIZER`. You find the service implementation in the [example](https://github.com/caos/zitadel-examples/tree/main/angular).

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

### Add Login to Your Application

To log a user in, you need a component or a guard.

- A component could provide a button, starting the login flow on click.

- A guard that starts a login flow once a user without a stored valid access token tries to access a protected route.

Using these components heavily depends on your application. In most cases, you need both.

Generate a component like this:

```bash
ng g component components/login
```

Inject the _AuthenticationService_ and call `authenticate()` on some click event.

Same for the guard:

```bash
ng g guard guards/auth
```

This code shows the _AuthGuard_ used in ZITADEL Console.

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

Add the guard to your _RouterModule_ similar to this:

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

### Add Logout to Your Application

Call `auth.signout()` for logging the current user out. Note that you can also configure a logout redirect URI if you want your users to be redirected after logout.

```ts
import { AuthenticationService } from 'src/app/services/authentication.service';

export class SomeComponentWithLogout {
    constructor(private authService: AuthenticationService){}

    public signout(): Promise<void> {
        return this.authService.signout();
    }
}
```

### Show User Information

To fetch user data, ZITADEL's user info endpoint has to be called. This data contains sensitive information and artifacts related to the current user's identity and the scopes you defined in your _AuthConfig_.
Our _AuthenticationService_ already includes a method called _getOIDCUser()_. You can call it wherever you need this information.

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

If you get stuck, consider checking out our [example](https://github.com/caos/zitadel-examples/tree/main/angular) application. It includes all the mentioned functionality of this quickstart. You can simply start by cloning the repository and replacing the _AuthConfig_ in the _AppModule_ by your own configuration. If you run into issues, contact us or raise an issue on [GitHub](https://github.com/caos/zitadel).

![App in console](/img/angular/app-screen.png)

### What's next?

Now that you have enabled authentication, it's time to add authorization to your application using ZITADEL APIs. Refer to the [docs](../../apis/introduction) or check out our ZITADEL Console code on [GitHub](https://github.com/caos/zitadel) which is using gRPC to access data.

For more information about creating an Angular application, refer to [Angular](https://angular.io/start) and for more information about the OAuth/OIDC library used above, consider reading their docs at [angular-oauth2-oidc](https://github.com/manfredsteyer/angular-oauth2-oidc).
