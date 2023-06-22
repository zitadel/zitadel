---
title: Cloud Service
custom_edit_url: null
--- 

This annex of the [Framework Agreement](terms-of-service) describes the service levels offered by us for our Services.

## Definitions

**Monthly quota** means the available usage per measure for one billing period. The quota is reset to zero with the start of a new billing period.

**Authenticated request** means any request to our API endpoints requiring a valid authorization header. We exclude requests with a server error, discovery endpoints, and endpoints to load UI assets.

**Action minutes** means execution time, rounded up to 1 second, of custom code execution via a customer defined Action.

**Adequate Country** means a country or territory recognized as providing an adequate level of protection for Personal Data under an adequacy decision made, from time to time, by (as applicable) (i) the Information Commissioner's Office and/or under applicable UK law (including the UK GDPR), or (ii) the [European Commission under the GDPR](https://ec.europa.eu/info/law/law-topic/data-protection/international-dimension-data-protection/adequacy-decisions_en).

## Data location

Data location refers to a region, consisting of one or many countries or territories, where the customer's data is stored in our database and processed by our systems.

We can not guarantee that during transit the data will only remain within this region. We take measures, as outlined in our [privacy policy](privacy-policy), to protect your data in transit and in rest.

The following regions will be available when using our cloud service. This list is for informational purposes and will be updated in due course, please refer to our website for all available regions at this time.

- **Global**: All available cloud regions offered by our cloud provider
- **Switzerland**: Exclusively on Swiss region
- **GDPR safe countries**: Exclusively [Adequate Countries](https://ec.europa.eu/info/law/law-topic/data-protection/international-dimension-data-protection/adequacy-decisions_en) as recognized by the European Commission under the GDPR

## Backup

Our backup strategy executes daily full backups and differential backups on much higher frequency.
In a disaster recovery scenario, our goal is to guarantee a recovery point objective (RPO) of 1h, and a higher but similar recovery time objective (RTO).
Under normal operations, RPO and RTO goals are below 1 minute.

If you you have different requirements we provide you with a flexible approach to backup, restore, and transfer data (f.e. to a self-hosted setup) through our APIs.
Please consult the [migration guides](../guides/migrate/introduction.md) for more information.

Last revised: June 21, 2023