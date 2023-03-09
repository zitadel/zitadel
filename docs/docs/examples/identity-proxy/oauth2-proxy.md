---
title: OAuth 2.0 Proxy
---

<!-- //TODO Florian update this to zitadel.cloud-->

[OAuth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy) is a project which allows services to delegate the authentication flow to a IDP, for example **ZITADEL**

## Configure ZITADEL

### Setup Application and get Keys

Before we can start building our application we have do do a few configuration steps in ZITADEL Console.
You will need to provide some information about your app. We recommend creating a new app to start from scratch. Navigate to your project and add a new application at the top of the page.
Select Web Application and continue.
We recommend that you use [Authorization Code](/apis/openidoauth/grant-types#authorization-code) for the OAuth 2.0 Proxy.

> Make sure Authentication Method is set to `BASIC` and the Application Type is set to `Web`.

### Redirect URLs

A redirect URL is a URL in your application where ZITADEL redirects the user after they have authenticated. Set your url to the domain the proxy will be deployed to or use the default one `http://127.0.0.1:4180/oauth2/callback`.

> If you are following along with the sample project you downloaded from our templates, you should set the Allowed Callback URL to <http://localhost:4200/auth/callback>. You will also have to set dev mode to `true` as this will enable unsecure http for the moment.

If you want to redirect the users back to a route on your application after they have logged out, add an optional redirect in the post redirectURI field.

Continue and Create the application.

### Client ID and Secret

After successful app creation a popup will appear showing you your clientID as well as a secret.
Copy your client ID and Secrets as it will be needed in the next step.

> Note: You will be able to regenerate the secret at a later time if you loose it.

## OAuth 2.0 Proxy Setup

### Authentication Example

```toml
provider = "oidc"
user_id_claim = "sub" #uses the subject as ID instead of the email
provider_display_name = "ZITADEL"
redirect_url = "http://127.0.0.1:4180/oauth2/callback"
oidc_issuer_url = "https://{your_domain}.zitadel.cloud"
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

> This was tested with version `oauth2-proxy v7.4.0 (built with go1.20.0)`

### Check for groups

If you want oauth2-proxy to check for roles in the tokens you have to add an [action](/docs/apis/actions/introduction) in ZITADEL to [complement the token](/docs/apis/actions/complement-token) according to [this example](https://github.com/zitadel/actions/blob/main/examples/custom_roles.js) and add the following configuration to the config:

```toml
oidc_groups_claim = "{your_actions_group_key}"
allowed_groups = ["list", "of", "allowed", "roles"]
```

## Completion

You have successfully integrated ZITADEL in your proxy!

### What next?
