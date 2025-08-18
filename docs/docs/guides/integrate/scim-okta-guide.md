---
title: SCIM Provisioning from Okta
---

This guide provides step-by-step instructions to configure SCIM provisioning from Okta into ZITADEL.

## Pre-requisites:

* Access to your ZITADEL Organization with an **Org Owner** role.

* Access to your Okta Admin dashboard.

* An existing **SAML app integration** between Okta (Identity Provider) and ZITADEL (Service Provider).

## Step 1: Set Up SCIM Provisioning in ZITADEL

SCIM provisioning in ZITADEL is accomplished by authenticating a Service User with appropriate permissions.

1. **Create a Service User**:

   * Follow [this guide](https://zitadel.com/docs/guides/manage/console/users) to create a Service User within your ZITADEL Organization.

2. **Assign the Role**:

   * Grant the Service User the **Org User Manager** role. No higher managerial role is required.

3. **Choose an Authentication Method**:

   * Select one of these two supported authentication methods:

     * Personal Access Token \- PAT

     * Client Credentials Grant

4. Detailed instructions to authenticate the Service User can be found [here](https://zitadel.com/docs/guides/integrate/service-users/authenticate-service-users).

## Step 2: Set Up SCIM in Okta

Follow these precise steps to configure SCIM provisioning in Okta:

1. Log in to your Okta Admin Console.

2. Navigate to **Applications** → **Application** and the existing **SAML app** linked to ZITADEL.

3. Select the **General** tab, then choose **Edit** for **App Settings**.

4. In the **Provisioning** section, select **SCIM** and then **Save**.

<img src="/docs/img/manage/users/enable-scim-provisioning.png" alt="Enable SCIM provisioning in Okta"/> 

5. Under the **General** tab, also confirm that [Federation Broker Mode](https://help.okta.com/en-us/content/topics/apps/apps-fbm-main.htm) is disabled.

6. Click on the **Provisioning** tab, then go to the **Integration** tab and select **Edit**.

<img src="/docs/img/manage/users/select-provisioning-actions.png" alt="Select provisioning actions in Okta"/>

7. Enter the **SCIM connector base URL** using this format:

```https://${ZITADEL_DOMAIN}/scim/v2/{orgId}```
Like the example in the above image: 
```https://test-domain-bkeog4.us1.zitadel.cloud/scim/v2/322355063156684166```

*(Find more details about endpoints [here](https://zitadel.com/docs/apis/scim2#supported-endpoints)).*

8. For **Unique identifier field for users**, enter **userName**.

9. Under **Supported provisioning actions**, select ***Push New Users*** and ***Push Profile Updates***.

10. Choose your authentication method under **Authentication Mode**:

    * **HTTP Header** if using a Personal Access Token (PAT).

    * **OAuth 2** if using Client Credentials Grant.

11. Provide the authentication details according to your chosen method:

    * For **HTTP Header (PAT)**, enter the PAT token generated from ZITADEL.

    * For **OAuth 2**, provide the client credentials (Client ID, Client Secret, token URL, authorization URL).

12. Click **Test Connection Configuration** to verify the integration (optional but recommended), then click **Save**.

13. Under the **Provisioning to App** settings, enable:

    * **Create Users**

    * **Update User Attributes**

    * **Deactivate Users**

14. Click **Save** to apply these settings.

<img src="/docs/img/manage/users/provisioning-to-app.png" alt="Enable provisioning to App in Okta"/>

## Step 3: Attribute Mapping (Recommended)

Review and adjust attribute mappings in Okta as needed:

* Ensure standard attributes such as `userName`, `email`, `name.givenName`, and `name.familyName` are correctly mapped.


## Step 4: Verify SCIM Provisioning

* Assign the configured application to test users/groups in Okta.

* Verify that users are automatically provisioned into ZITADEL by checking under **Users** in your ZITADEL console.

* Validate attribute synchronization and lifecycle management (activation, updates, deactivation).

## Helpful Reference Links

* [Authenticate users with SAML](https://zitadel.com/docs/guides/integrate/login/saml)

* [ZITADEL: Creating Service Users](https://zitadel.com/docs/guides/manage/console/users)

* [ZITADEL: Service User Authentication](https://zitadel.com/docs/guides/integrate/service-users/authenticate-service-users)

* [ZITADEL SCIM 2.0 API Endpoints](https://zitadel.com/docs/apis/scim2)

* [SCIM v2.0 (Preview) docs](https://zitadel.com/docs/guides/manage/user/scim2)