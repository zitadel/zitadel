---
title: ZITADEL with Symfony PHP
sidebar_label: Symfony
---

This integration guide demonstrates the recommended way to incorporate ZITADEL into your Symfony PHP application. 
It explains how to enable user login in your application and how to fetch data from the user info endpoint.

By the end of this guide, your application will have login functionality with basic role mapping and will be able to access the current user's profile.

> This documentation references our [example](https://github.com/zitadel/example-symfony-oidc) on GitHub.

## Set up application and obtain keys

![Create app in console](/img/symfony/app-create.png)
![Configure app authentication method in console](/img/symfony/app-auth-method.png)
![Configure app redirects console](/img/symfony/app-redirects.png)


## Symfony setup

Now that you have configured your web application on the ZITADEL side, you can proceed with the integration of your Symfony client.
The example is build on a [generated Symfony web app](https://symfony.com/doc/current/setup.html#creating-symfony-applications), using the following command:

> Skip this step if you are connecting ZITADEL to an existing application.

```bash
symfony new my_project_directory --version="7.0.*" --webapp
cd my_project_directory
```

> The remainder of this guide assumes a Symfony project which already includes all web app bundles, such as security, routing and ORM.
> If you are using this guide against an existing project you must make sure the required bundles are installed using the `composer require` command.

### Install Symfony dependencies

To connect with ZITADEL through OpenID connect, you need to install the [Symfony OIDC bundle](https://github.com/Drenso/symfony-oidc). Run the following command:

```bash
composer require drenso/symfony-oidc-bundle
```

## Define the Symfony app

### Create a User class

First we need to create a User class for the database, so we can persist User info between requests. In this case you don't need password authentication.
Email addresses are not unique for ZITADEL users. There can be multiple user accounts with the same email address.
See [User Constraints](https://zitadel.com/docs/concepts/structure/users#constraints) for more details.
We will use the User Info `sub` claim as unique "display" name for the user. `sub` equals the unique User ID from ZITADEL.
This creates a User Repository and Entity that implements the `UserInterface`:

> You can skip this step if you already have an existing User object in your project.

```bash
php bin/console make:user

 The name of the security user class (e.g. User) [User]:
 > User

 Do you want to store user data in the database (via Doctrine)? (yes/no) [yes]:
 > yes

 Enter a property name that will be the unique "display" name for the user (e.g. email, username, uuid) [email]:
 > sub

 Will this app need to hash/check user passwords? Choose No if passwords are not needed or will be checked/hashed by some other system (e.g. a single sign-on server).

 Does this app need to hash/check user passwords? (yes/no) [yes]:
 > no
```

Next, extend the User Entity with properties that we will obtain from ZITADEL and use in the application later.

> None of the following properties are required for authentication, but show how we can map User Info to a Symfony User Entity later. You can adjust the properties how you wish for your application.

```bash
php bin/console make:entity

 Class name of the entity to create or update (e.g. GrumpyElephant):
 > User

 Your entity already exists! So let's add some new fields!

 New property name (press <return> to stop adding fields):
 > display_name

 Field type (enter ? to see all types) [string]:
 > string

 Field length [255]:
 > 255

 Can this field be null in the database (nullable) (yes/no) [no]:
 > yes

 updated: src/Entity/User.php

 Add another property? Enter the property name (or press <return> to stop adding fields):
 > email

 Field type (enter ? to see all types) [string]:
 > string

 Field length [255]:
 > 255

 Can this field be null in the database (nullable) (yes/no) [no]:
 > yes

 updated: src/Entity/User.php

 Add another property? Enter the property name (or press <return> to stop adding fields):
 > email_verified

 Field type (enter ? to see all types) [string]:
 > boolean

 Can this field be null in the database (nullable) (yes/no) [no]:
 > yes

 updated: src/Entity/User.php

 Add another property? Enter the property name (or press <return> to stop adding fields):
 > created_at

 Field type (enter ? to see all types) [datetime_immutable]:
 > datetime_immutable

 Can this field be null in the database (nullable) (yes/no) [no]:
 > no

 updated: src/Entity/User.php

 Add another property? Enter the property name (or press <return> to stop adding fields):
 > updated_at

 Field type (enter ? to see all types) [datetime_immutable]:
 > datetime_immutable

 Can this field be null in the database (nullable) (yes/no) [no]:
 > no

 updated: src/Entity/User.php

 Add another property? Enter the property name (or press <return> to stop adding fields):
 > 
           
  Success!
```

Now edit `src/Entity/User.php` to add some methods to pretty-print user data later in this example. Add import near the top of the file:

```php
use DateTimeInterface;
```

And extend the User class with this methods:

```php
class User implements UserInterface
{
    ...

    public function implodeRoles(): string
    {
        return implode(', ', $this->getRoles());
    }

    public function formatCreatedAt(): string
    {
        return $this->created_at->format(DateTimeInterface::W3C);
    }

    public function formatUpdatedAt(): string
    {
        return $this->updated_at->format(DateTimeInterface::W3C);
    }
}
```

When you are done, the User Entity should look something like:

```php reference
https://github.com/zitadel/example-symfony-oidc/blob/main/src/Entity/User.php
```

Edit the User Repository to have a `findOneBySub` method, used later for OIDC User Info updates.

```php reference
https://github.com/zitadel/example-symfony-oidc/blob/main/src/Repository/UserRepository.php
```

### Create a Security Provider

Next you will need to create a Security Provider that integrates the OIDC flow between Symfony and ZITADEL.
Create a `ZitadelUserProvider` which implements `UserProviderInterface`, `OidcUserProviderInterface` and `LoggerAwareInterface`.
`LoggerAwareInterface` is optional if you want debug logging.

> We called this a `ZitadelUserProvider` because it carries a custom scope and claim mapping from ZITADEL roles to the Symfony role system.

```php reference
https://github.com/zitadel/example-symfony-oidc/blob/main/src/Security/ZitadelUserProvider.php
```

You can customize the User Info that is obtained and stored by adjusting the `SCOPES` constant, the `updateUserEntity` method and the User Entity.

### Controllers and templates

We need to create couple of Controllers and templates to define the app.

#### Index

The index controller serves a public page on the `/` route and provides some basic links to the authenticated sections of the app.

```php reference
https://github.com/zitadel/example-symfony-oidc/blob/main/src/Controller/IndexController.php
```

The index template:

```twig reference
https://github.com/zitadel/example-symfony-oidc/blob/main/templates/index.html.twig
```

#### Login

The login controller initiates the OIDC login flow by creating a Auth request and redirecting the user to ZITADEL.

```php reference
https://github.com/zitadel/example-symfony-oidc/blob/main/src/Controller/LoginController.php
```

#### Profile

The profile controller displays User Info of the currently authenticated user.
Any authenticated user will have access to this page.

```php reference
https://github.com/zitadel/example-symfony-oidc/blob/main/src/Controller/ProfileController.php
```

The profile template maps the User Entity to a HTML page.

```twig reference
https://github.com/zitadel/example-symfony-oidc/blob/main/templates/profile.html.twig
```

#### User list

The user list controller displays all users from the database that were created during OIDC login.
Only users with an admin role will have access to this page.

```php reference
https://github.com/zitadel/example-symfony-oidc/blob/main/src/Controller/UserListController.php
```

```twig reference
https://github.com/zitadel/example-symfony-oidc/blob/main/templates/user_list.html.twig
```

## Configure and run the application

> Never store and commit secrets in a `.env` file. Use a `env.local` file instead and make sure the file is in `.gitignore`.

### Database

Make sure you have a database configured in `.env` or `.env.local`. This example uses a local sqlite file to simplify setup:

```sh
DATABASE_URL="sqlite:///%kernel.project_dir%/var/data.db"
```

Create and run migrations:

```bash
php bin/console make:migration
php bin/console doctrine:migrations:migrate
```

### Security

A firewall needs to be defined along with roles based access control rules.
In the following example we define the `zitadel_user_provider` as the security class we wrote earlier. We configure the main firewall to use the `zitadel_user_provider` and listen for logout requests on the `/logout` path. We tell the oidc module to enable End Session support.

In the `access_control` section we protect the `/users` and `/profile` routes based on roles. Roles are mapped from ZITADEL to Symfony in the `ZitadelUserProvider` we wrote earlier.

```yaml reference
https://github.com/zitadel/example-symfony-oidc/blob/main/config/packages/security.yaml
```

### OIDC

The generated [`dresno_oidc.yaml`](https://github.com/zitadel/example-symfony-oidc/blob/main/config/packages/drenso_oidc.yaml) file can be edited to customize behavior of the OIDC bundle. For this example we stick with the default and use environment variables to connect to ZITADEL.

Edit `.env.local` to contain the details from the [Application setup section](#set-up-application-and-obtain-keys).

```sh
OIDC_WELL_KNOWN_URL="https://tims-zitadel-instance-oj7iry.zitadel.cloud/.well-known/openid-configuration"
OIDC_CLIENT_ID="248680248240075805@dev"
OIDC_CLIENT_SECRET="BJPhEJULSUXseC4geqg5Yg4wWMoy7RgZKar86mbIpt8ZekC5kixMzYGcXLDeeJv7"
```

> The well-known URL needs to be adjusted to your own instance domain.

Activate the route that is used as callback by the OIDC bundle:

```yaml reference
https://github.com/zitadel/example-symfony-oidc/blob/main/config/routes.yaml#L6-L7
```

### Run

Now we can use a local Symfony server to test the application.

```bash
symfony server:start --no-tls
```
