---
title: actions
---

currently support 3 triggers in the external authentication flow:

post authentication (user has authenticated externally and we retrieved / mapped the information as far as possible)
pre creation (before user is "autoregistered" / clicks on register on overview page after external authentication)
post creation (after the user was created from an external authentication)
context (information we provide you with):

accessToken (string for opaque and JWT)
idToken (string)
getClaim(string) interface{}: function which returns the requested claim
claimsJSON(): function which returns the complete payload of the id_token
api (information you can manipulate):

setFirstName(string)
setLastName(string)
setNickName(string)
setDisplayName(string)
setPreferredLanguage(string)
setGender(Gender)
setUsername(string) (currently only available on trigger pre creation!)
setPreferredUsername(string) (currently only available on trigger post authentication!)
setEmail(string)
setEmailVerified(bool)
setPhone(string)
setPhoneVerified(bool)
metadata (array of Metadata, push new entry)
userGrants (array of UserGrant, push new entry) (currently only available on trigger post creation!)
