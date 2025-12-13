---
title: Configure applications for your frontend and backend services and clients
---
import AppType from "../../manage/console/_application-types.mdx";

Applications are the entry point to your project.
[Users](./users) either login into one of your clients and interact with them directly or use one of your APIs.
All applications share the roles and authorizations of their [project](./projects).

## Supported application types

ZITADEL supports the following client types:

<AppType />

## Security considerations

Ensure the configuration of application settings is limited to authorized users only.

- Use [Manager roles](./managers) to limit permissions for your users to make changes to your applications
- When [granting projects](./granted_projects) to other organizations, the receiving organization can't see or change application configuration

## References

- [Configure Applications in the Console](../../manage/projects-apps/applications)
- [ZITADEL API: Applications](/docs/reference/api-v1/management/zitadel.management.v1.ManagementService.ListApps)
