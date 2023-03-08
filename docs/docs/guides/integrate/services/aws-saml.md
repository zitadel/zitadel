---
title: Connect with AWS through SAML 2.0
---

This guide shows how to enable login with ZITADEL on AWS SSO.

It covers how to:

- create and configure the application in your project
- create and configure the connection in your AWS SSO external IDP

Prerequisites:

- existing ZITADEL Instance, if not present follow [this guide](/guides/start/quickstart)
- existing ZITADEL Organization, if not present follow [this guide](/guides/manage/console/organizations)
- existing ZITADEL project, if not present follow the first 3 steps [here](/guides/manage/console/projects)
- prerequisites on AWS side [here](https://docs.aws.amazon.com/singlesignon/latest/userguide/prereqs.html).
- enabled AWS SSO [here](https://docs.aws.amazon.com/singlesignon/latest/userguide/step1.html?icmpid=docs_sso_console)

> We have to switch between ZITADEL and a AWS. If the headings begin with "ZITADEL" switch to the ZITADEL Console and if
> the headings start with "AWS" please switch to the AWS GUI.

## **AWS**: Change to external identity provider ZITADEL

As you have activated SSO you still have the possibility to use AWS itself to manage the users, but you can also use a
Microsoft AD or an external IDP.

Described [here](https://docs.aws.amazon.com/singlesignon/latest/userguide/manage-your-identity-source-idp.html) how you
can connect to ZITADEL as a SAML2 IDP.

1. Chose the External identity provider:
   ![Choose identity source](/img/saml/aws/change_idp.png)

2. Download the metadata file, to provide ZITADEL with all the information it needs, and save the AWS SSO Sign-in URL,
   which you use to login afterwards.

3. Fill out the fields as follows, to provide AWS with all the information it needs:
   ![Configure external identity provider](/img/saml/aws/configure_idp.png)

   To connect to another environment, change the domains, for example if you would use ZITADEL under the domain "
   example.com" you would have the URLs "https://accounts.example.com/saml/SSO"
   and "https://accounts.exmaple.com/saml/metadata".

4. Download the ZITADEL-used certificate to sign the responses, so that AWS can validation the signature.

   You can download the certificate from following
   URL: {your_instance_domain}/saml/v2/certificate

5. Then upload the ".crt"-file to AWS and click "next".

6. Lastly, accept to confirm the change and ZITADEL is used as the external identity provider for AWS SSO to provide
   connectivity to your AWS Accounts.

As for how the SSO users are then connected to the AWS accounts, you can find more information in the AWS documentation,
for example [here](https://docs.aws.amazon.com/singlesignon/latest/userguide/useraccess.html).

## **ZITADEL**: Create the application

The metadata used in this part is from "Change to external identity provider ZITADEL" step 2.

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

## **AWS**: Test the connection

The result, you can now login to you AWS account through your ZITADEL-login with the AWS SSO Sign-in URL, which you
should have saved in "Change to external identity provider ZITADEL" step 2.