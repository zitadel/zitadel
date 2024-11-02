---
title: ZITADEL with Flutter
sidebar_label: Flutter

---

This guide demonstrates how you integrate **ZITADEL** into a Flutter app. It refers to our example on [GitHub](https://github.com/zitadel/zitadel_flutter)

At the end of the guide you have a mobile application for **Android**, **iOS** and **Web** with the ability to authenticate users via ZITADEL.

If you need any other information about Flutter, head over to the [documentation](https://flutter.dev/).

## Setup Application

Before we can start building our application, we have to do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your Project, then add a new application at the top of the page.
Select **Native** application type and continue.

![Create app in console](/img/flutter/nativeapp.png)

### Redirect URIs

With the Redirect URIs field, you tell ZITADEL where it is allowed to redirect users to after authentication.
As our application will also support web, we have to make sure to set redirects for http and https, as well as a **custom-scheme** for our native Android and IOS Setup.

For our local web development, add a redirectURI for `http://localhost:4444/auth.html` with your custom port. For Android and IOS, add your **custom scheme**. In our case it is `com.example.zitadelflutter`.

:::caution Use Custom Redirect URI!

Your custom scheme has to be compliant with the OAuth 2.0
authentication for mobile devices ([RFC 8252 specification](https://tools.ietf.org/html/rfc8252)).
Otherwise your app might get rejected.

:::

For development, you need to set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.

If you want to redirect the users back to a route on your application after they have logged out, add an optional redirect in the Post Logout URIs field.

Continue and create the application.

After creation, go to **token settings** and check the refresh token checkbox. This allows us to request a refresh_token via `offline_access` scope.
Make sure to save the application.

### Client ID

After successful app creation, a pop-up will appear, showing the app's client ID. Copy the client ID, as you will need it to configure your Flutter application.

## Flutter Prerequisites

To move further in this quickstart, you'll need the following things prepared:

- Have Flutter (and Dart) installed ([how-to](https://flutter.dev/docs/get-started/install))
- Have an IDE set up for developing Flutter ([how-to](https://flutter.dev/docs/get-started/editor))
- Create a basic Flutter app ([how-to](https://flutter.dev/docs/get-started/codelab))
- Create a "Native" application in ZITADEL

After you created the starter Flutter app, the app will show a simple, templated Flutter app.

### Install Dependencies

To authenticate users with ZITADEL in a mobile application, some specific packages are needed.
The [RFC 8252 specification](https://tools.ietf.org/html/rfc8252) defines how
[OAUTH2.0 for mobile and native apps](https://oauth.net/2/native-apps/) works.
Basically, there are two major points in this specification:

1. It recommends to use [PKCE](https://oauth.net/2/pkce/)
2. It does not allow third party apps to use an embedded web view for the login process,
   the app must open the login page within the default browser

First install [http](https://pub.dev/packages/http) a library for making HTTP calls,
then [`flutter_web_auth_2`](https://pub.dev/packages/flutter_web_auth_2) and a secure storage to store the auth / refresh tokens [flutter_secure_storage](https://pub.dev/packages/flutter_secure_storage).

To install run:

```bash
flutter pub add http
flutter pub add flutter_web_auth_2
flutter pub add flutter_secure_storage
```

#### Setup for Android

Navigate to your `AndroidManifest.xml` at `<projectRoot>/android/app/src/main/AndroidManifest.xml` and add the following activity with your custom scheme.

```xml reference
https://github.com/zitadel/zitadel_flutter/blob/main/android/app/src/main/AndroidManifest.xml#L29-L38
```

Furthermore, for `secure_storage`, you need to set the minimum SDK version to 18
in `<projectRoot>/android/app/src/build.gradle`.

### Add Authentication

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

Note that we have to use our http redirect URL for web applications or otherwise use our custom scheme for Android and iOS devices.
To setup other platforms, read the documentation of the [Flutter Web Auth](https://pub.dev/packages/flutter_web_auth_2).

To ensure our application catches the callback URL, you have to create a `auth.html` file in the `/web`
folder with the following content:

```html reference
https://github.com/zitadel/zitadel_flutter/blob/main/web/auth.html
```

Now, you can run your application for iOS and Android devices with

```bash
flutter run
```

or by directly selecting your device

```bash
flutter run -d iphone
```

For Web make sure you run the application on your fixed port such that it matches your redirect URI in your ZITADEL application. We used 4444 as port before so the command would look like this:

```bash
flutter run -d chrome --web-port=4444
```

Our Android and iOS Application opens ZITADEL's login within a custom tab, on Web a new tab is opened.

### Result

If everything works out correctly, your applications should look like this:

<div style={{display: 'grid', 'gridColumnGap': '1rem', 'gridTemplateColumns': '1fr 1fr', 'maxWidth': '500px', 'margin': '0 auto'}}>
    <img src="/docs/img/flutter/not-authed.png" alt="Unauthenticated" height="500px" />
    <img src="/docs/img/flutter/authed.png" alt="Flutter Authenticated" height="500px" />
</div>

<div style={{display: 'grid', 'gridColumnGap': '1rem', 'gridTemplateColumns': '1fr 1fr', 'maxWidth': '800px', 'margin': '0 auto'}}>
    <img src="/docs/img/flutter/web-not-authed.png" alt="Unauthenticated" height="500px" />
    <img src="/docs/img/flutter/web-authed.png" alt="Flutter Authenticated" height="500px" />
</div>
