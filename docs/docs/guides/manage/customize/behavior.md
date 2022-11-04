---
title: Behavior Customization
---

In this guide, you will create a [ZITADEL action](../../../concepts/features/actions).
After users register using an external identity provider, the action assigns them a role.

## Prerequisites

Before you start, make sure you have everything set up correctly.

- You need to be at least a ZITADEL *ORG_OWNER*
- Your ZITADEL organization needs to have the actions feature enabled. <!-- TODO: How to enable it for SaaS ZITADEL? -->
- [Your ZITADEL organization needs to have at least one external identity provider enabled](../../integrate/identity-brokering)
- [You need to have at least one role configured for a project](../console/projects)

## Copy some information for the action

1. Select the **Projects** navigation item.
1. Select a project that has a role configured.
1. Copy the projects **Resource Id** on the screens top right.
1. Scroll to the **ROLES** section and note some roles key.

## Create the action

1. Select the **Actions** navigation item.
1. In the **Actions <i class="las la-code"></i>** section, select the **+ New** button.
1. Give the new action the name `addGrant`.
1. Paste this snippet into the multiline textfield.
1. Replace the snippets placeholders and select **Save**.


```js reference
https://github.com/zitadel/actions/blob/main/examples/add_user_grant.js
```

## Run the action when a user registers

Now, make the action hook into the [external authentication flow](../../../apis/actions#external-authentication-flow).

1. In the **Flows <i class="las la-exchange-alt"></i>** section, select the **+ New** button.
1. Select the **Flow Type** *External Authentication*.
1. Select the **Trigger Type** *Post Creation*.
1. In the **Actions** dropdown, check *addGrant*.
1. Select the **Save** button.

<!-- TODO: ## Test if your action works -->

New users automatically are assiged a role now if they register by authenticating with an external identity provider.

## What's next?

- [Read more about the concepts around actions](../../../concepts/features/actions)
- [Read more about all the options you have with actions](../../../apis/actions)