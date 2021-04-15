---
title: Angular
---

This Integration guide shows you the recommended way to integrate **ZITADEL** into your Angular Application.
It demonstrates how to add a user login to your application and fetch some data from the user info endpoint.

At the end of the guide you should have an application able to login a user and read the user profile.

> This documentation refers to our [Template](https://github.com/caos/zitadel-angular-template) in Github. Note that our **ZITADEL Console** is also written in Angular and can therefore be used as a reference.

## Setup Application and get Keys

Before we can start building our application we have do do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your [Project](https://console.zitadel.ch/projects) and add a new application at the top of the page.
Select Web Application and continue.
We recommend that you use [Authorization Code](../apis/openidoauth/grant-types#authorization-code) in combination with [Proof Key for Code Exchange](../apis/openidoauth/grant-types#proof-key-for-code-exchange) for all web applications.

![Create app in console](/img/angular/app-create-light.png)

### Redirect URLs

A redirect URL is a URL in your application where ZITADEL redirects the user after they have authenticated. Set your url to the domain the web app will be deployed to or use `localhost:4200` for development as Angular will be running on port 4200.

> If you are following along with the [sample](https://github.com/caos/zitadel-angular-template) project you downloaded from our templates, you should set the Allowed Callback URL to <http://localhost:4200/auth/callback>. You will also have to set dev mode to `true` as this will enable unsecure http for the moment.

If you want to redirect the users back to a route on your application after they have logged out, add an optional redirect in the post redirectURI field.

Continue and Create the application.

### Client ID and Secret

After successful app creation a popup will appear showing you your clientID as well as a secret.
Copy your client ID as it will be needed in the next step.

## Angular Setup

### Install Angular dependencies

You need to install an oauth / oidc client to connect with ZITADEL. Run the following command:

```bash
npm install angular-oauth2-oidc
```

This library helps integrating ZITADEL Authentication in your Angular Application.

### Create and configure Auth Module

Add the Auth module to your Angular imports in AppModule and setup the AuthConfig in a constant above.

```ts
...
import { AuthConfig, OAuthModule } from 'angular-oauth2-oidc';

const authConfig: AuthConfig = {
    scope: 'openid profile email',
    responseType: 'code',
    oidc: true,
    clientId: 'YOUR-CLIENT-ID', // replace with your appid
    dummyClientSecret: 'YOUR-SECRET', // required by library
    issuer: 'https://issuer.zitadel.ch',
    redirectUri: 'http://localhost:4200/auth/callback',
    postLogoutRedirectUri: 'http://localhost:4200/signedout', // optional
    requireHttps: false // required for running locally
};

@NgModule({
    declarations: [
        AppComponent,
        SignedoutComponent,
    ],
    imports: [
        OAuthModule..forRoot(),
...
```

Set **openid**, **profile** and **email** as scope, **code** as responseType, and oidc to **true**.
Then create a Authentication Service to provide the functions to authenticate your user.

You can use Angulars schematics to do so:

```bash
ng g component services/authentication
```

This will create an AuthenticationService automatically for you.

Copy the following code to your service. This code provides a function `authenticate()` which redirects the user to ZITADEL. After the user has logged in it will be redirected back to your redirectURI set in Auth Module and Console. Make sure both correspond, otherwise ZITADEL will throw an error.

```ts
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';

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
        console.log('auth');
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

Our template includes a statehandler service to redirect the user back to the route where he initially came from. It saves the route information to a unique id so that the user can be redirected back after a successful authentication.
If you don't need such a behaviour you can escape the following lines from the `authenticate()` method above.

```ts
...
const newState = setState ? await this.statehandler.createState().toPromise() : undefined;
...
```

If you decide to use it provide the service in the `app.module` and make sure it gets initialized first using angulars `APP_INITIALIZER`.

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

### Add Login in your application

To login a user, a component or a guard is needed.

- A component provides a button prompting the user to start the login flow.
`authenticate()` redirects your user to ZITADEL.ch for authentication. Upon successfull Authentication, ZITADEL will redirect the user back to your previously defined Redirect URL.

- A guard can be setup to check if the user has a valid **Access Token** to proceed. This will check if the user has a stored **accesstoken** in storage or otherwise prompt the user to login.

The use of this components totally depends on your application. In most cases you need both.

To create a component use:

```bash
ng g component components/login
```

and then inject the authService to call `authenticate()`.

Same for the guard:

```bash
ng g guard guards/auth
```

This code shows the AuthGuard used in our Console.

```ts
import { AuthService } from 'src/app/services/auth.service';

@Injectable({
    providedIn: 'root',
})
export class AuthGuard implements CanActivate {
    constructor(private auth: AuthService) { }

    public canActivate(
        _: ActivatedRouteSnapshot,
        state: RouterStateSnapshot,
    ): Observable<boolean> | Promise<boolean> | boolean {
        if (!this.auth.authenticated) {
            return this.auth.authenticate();
        }
        return this.auth.authenticated;
    }
}
```

it can easily be added to your RouterModule.

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

> Note: To complete the code flow, `authenticate()` needs to be called twice. You may have to add a guard to your callback url to make sure it will complete the flow.

```ts
    {
        path: 'auth/callback',
        canActivate: [AuthGuard],
        redirectTo: 'user',
    },
```

### Add Logout in your application

The authService and Library also provides a useful function for logging out your users. Just call `auth.signout()` to log out your user. Note that you can also configure your Logout Redirect URL if you want your Users to be redirected after logout.

```ts
import { AuthService } from 'src/app/services/auth.service.ts';

export class SomeComponentWithLogout {
    constructor(private authService: AuthService){}

    public signout(): Promise<void> {
        return this.authService.signout();
    }
}
```

### Show User Information

To fetch user data, ZITADELS user info endpoint has to be called. This data contains sensitive information and artifacts related to your users identity and the scopes you defined in your Auth Config.
Our AuthService already includes a function called getOIDCUser(). You can call it whereever you need this information.

```ts
import { AuthenticationService } from 'src/app/services/auth.service.ts';

public user$: Observable<any>;

constructor(private auth: AuthenticationService) {
    this.user$ = this.auth.getOIDCUser();
}
```

and in your html

```html
<div *ngIf="user$ | async as user">
    <p>{{user | json}}</p>
</div>
```

## Completion

You have successfully integrated ZITADEL in your Angular Application!

If you get stuck consider checking out our [template](https://github.com/caos/zitadel-angular-template) application which includes all the mentioned functionality of this quickstart. You can simply start by cloning the repo and replacing the AuthConfig in the app.module with your own configuration. If your run into issues don't hesitate to contact us or raise an issue on [Github](https://github.com/caos/zitadel).

![App in console](/img/angular/app-screen.png)

### Whats next?

Now you can proceed implementing our APIs to include Authorization. Refer to our [Docs](introduction) or checkout our Console Code on [Github](https://github.com/caos/zitadel) which is using GRPC to access data.

For more information about creating an angular application we refer to [Angular](https://angular.io/start) and for more information about the used oauth/oidc library consider reading their docs at [angular-oauth2-oidc](https://github.com/manfredsteyer/angular-oauth2-oidc).
