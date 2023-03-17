---
title: Actions
---

By using ZITADEL actions, you can manipulate ZITADELs behavior on specific Events.
This is useful when you have special business requirements that ZITADEL doesn't support out-of-the-box.

:::info
We're working on Actions continuously. In the [roadmap](https://zitadel.com/roadmap), you see how we are planning to expand and improve it. Please tell us about your needs and help us prioritize further fixes and features.
:::

## Why actions?
ZITADEL can't anticipate and solve every possible business rule and integration requirements from all ZITADEL users. Here are some examples:
- A business requires domain specific data validation before a user can be created or authenticated.
- A business needs to automate tasks. Roles should be assigned to users based on their ADFS 2016+ groups.
- A business needs to store metadata on a user that is used for integrating applications. 
- A business needs to restrict the users who are allowed to register to a certain organization by their email domains.

With actions, ZITADEL provides a way to solve such problems.

## How it works
Using the actions feature, *ORG_OWNERs* create a flow for each supported flow type.
Each flow type provides its own events.
You can hook into these events by assigning them an action.
An action is composed of
* a name,
* a custom JavaScript code snippet,
* an execution timeout in seconds,
* a switch that defines if its corresponding flow should fail if the action fails.

Within the JavaScript code, you can read and manipulate the state.

## Further reading

- [Assign users a role after they register using an external identity provider](/guides/manage/customize/behavior)
- [Actions reference](/apis/actions/introduction#action)
- [Actions Marketplace: Find example actions to use in ZITADEL](https://github.com/zitadel/actions)
