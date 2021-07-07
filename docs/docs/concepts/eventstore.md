---
title: Eventstore
---

ZITADEL is built on the [event sourcing pattern](architecture), where changes are stored as events in an eventstore.

## What is an eventstore?

Traditionally, data is stored in relations as a state

- A request needs to know the relations to select valid data
- If a relation changes, the requests need to change as well
- That is valid for actual, as well as for historical data

An Eventstore on the other hand stores events, meaning every change that happens to any piece of data relates to an event

- Think of it as a ledger that gets new entries over time, accumulative
- To request data, all you have to do is to sum the events as the summary reflects the actual state
- To investigate past changes to your system, you just select the events from your time range of interest
- That makes audit/analytics very powerful, due to the historical data available to build queries
