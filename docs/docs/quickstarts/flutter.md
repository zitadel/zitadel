---
title: Flutter
---

This guide demonstrates, how you integrate **ZITADEL** as an idendity provider to your Flutter app.

At the end of the guide you have a mobile application on Android and iOS with the ability
to authenticate users via ZITADEL.

If you need any other information about Flutter, head over to the [documentation](https://flutter.dev/).

## Prerequisites

To move further in this quickstart, you'll need the following things prepared:

- Have Flutter (and Dart) installed ([how-to](https://flutter.dev/docs/get-started/install))
- Have an IDE set up for developing Flutter ([how-to](https://flutter.dev/docs/get-started/editor))
- Create a basic Flutter app ([how-to](https://flutter.dev/docs/get-started/codelab))
- Create a "Native" application in a ZITADEL project

## Flutter with ZITADEL

### Hello World

After you created the basic Flutter app, the app will show the following screen:

![Flutter Hello World](/img/flutter/hello-world.png)

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

With this, you'll enable "nullable by default" mode in Flutter, as well as new language features.
For this quickstart, the minimal Flutter SDK version is set to the default: `sdk: ">=2.7.0 <3.0.0"`.

### Install Dependencies

To authenticate users with ZITADEL in a mobile application, some specific packages are needed.
The [RFC 8252 specification](https://tools.ietf.org/html/rfc8252) defines how
[OAUTH2.0 for mobile and native apps](https://oauth.net/2/native-apps/) works.
Basically, there are two major points in this specification:

1. It recommends to use [PKCE](https://oauth.net/2/pkce/)
2. It does not allow third party apps to open the browser for the login process,
   the app must open the login page within the embedded browser view

Install the [`appauth`](https://appauth.io/) package and a secure storage (to store the auth / refresh tokens):

```bash
flutter pub add http
flutter pub add flutter_appauth
flutter pub add flutter_secure_storage
```

#### Important on Android

To use this app auth method on Android 11, you'll need to add a `query` to the `AndroidManifest.xml`.
Go to `<projectRoot>/android/app/src/main/AndroidManifest.xml` and add to the `<manifest>` root:

```xml
<queries>
    <intent>
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.BROWSABLE" />
        <data android:scheme="https" />
    </intent>
</queries>
```

This allows the app to query for internal browser activities.

Furthermore, for `secure_storage`, you need to set the minimum SDK version to 18
in `<projectRoot>/android/app/src/build.gradle` and add the manifest placeholder
for your redirect url (the custom url scheme). In the end, the `defaultConfig`
section of the `build.gradle` file should look like this:

```gradle
defaultConfig {
    applicationId "ch.caos.zitadel_quickstart"
    minSdkVersion 18
    targetSdkVersion 30
    versionCode flutterVersionCode.toInteger()
    versionName flutterVersionName
    manifestPlaceholders = [
            'appAuthRedirectScheme': 'ch.caos.zitadel'
    ]
}
```

### Add Authentication

To reduce the commented default code, we will modify the `main.dart` file.

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

What we'll change now, is the `_MyHomePageState` class to enable
authentication via ZITADEL and remove the counter button of the hello
world application. We'll show the username of the authenticated user.

We define the needed elements for our state:

```dart
final _appAuth = FlutterAppAuth();
final _secureStorage = const FlutterSecureStorage();

var _busy = false;
var _authenticated = false;
var _username = '';
```

Then the builder method, which does show the login button if you're not
authenticated, a loading bar if the login process is going on and
your name if you are authenticated:

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

And finally the `_authenticate` method which calls the authorization endpoint,
then fetches the user info and stores the tokens into the secure storage.

```dart
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
```

Now, you should be able to login with a valid ZITADEL user.

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

![Flutter Unauthenticated](/img/flutter/not-authed.png)

![Flutter Authenticated](/img/flutter/authed.png)
