---
title: ZITADEL's In-built Audit Trail
sidebar_label: Audit Trail
---

ZITADEL provides you with an built-in audit trail to track all changes and events over an unlimited period of time.
Most other solutions replace a historic record and track changes in a separate log when information is updated.
ZITADEL only ever appends data in an [Eventstore](/docs/concepts/eventstore/overview), keeping all historic record.
The audit trail itself is identical to the state, since ZITADEL calculates the state from all the past changes.

![Example of events that happen for a profile change and a login](/img/concepts/audit-trail/audit-log-events.png)

This form of audit log has several benefits over storing classic audit logs.
You can view past data in-context of the whole system at a single point in time.
Reviewing a past state of the application can be important when tracing an incident that happened months back. Moreover the eventstore provides a truly complete and clean audit log.

:::info Future Plans
There will be three major areas for future development on the audit data

- [Metrics](https://github.com/zitadel/zitadel/issues/4458) and [standard reports](https://github.com/zitadel/zitadel/discussions/2162#discussioncomment-1153259)
- [Feedback loop](https://github.com/zitadel/zitadel/issues/5102) and threat detection
- Forensics and replay of events

:::

## Accessing the Audit Log

### Last changes of an object

You can check the last changes of most objects in the [Console](/docs/guides/manage/console/overview).
In the following screenshot you can see an example of last changes on an [user](/docs/guides/manage/console/users).
The same view is available on several other objects such as organization or project.

![Profile Self Manage](/img/guides/console/myprofile.png)

### Event View

Administrators can see all events across an instance and filter them directly in [Console](/docs/guides/manage/console/overview).
Go to your default settings and then click on the Tab **Events** to open the Event Viewer or browse to $CUSTOM-DOMAIN/ui/console/events

![Event viewer](/img/concepts/audit-trail/event-viewer.png)

### Event API

Since everything that is available in Console can also be called with our APIs, you can access all events and audit data trough our APIs:

- [Event API Guide](/docs/guides/integrate/zitadel-apis/event-api)
- [API Documentation](/docs/apis/resources/admin/events)

Access to the API is possible with a [Service User](/docs/guides/integrate/service-users/authenticate-service-users) account, allowing you to integrate the events with your own business logic.

## Using logs in external systems

You can use the events from the audit log in external systems such as a SOC/SIEM solution.

Follow our guide on how to [integrate ZITADEL with external systems for streaming events and audit logs](/docs/guides/integrate/external-audit-log).
