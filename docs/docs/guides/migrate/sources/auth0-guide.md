---
title: Migrating Users from Auth0 to ZITADEL (Including Password Hashes)
sidebar_label: Auth0 Migration Guide
---

## 1. Introduction

This guide will walk you through the steps to migrate users from Auth0 to ZITADEL, including password hashes (which requires Auth0's support assistance), so users don't need to reset their passwords.

**What you'll learn with this guide**
- How to prepare your data from Auth0
- Use of the ZITADEL migration tooling
- Performing the user import via ZITADEL's API
- Troubleshooting and validating the migration

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

## 3. Preparing Auth0 Data

### 3.1. Export User Profiles and Password Hashes from Auth0
You cannot bulk export user data from the Auth0 Dashboard. Instead, use the  [Auth0 Management API](https://auth0.com/docs/manage-users/user-migration#bulk-user-exports) or the [User Import/Export extension](https://auth0.com/docs/manage-users/user-migration/user-import-export-extension).

> **Important:** Password hashes cannot be obtained in a self-service way.  
> You must open a **support ticket** with Auth0 and request a password hash export.  
> If approved, Auth0 will provide an export containing the password hashes.  

Reference: [Export hashed passwords from Auth0](https://zitadel.com/docs/guides/migrate/sources/auth0#export-hashed-passwords)

---

## 4. Running the ZITADEL Migration Tool

### 4.1. Install the Migration Tool
Follow the installation instructions to set up the ZITADEL migration tool from [ZITADEL Tools](https://github.com/zitadel/zitadel-tools?tab=readme-ov-file#installation).

### 4.2. Generate Import JSON
Use the migration tool to convert the Auth0 export file to a ZITADEL-compatible JSON.
Step-by-step instructions: [Migration Tool for Auth0](https://github.com/zitadel/zitadel-tools/blob/main/cmd/migration/auth0/readme.md)

Typical steps:
- Run the migration tool with your exported Auth0 files as input.
- The tool generates a JSON file ready for import into ZITADEL.

Example:
After obtaining the 2 required input files (passwords and profile) in JSON lines format, you can run the following command:

Sample `passwords.ndjson` content, as obtained from the Auth0 Support team:
```json
{"_id":{"$oid":"emxdpVxozXeFb1HeEn5ThAK8"},"email_verified":true,"email":"tommie_krajcik85@hotmail.com","passwordHash":"$2b$10$d.GvZhGwTllA7OdAmsA75uGGzqr/mhdQoU88M3zD.fX3Vb8Rcf33.","password_set_date":{"$date":"2025-06-30T00:00:00.000Z"},"tenant":"test","connection":"Username-Password-Authentication","_tmp_is_unique":true}
```

Sample `profiles.json` content, as obtained from the Auth0 Management API:
```json
{"user_id":"auth0|emxdpVxozXeFb1HeEn5ThAK8","email_verified":true,"name":"Tommie Krajcik","email":"tommie_krajcik85@hotmail.com"}
```

Run the following command in your terminal (replace ORG_ID with your own organization ID):
```bash
zitadel-tools migrate auth0 --org=<ORG_ID> --users=./profiles.json --passwords=./passwords.ndjson --multiline --email-verified --output=./importBody.json --timeout=5m0s
```

The tool will merge both objects into a single one in the importBody.json output, this will be used in the next step to complete the import process.

## 5. Importing Users into ZITADEL

### 5.1. Obtain Access Token (or PAT) for API Access

To call the ZITADEL Admin API, you need to authenticate using a **Service User** with the `IAM_OWNER` Manager permissions.

There are two recommended authentication methods:

- **Client Credentials Flow**  
  [Learn how to authenticate with client credentials.](https://zitadel.com/docs/guides/integrate/service-users/client-credentials)

- **Personal Access Token (PAT)**  
  [Learn how to create and use a PAT.](https://zitadel.com/docs/guides/integrate/service-users/authenticate-service-users#personal-access-token)

**Reference:** [Service Users & API Authentication](https://zitadel.com/docs/guides/integrate/service-users/authenticate-service-users#authentication-methods)

---

### 5.2. Import Data with the ZITADEL API

- Use your **access token** or **PAT** to authenticate.
- Call the [Admin API – Import Data](https://zitadel.com/docs/apis/resources/admin/admin-service-import-data) endpoint, passing your generated JSON file.
- Verify that the users were imported successfully in the ZITADEL console.

**Import Endpoint:**

- `POST /admin/v1/import`
- `Authorization: Bearer <token>`
- **Body:** Generated in step 4.2

#### Example cURL request

```bash
curl --location 'https://<instance-domain>/admin/v1/import' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer <access-token>' \
--data-raw '{
  "dataOrgs": {
    "orgs": [
      {
        "orgId": "<your-org-id>",
        "humanUsers": [
          {
            "userId": "auth0|emxdpVxozXeFb1HeEn5ThAK8",
            "user": {
              "userName": "tommie_krajcik85@hotmail.com",
              "profile": {
                "firstName": "Tommie Krajcik",
                "lastName": "Tommie Krajcik"
              },
              "email": {
                "email": "tommie_krajcik85@hotmail.com",
                "isEmailVerified": true
              },
              "hashedPassword": {
                "value": "$2b$10$d.GvZhGwTllA7OdAmsA75uGGzqr/mhdQoU88M3zD.fX3Vb8Rcf33."
              }
            }
          }
        ]
      }
    ]
  },
  "timeout": "5m0s"
}'
```

## 6. Testing the Migration

### 6.1. Test User Login

Use the **ZITADEL login page** or your integrated app to test logging in with one of the imported users.

> **Password for the sample user:** `Password1!`

Confirm that the migrated password works as expected.

---

### 6.2. Troubleshooting

**Common issues:**

- Missing password hashes  
- Malformed JSON  
- Invalid or incomplete user data  

The import endpoint returns an `errors` array which can help you identify any issues with the import.

#### Where to check logs and get help

You can also verify that a user was imported by calling the **events endpoint** and checking for the following event type:

```json
"user.human.added"
```

## 7. Q&A and Further Resources

### Real-World Scenarios & Common Questions

**Q:** What is the maximum number of users that can be imported in a single batch?  
**A:** There is no hard limit on the number of users. However, there is a **timeout**.  
For **ZITADEL Cloud deployments**, the timeout is **5 minutes**, which typically allows for importing around **5,000 users per batch**.

---