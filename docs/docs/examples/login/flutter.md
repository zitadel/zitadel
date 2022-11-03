---
title: Flutter
---

This guide demonstrates how you integrate **ZITADEL** as an idendity provider to a Flutter app.

At the end of the guide you have a mobile application for **Android**, **iOS** and **Web** with the ability to authenticate users via ZITADEL.

If you need any other information about Flutter, head over to the [documentation](https://flutter.dev/).

## Setup Application

Before we can start building our application, we have to do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your Project, then add a new application at the top of the page.
Select **Native** application type and continue.

![Create app in console](/img/angular/app-create.png)

### Redirect URIs

With the Redirect URIs field, you tell ZITADEL where it is allowed to redirect users to after authentication. For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.

As our application will also support web, we have to make sure to set redirects for http and https, as well as a **custom-scheme** for our native Android and IOS Setup.

For local development, add a redirectURI for `http://localhost:4444/auth.html` and your custom scheme. In our case it is `com.example.zitadelflutter`.

If you want to redirect the users back to a route on your application after they have logged out, add an optional redirect in the Post Logout URIs field.

Continue and create the application.

After creation, go to **token settings** and check the refresh token checkbox. This allows us to request a refresh_token via `offline_access` scope.
Make sure to save the application.

### Client ID

After successful app creation, a pop-up will appear, showing the app's client ID. Copy the client ID, as you will need it to configure your Angular client.

## Flutter Prerequisites

To move further in this quickstart, you'll need the following things prepared:

- Have Flutter (and Dart) installed ([how-to](https://flutter.dev/docs/get-started/install))
- Have an IDE set up for developing Flutter ([how-to](https://flutter.dev/docs/get-started/editor))
- Create a basic Flutter app ([how-to](https://flutter.dev/docs/get-started/codelab))
- Create a "Native" application in a ZITADEL project

## Flutter with ZITADEL

In your native application on ZITADEL, you need to add a callback (redirect) uri
which matches the selected custom url scheme. As an example, if you intend to
use `ch.myexample.app://sign-me-in` as redirect URI on ZITADEL and in your app,
you need to register the `ch.myexample.app://` custom url scheme within Android and iOS.

:::caution Use Custom Redirect URI!

You'll need the custom redirect url to be compliant with the OAuth 2.0
authentication for mobile devices ([RFC 8252 specification](https://tools.ietf.org/html/rfc8252)).
Otherwise your app might get rejected.

:::

### Hello World

After you created the starter Flutter app, the app will show a simple, templated Flutter app.

### Install Dependencies

To authenticate users with ZITADEL in a mobile application, some specific packages are needed.
The [RFC 8252 specification](https://tools.ietf.org/html/rfc8252) defines how
[OAUTH2.0 for mobile and native apps](https://oauth.net/2/native-apps/) works.
Basically, there are two major points in this specification:

1. It recommends to use [PKCE](https://oauth.net/2/pkce/)
2. It does not allow third party apps to open the browser for the login process,
   the app must open the login page within the embedded browser view

First install [http](https://pub.dev/packages/http) a library for making HTTP calls,
then [`flutter_web_auth_2`](https://pub.dev/packages/flutter_web_auth_2) package and a secure storage to store the auth / refresh tokens [flutter_secure_storage](https://pub.dev/packages/flutter_secure_storage).

To install run:

```bash
flutter pub add http
flutter pub add flutter_web_auth_2
flutter pub add flutter_secure_storage
```

#### Setup for Android

Navigate to your `AndroidManifest.xml` at `<projectRoot>/android/app/src/main/AndroidManifest.xml` and add the following activity with your scheme.

```xml reference
https://github.com/zitadel/zitadel_flutter/blob/main/android/app/src/main/AndroidManifest.xml#L29-L38
```

Furthermore, for `secure_storage`, you need to set the minimum SDK version to 18
in `<projectRoot>/android/app/src/build.gradle`.

### Add Authentication

:::note

The auth redirect scheme "`ch.myexample.app`" does register all auth urls with the given
scheme for the app. So an url pointing to `ch.myexample.app://signin` and another one
for `ch.myexample.app://logout` will work with the same registration.

:::

To reduce the commented default code, we will modify the `main.dart` file.

First, the `MyApp` class: it remains a stateless widget:

```dart reference
https://github.com/zitadel/zitadel_flutter/blob/main/lib/main.dart#L14-L28
```

Second, the `MyHomePage` class will remain a stateful widget with
its title, we don't change any code here.

```dart reference
https://github.com/zitadel/zitadel_flutter/blob/main/lib/main.dart#L30-L37
```

What we'll change now, is the `_MyHomePageState` class to enable
authentication via ZITADEL and remove the counter button of the starter application. We'll show the username of the authenticated user.

We define the needed elements for our state:

```dart
var _busy = false;
var _authenticated = false;
var _username = '';
final storage = const FlutterSecureStorage();
```

Then the builder method, which does show the login button if you're not
authenticated, a loading bar if the login process is going on and
your name if you are authenticated:

```dart reference
https://github.com/zitadel/zitadel_flutter/blob/main/lib/main.dart#L119-L159
```

And finally the `_authenticate` method which calls the authorization endpoint,
then fetches the user info and stores the tokens into the secure storage.

```dart reference
https://github.com/zitadel/zitadel_flutter/blob/main/lib/main.dart#L45-L117
```

Now, you can run your application for IOS and Android devices with

```bash
flutter run
```

or by directly selecting your device

```bash
flutter run -d iphone
```

for Web make sure you run the application on a fixed port such that you can setup your redirect URI in ZITADEL console for testing.

```bash
flutter run -d chrome --web-port=4444
```

### Result

If everything works out correctly, you should see you applications like this:

<div style={{display: 'grid', 'grid-column-gap': '1rem', 'grid-template-columns': '1fr 1fr', 'max-width': '500px', 'margin': '0 auto'}}>
    <img src="/img/flutter/not-authed.png" alt="Unauthenticated" height="500px" />
    <img src="/img/flutter/authed.png" alt="Flutter Authenticated" height="500px" />
</div>

<div style={{display: 'grid', 'grid-column-gap': '1rem', 'grid-template-columns': '1fr 1fr', 'max-width': '800px', 'margin': '0 auto'}}>
    <img src="/img/flutter/web-not-authed.png" alt="Unauthenticated" height="500px" />
    <img src="/img/flutter/web-authed.png" alt="Flutter Authenticated" height="500px" />
</div>

Now the next step is ensuring our access tokens are valid on a next startup, and refreshing it via refresh_token if its not.
