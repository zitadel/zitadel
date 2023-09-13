---
title: Migrate from Keycloak
sidebar_label: From Keycloak
---

## Migrating from Keycloak to ZITADEL

This tutorial will use [Docker installation](https://www.docker.com/) for the prerequisites. However, both Keycloak and ZITADEL offer different installation methods. As a result, this guide won't include any required production tuning or security hardening for either system. However, it's advised you follow [recommended guidelines](https://docs.zitadel.com/docs/guides/manage/self-hosted/production) before putting those systems into production.

## Set up Keycloak
### Run Keycloak

To begin setting up Keycloak, you need to refer to the official [Keycloak Docker image](https://www.keycloak.org/getting-started/getting-started-docker). You'll use it to run a development version of the Keycloak server on your local machine:


```bash
docker run -p 8081:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:22.0.1 start-dev
```

In a few seconds, Keycloak will be available at [http://localhost:8081](http://localhost:8081). Access the **Administration Console** via the username `admin` and password `admin`:

<img src="/docs/img/guides/migrate/keycloak-01.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-02.png" alt="Migrating users from Keycloak to ZITADEL"/>


### Create a realm in Keycloak 

In order to configure Keycloak as the identity provider for your application, you need to create a new realm. This will allow users and authentication resources to be isolated from any other Keycloak usage. Click on the sidebar drop-down menu and select **Create Realm**. Then input the desired realm name and click **Create**:

<img src="/docs/img/guides/migrate/keycloak-03.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-04.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-05.png" alt="Migrating users from Keycloak to ZITADEL"/>


### Create an OAuth2/OIDC client in Keycloak

After you create the new realm, you need to set up a new client. A client is a core concept in the [OAuth 2 protocol](https://oauth.net/2/grant-types/) used in all flows. Here, it will be used by your application to connect to Keycloak and complete authentication lookups.

On the menu on the left, select **Clients** and click the button **Create client**. 

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

### Create user in Keycloak

The last thing you need to do in Keycloak is to create at least one new user. This user will be able to log into your application.

On the menu on the left, select **Users**, and click **Add user**. Fill in the username, email, and first and last names, and mark the email as verified. Click on **Create** to create a new user:

<img src="/docs/img/guides/migrate/keycloak-11.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-12.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-13.png" alt="Migrating users from Keycloak to ZITADEL"/>

Now you should attach a password to this user. Select the **Credentials** tab and click **Set password**.

<img src="/docs/img/guides/migrate/keycloak-14.png" alt="Migrating users from Keycloak to ZITADEL"/>

 On the new modal panel, input the desired password and select **Save**.

<img src="/docs/img/guides/migrate/keycloak-15.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-16.png" alt="Migrating users from Keycloak to ZITADEL"/>

### Export Keycloak users

Keycloak provides an [export](https://www.keycloak.org/server/importExport) functionality that allows user information to be extracted into JSON files. While it's intended to be used in another Keycloak instance, you can manipulate it to export users to a different user management system.

For example, in order to generate the export files with Keycloak, you will need to enter the Docker container, run the export command, and copy it outside the container:

```bash
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

## Configure the web application

Now that you have fully configured Keycloak, you need to configure an application to use Keycloak as its user management interface.

This stage focuses on ZITADEL's [sample Angular application](https://github.com/zitadel/zitadel-angular) (client-side application) as a sample application that requires a user login. This application uses [OAuth Authorization Code flow with PKCE](https://docs.zitadel.com/docs/guides/integrate/login-users#code-flow) for its authentication, which is supported by both Keycloak and ZITADEL.

In order to set up the ZITADEL sample Angular application, ensure you have [Node.js](https://nodejs.org/en/) installed and clone the example repository:

```bash
git clone https://github.com/zitadel/zitadel-angular.git
npm install -g @angular/cli
npm install
```

In this tutorial, the Keycloak realm is named `my-realm`, and the client ID is `test-client`. Edit the [src/app/app.module.ts file](https://github.com/zitadel/zitadel-angular/blob/main/src/app/app.module.ts) and update the `client ID` and `issuer`:

```js 
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

```bash
npm start
```

<img src="/docs/img/guides/migrate/keycloak-17.png" alt="Migrating users from Keycloak to ZITADEL"/>

After clicking the **Authenticate** button, users will be redirected to the login page hosted in Keycloak:

<img src="/docs/img/guides/migrate/keycloak-18.png" alt="Migrating users from Keycloak to ZITADEL"/>

After successfully logging into Keycloak, users are redirected back to the sample application:

<img src="/docs/img/guides/migrate/keycloak-19.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-20.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-21.png" alt="Migrating users from Keycloak to ZITADEL"/>

Now, you've successfully created a sample application that uses Keycloak for user and session management.

## Set up ZITADEL

After creating a sample application that connects to Keycloak, you need to set up ZITADEL in order to migrate the application and users from Keycloak to ZITADEL. For this, ZITADEL offers a [Docker Compose](https://zitadel.com/docs/self-hosting/deploy/compose) installation guide:

```bash
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

You ClientId will be displayed after you click **Create**. Make a note of it. 

<img src="/docs/img/guides/migrate/keycloak-34.png" alt="Migrating users from Keycloak to ZITADEL"/>

Click on Redirect Settings and select **Development Mode** via the toggle button. 

<img src="/docs/img/guides/migrate/keycloak-35.png" alt="Migrating users from Keycloak to ZITADEL"/>


Now that you've finalized the ZITADEL configuration for the project and application, your last required change is to modify the sample application to now use ZITADEL instead of Keycloak. Since both tools implement the same authentication flows, all you need to do is change the `issuer` URL and the `clientId`:

```js
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

```bash
npm start
```
Access [http://localhost:4200/](http://localhost:4200/). Now, you should have the sample application with ZITADEL for user access management. However, ZITADEL doesn't have any users other than the admin user yet. Try out the application with the ZITADEL admin user for now. 

<img src="/docs/img/guides/migrate/keycloak-36.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-37.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-38.png" alt="Migrating users from Keycloak to ZITADEL"/>


## Import Keycloak users into ZITADEL

While moving the functionality from Keycloak to ZITADEL might require only minimal changes, it's important to note that you may also want to move the users that already exist in Keycloak to ZITADEL.

As explained in this [ZITADEL user migration guide](https://zitadel.com/docs/guides/migrate/users), you can import users individually or in bulk. Since we are looking at importing a single user from Keycloak, migrating that  individual user to ZITADEL can be done with the [ImportHumanUser](https://zitadel.com/docs/apis/resources/mgmt/management-service-import-human-user) endpoint. 

> With this endpoint, an email will only be sent to the user if the email is marked as not verified or if there's no password set.

### Create a service user to consume ZITADEL API

But first of all, in order to use this ZITADEL API, you need to create a [service user](https://zitadel.com/docs/guides/integrate/serviceusers#exercise-create-a-service-user). 

Go to the **Users** menu and select the **Service Users** tab. And click the **+ New** button. 

<img src="/docs/img/guides/migrate/keycloak-39.png" alt="Migrating users from Keycloak to ZITADEL"/>

Fill in the details of the service user and click **Create**. 

<img src="/docs/img/guides/migrate/keycloak-40.png" alt="Migrating users from Keycloak to ZITADEL"/>

Your service user is now created and listed. 

<img src="/docs/img/guides/migrate/keycloak-41.png" alt="Migrating users from Keycloak to ZITADEL"/>

### Provide 'Org Owner' permissions to the service user

This service user needs to have elevated permissions in order to import users. For this example, you should make the service user an organization owner as explained in [this guide](https://zitadel.com/docs/guides/integrate/access-zitadel-apis#add-org_owner-to-service-user). 

Let's change the permissions as follows: 

Click on the button shown in the image below:

<img src="/docs/img/guides/migrate/keycloak-42.png" alt="Migrating users from Keycloak to ZITADEL"/>

Next, select your service user that you created and select the **Org Owner** checkbox to assign the permissions of an organization owner to the service user.  

<img src="/docs/img/guides/migrate/keycloak-43.png" alt="Migrating users from Keycloak to ZITADEL"/>

### Generate an access token for the service user

In order for the service user to access the API, they must be able to authenticate themselves. To authenticate the user, you can use either [JWT with Private Key](https://docs.zitadel.com/docs/guides/integrate/serviceusers#authenticating-a-service-user) flow (recommended for production) or [Personal Access Tokens](https://zitadel.com/docs/guides/integrate/pat)(PAT). In this tutorial we will choose the latter. 

Go to **Users** -> **Service Users** again and click on the service user, then select **Personal Access Tokens** on the left and click the **+ New** button. Copy the generated personal access token to use it later.  Click **Close** after copying the PAT. 

<img src="/docs/img/guides/migrate/keycloak-44.png" alt="Migrating users from Keycloak to ZITADEL"/>

### Import user to ZITADEL via ZITADEL API

if your Keycloak Realm has a single user, your `my-realm-users-0.json` file, into which you exported your Keycloak user previously, will look like this:

```js
{
  "realm" : "my-realm",
  "users" : [ {
    "id" : "826731b2-bf17-4bd9-b45c-6a26c76ddaae",
    "createdTimestamp" : 1693887631918,
    "username" : "test-user",
    "enabled" : true,
    "totp" : false,
    "emailVerified" : true,
    "firstName" : "John",
    "lastName" : "Doe",
    "email" : "test-user@mail.com",
    "credentials" : [ {
      "id" : "c3f3759e-9d8a-4628-aad9-09e66f28a4e2",
      "type" : "password",
      "userLabel" : "My password",
      "createdDate" : 1693888572700,
      "secretData" : "{\"value\":\"ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg=\",\"salt\":\"RaXjs4RiUKgJGkX6kp277w==\",\"additionalParameters\":{}}",
      "credentialData" : "{\"hashIterations\":27500,\"algorithm\":\"pbkdf2-sha256\",\"additionalParameters\":{}}"
    } ],
    "disableableCredentialTypes" : [ ],
    "requiredActions" : [ ],
    "realmRoles" : [ "default-roles-my-realm" ],
    "notBefore" : 0,
    "groups" : [ ]
  } ]
}
```

Now, you need to transform the JSON to the ZITADEL data format by adhering to the ZITADEL API [specification](https://docs.zitadel.com/docs/apis/proto/management#importhumanuser) to import a user. The minimal format would be as shown below: 

```js
{
    "userName": "test-user",
    "email": {
        "isEmailVerified": true,
        "email":  "test-user@mail.com"
    },
    "profile": {
        "firstName": "John",
        "lastName": "Doe"
    },
    "hashedPassword": {
        "value": "ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg=",
        "algorithm": "pbkdf2-sha256",
        "salt": "RaXjs4RiUKgJGkX6kp277w==",
        "hashIterations": "27500"
    }
}
```
You can generate a transformed JSON file with the following script (install jq before running the script): 

```bash
#!/bin/bash

# Check if jq is installed
if ! command -v jq &> /dev/null
then
    echo "jq could not be found. Please install jq first."
    exit
fi

# Convert the JSON
transformed=$(jq '
{
    "userName": .users[0].username,
    "email": {
        "isEmailVerified": .users[0].emailVerified,
        "email": .users[0].email
    },
    "profile": {
        "firstName": .users[0].firstName,
        "lastName": .users[0].lastName
    },
    "hashedPassword": {
        "value": (.users[0].credentials[0].secretData | fromjson).value,
        "algorithm": (.users[0].credentials[0].credentialData | fromjson).algorithm,
        "salt": (.users[0].credentials[0].secretData | fromjson).salt,
        "hashIterations": (.users[0].credentials[0].credentialData | fromjson).hashIterations | tostring
    }
}' my-realm-users-0.json)

echo "$transformed" > zitadel-users-file.json
```

If you want to run the provided bash script in one go, you can follow these steps:

1. Create a new file with any name, for example `keycloak-to-zitadel.sh`, and paste the script into it.

2. Make the script executable by running the following command to give the script execute permissions:

    ```bash
    chmod +x keycloak-to-zitadel.sh
    ```

3. Once the script has execute permissions, you can run it using:

   ```bash
   ./keycloak-to-zitadel.sh
   ```
    <img src="/docs/img/guides/migrate/keycloak-45.png" alt="Migrating users from Keycloak to ZITADEL"/>


Ensure `my-realm-users-0.json` is in the same directory for the script to process it, or modify the script paths accordingly if you store the files in different locations.

Now that we have the user details in the required JSON format, let’s call the ZITADEL API to add the user. 

Run the following cURL command to invoke the API and don't forget to replace `<service user access token>` with the service user's personal access token 

```bash
curl --request POST \
 --url http://localhost:8080/management/v1/users/human/_import \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <service user access token>' \
 --data @zitadel-users-file.json
```

A successful response would be as shown below: 

<img src="/docs/img/guides/migrate/keycloak-46.png" alt="Migrating users from Keycloak to ZITADEL"/>

> Note that the previous request imports a single user. If you're using ZITADEL Cloud and have a large number of users, you may hit its rate limit or may need to pay the excess number of API requests. If you experience this, reach out to the [ZITADEL support team](https://zitadel.com/contact), as they can provide an alternative migration tools to move a large number of users.


Now you have imported the Keycloak user into ZITADEL. To view your user go to [http://localhost:8080/ui/console/users](http://localhost:8080/ui/console/users) (or go to the **Users** tab to see the users). 

<img src="/docs/img/guides/migrate/keycloak-47.png" alt="Migrating users from Keycloak to ZITADEL"/>


