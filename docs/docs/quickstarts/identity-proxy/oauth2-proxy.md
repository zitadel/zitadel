---
title: OAuth 2.0 Proxy
---

The [OAuth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy) project lets services delegate the authentication flow to an IDP like **ZITADEL**

## Configure Zitadel

### Setup Application and get Keys

Before building your application, you'll need to head to the ZITADEL console and add some information about your app.
To start, we recommend creating a new app from scratch.
To do so:

1. Navigate to your [Project](https://console.zitadel.ch/projects).
1. Add a new application at the top of the page.
1. Select **Web Application** and continue.

For the OAuth 2.0 Proxy, we recommend using [Authorization Code](../../apis/openidoauth/grant-types#authorization-code).

> Make sure the Authentication Method is set to `BASIC` and the Application Type is set to `Web`.

### Redirect URLs

Set a redirect URL.
After users authenticate, ZITADEL will redirect them to this URL.
Set your URL to the domain the proxy will deploy to.
You can also use the default, `http://127.0.0.1:4180/oauth2/callback`.

> If you are following along with the sample project you downloaded from our templates,
> set the Allowed Callback URL to <http://localhost:4200/auth/callback>.
> You will also have to set dev mode to `true`.
> This enables unsecure http for the moment.

After users log out, you can redirect users back to a route on your application.
To do so, add an optional redirect in the post redirectURI field.

Continue and Create the application.

### Client ID and Secret

After you create your app, a popup will show your clientID and secret.
Copy these&mdash;
you'll use them in the next step.

> Note: If you lose your secret, you can regenerate it later.

## OAuth 2.0 Proxy Setup

### Authentication Example

```toml
provider = "oidc"
user_id_claim = "sub" #uses the subject as ID instead of the email
provider_display_name = "ZITADEL"
redirect_url = "http://127.0.0.1:4180/oauth2/callback"
oidc_issuer_url = "https://issuer.zitadel.ch"
upstreams = [
    "https://example.corp.com"
]
email_domains = [
    "*"
]
client_id = "{ZITADEL_GENERATED_CLIENT_ID}"
client_secret = "{ZITADEL_GENERATED_CLIENT_SECRET}"
pass_access_token = true
cookie_secret = "{SUPPLY_SOME_SECRET_HERE}"
skip_provider_button = true
cookie_secure = false #localdev only false
http_address = "127.0.0.1:4180" #localdev only
```

> This was tested with version `oauth2-proxy v6.1.1 (built with go1.14.2)`

## Completion

You have successfully integrated ZITADEL in your proxy!

### What next?
