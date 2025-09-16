---
title: Service Ping
sidebar_label: Service Ping 
---

Service Ping is a feature that periodically sends anonymized analytics and usage data from your ZITADEL system to a central endpoint.
This data helps improve ZITADEL by providing insights into its usage patterns.

The feature is enabled by default, but can be disabled either completely or for specific reports.
Checkout the configuration options below.

## Data Sent by Service Ping

### Base Information

If the feature is enabled, the base information will always be sent. To prevent that, you can opt out by disabling the entire Service Ping:

```yaml
ServicePing:
  Enabled: false # ZITADEL_SERVICEPING_ENABLED
```

The base information sent back includes the following:
- your systemID
- the currently run version of ZITADEL
- information on all instances
  - id
  - creation date
  - domains

### Resource Counts

Resource counts is a report that provides us with information about the number of resources in your ZITADEL instances.

The following resources are counted:
- Instances
- Organizations
- Projects per organization
- Users per organization
- Users of type machine per organization
- SCIM provisioned users per organization
- Instance Administrators
- Identity Providers
- LDAP Identity Providers
- Actions (V1)
- Targets and set up executions
- Login Policies
- MFA enforcement (if either MFA is required for local or all users through the login policy)
- Password Complexity Policies
- Password Expiry Policies
- Lockout Policies
- Notification Policies with option "Password change" enabled

The list might be extended in the future to include more resources.

To disable this report, set the following in your configuration file:

```yaml
ServicePing:
  Telemetry:
    ResourceCounts:
      Enabled: false # ZITADEL_SERVICEPING_TELEMETRY_RESOURCECOUNT_ENABLED
```

## Configuration

The Service Ping feature can be configured through the runtime configuration. Please check out the configuration file
for all available options. Below is a list of the most important options:

### Interval

This defines at which interval the Service Ping feature sends data to the central endpoint. It supports the extended cron syntax
and by default is set to `@daily`, which means it will send data every day. The time is randomized on startup to prevent
all systems from sending data at the same time.

You can adjust it to your needs to make sure there is no performance impact on your system.
For example, if you already have some scheduled job syncing data in and out of ZITADEL around a specific time or have regularly a
lot of traffic during the day, you might want to change it to a different time, e.g. `15 4 * * *` to send it every day at 4:15 AM.

The interval must be at least 30 minutes to prevent too frequent requests to the central endpoint.

### MaxAttempts

This defines how many attempts the Service Ping feature will make to send data to the central endpoint before giving up
for a specific interval and report. If one report fails, it will be retried up to this number of times. 
Other reports will still be handled in parallel and have their own retry count. This means if the base information 
only succeeded after three attempts, the resource count still has five attempts to be sent. 
