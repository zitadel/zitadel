---
title: Use ZITADEL in AWS as IDP
---

Prerequisites on the side of AWS are found [here](https://docs.aws.amazon.com/singlesignon/latest/userguide/prereqs.html).

First, to use any external IDP to log on to the AWS Console you need to activate the AWS SSO as described [here](https://docs.aws.amazon.com/singlesignon/latest/userguide/step1.html?icmpid=docs_sso_console).

As you have activated SSO you still have the possibility to use AWS itself to manage the users, but you can also use an Microsoft AD or an external IDP.

Described [here](https://docs.aws.amazon.com/singlesignon/latest/userguide/manage-your-identity-source-idp.html) how you can connect to ZITADEL as a SAML2 IDP,
you chose the External identity provider:
![Choose identity source](images/aws_change_idp.png)

To provide ZITADEL with all the information it needs you have to download the metadata file, which we need afterwards, and save the AWS SSO SIgn-in URL, which you use to login afterwards.

To provide AWS with all the information it needs you can fill out the fields as follows:
![Configure external identity provider](images/aws_configure_idp.png)
To connect to another environment, change the domains, for example if you would use ZITADEL under the domain "example.com" you would have the URLs "https://accounts.example.com/saml/SSO" and "https://accounts.exmaple.com/saml/metadata".

Last part of this step, you have to download the ZITADEL-used certificate to sign the responses, so that AWS can validation the signature.

You can download the certificate from following URL: [https://accounts.zitadel.ch/saml/certificate](https://accounts.zitadel.ch/saml/certificate)
Then just upload the ".crt"-file to AWS and click "next".

Lastly, you only have to accept to confirm the change and ZITADEL is used as the external identity provider for AWS SSO to provide connectivity to your AWS Accounts.

As for how the SSO users are then connected to the AWS accounts, you can find more information in the AWS documentation, for example [here](https://docs.aws.amazon.com/singlesignon/latest/userguide/useraccess.html).

The result, you can now login to you AWS account through your ZITADEL-login with the AWS SSO Sign-in URL, which you should have saved in a step before.