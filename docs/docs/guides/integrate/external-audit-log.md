---
title: Streaming audit logs to external systems (SIEM/SOC)
sidebar_label: Audit Logs (SIEM / SOC)
---

This document details integrating ZITADEL with external systems for streaming events and audit logs. 
This functionality allows you to centralize ZITADEL activity data alongside other security and operational information, facilitating comprehensive monitoring, analysis, and compliance reporting.

Integrating ZITADEL with external systems offers several advantages:
- **Centralized Monitoring**: Streamlining ZITADEL events into a single platform alongside data from various sources enables a consolidated view of your security posture. This comprehensive view simplifies threat detection, investigation, and incident response.
- **Enhanced Security Analytics**: External systems, such as Security Information and Event Management (SIEM) solutions, can leverage ZITADEL events to identify suspicious activities, potential security breaches, and user access anomalies.
- **Compliance Reporting**: ZITADEL events can be used to generate detailed audit trails, fulfilling regulatory compliance requirements for data access and user activity.

By integrating ZITADEL with external systems, you gain valuable insights into user behavior, system activity, and potential security threats, ultimately strengthening your overall security posture and regulatory compliance.

ZITADEL provides different solutions how to send events to external systems, the solution you choose might differ depending on your use case, your database and your environment (ZITADEL Cloud, Self-hosting).

The following table shows the available integration patterns for streaming audit logs to external systems.

|                                     | Description                                                                                                    | Self-hosting | ZITADEL Cloud |
|-------------------------------------|----------------------------------------------------------------------------------------------------------------|-------------|---------------|
| Events-API                          | Pulling events of all ZITADEL resources such as Users, Projects, Apps, etc. (Events = Change Log of Resources) | ✅           | ✅             |
| Cockroach Change Data Capture       | Sending events of all ZITADEL resources such as Users, Projects, Apps, etc. (Events = Change Log of Resources) | ✅           | ❌             |
| ZITADEL Actions Log to Stdout       | Custom log to messages possible on predefined triggers during login / register Flow                            | ✅           | ❌             |
| ZITADEL Actions trigger API/Webhook | Custom API/Webhook request on predefined triggers during login / register                                      | ✅           | ✅             |

### Events API

The ZITADEL Event API empowers you to proactively pull audit logs for comprehensive security and compliance monitoring, regardless of your environment (cloud or self-hosted). 
This API offers granular control through various filters, enabling you to:
- **Specify Event Types**: Focus on specific events of interest, such as user token created, password changed, or project added.
- **Target Aggregates**: Narrow down the data scope by filtering for events related to particular organizations, projects, or users.
- **Define Time Frames**: Retrieve audit logs for precise time periods, allowing you to schedule data retrieval at desired intervals (e.g., hourly) and analyze activity within specific windows.

You can find a comprehensive guide on how to use the events API for different use cases here: [Get Events from ZITADEL](/docs/guides/integrate/zitadel-apis/event-api)

### Cockroach Change Data Capture

