---
title: Saas Product with Authentication and Authorization
---

This is an example architecture for a typical SaaS product. 
To illustrate it a fictional organisation and project is taken.

## Example Case

The Timing Company has a product called Time.
They have two environments, the development and the production environment.
In this case Time uses authentication and authorizations from ZITADEL.
This means that the users and also their authorizations will be managed within ZITADEL.

![Architecture](/img/concepts/usecase/saas.png)

## Organisation

An organisation is the top level in ZITADEL. 
In an organisation projects and users are managed by the organisation.
You need at least one organisation for your own company in our case "The Timing Company".

For your customers you have different possibilities:
1. Your customer already owns an organisation ZITADEL
2. Your customer creates a new organisation in ZITADEL by itself
3. You create an organisation for your customer (If you like to verify the domain, the customer has to do it)

:::info
Subscriptions are organisation based. This means, that each organisation can choose her own tier based on the needed features.
:::

## Project

The idea of projects is to have a vessel for all components who are closely related to each other.

In our use case we would have two different projects, for each environment one. So lets call it "Time Dev" and "Time Prod".
These projects should be created in "The Timing Company" organisation, because this is the owner of the project.

In the project you will configure all your roles and applications (clients and APIs).

### Project Settings

You can configure `check roles on authentication` on the project, if you want to restrict access to users that have the correct authorization for the project.

### Project Grant

To give a customer permissions to a project, a project grant to the customer is needed (search the granted organization by its domain).
It is also possible to delegate only specific roles of the project to a customer.
As soon as a project grant exists, the customer will see the project in the granted projects section of his organisation and will be able to authorize his own users to the given proejct.

## Authorizations

To give a user permission to a project an authorization is needed.
All organisations which own the project or got it granted are able to authorize users.
It is also possible to authorize users outside the own company if the exact login name of the user is known.

## Project Login

There are some different use cases how the login should behave and look like:

1. Restrict Organisation

With the primary domain scope the organisation will be restricted to the requested domain, this means only users of the requestd organisation will be able to login.
The private labeling (branding) and the login policy of the requested organisation will trigger automatically.

:::note
More about the [Scopes](../../apis/openidoauth/scopes)
:::

2. Show private labeling (branding) of the project organisation

You can configure on project-level which branding should be shown to users.
In the default the design of ZITADEL will be shown, but as soon as the user is identified, the policy of the users organisation will be triggered.
If the setting is set to Ensure Project Resource Owner Setting, the private labeling of the project organisation will always be triggered.
The last possibility is to show the private labeling of the project organisation and as soon as the user is identitfied the user organisation settings will be triggered.
For this the Allow User Resource Owner Setting should be set.
:::note
More about [Private Labeling](../../guides/customization/branding)
:::