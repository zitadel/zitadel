---
title: ZITADEL with Go
sidebar_label: Go
---

This integration guide demonstrates the recommended way to incorporate ZITADEL into your Go web application. 
It explains how to enable user login in your application and how to fetch data from the user info endpoint.

By the end of this guide, your application will have login functionality and will be able to access the current user's profile.

> This documentation references our [example](https://github.com/zitadel/zitadel-go) on GitHub. 
> You can either create your own application or directly run the example by providing the necessary arguments.

## Set up application

Before we begin developing our application, we need to perform a few configuration steps in the ZITADEL Console.
You'll need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your Project, then add a new application at the top of the page.
Select the **Web** application type and continue.

![Create app in console](/img/go/app-create.png)

We recommend that you use [Proof Key for Code Exchange (PKCE)](/apis/openidoauth/grant-types#proof-key-for-code-exchange) for all applications.

![Create app in console - set auth method](/img/go/app-create-auth.png)

### Redirect URIs

The Redirect URIs field tells ZITADEL where it's allowed to redirect users after authentication. For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.
The Post-logout redirect send the users back to a route on your application after they have logged out.

> If you are following along with the [example](https://github.com/zitadel/zitadel-go), set the dev mode to `true`, the Redirect URIs to `http://localhost:8089/auth/callback` and Post-logout redirect URI to [http://localhost:8089/](http://localhost:8089/)>.

![Create app in console - set redirectURI](/img/go/app-create-redirect.png)

Continue and create the application.

### Client ID

After successful creation of the app, a pop-up will appear displaying the app's client ID. Copy the client ID, as you will need it to configure your Go client.

![Create app in console - copy client_id](/img/go/app-create-clientid.png)

## Go setup

Now that you have configured your web application on the ZITADEL side, you can proceed with the integration of your Go client.

### Install ZITADEL Go SDK

To connect with ZITADEL, you need to install an OAuth/OIDC client. Run the following command:

```bash
go get -u github.com/zitadel/zitadel-go/v3
```

### Create the application server

Create a new go file with the content below. This will create an application with a home and profile page.

```go reference
https://github.com/zitadel/zitadel-go/blob/next/example/app/app.go
```

This will basically set up everything. So let's look at some parts of the code.

**Register authentication handler**:

For the authentication to work, the SDK needs some handlers in your application.
In this example we will register them on the `/auth/` prefix.
The SDK itself will then register three routes on that to be able to:
 - start the authentication process and redirect to the Login UI (`/auth/login`)
 - continue with the authentication process after the login UI (`/auth/callback`)
 - terminate the session (`/auth/logout`)

```go
router.Handle("/auth/", z.Authentication)
```

***Authentication checks***

To ensure the user is authenticated before they are able to use your application, the middleware provides two options:
- You can either require the user to be authenticated. If they haven't already, they will be automatically redirected to the Login UI:
    ```go
    mw.RequireAuthentication()(handler)
    ```
- You can just check the user's authentication status, but still continue serving the page:
    ```go
    mw.CheckAuthentication()(handler)
    ```
  
***Authentication context***

If you used either of the authentication checks above, you can then access context information in your handler:
```go
mw.Context(req.Context())
```

### Add pages to your application

To be able to serve these pages create a `templates` directory in the same folder as you just created the go file.
Now create two HTML files in the new `templates` folder and copy the content of the examples:

**home.html**

The home page will display a short welcome message and allow the user to manually start the login process.

```go reference
https://github.com/zitadel/zitadel-go/blob/next/example/app/templates/home.html
```

**profile.html**

The profile page will display the Userinfo from the authentication context and allow the user to logout.

```go reference
https://github.com/zitadel/zitadel-go/blob/next/example/app/templates/profile.html
```

### Start your application

You will need to provide some values for the program to run:
- `domain`: Your ZITADEL instance domain, e.g. my-domain.zitadel.cloud
- `key`: Random secret string. Used for symmetric encryption of state parameters, cookies and PCKE. 
- `clientID`: The clientID provided by ZITADEL
- `redirectURI`: The redirectURI registered at ZITADEL
- `port`: The port on which the API will be accessible, default it 8089

```bash
go run main.go --domain <your domain> --key <key> -- clientID <clientID> --redirectURI <redirectURI>
```

This could look like:

```bash
go run main.go --domain my-domain.zitadel.cloud --key XKv2Lqd7YAq13NUZVUWZEWZeruqyzViM --clientID 243861220627644836@example --redirectURI http://localhost:8089/auth/callback
```

If you then visit on http://localhost:8089 you should get the following screen:

![Home Page](/img/go/app-home.png)

By clicking on `Login` you will be redirected to your ZITADEL instance. After login with your existing user you will be presented the profile page:

![Profile Page](/img/go/app-profile.png)

## Completion

Congratulations! You have successfully integrated your Go application with ZITADEL!

If you get stuck, consider checking out our [example](https://github.com/zitadel/zitadel-go) application. 
This application includes all the functionalities mentioned in this quickstart. 
You can directly start it with your own configuration. If you face issues, contact us or raise an issue on [GitHub](https://github.com/zitadel/zitadel-go/issues).

