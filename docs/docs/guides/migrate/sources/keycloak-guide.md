---

title: Migrating Users from Keycloak to ZITADEL (Including Password Hashes)
sidebar_label: Keycloak Migration Guide
---

## 1. Introduction

This guide will walk you through the steps to migrate users from **Keycloak** to **ZITADEL**, including password hashes, using the `zitadel-tools` CLI and the user import APIs.

**What you'll learn with this guide**

* How to export users from Keycloak
* Use of the ZITADEL migration tooling
* Performing the user import via ZITADEL's API
* Troubleshooting and validating the migration

---

## 2. Prerequisites

### 2.1. Install Go

The migration tool is written in Go. Download and install the latest version of Go from the [official Go website](https://go.dev/doc/install).

### 2.2. Create a ZITADEL Instance and Organization

You'll need a target organization in ZITADEL to import your users. You can create a new organization or use an existing one.

If you don't have a ZITADEL instance, you can [sign up for free here](https://zitadel.com) to create a new one for you.
See: [Managing Organizations in ZITADEL](https://zitadel.com/docs/guides/manage/console/organizations).

> **Note:** Copy your Organization ID (Resource ID) since you will use the id in the later steps.

---

## 3. Exporting User Data from Keycloak

### 3.1. Set up Keycloak Locally (Optional)

To run a local development Keycloak instance, use the official Docker image:

```bash
docker run -d -p 8081:8080 \
  -e KEYCLOAK_ADMIN=admin \
  -e KEYCLOAK_ADMIN_PASSWORD=admin \
  quay.io/keycloak/keycloak:22.0.1 start-dev
```

### 3.2. Export Users from Keycloak

Run the following command inside the Keycloak container to export your realm and users:

```bash
docker exec <container_name> \
  /opt/keycloak/bin/kc.sh export \
  --dir /tmp/export \
  --realm <your_realm_name> \
  --users realm_file
```

Then copy the exported file to your host machine:

```bash
docker cp <container_name>:/tmp/export/<your_realm_name>-realm.json .                                                                                                                       
```

This creates a file such as:

```
<your_realm_name>-realm.json
```

---

## 4. Running the ZITADEL Migration Tool

### 4.1. Install the Migration Tool

Follow the installation instructions to set up the ZITADEL migration tool from [ZITADEL Tools](https://github.com/zitadel/zitadel-tools?tab=readme-ov-file#installation).

### 4.2. Generate Import JSON

Use the migration tool to convert the Keycloak realm export into a ZITADEL-compatible JSON file:

```bash
zitadel-tools migrate keycloak \
  --org=<ORG_ID> \
  --realm=<your_realm_name>-realm.json \
  --output=./importBody.json \
  --timeout=5m0s \
  --multiline
```

The tool will generate `importBody.json`, which is ready for importing into ZITADEL.

---

## 5. Importing Users into ZITADEL

### 5.1. Obtain Access Token (or PAT) for API Access

To call the ZITADEL Management API, you need to authenticate using a **Service User** with the `IAM_OWNER` Manager permissions.

There are two recommended authentication methods:

* **Client Credentials Flow**
  [Learn how to authenticate with client credentials.](https://zitadel.com/docs/guides/integrate/service-users/client-credentials)

* **Personal Access Token (PAT)**
  [Learn how to create and use a PAT.](https://zitadel.com/docs/guides/integrate/service-users/personal-access-token)

**Reference:** [Service Users & API Authentication](https://zitadel.com/docs/guides/integrate/service-users/authenticate-service-users#authentication-methods)

---

### 5.2. Import Data with the ZITADEL API

Use your **access token** or **PAT** to authenticate, then call the [Management API – Human User Import](https://zitadel.com/docs/apis/resources/admin/admin-service-import-data) endpoint.

**Import Endpoint:**

* `POST /admin/v1/import`
* `Authorization: Bearer <token>`
* **Body:** Generated in step 4.2

#### Example cURL request

```bash
curl --request POST \
  --url https://<instance-domain>/admin/v1/import \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer <token>' \
  --data @importBody.json
```

Successful Response:
```bash
{
  "success": {
    "orgs": [
      {
        "orgId": "318900732864567390",
        "humanUserIds": [
          "da72ac13-6994-4498-8b27-3ff9555661b2",
          "4e987a01-34db-4393-b61c-1ce753baf69c",
          "1041d710-8a89-48f8-85b5-1ab9656190f3",
          "7b23b799-4f0f-4964-bc6d-95c534787d2c",
          "6f2f1b2f-b292-4431-932b-620124e065ec",
          "2c65045a-9de8-4d28-b686-b27bf3a70fc3",
          "aca2dd3e-689c-4ab6-b446-0990127b1e0d",
          "18a23e01-f0fe-443f-9f1c-2a8135cd22c2",
          "c49af4bf-0dbb-4994-b453-b8dd0d5006ea"
        ]
      }
    ]
  },
  "errors": [
    {
      "type": "org",
      "id": "318900732864567390",
      "message": "ID=ORG-lapo2m Message=Errors.Org.AlreadyExisting"
    }
  ]
}
```

ℹ️ Note: The above response indicates that the organization already existed, and users were successfully added. This is not an error, and you can consider the import successful as long as the HTTP status code is **200**.

---

## 6. Testing the Migration

### 6.1. Test User Login

Use the **ZITADEL login page** or your integrated app to test logging in with one of the imported users.

Confirm that the migrated password works as expected.

---

### 6.2. Troubleshooting

**Common issues:**

* Invalid Keycloak export format
* Malformed JSON
* Missing `orgId` or access token
* Timeout exceeded during import

The import API returns a detailed response with any errors encountered during the process.

#### Where to check logs and get help

You can verify that users were imported successfully by querying the **events API** and looking for the `user.human.added` event type.

Use the following request:

```bash
curl --location 'https://<instance-domain>/admin/v1/events/_search' \
--header 'Authorization: Bearer <token>' \
--header 'Content-Type: application/json' \
--data '{
  "asc": true,
  "limit": 1000,
  "event_types": [
    "user.human.added"
  ]
}'
```

This will return a list of user creation events including details such as email, username, and hashed password to help you confirm the imported data.
 
Successful Response
```bash
{
  "events": [
    {
      "type": {
        "type": "user.human.added",
        "localized": {
          "key": "EventTypes.user.human.added",
          "localizedMessage": "Person added"
        }
      },
      "payload": {
        "displayName": "test user",
        "email": "testuser@gmail.com",
        "userName": "testuser"
      },
      "aggregate": {
        "id": "da72ac13-6994-4498-8b27-3ff9555661b2",
        "resourceOwner": "318900732864567390"
      },
      "creationDate": "2025-07-22T15:16:06.364302Z"
    }
  ]
}
```

ℹ️ Note: If you see entries with "type": "user.human.added" and correct payload data, the import was successful.

---

## 7. Q&A and Further Resources

### Real-World Scenarios & Common Questions

**Q:** What is the maximum number of users that can be imported in a single batch?  
**A:** There is no hard limit on the number of users. However, there is a **timeout**.  
For **ZITADEL Cloud deployments**, the timeout is **5 minutes**, which typically allows for importing around **5,000 users per batch**.

---
