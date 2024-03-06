---
title: Streaming Audit Logs to External Systems (SIEM/SOC)
sidebar_label: Streaming Audit Logs to External Systems
---

This document details integrating ZITADEL with external systems for streaming events and audit logs. 
This functionality allows you to centralize ZITADEL activity data alongside other security and operational information, facilitating comprehensive monitoring, analysis, and compliance reporting.

Integrating ZITADEL with external systems offers several advantages:
- **Centralized Monitoring**: Streamlining ZITADEL events into a single platform alongside data from various sources enables a consolidated view of your security posture. This comprehensive view simplifies threat detection, investigation, and incident response.
- **Enhanced Security Analytics**: External systems, such as Security Information and Event Management (SIEM) solutions, can leverage ZITADEL events to identify suspicious activities, potential security breaches, and user access anomalies.
- **Compliance Reporting**: ZITADEL events can be used to generate detailed audit trails, fulfilling regulatory compliance requirements for data access and user activity.

By integrating ZITADEL with external systems, you gain valuable insights into user behavior, system activity, and potential security threats, ultimately strengthening your overall security posture and regulatory compliance.

ZITADEL does provide different solutions how to send events to external systems, the solution you choose might differ depending on your use case, database and environment (ZITADEL Cloud, Self-hosting) you are using.

## Change Log / Events

ZITADEL is based on an [event sourcing architecture](https://zitadel.com/docs/concepts/eventstore/overview), which means that each change happening on a resource is stored as event. 
This allowed you to get the change track of all resources.

There are different solutions how to get/send events to external systems, the solution you choose might differ depending on your use case, database and environment (ZITADEL Cloud, Self-hosting) you are using.

|                               | Self-hosting | ZITADEL Cloud |
|-------------------------------|--------------|---------------|
| Events-API                    | yes          | yes           |
| Cockroach Change Data Capture | yes          | no            |
| Postgres whatever             | yes          | no            |

### Events API

The ZITADEL Event API empowers you to proactively pull audit logs for comprehensive security and compliance monitoring, regardless of your deployment model (cloud or self-hosted). 
This API offers granular control through various filters, enabling you to:
- **Specify Event Types**: Focus on specific events of interest, such as user token created, password changed, or project added.
- **Target Aggregates**: Narrow down the data scope by filtering for events related to particular organizations, projects, or users.
- **Define Time Frames**: Retrieve audit logs for precise time periods, allowing you to schedule data retrieval at desired intervals (e.g., hourly) and analyze activity within specific windows.

You can find a comprehensive guide on how to use the events api for different use cases here: [Get Events from ZITADEL](/docs/guides/integrate/event-api)

### Cockroach Change Data Capture

For self-hosted ZITADEL deployments utilizing CockroachDB as the database, [CockroachDB's built-in Change Data Capture (CDC)](https://www.cockroachlabs.com/docs/stable/change-data-capture-overview) functionality provides a streamlined approach to integrate ZITADEL audit logs with external systems.

CDC captures row-level changes in your database and streams them as messages to a configurable destination, such as Google BigQuery or a SIEM/SOC solution. This real-time data stream enables:
- **Continuous monitoring**: Receive near-instantaneous updates on ZITADEL activity, facilitating proactive threat detection and response.
- **Simplified integration**: Leverage CockroachDB's native capabilities for real-time data transfer, eliminating the need for additional tools or configurations.

This approach is limited to self-hosted deployments using CockroachDB and requires expertise in managing the database and CDC configuration.

#### Sending events to Google Big Query

This example will show you how you can utilize CDC for sending all ZITADEL events to Google Big Query.
For a detailed description please read the [Get Started Guide](https://www.cockroachlabs.com/docs/v23.2/create-and-configure-changefeeds) and [Cloud Storage Authentication](https://www.cockroachlabs.com/docs/v23.2/cloud-storage-authentication?filters=gcs#set-up-google-cloud-storage-assume-role) from Cockroach.

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
```bash
  CREATE CHANGEFEED INTO 'gs://gc-storage-zitadel-data/events?partition_format=flat&AUTH=specified&CREDENTIALS=base64encodedkey' AS SELECT instance_id, aggregate_type, aggregate_id, owner, event_type, sequence, created_at FROM eventstore.events2;
  ```

In some cases you might want the payload of only some specific events.
This example shows you how to get all events and the instance domain events with the payload:
```bash
CREATE CHANGEFEED INTO 'gs://gc-storage-zitadel-data/events?partition_format=flat&AUTH=specified&CREDENTIALS=base64encodedkey' AS SELECT instance_id, aggregate_type, aggregate_id, owner, event_type, sequence, created_at CASE WHEN event_type IN ('instance.domain.added', 'instance.domain.removed', 'instance.domain.primary.set' ) THEN payload END AS payload FROM eventstore.events2;
```

The partition format in the example above is flat, this means that all files for each timestamp will be created in the same folder.
You will have files for different timestamps including the output for the events created in that time.
Each event is represented as a json row.

Example Output:
```json lines
{"aggregate_id": "26553987123463875", "aggregate_type": "user", "created_at": "2023-12-25T10:01:45.600913Z", "event_type": "user.human.added", "instance_id": "123456789012345667", "payload": null, "sequence": 1}
```