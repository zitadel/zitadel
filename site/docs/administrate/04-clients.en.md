---
title: Clients
---

### What are clients

Clients are applications that share the same security context and interface with an "authorization server" (issuer of access tokens).

For example you could have a software project existing out of a web app and a mobile app, both of these applications might consume the same roles because the end user might use both of them.

Typical types of applications are: 
* Web
* User Agent (Single-Page-Application)
* Native

Check out our [Integration Guide](integrate#Overview) for more information.

### Manage clients

Clients might use different protocols for integrating with an IAM. With ZITADEL it is possible to use OpenID Connect 1.0 / OAuth 2.0. In the future SAML 2.0 support is planned as well.

#### OIDC Configuration

> Document Settings

### Create a client

To make configuration of a client easy we provide a wizard which generates a specification conferment setup.
The wizard can be skipped for people who are needing special settings.

> For use cases where your configuration is not compliant we provide you a "dev mode" which disables conformance checks.

To create a new client start by browsing to your [project](administrate#Projects), this is normally something like [https://console.zitadel.ch/projects/78562301657017889](https://console.zitadel.ch/projects/78562301657017889)

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_projects_my_first_project.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_projects_my_first_project.png" itemprop="thumbnail" alt="Manage Clients" />
        </a>
        <figcaption itemprop="caption description">Manage Clients</figcaption>
    </figure>
</div>

Click the **New** button and a wizard will appear which will guide you through the process.

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_clients_my_first_spa_wizard_1.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_clients_my_first_spa_wizard_1.png" itemprop="thumbnail" alt="Client Wizard 1" />
        </a>
        <figcaption itemprop="caption description">Client Wizard 1</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_clients_my_first_spa_wizard_2.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_clients_my_first_spa_wizard_2.png" itemprop="thumbnail" alt="Client Wizard 2" />
        </a>
        <figcaption itemprop="caption description">Client Wizard 2</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_clients_my_first_spa_wizard_3.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_clients_my_first_spa_wizard_3.png" itemprop="thumbnail" alt="Client Wizard 3" />
        </a>
        <figcaption itemprop="caption description">Client Wizard 3</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_clients_my_first_spa_wizard_4.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_clients_my_first_spa_wizard_4.png" itemprop="thumbnail" alt="Client Wizard 4" />
        </a>
        <figcaption itemprop="caption description">Client Wizard 4</figcaption>
    </figure>
</div>

When the wizard is complete, the clients configuration will be displayed and you can now use this client.

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_clients_my_first_spa_config.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_clients_my_first_spa_config.png" itemprop="thumbnail" alt="Client Wizard Complete" />
        </a>
        <figcaption itemprop="caption description">Client Wizard Complete</figcaption>
    </figure>
</div>
