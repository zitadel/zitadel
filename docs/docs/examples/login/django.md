---
title: ZITADEL with Django Python
sidebar_label: Django
---

This integration guide demonstrates the recommended way to incorporate ZITADEL into your Django Python application.
It explains how to enable user login in your application and how to incorporate the ZITADEL users into the existing AuthenticationBackend.

By the end of this guide, your application will have login functionality with basic role mapping, admin console and polls as described in the Django guide.

:::info
This documentation references our [example](https://github.com/zitadel/example-django-python-oidc) on GitHub.
:::

## ZITADEL setup

Before we can start building our application, we have to do a few configuration steps in ZITADEL Console.

### Project roles

The Example expects [user roles](guides/integrate/retrieve-user-roles) to be returned after login.
This example expects 3 different roles:
- `admin`: superuser with permissions to use the admin console
- `staff`: user with permissions to see results of the polls
- `user`: normal user with permission to vote on the existing polls

In your project settings make sure the "Assert Roles On Authentication" is enabled.

![Project settings in console](/img/django/project-settings.png)

In the project Role tab, add 2 special roles:

- `admin`
- `staff`
- `user`

If none of the roles is provided as a user, the user in Django will not be created.

![Project roles in console](/img/django/project-roles.png)

Finally, we can assign the roles to users in the project's authorizations tab.

![Project authorizations in console](/img/django/project-authorizations.png)

### Set up application and obtain secrets

Next you will need to provide some information about your app.

In your Project, add a new application at the top of the page.
Select Web application type and continue.
We use [Authorization Code](/apis/openidoauth/grant-types#authorization-code)for our Symfony application.

![Create app in console](/img/django/app-create.png)

Select `CODE` in the next step. This makes sure you still get a secret. Note that the secret never gets exposed on the browser and is therefore kept in a confidential environment. Safe the generated

![Configure app authentication method in console](/img/django/app-auth-method.png)

With the Redirect URIs field, you tell ZITADEL where it is allowed to redirect users to after authentication. For development, you can set dev mode to `true` to enable insecure HTTP and redirect to a `localhost` URI.

For the example application we are writing use:

- `http://localhost:8000/oidc/callback/` as Redirect URI
- `http://localhost:8000/oidc/logout/` as post-logout URI.

![Configure app redirects console](/img/django/app-redirects.png)

After the final step you are presented with a client ID and secret.
Copy and paste them to a safe location for later use by the application.
The secret will not be displayed again, but you can regenerate one if you loose it.

## Setup new Django application

Now that you have configured your web application on the ZITADEL side, you can proceed with the integration of your Django client.
The used base is the "Writing your first Django app"-app from the Django documentation under [https://docs.djangoproject.com/en/5.0/intro/](https://docs.djangoproject.com/en/5.0/intro/), which has documented additional parts in to use [mozilla-django-oidc](https://github.com/mozilla/mozilla-django-oidc) to integrate ZITADEL as AuthenticationBackend.

:::info
Skip this step if you are connecting ZITADEL to an existing application.
:::

### Install Django depencies

To connect with ZITADEL through OpenID connect, you need to install a OIDC depencency, we will use the [mozilla-django-oidc](https://github.com/mozilla/mozilla-django-oidc). Run the following command:

```bash
python -m pip install mozilla-django-oidc
```

## Define the Django app

### Create the settings.py to include mozilla-django-oidc

To use the mozilla-django-oidc as AuthenticationBackend, there are several things to add to the settings.py, as described in the [documentation "Add settings to settings.py"](https://mozilla-django-oidc.readthedocs.io/en/stable/installation.html#add-settings-to-settings-py):

```python
INSTALLED_APPS = [
    ...
    "mozilla_django_oidc",  # Load after auth
    ...
]

MIDDLEWARE = [
    #...
    "mozilla_django_oidc.middleware.SessionRefresh",
]

AUTHENTICATION_BACKENDS = (
    "mysite.backend.PermissionBackend",
)

ZITADEL_PROJECT = "ID of the project you created the application in ZITADEL" 
OIDC_RP_CLIENT_ID = "ClientID provided by the created application in ZITADEL"
OIDC_RP_CLIENT_SECRET = "ClientSecret provided by the created application in ZITADEL"
OIDC_RP_SIGN_ALGO = "RS256"

OIDC_OP_JWKS_ENDPOINT = "https://example.zitadel.cloud/oauth/v2/keys"
OIDC_OP_AUTHORIZATION_ENDPOINT = "https://example.zitadel.cloud/oauth/v2/authorize"
OIDC_OP_TOKEN_ENDPOINT = "https://example.zitadel.cloud/oauth/v2/token"
OIDC_OP_USER_ENDPOINT = "https://example.zitadel.cloud/oidc/v1/userinfo"

LOGIN_REDIRECT_URL = "http://localhost:8000"
LOGOUT_REDIRECT_URL = "http://localhost:8000"
LOGIN_URL = "http://localhost:8000/oidc/authenticate/"
```

### AuthenticationBackend definition

To create and update the users regarding the roles given in the authentications in ZITADEL a Subclass of OIDCAuthenticationBackend has to be created:

backend.py
```
from mozilla_django_oidc.auth import OIDCAuthenticationBackend
from django.contrib.auth.models import Permission, User
from django.contrib import admin


class PermissionBackend(OIDCAuthenticationBackend):
    def create_user(self, claims):
        email = claims.get("email")
        username = self.get_username(claims)
        permClaim = (
            "urn:zitadel:iam:org:project:"
            + self.get_settings("ZITADEL_PROJECT")
            + ":roles"
        )

        if "admin" in claims[permClaim].keys():
            return self.UserModel.objects.create_user(
                username, email=email, is_superuser=True, is_staff=True
            )
        elif "staff" in claims[permClaim].keys():
            return self.UserModel.objects.create_user(
                username, email=email, is_staff=True
            )
        elif "user" in claims[permClaim].keys():
            return self.UserModel.objects.create_user(username, email=email)
        else:
            return self.UserModel.objects.none()

    def update_user(self, user, claims):
        permClaim = (
            "urn:zitadel:iam:org:project:"
            + self.get_settings("ZITADEL_PROJECT")
            + ":roles"
        )
        
        if "admin" in claims[permClaim].keys():
            user.is_superuser = True
            user.is_staff = True
        elif "staff" in claims[permClaim].keys():
            user.is_superuser = False
            user.is_staff = True
        elif "user" in claims[permClaim].keys():
            user.is_superuser = False
            user.is_staff = False
        return user

```

Which handles the users differently depending on if there are roles associated to:
- `admin` -> superuser
- `staff` -> staff
- `user` -> user
- `no role` -> no user gets created

### URLs 

To handle the callback and logout the urls have to be added to the urls.py:
```python
urlpatterns = [
    #...
    path("oidc/", include("mozilla_django_oidc.urls")),
]
```

## Configure and run the application

:::warning
Never store and commit secrets in the settings.py file
:::

### DB

Create and run migrations:

```bash
python manage.py migrate
```

### Run

You can use a local Django server to test the application.

```bash
python manage.py runserver
```

Visit http://localhost:8000/polls or http://localhost:8000/admin and click around.

## Completion

Congratulations! You have successfully integrated your Symfony application with ZITADEL!

If you get stuck, consider checking out our [example](https://github.com/zitadel/example-python-django-oidc) application. This application includes all the functionalities mentioned in this quick-start. You can start by cloning the repository and defining the settings in the settings.py. If you face issues, contact us or raise an issue on [GitHub](https://github.com/zitadel/example-python-django-oidc/issues).

### What's next?

Now that you have enabled authentication, it's time for you to add more authorizations to your application using ZITADEL APIs. To do this, you can refer to the [docs](/apis/introduction) or check out the ZITADEL Console code on [GitHub](https://github.com/zitadel/zitadel) which uses gRPC and OpenAPI to access data.
