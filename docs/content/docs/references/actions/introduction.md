---
title: Actions
---

This page describes the options you have when writing ZITADEL actions scripts.

## Language

ZITADEL interpretes the scripts as JavaScript.
Make sure your scripts are ECMAScript 5.1(+) compliant.
Go to the [goja GitHub page](https://github.com/dop251/goja) for detailed reference about the underlying library features and limitations.

Stuck customizing ZITADEL actions? Find samples for setting OIDC claims, SAML attributes, extending JIT provisioning data, calling external APIs, and more in [this repository](https://github.com/zitadel/actions).

Actions are a key feature to extend the functionality of ZITADEL and continuously improve the feature and expand the use cases. Check out our [roadmap](https://zitadel.com/roadmap) for more details.

## Action

The action object describes the logic defined by the user. Actions can be linked to multiple [triggers](#flows).

The name of the action must corelate with the javascript function in the code section. This function will be called by ZITADEL.

For reading and mutating state, the runtime executes the function that has the same name as the action.
The function receives the JavaScript objects `ctx` and `api`.

The object `ctx` provides readable information as object properties and by callable functions.
The object `api` provides mutable properties and state mutating functions.

The script of an action called **doSomething** should have a function called `doSomething` and look something like this:

```js
function doSomething(ctx, api){
    // read from ctx and manipulate with api
}
```

## Flows

Flows are the links between an [action](#action) and a specific point during a user interaction with ZITADEL. These specific point are called [Trigger Types](#trigger-types).

## Trigger Types

Trigger types define the point during execution of request. Each trigger defines which readable information (`ctx`) and mutual properties (`api`) are passed into the called function as well as which libraries are available for the function. Each trigger type is described in the [available flow types section](#available-flow-types)

## Available Flow Types

Currently ZITADEL provides the following flows:

- [Internal Authentication](./internal-authentication)
- [External Authentication](./external-authentication)
- [Complement Token](./complement-token)
- [Customize SAML Response](./customize-samlresponse)

## Available Modules inside Javascript

- [HTTP module](./modules#http) to call API's
- [Logging module](./modules#log) logs information to stdout
- [UUID module](./modules#uuid) generates uuids
