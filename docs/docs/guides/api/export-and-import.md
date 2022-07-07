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
* PAT's
* Application keys

Which results in that if you want to import, and you have no defined organization-specific custom policies, the experience for your users will not be exactly like in your old instance.

### Export from V1

***To use this requests you have to have an access token with enough permissions to export and import.***

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
    "with_otp": true
}' -o export.json
```

* "org_ids": to select which organizations should be exported
* "excluded_org_ids": to exclude several organization, if for example no organizations are selected
* "with_passwords": to include the hashed_passwords of the users in the export 
* "with_otp": to include the OTP-code of the users in the export 

### Import to V2 from V1-data

***To import the exported data into you new instance, you have to have an already existing instance on a ZITADEL V2, with all desired configuration and global resources.***

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

