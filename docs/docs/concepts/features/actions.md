---
title: Actions
---

By using ZITADEL actions, you can manipulate ZITADELs behavior on specific Events.
This is useful when you have special business requirements that ZITADEL doesn't support out-of-the-box.

:::caution
ZITADEL actions is in an early development stage.
In the [roadmap](https://zitadel.ch/roadmap), you see how we are planning to expand and improve it.
Please tell us about your needs and help us prioritize further fixes and features.
:::

## How actions are accessed
You can select the *Actions* navigation item if you
* select an organization that has the actions feature enabled,
* and are at least an *ORG_OWNER*

## How it works
Using the actions feature, *ORG_OWNERs* create a flow for each supported flow type.
Each flow type provides its own events.
You can hook into these events by assigning them an action.
An action is composed of
* a name,
* a custom JavaScript code snippet,
* an execution timeout in seconds,
* a switch that defines if its corresponding flow should fail if the action fails.

The JavaScript code has access to the two objects `ctx` and `api`.
The `ctx` object has readable context information.
The `api` object has methods for manipulating the ZITADEL state.
When actions are executed in an external authentication flow,
all api functions are used to manipulate the user that authenticates.    
