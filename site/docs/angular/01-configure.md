---
title: Configure Zitadel Application
---

### Setup Application and get Keys

We recommend that you use [Authorization Code](architecture#Authorization_Code) in combination with [Proof Key for Code Exchange](architecture#Proof_Key_for_Code_Exchange) for all web applications.
This flow has great support with most modern languages and frameworks and is the recommended default.

> In the OIDC and OAuth world this **client profile** is called "user-agent-based application"

Go to [ZITADEL Console Projects](https://console.zitadel.ch/projects), select your project, and click on `new application`. 
Enter a name for your new Application, select Web as Type, Basis as Auth Method, and then proceed.

### Redirect URLs

A redirect URL is a URL in your application where ZITADEL redirects the user after they have authenticated. Set your url, add an optional redirect after logout and then proceed and create.

> If you are following along with the sample project you downloaded from our Templates, you should set the Allowed Callback URL to http://localhost:4200.

Copy your client ID as it will be needed in the next step.
