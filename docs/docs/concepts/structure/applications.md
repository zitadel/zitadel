---
title: Configure applications for your frontend and backend services and clients
sidebar_label: Applications
sidebar_position: 5
---
import AppType from "../../guides/manage/console/_application-types.mdx";

Applications are the entry point to your project.
[Users](users.md) either login into one of your applications and interact with them directly or use one of your APIs.
All applications share the roles and authorizations of their [project](projects.md).

## Supported application types

ZITADEL supports the following application types:

<AppType />

## Security considerations

Ensure the configuration of application settings is limited to authorized users only.

- Use [Administrator roles](managers.mdx) to limit permissions for your users to make changes to your applications
- When [granting projects](granted_projects.md) to other organizations, the receiving organization can't see or change application configuration

## References

- [Configure Applications in the Console](../../guides/manage/console/applications-overview)
- [ZITADEL API: Applications](/docs/category/apis/resources/mgmt/applications)
