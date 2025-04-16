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
| ZITADEL Actions Log to Stdout       | Custom log to messages possible on predefined triggers during login / register Flow                            | ✅           | ❌             |
| ZITADEL Actions trigger API/Webhook | Custom API/Webhook request on predefined triggers during login / register                                      | ✅           | ✅             |

### Events API

The ZITADEL Event API empowers you to proactively pull audit logs for comprehensive security and compliance monitoring, regardless of your environment (cloud or self-hosted). 
This API offers granular control through various filters, enabling you to:
- **Specify Event Types**: Focus on specific events of interest, such as user token created, password changed, or project added.
- **Target Aggregates**: Narrow down the data scope by filtering for events related to particular organizations, projects, or users.
- **Define Time Frames**: Retrieve audit logs for precise time periods, allowing you to schedule data retrieval at desired intervals (e.g., hourly) and analyze activity within specific windows.

You can find a comprehensive guide on how to use the events API for different use cases here: [Get Events from ZITADEL](/docs/guides/integrate/zitadel-apis/event-api)

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
