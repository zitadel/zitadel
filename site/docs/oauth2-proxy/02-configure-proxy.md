---
title: OAuth 2.0 Proxy Setup
---

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