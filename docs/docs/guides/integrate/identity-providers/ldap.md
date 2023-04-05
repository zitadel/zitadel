---
title: Configure LDAP as Identity Provider
sidebar_label: LDAP
---

import GeneralConfigDescription from './_general_config_description.mdx';

This guides shows you how to connect LDAP as an identity provider in ZITADEL.

:::info
In ZITADEL you can connect an Identity Provider (IdP) like LDAP to your instance and provide it as default to all organizations or you can register the IdP to a specific organization only. This can also be done through your customers in a self-service fashion.
:::

## Prerequisite

To be able to use LDAP to authenticate your users you need an LDAP server available to ZITADEL, a user with permissions to read other users information, clear defined ObjectClass and attribute restrictions for available users that can login.

## ZITADEL Configuration

### Resulting process to connect LDAP

When you wnat to use a LDAP provider in ZITADEL, the following process is followed to login:

1. ZITADEL tries to connect to the LDAP server with or without TLS depending on the configuration
2. If the connection fails, the next server in the list will be used to try again.
3. ZITADEL tries a bind with the BindDN and BindPassword to check if it's possible to proceed
4. ZITADEL does a SearchQuery to find the UserDN with the provided configuration of base, filters and objectClasses 
5. ZITADEL tries a bind with the provided loginname and password
6. LDAP attributes get mapped to ZITADEL attributes as provided by the configuration

### Create new LDAP Provider

Go to the settings of your ZITADEL instance or the organization where you like to add a new LDAP provider.
Choose the LDAP provider template.

To configure the LDAP template please fill out the following fields:

**Name**: Name of the identity provider

**Servers**: List of servers in a format of "schema://host:port", as example "ldap://localhost:389", if TLS should be used then replace "ldap" with "ldaps" with the corresponding port.

**BaseDN**: BaseDN which will be used with each request to the LDAP server

**BindDn** and **BindPassword**: BindDN and password used to connect to the LDAP for the SearchQuery, should be an admin or user with enough permissions to search for the users to login.

**Userbase**: Base used for the user, normally "dn" but can also be configured.

**User filters**: Attributes of the user which are "or"-joined in the query for the user, used value is the input of the loginname, for example if you try to login with user@example.com and filters "uid" and "email" the resulting SearchQuery contains "(|(uid=user@example.com)(email=user@example.com))" 

**User Object Classes**: ObjectClasses which are "and"-joined in the SearchQuery and the user has to have in the LDAP.

**LDAP Attributes**: Mapping of LDAP attributes to ZITADEL attributes, the ID attributes is required, the rest depends on usage of the identity provider

**StartTLS**: If this setting is enabled after the initial connection ZITADEL tries to build a TLS connection.

**Timeout**: If this setting is set all connection run with a set timeout, if it is 0s the default timeout of 60s is used.

<GeneralConfigDescription name="GeneralConfigDescription" />

![GitHub Provider](/img/guides/zitadel_ldap_create_provider.png)

### Activate IdP

Once you created the IdP you need to activate it, to make it usable for your users.

![Activate the GitHub](/img/guides/zitadel_activate_ldap.png)

## Test the setup

To test the setup use incognito mode and browse to your login page.
If you succeeded you should see a new button which should redirect you the login side on ZITADEL for LDAP.

![GitHub Button](/img/guides/zitadel_login_ldap.png)

![GitHub Login](/img/guides/zitadel_login_ldap_input.png)
