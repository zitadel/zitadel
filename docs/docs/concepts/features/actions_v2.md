---
title: ZITADEL Actions v2
sidebar_label: Actions v2
---

By using ZITADEL actions V2, you can manipulate ZITADELs behavior on specific API calls, events or functions.
This is useful when you have special business requirements that ZITADEL doesn't support out-of-the-box.

:::info
We're working on Actions continuously. In the [roadmap](https://zitadel.com/roadmap), you see how we are planning to expand and improve it. Please tell us about your needs and help us prioritize further fixes and features.
:::

:::warning
To use Actions v2 activate the feature flag "Actions" [feature flag](/docs/apis/resources/feature_service_v2/feature-service-set-instance-features), to be able to manage the related resources.

The Actions v2 will always be executed if available, even if the feature flag is switched off, to remove any Actions v2 the related Execution has to be removed.
:::

## Why actions?
ZITADEL can't anticipate and solve every possible business rule and integration requirements from all ZITADEL users. Here are some examples:
- A business requires domain specific data validation before a user can be created or authenticated.
- A business needs to automate tasks. Roles should be assigned to users based on their ADFS 2016+ groups.
- A business needs to store metadata on a user that is used for integrating applications.
- A business needs to restrict the users who are allowed to register to a certain organization by their email domains.

With actions, ZITADEL provides a way to solve such problems.

## How it works
There are 3 components necessary:
- Endpoint, an external endpoint with the desired logic, can be whatever is necessary as long as it can receive an HTTP Post request.
- Target, a resource in ZITADEL with all necessary information how to trigger an endpoint
- Execution, a resource in ZITADEL with the information when to trigger which targets

The process is that ZITADEL decides at certain points that with a defined Execution a call to the defined Target(s) is triggered, 
so that everybody can implement their custom behaviour for as many processes as possible.

Possible conditions for the Execution:
- Request, to react to or manipulate requests to ZITADEL, for example add information to newly created users
- Response, to react to or manipulate responses to ZITADEL, for example to provision newly created users to other systems
- Function, to react to different functionality in ZITADEL, replaces [Actions](/concepts/features/actions).
- Event, to create to different events which get created in ZITADEL, for example to inform somebody if a user gets locked

:::info
Currently, the defined Actions v2 will be executed additionally to the defined [Actions](/concepts/features/actions).
:::

## Further reading

- [Actions v2 reference](/apis/actions/v3/usage)
- [Actions v2 example execution locally](/apis/actions/v3/testing-locally)