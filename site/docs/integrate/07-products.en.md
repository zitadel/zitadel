---
title:  Products
description: ...
---

### Grafana Example

#### Authenticate Grafana with ZITADEL

To authenticate **Grafana** with ZITADEL you can use the provided **Generic OAuth** plugin.

> We do not recommend that you rely on `allowed_domain` as means of authorizing subjects, but instead use **ZITADEL's** RBAC Assertion

1. Create Project with OpenID Connect / OAuth Client
2. Add config to your **Grafana** instance and restart it

```ini
[auth.generic_oauth]
enabled = true
name= ZITADEL
client_id = {ZITADEL_GENERATED_CLIENT_ID}
client_secret = {ZITADEL_GENERATED_CLIENT_SECRET}
scopes = openid profile email
auth_url = https://accounts.zitadel.ch/oauth/v2/authorize
token_url = https://api.zitadel.ch/oauth/v2/token
api_url = https://api.zitadel.ch/oauth/v2/userinfo
allow_sign_up = true
```

> Redirect URI https://<grafana domain>/login/generic_oauth

#### Authorizes Users with Roles in Grafana

**ZITADEL** provides projects with the option to provide Grafana with the users role.

1. Create Roles (Admin, Editor, Viewer) in **ZITADEL's** project which contains **Grafana**
2. Enable "Assert Roles on Authentication" so that the roles are asserted to the userinfo endpoint
3. (Optional) Enable "Check roles on Authentication", this will prevent that someone without any role can login to **Grafana** via **ZITADEL**
4. Append the config below to your **Grafana** instance and reload
5. Authorize the necessary users

```ini
[auth.generic_oauth]
...
role_attribute_path =  keys("urn:zitadel:iam:org:project:roles") | contains(@, 'Admin') && 'Admin' || contains(@, 'Editor') && 'Editor' || 'Viewer'
```

> Grafana can not directly use ZITADEL delegation feature but normal RBAC works fine
> Additional infos can be found in the [Grafana generic OAuth 2.0 documentation](https://grafana.com/docs/grafana/latest/auth/generic-oauth/)

### ArgoCD Example

> TODO

### Kubernetes Example

> TODO
