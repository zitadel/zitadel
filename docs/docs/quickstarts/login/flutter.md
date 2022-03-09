---
title: Flutter
---

This guide shows you how to integrate **ZITADEL** as an identity provider for a Flutter app.

At the end of the guide, you will have a mobile application on Android and iOS that can authenticate users via ZITADEL.

If you need any other information about Flutter, head to the [Flutter documentation](https://flutter.dev/).

## Prerequisites

Before starting, there are a few things you'll need to do:

- Have Flutter (and Dart) installed ([how-to](https://flutter.dev/docs/get-started/install))
- Have an IDE set up for developing Flutter ([how-to](https://flutter.dev/docs/get-started/editor))
- Create a basic Flutter app ([how-to](https://flutter.dev/docs/get-started/codelab))
- Create a "Native" application in a ZITADEL project

## Flutter with ZITADEL

In your native application on ZITADEL, add a callback (redirect) URI
that matches the selected custom URL scheme.
As an example, if you want to use `ch.myexample.app://sign-me-in` as a redirect URI on ZITADEL and in your app,
register the `ch.myexample.app://` custom URL within Android and iOS.

:::caution Use a Custom Redirect URI!

You'll need the custom redirect URI to be compliant with the OAuth 2.0
authentication for mobile devices ([RFC 8252 specification](https://tools.ietf.org/html/rfc8252)).
Otherwise your app might get rejected.

:::

### Hello World

After you create the basic Flutter app, the app shows the following screen:

<div style={{'text-align': 'center', 'margin-bottom': '1rem'}}>
  <img src="/img/flutter/hello-world.png" alt="Flutter Hello World" height="500px" />
</div>

You may want to change the Flutter SDK version in `pubspec.yaml` from

```yaml
environment:
  sdk: '>=2.7.0 <3.0.0'
```

to

```yaml
environment:
  sdk: '>=2.12.0 <3.0.0'
```

With this, you'll enable "nullable by default" mode in Flutter.
You'll also enable new language features.

For this quickstart, the minimal Flutter SDK version is set to the default: `sdk: ">=2.7.0 <3.0.0"`.

### Install Dependencies

To authenticate users with ZITADEL in a mobile application, you'll need some specific packages.

Install the [`appauth`](https://appauth.io/) package and a secure storage (to store the auth / refresh tokens):

```bash
flutter pub add http
flutter pub add flutter_appauth
flutter pub add flutter_secure_storage
```

#### Important on Android

To use this app auth method on Android 11, add a `query` to the `AndroidManifest.xml`.

Go to `<projectRoot>/android/app/src/main/AndroidManifest.xml` and add to the `<manifest>` root:

```xml title="<projectRoot>/android/app/src/main/AndroidManifest.xml"
<queries>
    <intent>
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.BROWSABLE" />
        <data android:scheme="https" />
    </intent>
    <intent>
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.APP_BROWSER" />
        <data android:scheme="https" />
    </intent>
</queries>
```

This allows the app to query for internal browser activities.

Furthermore, for `secure_storage`, you need to set the minimum SDK version to 18
in `<projectRoot>/android/app/src/build.gradle`.
Then, add the manifest placeholder for your redirect URL (the custom URL scheme).
In the end, the `defaultConfig` section of the `build.gradle` file should look like this:

```groovy title="<projectRoot>/android/app/src/build.gradle"
defaultConfig {
    applicationId "<<YOUR APPLICATION ID, for example ch.myexample.my_fancy_app>>"
    minSdkVersion 18
    targetSdkVersion 30
    versionCode flutterVersionCode.toInteger()
    versionName flutterVersionName
    manifestPlaceholders = [
            'appAuthRedirectScheme': '<<YOUR CUSTOM URL SCHEME, for example ch.myexample.app>>'
    ]
}
```

#### Important on iOS

Similar to Android,
to use custom redirect schemes in iOS, you need to register the custom URL.

In the `Info.plist` file of the Runner
project, you can add the `CFBundleTypeRole` and the `CFBundleUrlSchemes`.

```xml title="<projectRoot>/ios/Runner/Info.plist"
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleTypeRole</key>
        <string>Editor</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>YOUR CUSTOM URL SCHEME, for example ch.myexample.app</string>
        </array>
    </dict>
</array>
```

### Add Authentication

:::note

The auth redirect scheme "`ch.myexample.app`" registers all auth URLs with the given
scheme for the app.
So, a URL that points to `ch.myexample.app://signin` works with the same registration.

:::

To reduce the commented default code, modify the `main.dart` file.

First, the `MyApp` class: it remains a stateless widget:

```dart
class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter ZITADEL',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: MyHomePage(title: 'Flutter ZITADEL Quickstart'),
    );
  }
}
```

Second, the `MyHomePage` class will remain a stateful widget with
its title, we don't change any code here.

```dart
class MyHomePage extends StatefulWidget {
  MyHomePage({Key key, this.title}) : super(key: key);

  final String title;

  @override
  _MyHomePageState createState() => _MyHomePageState();
}
```

Now we will change the `_MyHomePageState` class.
This enables authentication via ZITADEL and removes the counter button of the hello
world application.
We'll show the username of the authenticated user.

1. Define the needed elements for our state:

```dart
final _appAuth = FlutterAppAuth();
final _secureStorage = const FlutterSecureStorage();

var _busy = false;
var _authenticated = false;
var _username = '';
```

2. Then add the builder method.
If you're not authenticated, this shows the login button.
If the login process is going on, it shows a loading bar.
If you are authenticated, it shows your name.

```dart
@override
Widget build(BuildContext context) {
  return Scaffold(
    appBar: AppBar(
      title: Text(widget.title),
    ),
    body: Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          if (!_authenticated && !_busy)
            Text(
              'You are not authenticated.',
            ),
          if (!_authenticated && !_busy)
            ElevatedButton.icon(
                icon: Icon(Icons.fingerprint),
                label: Text('login'),
                onPressed: _authenticate),
          if (_busy)
            Stack(
              children: [
                Center(child: Text("Busy, logging in.")),
                Opacity(
                  opacity: 0.5,
                  child: Container(
                    color: Colors.black,
                  ),
                ),
                Center(child: CircularProgressIndicator()),
              ],
            ),
          if (_authenticated && !_busy)
            Text(
              'Hello $_username!',
            ),
        ],
      ),
    ),
  );
}
```

3. Finally add the `_authenticate` method.
This calls the authorization endpoint,
then fetches the user info and stores the tokens in secure storage.

```dart
Future<void> _authenticate() async {
  setState(() {
    _busy = true;
  });

  try {
    final result = await _appAuth.authorizeAndExchangeCode(
      AuthorizationTokenRequest(
        '<<CLIENT_ID>>', // Client ID of the native application
        '<<CALLBACK_URL>>', // The registered url from zitadel (e.g. ch.myexample.app://signin)
        issuer: '<<ISSUER>>', // most of the cases: https://issuer.zitadel.ch
        scopes: [
          'openid',
          'profile',
          'email',
          'offline_access',
        ],
      ),
    );

    final userInfoResponse = await get(
      Uri.parse('https://api.zitadel.ch/oauth/v2/userinfo'),
      headers: {
        HttpHeaders.authorizationHeader: 'Bearer ${result.accessToken}',
        HttpHeaders.acceptHeader: 'application/json; charset=UTF-8'
      },
    );
    final userJson = jsonDecode(utf8.decode(userInfoResponse.bodyBytes));

    await _secureStorage.write(
        key: 'auth_access_token', value: result?.accessToken);
    await _secureStorage.write(
        key: 'refresh_token', value: result?.refreshToken);
    await _secureStorage.write(key: 'id_token', value: result?.idToken);

    setState(() {
      _busy = false;
      _authenticated = true;
      _username = '${userJson['given_name']} ${userJson['family_name']}';
    });
  } catch (e, s) {
    print('error on authorizeAndExchangeCode token: $e - stack: $s');
    setState(() {
      _busy = false;
      _authenticated = false;
    });
  }
}
```

Now, you can log in as a valid ZITADEL user.

#### Result

In the end, our state class looks like:

```dart
class _MyHomePageState extends State<MyHomePage> {
  final _appAuth = FlutterAppAuth();
  final _secureStorage = const FlutterSecureStorage();

  var _busy = false;
  var _authenticated = false;
  var _username = '';

  Future<void> _authenticate() async {
    setState(() {
      _busy = true;
    });

    try {
      final result = await _appAuth.authorizeAndExchangeCode(
        AuthorizationTokenRequest(
          'CLIENT_ID',
          'CALLBACK_URL',
          issuer: 'ISSUER',
          scopes: [
            'openid',
            'profile',
            'email',
            'offline_access',
          ],
        ),
      );

      final userInfoResponse = await get(
        Uri.parse('https://api.zitadel.ch/oauth/v2/userinfo'),
        headers: {
          HttpHeaders.authorizationHeader: 'Bearer ${result.accessToken}',
          HttpHeaders.acceptHeader: 'application/json; charset=UTF-8'
        },
      );
      final userJson = jsonDecode(utf8.decode(userInfoResponse.bodyBytes));

      await _secureStorage.write(
          key: 'auth_access_token', value: result?.accessToken);
      await _secureStorage.write(
          key: 'refresh_token', value: result?.refreshToken);
      await _secureStorage.write(key: 'id_token', value: result?.idToken);

      setState(() {
        _busy = false;
        _authenticated = true;
        _username = '${userJson['given_name']} ${userJson['family_name']}';
      });
    } catch (e, s) {
      print('error on authorizeAndExchangeCode token: $e - stack: $s');
      setState(() {
        _busy = false;
        _authenticated = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            if (!_authenticated && !_busy)
              Text(
                'You are not authenticated.',
              ),
            if (!_authenticated && !_busy)
              ElevatedButton.icon(
                  icon: Icon(Icons.fingerprint),
                  label: Text('login'),
                  onPressed: _authenticate),
            if (_busy)
              Stack(
                children: [
                  Center(child: Text("Busy, logging in.")),
                  Opacity(
                    opacity: 0.5,
                    child: Container(
                      color: Colors.black,
                    ),
                  ),
                  Center(child: CircularProgressIndicator()),
                ],
              ),
            if (_authenticated && !_busy)
              Text(
                'Hello $_username!',
              ),
          ],
        ),
      ),
    );
  }
}
```

If you run this application, you can authenticate with a valid ZITADEL user.

<div style={{display:'flex', 'justify-content': 'center'}}>
  <div style={{display:'flex', 'align-items': 'center'}}>
    <img src="/img/flutter/not-authed.png" alt="Unauthenticated" height="500px" />
    <span style={{padding:'1rem'}}>becomes</span>
    <img src="/img/flutter/authed.png" alt="Flutter Authenticated" height="500px" />
  </div>
</div>
