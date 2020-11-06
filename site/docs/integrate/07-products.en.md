---
title:  Products
description: ...
---

### Grafana Example

**Grafana** defines itself as "The open-source platform for monitoring and observability."

The source code is provided on [Grafana's Github Repository](https://github.com/grafana/grafana)

#### Authenticate Grafana with ZITADEL

To authenticate **Grafana** with ZITADEL you can use the provided **Generic OAuth** plugin.

> We do not recommend that you rely on `allowed_domain` as means of authorizing subjects, but instead use **ZITADEL's** RBAC Assertion

1. Create a new project or use an existing one
2. Add OpenID Connect / OAuth 2.0 client to the project (See screenshot for settings)
3. Add config to your **Grafana** instance and restart it
4. Login to **Grafana**

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

> Grafanas's redirect is URI https://yourdomain.tld/login/generic_oauth

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/grafana_project_settings.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/grafana_project_settings.png" itemprop="thumbnail" alt="Project Settings for Grafana" />
        </a>
        <figcaption itemprop="caption description">Project Settings for Grafana</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/grafana_client_settings.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/grafana_client_settings.png" itemprop="thumbnail" alt="Client Settings for Grafana" />
        </a>
        <figcaption itemprop="caption description">Client Settings for Grafana</figcaption>
    </figure>
</div>

#### Authorizes Users with Roles in Grafana

**ZITADEL** provides projects with the option to provide Grafana with the users role.

1. Create Roles (Admin, Editor, Viewer) in **ZITADEL's** project which contains **Grafana**
2. Enable "Assert Roles on Authentication" so that the roles are asserted to the userinfo endpoint
3. (Optional) Enable "Check roles on Authentication", this will prevent that someone without any role to login **Grafana** via **ZITADEL**
4. Append the config below to your **Grafana** instance and reload
5. Authorize the necessary users

```ini
[auth.generic_oauth]
...
role_attribute_path =  keys("urn:zitadel:iam:org:project:roles") | contains(@, 'Admin') && 'Admin' || contains(@, 'Editor') && 'Editor' || 'Viewer'
...
```

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/grafana_zitadel_authorization.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/grafana_zitadel_authorization.png" itemprop="thumbnail" alt="Authorization for Grafana Role in ZITADEL" />
        </a>
        <figcaption itemprop="caption description">Authorization for Grafana Role in ZITADEL</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/grafana_login_button.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/grafana_login_button.png" itemprop="thumbnail" alt="Grafana Login" />
        </a>
        <figcaption itemprop="caption description">Grafana Login</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/grafana_profile_settings.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/grafana_profile_settings.png" itemprop="thumbnail" alt="Grafana with Editor Role mapped from ZITADEL" />
        </a>
        <figcaption itemprop="caption description">Grafana with Editor Role mapped from ZITADEL</figcaption>
    </figure>
</div>

> Grafana can not directly use ZITADEL delegation feature but normal RBAC works fine
> Additional infos can be found in the [Grafana generic OAuth 2.0 documentation](https://grafana.com/docs/grafana/latest/auth/generic-oauth/)

### ArgoCD Example

> TODO

### Kubernetes Example

> TODO
