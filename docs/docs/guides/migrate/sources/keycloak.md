---
title: Migrate from Keycloak
sidebar_label: From Keycloak
---

## Migrating from Keycloak to ZITADEL

This guide will use [Docker installation](https://www.docker.com/) to run Keycloak and ZITADEL. However, both Keycloak and ZITADEL offer different installation methods. As a result, this guide won't include any required production tuning or security hardening for either system. However, it's advised you follow [recommended guidelines](https://zitadel.com/docs/guides/manage/self-hosted/production) before putting those systems into production. You can skip setting up Keycloak and ZITADEL if you already have running instances.  

## Set up Keycloak
### Run Keycloak

To begin setting up Keycloak, you need to refer to the official [Keycloak Docker image](https://www.keycloak.org/getting-started/getting-started-docker). You'll use it to run a development version of the Keycloak server on your local machine:


```bash
docker run -d -p 8081:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:22.0.1 start-dev
```

In a few seconds, Keycloak will be available at [http://localhost:8081](http://localhost:8081). Access the **Administration Console** via the username `admin` and password `admin`:

<img src="/docs/img/guides/migrate/keycloak-01.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-02.png" alt="Migrating users from Keycloak to ZITADEL"/>


### Create a realm in Keycloak 

In order to configure Keycloak as the identity provider for your application, you need to create a new realm. This will allow users and authentication resources to be isolated from any other Keycloak usage. Click on the sidebar drop-down menu and select **Create Realm**. Then input the desired realm name and click **Create**:

<img src="/docs/img/guides/migrate/keycloak-03.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-04.png" alt="Migrating users from Keycloak to ZITADEL"/>

<img src="/docs/img/guides/migrate/keycloak-05.png" alt="Migrating users from Keycloak to ZITADEL"/>


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

# Run the export command inside the Keycloak container
# use the container ID of Keycloak
docker exec  <keycloak container ID>  /opt/keycloak/bin/kc.sh export --dir /tmp

# copy generated files from docker container to local machine
docker cp <keycloak container ID>:/tmp/my-realm-users-0.json .
```

## Set up ZITADEL

After creating a sample application that connects to Keycloak, you need to set up ZITADEL in order to migrate the application and users from Keycloak to ZITADEL. For this, ZITADEL offers a [Docker Compose](https://zitadel.com/docs/self-hosting/deploy/compose) installation guide. Follow the instructions under the [Docker compose](https://zitadel.com/docs/self-hosting/deploy/compose#docker-compose) section to run a ZITADEL instance locally. 

Next, the application will be available at [http://localhost:8080/ui/console/](http://localhost:8080/ui/console/).

<img src="/docs/img/guides/migrate/keycloak-22.png" alt="Migrating users from Keycloak to ZITADEL"/>

Now you can access the console with the following default credentials:

* **Username**: `zitadel-admin@zitadel.localhost`
* **Password**: `Password1!`


## Import Keycloak users into ZITADEL

As explained in this [ZITADEL user migration guide](https://zitadel.com/docs/guides/migrate/users), you can import users individually or in bulk. Since we are looking at importing a single user from Keycloak, migrating that individual user to ZITADEL can be done with the [ImportHumanUser](https://zitadel.com/docs/apis/resources/mgmt/management-service-import-human-user) endpoint. 

> With this endpoint, an email will only be sent to the user if the email is marked as not verified or if there's no password set.

### Create a service user to consume ZITADEL API

But first of all, in order to use this ZITADEL API, you need to create a [service user](https://zitadel.com/docs/guides/integrate/service-users/authenticate-service-users#exercise-create-a-service-user). 

Go to the **Users** menu and select the **Service Users** tab. And click the **+ New** button. 

<img src="/docs/img/guides/migrate/keycloak-39.png" alt="Migrating users from Keycloak to ZITADEL"/>

Fill in the details of the service user and click **Create**. 

<img src="/docs/img/guides/migrate/keycloak-40.png" alt="Migrating users from Keycloak to ZITADEL"/>

Your service user is now created and listed. 

<img src="/docs/img/guides/migrate/keycloak-41.png" alt="Migrating users from Keycloak to ZITADEL"/>

### Provide 'Org Owner' permissions to the service user

This service user needs to have elevated permissions in order to import users. For this example, you should make the service user an organization owner as explained in [this guide](/docs/guides/integrate/zitadel-apis/access-zitadel-apis#add-org_owner-to-service-user). 

Let's change the permissions as follows: 

Click on the button shown in the image below:

<img src="/docs/img/guides/migrate/keycloak-42.png" alt="Migrating users from Keycloak to ZITADEL"/>

Next, select your service user that you created and select the **Org Owner** checkbox to assign the permissions of an organization owner to the service user.  

<img src="/docs/img/guides/migrate/keycloak-43.png" alt="Migrating users from Keycloak to ZITADEL"/>

### Generate an access token for the service user

In order for the service user to access the API, they must be able to authenticate themselves. To authenticate the user, you can use either [JWT with Private Key](/docs/guides/integrate/service-users/authenticate-service-users#authenticating-a-service-user) flow (recommended for production) or [Personal Access Tokens](/docs/guides/integrate/service-users/personal-access-token)(PAT). In this guide, we will choose the latter. 

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

Now, you need to transform the JSON to the ZITADEL data format by adhering to the ZITADEL API [specification](https://zitadel.com/docs/apis/resources/mgmt/management-service-import-human-user) to import a user. The minimal format would be as shown below: 

```js 
{
    "userName": "test-user",
    "profile": {
        "firstName": "John",
        "lastName": "Doe"
    },
    "email": {
        "email": "test-user@mail.com",
        "isEmailVerified": true
    },
    "hashedPassword": {
        "value": "$pbkdf2-sha256$27500$RaXjs4RiUKgJGkX6kp277w==$ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg="
    }
}
     
```

Next, you must install [`zitadel-tools`](https://github.com/zitadel/zitadel-tools/tree/main), which is a utility toolset designed to facilitate various interactions with the ZITADEL platform, mainly with tasks related to authentication, authorization, and data migration. We will be using the `migrate` command:

Purpose: Assists users in transforming exported data from other identity providers to be compatible with Zitadel's import schema.
Supported Providers: Currently, migrations from Auth0 and Keycloak are supported.
Usage: Users can get a list of available sub-commands and flags with the --help flag.

Install `zitadel-tools` using the command below. Ensure you have Go already installed on your machine. 

```bash
go install github.com/zitadel/zitadel-tools@main
```

Now you can run the migration tool for Keycloak as explained in this [guide](https://github.com/zitadel/zitadel-tools/blob/main/cmd/migration/keycloak/readme.md). Let's go through the steps: 

The Keycloak migration tool facilitates the transfer of data to ZITADEL by creating a JSON file tailored to serve as the body for an import request to the ZITADEL API. Note that it's essential that an organization already exists within ZITADEL/

To perform the migration, you'll need:

- The organization ID (--org)
- A realm.json file (in our case, `my-realm-users-0.json`) that houses your exported Keycloak realm with user details (--realm). 
- Output path via --output (default: ./importBody.json)
- Timeout duration for the data import request using --timeout (default: 30 minutes)
- Pretty printing the output JSON with --multiline.

Execute with: 

```bash
zitadel-tools migrate keycloak --org=<organisation id> --realm=./realm.json --output=./importBody.json --timeout=1h --multiline
```

Example: 

```bash
zitadel-tools migrate keycloak --org=233868910057750531 --realm=./my-realm-users-0.json --output=./importBody.json --timeout=1h --multiline
```

Ensure `my-realm-users-0.json` is in the same directory for the tool to process it, or provide the path to the file. 

`importBody.json` will now contain the transformed data as shown below: 

```bash
{
  "dataOrgs": {
    "orgs": [
      {
        "orgId": "233868910057750531",
        "humanUsers": [
          {
            "userId": "826731b2-bf17-4bd9-b45c-6a26c76ddaae",
            "user": {
              "userName": "test-user",
              "profile": {
                "firstName": "John",
                "lastName": "Doe"
              },
              "email": {
                "email": "test-user@mail.com",
                "isEmailVerified": true
              },
              "hashedPassword": {
                "value": "$pbkdf2-sha256$27500$RaXjs4RiUKgJGkX6kp277w==$ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg="
              }
            }
          }
        ]
      }
    ]
  },
  "timeout": "1h0m0s"
}
```

Now copy the following portion to a separate file and name the file `zitadel-users-file.json`.

```bash
{
  "userName": "test-user",
  "profile": {
    "firstName": "John",
    "lastName": "Doe"
  },
  "email": {
    "email": "test-user@mail.com",
    "isEmailVerified": true
  },
  "hashedPassword": {
    "value": "$pbkdf2-sha256$27500$RaXjs4RiUKgJGkX6kp277w==$ng6oDRung/pBLayd5ro7IU3mL/p86pg3WvQNQc+N1Eg="
  }
}
```

Now that we have the user details in the required JSON format, let’s call the ZITADEL API to add the user. 

Run the following cURL command to invoke the API and don't forget to replace `<service user access token>` with the service user's personal access token: 

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


You can now view the Keycloak user's details in ZITADEL. You can see that the password is available too. 
