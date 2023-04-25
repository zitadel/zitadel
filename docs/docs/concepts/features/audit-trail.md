---
title: Audit Trail
---

ZITADEL provides you with an built-in audit trail to track all changes and events over an unlimited period of time.
Most solutions replace a historic record and track changes in a separate log when information is updated.
ZITADEL only ever appends data in an [Eventstore](https://docs.zitadel.com/docs/concepts/eventstore), keeping all historic record.
The audit trail itself is identical to the state, since ZITADEL calculates the state from all the past changes.

![Example of events that happen for a profile change and a login](/img/concepts/audit-log-events.png)

This form of audit log has several real-life benefits.
You can view past data in-context of the whole system at a single point in time.
Reviewing a past state of the application can be important when tracing an incident that happened months back. Moreover the eventstore provides a truly complete and clean audit log.

## Accessing the Audit Log

### Last changes of an object

### Event viewer

### Event API

## Using logs in external systems

## Future plans

- How to access information in the audit trail via GUI?
- How to access audit information via APIs (incl. ~planned~ features like event-API)?

We could include the following as additional info: 
- Future plans for standard reports
- ~Future plans for~ APIs (and why pull vs. push to handle backpressure of an HA system)
- Using / sending data to external log system