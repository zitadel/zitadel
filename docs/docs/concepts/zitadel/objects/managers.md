---
title: Managers
---

ZITADEL Managers are Users who have permission to manage ZITADEL itself. There are some different levels for managers.

- **IAM Managers**: This is the highest level. Users with IAM Manager roles are able to manage the whole IAM.
- **Org Managers**: Managers in the Organization Level are able to manage everything within the granted Organization.
- **Project Mangers**: In this level the user is able to manage a project.
- **Project Grant Manager**: The project grant manager is for projects, which are granted of another organization.

To configure managers in ZITADEL go to the resource where you like to add it (e.g IAM, Organization, Project, GrantedProject).
In the right part of the console you can finde **MANAGERS** in the details part. Here you have a list of the current managers and can add a new one.

## Roles

| Role   | Description   |
|---|---|
| IAM_OWNER  | Manage the IAM, manage all organizations with their content  |
| IAM_OWNER_VIEWER  | View the IAM and view all organizations with their content |
| IAM_ORG_MANAGER  | Manage all organizations including their policies, projects and users |
| IAM_USER_MANAGER  | Manage all users and their authorizations over all organizations |
| ORG_OWNER  | Manage everything within an organization  |
| ORG_OWNER_VIEWER  | View everything within an organization  |
| ORG_USER_MANAGER  | Manage users and their authorizations within an organization |
| ORG_USER_PERMISSION_EDITOR  | Manage user grants and view everything needed for this  |
| ORG_PROJECT_PERMISSION_EDITOR  | Grant Projects to other organizations and view everything needed for this  |
| ORG_PROJECT_CREATOR  | This role is used for users in the global organization. They are allowed to create projects and manage them.  |
| PROJECT_OWNER  | Manage everything within a project. This includes to grant users for the project.  |
| PROJECT_OWNER_VIEWER  | View everything within a project.|
| PROJECT_OWNER_GLOBAL  | Same as PROJECT_OWNER, but in the global organization. |
| PROJECT_OWNER_VIEWER_GLOBAL  | Same as PROJECT_OWNER_VIEWER, but in the global organization. |
| PROJECT_GRANT_OWNER  | Same as PROJECT_OWNER but for a granted proejct. |
