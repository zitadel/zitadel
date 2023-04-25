---
title: Audit Trail
---

ZITADEL provides you with an built-in audit trail to track all changes and events over an unlimited period of time.
Most solutions replace a historic record and track changes in a separate log when information is updated.
ZITADEL only ever appends data in an [Eventstore](https://docs.zitadel.com/docs/concepts/eventstore), keeping all historic record.
The audit trail itself is identical to the state, since ZITADEL calculates the state from all the past changes.

![Example of events that happen for a profile change and a login](/img/concepts/audit-trail/audit-log-events.png)

This form of audit log has several real-life benefits.
You can view past data in-context of the whole system at a single point in time.
Reviewing a past state of the application can be important when tracing an incident that happened months back. Moreover the eventstore provides a truly complete and clean audit log.

## Accessing the Audit Log

### Last changes of an object

You can check the last changes of most objects in the [Console](docs/guides/manage/console/overview).
In the following screenshot you can see an example of last changes on an [user](/docs/guides/manage/console/users).
The same view is available on several other objects such as organization or project.

![Profile Self Manage](/img/guides/console/myprofile.png)

### Event viewer

Administrators can see all events across an instance and filter them directly in [Console](docs/guides/manage/console/overview).
Go to your instance settings and then click on the Tab **Events** to open the Event Viewer or browse to $YOUR_DOMAIN/ui/console/events  

![Profile Self Manage](/img/concepts/audit-trail/event-viewer.png)

### Event API

- Guide: https://zitadel.com/docs/guides/integrate/event-api
- API Docs: https://zitadel.com/docs/category/apis/admin/events

## Using logs in external systems

## Future plans

- How to access information in the audit trail via GUI?
- How to access audit information via APIs (incl. ~planned~ features like event-API)?

We could include the following as additional info: 
- Future plans for standard reports
- ~Future plans for~ APIs (and why pull vs. push to handle backpressure of an HA system)
- Using / sending data to external log system