import React from "react";

export function SubProcessorTable() {

  const country_list = {
    us: "USA",
    eu: "EU",
    ch: "Switzerland",
    fr: "France",
    in: "India",
    de: "Germany",
    ee: "Estonia",
    nl: "Netherlands",
    ro: "Romania",
  }
  const processors = [
    {
      entity: "Google LLC",
      purpose: "Cloud infrastructure provider (Google Cloud), business applications and collaboration (Workspace), Data warehouse services, Content delivery network, DDoS and bot prevention",
      hosting: "Region designated by Customer, United States",
      country: country_list.us,
      enduserdata: "Yes (transit)"
    },
    {
      entity: "Cockroach Labs, Inc.",
      purpose: "Managed database services: Dedicated CockroachDB clusters on Google Cloud",
      hosting: "Region designated by Customer",
      country: country_list.us,
      enduserdata: "Yes (at rest)"
    },
    {
      entity: "Datadog, Inc.",
      purpose: "Infrastructure monitoring, log analytics, and alerting",
      hosting: country_list.eu,
      country: country_list.us,
      enduserdata: "Yes (logs)"
    },
    {
      entity: "Github, Inc.",
      purpose: "Source code management, code scanning, dependency management, security advisory, issue management, continuous integration",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Stripe Payments Europe, Ltd.",
      purpose: "Subscription management, payment process",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Bexio AG",
      purpose: "Customer management, payment process",
      hosting: country_list.ch,
      country: country_list.ch,
      enduserdata: false
    },
    {
      entity: "Mailjet SAS",
      purpose: "Marketing automation",
      hosting: country_list.eu,
      country: country_list.fr,
      enduserdata: false
    },
    {
      entity: "Postmark (AC PM LLC)",
      purpose: "Transactional mails, if no customer owned SMTP service is configured",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: "Yes (opt-out)"
    },
    {
      entity: "Vercel, Inc.",
      purpose: "Website hosting",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Agolia SAS",
      purpose: "Documentation search engine (zitadel.com/docs)",
      hosting: country_list.us,
      country: country_list.in,
      enduserdata: false
    },
    {
      entity: "Discord Netherlands BV",
      purpose: "Community chat (zitadel.com/chat)",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Statuspal",
      purpose: "ZITADEL Cloud service status announcements",
      hosting: country_list.us,
      country: country_list.de,
      enduserdata: false
    },
    {
      entity: "Plausible Insights OÜ",
      purpose: "Privacy-friendly web analytics",
      hosting: country_list.de,
      country: country_list.ee,
      enduserdata: false,
      dpa: 'https://plausible.io/dpa'
    },
    {
      entity: "Twillio Inc.",
      purpose: "Messaging platform for SMS",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: "Yes (opt-out)"
    },
    {
      entity: "Mohlmann Solutions SRL",
      purpose: "Global payroll",
      hosting: undefined,
      country: country_list.ro,
      enduserdata: false
    },
    {
      entity: "Remote Europe Holding, B.V.",
      purpose: "Global payroll",
      hosting: undefined,
      country: country_list.nl,
      enduserdata: false
    },
    {
      entity: "Clickhouse, Inc.",
      purpose: "Data warehouse services",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
  ]

  return (
    <table className="text-xs">
      <tr>
        <th>Entity name</th>
        <th>Purpose</th>
        <th>End-user data</th>
        <th>Hosting location</th>
        <th>Country of registration</th>
      </tr>
      {
        processors
          .sort((a, b) => {
            if (a.entity < b.entity) return -1
            if (a.entity > b.entity) return 1
            else return 0
          })
          .map((processor, rowID) => {
          return (
            <tr>
              <td key={rowID}>{processor.entity}</td>
              <td>{processor.purpose}</td>
              <td>{processor.enduserdata ? processor.enduserdata  : 'No'}</td>
              <td>{processor.hosting ? processor.hosting  : 'n/a'}</td>
              <td>{processor.country}</td>
            </tr>
          )
        })
      }
    </table>
  );
}
