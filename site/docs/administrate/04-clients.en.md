---
title: Clients
---

### What are clients

Clients are applications who share the same security context and interface with an "authorization server".
For example you could have a software project existing out of a web app and a mobile app, both of these applications might consume the same roles because the end user might use both of them.

### Manage clients

Clients might use different protocols for integrating with an IAM. With ZITADEL it is possible to use OpenID Connect 1.0 / OAuth 2.0. In the future SAML 2.0 support is planned as well.

#### OIDC Configuration

> Document Settings

### Create a client

To make configuration of a client easy we provide a wizard which generates a specification conferment setup.
The wizard can be skipped for people who are needing special settings.

> For use cases where your configuration is not compliant we provide you a "dev mode" which disables conformance checks.

To create a new client start by browsing to your [project](administrate#Projects), this is normally something like [https://console.zitadel.ch/projects/78562301657017889](https://console.zitadel.ch/projects/78562301657017889)

<img src="img/console_projects_my_first_project.png" alt="Manage Clients" width="1000px" height="auto">

Click the **New** button and a wizard will appear which will guide you through the process.

<img src="img/console_clients_my_first_spa_wizard_1.png" alt="Client Wizard" width="1000px" height="auto">

<img src="img/console_clients_my_first_spa_wizard_2.png" alt="Client Wizard" width="1000px" height="auto">

<img src="img/console_clients_my_first_spa_wizard_3.png" alt="Client Wizard" width="1000px" height="auto">

<img src="img/console_clients_my_first_spa_wizard_4.png" alt="Client Wizard" width="1000px" height="auto">

When the wizard is complete, the clients configuration will be displayed and you can now use this client.

<img src="img/console_clients_my_first_spa_config.png" alt="Client Wizard" width="1000px" height="auto">