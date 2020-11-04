---
title:  Server Side Application
description: ...
---

### SSA Protocol and Flow recommendation

If your [client](administrate#Clients) is a single page application (SPA) we recommend that you use [Authorization Code](documentation#Authorization_Code) in combination with [Proof Key for Code Exchange](documentation#Proof_Key_for_Code_Exchange).

This flow has great support with most modern languages and frameworks and is the recommended default.

> In the OIDC and OAuth world this **client profile** is called "user-agent-based application"

---

With ZITADEL you can manage the [roles](administrate#Roles) a [project](administrate#Projects) supplies to your users in the form of authorizations.
On the [project](administrate#Projects) it can be configured how **project roles** are supplied to the [clients](administrate#Clients).
By default ZITADEL asserts the claim **urn:zitadel:iam:org:project:roles** to the [Userinfo Endpoint](documentation#userinfo_endpoint)

- Assert the claim **urn:zitadel:iam:org:project:roles** to **access_token**
- Assert the claim **urn:zitadel:iam:org:project:roles** to **id_token**

```JSON
 "urn:zitadel:iam:org:project:roles": {
    "user": {
      "id1": "acme.zitadel.ch",
      "id2": "caos.ch",
    }
  }
```
---

#### ASP.net core example

> Link here

