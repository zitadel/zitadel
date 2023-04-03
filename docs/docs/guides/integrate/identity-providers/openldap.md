---
title: Configure local OpenLDAP as Identity Provider
sidebar_label: Local OpenLDAP
---

This guides shows you how to connect OpenLDAP as an identity provider in ZITADEL.

:::info
In ZITADEL you can connect an Identity Provider (IdP) like LDAP to your instance and provide it as default to all organizations or you can register the IdP to a specific organization only. This can also be done through your customers in a self-service fashion.
:::

## Prerequisite

To be able to use LDAP to authenticate your users you need an LDAP server available to ZITADEL, a user with permissions to read other users information, clear defined ObjectClass and attribute restrictions for available users that can login.

## OpenLDAP Configuration

### Basic configuration 

To run LDAP locally to test it with ZITADEL please refer to [OpenLDAP](https://www.openldap.org/) with [slapd](https://www.openldap.org/software/man.cgi?query=slapd).

For a quickstart guide please refer to their [official documentation](https://www.openldap.org/doc/admin22/quickstart.html).

A basic configuration would be like this 
```
#
# See slapd.conf(5) for details on configuration options.
# This file should NOT be world readable.
#
include /usr/local/etc/openldap/schema/core.schema
include /usr/local/etc/openldap/schema/cosine.schema
include /usr/local/etc/openldap/schema/inetorgperson.schema
include /usr/local/etc/openldap/schema/nis.schema
include /usr/local/etc/openldap/schema/misc.schema

# Define global ACLs to disable default read access.

# Do not enable referrals until AFTER you have a working directory
# service AND an understanding of referrals.
#referral       ldap://root.openldap.org

pidfile         /usr/local/var/run/slapd.pid
argsfile        /usr/local/var/run/slapd.args

# Load dynamic backend modules:
modulepath      /usr/local/Cellar/openldap/2.4.53/libexec/openldap
moduleload      back_mdb.la
moduleload      back_ldap.la

# Sample security restrictions
#       Require integrity protection (prevent hijacking)
#       Require 112-bit (3DES or better) encryption for updates
#       Require 63-bit encryption for simple bind
# security ssf=1 update_ssf=112 simple_bind=64

# Sample access control policy:
#       Root DSE: allow anyone to read it
#       Subschema (sub)entry DSE: allow anyone to read it
#       Other DSEs:
#               Allow self write access
#               Allow authenticated users read access
#               Allow anonymous users to authenticate
#       Directives needed to implement policy:
# access to dn.base="" by * read
# access to dn.base="cn=Subschema" by * read
# access to *
#       by self write
#       by users read
#       by anonymous auth
#
# if no access controls are present, the default policy
# allows anyone and everyone to read anything but restricts
# updates to rootdn.  (e.g., "access to * by * read")
#
# rootdn can always read and write EVERYTHING!

#######################################################################
# MDB database definitions
#######################################################################

database        ldif
#maxsize                1073741824
suffix          "dc=example,dc=com"
rootdn          "cn=admin,dc=example,dc=com"
# Cleartext passwords, especially for the rootdn, should
# be avoid.  See slappasswd(8) and slapd.conf(5) for details.
# Use of strong authentication encouraged.
rootpw          {SSHA}6FTOTIITpkP9IAf22VjHqu4JisyBmW5A
# The database directory MUST exist prior to running slapd AND
# should only be accessible by the slapd and slap tools.
# Mode 700 recommended.
directory       /usr/local/var/openldap-data
# Indices to maintain
#index  objectClass     eq
```

Which is the default configuration with an admin user under the DN `cn=admin,dc=example,dc=com` and password `Password1!`, BaseDN `"dc=example,dc=com` and database set to `ldif`.
In addition, there are some schemas included which can be used to create the users.

### Example users

For a basic structure and an example user you can use this structure in a `.ldif` file:
```
dn: dc=example,dc=com
dc: example
description: Company
objectClass: dcObject
objectClass: organization
o: Example, Inc.

dn: ou=people, dc=example,dc=com
ou: people
description: All people in organisation
objectclass: organizationalunit

dn: cn=test,ou=people,dc=example,dc=com
objectclass: inetOrgPerson
cn: testuser
sn: test
uid: test
userpassword: {SHA}qUqP5cyxm6YcTAhz05Hph5gvu9M=
mail: test@example.com
description: Person
ou: Human Resources
```

Which in essence creates a user with DN `cn=test,ou=people,dc=example,dc=com`, uid `test` and password `test`.

The user can be applied after OpenLDAP is running with 
```bash
ldapadd -x -W -h localhost -D "cn=admin,dc=example,dc=com" -f example.ldif -w 'Password1!'
```

## ZITADEL Configuration

### Create new LDAP Provider

Go to the settings of your ZITADEL instance or the organization where you like to add a new LDAP provider.
Choose the LDAP provider template.

To configure the LDAP template to work with the before configured OpenLDAP, please fill out the following fields:

**Name**: OpenLDAP

**Servers**: "ldap://localhost:389"

**BaseDN**: "dc=example,dc=com"

**BindDn**: "cn=admin,dc=example,dc=com"

**BindPassword**: "Password1!"

**Userbase**: "dn"

**User filters**: "uid"

**User Object Classes**: "inetOrgPerson"

**LDAP Attributes**: id attributes = "uid"

**Automatic creation**: If this setting is enabled the user will be created automatically within ZITADEL, if it doesn't exist.

**Automatic update**: If this setting is enabled, the user will be updated within ZITADEL, if some user data are changed withing the provider. E.g if the lastname changes on the GitHub account, the information will be changed on the ZITADEL account on the next login. 

**Account creation allowed**: This setting determines if account creation within ZITADEL is allowed or not.

**Account linking allowed**: This setting determines if account linking is allowed. (E.g an account within ZITADEL should already be existing and the when login with GitHub an account should be linked)

:::info
Either account creation or account linking have to be enabled. Otherwise, the provider can't be used.
:::

![GitHub Provider](/img/guides/zitadel_ldap_create_provider.png)

### Activate IdP

Once you created the IdP you need to activate it, to make it usable for your users.

![Activate the GitHub](/img/guides/zitadel_activate_ldap.png)

## Test the setup

To test the setup use incognito mode and browse to your login page.
If you succeeded you should see a new button which should redirect you the login side on ZITADEL for LDAP.

![GitHub Button](/img/guides/zitadel_login_ldap.png)

![GitHub Login](/img/guides/zitadel_login_ldap_input.png)
