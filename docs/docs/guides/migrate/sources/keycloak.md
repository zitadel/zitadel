---
title: Migrate from Keycloak
sidebar_label: From Keycloak
---

## Migrating from Keycloak to ZITADEL

This tutorial will use [Docker installation](https://www.docker.com/) for the prerequisites. However, both Keycloak and ZITADEL offer different installation methods. As a result, this guide won't include any required production tuning or security hardening for either system. However, it's advised you follow [recommended guidelines](https://docs.zitadel.com/docs/guides/manage/self-hosted/production) before putting those systems into production.

## Setting Up Keycloak
### Run Keycloak

To begin setting up Keycloak, you need to refer to the official [Keycloak Docker image](https://www.keycloak.org/getting-started/getting-started-docker). You'll use it to run a development version of the Keycloak server on your local machine:


```
Bash Script

docker run -p 8081:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:22.0.1 start-dev

```

In a few seconds, Keycloak will be available at [http://localhost:8081](http://localhost:8081). Access the **Administration Console** via the username `admin` and password `admin`:

<img src="/docs/img/guides/migrate/keycloak-01.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-02.png" alt="Migrating users from Keycloak to ZITADEL"/>

> Please note that this is a development setup, and it's not secured for production usage.


### Create a Realm in Keycloak 

In order to configure Keycloak as the identity provider for your application, you need to create a new realm. This will allow users and authentication resources to be isolated from any other Keycloak usage. Click on the sidebar drop-down menu and select **Create Realm**. Then input the desired realm name and click **Create**:

<img src="/docs/img/guides/migrate/keycloak-03.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-04.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-05.png" alt="Migrating users from Keycloak to ZITADEL"/>


### Create an OAuth2/OIDC Client in Keycloak

After you create the new realm, you need to set up a new client. A client is a core concept in the [OAuth 2 protocol](https://oauth.net/2/grant-types/) used in all flows. Here, it will be used by your application to connect to Keycloak and complete authentication lookups.

On the menu on the left, select **Clients** and click the button **Create Client**. 

<img src="/docs/img/guides/migrate/keycloak-06.png" alt="Migrating users from Keycloak to ZITADEL"/>

Leave the client type as **OpenID Connect**, fill in the desired **Client ID** and **Name**, then click on **Next**:

<img src="/docs/img/guides/migrate/keycloak-07.png" alt="Migrating users from Keycloak to ZITADEL"/>

You don't need to make any changes to the **Capability config**. Go ahead and click **Next**.

<img src="/docs/img/guides/migrate/keycloak-08.png" alt="Migrating users from Keycloak to ZITADEL"/>

The web application used for this demo will run on `http://localhost:4200/`. In order to allow login and logout from the application, this client needs to be configured to accept your application URL in **Login settings**. Edit **Root URL**, **Valid redirect URI**, and **Valid post logout redirect URIs** to point to your application URLs. Without this configuration, Keycloak will refuse login and logout from your application due to security concerns.
 
Additionally, **Web origins** needs to be configured to support required cross-domain requests; otherwise, the request will be blocked on all browsers due to security concerns. To make this application work, fill in all the fields as shown below. 

<img src="/docs/img/guides/migrate/keycloak-09.png" alt="Migrating users from Keycloak to ZITADEL"/>

Finally, create the client by clicking **Save**.

<img src="/docs/img/guides/migrate/keycloak-10.png" alt="Migrating users from Keycloak to ZITADEL"/>

### Create User in Keycloak
The last thing you need to do in Keycloak is to create at least one new user. This user will be able to log into your application.

On the menu on the left, select **Users**, and click **Add user**. Fill in the username, email, and first and last names; and mark the email as verified. Click on **Create** to create a new user:

<img src="/docs/img/guides/migrate/keycloak-11.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-12.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-13.png" alt="Migrating users from Keycloak to ZITADEL"/>

Now you should attach a password to this user. Select the **Credentials** tab and click **Set password**.

<img src="/docs/img/guides/migrate/keycloak-14.png" alt="Migrating users from Keycloak to ZITADEL"/>

 On the new modal panel, input the desired password and select **Save**.

<img src="/docs/img/guides/migrate/keycloak-15.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-16.png" alt="Migrating users from Keycloak to ZITADEL"/>

### Export Keycloak Users
Keycloak provides an [export](https://www.keycloak.org/server/importExport) functionality that allows user information to be extracted into JSON files. While it's intended to be used in another Keycloak instance, you can manipulate it to export users to a different user management system.

For example, in order to generate the export files with Keycloak, you will need to enter the Docker container, run the export command, and copy it outside the container:

```Bash Script
# Recover the Container ID for Keycloak
docker ps

# Starting a bash command inside the Keycloak container
# use the container ID of Keycloak
docker exec -it <keycloak container ID> bash 

# And run the export command
 /opt/keycloak/bin/kc.sh export --dir /tmp

#exit
exit

# copy generated files from docker container to local machine
docker cp <keycloak container ID>:/tmp/my-realm-users-0.json .
```

## Configure the Web Application

Now that you've fully configured Keycloak, you need to configure an application to use Keycloak as its user management interface.

This stage focuses on ZITADEL's [sample Angular application](https://github.com/zitadel/zitadel-angular) (client-side application) as a sample application that requires a user login. This application uses [OAuth Authorization Code flow with PKCE](https://docs.zitadel.com/docs/guides/integrate/login-users#code-flow) for its authentication, which is supported by both Keycloak and ZITADEL.

In order to set up the ZITADEL sample Angular application, ensure you have [Node.js](https://nodejs.org/en/) installed and clone the example repository:

```Bash Script
git clone https://github.com/zitadel/zitadel-angular.git
npm install -g @angular/cli
npm install
```

In this tutorial, the Keycloak realm is named `my-realm`, and the client ID is `test-client`. Edit the [src/app/app.module.ts file](https://github.com/zitadel/zitadel-angular/blob/main/src/app/app.module.ts) and update the `client ID` and `issuer`:

```JavaScript 
const authConfig: AuthConfig = {
   scope: 'openid profile email',
   responseType: 'code',
   oidc: true,
   clientId: 'test-client',
   issuer: 'http://localhost:8081/realms/my-realm',
   redirectUri: 'http://localhost:4200/auth/callback',
   postLogoutRedirectUri: 'http://localhost:4200/signedout',
   requireHttps: false, // required for running locally
   disableAtHashCheck: true
};
```

The landing page of the application will be as follows: 


Per the sample [documentation](https://github.com/zitadel/zitadel-angular#readme), running the development server will serve the browser-side application at [http://localhost:4200/](http://localhost:4200/):

```Bash Script
npm start
```

<img src="/docs/img/guides/migrate/keycloak-17.png" alt="Migrating users from Keycloak to ZITADEL"/>

After clicking the **Authenticate** button, users will be redirected to the login page hosted in Keycloak:

<img src="/docs/img/guides/migrate/keycloak-18.png" alt="Migrating users from Keycloak to ZITADEL"/>

After successfully logging into Keycloak, users are redirected back to the sample application:


Now, you've successfully created a sample application that uses Keycloak for user and session management.

<img src="/docs/img/guides/migrate/keycloak-19.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-20.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-21.png" alt="Migrating users from Keycloak to ZITADEL"/>

Now, you've successfully created a sample application that uses Keycloak for user and session management.

## Setting Up ZITADEL

After creating a sample application that connects to Keycloak, you need to set up ZITADEL in order to migrate the application and users from Keycloak to ZITADEL. For this, ZITADEL offers a [Docker Compose](https://zitadel.com/docs/self-hosting/deploy/compose) installation guide:

```Bash Script
# Download the docker compose example configuration.

wget https://raw.githubusercontent.com/zitadel/zitadel/main/docs/docs/self-hosting/deploy/docker-compose.yaml


# Run the database and application containers. It will start both the database and the web-server components
docker compose up --detach
```

In a few seconds, the application will be available at [http://localhost:8080/ui/console/](http://localhost:8080/ui/console/).

<img src="/docs/img/guides/migrate/keycloak-22.png" alt="Migrating users from Keycloak to ZITADEL"/>

Now you can access the console with the following default credentials:

* **Username**: `zitadel-admin@zitadel.localhost`
* **Password**: `Password1!`

<img src="/docs/img/guides/migrate/keycloak-23.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-24.png" alt="Migrating users from Keycloak to ZITADEL"/>

Skip the 2-factor set up.

<img src="/docs/img/guides/migrate/keycloak-25.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-26.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-27.png" alt="Migrating users from Keycloak to ZITADEL"/>


Next, you need to create a new project and a new application in ZITADEL. ZITADEL projects are a similar concept to Keycloak realms, and an application is equivalent to Keycloak clients.

To create a new project, select the **Projects** tab and click **Create New Project**. Fill in the desired name of the project and click **Continue**:

<img src="/docs/img/guides/migrate/keycloak-28.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-29.png" alt="Migrating users from Keycloak to ZITADEL"/>

Within a project, click on the **+** button to create a new application. Fill in the desired name, and select **Web** for the type of application and click on **Continue**.

<img src="/docs/img/guides/migrate/keycloak-30.png" alt="Migrating users from Keycloak to ZITADEL"/>

Select **PKCE** and click **Continue**.

<img src="/docs/img/guides/migrate/keycloak-31.png" alt="Migrating users from Keycloak to ZITADEL"/>

Now configure **Redirect URIs** (http://localhost:4200/auth/callback
) and **Post Logout URIs** (http://localhost:4200/signedout
) to the sample application URL, and click **Continue**. 

<img src="/docs/img/guides/migrate/keycloak-32.png" alt="Migrating users from Keycloak to ZITADEL"/>

Click on **Create** to create the new PKCE application in ZITADEL. You will now have access to the randomly generated Client ID, which will be used later on to configure your application:

<img src="/docs/img/guides/migrate/keycloak-33.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-34.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-35.png" alt="Migrating users from Keycloak to ZITADEL"/>

You need to configure additional origins to prevent CORS errors. Inside the application configuration, select **Additional Origins** and add the required URLs:


Now that you've finalized the ZITADEL configuration for the project and application, your last required change is to modify the sample application to now use ZITADEL instead of Keycloak. Since both tools implement the same authentication flows, all you need to do is change the `issuer` URL and the `client ID`:

```JavaScript
const authConfig: AuthConfig = {
   scope: 'openid profile email',
   responseType: 'code',
   oidc: true,
   clientId: '<your_cliend_id>` // e.g. 230518162431475715@testproject
,
   issuer: 'http://localhost:8080', 
   redirectUri: 'http://localhost:4200/auth/callback',
   postLogoutRedirectUri: 'http://localhost:4200/signedout',
   requireHttps: false, // required for running locally
   disableAtHashCheck: true
};
```

After you change the sample configuration, when attempting to authenticate, the login page will be served by ZITADEL:

```Bash Script
npm start
```
Access [http://localhost:4200/](http://localhost:4200/):

<img src="/docs/img/guides/migrate/keycloak-36.png" alt="Migrating users from Keycloak to ZITADEL"/>

Now, you should have the sample application with ZITADEL for user access management. However, ZITADEL doesn't have any users other than the admin user yet.

<img src="/docs/img/guides/migrate/keycloak-37.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-38.png" alt="Migrating users from Keycloak to ZITADEL"/>
