---
title: Export and import with ZITADEL
---

## Export from V1 to Import into V2

To migrate from ZITADEL V1 to V2 the API provides you with a possibility to export all resources which are under your organizations.
Currently, this doesn't include the following points:

* Global policies
* IAM members
* Global IDPs
* Global second/multi factors
* Machine keys
* Personal Access Tokens
* Application keys
* Passwordless authentication

Which results in that if you want to import, and you have no defined organization-specific custom policies, the experience for your users will not be exactly like in your old instance.

:::note 
Note that the resources will be migrated without the event stream. This means that you will not have the audit trail for the imported objects.
:::

### Use the API

To export all necessary data you only have to use one request, as an example:

```bash
curl  --request POST \
  --url {your_domain}/admin/v1/export \
  --header 'Authorization: Bearer XXXX' \
  --header 'Content-Type: application/json' \
  --data '{    
    "org_ids": [ "70669144072186707", "70671105999825752" ],
    "excluded_org_ids": [ ],
    "with_passwords": true,
    "with_otp": true,
    "timeout": "30s",
    "response_output": true
}' -o export.json
```

* "org_ids": to select which organizations should be exported
* "excluded_org_ids": to exclude several organization, if for example no organizations are selected
* "with_passwords": to include the hashed_passwords of the users in the export 
* "with_otp": to include the OTP-code of the users in the export
* "timeout": timeout of the call to export the data
* "response_output": to output the export as response to the call

:::note 
To import the exported data into you new instance, you have to have an already existing instance on a ZITADEL V2, with all desired configuration and global resources.
:::

Then as an example you can use one request for the import:

```bash
curl --request POST \
    --url {your_domain}/admin/v1/import \
    --header 'Authorization: Bearer XXXX' \
    --header 'Content-Type: application/json' \
    --data '{
        "data_orgsv1": '$(cat export.json)'
}'
```

### Use a Google Cloud Storage

:::note 
To use this requests you have to have an access token with enough permissions to export and import.
The used serviceaccount has to have at least the role "Storage Object Creator" to create objects on GCS
:::

To export all necessary data you only have to use one request which results in a file in your GCS, as an example:

```bash
curl  --request POST \
  --url {your_domain}/admin/v1/export \
  --header 'Authorization: Bearer XXXX' \
  --header 'Content-Type: application/json' \
  --data '	"{
    "org_ids":  [ "70669144072186707", "70671105999825752" ],
    "excluded_org_ids": [ ],
    "with_passwords": true,
    "with_otp": true,
    "timeout": "10m",
    "gcs_output": {
        "path": "export.json",
        "bucket": "caos-zitadel-exports",
        "serviceaccount_json": "XXXX"
    }
}'
```

* "org_ids": to select which organizations should be exported
* "excluded_org_ids": to exclude several organization, if for example no organizations are selected
* "with_passwords": to include the hashed_passwords of the users in the export
* "with_otp": to include the OTP-code of the users in the export
* "timeout": timeout for the call to export the data
* "gcs_output": to write a file into GCS as output to the call
  * "path": path to the output file on GCS
  * "bucket": used bucket for output on GCS
  * "serviceaccount_json": base64-encoded serviceaccount.json used to output the file on GCS

:::note
To import the exported data into you new instance, you have to have an already existing instance on a ZITADEL V2, with all desired configuration and global resources.
The used serviceaccount has to have at least the role "Storage Object Viewer" to read objects from GCS
:::

Then as an example you can use one request for the import:

```bash
curl --request POST \
    --url {your_domain}/admin/v1/import \
    --header 'Authorization: Bearer XXXX' \
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

* "timeout": timeout for the import task
* "data_orgsv1_gcs": to read the export from GCS directly
    * "path": path to the exported file on GCS
    * "bucket": used bucket to read from GCS
    * "serviceaccount_json": base64-encoded serviceaccount.json used to read the file from GCS

