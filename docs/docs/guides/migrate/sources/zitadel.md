---
title: Migrate from ZITADEL
sidebar_label: From ZITADEL
---

This guide explains how to migrate from ZITADEL, this includes

- ZITADEL Cloud to self-hosted
- ZITADEL self-hosted to ZITADEL Cloud
- ZITADEL v1 (deprecated) to ZITADEL v2.x

## Considerations

The following scripts don't include:

- Global policies
- IAM members
- Global IDPs
- Global second/multi factors
- Machine keys
- Personal Access Tokens
- Application keys
- Passwordless authentication

Which results in that if you want to import, and you have no defined organization-specific custom policies, the experience for your users will not be exactly like in your old instance.

:::note
Note that the resources will be migrated without the event stream. This means that you will not have the audit trail for the imported objects.
:::

## Authorization

You need a PAT from a service user with IAM Owner permissions in both the source and target system.

### Source system

1. Go to your default organization
2. Create a service user "import_user" with Access Token Type "Bearer"
3. Create a [personal access token](/docs/guides/integrate/service-users/personal-access-token)
4. Go to the Default settings
5. Add the import_user as [manager](/docs/guides/manage/console/managers) with the role "IAM Owner"

Save the PAT to the environment variabel `PAT_EXPORT_TOKEN` and the source domain as `ZITADEL_EXPORT_DOMAIN` to run the following scripts.

### Target system

1. Go to your default organization
2. Create a service user "export_user" with Access Token Type "Bearer"
3. Create a [personal access token](/docs/guides/integrate/service-users/personal-access-token)
4. Go to the Default settings
5. Add the export_user as [manager](/docs/guides/manage/console/managers) with the role "IAM Owner"

Save the PAT to the environment variabel `PAT_IMPORT_TOKEN` and the source domain as `ZITADEL_IMPORT_DOMAIN` to run the following scripts.

:::warning Clean-up
You should let the PAT expire as soon as possible.
Make sure to delete the created users after you are done with the migration.
:::

## Use file

### Export to file

To export all necessary data you only have to use one request, as an example:

```bash
curl  --request POST \
  --url $ZITADEL_EXPORT_DOMAIN/admin/v1/export \
  --header "Authorization: Bearer $PAT_EXPORT_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
    "org_ids": [ ],
    "excluded_org_ids": [ ],
    "with_passwords": true,
    "with_otp": true,
    "timeout": "30s",
    "response_output": true
}' -o export.json
```

| Field            | Type            | Description                                                                                                                                                        |
| ---------------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| org_ids          | list of strings | provide a list of organizationIDs to select which organizations should be exported (eg, `[ "70669144072186707", "70671105999825752" ]`); leave empty to export all |
| excluded_org_ids | list of strings | to exclude several organization, if for example no organizations are selected                                                                                      |
| with_passwords   | bool            | to include the hashed_passwords of the users in the export                                                                                                         |
| with_otp         | bool            | to include the OTP-code of the users in the export                                                                                                                 |
| timeout          | duration string | timeout of the call to export the data                                                                                                                             |
| response_output  | bool            | to output the export as response to the call                                                                                                                       |

### Import from file

:::note
To import the exported data into you new instance, you have to have an already existing instance on a ZITADEL, with all desired configuration and global resources.
:::

Then as an example you can use one request for the import:

```bash
curl --request POST \
    --url $ZITADEL_IMPORT_DOMAIN/admin/v1/import \
    --header "Authorization: Bearer $PAT_IMPORT_TOKEN" \
    --header 'Content-Type: application/json' \
    --data '{
        "timeout": "10m",
        "data_orgsv1": '"$(cat export.json)"'
}'
```

| Field       | Type            | Description                             |
| ----------- | --------------- | --------------------------------------- |
| timeout     | duration string | timeout of the call to import the data  |
| data_orgsv1 | string          | data which was exported from ZITADEL V1 |

## Use Google Cloud Storage

### Export to GCS

:::note
To use this requests you have to have an access token with enough permissions to export and import.
The used serviceaccount has to have at least the role "Storage Object Creator" to create objects on GCS
:::

To export all necessary data you only have to use one request which results in a file in your GCS, as an example:

```bash
curl  --request POST \
  --url $ZITADEL_EXPORT_DOMAIN/admin/v1/export \
  --header "Authorization: Bearer $PAT_EXPORT_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
    "org_ids": [ ],
    "excluded_org_ids": [ ],
    "with_passwords": true,
    "with_otp": true,
    "timeout": "30s",
    "gcs_output": {
        "path": "export.json",
        "bucket": "caos-zitadel-exports",
        "serviceaccount_json": "XXXX"
    }
}' -o export.json
```

| Field            | Type                    | Description                                                                                                                                                        |
| ---------------- | ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| org_ids          | list of strings         | provide a list of organizationIDs to select which organizations should be exported (eg, `[ "70669144072186707", "70671105999825752" ]`); leave empty to export all |
| excluded_org_ids | list of strings         | to exclude several organization, if for example no organizations are selected                                                                                      |
| with_passwords   | bool                    | to include the hashed_passwords of the users in the export                                                                                                         |
| with_otp         | bool                    | to include the OTP-code of the users in the export                                                                                                                 |
| timeout          | duration string         | timeout of the call to export the data                                                                                                                             |
| gcs_output       | object(data_orgsv1_gcs) | to write a file into GCS as output to the call                                                                                                                     |

data_orgsv1_gcs object:

| Field               | Type   | Description                                                       |
| ------------------- | ------ | ----------------------------------------------------------------- |
| path                | string | path to the output file on GCS                                    |
| bucket              | string | used bucket for output on GCS                                     |
| serviceaccount_json | string | base64-encoded serviceaccount.json used to output the file on GCS |

### Import to GCS

:::note
To import the exported data into you new instance, you have to have an already existing instance on a ZITADEL, with all desired configuration and global resources.
The used serviceaccount has to have at least the role "Storage Object Viewer" to read objects from GCS
:::

Then as an example you can use one request for the import:

```bash
curl --request POST \
    --url $ZITADEL_IMPORT_DOMAIN/admin/v1/import \
    --header "Authorization: Bearer $PAT_IMPORT_TOKEN" \
    --header 'Content-Type: application/json' \
    --data '{
        "timeout": "10m",
        "data_orgsv1_gcs": {
          "path": "export.json",
          "bucket": "caos-zitadel-exports",
          "serviceaccount_json": "XXXX"
      }
}'
```

| Field           | Type                    | Description                            |
| --------------- | ----------------------- | -------------------------------------- |
| timeout         | duration string         | timeout of the call to import the data |
| data_orgsv1_gcs | object(data_orgsv1_gcs) | to read the export from GCS directly   |

data_orgsv1_gcs object:

| Field               | Type   | Description                                                       |
| ------------------- | ------ | ----------------------------------------------------------------- |
| path                | string | path to the exported file on GCS                                  |
| bucket              | string | used bucket to read from GCS                                      |
| serviceaccount_json | string | base64-encoded serviceaccount.json used to read the file from GCS |
