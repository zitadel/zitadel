---
title: Angular
---

This integration guide demonstrates the recommended way to incorporate ZITADEL into your Angular application. 
It explains how to enable user login in your application and how to fetch data from the user info endpoint.

By the end of this guide, your application will have login functionality and will be able to access the current user's profile.

> This documentation references our [example](https://github.com/zitadel/zitadel-angular) on GitHub. Please note that we wrote the ZITADEL Console in Angular, so you can also use that as a reference.

## Set up application and obtain keys

Before we begin developing our application, we need to perform a few configuration steps in the ZITADEL Console.
You'll need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your Project, then add a new application at the top of the page.
Select the **User Agent** application type and continue.
We recommend that you use [Proof Key for Code Exchange (PKCE)](/apis/openidoauth/grant-types#proof-key-for-code-exchange) for all SPA applications.

![Create app in console](/img/angular/app-create.png)


### Redirect URIs

The Redirect URIs field tells ZITADEL where it's allowed to redirect users after authentication. For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.
The Post-logout redirect send the users back to a route on your application after they have logged out.

> If you are following along with the [example](https://github.com/zitadel/zitadel-angular), set the dev mode to `true`, the Redirect URIs to <http://localhost:4200/auth/callback> and Post redirect URI to <http://localhost:4200/signedout>.

Continue and create the application.

### Client ID

After successful creation of the app, a pop-up will appear displaying the app's client ID. Copy the client ID, as you will need it to configure your Angular client.

## Angular setup

Now that you have configured your web application on the ZITADEL side, you can proceed with the integration of your Angular client.

### Install Angular dependencies

To connect with ZITADEL, you need to install an OAuth/OIDC client. Run the following command:

```bash
npm install angular-oauth2-oidc
```

### Create and configure the auth module

Add _OAuthModule_ to your Angular imports in _AppModule_ and provide the _AuthConfig_ in the providers' section. Also, ensure that you import the _HTTPClientModule_.

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/app.module.ts
```

Set _openid_, _profile_ and _email_ as scope, _code_ as responseType, and oidc to _true_. Then create an authentication service to provide the functions to authenticate your user.

You can use Angular’s schematics to do so:

```bash
ng g service services/authentication
```

Copy the following code to your service. This code provides a function `authenticate()`, which redirects the user to ZITADEL. After a successful login, ZITADEL redirects the user back to the redirect URI configured in _AuthModule_ and ZITADEL Console. Ensure that both correspond, otherwise ZITADEL will throw an error.

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/services/authentication.service.ts
```

Our example includes a _StatehandlerService_ that redirects the user back to the route from which they initially came.If you don't need such behavior, you can omit the following line from the `authenticate()` method above.

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/services/authentication.service.ts#L45
```
If you decide to use the _StatehandlerService_, include it in the `app.module`. Ensure it gets initialized first using Angular’s `APP_INITIALIZER`. You can find the service implementation in the [example](https://github.com/zitadel/zitadel-angular).

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/app.module.ts#L26-L30
```

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/app.module.ts#L55-L78
```

### Add login to your application

To log in a user, you need a component or a guard.

- A component could provide a button that initiates the login flow when clicked.
- A guard initiates a login flow when a user without a stored valid access token attempts to access a protected route.

The use of these components depends heavily on your application. In most cases, you need both.

Generate a component like this:

```bash
ng g component components/login
```

Inject the _AuthenticationService_ and call `authenticate()` on some click event.

Do the same for the guard:
```bash
ng g guard guards/auth
```

This code shows the _AuthGuard_ used in ZITADEL Console.

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/guards/auth.guard.ts
```

Add the guard to your _RouterModule_ similar to this:

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/app-routing.module.ts#L9-L31
```

> Note: Make sure you redirect the user from your callback URL to a guarded page, so the `authenticate()` method is called again, and the access token is stored.

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/app-routing.module.ts#L19-L21
```

### Add logout to your application

Call `auth.signout()` to log out the current user. Keep in mind that you can also configure a logout redirect URI if you want your users to be redirected after logout.

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/components/user/user.component.ts#L20-L22
```

### Display user information

To fetch user data, you need to call the ZITADEL's user info endpoint. This data contains sensitive information and artifacts related to the current user's identity and the scopes you defined in your _AuthConfig_.
Our _AuthenticationService_ already includes a method called _getOIDCUser()_. You can call it wherever you need this information.

```ts reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/components/user/user.component.ts
```

And in your HTML file:

```html reference
https://github.com/zitadel/zitadel-angular/blob/main/src/app/components/user/user.component.html
```

### Refresh token

If you want to add a refresh token to your application, navigate to the console application and check the box in the configuration section.
Then add `offline_access` to the scopes and add the following line:

```
this.oauthService.setupAutomaticSilentRefresh();
```

This line automatically refreshes a token before it expires.


## Completion

Congratulations! You have successfully integrated your Angular application with ZITADEL!

If you get stuck, consider checking out our [example](https://github.com/zitadel/zitadel-angular) application. This application includes all the funcationalities mentioned in this quickstart. You can start by cloning the repository and replacing the _AuthConfig_ in the _AppModule_ with your own configuration. If you face issues, contact us or raise an issue on [GitHub](https://github.com/zitadel/zitadel).

![App in console](/img/angular/app-screen.png)

### What's next?

Now that you have enabled authentication, it's time for you to add authorization to your application using ZITADEL APIs. To do this, you can refer to the [docs](/apis/introduction) or check out the ZITADEL Console code on [GitHub](https://github.com/zitadel/zitadel) which uses gRPC to access data.

For more information on how to create an Angular application, you can refer to [Angular](https://angular.io/start). If you want to learn more about the OAuth/OIDC library used above, consider reading the docs at [angular-oauth2-oidc](https://github.com/manfredsteyer/angular-oauth2-oidc).