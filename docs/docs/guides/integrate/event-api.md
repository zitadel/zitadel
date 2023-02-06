---
title: Get events from ZITADEL
---

ZITADEL leverages the power of eventsourcing, meaning every action and change within the system generates a corresponding event that is stored in the database. 
To provide you with greater flexibility and access to these events, ZITADEL has introduced an Event API. 
This API allows you to easily retrieve and utilize the events generated within the system, enabling you to integrate them into your own system and respond to specific events as they occur.

You need to give a user the [manager role](https://zitadel.com/docs/guides/manage/console/managers) IAM_OWNER_VIEWER or IAM_OWNER to access the Event API.

If you like to know more about eventsourcing/eventstore and how this works in ZITADEL, head over to our [concepts](../../concepts/eventstore/overview).
## Request Events

Call the [ListEvents](../../apis/proto/admin#listevents) enpoint in the Administration API to get all the events you need.
To further restrict your result you can add the following filters:
- sequence
- editor user id
- event types
- aggregate id
- aggregate types
- resource owner
- creation date

```bash
curl --request POST \
  --url $YOUR-DOMAIN/admin/v1/events/_search \
  --header "Authorization: Bearer $TOKEN"
```

## Get event types

To be able to filter for the different event types ZITADEL knows, you can request the [EventTypesList](../../apis/proto/admin#listeventtypes)

```bash
curl --request POST \
--url $YOUR-DOMAIN/admin/v1/events/types/_search \
--header "Authorization: Bearer $TOKEN" \
--header 'Content-Type: application/json' \
'
```

The response will give you a list of event types. The type is what the event is called in ZITADEL itself (technical).
You can also find a translation for the event to better reflect it for an enduser perspective.

The following example shows you the event types for a password check (failed/succeeded).

```bash
...
{
    "type": "user.human.password.check.failed",
    "localized": {
        "key": "EventTypes.user.human.password.check.failed",
        "localizedMessage": "Password check failed"
    }
},
{
    "type": "user.human.password.check.succeeded",
    "localized": {
        "key": "EventTypes.user.human.password.check.succeeded",
        "localizedMessage": "Password check succeeded"
    }
},
...
```

## Get aggregate types

To be able to filter for the different aggregate types (resources) ZITADEL knows, you can request the [AggregateTypesList](../../apis/proto/admin#listaggregatetypes)

```bash
curl --request POST \
  --url $YOUR-DOMAIN/admin/v1/aggregates/types/_search \
  --header "Authorization: Bearer $TOKEN" \
  --header 'Content-Type: application/json'
```

The response will give you a list of aggregate types. The type is what the aggregate is called in ZITADEL itself (technical).
You can also find a translation for the aggregae to better reflect it for an enduser perspective.

The following example shows you the aggregate type for the user.

```bash
...
{
    "type": "user",
    "localized": {
        "key": "AggregateTypes.user",
        "localizedMessage": "User"
    }
},
...
```

## Example: Get user changes since a specific date

Assuming you use ZITADEL as your single source of truth for your user data.
Now you like to react to changes on the users to update data in other your other systems.

This example shows you how to get all events from users, filtered with the creation_date (e.g since last day/hour, etc).

```bash
curl --request POST \
  --url $YOUR-DOMAIN/admin/v1/events/_search \
  --header "Authorization: Bearer $TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
	"asc": false,
	"limit": 1000,
	"creation_date": "2023-02-01T10:00:00.000000Z",
	"aggregate_types": [
		"user"
	]
}'
```

## Example: Find out when user have been authenticated

The following example shows you how you could use the events search to get all events where a token has been created.
Also we include the refresh tokens in this example to know when the user has become a new token.

```bash
curl --request POST \
  --url $YOUR-DOMAIN/admin/v1/events/_search \
  --header "Authorization: Bearer $TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
	"asc": true,
	"limit": 1000,
	"event_types": [
		"user.token.added",
		"user.refresh.token.added
	]
}'
```


## Example: Get failed login attempt

The following example shows you how you could use the events search to find out the failed login attempts of your users.
You have to include all the event types that tell you that a login attempt has failed.
In this case this are the following events:
- Password verification failed
- One-time-password (OTP) check failed (Authenticator Apps like Authy, Google Authenticator, etc)
- Universal-Second-Factor (U2F) check failed (FaceID, WindowsHello, FingerPrint, etc)
- Passwordless/Passkey check failed (FaceID, WindowsHello, FingerPrint, etc)

```bash
curl --request POST \
  --url $YOUR-DOMAIN/admin/v1/events/_search \
  --header "Authorization: Bearer $TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
	"asc": true,
	"limit": 1000,
	"event_types": [
		"user.human.password.check.failed",
		"user.mfa.otp.check.failed",
		"user.human.mfa.u2f.token.check.failed",
		"user.human.passwordless.token.check.failed"
	]
}'
```

