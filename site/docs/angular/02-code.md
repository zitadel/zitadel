---
title: Angular Setup
---

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

``` bash
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

To create a component use 
``` bash
ng g component components/login
```
and then inject the authService to call `authenticate()`.

Same for the guard:
``` bash
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