For self-hosted ZITADEL deployments utilizing CockroachDB as the database, [CockroachDB's built-in Change Data Capture (CDC)](https://www.cockroachlabs.com/docs/stable/change-data-capture-overview) functionality provides a streamlined approach to integrate ZITADEL audit logs with external systems.

CDC captures row-level changes in your database and streams them as messages to a configurable destination, such as Google BigQuery or a SIEM/SOC solution. This real-time data stream enables:
- **Continuous monitoring**: Receive near-instantaneous updates on ZITADEL activity, facilitating proactive threat detection and response.
- **Simplified integration**: Leverage CockroachDB's native capabilities for real-time data transfer, eliminating the need for additional tools or configurations.

This approach is limited to self-hosted deployments using CockroachDB and requires expertise in managing the database and CDC configuration.

#### Sending events to Google Cloud Storage using Change Data Capture

This example will show you how you can utilize CDC for sending all ZITADEL events to Google Cloud Storage.
For a detailed description please read [CockroachLab's Get Started Guide](https://www.cockroachlabs.com/docs/v23.2/create-and-configure-changefeeds) and [Cloud Storage Authentication](https://www.cockroachlabs.com/docs/v23.2/cloud-storage-authentication?filters=gcs#set-up-google-cloud-storage-assume-role) from Cockroach.

You will need a Google Cloud Storage Bucket and a service account.
1. [Create Google Cloud Storage Bucket](https://cloud.google.com/storage/docs/creating-buckets)
2. [Create Service Account](https://cloud.google.com/iam/docs/service-accounts-create)
3. Create a role with the `storage.objects.create` permission
4. Grant service account access to the bucket
5. Create key for service account and download it

Now we need to enable and create the changefeed in the cockroach DB.
1. [Enable rangefeeds on cockroach cluster](https://www.cockroachlabs.com/docs/v23.2/create-and-configure-changefeeds#enable-rangefeeds)
  ```bash
  SET CLUSTER SETTING kv.rangefeed.enabled = true;
  ```
2. Encode the keyfile from the service account with base64 and replace the placeholder it in the script below
3. Create Changefeed to send data into Google Cloud Storage
   The following example sends all events without payload to Google Cloud Storage
   Per default we do not want to send the payload of the events, as this could potentially include personally identifiable information (PII)
   If you want to include the payload, you can just add `payload` to the select list in the query.  
   ```sql
    CREATE CHANGEFEED INTO 'gs://gc-storage-zitadel-data/events?partition_format=flat&AUTH=specified&CREDENTIALS=base64encodedkey' 
    AS SELECT instance_id, aggregate_type, aggregate_id, owner, event_type, sequence, created_at 
    FROM eventstore.events2;
   ```

In some cases you might want the payload of only some specific events.
This example shows you how to get all events and the instance domain events with the payload:
   ```sql
    CREATE CHANGEFEED INTO 'gs://gc-storage-zitadel-data/events?partition_format=flat&AUTH=specified&CREDENTIALS=base64encodedkey' 
    AS SELECT instance_id, aggregate_type, aggregate_id, owner, event_type, sequence, created_at 
    CASE WHEN event_type IN ('instance.domain.added', 'instance.domain.removed', 'instance.domain.primary.set' ) 
    THEN payload END AS payload 
    FROM eventstore.events2;
   ```

The partition format in the example above is flat, this means that all files for each timestamp will be created in the same folder.
You will have files for different timestamps including the output for the events created in that time.
Each event is represented as a json row.

Example Output:
```json lines
{
   "aggregate_id": "26553987123463875", 
   "aggregate_type": "user",
   "created_at": "2023-12-25T10:01:45.600913Z",
   "event_type": "user.human.added",
   "instance_id": "123456789012345667", 
   "payload": null,
   "sequence": 1
}
```

## ZITADEL Actions

ZITADEL [Actions](/docs/concepts/features/actions) offer a powerful mechanism for extending the platform's capabilities and integrating with external systems tailored to your specific requirements. 
Actions are essentially custom JavaScript snippets that execute at predefined triggers during the registration or login flow of a user.

In the future ZITADEL Actions will be extended to allow to not only define them during the login and register flow, but also on each API Request, Event or Predefined Functions.

### Log to stdout

With the [log module](/docs/apis/actions/modules#log) you can log any custom message to stdout.
Those logs in stdout can be collected by your external system.

Example Use Case:
In my external system for example Splunk I want to be able to get an information each time a user has authenticated.

1. Define an action that logs successful and failed login to your stdout.
   Make sure the name of the action is the same as of the function in the script.
   ```ts reference
   https://github.com/zitadel/actions/blob/main/examples/post_auth_log.js
   ```
2. Add the action to the following Flows and Triggers
   - Flow: Internal Authentication - Trigger: Post Authentication
   - Flow: External Authentication - Trigger: Post Authentication
3. Authenticate User
4. Collect Data from stdout

### Webhook/API request

The [http module](/docs/apis/actions/modules#http) allows you to make a request to a REST API. 
This allows you to send a request at a specific point during the login or registration flow with the data you defined in your action.

Example use case:
You want to send a request to an endpoint each time after an authentication (successful or not).

1. Define an action that calls API endpoint.
   Make sure the name of the action is the same as of the function in the script.
   Example how to call an API endpoint:
   ```ts reference
   https://github.com/zitadel/actions/blob/main/examples/make_api_call.js
   ```
2. Add the action to the following flows and triggers
   - Flow: Internal Authentication - Trigger: Post Authentication
   - Flow: External Authentication - Trigger: Post Authentication
3. Authenticate user
4. Get data on your API
