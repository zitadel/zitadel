---
title: Saas Product with Authentication and Authorization
---

This is an example architecture for a typical saas product. 
To illustrate it a fictional organisation and project is taken.

## Example Case

The Timing Company has a product called Time.
They have two environments, the development and the production environment.
In this case the time uses authentication and authorizations from ZITADEL.
This means that the users and also their authorizations will be managed within ZITADEL.

![Architecture](/img/concepts/usecase/saas.png)

## Organisation

An organisation is the top level in ZITADEL. 
In an organisation projects and users are managed by the organisation.
You need at least one organisation for your own company in our case "The Timing Company".

For your costumers you have different possibilities:
1. Your customer already owns an organisation ZITADEL
2. Your customer creates a new organisation in ZITADEL by itself
3. You create an organisation for your customer (If you like to verify the domain, the customer has to do it)

:::note

Subscriptions are organisation based. This means, that each organisation can choose her own tier based on the needed features.

:::

## Project

The idea of projects is to have a vessel for all components who are closely related to each other.

In our use case we would have two different projects, for each environment one. So lets call it "Time Dev" and "Time Prod".
These projects should be created in "The Timing Company" organisation, because this is the owner of the project.

In the project you will configure all your roles and applications (clients and APIs).

To give a customer permissions to a project, a project grant to the customer is need (Search by domain).
It is also possible to only allow the customer to use specific roles.
As soon as a project grant exists, the customer will see the project in the granted projects section of his organisation, and will be able to authorize his own users to the given proejct.


### Project Grant


