---
title: Google Workspace SSO with ZITADEL
sidebar_label: Google Workspace
---

This guide shows how to enable login with ZITADEL on Google Workspace.

You can configure two types of SAML SSO on Google Workspace:

- [SSO profile for your organization](#sso-profile-for-your-organization)
- [Third-party SSO SAML profile](#third-party-sso-saml-profile)

Both profiles need to be configured differently.
Please make sure to configure your application for the correct type.
Please refer to Google Help to [Set up SSO for your organization](https://support.google.com/a/answer/12032922) in case you need additional information on the Workspace setup.

:::info OpenID Connect
At this time Google supports SSO with OpenID Connect only for few providers.
:::

Prerequisites:

- You need to have a domain registered with your Google Workspace account to configure SSO profiles
- Make sure that you [verify the same domain also in your ZITADEL organization and set it as primary domain](/docs/guides/manage/console/organizations#domain-verification-and-primary-domain)
- A user in Google Workspace (eg, road.runner@acme.com)
- A user in ZITADEL with the same username (eg, road.runner@acme.com); make sure you verify the domain to set the username. This is different than the user's email address

## SSO profile for your organization

### Configure SSO profile on Google Workspace

Open the Google settings for [SSO with third-party IdP](https://admin.google.com/u/1/ac/security/sso) and click on *ADD SSO PROFILE*.

![SSO with third-party IdP](/img/guides/integrate/services/google-workspace-sso-overview.png)

Download the public certificate from your ZITADEL instance by requesting `$YOUR_DOMAIN/saml/v2/certificate`

```bash
 wget $YOUR_DOMAIN/saml/v2/certificate -O idp.crt
```

Always replace `$YOUR_DOMAIN` with your instance domain.

Use the following configuration

| Setting | Value |
| --- | --- |
| Set up SSO with third-party identity provider | Enable (check) |
| Sign-in page URL | $YOUR_DOMAIN/saml/v2/SSO |
| Sign-out page URL | $YOUR_DOMAIN/saml/v2/SLO |
| Verification Certificate | Upload the certificate (idp.crt) |
| Use a domain-specific issuer | Enable (check) |
| Network masks | Leave blank|
| Change password URL | $YOUR_DOMAIN/ui/console/users/me?id=security |

### Create a SAML application in ZITADEL

Create a new .xml file with the following minimal SAML metadata contents:

```xml
<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="${ENTITYID}">
    <md:SPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol urn:oasis:names:tc:SAML:1.1:protocol">
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="${ACSURL}" index="0"/>
    </md:SPSSODescriptor>
</md:EntityDescriptor>
```

Set or replace the variables with the values from the next screen as follows:

- `${ENTITYID}`: google.com/a/<your_domain>
- `${ACSURL}`: https://www.google.com/a/<your_domain>/acs

`<your_domain>` is the domain you have verified in Google Workspace.

In your existing project:

Press the "+"-button to add an application
![Project](/img/saml/zitadel/project.png)

Fill in a name for the application and chose the SAML type, then click "Continue".
![New Application](/img/saml/zitadel/application_saml.png)

Either fill in the URL where ZITADEL can read the metadata from, or upload the metadata XML directly, then click "Continue".
![Add Metadata to Application](/img/saml/zitadel/application_saml_metadata.png)

Check your application, if everything is correct, press "Create".
![Create Application](/img/saml/zitadel/application_saml_create.png)

### Activate the SSO profile for your organization

Make sure to enable the SSO profile for your organization.

In the [domain-specific service URLs](https://admin.google.com/u/1/ac/security/sso/domain-specific-service-urls) settings select "Automatically redirect users to the third-party IdP in the following SSO profile" and select as SSO profile "SSO profile for your organization".

![domain-specific service URLs](/img/guides/integrate/services/google-workspace-domain-sepcific-urls.png)

Save the settings.

![SSO with third-party IdP lower part](/img/guides/integrate/services/google-workspace-sso-bottom.png)

### Verify the SSO profile for your organization

Now you should be all set to verify your setup:

- Open Gmail in an incognito session with the following link: https://mail.google.com/a/<your_domain>
- Enter your username and credentials
- You should be redirected to Gmail and logged in

`<your_domain>` is the domain you have verified in Google Workspace.

## Third-party SSO SAML profile

Configure a third party SSO SAML profile and login users with ZITADEL to Google Workspace.

### Add SAML profile on Google Workspace

Open the Google settings for [SSO with third-party IdP](https://admin.google.com/u/1/ac/security/sso) and click on *ADD SAML PROFILE*.

![SSO with third-party IdP](/img/guides/integrate/services/google-workspace-sso-overview.png)

Download the public certificate from your ZITADEL instance by requesting `$YOUR_DOMAIN/saml/v2/certificate`

```bash
 wget $YOUR_DOMAIN/saml/v2/certificate -O idp.crt
```

Always replace `$YOUR_DOMAIN` with your instance domain.

Use the following configuration

| Setting | Value |
| --- | --- |
| SSO profile name | ZITADEL SSO |
| IDP entity ID | $YOUR_DOMAIN/saml/v2/metadata |
| Sign-in page URL | $YOUR_DOMAIN/saml/v2/SSO |
| Sign-out page URL | $YOUR_DOMAIN/saml/v2/SLO |
| Change password URL | $YOUR_DOMAIN/ui/console/users/me?id=security |
| Verification Certificate | Upload the certificate (idp.crt) |

Now go ahead and click *SAVE*

### Entity ID and ACS URL

Open the Google settings for [SSO with third-party IdP](https://admin.google.com/u/1/ac/security/sso) and click on the SAML Profile *ZITADEL SSO*

![SSO Profile Overview](/img/guides/integrate/services/google-workspace-sso-zitadel.png)

You can copy the "Entity ID" and "ACS URL" from the "SP details" section.

![ZITADEL SSO Profile](/img/guides/integrate/services/google-workspace-zitadel-profile-configured.png)

### Create a SAML application in ZITADEL

Create a new .xml file with the following minimal SAML metadata contents:

```xml
<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="${ENTITYID}">
    <md:SPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol urn:oasis:names:tc:SAML:1.1:protocol">
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="${ACSURL}" index="0"/>
    </md:SPSSODescriptor>
</md:EntityDescriptor>
```

Set or replace the variables with the values from the next screen as follows:

- `${ENTITYID}`: https://accounts.google.com/samlrp/metadata?rpid=<your_value>
- `${ACSURL}`: https://accounts.google.com/samlrp/acs?rpid=<your_value>

Replace `<your_value>` with the values from the [SSO profile](#entity-id-and-acs-url).

In your existing project:

Press the "+"-button to add an application
![Project](/img/saml/zitadel/project.png)

Fill in a name for the application and chose the SAML type, then click "Continue".
![New Application](/img/saml/zitadel/application_saml.png)

Either fill in the URL where ZITADEL can read the metadata from, or upload the metadata XML directly, then click "Continue".
![Add Metadata to Application](/img/saml/zitadel/application_saml_metadata.png)

Check your application, if everything is correct, press "Create".
![Create Application](/img/saml/zitadel/application_saml_create.png)

### Activate the SSO profile

Make sure to enable the SSO profile.

In the [domain-specific service URLs](https://admin.google.com/u/1/ac/security/sso/domain-specific-service-urls) settings select "Automatically redirect users to the third-party IdP in the following SSO profile" and select as SSO profile "ZITADEL SSO".

![domain-specific service URLs with ZITADEL SSO](/img/guides/integrate/services/google-workspace-zitadel-set-profile.png)

Save the settings.

![SSO with third-party IdP lower part with ZITADEL SSO](/img/guides/integrate/services/google-workspace-sso-zitadel.png)

### Verify the SAML SSO profile

Now you should be all set to verify your setup:

- Open Gmail in an incognito session with the following link: https://mail.google.com/a/<your_domain>
- Enter your username and credentials
- You should be redirected to Gmail and logged in

`<your_domain>` is the domain you have verified in Google Workspace.

### Troubleshooting

Make sure you don't use a super admin account in Google Workspace to test SSO. Super Admin users are not allowed to login with SSO and you might receive an status code 500.
