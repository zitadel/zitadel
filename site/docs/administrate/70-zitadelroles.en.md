---
title: ZITADEL Roles
---

### ZITADEL's Roles

**ZITADEL's** own role model is built around the IAM resources. The roles have some hierarchies to them. For example a IAM_OWNER can view and edit every resource of the system. ORG_OWNERS can only manage their resources included within their organisation. This includes projects, clients, users, and so on.

#### System Roles

IAM_OWNER

IAM_OWNER_VIEWER

#### Organisation Roles

ORG_OWNER

ORG_OWNER_VIEWER

ORG_USER_PERMISSION_EDITOR

ORG_PROJECT_PERMISSION_EDITOR

ORG_PROJECT_CREATOR

#### Owned Project Roles

PROJECT_OWNER

PROJECT_OWNER_VIEWER

PROJECT_OWNER_GLOBAL

PROJECT_OWNER_VIEWER_GLOBAL

#### Granted Project Roles

PROJECT_GRANT_OWNER

PROJECT_GRANT_OWNER_VIEWER

### Manage ZITADEL Roles

You can grant ZITADEL Roles directly on a resource like organisation or project. Or, if the user is in your organisation, by applying the roles to the user directly:

- [Manage Organisation ZITADEL Roles](administrate#Manage_Organisation_ZITADEL_Roles)
- [Manage Project ZITADEL Roles](administrate#Manage_Organisation_ZITADEL_Roles)
- [Manage User ZITADEL Roles](administrate#Manage_Organisation_ZITADEL_Roles)