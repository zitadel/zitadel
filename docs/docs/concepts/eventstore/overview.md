---
title: Eventstore
sidebar_label: Overview
---

ZITADEL is built on the [Event Sourcing pattern](../architecture/software), where changes are stored as events in an Event Store.

## What is an Event Store?

Traditionally, data is stored in relations as a state

- A request needs to know the relations to select valid data
- If a relation changes, the requests need to change as well
- That is valid for actual, as well as for historical data

An Event Store on the other hand stores events, meaning every change that happens to any piece of data relates to an event.
The data is stored as events in an append-only log.

- Think of it as a ledger that gets new entries over time, accumulative
- To request data, all you have to do is to sum the events as the summary reflects the actual state
- To investigate past changes to your system, you just select the events from your time range of interest
- That makes audit/analytics very powerful, due to the historical data available to build queries

## Benefits

- Audit: You have a built-in audit trail that tracks all changes over an unlimited period of time.
- Travel back in time: With our way of storing data we can show you all of your resources at a given point in time. 
- Future Projections: It is easy to compute projections with new business logic by replaying all events since installation.

## Definitions

Event Sourcing has some specific terms that are often used in our documentation. To understand how ZITADEL works it is important to understand this key definitions.

### Events

An event is something that happens in the system and gets written to the database. This is the single source of truth.
Events are immutable and the current state of your system is derived from the events.

Possible Events:
- user.added
- user.changed
- product.added
- user.password.checked

### Aggregate

An aggregate consist of multiple events. All events together from an aggregate will lead to the current state of the aggregate.
The aggregate can be compared with an object or a resources. An aggregates should be used as transaction boundary.

### Projections

Projections contain the computed objects, that will be used on the query side for all the requests.
Think of this as a normalized view of specific events of one or multiple aggregates.
