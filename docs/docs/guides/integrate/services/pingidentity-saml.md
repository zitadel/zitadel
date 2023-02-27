---
title: Connect with Ping Identity through SAML 2.0
---

This guide shows how to enable login with ZITADEL on Auth0.

It covers how to:

- create and configure the application in your project
- create and configure the connection in your Ping Identity tenant

Prerequisites:

- existing ZITADEL Instance, if not present follow [this guide](/guides/start/quickstart)
- existing ZITADEL Organization, if not present follow [this guide](/guides/manage/console/organizations)
- existing ZITADEL project, if not present follow the first 3 steps [here](/guides/manage/console/projects)
- existing Pingidentity environment [here](https://docs.pingidentity.com/bundle/pingone/page/wqe1564020490538.html)

> We have to switch between ZITADEL and Ping Identity. If the headings begin with "ZITADEL" switch to the ZITADEL
> Console and
> if the headings start with "Ping" please switch to the PingIdentity GUI.

## **Ping**: Create a new external identity provider

To add an
additional [external identity provider](https://docs.pingidentity.com/bundle/pingone/page/jvz1567784210191.html), you
can follow the instructions [here](https://docs.pingidentity.com/bundle/pingone/page/ovy1567784211297.html)

1. As described you have to create a new provider, with a unique identifier:
   ![Create IDP Profile](/img/saml/pingidentity/create_idp_profile.png)

We recommend activating signing the auth request whenever possible:
![Configure PingOne Connection](/img/saml/pingidentity/conf_connection.png)

2. Manually enter the necessary information:

- SSO Endpoint, for example https://accounts.example.com/saml/SSO
- IDP EntityID, for example https://accounts.example.com/saml/metadata
- Binding, which is a decision which you can take yourself, we recommend HTTP POST as it has fewer restrictions
- Import certificate, provided from the certificate endpoint
  ![Configure IDP Connection](/img/saml/pingidentity/conf_idp_connection.png)

Everything you need to know about the attribute mapping you can find
in [Ping Identity's documentation](https://docs.pingidentity.com/bundle/pingone/page/pwv1567784207915.html)

3. With this you have defined to connection to ZITADEL as an external IDP, next is the policy to use ZITADEL as an IDP
   to
   connect to an application. The "How to" for that can be
   found [here](https://docs.pingidentity.com/bundle/pingone/page/zqd1616600404402.html).

## **ZITADEL**: Create the application

To add the connection to ZITADEL you have to build the metadata, which should minimalistic look like this, the necessary
information can be found on the External IDPs page under "P1Connection" and "IDP Configuration" :

```xml
ENTITYID="PINGONE (SP) ENTITY ID"
        ACSURL="ACS ENDPOINT"
        <?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="${ENTITYID}">
    <md:SPSSODescriptor
            protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol urn:oasis:names:tc:SAML:1.1:protocol">
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="${ACSURL}"
                                     index="0"/>
    </md:SPSSODescriptor>
</md:EntityDescriptor>
```

![Identity Providers P1 Connection](/img/saml/pingidentity/idp_p1_connection.png)
![Identity Providers IDP Configuration](/img/saml/pingidentity/idp_idp_configuration.png)

In your existing project:

1. Press the "+"-button to add an application
   ![Project](/img/saml/zitadel/project.png)
2. Fill in a name for the application and chose the SAML type, then click "Continue".
   ![New Application](/img/saml/zitadel/application_saml.png)
3. Either fill in the URL where ZITADEL can read the metadata from, or upload the metadata XML directly, then click "
   Continue".
   ![Add Metadata to Application](/img/saml/zitadel/application_saml_metadata.png)
4. Check your application, if everything is correct, press "Create".
   ![Create Application](/img/saml/zitadel/application_saml_create.png)

Everything on the side of ZITADEL is done if the application is correctly created.
